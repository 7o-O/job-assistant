// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"strings"

	"job-assistant/internal/llm"
	"job-assistant/internal/model"
	"job-assistant/internal/svc"
	"job-assistant/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AnalyzeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAnalyzeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AnalyzeLogic {
	return &AnalyzeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AnalyzeLogic) Analyze(req *types.AnalyzeRequest) (resp *types.CommonResponse, err error) {
	//校验参数
	if strings.TrimSpace(req.JobDescription) == "" {
		return &types.CommonResponse{
			Code:    400,
			Message: "描述不能为空",
			Success: false,
		}, nil
	}

	if strings.TrimSpace(req.Question) == "" {
		return &types.CommonResponse{
			Code:    400,
			Message: "问题不能为空",
			Success: false,
		}, nil
	}

	//发送给大模型
	messages := []llm.Message{
		{
			Role:    "system",
			Content: "你是一名专业的求职顾问，请使用中文回答，分析要具体、实用，适合求职者理解。",
		},
		{
			Role:    "user",
			Content: "岗位描述：\n" + req.JobDescription + "\n\n用户问题:\n" + req.Question,
		},
	}

	//调用
	answer, err := l.svcCtx.LLM.Chat(l.ctx, messages)
	if err != nil {
		return &types.CommonResponse{
			Code:    500,
			Message: "大模型调用失败: " + err.Error(),
			Success: false,
		}, nil
	}

	//保存到数据库
	record := model.AnalyzeRecord{
		JobDescription: req.JobDescription,
		Question:       req.Question,
		Answer:         answer,
	}

	if err := l.svcCtx.DB.Create(&record).Error; err != nil {
		return &types.CommonResponse{
			Code:    500,
			Message: "保存数据库失败: " + err.Error(),
			Success: false,
		}, nil
	}

	return &types.CommonResponse{
		Code:    200,
		Message: "success",
		Success: true,
		Data: map[string]interface{}{
			"id":     record.Id,
			"answer": record.Answer,
		},
	}, nil

}
