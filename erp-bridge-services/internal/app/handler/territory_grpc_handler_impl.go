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

func (h *BridgeGrpcHandler) GetTerritoryList(ctx context.Context, req *bridgeService.GetTerritoryListRequest) (res *bridgeService.GetTerritoryListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetTerritoryList")
	defer span.End()

	var territorys []dto.TerritoryResponse
	territorys, _, err = h.ServicesTerritory.Get(ctx, int(req.Offset), int(req.Limit), int(req.Status), req.Search, req.OrderBy, req.RegionId, req.SalespersonId, req.CustomerTypeId, req.SubDistrictId)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*bridgeService.Territory
	for _, territory := range territorys {
		data = append(data, &bridgeService.Territory{
			Id:             territory.ID,
			Code:           territory.Code,
			Description:    territory.Description,
			RegionId:       territory.RegionID,
			SalespersonId:  territory.SalespersonID,
			CustomerTypeId: territory.CustomerTypeID,
			SubDistrictId:  territory.SubDistrictID,
			CreatedAt:      timestamppb.New(territory.CreatedAt),
			UpdatedAt:      timestamppb.New(territory.UpdatedAt),
		})
	}

	res = &bridgeService.GetTerritoryListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *BridgeGrpcHandler) GetTerritoryDetail(ctx context.Context, req *bridgeService.GetTerritoryDetailRequest) (res *bridgeService.GetTerritoryDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetTerritoryDetail")
	defer span.End()

	var territory dto.TerritoryResponse
	territory, err = h.ServicesTerritory.GetDetail(ctx, req.Id, req.Code, req.SalespersonId)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.GetTerritoryDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &bridgeService.Territory{
			Id:             territory.ID,
			Code:           territory.Code,
			Description:    territory.Description,
			RegionId:       territory.RegionID,
			SalespersonId:  territory.SalespersonID,
			CustomerTypeId: territory.CustomerTypeID,
			SubDistrictId:  territory.SubDistrictID,
			CreatedAt:      timestamppb.New(territory.CreatedAt),
			UpdatedAt:      timestamppb.New(territory.UpdatedAt),
		},
	}
	return
}
