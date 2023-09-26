package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/dto"
	crmService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/crm_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *CrmGrpcHandler) GetSalesAssignmentItemList(ctx context.Context, req *crmService.GetSalesAssignmentItemListRequest) (res *crmService.GetSalesAssignmentItemListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "CrmGrpcHandler.GetSalesAssignmentItemList")
	defer span.End()

	// var items []*dto.SalesAssignmentItemResponse
	// items, _, err = h.ServicesSalesAssignmentItem.Get(ctx, int(req.Offset), int(req.Limit), int(req.Status), req.Search, req.OrderBy, int(req.TeritoryId), int(req.SalespersonId), req.SubmitDateFrom.AsTime(), req.SubmitDateTo.AsTime(), req.StartDateFrom.AsTime(), req.StartDateTo.AsTime(), req.EndDateFrom.AsTime(), req.EndDateTo.AsTime(), int(req.Task), int(req.OutOfRoute), req.CustomerType)
	// if err != nil {
	// 	err = status.New(codes.NotFound, err.Error()).Err()
	// 	h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
	// 	return
	// }

	// var data []*crmService.SalesAssignmentItem
	// for _, item := range items {
	// 	var objectiveValues []*crmService.SalesAssignmentObjective
	// 	for _, val := range item.ObjectiveValues {
	// 		objectiveValues = append(objectiveValues, &crmService.SalesAssignmentObjective{
	// 			Id:         val.ID,
	// 			Code:       val.Code,
	// 			Name:       val.Name,
	// 			Objective:  val.Objective,
	// 			SurveyLink: val.SurveyLink,
	// 			Status:     int32(val.Status),
	// 			CreatedAt:  timestamppb.New(val.CreatedAt),
	// 			CreatedBy:  val.CreatedBy.ID,
	// 			UpdatedAt:  timestamppb.New(val.UpdatedAt),
	// 		})
	// 	}
	// 	var finishDate *timestamppb.Timestamp

	// 	if item.FinishDate != nil {
	// 		finishDate = timestamppb.New(*item.FinishDate)
	// 	}
	// 	var (
	// 		address *crmService.Address
	// 		ca      *crmService.CustomerAcquisitionResponse
	// 	)
	// 	if item.Address != nil {
	// 		address = &crmService.Address{
	// 			Id:           item.Address.ID,
	// 			Code:         item.Address.Code,
	// 			CustomerName: item.Address.Name,
	// 		}
	// 	}
	// 	if item.CustomerAcquisition != nil {
	// 		ca = &crmService.CustomerAcquisitionResponse{
	// 			Id:   item.CustomerAcquisition.ID,
	// 			Code: item.CustomerAcquisition.Code,
	// 			Name: item.CustomerAcquisition.Name,
	// 		}
	// 	}
	// 	data = append(data, &crmService.SalesAssignmentItem{
	// 		Id:                    item.ID,
	// 		SalesAssignmentId:     item.SalesAssignmentID,
	// 		SalesPersonId:         item.SalesPerson.ID,
	// 		AddressId:             item.AddressId,
	// 		CustomerAcquisitionId: item.CustomerAcquisitionID,
	// 		Latitude:              item.Latitude,
	// 		Longitude:             item.Longitude,
	// 		Task:                  int32(item.Task),
	// 		CustomerType:          int32(item.CustomerType),
	// 		ObjectiveCodes:        item.ObjectiveCodes,
	// 		ActualDistance:        item.ActualDistance,
	// 		OutOfRoute:            int32(item.OutOfRoute),
	// 		StartDate:             timestamppb.New(item.StartDate),
	// 		EndDate:               timestamppb.New(item.EndDate),
	// 		FinishDate:            finishDate,
	// 		SubmitDate:            timestamppb.New(item.SubmitDate),
	// 		TaskImageUrl:          item.TaskImageUrls,
	// 		TaskAnswer:            int32(item.TaskAnswer),
	// 		Status:                int32(item.Status),
	// 		EffectiveCall:         int32(item.EffectiveCall),
	// 		Salesperson: &crmService.User{
	// 			Id:           item.SalesPerson.ID,
	// 			EmployeeCode: item.SalesPerson.Code,
	// 			Name:         item.SalesPerson.Name,
	// 		},
	// 		ObectiveValues:      objectiveValues,
	// 		Address:             address,
	// 		CustomerAcquisition: ca,
	// 	})
	// }

	// res = &crmService.GetSalesAssignmentItemListResponse{
	// 	Code:    int32(codes.OK),
	// 	Message: codes.OK.String(),
	// 	Data:    data,
	// }
	return
}

func (h *CrmGrpcHandler) GetSalesAssignmentItemDetail(ctx context.Context, req *crmService.GetSalesAssignmentItemDetailRequest) (res *crmService.GetSalesAssignmentItemDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetSalesAssignmentItemDetail")
	defer span.End()

	// var item dto.SalesAssignmentItemResponse
	// item, err = h.ServicesSalesAssignmentItem.GetByID(ctx, req.Id)
	// if err != nil {
	// 	err = status.New(codes.NotFound, err.Error()).Err()
	// 	h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
	// 	return
	// }

	// var objectiveValues []*crmService.SalesAssignmentObjective
	// for _, val := range item.ObjectiveValues {
	// 	objectiveValues = append(objectiveValues, &crmService.SalesAssignmentObjective{
	// 		Id:         val.ID,
	// 		Code:       val.Code,
	// 		Name:       val.Name,
	// 		Objective:  val.Objective,
	// 		SurveyLink: val.SurveyLink,
	// 		Status:     int32(val.Status),
	// 		CreatedAt:  timestamppb.New(val.CreatedAt),
	// 		CreatedBy:  val.CreatedBy.ID,
	// 		UpdatedAt:  timestamppb.New(val.UpdatedAt),
	// 	})
	// }

	// var finishDate *timestamppb.Timestamp
	// if item.FinishDate != nil {
	// 	finishDate = timestamppb.New(*item.FinishDate)
	// }
	// var (
	// 	address *crmService.Address
	// 	ca      *crmService.CustomerAcquisitionResponse
	// )
	// if item.Address != nil {
	// 	address = &crmService.Address{
	// 		Id:           item.Address.ID,
	// 		Code:         item.Address.Code,
	// 		CustomerName: item.Address.Name,
	// 	}
	// }
	// if item.CustomerAcquisition != nil {
	// 	ca = &crmService.CustomerAcquisitionResponse{
	// 		Id:   item.CustomerAcquisition.ID,
	// 		Code: item.CustomerAcquisition.Code,
	// 		Name: item.CustomerAcquisition.Name,
	// 	}
	// }
	// res = &crmService.GetSalesAssignmentItemDetailResponse{
	// 	Code:    int32(codes.OK),
	// 	Message: codes.OK.String(),
	// 	Data: &crmService.SalesAssignmentItem{
	// 		Id:                    item.ID,
	// 		SalesAssignmentId:     item.SalesAssignmentID,
	// 		SalesPersonId:         item.SalesPerson.ID,
	// 		AddressId:             item.Address.ID,
	// 		CustomerAcquisitionId: item.CustomerAcquisitionID,
	// 		Latitude:              item.Latitude,
	// 		Longitude:             item.Longitude,
	// 		Task:                  int32(item.Task),
	// 		CustomerType:          int32(item.CustomerType),
	// 		ObjectiveCodes:        item.ObjectiveCodes,
	// 		ActualDistance:        item.ActualDistance,
	// 		OutOfRoute:            int32(item.OutOfRoute),
	// 		StartDate:             timestamppb.New(item.StartDate),
	// 		EndDate:               timestamppb.New(item.EndDate),
	// 		FinishDate:            finishDate,
	// 		SubmitDate:            timestamppb.New(item.SubmitDate),
	// 		TaskImageUrl:          item.TaskImageUrls,
	// 		TaskAnswer:            int32(item.TaskAnswer),
	// 		Status:                int32(item.Status),
	// 		EffectiveCall:         int32(item.EffectiveCall),
	// 		Address:               address,
	// 		CustomerAcquisition:   ca,
	// 		Salesperson: &crmService.User{
	// 			Id:           item.SalesPerson.ID,
	// 			EmployeeCode: item.SalesPerson.Code,
	// 			Name:         item.SalesPerson.Name,
	// 		},
	// 		ObectiveValues: objectiveValues,
	// 	},
	// }
	return
}

func (h *CrmGrpcHandler) CheckTaskSalesAssignmentItemActive(ctx context.Context, req *crmService.CheckTaskSalesAssignmentItemRequest) (res *crmService.CheckTaskSalesAssignmentItemResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.CheckTaskSalesAssignmentItemActive")
	defer span.End()

	// var existed bool
	// existed, err = h.ServicesSalesAssignmentItem.CheckActiveTask(ctx, req.SalespersonId)
	// if err != nil {
	// 	err = status.New(codes.NotFound, err.Error()).Err()
	// 	h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
	// 	return
	// }

	// res = &crmService.CheckTaskSalesAssignmentItemResponse{
	// 	Code:    int32(codes.OK),
	// 	Message: codes.OK.String(),
	// 	Data: &crmService.BooleanResponse{
	// 		Existed: existed,
	// 	},
	// }
	return
}

func (h *CrmGrpcHandler) SubmitTaskVisitFU(ctx context.Context, req *crmService.UpdateSubmitTaskVisitFURequest) (res *crmService.UpdateSubmitTaskVisitFUResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.CheckTaskSalesAssignmentItemActive")
	defer span.End()

	// var sai *dto.SalesAssignmentItemResponse
	// sai, err = h.ServicesSalesAssignmentItem.SubmitTaskVisitFU(ctx, dto.UpdateSubmitTaskVisitFURequest{
	// 	SalesAssignmentItemResponse: dto.SalesAssignmentItemResponse{
	// 		ID:                    req.Id,
	// 		CustomerAcquisitionID: req.CustomerAcquisitionId,
	// 		AddressId:             req.AddressId,
	// 		TaskAnswer:            int8(req.TaskAnswer),
	// 		TaskImageUrls:         strings.Split(req.TaskImageUrls, ","),
	// 		Latitude:              req.Latitude,
	// 		Longitude:             req.Longitude,
	// 		ActualDistance:        req.ActualDistance,
	// 	},
	// })
	// if err != nil {
	// 	err = status.New(codes.NotFound, err.Error()).Err()
	// 	h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
	// 	return
	// }

	// var finishDate *timestamppb.Timestamp
	// if sai.FinishDate != nil {
	// 	finishDate = timestamppb.New(*sai.FinishDate)
	// }
	// var (
	// 	address *crmService.Address
	// 	ca      *crmService.CustomerAcquisitionResponse
	// )
	// if sai.Address != nil {
	// 	address = &crmService.Address{
	// 		Id:           sai.Address.ID,
	// 		Code:         sai.Address.Code,
	// 		CustomerName: sai.Address.Name,
	// 	}
	// }
	// if sai.CustomerAcquisition != nil {
	// 	ca = &crmService.CustomerAcquisitionResponse{
	// 		Id:   sai.CustomerAcquisition.ID,
	// 		Code: sai.CustomerAcquisition.Code,
	// 		Name: sai.CustomerAcquisition.Name,
	// 	}
	// }
	// res = &crmService.UpdateSubmitTaskVisitFUResponse{
	// 	Code:    int32(codes.OK),
	// 	Message: codes.OK.String(),
	// 	Data: &crmService.SalesAssignmentItem{
	// 		Id:                    sai.ID,
	// 		SalesAssignmentId:     sai.SalesAssignmentID,
	// 		AddressId:             sai.AddressId,
	// 		CustomerAcquisitionId: sai.CustomerAcquisitionID,
	// 		Latitude:              sai.Latitude,
	// 		Longitude:             sai.Longitude,
	// 		Task:                  int32(sai.Task),
	// 		CustomerType:          int32(sai.CustomerType),
	// 		ObjectiveCodes:        sai.ObjectiveCodes,
	// 		ActualDistance:        sai.ActualDistance,
	// 		OutOfRoute:            int32(sai.OutOfRoute),
	// 		StartDate:             timestamppb.New(sai.StartDate),
	// 		EndDate:               timestamppb.New(sai.EndDate),
	// 		FinishDate:            finishDate,
	// 		SubmitDate:            timestamppb.New(sai.SubmitDate),
	// 		TaskImageUrl:          sai.TaskImageUrls,
	// 		TaskAnswer:            int32(sai.TaskAnswer),
	// 		Status:                int32(sai.Status),
	// 		EffectiveCall:         int32(sai.EffectiveCall),
	// 		Address:               address,
	// 		CustomerAcquisition:   ca,
	// 	},
	// }
	return
}

func (h *CrmGrpcHandler) CheckoutTaskVisitFU(ctx context.Context, req *crmService.CheckoutTaskVisitFURequest) (res *crmService.CheckoutTaskVisitFUResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.CheckTaskSalesAssignmentItemActive")
	defer span.End()

	err = h.ServicesSalesAssignmentItem.CheckoutTaskVisitFU(ctx, dto.CheckoutTaskRequest{
		Id:                  req.Id,
		Task:                int8(req.Task),
		CustomerAcquisition: req.CustomerAcquisition,
	})
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &crmService.CheckoutTaskVisitFUResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    nil,
	}

	return
}

func (h *CrmGrpcHandler) BulkCheckoutTaskVisitFU(ctx context.Context, req *crmService.BulkCheckoutTaskVisitFURequest) (res *crmService.BulkCheckoutTaskVisitFUResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.CheckTaskSalesAssignmentItemActive")
	defer span.End()

	// err = h.ServicesSalesAssignmentItem.BulkCheckoutTaskVisitFU(ctx, dto.BulkCheckoutTaskRequest{
	// 	SalesPersonId: req.SalespersonId,
	// })
	// if err != nil {
	// 	err = status.New(codes.NotFound, err.Error()).Err()
	// 	h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
	// 	return
	// }

	// res = &crmService.BulkCheckoutTaskVisitFUResponse{
	// 	Code:    int32(codes.OK),
	// 	Message: codes.OK.String(),
	// 	Data:    nil,
	// }

	return
}

func (h *CrmGrpcHandler) SubmitTaskFailed(ctx context.Context, req *crmService.SubmitTaskFailedRequest) (res *crmService.SubmitTaskFailedResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.CheckTaskSalesAssignmentItemActive")
	defer span.End()

	// var sfv *dto.SalesFailedVisitResponse
	// sfv, err = h.ServicesSalesFailedVisit.SubmitTaskFailed(ctx, dto.SalesFailedVisitRequest{
	// 	SalesAssignmentItemId: req.SalesAssignmentItemId,
	// 	FailedStatus:          req.FailedStatus,
	// 	DescriptionFailed:     req.DescriptionFailed,
	// 	FailedImage:           req.FailedImage,
	// })
	// if err != nil {
	// 	err = status.New(codes.NotFound, err.Error()).Err()
	// 	h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
	// 	return
	// }

	// res = &crmService.SubmitTaskFailedResponse{
	// 	Code:    int32(codes.OK),
	// 	Message: codes.OK.String(),
	// 	Data: &crmService.SalesFailedVisitResponse{
	// 		Id:                    sfv.ID,
	// 		SalesAssignmentItemId: sfv.SalesAssignmentItemId,
	// 		FailedStatus:          sfv.FailedStatus,
	// 		DescriptionFailed:     sfv.DescriptionFailed,
	// 		FailedImage:           sfv.FailedImage,
	// 	},
	// }

	return
}

func (h *CrmGrpcHandler) CreateSalesAssignmentItem(ctx context.Context, req *crmService.CreateSalesAssignmentItemRequest) (res *crmService.GetSalesAssignmentItemDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.CreateSalesAssignmentItem")
	defer span.End()

	// var sai *dto.SalesAssignmentItemResponse
	// sai, err = h.ServicesSalesAssignmentItem.Create(ctx, &dto.SalesAssignmentItemRequest{
	// 	SalesAssignmentID:     req.SalesAssignmentId,
	// 	SalesPersonId:         req.SalesersonId,
	// 	AddressId:             req.AddressId,
	// 	CustomerAcquisitionID: req.CustomerAcquisitionId,
	// 	Latitude:              req.Latitude,
	// 	Longitude:             req.Longitude,
	// 	Task:                  int8(req.Task),
	// 	CustomerType:          int8(req.CustomerType),
	// 	ObjectiveCodes:        req.ObjectiveCodes,
	// 	ActualDistance:        req.ActualDistance,
	// 	OutOfRoute:            int8(req.OutOfRoute),
	// 	StartDate:             req.StartDate.AsTime(),
	// 	EndDate:               req.EndDate.AsTime(),
	// 	SubmitDate:            req.SubmitDate.AsTime(),
	// 	TaskImageUrls:         strings.Split(req.TaskImageUrls, ","),
	// 	TaskAnswer:            int8(req.TaskAnswer),
	// 	Status:                int8(req.Status),
	// 	EffectiveCall:         int8(req.EffectiveCall),
	// })
	// if err != nil {
	// 	err = status.New(codes.NotFound, err.Error()).Err()
	// 	h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
	// 	return
	// }

	// res = &crmService.GetSalesAssignmentItemDetailResponse{
	// 	Data: &crmService.SalesAssignmentItem{
	// 		Id:                    sai.ID,
	// 		SalesAssignmentId:     sai.SalesAssignmentID,
	// 		SalesPersonId:         sai.SalesPerson.ID,
	// 		AddressId:             sai.AddressId,
	// 		CustomerAcquisitionId: sai.CustomerAcquisitionID,
	// 		Latitude:              sai.Latitude,
	// 		Longitude:             sai.Longitude,
	// 		Task:                  int32(sai.Task),
	// 		CustomerType:          int32(sai.CustomerType),
	// 		ObjectiveCodes:        sai.ObjectiveCodes,
	// 		ActualDistance:        sai.ActualDistance,
	// 		OutOfRoute:            int32(sai.OutOfRoute),
	// 		StartDate:             timestamppb.New(sai.StartDate),
	// 		EndDate:               timestamppb.New(sai.EndDate),
	// 		SubmitDate:            timestamppb.New(sai.SubmitDate),
	// 		TaskImageUrl:          sai.TaskImageUrls,
	// 		TaskAnswer:            int32(sai.Task),
	// 		Status:                int32(sai.Status),
	// 		EffectiveCall:         int32(sai.EffectiveCall),
	// 	},
	// }
	return
}
