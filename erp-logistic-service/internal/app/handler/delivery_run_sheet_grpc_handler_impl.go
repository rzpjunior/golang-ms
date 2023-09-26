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

func (h *LogisticGrpcHandler) GetDeliveryRunSheetList(ctx context.Context, req *logisticService.GetDeliveryRunSheetListRequest) (res *logisticService.GetDeliveryRunSheetListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetDeliveryRunSheetList")
	defer span.End()

	var statusInt []int

	for _, status := range req.Status {
		statusInt = append(statusInt, int(status))
	}

	deliveryRunSheets, _, err := h.ServicesDeliveryRunSheet.Get(ctx, dto.DeliveryRunSheetGetRequest{
		Offset:        int(req.Offset),
		Limit:         int(req.Limit),
		OrderBy:       req.OrderBy,
		GroupBy:       req.GroupBy,
		Status:        statusInt,
		ArrCourierIDs: req.CourierId,
	})
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*logisticService.DeliveryRunSheet
	for _, drs := range deliveryRunSheets {
		data = append(data, &logisticService.DeliveryRunSheet{
			Id:                drs.ID,
			Code:              drs.Code,
			CourierId:         drs.CourierID,
			DeliveryDate:      timestamppb.New(drs.DeliveryDate),
			StartedAt:         timestamppb.New(drs.StartedAt),
			FinishedAt:        timestamppb.New(drs.FinishedAt),
			StartingLatitude:  drs.StartingLatitude,
			StartingLongitude: drs.StartingLongitude,
			FinishedLatitude:  drs.FinishedLatitude,
			FinishedLongitude: drs.FinishedLongitude,
			Status:            int32(drs.Status),
		})
	}

	res = &logisticService.GetDeliveryRunSheetListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}

	return
}

func (h *LogisticGrpcHandler) GetDeliveryRunSheetDetail(ctx context.Context, req *logisticService.GetDeliveryRunSheetDetailRequest) (res *logisticService.GetDeliveryRunSheetDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetDeliveryRunSheetDetail")
	defer span.End()

	deliveryRunSheet, err := h.ServicesDeliveryRunSheet.GetDetail(ctx, req.Id, req.Code)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &logisticService.GetDeliveryRunSheetDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &logisticService.DeliveryRunSheet{
			Id:                deliveryRunSheet.ID,
			Code:              deliveryRunSheet.Code,
			CourierId:         deliveryRunSheet.CourierID,
			DeliveryDate:      timestamppb.New(deliveryRunSheet.DeliveryDate),
			StartedAt:         timestamppb.New(deliveryRunSheet.StartedAt),
			FinishedAt:        timestamppb.New(deliveryRunSheet.FinishedAt),
			StartingLatitude:  deliveryRunSheet.StartingLatitude,
			StartingLongitude: deliveryRunSheet.StartingLongitude,
			FinishedLatitude:  deliveryRunSheet.FinishedLatitude,
			FinishedLongitude: deliveryRunSheet.FinishedLongitude,
			Status:            int32(deliveryRunSheet.Status),
		},
	}

	return
}

func (h *LogisticGrpcHandler) CreateDeliveryRunSheet(ctx context.Context, req *logisticService.CreateDeliveryRunSheetRequest) (res *logisticService.CreateDeliveryRunSheetResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.CreateDeliveryRunSheet")
	defer span.End()

	deliveryRunSheet, err := h.ServicesDeliveryRunSheet.Create(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &logisticService.CreateDeliveryRunSheetResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &logisticService.DeliveryRunSheet{
			Id:                deliveryRunSheet.ID,
			Code:              deliveryRunSheet.Code,
			CourierId:         deliveryRunSheet.CourierID,
			DeliveryDate:      timestamppb.New(deliveryRunSheet.DeliveryDate),
			StartedAt:         timestamppb.New(deliveryRunSheet.StartedAt),
			FinishedAt:        timestamppb.New(deliveryRunSheet.FinishedAt),
			StartingLatitude:  deliveryRunSheet.StartingLatitude,
			StartingLongitude: deliveryRunSheet.StartingLongitude,
			FinishedLatitude:  deliveryRunSheet.FinishedLatitude,
			FinishedLongitude: deliveryRunSheet.FinishedLongitude,
			Status:            int32(deliveryRunSheet.Status),
		},
	}

	return
}

func (h *LogisticGrpcHandler) FinishDeliveryRunSheet(ctx context.Context, req *logisticService.FinishDeliveryRunSheetRequest) (res *logisticService.FinishDeliveryRunSheetResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.FinishDeliveryRunSheet")
	defer span.End()

	deliveryRunSheet, err := h.ServicesDeliveryRunSheet.Finish(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &logisticService.FinishDeliveryRunSheetResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &logisticService.DeliveryRunSheet{
			Id:                deliveryRunSheet.ID,
			Code:              deliveryRunSheet.Code,
			CourierId:         deliveryRunSheet.CourierID,
			DeliveryDate:      timestamppb.New(deliveryRunSheet.DeliveryDate),
			StartedAt:         timestamppb.New(deliveryRunSheet.StartedAt),
			FinishedAt:        timestamppb.New(deliveryRunSheet.FinishedAt),
			StartingLatitude:  deliveryRunSheet.StartingLatitude,
			StartingLongitude: deliveryRunSheet.StartingLongitude,
			FinishedLatitude:  deliveryRunSheet.FinishedLatitude,
			FinishedLongitude: deliveryRunSheet.FinishedLongitude,
			Status:            int32(deliveryRunSheet.Status),
		},
	}

	return
}
