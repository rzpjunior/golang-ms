package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-inventory-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-inventory-service/internal/app/dto"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

type IUomService interface {
	Get(ctx context.Context, req dto.GetUomRequest) (res []dto.UomGPResponse, total int64, err error)
	GetDetail(ctx context.Context, id string) (res dto.UomGPResponse, err error)
}

type UomService struct {
	opt opt.Options
}

func NewServiceUom() IUomService {
	return &UomService{
		opt: global.Setup.Common,
	}
}

func (s *UomService) Get(ctx context.Context, req dto.GetUomRequest) (res []dto.UomGPResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "UomService.Get")
	defer span.End()

	var uom *bridgeService.GetUomGPResponse

	if uom, err = s.opt.Client.BridgeServiceGrpc.GetUomGPList(ctx, &bridgeService.GetUomGPListRequest{
		Limit:  int32(req.Limit),
		Offset: int32(req.Offset),
		Search: req.Search,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "Uom")
		return
	}

	for _, uom := range uom.Data {
		res = append(res, dto.UomGPResponse{
			ID:   uom.Uomschdl,
			Name: uom.Umschdsc,
		})
	}

	total = int64(len(uom.Data))

	return
}

func (s *UomService) GetDetail(ctx context.Context, id string) (res dto.UomGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "UomService.GetUom")
	defer span.End()

	var uom *bridgeService.GetUomGPResponse

	if uom, err = s.opt.Client.BridgeServiceGrpc.GetUomGPDetail(ctx, &bridgeService.GetUomGPDetailRequest{
		Id: id,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "Uom")
		return
	}

	res = dto.UomGPResponse{
		ID:   uom.Data[0].Uomschdl,
		Name: uom.Data[0].Umschdsc,
	}

	return
}
