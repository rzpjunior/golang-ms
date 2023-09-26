package service

import (
	"context"
	"fmt"
	"strconv"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/utils"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/catalog_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/sales_service"
)

type ITransactionHistoryService interface {
	Get(ctx context.Context, req *dto.RequestGetHistoryTransaction) (res dto.DataGetHistoryTransaction, total int64, err error)
	GetDetail(ctx context.Context, req *dto.RequestGetDetailSO) (res dto.SalesOrderDetailResponse, err error)
	GetInvoiceDetail(ctx context.Context, req *dto.RequestGetInvoiceDetail) (res *dto.SalesInvoiceDetailResponse, err error)
}

type TransactionHistoryService struct {
	opt opt.Options
}

func NewTransactionHistoryService() ITransactionHistoryService {
	return &TransactionHistoryService{
		opt: global.Setup.Common,
	}
}

func (s *TransactionHistoryService) Get(ctx context.Context, req *dto.RequestGetHistoryTransaction) (res dto.DataGetHistoryTransaction, total int64, err error) {
	// var total int64
	fmt.Println(total)
	ctx, span := s.opt.Trace.Start(ctx, "TransactionHistoryService.Get")
	defer span.End()
	customerID, _ := strconv.Atoi(req.Session.Customer.ID)
	if req.Data.AddressID != req.Session.Address.ID {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	addressID, _ := strconv.Atoi(req.Data.AddressID)
	var Type, category int
	if req.Data.Type != "" {
		Type, _ = strconv.Atoi(req.Data.Type)
	} else {
		Type = 0
	}
	if req.Data.Category != "" {
		category, _ = strconv.Atoi(req.Data.Category)
	} else {
		category = 0
	}
	so, err := s.opt.Client.SalesServiceGrpc.GetSalesOrderListMobile(ctx, &sales_service.GetSalesOrderListRequest{
		CustomerId:   int64(customerID),
		AddressId:    int64(addressID),
		AddressCode:  req.Data.AddressID,
		Type:         int64(Type),
		Category:     int64(category),
		CustomerCode: req.Session.Customer.Code,
		Limit:        int32(req.Limit),
		Offset:       int32(req.Offset),
		//OrderBy:    req.Data.OrderBy,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	layout := "2006-01-02"
	//loc, _ := time.LoadLocation("Asia/Jakarta")

	for _, v := range so.Data {
		// recognitionDate, _ := time.ParseInLocation(layout, v.OrderDate.AsTime().Format(layout),loc)
		// deliveryDate, _ := time.ParseInLocation(layout, v.RequestDeliveryDate.AsTime().Format(layout), loc)
		// rd := recognitionDate.Format(layout)
		// dd := deliveryDate.Format(layout)

		res.SalesOrder = append(res.SalesOrder, &dto.ListSalesOrder{
			ID:                strconv.Itoa(int(v.Id)),
			OrderCode:         v.SalesOrderNumber,
			OrderDate:         v.OrderDate.AsTime().Format(layout),
			OrderDeliveryDate: v.RequestDeliveryDate.AsTime().Format(layout),
			OrderStatus:       strconv.Itoa(int(v.Status)),
			TotalCharge:       strconv.Itoa(int(v.TotalCharge)),
			OrderTypeSlsID:    strconv.Itoa(int(v.OrderTypeSlsId)),
			TermPaymentSlsID:  strconv.Itoa(int(v.TermPaymentSlsId)),
		})
	}
	res.AddressID = req.Session.Address.ID
	res.Type = req.Data.Type
	res.Category = req.Data.Category
	total = int64(so.TotalRecords)
	return
}

func (s *TransactionHistoryService) GetDetail(ctx context.Context, req *dto.RequestGetDetailSO) (res dto.SalesOrderDetailResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "TransactionHistoryService.GetDetail")
	defer span.End()
	customerID, _ := strconv.Atoi(req.Session.Customer.ID)
	if req.Data.AddressID != req.Session.Address.ID {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	//addressID, _ := strconv.Atoi(req.Data.AddressID)

	// soID, _ := strconv.Atoi(req.Data.SalesOrderID)
	so, err := s.opt.Client.SalesServiceGrpc.GetSalesOrderDetail(ctx, &sales_service.GetSalesOrderDetailRequest{
		// Id: int64(soID),
		Code:         req.Data.SalesOrderID,
		CustomerId:   int64(customerID),
		CustomerIdGp: req.Session.Customer.Code,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("sales", "so detail")
		return
	}
	// if so.Data.CustomerIdGp != req.Session.Customer.Code {

	// }
	// soi, err := s.opt.Client.SalesServiceGrpc.GetSalesOrderItemList(ctx, &sales_service.GetSalesOrderItemListRequest{
	// 	SalesOrderId: int64(soID),
	// })
	// if err != nil {
	// 	span.RecordError(err)
	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 	return
	// }

	// admDivisionID, err := strconv.Atoi(req.Session.Address.AdmDivisionId)
	// admDivision, err := s.opt.Client.BridgeServiceGrpc.GetAdmDivisionDetail(ctx, &bridge_service.GetAdmDivisionDetailRequest{
	// 	Id: int64(admDivisionID),
	// })
	address, err := s.opt.Client.BridgeServiceGrpc.GetAddressGPDetail(ctx, &bridge_service.GetAddressGPDetailRequest{
		Id: so.Data.AddressIdGp,
	})
	admDivision, err := s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPList(ctx, &bridge_service.GetAdmDivisionGPListRequest{
		Limit:           1,
		Offset:          0,
		AdmDivisionCode: address.Data[0].GnL_Administrative_Code,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "adm division detail")
		return
	}
	Wrt, err := s.opt.Client.ConfigurationServiceGrpc.GetWrtList(ctx, &configuration_service.GetWrtListRequest{
		RegionId: admDivision.Data[0].Region,
		// Type:     int32(Type),
		Limit:  10,
		Offset: 1,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("configuration", "wrt list")
		return
	}
	layout := "2006-01-02"
	//loc, _ := time.LoadLocation("Asia/Jakarta")
	// recognitionDate, _ := time.ParseInLocation(layout, so.Data.RecognitionDate.AsTime().Format(layout), loc)
	// deliveryDate, _ := time.ParseInLocation(layout, so.Data.RequestDeliveryDate.AsTime().Format(layout), loc)
	// rd := recognitionDate.Format(layout)
	// dd := deliveryDate.Format(layout)
	ca := so.Data.RecognitionDate.AsTime().Format("2006-01-02 15:04:05")
	// for _, v := range soi.Data {

	// 	res.SalesOrderItems = append(res.SalesOrderItems, &dto.SalesOrderItemResponse{
	// 		ID: strconv.Itoa(int(v.Id)),
	// 		//SalesOrderItemID:   strconv.Itoa(int(v.Id)),
	// 		ProductID:          strconv.Itoa(int(v.ItemId)),
	// 		OrderQty:           strconv.Itoa(int(v.OrderQty)),
	// 		UnitPrice:          strconv.Itoa(int(v.UnitPrice)),
	// 		ShadowPrice:        strconv.Itoa(int(v.ShadowPrice)),
	// 		Subtotal:           strconv.Itoa(int(v.Subtotal)),
	// 		Weight:             strconv.Itoa(int(v.Weight)),
	// 		Note:               v.Note,
	// 		UomName:            uom.Data.Description,
	// 		Name:               item.Data.Description,
	// 		ImageUrl:           itemImage.Data.ImageUrl,
	// 		DiscountQty:        strconv.Itoa(int(v.DiscountQty)),
	// 		UnitPriceDiscount:  strconv.Itoa(int(v.UnitPriceDiscount)),
	// 		ItemDiscountAmount: strconv.Itoa(int(v.ItemDiscountAmount)),
	// 	})
	// }
	res = dto.SalesOrderDetailResponse{
		ID:                  strconv.Itoa(int(so.Data.Id)),
		Code:                so.Data.SalesOrderNumber,
		RecognitionDate:     so.Data.RecognitionDate.AsTime().Format(layout),
		DeliveryDate:        so.Data.RequestDeliveryDate.AsTime().Format(layout),
		ShippingAddress:     address.Data[0].AddresS1 + " " + address.Data[0].AddresS2 + " " + address.Data[0].AddresS3,
		ShippingAddressNote: so.Data.ShippingAddressNote,
		DeliveryFee:         strconv.FormatFloat(so.Data.DeliveryFee, 'f', 1, 64),
		VouDiscAmount:       strconv.FormatFloat(so.Data.VouDiscAmount, 'f', 1, 64),
		// VoucherType:            so.Data.vou,
		PointRedeemAmount: strconv.FormatFloat(so.Data.PointRedeemAmount, 'f', 1, 64),
		TotalPrice:        strconv.FormatFloat(so.Data.TotalPrice, 'f', 1, 64),
		TotalCharge:       strconv.FormatFloat(so.Data.TotalCharge, 'f', 1, 64),
		Note:              so.Data.Note,
		Status:            strconv.Itoa(int(so.Data.Status)),
		//SalesOrderItems:   []*dto.SalesOrderItemResponse{},
		WrtName:     Wrt.Data[0].Name,
		PhoneNumber: so.Data.Phone,
		CityName:    admDivision.Data[0].City,
		PicName:     so.Data.PicName,
		CreatedAt:   ca,
		AddressName: so.Data.AddressName,
	}

	for _, v := range so.Data.DataSalesOrderItem {
		item, err := s.opt.Client.CatalogServiceGrpc.GetItemDetail(ctx, &catalog_service.GetItemDetailRequest{
			Id: v.ItemIdGp,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			break
		}
		itemImage, err := s.opt.Client.CatalogServiceGrpc.GetItemImageDetail(ctx, &catalog_service.GetItemImageDetailRequest{
			ItemId:    v.ItemId,
			MainImage: 1,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			break
		}
		uom, err := s.opt.Client.BridgeServiceGrpc.GetUomDetail(ctx, &bridge_service.GetUomDetailRequest{
			Id: utils.ToInt64(item.Data.UomId),
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			break
		}
		res.SalesOrderItems = append(res.SalesOrderItems, &dto.SalesOrderItemResponse{
			ID:                 strconv.Itoa(int(v.Id)),
			ProductID:          strconv.Itoa(int(v.ItemId)),
			OrderQty:           strconv.FormatFloat(v.OrderQty, 'f', 1, 64),
			UnitPrice:          strconv.FormatFloat(v.UnitPrice, 'f', 1, 64),
			ShadowPrice:        strconv.FormatFloat(v.ShadowPrice, 'f', 1, 64),
			Subtotal:           strconv.FormatFloat(v.Subtotal, 'f', 1, 64),
			Weight:             strconv.FormatFloat(v.Weight, 'f', 1, 64),
			Note:               v.Note,
			UomName:            uom.Data.Description,
			Name:               item.Data.Description,
			ImageUrl:           itemImage.Data.ImageUrl,
			DiscountQty:        strconv.FormatFloat(v.DiscountQty, 'f', 1, 64),
			UnitPriceDiscount:  strconv.FormatFloat(v.UnitPriceDiscount, 'f', 1, 64),
			ItemDiscountAmount: strconv.FormatFloat(v.ItemDiscountAmount, 'f', 1, 64),
		})
	}
	// res.AddressID = req.Session.Address.ID
	// res.Type = req.Data.Type
	// res.Category = req.Data.Category

	return
}

func (s *TransactionHistoryService) GetInvoiceDetail(ctx context.Context, req *dto.RequestGetInvoiceDetail) (res *dto.SalesInvoiceDetailResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "TransactionHistoryService.GetInvoiceDetail")
	defer span.End()

	var (
		salesInvoiceDetail   *dto.SalesInvoice
		salesInvoiceItemList []*dto.SalesInvoiceItem
		salesPaymentList     []*dto.SalesPayment
	)

	// salesOrder, err := s.opt.Client.BridgeServiceGrpc.GetSalesOrderDetail(ctx, &bridge_service.GetSalesOrderDetailRequest{
	// 	Id: utils.ToInt64(req.Data.SalesOrderID),
	// })
	salesOrder, err := s.opt.Client.BridgeServiceGrpc.GetSalesOrderListGPAll(ctx, &bridge_service.GetSalesOrderGPListRequest{
		// Id: utils.ToInt64(req.Data.SalesOrderID),
		Limit:      1,
		Offset:     0,
		Custnumber: req.Session.Customer.Code,
		SoNumber:   req.Data.SalesOrderID,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// salesInvoice, err := s.opt.Client.BridgeServiceGrpc.GetSalesInvoiceDetail(ctx, &bridge_service.GetSalesInvoiceDetailRequest{
	// 	Id: utils.ToInt64(req.Data.SalesInvoiceID),
	// })
	salesInvoice, err := s.opt.Client.SalesServiceGrpc.GetSalesInvoiceGPMobileList(ctx, &sales_service.GetSalesInvoiceGPMobileListRequest{
		Limit:      1,
		Offset:     0,
		Custnumber: req.Session.Customer.Code,
		SoNumber:   req.Data.SalesOrderID,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	salesInvoiceDetail = &dto.SalesInvoice{
		ID:                salesInvoice.Data[0].SalesInvoice.Id,
		OrderCode:         salesOrder.Data[0].Sopnumbe,
		InvoiceCode:       salesInvoice.Data[0].SalesInvoice.InvoiceCode,
		OrderDate:         salesOrder.Data[0].Docdate,
		InvoiceDate:       salesInvoice.Data[0].SalesInvoice.InvoiceDate,
		TotalPrice:        salesInvoice.Data[0].SalesInvoice.TotalPrice,
		DeliveryFee:       salesInvoice.Data[0].SalesInvoice.DeliveryFee,
		VoucherAmount:     salesInvoice.Data[0].SalesInvoice.VoucherAmount,
		PointRedeemAmount: salesInvoice.Data[0].SalesInvoice.PointRedeemAmount,
		AdjustmentAmount:  salesInvoice.Data[0].SalesInvoice.AdjustmentAmount,
		TotalCharge:       salesInvoice.Data[0].SalesInvoice.TotalCharge,
	}

	// salesInvoiceItem, err := s.opt.Client.BridgeServiceGrpc.GetSalesInvoiceItemList(ctx, &bridge_service.GetSalesInvoiceItemListRequest{
	// 	SalesInvoiceId: utils.ToInt64(req.Data.SalesInvoiceID),
	// })
	// if err != nil {
	// 	span.RecordError(err)
	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 	return
	// }

	// for i, v := range salesInvoiceItem.Data {
	// 	salesInvoiceItemList = append(salesInvoiceItemList, &dto.SalesInvoiceItem{
	// 		ItemID:            utils.ToString(v.ItemId),
	// 		ItemName:          fmt.Sprintf("Dummy Item %d", i+1),
	// 		InvoiceQty:        v.InvoiceQty,
	// 		UomName:           fmt.Sprintf("Dummy Uom %d", i+1),
	// 		UnitPrice:         v.UnitPrice,
	// 		Subtotal:          v.Subtotal,
	// 		SkuDiscountAmount: v.SkuDiscAmount,
	// 	})
	// 	if i == 5 {
	// 		break
	// 	}
	// }
	for i, v := range salesInvoice.Data[0].InvoiceItem {
		salesInvoiceItemList = append(salesInvoiceItemList, &dto.SalesInvoiceItem{
			ItemID:     v.ItemId,
			ItemName:   v.ItemName,
			InvoiceQty: v.InvoiceQty,
			UomName:    v.UomName,
			UnitPrice:  v.UnitPrice,
			Subtotal:   v.Subtotal,
			// SkuDiscountAmount: v.SkuDiscAmount,
		})
		if i == 5 {
			break
		}
	}

	salesPayment, err := s.opt.Client.BridgeServiceGrpc.GetCashReceiptList(ctx, &bridge_service.GetCashReceiptListRequest{
		Limit:    1,
		Offset:   0,
		Custnmbr: req.Session.Customer.Code,
		Sopnumbe: salesInvoice.Data[0].SalesInvoice.Id,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	for i, v := range salesPayment.Data {
		salesPaymentList = append(salesPaymentList, &dto.SalesPayment{
			ID:             v.Docnumbr,
			Code:           v.Docnumbr,
			PaymentDate:    v.Docdate,
			PaymentTime:    v.Docdate,
			Amount:         utils.ToString(v.Ortrxamt),
			PaymentChannel: v.PaymentMethod,
			Status:         int8(v.Dcstatus),
		})
		if i == 2 {
			break
		}
	}

	res = &dto.SalesInvoiceDetailResponse{
		InvoiceDetail:  salesInvoiceDetail,
		InvoiceItem:    salesInvoiceItemList,
		InvoicePayment: salesPaymentList,
	}

	return
}
