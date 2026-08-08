package model

import "time"

type AnalyzeRecord struct {
	Id             int       `gorm:"id"`
	JobDescription string    `gorm:"job_description"`
	Question       string    `gorm:"question"`
	Answer         string    `gorm:"answer"`
	CreatedAt      time.Time `gorm:"created_at"`
	UpdatedAt      time.Time `gorm:"updated_at"`
}

func (AnalyzeRecord) TableName() string {
	return "analyze_records"
}
