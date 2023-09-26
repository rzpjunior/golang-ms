package service

import (
	"context"
	"regexp"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-helper-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-helper-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/site_service"
)

type IHelperAppService interface {
	Login(ctx context.Context, req dto.HelperAppLoginRequest) (res *dto.HelperAppLoginResponse, err error)
	GetPickingOrder(ctx context.Context, req dto.HelperAppGetPickingOrderRequest) (res *dto.HelperAppGetPickingOrderResponse, err error)
	PickerWidget(ctx context.Context, req dto.HelperAppPickerWidgetRequest) (res *dto.HelperAppPickerWidgetResponse, err error)
	GetPickingOrderProducts(ctx context.Context, req dto.HelperAppGetPickingOrderProductsRequest) (res *dto.HelperAppGetPickingOrderProductsResponse, err error)
	GetPickingOrderProductsSalesOrder(ctx context.Context, docNumber string, itemNumber string) (res *dto.HelperAppGetPickingOrderProductsSalesOrder, err error)
	StartPickingOrder(ctx context.Context, req dto.HelperAppDocNumberRequest) (res *dto.HelperAppSuccessResponse, err error)
	SubmitPicking(ctx context.Context, req dto.HelperAppSubmitPickingRequest) (res *dto.HelperAppSuccessResponse, err error)
	GetSalesOrderPicking(ctx context.Context, docNumber string) (res *dto.HelperAppGetSalesOrderPickingResponse, err error)
	GetSalesOrderPickingDetail(ctx context.Context, sopNumber string) (res *dto.HelperAppGetSalesOrderPickingDetailResponse, err error)
	SubmitSalesOrder(ctx context.Context, req dto.HelperAppSubmitSalesOrderRequest) (res *dto.HelperAppSuccessResponse, err error)
	History(ctx context.Context, req dto.HelperAppHistoryRequest) (res *dto.HelperAppHistoryResponse, err error)
	SPVGetSalesOrderList(ctx context.Context, req dto.HelperAppGetSalesOrderToCheckRequest) (res *dto.HelperAppGetSalesOrderToCheckResponse, err error)
	SPVWidget(ctx context.Context, req dto.HelperAppSPVWidgetRequest) (res *dto.HelperAppSPVWidgetResponse, err error)
	SPVGetSalesOrderDetail(ctx context.Context, sopNumber string) (res *dto.HelperAppGetSalesOrderToCheckDetailResponse, err error)
	SPVRejectSalesOrder(ctx context.Context, req dto.HelperAppSopNumberRequest) (res *dto.HelperAppSuccessResponse, err error)
	SPVAcceptSalesOrder(ctx context.Context, req dto.HelperAppSopNumberRequest) (res *dto.HelperAppSuccessResponse, err error)
	SPVGetWrtMonitoring(ctx context.Context, req dto.HelperAppGetWrtMonitoringRequest) (res *dto.HelperAppGetWrtMonitoringResponse, err error)
	SPVGetWrtMonitoringDetail(ctx context.Context, req dto.HelperAppWrtMonitoringDetailRequest) (res *dto.HelperAppGetWrtMonitoringDetailResponse, err error)
	CheckerGetSalesOrderList(ctx context.Context, req dto.HelperAppGetSalesOrderToCheckRequest) (res *dto.HelperAppGetSalesOrderToCheckResponse, err error)
	CheckerWidget(ctx context.Context, req dto.HelperAppCheckerWidgetRequest) (res *dto.HelperAppCheckerWidgetResponse, err error)
	CheckerGetSalesOrderDetail(ctx context.Context, req dto.HelperAppCheckerGetSalesOrderDetailRequest) (res *dto.HelperAppGetSalesOrderToCheckDetailResponse, err error)
	CheckerStartChecking(ctx context.Context, req dto.HelperAppCheckerStartCheckingRequest) (res *dto.HelperAppSuccessResponse, err error)
	CheckerSubmitChecking(ctx context.Context, req dto.HelperAppCheckerSubmitCheckingRequest) (res *dto.HelperAppSuccessResponse, err error)
	CheckerRejectSalesOrder(ctx context.Context, req dto.HelperAppCheckerRejectSalesOrderRequest) (res *dto.HelperAppSuccessResponse, err error)
	CheckerGetDeliveryKoli(ctx context.Context, sopNumber string) (res *dto.HelperAppCheckerGetDeliveryKoliResponse, err error)
	CheckerAcceptSalesOrder(ctx context.Context, req dto.HelperAppCheckerAcceptSalesOrderRequest) (res *dto.HelperAppCheckerAcceptSalesOrderResponse, err error)
	CheckerHistory(ctx context.Context, req dto.HelperAppCheckerHistoryRequest) (res *dto.HelperAppCheckerHistoryResponse, err error)
	CheckerHistoryDetail(ctx context.Context, req dto.HelperAppCheckerHistoryDetailRequest) (res *dto.HelperAppCheckerHistoryDetailResponse, err error)
}

type HelperAppService struct {
	opt opt.Options
}

func NewServiceHelperApp() IHelperAppService {
	return &HelperAppService{
		opt: global.Setup.Common,
	}
}

func (s *HelperAppService) Login(ctx context.Context, req dto.HelperAppLoginRequest) (res *dto.HelperAppLoginResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "HelperAppService.Login")
	defer span.End()

	var login *site_service.LoginHelperResponse
	if login, err = s.opt.Client.SiteServiceGrpc.LoginHelper(ctx, &site_service.LoginHelperRequest{
		Email:         req.Email,
		Password:      req.Password,
		Timezone:      req.Timezone,
		FirebaseToken: req.FirebaseToken,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("site", "Please recheck Email or Password")
		return
	}

	res = &dto.HelperAppLoginResponse{
		Code:    login.Code,
		Message: login.Message,
		User: &dto.HelperAppUser{
			Id:       login.User.Id,
			Name:     login.User.Name,
			SiteId:   login.User.SiteId,
			SiteName: login.User.SiteName,
			RoleName: login.User.RoleName,
		},
		Token:         login.Token,
		FirebaseToken: login.FirebaseToken,
	}

	return
}

func (s *HelperAppService) GetPickingOrder(ctx context.Context, req dto.HelperAppGetPickingOrderRequest) (res *dto.HelperAppGetPickingOrderResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "HelperAppService.GetPickingOrder")
	defer span.End()

	var pickingOrder *site_service.GetPickingOrderHeaderResponse
	if pickingOrder, err = s.opt.Client.SiteServiceGrpc.GetPickingOrderHeader(ctx, &site_service.GetPickingOrderHeaderRequest{
		Limit:            int32(req.Limit),
		Offset:           int32(req.Offset),
		Locncode:         req.LocationCode,
		Sopnumbe:         req.SopNumber,
		Docnumbr:         req.DocNumber,
		Itemnmbr:         req.ItemNumber,
		GnlHelperId:      req.HelperId,
		Custname:         req.CustomerName,
		WmsPickingStatus: int32(req.Status),
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("site", "Get Picking Order Header")
		return
	}

	var pickingOrderResponse []*dto.HelperAppGetPickingOrderPickingOrder
	for _, v := range pickingOrder.Data {
		pickingOrderResponse = append(pickingOrderResponse, &dto.HelperAppGetPickingOrderPickingOrder{
			Id:              v.DocNumber,
			DocDate:         v.DocDate,
			PickerId:        v.PickerId,
			Status:          int8(v.Status),
			TotalSalesOrder: v.TotalSalesOrder,
			Note:            v.Note,
		})
	}

	res = &dto.HelperAppGetPickingOrderResponse{
		PickingOrder: pickingOrderResponse,
	}

	return
}

func (s *HelperAppService) PickerWidget(ctx context.Context, req dto.HelperAppPickerWidgetRequest) (res *dto.HelperAppPickerWidgetResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "HelperAppService.PickerWidget")
	defer span.End()

	var widget *site_service.PickerWidgetResponse
	if widget, err = s.opt.Client.SiteServiceGrpc.PickerWidget(ctx, &site_service.PickerWidgetRequest{
		GnlHelperId: req.HelperId,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("site", "Picker Widget")
		return
	}

	res = &dto.HelperAppPickerWidgetResponse{
		TotalSalesOrder:             widget.TotalSalesOrder,
		TotalNew:                    widget.TotalNew,
		TotalOnProgress:             widget.TotalOnProgress,
		TotalOnProgressPercentage:   widget.TotalOnProgressPercentage,
		TotalPicked:                 widget.TotalPicked,
		TotalPickedPercentage:       widget.TotalPickedPercentage,
		TotalNeedApproval:           widget.TotalNeedApproval,
		TotalNeedApprovalPercentage: widget.TotalNeedApprovalPercentage,
	}

	return
}

func (s *HelperAppService) GetPickingOrderProducts(ctx context.Context, req dto.HelperAppGetPickingOrderProductsRequest) (res *dto.HelperAppGetPickingOrderProductsResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "HelperAppService.GetPickingOrderProducts")
	defer span.End()

	var pickingOrderDetail *site_service.GetPickingOrderDetailResponse
	if pickingOrderDetail, err = s.opt.Client.SiteServiceGrpc.GetPickingOrderDetail(ctx, &site_service.GetPickingOrderDetailRequest{
		Id:       req.DocNumber,
		ItemName: req.ItemNameSearch,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("site", "Get Picking Order Detail")
		return

	}

	var productsResponse []*dto.HelperAppAggregatedProduct
	for _, v := range pickingOrderDetail.Data.Product {
		productsResponse = append(productsResponse, &dto.HelperAppAggregatedProduct{
			ItemNumber:      v.ItemNumber,
			ItemName:        v.ItemName,
			Picture:         v.Picture,
			UomDescription:  v.UomDescription,
			TotalOrderQty:   v.TotalOrderQty,
			TotalPickedQty:  v.TotalPickedQty,
			TotalSalesOrder: v.TotalSalesOrder,
			Status:          int8(v.Status),
		})
	}

	res = &dto.HelperAppGetPickingOrderProductsResponse{
		Id:       pickingOrderDetail.Data.DocNumber,
		Status:   int8(pickingOrderDetail.Data.Status),
		Products: productsResponse,
	}

	return
}

func (s *HelperAppService) GetPickingOrderProductsSalesOrder(ctx context.Context, docNumber string, itemNumber string) (res *dto.HelperAppGetPickingOrderProductsSalesOrder, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "HelperAppService.GetPickingOrderProductsSalesOrder")
	defer span.End()

	var aggregatedProductSalesOrders *site_service.GetAggregatedProductSalesOrderResponse
	if aggregatedProductSalesOrders, err = s.opt.Client.SiteServiceGrpc.GetAggregatedProductSalesOrder(ctx, &site_service.GetAggregatedProductSalesOrderRequest{
		Id:         docNumber,
		ItemNumber: itemNumber,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("site", "Get Aggregated Product Sales Order")
		return
	}

	var salesOrdersResponse []*dto.HelperAppAggregatedProductSO
	for _, v := range aggregatedProductSalesOrders.Data.SalesOrder {
		salesOrdersResponse = append(salesOrdersResponse, &dto.HelperAppAggregatedProductSO{
			Id:            v.Id,
			SopNumber:     v.SopNumber,
			MerchantName:  v.MerchantName,
			Wrt:           v.Wrt,
			OrderQty:      v.OrderQty,
			PickedQty:     v.PickedQty,
			UnfulfillNote: v.UnfulfillNote,
			Status:        int8(v.Status),
		})
	}

	res = &dto.HelperAppGetPickingOrderProductsSalesOrder{
		ItemId:         aggregatedProductSalesOrders.Data.ItemNumber,
		ItemName:       aggregatedProductSalesOrders.Data.ItemName,
		UomDescription: aggregatedProductSalesOrders.Data.UomDescription,
		Picture:        aggregatedProductSalesOrders.Data.Picture,
		SalesOrders:    salesOrdersResponse,
	}
	return
}

func (s *HelperAppService) StartPickingOrder(ctx context.Context, req dto.HelperAppDocNumberRequest) (res *dto.HelperAppSuccessResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "HelperAppService.StartPickingOrder")
	defer span.End()

	var pickingOrder *site_service.SuccessResponse
	if pickingOrder, err = s.opt.Client.SiteServiceGrpc.StartPickingOrder(ctx, &site_service.StartPickingOrderRequest{
		DocNumber: req.DocNumber,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("site", "Start Picking Order")
		return
	}

	res = &dto.HelperAppSuccessResponse{
		Success: pickingOrder.Success,
	}

	return
}

func (s *HelperAppService) SubmitPicking(ctx context.Context, req dto.HelperAppSubmitPickingRequest) (res *dto.HelperAppSuccessResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "HelperAppService.SubmitPicking")
	defer span.End()

	var response *site_service.SuccessResponse

	var requestPicking []*site_service.SubmitPickingModel
	for _, v := range req.Request {
		requestPicking = append(requestPicking, &site_service.SubmitPickingModel{
			Id:            v.Id,
			PickQty:       v.PickQty,
			UnfulfillNote: v.UnfulfillNote,
		})
	}

	if response, err = s.opt.Client.SiteServiceGrpc.SubmitPicking(ctx, &site_service.SubmitPickingRequest{
		Request: requestPicking,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("site", "Submit Picking")
		return
	}

	res = &dto.HelperAppSuccessResponse{
		Success: response.Success,
	}

	return
}

func (s *HelperAppService) GetSalesOrderPicking(ctx context.Context, docNumber string) (res *dto.HelperAppGetSalesOrderPickingResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "HelperAppService.GetSalesOrderPicking")
	defer span.End()

	var response *site_service.GetSalesOrderPickingResponse
	if response, err = s.opt.Client.SiteServiceGrpc.GetSalesOrderPicking(ctx, &site_service.GetSalesOrderPickingRequest{
		DocNumber: docNumber,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("site", "Get Sales Order Picking")
		return
	}

	var salesOrderResponse []*dto.HelperAppSalesOrderPickingModel
	for _, v := range response.Data {
		salesOrderResponse = append(salesOrderResponse, &dto.HelperAppSalesOrderPickingModel{
			SopNumber:        v.SopNumber,
			MerchantName:     v.MerchantName,
			SopNote:          v.SopNote,
			TotalKoli:        v.TotalKoli,
			Status:           int8(v.Status),
			ReadyToPack:      v.ReadyToPack,
			ContainUnfulfill: v.ContainUnfulfill,
		})
	}

	res = &dto.HelperAppGetSalesOrderPickingResponse{
		SalesOrder: salesOrderResponse,
	}

	return
}

func (s *HelperAppService) GetSalesOrderPickingDetail(ctx context.Context, sopNumber string) (res *dto.HelperAppGetSalesOrderPickingDetailResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "HelperAppService.GetSalesOrderPickingDetail")
	defer span.End()

	var response *site_service.GetSalesOrderPickingDetailResponse
	if response, err = s.opt.Client.SiteServiceGrpc.GetSalesOrderPickingDetail(ctx, &site_service.GetSalesOrderPickingDetailRequest{
		SopNumber: sopNumber,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("site", "Get Sales Order Picking Detail")
		return
	}

	var itemResponse []*dto.HelperAppGetSalesOrderPickingDetailItem
	for _, v := range response.Item {
		itemResponse = append(itemResponse, &dto.HelperAppGetSalesOrderPickingDetailItem{
			Id:                   v.Id,
			PickingOrderAssignId: v.PickingOrderAssignId,
			ItemNumber:           v.ItemNumber,
			ItemName:             v.ItemName,
			Picture:              v.Picture,
			OrderQty:             v.OrderQty,
			PickQty:              v.PickQty,
			CheckQty:             v.CheckQty,
			ExcessQty:            v.ExcessQty,
			UnfulfillNote:        v.UnfulfillNote,
			Uom:                  v.Uom,
			Status:               int8(v.Status),
		})
	}

	res = &dto.HelperAppGetSalesOrderPickingDetailResponse{
		SopNumber:           response.SopNumber,
		MerchantName:        response.MerchantName,
		Wrt:                 response.Wrt,
		DeliveryDate:        response.DeliveryDate,
		TotalKoli:           response.TotalKoli,
		TotalItemOnProgress: response.TotalItemOnProgress,
		TotalItem:           response.TotalItem,
		SopNote:             response.SopNote,
		Status:              int64(response.Status),
		Item:                itemResponse,
		HelperName:          response.HelperName,
		HelperId:            response.HelperId,
	}

	return
}

func (s *HelperAppService) SubmitSalesOrder(ctx context.Context, req dto.HelperAppSubmitSalesOrderRequest) (res *dto.HelperAppSuccessResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "HelperAppService. SubmitSalesOrder")
	defer span.End()

	var response *site_service.SuccessResponse

	var requestKoli []*site_service.RequestDeliveryKoli
	for _, v := range req.Koli {
		requestKoli = append(requestKoli, &site_service.RequestDeliveryKoli{
			Id:       v.Id,
			Quantity: v.Quantity,
		})
	}

	if response, err = s.opt.Client.SiteServiceGrpc.SubmitSalesOrder(ctx, &site_service.SubmitSalesOrderRequest{
		SopNumber: req.SopNumber,
		Request:   requestKoli,
		PickerId:  req.PickerId,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("site", "Submit Sales Order")
		return
	}

	res = &dto.HelperAppSuccessResponse{
		Success: response.Success,
	}

	return
}

func (s *HelperAppService) History(ctx context.Context, req dto.HelperAppHistoryRequest) (res *dto.HelperAppHistoryResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "HelperAppService. History")
	defer span.End()

	var response *site_service.HistoryResponse
	if response, err = s.opt.Client.SiteServiceGrpc.History(ctx, &site_service.HistoryRequest{
		Limit:     int32(req.Limit),
		Offset:    int32(req.Offset),
		PickerId:  req.PickerId,
		SopNumber: req.SopNumber,
		Custname:  req.CustomerName,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("site", "History")
		return
	}

	var salesOrderResponse []*dto.HelperAppSalesOrderPickingModel
	for _, v := range response.Data {
		salesOrderResponse = append(salesOrderResponse, &dto.HelperAppSalesOrderPickingModel{
			SopNumber:    v.SopNumber,
			MerchantName: v.MerchantName,
			SopNote:      v.SopNote,
			TotalKoli:    v.TotalKoli,
			Status:       int8(v.Status),
			CountPrintDo: v.CountPrintDo,
			CountPrintSi: v.CountPrintSi,
		})
	}

	res = &dto.HelperAppHistoryResponse{
		SalesOrder: salesOrderResponse,
	}

	return
}

func (s *HelperAppService) SPVGetSalesOrderList(ctx context.Context, req dto.HelperAppGetSalesOrderToCheckRequest) (res *dto.HelperAppGetSalesOrderToCheckResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "HelperAppService.SPVGetSalesOrderList")
	defer span.End()

	var response *site_service.GetSalesOrderToCheckResponse
	if response, err = s.opt.Client.SiteServiceGrpc.GetSalesOrderToCheck(ctx, &site_service.GetSalesOrderToCheckRequest{
		Offset:    int32(req.Offset),
		Limit:     int32(req.Limit),
		SiteId:    req.SiteId,
		SopNumber: req.SopNumber,
		Status:    []int32{35},
		Custname:  req.CustomerName,
		WrtIds:    req.WrtIDs,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("site", "Get Sales Order To Check")
		return
	}

	var salesOrderResponse []*dto.HelperAppGetSalesOrderToCheckSOModel
	for _, v := range response.Data {
		salesOrderResponse = append(salesOrderResponse, &dto.HelperAppGetSalesOrderToCheckSOModel{
			SopNumber:           v.SopNumber,
			MerchantName:        v.MerchantName,
			DeliveryDate:        v.DeliveryDate,
			Wrt:                 v.Wrt,
			SopNote:             v.SopNote,
			TotalItemOnProgress: v.TotalItemOnProgress,
			TotalItem:           v.TotalItem,
			TotalKoli:           v.TotalKoli,
			CheckerName:         v.CheckerName,
			PickerName:          v.PickerName,
			Status:              int8(v.Status),
		})
	}

	res = &dto.HelperAppGetSalesOrderToCheckResponse{
		SalesOrder: salesOrderResponse,
	}

	return
}

func (s *HelperAppService) SPVWidget(ctx context.Context, req dto.HelperAppSPVWidgetRequest) (res *dto.HelperAppSPVWidgetResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "HelperAppService.SPVWidget")
	defer span.End()

	var widget *site_service.SPVWidgetResponse
	if widget, err = s.opt.Client.SiteServiceGrpc.SPVWidget(ctx, &site_service.SPVWidgetRequest{
		SiteIdGp: req.SiteIdGp,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("site", "SPV Widget")
		return
	}

	res = &dto.HelperAppSPVWidgetResponse{
		TotalSalesOrder:             widget.TotalSalesOrder,
		TotalNew:                    widget.TotalNew,
		TotalOnProgress:             widget.TotalOnProgress,
		TotalOnProgressPercentage:   widget.TotalOnProgressPercentage,
		TotalNeedApproval:           widget.TotalNeedApproval,
		TotalNeedApprovalPercentage: widget.TotalNeedApprovalPercentage,
		TotalFinished:               widget.TotalFinished,
		TotalFinishedPercentage:     widget.TotalFinishedPercentage,
	}

	return
}

func (s *HelperAppService) SPVGetSalesOrderDetail(ctx context.Context, sopNumber string) (res *dto.HelperAppGetSalesOrderToCheckDetailResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "HelperAppService.SPVGetSalesOrderDetail")
	defer span.End()

	var response *site_service.GetSalesOrderToCheckDetailResponse
	if response, err = s.opt.Client.SiteServiceGrpc.SPVGetSalesOrderToCheckDetail(ctx, &site_service.GetSalesOrderToCheckDetailRequest{
		SopNumber: sopNumber,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("site", "Spv Get Sales Order To Check Detail")
		return
	}

	var itemResponse []*dto.HelperAppGetSalesOrderPickingDetailItem
	for _, v := range response.Item {
		itemResponse = append(itemResponse, &dto.HelperAppGetSalesOrderPickingDetailItem{
			Id:                   v.Id,
			PickingOrderAssignId: v.PickingOrderAssignId,
			ItemNumber:           v.ItemNumber,
			ItemName:             v.ItemName,
			Picture:              v.Picture,
			OrderQty:             v.OrderQty,
			PickQty:              v.PickQty,
			CheckQty:             v.CheckQty,
			ExcessQty:            v.ExcessQty,
			UnfulfillNote:        v.UnfulfillNote,
			Uom:                  v.Uom,
			Status:               int8(v.Status),
		})
	}

	res = &dto.HelperAppGetSalesOrderToCheckDetailResponse{
		SopNumber:           response.SopNumber,
		MerchantName:        response.MerchantName,
		DeliveryDate:        response.DeliveryDate,
		Wrt:                 response.Wrt,
		SopNote:             response.SopNote,
		TotalItemOnProgress: response.TotalItemOnProgress,
		TotalItem:           response.TotalItem,
		TotalKoli:           response.TotalKoli,
		PickerName:          response.PickerName,
		Item:                itemResponse,
		Status:              int8(response.Status),
	}

	return
}

func (s *HelperAppService) SPVRejectSalesOrder(ctx context.Context, req dto.HelperAppSopNumberRequest) (res *dto.HelperAppSuccessResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "HelperAppService.SPVRejectSalesOrder")
	defer span.End()

	var response *site_service.SuccessResponse
	if response, err = s.opt.Client.SiteServiceGrpc.SPVRejectSalesOrder(ctx, &site_service.SPVRejectSalesOrderRequest{
		SopNumber: req.SopNumber,
		SpvId:     req.SpvId,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("site", "Spv Reject Sales Order")
		return
	}

	res = &dto.HelperAppSuccessResponse{
		Success: response.Success,
	}

	return
}

func (s *HelperAppService) SPVAcceptSalesOrder(ctx context.Context, req dto.HelperAppSopNumberRequest) (res *dto.HelperAppSuccessResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "HelperAppService.SPVAcceptSalesOrder")
	defer span.End()

	var response *site_service.SuccessResponse
	if response, err = s.opt.Client.SiteServiceGrpc.SPVAcceptSalesOrder(ctx, &site_service.SPVAcceptSalesOrderRequest{
		SopNumber: req.SopNumber,
		SpvId:     req.SpvId,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("site", "Spv Accept Sales Order")
		return
	}

	res = &dto.HelperAppSuccessResponse{
		Success: response.Success,
	}

	return
}

func (s *HelperAppService) SPVGetWrtMonitoring(ctx context.Context, req dto.HelperAppGetWrtMonitoringRequest) (res *dto.HelperAppGetWrtMonitoringResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "HelperAppService.SPVGetWrtMonitoring")
	defer span.End()

	var response *site_service.GetWrtMonitoringListResponse
	if response, err = s.opt.Client.SiteServiceGrpc.SPVWrtMonitoring(ctx, &site_service.GetWrtMonitoringListRequest{
		SiteId:   req.SiteId,
		Type:     req.Type,
		HelperId: req.HelperId,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("site", "SPV Wrt Monitoring")
		return
	}

	var responseData []*dto.HelperAppWrtMonitoring
	for _, v := range response.Data {
		responseData = append(responseData, &dto.HelperAppWrtMonitoring{
			WrtId:                v.WrtId,
			WrtDescription:       v.WrtDesc,
			CountSalesOrder:      v.CountSo,
			OnProgress:           v.OnProgress,
			OnProgressPercentage: v.OnProgressPercentage,
			Finished:             v.Finished,
			FinishedPercentage:   v.FinishedPercentage,
		})
	}

	res = &dto.HelperAppGetWrtMonitoringResponse{
		Data: responseData,
	}

	return
}

func (s *HelperAppService) SPVGetWrtMonitoringDetail(ctx context.Context, req dto.HelperAppWrtMonitoringDetailRequest) (res *dto.HelperAppGetWrtMonitoringDetailResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "HelperAppService.SPVGetWrtMonitoringDetail")
	defer span.End()

	var response *site_service.GetWrtMonitoringDetailResponse
	if response, err = s.opt.Client.SiteServiceGrpc.SPVWrtMonitoringDetail(ctx, &site_service.GetWrtMonitoringDetailRequest{
		SiteId:   req.SiteId,
		WrtId:    req.WrtId,
		Type:     req.Type,
		HelperId: req.HelperId,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("site", "SPV Wrt Monitoring Detail")
		return
	}

	var responseData []*dto.HelperAppWrtMonitoringDetail
	for _, v := range response.Data {
		responseData = append(responseData, &dto.HelperAppWrtMonitoringDetail{
			SopNumber:    v.SopNumber,
			MerchantName: v.MerchantName,
			TotalKoli:    v.TotalKoli,
			HelperCode:   v.HelperCode,
			HelperName:   v.HelperName,
			Status:       int8(v.Status),
		})
	}

	res = &dto.HelperAppGetWrtMonitoringDetailResponse{
		SalesOrder: responseData,
	}

	return
}

func (s *HelperAppService) CheckerGetSalesOrderList(ctx context.Context, req dto.HelperAppGetSalesOrderToCheckRequest) (res *dto.HelperAppGetSalesOrderToCheckResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "HelperAppService.CheckerGetSalesOrderList")
	defer span.End()

	var status32 []int32
	if len(req.Statuses) != 0 {
		for _, v := range req.Statuses {
			status32 = append(status32, int32(v))
		}
	}

	var response *site_service.GetSalesOrderToCheckResponse
	if response, err = s.opt.Client.SiteServiceGrpc.GetSalesOrderToCheck(ctx, &site_service.GetSalesOrderToCheckRequest{
		Offset:    int32(req.Offset),
		Limit:     int32(req.Limit),
		SiteId:    req.SiteId,
		SopNumber: req.SopNumber,
		Status:    status32,
		WrtIds:    req.WrtIDs,
		Custname:  req.CustomerName,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("site", "Get Sales Order To Check")
		return
	}

	var salesOrderResponse []*dto.HelperAppGetSalesOrderToCheckSOModel
	for _, v := range response.Data {
		salesOrderResponse = append(salesOrderResponse, &dto.HelperAppGetSalesOrderToCheckSOModel{
			SopNumber:           v.SopNumber,
			MerchantName:        v.MerchantName,
			DeliveryDate:        v.DeliveryDate,
			Wrt:                 v.Wrt,
			SopNote:             v.SopNote,
			TotalItemOnProgress: v.TotalItemOnProgress,
			TotalItem:           v.TotalItem,
			TotalKoli:           v.TotalKoli,
			CheckerName:         v.CheckerName,
			PickerName:          v.PickerName,
			Status:              int8(v.Status),
		})
	}

	res = &dto.HelperAppGetSalesOrderToCheckResponse{
		SalesOrder: salesOrderResponse,
	}

	return
}

func (s *HelperAppService) CheckerWidget(ctx context.Context, req dto.HelperAppCheckerWidgetRequest) (res *dto.HelperAppCheckerWidgetResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "HelperAppService.CheckerWidget")
	defer span.End()

	var widget *site_service.CheckerWidgetResponse
	if widget, err = s.opt.Client.SiteServiceGrpc.CheckerWidget(ctx, &site_service.CheckerWidgetRequest{
		CheckerId: req.CheckerId,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("site", "Checker Widget")
		return
	}

	res = &dto.HelperAppCheckerWidgetResponse{
		TotalSalesOrder:         widget.TotalSalesOrder,
		TotalPicked:             widget.TotalPicked,
		TotalChecking:           widget.TotalChecking,
		TotalCheckingPercentage: widget.TotalCheckingPercentage,
		TotalFinished:           widget.TotalFinished,
		TotalFinishedPercentage: widget.TotalFinishedPercentage,
	}

	return
}

func (s *HelperAppService) CheckerGetSalesOrderDetail(ctx context.Context, req dto.HelperAppCheckerGetSalesOrderDetailRequest) (res *dto.HelperAppGetSalesOrderToCheckDetailResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "HelperAppService.CheckerGetSalesOrderDetail")
	defer span.End()

	var response *site_service.GetSalesOrderToCheckDetailResponse
	if response, err = s.opt.Client.SiteServiceGrpc.CheckerGetSalesOrderToCheckDetail(ctx, &site_service.GetSalesOrderToCheckDetailRequest{
		SopNumber: req.SopNumber,
		CheckerId: req.CheckerId,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("site", "Checker Get Sales Order To Check Detail")
		return
	}

	var itemResponse []*dto.HelperAppGetSalesOrderPickingDetailItem
	for _, v := range response.Item {
		itemResponse = append(itemResponse, &dto.HelperAppGetSalesOrderPickingDetailItem{
			Id:                   v.Id,
			PickingOrderAssignId: v.PickingOrderAssignId,
			ItemNumber:           v.ItemNumber,
			ItemName:             v.ItemName,
			Picture:              v.Picture,
			OrderQty:             v.OrderQty,
			PickQty:              v.PickQty,
			CheckQty:             v.CheckQty,
			ExcessQty:            v.ExcessQty,
			UnfulfillNote:        v.UnfulfillNote,
			Uom:                  v.Uom,
			Status:               int8(v.Status),
		})
	}

	res = &dto.HelperAppGetSalesOrderToCheckDetailResponse{
		SopNumber:           response.SopNumber,
		MerchantName:        response.MerchantName,
		DeliveryDate:        response.DeliveryDate,
		Wrt:                 response.Wrt,
		SopNote:             response.SopNote,
		TotalItemOnProgress: response.TotalItemOnProgress,
		TotalItem:           response.TotalItem,
		TotalKoli:           response.TotalKoli,
		PickerName:          response.PickerName,
		Item:                itemResponse,
		Status:              int8(response.Status),
	}

	return
}

func (s *HelperAppService) CheckerStartChecking(ctx context.Context, req dto.HelperAppCheckerStartCheckingRequest) (res *dto.HelperAppSuccessResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "HelperAppService.CheckerStartChecking")
	defer span.End()

	var response *site_service.SuccessResponse
	if response, err = s.opt.Client.SiteServiceGrpc.CheckerStartChecking(ctx, &site_service.CheckerStartCheckingRequest{
		SopNumber: req.SopNumber,
		CheckerId: req.CheckerId,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("site", "Checker Start Checking")
		return
	}

	res = &dto.HelperAppSuccessResponse{
		Success: response.Success,
	}

	return
}

func (s *HelperAppService) CheckerSubmitChecking(ctx context.Context, req dto.HelperAppCheckerSubmitCheckingRequest) (res *dto.HelperAppSuccessResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "HelperAppService.CheckerSubmitChecking")
	defer span.End()

	var requestArray []*site_service.CheckerSubmitCheckingModel
	for _, v := range req.Request {
		requestArray = append(requestArray, &site_service.CheckerSubmitCheckingModel{
			ItemNumber: v.ItemNumber,
			CheckQty:   v.CheckQuantity,
		})
	}

	var response *site_service.SuccessResponse
	if response, err = s.opt.Client.SiteServiceGrpc.CheckerSubmitChecking(ctx, &site_service.CheckerSubmitCheckingRequest{
		SopNumber: req.SopNumber,
		Request:   requestArray,
		CheckerId: req.CheckerId,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("site", "checker submit checking")
		return
	}

	res = &dto.HelperAppSuccessResponse{
		Success: response.Success,
	}

	return
}

func (s *HelperAppService) CheckerRejectSalesOrder(ctx context.Context, req dto.HelperAppCheckerRejectSalesOrderRequest) (res *dto.HelperAppSuccessResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "HelperAppService.CheckerRejectSalesOrder")
	defer span.End()

	var response *site_service.SuccessResponse
	if response, err = s.opt.Client.SiteServiceGrpc.CheckerRejectSalesOrder(ctx, &site_service.CheckerRejectSalesOrderRequest{
		SopNumber:        req.SopNumber,
		ItemNumberReject: req.ItemNumberReject,
		CheckerId:        req.CheckerId,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("site", "Checker Reject Sales Order")
		return
	}

	res = &dto.HelperAppSuccessResponse{
		Success: response.Success,
	}

	return
}

func (s *HelperAppService) CheckerGetDeliveryKoli(ctx context.Context, sopNumber string) (res *dto.HelperAppCheckerGetDeliveryKoliResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "HelperAppService.CheckerGetDeliveryKoli")
	defer span.End()

	var response *site_service.CheckerGetDeliveryKoliResponse
	if response, err = s.opt.Client.SiteServiceGrpc.CheckerGetDeliveryKoli(ctx, &site_service.CheckerGetDeliveryKoliRequest{
		SopNumber: sopNumber,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("site", "Checker Get Delivery Koli")
		return
	}

	var deliveryKoliResponse []*dto.DeliveryKoliResponse
	for _, v := range response.Data {
		deliveryKoliResponse = append(deliveryKoliResponse, &dto.DeliveryKoliResponse{
			SalesOrderCode: v.SalesOrderCode,
			KoliId:         v.KoliId,
			Name:           v.Name,
			Quantity:       v.Quantity,
		})
	}

	res = &dto.HelperAppCheckerGetDeliveryKoliResponse{
		DeliveryKoli: deliveryKoliResponse,
	}

	return
}

func (s *HelperAppService) CheckerAcceptSalesOrder(ctx context.Context, req dto.HelperAppCheckerAcceptSalesOrderRequest) (res *dto.HelperAppCheckerAcceptSalesOrderResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "HelperAppService.CheckerAcceptSalesOrder")
	defer span.End()

	var requestArray []*site_service.RequestDeliveryKoli
	for _, v := range req.Koli {
		requestArray = append(requestArray, &site_service.RequestDeliveryKoli{
			Id:       v.Id,
			Quantity: v.Quantity,
		})
	}

	var response *site_service.CheckerAcceptSalesOrderResponse
	if response, err = s.opt.Client.SiteServiceGrpc.CheckerAcceptSalesOrder(ctx, &site_service.CheckerAcceptSalesOrderRequest{
		SopNumber: req.SopNumber,
		Koli:      requestArray,
		CheckerId: req.CheckerId,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		r := regexp.MustCompile(`"id":"(.*?)"`)
		match := r.FindStringSubmatch(err.Error())

		if len(match) > 1 {
			err = edenlabs.ErrorValidation("id", match[1]) // TODO: need add logic not only return error from bridge, but need get all error from site service"
		} else {
			err = edenlabs.ErrorRpcNotFound("site", "Checker Accept Sales Order")
		}

		return
	}

	res = &dto.HelperAppCheckerAcceptSalesOrderResponse{
		Success:       response.Success,
		DeliveryOrder: response.DeliveryOrder,
		SalesInvoice:  response.SalesInvoice,
	}

	return
}

func (s *HelperAppService) CheckerHistory(ctx context.Context, req dto.HelperAppCheckerHistoryRequest) (res *dto.HelperAppCheckerHistoryResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "HelperAppService.CheckerHistory")
	defer span.End()

	var response *site_service.CheckerHistoryResponse
	if response, err = s.opt.Client.SiteServiceGrpc.CheckerHistory(ctx, &site_service.CheckerHistoryRequest{
		Offset:    int32(req.Offset),
		Limit:     int32(req.Limit),
		CheckerId: req.CheckerId,
		WrtId:     req.WrtIdGP,
		SopNumber: req.SopNumber,
		Custname:  req.CustomerName,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("site", "Checker History")
		return
	}

	var salesOrderResponse []*dto.HelperAppGetSalesOrderToCheckSOModel
	for _, v := range response.Data {
		salesOrderResponse = append(salesOrderResponse, &dto.HelperAppGetSalesOrderToCheckSOModel{
			SopNumber:           v.SopNumber,
			MerchantName:        v.MerchantName,
			DeliveryDate:        v.DeliveryDate,
			Wrt:                 v.Wrt,
			SopNote:             v.SopNote,
			TotalItemOnProgress: v.TotalItemOnProgress,
			TotalItem:           v.TotalItem,
			TotalKoli:           v.TotalKoli,
			CheckerName:         v.CheckerName,
			PickerName:          v.PickerName,
			Status:              int8(v.Status),
			CountPrintDo:        v.CountPrintDo,
			CountPrintSi:        v.CountPrintSi,
		})
	}

	res = &dto.HelperAppCheckerHistoryResponse{
		SalesOrder: salesOrderResponse,
	}

	return
}

func (s *HelperAppService) CheckerHistoryDetail(ctx context.Context, req dto.HelperAppCheckerHistoryDetailRequest) (res *dto.HelperAppCheckerHistoryDetailResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "HelperAppService.CheckerHistoryDetail")
	defer span.End()

	var response *site_service.CheckerHistoryDetailResponse
	if response, err = s.opt.Client.SiteServiceGrpc.CheckerHistoryDetail(ctx, &site_service.CheckerHistoryDetailRequest{
		CheckerId: req.CheckerId,
		SopNumber: req.SopNumber,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("site", "Checker History Detail")
		return
	}

	var itemResponse []*dto.HelperAppGetSalesOrderPickingDetailItem
	for _, v := range response.Item {
		itemResponse = append(itemResponse, &dto.HelperAppGetSalesOrderPickingDetailItem{
			Id:                   v.Id,
			PickingOrderAssignId: v.PickingOrderAssignId,
			ItemNumber:           v.ItemNumber,
			ItemName:             v.ItemName,
			Picture:              v.Picture,
			OrderQty:             v.OrderQty,
			PickQty:              v.PickQty,
			CheckQty:             v.CheckQty,
			ExcessQty:            v.ExcessQty,
			UnfulfillNote:        v.UnfulfillNote,
			Uom:                  v.Uom,
			Status:               int8(v.Status),
		})
	}

	res = &dto.HelperAppCheckerHistoryDetailResponse{
		SopNumber:           response.SopNumber,
		MerchantName:        response.MerchantName,
		DeliveryDate:        response.DeliveryDate,
		Wrt:                 response.Wrt,
		SopNote:             response.SopNote,
		TotalItemOnProgress: response.TotalItemOnProgress,
		TotalItem:           response.TotalItem,
		TotalKoli:           response.TotalKoli,
		PickerName:          response.PickerName,
		Item:                itemResponse,
		Status:              int8(response.Status),
		CountPrintDo:        response.CountPrintDo,
		CountPrintSi:        response.CountPrintSi,
	}

	return
}
