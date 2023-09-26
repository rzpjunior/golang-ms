package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *BridgeGrpcHandler) GetDeliveryOrderListGP(ctx context.Context, req *bridgeService.GetDeliveryOrderGPListRequest) (res *bridgeService.GetDeliveryOrderGPListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetPurchasePlanGPList")
	defer span.End()

	res, err = h.ServiceDeliveryOrder.GetListGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *BridgeGrpcHandler) GetDeliveryOrderDetail(ctx context.Context, req *bridgeService.GetDeliveryOrderDetailRequest) (res *bridgeService.GetDeliveryOrderDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetDeliveryOrderDetail")
	defer span.End()

	var DeliveryOrder dto.DeliveryOrderResponse
	DeliveryOrder, err = h.ServiceDeliveryOrder.GetDetail(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.GetDeliveryOrderDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &bridgeService.DeliveryOrder{
			Id:        DeliveryOrder.ID,
			WrtId:     DeliveryOrder.WrtID,
			Status:    int32(DeliveryOrder.Status),
			CreatedAt: timestamppb.New(DeliveryOrder.CreatedDate),
			SiteId:    DeliveryOrder.SiteID,
			//SiteId:        DeliveryOrder.SiteID,
		},
	}
	return
}

func (h *BridgeGrpcHandler) CreateDeliveryOrder(ctx context.Context, req *bridgeService.CreateDeliveryOrderRequest) (res *bridgeService.CreateDeliveryOrderResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetDeliveryOrderDetail")
	defer span.End()

	res, err = h.ServiceDeliveryOrder.Create(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}
