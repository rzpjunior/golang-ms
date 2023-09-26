package service

import (
	"context"
	"strconv"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

type IPaymentMethodService interface {
	GetGP(ctx context.Context, req *pb.GetPaymentMethodGPListRequest) (res *pb.GetPaymentMethodGPResponse, err error)
	GetDetailGP(ctx context.Context, req *pb.GetPaymentMethodGPDetailRequest) (res *pb.GetPaymentMethodGPResponse, err error)
}

type PaymentMethodService struct {
	opt opt.Options
}

func NewPaymentMethodService() IPaymentMethodService {
	return &PaymentMethodService{
		opt: global.Setup.Common,
	}
}

func (s *PaymentMethodService) GetGP(ctx context.Context, req *pb.GetPaymentMethodGPListRequest) (res *pb.GetPaymentMethodGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PaymentMethodService.GetGP")
	defer span.End()

	params := map[string]string{
		"interid":    global.EnvDatabaseGP,
		"PageNumber": strconv.Itoa(int(req.Offset)),
		"PageSize":   strconv.Itoa(int(req.Limit)),
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "paymentmethod/getall", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *PaymentMethodService) GetDetailGP(ctx context.Context, req *pb.GetPaymentMethodGPDetailRequest) (res *pb.GetPaymentMethodGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PaymentMethodService.GetDetailGP")
	defer span.End()

	params := map[string]string{
		"interid": global.EnvDatabaseGP,
		"id":      req.Id,
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "paymentmethod/getbyid", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}
