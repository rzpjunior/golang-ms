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

type IPaymentChannelService interface {
	GetListGRPC(ctx context.Context, req *pb.GetPaymentChannelListRequest) (res []*model.PaymentChannel,  err error)
}

type PaymentChannelService struct {
	opt                      opt.Options
	RepositoryPaymentChannel repository.IPaymentChannelRepository
}

func NewPaymentChannelService() IPaymentChannelService {
	return &PaymentChannelService{
		opt:                      global.Setup.Common,
		RepositoryPaymentChannel: repository.NewPaymentChannelRepository(),
	}
}

func (s *PaymentChannelService) GetListGRPC(ctx context.Context, req *pb.GetPaymentChannelListRequest) (res []*model.PaymentChannel,  err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PaymentChannelService.Get")
	defer span.End()

	res, err = s.RepositoryPaymentChannel.GetPaymentChannelList(ctx, req)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}
