package handler

import (
	context "context"

	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/dto"
	crmService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/crm_service"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *CrmGrpcHandler) CheckTaskCustomerAcquisitionActive(ctx context.Context, req *crmService.CheckTaskCustomerAcquisitionRequest) (res *crmService.CheckTaskCustomerAcquisitionResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "CrmGrpcHandler.CheckTaskCustomerAcquisitionActive")
	defer span.End()

	// var existed bool
	// existed, err = h.ServicesCustomerAcquisition.CheckActiveTask(ctx, req.SalespersonId)
	// if err != nil {
	// 	err = status.New(codes.NotFound, err.Error()).Err()
	// 	h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
	// 	return
	// }

	// res = &crmService.CheckTaskCustomerAcquisitionResponse{
	// 	Data: &crmService.BooleanResponse{
	// 		Existed: existed,
	// 	},
	// }
	return
}

func (h *CrmGrpcHandler) SubmitTaskCustomerAcquisition(ctx context.Context, req *crmService.SubmitTaskCustomerAcquisitionRequest) (res *crmService.SubmitTaskCustomerAcquisitionResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "CrmGrpcHandler.CheckTaskCustomerAcquisitionActive")
	defer span.End()

	var reqProducts []*dto.CustomerAcqProduct
	for _, reqProduct := range req.Product {
		reqProducts = append(reqProducts, &dto.CustomerAcqProduct{
			Id:  reqProduct.Id,
			Top: int8(reqProduct.Top),
		})
	}

	// var result *dto.CustomerAcquisitionResponse
	// result, err = h.ServicesCustomerAcquisition.SubmitTask(ctx, dto.SubmitTaskCustomerAcqRequest{
	// 	SalesPersonId:            req.SalespersonId,
	// 	CustomerName:             req.CustomerName,
	// 	PhoneNumber:              req.PhoneNumber,
	// 	AddressDetail:            req.AddressDetail,
	// 	FoodApp:                  int8(req.FoodApp),
	// 	UserLatitude:             req.UserLatitude,
	// 	UserLongitude:            req.UserLongitude,
	// 	PotentialRevenue:         req.PotentialRevenue,
	// 	CustomerAcquisitionPhoto: req.CustomerAcquisitionPhoto,
	// 	Product:                  reqProducts,
	// })
	// if err != nil {
	// 	err = status.New(codes.Internal, err.Error()).Err()
	// 	h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
	// 	return
	// }

	// var caItems []*crmService.CustomerAcquisitionItemResponse
	// for _, cai := range result.CustomerAcquisitionItems {
	// 	caItems = append(caItems, &crmService.CustomerAcquisitionItemResponse{
	// 		Id:                    cai.ID,
	// 		CustomerAcquisitionId: cai.CustomerAcquisitionID,
	// 		ItemId:                cai.Item.ID,
	// 		IsTop:                 int32(cai.IsTop),
	// 		CreatedAt:             timestamppb.New(cai.CreatedAt),
	// 		UpdatedAt:             timestamppb.New(cai.UpdatedAt),
	// 	})
	// }

	// res = &crmService.SubmitTaskCustomerAcquisitionResponse{
	// 	CustomerAcquisition: &crmService.CustomerAcquisitionResponse{
	// 		Id:               result.ID,
	// 		Code:             result.Code,
	// 		Task:             int32(result.Task),
	// 		Name:             result.Name,
	// 		PhoneNumber:      result.PhoneNumber,
	// 		Latitude:         result.Latitude,
	// 		Longitude:        result.Longitude,
	// 		AddressName:      result.AddressName,
	// 		FoodApp:          int32(result.FoodApp),
	// 		PotentialRevenue: result.PotentialRevenue,
	// 		TaskImageUrl:     result.TaskImageUrl,
	// 		SalespersonId:    result.Salesperson.ID,
	// 	},
	// 	CustomerAcquisitionItem: caItems,
	// }
	return
}

func (h *CrmGrpcHandler) GetCustomerAcquisitionById(ctx context.Context, req *crmService.GetCustomerAcquisitionByIdRequest) (res *crmService.GetCustomerAcquisitionDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "CrmGrpcHandler.GetCustomerAcquisitionById")
	defer span.End()

	// var item dto.CustomerAcquisitionResponse
	// item, err = h.ServicesCustomerAcquisition.GetByID(ctx, req.Id)
	// if err != nil {
	// 	err = status.New(codes.NotFound, err.Error()).Err()
	// 	h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
	// 	return
	// }

	// res = &crmService.GetCustomerAcquisitionDetailResponse{
	// 	Code:    int32(codes.OK),
	// 	Message: codes.OK.String(),
	// 	Data: &crmService.CustomerAcquisitionResponse{
	// 		Id:               item.ID,
	// 		Code:             item.Code,
	// 		Name:             item.Name,
	// 		PhoneNumber:      item.PhoneNumber,
	// 		AddressName:      item.AddressName,
	// 		FoodApp:          int32(item.FoodApp),
	// 		PotentialRevenue: item.PotentialRevenue,
	// 		SalespersonId:    item.Salesperson.ID,
	// 		TerritoryId:      item.Territory.ID,
	// 		Latitude:         item.Latitude,
	// 		Longitude:        item.Longitude,
	// 		Task:             int32(item.Task),
	// 		FinishDate:       timestamppb.New(item.FinishDate),
	// 		SubmitDate:       timestamppb.New(item.SubmitDate),
	// 		TaskImageUrl:     item.TaskImageUrl,
	// 		Status:           int32(item.Status),
	// 		CreatedAt:        timestamppb.New(item.CreatedAt),
	// 		UpdatedAt:        timestamppb.New(item.UpdatedAt),
	// 	},
	// }
	return
}

func (h *CrmGrpcHandler) GetCustomerAcquisitionList(ctx context.Context, req *crmService.GetCustomerAcquisitionListRequest) (res *crmService.GetCustomerAcquisitionListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "CrmGrpcHandler.GetCustomerAcquisitionList")
	defer span.End()

	// var items []*dto.CustomerAcquisitionResponse
	// items, _, err = h.ServicesCustomerAcquisition.Get(ctx, int(req.Offset), int(req.Limit), int(req.Status), req.Search, req.OrderBy, req.TeritoryId, req.SalespersonId, req.SubmitDateFrom.AsTime(), req.SubmitDateTo.AsTime())
	// if err != nil {
	// 	err = status.New(codes.NotFound, err.Error()).Err()
	// 	h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
	// 	return
	// }
	var data []*crmService.CustomerAcquisitionResponse
	// for _, item := range items {
	// 	data = append(data, &crmService.CustomerAcquisitionResponse{
	// 		Id:               item.ID,
	// 		Code:             item.Code,
	// 		Name:             item.Name,
	// 		PhoneNumber:      item.PhoneNumber,
	// 		AddressName:      item.AddressName,
	// 		FoodApp:          int32(item.FoodApp),
	// 		PotentialRevenue: item.PotentialRevenue,
	// 		SalespersonId:    item.Salesperson.ID,
	// 		TerritoryId:      item.Territory.ID,
	// 		Latitude:         item.Latitude,
	// 		Longitude:        item.Longitude,
	// 		Task:             int32(item.Task),
	// 		FinishDate:       timestamppb.New(item.FinishDate),
	// 		SubmitDate:       timestamppb.New(item.SubmitDate),
	// 		TaskImageUrl:     item.TaskImageUrl,
	// 		Status:           int32(item.Status),
	// 		CreatedAt:        timestamppb.New(item.CreatedAt),
	// 		UpdatedAt:        timestamppb.New(item.UpdatedAt),
	// 	})
	// }

	res = &crmService.GetCustomerAcquisitionListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *CrmGrpcHandler) GetCustomerAcquisitionListWithExcludedIds(ctx context.Context, req *crmService.GetCustomerAcquisitionListWithExcludedIdsRequest) (res *crmService.GetCustomerAcquisitionListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "CrmGrpcHandler.GetCustomerAcquisitionListWithExcludedIds")
	defer span.End()

	// var items []*dto.CustomerAcquisitionResponse
	// items, _, err = h.ServicesCustomerAcquisition.GetWithExcludedIds(ctx, int(req.Offset), int(req.Limit), int(req.Status), req.Search, req.OrderBy, req.TeritoryId, req.SalespersonId, req.SubmitDateFrom.AsTime(), req.SubmitDateTo.AsTime(), req.ExcludedIds)
	// if err != nil {
	// 	err = status.New(codes.NotFound, err.Error()).Err()
	// 	h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
	// 	return
	// }
	var data []*crmService.CustomerAcquisitionResponse
	// for _, item := range items {
	// 	data = append(data, &crmService.CustomerAcquisitionResponse{
	// 		Id:               item.ID,
	// 		Code:             item.Code,
	// 		Name:             item.Name,
	// 		PhoneNumber:      item.PhoneNumber,
	// 		AddressName:      item.AddressName,
	// 		FoodApp:          int32(item.FoodApp),
	// 		PotentialRevenue: item.PotentialRevenue,
	// 		SalespersonId:    item.Salesperson.ID,
	// 		TerritoryId:      item.Territory.ID,
	// 		Latitude:         item.Latitude,
	// 		Longitude:        item.Longitude,
	// 		Task:             int32(item.Task),
	// 		FinishDate:       timestamppb.New(item.FinishDate),
	// 		SubmitDate:       timestamppb.New(item.SubmitDate),
	// 		TaskImageUrl:     item.TaskImageUrl,
	// 		Status:           int32(item.Status),
	// 		CreatedAt:        timestamppb.New(item.CreatedAt),
	// 		UpdatedAt:        timestamppb.New(item.UpdatedAt),
	// 	})
	// }

	res = &crmService.GetCustomerAcquisitionListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *CrmGrpcHandler) GetCountCustomerAcquisition(ctx context.Context, req *crmService.GetCountCustomerAcquisitionRequest) (res *crmService.GetCountCustomerAcquisitionResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "CrmGrpcHandler.GetCustomerAcquisitionListWithExcludedIds")
	defer span.End()

	var count int64
	count, err = h.ServicesCustomerAcquisition.CountCustomerAcq(ctx, req.SalespersonId, req.SubmitDateFrom.AsTime(), req.SubmitDateTo.AsTime())
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	res = &crmService.GetCountCustomerAcquisitionResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &crmService.CountCustomerAcquisitionResponse{
			Count: count,
		},
	}
	return
}
