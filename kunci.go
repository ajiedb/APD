package models

import (
	"fmt"
	"encoding/json"
	"gorm.io/gorm"
)

type KunciCP struct {
	ID          uint            `json:"id"`
	UUID        string          `json:"uuid"`
	Description string          `json:"description"`
	JawabanID   uint            `json:"jawaban_id"`
	Score       uint            `json:"score"`
	DeletedAt   *gorm.DeletedAt `json:"deletedAt" gorm:"column:deletedAt"`
}

func (KunciCP) TableName() string {
	return "kunci_isian"
}

func (KunciCP) GetBySoalIDs(soalIDs []string, results *[]map[string]interface{}) error {
	db := GetDBSiap()
	tableName := KunciCP{}.TableName()

	query := db.Table(tableName).
		Model(KunciCP{}).
		Where("soal_id IN ?", soalIDs).
		Where("deletedAt IS NULL").
		Find(results)

	return query.Error
}

type KunciPG struct {
	ID          uint            `json:"id"`
	UUID        string          `json:"uuid"`
	Description string          `json:"description,omitempty"`
	JawabanID   uint            `json:"jawaban_id"`
	Order       uint            `json:"order"`
	DeletedAt   *gorm.DeletedAt `json:"deletedAt" gorm:"column:deletedAt"`
}

func (KunciPG) TableName() string {
	return "kunci_multiple_choices"
}

func (KunciPG) GetBySoalIDs(soalIDs []string, results *[]map[string]interface{}) error {
	db := GetDBSiap()
	tableName := KunciPG{}.TableName()

	query := db.Table(tableName).
		Model(KunciPG{}).
		Where("soal_id IN ?", soalIDs).
		Where("deletedAt IS NULL").
		Find(results)

	return query.Error
}

type KunciPGK1 struct {
	ID        uint   `json:"id"`
	UUID      string `json:"uuid"`
	SoalID    uint   `json:"soal_id"`
	JawabanID uint   `json:"jawaban_id"`
	Order     uint   `json:"order"`
}

func (KunciPGK1) TableName() string {
	return "kunci_pgk_singular"
}

func (KunciPGK1) GetBySoalIDs(soalIDs []string, results *[]map[string]interface{}) error {
	db := GetDBSiap()
	tableName := KunciPGK1{}.TableName()

	query := db.Table(tableName).
		Model(KunciPGK1{}).
		Where("soal_id IN ?", soalIDs).
		Where("deletedAt IS NULL").
		Order(fmt.Sprintf("%s.order ASC", tableName)).
		Find(results)

	return query.Error
}

type KunciPGK2 struct {
	ID           uint   `json:"id"`
	UUID         string `json:"uuid"`
	SoalID       uint   `json:"soal_id"`
	JawabanID    uint   `json:"jawaban_id"`
	PertanyaanID uint   `json:"pertanyaan_id"`
}

func (KunciPGK2) TableName() string {
	return "kunci_pgk_multiple"
}

func (KunciPGK2) GetBySoalIDs(soalIDs []string, results *[]map[string]interface{}) error {
	db := GetDBSiap()
	tableName := KunciPGK2{}.TableName()

	query := db.Table(tableName).
		Model(KunciPGK2{}).
		Where("soal_id IN ?", soalIDs).
		Where("deletedAt IS NULL").
		Find(results)

	return query.Error
}

type KunciPC struct {
	ID        uint   `json:"id"`
	UUID      string `json:"uuid"`
	SoalID    uint   `json:"soal_id"`
	JawabanID uint   `json:"jawaban_id"`
	Order     uint   `json:"order"`
}

func (KunciPC) TableName() string {
	return "kunci_pencocokan"
}

func (KunciPC) GetBySoalIDs(soalIDs []string, results *[]map[string]interface{}) error {
	db := GetDBSiap()
	tableName := KunciPC{}.TableName()

	query := db.Table(tableName).
		Model(KunciPC{}).
		Where("soal_id IN ?", soalIDs).
		Where("deletedAt IS NULL").
		Order(fmt.Sprintf("%s.order ASC", tableName)).
		Find(results)

	return query.Error
}

type JawabanKunci struct {
    JawabanKunci	string	`json:"jawabanKunci"`
}

func (JawabanKunci) GetJawabanKunciPG(soalIDs []int,) (map[string]interface{},error) {
	db := GetDBSiap()
	var rawJSON string
	query := `SELECT
		CONCAT('{', GROUP_CONCAT(
			CONCAT('"', t.soal_id, '":', data)
			SEPARATOR ','
		), '}') AS jawabanKunci
	FROM (
		SELECT
			c.soal_id,
			CONCAT(
				'{',
				'"opsKunci":', 
					CONCAT('{', GROUP_CONCAT(
						CONCAT(
							'"', c.uuid, '":',
							JSON_OBJECT('order', c.order, 'answer', e.uuid)
						) SEPARATOR ','
					), '}'),
				',"opsNum":1',
				'}'
			) AS data
		FROM jawaban_multiple_choices c
		LEFT JOIN (
			SELECT * 
			FROM kunci_multiple_choices
			WHERE deletedAt IS NULL
		) d ON c.soal_id = d.soal_id AND c.order = d.order
		LEFT JOIN jawaban_multiple_choices e
			ON d.order = e.order AND d.soal_id = e.soal_id
		WHERE c.soal_id IN (?)
		AND c.deletedAt IS NULL
		GROUP BY c.soal_id
	) AS t`

	if err := db.Raw(query, soalIDs).Scan(&rawJSON).Error; err != nil {
		return map[string]interface{}{}, err
	}

	var result map[string]interface{}
    if err := json.Unmarshal([]byte(rawJSON), &result); err != nil {
        return map[string]interface{}{}, err
    }

    return result, nil

}

func (JawabanKunci) GetJawabanKunciPC(soalIDs []int,) (map[string]interface{},error) {
	db := GetDBSiap()

	var rawJSON string
	query := `
			SELECT 
			CONCAT('{', GROUP_CONCAT(
				CONCAT('"', soal_id, '":', data)
				SEPARATOR ','
			), '}') AS jawabanKunci
		FROM (
			SELECT 
				c.soal_id,
				CONCAT(
					'{',
					'"opsKunci":[',
						GROUP_CONCAT(
							JSON_OBJECT(e.uuid, c.uuid)
							SEPARATOR ','
						),
					'],',
					'"opsNum":', COUNT(d.uuid), ',',
					'"opsOrder":{',
						GROUP_CONCAT(
							CONCAT('"', c.uuid, '":', c.order)
							SEPARATOR ','
						),
					'}',
					'}'
				) AS data
			FROM jawaban_pencocokan c
			LEFT JOIN (
				SELECT * 
				FROM kunci_pencocokan
				WHERE deletedAt IS NULL
			) d ON c.id = d.jawaban_id
			LEFT JOIN jawaban_pencocokan e 
				ON d.order = e.order AND d.soal_id = e.soal_id
			WHERE c.soal_id IN (?)
			AND c.deletedAt IS NULL
			GROUP BY c.soal_id
		) AS t;
    `

	if err := db.Raw(query, soalIDs).Scan(&rawJSON).Error; err != nil {
		return map[string]interface{}{}, err
	}

	var result map[string]interface{}
    if err := json.Unmarshal([]byte(rawJSON), &result); err != nil {
        return map[string]interface{}{}, err
    }

	return result, nil

}

func (JawabanKunci) GetJawabanKunciPGK1(soalIDs []int,) (map[string]interface{},error) {
	db := GetDBSiap()

	var rawJSON string
	query := `
			SELECT 
			CONCAT('{', GROUP_CONCAT(
				CONCAT('"', soal_id, '":', data)
				SEPARATOR ','
			), '}') AS jawabanKunci
		FROM (
			SELECT 
				c.soal_id,
				CONCAT(
					'{',
					'"opsKunci":[',
						GROUP_CONCAT(
							JSON_OBJECT(c.uuid, e.uuid)
							SEPARATOR ','
						),
					'],',
					'"opsNum":', COUNT(c.uuid),
					'}'
				) AS data
			FROM jawaban_pgk_singular c
			LEFT JOIN (
				SELECT * 
				FROM kunci_pgk_singular
				WHERE deletedAt IS NULL
			) d ON c.id = d.jawaban_id
			LEFT JOIN jawaban_pgk_singular e 
				ON d.order = e.order AND d.soal_id = e.soal_id
			WHERE c.soal_id IN (?)
			AND c.deletedAt IS NULL
			AND e.deletedAt IS NULL
			GROUP BY c.soal_id
		) AS t;
    `

	if err := db.Raw(query, soalIDs).Scan(&rawJSON).Error; err != nil {
		return map[string]interface{}{}, err
	}

	var result map[string]interface{}
    if err := json.Unmarshal([]byte(rawJSON), &result); err != nil {
        return map[string]interface{}{}, err
    }

	return result, nil

}

func (JawabanKunci) GetJawabanKunciPGK2(soalIDs []int,) (map[string]interface{},error) {
	db := GetDBSiap()

	var rawJSON string
	query := `
			SELECT 
			CONCAT('{', GROUP_CONCAT(
				CONCAT('"', soal_id, '":', data)
				SEPARATOR ','
			), '}') AS jawabanKunci
		FROM (
			SELECT 
				c.soal_id,
				CONCAT(
					'{',
					'"opsKunci":[',
						GROUP_CONCAT(
							JSON_OBJECT(a.uuid, c.uuid)
							SEPARATOR ','
						),
					'],',
					'"opsNum":', COUNT(a.uuid), ',',
					'"opsOrder":', d.opsOrder,
					'}'
				) AS data
			FROM pertanyaan_pgk_multiple a
			LEFT JOIN kunci_pgk_multiple b
				ON a.soal_id = b.soal_id AND a.id = b.pertanyaan_id
			LEFT JOIN jawaban_pgk_multiple c
				ON a.soal_id = c.soal_id AND b.jawaban_id = c.id
			LEFT JOIN (
				SELECT CONCAT(
						'{',
						GROUP_CONCAT(
							CONCAT('"', a.uuid, '":', a.order)
							SEPARATOR ','
						),
						'}'
					) AS opsOrder,
				a.soal_id
				FROM jawaban_pgk_multiple a
				WHERE a.deletedAt IS NULL
				GROUP BY a.soal_id
			) d ON a.soal_id = d.soal_id
			WHERE a.soal_id IN (?)
			AND a.deletedAt IS NULL
			AND b.deletedAt IS NULL
			GROUP BY c.soal_id
		) AS t;
    `

	if err := db.Raw(query, soalIDs).Scan(&rawJSON).Error; err != nil {
		return map[string]interface{}{}, err
	}

	var result map[string]interface{}
    if err := json.Unmarshal([]byte(rawJSON), &result); err != nil {
        return map[string]interface{}{}, err
    }

	return result, nil

}

func (JawabanKunci) GetJawabanKunciShort(soalIDs []int,) (map[string]interface{},error) {
	db := GetDBSiap()

	var rawJSON string
	query := `
			SELECT 
			CONCAT('{', GROUP_CONCAT(
				CONCAT('"', soal_id, '":', data)
				SEPARATOR ','
			), '}') AS jawabanKunci
		FROM (
			SELECT 
				c.soal_id,
				CONCAT(
					'{',
					'"opsKunci":[',
						GROUP_CONCAT(
							JSON_OBJECT(c.uuid, d.description)
							SEPARATOR ','
						),
					'],',
					'"opsNum":', COUNT(c.uuid),
					'}'
				) AS data
			FROM jawaban_isian c
			LEFT JOIN (
				SELECT * 
				FROM kunci_isian
				WHERE deletedAt IS NULL
				AND score = 1
			) d ON c.soal_id = d.soal_id 
			AND c.id = d.jawaban_id
			WHERE c.soal_id IN (?)
			AND c.deletedAt IS NULL
			GROUP BY c.soal_id
		) AS t;
    `

	if err := db.Raw(query, soalIDs).Scan(&rawJSON).Error; err != nil {
		return map[string]interface{}{}, err
	}

	var result map[string]interface{}
    if err := json.Unmarshal([]byte(rawJSON), &result); err != nil {
        return map[string]interface{}{}, err
    }

	return result, nil

}



