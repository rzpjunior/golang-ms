package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/repository"
)

type IDivisionService interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (res []*dto.DivisionResponse, total int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (res *dto.DivisionResponse, err error)
}

type DivisionService struct {
	opt                opt.Options
	RepositoryDivision repository.IDivisionRepository
}

func NewDivisionService() IDivisionService {
	return &DivisionService{
		opt:                global.Setup.Common,
		RepositoryDivision: repository.NewDivisionRepository(),
	}
}

func (s *DivisionService) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (res []*dto.DivisionResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DivisionService.Get")
	defer span.End()

	var divisions []*model.Division
	divisions, total, err = s.RepositoryDivision.Get(ctx, offset, limit, status, search, orderBy)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, CustomerType := range divisions {
		res = append(res, &dto.DivisionResponse{
			ID:            CustomerType.ID,
			Code:          CustomerType.Code,
			Name:          CustomerType.Name,
			Note:          CustomerType.Note,
			Status:        CustomerType.Status,
			StatusConvert: statusx.ConvertStatusValue(CustomerType.Status),
		})
	}

	return
}

func (s *DivisionService) GetDetail(ctx context.Context, id int64, code string) (res *dto.DivisionResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DivisionService.GetDetail")
	defer span.End()

	var division *model.Division
	division, err = s.RepositoryDivision.GetDetail(ctx, id, code)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &dto.DivisionResponse{
		ID:            division.ID,
		Code:          division.Code,
		Name:          division.Name,
		Note:          division.Note,
		Status:        division.Status,
		StatusConvert: statusx.ConvertStatusValue(division.Status),
	}

	return
}
