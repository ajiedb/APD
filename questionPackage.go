package models

import (
	"fmt"
	"strconv"
	"strings"
)

type QuestionPackage struct {
	ID                    uint                    `json:"id"`
	UUID                  string                  `json:"uuid"`
	PackageName           string                  `json:"package_name"`
	PackageType           int                     `json:"-"`
	MataPelajaranID       int64                   `json:"mata_pelajaran_id"`
	NamaMataPelajaran     string                  `json:"nama_mata_pelajaran"`
	CreatedBy             int64                   `json:"-"`
	CreatedAt             string                  `gorm:"column:createdAt;" json:"createdAt"`
	UpdatedAt             string                  `gorm:"column:updatedAt;" json:"-"`
	ExpiredDate           string                  `gorm:"column:expired_date" json:"-"`
	QuestionPackageDetail []QuestionPackageDetail `gorm:"foreignKey:PackageID" json:"package_detail"`
}

type QuestionPackageResponse struct {
	ID              uint   `json:"id"`
	UUID            string `json:"uuid"`
	MataPelajaranId int    `json:"mata_pelajaran_id"`
	PackageName     string `json:"package_name"`
	JumlahSoal      int    `json:"jumlah_soal"`
	CreatedAt       string `gorm:"column:createdAt;" json:"createdAt"`
	UpdatedAt       string `gorm:"column:updatedAt;" json:"UpdatedAt"`
	Tahun           int    `json:"tahun"`
}

func (QuestionPackageResponse) TableName() string {
	return "question_package"
}

func (QuestionPackage) TableName() string {
	return "question_package"
}

func (QuestionPackage) GetByUUIDs(uuids []string) ([]QuestionPackage, error) {
	db := GetDBSiap()

	var results []QuestionPackage
	query := db.Where("uuid IN ?", uuids).Find(&results)
	if err := query.Error; err != nil {
		return results, err
	}

	return results, nil
}

func (QuestionPackage) GetByIDsWithMatpel(ids ...string) ([]QuestionPackage, error) {
	db := GetDBSiap()
	idsStr := strings.Join(ids, ",")
	strQery := fmt.Sprintf("question_package.id IN (%s)", idsStr)

	var questionPackages []QuestionPackage
	queryModel := db.Model(QuestionPackage{}).
		Preload("QuestionPackageDetail").
		Select("question_package.*, lib_mata_pelajaran.title as nama_mata_pelajaran").
		Joins("JOIN lib_mata_pelajaran ON question_package.mata_pelajaran_id = lib_mata_pelajaran.id").
		Where(strQery)

	err := queryModel.Find(&questionPackages).Error
	if err != nil {
		return nil, err
	}

	return questionPackages, nil
}

func (QuestionPackage) GetAll() ([]QuestionPackage, error) {
	db := GetDBSiap()

	var questionPackages []QuestionPackage
	queryModel := db.Model(QuestionPackage{}).
		Preload("QuestionPackageDetail").
		Select("question_package.*, lib_mata_pelajaran.title as nama_mata_pelajaran").
		Joins("JOIN lib_mata_pelajaran ON question_package.mata_pelajaran_id = lib_mata_pelajaran.id")

	err := queryModel.Find(&questionPackages).Error
	if err != nil {
		return nil, err
	}

	return questionPackages, nil
}

func (QuestionPackage) GetByID(id string) ([]QuestionPackage, error) {
	db := GetDBSiap()

	var questionPackages []QuestionPackage
	queryModel := db.Model(QuestionPackage{}).
		Preload("QuestionPackageDetail").
		Select("question_package.*, lib_mata_pelajaran.title as nama_mata_pelajaran").
		Joins("JOIN lib_mata_pelajaran ON question_package.mata_pelajaran_id = lib_mata_pelajaran.id").
		Where("question_package.id = ?", id)

	err := queryModel.Find(&questionPackages).Error
	if err != nil {
		return nil, err
	}

	return questionPackages, nil
}

func (QuestionPackage) GetPackageValidate(packages []int) (int64, error) {
	db := GetDBSiap()

	query := db.Model(&QuestionPackage{}).Where("id IN (?)", packages)

	var count int64
	err := query.Count(&count).Error
	if err != nil {
		return count, err
	}

	return count, err
}

func GetPerakitan(keyword string, page string, mataPelajaranId int, year int, pageSize int, reqUUID string) ([]QuestionPackageResponse, error) {
	db := GetDBSiap()

	var perakitan []QuestionPackageResponse
	row := pageSize
	if page == "" {
		page = "1"
	}

	offset, e := strconv.Atoi(page)
	if e == nil {
		offset = (offset - 1) * row
	}

	query := db.Model(&QuestionPackageResponse{}).Offset(offset).Limit(row)

	if keyword != "" {
		query.Where("package_name LIKE ? ", "%"+keyword+"%")
	}

	if mataPelajaranId != 0 {
		query.Where("mata_pelajaran_id = ? ", mataPelajaranId)
	}

	if year != 0 {
		yearString := strconv.Itoa(year)
		startDate := yearString + "-01-01"
		endDate := yearString + "-12-31"
		query.Where("createdAt Between ? and ?", startDate, endDate)
	}

	if reqUUID != "" {
		repPackages := ReportPackages{}
		dataRepPackages, _:= repPackages.GetByReqUUID(reqUUID)
		query.Where("question_package.uuid in (?)",dataRepPackages)
	}

	err := query.Debug().Select("id,uuid ,mata_pelajaran_id, package_name, createdAt, updatedAt, YEAR(createdAt) AS tahun, (select count(id) from question_package_detail where package_id = question_package.id and deletedAt is null) as jumlah_soal").Find(&perakitan).Error
	if err != nil {
		return perakitan, err
	}

	for i := range perakitan {
		perakitan[i].CreatedAt = convertDatetimeEvent(perakitan[i].CreatedAt)
		perakitan[i].UpdatedAt = convertDatetimeEvent(perakitan[i].UpdatedAt)
	}

	return perakitan, err
}

func GetPerakitanCount(keyword string,mataPelajaranId int, year int, reqUUID string) int64 {
	db := GetDBSiap()

	var perakitan QuestionPackage
	query := db.Model(&perakitan)

	if keyword != "" {
		query.Where("package_name LIKE ?", "%"+keyword+"%")
	}

	if reqUUID !=""{
		repPackages := ReportPackages{}
		dataRepPackages, _ := repPackages.GetByReqUUID(reqUUID)
		query.Where("question_package.uuid in (?)",dataRepPackages)
	}

	var count int64
	err := query.Count(&count).Error
	if err != nil {
		return count
	}

	return count
}

func (QuestionPackage) GetByIDs(ids []uint) ([]QuestionPackage, error) {
	db := GetDBSiap()
	var results []QuestionPackage
	query := db.Where("id in ?", ids)
	if err := query.Find(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}
