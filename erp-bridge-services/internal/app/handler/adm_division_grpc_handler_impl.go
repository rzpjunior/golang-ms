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

func (h *BridgeGrpcHandler) GetAdmDivisionList(ctx context.Context, req *bridgeService.GetAdmDivisionListRequest) (res *bridgeService.GetAdmDivisionListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetAdmDivisionList")
	defer span.End()

	var admDivisions []dto.AdmDivisionResponse
	admDivisions, _, err = h.ServicesAdmDivision.GetListGRPC(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*bridgeService.AdmDivision
	for _, admDivision := range admDivisions {
		data = append(data, &bridgeService.AdmDivision{
			Id:            admDivision.ID,
			Code:          admDivision.Code,
			RegionId:      admDivision.RegionID,
			City:          admDivision.City,
			District:      admDivision.District,
			SubDistrictId: admDivision.SubDistrictID,
			PostalCode:    admDivision.PostalCode,
			Status:        int32(admDivision.Status),
			CreatedAt:     timestamppb.New(admDivision.CreatedAt),
			UpdatedAt:     timestamppb.New(admDivision.UpdatedAt),
		})
	}

	res = &bridgeService.GetAdmDivisionListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *BridgeGrpcHandler) GetAdmDivisionDetail(ctx context.Context, req *bridgeService.GetAdmDivisionDetailRequest) (res *bridgeService.GetAdmDivisionDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetAdmDivisionDetail")
	defer span.End()

	var admDivision dto.AdmDivisionResponse
	admDivision, err = h.ServicesAdmDivision.GetDetailGRPC(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.GetAdmDivisionDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &bridgeService.AdmDivision{
			Id:            admDivision.ID,
			Code:          admDivision.Code,
			RegionId:      admDivision.RegionID,
			City:          admDivision.City,
			District:      admDivision.District,
			SubDistrictId: admDivision.SubDistrictID,
			PostalCode:    admDivision.PostalCode,
			Status:        int32(admDivision.Status),
			CreatedAt:     timestamppb.New(admDivision.CreatedAt),
			UpdatedAt:     timestamppb.New(admDivision.UpdatedAt),
		},
	}
	return
}

func (h *BridgeGrpcHandler) GetAdmDivisionGPList(ctx context.Context, req *bridgeService.GetAdmDivisionGPListRequest) (res *bridgeService.GetAdmDivisionGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetAdmDivisionGPList")
	defer span.End()

	res, err = h.ServicesAdmDivision.GetGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	return
}

func (h *BridgeGrpcHandler) GetAdmDivisionGPDetail(ctx context.Context, req *bridgeService.GetAdmDivisionGPDetailRequest) (res *bridgeService.GetAdmDivisionGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetAdmDivisionGPDetail")
	defer span.End()

	res, err = h.ServicesAdmDivision.GetDetailGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	return
}
