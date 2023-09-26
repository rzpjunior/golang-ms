package handler

import (
	context "context"
	"fmt"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-logistic-service/internal/app/dto"
	logisticService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/logistic_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *LogisticGrpcHandler) GetDeliveryRunSheetItemList(ctx context.Context, req *logisticService.GetDeliveryRunSheetItemListRequest) (res *logisticService.GetDeliveryRunSheetItemListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetDeliveryRunSheetItemList")
	defer span.End()
	fmt.Println("masuk", req.SalesOrderId)
	request := dto.DeliveryRunSheetItemGetRequest{
		Offset:               int(req.Offset),
		Limit:                int(req.Limit),
		OrderBy:              req.OrderBy,
		GroupBy:              req.GroupBy,
		DeliveryRunSheetIDs:  req.DeliveryRunSheetId,
		CourierIDs:           req.CourierId,
		ArrSalesOrderIDs:     req.SalesOrderId,
		SearchSalesOrderCode: req.Search,
	}

	for _, status := range req.Status {
		request.Status = append(request.Status, int(status))
	}
	for _, stepType := range req.StepType {
		request.StepType = append(request.StepType, int(stepType))
	}

	deliveryRunSheetItems, _, err := h.ServicesDeliveryRunSheetItem.Get(ctx, &request)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*logisticService.DeliveryRunSheetItem
	for _, drsi := range deliveryRunSheetItems {
		data = append(data, &logisticService.DeliveryRunSheetItem{
			Id:                          drsi.ID,
			DeliveryRunSheetId:          drsi.DeliveryRunSheetID,
			CourierId:                   drsi.CourierID,
			SalesOrderId:                drsi.SalesOrderID,
			StepType:                    int32(drsi.StepType),
			Latitude:                    drsi.Latitude,
			Longitude:                   drsi.Longitude,
			Status:                      int32(drsi.Status),
			Note:                        drsi.Note,
			RecipientName:               drsi.RecipientName,
			MoneyReceived:               drsi.MoneyReceived,
			DeliveryEvidenceImageUrl:    drsi.DeliveryEvidenceImageURL,
			TransactionEvidenceImageUrl: drsi.TransactionEvidenceImageURL,
			ArrivalTime:                 timestamppb.New(drsi.ArrivalTime),
			UnpunctualReason:            int32(drsi.UnpunctualReason),
			UnpunctualDetail:            int32(drsi.UnpunctualDetail),
			FarDeliveryReason:           drsi.FarDeliveryReason,
			CreatedAt:                   timestamppb.New(drsi.CreatedAt),
			StartedAt:                   timestamppb.New(drsi.StartedAt),
			FinishedAt:                  timestamppb.New(drsi.FinishedAt),
		})
	}

	res = &logisticService.GetDeliveryRunSheetItemListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}

	return
}

func (h *LogisticGrpcHandler) GetDeliveryRunSheetItemDetail(ctx context.Context, req *logisticService.GetDeliveryRunSheetItemDetailRequest) (res *logisticService.GetDeliveryRunSheetItemDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetDeliveryRunSheetItemDetail")
	defer span.End()

	deliveryRunSheetItem, err := h.ServicesDeliveryRunSheetItem.GetDetail(ctx, req.Id)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &logisticService.GetDeliveryRunSheetItemDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &logisticService.DeliveryRunSheetItem{
			Id:                          deliveryRunSheetItem.ID,
			DeliveryRunSheetId:          deliveryRunSheetItem.DeliveryRunSheetID,
			CourierId:                   deliveryRunSheetItem.CourierID,
			SalesOrderId:                deliveryRunSheetItem.SalesOrderID,
			StepType:                    int32(deliveryRunSheetItem.StepType),
			Latitude:                    deliveryRunSheetItem.Latitude,
			Longitude:                   deliveryRunSheetItem.Longitude,
			Status:                      int32(deliveryRunSheetItem.Status),
			Note:                        deliveryRunSheetItem.Note,
			RecipientName:               deliveryRunSheetItem.RecipientName,
			MoneyReceived:               deliveryRunSheetItem.MoneyReceived,
			DeliveryEvidenceImageUrl:    deliveryRunSheetItem.DeliveryEvidenceImageURL,
			TransactionEvidenceImageUrl: deliveryRunSheetItem.TransactionEvidenceImageURL,
			ArrivalTime:                 timestamppb.New(deliveryRunSheetItem.ArrivalTime),
			UnpunctualReason:            int32(deliveryRunSheetItem.UnpunctualReason),
			UnpunctualDetail:            int32(deliveryRunSheetItem.UnpunctualDetail),
			FarDeliveryReason:           deliveryRunSheetItem.FarDeliveryReason,
			CreatedAt:                   timestamppb.New(deliveryRunSheetItem.CreatedAt),
			StartedAt:                   timestamppb.New(deliveryRunSheetItem.StartedAt),
			FinishedAt:                  timestamppb.New(deliveryRunSheetItem.FinishedAt),
		},
	}

	return
}

func (h *LogisticGrpcHandler) CreateDeliveryRunSheetItemPickup(ctx context.Context, req *logisticService.CreateDeliveryRunSheetItemRequest) (res *logisticService.CreateDeliveryRunSheetItemResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.CreateDeliveryRunSheetPickup")
	defer span.End()

	deliveryRunSheetItem, err := h.ServicesDeliveryRunSheetItem.CreatePickup(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &logisticService.CreateDeliveryRunSheetItemResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &logisticService.DeliveryRunSheetItem{
			Id:                          deliveryRunSheetItem.ID,
			DeliveryRunSheetId:          deliveryRunSheetItem.DeliveryRunSheetID,
			CourierId:                   deliveryRunSheetItem.CourierID,
			SalesOrderId:                deliveryRunSheetItem.SalesOrderID,
			StepType:                    int32(deliveryRunSheetItem.StepType),
			Latitude:                    deliveryRunSheetItem.Latitude,
			Longitude:                   deliveryRunSheetItem.Longitude,
			Status:                      int32(deliveryRunSheetItem.Status),
			Note:                        deliveryRunSheetItem.Note,
			RecipientName:               deliveryRunSheetItem.RecipientName,
			MoneyReceived:               deliveryRunSheetItem.MoneyReceived,
			DeliveryEvidenceImageUrl:    deliveryRunSheetItem.DeliveryEvidenceImageURL,
			TransactionEvidenceImageUrl: deliveryRunSheetItem.TransactionEvidenceImageURL,
			ArrivalTime:                 timestamppb.New(deliveryRunSheetItem.ArrivalTime),
			UnpunctualReason:            int32(deliveryRunSheetItem.UnpunctualReason),
			UnpunctualDetail:            int32(deliveryRunSheetItem.UnpunctualDetail),
			FarDeliveryReason:           deliveryRunSheetItem.FarDeliveryReason,
			CreatedAt:                   timestamppb.New(deliveryRunSheetItem.CreatedAt),
			StartedAt:                   timestamppb.New(deliveryRunSheetItem.StartedAt),
			FinishedAt:                  timestamppb.New(deliveryRunSheetItem.FinishedAt),
		},
	}

	return
}

func (h *LogisticGrpcHandler) CreateDeliveryRunSheetItemDelivery(ctx context.Context, req *logisticService.CreateDeliveryRunSheetItemRequest) (res *logisticService.CreateDeliveryRunSheetItemResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.CreateDeliveryRunSheetDelivery")
	defer span.End()

	deliveryRunSheetItem, err := h.ServicesDeliveryRunSheetItem.CreateDelivery(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &logisticService.CreateDeliveryRunSheetItemResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &logisticService.DeliveryRunSheetItem{
			Id:                          deliveryRunSheetItem.ID,
			DeliveryRunSheetId:          deliveryRunSheetItem.DeliveryRunSheetID,
			CourierId:                   deliveryRunSheetItem.CourierID,
			SalesOrderId:                deliveryRunSheetItem.SalesOrderID,
			StepType:                    int32(deliveryRunSheetItem.StepType),
			Latitude:                    deliveryRunSheetItem.Latitude,
			Longitude:                   deliveryRunSheetItem.Longitude,
			Status:                      int32(deliveryRunSheetItem.Status),
			Note:                        deliveryRunSheetItem.Note,
			RecipientName:               deliveryRunSheetItem.RecipientName,
			MoneyReceived:               deliveryRunSheetItem.MoneyReceived,
			DeliveryEvidenceImageUrl:    deliveryRunSheetItem.DeliveryEvidenceImageURL,
			TransactionEvidenceImageUrl: deliveryRunSheetItem.TransactionEvidenceImageURL,
			ArrivalTime:                 timestamppb.New(deliveryRunSheetItem.ArrivalTime),
			UnpunctualReason:            int32(deliveryRunSheetItem.UnpunctualReason),
			UnpunctualDetail:            int32(deliveryRunSheetItem.UnpunctualDetail),
			FarDeliveryReason:           deliveryRunSheetItem.FarDeliveryReason,
			CreatedAt:                   timestamppb.New(deliveryRunSheetItem.CreatedAt),
			StartedAt:                   timestamppb.New(deliveryRunSheetItem.StartedAt),
			FinishedAt:                  timestamppb.New(deliveryRunSheetItem.FinishedAt),
		},
	}

	return
}

func (h *LogisticGrpcHandler) StartDeliveryRunSheetItem(ctx context.Context, req *logisticService.StartDeliveryRunSheetItemRequest) (res *logisticService.StartDeliveryRunSheetItemResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.StartDeliveryRunSheetItem")
	defer span.End()

	deliveryRunSheetItem, err := h.ServicesDeliveryRunSheetItem.Start(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &logisticService.StartDeliveryRunSheetItemResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &logisticService.DeliveryRunSheetItem{
			Id:                          deliveryRunSheetItem.ID,
			DeliveryRunSheetId:          deliveryRunSheetItem.DeliveryRunSheetID,
			CourierId:                   deliveryRunSheetItem.CourierID,
			SalesOrderId:                deliveryRunSheetItem.SalesOrderID,
			StepType:                    int32(deliveryRunSheetItem.StepType),
			Latitude:                    deliveryRunSheetItem.Latitude,
			Longitude:                   deliveryRunSheetItem.Longitude,
			Status:                      int32(deliveryRunSheetItem.Status),
			Note:                        deliveryRunSheetItem.Note,
			RecipientName:               deliveryRunSheetItem.RecipientName,
			MoneyReceived:               deliveryRunSheetItem.MoneyReceived,
			DeliveryEvidenceImageUrl:    deliveryRunSheetItem.DeliveryEvidenceImageURL,
			TransactionEvidenceImageUrl: deliveryRunSheetItem.TransactionEvidenceImageURL,
			ArrivalTime:                 timestamppb.New(deliveryRunSheetItem.ArrivalTime),
			UnpunctualReason:            int32(deliveryRunSheetItem.UnpunctualReason),
			UnpunctualDetail:            int32(deliveryRunSheetItem.UnpunctualDetail),
			FarDeliveryReason:           deliveryRunSheetItem.FarDeliveryReason,
			CreatedAt:                   timestamppb.New(deliveryRunSheetItem.CreatedAt),
			StartedAt:                   timestamppb.New(deliveryRunSheetItem.StartedAt),
			FinishedAt:                  timestamppb.New(deliveryRunSheetItem.FinishedAt),
		},
	}

	return
}

func (h *LogisticGrpcHandler) PostponeDeliveryRunSheetItem(ctx context.Context, req *logisticService.PostponeDeliveryRunSheetItemRequest) (res *logisticService.PostponeDeliveryRunSheetItemResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.PostponeDeliveryRunSheetItem")
	defer span.End()

	deliveryRunSheetItem, err := h.ServicesDeliveryRunSheetItem.Postpone(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	res = &logisticService.PostponeDeliveryRunSheetItemResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &logisticService.DeliveryRunSheetItem{
			Id:                          deliveryRunSheetItem.ID,
			DeliveryRunSheetId:          deliveryRunSheetItem.DeliveryRunSheetID,
			CourierId:                   deliveryRunSheetItem.CourierID,
			SalesOrderId:                deliveryRunSheetItem.SalesOrderID,
			StepType:                    int32(deliveryRunSheetItem.StepType),
			Latitude:                    deliveryRunSheetItem.Latitude,
			Longitude:                   deliveryRunSheetItem.Longitude,
			Status:                      int32(deliveryRunSheetItem.Status),
			Note:                        deliveryRunSheetItem.Note,
			RecipientName:               deliveryRunSheetItem.RecipientName,
			MoneyReceived:               deliveryRunSheetItem.MoneyReceived,
			DeliveryEvidenceImageUrl:    deliveryRunSheetItem.DeliveryEvidenceImageURL,
			TransactionEvidenceImageUrl: deliveryRunSheetItem.TransactionEvidenceImageURL,
			ArrivalTime:                 timestamppb.New(deliveryRunSheetItem.ArrivalTime),
			UnpunctualReason:            int32(deliveryRunSheetItem.UnpunctualReason),
			UnpunctualDetail:            int32(deliveryRunSheetItem.UnpunctualDetail),
			FarDeliveryReason:           deliveryRunSheetItem.FarDeliveryReason,
			CreatedAt:                   timestamppb.New(deliveryRunSheetItem.CreatedAt),
			StartedAt:                   timestamppb.New(deliveryRunSheetItem.StartedAt),
			FinishedAt:                  timestamppb.New(deliveryRunSheetItem.FinishedAt),
		},
	}

	return
}

func (h *LogisticGrpcHandler) FailPickupDeliveryRunSheetItem(ctx context.Context, req *logisticService.FailPickupDeliveryRunSheetItemRequest) (res *logisticService.FailPickupDeliveryRunSheetItemResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.FailPickupDeliveryRunSheetItem")
	defer span.End()

	deliveryRunSheetItem, err := h.ServicesDeliveryRunSheetItem.FailPickup(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &logisticService.FailPickupDeliveryRunSheetItemResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &logisticService.DeliveryRunSheetItem{
			Id:                          deliveryRunSheetItem.ID,
			DeliveryRunSheetId:          deliveryRunSheetItem.DeliveryRunSheetID,
			CourierId:                   deliveryRunSheetItem.CourierID,
			SalesOrderId:                deliveryRunSheetItem.SalesOrderID,
			StepType:                    int32(deliveryRunSheetItem.StepType),
			Latitude:                    deliveryRunSheetItem.Latitude,
			Longitude:                   deliveryRunSheetItem.Longitude,
			Status:                      int32(deliveryRunSheetItem.Status),
			Note:                        deliveryRunSheetItem.Note,
			RecipientName:               deliveryRunSheetItem.RecipientName,
			MoneyReceived:               deliveryRunSheetItem.MoneyReceived,
			DeliveryEvidenceImageUrl:    deliveryRunSheetItem.DeliveryEvidenceImageURL,
			TransactionEvidenceImageUrl: deliveryRunSheetItem.TransactionEvidenceImageURL,
			ArrivalTime:                 timestamppb.New(deliveryRunSheetItem.ArrivalTime),
			UnpunctualReason:            int32(deliveryRunSheetItem.UnpunctualReason),
			UnpunctualDetail:            int32(deliveryRunSheetItem.UnpunctualDetail),
			FarDeliveryReason:           deliveryRunSheetItem.FarDeliveryReason,
			CreatedAt:                   timestamppb.New(deliveryRunSheetItem.CreatedAt),
			StartedAt:                   timestamppb.New(deliveryRunSheetItem.StartedAt),
			FinishedAt:                  timestamppb.New(deliveryRunSheetItem.FinishedAt),
		},
	}

	return
}

func (h *LogisticGrpcHandler) FailDeliveryDeliveryRunSheetItem(ctx context.Context, req *logisticService.FailDeliveryDeliveryRunSheetItemRequest) (res *logisticService.FailDeliveryDeliveryRunSheetItemResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.FailDeliveryDeliveryRunSheetItem")
	defer span.End()

	deliveryRunSheetItem, err := h.ServicesDeliveryRunSheetItem.FailDelivery(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &logisticService.FailDeliveryDeliveryRunSheetItemResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &logisticService.DeliveryRunSheetItem{
			Id:                          deliveryRunSheetItem.ID,
			DeliveryRunSheetId:          deliveryRunSheetItem.DeliveryRunSheetID,
			CourierId:                   deliveryRunSheetItem.CourierID,
			SalesOrderId:                deliveryRunSheetItem.SalesOrderID,
			StepType:                    int32(deliveryRunSheetItem.StepType),
			Latitude:                    deliveryRunSheetItem.Latitude,
			Longitude:                   deliveryRunSheetItem.Longitude,
			Status:                      int32(deliveryRunSheetItem.Status),
			Note:                        deliveryRunSheetItem.Note,
			RecipientName:               deliveryRunSheetItem.RecipientName,
			MoneyReceived:               deliveryRunSheetItem.MoneyReceived,
			DeliveryEvidenceImageUrl:    deliveryRunSheetItem.DeliveryEvidenceImageURL,
			TransactionEvidenceImageUrl: deliveryRunSheetItem.TransactionEvidenceImageURL,
			ArrivalTime:                 timestamppb.New(deliveryRunSheetItem.ArrivalTime),
			UnpunctualReason:            int32(deliveryRunSheetItem.UnpunctualReason),
			UnpunctualDetail:            int32(deliveryRunSheetItem.UnpunctualDetail),
			FarDeliveryReason:           deliveryRunSheetItem.FarDeliveryReason,
			CreatedAt:                   timestamppb.New(deliveryRunSheetItem.CreatedAt),
			StartedAt:                   timestamppb.New(deliveryRunSheetItem.StartedAt),
			FinishedAt:                  timestamppb.New(deliveryRunSheetItem.FinishedAt),
		},
	}

	return
}

func (h *LogisticGrpcHandler) SuccessDeliveryRunSheetItem(ctx context.Context, req *logisticService.SuccessDeliveryRunSheetItemRequest) (res *logisticService.SuccessDeliveryRunSheetItemResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.SuccessDeliveryRunSheetItem")
	defer span.End()

	deliveryRunSheetItem, err := h.ServicesDeliveryRunSheetItem.Success(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &logisticService.SuccessDeliveryRunSheetItemResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &logisticService.DeliveryRunSheetItem{
			Id:                          deliveryRunSheetItem.ID,
			DeliveryRunSheetId:          deliveryRunSheetItem.DeliveryRunSheetID,
			CourierId:                   deliveryRunSheetItem.CourierID,
			SalesOrderId:                deliveryRunSheetItem.SalesOrderID,
			StepType:                    int32(deliveryRunSheetItem.StepType),
			Latitude:                    deliveryRunSheetItem.Latitude,
			Longitude:                   deliveryRunSheetItem.Longitude,
			Status:                      int32(deliveryRunSheetItem.Status),
			Note:                        deliveryRunSheetItem.Note,
			RecipientName:               deliveryRunSheetItem.RecipientName,
			MoneyReceived:               deliveryRunSheetItem.MoneyReceived,
			DeliveryEvidenceImageUrl:    deliveryRunSheetItem.DeliveryEvidenceImageURL,
			TransactionEvidenceImageUrl: deliveryRunSheetItem.TransactionEvidenceImageURL,
			ArrivalTime:                 timestamppb.New(deliveryRunSheetItem.ArrivalTime),
			UnpunctualReason:            int32(deliveryRunSheetItem.UnpunctualReason),
			UnpunctualDetail:            int32(deliveryRunSheetItem.UnpunctualDetail),
			FarDeliveryReason:           deliveryRunSheetItem.FarDeliveryReason,
			CreatedAt:                   timestamppb.New(deliveryRunSheetItem.CreatedAt),
			StartedAt:                   timestamppb.New(deliveryRunSheetItem.StartedAt),
			FinishedAt:                  timestamppb.New(deliveryRunSheetItem.FinishedAt),
		},
	}

	return
}

func (h *LogisticGrpcHandler) ArrivedDeliveryRunSheetItem(ctx context.Context, req *logisticService.ArrivedDeliveryRunSheetItemRequest) (res *logisticService.ArrivedDeliveryRunSheetItemResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.ArrivedDeliveryRunSheetItem")
	defer span.End()

	deliveryRunSheetItem, err := h.ServicesDeliveryRunSheetItem.Arrived(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &logisticService.ArrivedDeliveryRunSheetItemResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &logisticService.DeliveryRunSheetItem{
			Id:                          deliveryRunSheetItem.ID,
			DeliveryRunSheetId:          deliveryRunSheetItem.DeliveryRunSheetID,
			CourierId:                   deliveryRunSheetItem.CourierID,
			SalesOrderId:                deliveryRunSheetItem.SalesOrderID,
			StepType:                    int32(deliveryRunSheetItem.StepType),
			Latitude:                    deliveryRunSheetItem.Latitude,
			Longitude:                   deliveryRunSheetItem.Longitude,
			Status:                      int32(deliveryRunSheetItem.Status),
			Note:                        deliveryRunSheetItem.Note,
			RecipientName:               deliveryRunSheetItem.RecipientName,
			MoneyReceived:               deliveryRunSheetItem.MoneyReceived,
			DeliveryEvidenceImageUrl:    deliveryRunSheetItem.DeliveryEvidenceImageURL,
			TransactionEvidenceImageUrl: deliveryRunSheetItem.TransactionEvidenceImageURL,
			ArrivalTime:                 timestamppb.New(deliveryRunSheetItem.ArrivalTime),
			UnpunctualReason:            int32(deliveryRunSheetItem.UnpunctualReason),
			UnpunctualDetail:            int32(deliveryRunSheetItem.UnpunctualDetail),
			FarDeliveryReason:           deliveryRunSheetItem.FarDeliveryReason,
			CreatedAt:                   timestamppb.New(deliveryRunSheetItem.CreatedAt),
			StartedAt:                   timestamppb.New(deliveryRunSheetItem.StartedAt),
			FinishedAt:                  timestamppb.New(deliveryRunSheetItem.FinishedAt),
		},
	}

	return
}
