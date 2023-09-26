package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *BridgeGrpcHandler) GetWrtList(ctx context.Context, req *bridgeService.GetWrtListRequest) (res *bridgeService.GetWrtListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetWrtList")
	defer span.End()

	var Wrtes []dto.WrtResponse
	Wrtes, _, err = h.ServicesWrt.Get(ctx, int(req.Offset), int(req.Limit), req.RegionId, req.Search)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*bridgeService.Wrt
	for _, Wrt := range Wrtes {
		region := &bridgeService.Region{
			Id:          Wrt.Region.ID,
			Code:        Wrt.Region.Code,
			Description: Wrt.Region.Description,
		}

		data = append(data, &bridgeService.Wrt{
			Id:        Wrt.ID,
			RegionId:  Wrt.RegionID,
			Code:      Wrt.Code,
			StartTime: Wrt.StartTime,
			EndTime:   Wrt.EndTime,
			Region:    region,
		})
	}

	res = &bridgeService.GetWrtListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *BridgeGrpcHandler) GetWrtDetail(ctx context.Context, req *bridgeService.GetWrtDetailRequest) (res *bridgeService.GetWrtDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetWrtDetail")
	defer span.End()

	var Wrt dto.WrtResponse

	Wrt, err = h.ServicesWrt.GetDetail(ctx, req.Id, "")
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data *bridgeService.Wrt

	region := &bridgeService.Region{
		Id:          Wrt.Region.ID,
		Code:        Wrt.Region.Code,
		Description: Wrt.Region.Description,
	}

	data = &bridgeService.Wrt{
		Id:        Wrt.ID,
		RegionId:  Wrt.RegionID,
		Code:      Wrt.Code,
		StartTime: Wrt.StartTime,
		EndTime:   Wrt.EndTime,
		Region:    region,
	}

	res = &bridgeService.GetWrtDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *BridgeGrpcHandler) GetWrtGPList(ctx context.Context, req *bridgeService.GetWrtGPListRequest) (res *bridgeService.GetWrtGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetWrtGPList")
	defer span.End()

	res, err = h.ServicesWrt.GetGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	return
}

func (h *BridgeGrpcHandler) GetWrtGPDetail(ctx context.Context, req *bridgeService.GetWrtGPDetailRequest) (res *bridgeService.GetWrtGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetWrtGPDetail")
	defer span.End()

	res, err = h.ServicesWrt.GetDetailGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	return
}
