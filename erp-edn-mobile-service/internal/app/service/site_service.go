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

func NewServiceSite() ISiteService {
	m := new(SiteService)
	m.opt = global.Setup.Common
	return m
}

type ISiteService interface {
	GetSites(ctx context.Context, req dto.SiteListRequest) (res []*dto.SiteResponse, err error)
	GetSiteDetailById(ctx context.Context, req dto.SiteDetailRequest) (res *dto.SiteResponse, err error)
	GetGP(ctx context.Context, req dto.SiteListRequest) (res []*dto.SiteGP, total int64, err error)
	GetDetaiGPlById(ctx context.Context, id string) (res *dto.SiteGP, err error)
}

type SiteService struct {
	opt opt.Options
}

func (s *SiteService) GetSites(ctx context.Context, req dto.SiteListRequest) (res []*dto.SiteResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SiteService.GetSites")
	defer span.End()

	// get sites from bridge
	var siteRes *bridgeService.GetSiteListResponse
	siteRes, err = s.opt.Client.BridgeServiceGrpc.GetSiteList(ctx, &bridgeService.GetSiteListRequest{
		Limit:   req.Limit,
		Offset:  req.Offset,
		Status:  req.Status,
		Search:  req.Search,
		OrderBy: req.OrderBy,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "site")
		return
	}

	datas := []*dto.SiteResponse{}
	for _, site := range siteRes.Data {
		datas = append(datas, &dto.SiteResponse{
			// ID:            site.Id,
			Code:          site.Code,
			Name:          site.Description,
			Description:   site.Description,
			Status:        int8(site.Status),
			StatusConvert: statusx.ConvertStatusValue(int8(site.Status)),
			CreatedAt:     site.CreatedAt.AsTime(),
			UpdatedAt:     site.UpdatedAt.AsTime(),
		})
	}
	res = datas

	return
}

func (s *SiteService) GetSiteDetailById(ctx context.Context, req dto.SiteDetailRequest) (res *dto.SiteResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SiteService.GetSiteDetailById")
	defer span.End()

	// get Site from bridge
	var siteRes *bridgeService.GetSiteDetailResponse
	siteRes, err = s.opt.Client.BridgeServiceGrpc.GetSiteDetail(ctx, &bridgeService.GetSiteDetailRequest{
		Id: int64(req.Id),
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "site")
		return
	}

	res = &dto.SiteResponse{
		// ID:            siteRes.Data.Id,
		Code:          siteRes.Data.Code,
		Name:          siteRes.Data.Description,
		Description:   siteRes.Data.Description,
		Status:        int8(siteRes.Data.Status),
		StatusConvert: statusx.ConvertStatusValue(int8(siteRes.Data.Status)),
		CreatedAt:     siteRes.Data.CreatedAt.AsTime(),
		UpdatedAt:     siteRes.Data.UpdatedAt.AsTime(),
	}

	return
}

func (s *SiteService) GetGP(ctx context.Context, req dto.SiteListRequest) (res []*dto.SiteGP, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SiteService.GetGP")
	defer span.End()

	// get site from bridge
	var siteRes *bridgeService.GetSiteGPResponse
	siteRes, err = s.opt.Client.BridgeServiceGrpc.GetSiteGPList(ctx, &bridgeService.GetSiteGPListRequest{
		Limit:    req.Limit,
		Offset:   req.Offset,
		Locndscr: req.Search,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "site")
		return
	}

	datas := []*dto.SiteGP{}
	for _, site := range siteRes.Data {
		datas = append(datas, &dto.SiteGP{
			ID:               site.Locncode,
			GnL_Site_Type_ID: site.GnlSiteTypeId,
			Name:             site.Locndscr,
			Description:      site.Locndscr,

			Address: site.AddresS1 + site.AddresS2 + site.AddresS3,
			// AddresS2:                site.AddresS2,
			// AddresS3:                site.AddresS3,
			PhonE1:                  site.PhonE1,
			PhonE2:                  site.PhonE2,
			PhonE3:                  site.PhonE3,
			City:                    site.City,
			State:                   site.State,
			Faxnumbr:                site.Faxnumbr,
			Zipcode:                 site.Zipcode,
			Ccode:                   site.Ccode,
			Country:                 site.Country,
			Staxschd:                site.Staxschd,
			Pctaxsch:                site.Pctaxsch,
			Location_Segment:        site.LocationSegment,
			GnL_Administrative_Code: site.GnlAdministrativeCode,
			Inactive:                site.Inactive,
		})
	}

	total = int64(len(datas))
	res = datas

	return
}

func (s *SiteService) GetDetaiGPlById(ctx context.Context, id string) (res *dto.SiteGP, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SiteService.GetDetaiGPlById")
	defer span.End()

	// get site from bridge
	var site *bridgeService.GetSiteGPResponse
	site, err = s.opt.Client.BridgeServiceGrpc.GetSiteGPDetail(ctx, &bridgeService.GetSiteGPDetailRequest{
		Id: id,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "site")
		return
	}

	res = &dto.SiteGP{
		ID:               site.Data[0].Locncode,
		GnL_Site_Type_ID: site.Data[0].GnlSiteTypeId,
		Name:             site.Data[0].Locndscr,
		Description:      site.Data[0].Locndscr,
		Address:          site.Data[0].AddresS1 + site.Data[0].AddresS2 + site.Data[0].AddresS3,
		// AddresS2:                site.Data[0].AddresS2,
		// AddresS3:                site.Data[0].AddresS3,
		PhonE1:                  site.Data[0].PhonE1,
		PhonE2:                  site.Data[0].PhonE2,
		PhonE3:                  site.Data[0].PhonE3,
		City:                    site.Data[0].City,
		State:                   site.Data[0].State,
		Faxnumbr:                site.Data[0].Faxnumbr,
		Zipcode:                 site.Data[0].Zipcode,
		Ccode:                   site.Data[0].Ccode,
		Country:                 site.Data[0].Country,
		Staxschd:                site.Data[0].Staxschd,
		Pctaxsch:                site.Data[0].Pctaxsch,
		Location_Segment:        site.Data[0].LocationSegment,
		GnL_Administrative_Code: site.Data[0].GnlAdministrativeCode,
		Inactive:                site.Data[0].Inactive,
	}

	return
}
