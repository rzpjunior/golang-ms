package service

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/jwt"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/audit_service"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/catalog_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/notification_service"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/site_service"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-site-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/internal/app/repository"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IPickingOrderService interface {
	// Login
	Login(ctx context.Context, req *pb.LoginHelperRequest) (res *pb.LoginHelperResponse, err error)
	// Picker
	GetGrpc(ctx context.Context, req *pb.GetPickingOrderHeaderRequest) (res *pb.GetPickingOrderHeaderResponse, err error)
	GetDetailGrpc(ctx context.Context, req *pb.GetPickingOrderDetailRequest) (res *pb.GetPickingOrderDetailResponse, err error)
	GetDetailAggregatedProduct(ctx context.Context, req *pb.GetAggregatedProductSalesOrderRequest) (res *pb.GetAggregatedProductSalesOrderResponse, err error)
	StartPickingOrder(ctx context.Context, req *pb.StartPickingOrderRequest) (res *pb.SuccessResponse, err error)
	SubmitPicking(ctx context.Context, req *pb.SubmitPickingRequest) (res *pb.SuccessResponse, err error)
	GetSalesOrderPicking(ctx context.Context, req *pb.GetSalesOrderPickingRequest) (res *pb.GetSalesOrderPickingResponse, err error)
	GetSalesOrderPickingDetail(ctx context.Context, req *pb.GetSalesOrderPickingDetailRequest) (res *pb.GetSalesOrderPickingDetailResponse, err error)
	SubmitSalesOrder(ctx context.Context, req *pb.SubmitSalesOrderRequest) (res *pb.SuccessResponse, err error)
	History(ctx context.Context, req *pb.HistoryRequest) (res *pb.HistoryResponse, err error)
	HistoryDetail(ctx context.Context, req *pb.HistoryDetailRequest) (res *pb.HistoryDetailResponse, err error)
	PickerWidget(ctx context.Context, req *pb.PickerWidgetRequest) (res *pb.PickerWidgetResponse, err error)
	// SPV & Checker
	GetSalesOrderToCheck(ctx context.Context, req *pb.GetSalesOrderToCheckRequest) (res *pb.GetSalesOrderToCheckResponse, err error)
	// SPV
	SPVGetSalesOrderToCheckDetail(ctx context.Context, req *pb.GetSalesOrderToCheckDetailRequest) (res *pb.GetSalesOrderToCheckDetailResponse, err error)
	SPVRejectSalesOrder(ctx context.Context, req *pb.SPVRejectSalesOrderRequest) (res *pb.SuccessResponse, err error)
	SPVAcceptSalesOrder(ctx context.Context, req *pb.SPVAcceptSalesOrderRequest) (res *pb.SuccessResponse, err error)
	SPVWidget(ctx context.Context, req *pb.SPVWidgetRequest) (res *pb.SPVWidgetResponse, err error)
	SPVWrtMonitoring(ctx context.Context, req *pb.GetWrtMonitoringListRequest) (res *pb.GetWrtMonitoringListResponse, err error)
	SPVWrtMonitoringDetail(ctx context.Context, req *pb.GetWrtMonitoringDetailRequest) (res *pb.GetWrtMonitoringDetailResponse, err error)
	// Checker
	CheckerGetSalesOrderToCheckDetail(ctx context.Context, req *pb.GetSalesOrderToCheckDetailRequest) (res *pb.GetSalesOrderToCheckDetailResponse, err error)
	CheckerStartChecking(ctx context.Context, req *pb.CheckerStartCheckingRequest) (res *pb.SuccessResponse, err error)
	CheckerSubmitChecking(ctx context.Context, req *pb.CheckerSubmitCheckingRequest) (res *pb.SuccessResponse, err error)
	CheckerRejectSalesOrder(ctx context.Context, req *pb.CheckerRejectSalesOrderRequest) (res *pb.SuccessResponse, err error)
	CheckerGetDeliveryKoli(ctx context.Context, req *pb.CheckerGetDeliveryKoliRequest) (res *pb.CheckerGetDeliveryKoliResponse, err error)
	CheckerAcceptSalesOrder(ctx context.Context, req *pb.CheckerAcceptSalesOrderRequest) (res *pb.CheckerAcceptSalesOrderResponse, err error)
	CheckerHistory(ctx context.Context, req *pb.CheckerHistoryRequest) (res *pb.CheckerHistoryResponse, err error)
	CheckerHistoryDetail(ctx context.Context, req *pb.CheckerHistoryDetailRequest) (res *pb.CheckerHistoryDetailResponse, err error)
	CheckerWidget(ctx context.Context, req *pb.CheckerWidgetRequest) (res *pb.CheckerWidgetResponse, err error)
}

type PickingOrderService struct {
	opt                          opt.Options
	RepositoryPickingOrderItem   repository.IPickingOrderItemRepository
	RepositoryPickingOrderAssign repository.IPickingOrderAssignRepository
	RepositoryPickingOrder       repository.IPickingOrderRepository
	RepositoryDeliveryKoli       repository.IDeliveryKoliRepository
	RepositoryKoli               repository.IKoliRepository
	RepositoryHelperToken        repository.IHelperTokenRepository
}

func NewPickingOrderService() IPickingOrderService {
	return &PickingOrderService{
		opt:                          global.Setup.Common,
		RepositoryPickingOrderItem:   repository.NewPickingOrderItemRepository(),
		RepositoryPickingOrderAssign: repository.NewPickingOrderAssignRepository(),
		RepositoryPickingOrder:       repository.NewPickingOrderRepository(),
		RepositoryDeliveryKoli:       repository.NewDeliveryKoliRepository(),
		RepositoryKoli:               repository.NewKoliRepository(),
		RepositoryHelperToken:        repository.NewHelperTokenRepository(),
	}
}

func (s *PickingOrderService) Login(ctx context.Context, req *pb.LoginHelperRequest) (res *pb.LoginHelperResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PickingOrderService.LoginHelper")
	defer span.End()

	var login *bridgeService.LoginHelperResponse
	if login, err = s.opt.Client.BridgeServiceGrpc.LoginHelper(ctx, &bridgeService.LoginHelperRequest{
		Email:    req.Email,
		Password: req.Password,
	}); err != nil || login.Code != 200 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("bridge", "Please recheck Email or Password")
		return
	}

	if login.Code == 200 {
		var helper *bridgeService.GetHelperGPResponse
		if helper, err = s.opt.Client.BridgeServiceGrpc.GetHelperGPList(ctx, &bridgeService.GetHelperGPListRequest{
			Limit:  1,
			Userid: req.Email,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "helper")
			return
		}

		// if not helper not found
		if len(helper.Data) == 0 {
			res = &pb.LoginHelperResponse{
				Code:    404,
				Message: "helper not found",
				User:    &pb.LoginHelperResponse_User{},
			}
			return
		}

		var notificationToken *model.HelperToken
		if notificationToken, err = s.RepositoryHelperToken.GetByHelperId(ctx, helper.Data[0].GnlHelperId); err != nil {
			// not exist, and will create
			notificationToken = &model.HelperToken{
				HelperIdGp:        helper.Data[0].GnlHelperId,
				NotificationToken: req.FirebaseToken,
			}
			if err = s.RepositoryHelperToken.Create(ctx, notificationToken); err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorInvalid("notification token")
				return
			}
		} else {
			// exist so will only update
			notificationToken.NotificationToken = req.FirebaseToken
			if err = s.RepositoryHelperToken.Update(ctx, notificationToken, "NotificationToken"); err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorInvalid("notification token")
				return
			}
		}

		jwtInit := jwt.NewJWT([]byte(s.opt.Config.Jwt.Key))
		uc := jwt.UserHelperMobileClaim{
			UserID:   helper.Data[0].GnlHelperId,
			SiteId:   helper.Data[0].Sites[0].Locncode,
			Platform: "helper-mobile",
			Timezone: req.Timezone,
		}

		var jwtGenerate string
		jwtGenerate, err = jwtInit.Create(uc)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("jwt generate")
			return
		}

		var site *bridgeService.GetSiteGPResponse
		if site, err = s.opt.Client.BridgeServiceGrpc.GetSiteGPDetail(ctx, &bridgeService.GetSiteGPDetailRequest{
			Id: helper.Data[0].Sites[0].Locncode,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "site")
			return
		}

		res = &pb.LoginHelperResponse{
			Code:    login.Code,
			Message: login.Message,
			User: &pb.LoginHelperResponse_User{
				Id:       helper.Data[0].GnlHelperId,
				Name:     helper.Data[0].GnlHelperName,
				SiteId:   helper.Data[0].Sites[0].Locncode,
				SiteName: site.Data[0].Locndscr,
				RoleName: helper.Data[0].HelperTypeDesc,
			},
			Token:         jwtGenerate,
			FirebaseToken: notificationToken.NotificationToken,
		}

		return
	}

	res = &pb.LoginHelperResponse{
		Code:    login.Code,
		Message: login.Message,
		User:    &pb.LoginHelperResponse_User{},
	}

	return
}

func (s *PickingOrderService) GetGrpc(ctx context.Context, req *pb.GetPickingOrderHeaderRequest) (res *pb.GetPickingOrderHeaderResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PickingOrderService.GetGrpc")
	defer span.End()

	// TODO UNCOMMENT CODE

	currentTime := time.Now()
	yesterday := currentTime.Add(-time.Hour * 24)
	tomorrow := currentTime.Add(time.Hour * 24)

	/*
		Status GP
		(?):New;3:Picked;4:Finished
	*/
	var status int32
	switch req.WmsPickingStatus {
	case 6:
		status = 1
	case 21, 9:
		status = 3
	case 2:
		status = 4
	default:
		status = 0
	}
	var pickingOrder *bridgeService.GetPickingOrderGPHeaderResponse
	if pickingOrder, err = s.opt.Client.BridgeServiceGrpc.GetPickingOrderGPHeader(ctx, &bridgeService.GetPickingOrderGPHeaderRequest{
		Limit:            req.Limit,
		Offset:           req.Offset,
		Locncode:         req.Locncode,
		Sopnumbe:         req.Sopnumbe,
		Docnumbr:         req.Docnumbr,
		Itemnmbr:         req.Itemnmbr,
		DocdateFrom:      yesterday.Format("2006-01-02"),
		DocdateTo:        tomorrow.Format("2006-01-02"),
		GnlHelperId:      req.GnlHelperId,
		WmsPickingStatus: status,
		Custname:         req.Custname,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "picking order")
		return
	}

	var response []*pb.PickingOrderHeader
	for _, v := range pickingOrder.Data {
		if v.WmsPickingStatus == 4 {
			continue
		}
		// sync picking order data
		po := &model.PickingOrder{
			DocNumber: v.Docnumbr,
			PickerId:  v.WmsPickerId,
			Status:    6,
		}
		if err = s.RepositoryPickingOrder.SyncGP(ctx, po); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("picking order sync")
			return
		}

		response = append(response, &pb.PickingOrderHeader{
			DocNumber:       po.DocNumber,
			PickerId:        v.WmsPickerId,
			DocDate:         v.Docdate,
			Status:          int32(po.Status),
			TotalSalesOrder: int64(v.SalesOrderCount),
			Note:            v.Note,
		})
	}

	res = &pb.GetPickingOrderHeaderResponse{
		Data: response,
	}

	return
}

func (s *PickingOrderService) GetDetailGrpc(ctx context.Context, req *pb.GetPickingOrderDetailRequest) (res *pb.GetPickingOrderDetailResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PickingOrderService.GetDetailGrpc")
	defer span.End()

	var pickingOrderDetail *bridgeService.GetPickingOrderGPDetailResponse
	if pickingOrderDetail, err = s.opt.Client.BridgeServiceGrpc.GetPickingOrderGPDetail(ctx, &bridgeService.GetPickingOrderGPDetailRequest{
		Id: req.Id,
	}); err != nil {
		fmt.Println(err, "pickingOrderDetailpickingOrderDetailpickingOrderDetail----------2", pickingOrderDetail)
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "picking order")
		return
	}
	fmt.Println("pickingOrderDetailpickingOrderDetailpickingOrderDetail----------", pickingOrderDetail)

	var internalPickingOrder *model.PickingOrder
	if internalPickingOrder, err = s.RepositoryPickingOrder.GetDetail(ctx, 0, pickingOrderDetail.Data[0].Docnumbr); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order assign sync")
		return
	}

	// hold aggregated product map
	mapAggregateProduct := map[string]dto.AggregateProduct{}

	for _, v := range pickingOrderDetail.Data[0].Details {
		// TODO : GET SALES ORDER DELIVERY DATE AND WRT THEN SAVE TO DATABASE
		// NOT DONE BECAUSE EXISTING DATA FROM GP DOESNT HAVE THE SALES ORDER DATA CAUSING ERROR
		// get picking order assign's sales order information

		var salesOrder *bridgeService.GetSalesOrderGPListResponse
		if salesOrder, err = s.opt.Client.BridgeServiceGrpc.GetSalesOrderListGPByID(ctx, &bridgeService.GetSalesOrderGPListByIDRequest{
			Id: v.Sopnumbe,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "sales order")
			return
		}

		// get wrt
		var wrt *bridgeService.GetWrtGPResponse
		if wrt, err = s.opt.Client.BridgeServiceGrpc.GetWrtGPList(ctx, &bridgeService.GetWrtGPListRequest{
			Search:    salesOrder.Data[0].Wrt[0].GnL_WRT_ID,
			GnlRegion: salesOrder.Data[0].GnL_Region,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "wrt")
			return
		}

		deliveryDate, _ := time.Parse("2006-01-02", salesOrder.Data[0].ReqShipDate)
		poa := &model.PickingOrderAssign{
			PickingOrderId: internalPickingOrder.Id,
			SopNumber:      v.Sopnumbe,
			SiteID:         v.Locncode,
			// TODO UNCOMMENT DUMMY
			// DeliveryDate: time.Now(), // DUMMY
			// WrtIdGP:      "05-07",    // DUMMY
			DeliveryDate: deliveryDate,
			WrtIdGP:      wrt.Data[0].GnL_WRT_ID,
			Status:       6,
		}
		if err = s.RepositoryPickingOrderAssign.SyncGP(ctx, poa); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("picking order assign sync")
			return
		}

		// if status not finished,need approval,picked,checking
		if poa.Status != 2 && poa.Status != 35 && poa.Status != 16 && poa.Status != 20 {
			// sync picking order item data
			poi := &model.PickingOrderItem{
				IdGp:                 int64(v.Lnitmseq),
				PickingOrderAssignId: poa.ID,
				ItemNumber:           v.Itemnmbr,
				OrderQuantity:        v.IvmQty,
				Status:               6,
			}
			if err = s.RepositoryPickingOrderItem.SyncGP(ctx, poi); err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorInvalid("picking order item sync")
				return
			}

			var productDetail *bridgeService.GetItemGPResponse
			if productDetail, err = s.opt.Client.BridgeServiceGrpc.GetItemGPDetail(ctx, &bridgeService.GetItemGPDetailRequest{
				Id: v.Itemnmbr,
			}); err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("inventory", "product detail")
				return
			}

			// if item name search used and item doesn't contain the word
			if req.ItemName != "" {
				containString := strings.Contains(strings.ToLower(productDetail.Data[0].Itemdesc), strings.ToLower(req.ItemName))
				if !containString {
					continue
				}
			}

			// aggregate products
			var product dto.AggregateProduct = mapAggregateProduct[poi.ItemNumber]
			product = dto.AggregateProduct{
				ItemNumber:      poi.ItemNumber,
				ItemName:        productDetail.Data[0].Itemdesc,
				UomDescription:  v.IvmUofm,
				TotalOrderQty:   mapAggregateProduct[poi.ItemNumber].TotalOrderQty + poi.OrderQuantity,
				TotalPickedQty:  mapAggregateProduct[poi.ItemNumber].TotalPickedQty + poi.PickQuantity,
				TotalSalesOrder: mapAggregateProduct[poi.ItemNumber].TotalSalesOrder + 1,
				Status:          mapAggregateProduct[poi.ItemNumber].Status,
			}

			// initial status
			if product.Status == 0 {
				if poi.Status == 6 || poi.Status == 16 { // new, picked, rejected
					product.Status = poi.Status
				} else if poi.Status == 34 { // unfulfill will be counted as picked
					product.Status = 16
				} else if poi.Status == 9 { // rejected
					product.Status = 9
				}
			} else { // already has status
				if poi.Status == 16 { // status picked
					if product.Status != 9 {
						product.Status = 16
					}
				} else if poi.Status == 9 { // status rejected
					product.Status = 9
				}
			}

			mapAggregateProduct[poi.ItemNumber] = product
		}
	}

	var responseProduct []*pb.AggregatedProduct
	for _, v := range mapAggregateProduct {
		var productImage *catalog_service.GetItemDetailResponse
		fmt.Println("IAMGE nya nihhh>>>>>>>>", v.ItemNumber)
		if productImage, err = s.opt.Client.CatalogServiceGrpc.GetItemDetail(ctx, &catalog_service.GetItemDetailRequest{
			Id: v.ItemNumber,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("inventory", "product image")
			return
		}

		product := &pb.AggregatedProduct{
			ItemNumber:      v.ItemNumber,
			ItemName:        v.ItemName,
			Picture:         "", // filled below if image exist
			UomDescription:  v.UomDescription,
			TotalOrderQty:   v.TotalOrderQty,
			TotalPickedQty:  v.TotalPickedQty,
			TotalSalesOrder: v.TotalSalesOrder,
			Status:          int32(v.Status),
		}

		if len(productImage.Data.ItemImage) > 0 {
			product.Picture = productImage.Data.ItemImage[0].ImageUrl
		}

		responseProduct = append(responseProduct, product)
	}

	res = &pb.GetPickingOrderDetailResponse{
		Data: &pb.PickingOrderDetail{
			DocNumber: internalPickingOrder.DocNumber,
			Status:    int32(internalPickingOrder.Status),
			Product:   responseProduct,
		},
	}

	return
}

func (s *PickingOrderService) GetDetailAggregatedProduct(ctx context.Context, req *pb.GetAggregatedProductSalesOrderRequest) (res *pb.GetAggregatedProductSalesOrderResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PickingOrderService.GetDetailAggregatedProduct")
	defer span.End()

	// get internal id
	var pickingOrder *model.PickingOrder
	if pickingOrder, err = s.RepositoryPickingOrder.GetDetail(ctx, 0, req.Id); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order")
		return
	}

	// get sales order inside the picking order
	var pickingOrderAssign []*model.PickingOrderAssign
	if pickingOrderAssign, _, err = s.RepositoryPickingOrderAssign.Get(ctx, &dto.PickingOrderAssignGetRequest{
		Status:         []int{6, 21, 9}, // new , on progress only , and rejected
		PickingOrderId: []int64{pickingOrder.Id},
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order assign")
		return
	}

	// get product information
	var productDetail *bridgeService.GetItemGPResponse
	if productDetail, err = s.opt.Client.BridgeServiceGrpc.GetItemGPDetail(ctx, &bridgeService.GetItemGPDetailRequest{
		Id: req.ItemNumber,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("inventory", "product detail")
		return
	}

	// get product image
	var productImage *catalog_service.GetItemDetailResponse
	if productImage, err = s.opt.Client.CatalogServiceGrpc.GetItemDetail(ctx, &catalog_service.GetItemDetailRequest{
		Id: req.ItemNumber,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("inventory", "product image")
		return
	}

	// preprare product response
	var responseSalesOrder []*pb.SalesOrderInformation
	for _, v := range pickingOrderAssign {
		// get picking order item
		var pickingOrderItem []*model.PickingOrderItem
		if pickingOrderItem, _, err = s.RepositoryPickingOrderItem.Get(ctx, &dto.PickingOrderItemGetRequest{
			PickingOrderAssignId: []int64{v.ID},
			ItemNumber:           []string{req.ItemNumber},
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("picking order item")
			return
		}

		// get picking order assign's sales order information
		var salesOrder *bridgeService.GetSalesOrderGPListResponse
		if salesOrder, err = s.opt.Client.BridgeServiceGrpc.GetSalesOrderListGPByID(ctx, &bridgeService.GetSalesOrderGPListByIDRequest{
			Id: v.SopNumber,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "sales order")
			return
		}

		// get wrt
		var wrt *bridgeService.GetWrtGPResponse
		if wrt, err = s.opt.Client.BridgeServiceGrpc.GetWrtGPList(ctx, &bridgeService.GetWrtGPListRequest{
			Search:    salesOrder.Data[0].Wrt[0].GnL_WRT_ID,
			GnlRegion: salesOrder.Data[0].GnL_Region,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "wrt")
			return
		}

		// TODO DELETE CODE
		for _, v2 := range pickingOrderItem {
			responseSalesOrder = append(responseSalesOrder, &pb.SalesOrderInformation{
				Id:           v2.Id,
				SopNumber:    v.SopNumber,
				MerchantName: salesOrder.Data[0].Customer[0].Custname,
				// MerchantName: "DUMMY MERCHANT NAME",
				Wrt: wrt.Data[0].Strttime + "-" + wrt.Data[0].Endtime,
				// Wrt:           "DUMMY WRT",
				OrderQty:      v2.OrderQuantity,
				PickedQty:     v2.PickQuantity,
				UnfulfillNote: v2.UnfulfillNote,
				Status:        int32(v2.Status),
			})
		}
	}

	res = &pb.GetAggregatedProductSalesOrderResponse{
		Data: &pb.AggregatedProductSalesOrder{
			ItemNumber:     req.ItemNumber,
			ItemName:       productDetail.Data[0].Itemdesc,
			UomDescription: productDetail.Data[0].Uomschdl,
			Picture:        "", // filled below if image exist
			SalesOrder:     responseSalesOrder,
		},
	}

	if len(productImage.Data.ItemImage) > 0 {
		res.Data.Picture = productImage.Data.ItemImage[0].ImageUrl
	}

	return
}

func (s *PickingOrderService) StartPickingOrder(ctx context.Context, req *pb.StartPickingOrderRequest) (res *pb.SuccessResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PickingOrderService.StartPickingOrder")
	defer span.End()

	// get internal id
	var pickingOrder *model.PickingOrder
	if pickingOrder, err = s.RepositoryPickingOrder.GetDetail(ctx, 0, req.DocNumber); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order")
		return
	}

	// if status not new
	if pickingOrder.Status != 6 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorMustStatus("picking order", "new")
		return
	}

	// get sales order inside the picking order
	var pickingOrderAssign []*model.PickingOrderAssign
	if pickingOrderAssign, _, err = s.RepositoryPickingOrderAssign.Get(ctx, &dto.PickingOrderAssignGetRequest{
		PickingOrderId: []int64{pickingOrder.Id},
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order assign")
		return
	}

	// submit picking
	var pickingDetailRequest []*bridgeService.PickingDetails
	for _, v := range pickingOrderAssign {
		var pickingOrderItems []*model.PickingOrderItem
		if pickingOrderItems, _, err = s.RepositoryPickingOrderItem.Get(ctx, &dto.PickingOrderItemGetRequest{
			PickingOrderAssignId: []int64{v.ID},
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("picking order assign")
			return
		}

		for _, v2 := range pickingOrderItems {
			pickingDetailRequest = append(pickingDetailRequest, &bridgeService.PickingDetails{
				Sopnumbe:     v.SopNumber,
				Lnitmseq:     int32(v2.IdGp),
				IvmQtyPickso: v2.OrderQuantity,
			})
		}
	}

	pickingOrder.Status = 21
	pickingOrder.StartTime = time.Now()
	if _, err = s.opt.Client.BridgeServiceGrpc.SubmitPickingCheckingPickingOrder(ctx, &bridgeService.SubmitPickingCheckingRequest{
		Picking: &bridgeService.Picking{
			Docnumbr: pickingOrder.DocNumber,
			Strttime: pickingOrder.StartTime.Format("15:04:05"),
			Endtime:  "23:59:59",
			Details:  pickingDetailRequest,
		},
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "submit picking")
		return
	}
	// change status PO from new to on progress
	if err = s.RepositoryPickingOrder.Update(ctx, pickingOrder, "Status", "StartTime"); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order")
		return
	}

	for _, v := range pickingOrderAssign {
		v.Status = 21
		if err = s.RepositoryPickingOrderAssign.Update(ctx, v, "Status"); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("picking order assign")
			return
		}
	}

	res = &pb.SuccessResponse{
		Success: true,
	}

	return
}

func (s *PickingOrderService) SubmitPicking(ctx context.Context, req *pb.SubmitPickingRequest) (res *pb.SuccessResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PickingOrderService.SubmitPicking")
	defer span.End()

	for _, v := range req.Request {
		var pickingOrderItem *model.PickingOrderItem
		if pickingOrderItem, err = s.RepositoryPickingOrderItem.GetByID(ctx, &dto.PickingOrderItemGetDetailRequest{
			Id: v.Id,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("picking order item")
			return
		}

		// calculate excess qty
		var excessQuantity float64
		if v.PickQty > pickingOrderItem.OrderQuantity {
			excessQuantity = v.PickQty - pickingOrderItem.OrderQuantity
		}

		// update picking order assign to in progress
		var pickingOrderAssign *model.PickingOrderAssign
		if pickingOrderAssign, err = s.RepositoryPickingOrderAssign.GetByID(ctx, pickingOrderItem.PickingOrderAssignId, ""); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("picking order assign")
			return
		}

		// if status not new / on progress / rejected then refuse write
		if pickingOrderAssign.Status != 6 && pickingOrderAssign.Status != 21 && pickingOrderAssign.Status != 9 {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorMustStatus("picking order assign", "new / on progress / rejected")
			return
		}

		// update picking order to in progress
		var pickingOrder *model.PickingOrder
		if pickingOrder, err = s.RepositoryPickingOrder.GetDetail(ctx, pickingOrderAssign.PickingOrderId, ""); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("picking order")
			return
		}

		// change status picking order assign to on progress again
		pickingOrderAssign.Status = 21
		if err = s.RepositoryPickingOrderAssign.Update(ctx, pickingOrderAssign, "Status"); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("picking order assign")
			return
		}

		// change status picking order to on progress again
		pickingOrder.Status = 21
		if err = s.RepositoryPickingOrder.Update(ctx, pickingOrder, "Status"); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("picking order")
			return
		}

		// check note and write
		if v.PickQty < pickingOrderItem.OrderQuantity { // if unfulfill then need note
			if v.UnfulfillNote == "" {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRequired("unfulfill note")
				return
			}

			pickingOrderItem = &model.PickingOrderItem{
				Id:             v.Id,
				PickQuantity:   v.PickQty,
				UnfulfillNote:  v.UnfulfillNote,
				ExcessQuantity: excessQuantity,
				Status:         34,
			}

			if err = s.RepositoryPickingOrderItem.Update(ctx, pickingOrderItem, "PickQuantity", "UnfulfillNote", "ExcessQuantity", "Status"); err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorInvalid("picking order item")
				return
			}
		} else {
			pickingOrderItem = &model.PickingOrderItem{
				Id:             v.Id,
				PickQuantity:   v.PickQty,
				UnfulfillNote:  "",
				ExcessQuantity: excessQuantity,
				Status:         16,
			}

			if err = s.RepositoryPickingOrderItem.Update(ctx, pickingOrderItem, "PickQuantity", "UnfulfillNote", "ExcessQuantity", "Status"); err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorInvalid("picking order item")
				return
			}
		}
	}

	res = &pb.SuccessResponse{
		Success: true,
	}

	return
}

func (s *PickingOrderService) GetSalesOrderPicking(ctx context.Context, req *pb.GetSalesOrderPickingRequest) (res *pb.GetSalesOrderPickingResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PickingOrderService.GetSalesOrderPicking")
	defer span.End()

	var response []*pb.SalesOrderPicking

	var pickingOrder *model.PickingOrder
	if pickingOrder, err = s.RepositoryPickingOrder.GetDetail(ctx, 0, req.DocNumber); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order assign")
		return
	}

	// get sales order inside the picking order
	var pickingOrderAssign []*model.PickingOrderAssign
	if pickingOrderAssign, _, err = s.RepositoryPickingOrderAssign.Get(ctx, &dto.PickingOrderAssignGetRequest{
		PickingOrderId: []int64{pickingOrder.Id},
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order assign")
		return
	}

	for _, v := range pickingOrderAssign {
		var (
			isKoliProcessable bool = true
			unfulfillExist    bool
			totalKoli         float64
			pickingOrderItem  []*model.PickingOrderItem
			deliveryKoli      []*model.DeliveryKoli
		)

		// get items
		if pickingOrderItem, _, err = s.RepositoryPickingOrderItem.Get(ctx, &dto.PickingOrderItemGetRequest{
			PickingOrderAssignId: []int64{v.ID},
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("picking order item")
			return
		}

		// check if all item has been picked
		for _, v2 := range pickingOrderItem {
			// if contain new / rejected item, then user cannot fill koli
			if v2.Status == 6 || v2.Status == 9 {
				isKoliProcessable = false
			}

			// if status unfulfill
			if v2.Status == 34 {
				unfulfillExist = true
			}
		}

		// count koli
		if deliveryKoli, _, err = s.RepositoryDeliveryKoli.Get(ctx, &dto.DeliveryKoliGetRequest{
			SopNumber: v.SopNumber,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("delivery koli")
			return
		}
		for _, v2 := range deliveryKoli {
			totalKoli += v2.Quantity
		}

		// get picking order assign's sales order information
		var salesOrder *bridgeService.GetSalesOrderGPListResponse
		if salesOrder, err = s.opt.Client.BridgeServiceGrpc.GetSalesOrderListGPByID(ctx, &bridgeService.GetSalesOrderGPListByIDRequest{
			Id: v.SopNumber,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "sales order")
			return
		}

		// TODO DELETE CODE
		response = append(response, &pb.SalesOrderPicking{
			SopNumber:    v.SopNumber,
			MerchantName: salesOrder.Data[0].Customer[0].Custname,
			// MerchantName: "DUMMY MERCHANT",
			SopNote: salesOrder.Data[0].Commntid,
			// SopNote:          "DUMMY SOP NOTE",
			TotalKoli:        totalKoli,
			Status:           int32(v.Status),
			ReadyToPack:      isKoliProcessable,
			ContainUnfulfill: unfulfillExist,
		})
	}

	res = &pb.GetSalesOrderPickingResponse{
		Data: response,
	}

	return
}

func (s *PickingOrderService) GetSalesOrderPickingDetail(ctx context.Context, req *pb.GetSalesOrderPickingDetailRequest) (res *pb.GetSalesOrderPickingDetailResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PickingOrderService.GetSalesOrderPickingDetail")
	defer span.End()

	// get sales order inside the picking order
	var pickingOrderAssign *model.PickingOrderAssign
	if pickingOrderAssign, err = s.RepositoryPickingOrderAssign.GetByID(ctx, 0, req.SopNumber); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order assign")
		return
	}

	// this handler was used for showed for check Checker information when New, OnProgress(not finished yet)
	var helperName, helperID string
	if pickingOrderAssign.CheckerIdGp != "" {
		var helperDetail *bridgeService.GetHelperGPResponse
		if helperDetail, err = s.opt.Client.BridgeServiceGrpc.GetHelperGPDetail(ctx, &bridgeService.GetHelperGPDetailRequest{
			Id: pickingOrderAssign.CheckerIdGp,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("helper detail")
			return
		}

		helperName = helperDetail.Data[0].GnlHelperName
		helperID = helperDetail.Data[0].GnlHelperId
	}

	// get all item of the SOP
	var pickingOrderItem []*model.PickingOrderItem
	if pickingOrderItem, _, err = s.RepositoryPickingOrderItem.Get(ctx, &dto.PickingOrderItemGetRequest{
		PickingOrderAssignId: []int64{pickingOrderAssign.ID},
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order item")
		return
	}

	// count delivery koli
	var (
		deliveryKoli []*model.DeliveryKoli
		totalKoli    float64
	)
	if deliveryKoli, _, err = s.RepositoryDeliveryKoli.Get(ctx, &dto.DeliveryKoliGetRequest{
		SopNumber: req.SopNumber,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("delivery koli")
		return
	}
	for _, v2 := range deliveryKoli {
		totalKoli += v2.Quantity
	}

	// get picking order assign's sales order information
	var salesOrder *bridgeService.GetSalesOrderGPListResponse
	if salesOrder, err = s.opt.Client.BridgeServiceGrpc.GetSalesOrderListGPByID(ctx, &bridgeService.GetSalesOrderGPListByIDRequest{
		Id: pickingOrderAssign.SopNumber,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "sales order")
		return
	}

	// get wrt
	var wrt *bridgeService.GetWrtGPResponse
	if wrt, err = s.opt.Client.BridgeServiceGrpc.GetWrtGPList(ctx, &bridgeService.GetWrtGPListRequest{
		Search:    pickingOrderAssign.WrtIdGP,
		GnlRegion: salesOrder.Data[0].GnL_Region,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "wrt")
		return
	}

	// prepare item response
	var (
		itemResponse        []*pb.PickingOrderItem
		totalItemOnProgress int64
	)
	for _, v := range pickingOrderItem {
		// get product detail
		var productDetail *bridgeService.GetItemGPResponse
		if productDetail, err = s.opt.Client.BridgeServiceGrpc.GetItemGPDetail(ctx, &bridgeService.GetItemGPDetailRequest{
			Id: v.ItemNumber,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("inventory", "product detail")
			return
		}

		// get product image
		var productImage *catalog_service.GetItemDetailResponse
		if productImage, err = s.opt.Client.CatalogServiceGrpc.GetItemDetail(ctx, &catalog_service.GetItemDetailRequest{
			Id: v.ItemNumber,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("inventory", "product image")
			return
		}

		product := &pb.PickingOrderItem{
			Id:                   v.Id,
			PickingOrderAssignId: v.PickingOrderAssignId,
			ItemNumber:           v.ItemNumber,
			ItemName:             productDetail.Data[0].Itemdesc,
			Picture:              "", // filled below if picture exist
			OrderQty:             v.OrderQuantity,
			PickQty:              v.PickQuantity,
			CheckQty:             v.CheckQuantity,
			ExcessQty:            v.ExcessQuantity,
			UnfulfillNote:        v.UnfulfillNote,
			Uom:                  productDetail.Data[0].Uomschdl,
			Status:               int32(v.Status),
		}

		if len(productImage.Data.ItemImage) > 0 {
			product.Picture = productImage.Data.ItemImage[0].ImageUrl
		}

		itemResponse = append(itemResponse, product)
		if product.Status != 16 {
			totalItemOnProgress++
		}
	}

	// TODO DELETE CODE
	res = &pb.GetSalesOrderPickingDetailResponse{
		SopNumber:    pickingOrderAssign.SopNumber,
		MerchantName: salesOrder.Data[0].Customer[0].Custname,
		// MerchantName:        "DUMMY MERCHANT",
		Wrt:                 wrt.Data[0].Strttime + "-" + wrt.Data[0].Endtime,
		DeliveryDate:        salesOrder.Data[0].ReqShipDate,
		TotalKoli:           totalKoli,
		TotalItemOnProgress: totalItemOnProgress,
		TotalItem:           int64(len(itemResponse)),
		SopNote:             salesOrder.Data[0].Commntid,
		Item:                itemResponse,
		Status:              int32(pickingOrderAssign.Status),
		HelperName:          helperName,
		HelperId:            helperID,
	}

	return
}

func (s *PickingOrderService) SubmitSalesOrder(ctx context.Context, req *pb.SubmitSalesOrderRequest) (res *pb.SuccessResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PickingOrderService.SubmitSalesOrder")
	defer span.End()

	var (
		sendToSpv bool
		haveKoli  bool
	)

	// get sales order inside the picking order
	var pickingOrderAssign *model.PickingOrderAssign
	if pickingOrderAssign, err = s.RepositoryPickingOrderAssign.GetByID(ctx, 0, req.SopNumber); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order assign")
		return
	}

	// if stataus not on progress
	if pickingOrderAssign.Status != 21 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorMustStatus("picking order assign", "on progress")
		return
	}

	var pickingOrderItem []*model.PickingOrderItem
	if pickingOrderItem, _, err = s.RepositoryPickingOrderItem.Get(ctx, &dto.PickingOrderItemGetRequest{
		PickingOrderAssignId: []int64{pickingOrderAssign.ID},
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order item")
		return
	}

	for _, v := range pickingOrderItem {
		if v.Status == 6 || v.Status == 9 {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorValidation("picking status", "one or more item have not been picked")
			return
		}
		if v.Status == 34 {
			sendToSpv = true
		}
	}

	// validate koli
	for _, v := range req.Request {
		if _, err = s.RepositoryKoli.GetByID(ctx, v.Id); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("koli")
			return
		}

		if v.Quantity > 0 {
			haveKoli = true
		}
	}

	// if all koli quantity is zero
	if !haveKoli {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRequired("koli")
		return
	}

	// write delivery koli
	for _, v := range req.Request {
		if v.Quantity != 0 {
			if err = s.RepositoryDeliveryKoli.Create(ctx, &model.DeliveryKoli{
				SopNumber: pickingOrderAssign.SopNumber,
				KoliId:    v.Id,
				Quantity:  v.Quantity,
			}); err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorInvalid("delivery koli")
				return
			}
		}
	}
	// update picking order assign status
	if sendToSpv {
		pickingOrderAssign.Status = 35
		if err = s.RepositoryPickingOrderAssign.Update(ctx, pickingOrderAssign, "Status"); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("picking order assign")
			return
		}
	} else {
		pickingOrderAssign.Status = 16
		if err = s.RepositoryPickingOrderAssign.Update(ctx, pickingOrderAssign, "Status"); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("picking order assign")
			return
		}
	}

	// audit log
	if _, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &audit_service.CreateLogRequest{
		Log: &audit_service.Log{
			UserId:      0,
			ReferenceId: pickingOrderAssign.SopNumber,
			Type:        "picking order assign",
			Function:    "submit sales order",
			CreatedAt:   timestamppb.New(time.Now()),
			Note:        "Picker with id of " + req.PickerId + " finished " + pickingOrderAssign.SopNumber,
		},
	}); err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		span.RecordError(err)
		err = edenlabs.ErrorRpcNotFound("audit", "audit_log")
		return
	}

	// update picking order assign status time finished if no poa anymore
	var total int64
	if _, total, err = s.RepositoryPickingOrderAssign.Get(ctx, &dto.PickingOrderAssignGetRequest{
		Status:         []int{6, 21, 9},
		PickingOrderId: []int64{pickingOrderAssign.PickingOrderId},
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order assign")
		return
	}

	if total == 0 {
		var pickingOrder *model.PickingOrder
		if pickingOrder, err = s.RepositoryPickingOrder.GetDetail(ctx, pickingOrderAssign.PickingOrderId, ""); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("picking order")
			return
		}

		pickingOrder.EndTime = time.Now()
		if err = s.RepositoryPickingOrder.Update(ctx, pickingOrder, "EndTime"); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("picking order")
			return
		}
	}

	res = &pb.SuccessResponse{
		Success: true,
	}

	return
}

func (s *PickingOrderService) History(ctx context.Context, req *pb.HistoryRequest) (res *pb.HistoryResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PickingOrderService.History")
	defer span.End()

	// TODO REMOVE DUMMY
	// dummyRES := &pb.HistoryResponse{
	// 	Data: []*pb.SalesOrderToCheck{
	// 		{
	// 			SopNumber:           "DUMMY SOP NUMBER",
	// 			MerchantName:        "DUMMY MERCHANT NAME",
	// 			DeliveryDate:        "DUMMY DELIVERY DATE",
	// 			Wrt:                 "DUMMY WRT",
	// 			SopNote:             "DUMMY SOP NOTE",
	// 			TotalItemOnProgress: 1,
	// 			TotalItem:           1,
	// 			TotalKoli:           1,
	// 			CheckerName:         "DUMMY CHECKER NAME",
	// 			PickerName:          "DUMMY PICKER NAME",
	// 			Status:              16,
	// 		},
	// 	},
	// }
	// return dummyRES, nil

	var response []*pb.SalesOrderToCheck

	currentTime := time.Now()
	yesterday := currentTime.Add(-time.Hour * 48)
	tomorrow := currentTime.Add(time.Hour * 24)

	var pickingOrder *bridgeService.GetPickingOrderGPHeaderResponse
	if pickingOrder, err = s.opt.Client.BridgeServiceGrpc.GetPickingOrderGPHeader(ctx, &bridgeService.GetPickingOrderGPHeaderRequest{
		Limit:  req.Limit,
		Offset: req.Offset,
		// Sopnumbe:         req.SopNumber, // GP cannot search like for so number
		DocdateFrom:      yesterday.Format("2006-01-02"),
		DocdateTo:        tomorrow.Format("2006-01-02"),
		GnlHelperId:      req.PickerId,
		Custname:         req.Custname,
		WmsPickingStatus: 4,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "picking order")
		return
	}

	for _, v := range pickingOrder.Data {
		var pickingOrder *model.PickingOrder
		if pickingOrder, err = s.RepositoryPickingOrder.GetDetail(ctx, 0, v.Docnumbr); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("picking order assign")
			return
		}

		// get sales order inside the picking order
		var pickingOrderAssign []*model.PickingOrderAssign
		if pickingOrderAssign, _, err = s.RepositoryPickingOrderAssign.Get(ctx, &dto.PickingOrderAssignGetRequest{
			PickingOrderId: []int64{pickingOrder.Id},
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("picking order assign")
			return
		}
		for _, v := range pickingOrderAssign {
			var (
				totalKoli    float64
				deliveryKoli []*model.DeliveryKoli
			)
			if !strings.Contains(v.SopNumber, req.SopNumber) {
				continue
			}

			// count koli
			if deliveryKoli, _, err = s.RepositoryDeliveryKoli.Get(ctx, &dto.DeliveryKoliGetRequest{
				SopNumber: v.SopNumber,
			}); err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorInvalid("delivery koli")
				return
			}
			for _, v2 := range deliveryKoli {
				totalKoli += v2.Quantity
			}

			// get picking order assign's sales order information
			var salesOrder *bridgeService.GetSalesOrderGPListResponse
			if salesOrder, err = s.opt.Client.BridgeServiceGrpc.GetSalesOrderListGPByID(ctx, &bridgeService.GetSalesOrderGPListByIDRequest{
				Id: v.SopNumber,
			}); err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "sales order")
				return
			}

			// get picker name
			var pickingOrder *model.PickingOrder
			if pickingOrder, err = s.RepositoryPickingOrder.GetDetail(ctx, v.PickingOrderId, ""); err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorInvalid("picking order")
				return
			}

			var picker *bridgeService.GetHelperGPResponse
			if picker, err = s.opt.Client.BridgeServiceGrpc.GetHelperGPDetail(ctx, &bridgeService.GetHelperGPDetailRequest{
				Id: pickingOrder.PickerId,
			}); err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "helper")
				return
			}

			var checker *bridgeService.GetHelperGPResponse
			if checker, err = s.opt.Client.BridgeServiceGrpc.GetHelperGPDetail(ctx, &bridgeService.GetHelperGPDetailRequest{
				Id: v.CheckerIdGp,
			}); err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "helper")
				return
			}

			// get DO & SI
			var do *bridgeService.GetDeliveryOrderGPListResponse
			var si *bridgeService.GetSalesInvoiceGPListResponse
			var count_print_si, count_print_do int32

			if do, err = s.opt.Client.BridgeServiceGrpc.GetDeliveryOrderListGP(ctx, &bridgeService.GetDeliveryOrderGPListRequest{
				SopNumbe: v.SopNumber,
			}); err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "delivery order")
				return
			}

			if si, err = s.opt.Client.BridgeServiceGrpc.GetSalesInvoiceGPList(ctx, &bridgeService.GetSalesInvoiceGPListRequest{
				SoNumber: v.SopNumber,
			}); err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "sales invoice")
				return
			}

			if len(do.Data) > 0 {
				count_print_do = do.Data[0].DataAttachment.PrintCount
			}

			if len(si.Data) > 0 {
				count_print_si = si.Data[0].DataAttachment.PrintCount
			}
			// TODO fill proper data
			response = append(response, &pb.SalesOrderToCheck{
				SopNumber:           v.SopNumber,
				MerchantName:        salesOrder.Data[0].Customer[0].Custname,
				DeliveryDate:        salesOrder.Data[0].ReqShipDate,
				Wrt:                 "",
				SopNote:             salesOrder.Data[0].Commntid,
				TotalItemOnProgress: 0,
				TotalItem:           0,
				TotalKoli:           totalKoli,
				CheckerName:         checker.Data[0].GnlHelperName,
				PickerName:          picker.Data[0].GnlHelperName,
				Status:              int32(v.Status),
				CountPrintDo:        count_print_do,
				CountPrintSi:        count_print_si,
			})
		}
	}

	res = &pb.HistoryResponse{
		Data: response,
	}

	return
}

func (s *PickingOrderService) HistoryDetail(ctx context.Context, req *pb.HistoryDetailRequest) (res *pb.HistoryDetailResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PickingOrderService.HistoryDetail")
	defer span.End()

	// // TODO REMOVE DUMMY
	// dummyRES := &pb.HistoryDetailResponse{
	// 	SopNumber:           "DUMMY SOP NUMBER",
	// 	MerchantName:        "DUMMY MERCHANT",
	// 	Wrt:                 "DUMMY WRT",
	// 	DeliveryDate:        "DUMMY DELIVERY DATE",
	// 	TotalKoli:           1,
	// 	TotalItemOnProgress: 1,
	// 	TotalItem:           1,
	// 	SopNote:             "DUMMY SOP NOTE",
	// 	Item: []*pb.PickingOrderItem{
	// 		{
	// 			Id:                   1,
	// 			PickingOrderAssignId: 1,
	// 			ItemNumber:           "DUMMY ITEM NUMBER",
	// 			ItemName:             "DUMMY ITEM NAME",
	// 			Picture:              "https://www.google.com/url?sa=i&url=https%3A%2F%2Fwww.marca.com%2Fen%2Fwwe%2F2023%2F02%2F05%2F63dee56cca47419b198b456d.html&psig=AOvVaw1PriZ2jtFpIWL8zpJB_EUt&ust=1683603710713000&source=images&cd=vfe&ved=0CBEQjRxqFwoTCKj-xNrm5P4CFQAAAAAdAAAAABAI",
	// 			OrderQty:             1,
	// 			PickQty:              1,
	// 			CheckQty:             1,
	// 			ExcessQty:            0,
	// 			UnfulfillNote:        "",
	// 			Uom:                  "KG",
	// 			Status:               16,
	// 		},
	// 	},
	// 	Status: 2,
	// }

	// return dummyRES, nil

	// get sales order inside the picking order
	var pickingOrderAssign *model.PickingOrderAssign
	if pickingOrderAssign, err = s.RepositoryPickingOrderAssign.GetByID(ctx, 0, req.SopNumber); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order assign")
		return
	}

	// get all item of the SOP
	var pickingOrderItem []*model.PickingOrderItem
	if pickingOrderItem, _, err = s.RepositoryPickingOrderItem.Get(ctx, &dto.PickingOrderItemGetRequest{
		PickingOrderAssignId: []int64{pickingOrderAssign.ID},
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order item")
		return
	}

	// count delivery koli
	var (
		deliveryKoli []*model.DeliveryKoli
		totalKoli    float64
	)
	if deliveryKoli, _, err = s.RepositoryDeliveryKoli.Get(ctx, &dto.DeliveryKoliGetRequest{
		SopNumber: req.SopNumber,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("delivery koli")
		return
	}
	for _, v2 := range deliveryKoli {
		totalKoli += v2.Quantity
	}

	// get picking order assign's sales order information
	var salesOrder *bridgeService.GetSalesOrderGPListResponse
	if salesOrder, err = s.opt.Client.BridgeServiceGrpc.GetSalesOrderListGPByID(ctx, &bridgeService.GetSalesOrderGPListByIDRequest{
		Id: pickingOrderAssign.SopNumber,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "sales order")
		return
	}

	// get wrt
	var wrt *bridgeService.GetWrtGPResponse
	if wrt, err = s.opt.Client.BridgeServiceGrpc.GetWrtGPList(ctx, &bridgeService.GetWrtGPListRequest{
		Search:    pickingOrderAssign.WrtIdGP,
		GnlRegion: salesOrder.Data[0].GnL_Region,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "wrt")
		return
	}

	// prepare item response
	var (
		itemResponse        []*pb.PickingOrderItem
		totalItemOnProgress int64
	)
	for _, v := range pickingOrderItem {
		// get product detail
		var productDetail *bridgeService.GetItemGPResponse
		if productDetail, err = s.opt.Client.BridgeServiceGrpc.GetItemGPDetail(ctx, &bridgeService.GetItemGPDetailRequest{
			Id: v.ItemNumber,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("inventory", "product detail")
			return
		}

		// get product image
		var productImage *catalog_service.GetItemDetailResponse
		if productImage, err = s.opt.Client.CatalogServiceGrpc.GetItemDetail(ctx, &catalog_service.GetItemDetailRequest{
			Id: v.ItemNumber,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("inventory", "product image")
			return
		}

		product := &pb.PickingOrderItem{
			Id:                   v.Id,
			PickingOrderAssignId: v.PickingOrderAssignId,
			ItemNumber:           v.ItemNumber,
			ItemName:             productDetail.Data[0].Itemdesc,
			Picture:              "", // filled below if picture exist
			OrderQty:             v.OrderQuantity,
			PickQty:              v.PickQuantity,
			CheckQty:             v.CheckQuantity,
			ExcessQty:            v.ExcessQuantity,
			UnfulfillNote:        v.UnfulfillNote,
			Uom:                  productDetail.Data[0].Uomschdl,
			Status:               int32(v.Status),
		}

		if len(productImage.Data.ItemImage) > 0 {
			product.Picture = productImage.Data.ItemImage[0].ImageUrl
		}

		itemResponse = append(itemResponse, product)
		if product.Status != 16 {
			totalItemOnProgress++
		}
	}

	res = &pb.HistoryDetailResponse{
		SopNumber:           pickingOrderAssign.SopNumber,
		MerchantName:        salesOrder.Data[0].Customer[0].Custname,
		Wrt:                 wrt.Data[0].Strttime + "-" + wrt.Data[0].Endtime,
		DeliveryDate:        salesOrder.Data[0].ReqShipDate,
		TotalKoli:           totalKoli,
		TotalItemOnProgress: totalItemOnProgress,
		TotalItem:           int64(len(itemResponse)),
		SopNote:             salesOrder.Data[0].Commntid,
		Item:                itemResponse,
		Status:              int32(pickingOrderAssign.Status),
	}

	return
}

func (s *PickingOrderService) PickerWidget(ctx context.Context, req *pb.PickerWidgetRequest) (res *pb.PickerWidgetResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PickingOrderService.PickerWidget")
	defer span.End()

	var pickingOrder *bridgeService.GetPickingOrderGPHeaderResponse

	currentTime := time.Now()
	yesterday := currentTime.Add(-time.Hour * 24)
	tomorrow := currentTime.Add(time.Hour * 24)

	if pickingOrder, err = s.opt.Client.BridgeServiceGrpc.GetPickingOrderGPHeader(ctx, &bridgeService.GetPickingOrderGPHeaderRequest{
		Limit:       1000,
		DocdateFrom: yesterday.Format("2006-01-02"),
		DocdateTo:   tomorrow.Format("2006-01-02"),
		GnlHelperId: req.GnlHelperId,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "picking order")
		return
	}

	var pickingOrderFilter []int64
	for _, v := range pickingOrder.Data {
		// get internal picking order id
		var pickingOrder *model.PickingOrder
		if pickingOrder, err = s.RepositoryPickingOrder.GetDetail(ctx, 0, v.Docnumbr); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("picking order")
			return
		}

		pickingOrderFilter = append(pickingOrderFilter, pickingOrder.Id)
	}

	// get sales order inside the picking order
	var pickingOrderAssign []*model.PickingOrderAssign
	if pickingOrderAssign, _, err = s.RepositoryPickingOrderAssign.Get(ctx, &dto.PickingOrderAssignGetRequest{
		PickingOrderId: pickingOrderFilter,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order assign")
		return
	}

	res = &pb.PickerWidgetResponse{}

	for _, v := range pickingOrderAssign {
		if v.Status == 6 { // new
			res.TotalNew++
		}

		if v.Status == 2 { // finished
			res.TotalPicked++
		}

		if v.Status == 21 || v.Status == 9 { // on progress & rejected
			res.TotalOnProgress++
		}

		if v.Status == 35 { // need approval
			res.TotalNeedApproval++
		}

		res.TotalSalesOrder++
	}

	if res.TotalSalesOrder != 0 {
		res.TotalOnProgressPercentage = math.Round(float64(res.TotalOnProgress)/float64(res.TotalSalesOrder)*100) / 100
		res.TotalNeedApprovalPercentage = math.Round(float64(res.TotalNeedApproval)/float64(res.TotalSalesOrder)*100) / 100
		res.TotalPickedPercentage = math.Round(float64(res.TotalPicked)/float64(res.TotalSalesOrder)*100) / 100
	}

	return
}

func (s *PickingOrderService) GetSalesOrderToCheck(ctx context.Context, req *pb.GetSalesOrderToCheckRequest) (res *pb.GetSalesOrderToCheckResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PickingOrderService.GetSalesOrderToCheck")
	defer span.End()

	var (
		statusInt []int
		sopNumber []string
		reqRepo   *dto.PickingOrderAssignGetRequest
	)

	for _, v := range req.Status {
		statusInt = append(statusInt, int(v))
	}
	if req.SopNumber != "" {
		sopNumber = append(sopNumber, req.SopNumber)
	}

	// get sales order inside the picking order
	reqRepo = &dto.PickingOrderAssignGetRequest{
		Offset:    int(req.Offset),
		Limit:     int(req.Limit),
		Status:    statusInt,
		SopNumber: sopNumber,
		SiteID:    []string{req.SiteId},
	}
	if len(req.WrtIds) != 0 {
		reqRepo.WrtId = req.WrtIds
	}

	var pickingOrderAssign []*model.PickingOrderAssign
	if pickingOrderAssign, _, err = s.RepositoryPickingOrderAssign.Get(ctx, reqRepo); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order assign")
		return
	}

	var response []*pb.SalesOrderToCheck
	for _, v := range pickingOrderAssign {
		var (
			pickingOrderItem    []*model.PickingOrderItem
			totalItemOnProgress int64
			totalItem           int64
			deliveryKoli        []*model.DeliveryKoli
			totalKoli           float64
		)

		// get total item and item on progress
		if pickingOrderItem, _, err = s.RepositoryPickingOrderItem.Get(ctx, &dto.PickingOrderItemGetRequest{
			PickingOrderAssignId: []int64{v.ID},
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("picking order item")
			return
		}

		// count total item and item on progress
		for _, v2 := range pickingOrderItem {
			totalItem++
			if v2.Status != 16 {
				totalItemOnProgress++
			}
		}

		// count delivery koli
		if deliveryKoli, _, err = s.RepositoryDeliveryKoli.Get(ctx, &dto.DeliveryKoliGetRequest{
			SopNumber: v.SopNumber,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("delivery koli")
			return
		}
		for _, v2 := range deliveryKoli {
			totalKoli += v2.Quantity
		}

		// get picking order assign's sales order information
		var salesOrder *bridgeService.GetSalesOrderGPListResponse
		if salesOrder, err = s.opt.Client.BridgeServiceGrpc.GetSalesOrderListGPByID(ctx, &bridgeService.GetSalesOrderGPListByIDRequest{
			Id: v.SopNumber,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "sales order")
			return
		}
		fmt.Println(req.Custname != " ", ">>>>>>>", salesOrder.Data[0].Customer[0].Custname, "=============================", req.Custname, "----", strings.Contains(strings.ToLower(salesOrder.Data[0].Customer[0].Custname), strings.ToLower(req.Custname)))
		// filtering by customer name
		if !(req.Custname != " " && strings.Contains(strings.ToLower(salesOrder.Data[0].Customer[0].Custname), strings.ToLower(req.Custname))) {
			continue
		}

		// get wrt
		var wrt *bridgeService.GetWrtGPResponse
		if wrt, err = s.opt.Client.BridgeServiceGrpc.GetWrtGPList(ctx, &bridgeService.GetWrtGPListRequest{
			Search:    salesOrder.Data[0].Wrt[0].GnL_WRT_ID,
			GnlRegion: salesOrder.Data[0].GnL_Region,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "wrt")
			return
		}

		// get picker name
		var pickingOrder *model.PickingOrder
		if pickingOrder, err = s.RepositoryPickingOrder.GetDetail(ctx, v.PickingOrderId, ""); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("picking order")
			return
		}

		var picker *bridgeService.GetHelperGPResponse
		if picker, err = s.opt.Client.BridgeServiceGrpc.GetHelperGPDetail(ctx, &bridgeService.GetHelperGPDetailRequest{
			Id: pickingOrder.PickerId,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "helper")
			return
		}

		// TODO DELETE CODE
		salesOrderResponse := &pb.SalesOrderToCheck{
			SopNumber:    v.SopNumber,
			MerchantName: salesOrder.Data[0].Customer[0].Custname,
			DeliveryDate: salesOrder.Data[0].ReqShipDate,
			Wrt:          wrt.Data[0].Strttime + "-" + wrt.Data[0].Endtime,
			SopNote:      salesOrder.Data[0].Commntid,
			// MerchantName:        "DUMMY MERCHANT",
			// DeliveryDate:        "DUMMY DELIVERY DATE",
			// Wrt:                 "DUMMY WRT",
			// SopNote:             "DUMMY SOP NOTE",
			TotalItemOnProgress: totalItemOnProgress,
			TotalItem:           totalItem,
			TotalKoli:           totalKoli,
			CheckerName:         "", // filled below if exist
			PickerName:          picker.Data[0].GnlHelperName,
			Status:              int32(v.Status),
		}

		// get checker name
		var checker *bridgeService.GetHelperGPResponse
		if v.CheckerIdGp != "" {
			if checker, err = s.opt.Client.BridgeServiceGrpc.GetHelperGPDetail(ctx, &bridgeService.GetHelperGPDetailRequest{
				Id: v.CheckerIdGp,
			}); err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "helper")
				return
			}

			salesOrderResponse.CheckerName = checker.Data[0].GnlHelperName
		}

		response = append(response, salesOrderResponse)
	}

	res = &pb.GetSalesOrderToCheckResponse{
		Data: response,
	}

	return
}

func (s *PickingOrderService) SPVGetSalesOrderToCheckDetail(ctx context.Context, req *pb.GetSalesOrderToCheckDetailRequest) (res *pb.GetSalesOrderToCheckDetailResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PickingOrderService.SPVAcceptSalesOrder")
	defer span.End()

	// get sales order inside the picking order
	var pickingOrderAssign *model.PickingOrderAssign
	if pickingOrderAssign, err = s.RepositoryPickingOrderAssign.GetByID(ctx, 0, req.SopNumber); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order assign")
		return
	}

	// if status not need approval
	if pickingOrderAssign.Status != 35 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorMustStatus("picking order assign", "need approval")
		return
	}

	// get all item
	var pickingOrderItem []*model.PickingOrderItem
	if pickingOrderItem, _, err = s.RepositoryPickingOrderItem.Get(ctx, &dto.PickingOrderItemGetRequest{
		PickingOrderAssignId: []int64{pickingOrderAssign.ID},
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order item")
		return
	}

	// get picker name
	var pickingOrder *model.PickingOrder
	if pickingOrder, err = s.RepositoryPickingOrder.GetDetail(ctx, pickingOrderAssign.PickingOrderId, ""); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order")
		return
	}

	var picker *bridgeService.GetHelperGPResponse
	if picker, err = s.opt.Client.BridgeServiceGrpc.GetHelperGPDetail(ctx, &bridgeService.GetHelperGPDetailRequest{
		Id: pickingOrder.PickerId,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "helper")
		return
	}

	// prepare item response
	var (
		itemResponse        []*pb.PickingOrderItem
		totalItemOnProgress int64
	)
	for _, v := range pickingOrderItem {
		// get product detail
		var productDetail *bridgeService.GetItemGPResponse
		if productDetail, err = s.opt.Client.BridgeServiceGrpc.GetItemGPDetail(ctx, &bridgeService.GetItemGPDetailRequest{
			Id: v.ItemNumber,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("inventory", "product detail")
			return
		}

		// get product image
		var productImage *catalog_service.GetItemDetailResponse
		if productImage, err = s.opt.Client.CatalogServiceGrpc.GetItemDetail(ctx, &catalog_service.GetItemDetailRequest{
			Id: v.ItemNumber,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("inventory", "product image")
			return
		}

		item := &pb.PickingOrderItem{
			Id:                   v.Id,
			PickingOrderAssignId: v.PickingOrderAssignId,
			ItemNumber:           v.ItemNumber,
			ItemName:             productDetail.Data[0].Itemdesc,
			Picture:              "", // filled below if exist
			OrderQty:             v.OrderQuantity,
			PickQty:              v.PickQuantity,
			CheckQty:             v.CheckQuantity,
			ExcessQty:            v.ExcessQuantity,
			UnfulfillNote:        v.UnfulfillNote,
			Uom:                  productDetail.Data[0].Uomschdl,
			Status:               int32(v.Status),
		}

		if len(productImage.Data.ItemImage) > 0 {
			item.Picture = productImage.Data.ItemImage[0].ImageUrl
		}

		itemResponse = append(itemResponse, item)

		if item.Status != 16 {
			totalItemOnProgress++
		}
	}

	// count delivery koli
	var (
		deliveryKoli []*model.DeliveryKoli
		totalKoli    float64
	)
	if deliveryKoli, _, err = s.RepositoryDeliveryKoli.Get(ctx, &dto.DeliveryKoliGetRequest{
		SopNumber: pickingOrderAssign.SopNumber,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("delivery koli")
		return
	}
	for _, v2 := range deliveryKoli {
		totalKoli += v2.Quantity
	}

	// get picking order assign's sales order information
	var salesOrder *bridgeService.GetSalesOrderGPListResponse
	if salesOrder, err = s.opt.Client.BridgeServiceGrpc.GetSalesOrderListGPByID(ctx, &bridgeService.GetSalesOrderGPListByIDRequest{
		Id: pickingOrderAssign.SopNumber,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "sales order")
		return
	}

	// get wrt
	var wrt *bridgeService.GetWrtGPResponse
	if wrt, err = s.opt.Client.BridgeServiceGrpc.GetWrtGPList(ctx, &bridgeService.GetWrtGPListRequest{
		Search:    salesOrder.Data[0].Wrt[0].GnL_WRT_ID,
		GnlRegion: salesOrder.Data[0].GnL_Region,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "wrt")
		return
	}

	// TODO DELETE CODE
	res = &pb.GetSalesOrderToCheckDetailResponse{
		SopNumber:    pickingOrderAssign.SopNumber,
		DeliveryDate: salesOrder.Data[0].ReqShipDate,
		MerchantName: salesOrder.Data[0].Customer[0].Custname,
		// DeliveryDate: "DUMMY DELIVERY DATE",
		// MerchantName: "DUMMY MERCHANT",
		SopNote: salesOrder.Data[0].Commntid,
		Wrt:     wrt.Data[0].Strttime + "-" + wrt.Data[0].Endtime,
		// SopNote:             "DUMMY SOP NOTE",
		// Wrt:                 "DUMMY WRT",
		PickerName:          picker.Data[0].GnlHelperName,
		TotalKoli:           totalKoli,
		TotalItemOnProgress: totalItemOnProgress,
		TotalItem:           int64(len(itemResponse)),
		Item:                itemResponse,
		Status:              int32(pickingOrderAssign.Status),
	}

	return
}

func (s *PickingOrderService) SPVRejectSalesOrder(ctx context.Context, req *pb.SPVRejectSalesOrderRequest) (res *pb.SuccessResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PickingOrderService.SPVRejectSalesOrder")
	defer span.End()

	// get sales order inside the picking order
	var pickingOrderAssign *model.PickingOrderAssign
	if pickingOrderAssign, err = s.RepositoryPickingOrderAssign.GetByID(ctx, 0, req.SopNumber); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order assign")
		return
	}

	// if status not need approval
	if pickingOrderAssign.Status != 35 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorMustStatus("picking order assign", "need approval")
		return
	}

	var pickingOrder *model.PickingOrder
	if pickingOrder, err = s.RepositoryPickingOrder.GetDetail(ctx, pickingOrderAssign.PickingOrderId, ""); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order")
		return
	}

	var pickingOrderItem []*model.PickingOrderItem
	if pickingOrderItem, _, err = s.RepositoryPickingOrderItem.Get(ctx, &dto.PickingOrderItemGetRequest{
		Status:               []int{34},
		PickingOrderAssignId: []int64{pickingOrderAssign.ID},
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order item")
		return
	}

	// update item status that was unfulfill to rejected
	for _, v := range pickingOrderItem {
		v.Status = 9
		if err = s.RepositoryPickingOrderItem.Update(ctx, v, "status"); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("picking order item")
			return
		}
	}

	// update picking order assign status  to rejected
	pickingOrderAssign.Status = 9
	if err = s.RepositoryPickingOrderAssign.Update(ctx, pickingOrderAssign, "status"); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order assign")
		return
	}

	// update picking order status to rejected
	pickingOrder.Status = 9
	if err = s.RepositoryPickingOrder.Update(ctx, pickingOrder, "status"); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order")
		return
	}

	// delete all delivery koli of the SO
	// get delivery koli
	var deliveryKoli []*model.DeliveryKoli
	if deliveryKoli, _, err = s.RepositoryDeliveryKoli.Get(ctx, &dto.DeliveryKoliGetRequest{
		SopNumber: pickingOrderAssign.SopNumber,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("delivery koli")
		return
	}

	// delete delivery koli
	for _, v := range deliveryKoli {
		if err = s.RepositoryDeliveryKoli.Delete(ctx, v); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("delivery koli")
			return
		}
	}

	// send notification
	var notificationToken *model.HelperToken
	if notificationToken, err = s.RepositoryHelperToken.GetByHelperId(ctx, pickingOrder.PickerId); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("notification token")
		return
	}
	if _, err = s.opt.Client.NotificationServiceGrpc.SendNotificationHelper(ctx, &notification_service.SendNotificationHelperRequest{
		SendTo:    notificationToken.NotificationToken,
		NotifCode: "NOT0010",
		Type:      "4",
		RefId:     pickingOrderAssign.SopNumber,
		StaffId:   notificationToken.HelperIdGp,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("notification")
		return
	}

	// audit log
	if _, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &audit_service.CreateLogRequest{
		Log: &audit_service.Log{
			UserId:      0,
			ReferenceId: pickingOrderAssign.SopNumber,
			Type:        "picking order assign",
			Function:    "spv reject",
			CreatedAt:   timestamppb.New(time.Now()),
			Note:        "SPV with id of " + req.SpvId + " rejected " + pickingOrderAssign.SopNumber,
		},
	}); err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		span.RecordError(err)
		err = edenlabs.ErrorRpcNotFound("audit", "audit_log")
		return
	}

	res = &pb.SuccessResponse{
		Success: true,
	}

	return
}

func (s *PickingOrderService) SPVAcceptSalesOrder(ctx context.Context, req *pb.SPVAcceptSalesOrderRequest) (res *pb.SuccessResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PickingOrderService.SPVAcceptSalesOrder")
	defer span.End()

	// get sales order inside the picking order
	var pickingOrderAssign *model.PickingOrderAssign
	if pickingOrderAssign, err = s.RepositoryPickingOrderAssign.GetByID(ctx, 0, req.SopNumber); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order assign")
		return
	}

	// if status not need approval
	if pickingOrderAssign.Status != 35 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorMustStatus("picking order assign", "need approval")
		return
	}

	var pickingOrderItem []*model.PickingOrderItem
	if pickingOrderItem, _, err = s.RepositoryPickingOrderItem.Get(ctx, &dto.PickingOrderItemGetRequest{
		Status:               []int{34},
		PickingOrderAssignId: []int64{pickingOrderAssign.ID},
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order item")
		return
	}

	// update item status that was unfulfill to rejected
	for _, v := range pickingOrderItem {
		v.Status = 16
		if err = s.RepositoryPickingOrderItem.Update(ctx, v, "status"); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("picking order item")
			return
		}
	}

	// update picking order assign status  to rejected
	pickingOrderAssign.Status = 16
	if err = s.RepositoryPickingOrderAssign.Update(ctx, pickingOrderAssign, "status"); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order assign")
		return
	}

	// send notification
	var pickingOrder *model.PickingOrder
	if pickingOrder, err = s.RepositoryPickingOrder.GetDetail(ctx, pickingOrderAssign.PickingOrderId, ""); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order")
		return
	}

	var notificationToken *model.HelperToken
	if notificationToken, err = s.RepositoryHelperToken.GetByHelperId(ctx, pickingOrder.PickerId); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("notification token")
		return
	}
	if _, err = s.opt.Client.NotificationServiceGrpc.SendNotificationHelper(ctx, &notification_service.SendNotificationHelperRequest{
		SendTo:    notificationToken.NotificationToken,
		NotifCode: "NOT0011",
		Type:      "4",
		RefId:     pickingOrderAssign.SopNumber,
		StaffId:   notificationToken.HelperIdGp,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("notification")
		return
	}

	// audit log
	if _, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &audit_service.CreateLogRequest{
		Log: &audit_service.Log{
			UserId:      0,
			ReferenceId: pickingOrderAssign.SopNumber,
			Type:        "picking order assign",
			Function:    "spv accept",
			CreatedAt:   timestamppb.New(time.Now()),
			Note:        "SPV with id of " + req.SpvId + " accepted " + pickingOrderAssign.SopNumber,
		},
	}); err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		span.RecordError(err)
		err = edenlabs.ErrorRpcNotFound("audit", "audit_log")
		return
	}

	res = &pb.SuccessResponse{
		Success: true,
	}

	return
}

func (s *PickingOrderService) SPVWidget(ctx context.Context, req *pb.SPVWidgetRequest) (res *pb.SPVWidgetResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PickingOrderService.SPVWidget")
	defer span.End()

	// get sales order
	var pickingOrderAssign []*model.PickingOrderAssign
	if pickingOrderAssign, _, err = s.RepositoryPickingOrderAssign.Get(ctx, &dto.PickingOrderAssignGetRequest{
		SiteID:           []string{req.SiteIdGp},
		DeliveryDateFrom: time.Now().Add(-time.Hour * 48),
		DeliveryDateTo:   time.Now().Add(time.Hour * 24),
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order assign")
		return
	}

	res = &pb.SPVWidgetResponse{}

	for _, v := range pickingOrderAssign {
		if v.Status == 6 { // new
			res.TotalNew++
		}

		if v.Status == 21 || v.Status == 9 { // on progress & rejected
			res.TotalOnProgress++
		}

		if v.Status == 35 { // need approval
			res.TotalNeedApproval++
		}

		if v.Status == 2 { // finished
			res.TotalFinished++
		}

		res.TotalSalesOrder++
	}

	if res.TotalSalesOrder != 0 {
		res.TotalOnProgressPercentage = math.Round(float64(res.TotalOnProgress)/float64(res.TotalSalesOrder)*100) / 100
		res.TotalNeedApprovalPercentage = math.Round(float64(res.TotalNeedApproval)/float64(res.TotalSalesOrder)*100) / 100
		res.TotalFinishedPercentage = math.Round(float64(res.TotalFinished)/float64(res.TotalSalesOrder)*100) / 100
	}

	return
}

func (s *PickingOrderService) SPVWrtMonitoring(ctx context.Context, req *pb.GetWrtMonitoringListRequest) (res *pb.GetWrtMonitoringListResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PickingOrderService.SPVWrtMonitoring")
	defer span.End()

	// get site's wrt
	// get site detail
	var site *bridgeService.GetSiteGPResponse
	if site, err = s.opt.Client.BridgeServiceGrpc.GetSiteGPDetail(ctx, &bridgeService.GetSiteGPDetailRequest{
		Id: req.SiteId,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "site")
		return
	}

	// get region of the adm division
	var admDivision *bridgeService.GetAdmDivisionGPResponse
	if admDivision, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPList(ctx, &bridgeService.GetAdmDivisionGPListRequest{
		Limit:           1000,
		AdmDivisionCode: site.Data[0].GnlAdministrativeCode,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "adm division")
		return
	}

	// get wrt of the region
	var wrt *bridgeService.GetWrtGPResponse
	if wrt, err = s.opt.Client.BridgeServiceGrpc.GetWrtGPList(ctx, &bridgeService.GetWrtGPListRequest{
		Limit:     1000,
		GnlRegion: admDivision.Data[0].Region,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "wrt")
		return
	}

	var response []*pb.WrtMonitoring
	if req.Type == 1 { // filter picker
		currentTime := time.Now()
		yesterday := currentTime.Add(-time.Hour * 24)
		tomorrow := currentTime.Add(time.Hour * 24)

		// get internal id for picking order
		var pickingOrderArray []int64
		for _, helperId := range req.HelperId {
			var pickingOrder *bridgeService.GetPickingOrderGPHeaderResponse
			if pickingOrder, err = s.opt.Client.BridgeServiceGrpc.GetPickingOrderGPHeader(ctx, &bridgeService.GetPickingOrderGPHeaderRequest{
				Limit:       1000,
				DocdateFrom: yesterday.Format("2006-01-02"),
				DocdateTo:   tomorrow.Format("2006-01-02"),
				GnlHelperId: helperId,
			}); err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "picking order")
				return
			}

			for _, v := range pickingOrder.Data {
				var pickingOrder *model.PickingOrder
				if pickingOrder, err = s.RepositoryPickingOrder.GetDetail(ctx, 0, v.Docnumbr); err != nil {
					span.RecordError(err)
					s.opt.Logger.AddMessage(log.ErrorLevel, err)
					err = edenlabs.ErrorInvalid("picking order")
					return
				}
				pickingOrderArray = append(pickingOrderArray, pickingOrder.Id)
			}
		}

		for _, v := range wrt.Data {
			var totalSO int64
			if _, totalSO, err = s.RepositoryPickingOrderAssign.Get(ctx, &dto.PickingOrderAssignGetRequest{
				DeliveryDateFrom: time.Now().Add(-time.Hour * 48),
				DeliveryDateTo:   time.Now().Add(time.Hour * 24),
				SiteID:           []string{req.SiteId},
				WrtId:            []string{v.GnL_WRT_ID},
				PickingOrderId:   pickingOrderArray,
			}); err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorInvalid("picking order assign")
				return
			}

			var totalOnProgress int64
			if _, totalOnProgress, err = s.RepositoryPickingOrderAssign.Get(ctx, &dto.PickingOrderAssignGetRequest{
				DeliveryDateFrom: time.Now().Add(-time.Hour * 48),
				DeliveryDateTo:   time.Now().Add(time.Hour * 24),
				SiteID:           []string{req.SiteId},
				WrtId:            []string{v.GnL_WRT_ID},
				Status:           []int{6, 9, 21, 35},
				PickingOrderId:   pickingOrderArray,
			}); err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorInvalid("picking order assign")
				return
			}
			var onProgressPercentage float64
			if totalOnProgress != 0 {
				onProgressPercentage = math.Round(float64(totalOnProgress)/float64(totalSO)*100) / 100
			}

			var totalFinished int64
			if _, totalFinished, err = s.RepositoryPickingOrderAssign.Get(ctx, &dto.PickingOrderAssignGetRequest{
				DeliveryDateFrom: time.Now().Add(-time.Hour * 48),
				DeliveryDateTo:   time.Now().Add(time.Hour * 24),
				SiteID:           []string{req.SiteId},
				WrtId:            []string{v.GnL_WRT_ID},
				Status:           []int{16},
				PickingOrderId:   pickingOrderArray,
			}); err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorInvalid("picking order assign")
				return
			}
			var finishedPercentage float64
			if totalFinished != 0 {
				finishedPercentage = math.Round(float64(totalFinished)/float64(totalSO)*100) / 100
			}

			response = append(response, &pb.WrtMonitoring{
				WrtId:                v.GnL_WRT_ID,
				WrtDesc:              v.Strttime + "-" + v.Endtime,
				CountSo:              totalSO,
				OnProgress:           totalOnProgress,
				OnProgressPercentage: onProgressPercentage,
				Finished:             totalFinished,
				FinishedPercentage:   finishedPercentage,
			})
		}
	} else if req.Type == 2 { // filter checker
		for _, v := range wrt.Data {
			var totalSO int64
			if _, totalSO, err = s.RepositoryPickingOrderAssign.Get(ctx, &dto.PickingOrderAssignGetRequest{
				DeliveryDateFrom: time.Now().Add(-time.Hour * 48),
				DeliveryDateTo:   time.Now().Add(time.Hour * 24),
				SiteID:           []string{req.SiteId},
				WrtId:            []string{v.GnL_WRT_ID},
				CheckerId:        req.HelperId,
			}); err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorInvalid("picking order assign")
				return
			}

			var totalOnProgress int64
			if _, totalOnProgress, err = s.RepositoryPickingOrderAssign.Get(ctx, &dto.PickingOrderAssignGetRequest{
				DeliveryDateFrom: time.Now().Add(-time.Hour * 48),
				DeliveryDateTo:   time.Now().Add(time.Hour * 24),
				SiteID:           []string{req.SiteId},
				WrtId:            []string{v.GnL_WRT_ID},
				Status:           []int{20},
				CheckerId:        req.HelperId,
			}); err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorInvalid("picking order assign")
				return
			}
			var onProgressPercentage float64
			if totalOnProgress != 0 {
				onProgressPercentage = math.Round(float64(totalOnProgress)/float64(totalSO)*100) / 100
			}

			var totalFinished int64
			if _, totalFinished, err = s.RepositoryPickingOrderAssign.Get(ctx, &dto.PickingOrderAssignGetRequest{
				DeliveryDateFrom: time.Now().Add(-time.Hour * 48),
				DeliveryDateTo:   time.Now().Add(time.Hour * 24),
				SiteID:           []string{req.SiteId},
				WrtId:            []string{v.GnL_WRT_ID},
				Status:           []int{2},
				CheckerId:        req.HelperId,
			}); err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorInvalid("picking order assign")
				return
			}
			var finishedPercentage float64
			if totalFinished != 0 {
				finishedPercentage = math.Round(float64(totalFinished)/float64(totalSO)*100) / 100
			}

			response = append(response, &pb.WrtMonitoring{
				WrtId:                v.GnL_WRT_ID,
				WrtDesc:              v.Strttime + "-" + v.Endtime,
				CountSo:              totalSO,
				OnProgress:           totalOnProgress,
				OnProgressPercentage: onProgressPercentage,
				Finished:             totalFinished,
				FinishedPercentage:   finishedPercentage,
			})

		}
	}

	res = &pb.GetWrtMonitoringListResponse{
		Data: response,
	}

	return
}

func (s *PickingOrderService) SPVWrtMonitoringDetail(ctx context.Context, req *pb.GetWrtMonitoringDetailRequest) (res *pb.GetWrtMonitoringDetailResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PickingOrderService.SPVWrtMonitoringDetail")
	defer span.End()

	var response []*pb.WrtMonitoringDetail
	if req.Type == 1 {
		// get picker picking order
		currentTime := time.Now()
		yesterday := currentTime.Add(-time.Hour * 24)
		tomorrow := currentTime.Add(time.Hour * 24)

		// get internal id for picking order
		var pickingOrderArray []int64
		for _, helperId := range req.HelperId {
			var pickingOrder *bridgeService.GetPickingOrderGPHeaderResponse
			if pickingOrder, err = s.opt.Client.BridgeServiceGrpc.GetPickingOrderGPHeader(ctx, &bridgeService.GetPickingOrderGPHeaderRequest{
				Limit:       1000,
				DocdateFrom: yesterday.Format("2006-01-02"),
				DocdateTo:   tomorrow.Format("2006-01-02"),
				GnlHelperId: helperId,
			}); err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "picking order")
				return
			}

			for _, v := range pickingOrder.Data {
				var pickingOrder *model.PickingOrder
				if pickingOrder, err = s.RepositoryPickingOrder.GetDetail(ctx, 0, v.Docnumbr); err != nil {
					span.RecordError(err)
					s.opt.Logger.AddMessage(log.ErrorLevel, err)
					err = edenlabs.ErrorInvalid("picking order")
					return
				}

				pickingOrderArray = append(pickingOrderArray, pickingOrder.Id)
			}
		}

		// get sales order inside the picking order
		var pickingOrderAssign []*model.PickingOrderAssign
		if pickingOrderAssign, _, err = s.RepositoryPickingOrderAssign.Get(ctx, &dto.PickingOrderAssignGetRequest{
			DeliveryDateFrom: time.Now().Add(-time.Hour * 48),
			DeliveryDateTo:   time.Now().Add(time.Hour * 24),
			SiteID:           []string{req.SiteId},
			WrtId:            []string{req.WrtId},
			Status:           []int{6, 9, 21, 35, 16},
			PickingOrderId:   pickingOrderArray,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("picking order assign")
			return
		}

		for _, v := range pickingOrderAssign {
			// get picking order assign's sales order information
			var salesOrder *bridgeService.GetSalesOrderGPListResponse
			if salesOrder, err = s.opt.Client.BridgeServiceGrpc.GetSalesOrderListGPByID(ctx, &bridgeService.GetSalesOrderGPListByIDRequest{
				Id: v.SopNumber,
			}); err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "sales order")
				return
			}

			// total koli
			var (
				deliveryKoli []*model.DeliveryKoli
				totalKoli    float64
			)
			if deliveryKoli, _, err = s.RepositoryDeliveryKoli.Get(ctx, &dto.DeliveryKoliGetRequest{
				SopNumber: v.SopNumber,
			}); err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorInvalid("delivery koli")
				return
			}
			for _, v2 := range deliveryKoli {
				totalKoli += v2.Quantity
			}

			// get picker name
			var pickingOrder *model.PickingOrder
			if pickingOrder, err = s.RepositoryPickingOrder.GetDetail(ctx, v.PickingOrderId, ""); err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorInvalid("picking order")
				return
			}

			var picker *bridgeService.GetHelperGPResponse
			if picker, err = s.opt.Client.BridgeServiceGrpc.GetHelperGPDetail(ctx, &bridgeService.GetHelperGPDetailRequest{
				Id: pickingOrder.PickerId,
			}); err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "helper")
				return
			}

			response = append(response, &pb.WrtMonitoringDetail{
				SopNumber:    v.SopNumber,
				MerchantName: salesOrder.Data[0].Customer[0].Custname,
				TotalKoli:    totalKoli,
				HelperCode:   picker.Data[0].GnlHelperId,
				HelperName:   picker.Data[0].GnlHelperName,
				Status:       int32(v.Status),
			})
		}
	} else if req.Type == 2 {
		// get sales order inside the picking order
		var pickingOrderAssign []*model.PickingOrderAssign
		if pickingOrderAssign, _, err = s.RepositoryPickingOrderAssign.Get(ctx, &dto.PickingOrderAssignGetRequest{
			DeliveryDateFrom: time.Now().Add(-time.Hour * 48),
			DeliveryDateTo:   time.Now().Add(time.Hour * 24),
			SiteID:           []string{req.SiteId},
			WrtId:            []string{req.WrtId},
			Status:           []int{20, 2},
			CheckerId:        req.HelperId,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("picking order assign")
			return
		}

		for _, v := range pickingOrderAssign {
			// get picking order assign's sales order information
			var salesOrder *bridgeService.GetSalesOrderGPListResponse
			if salesOrder, err = s.opt.Client.BridgeServiceGrpc.GetSalesOrderListGPByID(ctx, &bridgeService.GetSalesOrderGPListByIDRequest{
				Id: v.SopNumber,
			}); err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "sales order")
				return
			}

			// total koli
			var (
				deliveryKoli []*model.DeliveryKoli
				totalKoli    float64
			)
			if deliveryKoli, _, err = s.RepositoryDeliveryKoli.Get(ctx, &dto.DeliveryKoliGetRequest{
				SopNumber: v.SopNumber,
			}); err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorInvalid("delivery koli")
				return
			}
			for _, v2 := range deliveryKoli {
				totalKoli += v2.Quantity
			}

			// get picker name
			var pickingOrder *model.PickingOrder
			if pickingOrder, err = s.RepositoryPickingOrder.GetDetail(ctx, v.PickingOrderId, ""); err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorInvalid("picking order")
				return
			}

			var picker *bridgeService.GetHelperGPResponse
			if picker, err = s.opt.Client.BridgeServiceGrpc.GetHelperGPDetail(ctx, &bridgeService.GetHelperGPDetailRequest{
				Id: pickingOrder.PickerId,
			}); err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "helper")
				return
			}

			response = append(response, &pb.WrtMonitoringDetail{
				SopNumber:    v.SopNumber,
				MerchantName: salesOrder.Data[0].Customer[0].Custname,
				TotalKoli:    totalKoli,
				HelperCode:   picker.Data[0].GnlHelperId,
				HelperName:   picker.Data[0].GnlHelperName,
				Status:       int32(v.Status),
			})
		}
	}

	res = &pb.GetWrtMonitoringDetailResponse{
		Data: response,
	}

	return
}

func (s *PickingOrderService) CheckerGetSalesOrderToCheckDetail(ctx context.Context, req *pb.GetSalesOrderToCheckDetailRequest) (res *pb.GetSalesOrderToCheckDetailResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PickingOrderService.CheckerGetSalesOrderToCheckDetail")
	defer span.End()

	// get sales order inside the picking order
	var pickingOrderAssign *model.PickingOrderAssign
	if pickingOrderAssign, err = s.RepositoryPickingOrderAssign.GetByID(ctx, 0, req.SopNumber); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order assign")
		return
	}

	// if status not picked / checking
	if pickingOrderAssign.Status != 16 && pickingOrderAssign.Status != 20 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorMustStatus("picking order assign", "picked / checking")
		return
	}

	if pickingOrderAssign.Status == 20 {
		if pickingOrderAssign.CheckerIdGp != req.CheckerId {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("picking order assign")
			return
		}
	}

	var pickingOrderItem []*model.PickingOrderItem
	if pickingOrderItem, _, err = s.RepositoryPickingOrderItem.Get(ctx, &dto.PickingOrderItemGetRequest{
		PickingOrderAssignId: []int64{pickingOrderAssign.ID},
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order item")
		return
	}

	// get picker name
	var pickingOrder *model.PickingOrder
	if pickingOrder, err = s.RepositoryPickingOrder.GetDetail(ctx, pickingOrderAssign.PickingOrderId, ""); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order")
		return
	}

	var picker *bridgeService.GetHelperGPResponse
	if picker, err = s.opt.Client.BridgeServiceGrpc.GetHelperGPDetail(ctx, &bridgeService.GetHelperGPDetailRequest{
		Id: pickingOrder.PickerId,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "helper")
		return
	}

	var (
		itemResponse        []*pb.PickingOrderItem
		totalItemOnProgress int64
	)
	for _, v := range pickingOrderItem {
		var productDetail *bridgeService.GetItemGPResponse
		if productDetail, err = s.opt.Client.BridgeServiceGrpc.GetItemGPDetail(ctx, &bridgeService.GetItemGPDetailRequest{
			Id: v.ItemNumber,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("inventory", "product detail")
			return
		}

		var productImage *catalog_service.GetItemDetailResponse
		if productImage, err = s.opt.Client.CatalogServiceGrpc.GetItemDetail(ctx, &catalog_service.GetItemDetailRequest{
			Id: v.ItemNumber,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("inventory", "product image")
			return
		}

		item := &pb.PickingOrderItem{
			Id:                   v.Id,
			PickingOrderAssignId: v.PickingOrderAssignId,
			ItemNumber:           v.ItemNumber,
			ItemName:             productDetail.Data[0].Itemdesc,
			Picture:              "", // filled below if exist
			OrderQty:             v.OrderQuantity,
			PickQty:              v.PickQuantity,
			CheckQty:             v.CheckQuantity,
			ExcessQty:            v.ExcessQuantity,
			UnfulfillNote:        v.UnfulfillNote,
			Uom:                  productDetail.Data[0].Uomschdl,
			Status:               int32(v.Status),
		}

		if len(productImage.Data.ItemImage) > 0 {
			item.Picture = productImage.Data.ItemImage[0].ImageUrl
		}

		itemResponse = append(itemResponse, item)

		if item.Status != 16 {
			totalItemOnProgress++
		}
	}

	// count delivery koli
	var (
		deliveryKoli []*model.DeliveryKoli
		totalKoli    float64
	)
	if deliveryKoli, _, err = s.RepositoryDeliveryKoli.Get(ctx, &dto.DeliveryKoliGetRequest{
		SopNumber: pickingOrderAssign.SopNumber,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("delivery koli")
		return
	}
	for _, v2 := range deliveryKoli {
		totalKoli += v2.Quantity
	}

	// TODO UNCOMMENT GET DATA
	// get picking order assign's sales order information
	var salesOrder *bridgeService.GetSalesOrderGPListResponse
	if salesOrder, err = s.opt.Client.BridgeServiceGrpc.GetSalesOrderListGPByID(ctx, &bridgeService.GetSalesOrderGPListByIDRequest{
		Id: pickingOrderAssign.SopNumber,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "sales order")
		return
	}

	// get wrt
	var wrt *bridgeService.GetWrtGPResponse
	if wrt, err = s.opt.Client.BridgeServiceGrpc.GetWrtGPList(ctx, &bridgeService.GetWrtGPListRequest{
		Search:    salesOrder.Data[0].Wrt[0].GnL_WRT_ID,
		GnlRegion: salesOrder.Data[0].GnL_Region,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "wrt")
		return
	}

	// TODO DELETE CODE
	res = &pb.GetSalesOrderToCheckDetailResponse{
		SopNumber:    pickingOrderAssign.SopNumber,
		DeliveryDate: salesOrder.Data[0].ReqShipDate,
		MerchantName: salesOrder.Data[0].Customer[0].Custname,
		// DeliveryDate: "DUMMY DELIVERY DATE",
		// MerchantName: "DUMMY MERCHANT",
		SopNote: salesOrder.Data[0].Commntid,
		Wrt:     wrt.Data[0].Strttime + "-" + wrt.Data[0].Endtime,
		// SopNote:             "DUMMY SOP NOTE",
		// Wrt:                 "DUMMY WRT",
		PickerName:          picker.Data[0].GnlHelperName,
		TotalKoli:           totalKoli,
		TotalItemOnProgress: totalItemOnProgress,
		TotalItem:           int64(len(itemResponse)),
		Item:                itemResponse,
		Status:              int32(pickingOrderAssign.Status),
	}

	return
}

func (s *PickingOrderService) CheckerStartChecking(ctx context.Context, req *pb.CheckerStartCheckingRequest) (res *pb.SuccessResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PickingOrderService.CheckerStartChecking")
	defer span.End()

	var pickingOrderAssign *model.PickingOrderAssign
	if pickingOrderAssign, err = s.RepositoryPickingOrderAssign.GetByID(ctx, 0, req.SopNumber); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order assign")
		return
	}

	// if status not picked
	if pickingOrderAssign.Status != 16 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorMustStatus("picking order assign", "picked")
		return
	}

	// change status POA from picked to checking
	pickingOrderAssign.CheckerIdGp = req.CheckerId
	pickingOrderAssign.CheckingStartTime = time.Now()
	pickingOrderAssign.Status = 20
	if err = s.RepositoryPickingOrderAssign.Update(ctx, pickingOrderAssign, "CheckerIdGp", "CheckingStartTime", "Status"); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order assign")
		return
	}

	res = &pb.SuccessResponse{
		Success: true,
	}

	return
}

func (s *PickingOrderService) CheckerSubmitChecking(ctx context.Context, req *pb.CheckerSubmitCheckingRequest) (res *pb.SuccessResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PickingOrderService.CheckerSubmitChecking")
	defer span.End()

	var pickingOrderAssign *model.PickingOrderAssign
	if pickingOrderAssign, err = s.RepositoryPickingOrderAssign.GetByID(ctx, 0, req.SopNumber); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order assign")
		return
	}

	// if status not checking
	if pickingOrderAssign.Status != 20 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorMustStatus("picking order assign", "checking")
		return
	}

	if pickingOrderAssign.CheckerIdGp != req.CheckerId {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("Authorization", "This is not your job")
		return
	}

	for _, v := range req.Request {
		var pickingOrderItem *model.PickingOrderItem
		if pickingOrderItem, err = s.RepositoryPickingOrderItem.GetByID(ctx, &dto.PickingOrderItemGetDetailRequest{
			PickingOrderAssignId: pickingOrderAssign.ID,
			ItemNumber:           v.ItemNumber,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("picking order item")
			return
		}

		pickingOrderItem.CheckQuantity = v.CheckQty
		if err = s.RepositoryPickingOrderItem.Update(ctx, pickingOrderItem, "CheckQuantity"); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("picking order item")
			return
		}
	}

	res = &pb.SuccessResponse{
		Success: true,
	}

	return
}

func (s *PickingOrderService) CheckerRejectSalesOrder(ctx context.Context, req *pb.CheckerRejectSalesOrderRequest) (res *pb.SuccessResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PickingOrderService.CheckerRejectSalesOrder")
	defer span.End()

	var pickingOrderAssign *model.PickingOrderAssign
	if pickingOrderAssign, err = s.RepositoryPickingOrderAssign.GetByID(ctx, 0, req.SopNumber); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order assign")
		return
	}

	// if status not checking
	if pickingOrderAssign.Status != 20 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorMustStatus("picking order assign", "checking")
		return
	}

	// picking order assign has to be checker job
	if pickingOrderAssign.CheckerIdGp != req.CheckerId {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("Authorization", "This is not your job")
		return
	}

	var pickingOrder *model.PickingOrder
	if pickingOrder, err = s.RepositoryPickingOrder.GetDetail(ctx, pickingOrderAssign.PickingOrderId, ""); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order")
		return
	}

	// map to hold item that was rejected by checker
	itemNumberMap := map[string]bool{}
	for _, v := range req.ItemNumberReject {
		itemNumberMap[v] = true
	}

	var pickingOrderItem []*model.PickingOrderItem
	if pickingOrderItem, _, err = s.RepositoryPickingOrderItem.Get(ctx, &dto.PickingOrderItemGetRequest{
		PickingOrderAssignId: []int64{pickingOrderAssign.ID},
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order item")
		return
	}

	for _, v := range pickingOrderItem {
		if itemNumberMap[v.ItemNumber] {
			v.Status = 9
		}
		v.CheckQuantity = 0
		v.ExcessQuantity = 0
		if err = s.RepositoryPickingOrderItem.Update(ctx, v, "CheckQuantity", "ExcessQuantity", "Status"); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("picking order item")
			return
		}
	}

	// update picking order to rejected
	pickingOrderAssign.Status = 9
	if err = s.RepositoryPickingOrderAssign.Update(ctx, pickingOrderAssign, "status"); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order assign")
		return
	}

	// update picking order to rejected
	pickingOrder.Status = 9
	if err = s.RepositoryPickingOrder.Update(ctx, pickingOrder, "status"); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order")
		return
	}

	// delete all delivery koli of the SO
	// get delivery koli
	var deliveryKoli []*model.DeliveryKoli
	if deliveryKoli, _, err = s.RepositoryDeliveryKoli.Get(ctx, &dto.DeliveryKoliGetRequest{
		SopNumber: pickingOrderAssign.SopNumber,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("delivery koli")
		return
	}

	// delete delivery koli
	for _, v := range deliveryKoli {
		if err = s.RepositoryDeliveryKoli.Delete(ctx, v); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("delivery koli")
			return
		}
	}

	// send notification
	var notificationToken *model.HelperToken
	if notificationToken, err = s.RepositoryHelperToken.GetByHelperId(ctx, pickingOrder.PickerId); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("notification token")
		return
	}
	if _, err = s.opt.Client.NotificationServiceGrpc.SendNotificationHelper(ctx, &notification_service.SendNotificationHelperRequest{
		SendTo:    notificationToken.NotificationToken,
		NotifCode: "NOT0015",
		Type:      "4",
		RefId:     pickingOrderAssign.SopNumber,
		StaffId:   notificationToken.HelperIdGp,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("notification")
		return
	}

	// audit log
	if _, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &audit_service.CreateLogRequest{
		Log: &audit_service.Log{
			UserId:      0,
			ReferenceId: pickingOrderAssign.SopNumber,
			Type:        "picking order assign",
			Function:    "checker reject",
			CreatedAt:   timestamppb.New(time.Now()),
			Note:        "Checker with id of " + req.CheckerId + " rejected " + pickingOrderAssign.SopNumber,
		},
	}); err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		span.RecordError(err)
		err = edenlabs.ErrorRpcNotFound("audit", "audit_log")
		return
	}

	res = &pb.SuccessResponse{
		Success: true,
	}

	return
}

func (s *PickingOrderService) CheckerGetDeliveryKoli(ctx context.Context, req *pb.CheckerGetDeliveryKoliRequest) (res *pb.CheckerGetDeliveryKoliResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PickingOrderService.CheckerGetDeliveryKoli")
	defer span.End()

	var pickingOrderAssign *model.PickingOrderAssign
	if pickingOrderAssign, err = s.RepositoryPickingOrderAssign.GetByID(ctx, 0, req.SopNumber); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order assign")
		return
	}

	var koli []*model.Koli
	if koli, _, err = s.RepositoryKoli.Get(ctx, &dto.KoliGetRequest{}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("koli")
		return
	}

	var deliveryKoli []*model.DeliveryKoli
	if deliveryKoli, _, err = s.RepositoryDeliveryKoli.Get(ctx, &dto.DeliveryKoliGetRequest{
		SopNumber: pickingOrderAssign.SopNumber,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("delivery koli")
		return
	}

	deliveryKoliMap := map[int64]int64{}
	for _, v := range deliveryKoli {
		deliveryKoliMap[v.KoliId] = int64(v.Quantity)
	}

	var responseDeliveryKoli []*pb.DeliveryKoli
	for _, v := range koli {
		responseDeliveryKoli = append(responseDeliveryKoli, &pb.DeliveryKoli{
			SalesOrderCode: req.SopNumber,
			KoliId:         v.Id,
			Name:           v.Name,
			Quantity:       float64(deliveryKoliMap[v.Id]),
		})
	}

	res = &pb.CheckerGetDeliveryKoliResponse{
		Data: responseDeliveryKoli,
	}

	return
}

func (s *PickingOrderService) CheckerAcceptSalesOrder(ctx context.Context, req *pb.CheckerAcceptSalesOrderRequest) (res *pb.CheckerAcceptSalesOrderResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PickingOrderService.CheckerAcceptSalesOrder")
	defer span.End()

	var (
		pickingOrderAssign *model.PickingOrderAssign
		koli               *model.Koli
		koliStr            string
	)
	if pickingOrderAssign, err = s.RepositoryPickingOrderAssign.GetByID(ctx, 0, req.SopNumber); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order assign")
		return
	}

	// if status not checking
	if pickingOrderAssign.Status != 20 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorMustStatus("picking order assign", "checking")
		return
	}

	// picking order assign has to be checker job
	if pickingOrderAssign.CheckerIdGp != req.CheckerId {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("Authorization", "This is not your job")
		return
	}

	// validate first then
	// delete then create koli again
	var haveKoli bool
	var qtyKoli string
	for _, v := range req.Koli {
		if koli, err = s.RepositoryKoli.GetByID(ctx, v.Id); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("koli")
			return
		}
		qtyKoli = fmt.Sprintf("%.0f", v.Quantity)
		koliStr += qtyKoli + koli.Value + ", "

		if v.Quantity > 0 {
			haveKoli = true
		}
	}

	// if all koli quantity is zero
	if !haveKoli {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRequired("koli")
		return
	}

	// get delivery koli
	var deliveryKoli []*model.DeliveryKoli
	if deliveryKoli, _, err = s.RepositoryDeliveryKoli.Get(ctx, &dto.DeliveryKoliGetRequest{
		SopNumber: pickingOrderAssign.SopNumber,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("delivery koli")
		return
	}

	// delete delivery koli
	for _, v := range deliveryKoli {
		if err = s.RepositoryDeliveryKoli.Delete(ctx, v); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("delivery koli")
			return
		}
	}

	// create delivery koli
	var totalKoli float64
	for _, v := range req.Koli {
		writeDeliveryKoli := &model.DeliveryKoli{
			SopNumber: pickingOrderAssign.SopNumber,
			KoliId:    v.Id,
			Quantity:  v.Quantity,
		}
		totalKoli += v.Quantity

		if err = s.RepositoryDeliveryKoli.Create(ctx, writeDeliveryKoli); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("delivery koli")
			return
		}
	}
	// send notification
	var pickingOrder *model.PickingOrder
	if pickingOrder, err = s.RepositoryPickingOrder.GetDetail(ctx, pickingOrderAssign.PickingOrderId, ""); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order")
		return
	}

	// submit checking
	var (
		checkingRequest []*bridgeService.Checking
		checkingDetails []*bridgeService.CheckingDetails
	)

	var pickingOrderItems []*model.PickingOrderItem
	if pickingOrderItems, _, err = s.RepositoryPickingOrderItem.Get(ctx, &dto.PickingOrderItemGetRequest{
		PickingOrderAssignId: []int64{pickingOrderAssign.ID},
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order assign")
		return
	}

	/*
		if CheckQuantity > OrderQuantity, GP should create exxcess qty
	*/
	for _, v := range pickingOrderItems {
		checkingDetails = append(checkingDetails, &bridgeService.CheckingDetails{
			Sopnumbe:     pickingOrderAssign.SopNumber,
			Lnitmseq:     int32(v.IdGp),
			IvmQtyPickso: v.CheckQuantity,
		})
	}
	checkingRequest = append(checkingRequest, &bridgeService.Checking{
		Docnumbr:     "",
		Sopnumbe:     pickingOrderAssign.SopNumber,
		Strttime:     pickingOrderAssign.CheckingStartTime.Format("15:04:05"),
		Endtime:      pickingOrderAssign.CheckingEndTime.Format("15:04:05"),
		WmsPickerId:  pickingOrderAssign.CheckerIdGp,
		IvmKoli:      int32(totalKoli),
		ImvJenisKoli: strings.TrimSuffix(koliStr, ", "),
		Details:      checkingDetails,
	})

	var submitCheckingResponse *bridgeService.SubmitPickingCheckingResponse
	if submitCheckingResponse, err = s.opt.Client.BridgeServiceGrpc.SubmitPickingCheckingPickingOrder(ctx, &bridgeService.SubmitPickingCheckingRequest{
		Bachnumb: "SUBMIT",
		Picking: &bridgeService.Picking{
			Docnumbr: pickingOrder.DocNumber,
		},
		Checking: checkingRequest,
	}); submitCheckingResponse.Code != 200 || err != nil {
		if err != nil {
			err = edenlabs.ErrorRpcNotFound("bridge", "submit checking")
		} else if submitCheckingResponse.Code != 200 {
			err = edenlabs.ErrorRpc(submitCheckingResponse.Message)
		}
		s.opt.Logger.AddMessage(log.ErrorLevel, err.Error())

		span.RecordError(err)

		return
	}

	// update picking order assign
	// so before GP success, it will not update date in eden
	pickingOrderAssign.CheckingEndTime = time.Now()
	pickingOrderAssign.Status = 2
	if err = s.RepositoryPickingOrderAssign.Update(ctx, pickingOrderAssign, "CheckingEndTime", "status"); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order assign")
		return
	}

	var notificationToken *model.HelperToken
	if notificationToken, err = s.RepositoryHelperToken.GetByHelperId(ctx, pickingOrder.PickerId); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("notification token")
		return
	}
	if _, err = s.opt.Client.NotificationServiceGrpc.SendNotificationHelper(ctx, &notification_service.SendNotificationHelperRequest{
		SendTo:    notificationToken.NotificationToken,
		NotifCode: "NOT0016",
		Type:      "4",
		RefId:     pickingOrderAssign.SopNumber,
		StaffId:   notificationToken.HelperIdGp,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("notification")
		return
	}

	// audit log
	if _, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &audit_service.CreateLogRequest{
		Log: &audit_service.Log{
			UserId:      0,
			ReferenceId: pickingOrderAssign.SopNumber,
			Type:        "picking order assign",
			Function:    "checker accept",
			CreatedAt:   timestamppb.New(time.Now()),
			Note:        "Checker with id of " + req.CheckerId + " accepted " + pickingOrderAssign.SopNumber,
		},
	}); err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		span.RecordError(err)
		err = edenlabs.ErrorRpcNotFound("audit", "audit_log")
		return
	}

	// update picking order assign status to finished if no poa anymore
	var total int64
	if _, total, err = s.RepositoryPickingOrderAssign.Get(ctx, &dto.PickingOrderAssignGetRequest{
		Status:         []int{6, 21, 35, 16, 20, 9},
		PickingOrderId: []int64{pickingOrderAssign.PickingOrderId},
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order assign")
		return
	}

	if total == 0 {
		var pickingOrder *model.PickingOrder
		if pickingOrder, err = s.RepositoryPickingOrder.GetDetail(ctx, pickingOrderAssign.PickingOrderId, ""); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("picking order")
			return
		}

		pickingOrder.Status = 2
		if err = s.RepositoryPickingOrder.Update(ctx, pickingOrder, "Status"); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("picking order")
			return
		}
	}

	res = &pb.CheckerAcceptSalesOrderResponse{
		Success:       true,
		DeliveryOrder: submitCheckingResponse.Data.DeliveryOrder.Docnumbr,
		SalesInvoice:  submitCheckingResponse.Data.SalesInvoice.Docnumbr,
	}

	return
}

func (s *PickingOrderService) CheckerHistory(ctx context.Context, req *pb.CheckerHistoryRequest) (res *pb.CheckerHistoryResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PickingOrderService.CheckerHistory")
	defer span.End()

	var sopNumber []string
	if req.SopNumber != "" {
		sopNumber = append(sopNumber, req.SopNumber)
	}

	// get sales order inside the picking order
	var pickingOrderAssign []*model.PickingOrderAssign
	if pickingOrderAssign, _, err = s.RepositoryPickingOrderAssign.Get(ctx, &dto.PickingOrderAssignGetRequest{
		Offset:           int(req.Offset),
		Limit:            int(req.Limit),
		Status:           []int{2},
		SopNumber:        sopNumber,
		DeliveryDateFrom: time.Now().Add(-time.Hour * 48),
		DeliveryDateTo:   time.Now().Add(time.Hour * 24),
		CheckerId:        []string{req.CheckerId},
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order assign")
		return
	}

	var response []*pb.SalesOrderToCheck
	for _, v := range pickingOrderAssign {
		var (
			pickingOrderItem    []*model.PickingOrderItem
			totalItemOnProgress int64
			totalItem           int64
			deliveryKoli        []*model.DeliveryKoli
			totalKoli           float64
		)

		// get total item and item on progress
		if pickingOrderItem, _, err = s.RepositoryPickingOrderItem.Get(ctx, &dto.PickingOrderItemGetRequest{
			PickingOrderAssignId: []int64{v.ID},
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("picking order item")
			return
		}

		// count total item and item on progress
		for _, v2 := range pickingOrderItem {
			totalItem++
			if v2.Status != 16 {
				totalItemOnProgress++
			}
		}

		// count delivery koli
		if deliveryKoli, _, err = s.RepositoryDeliveryKoli.Get(ctx, &dto.DeliveryKoliGetRequest{
			SopNumber: v.SopNumber,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("delivery koli")
			return
		}
		for _, v2 := range deliveryKoli {
			totalKoli += v2.Quantity
		}

		// get picking order assign's sales order information
		var salesOrder *bridgeService.GetSalesOrderGPListResponse
		if salesOrder, err = s.opt.Client.BridgeServiceGrpc.GetSalesOrderListGPByID(ctx, &bridgeService.GetSalesOrderGPListByIDRequest{
			Id: v.SopNumber,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "sales order")
			return
		}

		// filtering by customer name
		if !(req.Custname != " " && strings.Contains(strings.ToLower(salesOrder.Data[0].Customer[0].Custname), strings.ToLower(req.Custname))) {
			continue
		}

		// get wrt
		var wrt *bridgeService.GetWrtGPResponse
		if wrt, err = s.opt.Client.BridgeServiceGrpc.GetWrtGPList(ctx, &bridgeService.GetWrtGPListRequest{
			Search:    salesOrder.Data[0].Wrt[0].GnL_WRT_ID,
			GnlRegion: salesOrder.Data[0].GnL_Region,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "wrt")
			return
		}

		// get picker name
		var pickingOrder *model.PickingOrder
		if pickingOrder, err = s.RepositoryPickingOrder.GetDetail(ctx, v.PickingOrderId, ""); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("picking order")
			return
		}

		var picker *bridgeService.GetHelperGPResponse
		if picker, err = s.opt.Client.BridgeServiceGrpc.GetHelperGPDetail(ctx, &bridgeService.GetHelperGPDetailRequest{
			Id: pickingOrder.PickerId,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "helper")
			return
		}

		// get DO & SI
		var do *bridgeService.GetDeliveryOrderGPListResponse
		var si *bridgeService.GetSalesInvoiceGPListResponse
		var count_print_si, count_print_do int32

		if do, err = s.opt.Client.BridgeServiceGrpc.GetDeliveryOrderListGP(ctx, &bridgeService.GetDeliveryOrderGPListRequest{
			SopNumbe: v.SopNumber,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "delivery order")
			return
		}

		if si, err = s.opt.Client.BridgeServiceGrpc.GetSalesInvoiceGPList(ctx, &bridgeService.GetSalesInvoiceGPListRequest{
			SoNumber: v.SopNumber,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "sales invoice")
			return
		}

		if len(do.Data) > 0 {
			count_print_do = do.Data[0].DataAttachment.PrintCount
		}

		if len(si.Data) > 0 {
			count_print_si = si.Data[0].DataAttachment.PrintCount
		}

		salesOrderResponse := &pb.SalesOrderToCheck{
			SopNumber:           v.SopNumber,
			MerchantName:        salesOrder.Data[0].Customer[0].Custname,
			DeliveryDate:        salesOrder.Data[0].ReqShipDate,
			Wrt:                 wrt.Data[0].Strttime + "-" + wrt.Data[0].Endtime,
			SopNote:             salesOrder.Data[0].Commntid,
			TotalItemOnProgress: totalItemOnProgress,
			TotalItem:           totalItem,
			TotalKoli:           totalKoli,
			CheckerName:         "", // filled below if exist
			PickerName:          picker.Data[0].GnlHelperName,
			Status:              int32(v.Status),
			CountPrintDo:        count_print_do,
			CountPrintSi:        count_print_si,
		}

		// get checker name
		var checker *bridgeService.GetHelperGPResponse
		if v.CheckerIdGp != "" {
			if checker, err = s.opt.Client.BridgeServiceGrpc.GetHelperGPDetail(ctx, &bridgeService.GetHelperGPDetailRequest{
				Id: v.CheckerIdGp,
			}); err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "helper")
				return
			}

			salesOrderResponse.CheckerName = checker.Data[0].GnlHelperName
		}

		response = append(response, salesOrderResponse)
	}

	res = &pb.CheckerHistoryResponse{
		Data: response,
	}

	return
}

func (s *PickingOrderService) CheckerHistoryDetail(ctx context.Context, req *pb.CheckerHistoryDetailRequest) (res *pb.CheckerHistoryDetailResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PickingOrderService.CheckerHistoryDetail")
	defer span.End()

	// get sales order inside the picking order
	var pickingOrderAssign *model.PickingOrderAssign
	if pickingOrderAssign, err = s.RepositoryPickingOrderAssign.GetByID(ctx, 0, req.SopNumber); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order assign")
		return
	}

	// if status not finished
	if pickingOrderAssign.Status != 2 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorMustStatus("picking order assign", "finished")
		return
	}

	if pickingOrderAssign.CheckerIdGp != req.CheckerId {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("Authorization", "This is not your job")
		return
	}

	var pickingOrderItem []*model.PickingOrderItem
	if pickingOrderItem, _, err = s.RepositoryPickingOrderItem.Get(ctx, &dto.PickingOrderItemGetRequest{
		PickingOrderAssignId: []int64{pickingOrderAssign.ID},
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order item")
		return
	}

	var (
		itemResponse        []*pb.PickingOrderItem
		totalItemOnProgress int64
	)
	for _, v := range pickingOrderItem {
		var productDetail *bridgeService.GetItemGPResponse
		if productDetail, err = s.opt.Client.BridgeServiceGrpc.GetItemGPDetail(ctx, &bridgeService.GetItemGPDetailRequest{
			Id: v.ItemNumber,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("inventory", "product detail")
			return
		}

		var productImage *catalog_service.GetItemDetailResponse
		if productImage, err = s.opt.Client.CatalogServiceGrpc.GetItemDetail(ctx, &catalog_service.GetItemDetailRequest{
			Id: v.ItemNumber,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("inventory", "product image")
			return
		}

		item := &pb.PickingOrderItem{
			Id:                   v.Id,
			PickingOrderAssignId: v.PickingOrderAssignId,
			ItemNumber:           v.ItemNumber,
			ItemName:             productDetail.Data[0].Itemdesc,
			Picture:              "", // filled below if exist
			OrderQty:             v.OrderQuantity,
			PickQty:              v.PickQuantity,
			CheckQty:             v.CheckQuantity,
			ExcessQty:            v.ExcessQuantity,
			UnfulfillNote:        v.UnfulfillNote,
			Uom:                  productDetail.Data[0].Uomschdl,
			Status:               int32(v.Status),
		}

		if len(productImage.Data.ItemImage) > 0 {
			item.Picture = productImage.Data.ItemImage[0].ImageUrl
		}

		itemResponse = append(itemResponse, item)

		if item.Status != 16 {
			totalItemOnProgress++
		}
	}

	// count delivery koli
	var (
		deliveryKoli []*model.DeliveryKoli
		totalKoli    float64
	)
	if deliveryKoli, _, err = s.RepositoryDeliveryKoli.Get(ctx, &dto.DeliveryKoliGetRequest{
		SopNumber: pickingOrderAssign.SopNumber,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("delivery koli")
		return
	}
	for _, v2 := range deliveryKoli {
		totalKoli += v2.Quantity
	}

	// get picking order assign's sales order information
	var salesOrder *bridgeService.GetSalesOrderGPListResponse
	if salesOrder, err = s.opt.Client.BridgeServiceGrpc.GetSalesOrderListGPByID(ctx, &bridgeService.GetSalesOrderGPListByIDRequest{
		Id: pickingOrderAssign.SopNumber,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "sales order")
		return
	}

	// get wrt
	var wrt *bridgeService.GetWrtGPResponse
	if wrt, err = s.opt.Client.BridgeServiceGrpc.GetWrtGPList(ctx, &bridgeService.GetWrtGPListRequest{
		Search:    salesOrder.Data[0].Wrt[0].GnL_WRT_ID,
		GnlRegion: salesOrder.Data[0].GnL_Region,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "wrt")
		return
	}

	// get picker name
	var pickingOrder *model.PickingOrder
	if pickingOrder, err = s.RepositoryPickingOrder.GetDetail(ctx, pickingOrderAssign.PickingOrderId, ""); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order")
		return
	}

	var picker *bridgeService.GetHelperGPResponse
	if picker, err = s.opt.Client.BridgeServiceGrpc.GetHelperGPDetail(ctx, &bridgeService.GetHelperGPDetailRequest{
		Id: pickingOrder.PickerId,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "helper")
		return
	}

	res = &pb.CheckerHistoryDetailResponse{
		SopNumber:           pickingOrderAssign.SopNumber,
		DeliveryDate:        salesOrder.Data[0].ReqShipDate,
		MerchantName:        salesOrder.Data[0].Customer[0].Custname,
		SopNote:             salesOrder.Data[0].Commntid,
		Wrt:                 wrt.Data[0].Strttime + "-" + wrt.Data[0].Endtime,
		PickerName:          picker.Data[0].GnlHelperName,
		TotalKoli:           totalKoli,
		TotalItemOnProgress: totalItemOnProgress,
		TotalItem:           int64(len(itemResponse)),
		Item:                itemResponse,
		Status:              int32(pickingOrderAssign.Status),
	}

	return
}

func (s *PickingOrderService) CheckerWidget(ctx context.Context, req *pb.CheckerWidgetRequest) (res *pb.CheckerWidgetResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PickingOrderService.CheckerWidget")
	defer span.End()

	// get total Sales Order
	var totalSalesOrder int64
	if _, totalSalesOrder, err = s.RepositoryPickingOrderAssign.Get(ctx, &dto.PickingOrderAssignGetRequest{
		DeliveryDateFrom: time.Now().Add(-time.Hour * 48),
		DeliveryDateTo:   time.Now().Add(time.Hour * 24),
		Status:           []int{20, 2, 16},
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order assign")
		return
	}

	// get total picked
	var totalPicked int64
	if _, totalPicked, err = s.RepositoryPickingOrderAssign.Get(ctx, &dto.PickingOrderAssignGetRequest{
		DeliveryDateFrom: time.Now().Add(-time.Hour * 48),
		DeliveryDateTo:   time.Now().Add(time.Hour * 24),
		Status:           []int{16},
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order assign")
		return
	}

	// get total checking checker
	var totalChecking int64
	if _, totalChecking, err = s.RepositoryPickingOrderAssign.Get(ctx, &dto.PickingOrderAssignGetRequest{
		DeliveryDateFrom: time.Now().Add(-time.Hour * 48),
		DeliveryDateTo:   time.Now().Add(time.Hour * 24),
		CheckerId:        []string{req.CheckerId},
		Status:           []int{20},
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order assign")
		return
	}

	// get finished sales order checker
	var totalFinished int64
	if _, totalFinished, err = s.RepositoryPickingOrderAssign.Get(ctx, &dto.PickingOrderAssignGetRequest{
		DeliveryDateFrom: time.Now().Add(-time.Hour * 48),
		DeliveryDateTo:   time.Now().Add(time.Hour * 24),
		CheckerId:        []string{req.CheckerId},
		Status:           []int{2},
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order assign")
		return
	}

	res = &pb.CheckerWidgetResponse{
		TotalSalesOrder: totalSalesOrder,
		TotalPicked:     totalPicked,
		TotalFinished:   totalFinished,
		TotalChecking:   totalChecking,
	}

	if res.TotalSalesOrder != 0 {
		res.TotalFinishedPercentage = math.Round(float64(res.TotalFinished)/float64(res.TotalSalesOrder)*100) / 100
		res.TotalCheckingPercentage = math.Round(float64(res.TotalChecking)/float64(res.TotalSalesOrder)*100) / 100
	}

	return
}
