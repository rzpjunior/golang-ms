package handler

import (
	context "context"

	// bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	salesService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/sales_service"
	"git.edenfarm.id/project-version3/erp-services/erp-sales-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-sales-service/internal/app/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *SalesGrpcHandler) GetSalesOrderList(ctx context.Context, req *salesService.GetSalesOrderListRequest) (res *salesService.GetSalesOrderListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "SalesGrpcHandler.GetSalesOrderList")
	defer span.End()
	var so []dto.SalesOrderResponse
	//var total int64
	so, _, err = h.ServiceSalesOrder.GetListGRPC(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	var data []*salesService.SalesOrder
	for _, salesOrder := range so {
		data = append(data, &salesService.SalesOrder{
			Id:                  salesOrder.ID,
			SalesOrderNumber:    salesOrder.SalesOrderNumber,
			SalesOrderNumberGp:  salesOrder.SalesOrderNumberGP,
			Status:              int32(salesOrder.Status),
			RequestsShipDate:    timestamppb.New(salesOrder.RequestsShipDate),
			RecognitionDate:     timestamppb.New(salesOrder.RecognitionDate),
			CreatedBy:           salesOrder.CreatedBy,
			CreatedAt:           timestamppb.New(salesOrder.CreatedAt),
			IntegrationCode:     salesOrder.IntegrationCode,
			BillingAddress:      salesOrder.BillingAddress,
			ShippingAddress:     salesOrder.ShippingAddress,
			ShippingAddressNote: salesOrder.ShippingAddressNote,
			DeliveryFee:         salesOrder.DeliveryFee,
			VouDiscAmount:       salesOrder.VouDiscAmount,
			TotalPrice:          salesOrder.TotalPrice,
			TotalCharge:         salesOrder.TotalCharge,
			TotalWeight:         salesOrder.TotalWeight,
			Note:                salesOrder.Note,
			SubDistrictIdGp:     salesOrder.SubDistrictIDGP,
			SiteIdGp:            salesOrder.SiteIDGP,
			TermPaymentSlsIdGp:  salesOrder.TermPaymentSlsIDGP,
			PaymentGroupSlsId:   salesOrder.PaymentGroupSlsID,
			RegionIdGp:          salesOrder.RegionIDGP,
			PaymentReminder:     int32(salesOrder.PaymentReminder),
			CustomerIdGp:        salesOrder.CustomerIDGP,
			CustomerPointLogId:  salesOrder.CustomerPointLogID,
			CancelType:          int32(salesOrder.CancelType),
			EdenPointCampaignId: salesOrder.EdenPointCampaignID,
			PriceLevelIdGp:      salesOrder.PriceLevelIDGP,
			WrtIdGp:             salesOrder.WrtIDGP,
			AddressIdGp:         salesOrder.AddressIDGP,
			ArchetypeIdGp:       salesOrder.ArchetypeIDGP,
		})
	}

	res = &salesService.GetSalesOrderListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *SalesGrpcHandler) GetSalesOrderDetail(ctx context.Context, req *salesService.GetSalesOrderDetailRequest) (res *salesService.GetSalesOrderDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "SalesGrpcHandler.GetSalesOrderDetail")
	defer span.End()
	var salesOrder dto.SalesOrderResponse
	//var total int64
	salesOrder, err = h.ServiceSalesOrder.GetDetailGRPC(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data *salesService.SalesOrder
	data = &salesService.SalesOrder{
		Id:                  salesOrder.ID,
		SalesOrderNumber:    salesOrder.SalesOrderNumber,
		SalesOrderNumberGp:  salesOrder.SalesOrderNumberGP,
		Status:              int32(salesOrder.Status),
		RequestsShipDate:    timestamppb.New(salesOrder.RequestsShipDate),
		RecognitionDate:     timestamppb.New(salesOrder.RecognitionDate),
		CreatedBy:           salesOrder.CreatedBy,
		CreatedAt:           timestamppb.New(salesOrder.CreatedAt),
		IntegrationCode:     salesOrder.IntegrationCode,
		BillingAddress:      salesOrder.BillingAddress,
		ShippingAddress:     salesOrder.ShippingAddress,
		ShippingAddressNote: salesOrder.ShippingAddressNote,
		DeliveryFee:         salesOrder.DeliveryFee,
		VouDiscAmount:       salesOrder.VouDiscAmount,
		TotalPrice:          salesOrder.TotalPrice,
		TotalCharge:         salesOrder.TotalCharge,
		TotalWeight:         salesOrder.TotalWeight,
		Note:                salesOrder.Note,
		SubDistrictIdGp:     salesOrder.SubDistrictIDGP,
		SiteIdGp:            salesOrder.SiteIDGP,
		TermPaymentSlsIdGp:  salesOrder.TermPaymentSlsIDGP,
		PaymentGroupSlsId:   salesOrder.PaymentGroupSlsID,
		RegionIdGp:          salesOrder.RegionIDGP,
		PaymentReminder:     int32(salesOrder.PaymentReminder),
		CustomerIdGp:        salesOrder.CustomerIDGP,
		CustomerPointLogId:  salesOrder.CustomerPointLogID,
		CancelType:          int32(salesOrder.CancelType),
		EdenPointCampaignId: salesOrder.EdenPointCampaignID,
		PriceLevelIdGp:      salesOrder.PriceLevelIDGP,
		WrtIdGp:             salesOrder.WrtIDGP,
		AddressIdGp:         salesOrder.AddressIDGP,
		ArchetypeIdGp:       salesOrder.ArchetypeIDGP,
	}

	for _, v := range salesOrder.SalesOrderItem {
		data.SalesOrderItem = append(data.SalesOrderItem, &salesService.SalesOrderItem{
			Id:               v.ID,
			SalesOrderId:     v.SalesOrderID,
			ItemIdGp:         v.ItemIDGP,
			ItemName:         v.ItemName,
			PriceTieringIdGp: v.PriceTieringIDGP,
			OrderQty:         v.OrderQty,
			UnitPrice:        v.UnitPrice,
			Subtotal:         v.Subtotal,
			Weight:           v.Weight,
			UomIdGp:          v.UomIDGP,
			UomName:          v.UomName,
			ImageUrl:         v.ImageUrl,
		})
	}
	res = &salesService.GetSalesOrderDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *SalesGrpcHandler) GetSalesOrderItemList(ctx context.Context, req *salesService.GetSalesOrderItemListRequest) (res *salesService.GetSalesOrderItemListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "SalesGrpcHandler.GetSalesOrderItemList")
	defer span.End()
	var so []dto.SalesOrderItemResponse
	//var total int64
	so, _, err = h.ServiceSalesOrder.GetListItemGRPC(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	var data []*salesService.SalesOrderItem
	for _, salesOrderItem := range so {
		data = append(data, &salesService.SalesOrderItem{
			Id:               salesOrderItem.ID,
			SalesOrderId:     salesOrderItem.SalesOrderID,
			ItemIdGp:         salesOrderItem.ItemIDGP,
			ItemName:         salesOrderItem.ItemName,
			PriceTieringIdGp: salesOrderItem.PriceTieringIDGP,
			OrderQty:         salesOrderItem.OrderQty,
			UnitPrice:        salesOrderItem.UnitPrice,
			Subtotal:         salesOrderItem.Subtotal,
			Weight:           salesOrderItem.Weight,
			UomIdGp:          salesOrderItem.UomIDGP,
			UomName:          salesOrderItem.UomName,
			ImageUrl:         salesOrderItem.ImageUrl,
		})
	}

	res = &salesService.GetSalesOrderItemListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *SalesGrpcHandler) GetSalesOrderItemDetail(ctx context.Context, req *salesService.GetSalesOrderItemDetailRequest) (res *salesService.GetSalesOrderItemDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "SalesGrpcHandler.GetSalesOrderItemDetail")
	defer span.End()
	var salesOrderItem dto.SalesOrderItemResponse
	//var total int64
	salesOrderItem, err = h.ServiceSalesOrder.GetDetailItemGRPC(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	var data *salesService.SalesOrderItem
	data = &salesService.SalesOrderItem{
		Id:               salesOrderItem.ID,
		SalesOrderId:     salesOrderItem.SalesOrderID,
		ItemIdGp:         salesOrderItem.ItemIDGP,
		ItemName:         salesOrderItem.ItemName,
		PriceTieringIdGp: salesOrderItem.PriceTieringIDGP,
		OrderQty:         salesOrderItem.OrderQty,
		UnitPrice:        salesOrderItem.UnitPrice,
		Subtotal:         salesOrderItem.Subtotal,
		Weight:           salesOrderItem.Weight,
		UomIdGp:          salesOrderItem.UomIDGP,
		UomName:          salesOrderItem.UomName,
		ImageUrl:         salesOrderItem.ImageUrl,
	}

	res = &salesService.GetSalesOrderItemDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *SalesGrpcHandler) CreateSalesOrder(ctx context.Context, req *salesService.CreateSalesOrderRequest) (res *salesService.CreateSalesOrderResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "SalesGrpcHandler.GetSalesOrderItemDetail")
	defer span.End()
	resService, err := h.ServiceSalesOrder.CreateSalesOrder(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &salesService.CreateSalesOrderResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &salesService.SalesOrder{
			Id:               resService.ID,
			SalesOrderNumber: resService.SalesOrderNumber,
			TotalCharge:      resService.TotalCharge,
		},
	}
	return
}

func (h *SalesGrpcHandler) UpdateSalesOrder(ctx context.Context, req *salesService.UpdateSalesOrderRequest) (res *salesService.UpdateSalesOrderResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "SalesGrpcHandler.UpdateSalesOrderHeader")
	defer span.End()
	_, err = h.ServiceSalesOrder.UpdateSalesOrder(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &salesService.UpdateSalesOrderResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
	}
	return
}

func (h *SalesGrpcHandler) GetSalesOrderListMobile(ctx context.Context, req *salesService.GetSalesOrderListRequest) (res *salesService.GetSalesOrderListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "SalesGrpcHandler.GetSalesOrderList")
	defer span.End()
	var so []dto.SalesOrderResponse
	var total int64
	so, total, err = h.ServiceSalesOrder.GetListGRPCMobile(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	var data []*salesService.SalesOrder
	for _, salesOrder := range so {
		data = append(data, &salesService.SalesOrder{
			Id:                  salesOrder.ID,
			SalesOrderNumber:    salesOrder.SalesOrderNumber,
			SalesOrderNumberGp:  salesOrder.SalesOrderNumberGP,
			Status:              int32(salesOrder.Status),
			RequestsShipDate:    timestamppb.New(salesOrder.RequestsShipDate),
			RecognitionDate:     timestamppb.New(salesOrder.RecognitionDate),
			CreatedBy:           salesOrder.CreatedBy,
			CreatedAt:           timestamppb.New(salesOrder.CreatedAt),
			IntegrationCode:     salesOrder.IntegrationCode,
			BillingAddress:      salesOrder.BillingAddress,
			ShippingAddress:     salesOrder.ShippingAddress,
			ShippingAddressNote: salesOrder.ShippingAddressNote,
			DeliveryFee:         salesOrder.DeliveryFee,
			VouDiscAmount:       salesOrder.VouDiscAmount,
			TotalPrice:          salesOrder.TotalPrice,
			TotalCharge:         salesOrder.TotalCharge,
			TotalWeight:         salesOrder.TotalWeight,
			Note:                salesOrder.Note,
			SubDistrictIdGp:     salesOrder.SubDistrictIDGP,
			SiteIdGp:            salesOrder.SiteIDGP,
			TermPaymentSlsIdGp:  salesOrder.TermPaymentSlsIDGP,
			PaymentGroupSlsId:   salesOrder.PaymentGroupSlsID,
			RegionIdGp:          salesOrder.RegionIDGP,
			PaymentReminder:     int32(salesOrder.PaymentReminder),
			CustomerIdGp:        salesOrder.CustomerIDGP,
			CustomerPointLogId:  salesOrder.CustomerPointLogID,
			CancelType:          int32(salesOrder.CancelType),
			EdenPointCampaignId: salesOrder.EdenPointCampaignID,
			PriceLevelIdGp:      salesOrder.PriceLevelIDGP,
			WrtIdGp:             salesOrder.WrtIDGP,
			AddressIdGp:         salesOrder.AddressIDGP,
			ArchetypeIdGp:       salesOrder.ArchetypeIDGP,
		})
	}

	res = &salesService.GetSalesOrderListResponse{
		Code:         int32(codes.OK),
		Message:      codes.OK.String(),
		Data:         data,
		TotalRecords: int32(total),
	}
	return
}

func (h *SalesGrpcHandler) GetSalesOrderFeedbackList(ctx context.Context, req *salesService.GetSalesOrderFeedbackListRequest) (res *salesService.GetSalesOrderFeedbackListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "SalesGrpcHandler.GetSalesOrderList")
	defer span.End()
	var soFeedback []dto.SalesOrderFeedback
	//var total int64
	soFeedback, _, err = h.ServiceSalesOrder.GetListSalesOrderFeedbackGRPC(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	var data []*salesService.SalesOrderFeedback
	for _, sof := range soFeedback {
		data = append(data, &salesService.SalesOrderFeedback{
			SalesOrderCode: sof.SalesOrderCode,
			DeliveryDate:   sof.DeliveryDate,
			RatingScore:    int32(sof.RatingScore),
			Description:    sof.Description,
			Tags:           sof.Tags,
			TotalCharge:    sof.TotalCharge,
			SalesOrderId:   sof.SalesOrderID,
		})
	}

	res = &salesService.GetSalesOrderFeedbackListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *SalesGrpcHandler) CreateSalesOrderFeedback(ctx context.Context, req *salesService.CreateSalesOrderFeedbackRequest) (res *salesService.CreateSalesOrderFeedbackResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "SalesGrpcHandler.GetSalesOrderItemDetail")
	defer span.End()
	_, err = h.ServiceSalesOrder.CreateSalesOrderFeedback(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &salesService.CreateSalesOrderFeedbackResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
	}
	return
}

func (h *SalesGrpcHandler) GetSalesOrderGPList(ctx context.Context, req *salesService.CreateSalesOrderFeedbackRequest) (res *salesService.CreateSalesOrderFeedbackResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "SalesGrpcHandler.GetSalesOrderItemDetail")
	defer span.End()
	_, err = h.ServiceSalesOrder.CreateSalesOrderFeedback(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &salesService.CreateSalesOrderFeedbackResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
	}
	return
}

func (h *SalesGrpcHandler) GetSalesOrderListCronJob(ctx context.Context, req *salesService.GetSalesOrderListCronjobRequest) (res *salesService.GetSalesOrderListCronjobResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "SalesGrpcHandler.GetSalesOrderListCronJob")
	defer span.End()

	var salesOrders []*model.SalesOrder
	salesOrders, err = h.ServiceSalesOrder.GetSalesOrderListCronJob(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	var data []*salesService.SalesOrderCronjob
	for _, so := range salesOrders {
		data = append(data, &salesService.SalesOrderCronjob{
			CustomerIdGp:     so.CustomerIDGP,
			SalesOrderNumber: so.SalesOrderNumber,
			Id:               so.ID,
		})
	}
	res = &salesService.GetSalesOrderListCronjobResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}

	return
}

func (h *SalesGrpcHandler) UpdateSalesOrderRemindPayment(ctx context.Context, req *salesService.UpdateSalesOrderRemindPaymentRequest) (res *salesService.UpdateSalesOrderRemindPaymentResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "SalesGrpcHandler.UpdateSalesOrderRemindPayment")
	defer span.End()

	res, err = h.ServiceSalesOrder.UpdateSalesOrderRemindPayment(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &salesService.UpdateSalesOrderRemindPaymentResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
	}

	return
}

func (h *SalesGrpcHandler) ExpiredSalesOrder(ctx context.Context, req *salesService.ExpiredSalesOrderRequest) (res *salesService.ExpiredSalesOrderResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "SalesGrpcHandler.ExpiredSalesOrder")
	defer span.End()

	res, err = h.ServiceSalesOrder.ExpiredSalesOrder(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &salesService.ExpiredSalesOrderResponse{
		Code:         int32(codes.OK),
		Message:      codes.OK.String(),
		CustomerIdGp: res.CustomerIdGp,
	}

	return
}

func (h *SalesGrpcHandler) CreateSalesOrderPaid(ctx context.Context, req *salesService.CreateSalesOrderPaidRequest) (res *salesService.CreateSalesOrderPaidResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "SalesGrpcHandler.GetSalesOrderItemDetail")
	defer span.End()

	res, err = h.ServiceSalesOrder.CreateSalesOrderPaid(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &salesService.CreateSalesOrderPaidResponse{
		Code:         int32(codes.OK),
		Message:      codes.OK.String(),
		CustomerIdGp: res.CustomerIdGp,
	}

	return
}
