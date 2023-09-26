package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/sales_service"
	"git.edenfarm.id/project-version3/erp-services/erp-sales-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-sales-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-sales-service/internal/app/repository"
)

type IPaymentGroupCombService interface {
	GetListGRPC(ctx context.Context, req *pb.GetPaymentGroupCombListRequest) (res []*model.PaymentGroupComb,  err error)
}

type PaymentGroupCombService struct {
	opt                      opt.Options
	RepositoryPaymentGroupComb repository.IPaymentGroupCombRepository
}

func NewPaymentGroupCombService() IPaymentGroupCombService {
	return &PaymentGroupCombService{
		opt:                      global.Setup.Common,
		RepositoryPaymentGroupComb: repository.NewPaymentGroupCombRepository(),
	}
}

func (s *PaymentGroupCombService) GetListGRPC(ctx context.Context, req *pb.GetPaymentGroupCombListRequest) (res []*model.PaymentGroupComb,  err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PaymentGroupCombService.Get")
	defer span.End()

	res, err = s.RepositoryPaymentGroupComb.GetPaymentGroupCombList(ctx, req)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}
