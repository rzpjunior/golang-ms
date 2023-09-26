package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"git.edenfarm.id/project-version3/erp-services/erp-sales-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-sales-service/internal/app/dto"
)

type IPaymentTermService interface {
	Get(ctx context.Context, req *dto.GetPaymentTermRequest) (res []*dto.PaymentTermResponse, total int64, err error)
}

type PaymentTermService struct {
	opt opt.Options
}

func NewServicePaymentTerm() IPaymentTermService {
	return &PaymentTermService{
		opt: global.Setup.Common,
	}
}

func (s *PaymentTermService) Get(ctx context.Context, req *dto.GetPaymentTermRequest) (res []*dto.PaymentTermResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PaymentTermService.Get")
	defer span.End()

	var paymentTerms *bridgeService.GetPaymentTermGPResponse

	if paymentTerms, err = s.opt.Client.BridgeServiceGrpc.GetPaymentTermGPList(ctx, &bridgeService.GetPaymentTermGPListRequest{
		Limit:         int32(req.Limit),
		Offset:        int32(req.Offset),
		PaymentUsefor: req.PaymentUseFor,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "payment term")
		return
	}

	for _, paymentTerm := range paymentTerms.Data {
		res = append(res, &dto.PaymentTermResponse{
			ID:                       paymentTerm.Pymtrmid,
			Description:              paymentTerm.Pymtrmid,
			DueType:                  paymentTerm.DuetypeDesc,
			PaymentUseFor:            int(paymentTerm.GnlPaymentUsefor),
			PaymentUseForDescription: paymentTerm.GnlPaymentUseforDesc,
		})
	}

	total = int64(paymentTerms.TotalRecords)

	return
}
