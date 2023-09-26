package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *BridgeGrpcHandler) GetDivisionList(ctx context.Context, req *bridgeService.GetDivisionListRequest) (res *bridgeService.GetDivisionListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetDivisionList")
	defer span.End()

	var divisions []*dto.DivisionResponse
	divisions, _, err = h.ServicesDivision.Get(ctx, int(req.Offset), int(req.Limit), int(req.Status), req.Search, req.OrderBy)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*bridgeService.Division
	for _, division := range divisions {
		data = append(data, &bridgeService.Division{
			Id:     division.ID,
			Code:   division.Code,
			Name:   division.Name,
			Note:   division.Note,
			Status: int32(division.Status),
		})
	}

	res = &bridgeService.GetDivisionListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *BridgeGrpcHandler) GetDivisionDetail(ctx context.Context, req *bridgeService.GetDivisionDetailRequest) (res *bridgeService.GetDivisionDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetDivisionDetail")
	defer span.End()

	var address *dto.DivisionResponse
	address, err = h.ServicesDivision.GetDetail(ctx, req.Id, req.Code)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.GetDivisionDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &bridgeService.Division{
			Id:     address.ID,
			Code:   address.Code,
			Name:   address.Name,
			Note:   address.Note,
			Status: int32(address.Status),
		},
	}
	return
}
