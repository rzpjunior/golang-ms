package service

import (
	"context"
	"fmt"
	"math"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/constants"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/timex"
	"git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/internal/app/dto"
	auditService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/audit_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/catalog_service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewServiceSalesInvoice() ISalesInvoiceService {
	m := new(SalesInvoiceService)
	m.opt = global.Setup.Common
	return m
}

type ISalesInvoiceService interface {
	Get(ctx context.Context, req dto.SalesInvoiceListRequest) (res []*dto.SalesInvoiceResponse, total int64, err error)
	GetByID(ctx context.Context, req dto.SalesInvoiceDetailRequest) (res *dto.SalesInvoiceResponse, err error)
	CreateGP(ctx context.Context, req dto.CreateSalesInvoiceRequest) (res *dto.SalesInvoiceResponse, err error)
	GetListGP(ctx context.Context, req dto.GetSalesInvoiceGPRequest) (res []*bridgeService.SalesInvoiceGP, total int64, err error)
	GetOrderPerformance(ctx context.Context, req dto.SalesInvoiceListRequest) (res []*dto.SalesInvoiceResponse, summaryOrderPerformance *dto.OrderPerformance, err error)
}

type SalesInvoiceService struct {
	opt opt.Options
}

func (s *SalesInvoiceService) Get(ctx context.Context, req dto.SalesInvoiceListRequest) (res []*dto.SalesInvoiceResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesInvoiceService.Get")
	defer span.End()

	// get sales invoice from bridge
	var siRes *bridgeService.GetSalesInvoiceGPListResponse
	var recognitionDateFrom, recognitionDateTo string
	if timex.IsValid(req.RecognitionDateFrom) {
		recognitionDateFrom = req.RecognitionDateFrom.Format(timex.InFormatDate)
	}

	if timex.IsValid(req.RecognitionDateTo) {
		recognitionDateTo = req.RecognitionDateTo.Format(timex.InFormatDate)
	}

	siRes, err = s.opt.Client.BridgeServiceGrpc.GetSalesInvoiceGPList(ctx, &bridgeService.GetSalesInvoiceGPListRequest{
		Limit:  req.Limit,
		Offset: req.Offset,
		// Status:        req.Status,
		SiNumber:            req.Search,
		DocdateFrom:         recognitionDateFrom,
		DocdateTo:           recognitionDateTo,
		OrderBy:             "desc",
		Locncode:            req.SiteID,
		Custnumber:          req.CustomerID,
		GnlCustTypeId:       "BTY0015",
		RemainingAmountFlag: req.Status,
		// OrderBy: req.OrderBy,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "sales invoice")
		return
	}

	datas := []*dto.SalesInvoiceResponse{}
	for _, si := range siRes.Data {
		var docDate time.Time
		docDate, err = time.Parse("2006-01-02T15:04:05", si.Docdate)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("doc_date")
			return
		}
		var dueDate time.Time
		dueDate, err = time.Parse("2006-01-02T15:04:05", si.Duedate)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("due_date")
			return
		}
		var custType, custTypeDesc string
		var customerDetail *bridgeService.GetCustomerGPResponse
		if si.SalesOrder[0].Custnmbr != "" {
			customerDetail, _ = s.opt.Client.BridgeServiceGrpc.GetCustomerGPList(ctx, &bridgeService.GetCustomerGPListRequest{
				Id:             si.SalesOrder[0].Custnmbr,
				Limit:          1,
				Offset:         0,
				CustomerTypeId: "BTY0015",
				Inactive:       "0",
			})
			if len(customerDetail.Data) == 0 {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "sales invoice")
				return
			}
			// cust,_ =s.opt.Client.BridgeServiceGrpc.GetCustomerGPDetail(ctx, &bridgeService.GetCustomerGPDetailRequest{}
			if len(customerDetail.Data[0].CustomerType) != 0 {
				custType = customerDetail.Data[0].CustomerType[0].GnL_Cust_Type_ID
				custTypeDesc = customerDetail.Data[0].CustomerType[0].GnL_CustType_Description
			}

		} else {
			customerDetail = &bridgeService.GetCustomerGPResponse{}
			customerDetail.Data = append(customerDetail.Data, &bridgeService.CustomerGP{})
		}
		var tempSI *dto.SalesInvoiceResponse
		statusDoc := 1
		if si.RemainingAmount == 0 {
			statusDoc = 2
		}
		tempSI = &dto.SalesInvoiceResponse{
			// ID:                si.Id,
			Code:             si.Sopnumbe,
			CodeExt:          si.Sopnumbe,
			Status:           int8(statusDoc),
			RecognitionDate:  docDate,
			DueDate:          dueDate,
			BillingAddress:   si.Address,
			DeliveryFee:      si.Frtamnt,
			CustomerTypeID:   custType,
			CustomerTypeDesc: custTypeDesc,
			// VouRedeemCode:     si.VouRedeemCode,
			// VouDiscAmount:     si.VouDiscAmount,
			// PointRedeemAmount: si.PointRedeemAmount,
			// Adjustment:        int8(si.Adjustment),
			// AdjAmount:         si.AdjAmount,
			// AdjNote:           si.AdjNote,
			TotalPrice:  si.Subtotal,
			TotalCharge: si.Ordocamt,
			// DeltaPrint:        si.DeltaPrint,
			// VoucherID:         si.VoucherId,
			RemainingAmount: si.RemainingAmount,
			CustomerName:    customerDetail.Data[0].Custname,
			RegionID:        si.SalesOrder[0].GnlRegion,
			// Note:              si.Note,
			SalesOrder: &dto.SalesOrderResponse{
				// Code:       si.SalesOrder[0].Orignumb,
				CustomerID: customerDetail.Data[0].Custnmbr,
				AddressGP: &dto.AddressResponse{
					CustomerName: customerDetail.Data[0].Custname,
					Customer: &dto.CustomerResponse{
						Name: customerDetail.Data[0].Custname,
					},
				},
			},
			SiteID: si.Locncode,
		}

		for _, v := range si.Details {

			tempSI.SalesInvoiceItem = append(tempSI.SalesInvoiceItem, &dto.SalesInvoiceItemResponse{
				SalesInvoiceID: tempSI.Code,
				ItemID:         v.Itemnmbr,
				InvoiceQty:     v.Quantity,
				UnitPrice:      v.Unitprce,
				Subtotal:       v.Xtndprce,
				UomName:        v.Uofm,
			})
		}

		datas = append(datas, tempSI)

	}
	res = datas
	total = int64(siRes.TotalRecords)

	return
}

func (s *SalesInvoiceService) GetByID(ctx context.Context, req dto.SalesInvoiceDetailRequest) (res *dto.SalesInvoiceResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesInvoiceService.GetByID")
	defer span.End()

	// get SalesInvoice from bridge
	var siRes *bridgeService.GetSalesInvoiceGPDetailResponse
	siRes, err = s.opt.Client.BridgeServiceGrpc.GetSalesInvoiceGPDetail(ctx, &bridgeService.GetSalesInvoiceGPDetailRequest{
		Id: req.Id,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "sales invoice")
		return
	}

	var sii []*dto.SalesInvoiceItemResponse

	for _, si := range siRes.Data[0].Details {
		var item *bridgeService.GetItemGPResponse
		item, err = s.opt.Client.BridgeServiceGrpc.GetItemGPDetail(ctx, &bridgeService.GetItemGPDetailRequest{
			Id: si.Itemnmbr,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "sales invoice")
			return
		}
		sii = append(sii, &dto.SalesInvoiceItemResponse{
			UnitPrice:  si.Unitprce,
			InvoiceQty: si.Quantity,
			Subtotal:   si.Xtndprce,
			ItemID:     item.Data[0].Itemnmbr,
			ItemName:   item.Data[0].Itemdesc,
			// Item: &dto.ItemResponse{
			// 	// ID:          item.Data.Id,
			// 	Code:        item.Data[0].Itemnmbr,
			// 	Description: item.Data[0].Itemdesc,
			// },
			SalesInvoiceID: siRes.Data[0].Sopnumbe,
		})
	}

	// get CustomerGP from bridge
	var custGP *bridgeService.GetCustomerGPResponse
	custGP, err = s.opt.Client.BridgeServiceGrpc.GetCustomerGPList(ctx, &bridgeService.GetCustomerGPListRequest{
		Id:             siRes.Data[0].Custnmbr,
		Limit:          1,
		Offset:         0,
		CustomerTypeId: "BTY0015",
		Inactive:       "0",
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "address gp")
		return
	}

	//validasi edn
	if len(custGP.Data[0].CustomerType) > 0 {
		if custGP.Data[0].CustomerType[0].GnL_Cust_Type_ID != "BTY0015" {
			err = edenlabs.ErrorValidation("customer_type", "customer is not edn customer")
			return
		}
	} else {
		err = edenlabs.ErrorValidation("customer_type", "customer is not edn customer")
		return
	}
	// get Site from bridge
	var siteRes *bridgeService.GetSiteGPResponse
	siteRes, err = s.opt.Client.BridgeServiceGrpc.GetSiteGPDetail(ctx, &bridgeService.GetSiteGPDetailRequest{
		Id: siRes.Data[0].Locncode,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "site")
		return
	}

	var docDate time.Time
	docDate, err = time.Parse("2006-01-02T15:04:05", siRes.Data[0].Docdate)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("doc_date")
		return
	}
	var dueDate time.Time
	dueDate, err = time.Parse("2006-01-02T15:04:05", siRes.Data[0].Duedate)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("due_date")
		return
	}
	var totalPaid float64
	var sp []*dto.SalesPayment

	// for _, v := range siRes.Data[0].CashReceipt {
	// 	totalPaid += v.TotalPaid
	// 	var tempSP *bridgeService.GetSalesPaymentGPResponse

	// 	tempSP, err = s.opt.Client.BridgeServiceGrpc.GetSalesPaymentGPDetail(ctx, &bridgeService.GetSalesPaymentGPDetailRequest{
	// 		Id: v.Docnumbr,
	// 	})
	// 	if err != nil {
	// 		span.RecordError(err)
	// 		s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 		err = edenlabs.ErrorRpcNotFound("bridge", "sales payment")
	// 		return
	// 	}
	// 	var docDate time.Time
	// 	if tempSP.Data[0].Docdate != "" {
	// 		docDate, err = time.Parse("2006-01-02T15:04:05", tempSP.Data[0].Docdate)
	// 		if err != nil {
	// 			span.RecordError(err)
	// 			s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 			err = edenlabs.ErrorInvalid("doc_date")
	// 			return
	// 		}
	// 	}
	// 	var createdDate time.Time
	// 	if tempSP.Data[0].Creatddt != "" {
	// 		createdDate, err = time.Parse("2006-01-02T15:04:05", tempSP.Data[0].Creatddt)
	// 		if err != nil {
	// 			span.RecordError(err)
	// 			s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 			err = edenlabs.ErrorInvalid("created_date")
	// 			return
	// 		}
	// 	}

	// 	sp = append(sp, &dto.SalesPayment{
	// 		Code:            v.Docnumbr,
	// 		Amount:          v.TotalPaid,
	// 		Status:          int8(tempSP.Data[0].Dcstatus),
	// 		RecognitionDate: docDate,
	// 		CreatedAt:       createdDate,
	// 	})
	// }
	var listSP *bridgeService.GetSalesPaymentGPResponse
	listSP, err = s.opt.Client.BridgeServiceGrpc.GetSalesPaymentGPList(ctx, &bridgeService.GetSalesPaymentGPListRequest{
		Sopnumbe: req.Id,
		Limit:    int32(len(siRes.Data[0].CashReceipt)),
	})

	for _, v := range listSP.Data {
		totalPaid += v.Ortrxamt
		var tempSP *bridgeService.GetSalesPaymentGPResponse

		tempSP, err = s.opt.Client.BridgeServiceGrpc.GetSalesPaymentGPDetail(ctx, &bridgeService.GetSalesPaymentGPDetailRequest{
			Id: v.Docnumbr,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "sales payment")
			return
		}
		var docDate time.Time
		if tempSP.Data[0].Docdate != "" {
			docDate, err = time.Parse("2006-01-02T15:04:05", tempSP.Data[0].Docdate)
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorInvalid("doc_date")
				return
			}
		}
		var createdDate time.Time
		if tempSP.Data[0].Creatddt != "" {
			createdDate, err = time.Parse("2006-01-02T15:04:05", tempSP.Data[0].Creatddt)
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorInvalid("created_date")
				return
			}
		}

		sp = append(sp, &dto.SalesPayment{
			Code:            v.Docnumbr,
			Amount:          v.Ortrxamt,
			Status:          int8(tempSP.Data[0].Dcstatus),
			RecognitionDate: docDate,
			CreatedAt:       createdDate,
		})
	}
	statusDoc := 1
	if siRes.Data[0].RemainingAmount == 0 {
		statusDoc = 2
	}
	res = &dto.SalesInvoiceResponse{
		// ID:                siRes.Data[0].Id,
		Code:            siRes.Data[0].Sopnumbe,
		CodeExt:         siRes.Data[0].Sopnumbe,
		Status:          int8(statusDoc),
		RecognitionDate: docDate,
		DueDate:         dueDate,
		BillingAddress:  siRes.Data[0].Address,
		DeliveryFee:     siRes.Data[0].Frtamnt,
		// VouRedeemCode:     siRes.Data[0].VouRedeemCode,
		// VouDiscAmount:     siRes.Data[0].VouDiscAmount,
		// PointRedeemAmount: siRes.Data[0].PointRedeemAmount,
		// Adjustment:        int8(siRes.Data[0].Adjustment),
		// AdjAmount:         siRes.Data[0].AdjAmount,
		// AdjNote:           siRes.Data[0].AdjNote,
		TotalPrice:  siRes.Data[0].Subtotal,
		TotalCharge: siRes.Data[0].Ordocamt,
		// DeltaPrint:        siRes.Data[0].DeltaPrint,
		// VoucherID:         siRes.Data[0].VoucherId,
		RemainingAmount: siRes.Data[0].RemainingAmount,
		// Note:              siRes.Data[0].Note,
		TotalPaid: totalPaid,
		SalesOrder: &dto.SalesOrderResponse{
			AddressGP: &dto.AddressResponse{
				CustomerName: custGP.Data[0].Custname,
				Customer: &dto.CustomerResponse{
					Name: custGP.Data[0].Custname,
				},
			},
			Site: &dto.SiteResponse{
				// ID:            siteRes.Data.Id,
				Code:           siteRes.Data[0].Locncode,
				Name:           siteRes.Data[0].Locndscr,
				Description:    siteRes.Data[0].Locndscr,
				Status:         int8(siteRes.Data[0].Inactive),
				PhoneNumber:    siteRes.Data[0].PhonE1,
				AltPhoneNumber: siteRes.Data[0].PhonE2,
				// StatusConvert: statusx.ConvertStatusValue(int8(siteRes.Data.Status)),
				// CreatedAt:     siteRes.Data.CreatedAt.AsTime(),
				// UpdatedAt:     siteRes.Data.UpdatedAt.AsTime(),
			},
			// Code:       siRes.Data[0].SalesOrder[0].Orignumb,
			CustomerID: siRes.Data[0].Custnmbr,
		},
		SalesInvoiceItem: sii,
		CreatedAt:        docDate,
		SalesPayment:     sp,
	}

	return
}

func (s *SalesInvoiceService) CreateGP(ctx context.Context, req dto.CreateSalesInvoiceRequest) (res *dto.SalesInvoiceResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesInvoiceService.CreateGP")
	defer span.End()

	// create Purchase Order from bridge
	var (
		vourchers []*bridgeService.CreateSalesInvoiceGPRequest_VoucherApply
		// voucherList = make(map[string]bool)
		details    []*bridgeService.CreateSalesInvoiceGPRequest_DetailItem
		detailList = make(map[string]bool)
	)
	// for _, voucher := range req.VoucherApply {
	// 	if _, exist := voucherList[voucher.GnlVoucherId]; exist {
	// 		err = edenlabs.ErrorDuplicate("voucher id")
	// 		span.RecordError(err)
	// 		s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 		return
	// 	}

	// 	vourchers = append(vourchers, &bridgeService.CreateSalesInvoiceGPRequest_VoucherApply{
	// 		GnlVoucherType: voucher.GnlVoucherType,
	// 		GnlVoucherId:   voucher.GnlVoucherId,
	// 		Ordocamt:       voucher.Ordocamt,
	// 	})
	// 	voucherList[voucher.GnlVoucherId] = true
	// }
	customer, err := s.opt.Client.BridgeServiceGrpc.GetCustomerGPList(ctx, &bridge_service.GetCustomerGPListRequest{
		// Id: customerDetail.Data.CustomerIdGp,
		Limit:          1,
		Offset:         0,
		Id:             req.CustomerID,
		CustomerTypeId: "BTY0015",
	})
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	if len(customer.Data) == 0 {
		err = edenlabs.ErrorRpcNotFound("bridge", "customer")
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	address, err := s.opt.Client.BridgeServiceGrpc.GetAddressGPList(ctx, &bridge_service.GetAddressGPListRequest{
		// Id: customer.Data[0].Adrscode[0].Adrscode,
		Limit:          100,
		Offset:         0,
		CustomerNumber: customer.Data[0].Custnmbr,
		// Adrscode:       req.AddressID,
	})

	if err != nil {
		err = edenlabs.ErrorRpcNotFound("bridge", "address")
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	if len(address.Data) == 0 {
		err = edenlabs.ErrorRpcNotFound("bridge", "address")
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	tempAddress := &bridge_service.AddressGP{}
	for _, v := range address.Data {
		if v.TypeAddress == "ship_to" {
			tempAddress = v
		}
	}
	if tempAddress == nil {
		//span.RecordError(err)
		fmt.Println(address.Data)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return nil, err
	}
	var subtotal float64
	for _, detail := range req.Products {
		if _, exist := detailList[detail.ProductID]; exist {
			err = edenlabs.ErrorDuplicate("item")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		var itemDetail *catalog_service.GetItemDetailByInternalIdResponse
		itemDetail, err = s.opt.Client.CatalogServiceGrpc.GetItemDetailMasterComplexByInternalID(ctx, &catalog_service.GetItemDetailByInternalIdRequest{
			// Id:               detail.ProductID,
			ItemIdGp:         detail.ProductID,
			RegionIdGp:       tempAddress.AdministrativeDiv.GnlRegion,
			CustomerTypeIdGp: customer.Data[0].CustomerType[0].GnL_Cust_Type_ID,
			LocationCode:     tempAddress.Locncode,
			PriceLevel:       customer.Data[0].Prclevel,
			Salability:       1,
		})
		if err != nil {
			err = edenlabs.ErrorRpcNotFound("catalog", "item")
			return
		}
		if itemDetail.Data.ItemSite[0].TotalStock <= 0 {
			err = edenlabs.ErrorValidation("qty", detail.ProductID+" stock is less than 0")
			return
		}
		qtyModulo := math.Mod(detail.Quantity, itemDetail.Data.OrderMinQty)
		if qtyModulo != 0 {
			err = edenlabs.ErrorValidation("qty", "invalid product quantity")
			return
		}

		fmt.Print(itemDetail)

		details = append(details, &bridgeService.CreateSalesInvoiceGPRequest_DetailItem{
			// Lnitmseq:   detail.Lnitmseq,
			Itemnmbr:   detail.ProductID,
			Locncode:   tempAddress.Locncode,
			Uofm:       itemDetail.Data.UomId,
			Pricelvl:   customer.Data[0].Prclevel,
			Quantity:   detail.Quantity,
			Unitprce:   float64(detail.UnitPrice),
			Xtndprce:   detail.Quantity * float64(detail.UnitPrice),
			GnL_Weight: int32(detail.Weight),
		})
		subtotal += (detail.Quantity * float64(detail.UnitPrice))
		detailList[detail.ProductID] = true
	}
	layout := "2006-01-02"
	if subtotal > customer.Data[0].RemainingCreditLimit {
		err = edenlabs.ErrorValidation("credit_limit", "The amount exceeds the Customer "+customer.Data[0].Custname+"'s credit limit")
		return

	}

	resGP, err := s.opt.Client.BridgeServiceGrpc.CreateSalesInvoiceGP(ctx, &bridgeService.CreateSalesInvoiceGPRequest{
		Docdate: time.Now().Format(layout),
		// Orignumb:           req.Orignumb,
		// Sopnumbe:           req.Sopnumbe,
		// Docid:              req.Docid,
		// Freight:            req.Freight,
		// Docamnt:            req.Docamnt,
		Prstadcd: customer.Data[0].Adrscode[0].Adrscode,
		Custnmbr: customer.Data[0].Custnmbr,
		Custname: customer.Data[0].Custname,
		Curncyid: " ",
		Subtotal: subtotal,
		Docamnt:  subtotal,
		// Trdisamt:           req.Trdisamt,
		// Miscamnt:           req.Miscamnt,
		// Taxamnt:            req.Taxamnt,
		GnlRequestShipDate: req.DeliveryDateStr,
		GnlRegion:          tempAddress.AdministrativeDiv.GnlRegion,
		GnlWrtId:           req.WrtID,
		GnlArchetypeId:     tempAddress.GnL_Archetype_ID,
		// GnlOrderChannel:    req.GnlOrderChannel,
		// GnlSoCodeApps:      req.GnlSoCodeApps,
		// GnlTotalweight:     req.GnlTotalWeight,
		// Userid:             req.UserID,
		AmountReceived: &bridgeService.CreateSalesInvoiceGPRequest_AmountReceived{
			// Amount:   req.AmountReceived.Amount,
			// Chekbkid: req.AmountReceived.Chekbkid,
		},
		VoucherApply: vourchers,
		Detailitems:  details,
		Pymtrmid:     customer.Data[0].Pymtrmid[0].Pymtrmid,
		Locncode:     tempAddress.Locncode,
		Docid:        "SIN",
		// Shipmthd: ,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "sales invoice")
		return
	}

	userID := ctx.Value(constants.KeyUserID).(int64)
	_, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &auditService.CreateLogRequest{
		Log: &auditService.Log{
			UserId:      userID,
			ReferenceId: resGP.Sopnumbe,
			Type:        "sales_invoice",
			Function:    "create",
			CreatedAt:   timestamppb.New(time.Now()),
		},
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpc("audit")
		return
	}
	res = &dto.SalesInvoiceResponse{
		// ID:                userID,
		Code:    resGP.Sopnumbe,
		CodeExt: resGP.Sopnumbe,
	}
	return
}

func (s *SalesInvoiceService) GetListGP(ctx context.Context, req dto.GetSalesInvoiceGPRequest) (res []*bridgeService.SalesInvoiceGP, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesInvoiceService.GetListGP")
	defer span.End()

	var si *bridgeService.GetSalesInvoiceGPListResponse

	if si, err = s.opt.Client.BridgeServiceGrpc.GetSalesInvoiceGPList(ctx, &bridgeService.GetSalesInvoiceGPListRequest{
		Limit:    int32(req.Limit),
		Offset:   int32(req.Offset),
		SoNumber: req.Sopnumbe,
	}); err != nil || !si.Succeeded {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "sales invoice")
		return
	}

	total = int64(len(si.Data))
	res = si.Data

	return
}

func (s *SalesInvoiceService) GetOrderPerformance(ctx context.Context, req dto.SalesInvoiceListRequest) (res []*dto.SalesInvoiceResponse, summaryOrderPerformance *dto.OrderPerformance, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesInvoiceService.Get")
	defer span.End()

	// get sales invoice from bridge
	var siRes *bridgeService.GetSalesInvoiceGPListResponse
	var recognitionDateFrom, recognitionDateTo string
	if timex.IsValid(req.RecognitionDateFrom) {
		recognitionDateFrom = req.RecognitionDateFrom.Format(timex.InFormatDate)
	}

	if timex.IsValid(req.RecognitionDateTo) {
		recognitionDateTo = req.RecognitionDateTo.Format(timex.InFormatDate)
	}

	siRes, err = s.opt.Client.BridgeServiceGrpc.GetSalesInvoiceGPList(ctx, &bridgeService.GetSalesInvoiceGPListRequest{
		Limit:         req.Limit,
		Offset:        req.Offset,
		DocdateFrom:   recognitionDateFrom,
		DocdateTo:     recognitionDateTo,
		Locncode:      req.SiteID,
		Custnumber:    req.CustomerID,
		GnlCustTypeId: "BTY0015",
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "sales invoice")
		return
	}

	datas := []*dto.SalesInvoiceResponse{}
	for _, si := range siRes.Data {
		var docDate time.Time
		docDate, err = time.Parse("2006-01-02T15:04:05", si.Docdate)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("doc_date")
			return
		}
		var dueDate time.Time
		dueDate, err = time.Parse("2006-01-02T15:04:05", si.Duedate)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("due_date")
			return
		}
		var custType, custTypeDesc string
		var customerDetail *bridgeService.GetCustomerGPResponse
		if si.SalesOrder[0].Custnmbr != "" {
			customerDetail, _ = s.opt.Client.BridgeServiceGrpc.GetCustomerGPList(ctx, &bridgeService.GetCustomerGPListRequest{
				Id:             si.SalesOrder[0].Custnmbr,
				Limit:          1,
				Offset:         0,
				CustomerTypeId: "BTY0015",
				Inactive:       "0",
			})
			if len(customerDetail.Data) == 0 {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "sales invoice")
				return
			}
			if len(customerDetail.Data[0].CustomerType) != 0 {
				custType = customerDetail.Data[0].CustomerType[0].GnL_Cust_Type_ID
				custTypeDesc = customerDetail.Data[0].CustomerType[0].GnL_CustType_Description
			}

		} else {
			customerDetail = &bridgeService.GetCustomerGPResponse{}
			customerDetail.Data = append(customerDetail.Data, &bridgeService.CustomerGP{})
		}
		var tempSI *dto.SalesInvoiceResponse

		tempSI = &dto.SalesInvoiceResponse{
			Code:             si.Sopnumbe,
			CodeExt:          si.Sopnumbe,
			Status:           int8(si.Sopstatus),
			RecognitionDate:  docDate,
			DueDate:          dueDate,
			BillingAddress:   si.Address,
			DeliveryFee:      si.Frtamnt,
			CustomerTypeID:   custType,
			CustomerTypeDesc: custTypeDesc,
			TotalPrice:       si.Subtotal,
			TotalCharge:      si.Ordocamt,
			RemainingAmount:  si.RemainingAmount,
			SalesOrder: &dto.SalesOrderResponse{
				CustomerID: customerDetail.Data[0].Custnmbr,
				AddressGP: &dto.AddressResponse{
					CustomerName: customerDetail.Data[0].Custname,
					Customer: &dto.CustomerResponse{
						Name: customerDetail.Data[0].Custname,
					},
				},
			},
			SiteID: si.Locncode,
		}

		for _, v := range si.Details {

			tempSI.SalesInvoiceItem = append(tempSI.SalesInvoiceItem, &dto.SalesInvoiceItemResponse{
				SalesInvoiceID: tempSI.Code,
				ItemID:         v.Itemnmbr,
				InvoiceQty:     v.Quantity,
				UnitPrice:      v.Unitprce,
				Subtotal:       v.Xtndprce,
				UomName:        v.Uofm,
			})
		}

		datas = append(datas, tempSI)

	}
	rest := datas

	// Map to put item with their qty
	filterProduct := make(map[string]float64)

	for _, n := range rest {
		for _, t := range n.SalesInvoiceItem {
			filterProduct[t.ItemID] += t.InvoiceQty
		}

	}

	// Find the product with the highest quantity (TOP PRODUCT)
	var highestProduct string
	var highestQty float64

	for ItemID, InvoiceQty := range filterProduct {
		if InvoiceQty > highestQty {
			highestProduct = ItemID
			highestQty = InvoiceQty
		}
	}

	// Order Performance
	selectedItems := []*dto.SalesInvoiceResponse{}

	// Select only SI object that has top product
	var (
		totalWeights, totalCharges float64
		totalOrder                 int64
	)

	for _, n := range rest {
		for _, t := range n.SalesInvoiceItem {
			if t.ItemID == highestProduct {
				selectedItems = append(selectedItems, n)
				break
			}
		}
	}

	for _, objectSI := range selectedItems {
		for _, objectDetailSI := range objectSI.SalesInvoiceItem {
			totalCharges += objectDetailSI.Subtotal
			// only count the qty of top product
			if objectDetailSI.ItemID == highestProduct {
				totalWeights += objectDetailSI.InvoiceQty
			}
		}
		totalOrder += 1
	}

	avgSales := float64(0)
	if totalOrder > 0 {
		avgSales = totalCharges / float64(totalOrder)
	}

	// Get detail item
	var itemDetail *catalog_service.GetItemDetailByInternalIdResponse
	itemDetail, err = s.opt.Client.CatalogServiceGrpc.GetItemDetailMasterComplexByInternalID(ctx, &catalog_service.GetItemDetailByInternalIdRequest{
		ItemIdGp: highestProduct,
	})
	if err != nil {
		err = edenlabs.ErrorRpcNotFound("catalog", "item")
		return
	}

	summaryOrderPerformance = &dto.OrderPerformance{
		ProductID:   highestProduct,
		ProductName: itemDetail.Data.Description,
		QtySell:     totalWeights,
		AvgSales:    avgSales,
		OrderTotal:  totalOrder,
	}

	return
}
