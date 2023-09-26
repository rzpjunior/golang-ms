package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-logistic-service/internal/app/dto"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

type ICourierVendorService interface {
	Get(ctx context.Context, req dto.GetCourierVendorRequest) (res []dto.CourierVendorResponse, total int64, err error)
	GetDetail(ctx context.Context, id string) (res dto.CourierVendorResponse, err error)
}

type CourierVendorService struct {
	opt opt.Options
}

func NewServiceCourierVendor() ICourierVendorService {
	return &CourierVendorService{
		opt: global.Setup.Common,
	}
}

func (s *CourierVendorService) Get(ctx context.Context, req dto.GetCourierVendorRequest) (res []dto.CourierVendorResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CourierVendorService.Get")
	defer span.End()

	var courierVendor *bridgeService.GetCourierVendorGPResponse

	if courierVendor, err = s.opt.Client.BridgeServiceGrpc.GetCourierVendorGPList(ctx, &bridgeService.GetCourierVendorGPListRequest{
		Limit:                int32(req.Limit),
		Offset:               int32(req.Offset),
		Locncode:             req.SiteId,
		GnlCourierVendorName: req.CourierVendorName,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "courier vendor")
		return
	}

	for _, courierVendor := range courierVendor.Data {
		res = append(res, dto.CourierVendorResponse{
			ID:     courierVendor.GnlCourierVendorId,
			Name:   courierVendor.GnlCourierVendorName,
			SiteId: courierVendor.Locncode,
			Status: courierVendor.Inactive,
		})
	}

	total = int64(len(courierVendor.Data))

	return
}

func (s *CourierVendorService) GetDetail(ctx context.Context, id string) (res dto.CourierVendorResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CourierVendorService.GetDetail")
	defer span.End()

	var courierVendor *bridgeService.GetCourierVendorGPResponse

	if courierVendor, err = s.opt.Client.BridgeServiceGrpc.GetCourierVendorGPDetail(ctx, &bridgeService.GetCourierVendorGPDetailRequest{
		Id: id,
	}); err != nil || !courierVendor.Succeeded {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "courier vendor")
		return
	}

	res = dto.CourierVendorResponse{
		ID:     courierVendor.Data[0].GnlCourierVendorId,
		Name:   courierVendor.Data[0].GnlCourierVendorName,
		SiteId: courierVendor.Data[0].Locncode,
		Status: courierVendor.Data[0].Inactive,
	}

	return
}
