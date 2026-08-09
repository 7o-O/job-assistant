// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

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
	// Default to the first page when page is missing or invalid.
	page := req.Page
	if page <= 0 {
		page = 1
	}

	// Default to 10 records per page when pageSize is missing or invalid.
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	// Limit pageSize to protect the database from large queries.
	if pageSize > 50 {
		pageSize = 50
	}

	// Count total records so the client can calculate total pages.
	var total int64
	if err := l.svcCtx.DB.
		Model(&model.AnalyzeRecord{}).
		Count(&total).Error; err != nil {
		return &types.CommonResponse{
			Code:    500,
			Message: "query total failed: " + err.Error(),
			Success: false,
		}, nil
	}

	// Skip records from previous pages and read only the current page.
	offset := (page - 1) * pageSize

	// Query current-page records with newest records first.
	var records []model.AnalyzeRecord
	if err := l.svcCtx.DB.
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&records).Error; err != nil {
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
