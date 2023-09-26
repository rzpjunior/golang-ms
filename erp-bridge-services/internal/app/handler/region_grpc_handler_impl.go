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

func (h *BridgeGrpcHandler) GetRegionList(ctx context.Context, req *bridgeService.GetRegionListRequest) (res *bridgeService.GetRegionListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetRegionList")
	defer span.End()

	var regions []dto.RegionResponse
	regions, _, err = h.ServicesRegion.Get(ctx, int(req.Offset), int(req.Limit), int(req.Status), req.Search, req.OrderBy)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*bridgeService.Region
	for _, region := range regions {
		data = append(data, &bridgeService.Region{
			Id:          region.ID,
			Code:        region.Code,
			Description: region.Description,
			Status:      int32(region.Status),
			CreatedAt:   timestamppb.New(region.CreatedAt),
			UpdatedAt:   timestamppb.New(region.UpdatedAt),
		})
	}

	res = &bridgeService.GetRegionListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *BridgeGrpcHandler) GetRegionDetail(ctx context.Context, req *bridgeService.GetRegionDetailRequest) (res *bridgeService.GetRegionDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetRegionDetail")
	defer span.End()

	var region dto.RegionResponse
	region, err = h.ServicesRegion.GetDetail(ctx, req.Id, req.Code)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.GetRegionDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &bridgeService.Region{
			Id:          region.ID,
			Code:        region.Code,
			Description: region.Description,
			Status:      int32(region.Status),
			CreatedAt:   timestamppb.New(region.CreatedAt),
			UpdatedAt:   timestamppb.New(region.UpdatedAt),
		},
	}
	return
}
