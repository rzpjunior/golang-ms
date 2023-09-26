package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/internal/app/repository"
	logisticService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/logistic_service"
)

type IAddressCoordinateLogService interface {
	Get(ctx context.Context, req dto.AddressCoordinateLogGetRequest) (res []*dto.AddressCoordinateLogResponse, total int64, err error)
	GetDetail(ctx context.Context, id int64, deliveryRunSheetItemId int64) (res *dto.AddressCoordinateLogResponse, err error)
	Create(ctx context.Context, req *logisticService.CreateAddressCoordinateLogRequest) (res *dto.AddressCoordinateLogResponse, err error)
	GetMostTrusted(ctx context.Context, addressId string) (res *dto.AddressCoordinateLogResponse, err error)
}

type AddressCoordinateLogService struct {
	opt                            opt.Options
	RepositoryAddressCoordinateLog repository.IAddressCoordinateLogRepository
}

func NewAddressCoordinateLogService() IAddressCoordinateLogService {
	return &AddressCoordinateLogService{
		opt:                            global.Setup.Common,
		RepositoryAddressCoordinateLog: repository.NewAddressCoordinateLogRepository(),
	}
}

func (s *AddressCoordinateLogService) Get(ctx context.Context, req dto.AddressCoordinateLogGetRequest) (res []*dto.AddressCoordinateLogResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AddressCoordinateLogService.Get")
	defer span.End()

	var addressCoordinateLogs []*model.AddressCoordinateLog
	addressCoordinateLogs, total, err = s.RepositoryAddressCoordinateLog.Get(ctx, dto.AddressCoordinateLogGetRequest{
		OrderBy:          req.OrderBy,
		GroupBy:          req.GroupBy,
		ArrAddressIDs:    req.ArrAddressIDs,
		ArrSalesOrderIDs: req.ArrSalesOrderIDs,
	})
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, address := range addressCoordinateLogs {
		res = append(res, &dto.AddressCoordinateLogResponse{
			ID:             address.ID,
			Latitude:       address.Latitude,
			Longitude:      address.Longitude,
			LogChannelID:   address.LogChannelID,
			MainCoordinate: address.MainCoordinate,
			CreatedAt:      address.CreatedAt,
			CreatedBy:      address.CreatedBy,
			AddressID:      address.AddressID,
			SalesOrderID:   address.SalesOrderID,
		})
	}

	return
}

func (s *AddressCoordinateLogService) GetDetail(ctx context.Context, id int64, deliveryRunSheetItemId int64) (res *dto.AddressCoordinateLogResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AddressCoordinateLogService.GetDetail")
	defer span.End()

	var addressCoordinateLog *model.AddressCoordinateLog
	addressCoordinateLog, err = s.RepositoryAddressCoordinateLog.GetByID(ctx, id)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &dto.AddressCoordinateLogResponse{
		ID:             addressCoordinateLog.ID,
		Latitude:       addressCoordinateLog.Latitude,
		Longitude:      addressCoordinateLog.Longitude,
		LogChannelID:   addressCoordinateLog.LogChannelID,
		MainCoordinate: addressCoordinateLog.MainCoordinate,
		CreatedAt:      addressCoordinateLog.CreatedAt,
		CreatedBy:      addressCoordinateLog.CreatedBy,
		AddressID:      addressCoordinateLog.AddressID,
		SalesOrderID:   addressCoordinateLog.SalesOrderID,
	}

	return
}

func (s *AddressCoordinateLogService) Create(ctx context.Context, req *logisticService.CreateAddressCoordinateLogRequest) (res *dto.AddressCoordinateLogResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AddressCoordinateLogService.Create")
	defer span.End()

	model := &model.AddressCoordinateLog{
		Latitude:       *req.Model.Latitude,
		Longitude:      *req.Model.Longitude,
		LogChannelID:   int8(req.Model.LogChannelId),
		MainCoordinate: int8(req.Model.MainCoordinate),
		CreatedAt:      req.Model.CreatedAt.AsTime(),
		CreatedBy:      req.Model.CreatedBy,
		AddressID:      req.Model.AddressId,
		SalesOrderID:   req.Model.SalesOrderId,
	}

	if err = s.RepositoryAddressCoordinateLog.Create(ctx, model); err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &dto.AddressCoordinateLogResponse{
		ID:             model.ID,
		Latitude:       model.Latitude,
		Longitude:      model.Longitude,
		LogChannelID:   model.LogChannelID,
		MainCoordinate: model.MainCoordinate,
		CreatedAt:      model.CreatedAt,
		CreatedBy:      model.CreatedBy,
		AddressID:      model.AddressID,
		SalesOrderID:   model.SalesOrderID,
	}

	return
}

func (s *AddressCoordinateLogService) GetMostTrusted(ctx context.Context, addressId string) (res *dto.AddressCoordinateLogResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AddressCoordinateLogService.GetMostTrusted")
	defer span.End()

	var addressCoordinateLog *model.AddressCoordinateLog
	addressCoordinateLog, err = s.RepositoryAddressCoordinateLog.GetMostTrusted(ctx, addressId)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &dto.AddressCoordinateLogResponse{
		ID:             addressCoordinateLog.ID,
		Latitude:       addressCoordinateLog.Latitude,
		Longitude:      addressCoordinateLog.Longitude,
		LogChannelID:   addressCoordinateLog.LogChannelID,
		MainCoordinate: addressCoordinateLog.MainCoordinate,
		CreatedAt:      addressCoordinateLog.CreatedAt,
		CreatedBy:      addressCoordinateLog.CreatedBy,
		AddressID:      addressCoordinateLog.AddressID,
		SalesOrderID:   addressCoordinateLog.SalesOrderID,
	}

	return
}
