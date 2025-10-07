package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"
	"ujicoba-go/utils"

	"ujicoba-go/config"

	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TestParticipantReport struct {
	KdProv		string					`json:"KdProv"`
	KdRayon     string      			`json:"KdRayon"`
	KdSek		string      			`json:"KdSek"`
	NPSN		string      			`json:"NPSN"`
	KdJenjang	string      			`json:"KdJenjang"`
	StsSek		string      			`json:"StsSek"`
	IDSekolah	string      			`json:"IDSekolah"`
	Kelas		string      			`json:"Kelas"`
	IDSiswa   	int         			`json:"IDSiswa"`
	NamaSiswa 	string      			`json:"NamaSiswa"`
	PaketSoal	string      			`json:"PaketSoal"`
	Event       string      			`json:"Event"`
	Ujian       string      			`json:"Ujian"`
}

type TestParticipant struct {
	ID                         uint            `json:"id"`
	UUID                       string          `json:"uuid"`
	EventID                    uint            `json:"event_id"`
	TestId                     uint            `json:"test_id"`
	ExamID                     uint            `json:"exam_id"`
	ParticipantType            uint            `json:"participant_type"`
	SiapParticipantUUID        string          `json:"siap_participant_uuid"`
	SiapParticipantName        string          `json:"siap_participant_name"`
	SiapParticipantSchoolId    string          `json:"siap_participant_school_id"`
	SiapParticipantSchoolTitle string          `json:"siap_participant_school_title"`
	Status                     int             `json:"status"`
	SiapTestPackageUuid        string          `json:"siap_test_package_uuid"`
	SiapTestPackageTitle       string          `json:"siap_test_package_title"`
	Answer                     *string         `json:"answer"`
	Score                      *int            `json:"score"`
	TotalRightAnswers          *int            `json:"total_right_answers"`
	TotalWrongAnswers          *int            `json:"total_wrong_answers"`
	CreatedAt                  string          `json:"created_at"`
	UpdatedAt                  string          `json:"updated_at"`
	DeletedAt                  *gorm.DeletedAt `json:"deleted_at"`
	LastOnline                 *string         `json:"last_online"`
	PendataanNisNip            *string         `json:"pendataan_nis_nip"`
	PendataanEmail             *string         `json:"pendataan_email"`
	PendataanUserCode          *string         `json:"pendataan_user_code"`
	PendataanMajor             *string         `json:"pendataan_major"`
	PendataanClass             *string         `json:"pendataan_class"`
	PendataanSubject           *string         `json:"pendataan_subject"`
	PendataanEducationLevel    *string         `json:"pendataan_education_level"`
	PackageID                  int             `json:"package_id" gorm:"package_id"`
	RoomName                   *string         `json:"room_name"`
	Package                    []Package       `json:"package" gorm:"-"`
}

type Package struct {
	ID              uint   `json:"id"`
	MataPelajaranID int    `json:"mata_pelajaran_id"`
	PackageName     string `json:"package_name"`
	JumlahSoal      int    `json:"jumlah_soal"`
}

func (TestParticipant) TableName() string {
	return "test_participants"
}

type TestParticipantListQuery struct {
	Page        string `form:"page,default=1"`
	PageSize    int    `form:"pagesize,default=10"`
	SchoolID    int    `form:"schoolid"`
	SchoolID2   int    `form:"school_id"`
	ExamID      int    `form:"examid"`
	EventID     string `form:"eventuuid"`
	Order       string `form:"order,default=1"`
	Orderby     string `form:"orderby,default=updated_at"`
	HasAnswer   int    `form:"hasAnswer"`
	IsActive    int    `form:"is_active,defult=0"`
	NotInStatus []int
}

type TestParticipantListResult struct {
	TestParticipant []TestParticipant
	Total           int64
}

type TestParticipantExam struct {
	ID                         uint            `json:"id"`
	UUID                       sql.NullString  `json:"uuid"`
	EventID                    uint            `json:"event_id"`
	TestId                     uint            `json:"test_id"`
	ExamID                     uint            `json:"exam_id"`
	ParticipantType            uint            `json:"participant_type"`
	SiapParticipantUUID        sql.NullString  `json:"siap_participant_uuid"`
	SiapParticipantName        sql.NullString  `json:"siap_participant_name"`
	SiapParticipantSchoolId    sql.NullString  `json:"siap_participant_school_id"`
	SiapParticipantSchoolTitle sql.NullString  `json:"siap_participant_school_title"`
	Status                     sql.NullString  `json:"status"`
	SiapTestPackageUuid        sql.NullString  `json:"siap_test_package_uuid"`
	SiapTestPackageTitle       sql.NullString  `json:"siap_test_package_title"`
	Answer                     sql.NullString  `json:"answer"`
	Score                      sql.NullString  `json:"score"`
	TotalRightAnswers          sql.NullString  `json:"total_right_answers"`
	TotalWrongAnswers          sql.NullString  `json:"total_wrong_answers"`
	CreatedAt                  sql.NullString  `json:"created_at"`
	UpdatedAt                  sql.NullString  `json:"updated_at"`
	DeletedAt                  *gorm.DeletedAt `json:"deleted_at"`
	LastOnline                 sql.NullString  `json:"last_online"`
	PendataanNisNip            sql.NullString  `json:"pendataan_nis_nip"`
	PendataanEmail             string          `json:"pendataan_email"`
	PendataanUserCode          sql.NullString  `json:"pendataan_user_code"`
	PendataanMajor             sql.NullString  `json:"pendataan_major"`
	PendataanClass             sql.NullString  `json:"pendataan_class"`
	PendataanSubject           sql.NullString  `json:"pendataan_subject"`
	PendataanEducationLevel    sql.NullString  `json:"pendataan_education_level"`
	PackageID                  int             `json:"package_id" gorm:"package_id"`
}

func (TestParticipantExam) TableName() string {
	return "test_participants"
}

type StudentRoom struct {
	ID        uint            `json:"id"`
	Order     int             `json:"order"`
	SekolahID int             `json:"sekolah_id"`
	RuanganID int             `json:"ruangan_id"`
	StudentID int             `json:"student_id"`
	CreatedAt sql.NullString  `json:"created_at"`
	UpdatedAt sql.NullString  `json:"updated_at"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at"`
}

func (StudentRoom) TableName() string {
	return "ruangan_student"
}

type roomDetail struct {
	RoomName *string `gorm:"column:nama_ruang" json:"nama_ruang"`
}

func (TestParticipant) ListStudent(query *TestParticipantListQuery, examID string) (TestParticipantListResult, error) {
	db := GetDBApd()
	var participant []TestParticipant
	page, err := strconv.Atoi(query.Page)

	if err != nil {
		page = 1
	}
	offset := (page - 1) * query.PageSize

	orderValue := "ASC"
	if query.Order == "0" {
		orderValue = "ASC"
	} else if query.Order == "1" {
		orderValue = "DESC"
	}

	orderStr := fmt.Sprintf("%s %s", query.Orderby, orderValue)

	queryModel := db.Model(&TestParticipant{}).Offset(offset).Limit(query.PageSize).Order(orderStr).Where("exam_id = ?", examID)

	if len(query.NotInStatus) > 0 {
		queryModel = queryModel.Where("status NOT IN ?", query.NotInStatus)
	}

	if schoolId := query.SchoolID2; schoolId != 0 {
		queryModel = queryModel.Where("siap_participant_school_id = ?", schoolId)
	}

	if err := queryModel.Find(&participant).Error; err != nil {
		return TestParticipantListResult{}, err
	}

	var total int64
	queryCount := db.Model(&TestParticipant{}).Where("exam_id = ?", examID)

	if schoolId := query.SchoolID2; schoolId != 0 {
		queryCount = queryCount.Where("siap_participant_school_id = ?", schoolId)
	}

	if len(query.NotInStatus) > 0 {
		queryCount = queryCount.Where("status NOT IN ?", query.NotInStatus)
	}

	if err := queryCount.Count(&total).Error; err != nil {
		return TestParticipantListResult{}, err
	}

	if config.EnvVariable("SOAL_ONLINE") == "true" {
		for i := range participant {
			dbSiap := GetDBSiap()
			var packageDetail Package

			if err := dbSiap.Select("id, mata_pelajaran_id, package_name, (select count(id) from question_package_detail where package_id = question_package.id and deletedAt is null) as jumlah_soal").
				Model(&QuestionPackage{}).
				First(&packageDetail, participant[i].PackageID).
				Error; err != nil {
			}
			participant[i].Package = append(participant[i].Package, packageDetail)

			dbPdt := GetDBPdt()

			var room roomDetail
			error := dbPdt.Select("(select name from ruangan where id = ruangan_student.ruangan_id) as nama_ruang").Where("sekolah_id = ? and student_id = ?", participant[i].SiapParticipantSchoolId, participant[i].SiapParticipantUUID).Model(&StudentRoom{}).First(&room).Error
			if error == nil {
				participant[i].RoomName = room.RoomName
			}

			participant[i].CreatedAt = convertDatetimeEvent(participant[i].CreatedAt)
			participant[i].UpdatedAt = convertDatetimeEvent(participant[i].UpdatedAt)
		}
	}

	result := TestParticipantListResult{
		TestParticipant: participant,
		Total:           total,
	}

	return result, nil
}

type TestParticipantListBasic struct {
	ID                         uint            `json:"id"`
	UUID                       string          `json:"uuid"`
	EventID                    uint            `json:"event_id"`
	TestId                     uint            `json:"test_id"`
	ExamID                     uint            `json:"exam_id"`
	ParticipantType            uint            `json:"participant_type"`
	SiapParticipantUUID        string          `json:"siap_participant_uuid"`
	SiapParticipantName        string          `json:"siap_participant_name"`
	SiapParticipantSchoolId    string          `json:"siap_participant_school_id"`
	SiapParticipantSchoolTitle string          `json:"siap_participant_school_title"`
	Status                     int             `json:"status"`
	SiapTestPackageUuid        string          `json:"siap_test_package_uuid"`
	SiapTestPackageTitle       string          `json:"siap_test_package_title"`
	Answer                     *string         `json:"answer"`
	Score                      *int            `json:"score"`
	TotalRightAnswers          *int            `json:"total_right_answers"`
	TotalWrongAnswers          *int            `json:"total_wrong_answers"`
	CreatedAt                  *string         `json:"created_at"`
	UpdatedAt                  *string         `json:"updated_at"`
	DeletedAt                  *gorm.DeletedAt `json:"deleted_at"`
	LastOnline                 *string         `json:"last_online"`
	PendataanNisNip            *string         `json:"pendataan_nis_nip"`
	PendataanEmail             *string         `json:"pendataan_email"`
	PendataanUserCode          *string         `json:"pendataan_user_code"`
	PendataanMajor             *string         `json:"pendataan_major"`
	PendataanClass             *string         `json:"pendataan_class"`
	PendataanSubject           *string         `json:"pendataan_subject"`
	PendataanEducationLevel    *string         `json:"pendataan_education_level"`
	PackageID                  int             `json:"package_id" gorm:"package_id"`
	Event                      *Event          `json:"event" gorm:"-"`
	Exam                       *Exams          `json:"exam" gorm:"-"`
	Log                        *[]Log          `json:"log" gorm:"-"`
}

type TestParticipantList struct {
	TestParticipant []TestParticipantListBasic
	Total           int64
}

type Exams struct {
	ID                uint            `json:"id"`
	EventID           uint            `json:"event_id"`
	JenisPendidikanID int             `json:"jenis_pendidikan_id"`
	ClassID           int             `json:"class_id"`
	JurusanID         int             `json:"jurusan_id"`
	MataPelajaranID   int             `json:"mata_pelajaran_id"`
	Name              string          `json:"name"`
	CreatedAt         string          `json:"created_at"`
	UpdatedAt         string          `json:"updated_at"`
	DeletedAt         *gorm.DeletedAt `json:"deleted_at"`
}

type Log struct {
	ID                 uint    `json:"id"`
	TestParticipantID  uint    `json:"test_participant_id"`
	TestAnswerID       *int64  `json:"test_answer_id"`
	Action             int16   `json:"action"`
	CreatedAt          string  `json:"created_at"`
	UpdatedAt          string  `json:"updated_at"`
	CreatedTime        string  `json:"created_time"`
	UpdatedTime        string  `json:"updated_time"`
	IpClient           string  `json:"ip_client"`
	IpServer           string  `json:"ip_server"`
	SiapQuestionCode   *string `json:"siap_question_code"`
	NextQuestionCode   *string `json:"next_question_code"`
	SiapQuestionNumber *int16  `json:"siap_question_number"`
	NextQuestionNumber *int16  `json:"next_question_number"`
	IsInDoubt          *int16  `json:"is_in_doubt"`
	RemainingTime      *int64  `json:"remaining_time"`
	Answer             *string `json:"answer"`
	IntervalNow        string  `json:"interval_now,omitempty"`
	IsOnline           bool    `json:"is_online,omitempty"`
}

func convertDatetimeLog(input string) string {
	t, _ := time.Parse("2006-01-02T15:04:05Z", input)
	tz, _ := time.LoadLocation("UTC")
	zoneTime := t.In(tz)
	return zoneTime.Format("2006-01-02 15:04:05")
}

func convertDatetimeLogSlash(input string) string {
	t, _ := time.Parse("2006-01-02T15:04:05Z", input)
	tz, _ := time.LoadLocation("UTC")
	zoneTime := t.In(tz)
	return zoneTime.Format("02/01/2006 15:04:05")
}

func (TestParticipant) List(query *TestParticipantListQuery, examID string) (TestParticipantList, error) {
	db := GetDBApd()
	var participant []TestParticipantListBasic
	page, err := strconv.Atoi(query.Page)
	pageSize := query.PageSize

	if err != nil {
		page = 1
	}

	offset := (page - 1) * pageSize
	orderValue := "ASC"
	if query.Order == "0" {
		orderValue = "ASC"
	} else if query.Order == "1" {
		orderValue = "DESC"
	}

	orderStr := fmt.Sprintf("%s %s", query.Orderby, orderValue)

	queryModel := db.Select("test_participants.id,test_participants.uuid,event_id,test_participants.test_id,test_participants.exam_id,participant_type,siap_participant_uuid,siap_participant_name,siap_participant_school_id,siap_participant_school_title,test_participants.status,test_participants.siap_test_package_uuid,test_participants.siap_test_package_title,test_participants.answer,test_participants.score,test_participants.total_right_answers,test_participants.total_wrong_answers,test_participants.created_at,test_participants.updated_at,test_participants.deleted_at,last_online,pendataan_nis_nip,pendataan_email,pendataan_user_code,pendataan_major,pendataan_class,pendataan_subject,pendataan_education_level,package_id").
		Model(&TestParticipant{}).
		Offset(offset).
		Limit(pageSize).
		Order(orderStr).
		Group("test_participants.id")
	var total int64
	queryCount := db.Model(&TestParticipant{}).Group("test_participants.id")

	if schoolID := query.SchoolID; schoolID != 0 {
		queryModel = queryModel.Where("siap_participant_school_id = ?", schoolID)
		queryCount = queryCount.Where("siap_participant_school_id = ?", schoolID)
	}

	if examID := query.ExamID; examID != 0 {
		queryModel = queryModel.Where("test_participants.exam_id = ?", examID)
		queryCount = queryCount.Where("test_participants.exam_id = ?", examID)
	}

	if isActive := query.IsActive; isActive != 0 {
		queryModel = queryModel.Where("test_participants.status", isActive)
		queryCount = queryCount.Where("test_participants.status", isActive)
	}

	var eventDetail Event

	if eventID := query.EventID; eventID != "" {
		err = db.Select("id").Model(&Event{}).Where("uuid = ?", eventID).First(&eventDetail).Error
		if err != nil {

		}
		queryModel = queryModel.Where("event_id = ?", eventDetail.ID)
		queryCount = queryCount.Where("event_id = ?", eventDetail.ID)
	}

	if hasAnswer := query.HasAnswer; hasAnswer == 1 {
		queryModel = queryModel.Joins("JOIN test_answers ON test_answers.test_participant_id = test_participants.id")
		queryCount = queryCount.Joins("JOIN test_answers ON test_answers.test_participant_id = test_participants.id")
	}

	if err := queryModel.Find(&participant).Error; err != nil {
		return TestParticipantList{}, err
	}

	// Collect unique event IDs and exam IDs
	eventIDs := make([]int, 0, len(participant))
	examIDs := make([]int, 0, len(participant))
	for _, data := range participant {
		eventIDs = append(eventIDs, int(data.EventID))
		examIDs = append(examIDs, int(data.ExamID))
	}
	eventIDs = uniqueIntSlice(eventIDs)
	examIDs = uniqueIntSlice(examIDs)

	// Fetch events
	events := make(map[int]Event)
	if len(eventIDs) > 0 {
		var eventsList []Event
		errEvent := db.Select("id, uuid, name, start_date, duration, status, created_at, updated_at, is_showing_result, timezone, deleted_at").Model(&Event{}).Where("id IN (?)", eventIDs).Find(&eventsList).Error
		if errEvent == nil {
			for _, event := range eventsList {
				event.CreatedAt = convertDatetimeLog(event.CreatedAt)
				event.UpdatedAt = convertDatetimeLog(event.UpdatedAt)
				event.StartDate = convertDatetimeLog(event.StartDate)
				events[int(event.ID)] = event
			}
		}
	}

	// Fetch exams
	exams := make(map[int]Exam)
	if len(examIDs) > 0 {
		var examsList []Exam
		errExam := db.Select("id,event_id,jenis_pendidikan_id,class_id,jurusan_id,mata_pelajaran_id,name,created_at,updated_at,deleted_at").Model(&Exam{}).Where("id IN (?)", examIDs).Find(&examsList).Error
		if errExam == nil {
			for _, exam := range examsList {
				exam.CreatedAt = convertDatetimeLog(exam.CreatedAt)
				exam.UpdatedAt = convertDatetimeLog(exam.UpdatedAt)
				exams[int(exam.ID)] = exam
			}
		}
	}

	// Fetch logs
	for i, data := range participant {
		if participant[i].CreatedAt != nil {
			createdAt := *participant[i].CreatedAt
			createdAt = convertDatetimeLog(createdAt)
			participant[i].CreatedAt = &createdAt
		}

		if participant[i].UpdatedAt != nil {
			UpdatedAt := *participant[i].UpdatedAt
			UpdatedAt = convertDatetimeLog(UpdatedAt)
			participant[i].UpdatedAt = &UpdatedAt
		}

		if participant[i].LastOnline != nil {
			online := *participant[i].LastOnline
			lastonline := convertDatetimeLog(online)
			participant[i].LastOnline = &lastonline
		}

		event, eventExists := events[int(data.EventID)]
		if eventExists {
			eventCopy := event
			participant[i].Event = &eventCopy
		}

		exam, examExists := exams[int(data.ExamID)]
		if examExists {
			examCopy := exam
			participant[i].Exam = &Exams{
				ID:                examCopy.ID,
				EventID:           examCopy.EventID,
				JenisPendidikanID: examCopy.JenisPendidikanID,
				ClassID:           examCopy.ClassID,
				JurusanID:         examCopy.JurusanID,
				MataPelajaranID:   examCopy.MataPelajaranID,
				Name:              examCopy.Name,
				CreatedAt:         examCopy.CreatedAt,
				UpdatedAt:         examCopy.UpdatedAt,
				DeletedAt:         examCopy.DeletedAt,
			}
		}

		var logdata []Log
		errLog := db.Select("id,test_participant_id,test_answer_id,action,created_at,created_at as created_time,updated_at,updated_at as updated_time,ip_client,ip_server,siap_question_code,next_question_code,siap_question_number,next_question_number,is_in_doubt,remaining_time,answer, '' as interval_now, false as is_online").
			Model(&TestParticipantLogView{}).
			Where("test_participant_id = ?", data.ID).
			Find(&logdata).Error
		if errLog != nil || len(logdata) != 0 {
			participant[i].Log = nil
		}

		isOnline := getIsOnline(int(participant[i].ID))

		for y, _ := range logdata {
			date := carbon.Parse(logdata[y].CreatedAt)
			now := carbon.Now()
			diff := now.SetLocale("id").DiffForHumans(date)

			logdata[y].IsOnline = isOnline
			logdata[y].CreatedAt = convertDatetimeLog(logdata[y].CreatedAt)
			logdata[y].IntervalNow = diff
			logdata[y].CreatedTime = convertDatetimeLogSlash(logdata[y].CreatedTime)
			logdata[y].UpdatedAt = convertDatetimeLog(logdata[y].UpdatedAt)
			logdata[y].UpdatedTime = convertDatetimeLogSlash(logdata[y].UpdatedTime)
		}
		participant[i].Log = &logdata
	}

	if err := queryCount.Count(&total).Error; err != nil {
		return TestParticipantList{}, err
	}

	result := TestParticipantList{
		TestParticipant: participant,
		Total:           total,
	}

	return result, nil
}

func uniqueIntSlice(ints []int) []int {
	uniqueMap := make(map[int]bool)
	uniqueSlice := make([]int, 0, len(ints))
	for _, val := range ints {
		if !uniqueMap[val] {
			uniqueMap[val] = true
			uniqueSlice = append(uniqueSlice, val)
		}
	}
	return uniqueSlice
}

func getIsOnline(testParticipantID int) bool {
	db := GetDBApd()
	lastLogin := TestParticipantLogs{}
	db.Select("created_at").Where("test_participant_id = ? AND action = ?", testParticipantID, utils.ACTION_LOG_LOGIN).
		Where("DATE(created_at) = CURDATE()").
		Order("created_at DESC").
		First(&lastLogin)

	lastLogout := TestParticipantLogs{}
	db.Select("created_at").Where("test_participant_id = ? AND action = ?", testParticipantID, utils.ACTION_LOG_LOGOUT).
		Where("DATE(created_at) = CURDATE()").
		Order("created_at DESC").
		First(&lastLogout)

	if lastLogout.ID == 0 {
		return true
	}

	loginTime, _ := time.Parse("2006-01-02T15:04:05Z", lastLogin.CreatedAt)
	logoutTime, _ := time.Parse("2006-01-02T15:04:05Z", lastLogout.CreatedAt)

	duration := logoutTime.Sub(loginTime)
	minutesDiff := int(duration.Minutes())

	if minutesDiff < 1 {
		return true
	}

	return false
}

func (TestParticipant) GetParticipantCount(participant int) (int64, error) {
	db := GetDBApd()

	query := db.Model(&TestParticipant{}).Where("id = ?", participant)

	var count int64
	err := query.Count(&count).Error
	if err != nil {
		return count, err
	}

	return count, err
}

type TestParticipantUpdate struct {
	ID                   uint   `json:"id"`
	Status               string `json:"status"`
	SiapTestPackageUuid  string `json:"siap_test_package_uuid"`
	SiapTestPackageTitle string `json:"siap_test_package_title"`
	UpdatedAt            string `json:"updated_at"`
	PackageID            int    `json:"package_id" gorm:"package_id"`
}

func (t *TestParticipant) Update(id string) error {
	db := GetDBApd()
	var testParticipant TestParticipant

	if err := db.Where("id = ?", id).First(&testParticipant).Error; err != nil {
		return err
	}

	status := strconv.Itoa(t.Status)

	testParticipantsUpdate := TestParticipantUpdate{
		SiapTestPackageUuid:  t.SiapTestPackageUuid,
		SiapTestPackageTitle: t.SiapTestPackageTitle,
		PackageID:            t.PackageID,
		Status:               status,
		UpdatedAt:            time.Now().Format("2006-01-02 15:04:05"),
	}

	if err := db.Model(TestParticipant{}).Where("id = ?", id).Updates(testParticipantsUpdate).Error; err != nil {
		return err
	}

	testParticipant.CreatedAt = convertDatetimeEvent(testParticipant.CreatedAt)

	return nil
}

func (TestParticipant) BulkUpsertParticipantExam(data *[]TestParticipantExam, isUpdate bool, txApd *gorm.DB) error {
	for i := range *data {
		(*data)[i].CreatedAt = sql.NullString{String: time.Now().Format("2006-01-02 15:04:05"), Valid: true}
		(*data)[i].Status = sql.NullString{
			Valid:  true,
			String: "-1",
		}
		if (*data)[i].PackageID != 0 {
			(*data)[i].Status.String = "0"
		}
	}
	err := txApd.Clauses(clause.OnConflict{
		DoNothing: !isUpdate,
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"package_id",
			"siap_test_package_uuid",
			"siap_test_package_title",
			"updated_at",
			"created_at",
			"status",
		}),
	}).
		Create(&data).Error

	return err
}

func (tp *TestParticipant) GetParticipantByExam(participants *[]TestParticipant, txApd *gorm.DB) error {
	if err := txApd.Where("exam_id", tp.ExamID).Find(&participants).Error; err != nil {
		return err
	}
	return nil
}

type TestParticipantByExamID struct {
	ID                         uint   `json:"id"`
	SiapParticipantSchoolTitle string `json:"siap_participant_school_title"`
	SiapParticipantName        string `json:"siap_participant_name"`
	SiapParticipantUUID        string `json:"siap_participant_uuid"`
	PackageID                  uint   `json:"package_id"`
}

type PdtPartcipant struct {
	ID string `gorm:"column:id" json:"id"`
}

func (TestParticipant) GetByExamID(id string) ([]TestParticipantByExamID, error) {
	db := GetDBApd()

	query := db.Select("siap_participant_uuid").Model(&TestParticipant{}).Where("exam_id = ?", id)

	var testParticipant []TestParticipantByExamID
	err := query.Find(&testParticipant).Error
	if err != nil {
		return testParticipant, err
	}

	var idTestParticipant []string
	for _, data := range testParticipant {
		idTestParticipant = append(idTestParticipant, data.SiapParticipantUUID)
	}

	dbPdt := GetDBPdt()
	var pdtParticipant []PdtPartcipant
	errPdt := dbPdt.Select("id").
		Where("id in (?)", idTestParticipant).
		Where("status = ?", utils.EVENT_SISWA_ACTIVE).
		Model(&Student{}).
		Find(&pdtParticipant).Error

	if errPdt != nil {
		return testParticipant, err
	}

	var idTestParticipantActive []string
	for _, data := range pdtParticipant {
		idTestParticipantActive = append(idTestParticipantActive, data.ID)
	}

	queryActiveParticipant := db.Model(&TestParticipant{}).Where("exam_id = ?", id).Where("siap_participant_uuid in (?)", idTestParticipantActive)
	var activeParticipant []TestParticipantByExamID
	errActiveParticipant := queryActiveParticipant.Find(&activeParticipant).Error
	if errActiveParticipant != nil {
		return activeParticipant, errActiveParticipant
	}

	return activeParticipant, errActiveParticipant
}

func (TestParticipant) GetActiveByUUID(uuid string) (TestParticipant, error) {
	db := GetDBApd()
	var result TestParticipant
	query := db.Where("uuid = ?", uuid).Where("status != ?", utils.PARTICIPANT_STATUS_NOTACTIVE)

	if err := query.First(&result).Error; err != nil {
		return result, err
	}

	if result.CreatedAt != "" {
		result.CreatedAt = ConvertDatetime(result.CreatedAt)
	}

	if result.UpdatedAt != "" {
		result.UpdatedAt = ConvertDatetime(result.UpdatedAt)
	}

	return result, nil
}

func (participant *TestParticipant) GetByExam() (TestParticipant, error) {
	db := GetDBApd()
	var result TestParticipant
	query := db.
		Where("uuid = ?", participant.UUID).
		Where("exam_id", participant.ExamID)

	if err := query.First(&result).Error; err != nil {
		return result, err
	}

	if result.CreatedAt != "" {
		result.CreatedAt = ConvertDatetime(result.CreatedAt)
	}

	if result.UpdatedAt != "" {
		result.UpdatedAt = ConvertDatetime(result.UpdatedAt)
	}

	return result, nil
}

func (participant *TestParticipant) GetByExams() ([]TestParticipant, error) {
	db := GetDBApd()
	var result []TestParticipant
	query := db.
		Where("siap_participant_school_id", participant.SiapParticipantSchoolId)

	if participant.ExamID != 0 {
		query = query.Where("exam_id", participant.ExamID)
	}

	if err := query.Find(&result).Error; err != nil {
		return result, err
	}

	for i := range result {
		if result[i].CreatedAt != "" {
			result[i].CreatedAt = ConvertDatetime(result[i].CreatedAt)
		}

		if result[i].UpdatedAt != "" {
			result[i].UpdatedAt = ConvertDatetime(result[i].UpdatedAt)
		}
	}

	return result, nil
}

func (TestParticipant) GetByUUID(uuid string) (TestParticipantExam, error) {
	db := GetDBApd()

	var result TestParticipantExam
	query := db.Where("uuid = ?", uuid)

	if err := query.First(&result).Error; err != nil {
		return result, err
	}

	return result, nil
}

// update with transaction
func (t *TestParticipant) UpdateWithTx(id string, txApd *gorm.DB) error {
	var testParticipant TestParticipant

	if err := txApd.Where("id = ?", id).First(&testParticipant).Error; err != nil {
		return err
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	t.UpdatedAt = now

	if t.CreatedAt == "" {
		t.CreatedAt = now
	}

	if err := txApd.Model(&testParticipant).Updates(t).Error; err != nil {
		return err
	}

	*t = testParticipant

	return nil
}

func (TestParticipant) ToggleActive(participant *TestParticipantExam) error {
	db := GetDBApd()

	now := time.Now().Format("2006-01-02 15:04:05")

	var createdAt time.Time
	var err error
	if participant.CreatedAt.String != "" {
		createdAt, err = time.Parse("2006-01-02T15:04:05Z", participant.CreatedAt.String)
		if err != nil {
			return err
		}
	} else {
		createdAt = time.Now()
	}

	participant.UpdatedAt = sql.NullString{
		Valid:  true,
		String: now,
	}
	participant.CreatedAt = sql.NullString{
		Valid:  true,
		String: createdAt.Format("2006-01-02 15:04:05"),
	}

	participant.LastOnline = sql.NullString{
		Valid: false,
	}

	if participant.Status.String == "-1" {
		participant.Status.String = "0"
	} else if participant.Status.String == "0" {
		participant.Status.String = "-1"
	} else {
		return errors.New("error_test_participant_status")
	}

	err = db.Clauses(clause.OnConflict{
		DoNothing: false,
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"status", "updated_at", "created_at"}),
	}).
		Create(&participant).Error

	return err
}

func (participant *TestParticipant) GetMyPerticipantWithExamID(examId string) (TestParticipantExam, error) {
	db := GetDBApd()

	var result TestParticipantExam
	query := db.Where("exam_id = ?", examId).Where("siap_participant_uuid = ?", participant.SiapParticipantUUID)

	if err := query.First(&result).Error; err != nil {
		return result, err
	}

	if result.CreatedAt.Valid {
		result.CreatedAt.String = ConvertDatetime(result.CreatedAt.String)
	}

	if result.UpdatedAt.Valid {
		result.UpdatedAt.String = ConvertDatetime(result.UpdatedAt.String)
	}

	if result.LastOnline.Valid {
		result.LastOnline.String = ConvertDatetime(result.LastOnline.String)
	}

	return result, nil
}

func (TestParticipant) ResetLoginWithTx(participant *TestParticipantExam, txApd *gorm.DB) error {
	query := txApd.Where("exam_id = ?", participant.ExamID).
		Where("siap_participant_uuid = ?", participant.SiapParticipantUUID.String)

	participant.LastOnline = sql.NullString{
		Valid:  false,
		String: "",
	}

	if err := query.Updates(&participant).Error; err != nil {
		return err
	}

	return nil
}

func (t *TestParticipant) DeleteBySchoolWithTx(schoolId int, txApd *gorm.DB) error {
	if err := txApd.
		Where("siap_participant_school_id = ?", schoolId).
		Where("exam_id = ?", t.ExamID).
		Delete(&TestParticipantExam{}).Error; err != nil {
		return err
	}
	return nil
}

func (TestParticipant) GetByIDs(ids ...string) ([]TestParticipant, error) {
	db := GetDBApd()

	var result []TestParticipant
	query := db.Where("id in (?)", ids)

	if err := query.Find(&result).Error; err != nil {
		return result, err
	}

	return result, nil
}

func (TestParticipant) GetByIDsInt(ids []int) ([]TestParticipant, error) {
	db := GetDBApd()

	var result []TestParticipant
	query := db.Where("id in (?)", ids)

	if err := query.Find(&result).Error; err != nil {
		return result, err
	}

	return result, nil
}

func (t *TestParticipant) AssignParticipantPackages(data []TestParticipant, txApd *gorm.DB) error {
	for _, value := range data {
		var testParticipant TestParticipant
		if err := txApd.Where("id = ?", value.ID).First(&testParticipant).Error; err != nil {
			return err
		}

		now := time.Now().Format("2006-01-02 15:04:05")
		t.UpdatedAt = now
		t.PackageID = value.PackageID
		t.SiapTestPackageUuid = value.SiapTestPackageUuid
		t.SiapTestPackageTitle = value.SiapTestPackageTitle

		if err := txApd.Model(&testParticipant).Updates(t).Error; err != nil {
			return err
		}
	}

	return nil
}

func (t *TestParticipant) GetParticipantSchoolExam(schoolId int) ([]TestParticipantExam, error) {
	db := GetDBApd()

	var results []TestParticipantExam

	if err := db.
		Where("siap_participant_school_id = ?", schoolId).
		Where("exam_id = ?", t.ExamID).Find(&results).Error; err != nil {
		return results, err
	}

	return results, nil
}

func (t *TestParticipant) GetByLastOnlineByToday() (TestParticipantExam, error) {
	db := GetDBApd()
	now := time.Now().Format("2006-01-02")

	var result TestParticipantExam
	if err := db.
		Where("siap_participant_uuid = ?", t.SiapParticipantUUID).
		Where("last_online >= ? AND last_online < ?", now, now+" 24:00:00").
		First(&result).Error; err != nil {
		return result, err
	}

	return result, nil
}

func (t *TestParticipant) SetNullLastOnlineWithTx(txApd *gorm.DB) error {
	query := txApd.Model(TestParticipant{}).
		Where("id = ?", t.ID).
		Updates(map[string]interface{}{
			"last_online": gorm.Expr("NULL"),
		})

	if err := query.Error; err != nil {
		return err
	}

	return nil
}

func (t *TestParticipant) UpdateAllByTx(id []int, txApd *gorm.DB) error {
	var testParticipant TestParticipant

	now := time.Now().Format("2006-01-02 15:04:05")
	t.UpdatedAt = now

	if t.CreatedAt == "" {
		t.CreatedAt = now
	}

	if err := txApd.Model(&testParticipant).Where("id in (?)", id).Updates(t).Error; err != nil {
		return err
	}

	*t = testParticipant

	return nil
}

func (TestParticipant) GetDataExport(UUIDs []string, eventIDs []string, examIDs []string) ([]TestParticipant, error) {
	db := GetDBApd()
	query := db.
		Where("siap_test_package_uuid IN ?", UUIDs)

	if len(eventIDs) > 0 {
		query = query.Where("event_id IN ?", eventIDs)
	}
	if len(examIDs) > 0 {
		query = query.Where("exam_id IN ?", examIDs)
	}

	var results []TestParticipant
	if err := query.Find(&results).Error; err != nil {
		return results, err
	}

	return results, nil
}

type byEventReponse struct {
	Total     int64    `gorm:"column:id" json:"id"`
	PackageID []string `gorm:"column:package_id" json:"package_id"`
}

func (participant *TestParticipant) GetByEventCount(schoolId int, eventId int) (*byEventReponse, error) {
	db := GetDBApd()
	var result byEventReponse
	var total int64
	var packages []TestParticipant
	var packagesId []string
	query := db.Where("event_id", eventId).Where("siap_participant_school_id", schoolId)

	if err := query.Model(&TestParticipant{}).Count(&total).Error; err != nil {
		return nil, err
	}

	queryDetail := db.Model(&TestParticipant{}).Where("event_id", eventId).Where("siap_participant_school_id", schoolId)

	if errDetail := queryDetail.Find(&packages).Error; errDetail != nil {
		return nil, errDetail
	}

	for _, value := range packages {
		packagesId = append(packagesId, value.SiapTestPackageUuid)
	}

	result = byEventReponse{
		Total:     total,
		PackageID: packagesId,
	}

	return &result, nil
}

func (TestParticipant) GetByExamSchool(examIDs []string, schoolID uint) ([]TestParticipant, error) {
	db := GetDBApd()
	var results []TestParticipant
	query := db.Where("exam_id in ?", examIDs).
		Where("siap_participant_school_id = ?", schoolID).
		Where("test_id IS NOT NULL")
	if err := query.Find(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

type TestParticipantBySchool struct {
	ID                uint    `json:"id"`
	Status            int     `json:"status"`
	Answer            *string `json:"answer"`
	Score             *int    `json:"score"`
	TotalRightAnswers *int    `json:"total_right_answers"`
	TotalWrongAnswers *int    `json:"total_wrong_answers"`
	LastOnline        *string `json:"last_online"`
}

func (TestParticipant) GetBySchoolId(schoolId ...int) ([]TestParticipantBySchool, error) {
	db := GetDBApd()

	var result []TestParticipantBySchool
	query := db.Model(TestParticipant{}).Where("siap_participant_school_id in (?)", schoolId)

	if err := query.Find(&result).Error; err != nil {
		return result, err
	}

	return result, nil
}

func (tp TestParticipant) UpsertRawWithTx(txApd *gorm.DB, data []map[string]interface{}) error {
	for _, row := range data {
		query := txApd.Table(tp.TableName()).Clauses(clause.OnConflict{
			DoNothing: false,
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.Assignments(row),
		}).Create(row)

		if err := query.Error; err != nil {
			return err
		}
	}

	return nil
}

func (TestParticipant) ParseJSON(data TestParticipantExam) map[string]interface{} {
	result := map[string]interface{}{}
	result["id"] = data.ID
	result["uuid"] = data.UUID.String
	result["event_id"] = data.EventID
	result["test_id"] = data.TestId
	result["exam_id"] = data.ExamID
	result["participant_type"] = data.ParticipantType
	result["siap_participant_uuid"] = data.SiapParticipantUUID.String
	result["siap_participant_name"] = data.SiapParticipantName.String
	result["siap_participant_school_id"] = data.SiapParticipantSchoolId.String
	result["siap_participant_school_title"] = data.SiapParticipantSchoolTitle.String
	result["status"] = data.Status
	result["siap_test_package_uuid"] = data.SiapTestPackageUuid.String
	result["siap_test_package_title"] = data.SiapTestPackageTitle.String
	result["answer"] = data.Answer.String
	result["score"] = data.Score.String
	result["total_right_answers"] = data.TotalRightAnswers.String
	result["total_wrong_answers"] = data.TotalWrongAnswers.String
	result["pendataan_nis_nip"] = data.PendataanNisNip.String
	result["pendataan_email"] = data.PendataanEmail
	result["pendataan_user_code"] = data.PendataanUserCode.String
	result["pendataan_major"] = data.PendataanMajor.String
	result["pendataan_class"] = data.PendataanClass.String
	result["pendataan_subject"] = data.PendataanSubject.String
	result["pendataan_education_level"] = data.PendataanEducationLevel.String
	result["package_id"] = data.PackageID
	result["deleted_at"] = data.DeletedAt

	if data.CreatedAt.Valid {
		result["created_at"] = ConvertDatetime(data.CreatedAt.String)
	} else {
		result["created_at"] = nil
	}

	if data.UpdatedAt.Valid {
		result["updated_at"] = ConvertDatetime(data.UpdatedAt.String)
	} else {
		result["updated_at"] = nil
	}

	if data.LastOnline.Valid {
		result["last_online"] = ConvertDatetime(data.LastOnline.String)
	} else {
		result["last_online"] = nil
	}

	return result
}

func (TestParticipant) UpdateStatuses(ids []int, status string) error {
	db := GetDBApd()
	query := db.Model(TestParticipantExam{}).Where("id IN ?", ids).Update("status", status)
	if err := query.Error; err != nil {
		return err
	}
	return nil
}

func (tp *TestParticipant) UpsertTableWithTx(data []map[string]interface{}, txApd *gorm.DB) error {
	if len(data) > 0 {
		for _, row := range data {
			query := txApd.Table(tp.TableName()).Clauses(clause.OnConflict{
				DoNothing: false,
				Columns:   []clause.Column{{Name: "id"}},
				DoUpdates: clause.Assignments(row),
			}).Create(row)

			if err := query.Error; err != nil {
				return err
			}
		}
	}

	return nil
}

func (tp *TestParticipant) UpdateAllNonActiveByEvent(eventID uint) error {
	examM := Exam{}
	exams, err := examM.GetByEventID(int64(eventID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}
	examIDs := func() []uint {
		var ids []uint
		for _, item := range exams {
			ids = append(ids, item.ID)
		}
		return ids
	}()

	err = tp.UpdateAllNonActiveByExams(examIDs)
	if err != nil {
		return err
	}

	return nil
}

func (TestParticipant) UpdateAllNonActiveByExams(examIDs []uint) error {
	db := GetDBApd()

	query := db.Where("exam_id IN ?", examIDs).Find(&TestParticipant{})
	if err := query.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}

	query = db.Exec("UPDATE test_participants SET status = ? WHERE exam_id IN ?", utils.PARTICIPANT_STATUS_NOTACTIVE, examIDs)
	if err := query.Error; err != nil {
		return err
	}

	return nil
}

func (TestParticipantReport) GetParticipantAnswer(packageIDs []int) ([]TestParticipantReport, error) {
	db := GetDBApd()
	var results []TestParticipantReport
	schema := config.EnvVariable("PENDATAAN_DB_DATABASE")
	query := fmt.Sprintf(`
		select 
			%[1]s.sekolah.kode_provinsi AS KdProv,
			%[1]s.sekolah.kode_rayon AS KdRayon,
			%[1]s.sekolah.kode_sekolah AS KdSek,
			%[1]s.sekolah.npsn AS NPSN,
			test_participants.pendataan_education_level AS KdJenjang,
			CASE WHEN %[1]s.sekolah.status_sekolah = 'N' THEN 'Negeri' ELSE 'Swasta' END AS StsSek,
			test_participants.siap_participant_school_id AS IDSekolah,
			test_participants.pendataan_class AS Kelas,
			test_participants.id AS IDSiswa,
			test_participants.name AS NamaSiswa,
			test_participants.siap_test_package_title AS PaketSoal,
			events.name AS Event,
			exam.name AS Ujian
			FROM 
			(
				SELECT 
					id,
					event_id,
					CONCAT(siap_participant_name, ' (', pendataan_user_code, ')') AS name,
					siap_participant_school_id,
					siap_test_package_title,
					pendataan_class,
					pendataan_education_level,
					exam_id
				FROM test_participants
				where package_id in (?)
			)test_participants
			LEFT JOIN events ON test_participants.event_id = events.id
			LEFT JOIN %[1]s.sekolah ON test_participants.siap_participant_school_id = %[1]s.sekolah.id
			LEFT JOIN exam ON test_participants.exam_id = exam.id
			ORDER BY %[1]s.sekolah.kode_provinsi, %[1]s.sekolah.kode_sekolah, test_participants.name,test_participants.siap_test_package_title
	`,schema)

	if 	err := db.Raw(query, packageIDs).Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}
