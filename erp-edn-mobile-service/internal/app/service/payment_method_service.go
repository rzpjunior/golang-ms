package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/sales_service"
)

func NewServicePaymentMethod() IPaymentMethodService {
	m := new(PaymentMethodService)
	m.opt = global.Setup.Common
	return m
}

type IPaymentMethodService interface {
	GetGP(ctx context.Context, req dto.PaymentMethodListRequest) (res []*dto.PaymentMethodGP, total int64, err error)
	Get(ctx context.Context, req dto.PaymentMethodListRequest) (res []*dto.PaymentMethodGP, total int64, err error)
	GetDetaiGPlById(ctx context.Context, id string) (res *dto.PaymentMethodGP, err error)
}

type PaymentMethodService struct {
	opt opt.Options
}

func (s *PaymentMethodService) GetGP(ctx context.Context, req dto.PaymentMethodListRequest) (res []*dto.PaymentMethodGP, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PaymentMethodService.GetGP")
	defer span.End()

	// get payment method from bridge
	// var pmRes *bridgeService.GetPaymentMethodGPResponse
	// pmRes, err = s.opt.Client.BridgeServiceGrpc.GetPaymentMethodGPList(ctx, &bridgeService.GetPaymentMethodGPListRequest{
	// 	Limit:  req.Limit,
	// 	Offset: req.Offset,
	// })
	// if err != nil {
	// 	span.RecordError(err)
	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 	err = edenlabs.ErrorRpcNotFound("bridge", "payment method")
	// 	return
	// }

	// datas := []*dto.PaymentMethodGP{}
	// for _, pm := range pmRes.Data {
	// 	datas = append(datas, &dto.PaymentMethodGP{
	// 		Pmntnmbr:             pm.Pmntnmbr,
	// 		Vendorid:             pm.Vendorid,
	// 		Docnumbr:             pm.Docnumbr,
	// 		Docdate:              pm.Docdate,
	// 		Chekbkid:             pm.Chekbkid,
	// 		Pyenttyp:             pm.Pyenttyp,
	// 		PaymentMethodDesc:    pm.PaymentMethodDesc,
	// 		PrpPaymentMethodId:   pm.PrpPaymentMethodId,
	// 		PrpPaymentMethodDesc: pm.PrpPaymentMethodDesc,
	// 		Inactive:             pm.Inactive,
	// 		InactiveDesc:         pm.InactiveDesc,
	// 	})
	// }
	// paymentChannel, err := s.opt.Client.SalesServiceGrpc.GetPaymentChannelList(ctx, &sales_service.GetPaymentChannelListRequest{
	// 	Status:     1,
	// 	PublishIva: 1,
	// 	// PublishFva:      1,
	// 	PaymentMethodId: 2,
	// })

	paymentMethod, err := s.opt.Client.SalesServiceGrpc.GetPaymentMethodList(ctx, &sales_service.GetPaymentMethodListRequest{
		Status: 1,
		Search: req.Search,
	})

	for _, v := range paymentMethod.Data {
		res = append(res, &dto.PaymentMethodGP{
			ID:                v.Code,
			PaymentMethodCode: v.Code,
			PaymentMethodDesc: v.Name,
		})
	}

	total = int64(len(res))

	return
}

func (s *PaymentMethodService) GetDetaiGPlById(ctx context.Context, id string) (res *dto.PaymentMethodGP, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PaymentMethodService.GetDetaiGPlById")
	defer span.End()

	// get paymentMethod from bridge
	paymentMethod, err := s.opt.Client.SalesServiceGrpc.GetPaymentMethodList(ctx, &sales_service.GetPaymentMethodListRequest{
		Status: 1,
		Id:     id,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "payment method")
		return
	}

	if len(paymentMethod.Data) > 0 {
		res = &dto.PaymentMethodGP{
			ID:                paymentMethod.Data[0].Code,
			PaymentMethodDesc: paymentMethod.Data[0].Name,
			PaymentMethodCode: paymentMethod.Data[0].Code,
		}
	}

	return
}

func (s *PaymentMethodService) Get(ctx context.Context, req dto.PaymentMethodListRequest) (res []*dto.PaymentMethodGP, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PaymentMethodService.GetGP")
	defer span.End()

	// paymentMethod, _ := s.opt.Client.ConfigurationServiceGrpc.GetConfigAppList(ctx, &configuration_service.GetConfigAppListRequest{
	// 	Offset:    0,
	// 	Limit:     100,
	// 	Attribute: req.Search,
	// 	Field:     "EDN Payment Method",
	// })

	// for _, v := range paymentMethod.Data {
	// 	res = append(res, &dto.PaymentMethodGP{
	// 		ID:                v.Value,
	// 		PaymentMethodCode: v.Attribute,
	// 		PaymentMethodDesc: v.Attribute,
	// 	})

	// }

	res = append(res, &dto.PaymentMethodGP{
		ID:                "0",
		PaymentMethodCode: "Check",
		PaymentMethodDesc: "Check",
	})
	res = append(res, &dto.PaymentMethodGP{
		ID:                "1",
		PaymentMethodCode: "Cash",
		PaymentMethodDesc: "Cash",
	})
	total = int64(len(res))

	return
}
