package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *BridgeGrpcHandler) GetVehicleProfileList(ctx context.Context, req *bridgeService.GetVehicleProfileListRequest) (res *bridgeService.GetVehicleProfileListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetVehicleProfileList")
	defer span.End()

	var vehicleProfiles []dto.VehicleProfileResponse
	vehicleProfiles, _, err = h.ServicesVehicleProfile.Get(ctx, int(req.Offset), int(req.Limit), int(req.Status), req.Search, req.OrderBy, req.CourierVendorId)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*bridgeService.VehicleProfile
	for _, vProfile := range vehicleProfiles {
		data = append(data, &bridgeService.VehicleProfile{
			Id:                  vProfile.ID,
			CourierVendorId:     vProfile.CourierVendorID,
			Code:                vProfile.Code,
			Name:                vProfile.Name,
			MaxKoli:             vProfile.MaxKoli,
			MaxWeight:           vProfile.MaxWeight,
			MaxFragile:          vProfile.MaxFragile,
			SpeedFactor:         vProfile.SpeedFactor,
			RoutingProfile:      int32(vProfile.RoutingProfile),
			Status:              int32(vProfile.Status),
			Skills:              vProfile.Skills,
			InitialCost:         vProfile.InitialCost,
			SubsequentCost:      vProfile.SubsequentCost,
			MaxAvailableVehicle: vProfile.MaxAvailableVehicle,
		})
	}

	res = &bridgeService.GetVehicleProfileListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}

	return
}

func (h *BridgeGrpcHandler) GetVehicleProfileDetail(ctx context.Context, req *bridgeService.GetVehicleProfileDetailRequest) (res *bridgeService.GetVehicleProfileDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetVehicleProfileDetail")
	defer span.End()

	var vehicleProfile dto.VehicleProfileResponse
	vehicleProfile, err = h.ServicesVehicleProfile.GetDetail(ctx, req.Id, req.Code)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.GetVehicleProfileDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &bridgeService.VehicleProfile{
			Id:                  vehicleProfile.ID,
			CourierVendorId:     vehicleProfile.CourierVendorID,
			Code:                vehicleProfile.Code,
			Name:                vehicleProfile.Name,
			MaxKoli:             vehicleProfile.MaxKoli,
			MaxWeight:           vehicleProfile.MaxWeight,
			MaxFragile:          vehicleProfile.MaxFragile,
			SpeedFactor:         vehicleProfile.SpeedFactor,
			RoutingProfile:      int32(vehicleProfile.RoutingProfile),
			Status:              int32(vehicleProfile.Status),
			Skills:              vehicleProfile.Skills,
			InitialCost:         vehicleProfile.InitialCost,
			SubsequentCost:      vehicleProfile.SubsequentCost,
			MaxAvailableVehicle: vehicleProfile.MaxAvailableVehicle,
		},
	}

	return
}

func (h *BridgeGrpcHandler) GetVehicleProfileGPList(ctx context.Context, req *bridgeService.GetVehicleProfileGPListRequest) (res *bridgeService.GetVehicleProfileGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetVehicleProfileGPList")
	defer span.End()

	res, err = h.ServicesVehicleProfile.GetGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *BridgeGrpcHandler) GetVehicleProfileGPDetail(ctx context.Context, req *bridgeService.GetVehicleProfileGPDetailRequest) (res *bridgeService.GetVehicleProfileGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetVehicleProfileGPDetail")
	defer span.End()

	res, err = h.ServicesVehicleProfile.GetDetailGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}
