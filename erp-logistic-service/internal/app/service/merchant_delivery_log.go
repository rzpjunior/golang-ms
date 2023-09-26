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

type IMerchantDeliveryLogService interface {
	Create(ctx context.Context, req *logisticService.CreateMerchantDeliveryLogRequest) (res *dto.MerchantDeliveryLogResponse, err error)
	GetFirst(ctx context.Context, deliveryRunSheetItemId int64) (res *dto.MerchantDeliveryLogResponse, err error)
}

type MerchantDeliveryLogService struct {
	opt                           opt.Options
	RepositoryMerchantDeliveryLog repository.IMerchantDeliveryLogRepository
}

func NewMerchantDeliveryLogService() IMerchantDeliveryLogService {
	return &MerchantDeliveryLogService{
		opt:                           global.Setup.Common,
		RepositoryMerchantDeliveryLog: repository.NewMerchantDeliveryLogRepository(),
	}
}

func (s *MerchantDeliveryLogService) Create(ctx context.Context, req *logisticService.CreateMerchantDeliveryLogRequest) (res *dto.MerchantDeliveryLogResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "MerchantDeliveryLogService.Create")
	defer span.End()

	model := &model.MerchantDeliveryLog{
		Latitude:               req.Model.Latitude,
		Longitude:              req.Model.Longitude,
		CreatedAt:              req.Model.CreatedAt.AsTime().Unix(),
		DeliveryRunSheetItemId: req.Model.DeliveryRunSheetItemId,
	}

	if err = s.RepositoryMerchantDeliveryLog.Create(ctx, model); err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &dto.MerchantDeliveryLogResponse{
		Id:                     model.Id,
		Latitude:               model.Latitude,
		Longitude:              model.Longitude,
		CreatedAt:              time.Unix(model.CreatedAt, 0),
		DeliveryRunSheetItemId: model.DeliveryRunSheetItemId,
	}

	return
}

func (s *MerchantDeliveryLogService) GetFirst(ctx context.Context, deliveryRunSheetItemId int64) (res *dto.MerchantDeliveryLogResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "MerchantDeliveryLogService.GetFirst")
	defer span.End()

	var merchantDeliveryLog *model.MerchantDeliveryLog
	merchantDeliveryLog, err = s.RepositoryMerchantDeliveryLog.GetFirst(ctx, deliveryRunSheetItemId)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &dto.MerchantDeliveryLogResponse{
		Id:                     merchantDeliveryLog.Id,
		Latitude:               merchantDeliveryLog.Latitude,
		Longitude:              merchantDeliveryLog.Longitude,
		CreatedAt:              time.Unix(merchantDeliveryLog.CreatedAt, 0),
		DeliveryRunSheetItemId: merchantDeliveryLog.DeliveryRunSheetItemId,
	}

	return
}
