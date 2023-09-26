package handler

import (
	context "context"
	"fmt"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *BridgeGrpcHandler) GetDistrictList(ctx context.Context, req *bridgeService.GetDistrictListRequest) (res *bridgeService.GetDistrictListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetDistrictList")
	defer span.End()

	var districts []*dto.DistrictResponse
	districts, _, err = h.ServicesDistrict.Get(ctx, int(req.Offset), int(req.Limit), int(req.Status), req.Search, req.OrderBy)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*bridgeService.District
	for _, district := range districts {
		data = append(data, &bridgeService.District{
			Id:            district.ID,
			Code:          district.Code,
			Value:         district.Value,
			Name:          district.Name,
			Note:          district.Note,
			Status:        int32(district.Status),
			StatusConvert: district.StatusConvert,
		})
	}

	res = &bridgeService.GetDistrictListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *BridgeGrpcHandler) GetDistrictDetail(ctx context.Context, req *bridgeService.GetDistrictDetailRequest) (res *bridgeService.GetDistrictDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetSubDistrictDetail")
	defer span.End()

	var district *dto.DistrictResponse
	district, err = h.ServicesDistrict.GetDetail(ctx, req.Id, req.Code)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.GetDistrictDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &bridgeService.District{
			Id:            district.ID,
			Code:          district.Code,
			Value:         district.Value,
			Name:          district.Name,
			Note:          district.Note,
			Status:        int32(district.Status),
			StatusConvert: district.StatusConvert,
		},
	}
	return
}

func (h *BridgeGrpcHandler) GetDistrictInIdsList(ctx context.Context, req *bridgeService.GetDistrictInIdsListRequest) (res *bridgeService.GetDistrictListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetDistrictList")
	defer span.End()

	var districts []*dto.DistrictResponse
	fmt.Println(h.ServicesDistrict)
	districts, _, err = h.ServicesDistrict.GetInIds(ctx, int(req.Offset), int(req.Limit), req.Ids, int(req.Status), req.Search, req.OrderBy)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*bridgeService.District
	for _, district := range districts {
		data = append(data, &bridgeService.District{
			Id:            district.ID,
			Code:          district.Code,
			Value:         district.Value,
			Name:          district.Name,
			Note:          district.Note,
			Status:        int32(district.Status),
			StatusConvert: district.StatusConvert,
		})
	}

	res = &bridgeService.GetDistrictListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}
