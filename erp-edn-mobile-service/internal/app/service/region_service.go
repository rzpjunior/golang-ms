package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/global"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	dto "git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/internal/app/dto"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

func NewServiceRegion() IRegionService {
	m := new(RegionService)
	m.opt = global.Setup.Common
	return m
}

type IRegionService interface {
	GetRegions(ctx context.Context, req dto.RegionListRequest) (res []*dto.RegionResponse, err error)
	GetRegionDetailById(ctx context.Context, req dto.RegionDetailRequest) (res *dto.RegionResponse, err error)
}

type RegionService struct {
	opt opt.Options
}

func (s *RegionService) GetRegions(ctx context.Context, req dto.RegionListRequest) (res []*dto.RegionResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "RegionService.GetRegions")
	defer span.End()

	// get region from bridge
	var regionRes *bridgeService.GetRegionListResponse
	regionRes, err = s.opt.Client.BridgeServiceGrpc.GetRegionList(ctx, &bridgeService.GetRegionListRequest{
		Limit:   req.Limit,
		Offset:  req.Offset,
		Status:  req.Status,
		Search:  req.Search,
		OrderBy: req.OrderBy,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "region")
		return
	}

	datas := []*dto.RegionResponse{}
	for _, region := range regionRes.Data {
		datas = append(datas, &dto.RegionResponse{
			// ID:            region.Id,
			Code:          region.Code,
			Description:   region.Description,
			Status:        int8(region.Status),
			StatusConvert: statusx.ConvertStatusValue(int8(region.Status)),
			CreatedAt:     region.CreatedAt.AsTime(),
			UpdatedAt:     region.UpdatedAt.AsTime(),
		})
	}
	res = datas

	return
}

func (s *RegionService) GetRegionDetailById(ctx context.Context, req dto.RegionDetailRequest) (res *dto.RegionResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "RegionService.GetRegionDetailById")
	defer span.End()

	// get Region from bridge
	var regionRes *bridgeService.GetRegionDetailResponse
	regionRes, err = s.opt.Client.BridgeServiceGrpc.GetRegionDetail(ctx, &bridgeService.GetRegionDetailRequest{
		Id: int64(req.Id),
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "site")
		return
	}

	res = &dto.RegionResponse{
		// ID:            regionRes.Data.Id,
		Code:          regionRes.Data.Code,
		Description:   regionRes.Data.Description,
		Status:        int8(regionRes.Data.Status),
		StatusConvert: statusx.ConvertStatusValue(int8(regionRes.Data.Status)),
		CreatedAt:     regionRes.Data.CreatedAt.AsTime(),
		UpdatedAt:     regionRes.Data.UpdatedAt.AsTime(),
	}

	return
}
