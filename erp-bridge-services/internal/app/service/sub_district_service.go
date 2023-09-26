package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/repository"
)

type ISubDistrictService interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (res []dto.SubDistrictResponse, total int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (res dto.SubDistrictResponse, err error)
}

type SubDistrictService struct {
	opt                   opt.Options
	RepositorySubDistrict repository.ISubDistrictRepository
}

func NewSubDistrictService() ISubDistrictService {
	return &SubDistrictService{
		opt:                   global.Setup.Common,
		RepositorySubDistrict: repository.NewSubDistrictRepository(),
	}
}

func (s *SubDistrictService) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (res []dto.SubDistrictResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SubDistrictService.Get")
	defer span.End()

	var subDistricts []*model.SubDistrict
	subDistricts, total, err = s.RepositorySubDistrict.Get(ctx, offset, limit, status, search, orderBy)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, subDistrict := range subDistricts {
		res = append(res, dto.SubDistrictResponse{
			ID:        subDistrict.ID,
			Status:    subDistrict.Status,
			CreatedAt: subDistrict.CreatedAt,
			UpdatedAt: subDistrict.UpdatedAt,
		})
	}

	return
}

func (s *SubDistrictService) GetDetail(ctx context.Context, id int64, code string) (res dto.SubDistrictResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SubDistrictService.GetDetail")
	defer span.End()

	var subDistrict *model.SubDistrict
	subDistrict, err = s.RepositorySubDistrict.GetDetail(ctx, id, code)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.SubDistrictResponse{
		ID:          subDistrict.ID,
		Code:        subDistrict.Code,
		Description: subDistrict.Description,
		Status:      subDistrict.Status,
		CreatedAt:   subDistrict.CreatedAt,
		UpdatedAt:   subDistrict.UpdatedAt,
	}

	return
}
