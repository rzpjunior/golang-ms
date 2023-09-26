package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-site-service/internal/app/dto"
)

type ISiteService interface {
	Get(ctx context.Context, req dto.GetSiteRequest) (res []dto.FilterSiteResponse, total int64, err error)
	GetDetail(ctx context.Context, id string) (res dto.FilterSiteResponse, err error)
}

type SiteService struct {
	opt opt.Options
}

func NewServiceSite() ISiteService {
	return &SiteService{
		opt: global.Setup.Common,
	}
}

func (s *SiteService) Get(ctx context.Context, req dto.GetSiteRequest) (res []dto.FilterSiteResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SiteService.Get")
	defer span.End()

	var site *bridgeService.GetSiteGPResponse

	if site, err = s.opt.Client.BridgeServiceGrpc.GetSiteGPList(ctx, &bridgeService.GetSiteGPListRequest{
		Limit:    int32(req.Limit),
		Offset:   int32(req.Offset),
		Locncode: req.Search,
		Locndscr: req.Search,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "site")
		return
	}

	for _, site := range site.Data {
		res = append(res, dto.FilterSiteResponse{
			ID:   site.Locncode,
			Name: site.Locndscr,
		})
	}

	total = int64(len(site.Data))

	return
}

func (s *SiteService) GetDetail(ctx context.Context, id string) (res dto.FilterSiteResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SiteService.Get")
	defer span.End()

	var site *bridgeService.GetSiteGPResponse
	if site, err = s.opt.Client.BridgeServiceGrpc.GetSiteGPDetail(ctx, &bridgeService.GetSiteGPDetailRequest{
		Id: id,
	}); err != nil || !site.Succeeded {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "site")
		return
	}

	res = dto.FilterSiteResponse{
		ID:   site.Data[0].Locncode,
		Name: site.Data[0].Locndscr,
	}

	return
}
