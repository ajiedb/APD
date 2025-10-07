package models

import (
	"database/sql"
	"time"
	"gorm.io/gorm"
)

type TestParticipantLogs struct {
	ID                 uint           `json:"id"`
	TestParticipantID  uint           `json:"test_participant_id"`
	TestAnswerID       sql.NullInt64  `json:"test_answer_id"`
	Action             sql.NullInt16  `json:"action"`
	CreatedAt          string         `json:"created_at"`
	UpdatedAt          string         `json:"updated_at"`
	IpClient           sql.NullString `json:"ip_client"`
	IpServer           sql.NullString `json:"ip_server"`
	SiapQuestionCode   sql.NullString `json:"siap_question_code"`
	NextQuestionCode   sql.NullString `json:"next_question_code"`
	SiapQuestionNumber sql.NullInt16  `json:"siap_question_number"`
	NextQuestionNumber sql.NullInt16  `json:"next_question_number"`
	IsInDoubt          sql.NullInt16  `json:"is_in_doubt"`
	RemainingTime      sql.NullInt64  `json:"remaining_time"`
	Second             sql.NullInt64  `json:"second"`
	Answer             *string        `json:"answer"`
}

func (TestParticipantLogs) TableName() string {
	return "test_participant_logs"
}

type TestParticipantLogView struct {
	ID                 uint   `json:"id"`
	TestParticipantID  uint   `json:"test_participant_id"`
	TestAnswerID       int64  `json:"test_answer_id"`
	Action             int16  `json:"action"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
	CreatedTime        string `json:"created_time"`
	UpdatedTime        string `json:"updated_time"`
	IpClient           string `json:"ip_client"`
	IpServer           string `json:"ip_server"`
	SiapQuestionCode   string `json:"siap_question_code"`
	NextQuestionCode   string `json:"next_question_code"`
	SiapQuestionNumber int16  `json:"siap_question_number"`
	NextQuestionNumber int16  `json:"next_question_number"`
	IsInDoubt          int16  `json:"is_in_doubt"`
	RemainingTime      int16  `json:"remaining_time"`
	Answer             string `json:"answer"`
}

type TestParticipantLogReport struct {
	TestParticipantID 	int 	`json:"test_participant_id"`
	SiapQuestionCode	string	`json:"siap_question_code"`
    NumOfSee			int		`json:"numOfSee"`
	TimeOfSee			int		`json:"timeOfSee"`	
}

func (TestParticipantLogView) TableName() string {
	return "test_participant_logs"
}

func (tLogs *TestParticipantLogs) Save() error {
	db := GetDBApd()

	now := time.Now().Format("2006-01-02 15:04:05")
	tLogs.CreatedAt = now
	tLogs.UpdatedAt = now

	if err := db.Create(&tLogs).Error; err != nil {
		return err
	}
	return nil
}

func (tLogs *TestParticipantLogs) SaveWithTx(txApd *gorm.DB) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	tLogs.CreatedAt = now
	tLogs.UpdatedAt = now

	if err := txApd.Create(&tLogs).Error; err != nil {
		return err
	}
	return nil
}

func (tLogs *TestParticipantLogs) UpdateWithTx(txApd *gorm.DB) error {
	now := time.Now().Format("2006-01-02 15:04:05")

	tLogs.UpdatedAt = now

	if err := txApd.Where("id = ?", tLogs.ID).Updates(&tLogs).Error; err != nil {
		return err
	}

	return nil
}

func (tLogs *TestParticipantLogs) GetByParticipant() ([]TestParticipantLogs, error) {
	db := GetDBApd()
	var results []TestParticipantLogs
	if err := db.
		Where("test_participant_id", tLogs.TestParticipantID).
		Order("id DESC").
		Find(&results).Error; err != nil {
		return results, err
	}

	return results, nil
}

func (tLogs *TestParticipantLogs) GetByParticipants(testparticipantIds []int) ([]TestParticipantLogs, error) {
	db := GetDBApd()
	var results []TestParticipantLogs
	if err := db.
		Where("test_participant_id in (?)", testparticipantIds).
		Order("id DESC").
		Find(&results).Error; err != nil {
		return results, err
	}

	return results, nil
}

func (TestParticipantLogReport) GetParticipantLog(siswaID []int,) (/* [map[string]interface{}] */[]TestParticipantLogReport,error) {
	db := GetDBApd()
	
	var res []TestParticipantLogReport
	
	query := ` SELECT 
			l.test_participant_id,
			b.siap_question_code,
			COUNT(l.id) AS NumOfSee,
			SUM(l.second) AS TimeOfSee
	FROM test_participant_logs l
	JOIN test_answers b 
	ON l.test_answer_id = b.id
	WHERE l.test_participant_id IN (?)
	GROUP BY l.test_participant_id, b.siap_question_code`

	if err := db.Raw(query, siswaID).Scan(&res).Error; err != nil {
		return nil, err
	}

    return res, nil

}
