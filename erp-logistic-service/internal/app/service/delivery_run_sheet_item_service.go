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
	logisticService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/logistic_service"
)

type IDeliveryRunSheetItemService interface {
	Get(ctx context.Context, req *dto.DeliveryRunSheetItemGetRequest) (res []*dto.DeliveryRunSheetItemResponse, total int64, err error)
	GetDetail(ctx context.Context, id int64) (res *dto.DeliveryRunSheetItemResponse, err error)
	CreatePickup(ctx context.Context, req *logisticService.CreateDeliveryRunSheetItemRequest) (res *dto.DeliveryRunSheetItemResponse, err error)
	CreateDelivery(ctx context.Context, req *logisticService.CreateDeliveryRunSheetItemRequest) (res *dto.DeliveryRunSheetItemResponse, err error)
	Start(ctx context.Context, req *logisticService.StartDeliveryRunSheetItemRequest) (res *dto.DeliveryRunSheetItemResponse, err error)
	Postpone(ctx context.Context, req *logisticService.PostponeDeliveryRunSheetItemRequest) (res *dto.DeliveryRunSheetItemResponse, err error)
	FailPickup(ctx context.Context, req *logisticService.FailPickupDeliveryRunSheetItemRequest) (res *dto.DeliveryRunSheetItemResponse, err error)
	FailDelivery(ctx context.Context, req *logisticService.FailDeliveryDeliveryRunSheetItemRequest) (res *dto.DeliveryRunSheetItemResponse, err error)
	Success(ctx context.Context, req *logisticService.SuccessDeliveryRunSheetItemRequest) (res *dto.DeliveryRunSheetItemResponse, err error)
	Arrived(ctx context.Context, req *logisticService.ArrivedDeliveryRunSheetItemRequest) (res *dto.DeliveryRunSheetItemResponse, err error)
}

type DeliveryRunSheetItemService struct {
	opt                            opt.Options
	RepositoryDeliveryRunSheetItem repository.IDeliveryRunSheetItemRepository
}

func NewDeliveryRunSheetItemService() IDeliveryRunSheetItemService {
	return &DeliveryRunSheetItemService{
		opt:                            global.Setup.Common,
		RepositoryDeliveryRunSheetItem: repository.NewDeliveryRunSheetItemRepository(),
	}
}

func (s *DeliveryRunSheetItemService) Get(ctx context.Context, req *dto.DeliveryRunSheetItemGetRequest) (res []*dto.DeliveryRunSheetItemResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DeliveryRunSheetItemService.Get")
	defer span.End()

	var deliveryRunSheetItems []*model.DeliveryRunSheetItem
	deliveryRunSheetItems, total, err = s.RepositoryDeliveryRunSheetItem.Get(ctx, dto.DeliveryRunSheetItemGetRequest{
		Offset:               req.Offset,
		Limit:                req.Limit,
		OrderBy:              req.OrderBy,
		GroupBy:              req.GroupBy,
		StepType:             req.StepType,
		Status:               req.Status,
		DeliveryRunSheetIDs:  req.DeliveryRunSheetIDs,
		CourierIDs:           req.CourierIDs,
		ArrSalesOrderIDs:     req.ArrSalesOrderIDs,
		SearchSalesOrderCode: req.SearchSalesOrderCode,
	})
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, drsi := range deliveryRunSheetItems {
		res = append(res, &dto.DeliveryRunSheetItemResponse{
			ID:                          drsi.ID,
			StepType:                    drsi.StepType,
			Latitude:                    drsi.Latitude,
			Longitude:                   drsi.Longitude,
			Status:                      drsi.Status,
			Note:                        drsi.Note,
			RecipientName:               drsi.RecipientName,
			MoneyReceived:               drsi.MoneyReceived,
			DeliveryEvidenceImageURL:    drsi.DeliveryEvidenceImageURL,
			TransactionEvidenceImageURL: drsi.TransactionEvidenceImageURL,
			ArrivalTime:                 drsi.ArrivalTime,
			UnpunctualReason:            drsi.UnpunctualReason,
			UnpunctualDetail:            drsi.UnpunctualDetail,
			FarDeliveryReason:           drsi.FarDeliveryReason,
			CreatedAt:                   drsi.CreatedAt,
			StartedAt:                   drsi.StartedAt,
			FinishedAt:                  drsi.FinishedAt,
			DeliveryRunSheetID:          drsi.DeliveryRunSheetID,
			CourierID:                   drsi.CourierID,
			SalesOrderID:                drsi.SalesOrderID,
		})
	}

	return
}

func (s *DeliveryRunSheetItemService) GetDetail(ctx context.Context, id int64) (res *dto.DeliveryRunSheetItemResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DeliveryRunSheetItemService.GetDetail")
	defer span.End()

	var deliveryRunSheetItem *model.DeliveryRunSheetItem
	deliveryRunSheetItem, err = s.RepositoryDeliveryRunSheetItem.GetByID(ctx, id)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &dto.DeliveryRunSheetItemResponse{
		ID:                          deliveryRunSheetItem.ID,
		StepType:                    deliveryRunSheetItem.StepType,
		Latitude:                    deliveryRunSheetItem.Latitude,
		Longitude:                   deliveryRunSheetItem.Longitude,
		Status:                      deliveryRunSheetItem.Status,
		Note:                        deliveryRunSheetItem.Note,
		RecipientName:               deliveryRunSheetItem.RecipientName,
		MoneyReceived:               deliveryRunSheetItem.MoneyReceived,
		DeliveryEvidenceImageURL:    deliveryRunSheetItem.DeliveryEvidenceImageURL,
		TransactionEvidenceImageURL: deliveryRunSheetItem.TransactionEvidenceImageURL,
		ArrivalTime:                 deliveryRunSheetItem.ArrivalTime,
		UnpunctualReason:            deliveryRunSheetItem.UnpunctualReason,
		UnpunctualDetail:            deliveryRunSheetItem.UnpunctualDetail,
		FarDeliveryReason:           deliveryRunSheetItem.FarDeliveryReason,
		CreatedAt:                   deliveryRunSheetItem.CreatedAt,
		StartedAt:                   deliveryRunSheetItem.StartedAt,
		FinishedAt:                  deliveryRunSheetItem.FinishedAt,
		DeliveryRunSheetID:          deliveryRunSheetItem.DeliveryRunSheetID,
		CourierID:                   deliveryRunSheetItem.CourierID,
		SalesOrderID:                deliveryRunSheetItem.SalesOrderID,
	}

	return
}

func (s *DeliveryRunSheetItemService) CreatePickup(ctx context.Context, req *logisticService.CreateDeliveryRunSheetItemRequest) (res *dto.DeliveryRunSheetItemResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DeliveryRunSheetItemService.CreatePickup")
	defer span.End()

	model := &model.DeliveryRunSheetItem{
		StepType:           1,
		Latitude:           req.Model.Latitude,
		Longitude:          req.Model.Longitude,
		Status:             3,
		CreatedAt:          req.Model.CreatedAt.AsTime(),
		ArrivalTime:        req.Model.ArrivalTime.AsTime(),
		StartedAt:          req.Model.StartedAt.AsTime(),
		FinishedAt:         req.Model.FinishedAt.AsTime(),
		DeliveryRunSheetID: req.Model.DeliveryRunSheetId,
		CourierID:          req.Model.CourierId,
		SalesOrderID:       req.Model.SalesOrderId,
	}

	if err = s.RepositoryDeliveryRunSheetItem.Create(ctx, model); err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &dto.DeliveryRunSheetItemResponse{
		ID:                          model.ID,
		StepType:                    model.StepType,
		Latitude:                    model.Latitude,
		Longitude:                   model.Longitude,
		Status:                      model.Status,
		Note:                        model.Note,
		RecipientName:               model.RecipientName,
		MoneyReceived:               model.MoneyReceived,
		DeliveryEvidenceImageURL:    model.DeliveryEvidenceImageURL,
		TransactionEvidenceImageURL: model.TransactionEvidenceImageURL,
		ArrivalTime:                 model.ArrivalTime,
		UnpunctualReason:            model.UnpunctualReason,
		UnpunctualDetail:            model.UnpunctualDetail,
		FarDeliveryReason:           model.FarDeliveryReason,
		CreatedAt:                   model.CreatedAt,
		StartedAt:                   model.StartedAt,
		FinishedAt:                  model.FinishedAt,
		DeliveryRunSheetID:          model.DeliveryRunSheetID,
		CourierID:                   model.CourierID,
		SalesOrderID:                model.SalesOrderID,
	}

	return
}

func (s *DeliveryRunSheetItemService) CreateDelivery(ctx context.Context, req *logisticService.CreateDeliveryRunSheetItemRequest) (res *dto.DeliveryRunSheetItemResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DeliveryRunSheetItemService.CreateDelivery")
	defer span.End()

	model := &model.DeliveryRunSheetItem{
		StepType:           2,
		Latitude:           req.Model.Latitude,
		Longitude:          req.Model.Longitude,
		Status:             1,
		CreatedAt:          req.Model.CreatedAt.AsTime(),
		DeliveryRunSheetID: req.Model.DeliveryRunSheetId,
		CourierID:          req.Model.CourierId,
		SalesOrderID:       req.Model.SalesOrderId,
	}

	if err = s.RepositoryDeliveryRunSheetItem.Create(ctx, model); err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &dto.DeliveryRunSheetItemResponse{
		ID:                          model.ID,
		StepType:                    model.StepType,
		Latitude:                    model.Latitude,
		Longitude:                   model.Longitude,
		Status:                      model.Status,
		Note:                        model.Note,
		RecipientName:               model.RecipientName,
		MoneyReceived:               model.MoneyReceived,
		DeliveryEvidenceImageURL:    model.DeliveryEvidenceImageURL,
		TransactionEvidenceImageURL: model.TransactionEvidenceImageURL,
		ArrivalTime:                 model.ArrivalTime,
		UnpunctualReason:            model.UnpunctualReason,
		UnpunctualDetail:            model.UnpunctualDetail,
		FarDeliveryReason:           model.FarDeliveryReason,
		CreatedAt:                   model.CreatedAt,
		StartedAt:                   model.StartedAt,
		FinishedAt:                  model.FinishedAt,
		DeliveryRunSheetID:          model.DeliveryRunSheetID,
		CourierID:                   model.CourierID,
		SalesOrderID:                model.SalesOrderID,
	}

	return
}

func (s *DeliveryRunSheetItemService) Start(ctx context.Context, req *logisticService.StartDeliveryRunSheetItemRequest) (res *dto.DeliveryRunSheetItemResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DeliveryRunSheetItemService.Start")
	defer span.End()

	var model = &model.DeliveryRunSheetItem{
		ID:        req.Id,
		Status:    2,
		Latitude:  &req.Latitude,
		Longitude: &req.Longitude,
		StartedAt: time.Now(),
	}

	if err = s.RepositoryDeliveryRunSheetItem.Update(ctx, model, "Status", "Latitude", "Longitude", "StartedAt"); err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &dto.DeliveryRunSheetItemResponse{
		ID:        model.ID,
		Status:    model.Status,
		StartedAt: model.StartedAt,
	}

	return
}

func (s *DeliveryRunSheetItemService) Postpone(ctx context.Context, req *logisticService.PostponeDeliveryRunSheetItemRequest) (res *dto.DeliveryRunSheetItemResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DeliveryRunSheetItemService.Postpone")
	defer span.End()

	var model = &model.DeliveryRunSheetItem{
		ID:     req.Id,
		Status: 4,
		Note:   req.Note,
	}

	if err = s.RepositoryDeliveryRunSheetItem.Update(ctx, model, "Status", "Note"); err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &dto.DeliveryRunSheetItemResponse{
		ID:     model.ID,
		Status: model.Status,
		Note:   model.Note,
	}

	return
}

func (s *DeliveryRunSheetItemService) FailPickup(ctx context.Context, req *logisticService.FailPickupDeliveryRunSheetItemRequest) (res *dto.DeliveryRunSheetItemResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DeliveryRunSheetItemService.FailPickup")
	defer span.End()

	var model = &model.DeliveryRunSheetItem{
		ID:         req.Id,
		Status:     5,
		Note:       req.Note,
		FinishedAt: time.Now(),
	}

	if err = s.RepositoryDeliveryRunSheetItem.Update(ctx, model, "Status", "Note", "FinishedAt"); err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &dto.DeliveryRunSheetItemResponse{
		ID:         model.ID,
		Status:     model.Status,
		Note:       model.Note,
		FinishedAt: model.FinishedAt,
	}

	return
}

func (s *DeliveryRunSheetItemService) FailDelivery(ctx context.Context, req *logisticService.FailDeliveryDeliveryRunSheetItemRequest) (res *dto.DeliveryRunSheetItemResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DeliveryRunSheetItemService.FailDelivery")
	defer span.End()

	var model = &model.DeliveryRunSheetItem{
		ID:         req.Id,
		Status:     5,
		Latitude:   req.Latitude,
		Longitude:  req.Longitude,
		Note:       req.Note,
		FinishedAt: time.Now(),
	}

	if err = s.RepositoryDeliveryRunSheetItem.Update(ctx, model, "Status", "Latitude", "Longitude", "Note", "FinishedAt"); err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &dto.DeliveryRunSheetItemResponse{
		ID:         model.ID,
		Status:     model.Status,
		Latitude:   model.Latitude,
		Longitude:  model.Longitude,
		Note:       model.Note,
		FinishedAt: model.FinishedAt,
	}

	return
}

func (s *DeliveryRunSheetItemService) Success(ctx context.Context, req *logisticService.SuccessDeliveryRunSheetItemRequest) (res *dto.DeliveryRunSheetItemResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DeliveryRunSheetItemService.Success")
	defer span.End()

	var model = &model.DeliveryRunSheetItem{
		ID:                          req.Id,
		Latitude:                    req.Latitude,
		Longitude:                   req.Longitude,
		Status:                      3,
		Note:                        req.Note,
		RecipientName:               req.RecipientName,
		MoneyReceived:               req.MoneyReceived,
		DeliveryEvidenceImageURL:    req.DeliveryEvidenceImageUrl,
		TransactionEvidenceImageURL: req.TransactionEvidenceImageUrl,
		UnpunctualReason:            int8(req.UnpunctualReason),
		UnpunctualDetail:            int8(req.UnpunctualDetail),
		FarDeliveryReason:           req.FarDeliveryReason,
		FinishedAt:                  time.Now(),
	}

	if err = s.RepositoryDeliveryRunSheetItem.Update(ctx, model, "Latitude", "Longitude", "Status", "Note", "RecipientName", "MoneyReceived", "DeliveryEvidenceImageURL", "TransactionEvidenceImageURL", "UnpunctualReason", "UnpunctualDetail", "FarDeliveryReason", "FinishedAt"); err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &dto.DeliveryRunSheetItemResponse{
		ID:                          model.ID,
		Latitude:                    model.Latitude,
		Longitude:                   model.Longitude,
		Status:                      model.Status,
		Note:                        model.Note,
		RecipientName:               model.RecipientName,
		MoneyReceived:               model.MoneyReceived,
		DeliveryEvidenceImageURL:    model.DeliveryEvidenceImageURL,
		TransactionEvidenceImageURL: model.TransactionEvidenceImageURL,
		UnpunctualReason:            model.UnpunctualReason,
		UnpunctualDetail:            model.UnpunctualDetail,
		FarDeliveryReason:           model.FarDeliveryReason,
		FinishedAt:                  model.FinishedAt,
	}

	return
}

func (s *DeliveryRunSheetItemService) Arrived(ctx context.Context, req *logisticService.ArrivedDeliveryRunSheetItemRequest) (res *dto.DeliveryRunSheetItemResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DeliveryRunSheetItemService.Arrived")
	defer span.End()

	var model = &model.DeliveryRunSheetItem{
		ID:          req.Id,
		ArrivalTime: time.Now(),
	}

	if err = s.RepositoryDeliveryRunSheetItem.Update(ctx, model, "ArrivalTime"); err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &dto.DeliveryRunSheetItemResponse{
		ID:          model.ID,
		ArrivalTime: model.ArrivalTime,
	}

	return
}
