package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *BridgeGrpcHandler) GetCourierList(ctx context.Context, req *bridgeService.GetCourierListRequest) (res *bridgeService.GetCourierListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetCourierList")
	defer span.End()

	var couriers []dto.CourierResponse
	couriers, _, err = h.ServicesCourier.Get(ctx, int(req.Offset), int(req.Limit), int(req.Status), req.Search, req.OrderBy, req.VehicleProfileId, req.EmergencyMode)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*bridgeService.Courier
	for _, courier := range couriers {
		data = append(data, &bridgeService.Courier{
			Id:                courier.ID,
			RoleId:            courier.RoleID,
			UserId:            courier.UserID,
			Code:              courier.Code,
			Name:              courier.Name,
			PhoneNumber:       courier.PhoneNumber,
			VehicleProfileId:  courier.VehicleProfileID,
			LicensePlate:      courier.LicensePlate,
			EmergencyMode:     int32(courier.EmergencyMode),
			LastEmergencyTime: timestamppb.New(courier.LastEmergencyTime),
			Status:            int32(courier.Status),
		})
	}

	res = &bridgeService.GetCourierListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}

	return
}

func (h *BridgeGrpcHandler) GetCourierDetail(ctx context.Context, req *bridgeService.GetCourierDetailRequest) (res *bridgeService.GetCourierDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetCourierDetail")
	defer span.End()

	var courier dto.CourierResponse
	courier, err = h.ServicesCourier.GetDetail(ctx, req.Id, req.Code, req.UserId)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.GetCourierDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &bridgeService.Courier{
			Id:                courier.ID,
			RoleId:            courier.RoleID,
			UserId:            courier.UserID,
			Code:              courier.Code,
			Name:              courier.Name,
			PhoneNumber:       courier.PhoneNumber,
			VehicleProfileId:  courier.VehicleProfileID,
			LicensePlate:      courier.LicensePlate,
			EmergencyMode:     int32(courier.EmergencyMode),
			LastEmergencyTime: timestamppb.New(courier.LastEmergencyTime),
			Status:            int32(courier.Status),
		},
	}

	return
}

func (h *BridgeGrpcHandler) GetCourierGPList(ctx context.Context, req *bridgeService.GetCourierGPListRequest) (res *bridgeService.GetCourierGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetCourierGPList")
	defer span.End()

	res, err = h.ServicesCourier.GetGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *BridgeGrpcHandler) GetCourierGPDetail(ctx context.Context, req *bridgeService.GetCourierGPDetailRequest) (res *bridgeService.GetCourierGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetCourierGPDetail")
	defer span.End()

	res, err = h.ServicesCourier.GetDetailGP(ctx, req)

	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *BridgeGrpcHandler) ActivateEmergencyCourier(ctx context.Context, req *bridgeService.EmergencyCourierRequest) (res *bridgeService.EmergencyCourierResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.ActivateEmergencyCourier")
	defer span.End()

	_, err = h.ServicesCourier.ActivateEmergency(ctx, req.Id)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.EmergencyCourierResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
	}
	return
}

func (h *BridgeGrpcHandler) DeactivateEmergencyCourier(ctx context.Context, req *bridgeService.EmergencyCourierRequest) (res *bridgeService.EmergencyCourierResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.DeactivateEmergencyCourier")
	defer span.End()

	_, err = h.ServicesCourier.DeactivateEmergency(ctx, req.Id)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.EmergencyCourierResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
	}

	return
}
