// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"strings"

	"job-assistant/internal/model"
	"job-assistant/internal/svc"
	"job-assistant/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HistoryLogic {
	return &HistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HistoryLogic) History(req *types.HistoryRequest) (resp *types.CommonResponse, err error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}

	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	if pageSize > 50 {
		pageSize = 50
	}

	query := l.svcCtx.DB.Model(&model.AnalyzeRecord{})

	//关键字查询
	if strings.TrimSpace(req.KeyWord) != "" {
		keyword := "%" + strings.TrimSpace(req.KeyWord) + "%"
		query = query.Where("job_description LIKE ? OR question LIKE ? OR answer LIKE ?", keyword, keyword, keyword)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return &types.CommonResponse{
			Code:    500,
			Message: "query total failed: " + err.Error(),
			Success: false,
		}, nil
	}

	offset := (page - 1) * pageSize

	// 查询历史
	var records []model.AnalyzeRecord
	if err := query.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&records).Error; err != nil {
		return &types.CommonResponse{
			Code:    500,
			Message: "query history failed: " + err.Error(),
			Success: false,
		}, nil
	}

	return &types.CommonResponse{
		Code:    200,
		Message: "success",
		Success: true,
		Data: map[string]interface{}{
			"list":     records,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	}, nil
}
