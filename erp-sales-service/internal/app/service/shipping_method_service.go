package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/utils"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"git.edenfarm.id/project-version3/erp-services/erp-sales-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-sales-service/internal/app/dto"
)

type IShippingMethodService interface {
	Get(ctx context.Context, req *dto.GetShippingMethodRequest) (res []*dto.ShippingMethodResponse, total int64, err error)
}

type ShippingMethodService struct {
	opt opt.Options
}

func NewServiceShippingMethod() IShippingMethodService {
	return &ShippingMethodService{
		opt: global.Setup.Common,
	}
}

func (s *ShippingMethodService) Get(ctx context.Context, req *dto.GetShippingMethodRequest) (res []*dto.ShippingMethodResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ShippingMethodService.Get")
	defer span.End()

	var (
		shippingMethods *bridgeService.GetShippingMethodResponse
		typeGP          string
	)

	if req.Type != "" {
		typeGP = utils.ToString(req.Type)
	}

	if shippingMethods, err = s.opt.Client.BridgeServiceGrpc.GetShippingMethodList(ctx, &bridgeService.GetShippingMethodListRequest{
		Limit:    req.Limit,
		Offset:   req.Offset,
		Shiptype: typeGP,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "sales Person")
		return
	}

	for _, shippingMethod := range shippingMethods.Data {
		res = append(res, &dto.ShippingMethodResponse{
			ID:              shippingMethod.Shipmthd,
			Description:     shippingMethod.Shmthdsc,
			Type:            int8(shippingMethod.Shiptype),
			TypeDescription: shippingMethod.ShiptypeDesc,
		})
	}

	total = int64(shippingMethods.TotalRecords)

	return
}
