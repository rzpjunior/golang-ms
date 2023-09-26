package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *BridgeGrpcHandler) GetDeliveryFeeList(ctx context.Context, req *bridgeService.GetDeliveryFeeListRequest) (res *bridgeService.GetDeliveryFeeListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetDeliveryFeeList")
	defer span.End()

	var DeliveryFeees []dto.DeliveryFeeResponse
	DeliveryFeees, _, err = h.ServicesDeliveryFee.Get(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*bridgeService.DeliveryFee
	for _, DeliveryFee := range DeliveryFeees {
		data = append(data, &bridgeService.DeliveryFee{
			Id:             DeliveryFee.ID,
			Code:           DeliveryFee.Code,
			Name:           DeliveryFee.Name,
			Note:           DeliveryFee.Note,
			Status:         DeliveryFee.Status,
			MinimumOrder:   DeliveryFee.MinimumOrder,
			DeliveryFee:    DeliveryFee.DeliveryFee,
			RegionId:       DeliveryFee.RegionId,
			CustomerTypeId: DeliveryFee.CutomerTypeId,
		})
	}

	res = &bridgeService.GetDeliveryFeeListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *BridgeGrpcHandler) GetDeliveryFeeDetail(ctx context.Context, req *bridgeService.GetDeliveryFeeDetailRequest) (res *bridgeService.GetDeliveryFeeDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetDeliveryFeeDetail")
	defer span.End()

	var DeliveryFee dto.DeliveryFeeResponse
	DeliveryFee, err = h.ServicesDeliveryFee.GetDetail(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.GetDeliveryFeeDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &bridgeService.DeliveryFee{
			Id:             DeliveryFee.ID,
			Code:           DeliveryFee.Code,
			Name:           DeliveryFee.Name,
			Note:           DeliveryFee.Note,
			Status:         DeliveryFee.Status,
			MinimumOrder:   DeliveryFee.MinimumOrder,
			DeliveryFee:    DeliveryFee.DeliveryFee,
			RegionId:       DeliveryFee.RegionId,
			CustomerTypeId: DeliveryFee.CutomerTypeId,
		},
	}
	return
}

func (h *BridgeGrpcHandler) GetDeliveryFeeGPList(ctx context.Context, req *bridgeService.GetDeliveryFeeGPListRequest) (res *bridgeService.GetDeliveryFeeGPListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetCustomerTypeGPList")
	defer span.End()

	res, err = h.ServicesDeliveryFee.GetGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	return
}
