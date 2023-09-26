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

type IItemClassService interface {
	Get(ctx context.Context, req dto.GetItemClassRequest) (res []dto.ItemClassResponse, total int64, err error)
	GetDetail(ctx context.Context, id string) (res dto.ItemClassResponse, err error)
}

type ItemClassService struct {
	opt opt.Options
}

func NewServiceItemClass() IItemClassService {
	return &ItemClassService{
		opt: global.Setup.Common,
	}
}

func (s *ItemClassService) Get(ctx context.Context, req dto.GetItemClassRequest) (res []dto.ItemClassResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemClassService.Get")
	defer span.End()

	var itemClass *bridgeService.GetItemClassGPResponse

	if itemClass, err = s.opt.Client.BridgeServiceGrpc.GetItemClassGPList(ctx, &bridgeService.GetItemClassGPListRequest{
		Limit:  int32(req.Limit),
		Offset: int32(req.Offset),
		Search: req.Search,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "ItemClass")
		return
	}

	for _, itemClass := range itemClass.Data {
		res = append(res, dto.ItemClassResponse{
			ID:   itemClass.Itmclscd,
			Name: itemClass.Itmclsdc,
		})
	}

	total = int64(len(itemClass.Data))

	return
}

func (s *ItemClassService) GetDetail(ctx context.Context, id string) (res dto.ItemClassResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemClassService.GetItemClass")
	defer span.End()

	var itemClass *bridgeService.GetItemClassGPResponse

	if itemClass, err = s.opt.Client.BridgeServiceGrpc.GetItemClassGPDetail(ctx, &bridgeService.GetItemClassGPDetailRequest{
		Id: id,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "ItemClass")
		return
	}

	res = dto.ItemClassResponse{
		ID:   itemClass.Data[0].Itmclscd,
		Name: itemClass.Data[0].Itmclsdc,
	}

	return
}
