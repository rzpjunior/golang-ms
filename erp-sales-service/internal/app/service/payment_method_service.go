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

type IPaymentMethodService interface {
	GetListGRPC(ctx context.Context, req *pb.GetPaymentMethodListRequest) (res []*model.PaymentMethod, err error)
}

type PaymentMethodService struct {
	opt                     opt.Options
	RepositoryPaymentMethod repository.IPaymentMethodRepository
}

func NewPaymentMethodService() IPaymentMethodService {
	return &PaymentMethodService{
		opt:                     global.Setup.Common,
		RepositoryPaymentMethod: repository.NewPaymentMethodRepository(),
	}
}

func (s *PaymentMethodService) GetListGRPC(ctx context.Context, req *pb.GetPaymentMethodListRequest) (res []*model.PaymentMethod, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PaymentMethodService.Get")
	defer span.End()

	res,err = s.RepositoryPaymentMethod.GetPaymentMethodList(ctx, req)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}
