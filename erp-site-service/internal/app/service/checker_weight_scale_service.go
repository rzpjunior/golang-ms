package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/catalog_service"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-site-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/internal/app/repository"
)

type ICheckerWeightScaleService interface {
	Get(ctx context.Context, req *dto.CheckerWeightScaleGetRequest) (res *dto.CheckerWeightScaleGetResponse, err error)
	Update(ctx context.Context, req *dto.CheckerWeightScaleUpdateRequest) (res *dto.CheckerSuccessResponse, err error)
}

type CheckerWeightScaleService struct {
	opt                          opt.Options
	RepositoryPickingOrderItem   repository.IPickingOrderItemRepository
	RepositoryPickingOrderAssign repository.IPickingOrderAssignRepository
}

func NewCheckerWeightScaleService() ICheckerWeightScaleService {
	return &CheckerWeightScaleService{
		opt:                          global.Setup.Common,
		RepositoryPickingOrderItem:   repository.NewPickingOrderItemRepository(),
		RepositoryPickingOrderAssign: repository.NewPickingOrderAssignRepository(),
	}
}

func (s *CheckerWeightScaleService) Get(ctx context.Context, req *dto.CheckerWeightScaleGetRequest) (res *dto.CheckerWeightScaleGetResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CheckerWeightScaleService.Get")
	defer span.End()

	var pickingOrderItem *model.PickingOrderItem
	if pickingOrderItem, err = s.RepositoryPickingOrderItem.GetByID(ctx, &dto.PickingOrderItemGetDetailRequest{
		Id: req.PickingOrderItemId,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order item")
		return
	}

	var productDetail *bridgeService.GetItemGPResponse
	if productDetail, err = s.opt.Client.BridgeServiceGrpc.GetItemGPDetail(ctx, &bridgeService.GetItemGPDetailRequest{
		Id: pickingOrderItem.ItemNumber,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("inventory", "product detail")
		return
	}

	var productImage *catalog_service.GetItemDetailResponse
	if productImage, err = s.opt.Client.CatalogServiceGrpc.GetItemDetail(ctx, &catalog_service.GetItemDetailRequest{
		Id: pickingOrderItem.ItemNumber,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("inventory", "product image")
		return
	}

	// TODO : UNCOMMENT DUMMY
	res = &dto.CheckerWeightScaleGetResponse{
		PickingOrderItemId: pickingOrderItem.Id,
		ProductPicture:     "", // filled below if exist
		ProductName:        productDetail.Data[0].Itmgedsc,
		ProductId:          pickingOrderItem.ItemNumber,
		OrderQty:           pickingOrderItem.OrderQuantity,
		OrderMinQty:        1,
		// OrderMinQty:        productDetail.Data[0].Minorqty,
	}

	if len(productImage.Data.ItemImage) > 0 {
		res.ProductPicture = productImage.Data.ItemImage[0].ImageUrl
	}

	return
}

func (s *CheckerWeightScaleService) Update(ctx context.Context, req *dto.CheckerWeightScaleUpdateRequest) (res *dto.CheckerSuccessResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CheckerWeightScaleService.Update")
	defer span.End()

	var pickingOrderItem *model.PickingOrderItem
	if pickingOrderItem, err = s.RepositoryPickingOrderItem.GetByID(ctx, &dto.PickingOrderItemGetDetailRequest{
		Id: req.PickingOrderItemId,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order item")
		return
	}

	// check picking order assign if status on checking
	var pickingOrderAssign *model.PickingOrderAssign
	if pickingOrderAssign, err = s.RepositoryPickingOrderAssign.GetByID(ctx, pickingOrderItem.PickingOrderAssignId, ""); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order assign")
		return
	}

	if pickingOrderAssign.Status != 20 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorMustStatus("picking order assign", "checking")
		return
	}

	pickingOrderItem.CheckQuantity = req.CheckQty

	if err = s.RepositoryPickingOrderItem.Update(ctx, pickingOrderItem, "CheckQuantity"); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order item")
		return
	}

	res = &dto.CheckerSuccessResponse{
		Success: true,
	}

	return
}
