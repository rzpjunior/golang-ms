package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/global"

	dto "git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/internal/app/dto"
)

type IPaymentTermService interface {
	Get(ctx context.Context, req dto.PaymentTermListRequest) (res []*dto.PaymentTermResponse, err error)
	GetByID(ctx context.Context, id string) (res *dto.PaymentTermResponse, err error)
}

type PaymentTermService struct {
	opt opt.Options
}

func NewPaymentTermService() IPaymentTermService {
	return &PaymentTermService{
		opt: global.Setup.Common,
	}
}

func (s *PaymentTermService) Get(ctx context.Context, req dto.PaymentTermListRequest) (res []*dto.PaymentTermResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PaymentTermService.Get")
	defer span.End()

	var paymentTerms *bridgeService.GetPaymentTermGPResponse
	paymentTerms, err = s.opt.Client.BridgeServiceGrpc.GetPaymentTermGPList(ctx, &bridge_service.GetPaymentTermGPListRequest{
		Limit:         req.Limit,
		Offset:        req.Offset,
		PaymentUsefor: "1",
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "payment_term")
		return
	}

	for _, paymentTerm := range paymentTerms.Data {
		res = append(res, &dto.PaymentTermResponse{
			Id:          paymentTerm.Pymtrmid,
			Code:        paymentTerm.Pymtrmid,
			Description: paymentTerm.Pymtrmid,
		})
	}

	return
}

func (s *PaymentTermService) GetByID(ctx context.Context, id string) (res *dto.PaymentTermResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PaymentTermService.GetByID")
	defer span.End()

	var paymentTerm *bridgeService.GetPaymentTermGPResponse
	paymentTerm, err = s.opt.Client.BridgeServiceGrpc.GetPaymentTermGPDetail(ctx, &bridge_service.GetPaymentTermGPDetailRequest{
		Id: id,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "payment_method")
		return
	}

	res = &dto.PaymentTermResponse{
		Id:          paymentTerm.Data[0].Pymtrmid,
		Code:        paymentTerm.Data[0].Pymtrmid,
		Description: paymentTerm.Data[0].Pymtrmid,
	}

	return
}
