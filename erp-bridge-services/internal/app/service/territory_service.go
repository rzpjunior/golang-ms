package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/timex"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/repository"
)

type ITerritoryService interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, regionID int64, salespersonID int64, CustomerTypeID int64, subDistrictID int64) (res []dto.TerritoryResponse, total int64, err error)
	GetDetail(ctx context.Context, id int64, code string, salespersonID int64) (res dto.TerritoryResponse, err error)
}

type TerritoryService struct {
	opt                 opt.Options
	RepositoryTerritory repository.ITerritoryRepository
}

func NewTerritoryService() ITerritoryService {
	return &TerritoryService{
		opt:                 global.Setup.Common,
		RepositoryTerritory: repository.NewTerritoryRepository(),
	}
}

func (s *TerritoryService) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, regionID int64, salespersonID int64, CustomerTypeID int64, subDistrictID int64) (res []dto.TerritoryResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "TerritoryService.Get")
	defer span.End()

	var territories []*model.Territory
	territories, total, err = s.RepositoryTerritory.Get(ctx, offset, limit, status, search, orderBy, regionID, salespersonID, CustomerTypeID, subDistrictID)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, territory := range territories {
		res = append(res, dto.TerritoryResponse{
			ID:             territory.ID,
			Code:           territory.Code,
			Description:    territory.Description,
			RegionID:       territory.RegionID,
			SalespersonID:  territory.SalespersonID,
			CustomerTypeID: territory.CustomerTypeID,
			SubDistrictID:  territory.SubDistrictID,
			CreatedAt:      timex.ToLocTime(ctx, territory.CreatedAt),
			UpdatedAt:      timex.ToLocTime(ctx, territory.UpdatedAt),
		})
	}

	return
}

func (s *TerritoryService) GetDetail(ctx context.Context, id int64, code string, salespersonID int64) (res dto.TerritoryResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "TerritoryService.GetDetail")
	defer span.End()

	var territory *model.Territory
	territory, err = s.RepositoryTerritory.GetDetail(ctx, id, code, salespersonID)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.TerritoryResponse{
		ID:             territory.ID,
		Code:           territory.Code,
		Description:    territory.Description,
		RegionID:       territory.RegionID,
		SalespersonID:  territory.SalespersonID,
		CustomerTypeID: territory.CustomerTypeID,
		SubDistrictID:  territory.SubDistrictID,
		CreatedAt:      timex.ToLocTime(ctx, territory.CreatedAt),
		UpdatedAt:      timex.ToLocTime(ctx, territory.UpdatedAt),
	}

	return
}
