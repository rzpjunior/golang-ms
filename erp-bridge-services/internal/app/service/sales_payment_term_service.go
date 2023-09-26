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

type ISalesPaymentTermService interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (res []dto.SalesPaymentTermResponse, total int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (res dto.SalesPaymentTermResponse, err error)
}

type SalesPaymentTermService struct {
	opt                        opt.Options
	RepositorySalesPaymentTerm repository.ISalesPaymentTermRepository
}

func NewSalesPaymentTermService() ISalesPaymentTermService {
	return &SalesPaymentTermService{
		opt:                        global.Setup.Common,
		RepositorySalesPaymentTerm: repository.NewSalesPaymentTermRepository(),
	}
}

func (s *SalesPaymentTermService) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (res []dto.SalesPaymentTermResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesPaymentTermService.Get")
	defer span.End()

	var salesPaymentTerms []*model.SalesPaymentTerm
	salesPaymentTerms, total, err = s.RepositorySalesPaymentTerm.Get(ctx, offset, limit, status, search, orderBy)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, salesPaymentTerm := range salesPaymentTerms {
		res = append(res, dto.SalesPaymentTermResponse{
			ID:            salesPaymentTerm.ID,
			Code:          salesPaymentTerm.Code,
			Description:   salesPaymentTerm.Description,
			Status:        salesPaymentTerm.Status,
			StatusConvert: statusx.ConvertStatusValue(salesPaymentTerm.Status),
			CreatedAt:     timex.ToLocTime(ctx, salesPaymentTerm.CreatedAt),
			UpdatedAt:     timex.ToLocTime(ctx, salesPaymentTerm.UpdatedAt),
		})
	}

	return
}

func (s *SalesPaymentTermService) GetDetail(ctx context.Context, id int64, code string) (res dto.SalesPaymentTermResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesPaymentTermService.GetDetail")
	defer span.End()

	var salesPaymentTerm *model.SalesPaymentTerm
	salesPaymentTerm, err = s.RepositorySalesPaymentTerm.GetDetail(ctx, id, code)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.SalesPaymentTermResponse{
		ID:            salesPaymentTerm.ID,
		Code:          salesPaymentTerm.Code,
		Description:   salesPaymentTerm.Description,
		Status:        salesPaymentTerm.Status,
		StatusConvert: statusx.ConvertStatusValue(salesPaymentTerm.Status),
		CreatedAt:     timex.ToLocTime(ctx, salesPaymentTerm.CreatedAt),
		UpdatedAt:     timex.ToLocTime(ctx, salesPaymentTerm.UpdatedAt),
	}

	return
}
