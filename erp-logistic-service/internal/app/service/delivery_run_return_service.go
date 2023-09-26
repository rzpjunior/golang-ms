package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-logistic-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/internal/app/repository"

	configurationService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
	logisticService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/logistic_service"
)

type IDeliveryRunReturnService interface {
	Get(ctx context.Context, req dto.DeliveryRunReturnGetRequest) (res []*dto.DeliveryRunReturnResponse, total int64, err error)
	GetDetail(ctx context.Context, id int64, code string, deliveryRunSheetItemId int64) (res *dto.DeliveryRunReturnResponse, err error)
	Create(ctx context.Context, req *logisticService.CreateDeliveryRunReturnRequest) (res *dto.DeliveryRunReturnResponse, err error)
	Update(ctx context.Context, req *logisticService.UpdateDeliveryRunReturnRequest) (res *dto.DeliveryRunReturnResponse, err error)
	Delete(ctx context.Context, id int64) (res *dto.DeliveryRunReturnResponse, err error)
}

type DeliveryRunReturnService struct {
	opt                         opt.Options
	RepositoryDeliveryRunReturn repository.IDeliveryRunReturnRepository
}

func NewDeliveryRunReturnService() IDeliveryRunReturnService {
	return &DeliveryRunReturnService{
		opt:                         global.Setup.Common,
		RepositoryDeliveryRunReturn: repository.NewDeliveryRunReturnRepository(),
	}
}

func (s *DeliveryRunReturnService) Get(ctx context.Context, req dto.DeliveryRunReturnGetRequest) (res []*dto.DeliveryRunReturnResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DeliveryRunReturnService.Get")
	defer span.End()

	var deliveryRunSheetReturns []*model.DeliveryRunReturn
	deliveryRunSheetReturns, total, err = s.RepositoryDeliveryRunReturn.Get(ctx, dto.DeliveryRunReturnGetRequest{
		Offset:                     req.Offset,
		Limit:                      req.Limit,
		OrderBy:                    req.OrderBy,
		ArrDeliveryRunSheetItemIDs: req.ArrDeliveryRunSheetItemIDs,
	})
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, drr := range deliveryRunSheetReturns {
		res = append(res, &dto.DeliveryRunReturnResponse{
			ID:                     drr.ID,
			Code:                   drr.Code,
			TotalPrice:             drr.TotalPrice,
			TotalCharge:            drr.TotalCharge,
			CreatedAt:              drr.CreatedAt,
			DeliveryRunSheetItemID: drr.DeliveryRunSheetItemID,
		})
	}

	return
}

func (s *DeliveryRunReturnService) GetDetail(ctx context.Context, id int64, code string, deliveryRunSheetItemId int64) (res *dto.DeliveryRunReturnResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DeliveryRunReturnService.GetDetail")
	defer span.End()

	var deliveryRunReturn *model.DeliveryRunReturn
	deliveryRunReturn, err = s.RepositoryDeliveryRunReturn.GetByID(ctx, id, code, deliveryRunSheetItemId)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	res = &dto.DeliveryRunReturnResponse{
		ID:                     deliveryRunReturn.ID,
		Code:                   deliveryRunReturn.Code,
		TotalPrice:             deliveryRunReturn.TotalPrice,
		TotalCharge:            deliveryRunReturn.TotalCharge,
		CreatedAt:              deliveryRunReturn.CreatedAt,
		DeliveryRunSheetItemID: deliveryRunReturn.DeliveryRunSheetItemID,
	}

	return
}

func (s *DeliveryRunReturnService) Create(ctx context.Context, req *logisticService.CreateDeliveryRunReturnRequest) (res *dto.DeliveryRunReturnResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DeliveryRunReturnService.Create")
	defer span.End()

	model := &model.DeliveryRunReturn{
		TotalPrice:             req.Model.TotalPrice,
		TotalCharge:            req.Model.TotalCharge,
		CreatedAt:              req.Model.CreatedAt.AsTime(),
		DeliveryRunSheetItemID: req.Model.DeliveryRunSheetItemId,
	}

	var codeGenerator *configurationService.GetGenerateCodeResponse
	if codeGenerator, err = s.opt.Client.ConfigurationServiceGrpc.GetGenerateCode(ctx, &configurationService.GetGenerateCodeRequest{
		Format: "DRR",
		Domain: "delivery_run_return",
		Length: 6,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	model.Code = codeGenerator.Data.Code

	if err = s.RepositoryDeliveryRunReturn.Create(ctx, model); err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &dto.DeliveryRunReturnResponse{
		ID:                     model.ID,
		Code:                   codeGenerator.Data.Code,
		TotalPrice:             model.TotalPrice,
		TotalCharge:            model.TotalCharge,
		CreatedAt:              model.CreatedAt,
		DeliveryRunSheetItemID: model.DeliveryRunSheetItemID,
	}

	return
}

func (s *DeliveryRunReturnService) Update(ctx context.Context, req *logisticService.UpdateDeliveryRunReturnRequest) (res *dto.DeliveryRunReturnResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DeliveryRunReturnService.Create")
	defer span.End()

	var model = &model.DeliveryRunReturn{
		ID:          req.Id,
		TotalPrice:  req.TotalPrice,
		TotalCharge: req.TotalCharge,
	}

	if err = s.RepositoryDeliveryRunReturn.Update(ctx, model, "TotalPrice", "TotalCharge"); err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &dto.DeliveryRunReturnResponse{
		ID:          model.ID,
		TotalPrice:  model.TotalPrice,
		TotalCharge: model.TotalCharge,
	}

	return
}

func (s *DeliveryRunReturnService) Delete(ctx context.Context, id int64) (res *dto.DeliveryRunReturnResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DeliveryRunReturnService.Delete")
	defer span.End()

	var deliveryRunReturn *model.DeliveryRunReturn
	deliveryRunReturn, err = s.RepositoryDeliveryRunReturn.GetByID(ctx, id, "", 0)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	err = s.RepositoryDeliveryRunReturn.Delete(ctx, id)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &dto.DeliveryRunReturnResponse{
		ID:                     deliveryRunReturn.ID,
		Code:                   deliveryRunReturn.Code,
		TotalPrice:             deliveryRunReturn.TotalPrice,
		TotalCharge:            deliveryRunReturn.TotalCharge,
		CreatedAt:              deliveryRunReturn.CreatedAt,
		DeliveryRunSheetItemID: deliveryRunReturn.DeliveryRunSheetItemID,
	}

	return
}
