package service

import (
	"context"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/edenlabs/edenlabs/utils"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/repository"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	configurationService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/crm_service"
)

type ICustomerService interface {
	GetDetail(ctx context.Context, req *dto.CustomerRequestGetDetail) (res *dto.CustomerResponseGet, err error)
	Update(ctx context.Context, req *dto.CustomerRequestUpdate) (err error)
	Get(ctx context.Context, req *dto.CustomerGetListRequest) (res []*dto.CustomerResponseGet, total int64, err error)
	Create(ctx context.Context, req *crm_service.CreateCustomerRequest) (customer *model.Customer, err error)
	GetDetailComplex(ctx context.Context, req *dto.CustomerRequestGetDetail) (res *dto.CustomerResponseGet, err error)
	GetCustomerID(ctx context.Context, req *crm_service.GetCustomerIDRequest) (res []*model.Customer, err error)
}

type CustomerService struct {
	opt                opt.Options
	RepositoryCustomer repository.ICustomerRepository
}

func NewCustomerService() ICustomerService {
	return &CustomerService{
		opt:                global.Setup.Common,
		RepositoryCustomer: repository.NewCustomerRepository(),
	}
}

func (s *CustomerService) GetDetail(ctx context.Context, req *dto.CustomerRequestGetDetail) (res *dto.CustomerResponseGet, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "Customer.GetDetail")
	defer span.End()

	var customer *model.Customer
	customer, err = s.RepositoryCustomer.GetDetail(ctx, req)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &dto.CustomerResponseGet{
		ID:                      customer.ID,
		Code:                    customer.CustomerIDGP,
		ProspectiveCustomerID:   customer.ProspectiveCustomerID,
		MembershipLevelID:       customer.MembershipLevelID,
		MembershipCheckpointID:  customer.MembershipCheckpointID,
		TotalPoint:              customer.TotalPoint,
		ProfileCode:             customer.ProfileCode,
		Email:                   customer.Email,
		ReferenceInfo:           customer.ReferenceInfo,
		UpgradeStatus:           customer.UpgradeStatus,
		KtpPhotosUrl:            customer.KtpPhotosUrl,
		CustomerPhotosUrl:       customer.CustomerPhotosUrl,
		CustomerSelfieUrl:       customer.CustomerSelfieUrl,
		CreatedAt:               customer.CreatedAt,
		UpdatedAt:               customer.UpdatedAt,
		MembershipRewardID:      customer.MembershipRewardID,
		MembershipRewardAmmount: customer.MembershipRewardAmmount,
		ReferrerID:              customer.ReferrerID,
		ReferrerCode:            customer.ReferrerCode,
		ReferralCode:            customer.ReferralCode,
	}

	return
}

func (s *CustomerService) Update(ctx context.Context, req *dto.CustomerRequestUpdate) (err error) {
	ctx, span := s.opt.Trace.Start(ctx, "Customer.GetDetail")
	defer span.End()

	param := &dto.CustomerRequestGetDetail{
		ID:           req.ID,
		CustomerIDGP: req.CustomerIDGP,
	}

	var customer *model.Customer
	customer, err = s.RepositoryCustomer.GetDetail(ctx, param)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	customer = &model.Customer{
		ID:                     customer.ID,
		CustomerIDGP:           req.CustomerIDGP,
		ProspectiveCustomerID:  req.ProspectiveCustomerID,
		MembershipLevelID:      req.MembershipLevelID,
		MembershipCheckpointID: req.MembershipCheckpointID,
		TotalPoint:             req.TotalPoint,
		ProfileCode:            req.ProfileCode,
		ReferenceInfo:          req.ReferenceInfo,
		UpgradeStatus:          req.UpgradeStatus,
		UpdatedAt:              time.Now(),
	}

	err = s.RepositoryCustomer.Update(ctx, customer, req.FieldUpdate...)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *CustomerService) Get(ctx context.Context, req *dto.CustomerGetListRequest) (res []*dto.CustomerResponseGet, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "Customer.Get")
	defer span.End()

	var (
		statusGP, customerTypeID string
	)

	if req.Status != 0 {
		switch req.Status {
		case 1:
			statusGP = "0"
		case 7:
			statusGP = "1"
		default:
			statusGP = utils.ToString(req.Status)
		}
	}

	if req.CustomerType != "" {
		var customerType *bridge_service.GetCustomerTypeGPResponse
		customerType, err = s.opt.Client.BridgeServiceGrpc.GetCustomerTypeGPList(ctx, &bridge_service.GetCustomerTypeGPListRequest{
			Limit:       int32(req.Limit),
			Offset:      int32(req.Offset),
			Description: req.CustomerType,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		if len(customerType.Data) != 0 {
			customerTypeID = customerType.Data[0].GnL_Cust_Type_ID
		} else {
			return
		}
	}

	var customerGP *bridge_service.GetCustomerGPResponse
	customerGP, err = s.opt.Client.BridgeServiceGrpc.GetCustomerGPList(ctx, &bridge_service.GetCustomerGPListRequest{
		Limit:          int32(req.Limit),
		Offset:         int32(req.Offset),
		Inactive:       statusGP,
		Name:           req.Search,
		CustomerTypeId: customerTypeID,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, v := range customerGP.Data {
		var customerTypeResponse *dto.CustomerTypeResponse
		if v.CustomerType != nil {
			customerTypeResponse = &dto.CustomerTypeResponse{
				ID:          v.CustomerType[0].GnL_Cust_Type_ID,
				Code:        v.CustomerType[0].GnL_Cust_Type_ID,
				Description: v.CustomerType[0].GnL_CustType_Description,
			}

		}

		customer := &model.Customer{
			CustomerIDGP: v.Custnmbr,
			CreatedAt:    time.Now(),
		}
		if err = s.RepositoryCustomer.SyncGP(ctx, customer); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("customer sync")
			return
		}

		res = append(res, &dto.CustomerResponseGet{
			ID:                         customer.ID,
			Code:                       customer.CustomerIDGP,
			ProspectiveCustomerID:      customer.ProspectiveCustomerID,
			MembershipLevelID:          customer.MembershipLevelID,
			MembershipCheckpointID:     customer.MembershipCheckpointID,
			TotalPoint:                 customer.TotalPoint,
			ProfileCode:                customer.ProfileCode,
			Email:                      customer.Email,
			ReferenceInfo:              customer.ReferenceInfo,
			UpgradeStatus:              customer.UpgradeStatus,
			KtpPhotosUrl:               customer.KtpPhotosUrl,
			CustomerPhotosUrl:          customer.CustomerPhotosUrl,
			CustomerSelfieUrl:          customer.CustomerSelfieUrl,
			CreatedAt:                  customer.CreatedAt,
			UpdatedAt:                  customer.UpdatedAt,
			MembershipRewardID:         customer.MembershipRewardID,
			MembershipRewardAmmount:    customer.MembershipRewardAmmount,
			Name:                       v.Custname,
			PhonE1:                     v.PhonE1,
			CorporateCustomerNumber:    v.Cprcstnm,
			ContactPerson:              v.Cntcprsn,
			StatementName:              v.Stmtname,
			ShortName:                  v.Shrtname,
			Upszone:                    v.Upszone,
			TaxScheduleID:              v.Taxschid,
			AddresS1:                   v.AddresS1,
			AddresS2:                   v.AddresS2,
			AddresS3:                   v.AddresS3,
			Country:                    v.Country,
			City:                       v.City,
			State:                      v.State,
			Zip:                        v.Zip,
			PhonE2:                     v.PhonE2,
			PhonE3:                     v.PhonE3,
			Fax:                        v.Fax,
			PrimaryAddressCode:         v.Prbtadcd,
			StreetAddressCode:          v.Staddrcd,
			SalesPersonID:              v.Slprsnid,
			CheckbookID:                v.Chekbkid,
			CreditLimitType:            v.Crlmttyp,
			CreditLimitTypeDescription: v.CreditLimitTypeDesc,
			CreditLimitAmount:          v.Crlmtamt,
			CurrencyID:                 v.Curncyid,
			RateTypeID:                 v.Ratetpid,
			CustomerDiscount:           v.Custdisc,
			MinimumPaymentType:         v.Minpytyp,
			MinimumPaymentTypeDesc:     v.MinimumPaymentTypeDesc,
			MinimumPaymentDollarAmount: v.Minpydlr,
			MinimumPaymentPercent:      v.Minpypct,
			FinanceChargeType:          v.Fnchatyp,
			FinanceChargeTypeDesc:      v.FinanceChargeAmtTypeDesc,
			FinanceChargePercent:       v.Fnchpcnt,
			FinanceChargeDollarAmount:  v.Finchdlr,
			MaximumWriteoffType:        v.Mxwoftyp,
			MaximumWriteoffTypeDesc:    v.MaximumWriteoffTypeDesc,
			MaximumWriteoffAmount:      v.Mxwrofam,
			Comment1:                   v.CommenT1,
			Comment2:                   v.CommenT2,
			UserDefined1:               v.UserdeF1,
			UserDefined2:               v.UserdeF2,
			TaxExempt1:                 v.TaxexmT1,
			TaxExempt2:                 v.TaxexmT2,
			TaxRegistrationNumber:      v.Txrgnnum,
			BalanceType:                v.Balnctyp,
			BalanceTypeDesc:            v.BalanceTypeDesc,
			StatementCycle:             v.Stmtcycl,
			StatementCycleDesc:         v.StatementCycleDesc,
			BankName:                   v.Bankname,
			BankBranch:                 v.Bnkbrnch,
			Hold:                       v.Hold,
			CreditCardID:               v.Crcardid,
			CreditCardNumber:           v.Crcrdnum,
			CreditCardExpDate:          v.Ccrdxpdt,
			Inactive:                   v.Inactive,
			ReferrerID:                 customer.ReferrerID,
			ReferrerCode:               v.GnlReferrerCode,
			ReferralCode:               v.GnlReferralCode,
			CustomerType:               customerTypeResponse,
		})
	}
	total = int64(len(res))

	return
}

func (s *CustomerService) Create(ctx context.Context, req *crm_service.CreateCustomerRequest) (customer *model.Customer, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "Customer.GetDetail")
	defer span.End()

	var birthDate time.Time
	if req.Data.BirthDate != "" {
		layout := "2006-01-02"
		birthDate, err = time.Parse(layout, req.Data.BirthDate)
		if err != nil {
			return
		}
	}

	customer = &model.Customer{
		CustomerIDGP:           req.Data.CustomerIdGp,
		ProspectiveCustomerID:  req.Data.ProspectiveCustomerId,
		MembershipLevelID:      req.Data.MembershipLevelId,
		MembershipCheckpointID: req.Data.MembershipCheckpointId,
		TotalPoint:             req.Data.TotalPoint,
		ProfileCode:            req.Data.ProfileCode,
		Email:                  req.Data.Email,
		// Password:                "",
		ReferenceInfo:           req.Data.ReferenceInfo,
		UpgradeStatus:           int8(req.Data.UpgradeStatus),
		KtpPhotosUrl:            req.Data.KtpPhotosUrl,
		CustomerPhotosUrl:       req.Data.CustomerPhotosUrl,
		CustomerSelfieUrl:       req.Data.CustomerSelfieUrl,
		CreatedAt:               time.Now(),
		UpdatedAt:               time.Now(),
		MembershipRewardID:      req.Data.MembershipRewardId,
		MembershipRewardAmmount: req.Data.MembershipRewardAmmount,
		ReferrerID:              req.Data.ReferrerId,
		ReferrerCode:            req.Data.ReferrerCode,
		ReferralCode:            req.Data.ReferralCode,
		Gender:                  int(req.Data.Gender),
		BirthDate:               birthDate,
	}

	cust, err := s.RepositoryCustomer.Create(ctx, customer)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	customer.ID = cust.ID

	return
}

// This function be created cause performance considering
func (s *CustomerService) GetDetailComplex(ctx context.Context, req *dto.CustomerRequestGetDetail) (res *dto.CustomerResponseGet, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "Customer.Get")
	defer span.End()

	var (
		customerGP                                   *bridge_service.GetCustomerGPResponse
		salespersonResponse                          *dto.SalespersonResponse
		archetypeResponse                            *dto.ArchetypeResponse
		customerTypeResponse                         *dto.CustomerTypeResponse
		companyAddress, shipToAddress, billToAddress *dto.AddressCustomerResponse
		paymentTermResponse                          *dto.PaymentTermResponse
		shippingMethodResponse                       *dto.ShippingMethodResponse
		priceLevelResponse                           *dto.PriceLevelResponse
		salesTerritoryResponse                       *dto.SalesTerritoryResponse
		siteResponse                                 *dto.SiteResponse
		customerClassResponse                        *dto.CustomerClassResponse
		businessTypeResponse                         *dto.GlossaryResponse
		// timeConcentResponse, referenceInfoResponse, invoiceTermResponse *dto.GlossaryResponse
		statusArchetype, statusCustomerType int8
	)

	var customer *model.Customer
	customer, err = s.RepositoryCustomer.GetDetail(ctx, req)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	customerGP, err = s.opt.Client.BridgeServiceGrpc.GetCustomerGPList(ctx, &bridge_service.GetCustomerGPListRequest{
		Id:     customer.CustomerIDGP,
		Limit:  1,
		Offset: 0,
	})
	if err != nil || len(customerGP.Data) == 0 {
		err = edenlabs.ErrorInvalid("customer_id_gp")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if customerGP.Data[0].Slprsnid != "" {
		var salesPerson *bridgeService.GetSalesPersonGPResponse
		salesPerson, err = s.opt.Client.BridgeServiceGrpc.GetSalesPersonGPDetail(ctx, &bridgeService.GetSalesPersonGPDetailRequest{
			Id: customerGP.Data[0].Slprsnid,
		})
		if err != nil {
			err = edenlabs.ErrorInvalid("sales_person_id")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		salespersonResponse = &dto.SalespersonResponse{
			ID:   salesPerson.Data[0].Slprsnid,
			Code: salesPerson.Data[0].Slprsnid,
			Name: salesPerson.Data[0].Slprsnfn + " " + salesPerson.Data[0].Sprsnsmn + " " + salesPerson.Data[0].Sprsnsln,
		}
	}

	if len(customerGP.Data[0].CustomerType) != 0 {
		if customerGP.Data[0].CustomerType[0].Inactive == 0 {
			statusCustomerType = statusx.ConvertStatusName(statusx.Active)
		} else {
			statusCustomerType = statusx.ConvertStatusName(statusx.Archived)
		}

		customerTypeResponse = &dto.CustomerTypeResponse{
			ID:            customerGP.Data[0].CustomerType[0].GnL_Cust_Type_ID,
			Code:          customerGP.Data[0].CustomerType[0].GnL_Cust_Type_ID,
			Description:   customerGP.Data[0].CustomerType[0].GnL_CustType_Description,
			Status:        statusCustomerType,
			ConvertStatus: statusx.ConvertStatusValue(statusCustomerType),
		}
	}

	var addressList *bridgeService.GetAddressGPResponse
	addressList, err = s.opt.Client.BridgeServiceGrpc.GetAddressGPList(ctx, &bridgeService.GetAddressGPListRequest{
		Limit:          10,
		Offset:         0,
		CustomerNumber: customer.CustomerIDGP,
	})
	if err != nil {
		err = edenlabs.ErrorInvalid("customer_id_gp")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, address := range addressList.Data {
		var admDivision *bridgeService.GetAdmDivisionGPResponse
		admDivision, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPList(ctx, &bridgeService.GetAdmDivisionGPListRequest{
			AdmDivisionCode: address.AdministrativeDiv.GnlAdministrativeCode,
			Limit:           1,
			Offset:          0,
		})

		if err != nil || len(admDivision.Data) == 0 {
			err = edenlabs.ErrorInvalid("adm_division_id")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		if address.TypeAddress == "statement_to" {
			companyAddress = &dto.AddressCustomerResponse{
				ID:            address.Adrscode,
				AddressName:   address.ShipToName,
				AddressType:   address.TypeAddress,
				Address1:      address.AddresS1,
				Address2:      address.AddresS2,
				Address3:      address.AddresS3,
				AdmDivisionID: admDivision.Data[0].Code,
				Region:        admDivision.Data[0].Region,
				Province:      admDivision.Data[0].State,
				City:          admDivision.Data[0].City,
				District:      admDivision.Data[0].District,
				SubDistrict:   admDivision.Data[0].Subdistrict,
				PostalCode:    admDivision.Data[0].Zipcode,
			}
		}

		if address.TypeAddress == "ship_to" {
			shipToAddress = &dto.AddressCustomerResponse{
				ID:            address.Adrscode,
				AddressName:   address.ShipToName,
				AddressType:   address.TypeAddress,
				Address1:      address.AddresS1,
				Address2:      address.AddresS2,
				Address3:      address.AddresS3,
				AdmDivisionID: admDivision.Data[0].Code,
				Region:        admDivision.Data[0].Region,
				Province:      admDivision.Data[0].State,
				City:          admDivision.Data[0].City,
				District:      admDivision.Data[0].District,
				SubDistrict:   admDivision.Data[0].Subdistrict,
				PostalCode:    admDivision.Data[0].Zipcode,
			}

			if address.Locncode != "" {
				var site *bridgeService.GetSiteGPResponse
				site, err = s.opt.Client.BridgeServiceGrpc.GetSiteGPDetail(ctx, &bridgeService.GetSiteGPDetailRequest{
					Id: address.Locncode,
				})
				if err != nil || len(site.Data) == 0 {
					span.RecordError(err)
					s.opt.Logger.AddMessage(log.ErrorLevel, err)
					return
				}

				siteResponse = &dto.SiteResponse{
					ID:          site.Data[0].Locncode,
					Description: site.Data[0].Locndscr,
					Type:        site.Data[0].GnlSiteTypeId,
				}
			}
			if address.GnL_Archetype_ID != "" {

				var archetype *bridgeService.GetArchetypeGPResponse
				archetype, err = s.opt.Client.BridgeServiceGrpc.GetArchetypeGPDetail(ctx, &bridgeService.GetArchetypeGPDetailRequest{
					Id: address.GnL_Archetype_ID,
				})
				if err != nil {
					err = edenlabs.ErrorInvalid("archetype_id")
					span.RecordError(err)
					s.opt.Logger.AddMessage(log.ErrorLevel, err)
					return
				}

				if archetype.Data[0].Inactive == 0 {
					statusArchetype = statusx.ConvertStatusName(statusx.Active)
				} else {
					statusArchetype = statusx.ConvertStatusName(statusx.Archived)
				}

				archetypeResponse = &dto.ArchetypeResponse{
					ID:             archetype.Data[0].GnlArchetypeId,
					Code:           archetype.Data[0].GnlArchetypeId,
					Description:    archetype.Data[0].GnlArchetypedescription,
					CustomerTypeID: archetype.Data[0].GnlCustTypeId,
					Status:         statusArchetype,
					ConvertStatus:  statusx.ConvertStatusValue(statusArchetype),
				}
			}
		}

		if address.TypeAddress == "bill_to" {
			billToAddress = &dto.AddressCustomerResponse{
				ID:            address.Adrscode,
				AddressName:   address.ShipToName,
				AddressType:   address.TypeAddress,
				Address1:      address.AddresS1,
				Address2:      address.AddresS2,
				Address3:      address.AddresS3,
				AdmDivisionID: admDivision.Data[0].Code,
				Region:        admDivision.Data[0].Region,
				Province:      admDivision.Data[0].State,
				City:          admDivision.Data[0].City,
				District:      admDivision.Data[0].District,
				SubDistrict:   admDivision.Data[0].Subdistrict,
				PostalCode:    admDivision.Data[0].Zipcode,
			}
		}
	}

	if customerGP.Data[0].Pymtrmid != nil {
		var paymentTerm *bridgeService.GetPaymentTermGPResponse
		paymentTerm, err = s.opt.Client.BridgeServiceGrpc.GetPaymentTermGPDetail(ctx, &bridgeService.GetPaymentTermGPDetailRequest{
			Id: customerGP.Data[0].Pymtrmid[0].Pymtrmid,
		})
		if err != nil || len(paymentTerm.Data) == 0 {
			err = edenlabs.ErrorInvalid("payment_term_id")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		paymentTermResponse = &dto.PaymentTermResponse{
			ID:                       paymentTerm.Data[0].Pymtrmid,
			Description:              paymentTerm.Data[0].Pymtrmid,
			DueType:                  paymentTerm.Data[0].DuetypeDesc,
			PaymentUseFor:            int(paymentTerm.Data[0].GnlPaymentUsefor),
			PaymentUseForDescription: paymentTerm.Data[0].GnlPaymentUseforDesc,
		}
	}

	if customerGP.Data[0].Shipmthd != "" {
		var shippingMethod *bridgeService.GetShippingMethodResponse
		shippingMethod, err = s.opt.Client.BridgeServiceGrpc.GetShippingMethodDetail(ctx, &bridgeService.GetShippingMethodDetailRequest{
			Id: customerGP.Data[0].Shipmthd,
		})
		if err != nil || len(shippingMethod.Data) == 0 {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		shippingMethodResponse = &dto.ShippingMethodResponse{
			ID:              shippingMethod.Data[0].Shipmthd,
			Description:     shippingMethod.Data[0].Shmthdsc,
			Type:            int8(shippingMethod.Data[0].Shiptype),
			TypeDescription: shippingMethod.Data[0].ShiptypeDesc,
		}
	}

	if customerGP.Data[0].Prclevel != "" {
		var priceLevel *bridgeService.GetSalesPriceLevelResponse
		priceLevel, err = s.opt.Client.BridgeServiceGrpc.GetSalesPriceLevelDetail(ctx, &bridgeService.GetSalesPriceLevelDetailRequest{
			Id: customerGP.Data[0].Prclevel,
		})
		if err != nil || len(priceLevel.Data) == 0 {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		priceLevelResponse = &dto.PriceLevelResponse{
			ID:             priceLevel.Data[0].Prclevel,
			Description:    priceLevel.Data[0].Prclevel,
			CustomerTypeID: priceLevel.Data[0].GnlCustTypeId,
			RegionID:       priceLevel.Data[0].GnlRegion,
		}
	}

	if customerGP.Data[0].Salsterr != "" {
		var salesTerritory *bridgeService.GetSalesTerritoryGPResponse
		salesTerritory, err = s.opt.Client.BridgeServiceGrpc.GetSalesTerritoryGPDetail(ctx, &bridgeService.GetSalesTerritoryGPDetailRequest{
			Id: customerGP.Data[0].Salsterr,
		})
		if err != nil || len(salesTerritory.Data) == 0 {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		layout := "2006-01-02T15:04:05"
		createdAt, _ := time.Parse(layout, salesTerritory.Data[0].Creatddt)
		updatedAt, _ := time.Parse(layout, salesTerritory.Data[0].Modifdt)

		salesTerritoryResponse = &dto.SalesTerritoryResponse{
			ID:            salesTerritory.Data[0].Salsterr,
			Description:   salesTerritory.Data[0].Slterdsc,
			SalespersonID: salesTerritory.Data[0].Slprsnid,
			CreatedAt:     createdAt,
			UpdatedAt:     updatedAt,
		}
	}

	if customerGP.Data[0].Custclas != "" {
		var customerClass *bridgeService.GetCustomerClassResponse
		customerClass, err = s.opt.Client.BridgeServiceGrpc.GetCustomerClassDetail(ctx, &bridgeService.GetCustomerClassDetailRequest{
			Id: customerGP.Data[0].Custclas,
		})
		if err != nil || len(customerClass.Data) == 0 {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		customerClassResponse = &dto.CustomerClassResponse{
			ID:          customerClass.Data[0].Classid,
			Code:        customerClass.Data[0].Classid,
			Description: customerClass.Data[0].Clasdscr,
		}
	}

	if customerGP.Data[0].GnlBusinessType != 0 {
		var businessType *configurationService.GetGlossaryDetailResponse
		businessType, err = s.opt.Client.ConfigurationServiceGrpc.GetGlossaryDetail(ctx, &configurationService.GetGlossaryDetailRequest{
			Table:     "customer",
			Attribute: "business_type",
			ValueInt:  int32(customerGP.Data[0].GnlBusinessType),
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		businessTypeResponse = &dto.GlossaryResponse{
			ID:        int64(businessType.Data.Id),
			Table:     businessType.Data.Table,
			Attribute: businessType.Data.Attribute,
			ValueInt:  int8(businessType.Data.ValueInt),
			ValueName: businessType.Data.ValueName,
			Note:      businessType.Data.Note,
		}
	}

	res = &dto.CustomerResponseGet{
		ID:                         customer.ID,
		Code:                       customer.CustomerIDGP,
		ProspectiveCustomerID:      customer.ProspectiveCustomerID,
		MembershipLevelID:          customer.MembershipLevelID,
		MembershipCheckpointID:     customer.MembershipCheckpointID,
		TotalPoint:                 customer.TotalPoint,
		ProfileCode:                customer.ProfileCode,
		Email:                      customer.Email,
		ReferenceInfo:              customer.ReferenceInfo,
		UpgradeStatus:              customer.UpgradeStatus,
		KtpPhotosUrl:               customer.KtpPhotosUrl,
		CustomerPhotosUrl:          customer.CustomerPhotosUrl,
		CustomerSelfieUrl:          customer.CustomerSelfieUrl,
		CreatedAt:                  customer.CreatedAt,
		UpdatedAt:                  customer.UpdatedAt,
		MembershipRewardID:         customer.MembershipRewardID,
		MembershipRewardAmmount:    customer.MembershipRewardAmmount,
		Name:                       customerGP.Data[0].Custname,
		PhonE1:                     customerGP.Data[0].PhonE1,
		CorporateCustomerNumber:    customerGP.Data[0].Cprcstnm,
		ContactPerson:              customerGP.Data[0].Cntcprsn,
		StatementName:              customerGP.Data[0].Stmtname,
		ShortName:                  customerGP.Data[0].Shrtname,
		Upszone:                    customerGP.Data[0].Upszone,
		TaxScheduleID:              customerGP.Data[0].Taxschid,
		AddresS1:                   customerGP.Data[0].AddresS1,
		AddresS2:                   customerGP.Data[0].AddresS2,
		AddresS3:                   customerGP.Data[0].AddresS3,
		Country:                    customerGP.Data[0].Country,
		City:                       customerGP.Data[0].City,
		State:                      customerGP.Data[0].State,
		Zip:                        customerGP.Data[0].Zip,
		PhonE2:                     customerGP.Data[0].PhonE2,
		PhonE3:                     customerGP.Data[0].PhonE3,
		Fax:                        customerGP.Data[0].Fax,
		PrimaryAddressCode:         customerGP.Data[0].Prbtadcd,
		StreetAddressCode:          customerGP.Data[0].Staddrcd,
		SalesPersonID:              customerGP.Data[0].Slprsnid,
		CheckbookID:                customerGP.Data[0].Chekbkid,
		CreditLimitType:            customerGP.Data[0].Crlmttyp,
		CreditLimitTypeDescription: customerGP.Data[0].CreditLimitTypeDesc,
		CreditLimitAmount:          customerGP.Data[0].Crlmtamt,
		CurrencyID:                 customerGP.Data[0].Curncyid,
		RateTypeID:                 customerGP.Data[0].Ratetpid,
		CustomerDiscount:           customerGP.Data[0].Custdisc,
		MinimumPaymentType:         customerGP.Data[0].Minpytyp,
		MinimumPaymentTypeDesc:     customerGP.Data[0].MinimumPaymentTypeDesc,
		MinimumPaymentDollarAmount: customerGP.Data[0].Minpydlr,
		MinimumPaymentPercent:      customerGP.Data[0].Minpypct,
		FinanceChargeType:          customerGP.Data[0].Fnchatyp,
		FinanceChargeTypeDesc:      customerGP.Data[0].FinanceChargeAmtTypeDesc,
		FinanceChargePercent:       customerGP.Data[0].Fnchpcnt,
		FinanceChargeDollarAmount:  customerGP.Data[0].Finchdlr,
		MaximumWriteoffType:        customerGP.Data[0].Mxwoftyp,
		MaximumWriteoffTypeDesc:    customerGP.Data[0].MaximumWriteoffTypeDesc,
		MaximumWriteoffAmount:      customerGP.Data[0].Mxwrofam,
		Comment1:                   customerGP.Data[0].CommenT1,
		Comment2:                   customerGP.Data[0].CommenT2,
		UserDefined1:               customerGP.Data[0].UserdeF1,
		UserDefined2:               customerGP.Data[0].UserdeF2,
		TaxExempt1:                 customerGP.Data[0].TaxexmT1,
		TaxExempt2:                 customerGP.Data[0].TaxexmT2,
		TaxRegistrationNumber:      customerGP.Data[0].Txrgnnum,
		BalanceType:                customerGP.Data[0].Balnctyp,
		BalanceTypeDesc:            customerGP.Data[0].BalanceTypeDesc,
		StatementCycle:             customerGP.Data[0].Stmtcycl,
		StatementCycleDesc:         customerGP.Data[0].StatementCycleDesc,
		BankName:                   customerGP.Data[0].Bankname,
		BankBranch:                 customerGP.Data[0].Bnkbrnch,
		Hold:                       customerGP.Data[0].Hold,
		CreditCardID:               customerGP.Data[0].Crcardid,
		CreditCardNumber:           customerGP.Data[0].Crcrdnum,
		CreditCardExpDate:          customerGP.Data[0].Ccrdxpdt,
		Inactive:                   customerGP.Data[0].Inactive,
		ReferrerID:                 customer.ReferrerID,
		ReferrerCode:               customer.ReferrerCode,
		ReferralCode:               customer.ReferralCode,
		CustomerClass:              customerClassResponse,
		Salesperson:                salespersonResponse,
		Archetype:                  archetypeResponse,
		CustomerType:               customerTypeResponse,
		SalesTerritory:             salesTerritoryResponse,
		PriceLevel:                 priceLevelResponse,
		Site:                       siteResponse,
		ShippingMethod:             shippingMethodResponse,
		CompanyAddress:             companyAddress,
		ShipToAddress:              shipToAddress,
		BillToAddress:              billToAddress,
		PaymentTerm:                paymentTermResponse,
		BusinessType:               businessTypeResponse,
	}

	return
}

func (s *CustomerService) GetCustomerID(ctx context.Context, req *crm_service.GetCustomerIDRequest) (res []*model.Customer, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CustomerService.GetFirebaseToken")
	defer span.End()

	res, err = s.RepositoryCustomer.GetCustomerID(ctx, req)

	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	return
}
