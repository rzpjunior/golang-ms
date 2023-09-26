package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/edenlabs/edenlabs/timex"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/repository"
)

type IRegionService interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (res []dto.RegionResponse, total int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (res dto.RegionResponse, err error)
}

type RegionService struct {
	opt              opt.Options
	RepositoryRegion repository.IRegionRepository
}

func NewRegionService() IRegionService {
	return &RegionService{
		opt:              global.Setup.Common,
		RepositoryRegion: repository.NewRegionRepository(),
	}
}

func (s *RegionService) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (res []dto.RegionResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "RegionService.Get")
	defer span.End()

	var regions []*model.Region
	regions, total, err = s.RepositoryRegion.Get(ctx, offset, limit, status, search, orderBy)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, region := range regions {
		res = append(res, dto.RegionResponse{
			ID:            region.ID,
			Code:          region.Code,
			Description:   region.Description,
			Status:        region.Status,
			StatusConvert: statusx.ConvertStatusValue(region.Status),
			CreatedAt:     timex.ToLocTime(ctx, region.CreatedAt),
			UpdatedAt:     timex.ToLocTime(ctx, region.UpdatedAt),
		})
	}

	return
}

func (s *RegionService) GetDetail(ctx context.Context, id int64, code string) (res dto.RegionResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "RegionService.GetDetail")
	defer span.End()

	var region *model.Region
	region, err = s.RepositoryRegion.GetDetail(ctx, id, code)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.RegionResponse{
		ID:            region.ID,
		Code:          region.Code,
		Description:   region.Description,
		Status:        region.Status,
		StatusConvert: statusx.ConvertStatusValue(region.Status),
		CreatedAt:     timex.ToLocTime(ctx, region.CreatedAt),
		UpdatedAt:     timex.ToLocTime(ctx, region.UpdatedAt),
	}

	return
}
