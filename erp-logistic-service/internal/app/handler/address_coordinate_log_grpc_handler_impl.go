package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/internal/app/dto"
	logisticService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/logistic_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *LogisticGrpcHandler) GetAddressCoordinateLogList(ctx context.Context, req *logisticService.GetAddressCoordinateLogListRequest) (res *logisticService.GetAddressCoordinateLogListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetAddressCoordinateLogList")
	defer span.End()

	addressCoordinateLogs, _, err := h.ServicesAddressCoordinateLog.Get(ctx, dto.AddressCoordinateLogGetRequest{
		OrderBy:          req.OrderBy,
		GroupBy:          req.GroupBy,
		ArrAddressIDs:    req.AddressId,
		ArrSalesOrderIDs: req.SalesOrderId,
	})
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*logisticService.AddressCoordinateLog
	for _, address := range addressCoordinateLogs {
		data = append(data, &logisticService.AddressCoordinateLog{
			Id:             address.ID,
			AddressId:      address.AddressID,
			SalesOrderId:   address.SalesOrderID,
			Latitude:       &address.Latitude,
			Longitude:      &address.Longitude,
			LogChannelId:   int32(address.LogChannelID),
			MainCoordinate: int32(address.MainCoordinate),
			CreatedAt:      timestamppb.New(address.CreatedAt),
			CreatedBy:      address.CreatedBy,
		})
	}

	res = &logisticService.GetAddressCoordinateLogListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}

	return
}

func (h *LogisticGrpcHandler) GetAddressCoordinateLogDetail(ctx context.Context, req *logisticService.GetAddressCoordinateLogDetailRequest) (res *logisticService.GetAddressCoordinateLogDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetAddressCoordinateLogDetail")
	defer span.End()

	addressCoordinateLog, err := h.ServicesAddressCoordinateLog.GetDetail(ctx, req.Id, 0)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &logisticService.GetAddressCoordinateLogDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &logisticService.AddressCoordinateLog{
			Id:             addressCoordinateLog.ID,
			AddressId:      addressCoordinateLog.AddressID,
			SalesOrderId:   addressCoordinateLog.SalesOrderID,
			Latitude:       &addressCoordinateLog.Latitude,
			Longitude:      &addressCoordinateLog.Longitude,
			LogChannelId:   int32(addressCoordinateLog.LogChannelID),
			MainCoordinate: int32(addressCoordinateLog.MainCoordinate),
			CreatedAt:      timestamppb.New(addressCoordinateLog.CreatedAt),
			CreatedBy:      addressCoordinateLog.CreatedBy,
		},
	}

	return
}

func (h *LogisticGrpcHandler) CreateAddressCoordinateLog(ctx context.Context, req *logisticService.CreateAddressCoordinateLogRequest) (res *logisticService.CreateAddressCoordinateLogResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.CreateAddressCoordinateLog")
	defer span.End()

	addressCoordinateLog, err := h.ServicesAddressCoordinateLog.Create(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &logisticService.CreateAddressCoordinateLogResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &logisticService.AddressCoordinateLog{
			Id:             addressCoordinateLog.ID,
			AddressId:      addressCoordinateLog.AddressID,
			SalesOrderId:   addressCoordinateLog.SalesOrderID,
			Latitude:       &addressCoordinateLog.Latitude,
			Longitude:      &addressCoordinateLog.Longitude,
			LogChannelId:   int32(addressCoordinateLog.LogChannelID),
			MainCoordinate: int32(addressCoordinateLog.MainCoordinate),
			CreatedAt:      timestamppb.New(addressCoordinateLog.CreatedAt),
			CreatedBy:      addressCoordinateLog.CreatedBy,
		},
	}

	return
}

func (h *LogisticGrpcHandler) GetMostTrustedAddressCoordinateLog(ctx context.Context, req *logisticService.GetMostTrustedAddressCoordinateLogRequest) (res *logisticService.GetMostTrustedAddressCoordinateLogResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetMostTrustedAddressCoordinateLog")
	defer span.End()

	addressCoordinateLog, err := h.ServicesAddressCoordinateLog.GetMostTrusted(ctx, req.AddressId)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &logisticService.GetMostTrustedAddressCoordinateLogResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &logisticService.AddressCoordinateLog{
			Id:             addressCoordinateLog.ID,
			AddressId:      addressCoordinateLog.AddressID,
			SalesOrderId:   addressCoordinateLog.SalesOrderID,
			Latitude:       &addressCoordinateLog.Latitude,
			Longitude:      &addressCoordinateLog.Longitude,
			LogChannelId:   int32(addressCoordinateLog.LogChannelID),
			MainCoordinate: int32(addressCoordinateLog.MainCoordinate),
			CreatedAt:      timestamppb.New(addressCoordinateLog.CreatedAt),
			CreatedBy:      addressCoordinateLog.CreatedBy,
		},
	}

	return
}
