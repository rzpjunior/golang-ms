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

func (h *BridgeGrpcHandler) GetUomList(ctx context.Context, req *bridgeService.GetUomListRequest) (res *bridgeService.GetUomListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetUomList")
	defer span.End()

	var uoms []dto.UomResponse
	uoms, _, err = h.ServicesUom.Get(ctx, int(req.Offset), int(req.Limit), int(req.Status), req.Search, req.OrderBy)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*bridgeService.Uom
	for _, uom := range uoms {
		data = append(data, &bridgeService.Uom{
			Id:             uom.ID,
			Code:           uom.Code,
			Description:    uom.Description,
			Status:         int32(uom.Status),
			DecimalEnabled: int32(uom.DecimalEnabled),
			CreatedAt:      timestamppb.New(uom.CreatedAt),
			UpdatedAt:      timestamppb.New(uom.UpdatedAt),
		})
	}

	res = &bridgeService.GetUomListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *BridgeGrpcHandler) GetUomDetail(ctx context.Context, req *bridgeService.GetUomDetailRequest) (res *bridgeService.GetUomDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetUomDetail")
	defer span.End()

	var uom dto.UomResponse
	uom, err = h.ServicesUom.GetDetail(ctx, req.Id, req.Code)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.GetUomDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &bridgeService.Uom{
			Id:             uom.ID,
			Code:           uom.Code,
			Description:    uom.Description,
			Status:         int32(uom.Status),
			DecimalEnabled: int32(uom.DecimalEnabled),
			CreatedAt:      timestamppb.New(uom.CreatedAt),
			UpdatedAt:      timestamppb.New(uom.UpdatedAt),
		},
	}
	return
}

func (h *BridgeGrpcHandler) GetUomGPList(ctx context.Context, req *bridgeService.GetUomGPListRequest) (res *bridgeService.GetUomGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetUOMGPList")
	defer span.End()

	res, err = h.ServicesUom.GetGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	return
}

func (h *BridgeGrpcHandler) GetUomGPDetail(ctx context.Context, req *bridgeService.GetUomGPDetailRequest) (res *bridgeService.GetUomGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetUOMGPDetail")
	defer span.End()

	res, err = h.ServicesUom.GetDetailGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	return
}
