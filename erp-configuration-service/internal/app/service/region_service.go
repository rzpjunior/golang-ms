package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-configuration-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

type IRegionService interface {
	Get(ctx context.Context, req *dto.RegionGetRequest) (res []*dto.RegionResponse, total int64, err error)
}

type RegionService struct {
	opt opt.Options
}

func NewRegionService() IRegionService {
	return &RegionService{
		opt: global.Setup.Common,
	}
}

func (s *RegionService) Get(ctx context.Context, req *dto.RegionGetRequest) (res []*dto.RegionResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "RegionService.Get")
	defer span.End()

	var regions *bridge_service.GetAdmDivisionGPResponse
	regions, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPDetail(ctx, &bridge_service.GetAdmDivisionGPDetailRequest{
		Type: "region",
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, region := range regions.Data {
		res = append(res, &dto.RegionResponse{
			ID:          region.Region,
			Code:        region.Region,
			Description: region.Region,
		})
	}

	total = int64(len(res))

	return
}
