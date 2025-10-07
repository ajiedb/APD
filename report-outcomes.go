package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"ujicoba-go/utils"
	"strings"
)

type ReportOutcomes struct {
	ID              uint            `json:"id"`
	UUID            string          `json:"uuid"`
	ReqID           uint          	`json:"req_id"`
	Type 			int				`json:"type"`
	Name            string          `json:"name"`
	Status          uint          	`json:"status"`
	CreatedAt       string          `json:"created_at"`
	UpdatedAt       string          `json:"updated_at"`
	DeletedAt       *gorm.DeletedAt `json:"deleted_at"`
}

type ReportOutComesList struct {
	Type          	int          	`json:"type"`
	TypeName        string          `json:"typeName"`
	Name       		string          `json:"name"`
	Status       	string         	`json:"status"`
	UUID			string			`json:"uuid"`
}

type InCompleteReports struct {
	UUID          	string         	`json:"uuid"`
	Type       		int          	`json:"type"`
}


func (ReportOutcomes) CreateReportOutcomes(data []ReportOutcomes) (error) {
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

func (ReportOutcomes) UpdateReportOutcomes(fileName string, req_id string, report_type int) (error) {
	db := GetDBApd()
	dataNow := time.Now().Format(utils.DATE_FORMAT)
	query := `
		UPDATE report_outcomes ro
		JOIN report_reqs rr ON ro.req_id = rr.id
		SET ro.status = 1, 
			ro.name = ?,
			ro.updated_at = ?
		WHERE rr.uuid = ? 
		AND ro.type = ?
	`

	if err := db.Exec(query, fileName,dataNow,req_id,report_type).Error; err != nil {
		return err
	}
	return nil
}

func (ReportOutcomes) GetReportOutcomes(reqUUID string) ([]ReportOutComesList, error) {
	db := GetDBApd()
	var outcomes []ReportOutComesList
	query := `
		select 
		a.type, 
		a.name,
		a.status,
		b.uuid
		from report_outcomes a
		left join report_reqs b
		on a.req_id = b.id
		where b.uuid = ?
	`
	if err := db.Debug().Raw(query, reqUUID).Scan(&outcomes).Error; err != nil {
		return nil, err
	}
	reportTypeConstant := utils.ReportTypeConstant{}

	for i  := range outcomes {
		outcomes[i].TypeName = strings.ReplaceAll(*reportTypeConstant.ReportType(outcomes[i].Type), "_", " ")
	}

	return outcomes,nil
}

func (ReportOutcomes) GetInCompleteReport() ([]InCompleteReports, error) {
	db := GetDBApd()
	var outcomes []InCompleteReports
	query := `
		select 
		b.uuid,
		a.type
		from report_outcomes a
		left join report_reqs b
		on a.req_id  = b.id
		where a.status = 0
		order by b.created_at asc
		limit 1
	`
	if err := db.Raw(query).Scan(&outcomes).Error; err != nil {
		return nil, err
	}

	var uuid []string
	var repType []int
	for _, data := range outcomes {
		uuid = append(uuid,data.UUID)
		repType = append(repType, data.Type)
	}

	dataNow := time.Now().Format(utils.DATE_FORMAT)
	queryUpdate := `
		update report_outcomes ro
		left join report_reqs rr
		on ro.req_id =  rr.id
		set status =2,
		ro.updated_at = ?
		where rr.uuid in (?)
		and ro.type in (?)
	`

	if err := db.Exec(queryUpdate, dataNow,uuid,repType).Error; err != nil {
		return nil, err
	}

	return outcomes,nil
}