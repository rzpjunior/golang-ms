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

	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

type IBankService interface {
	Get(ctx context.Context, req *pb.GetBankListRequest) (res []dto.BankResponse, total int64, err error)
	GetDetail(ctx context.Context, req *pb.GetBankDetailRequest) (res dto.BankResponse, err error)
}

type BankService struct {
	opt            opt.Options
	RepositoryBank repository.IBankRepository
}

func NewBankService() IBankService {
	return &BankService{
		opt:            global.Setup.Common,
		RepositoryBank: repository.NewBankRepository(),
	}
}

func (s *BankService) Get(ctx context.Context, req *pb.GetBankListRequest) (res []dto.BankResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "BankService.Get")
	defer span.End()

	var Bankes []*model.Bank
	Bankes, total, err = s.RepositoryBank.Get(ctx, req)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, Bank := range Bankes {
		res = append(res, dto.BankResponse{
			ID:              Bank.ID,
			Code:            Bank.Code,
			Description:     Bank.Description,
			Value:           Bank.Value,
			ImageUrl:        Bank.ImageUrl,
			PaymentGuideUrl: Bank.PaymentGuideUrl,
			PublishIVA:      Bank.PublishIVA,
			PublishFVA:      Bank.PublishFVA,
			Status:          Bank.Status,
			StatusConvert:   statusx.ConvertStatusValue(Bank.Status),
			CreatedAt:       timex.ToLocTime(ctx, Bank.CreatedAt),
			UpdatedAt:       timex.ToLocTime(ctx, Bank.UpdatedAt),
		})
	}

	return
}

func (s *BankService) GetDetail(ctx context.Context, req *pb.GetBankDetailRequest) (res dto.BankResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "BankService.GetDetail")
	defer span.End()

	var Bank *model.Bank
	Bank, err = s.RepositoryBank.GetDetail(ctx, req)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.BankResponse{
		ID:              Bank.ID,
		Code:            Bank.Code,
		Description:     Bank.Description,
		Value:           Bank.Value,
		ImageUrl:        Bank.ImageUrl,
		PaymentGuideUrl: Bank.PaymentGuideUrl,
		PublishIVA:      Bank.PublishIVA,
		PublishFVA:      Bank.PublishFVA,
		Status:          Bank.Status,
		StatusConvert:   statusx.ConvertStatusValue(Bank.Status),
		CreatedAt:       timex.ToLocTime(ctx, Bank.CreatedAt),
		UpdatedAt:       timex.ToLocTime(ctx, Bank.UpdatedAt),
	}

	return
}
