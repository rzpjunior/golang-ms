package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/utils"
	"git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

func NewServicePaymentTerm() IPaymentTermService {
	m := new(PaymentTermService)
	m.opt = global.Setup.Common
	return m
}

type IPaymentTermService interface {
	GetGP(ctx context.Context, req dto.PaymentTermListRequest) (res []*dto.PaymentTermGP, total int64, err error)
}

type PaymentTermService struct {
	opt opt.Options
}

func (s *PaymentTermService) GetGP(ctx context.Context, req dto.PaymentTermListRequest) (res []*dto.PaymentTermGP, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PaymentTermService.GetGP")
	defer span.End()

	paymentTerm, err := s.opt.Client.BridgeServiceGrpc.GetPaymentTermGPList(ctx, &bridge_service.GetPaymentTermGPListRequest{
		PaymentUsefor: "2",
		PaymentTermId: req.Search,
		Limit:         req.Limit,
		Offset:        req.Offset,
	})

	for _, v := range paymentTerm.Data {
		if v.Pymtrmid == "PBD" {
			continue
		}
		if v.Pymtrmid == "PWD" {
			continue
		}
		if v.Pymtrmid == "BNS" {
			continue
		}
		res = append(res, &dto.PaymentTermGP{
			ID:                v.Pymtrmid,
			PaymentTermCode:   v.Pymtrmid,
			CalculateFromDays: utils.ToString(v.CalculateDateFromDays),
		})
	}

	total = int64(len(res))

	return
}
