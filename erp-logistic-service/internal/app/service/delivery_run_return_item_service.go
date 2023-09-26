package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-logistic-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/internal/app/repository"

	logisticService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/logistic_service"
)

type IDeliveryRunReturnItemService interface {
	Get(ctx context.Context, req dto.DeliveryRunReturnItemGetRequest) (res []*dto.DeliveryRunReturnItemResponse, total int64, err error)
	GetDetail(ctx context.Context, id int64, deliveryOrderItemId string) (res *dto.DeliveryRunReturnItemResponse, err error)
	Create(ctx context.Context, req *logisticService.CreateDeliveryRunReturnItemRequest) (res *dto.DeliveryRunReturnItemResponse, err error)
	Update(ctx context.Context, req *logisticService.UpdateDeliveryRunReturnItemRequest) (res *dto.DeliveryRunReturnItemResponse, err error)
	Delete(ctx context.Context, id int64) (res *dto.DeliveryRunReturnItemResponse, err error)
}

type DeliveryRunReturnItemService struct {
	opt                             opt.Options
	RepositoryDeliveryRunReturnItem repository.IDeliveryRunReturnItemRepository
}

func NewDeliveryRunReturnItemService() IDeliveryRunReturnItemService {
	return &DeliveryRunReturnItemService{
		opt:                             global.Setup.Common,
		RepositoryDeliveryRunReturnItem: repository.NewDeliveryRunReturnItemRepository(),
	}
}

func (s *DeliveryRunReturnItemService) Get(ctx context.Context, req dto.DeliveryRunReturnItemGetRequest) (res []*dto.DeliveryRunReturnItemResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DeliveryRunReturnItemService.Get")
	defer span.End()

	var deliveryRunSheetReturnItems []*model.DeliveryRunReturnItem
	deliveryRunSheetReturnItems, total, err = s.RepositoryDeliveryRunReturnItem.Get(ctx, dto.DeliveryRunReturnItemGetRequest{
		Offset:                  req.Offset,
		Limit:                   req.Limit,
		OrderBy:                 req.OrderBy,
		ArrDeliveryRunReturnIDs: req.ArrDeliveryRunReturnIDs,
		ArrDeliveryOrderItemIDs: req.ArrDeliveryOrderItemIDs,
	})
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, drri := range deliveryRunSheetReturnItems {
		res = append(res, &dto.DeliveryRunReturnItemResponse{
			ID:                  drri.ID,
			ReceiveQty:          drri.ReceiveQty,
			ReturnReason:        drri.ReturnReason,
			ReturnEvidence:      drri.ReturnEvidence,
			Subtotal:            drri.Subtotal,
			DeliveryRunReturnID: drri.DeliveryRunReturnID,
			DeliveryOrderItemID: drri.DeliveryOrderItemID,
		})
	}

	return
}

func (s *DeliveryRunReturnItemService) GetDetail(ctx context.Context, id int64, deliveryOrderItemId string) (res *dto.DeliveryRunReturnItemResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DeliveryRunReturnItemService.GetDetail")
	defer span.End()

	var deliveryRunReturnItem *model.DeliveryRunReturnItem
	deliveryRunReturnItem, err = s.RepositoryDeliveryRunReturnItem.GetByID(ctx, id, deliveryOrderItemId)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &dto.DeliveryRunReturnItemResponse{
		ID:                  deliveryRunReturnItem.ID,
		ReceiveQty:          deliveryRunReturnItem.ReceiveQty,
		ReturnReason:        deliveryRunReturnItem.ReturnReason,
		ReturnEvidence:      deliveryRunReturnItem.ReturnEvidence,
		Subtotal:            deliveryRunReturnItem.Subtotal,
		DeliveryRunReturnID: deliveryRunReturnItem.DeliveryRunReturnID,
		DeliveryOrderItemID: deliveryRunReturnItem.DeliveryOrderItemID,
	}

	return
}

func (s *DeliveryRunReturnItemService) Create(ctx context.Context, req *logisticService.CreateDeliveryRunReturnItemRequest) (res *dto.DeliveryRunReturnItemResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DeliveryRunReturnItemService.Create")
	defer span.End()

	model := &model.DeliveryRunReturnItem{
		ReceiveQty:          req.Model.ReceiveQty,
		ReturnReason:        int8(req.Model.ReturnReason),
		ReturnEvidence:      req.Model.ReturnEvidence,
		Subtotal:            req.Model.Subtotal,
		DeliveryRunReturnID: req.Model.DeliveryRunReturnId,
		DeliveryOrderItemID: req.Model.DeliveryOrderItemId,
	}

	if err = s.RepositoryDeliveryRunReturnItem.Create(ctx, model); err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &dto.DeliveryRunReturnItemResponse{
		ID:                  model.ID,
		ReceiveQty:          model.ReceiveQty,
		ReturnReason:        int8(model.ReturnReason),
		ReturnEvidence:      model.ReturnEvidence,
		Subtotal:            model.Subtotal,
		DeliveryRunReturnID: model.DeliveryRunReturnID,
		DeliveryOrderItemID: model.DeliveryOrderItemID,
	}

	return
}

func (s *DeliveryRunReturnItemService) Update(ctx context.Context, req *logisticService.UpdateDeliveryRunReturnItemRequest) (res *dto.DeliveryRunReturnItemResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DeliveryRunReturnItemService.Create")
	defer span.End()

	var model = &model.DeliveryRunReturnItem{
		ID:             req.Id,
		ReceiveQty:     req.ReceiveQty,
		ReturnReason:   int8(req.ReturnReason),
		ReturnEvidence: req.ReturnEvidence,
		Subtotal:       req.Subtotal,
	}

	if err = s.RepositoryDeliveryRunReturnItem.Update(ctx, model, "ReceiveQty", "ReturnReason", "ReturnEvidence", "Subtotal"); err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &dto.DeliveryRunReturnItemResponse{
		ID:             model.ID,
		ReceiveQty:     model.ReceiveQty,
		ReturnReason:   int8(model.ReturnReason),
		ReturnEvidence: model.ReturnEvidence,
		Subtotal:       model.Subtotal,
	}

	return
}

func (s *DeliveryRunReturnItemService) Delete(ctx context.Context, id int64) (res *dto.DeliveryRunReturnItemResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DeliveryRunReturnItemService.Delete")
	defer span.End()

	var deliveryRunReturnItem *model.DeliveryRunReturnItem
	deliveryRunReturnItem, err = s.RepositoryDeliveryRunReturnItem.GetByID(ctx, id, "")
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	err = s.RepositoryDeliveryRunReturnItem.Delete(ctx, id)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &dto.DeliveryRunReturnItemResponse{
		ID:                  deliveryRunReturnItem.ID,
		ReceiveQty:          deliveryRunReturnItem.ReceiveQty,
		ReturnReason:        int8(deliveryRunReturnItem.ReturnReason),
		ReturnEvidence:      deliveryRunReturnItem.ReturnEvidence,
		Subtotal:            deliveryRunReturnItem.Subtotal,
		DeliveryRunReturnID: deliveryRunReturnItem.DeliveryRunReturnID,
		DeliveryOrderItemID: deliveryRunReturnItem.DeliveryOrderItemID,
	}

	return
}
