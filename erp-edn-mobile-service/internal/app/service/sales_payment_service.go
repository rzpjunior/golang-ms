package service

import (
	"context"
	"fmt"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/timex"
	"git.edenfarm.id/edenlabs/edenlabs/utils"
	"git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/global"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/internal/app/dto"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

func NewServiceSalesPayment() ISalesPaymentService {
	m := new(SalesPaymentService)
	m.opt = global.Setup.Common
	return m
}

type ISalesPaymentService interface {
	Get(ctx context.Context, req dto.SalesPaymentListRequest) (res []*dto.SalesPaymentResponse, total int64, err error)
	GetByID(ctx context.Context, req dto.SalesPaymentDetailRequest) (res *dto.SalesPaymentResponse, err error)
	CreateGP(ctx context.Context, req *dto.CreateSalesPaymentRequest) (res *dto.SalesPaymentResponse, err error)
	GetListGP(ctx context.Context, req dto.GetSalesPaymentGPListRequest) (res []*dto.SalesPaymentGP, total int64, err error)
	GetDetailGP(ctx context.Context, id string) (res *dto.SalesPaymentGP, err error)
	GetPaymentPerformance(ctx context.Context, req dto.PerformancePaymentRequest) (res []*dto.SalesPaymentResponse, summaryPaymentPerformance *dto.PaymentPerformance, err error)
}

type SalesPaymentService struct {
	opt opt.Options
}

func (s *SalesPaymentService) Get(ctx context.Context, req dto.SalesPaymentListRequest) (res []*dto.SalesPaymentResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesPaymentService.Get")
	defer span.End()

	// get sales invoice from bridge
	var spRes *bridgeService.GetSalesPaymentGPResponse
	var recognitionDateFrom, recognitionDateTo string
	if timex.IsValid(req.RecognitionDateFrom) {
		recognitionDateFrom = req.RecognitionDateFrom.Format(timex.InFormatDate)
	}

	if timex.IsValid(req.RecognitionDateTo) {
		recognitionDateTo = req.RecognitionDateTo.Format(timex.InFormatDate)
	}

	spRes, err = s.opt.Client.BridgeServiceGrpc.GetSalesPaymentGPList(ctx, &bridgeService.GetSalesPaymentGPListRequest{
		Limit:         req.Limit,
		Offset:        req.Offset,
		DocdateFrom:   recognitionDateFrom,
		DocdateTo:     recognitionDateTo,
		Docnumbr:      req.Search,
		Locncode:      req.SiteID,
		OrderBy:       "desc",
		Custnmbr:      req.CustomerID,
		Status:        req.Status,
		GnlCustTypeId: "BTY0015",
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "sales payment")
		return
	}

	datas := []*dto.SalesPaymentResponse{}
	for _, sp := range spRes.Data {
		var docDate time.Time
		docDate, err = time.Parse("2006-01-02", sp.Docdate)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("doc_date")
			return
		}
		var createdDate time.Time
		createdDate, err = time.Parse("2006-01-02", sp.Creatddt)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("created_date")
			return
		}
		var tempSI *dto.SalesInvoiceResponse

		tempSI = &dto.SalesInvoiceResponse{}
		// get SalesInvoice from bridge
		if len(sp.SalesInvoice) > 0 {
			var siRes *bridgeService.GetSalesInvoiceGPDetailResponse
			siRes, err = s.opt.Client.BridgeServiceGrpc.GetSalesInvoiceGPDetail(ctx, &bridgeService.GetSalesInvoiceGPDetailRequest{
				Id: sp.SalesInvoice[0].Sopnumbe,
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "sales invoice")
				return
			}
			var docDateSI time.Time
			docDateSI, err = time.Parse("2006-01-02T15:04:05", siRes.Data[0].Docdate)
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
			tempSI = &dto.SalesInvoiceResponse{
				Code:            siRes.Data[0].Sopnumbe,
				TotalCharge:     siRes.Data[0].Ordocamt,
				RecognitionDate: docDateSI,
				DueDate:         dueDate,
				RemainingAmount: sp.SalesInvoice[0].RemainingAmount,
			}
		}
		paidOff := 0
		if tempSI.RemainingAmount == 0 {
			paidOff = 1
		}
		datas = append(datas, &dto.SalesPaymentResponse{
			ID:              sp.Docnumbr,
			Code:            sp.Docnumbr,
			RecognitionDate: docDate,
			Amount:          sp.Ortrxamt,
			// BankReceiveNum:   sp.,
			PaidOff: int8(paidOff),
			// ImageUrl:         sp.ImageUrl,
			CreatedAt: createdDate,
			// CreatedBy:        sp.CreatedBy,
			// CancellationNote: sp.,
			// ReceivedDate:     sp.ReceivedDate.AsTime(),
			Status:       int8(sp.Dcstatus),
			CustomerID:   sp.Custnmbr,
			CustomerName: sp.Custname,
			// Note:             sp.Note,
			SalesInvoice: tempSI,
		})
	}
	res = datas
	total = int64(spRes.TotalRecords)
	return
}

func (s *SalesPaymentService) GetByID(ctx context.Context, req dto.SalesPaymentDetailRequest) (res *dto.SalesPaymentResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesPaymentService.GetByID")
	defer span.End()

	// get SalesPayment from bridge
	var (
		sp   *bridgeService.GetSalesPaymentGPResponse
		cust *bridgeService.GetCustomerGPResponse
		// site *bridgeService.GetSiteDetailResponse
	)
	sp, err = s.opt.Client.BridgeServiceGrpc.GetSalesPaymentGPList(ctx, &bridgeService.GetSalesPaymentGPListRequest{
		Limit:  1,
		Offset: 0,
		// DocdateFrom: recognitionDateFrom,
		// DocdateTo:   recognitionDateTo,
		Docnumbr: req.Id,
		// Locncode:    req.SiteID,
		OrderBy: "desc",
		// Custnmbr:    req.CustomerID,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "sales payment")
		return
	}
	// sp, err = s.opt.Client.BridgeServiceGrpc.GetSalesPaymentGPDetail(ctx, &bridgeService.GetSalesPaymentGPDetailRequest{
	// 	Id: req.Id,
	// })
	// if err != nil {
	// 	span.RecordError(err)
	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 	err = edenlabs.ErrorRpcNotFound("bridge", "sales payment")
	// 	return
	// }

	cust, err = s.opt.Client.BridgeServiceGrpc.GetCustomerGPList(ctx, &bridgeService.GetCustomerGPListRequest{
		Id:             sp.Data[0].Custnmbr,
		Limit:          1,
		Offset:         0,
		CustomerTypeId: "BTY0015",
		Inactive:       "0",
	})

	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "sales invoice")
		return
	}
	//validasi edn
	if len(cust.Data[0].CustomerType) > 0 {
		if cust.Data[0].CustomerType[0].GnL_Cust_Type_ID != "BTY0015" {
			err = edenlabs.ErrorValidation("customer_type", "customer is not edn customer")
			return
		}
	} else {
		err = edenlabs.ErrorValidation("customer_type", "customer is not edn customer")
		return
	}
	// site, err = s.opt.Client.BridgeServiceGrpc.GetSiteDetail(ctx, &bridgeService.GetSiteDetailRequest{
	// 	// Id: 1,
	// })
	// if err != nil {
	// 	span.RecordError(err)
	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 	err = edenlabs.ErrorRpcNotFound("bridge", "sales invoice")
	// 	return
	// }

	var docDate time.Time
	if sp.Data[0].Docdate != "" {
		docDate, err = time.Parse("2006-01-02", sp.Data[0].Docdate)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("doc_date")
			return
		}
	}

	var createdDate time.Time
	if sp.Data[0].Creatddt != "" {
		createdDate, err = time.Parse("2006-01-02", sp.Data[0].Creatddt)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("created_date")
			return
		}
	}

	var tempSI *dto.SalesInvoiceResponse
	tempSI = &dto.SalesInvoiceResponse{}
	// get SalesInvoice from bridge
	if len(sp.Data[0].SalesInvoice) > 0 {
		var siRes *bridgeService.GetSalesInvoiceGPDetailResponse
		siRes, err = s.opt.Client.BridgeServiceGrpc.GetSalesInvoiceGPDetail(ctx, &bridgeService.GetSalesInvoiceGPDetailRequest{
			Id: sp.Data[0].SalesInvoice[0].Sopnumbe,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "sales invoice")
			return
		}
		var docDateSI time.Time
		docDateSI, err = time.Parse("2006-01-02T15:04:05", siRes.Data[0].Docdate)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("doc_date")
			return
		}
		var dueDateSI time.Time
		dueDateSI, err = time.Parse("2006-01-02T15:04:05", siRes.Data[0].Duedate)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("due_date")
			return
		}

		tempSI = &dto.SalesInvoiceResponse{
			Code:            siRes.Data[0].Sopnumbe,
			TotalCharge:     siRes.Data[0].Ordocamt,
			RecognitionDate: docDateSI,
			DueDate:         dueDateSI,
			SiteID:          siRes.Data[0].Locncode,
			RemainingAmount: sp.Data[0].SalesInvoice[0].RemainingAmount,
		}
	}
	var site *bridgeService.GetSiteGPResponse
	if tempSI.SiteID != "" {
		site, err = s.opt.Client.BridgeServiceGrpc.GetSiteGPDetail(ctx, &bridgeService.GetSiteGPDetailRequest{
			Id: tempSI.SiteID,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "site")
			return
		}
	} else {
		site = &bridgeService.GetSiteGPResponse{}
		site.Data = append(site.Data, &bridgeService.SiteGP{})
	}
	paidOff := 0
	if tempSI.RemainingAmount == 0 {
		paidOff = 1
	}
	res = &dto.SalesPaymentResponse{
		// ID:               sp.Data[0].Id,
		Code:            sp.Data[0].Docnumbr,
		RecognitionDate: docDate,
		Amount:          sp.Data[0].Ortrxamt,
		// BankReceiveNum:   sp.Data[0].BankReceiveNum,
		PaidOff: int8(paidOff),
		// ImageUrl:         sp.Data[0].ImageUrl,
		CreatedAt: createdDate,
		// CreatedBy:        sp.Data[0].CreatedBy,
		// CancellationNote: sp.Data[0].CancellationNote,
		// ReceivedDate:     sp.Data[0].ReceivedDate.AsTime(),
		Status: int8(sp.Data[0].Dcstatus),
		// Note:             sp.Data[0].Note,
		PaymentMethod: sp.Data[0].PaymentMethod,
		CustomerName:  cust.Data[0].Custname,
		SalesInvoice: &dto.SalesInvoiceResponse{
			SalesOrder: &dto.SalesOrderResponse{
				// AddressGP: &dto.AddressGP{
				// 	Customer: &dto.CustomerGP{
				// 		Custnmbr: cust.Data[0].Custnmbr,
				// 		Custclas: cust.Data[0].Custclas,
				// 		Custname: cust.Data[0].Custname,
				// 		Cprcstnm: cust.Data[0].Cprcstnm,
				// 		Cntcprsn: cust.Data[0].Cntcprsn,
				// 		Stmtname: cust.Data[0].Stmtname,
				// 		Shrtname: cust.Data[0].Shrtname,
				// 		Upszone:  cust.Data[0].Upszone,
				// 		Shipmthd: cust.Data[0].Shipmthd,
				// 	},
				// },
				Site: &dto.SiteResponse{
					// ID:            site.Data.Id,
					Code:           site.Data[0].Locncode,
					Name:           site.Data[0].Locndscr,
					Description:    site.Data[0].Locndscr,
					Address:        site.Data[0].AddresS1 + " " + site.Data[0].AddresS2 + " " + site.Data[0].AddresS3,
					PhoneNumber:    site.Data[0].PhonE1,
					AltPhoneNumber: site.Data[0].PhonE2,
					// Name: site.Data[0].,
					// Status:        int8(site.Data.Status),
					// StatusConvert: statusx.ConvertStatusValue(int8(site.Data.Status)),
					// CreatedAt:     site.Data.CreatedAt.AsTime(),
					// UpdatedAt:     site.Data.UpdatedAt.AsTime(),

				},
			},
			Code:            tempSI.Code,
			TotalCharge:     tempSI.TotalCharge,
			RecognitionDate: tempSI.RecognitionDate,
			DueDate:         tempSI.DueDate,
			RemainingAmount: tempSI.RemainingAmount,
		},
	}

	fmt.Print("tes")
	return
}

func (s *SalesPaymentService) CreateGP(ctx context.Context, req *dto.CreateSalesPaymentRequest) (res *dto.SalesPaymentResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesPaymentService.CreateGP")
	defer span.End()

	si, _ := s.opt.Client.BridgeServiceGrpc.GetSalesInvoiceGPDetail(ctx, &bridgeService.GetSalesInvoiceGPDetailRequest{
		Id: req.SalesInvoiceID,
	})

	if len(si.Data) == 0 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "sales invoice")
		return
	}

	// checkbkid, _ := s.opt.Client.ConfigurationServiceGrpc.GetConfigAppList(ctx, &configuration_service.GetConfigAppListRequest{
	// 	Offset:    0,
	// 	Limit:     1,
	// 	Attribute: req.RegionID,
	// 	Field:     "EDN App Checkbook ID",
	// })

	// if len(checkbkid.Data) == 0 {
	// 	span.RecordError(err)
	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 	err = edenlabs.ErrorRpcNotFound("configuration", "checkbook id")
	// 	return
	// }

	layout := "2006-01-02"
	var payloadGP *bridgeService.CreateSalesPaymentGPnonPBDRequest
	docDate := time.Now()
	var applyTo []*bridgeService.CreateSalesPaymentGPnonPBDRequest_ApplyTo
	applyTo = append(applyTo, &bridgeService.CreateSalesPaymentGPnonPBDRequest_ApplyTo{
		Sopnumbe:    req.SalesInvoiceID,
		ApplyAmount: req.Amount,
	})
	payloadGP = &bridgeService.CreateSalesPaymentGPnonPBDRequest{
		Chekbkid: req.CheckbookID,
		Docdate:  docDate.Format(layout),
		Custnmbr: si.Data[0].Custnmbr,
		Cshrctyp: int32(utils.ToInt(req.PaymentMethodID)),
		Ortrxamt: req.Amount,
		Trxdscrn: "",
		ApplyTo:  applyTo,
	}
	resGP, err := s.opt.Client.BridgeServiceGrpc.CreateSalesPaymentGPnonPBD(ctx, payloadGP)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "sales payment")
		return
	}

	res = &dto.SalesPaymentResponse{
		Code:            resGP.Paymentnumber,
		Amount:          req.Amount,
		BankReceiveNum:  req.CheckbookID,
		RecognitionDate: docDate,
		CreatedAt:       docDate,
	}
	return
}

func (s *SalesPaymentService) GetListGP(ctx context.Context, req dto.GetSalesPaymentGPListRequest) (res []*dto.SalesPaymentGP, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesPaymentService.GetListGP")
	defer span.End()

	var sp *bridgeService.GetSalesPaymentGPResponse

	if sp, err = s.opt.Client.BridgeServiceGrpc.GetSalesPaymentGPList(ctx, &bridgeService.GetSalesPaymentGPListRequest{
		Limit:       int32(req.Limit),
		Offset:      int32(req.Offset),
		Docnumbr:    req.Custnmbr,
		DocdateFrom: req.DocdateFrom,
		DocdateTo:   req.DocdateTo,
	}); err != nil || !sp.Succeeded {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "sales payment")
		return
	}

	for _, payment := range sp.Data {
		var (
			soList []*dto.SO_SalesPaymentGP
			siList []*dto.SI_SalesPaymentGP
		)
		for _, so := range payment.SalesOrder {
			soList = append(soList, &dto.SO_SalesPaymentGP{
				Orignumb: so.Orignumb,
				Ordrdate: so.Ordrdate,
			})
		}
		for _, si := range payment.SalesInvoice {
			siList = append(siList, &dto.SI_SalesPaymentGP{
				Sopnumbe:  si.Sopnumbe,
				Docdate:   si.Docdate,
				GnlRegion: si.GnlRegion,
				Locncode:  si.Locncode,
			})
		}
		res = append(res, &dto.SalesPaymentGP{
			Docnumbr:      payment.Docnumbr,
			Docdate:       payment.Docdate,
			Custnmbr:      payment.Custnmbr,
			Custname:      payment.Custname,
			Curncyid:      payment.Curncyid,
			Cshrctyp:      int(payment.Cshrctyp),
			PaymentMethod: payment.PaymentMethod,
			Dcstatus:      int(payment.Dcstatus),
			Ortrxamt:      payment.Ortrxamt,
			Creatddt:      payment.Creatddt,
			SalesInvoice:  siList,
			SalesOrder:    soList,
		})
	}

	total = int64(sp.TotalRecords)
	// total = int64(len(sp.Data))
	return
}

func (s *SalesPaymentService) GetDetailGP(ctx context.Context, id string) (res *dto.SalesPaymentGP, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesPaymentService.GetDetailGP")
	defer span.End()

	var sp *bridgeService.GetSalesPaymentGPResponse

	if sp, err = s.opt.Client.BridgeServiceGrpc.GetSalesPaymentGPDetail(ctx, &bridgeService.GetSalesPaymentGPDetailRequest{
		Id: id,
	}); err != nil || !sp.Succeeded {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "purchase order")
		return
	}

	if len(sp.Data) <= 0 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "purchase order")
		return
	}

	var (
		soList []*dto.SO_SalesPaymentGP
		siList []*dto.SI_SalesPaymentGP
	)
	for _, so := range sp.Data[0].SalesOrder {
		soList = append(soList, &dto.SO_SalesPaymentGP{
			Orignumb: so.Orignumb,
			Ordrdate: so.Ordrdate,
		})
	}
	for _, si := range sp.Data[0].SalesInvoice {
		siList = append(siList, &dto.SI_SalesPaymentGP{
			Sopnumbe:  si.Sopnumbe,
			Docdate:   si.Docdate,
			GnlRegion: si.GnlRegion,
			Locncode:  si.Locncode,
		})
	}
	res = &dto.SalesPaymentGP{
		Docnumbr:      sp.Data[0].Docnumbr,
		Docdate:       sp.Data[0].Docdate,
		Custnmbr:      sp.Data[0].Custnmbr,
		Custname:      sp.Data[0].Custname,
		Curncyid:      sp.Data[0].Curncyid,
		Cshrctyp:      int(sp.Data[0].Cshrctyp),
		PaymentMethod: sp.Data[0].PaymentMethod,
		Dcstatus:      int(sp.Data[0].Dcstatus),
		Ortrxamt:      sp.Data[0].Ortrxamt,
		Creatddt:      sp.Data[0].Creatddt,
		SalesInvoice:  siList,
		SalesOrder:    soList,
	}
	return
}

func (s *SalesPaymentService) GetPaymentPerformance(ctx context.Context, req dto.PerformancePaymentRequest) (res []*dto.SalesPaymentResponse, summaryPaymentPerformance *dto.PaymentPerformance, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesPaymentService.Get")
	defer span.End()

	// get sales invoice & sales payment before calculation
	var salesInvoicesFinish, salesInvoicesOnProgress []*dto.SalesInvoiceResponse
	var salesPayments []*dto.SalesPaymentResponse

	// get customer from bridge
	var customer *bridgeService.GetCustomerGPResponse
	if customer, err = s.opt.Client.BridgeServiceGrpc.GetCustomerGPDetail(ctx, &bridgeService.GetCustomerGPDetailRequest{
		Id: req.CustomerID,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "customer")
		return
	}

	// get sales invoice from bridge
	var siRes *bridgeService.GetSalesInvoiceGPListResponse
	// if timex.IsValid(req.RecognitionDateFrom) {
	// 	recognitionDateFrom := req.RecognitionDateFrom.Format(timex.InFormatDate)
	// }

	// if timex.IsValid(req.RecognitionDateTo) {
	// 	recognitionDateTo := req.RecognitionDateTo.Format(timex.InFormatDate)
	// }

	siRes, err = s.opt.Client.BridgeServiceGrpc.GetSalesInvoiceGPList(ctx, &bridgeService.GetSalesInvoiceGPListRequest{
		Limit:  req.Limit,
		Offset: req.Offset,
		// DocdateFrom:   recognitionDateFrom,
		// DocdateTo:     recognitionDateTo,
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

	// declare variable for calculation
	datasInvoiceFinish := []*dto.SalesInvoiceResponse{}
	datasInvoiceOnProgress := []*dto.SalesInvoiceResponse{}

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
			err = edenlabs.ErrorInvalid("doc_date")
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

		// seperate invoice finished and unfinished
		if si.RemainingAmount > 0 {
			datasInvoiceOnProgress = append(datasInvoiceOnProgress, tempSI)
		} else {
			datasInvoiceFinish = append(datasInvoiceFinish, tempSI)
		}

	}

	var countPaymentFinish, countPaymentOnProgress float64
	var totalFinalPaymentFinish, totalFinalPaymentOnProgress float64

	// AMOUNT SALES INVOICE FINISH
	var amountSalesInvoicesFinish, amountSalesInvoicesOnProgress, amountSalesInvoicesDueDate float64
	for _, siFinish := range datasInvoiceFinish {
		amountSalesInvoicesFinish += siFinish.TotalCharge

		// for total payment percentage
		// totalPaymentFinish += (siFinish.TotalCharge - siFinish.RemainingAmount)/siFinish.TotalCharge*100 --> FORMULA
		if siFinish.TotalCharge == 0 {
			totalFinalPaymentFinish += 100 // if total charge is 0, then the payment is 100%
		} else {
			totalFinalPaymentFinish += (siFinish.TotalCharge - siFinish.RemainingAmount) / siFinish.TotalCharge * 100
		}
		countPaymentFinish += 1
	}

	// AMOUNT SALES INVOICE ON PROGRESS
	var totalDiffPayment int
	for _, siOnProgress := range datasInvoiceOnProgress {
		amountSalesInvoicesOnProgress += siOnProgress.TotalCharge

		// for total payment percentage
		totalFinalPaymentOnProgress += (siOnProgress.TotalCharge - siOnProgress.RemainingAmount) / siOnProgress.TotalCharge * 100
		countPaymentOnProgress += 1

		// for total diff payment
		var earliestDate time.Time

		if earliestDate.IsZero() || siOnProgress.DueDate.Before(earliestDate) {
			earliestDate = earliestDate
		}

		duration := earliestDate.Sub(time.Now())
		totalDiffPayment = int(duration.Hours() / 24)
	}

	// AMOUNT SALES INVOICE DUE DATE
	for _, siOnProgressDueDate := range salesInvoicesOnProgress {
		if siOnProgressDueDate.DueDate.Before(time.Now()) {
			amountSalesInvoicesDueDate += siOnProgressDueDate.TotalCharge
		}
	}

	// get sales payment from bridge
	var spRes *bridgeService.GetSalesPaymentGPResponse
	var recognitionDateFrom, recognitionDateTo string
	if timex.IsValid(req.RecognitionDateFrom) {
		recognitionDateFrom = req.RecognitionDateFrom.Format(timex.InFormatDate)
	}

	if timex.IsValid(req.RecognitionDateTo) {
		recognitionDateTo = req.RecognitionDateTo.Format(timex.InFormatDate)
	}

	spRes, err = s.opt.Client.BridgeServiceGrpc.GetSalesPaymentGPList(ctx, &bridgeService.GetSalesPaymentGPListRequest{
		Limit:       req.Limit,
		Offset:      req.Offset,
		DocdateFrom: recognitionDateFrom,
		DocdateTo:   recognitionDateTo,
		Docnumbr:    req.Search,
		Locncode:    req.SiteID,
		OrderBy:     "desc",
		Custnmbr:    req.CustomerID,
		Status:      req.Status,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "sales payment")
		return
	}

	datasPayment := []*dto.SalesPaymentResponse{}
	for _, sp := range spRes.Data {
		// var docDate time.Time
		// docDate, err = time.Parse("2006-01-02T15:04:05", sp.Docdate)
		// if err != nil {
		// 	span.RecordError(err)
		// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
		// 	err = edenlabs.ErrorInvalid("doc_date")
		// 	return
		// }
		// var createdDate time.Time
		// createdDate, err = time.Parse("2006-01-02T15:04:05", sp.Creatddt)
		// if err != nil {
		// 	span.RecordError(err)
		// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
		// 	err = edenlabs.ErrorInvalid("doc_date")
		// 	return
		// }
		var tempSI *dto.SalesInvoiceResponse

		// get SalesInvoice from bridge
		if len(sp.SalesInvoice) > 0 {
			var siRes *bridgeService.GetSalesInvoiceGPDetailResponse
			siRes, err = s.opt.Client.BridgeServiceGrpc.GetSalesInvoiceGPDetail(ctx, &bridgeService.GetSalesInvoiceGPDetailRequest{
				Id: sp.SalesInvoice[0].Sopnumbe,
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "sales invoice")
				return
			}
			var docDateSI time.Time
			docDateSI, err = time.Parse("2006-01-02T15:04:05", siRes.Data[0].Docdate)
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
				err = edenlabs.ErrorInvalid("doc_date")
				return
			}
			tempSI = &dto.SalesInvoiceResponse{
				Code:            siRes.Data[0].Sopnumbe,
				TotalCharge:     siRes.Data[0].Ordocamt,
				RecognitionDate: docDateSI,
				DueDate:         dueDate,
			}
		}

		datasPayment = append(datasPayment, &dto.SalesPaymentResponse{
			// ID:               sp.Id,
			Code: sp.Docnumbr,
			// RecognitionDate: docDate,
			Amount:       sp.Ordocamt,
			PaidOff:      int8(sp.Ortrxamt),
			Status:       int8(sp.Dcstatus),
			CustomerID:   sp.Custnmbr,
			CustomerName: sp.Custname,
			SalesInvoice: tempSI,
		})
	}
	salesPayments = datasPayment

	var amountSalesPaymentsFinish float64

	// AMOUNT SALES PAYMENT FINISH
	for _, spFinish := range salesPayments {
		amountSalesPaymentsFinish += spFinish.Amount
	}

	var creditLimitUsageRemPer, avgPayAmount, overdueDebtRemPer float64
	creditLimitAmount := customer.Data[0].Crlmtamt
	remainingAmount := customer.Data[0].RemainingAmount

	// if credit limit amount is zero / not exist
	if creditLimitAmount == 0 {
		creditLimitUsageRemPer = 0
	} else {
		creditLimitUsageRemPer = (remainingAmount / creditLimitAmount) * 100
	}

	// if there is no finish sales invoices
	if len(salesInvoicesFinish) == 0 {
		avgPayAmount = 0
	} else {
		avgPayAmount = amountSalesPaymentsFinish / float64(len(salesInvoicesFinish))
	}

	if amountSalesPaymentsFinish == 0 {
		overdueDebtRemPer = 0
	} else {
		overdueDebtRemPer = (amountSalesInvoicesDueDate/amountSalesPaymentsFinish - (amountSalesInvoicesFinish + amountSalesInvoicesOnProgress)) * 100
	}

	// for total payment percentage
	var totalPaymentPercentage, totalPaymentFinishAndOnProgress, countPaymentFinishAndOnProgress float64
	totalPaymentFinishAndOnProgress = totalFinalPaymentFinish + totalFinalPaymentOnProgress
	countPaymentFinishAndOnProgress = countPaymentFinish + countPaymentOnProgress
	// set total payment percentage to 0 if there is no payment record
	if countPaymentFinishAndOnProgress == 0 {
		totalPaymentPercentage = 0
	} else {
		totalPaymentPercentage = totalPaymentFinishAndOnProgress / countPaymentFinishAndOnProgress
	}

	summaryPaymentPerformance = &dto.PaymentPerformance{
		CreditLimitAmount:                   customer.Data[0].Crlmtamt,
		CreditLimitRemainingAmount:          customer.Data[0].RemainingAmount,
		RemainingOutstanding:                amountSalesPaymentsFinish - (amountSalesInvoicesFinish + amountSalesInvoicesOnProgress),
		CreditLimitUsageRemainingPercentage: creditLimitUsageRemPer,
		OverdueDebtAmount:                   amountSalesInvoicesDueDate,
		OverdueDebtRemainingPercentage:      overdueDebtRemPer,
		AveragePaymentAmount:                avgPayAmount,
		AveragePaymentPercentage:            totalPaymentPercentage,
		AveragePaymentPeriod:                totalDiffPayment,
	}

	return
}
