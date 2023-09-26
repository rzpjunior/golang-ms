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

type ISalespersonService interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (res []dto.SalespersonResponse, total int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (res dto.SalespersonResponse, err error)
}

type SalespersonService struct {
	opt                   opt.Options
	RepositorySalesperson repository.ISalespersonRepository
}

func NewSalespersonService() ISalespersonService {
	return &SalespersonService{
		opt:                   global.Setup.Common,
		RepositorySalesperson: repository.NewSalespersonRepository(),
	}
}

func (s *SalespersonService) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (res []dto.SalespersonResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalespersonService.Get")
	defer span.End()

	var salespersons []*model.Salesperson
	salespersons, total, err = s.RepositorySalesperson.Get(ctx, offset, limit, status, search, orderBy)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, salesperson := range salespersons {
		res = append(res, dto.SalespersonResponse{
			ID:            salesperson.ID,
			Code:          salesperson.Code,
			FirstName:     salesperson.FirstName,
			MiddleName:    salesperson.MiddleName,
			LastName:      salesperson.LastName,
			Status:        salesperson.Status,
			StatusConvert: statusx.ConvertStatusValue(salesperson.Status),
			CreatedAt:     timex.ToLocTime(ctx, salesperson.CreatedAt),
			UpdatedAt:     timex.ToLocTime(ctx, salesperson.UpdatedAt),
		})
	}

	return
}

func (s *SalespersonService) GetDetail(ctx context.Context, id int64, code string) (res dto.SalespersonResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalespersonService.GetDetail")
	defer span.End()

	var salesperson *model.Salesperson
	salesperson, err = s.RepositorySalesperson.GetDetail(ctx, id, code)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.SalespersonResponse{
		ID:            salesperson.ID,
		Code:          salesperson.Code,
		FirstName:     salesperson.FirstName,
		MiddleName:    salesperson.MiddleName,
		LastName:      salesperson.LastName,
		Status:        salesperson.Status,
		StatusConvert: statusx.ConvertStatusValue(salesperson.Status),
		CreatedAt:     timex.ToLocTime(ctx, salesperson.CreatedAt),
		UpdatedAt:     timex.ToLocTime(ctx, salesperson.UpdatedAt),
	}

	return
}
