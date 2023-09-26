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

type ICourierService interface {
	Get(ctx context.Context, req dto.GetCourierRequest) (res []dto.CourierResponse, total int64, err error)
	GetDetail(ctx context.Context, id string) (res dto.CourierResponse, err error)
}

type CourierService struct {
	opt opt.Options
}

func NewServiceCourier() ICourierService {
	return &CourierService{
		opt: global.Setup.Common,
	}
}

func (s *CourierService) Get(ctx context.Context, req dto.GetCourierRequest) (res []dto.CourierResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CourierService.Get")
	defer span.End()

	var courier *bridgeService.GetCourierGPResponse

	if courier, err = s.opt.Client.BridgeServiceGrpc.GetCourierGPList(ctx, &bridgeService.GetCourierGPListRequest{
		Limit:              int32(req.Limit),
		Offset:             int32(req.Offset),
		GnlCourierName:     req.Name,
		GnlCourierVendorId: req.CourierVendorID,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "courier")
		return
	}

	for _, courier := range courier.Data {
		res = append(res, dto.CourierResponse{
			ID:                courier.GnlCourierId,
			Name:              courier.GnlCourierName,
			PhoneNumber:       courier.Phonname,
			VehicleProfileId:  courier.GnlVehicleProfileId,
			LicensePlate:      courier.GnlLicensePlate,
			EmergencyMode:     2,
			LastEmergencyTime: nil,
			Status:            courier.Inactive,
		})
	}

	total = int64(len(courier.Data))

	return
}

func (s *CourierService) GetDetail(ctx context.Context, id string) (res dto.CourierResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CourierService.GetDetail")
	defer span.End()

	var courier *bridgeService.GetCourierGPResponse

	if courier, err = s.opt.Client.BridgeServiceGrpc.GetCourierGPDetail(ctx, &bridgeService.GetCourierGPDetailRequest{
		Id: id,
	}); err != nil || !courier.Succeeded {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "courier")
		return
	}

	res = dto.CourierResponse{
		ID:                courier.Data[0].GnlCourierId,
		Name:              courier.Data[0].GnlCourierName,
		PhoneNumber:       courier.Data[0].Phonname,
		VehicleProfileId:  courier.Data[0].GnlVehicleProfileId,
		LicensePlate:      courier.Data[0].GnlLicensePlate,
		EmergencyMode:     2,
		LastEmergencyTime: nil,
		Status:            courier.Data[0].Inactive,
	}

	return
}
