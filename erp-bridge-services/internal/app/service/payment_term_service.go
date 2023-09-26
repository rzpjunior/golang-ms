package service

import (
	"context"
	"net/url"
	"strconv"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

type IPaymentTermService interface {
	GetGP(ctx context.Context, req *pb.GetPaymentTermGPListRequest) (res *pb.GetPaymentTermGPResponse, err error)
	GetDetailGP(ctx context.Context, req *pb.GetPaymentTermGPDetailRequest) (res *pb.GetPaymentTermGPResponse, err error)
}

type PaymentTermService struct {
	opt opt.Options
}

func NewPaymentTermService() IPaymentTermService {
	return &PaymentTermService{
		opt: global.Setup.Common,
	}
}

func (s *PaymentTermService) GetGP(ctx context.Context, req *pb.GetPaymentTermGPListRequest) (res *pb.GetPaymentTermGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PaymentTermService.GetGP")
	defer span.End()

	params := map[string]string{
		"interid":    global.EnvDatabaseGP,
		"PageNumber": strconv.Itoa(int(req.Offset)),
		"PageSize":   strconv.Itoa(int(req.Limit)),
	}

	if req.PaymentUsefor != "" {
		params["gnl_payment_usefor"] = req.PaymentUsefor
	}

	if req.PaymentTermId != "" {
		req.PaymentTermId = url.PathEscape(req.PaymentTermId)
		params["payment_term_id"] = req.PaymentTermId
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "paymentterm/getall", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *PaymentTermService) GetDetailGP(ctx context.Context, req *pb.GetPaymentTermGPDetailRequest) (res *pb.GetPaymentTermGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PaymentTermService.GetDetailGP")
	defer span.End()

	req.Id = url.PathEscape(req.Id)
	params := map[string]string{
		"interid": global.EnvDatabaseGP,
		"id":      req.Id,
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "paymentterm/getbyid", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}
