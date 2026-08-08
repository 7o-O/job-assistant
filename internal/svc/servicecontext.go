// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"job-assistant/internal/config"
	"job-assistant/internal/llm"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
	LLM    *llm.Client
}

func NewServiceContext(c config.Config) (*ServiceContext, error) {
	db, err := gorm.Open(mysql.Open(c.MySQL.DataSource), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	apiKey := os.Getenv("LLM_API_KEY")
	llmClient := llm.NewClient(c.LLM.BaseURL, apiKey, c.LLM.Model)
	return &ServiceContext{
		Config: c,
		DB:     db,
		LLM:    llmClient,
	}, nil
}
