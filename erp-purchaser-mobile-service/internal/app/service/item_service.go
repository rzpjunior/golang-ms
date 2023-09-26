package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/global"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	dto "git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/internal/app/dto"
)

type IItemService interface {
	Get(ctx context.Context, req dto.ItemListRequest) (res []*dto.ItemResponse, err error)
	GetByID(ctx context.Context, req dto.ItemDetailRequest) (res *dto.ItemResponse, err error)
}

type ItemService struct {
	opt opt.Options
}

func NewItemService() IItemService {
	return &ItemService{
		opt: global.Setup.Common,
	}
}

func (s *ItemService) Get(ctx context.Context, req dto.ItemListRequest) (res []*dto.ItemResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemService.Get")
	defer span.End()

	// get adddress from bridge
	var itemResponse *bridgeService.GetItemGPResponse
	itemResponse, err = s.opt.Client.BridgeServiceGrpc.GetItemGPList(ctx, &bridgeService.GetItemGPListRequest{
		Limit:  req.Limit,
		Offset: req.Offset,
		// Status:  req.Status,
		// Search:  req.Search,
		// OrderBy: req.OrderBy,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "address")
		return
	}

	for _, item := range itemResponse.Data {

		var unitPrice float64
		unitPrice = item.Currcost

		var uom *bridgeService.GetUomGPResponse
		uom, err = s.opt.Client.BridgeServiceGrpc.GetUomGPDetail(ctx, &bridgeService.GetUomGPDetailRequest{
			Id: item.Uomschdl,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "uom")
			return
		}

		res = append(res, &dto.ItemResponse{
			ID:   item.Itemnmbr,
			Code: item.Itemnmbr,
			Uom: &dto.UomResponse{
				ID:   uom.Data[0].Uomschdl,
				Code: uom.Data[0].Uomschdl,
				Name: uom.Data[0].Umschdsc,
			},
			Description:          item.Itemdesc,
			UnitWeightConversion: item.GnlWeighttolerance,
			OrderMinQty:          item.Minorqty,
			OrderMaxQty:          item.Maxordqty,
			ItemType:             item.ItemTypeDesc,
			UnitPrice:            unitPrice,
		})
	}

	return
}

func (s *ItemService) GetByID(ctx context.Context, req dto.ItemDetailRequest) (res *dto.ItemResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemService.GetByID")
	defer span.End()

	// get item from bridge
	var item *bridgeService.GetItemGPResponse
	item, err = s.opt.Client.BridgeServiceGrpc.GetItemGPDetail(ctx, &bridgeService.GetItemGPDetailRequest{
		Id: req.Id,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("catalog", "item")
		return
	}

	var unitPrice float64
	unitPrice = item.Data[0].Currcost

	var uom *bridgeService.GetUomGPResponse
	uom, err = s.opt.Client.BridgeServiceGrpc.GetUomGPDetail(ctx, &bridgeService.GetUomGPDetailRequest{
		Id: item.Data[0].Uomschdl,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "uom")
		return
	}

	res = &dto.ItemResponse{
		ID:   item.Data[0].Itemnmbr,
		Code: item.Data[0].Itemnmbr,
		Uom: &dto.UomResponse{
			ID:   uom.Data[0].Uomschdl,
			Code: uom.Data[0].Uomschdl,
			Name: uom.Data[0].Umschdsc,
		},
		Description:          item.Data[0].Itemdesc,
		UnitWeightConversion: item.Data[0].GnlWeighttolerance,
		OrderMinQty:          item.Data[0].Minorqty,
		OrderMaxQty:          item.Data[0].Maxordqty,
		ItemType:             item.Data[0].ItemTypeDesc,
		UnitPrice:            unitPrice,
	}

	return
}
