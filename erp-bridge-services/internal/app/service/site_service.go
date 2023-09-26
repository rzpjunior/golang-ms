package service

import (
	"context"
	"net/url"
	"strconv"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/repository"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

type ISiteService interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (res []dto.SiteResponse, total int64, err error)
	GetInIds(ctx context.Context, offset int, limit int, status int, ids []int64, orderBy string) (res []dto.SiteResponse, total int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (res dto.SiteResponse, err error)
	GetGP(ctx context.Context, req *pb.GetSiteGPListRequest) (res *pb.GetSiteGPResponse, err error)
	GetDetailGP(ctx context.Context, req *pb.GetSiteGPDetailRequest) (res *pb.GetSiteGPResponse, err error)
}

type SiteService struct {
	opt              opt.Options
	RepositorySite   repository.ISiteRepository
	RepositoryRegion repository.IRegionRepository
}

func NewSiteService() ISiteService {
	return &SiteService{
		opt:              global.Setup.Common,
		RepositorySite:   repository.NewSiteRepository(),
		RepositoryRegion: repository.NewRegionRepository(),
	}
}

func (s *SiteService) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (res []dto.SiteResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SiteService.Get")
	defer span.End()

	var sites []*model.Site
	sites, total, err = s.RepositorySite.Get(ctx, offset, limit, status, search, orderBy)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, site := range sites {
		res = append(res, dto.SiteResponse{
			ID:          site.ID,
			Code:        site.Code,
			Description: site.Description,
			Status:      site.Status,
			CreatedAt:   site.CreatedAt,
			UpdatedAt:   site.UpdatedAt,
		})
	}

	return
}

func (s *SiteService) GetInIds(ctx context.Context, offset int, limit int, status int, ids []int64, orderBy string) (res []dto.SiteResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SiteService.Get")
	defer span.End()

	var sites []*model.Site
	sites, total, err = s.RepositorySite.GetInIds(ctx, offset, limit, status, ids, orderBy)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, site := range sites {
		res = append(res, dto.SiteResponse{
			ID:          site.ID,
			Code:        site.Code,
			Description: site.Description,
			Status:      site.Status,
			CreatedAt:   site.CreatedAt,
			UpdatedAt:   site.UpdatedAt,
		})
	}

	return
}

func (s *SiteService) GetDetail(ctx context.Context, id int64, code string) (res dto.SiteResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SiteService.GetDetail")
	defer span.End()

	var (
		site   *model.Site
		region *model.Region
	)
	site, err = s.RepositorySite.GetDetail(ctx, id, code)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	region, err = s.RepositoryRegion.GetDetail(ctx, site.RegionId, "")
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.SiteResponse{
		ID:          site.ID,
		Code:        site.Code,
		Description: site.Description,
		Status:      site.Status,
		CreatedAt:   site.CreatedAt,
		UpdatedAt:   site.UpdatedAt,
		Region: &dto.RegionResponse{
			ID:            region.ID,
			Code:          region.Code,
			Description:   region.Description,
			Status:        region.Status,
			StatusConvert: statusx.ConvertStatusValue(region.Status),
			CreatedAt:     region.CreatedAt,
			UpdatedAt:     region.UpdatedAt,
		},
	}

	return
}

func (s *SiteService) GetGP(ctx context.Context, req *pb.GetSiteGPListRequest) (res *pb.GetSiteGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SiteService.GetGP")
	defer span.End()

	params := map[string]string{
		"interid":    global.EnvDatabaseGP,
		"PageNumber": strconv.Itoa(int(req.Offset)),
		"PageSize":   strconv.Itoa(int(req.Limit)),
	}

	if req.Locncode != "" {
		params["locncode_like"] = url.PathEscape(req.Locncode)
	}

	if req.Locndscr != "" {
		params["locndscr"] = url.PathEscape(req.Locndscr)
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "warehouselist/getall", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *SiteService) GetDetailGP(ctx context.Context, req *pb.GetSiteGPDetailRequest) (res *pb.GetSiteGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SiteService.GetDetailGP")
	defer span.End()

	params := map[string]string{
		"interid": global.EnvDatabaseGP,
		"id":      req.Id,
	}

	// fetch site from gp
	err = global.HttpRestApiToMicrosoftGP("GET", "warehouselist/getbyid", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}
