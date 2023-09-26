package service

import (
	"context"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/internal/app/repository"
	logisticService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/logistic_service"
)

type IPostponeDeliveryLogService interface {
	Create(ctx context.Context, req *logisticService.CreatePostponeDeliveryLogRequest) (res *dto.PostponeDeliveryLogResponse, err error)
}

type PostponeDeliveryLogService struct {
	opt                           opt.Options
	RepositoryPostponeDeliveryLog repository.IPostponeDeliveryLogRepository
}

func NewPostponeDeliveryLogService() IPostponeDeliveryLogService {
	return &PostponeDeliveryLogService{
		opt:                           global.Setup.Common,
		RepositoryPostponeDeliveryLog: repository.NewPostponeDeliveryLogRepository(),
	}
}

func (s *PostponeDeliveryLogService) Create(ctx context.Context, req *logisticService.CreatePostponeDeliveryLogRequest) (res *dto.PostponeDeliveryLogResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PostponeDeliveryLogService.Create")
	defer span.End()

	model := &model.PostponeDeliveryLog{
		PostponeReason:         req.Model.PostponeReason,
		StartedAtUnix:          req.Model.StartedAtUnix,
		PostponedAtUnix:        req.Model.PostponedAtUnix,
		PostponeEvidence:       req.Model.PostponeEvidence,
		DeliveryRunSheetItemID: req.Model.DeliveryRunSheetItemId,
	}

	if err = s.RepositoryPostponeDeliveryLog.Create(ctx, model); err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &dto.PostponeDeliveryLogResponse{
		ID:                     model.ID,
		PostponeReason:         model.PostponeReason,
		StartedAt:              time.Unix(model.StartedAtUnix, 0),
		PostponedAt:            time.Unix(model.PostponedAtUnix, 0),
		PostponeEvidence:       model.PostponeEvidence,
		DeliveryRunSheetItemID: model.DeliveryRunSheetItemID,
	}

	return
}
