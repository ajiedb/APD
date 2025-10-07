package models

import (
	"database/sql"
	"encoding/json"
	"fmt"

	// "strconv"
	"ujicoba-go/config"
	"ujicoba-go/utils"

	"gorm.io/gorm"
)

type PackageDetailSoal struct {
	ID                    uint            `json:"id"`
	PackageID             uint            `json:"package_id"`
	SoalID                uint            `json:"soal_id"`
	Order                 int             `json:"order"`
	CreatedAt             string          `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt             string          `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt             *gorm.DeletedAt `gorm:"column:deletedAt" json:"-"`
	SoalUUID              string          `json:"soal_uuid"`
	Description           *string         `json:"description"`
	AssignmentID          *int16          `json:"assignment_id"`
	SoalTypeID            *int16          `json:"soal_type_id"`
	SoalTypeCode          *string         `json:"soal_type_code"`
	SoalParent            *int16          `json:"soal_parent"`
	MataPelajaranID       *int16          `json:"mata_pelajaran_id"`
	FrameworkContentID    *int16          `json:"framework_content_id"`
	KonteksID             *int16          `json:"konteks_id"`
	KognitifID            *int16          `json:"kognitif_id"`
	StimulusID            *int64          `json:"stimulus_id"`
	StimulusTabID         *int16          `json:"stimulus_tab_id"`
	Variant               *int16          `json:"variant"`
	BookSource            *string         `json:"book_source"`
	Code                  *string         `json:"code"`
	Indicator             *string         `json:"indicator"`
	IndicatorID           *int64          `json:"indicator_id"`
	Rating                *float64        `json:"rating"`
	Status                *int16          `json:"status"`
	StatusRevisi          *int16          `json:"status_revisi"`
	IsRevisi              *int16          `json:"is_revisi"`
	IsRepository          *int16          `json:"is_repository"`
	WriterID              *int16          `json:"writer_id"`
	WriterName            *string         `json:"writer_name"`
	LockStatus            *string         `json:"lock_status"`
	LockUser              *string         `json:"lock_user"`
	LastLockStatus        *string         `json:"last_lock_status"`
	GlobalLocked          *int16          `json:"global_locked"`
	GlobalLockedBy        *int16          `json:"global_locked_by"`
	GlobalLockedNotes     *string         `json:"global_locked_notes"`
	PermanentLocked       *int16          `json:"permanent_locked"`
	PermanentLockedReason *string         `json:"permanent_locked_reason"`
	SoalCreatedAt         *string         `json:"soal_created_at"`
	SoalUpdatedAt         *string         `json:"soal_updated_at"`
	UpdatedByID           *string         `json:"updated_by_id" gorm:"column:updatedById"`
	UpdatedByName         *string         `json:"update_by_name" gorm:"column:updatedByName"`
	RevisionByID          *int16          `json:"revision_by_id"`
	RevisionByName        *string         `json:"revision_by_name"`
	SoalDeletedAt         *gorm.DeletedAt `json:"soal_deleted_at"`
}
type QuestionPackageSoal struct {
	ID                    uint                `json:"id"`
	UUID                  string              `json:"uuid"`
	PackageName           sql.NullString      `json:"package_name"`
	MataPelajaranID       sql.NullInt64       `json:"MataPelajaranID"`
	CreatedBy             sql.NullInt64       `gorm:"column:createdBy" json:"-"`
	CreatedAt             sql.NullString      `gorm:"column:createdAt;" json:"createdAt"`
	UpdatedAt             sql.NullString      `gorm:"column:updatedAt;" json:"-"`
	ExpiredDate           sql.NullString      `json:"-"`
	QuestionPackageDetail []PackageDetailSoal `gorm:"foreignKey:PackageID" json:"package_detail"`
}

type QuestionPackageSoalReport struct {
	SoalID                int             	 	`json:"soal_id"`
	SoalUUID           	  string	      		`json:"soal_UUID"`
	SoalTypeID       	  int       			`json:"soal_type_id"`
	SoalCode              string      			`json:"soal_code"`
}

func (QuestionPackageSoal) TableName() string {
	return "question_package"
}

func (PackageDetailSoal) TableName() string {
	return "question_package_detail"
}

func (qps *QuestionPackageSoal) ListQuestions() (QuestionPackageSoal, error) {
	db := GetDBSiap()
	var result QuestionPackageSoal

	query := db.
		Where("question_package.uuid = ?", qps.UUID).
		Preload("QuestionPackageDetail").
		Find(&result)

	if err := query.Error; err != nil {
		return result, err
	}

	return result, nil
}

func (qps *QuestionPackageSoal) GetByUUID() (QuestionPackageSoal, error) {
	db := GetDBSiap()
	var result QuestionPackageSoal

	query := db.
		Where("question_package.uuid = ?", qps.UUID).
		First(&result)

	if err := query.Error; err != nil {
		return result, err
	}

	return result, nil
}

func (QuestionPackageSoal) GetPackageDetail(packageId uint) ([]PackageDetailSoal, error) {
	db := GetDBSiap()
	var results []PackageDetailSoal

	query := db.
		Select("question_package_detail.*, soals.uuid as soal_uuid, soals.description, soals.description64, soals.assignment_id, soals.soal_type_id, soals.soal_parent, soals.mata_pelajaran_id, soals.framework_content_id, soals.konteks_id, soals.kognitif_id, soals.stimulus_id, soals.stimulus_tab_id, soals.variant, soals.book_source, soals.code, soals.indicator, soals.indicator_id, soals.rating, soals.status, soals.status_revisi, soals.is_revisi, soals.is_repository, soals.writer_id, soals.writer_name, soals.lock_status, soals.lock_user, soals.last_lock_status, soals.global_locked, soals.global_locked_by, soals.global_locked_notes, soals.permanent_locked, soals.permanent_locked_reason, soals.createdAt as soal_created_at, soals.updatedAt as soal_updated_at, soals.updatedById, soals.updatedByName, soals.revision_by_id, soals.revision_by_name, soals.deletedAt as soal_deleted_at, soal_types.code as soal_type_code").
		Joins("JOIN soals ON soals.id = question_package_detail.soal_id").
		Joins("JOIN soal_types ON soals.soal_type_id = soal_types.id").
		Where("question_package_detail.package_id = ?", packageId).
		Order("question_package_detail.order ASC").
		Find(&results)

	if err := query.Error; err != nil {
		return results, err
	}

	return results, nil
}

func (QuestionPackageSoal) GetSoalPackage(packageId uint) ([]PackageDetailSoal, error) {
	db := GetDBSiap()
	var results []PackageDetailSoal

	query := db.
		Where("question_package_detail.package_id = ?", packageId).
		Order("question_package_detail.order ASC").
		Find(&results)

	if err := query.Error; err != nil {
		return results, err
	}

	return results, nil
}

type resultParsingQuestions struct {
	Kode  string           `json:"kode"`
	Code  string           `json:"code"`
	Nomor int              `json:"nomor"`
	// Soal  ParseSoals `json:"soal"`
}

func (QuestionPackageSoal) GetCountSoalByParticipants(packageIDs []string, eventID string) (int64, error) {
	if config.EnvVariable("APP_ENV") != "vm" {
		var total int64

		db := GetDBSiap()
		for _, item := range packageIDs {
			var result map[string]interface{}
			query := db.
				Model(QuestionPackage{}).
				Select("COUNT(qpd.id) as total").
				Joins("LEFT JOIN question_package_detail as qpd ON qpd.package_id = question_package.id").
				Where("question_package.uuid = ?", item).
				Where("qpd.deletedAt is NULL").
				First(&result)
			if err := query.Error; err != nil {
				return 0, err
			}
			if result["total"].(int64) > 0 {
				total += result["total"].(int64)
			}
		}

		return total, nil
	}
	var totalSoal int
	redisService := utils.RedisService{}
	for _, items := range packageIDs {
		key := fmt.Sprintf("%s:%s", utils.PREFIX_QUESTIONS, items)
		packageDetail, _ := redisService.Get(key)
		//
		var parsePackageQuestions []resultParsingQuestions
		json.Unmarshal([]byte(packageDetail), &parsePackageQuestions)

		// cek detail per package
		for _, question := range parsePackageQuestions {
			questionKey := fmt.Sprintf("%s:%s", utils.PREFIX_QUESTION_DETAIL, question.Code)
			questionDetail, _ := redisService.CheckKey(questionKey)
			if(questionDetail > 0) {
				totalSoal++
			}
		}
		// totalSoal += len(parsePackageQuestions)
	}
	return int64(totalSoal), nil
}

func inSlice(needle int, haystack []int) bool {
	for _, value := range haystack {
		if value == needle {
			return true
		}
	}
	return false
}

func (QuestionPackageSoalReport) GetQuestionPackageSoalReport(packageList []string, typeList []int) ([]QuestionPackageSoalReport, error) {
	db := GetDBSiap().Debug()
	var results []QuestionPackageSoalReport
	query := `
	SELECT distinct 
		question_package_detail.soal_id,
		soals.uuid as soal_uuid,
		soals.soal_type_id,
		soals.code as soal_code
	FROM question_package
	LEFT JOIN question_package_detail
		ON question_package.id = question_package_detail.package_id
	LEFT JOIN soals
		ON soals.id = question_package_detail.soal_id
	WHERE question_package.uuid IN (?)
		AND question_package.deletedAt IS NULL
		AND question_package_detail.deletedAt IS NULL
		AND soals.deletedAt IS NULL
	`
	args := []interface{}{packageList}

	if(len(typeList) > 0){
		query += " AND soals.soal_type_id IN ?"
    	args = append(args, typeList)
	}

	if 	err :=  db.Raw(query,args...).Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}
