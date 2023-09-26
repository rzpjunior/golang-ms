package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *BridgeGrpcHandler) GetSalesOrderList(ctx context.Context, req *bridgeService.GetSalesOrderListRequest) (res *bridgeService.GetSalesOrderListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetSalesOrderList")
	defer span.End()

	var salesOrders []dto.SalesOrderResponse
	salesOrders, _, err = h.ServicesSalesOrder.Get(ctx, int(req.Offset), int(req.Limit), int(req.Status), req.Search, req.OrderBy, req.AddressId, req.CustomerId, req.SalespersonId, req.OrderDateFrom.AsTime(), req.OrderDateTo.AsTime())
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*bridgeService.SalesOrder
	for _, salesOrder := range salesOrders {
		data = append(data, &bridgeService.SalesOrder{
			Id:            salesOrder.ID,
			Code:          salesOrder.Code,
			DocNumber:     salesOrder.DocNumber,
			AddressId:     "ADVANCED0001",
			CustomerId:    "ADVANCED0001",
			SalespersonId: "GREG E.",
			WrtId:         "",
			Application:   int32(salesOrder.Application),
			Status:        int32(salesOrder.Status),
			OrderTypeId:   salesOrder.OrderTypeID,
			OrderDate:     timestamppb.New(salesOrder.OrderDate),
			Total:         salesOrder.Total,
			CreatedDate:   timestamppb.New(salesOrder.CreatedDate),
			ModifiedDate:  timestamppb.New(salesOrder.ModifiedDate),
			FinishedDate:  timestamppb.New(salesOrder.FinishedDate),
			CreatedAt:     timestamppb.New(salesOrder.CreatedAt),
			UpdatedAt:     timestamppb.New(salesOrder.UpdatedAt),
			SiteId:        "NORTH",
		})
	}

	res = &bridgeService.GetSalesOrderListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *BridgeGrpcHandler) GetSalesOrderDetail(ctx context.Context, req *bridgeService.GetSalesOrderDetailRequest) (res *bridgeService.GetSalesOrderDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetSalesOrderDetail")
	defer span.End()

	var salesOrder dto.SalesOrderResponse
	salesOrder, err = h.ServicesSalesOrder.GetDetail(ctx, req.Id, req.Code)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.GetSalesOrderDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &bridgeService.SalesOrder{
			Id:            salesOrder.ID,
			Code:          salesOrder.Code,
			DocNumber:     salesOrder.DocNumber,
			AddressId:     "ADVANCED0001",
			CustomerId:    "ADVANCED0001",
			SalespersonId: "GREG E.",
			WrtId:         "",
			Application:   int32(salesOrder.Application),
			Status:        int32(salesOrder.Status),
			OrderTypeId:   salesOrder.OrderTypeID,
			OrderDate:     timestamppb.New(salesOrder.OrderDate),
			Total:         salesOrder.Total,
			CreatedDate:   timestamppb.New(salesOrder.CreatedDate),
			ModifiedDate:  timestamppb.New(salesOrder.ModifiedDate),
			FinishedDate:  timestamppb.New(salesOrder.FinishedDate),
			CreatedAt:     timestamppb.New(salesOrder.CreatedAt),
			UpdatedAt:     timestamppb.New(salesOrder.UpdatedAt),
			SiteId:        "NORTH",
		},
	}
	return
}

func (h *BridgeGrpcHandler) CreateSalesOrder(ctx context.Context, req *bridgeService.CreateSalesOrderRequest) (res *bridgeService.CreateSalesOrderResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.CreateSalesOrder")
	defer span.End()

	var salesOrder dto.SalesOrderResponse
	salesOrder, err = h.ServicesSalesOrder.CreateSalesOrder(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.CreateSalesOrderResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &bridgeService.SalesOrder{
			Id:            salesOrder.ID,
			Code:          salesOrder.Code,
			DocNumber:     salesOrder.DocNumber,
			AddressId:     "ADVANCED0001",
			CustomerId:    "ADVANCED0001",
			SalespersonId: "GREG E.",
			WrtId:         "",
			Application:   int32(salesOrder.Application),
			Status:        int32(salesOrder.Status),
			OrderTypeId:   salesOrder.OrderTypeID,
			OrderDate:     timestamppb.New(salesOrder.OrderDate),
			Total:         salesOrder.Total,
			CreatedDate:   timestamppb.New(salesOrder.CreatedDate),
			ModifiedDate:  timestamppb.New(salesOrder.ModifiedDate),
			FinishedDate:  timestamppb.New(salesOrder.FinishedDate),
			CreatedAt:     timestamppb.New(salesOrder.CreatedAt),
			UpdatedAt:     timestamppb.New(salesOrder.UpdatedAt),
			SiteId:        "NORTH",
		},
	}
	return
}

func (h *BridgeGrpcHandler) GetSalesOrderListGPAll(ctx context.Context, req *bridgeService.GetSalesOrderGPListRequest) (res *bridgeService.GetSalesOrderGPListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetSalesOrderListGPAll")
	defer span.End()

	res, err = h.ServicesSalesOrder.GetGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	return
}

func (h *BridgeGrpcHandler) GetSalesOrderListGPByID(ctx context.Context, req *bridgeService.GetSalesOrderGPListByIDRequest) (res *bridgeService.GetSalesOrderGPListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetSalesOrderListGPByID")
	defer span.End()

	res, err = h.ServicesSalesOrder.GetGPByID(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	return
}

func (h *BridgeGrpcHandler) CreateSalesOrderGP(ctx context.Context, req *bridgeService.CreateSalesOrderGPRequest) (res *bridgeService.CreateSalesOrderGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.CreateSalesOrderGP")
	defer span.End()
	req.Interid = global.EnvDatabaseGP

	response, err := h.ServicesSalesOrder.CreateSalesOrderGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.CreateSalesOrderGPResponse{
		Code:     int32(response.Code),
		Message:  response.Message,
		Sopnumbe: response.Sopnumbe,
	}
	return
}

func (h *BridgeGrpcHandler) GetSalesMovementGP(ctx context.Context, req *bridgeService.GetSalesMovementGPRequest) (res *bridgeService.GetSalesMovementGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetSalesMovementGP")
	defer span.End()
	req.Interid = global.EnvDatabaseGP

	response, err := h.ServicesSalesOrder.GetSalesMovementGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*bridgeService.SalesMovement
	for _, salesOrder := range response.Data {
		data = append(data, &bridgeService.SalesMovement{
			SoNumber:      salesOrder.SoNumber,
			Docdate:       salesOrder.Docdate,
			SoStatus:      salesOrder.SoStatus,
			Picking:       salesOrder.Picking,
			Checking:      salesOrder.Checking,
			DeliveryOrder: salesOrder.DeliveryOrder,
			SiNumber:      salesOrder.SiNumber,
			SiStatus:      salesOrder.SiStatus,
			CashReceipt:   salesOrder.CashReceipt,
			SalesReturn:   salesOrder.SalesReturn,
		})
	}

	res = &bridgeService.GetSalesMovementGPResponse{
		Code:    int32(response.Code),
		Message: response.Message,
		Data:    data,
		// Sopnumbe: response.Sopnumbe,
	}
	
	return
}
