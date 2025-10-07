package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"fmt"
)

type ReportPackages struct {
	ID              uint            `json:"id"`
	UUID            string          `json:"uuid"`
	ReqID           uint          	`json:"req_id"`
	SiapTestPackageUUID string			`json:"siap_test_package_uuid"`
	CreatedAt       string          `json:"created_at"`
	UpdatedAt       string          `json:"updated_at"`
	DeletedAt       *gorm.DeletedAt `json:"deleted_at"`
}

func (ReportPackages) CreateReportPackages(data []ReportPackages) (error) {
	db := GetDBApd()
	dataNow := time.Now().Format("2006-01-02 15:04:05")

	for i, _ := range data {
		uuid := uuid.New().String()
		data[i].CreatedAt = dataNow
		data[i].UpdatedAt = dataNow
		data[i].UUID = uuid
	}
	
	if err := db.Create(&data).Error; err != nil {
		return err
	}

	return nil
}

func (ReportPackages) GetByReqUUID(uuid string) ([]string, error) {
	db := GetDBApd()
	var reportPackages []string

	if err := db.
		Model(&ReportPackages{}).
		Select("report_packages.siap_test_package_uuid").
		Joins("left join `report_reqs` on report_packages.req_id = report_reqs.id").
		Where("report_reqs.uuid = ?", uuid).
		Where("report_packages.deleted_at is null").
		Find(&reportPackages).
		Error; err != nil {
		return nil, err
	}

	fmt.Println(reportPackages)

	return reportPackages, nil
}
