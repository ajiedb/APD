package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"ujicoba-go/utils"
	"fmt"
)


type ReportReqs struct {
	ID              uint            `json:"id"`
	UUID            string          `json:"uuid"`
	CreatedAt       string          `json:"created_at"`
	UpdatedAt       string          `json:"updated_at"`
	DeletedAt       *gorm.DeletedAt `json:"deleted_at"`
}

type ReportReqsList struct {
	UUID            string          `json:"uuid"`
	CreatedAt       string          `json:"created_at"`
	PackageNum      string         	`json:"packageNum"`
	Status       	string         	`json:"status"`
}

type ReportReqsListQuery struct {
	utils.ListQueryParams
	CreatedAt string `form:"created_at"`
	Status int `form:"status,default=99"`
}

type ReportReqsListResult struct {
	Reports []ReportReqsList
	Total  int64
}

func (ReportReqs) CreateReportReqs(data *ReportReqs) (string,error) {
	db := GetDBApd()
	uuid := uuid.New().String()
	dataNow := time.Now().Format("2006-01-02 15:04:05")

	data.CreatedAt = dataNow
	data.UpdatedAt = dataNow
	data.UUID = uuid

	if err := db.Create(&data).Error; err != nil {
		return "",err
	}

	return data.UUID,nil
}

func (ReportReqsList) List(query *ReportReqsListQuery) (ReportReqsListResult, error) {
	db := GetDBApd()
	var reports []ReportReqsList
	offset := (query.Page - 1) * query.PageSize

	var orderValue string
	if query.Order == "0" {
		orderValue = "ASC"
	} else if query.Order == "1" {
		orderValue = "DESC"
	}
	orderStr := fmt.Sprintf("%s %s", query.Orderby, orderValue)

	subquery := db.
		Table("report_outcomes").
		Select("req_id, CASE WHEN SUM(case when status = 1 then 1 else 0 end) > 0 THEN 1 ELSE 0 END AS status").
		Group("req_id")
	
	queryModel := db.
		Table("report_reqs AS a").
		Select(`
			a.uuid,
			a.created_at,
			COUNT(b.siap_test_package_uuid) AS package_num,
			COALESCE(c.status, 0) AS status
		`).
		Joins("LEFT JOIN report_packages b ON a.id = b.req_id").
		Joins("LEFT JOIN (?) c ON a.id = c.req_id", subquery)
	
	if query.Status != 99 {
		queryModel = queryModel.Where("COALESCE(c.status, 0) = ?",query.Status)
	}

	if query.CreatedAt != "" {
		queryModel = queryModel.Where("DATE(a.created_at) = ?",query.CreatedAt)
	}

	// Main query with joins and aggregation
	err := queryModel.Debug().
		Group("a.uuid, a.created_at, c.status").
		Order("a." + orderStr).
		Limit(query.PageSize).
		Offset(offset).
		Scan(&reports).Error

	if err != nil {
		return ReportReqsListResult{}, err
	}

	var total int64
	countQuery := db.
		Table("report_reqs AS a").
		Select("a.uuid").
		Joins("LEFT JOIN report_packages b ON a.id = b.req_id").
		Joins("LEFT JOIN (?) c ON a.id = c.req_id", subquery)
		

	if query.Status != 99 {
		countQuery = countQuery.Where("COALESCE(c.status, 0) = ?",query.Status)
	}

	if query.CreatedAt != "" {
		queryModel = queryModel.Where("DATE(a.created_at) = ?",query.CreatedAt)
	}

	countQuery = countQuery.Group("a.uuid")

	err = db.Table("(?) as sub", countQuery).
		Count(&total).Error

	if err != nil {
		return ReportReqsListResult{}, err
	}

	return ReportReqsListResult{
		Reports: reports,
		Total: total,
	}, nil
}

