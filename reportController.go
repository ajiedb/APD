package controllers

import (
	"net/http"
	"ujicoba-go/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"fmt"
	"ujicoba-go/models"
	"strconv"
	"time"
	"encoding/json"
	"bytes"
	"regexp"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

type ReportController struct{}

type packageForm struct {
	Package []string `form:"package" validate:"required,required-dynamic"`
}

var fixedHeaders = []string{
		"Kd Prov",
		"Kd Rayon",
		"Kd Sek",
		"NPSN",
		"Kd Jenjang",
		"Sts Sek",
		"ID Sekolah",
		"Kelas",
		"ID Siswa",
		"Nama Siswa (Kode user)",
		"Paket Soal",
		"Event",
		"Ujian",
	}

var fixedHeadersResponseJawaban = []string{
		"Nama Siswa (Kode user)",
		"Paket Soal",
		"Event",
		"Ujian",
	}

func ConvAnsToMap(answer []models.TestAnswers) map[string]map[string]map[string]interface{}{
	mapTestAnswers := make(map[string]map[string]map[string]interface{})

	for _, ans := range answer {
		testParticipantID := strconv.Itoa(int(ans.TestParticipantID))
		if _, ok := mapTestAnswers[testParticipantID]; !ok {
			mapTestAnswers[testParticipantID] = make(map[string]map[string]interface{})
		}

		if _, ok := mapTestAnswers[testParticipantID][ans.SiapQuestionUUID]; !ok {
			mapTestAnswers[testParticipantID][ans.SiapQuestionUUID] = make(map[string]interface{})
		}

		mapTestAnswers[testParticipantID][ans.SiapQuestionUUID]["siap_answer_uuid"] = ans.SiapAnswerUUID
		val := ans.Answer
		if val == nil {
			mapTestAnswers[testParticipantID][ans.SiapQuestionUUID]["answer"] = nil
		} else {
			var data []map[string]interface{}
			json.Unmarshal([]byte(*val), &data)
			if(len(data) > 0){
				mapData := make(map[string]interface{})
				for _, dat := range data {
					id, ok := dat["id"].(string)
					if ok {
						mapData[id] = dat["jawaban_cp"]
					}
				}
				mapTestAnswers[testParticipantID][ans.SiapQuestionUUID]["answer"] = mapData
			 
			} else {
				if(*val != "[]" && *val != "null" && *val != ""){
					mapTestAnswers[testParticipantID][ans.SiapQuestionUUID]["answer"] = *val
					
				} else {
					mapTestAnswers[testParticipantID][ans.SiapQuestionUUID]["answer"] = nil
				}	
			}
			
		}
		
	}

	return mapTestAnswers
}

func ConvAnsToMapGroup(answer []models.TestAnswers, participant []models.TestParticipantReport) map[string]int{
	mapGroupTestAnswers := make(map[string]int)

	for _, ans := range answer {
		siapQuestionCode := ans.SiapQuestionCode
		val := ans.Answer
		if _, ok := mapGroupTestAnswers[siapQuestionCode]; !ok {
			mapGroupTestAnswers[siapQuestionCode] = 0
		}
		for _,participant := range participant{ 
			if(ans.TestParticipantID == uint(participant.IDSiswa)){
				if((ans.SiapAnswerUUID != "" && ans.SiapAnswerUUID != "null") || (val != nil && *val != "[]" && *val != "null" && *val != "")){
					mapGroupTestAnswers[siapQuestionCode] +=1
				}
			}
		}
	}

	return mapGroupTestAnswers
}


func ReportFileName(reportType string, t string, uuid string) string {
    return fmt.Sprintf("%s_%s_%s.xlsx", reportType, t, uuid)
}

func ExportSheetName(t string) string {
    return fmt.Sprintf("Export Jawaban %s", t)
}

func (ReportController) CreateReportRequest(c *gin.Context) {
	defer utils.EndProccess()
	var validate *validator.Validate = utils.GetValidate()
	var dataForm packageForm
	c.BindJSON(&dataForm)

	err := validate.Struct(dataForm)
	if err != nil {
		utils.ValidateErrorHandle(err, c)
		return
	}

	dataReportReqs := models.ReportReqs{}
	UUID, err := dataReportReqs.CreateReportReqs(&dataReportReqs)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"UUID":0,"message": err.Error()})
	}

	var dataReportPackages []models.ReportPackages
	for _, packages := range dataForm.Package {
		dataReportPackages = append(dataReportPackages, models.ReportPackages{
			ReqID:				dataReportReqs.ID,
			SiapTestPackageUUID:	packages,
		})
	}

	models.ReportPackages{}.CreateReportPackages(dataReportPackages)


	var dataReportOutcomes []models.ReportOutcomes
	for i := 0; i < 5; i++ {
		dataReportOutcomes = append(dataReportOutcomes, models.ReportOutcomes{
			ReqID:				dataReportReqs.ID,
			Type:				i+1,
			Status:				0,
		})
	}

	models.ReportOutcomes{}.CreateReportOutcomes(dataReportOutcomes)
	c.JSON(http.StatusOK, gin.H{"UUID":UUID,"message": "request successfully submitted"})
}

type ReportResult struct {
	FileName string
	UUID     string
	Status   string
}

func GenerateScoreReport(uuid string) (*ReportResult, error) {
	repPackages := models.ReportPackages{}
	dataRepPackages, err := repPackages.GetByReqUUID(uuid)
	if err != nil {
		return nil, err
	}

	repPackagesIDs := models.QuestionPackage{}
	dataRepPackagesIDs, _:= repPackagesIDs.GetByUUIDs(dataRepPackages)

	var packageIDs []int
	for _, packageID := range dataRepPackagesIDs {
		packageIDs = append(packageIDs, int(packageID.ID))
	}

	testParticipants := models.TestParticipantReport{}
	dataTestParticipants, err := testParticipants.GetParticipantAnswer(packageIDs)

	testSoals := models.QuestionPackageSoalReport{}
	testSoalsType := []int{
		utils.QUESTION_TYPE_MULTIPLE_CHOICE_TABLE,
		utils.QUESTION_TYPE_MATCHING_TABLE,
		utils.QUESTION_TYPE_COMPLEX_MULTIPLE_CHOICE_MCMA_TABLE,
		utils.QUESTION_TYPE_COMPLEX_MULTIPLE_CHOICE_RADIO_TABLE,
	}
	dataTestSoals, err := testSoals.GetQuestionPackageSoalReport(dataRepPackages, testSoalsType)

	testAnswers := models.TestAnswers{}
	dataTestAnswers, err := testAnswers.GetDataExport(dataRepPackages, []string{})

	mapTestAnswers := ConvAnsToMap(dataTestAnswers)

	var soalID1 []int
	var soalID3 []int
	var soalID8 []int
	var soalID9 []int
	for i := 0; i < len(dataTestSoals); i++ {
		if(dataTestSoals[i].SoalTypeID == utils.QUESTION_TYPE_MULTIPLE_CHOICE_TABLE ){
			soalID1 = append(soalID1, dataTestSoals[i].SoalID)
		} else if (dataTestSoals[i].SoalTypeID == utils.QUESTION_TYPE_MATCHING_TABLE) {
			soalID3 = append(soalID3, dataTestSoals[i].SoalID)
		} else if (dataTestSoals[i].SoalTypeID == utils.QUESTION_TYPE_COMPLEX_MULTIPLE_CHOICE_MCMA_TABLE) {
			soalID8 = append(soalID8, dataTestSoals[i].SoalID)
		} else if (dataTestSoals[i].SoalTypeID == utils.QUESTION_TYPE_COMPLEX_MULTIPLE_CHOICE_RADIO_TABLE) {
			soalID9 = append(soalID9, dataTestSoals[i].SoalID)
		} 
	}

	jawabanKunci := models.JawabanKunci{}
	var jawabanKunciPG,jawabanKunciPC, jawabanKunciPGK1, jawabanKunciPGK2 map[string]interface{}
	if(len(soalID1) > 0){
		jawabanKunciPG,err = jawabanKunci.GetJawabanKunciPG(soalID1)
	}

	if(len(soalID3) > 0){
		jawabanKunciPC,err = jawabanKunci.GetJawabanKunciPC(soalID3)
	}

	if(len(soalID8) > 0){
		jawabanKunciPGK1,err = jawabanKunci.GetJawabanKunciPGK1(soalID8)
	}

	if(len(soalID9) > 0){
		jawabanKunciPGK2,err = jawabanKunci.GetJawabanKunciPGK2(soalID9)
	}
	lenSoalOps := 0

	for _, soal := range dataTestSoals {
		soalID := strconv.Itoa(soal.SoalID)
		soalType := soal.SoalTypeID
		val := 0
		if(soalType == 1) {
			val = int(jawabanKunciPG[soalID].(map[string]interface{})["opsNum"].(float64))
		} else if (soalType == 3){
			val = int(jawabanKunciPC[soalID].(map[string]interface{})["opsNum"].(float64))
		} else if (soalType == 8){
			val = int(jawabanKunciPGK1[soalID].(map[string]interface{})["opsNum"].(float64))
		} else if (soalType == 9){
			val = int(jawabanKunciPGK2[soalID].(map[string]interface{})["opsNum"].(float64))
		}
		
		num := int(val)
		if err == nil {
			lenSoalOps += num
		}
	}

	var headers []interface{}
	var subHeaders []interface{}
	for _, h := range fixedHeaders {
		headers = append(headers, h)
		subHeaders = append(subHeaders,"")
	}
	
	soalTypeConstant := utils.SoalTypeConstant{}

	for _, soal := range dataTestSoals {
		soalID := strconv.Itoa(soal.SoalID)
		soalType := soal.SoalTypeID
		val := 0
		if(soalType == 1) {
			val = int(jawabanKunciPG[soalID].(map[string]interface{})["opsNum"].(float64))
		} else if (soalType == 3){
			val = int(jawabanKunciPC[soalID].(map[string]interface{})["opsNum"].(float64))
		} else if (soalType == 8){
			val = int(jawabanKunciPGK1[soalID].(map[string]interface{})["opsNum"].(float64))
		} else if (soalType == 9){
			val = int(jawabanKunciPGK2[soalID].(map[string]interface{})["opsNum"].(float64))
		}
		num := int(val)
		if err == nil {
			for i := 0; i < num; i++ {
				headers = append(headers, soal.SoalCode + " " +strconv.Itoa(i+1))
				subHeaders = append(subHeaders,soalTypeConstant.SoalType(soal.SoalTypeID))
			}
		}
	}

	t := time.Now().Format("02012006_1504")
	fileName :=  ReportFileName(utils.REPORT_FOURTH, t, utils.ShortUUID(uuid))
	f := excelize.NewFile()
	sheetName1 := ExportSheetName(t)
	f.SetSheetName("Sheet1", sheetName1)
	streamWriter, _ := f.NewStreamWriter(sheetName1)
	currentRow := 1
	for col, header := range headers {
		colName, _ := excelize.ColumnNumberToName(col + 1)
		strHeader := fmt.Sprintf("%v", header)
		charWidth := len(strHeader) + 2   
		colWidth := float64(charWidth) * 1.2
		f.SetColWidth(sheetName1, colName, colName, colWidth)
	}

	err = streamWriter.SetRow(fmt.Sprintf("A%d", currentRow), headers)

	currentRow += 1

	err = streamWriter.SetRow(fmt.Sprintf("A%d", currentRow),subHeaders)

	currentRow += 1
	data := func() {
		values := make([]interface{}, 0, len(headers))
		for row, participant := range dataTestParticipants {
			values = values[:0]
			values = append(values,
				participant.KdProv,
				participant.KdRayon,
				participant.KdSek,
				participant.NPSN,
				participant.KdJenjang,
				participant.StsSek,
				participant.IDSekolah,
				participant.Kelas,
				participant.IDSiswa,
				participant.NamaSiswa,
				participant.PaketSoal,
				participant.Event,
				participant.Ujian,
			)
			for _, soal := range dataTestSoals {
				soalID := strconv.Itoa(soal.SoalID)
				participantID := strconv.Itoa(participant.IDSiswa)
				soalType :=soal.SoalTypeID
				val := 0
				if(soalType == 1) {
					val = int(jawabanKunciPG[soalID].(map[string]interface{})["opsNum"].(float64))
				} else if (soalType == 3){
					val = int(jawabanKunciPC[soalID].(map[string]interface{})["opsNum"].(float64))
				} else if (soalType == 8){
					val = int(jawabanKunciPGK1[soalID].(map[string]interface{})["opsNum"].(float64))
				} else if (soalType == 9){
					val = int(jawabanKunciPGK2[soalID].(map[string]interface{})["opsNum"].(float64))
				}
				num := int(val)
				ansPartiSoal := mapTestAnswers[participantID][soal.SoalUUID]
				if err == nil {
					for i := 0; i < num; i++ {

						if(ansPartiSoal == nil){
							values = append(values, utils.QUESTION_ASSIGN_SEE_NOT_ANSWER)
						} else {
							if(soalType == utils.QUESTION_TYPE_MULTIPLE_CHOICE_TABLE){
								answerRaw :=ansPartiSoal["siap_answer_uuid"]
								jawabanKunci := jawabanKunciPG[soalID].(map[string]interface{})["opsKunci"]
								if(answerRaw == "" ){
									values = append(values, utils.QUESTION_ASSIGN_NOT_SEE)
								} else {
									if(jawabanKunci.(map[string]interface{})[answerRaw.(string)] != nil){
										jawabOrder := jawabanKunci.(map[string]interface{})[answerRaw.(string)].(map[string]interface{})["order"]
										values = append(values, alphabet[int(jawabOrder.(float64))-1])
									} 
								}
							} else {
								answerRaw :=ansPartiSoal["answer"]
								if(answerRaw == nil){
									values = append(values, utils.QUESTION_ASSIGN_NOT_SEE)
								} else {
									if(soalType == utils.QUESTION_TYPE_MATCHING_TABLE){
										jawabKunci := jawabanKunciPC[soalID].(map[string]interface{})["opsKunci"].([]interface{})[i]
										var keyJawabKunci string 
										var valJawabKunci interface{}
										for k, v := range jawabKunci.(map[string]interface{}) {
											keyJawabKunci = k
											valJawabKunci = v
											break 
										}
									
										if _, ok := answerRaw.(map[string]interface{}); ok {
											jawab:=answerRaw.(map[string]interface{})[valJawabKunci.(string)]
											if(jawab == keyJawabKunci){
												values = append(values, utils.QUESTION_ASSIGN_CORRECT)
											} else {
												values = append(values, utils.QUESTION_ASSIGN_NOT_CORRECT)
											}
										} else if _, ok := answerRaw.(string); ok {
											values = append(values, utils.QUESTION_ASSIGN_NOT_SEE)
										} else {
											values = append(values, utils.QUESTION_ASSIGN_NOT_SEE)
										}
										
									} else if (soalType == utils.QUESTION_TYPE_COMPLEX_MULTIPLE_CHOICE_MCMA_TABLE){
										jawabKunci := jawabanKunciPGK1[soalID].(map[string]interface{})["opsKunci"].([]interface{})[i]
										var keyJawabKunci string 
										var valJawabKunci interface{} 
										for k, v := range jawabKunci.(map[string]interface{}) {
											keyJawabKunci = k
											valJawabKunci = v
											break 
										}
									
										if _, ok := answerRaw.(map[string]interface{}); ok {
											if(valJawabKunci != nil){
												jawab:=answerRaw.(map[string]interface{})[valJawabKunci.(string)]
												if(jawab != nil){
													values = append(values, utils.QUESTION_ASSIGN_CORRECT)
												} else {
													values = append(values,utils.QUESTION_ASSIGN_NOT_CORRECT)
												}
											} else {
												jawab:=answerRaw.(map[string]interface{})[keyJawabKunci]
												if(jawab != nil){
													values = append(values, utils.QUESTION_ASSIGN_NOT_CORRECT)
												} else {
													values = append(values, utils.QUESTION_ASSIGN_CORRECT)
												}
											}
										} else if _, ok := answerRaw.(string); ok {
											values = append(values, utils.QUESTION_ASSIGN_NOT_SEE)
										} else {
											values = append(values, utils.QUESTION_ASSIGN_NOT_SEE)
										}
										
									} else if (soalType == utils.QUESTION_TYPE_COMPLEX_MULTIPLE_CHOICE_RADIO_TABLE) {
										jawabKunci := jawabanKunciPGK2[soalID].(map[string]interface{})["opsKunci"].([]interface{})[i]
										var keyJawabKunci string 
										var valJawabKunci interface{} 
										for k, v := range jawabKunci.(map[string]interface{}) {
											keyJawabKunci = k
											valJawabKunci = v
											break 
										}
							
										if _, ok := answerRaw.(map[string]interface{}); ok {
											jawab:=answerRaw.(map[string]interface{})[keyJawabKunci]
											if (jawab == valJawabKunci){	
												values = append(values, utils.QUESTION_ASSIGN_CORRECT)
											} else {
												values = append(values, utils.QUESTION_ASSIGN_NOT_CORRECT)
											}
										} else if _, ok := answerRaw.(string); ok {
											values = append(values, utils.QUESTION_ASSIGN_NOT_SEE)
										} else {
											values = append(values, utils.QUESTION_ASSIGN_NOT_SEE)
										}
										
									}
								}
							}
						}		
					}
				}
			}
			err = streamWriter.SetRow(fmt.Sprintf("A%d", row+currentRow), values)
		}
	}

	data()
	_ = streamWriter.Flush()

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}

	_, err = utils.UploadToMinio(fileName, &buf,uuid)
	if err != nil {
		return nil, err
	}

	err = models.ReportOutcomes{}.UpdateReportOutcomes(fileName, uuid, 4)
	if err != nil {
		return nil, err
	}

	return &ReportResult{
		FileName: fileName,
		UUID:     uuid,
		Status:   "success",
	}, nil
}

func (ReportController) ScoreOption(c *gin.Context) {
	defer utils.EndProccess()
	uuid := c.Param("uuid")

	result, err := GenerateScoreReport(uuid)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   result.Status,
		"reqUUID":  result.UUID,
		"fileName": result.FileName,
	})
}

func GenerateQuestionReport(uuid string) (*ReportResult, error) {
	repPackages := models.ReportPackages{}
	dataRepPackages, err := repPackages.GetByReqUUID(uuid)
	if err != nil {
		return nil, err
	}

	repPackagesIDs := models.QuestionPackage{}
	dataRepPackagesIDs, _:= repPackagesIDs.GetByUUIDs(dataRepPackages)

	var packageIDs []int
	for _, packageID := range dataRepPackagesIDs {
		packageIDs = append(packageIDs, int(packageID.ID))
	}

	testParticipants := models.TestParticipantReport{}
	dataTestParticipants, err := testParticipants.GetParticipantAnswer(packageIDs)

	testSoals := models.QuestionPackageSoalReport{}
	testSoalsType := []int{
		utils.QUESTION_TYPE_MULTIPLE_CHOICE_TABLE,
		utils.QUESTION_TYPE_MATCHING_TABLE,
		utils.QUESTION_TYPE_COMPLEX_MULTIPLE_CHOICE_MCMA_TABLE,
		utils.QUESTION_TYPE_COMPLEX_MULTIPLE_CHOICE_RADIO_TABLE,
	}
	dataTestSoals, err := testSoals.GetQuestionPackageSoalReport(dataRepPackages,testSoalsType)

	testAnswers := models.TestAnswers{}
	dataTestAnswers, err := testAnswers.GetDataExport(dataRepPackages,[]string{})

	mapTestAnswers := ConvAnsToMap(dataTestAnswers)

	var soalID1 []int
	var soalID3 []int
	var soalID8 []int
	var soalID9 []int
	for i := 0; i < len(dataTestSoals); i++ {
		if(dataTestSoals[i].SoalTypeID == utils.QUESTION_TYPE_MULTIPLE_CHOICE_TABLE ){
			soalID1 = append(soalID1, dataTestSoals[i].SoalID)
		} else if (dataTestSoals[i].SoalTypeID == utils.QUESTION_TYPE_MATCHING_TABLE) {
			soalID3 = append(soalID3, dataTestSoals[i].SoalID)
		} else if (dataTestSoals[i].SoalTypeID == utils.QUESTION_TYPE_COMPLEX_MULTIPLE_CHOICE_MCMA_TABLE) {
			soalID8 = append(soalID8, dataTestSoals[i].SoalID)
		} else if (dataTestSoals[i].SoalTypeID == utils.QUESTION_TYPE_COMPLEX_MULTIPLE_CHOICE_RADIO_TABLE) {
			soalID9 = append(soalID9, dataTestSoals[i].SoalID)
		} 
	}

	jawabanKunci := models.JawabanKunci{}
	var jawabanKunciPG,jawabanKunciPC, jawabanKunciPGK1, jawabanKunciPGK2 map[string]interface{}
	if(len(soalID1) > 0){
		jawabanKunciPG,err = jawabanKunci.GetJawabanKunciPG(soalID1)
	}

	if(len(soalID3) > 0){
		jawabanKunciPC,err = jawabanKunci.GetJawabanKunciPC(soalID3)
	}

	if(len(soalID8) > 0){
		jawabanKunciPGK1,err = jawabanKunci.GetJawabanKunciPGK1(soalID8)
	}

	if(len(soalID9) > 0){
		jawabanKunciPGK2,err = jawabanKunci.GetJawabanKunciPGK2(soalID9)
	}

	var headers []interface{}
    var subHeaders []interface{}
    for _, h := range fixedHeaders {
        headers = append(headers, h)
        subHeaders = append(subHeaders,"")
    }
	soalTypeConstant := utils.SoalTypeConstant{}

	for _, soal := range dataTestSoals {
		headers = append(headers, soal.SoalCode)
		subHeaders = append(subHeaders,soalTypeConstant.SoalType(soal.SoalTypeID))
	}

	t := time.Now().Format("02012006_1504") // dmY_Hi
	fileName :=  ReportFileName(utils.REPORT_FIFTH, t, utils.ShortUUID(uuid))
	f := excelize.NewFile()
	sheetName1 := ExportSheetName(t)
	f.SetSheetName("Sheet1", sheetName1)
	streamWriter, _ := f.NewStreamWriter(sheetName1)
	currentRow := 1

	
	for col, header := range headers {
		colName, _ := excelize.ColumnNumberToName(col + 1)
        strHeader := fmt.Sprintf("%v", header)
        charWidth := len(strHeader) + 2   
        colWidth := float64(charWidth) * 1.2
        f.SetColWidth(sheetName1, colName, colName, colWidth)
	}

	err = streamWriter.SetRow(fmt.Sprintf("A%d", currentRow), headers)

    currentRow += 1

    err = streamWriter.SetRow(fmt.Sprintf("A%d", currentRow),subHeaders)

    currentRow += 1

	data := func() {
		values := make([]interface{}, 0, len(headers))
		for row, participant := range dataTestParticipants {
			values = values[:0]
			values = append(values,
				participant.KdProv,
				participant.KdRayon,
				participant.KdSek,
				participant.NPSN,
				participant.KdJenjang,
				participant.StsSek,
				participant.IDSekolah,
				participant.Kelas,
				participant.IDSiswa,
				participant.NamaSiswa,
				participant.PaketSoal,
				participant.Event,
				participant.Ujian,
			)
			for _, soal := range dataTestSoals {
				soalID := strconv.Itoa(soal.SoalID)
				participantID := strconv.Itoa(participant.IDSiswa)
				soalType :=soal.SoalTypeID
				val := 0
                if(soalType == 1) {
                    val = int(jawabanKunciPG[soalID].(map[string]interface{})["opsNum"].(float64))
                } else if (soalType == 3){
                    val = int(jawabanKunciPC[soalID].(map[string]interface{})["opsNum"].(float64))
                } else if (soalType == 8){
                    val = int(jawabanKunciPGK1[soalID].(map[string]interface{})["opsNum"].(float64))
                } else if (soalType == 9){
                    val = int(jawabanKunciPGK2[soalID].(map[string]interface{})["opsNum"].(float64))
                }
				num := int(val)
				status := utils.QUESTION_ASSIGN_SEE_NOT_ANSWER
				allTrue := true
				ansPartiSoal := mapTestAnswers[participantID][soal.SoalUUID]
				if err == nil {
					for i := 0; i < num; i++ {
						if(ansPartiSoal == nil){
							//values = append(values, 8)
						} else {
							if(soalType == utils.QUESTION_TYPE_MULTIPLE_CHOICE_TABLE ){
								answerRaw :=ansPartiSoal["siap_answer_uuid"]
								jawabanKunci := jawabanKunciPG[soalID].(map[string]interface{})["opsKunci"]
								if(answerRaw == "" ){
									status = utils.QUESTION_ASSIGN_NOT_SEE
									//values = append(values, 9)
								} else {
									if(jawabanKunci.(map[string]interface{})[answerRaw.(string)] != nil){
										jawabOrder := jawabanKunci.(map[string]interface{})[answerRaw.(string)].(map[string]interface{})["order"]
										status = int(jawabOrder.(float64))
										//values = append(values, alphabet[int(jawabOrder.(float64))-1])
									} 
								}
							} else {
								answerRaw :=ansPartiSoal["answer"]
								if(answerRaw == nil){
									status = utils.QUESTION_ASSIGN_NOT_SEE
									//values = append(values, 9)
								} else {
									status = utils.QUESTION_ASSIGN_SEE_NOT_ANSWER_NOT_SEE
									if(soalType == utils.QUESTION_TYPE_MATCHING_TABLE){
										jawabKunci := jawabanKunciPC[soalID].(map[string]interface{})["opsKunci"].([]interface{})[i]
										var keyJawabKunci string 
										var valJawabKunci interface{}
										for k, v := range jawabKunci.(map[string]interface{}) {
											keyJawabKunci = k
											valJawabKunci = v
											break 
										}
										if _, ok := answerRaw.(map[string]interface{}); ok {
											jawab:=answerRaw.(map[string]interface{})[valJawabKunci.(string)]
											if(jawab == keyJawabKunci){
												//values = append(values, 1)
											} else {
												allTrue = false
												//values = append(values, 0)
											}
										} else if _, ok := answerRaw.(string); ok {
											status = utils.QUESTION_ASSIGN_NOT_SEE
										} else {
											status = utils.QUESTION_ASSIGN_NOT_SEE
										}
									} else if (soalType == utils.QUESTION_TYPE_COMPLEX_MULTIPLE_CHOICE_MCMA_TABLE){
										jawabKunci := jawabanKunciPGK1[soalID].(map[string]interface{})["opsKunci"].([]interface{})[i]
										var keyJawabKunci string 
										var valJawabKunci interface{} 
										for k, v := range jawabKunci.(map[string]interface{}) {
											keyJawabKunci = k
											valJawabKunci = v
											break 
										}
										if _, ok := answerRaw.(map[string]interface{}); ok {
											if(valJawabKunci != nil){
												jawab:=answerRaw.(map[string]interface{})[valJawabKunci.(string)]
												if(jawab != nil){
													//values = append(values, 1)
												} else {
													allTrue = false
													//values = append(values, 0)
												}
											} else {
												jawab:=answerRaw.(map[string]interface{})[keyJawabKunci]
												if(jawab != nil){
													allTrue = false
													//values = append(values, 0)
												} else {
													//values = append(values, 1)
												}
											}
										} else if _, ok := answerRaw.(string); ok {
											status = utils.QUESTION_ASSIGN_NOT_SEE
										} else {
											status = utils.QUESTION_ASSIGN_NOT_SEE
										}
										
									} else if (soalType == utils.QUESTION_TYPE_COMPLEX_MULTIPLE_CHOICE_RADIO_TABLE) {
										jawabKunci := jawabanKunciPGK2[soalID].(map[string]interface{})["opsKunci"].([]interface{})[i]
										var keyJawabKunci string 
										var valJawabKunci interface{} 
										for k, v := range jawabKunci.(map[string]interface{}) {
											keyJawabKunci = k
											valJawabKunci = v
											break 
										}
										if _, ok := answerRaw.(map[string]interface{}); ok {
											jawab:=answerRaw.(map[string]interface{})[keyJawabKunci]
											if (jawab == valJawabKunci){	
												//values = append(values, 1)
											} else {
												allTrue = false
												//values = append(values, 0)
											}
										} else if _, ok := answerRaw.(string); ok {
											status = utils.QUESTION_ASSIGN_NOT_SEE
										} else {
											status = utils.QUESTION_ASSIGN_NOT_SEE
										}
									}
								}
							}
						}		
					}
				}
				if(soalType == utils.QUESTION_TYPE_MULTIPLE_CHOICE_TABLE){
					if(status >=  utils.QUESTION_ASSIGN_SEE_NOT_ANSWER){
						values = append(values, status)
					} else {
						values = append(values, alphabet[status-1])
					}
					
				} else {
					if(status == utils.QUESTION_ASSIGN_SEE_NOT_ANSWER_NOT_SEE) {
						if(allTrue == true){
							values = append(values, utils.QUESTION_ASSIGN_CORRECT)
						} else {
							values = append(values, utils.QUESTION_ASSIGN_NOT_CORRECT)
						}
					} else {
						values = append(values, status)
					}
					
				}
				
			}
			err = streamWriter.SetRow(fmt.Sprintf("A%d", row+currentRow), values)
		}
	}

	data()
    _ = streamWriter.Flush()

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}

	_, err = utils.UploadToMinio(fileName, &buf,uuid)
	if err != nil {
		return nil, err
	}

	err = models.ReportOutcomes{}.UpdateReportOutcomes(fileName, uuid, 5)
	if err != nil {
		return nil, err
	}

	return &ReportResult{
		FileName: fileName,
		UUID:     uuid,
		Status:   "success",
	}, nil
}

func (ReportController) ScoreQuestion(c *gin.Context) {
	defer utils.EndProccess()
	uuid := c.Param("uuid")

	result, err := GenerateQuestionReport(uuid)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   result.Status,
		"reqUUID":  result.UUID,
		"fileName": result.FileName,
	})
}

func GenerateRawNonEssayShortReport(uuid string) (*ReportResult, error) {
	repPackages := models.ReportPackages{}
	dataRepPackages, err := repPackages.GetByReqUUID(uuid)
	if err != nil {
		return nil, err
	}

	repPackagesIDs := models.QuestionPackage{}
	dataRepPackagesIDs, _:= repPackagesIDs.GetByUUIDs(dataRepPackages)

	var packageIDs []int
	for _, packageID := range dataRepPackagesIDs {
		packageIDs = append(packageIDs, int(packageID.ID))
	}

	testParticipants := models.TestParticipantReport{}
	dataTestParticipants, err := testParticipants.GetParticipantAnswer(packageIDs)

	testSoals := models.QuestionPackageSoalReport{}
	testSoalsType := []int{
		utils.QUESTION_TYPE_MULTIPLE_CHOICE_TABLE,
		utils.QUESTION_TYPE_MATCHING_TABLE,
		utils.QUESTION_TYPE_COMPLEX_MULTIPLE_CHOICE_MCMA_TABLE,
		utils.QUESTION_TYPE_COMPLEX_MULTIPLE_CHOICE_RADIO_TABLE,
	}
	dataTestSoals, err := testSoals.GetQuestionPackageSoalReport(dataRepPackages,testSoalsType)

	testAnswers := models.TestAnswers{}
	dataTestAnswers, err := testAnswers.GetDataExport(dataRepPackages,[]string{})

	mapTestAnswers := ConvAnsToMap(dataTestAnswers)

	var soalID1 []int
	var soalID3 []int
	var soalID8 []int
	var soalID9 []int
	for i := 0; i < len(dataTestSoals); i++ {
		if(dataTestSoals[i].SoalTypeID == utils.QUESTION_TYPE_MULTIPLE_CHOICE_TABLE ){
			soalID1 = append(soalID1, dataTestSoals[i].SoalID)
		} else if (dataTestSoals[i].SoalTypeID == utils.QUESTION_TYPE_MATCHING_TABLE) {
			soalID3 = append(soalID3, dataTestSoals[i].SoalID)
		} else if (dataTestSoals[i].SoalTypeID == utils.QUESTION_TYPE_COMPLEX_MULTIPLE_CHOICE_MCMA_TABLE) {
			soalID8 = append(soalID8, dataTestSoals[i].SoalID)
		} else if (dataTestSoals[i].SoalTypeID == utils.QUESTION_TYPE_COMPLEX_MULTIPLE_CHOICE_RADIO_TABLE) {
			soalID9 = append(soalID9, dataTestSoals[i].SoalID)
		} 
	}

	jawabanKunci := models.JawabanKunci{}
	var jawabanKunciPG,jawabanKunciPC, jawabanKunciPGK1, jawabanKunciPGK2 map[string]interface{}
	if(len(soalID1) > 0){
		jawabanKunciPG,err = jawabanKunci.GetJawabanKunciPG(soalID1)
	}

	if(len(soalID3) > 0){
		jawabanKunciPC,err = jawabanKunci.GetJawabanKunciPC(soalID3)
	}

	if(len(soalID8) > 0){
		jawabanKunciPGK1,err = jawabanKunci.GetJawabanKunciPGK1(soalID8)
	}

	if(len(soalID9) > 0){
		jawabanKunciPGK2,err = jawabanKunci.GetJawabanKunciPGK2(soalID9)
	}

	lenSoalOps := 0

	for _, soal := range dataTestSoals {
		soalID := strconv.Itoa(soal.SoalID)
		soalType := soal.SoalTypeID
        val := 0
        if(soalType == 1) {
            val = int(jawabanKunciPG[soalID].(map[string]interface{})["opsNum"].(float64))
        } else if (soalType == 3){
            val = int(jawabanKunciPC[soalID].(map[string]interface{})["opsNum"].(float64))
        } else if (soalType == 8){
            val = int(jawabanKunciPGK1[soalID].(map[string]interface{})["opsNum"].(float64))
        } else if (soalType == 9){
            val = int(jawabanKunciPGK2[soalID].(map[string]interface{})["opsNum"].(float64))
        }
		num := int(val)
		if err == nil {
			lenSoalOps += num
		}
	}

	var headers []interface{}
    var subHeaders []interface{}
    for _, h := range fixedHeaders {
        headers = append(headers, h)
        subHeaders = append(subHeaders,"")
    }

	soalTypeConstant := utils.SoalTypeConstant{}

	for _, soal := range dataTestSoals {
		soalID := strconv.Itoa(soal.SoalID)
		soalType := soal.SoalTypeID
        val := 0
        if(soalType == 1) {
            val = int(jawabanKunciPG[soalID].(map[string]interface{})["opsNum"].(float64))
        } else if (soalType == 3){
            val = int(jawabanKunciPC[soalID].(map[string]interface{})["opsNum"].(float64))
        } else if (soalType == 8){
            val = int(jawabanKunciPGK1[soalID].(map[string]interface{})["opsNum"].(float64))
        } else if (soalType == 9){
            val = int(jawabanKunciPGK2[soalID].(map[string]interface{})["opsNum"].(float64))
        }
		num := int(val)
		if err == nil {
			for i := 0; i < num; i++ {
				headers = append(headers, soal.SoalCode + " " +strconv.Itoa(i+1))
                subHeaders = append(subHeaders,soalTypeConstant.SoalType(soal.SoalTypeID))
			}
		}
	}

	t := time.Now().Format("02012006_1504") // dmY_Hi
	fileName :=  ReportFileName(utils.REPORT_SECOND, t, utils.ShortUUID(uuid))
	f := excelize.NewFile()
	sheetName1 := ExportSheetName(t)
	f.SetSheetName("Sheet1", sheetName1)
	streamWriter, _ := f.NewStreamWriter(sheetName1)
	currentRow := 1

	for col, header := range headers {
		colName, _ := excelize.ColumnNumberToName(col + 1)
        strHeader := fmt.Sprintf("%v", header)
        charWidth := len(strHeader) + 2   
        colWidth := float64(charWidth) * 1.2
        f.SetColWidth(sheetName1, colName, colName, colWidth)
	}

	err = streamWriter.SetRow(fmt.Sprintf("A%d", currentRow), headers)

    currentRow += 1

    err = streamWriter.SetRow(fmt.Sprintf("A%d", currentRow),subHeaders)

    currentRow += 1

	subHeader := func() {
		values := make([]interface{}, 0, len(headers))
		values = values[:0]
		values = append(values,
				"",
				"",
				"",
				"",
				"",
				"",
				"",
				"",
				"",
				"",
				"",
				"",
				"",
		)
		for _, soal := range dataTestSoals {
				soalID := strconv.Itoa(soal.SoalID)
				soalType :=soal.SoalTypeID
				 val := 0
                if(soalType == 1) {
                    val = int(jawabanKunciPG[soalID].(map[string]interface{})["opsNum"].(float64))
                } else if (soalType == 3){
                    val = int(jawabanKunciPC[soalID].(map[string]interface{})["opsNum"].(float64))
                } else if (soalType == 8){
                    val = int(jawabanKunciPGK1[soalID].(map[string]interface{})["opsNum"].(float64))
                } else if (soalType == 9){
                    val = int(jawabanKunciPGK2[soalID].(map[string]interface{})["opsNum"].(float64))
                }
				num := int(val)
				if err == nil {
					for i := 0; i < num; i++ {
						if(soalType == utils.QUESTION_TYPE_MULTIPLE_CHOICE_TABLE){
							data:=jawabanKunciPG[soalID].(map[string]interface{})["opsKunci"].(map[string]interface{})
							var kunci int
							for _, innerMap := range data {
								if answer, ok := innerMap.(map[string]interface{})["answer"]; ok && answer != nil {
									kunci = int(innerMap.(map[string]interface{})["order"].(float64))
									break
								}
							}
							values = append(values, alphabet[kunci-1])
						} else {
							if(soalType == utils.QUESTION_TYPE_MATCHING_TABLE){
								jawab := jawabanKunciPC[soalID].(map[string]interface{})
								jawabKunci := jawab["opsKunci"].([]interface{})[i]
								jawabOrder := jawab["opsOrder"].(map[string]interface{})
								var keyJawabKunci string 
								var valJawabKunci interface{}
								for k, v := range jawabKunci.(map[string]interface{}) {
									keyJawabKunci = k
									valJawabKunci = v
									break 
								}
								values = append(values, fmt.Sprintf("L%vR%v",jawabOrder[valJawabKunci.(string)] , jawabOrder[keyJawabKunci]))
						
							} else if (soalType == utils.QUESTION_TYPE_COMPLEX_MULTIPLE_CHOICE_MCMA_TABLE){
								
								jawabKunci := jawabanKunciPGK1[soalID].(map[string]interface{})["opsKunci"].([]interface{})[i]
								
								var valJawabKunci interface{}
								for _, v := range jawabKunci.(map[string]interface{}) {
									valJawabKunci = v
									break 
								}
								
								if(valJawabKunci != nil){
									values = append(values, fmt.Sprintf("%vA", i+1))
								} else {
									values = append(values, "")
								}
								
							} else if (soalType == utils.QUESTION_TYPE_COMPLEX_MULTIPLE_CHOICE_RADIO_TABLE) {
								jawab := jawabanKunciPGK2[soalID].(map[string]interface{})
								jawabKunci := jawab["opsKunci"].([]interface{})[i]
								jawabOrder := jawab["opsOrder"].(map[string]interface{})
								
								var valJawabKunci interface{}
								for _, v:= range jawabKunci.(map[string]interface{}) {
									valJawabKunci = v
									break 
								}
								values = append(values, fmt.Sprintf("%v%v",alphabet[int(jawabOrder[valJawabKunci.(string)].(float64))-1] ,i+1))

							}
						}	
					}
				}
			}
		err = streamWriter.SetRow(fmt.Sprintf("A%d", currentRow), values)
	}

	data := func() {
		values := make([]interface{}, 0, len(headers))
		for row, participant := range dataTestParticipants {
			values = values[:0]
			values = append(values,
				participant.KdProv,
				participant.KdRayon,
				participant.KdSek,
				participant.NPSN,
				participant.KdJenjang,
				participant.StsSek,
				participant.IDSekolah,
				participant.Kelas,
				participant.IDSiswa,
				participant.NamaSiswa,
				participant.PaketSoal,
				participant.Event,
				participant.Ujian,
			)
			for _, soal := range dataTestSoals {
				soalID := strconv.Itoa(soal.SoalID)
				participantID := strconv.Itoa(participant.IDSiswa)
				soalType :=soal.SoalTypeID
				val := 0
                if(soalType == 1) {
                    val = int(jawabanKunciPG[soalID].(map[string]interface{})["opsNum"].(float64))
                } else if (soalType == 3){
                    val = int(jawabanKunciPC[soalID].(map[string]interface{})["opsNum"].(float64))
                } else if (soalType == 8){
                    val = int(jawabanKunciPGK1[soalID].(map[string]interface{})["opsNum"].(float64))
                } else if (soalType == 9){
                    val = int(jawabanKunciPGK2[soalID].(map[string]interface{})["opsNum"].(float64))
                }
				num := int(val)
				ansPartiSoal := mapTestAnswers[participantID][soal.SoalUUID]
				if err == nil {
					for i := 0; i < num; i++ {
						if(ansPartiSoal == nil){
							values = append(values, utils.QUESTION_ASSIGN_SEE_NOT_ANSWER)
						} else {
							if(soalType == utils.QUESTION_TYPE_MULTIPLE_CHOICE_TABLE){
								answerRaw :=ansPartiSoal["siap_answer_uuid"]
                                jawabanKunci := jawabanKunciPG[soalID].(map[string]interface{})["opsKunci"]
								if(answerRaw == "" ){
									values = append(values, utils.QUESTION_ASSIGN_NOT_SEE)
								} else {
									if(jawabanKunci.(map[string]interface{})[answerRaw .(string)] != nil){
										jawabOrder := jawabanKunci.(map[string]interface{})[answerRaw .(string)].(map[string]interface{})["order"]
										values = append(values, alphabet[int(jawabOrder.(float64))-1])
									} 
								}
							} else {
								answerRaw :=ansPartiSoal["answer"]
								if(answerRaw == nil){
									values = append(values, utils.QUESTION_ASSIGN_NOT_SEE)
								} else {
									if(soalType == utils.QUESTION_TYPE_MATCHING_TABLE){
										jawab:=jawabanKunciPC[soalID].(map[string]interface{})
										jawabKunci := jawab["opsKunci"].([]interface{})[i]
										jawabOrder := jawab["opsOrder"].(map[string]interface{})
										var valJawabKunci interface{}
										for _, v := range jawabKunci.(map[string]interface{}) {
											valJawabKunci = v
											break 
										}
										if _, ok := answerRaw.(map[string]interface{}); ok {
											jawab:=answerRaw.(map[string]interface{})[valJawabKunci.(string)]
											if(jawab == nil) {
												values = append(values, utils.QUESTION_ASSIGN_NOT_SEE)
											} else {
												values = append(values, fmt.Sprintf("L%vR%v", jawabOrder[valJawabKunci.(string)], jawabOrder[jawab.(string)]))
											}
										} else if _, ok := answerRaw.(string); ok {
											values = append(values, utils.QUESTION_ASSIGN_NOT_SEE)
										} else {
											values = append(values, utils.QUESTION_ASSIGN_NOT_SEE)
										}
									} else if (soalType == utils.QUESTION_TYPE_COMPLEX_MULTIPLE_CHOICE_MCMA_TABLE){
										jawabKunci := jawabanKunciPGK1[soalID].(map[string]interface{})["opsKunci"].([]interface{})[i]
										var keyJawabKunci string 
										
										for k, _ := range jawabKunci.(map[string]interface{}) {
											keyJawabKunci = k
											break 
										}
										if _, ok := answerRaw.(map[string]interface{}); ok {
											jawab:=answerRaw.(map[string]interface{})[keyJawabKunci]
										
											if(jawab != nil){
												values = append(values, fmt.Sprintf("%vA", i+1))
											} else {
												values = append(values, "")
											}
										} else if _, ok := answerRaw.(string); ok {
											values = append(values, "")
										} else {
											values = append(values, "")
										}
									} else if (soalType == utils.QUESTION_TYPE_COMPLEX_MULTIPLE_CHOICE_RADIO_TABLE) {
										jawab:=jawabanKunciPGK2[soalID].(map[string]interface{})
										jawabKunci := jawab["opsKunci"].([]interface{})[i]
										jawabOrder := jawab["opsOrder"].(map[string]interface{})
										var keyJawabKunci string 
										for k, _:= range jawabKunci.(map[string]interface{}) {
											keyJawabKunci = k
											break 
										}
										if _, ok := answerRaw.(map[string]interface{}); ok {
											jawab:=answerRaw.(map[string]interface{})[keyJawabKunci]
											if(jawab == nil){
												values = append(values, utils.QUESTION_ASSIGN_NOT_SEE)
											} else {
												values = append(values, fmt.Sprintf("%v%v",alphabet[int(jawabOrder[jawab.(string)].(float64))-1] ,i+1))
											}
										} else if _, ok := answerRaw.(string); ok {
											values = append(values, utils.QUESTION_ASSIGN_NOT_SEE)
										} else {
											values = append(values, utils.QUESTION_ASSIGN_NOT_SEE)
										}
									}
								}
							}
						}		
					}
				}
			}
			err = streamWriter.SetRow(fmt.Sprintf("A%d", row+currentRow), values)
		}
	}

	subHeader()

	currentRow += 1
	data()
	_ = streamWriter.Flush()

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}

	_, err = utils.UploadToMinio(fileName, &buf,uuid)
	if err != nil {
		return nil, err
	}

	err = models.ReportOutcomes{}.UpdateReportOutcomes(fileName, uuid, 2)
	if err != nil {
		return nil, err
	}

	return &ReportResult{
		FileName: fileName,
		UUID:     uuid,
		Status:   "success",
	}, nil
}

func (ReportController) RawNonEssayShort(c *gin.Context) {
	defer utils.EndProccess()	
	uuid := c.Param("uuid")

	result, err := GenerateRawNonEssayShortReport(uuid)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   result.Status,
		"reqUUID":  result.UUID,
		"fileName": result.FileName,
	})
}

func GenerateRawEssayShortReport(uuid string) (*ReportResult, error) {
	repPackages := models.ReportPackages{}
	dataRepPackages, err := repPackages.GetByReqUUID(uuid)
	if err != nil {
		return nil, err
	}

	repPackagesIDs := models.QuestionPackage{}
	dataRepPackagesIDs, _:= repPackagesIDs.GetByUUIDs(dataRepPackages)

	var packageIDs []int
	for _, packageID := range dataRepPackagesIDs {
		packageIDs = append(packageIDs, int(packageID.ID))
	}

	testParticipants := models.TestParticipantReport{}
	dataTestParticipants, err := testParticipants.GetParticipantAnswer(packageIDs)

	testSoals := models.QuestionPackageSoalReport{}
	testSoalsType := []int{
		utils.QUESTION_TYPE_ESSAY_TABLE,
		utils.QUESTION_TYPE_CLOZE_PROCEDURE_TABLE,
	}
	dataTestSoals, err := testSoals.GetQuestionPackageSoalReport(dataRepPackages,testSoalsType)
	

	testAnswers := models.TestAnswers{}
	dataTestAnswers, err := testAnswers.GetDataExport(dataRepPackages,[]string{})

	mapTestAnswers := ConvAnsToMap(dataTestAnswers)

	var soalID4 []int
	for i := 0; i < len(dataTestSoals); i++ {
		if(dataTestSoals[i].SoalTypeID == utils.QUESTION_TYPE_CLOZE_PROCEDURE_TABLE  ){
			soalID4 = append(soalID4, dataTestSoals[i].SoalID)
		}
	}

	jawabanKunci := models.JawabanKunci{}
	var jawabanKunciShort map[string]interface{}
	if(len(soalID4) > 0){
		jawabanKunciShort,err = jawabanKunci.GetJawabanKunciShort(soalID4)
	}

	lenSoalOps := 0

	for _, soal := range dataTestSoals {
		soalID := strconv.Itoa(soal.SoalID)
		soalType := soal.SoalTypeID
		if(soalType == utils.QUESTION_TYPE_ESSAY_TABLE){
			lenSoalOps += 1
		} else if (soalType == utils.QUESTION_TYPE_CLOZE_PROCEDURE_TABLE){
			val := jawabanKunciShort[soalID].(map[string]interface{})["opsNum"].(float64)
			num := int(val)
			if err == nil {
				lenSoalOps += num
			}
		}
	}

	var headers []interface{}
    var subHeaders []interface{}
    for _, h := range fixedHeaders {
        headers = append(headers, h)
        subHeaders = append(subHeaders,"")
    }

	soalTypeConstant := utils.SoalTypeConstant{}

	for _, soal := range dataTestSoals {
		soalID := strconv.Itoa(soal.SoalID)
		soalType := soal.SoalTypeID
		if(soalType == utils.QUESTION_TYPE_ESSAY_TABLE){
			headers = append(headers, soal.SoalCode + " 1")
            subHeaders = append(subHeaders,soalTypeConstant.SoalType(soal.SoalTypeID))
		} else if (soalType == utils.QUESTION_TYPE_CLOZE_PROCEDURE_TABLE){
			val := jawabanKunciShort[soalID].(map[string]interface{})["opsNum"].(float64)
			num := int(val)
			if err == nil {
				for i := 0; i < num; i++ {
					headers = append(headers, soal.SoalCode + " " +strconv.Itoa(i+1))
					subHeaders = append(subHeaders,soalTypeConstant.SoalType(soal.SoalTypeID))
				}
			}
		}
	}

	t := time.Now().Format("02012006_1504") // dmY_Hi
	fileName :=  ReportFileName(utils.REPORT_ONE, t, utils.ShortUUID(uuid))
	f := excelize.NewFile()
	sheetName1 := ExportSheetName(t)
	f.SetSheetName("Sheet1", sheetName1)
	streamWriter, _ := f.NewStreamWriter(sheetName1)
	currentRow := 1

	for col, header := range headers {
		colName, _ := excelize.ColumnNumberToName(col + 1)
        strHeader := fmt.Sprintf("%v", header)
        charWidth := len(strHeader) + 2   
        colWidth := float64(charWidth) * 1.2
        f.SetColWidth(sheetName1, colName, colName, colWidth)
	}

	err = streamWriter.SetRow(fmt.Sprintf("A%d", currentRow), headers)

    currentRow += 1

    err = streamWriter.SetRow(fmt.Sprintf("A%d", currentRow),subHeaders)

	data := func() {
		values := make([]interface{}, 0, len(headers))
		for row, participant := range dataTestParticipants {
			values = values[:0]
			values = append(values,
				participant.KdProv,
				participant.KdRayon,
				participant.KdSek,
				participant.NPSN,
				participant.KdJenjang,
				participant.StsSek,
				participant.IDSekolah,
				participant.Kelas,
				participant.IDSiswa,
				participant.NamaSiswa,
				participant.PaketSoal,
				participant.Event,
				participant.Ujian,
			)
			for _, soal := range dataTestSoals {
				soalID := strconv.Itoa(soal.SoalID)
				participantID := strconv.Itoa(participant.IDSiswa)
				soalType :=soal.SoalTypeID
				var num int
				if(soal.SoalTypeID ==  utils.QUESTION_TYPE_ESSAY_TABLE){
					num = 1
				} else if (soal.SoalTypeID == utils.QUESTION_TYPE_CLOZE_PROCEDURE_TABLE){
					val := jawabanKunciShort[soalID].(map[string]interface{})["opsNum"].(float64)
					num = int(val)
				}
				ansPartiSoal := mapTestAnswers[participantID][soal.SoalUUID]
				for i := 0; i < num; i++ {
					if(ansPartiSoal == nil){
						values = append(values, utils.QUESTION_ASSIGN_SEE_NOT_ANSWER)	
					} else {
						answerRaw :=ansPartiSoal["answer"]
						if(soalType == utils.QUESTION_TYPE_ESSAY_TABLE){
							if(answerRaw == nil || answerRaw == ""){
								values = append(values, utils.QUESTION_ASSIGN_NOT_SEE)
							} else {
								ans := answerRaw
								re := regexp.MustCompile(`<[^>]+>`)
								str, ok := ans.(string)
								if ok {
									val := re.ReplaceAllString(str, "")
									values = append(values, val)
								} else {
									values = append(values, ans)
								}
							}
						} else if(soalType == utils.QUESTION_TYPE_CLOZE_PROCEDURE_TABLE){
							if(answerRaw == nil || answerRaw == ""){
								values = append(values, utils.QUESTION_ASSIGN_NOT_SEE)
							} else {
								jawabKunci := jawabanKunciShort[soalID].(map[string]interface{})["opsKunci"].([]interface{})[i]
								var keyJawabKunci string 
								for k, _ := range jawabKunci.(map[string]interface{}) {
									keyJawabKunci = k
									break 
								}
								jawab:=answerRaw.(map[string]interface{})[keyJawabKunci]
								values = append(values, jawab)
							}
						}
					}
				}
				
			}
			err = streamWriter.SetRow(fmt.Sprintf("A%d", row+currentRow), values)
		}
	}

	currentRow += 1

	data()
	_ = streamWriter.Flush()

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}

	_, err = utils.UploadToMinio(fileName, &buf,uuid)
	if err != nil {
		return nil, err
	}

	err = models.ReportOutcomes{}.UpdateReportOutcomes(fileName, uuid, 1)
	if err != nil {
		return nil, err
	}

	return &ReportResult{
		FileName: fileName,
		UUID:     uuid,
		Status:   "success",
	}, nil
}

func (ReportController) RawEssayShort(c *gin.Context) {
	defer utils.EndProccess()	
	uuid := c.Param("uuid")

	result, err := GenerateRawEssayShortReport(uuid)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   result.Status,
		"reqUUID":  result.UUID,
		"fileName": result.FileName,
	})
}

func GenerateResponseJawabanReport(uuid string) (*ReportResult, error) {
	repPackages := models.ReportPackages{}
	dataRepPackages, err := repPackages.GetByReqUUID(uuid)
	if err != nil {
		return nil, err
	}

	repPackagesIDs := models.QuestionPackage{}
	dataRepPackagesIDs, _:= repPackagesIDs.GetByUUIDs(dataRepPackages)

	var packageIDs []int
	for _, packageID := range dataRepPackagesIDs {
		packageIDs = append(packageIDs, int(packageID.ID))
	}
	
	testParticipants := models.TestParticipantReport{}
	dataTestParticipants, err := testParticipants.GetParticipantAnswer(packageIDs)

	var siswaID []int
	for _, participant := range dataTestParticipants {
		siswaID = append(siswaID, participant.IDSiswa)
	}

	testSoals := models.QuestionPackageSoalReport{}
	testSoalsType := []int{
		utils.QUESTION_TYPE_MULTIPLE_CHOICE_TABLE,
		utils.QUESTION_TYPE_MATCHING_TABLE,
		utils.QUESTION_TYPE_COMPLEX_MULTIPLE_CHOICE_MCMA_TABLE,
		utils.QUESTION_TYPE_COMPLEX_MULTIPLE_CHOICE_RADIO_TABLE,
		utils.QUESTION_TYPE_CLOZE_PROCEDURE_TABLE,
		utils.QUESTION_TYPE_ESSAY_TABLE,
	}
	dataTestSoals, err := testSoals.GetQuestionPackageSoalReport(dataRepPackages,testSoalsType)

	testAnswers := models.TestAnswers{}
	dataTestAnswers, err := testAnswers.GetDataExport(dataRepPackages,[]string{})

	participantLog := models.TestParticipantLogReport{}
	testParticipantLog,err := participantLog.GetParticipantLog(siswaID)

	testParticipantLogMap := make(map[string]map[string]map[string]interface{})

	for _, log := range testParticipantLog {
		testParticipantID := strconv.Itoa(int(log.TestParticipantID))
		if _, ok := testParticipantLogMap[testParticipantID]; !ok {
			testParticipantLogMap[testParticipantID] = make(map[string]map[string]interface{})
		}

		if _, ok := testParticipantLogMap[testParticipantID][log.SiapQuestionCode]; !ok {
			testParticipantLogMap[testParticipantID][log.SiapQuestionCode] = make(map[string]interface{})
		}
		
		testParticipantLogMap[testParticipantID][log.SiapQuestionCode]["numOfSee"] = log.NumOfSee
		testParticipantLogMap[testParticipantID][log.SiapQuestionCode]["timeOfSee"] = log.TimeOfSee
		
	}

	mapGroupTestAnswers := ConvAnsToMapGroup(dataTestAnswers,dataTestParticipants)
	var headers []interface{}
    var subHeaders []interface{}
	var subSubHeaders []interface{}
    for _, h := range fixedHeadersResponseJawaban {
        headers = append(headers, h)
        subHeaders = append(subHeaders,"")
		subSubHeaders = append(subSubHeaders,"")
    }
	soalTypeConstant := utils.SoalTypeConstant{}


	for _, soal := range dataTestSoals {
		headers = append(headers, soal.SoalCode)
		headers = append(headers, "")
		subHeaders = append(subHeaders,soalTypeConstant.SoalType(soal.SoalTypeID))
		subHeaders = append(subHeaders,mapGroupTestAnswers[soal.SoalCode])
		subSubHeaders = append(subSubHeaders,"Jumlah Diliat (kali)")
		subSubHeaders = append(subSubHeaders,"Total Waktu (detik)")
	}

	t := time.Now().Format("02012006_1504") // dmY_Hi
	fileName :=  ReportFileName(utils.REPORT_THIRD, t, utils.ShortUUID(uuid))
	f := excelize.NewFile()
	sheetName1 := ExportSheetName(t)
	f.SetSheetName("Sheet1", sheetName1)
	streamWriter, _ := f.NewStreamWriter(sheetName1)
	currentRow := 1
	currentPos := len(fixedHeadersResponseJawaban)

	styleJSON := &excelize.Style{
	Alignment: &excelize.Alignment{
		Horizontal: "center",
		Vertical:   "center",
	},
}
	styleCenter, _ := f.NewStyle(styleJSON)

	for col, header := range headers {
		if(col <= (len(fixedHeadersResponseJawaban)-1)){	
			colName, _ := excelize.ColumnNumberToName(col + 1)
			strHeader := fmt.Sprintf("%v", header)
			charWidth := len(strHeader) + 2   
			colWidth := float64(charWidth) * 1.2
			f.SetColWidth(sheetName1, colName, colName, colWidth)
		} else {
			startNumb := utils.NumberToExcelColumn(currentPos+1)
			endNumb := utils.NumberToExcelColumn(currentPos + 2)
			
			startColumn := startNumb + fmt.Sprint(currentRow)
			endColumn := endNumb + fmt.Sprint(currentRow)

			width := float64(18)

			f.MergeCell(sheetName1, startColumn, endColumn)
			f.SetCellStyle(sheetName1, startColumn, endColumn, styleCenter)
			f.SetColWidth(sheetName1, utils.NumberToExcelColumn(currentPos+1), utils.NumberToExcelColumn(currentPos+1), width)
			f.SetColWidth(sheetName1, utils.NumberToExcelColumn(currentPos+2), utils.NumberToExcelColumn(currentPos+2), width)

			currentPos = currentPos+2
		}
	}

	err = streamWriter.SetRow(fmt.Sprintf("A%d", currentRow), headers)

    currentRow += 1

    err = streamWriter.SetRow(fmt.Sprintf("A%d", currentRow),subHeaders)

    currentRow += 1

	err = streamWriter.SetRow(fmt.Sprintf("A%d", currentRow),subSubHeaders)

    currentRow += 1

	data := func()  {
		values := make([]interface{}, 0, len(headers))
		for row, participant := range dataTestParticipants {
			values = values[:0]
			values = append(values,
				participant.NamaSiswa,
				participant.PaketSoal,
				participant.Event,
				participant.Ujian,
			)
			for _, soal := range dataTestSoals {
				participantID := strconv.Itoa(participant.IDSiswa)
				soalCode := soal.SoalCode
				numOfSee := 0
				timeOfSee := 0
				if(testParticipantLogMap[participantID] == nil || testParticipantLogMap[participantID][soalCode] == nil) {
				} else {
					numOfSee = testParticipantLogMap[participantID][soalCode]["numOfSee"].(int)
					timeOfSee = testParticipantLogMap[participantID][soalCode]["timeOfSee"].(int)
				}
				values = append(values, numOfSee)
				values = append(values, timeOfSee)
			}
			err = streamWriter.SetRow(fmt.Sprintf("A%d", row+currentRow), values)
		}
	}

	data()
    _ = streamWriter.Flush()

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}

	_, err = utils.UploadToMinio(fileName, &buf,uuid)
	if err != nil {
		return nil, err
	}

	err = models.ReportOutcomes{}.UpdateReportOutcomes(fileName, uuid, 3)
	if err != nil {
		return nil, err
	}
	
	return &ReportResult{
		FileName: fileName,
		UUID:     uuid,
		Status:   "success",
	}, nil
}

func (ReportController) ResponseJawaban(c *gin.Context) {
	defer utils.EndProccess()	
	uuid := c.Param("uuid")

	result, err := GenerateResponseJawabanReport(uuid)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   result.Status,
		"reqUUID":  result.UUID,
		"fileName": result.FileName,
	})
}

func (ReportController) List(c *gin.Context) {
	var listQuery models.ReportReqsListQuery

	if err := c.BindQuery(&listQuery); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reportReqsM := models.ReportReqsList{}

	data, err := reportReqsM.List(&listQuery)

	if err != nil {
		c.JSON(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	pagination := utils.GetPagination(c, int(data.Total), data.Reports)
	c.JSON(http.StatusOK, pagination)
}

func (ReportController) Outcomes(c *gin.Context) {	
	uuid := c.Param("uuid")
	outcomesM := models.ReportOutcomes{}
	dataOutcomes, err := outcomesM.GetReportOutcomes(uuid)
	if err != nil {
		utils.RespondWithError(c, http.StatusNoContent, "")
		return
	}

	c.JSON(http.StatusOK, dataOutcomes)
}