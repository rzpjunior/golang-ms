package service

import (
	"context"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-logistic-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/internal/app/repository"
	configurationService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
	logisticService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/logistic_service"
)

type IDeliveryRunSheetService interface {
	Get(ctx context.Context, req dto.DeliveryRunSheetGetRequest) (res []*dto.DeliveryRunSheetResponse, total int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (res *dto.DeliveryRunSheetResponse, err error)
	Create(ctx context.Context, req *logisticService.CreateDeliveryRunSheetRequest) (res *dto.DeliveryRunSheetResponse, err error)
	Finish(ctx context.Context, req *logisticService.FinishDeliveryRunSheetRequest) (res *dto.DeliveryRunSheetResponse, err error)
}

type DeliveryRunSheetService struct {
	opt                        opt.Options
	RepositoryDeliveryRunSheet repository.IDeliveryRunSheetRepository
}

func NewDeliveryRunSheetService() IDeliveryRunSheetService {
	return &DeliveryRunSheetService{
		opt:                        global.Setup.Common,
		RepositoryDeliveryRunSheet: repository.NewDeliveryRunSheetRepository(),
	}
}

func (s *DeliveryRunSheetService) Get(ctx context.Context, req dto.DeliveryRunSheetGetRequest) (res []*dto.DeliveryRunSheetResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DeliveryRunSheetService.Get")
	defer span.End()

	var deliveryRunSheets []*model.DeliveryRunSheet

	deliveryRunSheets, total, err = s.RepositoryDeliveryRunSheet.Get(ctx, dto.DeliveryRunSheetGetRequest{
		Offset:        req.Offset,
		Limit:         req.Limit,
		OrderBy:       req.OrderBy,
		GroupBy:       req.GroupBy,
		Status:        req.Status,
		ArrCourierIDs: req.ArrCourierIDs,
	})
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, drs := range deliveryRunSheets {
		res = append(res, &dto.DeliveryRunSheetResponse{
			ID:                drs.ID,
			Code:              drs.Code,
			DeliveryDate:      drs.DeliveryDate,
			StartedAt:         drs.StartedAt,
			FinishedAt:        drs.FinishedAt,
			StartingLatitude:  drs.StartingLatitude,
			StartingLongitude: drs.StartingLongitude,
			FinishedLatitude:  drs.FinishedLatitude,
			FinishedLongitude: drs.FinishedLongitude,
			Status:            drs.Status,
			CourierID:         drs.CourierID,
		})
	}

	return
}

func (s *DeliveryRunSheetService) GetDetail(ctx context.Context, id int64, code string) (res *dto.DeliveryRunSheetResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DeliveryRunSheetService.GetDetail")
	defer span.End()

	var deliveryRunSheet *model.DeliveryRunSheet
	deliveryRunSheet, err = s.RepositoryDeliveryRunSheet.GetByID(ctx, id, code)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &dto.DeliveryRunSheetResponse{
		ID:                deliveryRunSheet.ID,
		Code:              deliveryRunSheet.Code,
		DeliveryDate:      deliveryRunSheet.DeliveryDate,
		StartedAt:         deliveryRunSheet.StartedAt,
		FinishedAt:        deliveryRunSheet.FinishedAt,
		StartingLatitude:  deliveryRunSheet.StartingLatitude,
		StartingLongitude: deliveryRunSheet.StartingLongitude,
		FinishedLatitude:  deliveryRunSheet.FinishedLatitude,
		FinishedLongitude: deliveryRunSheet.FinishedLongitude,
		Status:            deliveryRunSheet.Status,
		CourierID:         deliveryRunSheet.CourierID,
	}

	return
}

func (s *DeliveryRunSheetService) Create(ctx context.Context, req *logisticService.CreateDeliveryRunSheetRequest) (res *dto.DeliveryRunSheetResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DeliveryRunSheetService.Create")
	defer span.End()

	var codeGenerator *configurationService.GetGenerateCodeResponse
	if codeGenerator, err = s.opt.Client.ConfigurationServiceGrpc.GetGenerateCode(ctx, &configurationService.GetGenerateCodeRequest{
		Format: "DRS",
		Domain: "delivery_run_sheet",
		Length: 6,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	model := &model.DeliveryRunSheet{
		Code:              codeGenerator.Data.Code,
		DeliveryDate:      req.Model.DeliveryDate.AsTime(),
		StartedAt:         req.Model.StartedAt.AsTime(),
		StartingLatitude:  req.Model.StartingLatitude,
		StartingLongitude: req.Model.StartingLongitude,
		Status:            int8(req.Model.Status),
		CourierID:         req.Model.CourierId,
	}

	if err = s.RepositoryDeliveryRunSheet.Create(ctx, model); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &dto.DeliveryRunSheetResponse{
		ID:                model.ID,
		Code:              model.Code,
		DeliveryDate:      model.DeliveryDate,
		StartedAt:         model.StartedAt,
		FinishedAt:        model.FinishedAt,
		StartingLatitude:  model.StartingLatitude,
		StartingLongitude: model.StartingLongitude,
		FinishedLatitude:  model.FinishedLatitude,
		FinishedLongitude: model.FinishedLongitude,
		Status:            model.Status,
		CourierID:         model.CourierID,
	}

	return
}

func (s *DeliveryRunSheetService) Finish(ctx context.Context, req *logisticService.FinishDeliveryRunSheetRequest) (res *dto.DeliveryRunSheetResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DeliveryRunSheetService.Finish")
	defer span.End()

	var model = &model.DeliveryRunSheet{
		ID:                req.Id,
		Status:            3,
		FinishedLatitude:  &req.Latitude,
		FinishedLongitude: &req.Longitude,
		FinishedAt:        time.Now(),
	}

	if err = s.RepositoryDeliveryRunSheet.Update(ctx, model, "Status", "FinishedLatitude", "FinishedLongitude", "FinishedAt"); err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &dto.DeliveryRunSheetResponse{
		ID:                model.ID,
		Status:            model.Status,
		FinishedLatitude:  &req.Latitude,
		FinishedLongitude: &req.Longitude,
		FinishedAt:        model.FinishedAt,
	}

	return
}
