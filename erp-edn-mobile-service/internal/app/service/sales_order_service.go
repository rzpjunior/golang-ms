package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/global"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/internal/app/dto"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

func NewServiceSalesOrder() ISalesOrderService {
	m := new(SalesOrderService)
	m.opt = global.Setup.Common
	return m
}

type ISalesOrderService interface {
	Create(ctx context.Context, req *dto.CreateSalesOrderRequest) (res *dto.SalesOrderResponse, err error)
}

type SalesOrderService struct {
	opt opt.Options
}

func (s *SalesOrderService) Create(ctx context.Context, req *dto.CreateSalesOrderRequest) (res *dto.SalesOrderResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesOrderService.Create")
	defer span.End()

	// get SalesOrder from bridge
	var (
		soRes *bridgeService.CreateSalesOrderResponse
		soi   []*bridgeService.SalesOrderItem
	)

	for _, item := range req.Items {
		soi = append(soi, &bridgeService.SalesOrderItem{
			ItemId:    item.ItemID,
			OrderQty:  item.Quantity,
			UnitPrice: item.UnitPrice,
			Subtotal:  item.UnitPrice * item.Quantity,
		})
	}

	soRes, err = s.opt.Client.BridgeServiceGrpc.CreateSalesOrder(ctx, &bridgeService.CreateSalesOrderRequest{
		Data: &bridgeService.SalesOrder{
			AddressId:     req.AddressID,
			CustomerId:    req.CustomerID,
			SalespersonId: req.SalespersonID,
			WrtId:         req.WrtID,
			SiteId:        req.SiteID,
		},
		Dataitem: soi,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "sales invoice")
		return
	}

	res = &dto.SalesOrderResponse{
		ID:            "SO1",
		Code:          "ABCD123",
		DocNumber:     soRes.Data.DocNumber,
		AddressID:     req.AddressID,
		CustomerID:    req.CustomerID,
		SalespersonID: req.SalespersonID,
		WrtID:         req.WrtID,
		OrderTypeID:   req.OrderTypeID,
		SiteID:        req.SiteID,
		Application:   int8(soRes.Data.Application),
		Status:        11,
		StatusConvert: statusx.ConvertStatusValue(11),
		Total:         soRes.Data.Total,
		CreatedAt:     soRes.Data.CreatedAt.AsTime(),
		CreatedDate:   soRes.Data.CreatedDate.AsTime(),
		UpdatedAt:     soRes.Data.UpdatedAt.AsTime(),
		ModifiedDate:  soRes.Data.ModifiedDate.AsTime(),
	}

	return
}
