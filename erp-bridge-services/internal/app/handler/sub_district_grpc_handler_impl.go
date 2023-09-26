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

func (h *BridgeGrpcHandler) GetSubDistrictList(ctx context.Context, req *bridgeService.GetSubDistrictListRequest) (res *bridgeService.GetSubDistrictListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetSubDistrictList")
	defer span.End()

	var subDistricts []dto.SubDistrictResponse
	subDistricts, _, err = h.ServicesSubDistrict.Get(ctx, int(req.Offset), int(req.Limit), int(req.Status), req.Search, req.OrderBy)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*bridgeService.SubDistrict
	for _, subDistrict := range subDistricts {
		data = append(data, &bridgeService.SubDistrict{
			Id:          subDistrict.ID,
			Code:        subDistrict.Code,
			Description: subDistrict.Description,
			Status:      int32(subDistrict.Status),
			CreatedAt:   timestamppb.New(subDistrict.CreatedAt),
			UpdatedAt:   timestamppb.New(subDistrict.UpdatedAt),
		})
	}

	res = &bridgeService.GetSubDistrictListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *BridgeGrpcHandler) GetSubDistrictDetail(ctx context.Context, req *bridgeService.GetSubDistrictDetailRequest) (res *bridgeService.GetSubDistrictDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetSubDistrictDetail")
	defer span.End()

	var subDistrict dto.SubDistrictResponse
	subDistrict, err = h.ServicesSubDistrict.GetDetail(ctx, req.Id, req.Code)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.GetSubDistrictDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &bridgeService.SubDistrict{
			Id:          subDistrict.ID,
			Code:        subDistrict.Code,
			Description: subDistrict.Description,
			Status:      int32(subDistrict.Status),
			CreatedAt:   timestamppb.New(subDistrict.CreatedAt),
			UpdatedAt:   timestamppb.New(subDistrict.UpdatedAt),
		},
	}
	return
}
