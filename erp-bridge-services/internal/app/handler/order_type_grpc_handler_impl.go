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

func (h *BridgeGrpcHandler) GetOrderTypeList(ctx context.Context, req *bridgeService.GetOrderTypeListRequest) (res *bridgeService.GetOrderTypeListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetOrderTypeList")
	defer span.End()

	var orderTypes []dto.OrderTypeResponse
	orderTypes, _, err = h.ServicesOrderType.Get(ctx, int(req.Offset), int(req.Limit), int(req.Status), req.Search, req.OrderBy)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*bridgeService.OrderType
	for _, orderType := range orderTypes {
		data = append(data, &bridgeService.OrderType{
			Id:          orderType.ID,
			Code:        orderType.Code,
			Description: orderType.Description,
			Status:      int32(orderType.Status),
			CreatedAt:   timestamppb.New(orderType.CreatedAt),
			UpdatedAt:   timestamppb.New(orderType.UpdatedAt),
		})
	}

	res = &bridgeService.GetOrderTypeListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *BridgeGrpcHandler) GetOrderTypeDetail(ctx context.Context, req *bridgeService.GetOrderTypeDetailRequest) (res *bridgeService.GetOrderTypeDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetOrderTypeDetail")
	defer span.End()

	var orderType dto.OrderTypeResponse
	orderType, err = h.ServicesOrderType.GetDetail(ctx, req.Id, req.Code)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.GetOrderTypeDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &bridgeService.OrderType{
			Id:          orderType.ID,
			Code:        orderType.Code,
			Description: orderType.Description,
			Status:      int32(orderType.Status),
			CreatedAt:   timestamppb.New(orderType.CreatedAt),
			UpdatedAt:   timestamppb.New(orderType.UpdatedAt),
		},
	}
	return
}

func (h *BridgeGrpcHandler) GetOrderTypeGPList(ctx context.Context, req *bridgeService.GetOrderTypeGPListRequest) (res *bridgeService.GetOrderTypeGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetOrderTypeGPList")
	defer span.End()

	res, err = h.ServicesOrderType.GetGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	return
}

func (h *BridgeGrpcHandler) GetOrderTypeGPDetail(ctx context.Context, req *bridgeService.GetOrderTypeGPDetailRequest) (res *bridgeService.GetOrderTypeGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetOrderTypeGPDetail")
	defer span.End()

	res, err = h.ServicesOrderType.GetDetailGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	return
}
