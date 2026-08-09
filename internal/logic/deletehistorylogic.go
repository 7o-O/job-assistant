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

type DeleteHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteHistoryLogic {
	return &DeleteHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteHistoryLogic) DeleteHistory(req *types.DeleteHistoryRequest) (resp *types.CommonResponse, err error) {
	if req.Id <= 0 {
		return &types.CommonResponse{
			Code:    400,
			Message: "没有该id",
			Success: false,
		}, nil
	}

	result := l.svcCtx.DB.Delete(&model.AnalyzeRecord{}, req.Id)

	if result.Error != nil {
		return &types.CommonResponse{
			Code:    500,
			Message: "错误: " + result.Error.Error(),
			Success: false,
		}, nil
	}

	if result.RowsAffected == 0 {
		return &types.CommonResponse{
			Code:    404,
			Message: "没有找到该历史",
			Success: false,
		}, nil
	}

	return &types.CommonResponse{
		Code:    200,
		Message: "删除成功",
		Success: true,
	}, nil
}
