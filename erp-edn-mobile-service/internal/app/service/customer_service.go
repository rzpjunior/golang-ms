package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/global"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/edenlabs/edenlabs/utils"
	dto "git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/crm_service"
)

func NewServiceCustomer() ICustomerService {
	m := new(CustomerService)
	m.opt = global.Setup.Common
	return m
}

type ICustomerService interface {
	GetCustomers(ctx context.Context, req dto.CustomerListRequest) (res []*dto.CustomerResponse, total int64, err error)
	GetCustomerDetailById(ctx context.Context, req dto.CustomerDetailRequest) (res *dto.CustomerResponse, err error)
	GetListGp(ctx context.Context, req dto.GetCustomerGPListRequest) (res []*dto.CustomerGP, total int64, err error)
	GetDetailGp(ctx context.Context, id string) (res *dto.CustomerGP, err error)
	CreateGP(ctx context.Context, req dto.CreateCustomerGPRequest) (res *bridgeService.CreateCustomerGPResponse, err error)
	GetOverdueMitra(ctx context.Context, req dto.CustomerListRequest) (res []*dto.CustomerResponse, total int64, err error)
}

type CustomerService struct {
	opt opt.Options
}

func (s *CustomerService) GetCustomers(ctx context.Context, req dto.CustomerListRequest) (res []*dto.CustomerResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CustomerService.GetCustomers")
	defer span.End()

	var statusGP int32

	switch statusx.ConvertStatusValue(int8(req.Status)) {
	// status convert
	case statusx.Active:
		statusGP = 0
	case statusx.NotActive:
		statusGP = 1
	}

	var custResponse *bridgeService.GetCustomerGPResponse
	custResponse, err = s.opt.Client.BridgeServiceGrpc.GetCustomerGPList(ctx, &bridgeService.GetCustomerGPListRequest{
		Limit:          req.Limit,
		Offset:         req.Offset,
		Inactive:       utils.ToString(statusGP),
		Name:           req.Search,
		Orderby:        req.OrderBy,
		CustomerTypeId: "BTY0015",
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "customer")
		return
	}

	datas := []*dto.CustomerResponse{}

	for _, cust := range custResponse.Data {

		var custType, custTypeDesc string

		if len(cust.CustomerType) != 0 {
			custType = cust.CustomerType[0].GnL_Cust_Type_ID
			custTypeDesc = cust.CustomerType[0].GnL_CustType_Description
		}

		// get region customer from bridge
		var regionCustomer string
		if cust.Adrscode != nil {
			admDivision, err := s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPList(ctx, &bridge_service.GetAdmDivisionGPListRequest{
				AdmDivisionCode: cust.Adrscode[0].AdministrativeDivision[0].CustAddrTypeBillTo[0].GnL_Administrative_Code,
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorValidation("bridge", "customer region")
			}
			regionCustomer = admDivision.Data[0].Region
		}
		var tempPaymentTerm *dto.CustomerGPpymtrmid
		if len(cust.Pymtrmid) > 0 {
			tempPaymentTerm = &dto.CustomerGPpymtrmid{
				Pymtrmid:              cust.Pymtrmid[0].Pymtrmid,
				Code:                  cust.Pymtrmid[0].Pymtrmid,
				CalculateDateFromDays: cust.Pymtrmid[0].CalculateDateFromDays,
			}
		} else {
			tempPaymentTerm = &dto.CustomerGPpymtrmid{}
		}

		// get customer class from bridge
		var (
			custClsDes        string
			creditLimitType   string
			creditLimitAmount float64
		)

		if cust.Custclas != "" {
			custCls, err := s.opt.Client.BridgeServiceGrpc.GetCustomerClassDetail(ctx, &bridge_service.GetCustomerClassDetailRequest{
				Id: cust.Custclas,
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorValidation("bridge", "customer region")
			}
			custClsDes = custCls.Data[0].Clasdscr
			creditLimitAmount = custCls.Data[0].Crlmtamt
		}

		// define credit limit type based on business type
		if cust.GnlBusinessType == 1 {
			creditLimitType = "Badan Usaha"
		} else {
			creditLimitType = "Personal"
		}

		var (
			statusGP   int8
			statusDesc string
		)

		switch statusx.ConvertStatusValue(int8(req.Status)) {
		// status convert
		case statusx.Active:
			statusGP = 1
			statusDesc = "active"
		case statusx.NotActive:
			statusGP = 26
			statusDesc = "not active"
		}

		var paymentGroupFinal string
		if cust.Pymtrmid[0].Pymtrmid == "COD" || cust.Pymtrmid[0].Pymtrmid == "BNS" || cust.Pymtrmid[0].Pymtrmid == "PWD" {
			paymentGroupFinal = "on-delivery"
		} else if cust.Pymtrmid[0].Pymtrmid == "PBD" {
			paymentGroupFinal = "advance"
		} else {
			paymentGroupFinal = "in-term"
		}

		// Get data from repo CRM
		var email string
		customerDetail, err := s.opt.Client.CrmServiceGrpc.GetCustomerDetail(ctx, &crm_service.GetCustomerDetailRequest{
			CustomerIdGp: cust.Custnmbr,
		})
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("crm", "customer")
			// put this because if customer not found in crm, we still want to show the data
			continue
		}
		if customerDetail.Data.Email != "" {
			email = customerDetail.Data.Email
		}

		// Get salesperson
		var (
			salesPersonId   string
			salesPersonName string
		)
		if cust.Slprsnid != "" {

			var salesPersonDetail *bridgeService.GetSalesPersonGPResponse
			salesPersonDetail, err = s.opt.Client.BridgeServiceGrpc.GetSalesPersonGPDetail(ctx, &bridgeService.GetSalesPersonGPDetailRequest{
				Id: cust.Slprsnid,
			})
			if err != nil || len(salesPersonDetail.Data) == 0 {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "Sales person")
			}
			salesPersonId = salesPersonDetail.Data[0].Slprsnid
			salesPersonName = salesPersonDetail.Data[0].Slprsnfn + " " + salesPersonDetail.Data[0].Sprsnsln
		}

		// Get Site from address ship to
		var (
			siteAddressDetail  *bridgeService.GetSiteGPResponse
			siteAddress        *bridgeService.GetAddressGPResponse
			siteCode, siteName string
		)

		if len(cust.Adrscode) > 0 {
			siteAddress, err = s.opt.Client.BridgeServiceGrpc.GetAddressGPDetail(ctx, &bridge_service.GetAddressGPDetailRequest{
				Id: cust.Adrscode[0].Adrscode,
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorValidation("bridge", "address site")
			}

			siteAddressDetail, err = s.opt.Client.BridgeServiceGrpc.GetSiteGPDetail(ctx, &bridge_service.GetSiteGPDetailRequest{
				Id: siteAddress.Data[0].Locncode,
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorValidation("bridge", "address site")
			}
			siteCode = siteAddressDetail.Data[0].Locncode
			siteName = siteAddressDetail.Data[0].Locndscr
		}

		// Get sales invoice based on customer
		var (
			siRes                       *bridgeService.GetSalesInvoiceGPListResponse
			totalRemainingAmountInvoice float64
			earliestDate                time.Time
			earliestDateStr             string
		)

		siRes, err = s.opt.Client.BridgeServiceGrpc.GetSalesInvoiceGPList(ctx, &bridgeService.GetSalesInvoiceGPListRequest{
			Limit:         20,
			Offset:        0,
			OrderBy:       "desc",
			Custnumber:    cust.Custnmbr,
			GnlCustTypeId: "BTY0015",
		})
		if err != nil {
			continue
		}

		if len(siRes.Data) > 0 {
			for _, si := range siRes.Data {
				if si.RemainingAmount > 0 {
					totalRemainingAmountInvoice += si.RemainingAmount
				}
				date, _ := time.Parse("2006-01-02T15:04:05Z", si.Duedate)
				if earliestDate.IsZero() || date.Before(earliestDate) {
					earliestDate = date
					earliestDateStr = earliestDate.Format("2006-01-02")
				}
			}
		}

		datas = append(datas, &dto.CustomerResponse{
			ID:                cust.Custnmbr,
			Name:              cust.Custname,
			Code:              cust.Custnmbr,
			PhoneNumber:       cust.PhonE1,
			AltPhoneNumber:    cust.PhonE2,
			BillingAddress:    cust.AddresS1 + " " + cust.AddresS2 + " " + cust.AddresS3,
			Status:            statusGP,
			StatusDescription: statusDesc,
			ReferralCode:      cust.GnlReferralCode,
			ReferrerCode:      cust.GnlReferrerCode,
			PaymentTerm:       tempPaymentTerm,
			CustomerTypeID:    custType,
			CustomerTypeDesc:  custTypeDesc,

			// Additional
			Region:            regionCustomer,
			ShippingAddress:   cust.AddresS1,
			OverdueDebt:       float64(cust.RemainingAmount),
			CustomerGroup:     "B2B",
			PaymentGroup:      paymentGroupFinal,
			PicName:           cust.Stmtname,
			CustomerClass:     cust.Custclas,
			CustomerClassDesc: custClsDes,
			Email:             email,
			SalesPerson:       salesPersonId,
			SalesPersonName:   salesPersonName,
			SiteCode:          siteCode,
			SiteName:          siteName,

			// Credit Limit
			CreditLimitType:      int32(cust.GnlBusinessType),
			CreditLimitTypeDesc:  creditLimitType,
			CreditLimitAmount:    creditLimitAmount,
			RemainingCreditLimit: cust.RemainingCreditLimit,

			DueDate:              earliestDateStr,
			TotalRemainingAmount: totalRemainingAmountInvoice,
		})
	}
	total = int64(len(custResponse.Data))

	res = datas

	return
}

func (s *CustomerService) GetCustomerDetailById(ctx context.Context, req dto.CustomerDetailRequest) (res *dto.CustomerResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CustomerService.GetCustomerDetailById")
	defer span.End()

	var custResponse *bridgeService.GetCustomerGPResponse
	custResponse, err = s.opt.Client.BridgeServiceGrpc.GetCustomerGPList(ctx, &bridgeService.GetCustomerGPListRequest{
		Limit:  1,
		Offset: 0,
		Id:     req.Id,
	})

	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "customer")
		return
	}

	datas := []*dto.CustomerResponse{}

	for _, cust := range custResponse.Data {

		// get region customer from bridge
		var (
			statusGP            int8
			statusDesc          string
			regionCustomer      string
			stateCustomer       string
			cityCustomer        string
			districtCustomer    string
			subDistrictCustomer string
			zipCustomer         string
		)

		if cust.Inactive == 0 {
			statusGP = 1
			statusDesc = "active"
		} else {
			statusGP = 26
			statusDesc = "not active"
		}

		if cust.Adrscode != nil {
			admDivision, err := s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPList(ctx, &bridge_service.GetAdmDivisionGPListRequest{
				AdmDivisionCode: cust.Adrscode[0].AdministrativeDivision[0].CustAddrTypeShipTo[0].GnL_Administrative_Code,
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorValidation("bridge", "customer region")
			}
			regionCustomer = admDivision.Data[0].Region
			stateCustomer = admDivision.Data[0].State
			cityCustomer = admDivision.Data[0].City
			districtCustomer = admDivision.Data[0].District
			subDistrictCustomer = admDivision.Data[0].Subdistrict
			zipCustomer = admDivision.Data[0].Zipcode
		}

		// get customer class from bridge
		var (
			custClsDes        string
			creditLimitType   string
			creditLimitAmount float64
		)
		if cust.Custclas != "" {
			custCls, err := s.opt.Client.BridgeServiceGrpc.GetCustomerClassDetail(ctx, &bridge_service.GetCustomerClassDetailRequest{
				Id: cust.Custclas,
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorValidation("bridge", "customer region")
			}
			custClsDes = custCls.Data[0].Clasdscr
			creditLimitAmount = custCls.Data[0].Crlmtamt
		}

		// define credit limit type based on business type
		if cust.GnlBusinessType == 1 {
			creditLimitType = "Badan Usaha"
		} else {
			creditLimitType = "Personal"
		}

		// Get Site from address ship to
		var siteAddress *bridgeService.GetAddressGPResponse
		siteAddress, err = s.opt.Client.BridgeServiceGrpc.GetAddressGPDetail(ctx, &bridge_service.GetAddressGPDetailRequest{
			Id: cust.Adrscode[0].Adrscode,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorValidation("bridge", "address site")
		}

		var siteAddressDetail *bridgeService.GetSiteGPResponse
		siteAddressDetail, err = s.opt.Client.BridgeServiceGrpc.GetSiteGPDetail(ctx, &bridge_service.GetSiteGPDetailRequest{
			Id: siteAddress.Data[0].Locncode,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorValidation("bridge", "address site")
		}

		// Get salesperson
		var salesPersonDetail *bridgeService.GetSalesPersonGPResponse
		salesPersonDetail, err = s.opt.Client.BridgeServiceGrpc.GetSalesPersonGPDetail(ctx, &bridgeService.GetSalesPersonGPDetailRequest{
			Id: cust.Slprsnid,
		})
		if err != nil || len(salesPersonDetail.Data) == 0 {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "Sales person")
			return
		}

		// Get data from repo CRM
		customerDetail, err := s.opt.Client.CrmServiceGrpc.GetCustomerDetail(ctx, &crm_service.GetCustomerDetailRequest{
			CustomerIdGp: cust.Custnmbr,
		})
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("crm", "customer")
		}

		var custType, custTypeDesc string

		if len(cust.CustomerType) != 0 {
			custType = cust.CustomerType[0].GnL_Cust_Type_ID
			custTypeDesc = cust.CustomerType[0].GnL_CustType_Description
		}

		var paymentGroupFinal string
		if cust.Pymtrmid[0].Pymtrmid == "COD" || cust.Pymtrmid[0].Pymtrmid == "BNS" || cust.Pymtrmid[0].Pymtrmid == "PWD" {
			paymentGroupFinal = "on-delivery"
		} else if cust.Pymtrmid[0].Pymtrmid == "PBD" {
			paymentGroupFinal = "advance"
		} else {
			paymentGroupFinal = "in-term"
		}

		// Get address from bridge
		var noteAddress string
		addrDetail, err := s.opt.Client.BridgeServiceGrpc.GetAddressGPList(ctx, &bridge_service.GetAddressGPListRequest{
			ExcludeType:    "bill_to,statement_to",
			CustomerNumber: cust.Custnmbr,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorValidation("bridge", "address detail")
		}
		// note: index [0] is bill to, index [1] is shipment to, index [2] is statement to
		noteAddress = addrDetail.Data[0].GnL_Address_Note

		// Get sales invoice based on customer
		var (
			siRes                       *bridgeService.GetSalesInvoiceGPListResponse
			totalRemainingAmountInvoice float64
			earliestDate                time.Time
			earliestDateStr             string
		)

		siRes, err = s.opt.Client.BridgeServiceGrpc.GetSalesInvoiceGPList(ctx, &bridgeService.GetSalesInvoiceGPListRequest{
			Limit:         20,
			Offset:        0,
			OrderBy:       "desc",
			Custnumber:    cust.Custnmbr,
			GnlCustTypeId: "BTY0015",
		})
		if err != nil {
			continue
		}

		if len(siRes.Data) > 0 {
			for _, si := range siRes.Data {
				if si.RemainingAmount > 0 {
					totalRemainingAmountInvoice += si.RemainingAmount
				}
				date, _ := time.Parse("2006-01-02T15:04:05Z", si.Duedate)
				if earliestDate.IsZero() || date.Before(earliestDate) {
					earliestDate = date
					earliestDateStr = earliestDate.Format("2006-01-02")
				}
			}
		}

		customerPhotos := strings.Split(customerDetail.Data.CustomerPhotosUrl, ",")

		datas = append(datas, &dto.CustomerResponse{
			ID:                cust.Custnmbr,
			Name:              cust.Custname,
			Code:              cust.Custnmbr,
			PhoneNumber:       cust.PhonE1,
			AltPhoneNumber:    cust.PhonE2,
			BillingAddress:    cust.AddresS1 + " " + cust.AddresS2 + " " + cust.AddresS3,
			Status:            statusGP,
			StatusDescription: statusDesc,
			ReferralCode:      cust.GnlReferralCode,
			ReferrerCode:      cust.GnlReferrerCode,
			PaymentTerm: &dto.CustomerGPpymtrmid{
				Pymtrmid:              cust.Pymtrmid[0].Pymtrmid,
				Code:                  cust.Pymtrmid[0].Pymtrmid,
				CalculateDateFromDays: cust.Pymtrmid[0].CalculateDateFromDays,
			},

			// Additional
			Region:                  regionCustomer,
			ShippingAddress:         cust.AddresS2,
			OverdueDebt:             float64(cust.RemainingAmount),
			Email:                   customerDetail.Data.Email,
			Note:                    noteAddress,
			PicName:                 cust.Stmtname,
			CustomerGroup:           "B2B",
			PaymentGroup:            paymentGroupFinal,
			CustomerTypeID:          custType,
			CustomerTypeDesc:        custTypeDesc,
			KTPPhotosUrlArr:         []string{customerDetail.Data.KtpPhotosUrl},
			MerchantPhotosUrlArr:    customerPhotos,
			PriceLevel:              cust.Prclevel,
			CustomerTypeCreditLimit: cust.Crlmttyp,
			SalesPerson:             salesPersonDetail.Data[0].Slprsnid,
			SalesPersonName:         salesPersonDetail.Data[0].Slprsnfn + " " + salesPersonDetail.Data[0].Sprsnsln,
			SiteCode:                siteAddressDetail.Data[0].Locncode,
			SiteName:                siteAddressDetail.Data[0].Locndscr,
			CustomerClass:           cust.Custclas,
			CustomerClassDesc:       custClsDes,

			// Credit Limit
			CreditLimitType:      int32(cust.GnlBusinessType),
			CreditLimitTypeDesc:  creditLimitType,
			CreditLimitAmount:    creditLimitAmount,
			RemainingCreditLimit: cust.RemainingCreditLimit,

			// Adm Division
			Province:    stateCustomer,
			City:        cityCustomer,
			District:    districtCustomer,
			SubDistrict: subDistrictCustomer,
			Zip:         zipCustomer,

			DueDate:              earliestDateStr,
			TotalRemainingAmount: totalRemainingAmountInvoice,
		})
	}

	res = datas[0]

	return
}

func (s *CustomerService) GetListGp(ctx context.Context, req dto.GetCustomerGPListRequest) (res []*dto.CustomerGP, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CustomerService.GetListGp")
	defer span.End()

	var customer *bridgeService.GetCustomerGPResponse

	if customer, err = s.opt.Client.BridgeServiceGrpc.GetCustomerGPList(ctx, &bridgeService.GetCustomerGPListRequest{
		Limit:  int32(req.Limit),
		Offset: int32(req.Offset),
	}); err != nil || !customer.Succeeded {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "customer")
		return
	}

	for _, cust := range customer.Data {
		pymtrmids := []*dto.CustomerGPpymtrmid{}
		for _, pymtrmid := range cust.Pymtrmid {
			pymtrmids = append(pymtrmids, &dto.CustomerGPpymtrmid{
				Pymtrmid:              pymtrmid.Pymtrmid,
				CalculateDateFromDays: pymtrmid.CalculateDateFromDays,
			})
		}
		res = append(res, &dto.CustomerGP{
			Custnmbr:                 cust.Custnmbr,
			Custclas:                 cust.Custclas,
			Custname:                 cust.Custname,
			Cprcstnm:                 cust.Cprcstnm,
			Cntcprsn:                 cust.Cntcprsn,
			Stmtname:                 cust.Stmtname,
			Shrtname:                 cust.Shrtname,
			Upszone:                  cust.Upszone,
			Shipmthd:                 cust.Shipmthd,
			Taxschid:                 cust.Taxschid,
			AddresS1:                 cust.AddresS1,
			AddresS2:                 cust.AddresS2,
			AddresS3:                 cust.AddresS3,
			Country:                  cust.Country,
			City:                     cust.City,
			State:                    cust.State,
			Zip:                      cust.Zip,
			PhonE1:                   cust.PhonE1,
			PhonE2:                   cust.PhonE2,
			PhonE3:                   cust.PhonE3,
			Fax:                      cust.Fax,
			Prbtadcd:                 cust.Prbtadcd,
			Prstadcd:                 cust.Prstadcd,
			Staddrcd:                 cust.Staddrcd,
			Slprsnid:                 cust.Slprsnid,
			Chekbkid:                 cust.Chekbkid,
			Pymtrmid:                 pymtrmids,
			Crlmttyp:                 cust.Crlmttyp,
			Crlmtamt:                 cust.Crlmtamt,
			Curncyid:                 cust.Curncyid,
			Ratetpid:                 cust.Ratetpid,
			Custdisc:                 cust.Custdisc,
			Prclevel:                 cust.Prclevel,
			Minpytyp:                 cust.Minpytyp,
			MinimumPaymentTypeDesc:   cust.MinimumPaymentTypeDesc,
			Minpydlr:                 cust.Minpydlr,
			Minpypct:                 cust.Minpypct,
			Fnchatyp:                 cust.Fnchatyp,
			FinanceChargeAmtTypeDesc: cust.FinanceChargeAmtTypeDesc,
			Fnchpcnt:                 cust.Fnchpcnt,
			Finchdlr:                 cust.Finchdlr,
			Mxwoftyp:                 cust.Mxwoftyp,
			MaximumWriteoffTypeDesc:  cust.MaximumWriteoffTypeDesc,
			Mxwrofam:                 cust.Mxwrofam,
			CommenT1:                 cust.CommenT1,
			CommenT2:                 cust.CommenT2,
			UserdeF1:                 cust.UserdeF1,
			UserdeF2:                 cust.UserdeF2,
			TaxexmT1:                 cust.TaxexmT1,
			TaxexmT2:                 cust.TaxexmT2,
			Txrgnnum:                 cust.Txrgnnum,
			Balnctyp:                 cust.Balnctyp,
			BalanceTypeDesc:          cust.BalanceTypeDesc,
			Stmtcycl:                 cust.Stmtcycl,
			StatementCycleDesc:       cust.StatementCycleDesc,
			Bankname:                 cust.Bankname,
			Bnkbrnch:                 cust.Bnkbrnch,
			Salsterr:                 cust.Salsterr,
			Inactive:                 cust.Inactive,
			Hold:                     cust.Hold,
			Crcardid:                 cust.Crcardid,
			Crcrdnum:                 cust.Crcrdnum,
			Ccrdxpdt:                 cust.Ccrdxpdt,
		})
	}

	total = int64(len(customer.Data))

	return
}

func (s *CustomerService) GetDetailGp(ctx context.Context, id string) (res *dto.CustomerGP, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CustomerService.GetDetailGp")
	defer span.End()

	var customer *bridgeService.GetCustomerGPResponse

	if customer, err = s.opt.Client.BridgeServiceGrpc.GetCustomerGPDetail(ctx, &bridgeService.GetCustomerGPDetailRequest{
		Id: id,
	}); err != nil || !customer.Succeeded {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "customer")
		return
	}

	pymtrmids := []*dto.CustomerGPpymtrmid{}
	for _, pymtrmid := range customer.Data[0].Pymtrmid {
		pymtrmids = append(pymtrmids, &dto.CustomerGPpymtrmid{
			Pymtrmid:              pymtrmid.Pymtrmid,
			CalculateDateFromDays: pymtrmid.CalculateDateFromDays,
		})
	}

	res = &dto.CustomerGP{
		Custnmbr:                 customer.Data[0].Custnmbr,
		Custclas:                 customer.Data[0].Custclas,
		Custname:                 customer.Data[0].Custname,
		Cprcstnm:                 customer.Data[0].Cprcstnm,
		Cntcprsn:                 customer.Data[0].Cntcprsn,
		Stmtname:                 customer.Data[0].Stmtname,
		Shrtname:                 customer.Data[0].Shrtname,
		Upszone:                  customer.Data[0].Upszone,
		Shipmthd:                 customer.Data[0].Shipmthd,
		Taxschid:                 customer.Data[0].Taxschid,
		AddresS1:                 customer.Data[0].AddresS1,
		AddresS2:                 customer.Data[0].AddresS2,
		AddresS3:                 customer.Data[0].AddresS3,
		Country:                  customer.Data[0].Country,
		City:                     customer.Data[0].City,
		State:                    customer.Data[0].State,
		Zip:                      customer.Data[0].Zip,
		PhonE1:                   customer.Data[0].PhonE1,
		PhonE2:                   customer.Data[0].PhonE2,
		PhonE3:                   customer.Data[0].PhonE3,
		Fax:                      customer.Data[0].Fax,
		Prbtadcd:                 customer.Data[0].Prbtadcd,
		Prstadcd:                 customer.Data[0].Prstadcd,
		Staddrcd:                 customer.Data[0].Staddrcd,
		Slprsnid:                 customer.Data[0].Slprsnid,
		Chekbkid:                 customer.Data[0].Chekbkid,
		Pymtrmid:                 pymtrmids,
		Crlmttyp:                 customer.Data[0].Crlmttyp,
		Crlmtamt:                 customer.Data[0].Crlmtamt,
		Curncyid:                 customer.Data[0].Curncyid,
		Ratetpid:                 customer.Data[0].Ratetpid,
		Custdisc:                 customer.Data[0].Custdisc,
		Prclevel:                 customer.Data[0].Prclevel,
		Minpytyp:                 customer.Data[0].Minpytyp,
		MinimumPaymentTypeDesc:   customer.Data[0].MinimumPaymentTypeDesc,
		Minpydlr:                 customer.Data[0].Minpydlr,
		Minpypct:                 customer.Data[0].Minpypct,
		Fnchatyp:                 customer.Data[0].Fnchatyp,
		FinanceChargeAmtTypeDesc: customer.Data[0].FinanceChargeAmtTypeDesc,
		Fnchpcnt:                 customer.Data[0].Fnchpcnt,
		Finchdlr:                 customer.Data[0].Finchdlr,
		Mxwoftyp:                 customer.Data[0].Mxwoftyp,
		MaximumWriteoffTypeDesc:  customer.Data[0].MaximumWriteoffTypeDesc,
		Mxwrofam:                 customer.Data[0].Mxwrofam,
		CommenT1:                 customer.Data[0].CommenT1,
		CommenT2:                 customer.Data[0].CommenT2,
		UserdeF1:                 customer.Data[0].UserdeF1,
		UserdeF2:                 customer.Data[0].UserdeF2,
		TaxexmT1:                 customer.Data[0].TaxexmT1,
		TaxexmT2:                 customer.Data[0].TaxexmT2,
		Txrgnnum:                 customer.Data[0].Txrgnnum,
		Balnctyp:                 customer.Data[0].Balnctyp,
		BalanceTypeDesc:          customer.Data[0].BalanceTypeDesc,
		Stmtcycl:                 customer.Data[0].Stmtcycl,
		StatementCycleDesc:       customer.Data[0].StatementCycleDesc,
		Bankname:                 customer.Data[0].Bankname,
		Bnkbrnch:                 customer.Data[0].Bnkbrnch,
		Salsterr:                 customer.Data[0].Salsterr,
		Inactive:                 customer.Data[0].Inactive,
		Hold:                     customer.Data[0].Hold,
		Crcardid:                 customer.Data[0].Crcardid,
		Crcrdnum:                 customer.Data[0].Crcrdnum,
		Ccrdxpdt:                 customer.Data[0].Ccrdxpdt,
	}

	return
}

func (s *CustomerService) CreateGP(ctx context.Context, req dto.CreateCustomerGPRequest) (res *bridgeService.CreateCustomerGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CustomerService.CreateGP")
	defer span.End()

	// New Customer
	if req.CustNmbr == "" {

		// check salesperson
		var salesPersonDetail *bridgeService.GetSalesPersonGPResponse
		salesPersonDetail, err = s.opt.Client.BridgeServiceGrpc.GetSalesPersonGPDetail(ctx, &bridgeService.GetSalesPersonGPDetailRequest{
			Id: req.SalesPerson,
		})
		if err != nil || len(salesPersonDetail.Data) == 0 {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "Sales person")
			return
		}

		// Image address validation
		if len(req.ImageAddress) > 3 && len(req.ImageAddress) < 1 {
			err = edenlabs.ErrorMustEqualOrLess("address image", "3")
			return
		}

		// check referrer code
		if req.ReferrerCode != "" {
			var customerGP *bridge_service.GetCustomerGPResponse
			customerGP, err = s.opt.Client.BridgeServiceGrpc.GetCustomerGPList(ctx, &bridge_service.GetCustomerGPListRequest{
				ReferralCode: req.ReferrerCode,
				Limit:        1,
				Offset:       0,
				Inactive:     "0",
			})
			if err != nil || len(customerGP.Data) == 0 {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorValidation("bridge", "referrer code")
				return
			}
		}

		// check customer class
		var creditLimitAmount float64
		if req.CustomerClass != "" {
			var custCls *bridge_service.GetCustomerClassResponse
			custCls, err = s.opt.Client.BridgeServiceGrpc.GetCustomerClassList(ctx, &bridge_service.GetCustomerClassListRequest{
				Classid: req.CustomerClass,
			})
			if err != nil || len(custCls.Data) == 0 {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "customer class")
				return
			}
			creditLimitAmount = custCls.Data[0].Crlmtamt
		}

		if len(req.Phone1) > 15 {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorMustEqualOrLess("customer_phone_number", "15")
			return
		}

		if len(req.Phone2) > 15 {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorMustEqualOrLess("customer_alt_phone_number", "15")
			return
		}

		if len(req.AddressAddr1) > 60 {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorMustEqualOrLess("address_addr1", "60")
			return
		}
		if len(req.AddressAddr2) > 60 {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorMustEqualOrLess("address_addr2", "60")
			return
		}
		if len(req.AddressAddr3) > 60 {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorMustEqualOrLess("address_addr3", "60")
			return
		}

		if len(req.CustName) > 60 {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorMustEqualOrLess("customer_name", "60")
			return
		}

		if len(req.AddressName) > 60 {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorMustEqualOrLess("adress_name", "60")
			return
		}

		if len(req.AddressNote) > 150 {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorMustEqualOrLess("address_note", "150")
			return
		}

		if len(req.PICName) > 60 {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorMustEqualOrLess("address_pic_name", "60")
			return
		}

		if len(req.CntcPrsn) > 60 {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorMustEqualOrLess("customer_pic_name", "60")
			return
		}

		bodyCreateCustomer := &bridge_service.CreateCustomerGPRequest{
			Custname:        req.CustName,
			Cprcstnm:        req.Phone1,
			Custpriority:    "1",
			GnlReferrerCode: req.ReferrerCode,
			GnlCustTypeId:   "BTY0015",
			GnlBusinessType: req.BusinessTypeCreditLimit,
			Pymtrmid:        req.PymtrmID,
			Prclevel:        req.PrcLevel,
			Comment1:        "",
			Cntcprsn:        req.PICName,
			Stmtname:        req.CntcPrsn,
			Custclas:        req.CustomerClass,
			Crlmttyp:        req.BusinessTypeCreditLimit,
			Crlmtamt:        creditLimitAmount,
		}

		// Create customer to GP
		var responseCreateCustomer *bridgeService.CreateCustomerGPResponse
		responseCreateCustomer, err = s.opt.Client.BridgeServiceGrpc.CreateCustomerGP(ctx, bodyCreateCustomer)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorValidation("bridge", "Error Create customer GP")
			return
		}

		// Generate Address Code
		codeGenerator, err := s.opt.Client.ConfigurationServiceGrpc.GetGenerateCode(ctx, &configuration_service.GetGenerateCodeRequest{
			Format: responseCreateCustomer.Custnmbr + "-",
			Domain: "address",
			Length: 3,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
		}
		req.CodeAddress = codeGenerator.Data.Code

		// Get adm division from subdistrict selection
		admDivision, err := s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPList(ctx, &bridge_service.GetAdmDivisionGPListRequest{
			Region:      req.Region,
			State:       req.State,
			City:        req.City,
			District:    req.District,
			SubDistrict: req.SubDistrict,
			Limit:       1,
			Offset:      0,
		})
		if err != nil || len(admDivision.Data) == 0 {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorValidation("bridge", "sub_district")
		}

		if len(req.AddressAddr1) > 60 {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "lebih dari 60")
			return nil, err
		}

		bodyCreateAddress := &bridge_service.CreateAddressRequest{
			Custnmbr:                responseCreateCustomer.Custnmbr,
			Adrscode:                req.CodeAddress,
			AddresS1:                req.AddressAddr1,
			AddresS2:                req.AddressAddr2,
			AddresS3:                req.AddressAddr3,
			Cntcprsn:                req.CustName,
			State:                   req.State,
			City:                    req.City,
			PhonE1:                  req.Phone1,
			PhonE2:                  req.Phone2,
			Zip:                     req.Zip,
			CCode:                   "ID",
			Country:                 "Indonesia",
			GnL_Address_Note:        req.AddressNote,
			ShipToName:              req.AddressName,
			GnL_Archetype_ID:        req.Archetype,
			GnL_Administrative_Code: admDivision.Data[0].Code,
			Slprsnid:                salesPersonDetail.Data[0].Slprsnid,
			Salsterr:                salesPersonDetail.Data[0].Salsterr,
			Shipmthd:                "DELIVERY",
			Locncode:                req.Site,
		}

		_, err = s.opt.Client.BridgeServiceGrpc.CreateAddress(ctx, bodyCreateAddress)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("Bridge", "Update address GP")
		}

		bodyUpdateCustomer := &bridge_service.UpdateCustomerGPRequest{
			Custnmbr: responseCreateCustomer.Custnmbr,
			Prstadcd: bodyCreateAddress.Adrscode,
			Prbtadcd: bodyCreateAddress.Adrscode,
			Staddrcd: bodyCreateAddress.Adrscode,
			Slprsnid: salesPersonDetail.Data[0].Slprsnid,
			Salsterr: salesPersonDetail.Data[0].Salsterr,
			Shipmthd: "DELIVERY",
			Address: &bridgeService.UpdateAddressRequest{
				Custnmbr:                responseCreateCustomer.Custnmbr,
				Adrscode:                req.CodeAddress,
				AddresS1:                req.AddressAddr1,
				AddresS2:                req.AddressAddr2,
				AddresS3:                req.AddressAddr3,
				Cntcprsn:                req.CustName,
				State:                   req.State,
				City:                    req.City,
				PhonE1:                  req.Phone1,
				PhonE2:                  req.Phone2,
				Zip:                     req.Zip,
				CCode:                   "ID",
				Country:                 "Indonesia",
				GnL_Address_Note:        req.AddressNote,
				ShipToName:              req.AddressName,
				GnL_Archetype_ID:        req.Archetype,
				GnL_Administrative_Code: admDivision.Data[0].Code,
				Slprsnid:                salesPersonDetail.Data[0].Slprsnid,
				Salsterr:                salesPersonDetail.Data[0].Salsterr,
				Shipmthd:                "DELIVERY",
				Locncode:                req.Site,
			},
		}

		// Update customer to GP
		_, err = s.opt.Client.BridgeServiceGrpc.UpdateCustomerGP(ctx, bodyUpdateCustomer)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorValidation("bridge", "Error Update customer GP")
		}

		// Add customer to Eden table
		reqCreateRequestCRM := &crm_service.CreateCustomerRequest{
			Data: &crm_service.Customer{
				CustomerIdGp:            responseCreateCustomer.Custnmbr,
				ProspectiveCustomerId:   0,
				MembershipLevelId:       0,
				MembershipCheckpointId:  0,
				TotalPoint:              0,
				ProfileCode:             responseCreateCustomer.Custnmbr,
				Email:                   req.Email,
				ReferenceInfo:           "",
				UpgradeStatus:           0,
				KtpPhotosUrl:            req.ImageKtp,
				CustomerPhotosUrl:       strings.Join(req.ImageAddress, ","),
				CustomerSelfieUrl:       "",
				MembershipRewardId:      0,
				MembershipRewardAmmount: 0,
				ReferralCode:            responseCreateCustomer.GnLReferralCode,
				ReferrerCode:            req.ReferrerCode,
				Gender:                  0,
				BirthDate:               "",
			},
		}

		_, err = s.opt.Client.CrmServiceGrpc.CreateCustomer(ctx, reqCreateRequestCRM)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
		}

	} else { // if customer already exist

		// check customer
		var customerDetail *bridgeService.GetCustomerGPResponse
		customerDetail, err = s.opt.Client.BridgeServiceGrpc.GetCustomerGPDetail(ctx, &bridgeService.GetCustomerGPDetailRequest{
			Id: req.CustNmbr,
		})
		if err != nil || len(customerDetail.Data) == 0 {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "customer")
			return
		}

		// generate customer code
		codeGenerator, err := s.opt.Client.ConfigurationServiceGrpc.GetGenerateCode(ctx, &configuration_service.GetGenerateCodeRequest{
			Format: customerDetail.Data[0].Custnmbr + "-",
			Domain: "address",
			Length: 3,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
		}

		// check salesperson
		var salesPersonDetail *bridgeService.GetSalesPersonGPResponse
		salesPersonDetail, err = s.opt.Client.BridgeServiceGrpc.GetSalesPersonGPDetail(ctx, &bridgeService.GetSalesPersonGPDetailRequest{
			Id: req.SalesPerson,
		})
		if err != nil || len(salesPersonDetail.Data) == 0 {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "Sales person")
		}

		// check adm division
		admDivision, err := s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPList(ctx, &bridge_service.GetAdmDivisionGPListRequest{
			State:       req.State,
			City:        req.City,
			District:    req.District,
			SubDistrict: req.SubDistrict,
			Limit:       1,
			Offset:      0,
		})
		if err != nil || len(admDivision.Data) == 0 {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorValidation("bridge", "sub_district")
		}

		if len(req.AddressAddr1) > 60 {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorMustEqualOrLess("address_addr1", "60")
		}
		if len(req.AddressAddr2) > 60 {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorMustEqualOrLess("address_addr2", "60")
		}
		if len(req.AddressAddr3) > 60 {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorMustEqualOrLess("address_addr3", "60")
		}

		if len(req.AddressName) > 60 {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorMustEqualOrLess("adress_name", "60")
		}

		if len(req.AddressNote) > 150 {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorMustEqualOrLess("address_note", "150")
		}

		if len(req.PICName) > 60 {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorMustEqualOrLess("address_pic_name", "60")
		}

		_, err = s.opt.Client.BridgeServiceGrpc.CreateAddress(ctx, &bridge_service.CreateAddressRequest{
			Custnmbr:                customerDetail.Data[0].Custnmbr,
			Custname:                customerDetail.Data[0].Custname,
			Cntcprsn:                req.CustName,
			ShipToName:              req.AddressName,
			Adrscode:                codeGenerator.Data.Code,
			AddresS1:                req.AddressAddr1,
			AddresS2:                req.AddressAddr2,
			AddresS3:                req.AddressAddr3,
			GnL_Address_Note:        req.AddressNote,
			Country:                 "Indonesia",
			City:                    admDivision.Data[0].City,
			State:                   admDivision.Data[0].State,
			PhonE1:                  req.Phone1,
			PhonE2:                  req.Phone2,
			Inactive:                "0",
			GnL_Administrative_Code: admDivision.Data[0].Code,
			Locncode:                req.Site,
			TypeAddress:             "other",
			GnL_Archetype_ID:        req.Archetype,
			Salsterr:                salesPersonDetail.Data[0].Salsterr,
			Slprsnid:                salesPersonDetail.Data[0].Slprsnid,
		})

		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
		}

	}

	// userID := ctx.Value(constants.KeyUserID).(int64)
	// _, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &auditService.CreateLogRequest{
	// 	Log: &auditService.Log{
	// 		UserId:      userID,
	// 		ReferenceId: "0",
	// 		Type:        "purchase_order",
	// 		Function:    "create",
	// 		CreatedAt:   timestamppb.New(time.Now()),
	// 	},
	// })
	// if err != nil {
	// 	span.RecordError(err)
	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 	err = edenlabs.ErrorRpc("audit")
	// 	return
	// }

	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "customer")
		return
	}

	return
}

func (s *CustomerService) GetOverdueMitra(ctx context.Context, req dto.CustomerListRequest) (res []*dto.CustomerResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CustomerService.GetCustomers")
	defer span.End()

	// Get sales invoice from bridge
	var siRes *bridgeService.GetSalesInvoiceGPListResponse

	siRes, err = s.opt.Client.BridgeServiceGrpc.GetSalesInvoiceGPList(ctx, &bridgeService.GetSalesInvoiceGPListRequest{
		Limit:               100,
		Offset:              0,
		OrderBy:             "desc",
		RemainingAmountFlag: "1",
	})
	if err != nil {
		fmt.Println("err", err)
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "sales invoice")
		return
	}

	uniqueCustomers := make(map[string]bool)
	uniqueCustomerData := []*bridgeService.SalesInvoiceGP{}

	for _, si := range siRes.Data {
		if _, exists := uniqueCustomers[si.Custnmbr]; !exists {
			uniqueCustomers[si.Custnmbr] = true
			uniqueCustomerData = append(uniqueCustomerData, si)
		}
	}

	datas := []*dto.CustomerResponse{}

	for _, data := range uniqueCustomerData {
		var custResponse *bridgeService.GetCustomerGPResponse

		custResponse, err = s.opt.Client.BridgeServiceGrpc.GetCustomerGPList(ctx, &bridgeService.GetCustomerGPListRequest{
			CustomerTypeId: "BTY0015",
			Id:             data.Custnmbr,
			Limit:          req.Limit,
			Offset:         req.Offset,
			Inactive:       "0",
		})
		if err != nil {
			continue
		}

		// Get sales invoice based on customer
		var (
			siRes                       *bridgeService.GetSalesInvoiceGPListResponse
			totalRemainingAmountInvoice float64
			earliestDate                time.Time
			earliestDateStr             string
		)

		siRes, err = s.opt.Client.BridgeServiceGrpc.GetSalesInvoiceGPList(ctx, &bridgeService.GetSalesInvoiceGPListRequest{
			Limit:         20,
			Offset:        0,
			OrderBy:       "desc",
			Custnumber:    data.Custnmbr,
			GnlCustTypeId: "BTY0015",
		})
		if err != nil {
			continue
		}

		if len(siRes.Data) > 0 {
			for _, si := range siRes.Data {
				if si.RemainingAmount > 0 {
					totalRemainingAmountInvoice += si.RemainingAmount
				}
				date, _ := time.Parse("2006-01-02T15:04:05Z", si.Duedate)
				if earliestDate.IsZero() || date.Before(earliestDate) {
					earliestDate = date
					earliestDateStr = earliestDate.Format("2006-01-02")
				}
			}
		}

		for _, cust := range custResponse.Data {
			datas = append(datas, &dto.CustomerResponse{
				ID:             cust.Custnmbr,
				Name:           cust.Custname,
				Code:           cust.Custnmbr,
				PhoneNumber:    cust.PhonE1,
				AltPhoneNumber: cust.PhonE2,
				BillingAddress: cust.AddresS1 + " " + cust.AddresS2 + " " + cust.AddresS3,
				// Status:            statusGP,
				// StatusDescription: statusDesc,
				ReferralCode: cust.GnlReferralCode,
				ReferrerCode: cust.GnlReferrerCode,
				// PaymentTerm:       tempPaymentTerm,
				// CustomerTypeID:    custType,
				// CustomerTypeDesc:  custTypeDesc,

				// Additional
				// Region:            regionCustomer,
				ShippingAddress: cust.AddresS1,
				OverdueDebt:     float64(cust.RemainingAmount),
				CustomerGroup:   "B2B",
				// PaymentGroup:      paymentGroupFinal,
				PicName:       cust.Stmtname,
				CustomerClass: cust.Custclas,
				// CustomerClassDesc: custClsDes,
				// Email:             email,
				// SalesPerson:       salesPersonId,
				// SalesPersonName:   salesPersonName,
				// SiteCode:          siteCode,
				// SiteName:          siteName,

				// Credit Limit
				CreditLimitType: int32(cust.GnlBusinessType),
				// CreditLimitTypeDesc:  creditLimitType,
				// CreditLimitAmount:    creditLimitAmount,
				RemainingCreditLimit: cust.RemainingCreditLimit,

				DueDate:              earliestDateStr,
				TotalRemainingAmount: totalRemainingAmountInvoice,
			})
		}
		total = int64(len(custResponse.Data))

		res = datas
	}

	return
}
