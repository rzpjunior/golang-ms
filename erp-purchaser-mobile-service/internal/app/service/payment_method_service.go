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

type IPaymentMethodService interface {
	Get(ctx context.Context, req dto.PaymentMethodListRequest) (res []*dto.PaymentMethodResponse, err error)
	GetByID(ctx context.Context, id string) (res *dto.PaymentMethodResponse, err error)
}

type PaymentMethodService struct {
	opt opt.Options
}

func NewPaymentMethodService() IPaymentMethodService {
	return &PaymentMethodService{
		opt: global.Setup.Common,
	}
}

func (s *PaymentMethodService) Get(ctx context.Context, req dto.PaymentMethodListRequest) (res []*dto.PaymentMethodResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PaymentMethodService.Get")
	defer span.End()

	var paymentMethods *bridgeService.GetPaymentMethodGPResponse
	paymentMethods, err = s.opt.Client.BridgeServiceGrpc.GetPaymentMethodGPList(ctx, &bridge_service.GetPaymentMethodGPListRequest{
		Limit:  req.Limit,
		Offset: req.Offset,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "payment_method")
		return
	}

	for _, paymentMethod := range paymentMethods.Data {
		res = append(res, &dto.PaymentMethodResponse{
			ID:   paymentMethod.PrpPaymentMethodId,
			Code: paymentMethod.PrpPaymentMethodId,
			Name: paymentMethod.PrpPaymentMethodDesc,
		})
	}

	return
}

func (s *PaymentMethodService) GetByID(ctx context.Context, id string) (res *dto.PaymentMethodResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PaymentMethodService.GetByID")
	defer span.End()

	var paymentMethod *bridgeService.GetPaymentMethodGPResponse
	paymentMethod, err = s.opt.Client.BridgeServiceGrpc.GetPaymentMethodGPDetail(ctx, &bridge_service.GetPaymentMethodGPDetailRequest{
		Id: id,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "payment_method")
		return
	}

	res = &dto.PaymentMethodResponse{
		ID:   paymentMethod.Data[0].PrpPaymentMethodId,
		Code: paymentMethod.Data[0].PrpPaymentMethodId,
		Name: paymentMethod.Data[0].PrpPaymentMethodDesc,
	}

	return
}
