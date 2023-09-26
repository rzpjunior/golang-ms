package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-promotion-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-promotion-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-promotion-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-promotion-service/internal/app/repository"
)

type IVoucherItemService interface {
	Get(ctx context.Context, req *dto.VoucherItemRequestGet) (res []*dto.VoucherItemResponse, total int64, err error)
}

type VoucherItemService struct {
	opt                   opt.Options
	RepositoryVoucherItem repository.IVoucherItemRepository
}

func NewVoucherItemService() IVoucherItemService {
	return &VoucherItemService{
		opt:                   global.Setup.Common,
		RepositoryVoucherItem: repository.NewVoucherItemRepository(),
	}
}

func (s *VoucherItemService) Get(ctx context.Context, req *dto.VoucherItemRequestGet) (res []*dto.VoucherItemResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "VoucherItemService.Get")
	defer span.End()

	var voucherItems []*model.VoucherItem
	voucherItems, total, err = s.RepositoryVoucherItem.Get(ctx, req)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, voucherItem := range voucherItems {
		res = append(res, &dto.VoucherItemResponse{
			ID:         voucherItem.ID,
			VoucherID:  voucherItem.VoucherID,
			ItemID:     voucherItem.ItemID,
			MinQtyDisc: voucherItem.MinQtyDisc,
			CreatedAt:  voucherItem.CreatedAt,
		})
	}

	return
}
