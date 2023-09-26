package service

import (
	"context"
	"strings"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/constants"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/edenlabs/edenlabs/timex"
	"git.edenfarm.id/edenlabs/edenlabs/utils"
	"git.edenfarm.id/edenlabs/edenlabs/validation"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/repository"
	accountService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/account_service"
	auditService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/audit_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	configurationService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
	crmService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/crm_service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IProspectiveCustomerService interface {
	Get(ctx context.Context, req *dto.ProspectiveCustomerGetRequest) (res []*dto.ProspectiveCustomerResponse, total int64, err error)
	GetDetail(ctx context.Context, req *dto.ProspectiveCustomerGetDetailRequest) (res *dto.ProspectiveCustomerResponse, err error)
	Decline(ctx context.Context, req dto.ProspectiveCustomerDecineRequest, id int64) (res *dto.ProspectiveCustomerResponse, err error)
	Delete(ctx context.Context, req *crmService.DeleteProspectiveCustomerRequest) (res *dto.ProspectiveCustomerResponse, err error)
	Create(ctx context.Context, req *dto.ProspectiveCustomerCreateRequest) (res *dto.ProspectiveCustomerResponse, err error)
	Upgrade(ctx context.Context, req *dto.ProspectiveCustomerUpgradeRequest) (res *dto.ProspectiveCustomerResponse, err error)
}

type ProspectiveCustomerService struct {
	opt                                  opt.Options
	RepositoryProspectiveCustomer        repository.IProspectiveCustomerRepository
	RepositoryProspectiveCustomerAddress repository.IProspectiveCustomerAddressRepository
	RepositoryCustomer                   repository.ICustomerRepository
}

func NewProspectiveCustomerService() IProspectiveCustomerService {
	return &ProspectiveCustomerService{
		opt:                                  global.Setup.Common,
		RepositoryProspectiveCustomer:        repository.NewProspectiveCustomerRepository(),
		RepositoryProspectiveCustomerAddress: repository.NewProspectiveCustomerAddressRepository(),
		RepositoryCustomer:                   repository.NewCustomerRepository(),
	}
}

func (s *ProspectiveCustomerService) Get(ctx context.Context, req *dto.ProspectiveCustomerGetRequest) (res []*dto.ProspectiveCustomerResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ProspectiveCustomerService.Get")
	defer span.End()

	var prospectiveCustomers []*model.ProspectiveCustomer
	prospectiveCustomers, total, err = s.RepositoryProspectiveCustomer.Get(ctx, req)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, prospectiveCustomer := range prospectiveCustomers {
		var (
			salespersonResponse                                                                                        *dto.SalespersonResponse
			archetypeResponse                                                                                          *dto.ArchetypeResponse
			customerTypeResponse                                                                                       *dto.CustomerTypeResponse
			companyAddressResponse, shippingAddressResponse, billingAddressResponse                                    *dto.ProspectiveCustomerAddressResponse
			createdByResponse                                                                                          *dto.CreatedByResponse
			paymentTermResponse                                                                                        *dto.PaymentTermResponse
			customerResponse                                                                                           *dto.CustomerResponse
			timeConcentResponse, referenceInfoResponse, invoiceTermResponse, applicationResponse, businessTypeResponse *dto.GlossaryResponse
			prospectiveCustomerResponse                                                                                *dto.ProspectiveCustomerResponse
			statusArchetype, statusCustomerType                                                                        int8
		)

		if prospectiveCustomer.SalespersonIDGP != "" {
			var salesPerson *bridgeService.GetSalesPersonGPResponse
			salesPerson, err = s.opt.Client.BridgeServiceGrpc.GetSalesPersonGPDetail(ctx, &bridgeService.GetSalesPersonGPDetailRequest{
				Id: prospectiveCustomer.SalespersonIDGP,
			})
			if err != nil || len(salesPerson.Data) == 0 {
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

		if prospectiveCustomer.ArchetypeIDGP != "" {
			var archetype *bridgeService.GetArchetypeGPResponse
			archetype, err = s.opt.Client.BridgeServiceGrpc.GetArchetypeGPDetail(ctx, &bridgeService.GetArchetypeGPDetailRequest{
				Id: prospectiveCustomer.ArchetypeIDGP,
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

		if prospectiveCustomer.CustomerTypeIDGP != "" {
			var customerType *bridgeService.GetCustomerTypeGPResponse
			customerType, err = s.opt.Client.BridgeServiceGrpc.GetCustomerTypeGPDetail(ctx, &bridgeService.GetCustomerTypeGPDetailRequest{
				Id: prospectiveCustomer.CustomerTypeIDGP,
			})
			if err != nil {
				err = edenlabs.ErrorInvalid("customer_type_id")
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}
			if customerType.Data[0].Inactive == 0 {
				statusCustomerType = statusx.ConvertStatusName(statusx.Active)
			} else {
				statusCustomerType = statusx.ConvertStatusName(statusx.Archived)
			}

			customerTypeResponse = &dto.CustomerTypeResponse{
				ID:            customerType.Data[0].GnL_Cust_Type_ID,
				Code:          customerType.Data[0].GnL_Cust_Type_ID,
				Description:   customerType.Data[0].GnL_CustType_Description,
				CustomerGroup: customerType.Data[0].GnL_Cust_GroupDesc,
				Status:        statusCustomerType,
				ConvertStatus: statusx.ConvertStatusValue(statusCustomerType),
			}
		}

		var prospcitiveCustomerAddress []*model.ProspectiveCustomerAddress

		prospcitiveCustomerAddress, _, err = s.RepositoryProspectiveCustomerAddress.Get(ctx, prospectiveCustomer.ID)
		if err != nil {
			err = edenlabs.ErrorInvalid("prospective customer id")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		for _, address := range prospcitiveCustomerAddress {
			var admDivisionID, region, province, city, district, subDistrict, postalCode string
			if address.AdmDivisionIDGP != "" {
				var admDivision *bridgeService.GetAdmDivisionGPResponse
				admDivision, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPList(ctx, &bridgeService.GetAdmDivisionGPListRequest{
					AdmDivisionCode: address.AdmDivisionIDGP,
					Limit:           1,
					Offset:          0,
				})

				if err != nil || len(admDivision.Data) == 0 {
					err = edenlabs.ErrorInvalid("adm_division_id")
					span.RecordError(err)
					s.opt.Logger.AddMessage(log.ErrorLevel, err)
					return
				}

				admDivisionID = admDivision.Data[0].Code
				region = admDivision.Data[0].Region
				province = admDivision.Data[0].State
				city = admDivision.Data[0].City
				district = admDivision.Data[0].District
				subDistrict = admDivision.Data[0].Subdistrict
				postalCode = admDivision.Data[0].Zipcode
			}

			if address.AddressType == "statement_to" {
				companyAddressResponse = &dto.ProspectiveCustomerAddressResponse{
					ID:                    address.ID,
					ProspectiveCustomerID: address.ProspectiveCustomerID,
					AddressName:           address.AddressName,
					AddressType:           address.AddressType,
					Address1:              address.Address1,
					Address2:              address.Address2,
					Address3:              address.Address3,
					AdmDivisionID:         admDivisionID,
					Region:                region,
					Province:              province,
					City:                  city,
					District:              district,
					SubDistrict:           subDistrict,
					PostalCode:            postalCode,
					CreatedAt:             address.CreatedAt,
					UpdatedAt:             address.UpdatedAt,
					ReferTo:               address.ReferTo,
					Note:                  address.Note,
					Latitude:              address.Latitude,
					Longitude:             address.Longitude,
				}
			}

			if address.AddressType == "ship_to" {
				shippingAddressResponse = &dto.ProspectiveCustomerAddressResponse{
					ID:                    address.ID,
					ProspectiveCustomerID: address.ProspectiveCustomerID,
					AddressName:           address.AddressName,
					AddressType:           address.AddressType,
					Address1:              address.Address1,
					Address2:              address.Address2,
					Address3:              address.Address3,
					AdmDivisionID:         admDivisionID,
					Region:                region,
					Province:              province,
					City:                  city,
					District:              district,
					SubDistrict:           subDistrict,
					PostalCode:            postalCode,
					CreatedAt:             address.CreatedAt,
					UpdatedAt:             address.UpdatedAt,
					ReferTo:               address.ReferTo,
					Note:                  address.Note,
					Latitude:              address.Latitude,
					Longitude:             address.Longitude,
				}
			}

			if address.AddressType == "bill_to" {
				billingAddressResponse = &dto.ProspectiveCustomerAddressResponse{
					ID:                    address.ID,
					ProspectiveCustomerID: address.ProspectiveCustomerID,
					AddressName:           address.AddressName,
					AddressType:           address.AddressType,
					Address1:              address.Address1,
					Address2:              address.Address2,
					Address3:              address.Address3,
					AdmDivisionID:         admDivisionID,
					Region:                region,
					Province:              province,
					City:                  city,
					District:              district,
					SubDistrict:           subDistrict,
					PostalCode:            postalCode,
					CreatedAt:             address.CreatedAt,
					UpdatedAt:             address.UpdatedAt,
					ReferTo:               address.ReferTo,
					Note:                  address.Note,
					Latitude:              address.Latitude,
					Longitude:             address.Longitude,
				}
			}
		}

		if prospectiveCustomer.ProcessedBy != 0 {
			var processedBy *accountService.GetUserDetailResponse
			processedBy, err = s.opt.Client.AccountServiceGrpc.GetUserDetail(ctx, &accountService.GetUserDetailRequest{
				Id: prospectiveCustomer.ProcessedBy,
			})
			if err != nil {
				err = edenlabs.ErrorInvalid("user_id")
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}

			createdByResponse = &dto.CreatedByResponse{
				ID:   processedBy.Data.Id,
				Name: processedBy.Data.Name,
			}
		}

		if prospectiveCustomer.CustomerIDGP != "" {
			var customer *bridgeService.GetCustomerGPResponse
			customer, err = s.opt.Client.BridgeServiceGrpc.GetCustomerGPList(ctx, &bridgeService.GetCustomerGPListRequest{
				Limit:  1,
				Offset: 0,
				Id:     prospectiveCustomer.CustomerIDGP,
			})
			if err != nil || len(customer.Data) == 0 {
				err = edenlabs.ErrorInvalid("customer_id")
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}
			var customerInternal *model.Customer
			customerInternal, err = s.RepositoryCustomer.GetDetail(ctx, &dto.CustomerRequestGetDetail{CustomerIDGP: prospectiveCustomer.CustomerIDGP})
			if err != nil {
				err = edenlabs.ErrorInvalid("customer_id")
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}
			customerResponse = &dto.CustomerResponse{
				ID:   customerInternal.ID,
				Code: customer.Data[0].Custnmbr,
				Name: customer.Data[0].Custname,
			}
		}

		if prospectiveCustomer.ReferenceInfo != 0 {
			var referenceInfo *configurationService.GetGlossaryDetailResponse
			referenceInfo, err = s.opt.Client.ConfigurationServiceGrpc.GetGlossaryDetail(ctx, &configurationService.GetGlossaryDetailRequest{
				Table:     "all",
				Attribute: "reference_info",
				ValueInt:  int32(prospectiveCustomer.ReferenceInfo),
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}

			referenceInfoResponse = &dto.GlossaryResponse{
				ID:        int64(referenceInfo.Data.Id),
				Table:     referenceInfo.Data.Table,
				Attribute: referenceInfo.Data.Attribute,
				ValueInt:  int8(referenceInfo.Data.ValueInt),
				ValueName: referenceInfo.Data.ValueName,
				Note:      referenceInfo.Data.Note,
			}
		}

		if prospectiveCustomer.TimeConsent != 0 {
			var timeConcent *configurationService.GetGlossaryDetailResponse
			timeConcent, err = s.opt.Client.ConfigurationServiceGrpc.GetGlossaryDetail(ctx, &configurationService.GetGlossaryDetailRequest{
				Table:     "prospect_customer",
				Attribute: "time_consent",
				ValueInt:  int32(prospectiveCustomer.TimeConsent),
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}

			timeConcentResponse = &dto.GlossaryResponse{
				ID:        int64(timeConcent.Data.Id),
				Table:     timeConcent.Data.Table,
				Attribute: timeConcent.Data.Attribute,
				ValueInt:  int8(timeConcent.Data.ValueInt),
				ValueName: timeConcent.Data.ValueName,
				Note:      timeConcent.Data.Note,
			}
		}

		if prospectiveCustomer.InvoiceTerm != 0 {
			var invoiceTerm *configurationService.GetGlossaryDetailResponse
			invoiceTerm, err = s.opt.Client.ConfigurationServiceGrpc.GetGlossaryDetail(ctx, &configurationService.GetGlossaryDetailRequest{
				Table:     "prospect_customer",
				Attribute: "invoice_term",
				ValueInt:  int32(prospectiveCustomer.InvoiceTerm),
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}

			invoiceTermResponse = &dto.GlossaryResponse{
				ID:        int64(invoiceTerm.Data.Id),
				Table:     invoiceTerm.Data.Table,
				Attribute: invoiceTerm.Data.Attribute,
				ValueInt:  int8(invoiceTerm.Data.ValueInt),
				ValueName: invoiceTerm.Data.ValueName,
				Note:      invoiceTerm.Data.Note,
			}
		}

		if prospectiveCustomer.Application != 0 {
			var application *configurationService.GetGlossaryDetailResponse
			application, err = s.opt.Client.ConfigurationServiceGrpc.GetGlossaryDetail(ctx, &configurationService.GetGlossaryDetailRequest{
				Table:     "all",
				Attribute: "application",
				ValueInt:  int32(prospectiveCustomer.Application),
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}

			applicationResponse = &dto.GlossaryResponse{
				ID:        int64(application.Data.Id),
				Table:     application.Data.Table,
				Attribute: application.Data.Attribute,
				ValueInt:  int8(application.Data.ValueInt),
				ValueName: application.Data.ValueName,
				Note:      application.Data.Note,
			}
		}

		if prospectiveCustomer.BusinessTypeIDGP != 0 {
			var businessType *configurationService.GetGlossaryDetailResponse
			businessType, err = s.opt.Client.ConfigurationServiceGrpc.GetGlossaryDetail(ctx, &configurationService.GetGlossaryDetailRequest{
				Table:     "customer",
				Attribute: "business_type",
				ValueInt:  int32(prospectiveCustomer.BusinessTypeIDGP),
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

		// Handling for data address is nil
		if shippingAddressResponse == nil {
			shippingAddressResponse = &dto.ProspectiveCustomerAddressResponse{}
		}
		if billingAddressResponse == nil {
			billingAddressResponse = &dto.ProspectiveCustomerAddressResponse{}
		}

		prospectiveCustomerResponse = &dto.ProspectiveCustomerResponse{
			ID:                         prospectiveCustomer.ID,
			Code:                       prospectiveCustomer.Code,
			Salesperson:                salespersonResponse,
			Archetype:                  archetypeResponse,
			CustomerType:               customerTypeResponse,
			Customer:                   customerResponse,
			BusinessName:               prospectiveCustomer.BusinessName,
			RegStatus:                  prospectiveCustomer.RegStatus,
			RegStatusConvert:           statusx.ConvertStatusValue(prospectiveCustomer.RegStatus),
			CreatedAt:                  timex.ToLocTime(ctx, prospectiveCustomer.CreatedAt),
			UpdatedAt:                  timex.ToLocTime(ctx, prospectiveCustomer.UpdatedAt),
			ProcessedAt:                timex.ToLocTime(ctx, prospectiveCustomer.ProcessedAt),
			ProcessedBy:                createdByResponse,
			DeclineType:                prospectiveCustomer.DeclineType,
			DeclineNote:                prospectiveCustomer.DeclineNote,
			BrandName:                  prospectiveCustomer.BrandName,
			Application:                applicationResponse,
			ShippingAddressReferTo:     shippingAddressResponse.ReferTo,
			ShippingAddressID:          shippingAddressResponse.ID,
			ShippingAddressName:        shippingAddressResponse.AddressName,
			ShippingAddressRegion:      shippingAddressResponse.Region,
			ShippingAddressDetail1:     shippingAddressResponse.Address1,
			ShippingAddressDetail2:     shippingAddressResponse.Address2,
			ShippingAddressDetail3:     shippingAddressResponse.Address3,
			ShippingAddressProvince:    shippingAddressResponse.Province,
			ShippingAddressCity:        shippingAddressResponse.City,
			ShippingAddressDistrict:    shippingAddressResponse.District,
			ShippingAddressSubDistrict: shippingAddressResponse.SubDistrict,
			ShippingAddressPostalCode:  shippingAddressResponse.PostalCode,
			ShippingAddressNote:        shippingAddressResponse.Note,
			ShippingAddressLatitude:    utils.ToString(shippingAddressResponse.Latitude),
			ShippingAddressLongitude:   utils.ToString(shippingAddressResponse.Longitude),
			BillingAddressReferTo:      billingAddressResponse.ReferTo,
			BillingAddressID:           billingAddressResponse.ID,
			BillingAddressName:         billingAddressResponse.AddressName,
			BillingAddressRegion:       billingAddressResponse.Region,
			BillingAddressDetail1:      billingAddressResponse.Address1,
			BillingAddressDetail2:      billingAddressResponse.Address2,
			BillingAddressDetail3:      billingAddressResponse.Address3,
			BillingAddressProvince:     billingAddressResponse.Province,
			BillingAddressCity:         billingAddressResponse.City,
			BillingAddressDistrict:     billingAddressResponse.District,
			BillingAddressSubDistrict:  billingAddressResponse.SubDistrict,
			BillingAddressPostalCode:   billingAddressResponse.PostalCode,
			BillingAddressNote:         billingAddressResponse.Note,
			BillingAddressLatitude:     utils.ToString(billingAddressResponse.Latitude),
			BillingAddressLongitude:    utils.ToString(billingAddressResponse.Longitude),
			OutletImage:                utils.StringToStringArray(prospectiveCustomer.OutletImage),
			TimeConsent:                timeConcentResponse,
			ReferenceInfo:              referenceInfoResponse,
			ReferrerCode:               prospectiveCustomer.ReferrerCode,
			OwnerName:                  prospectiveCustomer.OwnerName,
			OwnerRole:                  prospectiveCustomer.OwnerRole,
			Email:                      prospectiveCustomer.Email,
			BusinessType:               businessTypeResponse,
			PicOrderContact:            prospectiveCustomer.PicOrderContact,
			PicOrderName:               prospectiveCustomer.PicOrderName,
			PicFinanceContact:          prospectiveCustomer.PicFinanceContact,
			PicFinanceName:             prospectiveCustomer.PicFinanceName,
			IDCardDocName:              prospectiveCustomer.IDCardDocName,
			IDCardDocNumber:            prospectiveCustomer.IDCardDocNumber,
			IDCardDocURL:               prospectiveCustomer.IDCardDocURL,
			TaxpayerDocName:            prospectiveCustomer.TaxpayerDocName,
			TaxpayerDocNumber:          prospectiveCustomer.TaxpayerDocNumber,
			TaxpayerDocURL:             prospectiveCustomer.TaxpayerDocURL,
			CompanyContractDocName:     prospectiveCustomer.CompanyContractDocName,
			CompanyContractDocURL:      prospectiveCustomer.CompanyContractDocURL,
			NotarialDeedDocName:        prospectiveCustomer.NotarialDeedDocName,
			NotarialDeedDocURL:         prospectiveCustomer.NotarialDeedDocURL,
			TaxableEntrepeneurDocName:  prospectiveCustomer.TaxableEntrepeneurDocName,
			TaxableEntrepeneurDocURL:   prospectiveCustomer.TaxableEntrepeneurDocURL,
			CompanyCertificateRegName:  prospectiveCustomer.CompanyCertificateRegName,
			CompanyCertificateRegURL:   prospectiveCustomer.CompanyCertificateRegURL,
			BusinessLicenseDocName:     prospectiveCustomer.BusinessLicenseDocName,
			BusinessLicenseDocURL:      prospectiveCustomer.BusinessLicenseDocURL,
			PaymentTerm:                paymentTermResponse,
			ExchangeInvoice:            prospectiveCustomer.ExchangeInvoice,
			ExchangeInvoiceTime:        prospectiveCustomer.ExchangeInvoiceTime,
			FinanceEmail:               prospectiveCustomer.FinanceEmail,
			InvoiceTerm:                invoiceTermResponse,
			Comment1:                   prospectiveCustomer.Comment1,
			Comment2:                   prospectiveCustomer.Comment2,
			PicOperationName:           prospectiveCustomer.PicOperationName,
			PicOperationContact:        prospectiveCustomer.PicOperationContact,
			OwnerContact:               prospectiveCustomer.OwnerContact,
		}

		if companyAddressResponse != nil {
			// Handling nil coordinate for didn't give value 0
			var latitude, longitude string
			if utils.ToString(companyAddressResponse.Latitude) != "0" {
				latitude = utils.ToString(companyAddressResponse.Latitude)
			}
			if utils.ToString(companyAddressResponse.Longitude) != "0" {
				longitude = utils.ToString(companyAddressResponse.Longitude)
			}
			prospectiveCustomerResponse.CompanyAddressID = companyAddressResponse.ID
			prospectiveCustomerResponse.CompanyAddressName = companyAddressResponse.AddressName
			prospectiveCustomerResponse.CompanyAddressRegion = companyAddressResponse.Region
			prospectiveCustomerResponse.CompanyAddressDetail1 = companyAddressResponse.Address1
			prospectiveCustomerResponse.CompanyAddressDetail2 = companyAddressResponse.Address2
			prospectiveCustomerResponse.CompanyAddressDetail3 = companyAddressResponse.Address3
			prospectiveCustomerResponse.CompanyAddressProvince = companyAddressResponse.Province
			prospectiveCustomerResponse.CompanyAddressCity = companyAddressResponse.City
			prospectiveCustomerResponse.CompanyAddressDistrict = companyAddressResponse.District
			prospectiveCustomerResponse.CompanyAddressSubDistrict = companyAddressResponse.SubDistrict
			prospectiveCustomerResponse.CompanyAddressPostalCode = companyAddressResponse.PostalCode
			prospectiveCustomerResponse.CompanyAddressNote = companyAddressResponse.Note
			prospectiveCustomerResponse.CompanyAddressLatitude = latitude
			prospectiveCustomerResponse.CompanyAddressLongitude = longitude
		}

		// var customer *model.Customer
		res = append(res, prospectiveCustomerResponse)
	}

	return
}

func (s *ProspectiveCustomerService) GetDetail(ctx context.Context, req *dto.ProspectiveCustomerGetDetailRequest) (res *dto.ProspectiveCustomerResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ProspectiveCustomerService.GetDetail")
	defer span.End()

	var (
		prospectiveCustomer                                                                                        *model.ProspectiveCustomer
		salespersonResponse                                                                                        *dto.SalespersonResponse
		archetypeResponse                                                                                          *dto.ArchetypeResponse
		customerTypeResponse                                                                                       *dto.CustomerTypeResponse
		companyAddressResponse, shippingAddressResponse, billingAddressResponse                                    *dto.ProspectiveCustomerAddressResponse
		createdByResponse                                                                                          *dto.CreatedByResponse
		paymentTermResponse                                                                                        *dto.PaymentTermResponse
		customerResponse                                                                                           *dto.CustomerResponse
		shippingMethodResponse                                                                                     *dto.ShippingMethodResponse
		priceLevelResponse                                                                                         *dto.PriceLevelResponse
		salesTerritoryResponse                                                                                     *dto.SalesTerritoryResponse
		siteResponse                                                                                               *dto.SiteResponse
		customerClassResponse                                                                                      *dto.CustomerClassResponse
		timeConcentResponse, referenceInfoResponse, invoiceTermResponse, applicationResponse, businessTypeResponse *dto.GlossaryResponse
		statusArchetype, statusCustomerType                                                                        int8
		outletImageList                                                                                            []string
	)

	prospectiveCustomer, err = s.RepositoryProspectiveCustomer.GetDetail(ctx, req)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if prospectiveCustomer.SalespersonIDGP != "" {
		var salesPerson *bridgeService.GetSalesPersonGPResponse
		salesPerson, err = s.opt.Client.BridgeServiceGrpc.GetSalesPersonGPDetail(ctx, &bridgeService.GetSalesPersonGPDetailRequest{
			Id: prospectiveCustomer.SalespersonIDGP,
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

	if prospectiveCustomer.ArchetypeIDGP != "" {

		var archetype *bridgeService.GetArchetypeGPResponse
		archetype, err = s.opt.Client.BridgeServiceGrpc.GetArchetypeGPDetail(ctx, &bridgeService.GetArchetypeGPDetailRequest{
			Id: prospectiveCustomer.ArchetypeIDGP,
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

	if prospectiveCustomer.CustomerTypeIDGP != "" {

		var customerType *bridgeService.GetCustomerTypeGPResponse
		customerType, err = s.opt.Client.BridgeServiceGrpc.GetCustomerTypeGPDetail(ctx, &bridgeService.GetCustomerTypeGPDetailRequest{
			Id: prospectiveCustomer.CustomerTypeIDGP,
		})
		if err != nil {
			err = edenlabs.ErrorInvalid("customer_type_id")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		if customerType.Data[0].Inactive == 0 {
			statusCustomerType = statusx.ConvertStatusName(statusx.Active)
		} else {
			statusCustomerType = statusx.ConvertStatusName(statusx.Archived)
		}

		customerTypeResponse = &dto.CustomerTypeResponse{
			ID:            customerType.Data[0].GnL_Cust_Type_ID,
			Code:          customerType.Data[0].GnL_Cust_Type_ID,
			Description:   customerType.Data[0].GnL_CustType_Description,
			CustomerGroup: customerType.Data[0].GnL_Cust_GroupDesc,
			Status:        statusCustomerType,
			ConvertStatus: statusx.ConvertStatusValue(statusCustomerType),
		}
	}

	if prospectiveCustomer.ProcessedBy != 0 {

		var processedBy *accountService.GetUserDetailResponse
		processedBy, err = s.opt.Client.AccountServiceGrpc.GetUserDetail(ctx, &accountService.GetUserDetailRequest{
			Id: prospectiveCustomer.ProcessedBy,
		})
		if err != nil {
			err = edenlabs.ErrorInvalid("user_id")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		createdByResponse = &dto.CreatedByResponse{
			ID:   processedBy.Data.Id,
			Name: processedBy.Data.Name,
		}
	}

	var prosepcitiveCustomerAddress []*model.ProspectiveCustomerAddress

	prosepcitiveCustomerAddress, _, err = s.RepositoryProspectiveCustomerAddress.Get(ctx, prospectiveCustomer.ID)
	if err != nil {
		err = edenlabs.ErrorInvalid("prospective_customer_id")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, address := range prosepcitiveCustomerAddress {
		var admDivisionID, region, province, city, district, subDistrict, postalCode string
		if address.AdmDivisionIDGP != "" {
			var admDivision *bridgeService.GetAdmDivisionGPResponse
			admDivision, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPList(ctx, &bridgeService.GetAdmDivisionGPListRequest{
				AdmDivisionCode: address.AdmDivisionIDGP,
				Limit:           1,
				Offset:          0,
			})

			if err != nil || len(admDivision.Data) == 0 {
				err = edenlabs.ErrorInvalid("adm_division_id")
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}

			admDivisionID = admDivision.Data[0].Code
			region = admDivision.Data[0].Region
			province = admDivision.Data[0].State
			city = admDivision.Data[0].City
			district = admDivision.Data[0].District
			subDistrict = admDivision.Data[0].Subdistrict
			postalCode = admDivision.Data[0].Zipcode
		}

		if address.AddressType == "statement_to" {
			companyAddressResponse = &dto.ProspectiveCustomerAddressResponse{
				ID:                    address.ID,
				ProspectiveCustomerID: address.ProspectiveCustomerID,
				AddressName:           address.AddressName,
				AddressType:           address.AddressType,
				Address1:              address.Address1,
				Address2:              address.Address2,
				Address3:              address.Address3,
				AdmDivisionID:         admDivisionID,
				Region:                region,
				Province:              province,
				City:                  city,
				District:              district,
				SubDistrict:           subDistrict,
				PostalCode:            postalCode,
				CreatedAt:             address.CreatedAt,
				UpdatedAt:             address.UpdatedAt,
				ReferTo:               address.ReferTo,
				Note:                  address.Note,
				Latitude:              address.Latitude,
				Longitude:             address.Longitude,
			}
		}

		if address.AddressType == "ship_to" {
			shippingAddressResponse = &dto.ProspectiveCustomerAddressResponse{
				ID:                    address.ID,
				ProspectiveCustomerID: address.ProspectiveCustomerID,
				AddressName:           address.AddressName,
				AddressType:           address.AddressType,
				Address1:              address.Address1,
				Address2:              address.Address2,
				Address3:              address.Address3,
				AdmDivisionID:         admDivisionID,
				Region:                region,
				Province:              province,
				City:                  city,
				District:              district,
				SubDistrict:           subDistrict,
				PostalCode:            postalCode,
				CreatedAt:             address.CreatedAt,
				UpdatedAt:             address.UpdatedAt,
				ReferTo:               address.ReferTo,
				Note:                  address.Note,
				Latitude:              address.Latitude,
				Longitude:             address.Longitude,
			}
		}

		if address.AddressType == "bill_to" {
			billingAddressResponse = &dto.ProspectiveCustomerAddressResponse{
				ID:                    address.ID,
				ProspectiveCustomerID: address.ProspectiveCustomerID,
				AddressName:           address.AddressName,
				AddressType:           address.AddressType,
				Address1:              address.Address1,
				Address2:              address.Address2,
				Address3:              address.Address3,
				AdmDivisionID:         admDivisionID,
				Region:                region,
				Province:              province,
				City:                  city,
				District:              district,
				SubDistrict:           subDistrict,
				PostalCode:            postalCode,
				CreatedAt:             address.CreatedAt,
				UpdatedAt:             address.UpdatedAt,
				ReferTo:               address.ReferTo,
				Note:                  address.Note,
				Latitude:              address.Latitude,
				Longitude:             address.Longitude,
			}
		}
	}

	if prospectiveCustomer.PaymentTermIDGP != "" {

		var paymentTerm *bridgeService.GetPaymentTermGPResponse
		paymentTerm, err = s.opt.Client.BridgeServiceGrpc.GetPaymentTermGPDetail(ctx, &bridgeService.GetPaymentTermGPDetailRequest{
			Id: prospectiveCustomer.PaymentTermIDGP,
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

	if prospectiveCustomer.CustomerIDGP != "" {
		var customer *bridgeService.GetCustomerGPResponse
		customer, err = s.opt.Client.BridgeServiceGrpc.GetCustomerGPList(ctx, &bridgeService.GetCustomerGPListRequest{
			Limit:  1,
			Offset: 0,
			Id:     prospectiveCustomer.CustomerIDGP,
		})
		if err != nil || len(customer.Data) == 0 {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		var customerInternal *model.Customer
		customerInternal, err = s.RepositoryCustomer.GetDetail(ctx, &dto.CustomerRequestGetDetail{CustomerIDGP: prospectiveCustomer.CustomerIDGP})
		if err != nil || len(customer.Data) == 0 {
			err = edenlabs.ErrorInvalid("customer_id")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		customerResponse = &dto.CustomerResponse{
			ID:   customerInternal.ID,
			Code: customer.Data[0].Custnmbr,
			Name: customer.Data[0].Custname,
		}
	}

	if prospectiveCustomer.ShippingMethodIDGP != "" {
		var shippingMethod *bridgeService.GetShippingMethodResponse
		shippingMethod, err = s.opt.Client.BridgeServiceGrpc.GetShippingMethodDetail(ctx, &bridgeService.GetShippingMethodDetailRequest{
			Id: prospectiveCustomer.ShippingMethodIDGP,
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

	if prospectiveCustomer.SalesPriceLevelIDGP != "" {
		var priceLevel *bridgeService.GetSalesPriceLevelResponse
		priceLevel, err = s.opt.Client.BridgeServiceGrpc.GetSalesPriceLevelDetail(ctx, &bridgeService.GetSalesPriceLevelDetailRequest{
			Id: prospectiveCustomer.SalesPriceLevelIDGP,
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

	if prospectiveCustomer.SalesTerritoryIDGP != "" {
		var salesTerritory *bridgeService.GetSalesTerritoryGPResponse
		salesTerritory, err = s.opt.Client.BridgeServiceGrpc.GetSalesTerritoryGPDetail(ctx, &bridgeService.GetSalesTerritoryGPDetailRequest{
			Id: prospectiveCustomer.SalesTerritoryIDGP,
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

	if prospectiveCustomer.SiteIDGP != "" {
		var site *bridgeService.GetSiteGPResponse
		site, err = s.opt.Client.BridgeServiceGrpc.GetSiteGPDetail(ctx, &bridgeService.GetSiteGPDetailRequest{
			Id: prospectiveCustomer.SiteIDGP,
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

	if prospectiveCustomer.CustomerClassIDGP != "" {
		var customerClass *bridgeService.GetCustomerClassResponse
		customerClass, err = s.opt.Client.BridgeServiceGrpc.GetCustomerClassDetail(ctx, &bridgeService.GetCustomerClassDetailRequest{
			Id: prospectiveCustomer.CustomerClassIDGP,
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

	if prospectiveCustomer.ReferenceInfo != 0 {
		var referenceInfo *configurationService.GetGlossaryDetailResponse
		referenceInfo, err = s.opt.Client.ConfigurationServiceGrpc.GetGlossaryDetail(ctx, &configurationService.GetGlossaryDetailRequest{
			Table:     "all",
			Attribute: "reference_info",
			ValueInt:  int32(prospectiveCustomer.ReferenceInfo),
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		referenceInfoResponse = &dto.GlossaryResponse{
			ID:        int64(referenceInfo.Data.Id),
			Table:     referenceInfo.Data.Table,
			Attribute: referenceInfo.Data.Attribute,
			ValueInt:  int8(referenceInfo.Data.ValueInt),
			ValueName: referenceInfo.Data.ValueName,
			Note:      referenceInfo.Data.Note,
		}
	}

	if prospectiveCustomer.TimeConsent != 0 {
		var timeConcent *configurationService.GetGlossaryDetailResponse
		timeConcent, err = s.opt.Client.ConfigurationServiceGrpc.GetGlossaryDetail(ctx, &configurationService.GetGlossaryDetailRequest{
			Table:     "prospect_customer",
			Attribute: "time_consent",
			ValueInt:  int32(prospectiveCustomer.TimeConsent),
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		timeConcentResponse = &dto.GlossaryResponse{
			ID:        int64(timeConcent.Data.Id),
			Table:     timeConcent.Data.Table,
			Attribute: timeConcent.Data.Attribute,
			ValueInt:  int8(timeConcent.Data.ValueInt),
			ValueName: timeConcent.Data.ValueName,
			Note:      timeConcent.Data.Note,
		}
	}

	if prospectiveCustomer.InvoiceTerm != 0 {
		var invoiceTerm *configurationService.GetGlossaryDetailResponse
		invoiceTerm, err = s.opt.Client.ConfigurationServiceGrpc.GetGlossaryDetail(ctx, &configurationService.GetGlossaryDetailRequest{
			Table:     "prospect_customer",
			Attribute: "invoice_term",
			ValueInt:  int32(prospectiveCustomer.InvoiceTerm),
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		invoiceTermResponse = &dto.GlossaryResponse{
			ID:        int64(invoiceTerm.Data.Id),
			Table:     invoiceTerm.Data.Table,
			Attribute: invoiceTerm.Data.Attribute,
			ValueInt:  int8(invoiceTerm.Data.ValueInt),
			ValueName: invoiceTerm.Data.ValueName,
			Note:      invoiceTerm.Data.Note,
		}
	}

	if prospectiveCustomer.Application != 0 {
		var application *configurationService.GetGlossaryDetailResponse
		application, err = s.opt.Client.ConfigurationServiceGrpc.GetGlossaryDetail(ctx, &configurationService.GetGlossaryDetailRequest{
			Table:     "all",
			Attribute: "application",
			ValueInt:  int32(prospectiveCustomer.Application),
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		applicationResponse = &dto.GlossaryResponse{
			ID:        int64(application.Data.Id),
			Table:     application.Data.Table,
			Attribute: application.Data.Attribute,
			ValueInt:  int8(application.Data.ValueInt),
			ValueName: application.Data.ValueName,
			Note:      application.Data.Note,
		}
	}

	if prospectiveCustomer.BusinessTypeIDGP != 0 {
		var businessType *configurationService.GetGlossaryDetailResponse
		businessType, err = s.opt.Client.ConfigurationServiceGrpc.GetGlossaryDetail(ctx, &configurationService.GetGlossaryDetailRequest{
			Table:     "customer",
			Attribute: "business_type",
			ValueInt:  int32(prospectiveCustomer.BusinessTypeIDGP),
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

	// Handling for data address is nil
	if shippingAddressResponse == nil {
		shippingAddressResponse = &dto.ProspectiveCustomerAddressResponse{}
	}
	if billingAddressResponse == nil {
		billingAddressResponse = &dto.ProspectiveCustomerAddressResponse{}
	}

	// handling for outlet image null
	if prospectiveCustomer.OutletImage != "" {
		outletImageList = utils.StringToStringArray(prospectiveCustomer.OutletImage)
	}

	// var customer *model.Customer
	res = &dto.ProspectiveCustomerResponse{
		ID:                         prospectiveCustomer.ID,
		Code:                       prospectiveCustomer.Code,
		CustomerClass:              customerClassResponse,
		Salesperson:                salespersonResponse,
		Archetype:                  archetypeResponse,
		CustomerType:               customerTypeResponse,
		SalesTerritory:             salesTerritoryResponse,
		PriceLevel:                 priceLevelResponse,
		Site:                       siteResponse,
		ShippingMethod:             shippingMethodResponse,
		Customer:                   customerResponse,
		BusinessName:               prospectiveCustomer.BusinessName,
		RegStatus:                  prospectiveCustomer.RegStatus,
		RegStatusConvert:           statusx.ConvertStatusValue(prospectiveCustomer.RegStatus),
		CreatedAt:                  timex.ToLocTime(ctx, prospectiveCustomer.CreatedAt),
		UpdatedAt:                  timex.ToLocTime(ctx, prospectiveCustomer.UpdatedAt),
		ProcessedAt:                timex.ToLocTime(ctx, prospectiveCustomer.ProcessedAt),
		ProcessedBy:                createdByResponse,
		DeclineType:                prospectiveCustomer.DeclineType,
		DeclineNote:                prospectiveCustomer.DeclineNote,
		BrandName:                  prospectiveCustomer.BrandName,
		Application:                applicationResponse,
		ShippingAddressReferTo:     shippingAddressResponse.ReferTo,
		ShippingAddressID:          shippingAddressResponse.ID,
		ShippingAddressName:        shippingAddressResponse.AddressName,
		ShippingAddressRegion:      shippingAddressResponse.Region,
		ShippingAddressDetail1:     shippingAddressResponse.Address1,
		ShippingAddressDetail2:     shippingAddressResponse.Address2,
		ShippingAddressDetail3:     shippingAddressResponse.Address3,
		ShippingAddressProvince:    shippingAddressResponse.Province,
		ShippingAddressCity:        shippingAddressResponse.City,
		ShippingAddressDistrict:    shippingAddressResponse.District,
		ShippingAddressSubDistrict: shippingAddressResponse.SubDistrict,
		ShippingAddressPostalCode:  shippingAddressResponse.PostalCode,
		ShippingAddressNote:        shippingAddressResponse.Note,
		ShippingAddressLatitude:    utils.ToString(shippingAddressResponse.Latitude),
		ShippingAddressLongitude:   utils.ToString(shippingAddressResponse.Longitude),
		BillingAddressReferTo:      billingAddressResponse.ReferTo,
		BillingAddressID:           billingAddressResponse.ID,
		BillingAddressName:         billingAddressResponse.AddressName,
		BillingAddressRegion:       billingAddressResponse.Region,
		BillingAddressDetail1:      billingAddressResponse.Address1,
		BillingAddressDetail2:      billingAddressResponse.Address2,
		BillingAddressDetail3:      billingAddressResponse.Address3,
		BillingAddressProvince:     billingAddressResponse.Province,
		BillingAddressCity:         billingAddressResponse.City,
		BillingAddressDistrict:     billingAddressResponse.District,
		BillingAddressSubDistrict:  billingAddressResponse.SubDistrict,
		BillingAddressPostalCode:   billingAddressResponse.PostalCode,
		BillingAddressNote:         billingAddressResponse.Note,
		BillingAddressLatitude:     utils.ToString(billingAddressResponse.Latitude),
		BillingAddressLongitude:    utils.ToString(billingAddressResponse.Longitude),
		OutletImage:                outletImageList,
		TimeConsent:                timeConcentResponse,
		ReferenceInfo:              referenceInfoResponse,
		ReferrerCode:               prospectiveCustomer.ReferrerCode,
		OwnerName:                  prospectiveCustomer.OwnerName,
		OwnerRole:                  prospectiveCustomer.OwnerRole,
		Email:                      prospectiveCustomer.Email,
		BusinessType:               businessTypeResponse,
		PicOrderName:               prospectiveCustomer.PicOrderName,
		PicOrderContact:            prospectiveCustomer.PicOrderContact,
		PicFinanceContact:          prospectiveCustomer.PicFinanceContact,
		PicFinanceName:             prospectiveCustomer.PicFinanceName,
		IDCardDocName:              prospectiveCustomer.IDCardDocName,
		IDCardDocNumber:            prospectiveCustomer.IDCardDocNumber,
		IDCardDocURL:               prospectiveCustomer.IDCardDocURL,
		TaxpayerDocName:            prospectiveCustomer.TaxpayerDocName,
		TaxpayerDocNumber:          prospectiveCustomer.TaxpayerDocNumber,
		TaxpayerDocURL:             prospectiveCustomer.TaxpayerDocURL,
		CompanyContractDocName:     prospectiveCustomer.CompanyContractDocName,
		CompanyContractDocURL:      prospectiveCustomer.CompanyContractDocURL,
		NotarialDeedDocName:        prospectiveCustomer.NotarialDeedDocName,
		NotarialDeedDocURL:         prospectiveCustomer.NotarialDeedDocURL,
		TaxableEntrepeneurDocName:  prospectiveCustomer.TaxableEntrepeneurDocName,
		TaxableEntrepeneurDocURL:   prospectiveCustomer.TaxableEntrepeneurDocURL,
		CompanyCertificateRegName:  prospectiveCustomer.CompanyCertificateRegName,
		CompanyCertificateRegURL:   prospectiveCustomer.CompanyCertificateRegURL,
		BusinessLicenseDocName:     prospectiveCustomer.BusinessLicenseDocName,
		BusinessLicenseDocURL:      prospectiveCustomer.BusinessLicenseDocURL,
		PaymentTerm:                paymentTermResponse,
		ExchangeInvoice:            prospectiveCustomer.ExchangeInvoice,
		ExchangeInvoiceTime:        prospectiveCustomer.ExchangeInvoiceTime,
		FinanceEmail:               prospectiveCustomer.FinanceEmail,
		InvoiceTerm:                invoiceTermResponse,
		Comment1:                   prospectiveCustomer.Comment1,
		Comment2:                   prospectiveCustomer.Comment2,
		PicOperationName:           prospectiveCustomer.PicOperationName,
		PicOperationContact:        prospectiveCustomer.PicOperationContact,
		OwnerContact:               prospectiveCustomer.OwnerContact,
	}

	if companyAddressResponse != nil {
		// Handling nil coordinate for didn't give value 0
		var latitude, longitude string
		if utils.ToString(companyAddressResponse.Latitude) != "0" {
			latitude = utils.ToString(companyAddressResponse.Latitude)
		}
		if utils.ToString(companyAddressResponse.Longitude) != "0" {
			longitude = utils.ToString(companyAddressResponse.Longitude)
		}
		res.CompanyAddressID = companyAddressResponse.ID
		res.CompanyAddressName = companyAddressResponse.AddressName
		res.CompanyAddressRegion = companyAddressResponse.Region
		res.CompanyAddressDetail1 = companyAddressResponse.Address1
		res.CompanyAddressDetail2 = companyAddressResponse.Address2
		res.CompanyAddressDetail3 = companyAddressResponse.Address3
		res.CompanyAddressProvince = companyAddressResponse.Province
		res.CompanyAddressCity = companyAddressResponse.City
		res.CompanyAddressDistrict = companyAddressResponse.District
		res.CompanyAddressSubDistrict = companyAddressResponse.SubDistrict
		res.CompanyAddressPostalCode = companyAddressResponse.PostalCode
		res.CompanyAddressNote = companyAddressResponse.Note
		res.CompanyAddressLatitude = latitude
		res.CompanyAddressLongitude = longitude
	}

	return
}

func (s *ProspectiveCustomerService) Decline(ctx context.Context, req dto.ProspectiveCustomerDecineRequest, id int64) (res *dto.ProspectiveCustomerResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ProspectiveCustomerService.GetDetail")
	defer span.End()

	// validate is exist
	var prospectiveCustomerOld *model.ProspectiveCustomer
	prospectiveCustomerOld, err = s.RepositoryProspectiveCustomer.GetDetail(ctx, &dto.ProspectiveCustomerGetDetailRequest{ID: id})
	if err != nil {
		err = edenlabs.ErrorMustStatus("status", statusx.New)
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if prospectiveCustomerOld.RegStatus != statusx.ConvertStatusName(statusx.New) {
		err = edenlabs.ErrorMustStatus("status", statusx.New)
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	prospectiveCustomer := &model.ProspectiveCustomer{
		ID:          id,
		RegStatus:   statusx.ConvertStatusName(statusx.Declined),
		DeclineType: req.DeclineType,
		DeclineNote: req.DeclineNote,
	}

	err = s.RepositoryProspectiveCustomer.Update(ctx, prospectiveCustomer, "RegStatus", "DeclineType", "DeclineNote")
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// Update field upgrade status on customer
	if prospectiveCustomerOld.CustomerIDGP != "" {
		var customer *model.Customer
		customer, err = s.RepositoryCustomer.GetDetail(ctx, &dto.CustomerRequestGetDetail{CustomerIDGP: prospectiveCustomerOld.CustomerIDGP})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		customer.UpgradeStatus = statusx.ConvertStatusName(statusx.Declined)
		customer.UpdatedAt = time.Now()
		err = s.RepositoryCustomer.Update(ctx, customer, "UpgradeStatus", "UpdatedAt")
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	}

	userID := ctx.Value(constants.KeyUserID).(int64)
	_, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &auditService.CreateLogRequest{
		Log: &auditService.Log{
			UserId:      userID,
			ReferenceId: utils.ToString(prospectiveCustomer.ID),
			Type:        "prospective_customer",
			Function:    "decline",
			CreatedAt:   timestamppb.New(time.Now()),
		},
	})

	res = &dto.ProspectiveCustomerResponse{
		ID:                        prospectiveCustomerOld.ID,
		Code:                      prospectiveCustomerOld.Code,
		BusinessName:              prospectiveCustomerOld.BusinessName,
		RegStatus:                 prospectiveCustomer.RegStatus,
		CreatedAt:                 timex.ToLocTime(ctx, prospectiveCustomerOld.CreatedAt),
		UpdatedAt:                 timex.ToLocTime(ctx, prospectiveCustomerOld.UpdatedAt),
		ProcessedAt:               timex.ToLocTime(ctx, prospectiveCustomerOld.ProcessedAt),
		DeclineType:               prospectiveCustomer.DeclineType,
		DeclineNote:               prospectiveCustomer.DeclineNote,
		BrandName:                 prospectiveCustomerOld.BrandName,
		OutletImage:               utils.StringToStringArray(prospectiveCustomerOld.OutletImage),
		ReferrerCode:              prospectiveCustomerOld.ReferrerCode,
		OwnerName:                 prospectiveCustomerOld.OwnerName,
		OwnerRole:                 prospectiveCustomerOld.OwnerRole,
		Email:                     prospectiveCustomerOld.Email,
		PicOrderName:              prospectiveCustomerOld.PicOrderName,
		PicOrderContact:           prospectiveCustomerOld.PicOrderContact,
		PicFinanceContact:         prospectiveCustomerOld.PicFinanceContact,
		PicFinanceName:            prospectiveCustomerOld.PicFinanceName,
		IDCardDocName:             prospectiveCustomerOld.IDCardDocName,
		IDCardDocNumber:           prospectiveCustomerOld.IDCardDocNumber,
		IDCardDocURL:              prospectiveCustomerOld.IDCardDocURL,
		TaxpayerDocName:           prospectiveCustomerOld.TaxpayerDocName,
		TaxpayerDocNumber:         prospectiveCustomerOld.TaxpayerDocNumber,
		TaxpayerDocURL:            prospectiveCustomerOld.TaxpayerDocURL,
		CompanyContractDocName:    prospectiveCustomerOld.CompanyContractDocName,
		CompanyContractDocURL:     prospectiveCustomerOld.CompanyContractDocURL,
		NotarialDeedDocName:       prospectiveCustomerOld.NotarialDeedDocName,
		NotarialDeedDocURL:        prospectiveCustomerOld.NotarialDeedDocURL,
		TaxableEntrepeneurDocName: prospectiveCustomerOld.TaxableEntrepeneurDocName,
		TaxableEntrepeneurDocURL:  prospectiveCustomerOld.TaxableEntrepeneurDocURL,
		CompanyCertificateRegName: prospectiveCustomerOld.CompanyCertificateRegName,
		CompanyCertificateRegURL:  prospectiveCustomerOld.CompanyCertificateRegURL,
		BusinessLicenseDocName:    prospectiveCustomerOld.BusinessLicenseDocName,
		BusinessLicenseDocURL:     prospectiveCustomerOld.BusinessLicenseDocURL,
		ExchangeInvoice:           prospectiveCustomerOld.ExchangeInvoice,
		ExchangeInvoiceTime:       prospectiveCustomerOld.ExchangeInvoiceTime,
		FinanceEmail:              prospectiveCustomerOld.FinanceEmail,
		Comment1:                  prospectiveCustomerOld.Comment1,
		Comment2:                  prospectiveCustomerOld.Comment2,
		PicOperationName:          prospectiveCustomerOld.PicOperationName,
		PicOperationContact:       prospectiveCustomerOld.PicOperationContact,
		OwnerContact:              prospectiveCustomer.OwnerContact,
	}

	return
}

func (s *ProspectiveCustomerService) Delete(ctx context.Context, req *crmService.DeleteProspectiveCustomerRequest) (res *dto.ProspectiveCustomerResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ProspectiveCustomerService.GetDetail")
	defer span.End()

	prospectiveCustomer := &model.ProspectiveCustomer{
		ID:        req.Id,
		RegStatus: statusx.ConvertStatusName(statusx.Deleted),
	}
	err = s.RepositoryProspectiveCustomer.Update(ctx, prospectiveCustomer, "RegStatus")
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *ProspectiveCustomerService) Create(ctx context.Context, req *dto.ProspectiveCustomerCreateRequest) (res *dto.ProspectiveCustomerResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "Customer.Create")
	defer span.End()

	var (
		prospectiveCustomer, prospectiveCustomerOld                                                                                                         *model.ProspectiveCustomer
		salespersonResponse                                                                                                                                 *dto.SalespersonResponse
		archetypeResponse                                                                                                                                   *dto.ArchetypeResponse
		customerTypeResponse                                                                                                                                *dto.CustomerTypeResponse
		createdByResponse                                                                                                                                   *dto.CreatedByResponse
		paymentTermResponse                                                                                                                                 *dto.PaymentTermResponse
		customerResponse                                                                                                                                    *dto.CustomerResponse
		shippingMethodResponse                                                                                                                              *dto.ShippingMethodResponse
		priceLevelResponse                                                                                                                                  *dto.PriceLevelResponse
		salesTerritoryResponse                                                                                                                              *dto.SalesTerritoryResponse
		siteResponse                                                                                                                                        *dto.SiteResponse
		customerClassResponse                                                                                                                               *dto.CustomerClassResponse
		timeConcentResponse, referenceInfoResponse, invoiceTermResponse, applicationResponse, businessTypeResponse                                          *dto.GlossaryResponse
		statusArchetype, statusCustomerType                                                                                                                 int8
		admDivisionCompany, admDivisionShippingAddress, admDivisionBillingAddress                                                                           *bridgeService.GetAdmDivisionGPResponse
		attributeConfig, regionIDGP, admDivisionCompanyCode                                                                                                 string
		addressList                                                                                                                                         []*model.ProspectiveCustomerAddress
		companyAddress, shippingAddress, billingAddress                                                                                                     *model.ProspectiveCustomerAddress
		companyAddressLatitude, companyAddressLongitude, shippingAddressLatitude, shippingAddressLongitude, billingAddressLatitude, billingAddressLongitude float64
	)

	var count int64
	_, count, _ = s.RepositoryProspectiveCustomer.Get(ctx, &dto.ProspectiveCustomerGetRequest{CustomerID: req.CustomerCode, Status: statusx.ConvertStatusName(statusx.New)})
	if count != 0 && req.ProspectiveCustomerID == 0 {
		err = edenlabs.ErrorExists("customer_code")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// Validation brandname and set up value of attribute config available customer type
	if req.BusinessTypeID == 1 {
		// Set attribut config for individual business
		attributeConfig = "business_entity_customer_type"
	} else {
		// Set attribut config for business entity
		attributeConfig = "individual_business_customer_type"
	}

	// Validation Email
	if req.Email != "" && !validation.EmailOnly(req.Email) {
		err = edenlabs.ErrorInvalid("email")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// Validation required if exchange invoice value is 1
	if req.ExchangeInvoice == 1 {
		if req.InvoiceTerm == 0 {
			err = edenlabs.ErrorRequired("invoice_term")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		if req.FinanceEmail != "" && !validation.EmailOnly(req.FinanceEmail) {
			err = edenlabs.ErrorInvalid("financial_email")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	} else {
		// Default to direct invoice
		req.InvoiceTerm = 1
	}

	// Validate max lengt characters
	if len(req.BusinessName) > 64 {
		err = edenlabs.ErrorMustEqualOrLess("business_name", "64 characters")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// Validate max lengt characters
	if len(req.BrandName) > 64 {
		err = edenlabs.ErrorMustEqualOrLess("brand_name", "64 characters")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// Validate max lengt characters
	if len(req.ReferrerCode) > 30 {
		err = edenlabs.ErrorMustEqualOrLess("referrer_code", "30 characters")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// Validate prospective customer status must new
	if req.ProspectiveCustomerID != 0 {
		prospectiveCustomerOld, err = s.RepositoryProspectiveCustomer.GetDetail(ctx, &dto.ProspectiveCustomerGetDetailRequest{ID: req.ProspectiveCustomerID})
		if err != nil {
			err = edenlabs.ErrorInvalid("prospective_customer_id")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		if prospectiveCustomerOld.RegStatus != statusx.ConvertStatusName(statusx.New) {
			err = edenlabs.ErrorMustStatus("prospective_customer_id", statusx.New)
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	}

	// Validation Refferer Code
	if req.ReferrerCode != "" {
		var customer *bridgeService.GetCustomerGPResponse
		customer, err = s.opt.Client.BridgeServiceGrpc.GetCustomerGPList(ctx, &bridgeService.GetCustomerGPListRequest{
			Limit:        1,
			Offset:       0,
			ReferralCode: req.ReferrerCode,
			Inactive:     "0",
		})
		if err != nil || len(customer.Data) == 0 {
			err = edenlabs.ErrorInvalid("referrer_code")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		if customer.Data[0].GnlReferralCode != req.ReferrerCode {
			err = edenlabs.ErrorInvalid("referrer_code")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

	}

	if req.SalespersonID != "" {
		var salesPerson *bridgeService.GetSalesPersonGPResponse
		salesPerson, err = s.opt.Client.BridgeServiceGrpc.GetSalesPersonGPDetail(ctx, &bridgeService.GetSalesPersonGPDetailRequest{
			Id: req.SalespersonID,
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

	if req.ArchetypeID != "" {
		var archetype *bridgeService.GetArchetypeGPResponse
		archetype, err = s.opt.Client.BridgeServiceGrpc.GetArchetypeGPDetail(ctx, &bridgeService.GetArchetypeGPDetailRequest{
			Id: req.ArchetypeID,
		})
		if err != nil {
			err = edenlabs.ErrorInvalid("archetype_id")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		if archetype.Data[0].GnlCusttypeDescription == "Personal" {
			err = edenlabs.ErrorValidation("archetype_id", "cannot upgrade to personal or internal archetype")
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

	var businessType *configurationService.GetGlossaryDetailResponse
	businessType, err = s.opt.Client.ConfigurationServiceGrpc.GetGlossaryDetail(ctx, &configurationService.GetGlossaryDetailRequest{
		Table:     "customer",
		Attribute: "business_type",
		ValueInt:  int32(req.BusinessTypeID),
	})
	if err != nil {
		err = edenlabs.ErrorInvalid("business_type")
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

	if req.CustomerTypeID != "" {
		var customerType *bridgeService.GetCustomerTypeGPResponse
		customerType, err = s.opt.Client.BridgeServiceGrpc.GetCustomerTypeGPDetail(ctx, &bridgeService.GetCustomerTypeGPDetailRequest{
			Id: req.CustomerTypeID,
		})
		if err != nil {
			err = edenlabs.ErrorInvalid("customer_type_id")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		if customerType.Data[0].GnL_CustType_Description == "Personal" {
			err = edenlabs.ErrorValidation("customer_type_id", "cannot upgrade to personal")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		var businessConfig *configurationService.GetConfigAppDetailResponse
		businessConfig, err = s.opt.Client.ConfigurationServiceGrpc.GetConfigAppDetail(ctx, &configurationService.GetConfigAppDetailRequest{
			Attribute: attributeConfig,
		})
		if err != nil {
			err = edenlabs.ErrorInvalid("config_app")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		if !strings.Contains(businessConfig.Data.Value, req.CustomerTypeID) {
			err = edenlabs.ErrorValidation("customer_type_id", "customer type not available for "+businessType.Data.Note)
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		if customerType.Data[0].Inactive == 0 {
			statusCustomerType = statusx.ConvertStatusName(statusx.Active)
		} else {
			statusCustomerType = statusx.ConvertStatusName(statusx.Archived)
		}

		customerTypeResponse = &dto.CustomerTypeResponse{
			ID:            customerType.Data[0].GnL_Cust_Type_ID,
			Code:          customerType.Data[0].GnL_Cust_Type_ID,
			Description:   customerType.Data[0].GnL_CustType_Description,
			CustomerGroup: customerType.Data[0].GnL_Cust_GroupDesc,
			Status:        statusCustomerType,
			ConvertStatus: statusx.ConvertStatusValue(statusCustomerType),
		}
	}

	userID := ctx.Value(constants.KeyUserID).(int64)

	var processedBy *accountService.GetUserDetailResponse
	processedBy, err = s.opt.Client.AccountServiceGrpc.GetUserDetail(ctx, &accountService.GetUserDetailRequest{
		Id: userID,
	})
	if err != nil {
		err = edenlabs.ErrorInvalid("user_id")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	createdByResponse = &dto.CreatedByResponse{
		ID:   processedBy.Data.Id,
		Name: processedBy.Data.Name,
	}

	if req.PaymentTermID != "" {
		var paymentTerm *bridgeService.GetPaymentTermGPResponse
		paymentTerm, err = s.opt.Client.BridgeServiceGrpc.GetPaymentTermGPDetail(ctx, &bridgeService.GetPaymentTermGPDetailRequest{
			Id: req.PaymentTermID,
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

	var customer *bridgeService.GetCustomerGPResponse
	customer, err = s.opt.Client.BridgeServiceGrpc.GetCustomerGPList(ctx, &bridgeService.GetCustomerGPListRequest{
		Limit:  1,
		Offset: 0,
		Id:     req.CustomerCode,
	})
	if err != nil || len(customer.Data) == 0 {
		err = edenlabs.ErrorInvalid("customer_code")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if customer.Data[0].CustomerType[0].GnL_CustType_Description != "Personal" {
		err = edenlabs.ErrorValidation("customer_code", "Customer type of customer must be personal")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// Validation for referrer code
	if req.ReferrerCode == customer.Data[0].GnlReferralCode {
		err = edenlabs.ErrorValidation("referrer_code", "Referrer code cannot be same with referral code of customer")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	var customerInternal *model.Customer
	customerInternal, err = s.RepositoryCustomer.GetDetail(ctx, &dto.CustomerRequestGetDetail{CustomerIDGP: req.CustomerCode})
	if err != nil || len(customer.Data) == 0 {
		err = edenlabs.ErrorInvalid("customer_id")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	customerResponse = &dto.CustomerResponse{
		ID:   customerInternal.ID,
		Code: customer.Data[0].Custnmbr,
		Name: customer.Data[0].Custname,
	}

	if req.ShippingMethodID != "" {
		var shippingMethod *bridgeService.GetShippingMethodResponse
		shippingMethod, err = s.opt.Client.BridgeServiceGrpc.GetShippingMethodDetail(ctx, &bridgeService.GetShippingMethodDetailRequest{
			Id: req.ShippingMethodID,
		})
		if err != nil || len(shippingMethod.Data) == 0 {
			err = edenlabs.ErrorInvalid("shipping_method_id")
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

	if req.PriceLevelID != "" {
		var priceLevel *bridgeService.GetSalesPriceLevelResponse
		priceLevel, err = s.opt.Client.BridgeServiceGrpc.GetSalesPriceLevelDetail(ctx, &bridgeService.GetSalesPriceLevelDetailRequest{
			Id: req.PriceLevelID,
		})
		if err != nil || len(priceLevel.Data) == 0 {
			err = edenlabs.ErrorInvalid("price_level_id")
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

	if req.SalesTerritoryID != "" {
		var salesTerritory *bridgeService.GetSalesTerritoryGPResponse
		salesTerritory, err = s.opt.Client.BridgeServiceGrpc.GetSalesTerritoryGPDetail(ctx, &bridgeService.GetSalesTerritoryGPDetailRequest{
			Id: req.SalesTerritoryID,
		})
		if err != nil || len(salesTerritory.Data) == 0 {
			err = edenlabs.ErrorInvalid("sales_territory_id")
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

	if req.CustomerClassID != "" {
		var customerClass *bridgeService.GetCustomerClassResponse
		customerClass, err = s.opt.Client.BridgeServiceGrpc.GetCustomerClassDetail(ctx, &bridgeService.GetCustomerClassDetailRequest{
			Id: req.CustomerClassID,
		})
		if err != nil || len(customerClass.Data) == 0 {
			err = edenlabs.ErrorInvalid("customer_class_id")
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

	if req.ReferenceInfo != 0 {
		var referenceInfo *configurationService.GetGlossaryDetailResponse
		referenceInfo, err = s.opt.Client.ConfigurationServiceGrpc.GetGlossaryDetail(ctx, &configurationService.GetGlossaryDetailRequest{
			Table:     "all",
			Attribute: "reference_info",
			ValueInt:  int32(req.ReferenceInfo),
		})
		if err != nil {
			err = edenlabs.ErrorInvalid("reference_info")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		referenceInfoResponse = &dto.GlossaryResponse{
			ID:        int64(referenceInfo.Data.Id),
			Table:     referenceInfo.Data.Table,
			Attribute: referenceInfo.Data.Attribute,
			ValueInt:  int8(referenceInfo.Data.ValueInt),
			ValueName: referenceInfo.Data.ValueName,
			Note:      referenceInfo.Data.Note,
		}
	}

	if req.TimeConsent != 0 {
		var timeConcent *configurationService.GetGlossaryDetailResponse
		timeConcent, err = s.opt.Client.ConfigurationServiceGrpc.GetGlossaryDetail(ctx, &configurationService.GetGlossaryDetailRequest{
			Table:     "prospect_customer",
			Attribute: "time_consent",
			ValueInt:  int32(req.TimeConsent),
		})
		if err != nil {
			err = edenlabs.ErrorInvalid("time_consent")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		timeConcentResponse = &dto.GlossaryResponse{
			ID:        int64(timeConcent.Data.Id),
			Table:     timeConcent.Data.Table,
			Attribute: timeConcent.Data.Attribute,
			ValueInt:  int8(timeConcent.Data.ValueInt),
			ValueName: timeConcent.Data.ValueName,
			Note:      timeConcent.Data.Note,
		}
	}

	if req.InvoiceTerm != 0 {
		var invoiceTerm *configurationService.GetGlossaryDetailResponse
		invoiceTerm, err = s.opt.Client.ConfigurationServiceGrpc.GetGlossaryDetail(ctx, &configurationService.GetGlossaryDetailRequest{
			Table:     "prospect_customer",
			Attribute: "invoice_term",
			ValueInt:  int32(req.InvoiceTerm),
		})
		if err != nil {
			err = edenlabs.ErrorInvalid("invoice_term")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		invoiceTermResponse = &dto.GlossaryResponse{
			ID:        int64(invoiceTerm.Data.Id),
			Table:     invoiceTerm.Data.Table,
			Attribute: invoiceTerm.Data.Attribute,
			ValueInt:  int8(invoiceTerm.Data.ValueInt),
			ValueName: invoiceTerm.Data.ValueName,
			Note:      invoiceTerm.Data.Note,
		}
	}

	if req.RegistrationChannel != 0 {
		var application *configurationService.GetGlossaryDetailResponse
		application, err = s.opt.Client.ConfigurationServiceGrpc.GetGlossaryDetail(ctx, &configurationService.GetGlossaryDetailRequest{
			Table:     "all",
			Attribute: "application",
			ValueInt:  int32(req.RegistrationChannel),
		})
		if err != nil {
			err = edenlabs.ErrorInvalid("application")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		applicationResponse = &dto.GlossaryResponse{
			ID:        int64(application.Data.Id),
			Table:     application.Data.Table,
			Attribute: application.Data.Attribute,
			ValueInt:  int8(application.Data.ValueInt),
			ValueName: application.Data.ValueName,
			Note:      application.Data.Note,
		}
	}

	var prospectiveCustomerAddress *model.ProspectiveCustomerAddress
	if req.CompanyAddressID != 0 {
		prospectiveCustomerAddress, err = s.RepositoryProspectiveCustomerAddress.GetDetail(ctx, &dto.ProspectiveCustomerAddressGetDetailRequest{ID: req.CompanyAddressID})
		if err != nil || prospectiveCustomerAddress.ProspectiveCustomerID != req.ProspectiveCustomerID {
			err = edenlabs.ErrorInvalid("company_address_id")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	}

	if req.ShippingAddressID != 0 {
		prospectiveCustomerAddress, err = s.RepositoryProspectiveCustomerAddress.GetDetail(ctx, &dto.ProspectiveCustomerAddressGetDetailRequest{ID: req.ShippingAddressID})
		if err != nil || prospectiveCustomerAddress.ProspectiveCustomerID != req.ProspectiveCustomerID {
			err = edenlabs.ErrorInvalid("shipping_address_id")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	}

	if req.BillingAddressID != 0 {
		prospectiveCustomerAddress, err = s.RepositoryProspectiveCustomerAddress.GetDetail(ctx, &dto.ProspectiveCustomerAddressGetDetailRequest{ID: req.BillingAddressID})
		if err != nil || prospectiveCustomerAddress.ProspectiveCustomerID != req.ProspectiveCustomerID {
			err = edenlabs.ErrorInvalid("billing_address_id")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	}

	// Change character comma to point
	req.CompanyAddressLatitude = strings.ReplaceAll(req.CompanyAddressLatitude, ",", ".")
	req.CompanyAddressLongitude = strings.ReplaceAll(req.CompanyAddressLongitude, ",", ".")
	req.ShippingAddressLatitude = strings.ReplaceAll(req.ShippingAddressLatitude, ",", ".")
	req.ShippingAddressLongitude = strings.ReplaceAll(req.ShippingAddressLongitude, ",", ".")
	req.BillingAddressLatitude = strings.ReplaceAll(req.BillingAddressLatitude, ",", ".")
	req.BillingAddressLongitude = strings.ReplaceAll(req.BillingAddressLongitude, ",", ".")

	//
	companyAddressLatitude = utils.ToFloat(req.CompanyAddressLatitude)
	companyAddressLongitude = utils.ToFloat(req.CompanyAddressLongitude)
	shippingAddressLatitude = utils.ToFloat(req.ShippingAddressLatitude)
	shippingAddressLongitude = utils.ToFloat(req.ShippingAddressLongitude)
	billingAddressLatitude = utils.ToFloat(req.BillingAddressLatitude)
	billingAddressLongitude = utils.ToFloat(req.BillingAddressLongitude)

	// Validation Value Latitude Shipping Address
	if shippingAddressLatitude < -90 || shippingAddressLatitude > 90 {
		err = edenlabs.ErrorValidation("shipping_address_latitude", "Latitude value must more than equal -90 and less than equal 90")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	// Validation Value Longitude Shipping Address
	if shippingAddressLongitude < -180 || shippingAddressLongitude > 180 {
		err = edenlabs.ErrorValidation("shipping_address_longitude", "Longitude value must more than equal -180 and less than equal 180")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// Validation Value Latitude Billing Address
	if billingAddressLatitude < -90 || billingAddressLatitude > 90 {
		err = edenlabs.ErrorValidation("billing_address_latitude", "Latitude value must more than equal -90 and less than equal 90")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	// Validation Value Longitude Billing Address
	if billingAddressLongitude < -180 || billingAddressLongitude > 180 {
		err = edenlabs.ErrorValidation("billing_address_longitude", "Longitude value must more than equal -180 and less than equal 180")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if req.BusinessTypeID == 1 {
		// Validation Value Latitude Company Address
		if companyAddressLatitude < -90 || companyAddressLatitude > 90 {
			err = edenlabs.ErrorValidation("company_address_latitude", "Latitude value must more than equal -90 and less than equal 90")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		// Validation Value Longitude Company Address
		if companyAddressLongitude < -180 || companyAddressLongitude > 180 {
			err = edenlabs.ErrorValidation("company_address_longitude", "Longitude value must more than equal -180 and less than equal 180")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		if req.CompanyAddressRegion != "" || req.CompanyAddressSubDistrict != "" {
			// Get Adm Division
			admDivisionCompany, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPList(ctx, &bridgeService.GetAdmDivisionGPListRequest{
				Limit:       1,
				Offset:      0,
				State:       req.CompanyAddressProvince,
				City:        req.CompanyAddressCity,
				District:    req.CompanyAddressDistrict,
				SubDistrict: req.CompanyAddressSubDistrict,
			})

			if err != nil || len(admDivisionCompany.Data) == 0 {
				err = edenlabs.ErrorInvalid("company_address_sub_district")
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}

			admDivisionCompanyCode = admDivisionCompany.Data[0].Code
			// Set region id gp
			regionIDGP = admDivisionCompany.Data[0].Region
		}

		companyAddress = &model.ProspectiveCustomerAddress{
			ID:              int64(req.CompanyAddressID),
			AddressName:     req.CompanyAddressName,
			AddressType:     "statement_to",
			Address1:        req.CompanyAddressDetail1,
			Address2:        req.CompanyAddressDetail2,
			Address3:        req.CompanyAddressDetail3,
			AdmDivisionIDGP: admDivisionCompanyCode,
			Latitude:        utils.ToFloat(req.CompanyAddressLatitude),
			Longitude:       utils.ToFloat(req.CompanyAddressLongitude),
			Note:            req.CompanyAddressNote,
		}

		addressList = append(addressList, companyAddress)
	}

	if req.ShippingAddressReferTo == 1 && req.BusinessTypeID == 1 {
		shippingAddress = &model.ProspectiveCustomerAddress{
			ID:              int64(req.ShippingAddressID),
			AddressName:     req.CompanyAddressName,
			AddressType:     "ship_to",
			Address1:        req.CompanyAddressDetail1,
			Address2:        req.CompanyAddressDetail2,
			Address3:        req.CompanyAddressDetail3,
			AdmDivisionIDGP: admDivisionCompanyCode,
			Latitude:        utils.ToFloat(req.CompanyAddressLatitude),
			Longitude:       utils.ToFloat(req.CompanyAddressLongitude),
			ReferTo:         req.ShippingAddressReferTo,
			Note:            req.CompanyAddressNote,
		}
		addressList = append(addressList, shippingAddress)
	} else {
		admDivisionShippingAddress, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPList(ctx, &bridgeService.GetAdmDivisionGPListRequest{
			Limit:       1,
			Offset:      0,
			State:       req.ShippingAddressProvince,
			City:        req.ShippingAddressCity,
			District:    req.ShippingAddressDistrict,
			SubDistrict: req.ShippingAddressSubDistrict,
		})

		if err != nil || len(admDivisionShippingAddress.Data) == 0 {
			err = edenlabs.ErrorInvalid("shipping_address_sub_district")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		// Set region id gp
		regionIDGP = admDivisionShippingAddress.Data[0].Region

		shippingAddress = &model.ProspectiveCustomerAddress{
			ID:              int64(req.ShippingAddressID),
			AddressName:     req.ShippingAddressName,
			AddressType:     "ship_to",
			Address1:        req.ShippingAddressDetail1,
			Address2:        req.ShippingAddressDetail2,
			Address3:        req.ShippingAddressDetail3,
			AdmDivisionIDGP: admDivisionShippingAddress.Data[0].Code,
			Latitude:        utils.ToFloat(req.ShippingAddressLatitude),
			Longitude:       utils.ToFloat(req.ShippingAddressLongitude),
			Note:            req.ShippingAddressNote,
		}

		addressList = append(addressList, shippingAddress)
	}

	var admDivisionCoverage *bridgeService.GetAdmDivisionCoverageGPResponse
	admDivisionCoverage, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionCoverageGPList(ctx, &bridgeService.GetAdmDivisionCoverageGPListRequest{
		Limit:                 1,
		Offset:                0,
		GnlAdministrativeCode: shippingAddress.AdmDivisionIDGP,
	})

	if err != nil || len(admDivisionCoverage.Data) == 0 {
		err = edenlabs.ErrorInvalid("shipping_address_sub_district")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	siteResponse = &dto.SiteResponse{
		ID: admDivisionCoverage.Data[0].Locncode,
	}

	if (req.BillingAddressReferTo == 1 || req.BillingAddressReferTo == 3) && req.BusinessTypeID == 1 {
		billingAddress = &model.ProspectiveCustomerAddress{
			ID:              int64(req.BillingAddressID),
			AddressName:     req.CompanyAddressName,
			AddressType:     "bill_to",
			Address1:        req.CompanyAddressDetail1,
			Address2:        req.CompanyAddressDetail2,
			Address3:        req.CompanyAddressDetail3,
			AdmDivisionIDGP: admDivisionCompanyCode,
			Latitude:        utils.ToFloat(req.CompanyAddressLatitude),
			Longitude:       utils.ToFloat(req.CompanyAddressLongitude),
			ReferTo:         req.BillingAddressReferTo,
			Note:            req.CompanyAddressNote,
		}
		addressList = append(addressList, billingAddress)
	} else if req.BillingAddressReferTo == 2 {
		billingAddress = &model.ProspectiveCustomerAddress{
			ID:              int64(req.BillingAddressID),
			AddressName:     req.ShippingAddressName,
			AddressType:     "bill_to",
			Address1:        req.ShippingAddressDetail1,
			Address2:        req.ShippingAddressDetail2,
			Address3:        req.ShippingAddressDetail3,
			AdmDivisionIDGP: admDivisionShippingAddress.Data[0].Code,
			Latitude:        utils.ToFloat(req.ShippingAddressLatitude),
			Longitude:       utils.ToFloat(req.ShippingAddressLongitude),
			ReferTo:         req.BillingAddressReferTo,
			Note:            req.ShippingAddressNote,
		}
		addressList = append(addressList, billingAddress)
	} else {
		// Bill to Address
		admDivisionBillingAddress, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPList(ctx, &bridgeService.GetAdmDivisionGPListRequest{
			Limit:       1,
			Offset:      0,
			State:       req.BillingAddressProvince,
			SubDistrict: req.BillingAddressSubDistrict,
			District:    req.BillingAddressDistrict,
			City:        req.BillingAddressCity,
		})

		if err != nil || len(admDivisionBillingAddress.Data) == 0 {
			err = edenlabs.ErrorInvalid("billing_address_sub_district")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		billingAddress = &model.ProspectiveCustomerAddress{
			ID:              int64(req.BillingAddressID),
			AddressName:     req.BillingAddressName,
			AddressType:     "bill_to",
			Address1:        req.BillingAddressDetail1,
			Address2:        req.BillingAddressDetail2,
			Address3:        req.BillingAddressDetail3,
			AdmDivisionIDGP: admDivisionBillingAddress.Data[0].Code,
			Latitude:        utils.ToFloat(req.BillingAddressLatitude),
			Longitude:       utils.ToFloat(req.BillingAddressLongitude),
			Note:            req.BillingAddressNote,
		}
		addressList = append(addressList, billingAddress)
	}

	prospectiveCustomer = &model.ProspectiveCustomer{
		CustomerClassIDGP:         req.CustomerClassID,
		SalespersonIDGP:           req.SalespersonID,
		ArchetypeIDGP:             req.ArchetypeID,
		CustomerTypeIDGP:          req.CustomerTypeID,
		SalesTerritoryIDGP:        req.SalesTerritoryID,
		SalesPriceLevelIDGP:       req.PriceLevelID,
		RegionIDGP:                regionIDGP,
		SiteIDGP:                  siteResponse.ID,
		ShippingMethodIDGP:        req.ShippingMethodID,
		CustomerIDGP:              req.CustomerCode,
		BusinessName:              req.BusinessName,
		RegStatus:                 statusx.ConvertStatusName(statusx.New),
		ProcessedAt:               time.Now(),
		ProcessedBy:               userID,
		BrandName:                 req.BrandName,
		Application:               req.RegistrationChannel,
		OutletImage:               utils.ArrayStringToString(req.OutletImage),
		TimeConsent:               req.TimeConsent,
		ReferenceInfo:             req.ReferenceInfo,
		ReferrerCode:              req.ReferrerCode,
		OwnerName:                 req.OwnerName,
		OwnerRole:                 req.OwnerRole,
		Email:                     req.Email,
		BusinessTypeIDGP:          req.BusinessTypeID,
		PicOrderName:              req.PicOrderName,
		PicOrderContact:           req.PicOrderContact,
		PicFinanceContact:         req.PicFinanceContact,
		PicFinanceName:            req.PicFinanceName,
		PicOperationName:          req.PicOperationName,
		PicOperationContact:       req.PicOperationContact,
		IDCardDocName:             "ID-Card.pdf",
		IDCardDocNumber:           req.IDCardDocNumber,
		IDCardDocURL:              req.IDCardDocURL,
		TaxpayerDocName:           "NPWP.pdf",
		TaxpayerDocNumber:         req.TaxpayerDocNumber,
		TaxpayerDocURL:            req.TaxpayerDocURL,
		CompanyContractDocName:    "Company-Contract.pdf",
		CompanyContractDocURL:     req.CompanyContractDocURL,
		NotarialDeedDocName:       "Notarial-Deed.pdf",
		NotarialDeedDocURL:        req.NotarialDeedDocURL,
		TaxableEntrepeneurDocName: "Tax-Entrepreneur.pdf",
		TaxableEntrepeneurDocURL:  req.TaxableEntrepeneurDocURL,
		CompanyCertificateRegName: "Company-Certificate.pdf",
		CompanyCertificateRegURL:  req.CompanyCertificateRegURL,
		BusinessLicenseDocName:    "Business-License.pdf",
		BusinessLicenseDocURL:     req.BusinessLicenseDocURL,
		PaymentTermIDGP:           req.PaymentTermID,
		ExchangeInvoice:           req.ExchangeInvoice,
		ExchangeInvoiceTime:       req.ExchangeInvoiceTime,
		FinanceEmail:              req.FinanceEmail,
		InvoiceTerm:               req.InvoiceTerm,
		Comment1:                  req.Comment1,
		Comment2:                  req.Comment2,
		OwnerContact:              req.OwnerContact,
	}

	if req.ProspectiveCustomerID != 0 {
		prospectiveCustomer.ID = prospectiveCustomerOld.ID
		prospectiveCustomer.Code = prospectiveCustomerOld.Code
		prospectiveCustomer.CreatedAt = prospectiveCustomerOld.CreatedAt
		prospectiveCustomer.UpdatedAt = time.Now()
		err = s.RepositoryProspectiveCustomer.Update(ctx, prospectiveCustomer)
		if err != nil {
			err = edenlabs.ErrorInvalid("prospective_customer_id")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	} else {
		var codeGenerator *configurationService.GetGenerateCodeResponse
		codeGenerator, err = s.opt.Client.ConfigurationServiceGrpc.GetGenerateCode(ctx, &configurationService.GetGenerateCodeRequest{
			Format: "PCT",
			Domain: "prospective_customer",
			Length: 6,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("configuration", "generate_code")
			return
		}

		prospectiveCustomer.Code = codeGenerator.Data.Code
		prospectiveCustomer.CreatedAt = time.Now()
		_, err = s.RepositoryProspectiveCustomer.Create(ctx, prospectiveCustomer)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	}

	// Create New Address
	for _, address := range addressList {
		address.ProspectiveCustomerID = prospectiveCustomer.ID
		if address.ID != 0 {
			address.ProspectiveCustomerID = req.ProspectiveCustomerID
			address.UpdatedAt = time.Now()
			err = s.RepositoryProspectiveCustomerAddress.Update(ctx, address)
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}
		} else {
			address.CreatedAt = time.Now()
			_, err = s.RepositoryProspectiveCustomerAddress.Create(ctx, address)
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}
		}
	}
	// Update field upgrade status on customer
	if req.CustomerCode != "" {
		customerInternal.UpgradeStatus = statusx.ConvertStatusName(statusx.Requested)
		customerInternal.UpdatedAt = time.Now()
		err = s.RepositoryCustomer.Update(ctx, customerInternal, "UpgradeStatus", "UpdatedAt")
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	}

	_, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &auditService.CreateLogRequest{
		Log: &auditService.Log{
			UserId:      userID,
			ReferenceId: utils.ToString(prospectiveCustomer.ID),
			Type:        "prospective_customer",
			Function:    "save",
			CreatedAt:   timestamppb.New(time.Now()),
		},
	})

	res = &dto.ProspectiveCustomerResponse{
		ID:                        prospectiveCustomer.ID,
		Code:                      prospectiveCustomer.Code,
		CustomerClass:             customerClassResponse,
		Salesperson:               salespersonResponse,
		Archetype:                 archetypeResponse,
		CustomerType:              customerTypeResponse,
		SalesTerritory:            salesTerritoryResponse,
		PriceLevel:                priceLevelResponse,
		Site:                      siteResponse,
		ShippingMethod:            shippingMethodResponse,
		Customer:                  customerResponse,
		BusinessName:              prospectiveCustomer.BusinessName,
		RegStatus:                 prospectiveCustomer.RegStatus,
		RegStatusConvert:          statusx.ConvertStatusValue(prospectiveCustomer.RegStatus),
		CreatedAt:                 timex.ToLocTime(ctx, prospectiveCustomer.CreatedAt),
		UpdatedAt:                 timex.ToLocTime(ctx, prospectiveCustomer.UpdatedAt),
		ProcessedAt:               timex.ToLocTime(ctx, prospectiveCustomer.ProcessedAt),
		ProcessedBy:               createdByResponse,
		DeclineType:               prospectiveCustomer.DeclineType,
		DeclineNote:               prospectiveCustomer.DeclineNote,
		BrandName:                 prospectiveCustomer.BrandName,
		Application:               applicationResponse,
		OutletImage:               utils.StringToStringArray(prospectiveCustomer.OutletImage),
		TimeConsent:               timeConcentResponse,
		ReferenceInfo:             referenceInfoResponse,
		ReferrerCode:              prospectiveCustomer.ReferrerCode,
		OwnerName:                 prospectiveCustomer.OwnerName,
		OwnerRole:                 prospectiveCustomer.OwnerRole,
		Email:                     prospectiveCustomer.Email,
		BusinessType:              businessTypeResponse,
		PicOrderName:              prospectiveCustomer.PicOrderName,
		PicOrderContact:           prospectiveCustomer.PicOrderContact,
		PicFinanceContact:         prospectiveCustomer.PicFinanceContact,
		PicFinanceName:            prospectiveCustomer.PicFinanceName,
		IDCardDocName:             prospectiveCustomer.IDCardDocName,
		IDCardDocNumber:           prospectiveCustomer.IDCardDocNumber,
		IDCardDocURL:              prospectiveCustomer.IDCardDocURL,
		TaxpayerDocName:           prospectiveCustomer.TaxpayerDocName,
		TaxpayerDocNumber:         prospectiveCustomer.TaxpayerDocNumber,
		TaxpayerDocURL:            prospectiveCustomer.TaxpayerDocURL,
		CompanyContractDocName:    prospectiveCustomer.CompanyContractDocName,
		CompanyContractDocURL:     prospectiveCustomer.CompanyContractDocURL,
		NotarialDeedDocName:       prospectiveCustomer.NotarialDeedDocName,
		NotarialDeedDocURL:        prospectiveCustomer.NotarialDeedDocURL,
		TaxableEntrepeneurDocName: prospectiveCustomer.TaxableEntrepeneurDocName,
		TaxableEntrepeneurDocURL:  prospectiveCustomer.TaxableEntrepeneurDocURL,
		CompanyCertificateRegName: prospectiveCustomer.CompanyCertificateRegName,
		CompanyCertificateRegURL:  prospectiveCustomer.CompanyCertificateRegURL,
		BusinessLicenseDocName:    prospectiveCustomer.BusinessLicenseDocName,
		BusinessLicenseDocURL:     prospectiveCustomer.BusinessLicenseDocURL,
		PaymentTerm:               paymentTermResponse,
		ExchangeInvoice:           prospectiveCustomer.ExchangeInvoice,
		ExchangeInvoiceTime:       prospectiveCustomer.ExchangeInvoiceTime,
		FinanceEmail:              prospectiveCustomer.FinanceEmail,
		InvoiceTerm:               invoiceTermResponse,
		Comment1:                  prospectiveCustomer.Comment1,
		Comment2:                  prospectiveCustomer.Comment2,
	}

	return
}

func (s *ProspectiveCustomerService) Upgrade(ctx context.Context, req *dto.ProspectiveCustomerUpgradeRequest) (res *dto.ProspectiveCustomerResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ProspectiveCustomer.Upgrade")
	defer span.End()

	var (
		prospectiveCustomer, prospectiveCustomerOld *model.ProspectiveCustomer

		salespersonResponse                                                                                                                                 *dto.SalespersonResponse
		archetypeResponse                                                                                                                                   *dto.ArchetypeResponse
		customerTypeResponse                                                                                                                                *dto.CustomerTypeResponse
		createdByResponse                                                                                                                                   *dto.CreatedByResponse
		paymentTermResponse                                                                                                                                 *dto.PaymentTermResponse
		customerResponse                                                                                                                                    *dto.CustomerResponse
		shippingMethodResponse                                                                                                                              *dto.ShippingMethodResponse
		priceLevelResponse                                                                                                                                  *dto.PriceLevelResponse
		salesTerritoryResponse                                                                                                                              *dto.SalesTerritoryResponse
		siteResponse                                                                                                                                        *dto.SiteResponse
		customerClassResponse                                                                                                                               *dto.CustomerClassResponse
		timeConcentResponse, referenceInfoResponse, invoiceTermResponse, applicationResponse, businessTypeResponse                                          *dto.GlossaryResponse
		statusArchetype, statusCustomerType                                                                                                                 int8
		addressCodeCompanyAddress, addressCodeShippingAddress, AddressCodeBillingAddress, attributeConfig, regionIDGP                                       string
		addressList                                                                                                                                         []*model.ProspectiveCustomerAddress
		companyAddress, shippingAddress, billingAddress                                                                                                     *model.ProspectiveCustomerAddress
		admDivisionCompany, admDivisionShippingAddress, admDivisionBillingAddress                                                                           *bridgeService.GetAdmDivisionGPResponse
		companyAddressLatitude, companyAddressLongitude, shippingAddressLatitude, shippingAddressLongitude, billingAddressLatitude, billingAddressLongitude float64
	)

	// Validation exist data
	var count int64
	_, count, _ = s.RepositoryProspectiveCustomer.Get(ctx, &dto.ProspectiveCustomerGetRequest{CustomerID: req.CustomerCode, Status: statusx.ConvertStatusName(statusx.New)})
	if count != 0 && req.ProspectiveCustomerID == 0 {
		err = edenlabs.ErrorValidation("customer_code", "The data already exists, upgrade must be from the list")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// Validation Email
	if req.Email != "" && !validation.EmailOnly(req.Email) {
		err = edenlabs.ErrorInvalid("email")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if req.BusinessTypeID == 1 {
		if req.BrandName == "" {
			err = edenlabs.ErrorRequired("brand_name")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		// Set attribut config for business entity
		attributeConfig = "business_entity_customer_type"

		// Validate Company address
		if req.CompanyAddressName == "" {
			err = edenlabs.ErrorRequired("company_address_name")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		if req.CompanyAddressRegion == "" {
			err = edenlabs.ErrorRequired("company_address_region")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		if req.CompanyAddressProvince == "" {
			err = edenlabs.ErrorRequired("company_address_province")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		if req.CompanyAddressCity == "" {
			err = edenlabs.ErrorRequired("company_address_city")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		if req.CompanyAddressDistrict == "" {
			err = edenlabs.ErrorRequired("company_address_District")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		if req.CompanyAddressSubDistrict == "" {
			err = edenlabs.ErrorRequired("company_address_sub_district")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		if req.CompanyAddressPostalCode == "" {
			err = edenlabs.ErrorRequired("company_address_postal_code")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		if req.CompanyAddressLatitude == "" {
			err = edenlabs.ErrorRequired("company_address_latitude")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		if req.CompanyAddressLongitude == "" {
			err = edenlabs.ErrorRequired("company_address_longitude")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		if req.CompanyAddressDetail1 == "" {
			err = edenlabs.ErrorRequired("company_address_detail_1")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		if req.OwnerRole == "" {
			err = edenlabs.ErrorValidation("owner_role", "The signing role position is required")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		if req.CompanyContractDocURL == "" {
			err = edenlabs.ErrorRequired("company_contract_doc_url")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		if req.NotarialDeedDocURL == "" {
			err = edenlabs.ErrorRequired("notarial_deed_doc_url")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		if req.TaxableEntrepeneurDocURL == "" {
			err = edenlabs.ErrorRequired("taxable_entrepeneur_doc_url")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		if req.BusinessLicenseDocURL == "" {
			err = edenlabs.ErrorRequired("business_license_doc_url")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		if req.CompanyCertificateRegURL == "" {
			err = edenlabs.ErrorRequired("company_certificate_reg_url")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	} else {
		if req.OwnerContact == "" {
			err = edenlabs.ErrorValidation("owner_contact", "The business owner contact is required")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		// Set attribut config for individual business
		attributeConfig = "individual_business_customer_type"
	}

	if req.ExchangeInvoice == 1 {
		if req.ExchangeInvoiceTime == "" {
			err = edenlabs.ErrorRequired("exchange_invoice_time")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		if req.InvoiceTerm == 0 {
			err = edenlabs.ErrorRequired("invoice_term")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		if req.FinanceEmail == "" {
			err = edenlabs.ErrorRequired("financial_email")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		} else if !validation.EmailOnly(req.FinanceEmail) {
			err = edenlabs.ErrorInvalid("financial_email")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	} else {
		// Default to direct invoice
		req.InvoiceTerm = 1
	}

	// Validate max lengt characters
	if len(req.BusinessName) > 64 {
		err = edenlabs.ErrorMustEqualOrLess("business_name", "64 characters")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// Validate max lengt characters
	if len(req.BrandName) > 64 {
		err = edenlabs.ErrorMustEqualOrLess("brand_name", "64 characters")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// Validate max lengt characters
	if len(req.ReferrerCode) > 30 {
		err = edenlabs.ErrorMustEqualOrLess("referrer_code", "30 characters")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// Validate max lengt characters
	if len(req.CompanyAddressName) > 64 {
		err = edenlabs.ErrorMustEqualOrLess("company_address.address_name", "64 characters")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// Validate max lengt characters
	if len(req.ShippingAddressName) > 64 {
		err = edenlabs.ErrorMustEqualOrLess("ship_to_address.address_name", "64 characters")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// Validate max lengt characters
	if len(req.BillingAddressName) > 64 {
		err = edenlabs.ErrorMustEqualOrLess("bill_to_address.address_name", "64 characters")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// Validate prospective customer status must new
	if req.ProspectiveCustomerID != 0 {
		prospectiveCustomerOld, err = s.RepositoryProspectiveCustomer.GetDetail(ctx, &dto.ProspectiveCustomerGetDetailRequest{ID: req.ProspectiveCustomerID})
		if err != nil {
			err = edenlabs.ErrorInvalid("prospective_customer_id")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		if prospectiveCustomerOld.RegStatus != statusx.ConvertStatusName(statusx.New) {
			err = edenlabs.ErrorMustStatus("prospective_customer_id", statusx.New)
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	}

	var businessType *configurationService.GetGlossaryDetailResponse
	businessType, err = s.opt.Client.ConfigurationServiceGrpc.GetGlossaryDetail(ctx, &configurationService.GetGlossaryDetailRequest{
		Table:     "customer",
		Attribute: "business_type",
		ValueInt:  int32(req.BusinessTypeID),
	})
	if err != nil {
		err = edenlabs.ErrorInvalid("business_type_id")
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

	// Validation Refferer Code
	if req.ReferrerCode != "" {
		var customer *bridgeService.GetCustomerGPResponse
		customer, err = s.opt.Client.BridgeServiceGrpc.GetCustomerGPList(ctx, &bridgeService.GetCustomerGPListRequest{
			Limit:        1,
			Offset:       0,
			ReferralCode: req.ReferrerCode,
			Inactive:     "0",
		})
		if err != nil || len(customer.Data) == 0 {
			err = edenlabs.ErrorInvalid("referrer_code")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		if customer.Data[0].GnlReferralCode != req.ReferrerCode {
			err = edenlabs.ErrorInvalid("referrer_code")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	}

	var salesPerson *bridgeService.GetSalesPersonGPResponse
	salesPerson, err = s.opt.Client.BridgeServiceGrpc.GetSalesPersonGPDetail(ctx, &bridgeService.GetSalesPersonGPDetailRequest{
		Id: req.SalespersonID,
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

	var archetype *bridgeService.GetArchetypeGPResponse
	archetype, err = s.opt.Client.BridgeServiceGrpc.GetArchetypeGPDetail(ctx, &bridgeService.GetArchetypeGPDetailRequest{
		Id: req.ArchetypeID,
	})
	if err != nil {
		err = edenlabs.ErrorInvalid("archetype_id")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if archetype.Data[0].GnlCusttypeDescription == "Personal" {
		err = edenlabs.ErrorValidation("archetype_id", "cannot upgrade to personal or internal archetype")
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

	var customerType *bridgeService.GetCustomerTypeGPResponse
	customerType, err = s.opt.Client.BridgeServiceGrpc.GetCustomerTypeGPDetail(ctx, &bridgeService.GetCustomerTypeGPDetailRequest{
		Id: req.CustomerTypeID,
	})
	if err != nil {
		err = edenlabs.ErrorInvalid("customer_type_id")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if customerType.Data[0].GnL_CustType_Description == "Personal" {
		err = edenlabs.ErrorValidation("customer_type_id", "cannot upgrade to personal")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	var businessConfig *configurationService.GetConfigAppDetailResponse
	businessConfig, err = s.opt.Client.ConfigurationServiceGrpc.GetConfigAppDetail(ctx, &configurationService.GetConfigAppDetailRequest{
		Attribute: attributeConfig,
	})
	if err != nil {
		err = edenlabs.ErrorInvalid("config_app")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if !strings.Contains(businessConfig.Data.Value, req.CustomerTypeID) {
		err = edenlabs.ErrorValidation("customer_type_id", "customer type not available for "+businessType.Data.Note)
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if customerType.Data[0].Inactive == 0 {
		statusCustomerType = statusx.ConvertStatusName(statusx.Active)
	} else {
		statusCustomerType = statusx.ConvertStatusName(statusx.Archived)
	}

	customerTypeResponse = &dto.CustomerTypeResponse{
		ID:            customerType.Data[0].GnL_Cust_Type_ID,
		Code:          customerType.Data[0].GnL_Cust_Type_ID,
		Description:   customerType.Data[0].GnL_CustType_Description,
		CustomerGroup: customerType.Data[0].GnL_Cust_GroupDesc,
		Status:        statusCustomerType,
		ConvertStatus: statusx.ConvertStatusValue(statusCustomerType),
	}

	userID := ctx.Value(constants.KeyUserID).(int64)

	var processedBy *accountService.GetUserDetailResponse
	processedBy, err = s.opt.Client.AccountServiceGrpc.GetUserDetail(ctx, &accountService.GetUserDetailRequest{
		Id: userID,
	})
	if err != nil {
		err = edenlabs.ErrorInvalid("user_id")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	createdByResponse = &dto.CreatedByResponse{
		ID:   processedBy.Data.Id,
		Name: processedBy.Data.Name,
	}

	var paymentTerm *bridgeService.GetPaymentTermGPResponse
	paymentTerm, err = s.opt.Client.BridgeServiceGrpc.GetPaymentTermGPDetail(ctx, &bridgeService.GetPaymentTermGPDetailRequest{
		Id: req.PaymentTermID,
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

	var customer *bridgeService.GetCustomerGPResponse
	customer, err = s.opt.Client.BridgeServiceGrpc.GetCustomerGPList(ctx, &bridgeService.GetCustomerGPListRequest{
		Limit:  1,
		Offset: 0,
		Id:     req.CustomerCode,
	})

	if err != nil || len(customer.Data) == 0 {
		err = edenlabs.ErrorInvalid("customer_code")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if customer.Data[0].CustomerType[0].GnL_CustType_Description != "Personal" {
		err = edenlabs.ErrorValidation("customer_code", "Customer type of customer must be personal")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// Validation for referrer code
	if req.ReferrerCode == customer.Data[0].GnlReferralCode {
		err = edenlabs.ErrorValidation("referrer_code", "Referrer code cannot be same with referral code of customer")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	var customerInternal *model.Customer
	customerInternal, err = s.RepositoryCustomer.GetDetail(ctx, &dto.CustomerRequestGetDetail{CustomerIDGP: req.CustomerCode})
	if err != nil {
		err = edenlabs.ErrorInvalid("customer_id")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	customerResponse = &dto.CustomerResponse{
		ID:   customerInternal.ID,
		Code: customer.Data[0].Custnmbr,
		Name: customer.Data[0].Custname,
	}

	var shippingMethod *bridgeService.GetShippingMethodResponse
	shippingMethod, err = s.opt.Client.BridgeServiceGrpc.GetShippingMethodDetail(ctx, &bridgeService.GetShippingMethodDetailRequest{
		Id: req.ShippingMethodID,
	})
	if err != nil || len(shippingMethod.Data) == 0 {
		err = edenlabs.ErrorInvalid("shipping_method_id")
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

	var priceLevel *bridgeService.GetSalesPriceLevelResponse
	priceLevel, err = s.opt.Client.BridgeServiceGrpc.GetSalesPriceLevelDetail(ctx, &bridgeService.GetSalesPriceLevelDetailRequest{
		Id: req.PriceLevelID,
	})
	if err != nil || len(priceLevel.Data) == 0 {
		err = edenlabs.ErrorInvalid("price_level_id")
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

	var salesTerritory *bridgeService.GetSalesTerritoryGPResponse
	salesTerritory, err = s.opt.Client.BridgeServiceGrpc.GetSalesTerritoryGPDetail(ctx, &bridgeService.GetSalesTerritoryGPDetailRequest{
		Id: req.SalesTerritoryID,
	})
	if err != nil || len(salesTerritory.Data) == 0 {
		err = edenlabs.ErrorInvalid("sales_territory_id")
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

	var customerClass *bridgeService.GetCustomerClassResponse
	customerClass, err = s.opt.Client.BridgeServiceGrpc.GetCustomerClassDetail(ctx, &bridgeService.GetCustomerClassDetailRequest{
		Id: req.CustomerClassID,
	})
	if err != nil || len(customerClass.Data) == 0 {
		err = edenlabs.ErrorInvalid("customer_class_id")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	customerClassResponse = &dto.CustomerClassResponse{
		ID:          customerClass.Data[0].Classid,
		Code:        customerClass.Data[0].Classid,
		Description: customerClass.Data[0].Clasdscr,
	}

	var referenceInfo *configurationService.GetGlossaryDetailResponse
	referenceInfo, err = s.opt.Client.ConfigurationServiceGrpc.GetGlossaryDetail(ctx, &configurationService.GetGlossaryDetailRequest{
		Table:     "all",
		Attribute: "reference_info",
		ValueInt:  int32(req.ReferenceInfo),
	})
	if err != nil {
		err = edenlabs.ErrorInvalid("reference_info")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	referenceInfoResponse = &dto.GlossaryResponse{
		ID:        int64(referenceInfo.Data.Id),
		Table:     referenceInfo.Data.Table,
		Attribute: referenceInfo.Data.Attribute,
		ValueInt:  int8(referenceInfo.Data.ValueInt),
		ValueName: referenceInfo.Data.ValueName,
		Note:      referenceInfo.Data.Note,
	}

	var timeConcent *configurationService.GetGlossaryDetailResponse
	timeConcent, err = s.opt.Client.ConfigurationServiceGrpc.GetGlossaryDetail(ctx, &configurationService.GetGlossaryDetailRequest{
		Table:     "prospect_customer",
		Attribute: "time_consent",
		ValueInt:  int32(req.TimeConsent),
	})
	if err != nil {
		err = edenlabs.ErrorInvalid("time_consent")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	timeConcentResponse = &dto.GlossaryResponse{
		ID:        int64(timeConcent.Data.Id),
		Table:     timeConcent.Data.Table,
		Attribute: timeConcent.Data.Attribute,
		ValueInt:  int8(timeConcent.Data.ValueInt),
		ValueName: timeConcent.Data.ValueName,
		Note:      timeConcent.Data.Note,
	}

	var invoiceTerm *configurationService.GetGlossaryDetailResponse
	invoiceTerm, err = s.opt.Client.ConfigurationServiceGrpc.GetGlossaryDetail(ctx, &configurationService.GetGlossaryDetailRequest{
		Table:     "prospect_customer",
		Attribute: "invoice_term",
		ValueInt:  int32(req.InvoiceTerm),
	})
	if err != nil {
		err = edenlabs.ErrorInvalid("invoice_term")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	invoiceTermResponse = &dto.GlossaryResponse{
		ID:        int64(invoiceTerm.Data.Id),
		Table:     invoiceTerm.Data.Table,
		Attribute: invoiceTerm.Data.Attribute,
		ValueInt:  int8(invoiceTerm.Data.ValueInt),
		ValueName: invoiceTerm.Data.ValueName,
		Note:      invoiceTerm.Data.Note,
	}

	var application *configurationService.GetGlossaryDetailResponse
	application, err = s.opt.Client.ConfigurationServiceGrpc.GetGlossaryDetail(ctx, &configurationService.GetGlossaryDetailRequest{
		Table:     "all",
		Attribute: "application",
		ValueInt:  int32(req.RegistrationChannel),
	})
	if err != nil {
		err = edenlabs.ErrorInvalid("application")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	applicationResponse = &dto.GlossaryResponse{
		ID:        int64(application.Data.Id),
		Table:     application.Data.Table,
		Attribute: application.Data.Attribute,
		ValueInt:  int8(application.Data.ValueInt),
		ValueName: application.Data.ValueName,
		Note:      application.Data.Note,
	}

	var prospectiveCustomerAddress *model.ProspectiveCustomerAddress
	if req.CompanyAddressID != 0 {
		prospectiveCustomerAddress, err = s.RepositoryProspectiveCustomerAddress.GetDetail(ctx, &dto.ProspectiveCustomerAddressGetDetailRequest{ID: req.CompanyAddressID})
		if err != nil || prospectiveCustomerAddress.ProspectiveCustomerID != req.ProspectiveCustomerID {
			err = edenlabs.ErrorInvalid("company_address_id")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	}

	if req.ShippingAddressID != 0 {
		prospectiveCustomerAddress, err = s.RepositoryProspectiveCustomerAddress.GetDetail(ctx, &dto.ProspectiveCustomerAddressGetDetailRequest{ID: req.ShippingAddressID})
		if err != nil || prospectiveCustomerAddress.ProspectiveCustomerID != req.ProspectiveCustomerID {
			err = edenlabs.ErrorInvalid("shipping_address_id")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	}

	if req.BillingAddressID != 0 {
		prospectiveCustomerAddress, err = s.RepositoryProspectiveCustomerAddress.GetDetail(ctx, &dto.ProspectiveCustomerAddressGetDetailRequest{ID: req.BillingAddressID})
		if err != nil || prospectiveCustomerAddress.ProspectiveCustomerID != req.ProspectiveCustomerID {
			err = edenlabs.ErrorInvalid("billing_address_id")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	}

	// Change character comma to point
	req.CompanyAddressLatitude = strings.ReplaceAll(req.CompanyAddressLatitude, ",", ".")
	req.CompanyAddressLongitude = strings.ReplaceAll(req.CompanyAddressLongitude, ",", ".")
	req.ShippingAddressLatitude = strings.ReplaceAll(req.ShippingAddressLatitude, ",", ".")
	req.ShippingAddressLongitude = strings.ReplaceAll(req.ShippingAddressLongitude, ",", ".")
	req.BillingAddressLatitude = strings.ReplaceAll(req.BillingAddressLatitude, ",", ".")
	req.BillingAddressLongitude = strings.ReplaceAll(req.BillingAddressLongitude, ",", ".")

	//
	companyAddressLatitude = utils.ToFloat(req.CompanyAddressLatitude)
	companyAddressLongitude = utils.ToFloat(req.CompanyAddressLongitude)
	shippingAddressLatitude = utils.ToFloat(req.ShippingAddressLatitude)
	shippingAddressLongitude = utils.ToFloat(req.ShippingAddressLongitude)
	billingAddressLatitude = utils.ToFloat(req.BillingAddressLatitude)
	billingAddressLongitude = utils.ToFloat(req.BillingAddressLongitude)

	// Validation Value Latitude Shipping Address
	if shippingAddressLatitude < -90 || shippingAddressLatitude > 90 {
		err = edenlabs.ErrorValidation("shipping_address_latitude", "Latitude value must more than equal -90 and less than equal 90")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	// Validation Value Longitude Shipping Address
	if shippingAddressLongitude < -180 || shippingAddressLongitude > 180 {
		err = edenlabs.ErrorValidation("shipping_address_longitude", "Longitude value must more than equal -180 and less than equal 180")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// Validation Value Latitude Billing Address
	if billingAddressLatitude < -90 || billingAddressLatitude > 90 {
		err = edenlabs.ErrorValidation("billing_address_latitude", "Latitude value must more than equal -90 and less than equal 90")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// Validation Value Longitude Billing Address
	if billingAddressLongitude < -180 || billingAddressLongitude > 180 {
		err = edenlabs.ErrorValidation("billing_address_longitude", "Longitude value must more than equal -180 and less than equal 180")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if req.BusinessTypeID == 1 {
		// Validation Value Latitude Company Address
		if companyAddressLatitude < -90 || companyAddressLatitude > 90 {
			err = edenlabs.ErrorValidation("company_address_latitude", "Latitude value must more than equal -90 and less than equal 90")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		// Validation Value Longitude Company Address
		if companyAddressLongitude < -180 || companyAddressLongitude > 180 {
			err = edenlabs.ErrorValidation("company_address_longitude", "Longitude value must more than equal -180 and less than equal 180")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		// Company Address
		admDivisionCompany, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPList(ctx, &bridgeService.GetAdmDivisionGPListRequest{
			Limit:       1,
			Offset:      0,
			State:       req.CompanyAddressProvince,
			SubDistrict: req.CompanyAddressSubDistrict,
			District:    req.CompanyAddressDistrict,
			City:        req.CompanyAddressCity})

		if err != nil || len(admDivisionCompany.Data) == 0 {
			err = edenlabs.ErrorInvalid("company_address_sub_district")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		// Set region id gp
		regionIDGP = admDivisionCompany.Data[0].Region

		companyAddress = &model.ProspectiveCustomerAddress{
			ID:                  int64(req.CompanyAddressID),
			AddressName:         req.CompanyAddressName,
			AddressType:         "statement_to",
			Address1:            req.CompanyAddressDetail1,
			Address2:            req.CompanyAddressDetail2,
			Address3:            req.CompanyAddressDetail3,
			AdmDivisionIDGP:     admDivisionCompany.Data[0].Code,
			Latitude:            utils.ToFloat(req.CompanyAddressLatitude),
			Longitude:           utils.ToFloat(req.CompanyAddressLongitude),
			City:                admDivisionCompany.Data[0].City,
			State:               admDivisionCompany.Data[0].State,
			IsCreateAddressToGP: true,
			PostalCode:          admDivisionCompany.Data[0].Zipcode,
			Note:                req.CompanyAddressNote,
		}

		addressList = append(addressList, companyAddress)
	}

	if req.ShippingAddressReferTo == 1 && req.BusinessTypeID == 1 {
		shippingAddress = &model.ProspectiveCustomerAddress{
			ID:                  int64(req.ShippingAddressID),
			AddressName:         req.CompanyAddressName,
			AddressType:         "ship_to",
			Address1:            req.CompanyAddressDetail1,
			Address2:            req.CompanyAddressDetail2,
			Address3:            req.CompanyAddressDetail3,
			AdmDivisionIDGP:     admDivisionCompany.Data[0].Code,
			Latitude:            utils.ToFloat(req.CompanyAddressLatitude),
			Longitude:           utils.ToFloat(req.CompanyAddressLongitude),
			City:                admDivisionCompany.Data[0].City,
			State:               admDivisionCompany.Data[0].State,
			ReferTo:             req.ShippingAddressReferTo,
			IsCreateAddressToGP: false,
			PostalCode:          admDivisionCompany.Data[0].Zipcode,
			Note:                req.CompanyAddressNote,
		}
		addressList = append(addressList, shippingAddress)
	} else {
		admDivisionShippingAddress, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPList(ctx, &bridgeService.GetAdmDivisionGPListRequest{
			Limit:       1,
			Offset:      0,
			State:       req.ShippingAddressProvince,
			City:        req.ShippingAddressCity,
			District:    req.ShippingAddressDistrict,
			SubDistrict: req.ShippingAddressSubDistrict,
		})

		if err != nil || len(admDivisionShippingAddress.Data) == 0 {
			err = edenlabs.ErrorInvalid("shipping_address_sub_district")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		// Set region id gp
		regionIDGP = admDivisionShippingAddress.Data[0].Region

		shippingAddress = &model.ProspectiveCustomerAddress{
			ID:                  int64(req.ShippingAddressID),
			AddressName:         req.ShippingAddressName,
			AddressType:         "ship_to",
			Address1:            req.ShippingAddressDetail1,
			Address2:            req.ShippingAddressDetail2,
			Address3:            req.ShippingAddressDetail3,
			AdmDivisionIDGP:     admDivisionShippingAddress.Data[0].Code,
			Latitude:            utils.ToFloat(req.ShippingAddressLatitude),
			Longitude:           utils.ToFloat(req.ShippingAddressLongitude),
			City:                admDivisionShippingAddress.Data[0].City,
			State:               admDivisionShippingAddress.Data[0].State,
			IsCreateAddressToGP: true,
			PostalCode:          admDivisionShippingAddress.Data[0].Zipcode,
			Note:                req.ShippingAddressNote,
		}

		addressList = append(addressList, shippingAddress)
	}

	var admDivisionCoverage *bridgeService.GetAdmDivisionCoverageGPResponse
	admDivisionCoverage, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionCoverageGPList(ctx, &bridgeService.GetAdmDivisionCoverageGPListRequest{
		Limit:                 1,
		Offset:                0,
		GnlAdministrativeCode: shippingAddress.AdmDivisionIDGP,
	})

	if err != nil || len(admDivisionCoverage.Data) == 0 {
		err = edenlabs.ErrorInvalid("shipping_address_sub_district")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	siteResponse = &dto.SiteResponse{
		ID: admDivisionCoverage.Data[0].Locncode,
	}

	if (req.BillingAddressReferTo == 1 || req.BillingAddressReferTo == 3) && req.BusinessTypeID == 1 {
		billingAddress = &model.ProspectiveCustomerAddress{
			ID:                  int64(req.BillingAddressID),
			AddressName:         req.CompanyAddressName,
			AddressType:         "bill_to",
			Address1:            req.CompanyAddressDetail1,
			Address2:            req.CompanyAddressDetail2,
			Address3:            req.CompanyAddressDetail3,
			AdmDivisionIDGP:     admDivisionCompany.Data[0].Code,
			Latitude:            utils.ToFloat(req.CompanyAddressLatitude),
			Longitude:           utils.ToFloat(req.CompanyAddressLongitude),
			City:                admDivisionCompany.Data[0].City,
			State:               admDivisionCompany.Data[0].State,
			ReferTo:             req.BillingAddressReferTo,
			IsCreateAddressToGP: false,
			PostalCode:          admDivisionCompany.Data[0].Zipcode,
			Note:                req.CompanyAddressNote,
		}
		addressList = append(addressList, billingAddress)
	} else if req.BillingAddressReferTo == 2 {
		billingAddress = &model.ProspectiveCustomerAddress{
			ID:                  int64(req.BillingAddressID),
			AddressName:         req.ShippingAddressName,
			AddressType:         "bill_to",
			Address1:            req.ShippingAddressDetail1,
			Address2:            req.ShippingAddressDetail2,
			Address3:            req.ShippingAddressDetail3,
			AdmDivisionIDGP:     admDivisionShippingAddress.Data[0].Code,
			Latitude:            utils.ToFloat(req.ShippingAddressLatitude),
			Longitude:           utils.ToFloat(req.ShippingAddressLongitude),
			City:                admDivisionShippingAddress.Data[0].City,
			State:               admDivisionShippingAddress.Data[0].State,
			ReferTo:             req.BillingAddressReferTo,
			IsCreateAddressToGP: false,
			PostalCode:          admDivisionShippingAddress.Data[0].Zipcode,
			Note:                req.ShippingAddressNote,
		}
		addressList = append(addressList, billingAddress)
	} else {
		// Bill to Address
		admDivisionBillingAddress, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPList(ctx, &bridgeService.GetAdmDivisionGPListRequest{
			Limit:       1,
			Offset:      0,
			State:       req.BillingAddressProvince,
			City:        req.BillingAddressCity,
			District:    req.BillingAddressDistrict,
			SubDistrict: req.BillingAddressSubDistrict,
		})

		if err != nil || len(admDivisionBillingAddress.Data) == 0 {
			err = edenlabs.ErrorInvalid("billing_address_sub_district")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		billingAddress = &model.ProspectiveCustomerAddress{
			ID:                  int64(req.BillingAddressID),
			AddressName:         req.BillingAddressName,
			AddressType:         "bill_to",
			Address1:            req.BillingAddressDetail1,
			Address2:            req.BillingAddressDetail2,
			Address3:            req.BillingAddressDetail3,
			AdmDivisionIDGP:     admDivisionBillingAddress.Data[0].Code,
			Latitude:            utils.ToFloat(req.BillingAddressLatitude),
			Longitude:           utils.ToFloat(req.BillingAddressLongitude),
			City:                admDivisionBillingAddress.Data[0].City,
			State:               admDivisionBillingAddress.Data[0].State,
			IsCreateAddressToGP: true,
			PostalCode:          admDivisionBillingAddress.Data[0].Zipcode,
			Note:                req.BillingAddressNote,
		}
		addressList = append(addressList, billingAddress)
	}

	prospectiveCustomer = &model.ProspectiveCustomer{
		CustomerClassIDGP:         req.CustomerClassID,
		SalespersonIDGP:           req.SalespersonID,
		ArchetypeIDGP:             req.ArchetypeID,
		CustomerTypeIDGP:          req.CustomerTypeID,
		SalesTerritoryIDGP:        req.SalesTerritoryID,
		SalesPriceLevelIDGP:       req.PriceLevelID,
		SiteIDGP:                  siteResponse.ID,
		RegionIDGP:                regionIDGP,
		ShippingMethodIDGP:        req.ShippingMethodID,
		CustomerIDGP:              req.CustomerCode,
		BusinessName:              req.BusinessName,
		RegStatus:                 statusx.ConvertStatusName(statusx.Registered),
		ProcessedAt:               time.Now(),
		ProcessedBy:               userID,
		BrandName:                 req.BrandName,
		Application:               req.RegistrationChannel,
		OutletImage:               utils.ArrayStringToString(req.OutletImage),
		TimeConsent:               req.TimeConsent,
		ReferenceInfo:             req.ReferenceInfo,
		ReferrerCode:              req.ReferrerCode,
		OwnerName:                 req.OwnerName,
		OwnerRole:                 req.OwnerRole,
		Email:                     req.Email,
		BusinessTypeIDGP:          req.BusinessTypeID,
		PicOrderName:              req.PicOrderName,
		PicOrderContact:           req.PicOrderContact,
		PicFinanceContact:         req.PicFinanceContact,
		PicFinanceName:            req.PicFinanceName,
		IDCardDocName:             "ID-Card.pdf",
		IDCardDocNumber:           req.IDCardDocNumber,
		IDCardDocURL:              req.IDCardDocURL,
		TaxpayerDocName:           "NPWP.pdf",
		TaxpayerDocNumber:         req.TaxpayerDocNumber,
		TaxpayerDocURL:            req.TaxpayerDocURL,
		CompanyContractDocName:    "Company-Contract.pdf",
		CompanyContractDocURL:     req.CompanyContractDocURL,
		NotarialDeedDocName:       "Notarial-Deed.pdf",
		NotarialDeedDocURL:        req.NotarialDeedDocURL,
		TaxableEntrepeneurDocName: "Tax-Entrepreneur.pdf",
		TaxableEntrepeneurDocURL:  req.TaxableEntrepeneurDocURL,
		CompanyCertificateRegName: "Company-Certificate.pdf",
		CompanyCertificateRegURL:  req.CompanyCertificateRegURL,
		BusinessLicenseDocName:    "Business-License.pdf",
		BusinessLicenseDocURL:     req.BusinessLicenseDocURL,
		PaymentTermIDGP:           req.PaymentTermID,
		ExchangeInvoice:           req.ExchangeInvoice,
		ExchangeInvoiceTime:       req.ExchangeInvoiceTime,
		FinanceEmail:              req.FinanceEmail,
		InvoiceTerm:               req.InvoiceTerm,
		Comment1:                  req.Comment1,
		Comment2:                  req.Comment2,
		PicOperationName:          req.PicOperationName,
		PicOperationContact:       req.PicOperationContact,
		OwnerContact:              req.OwnerContact,
	}

	if req.ProspectiveCustomerID != 0 {
		prospectiveCustomer.ID = prospectiveCustomerOld.ID
		prospectiveCustomer.Code = prospectiveCustomerOld.Code
		prospectiveCustomer.UpdatedAt = time.Now()
		prospectiveCustomer.CreatedAt = prospectiveCustomerOld.CreatedAt
		err = s.RepositoryProspectiveCustomer.Update(ctx, prospectiveCustomer)
		if err != nil {
			err = edenlabs.ErrorInvalid("prospective_customer_id")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	} else {
		var codeGenerator *configurationService.GetGenerateCodeResponse
		codeGenerator, err = s.opt.Client.ConfigurationServiceGrpc.GetGenerateCode(ctx, &configurationService.GetGenerateCodeRequest{
			Format: "PCT",
			Domain: "prospective_customer",
			Length: 6,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("configuration", "generate_code")
			return
		}

		prospectiveCustomer.Code = codeGenerator.Data.Code
		prospectiveCustomer.CreatedAt = time.Now()
		_, err = s.RepositoryProspectiveCustomer.Create(ctx, prospectiveCustomer)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	}

	var addressListGP *bridgeService.GetAddressGPResponse
	addressListGP, err = s.opt.Client.BridgeServiceGrpc.GetAddressGPList(ctx, &bridgeService.GetAddressGPListRequest{
		Limit:          100,
		Offset:         0,
		Status:         "0",
		CustomerNumber: req.CustomerCode,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// Non Active Existing Address
	for _, v := range addressListGP.Data {
		_, err = s.opt.Client.BridgeServiceGrpc.UpdateAddress(ctx, &bridge_service.UpdateAddressRequest{
			Custnmbr:                v.Custnmbr,
			Custname:                v.Custname,
			Adrscode:                v.Adrscode,
			Cntcprsn:                v.Cntcprsn,
			ShipToName:              v.ShipToName,
			AddresS1:                v.AddresS1,
			AddresS2:                v.AddresS2,
			AddresS3:                v.AddresS3,
			GnL_Address_Note:        v.GnL_Address_Note,
			Country:                 v.Country,
			City:                    v.City,
			State:                   v.State,
			PhonE1:                  v.PhonE1,
			Inactive:                1,
			GnL_Administrative_Code: v.AdministrativeDiv.GnlAdministrativeCode,
			Locncode:                v.Locncode,
			TypeAddress:             v.TypeAddress,
			GnL_Latitude:            v.GnL_Latitude,
			GnL_Longitude:           v.GnL_Longitude,
			GnL_Archetype_ID:        v.GnL_Archetype_ID,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	}

	// Create New Address
	for _, address := range addressList {
		address.ProspectiveCustomerID = prospectiveCustomer.ID
		if address.ID != 0 {
			address.ProspectiveCustomerID = req.ProspectiveCustomerID
			address.UpdatedAt = time.Now()
			err = s.RepositoryProspectiveCustomerAddress.Update(ctx, address)
			if err != nil {
				err = edenlabs.ErrorInvalid("address_id")
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}
		} else {
			address.CreatedAt = time.Now()
			_, err = s.RepositoryProspectiveCustomerAddress.Create(ctx, address)
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}
		}

		if address.IsCreateAddressToGP {
			var (
				codeGenerator              *configurationService.GetGenerateCodeResponse
				contactPerson, phonePerson string
			)
			codeGenerator, err = s.opt.Client.ConfigurationServiceGrpc.GetGenerateCode(ctx, &configurationService.GetGenerateCodeRequest{
				Format: customer.Data[0].Custnmbr + "-",
				Domain: "address",
				Length: 3,
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("configuration", "generate_code")
				return
			}

			switch address.AddressType {
			case "statement_to":
				addressCodeCompanyAddress = codeGenerator.Data.Code
				phonePerson = req.PicOrderContact
				contactPerson = req.PicOrderName
			case "ship_to":
				addressCodeShippingAddress = codeGenerator.Data.Code
				phonePerson = req.PicOrderContact
				contactPerson = req.PicOrderName
			default:
				AddressCodeBillingAddress = codeGenerator.Data.Code
				phonePerson = req.PicFinanceContact
				contactPerson = req.PicFinanceName
			}

			_, err = s.opt.Client.BridgeServiceGrpc.CreateAddress(ctx, &bridge_service.CreateAddressRequest{
				Custnmbr:                customer.Data[0].Custnmbr,
				Custname:                req.BusinessName,
				Cntcprsn:                contactPerson,
				ShipToName:              address.AddressName,
				Adrscode:                codeGenerator.Data.Code,
				AddresS1:                address.Address1,
				AddresS2:                address.Address2,
				AddresS3:                address.Address3,
				GnL_Address_Note:        address.Note,
				Country:                 "Indonesia",
				City:                    address.City,
				State:                   address.State,
				PhonE1:                  customer.Data[0].PhonE1,
				PhonE2:                  phonePerson,
				Inactive:                "0",
				GnL_Administrative_Code: address.AdmDivisionIDGP,
				Locncode:                admDivisionCoverage.Data[0].Locncode,
				TypeAddress:             address.AddressType,
				GnL_Latitude:            address.Latitude,
				GnL_Longitude:           address.Longitude,
				GnL_Archetype_ID:        req.ArchetypeID,
				Slprsnid:                req.SalespersonID,
				Shipmthd:                req.ShippingMethodID,
				Salsterr:                req.SalesTerritoryID,
				Zip:                     address.PostalCode,
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}
		}
	}

	if req.ShippingAddressReferTo == 1 {
		addressCodeShippingAddress = addressCodeCompanyAddress
	}

	if req.BillingAddressReferTo == 1 || req.BillingAddressReferTo == 3 {
		AddressCodeBillingAddress = addressCodeCompanyAddress
	} else if req.BillingAddressReferTo == 2 {
		AddressCodeBillingAddress = addressCodeShippingAddress
	}

	var shipComplete int8
	if req.InvoiceTerm == 1 {
		shipComplete = 1
	} else {
		shipComplete = 0
	}

	_, err = s.opt.Client.BridgeServiceGrpc.UpdateCustomerGP(ctx, &bridgeService.UpdateCustomerGPRequest{
		Custnmbr: req.CustomerCode,
		Custname: req.BusinessName,
		Custclas: req.CustomerClassID,
		Address: &bridge_service.UpdateAddressRequest{
			Cntcprsn: req.PicOrderName,
			Adrscode: addressCodeShippingAddress,
			AddresS1: shippingAddress.Address1,
			AddresS2: shippingAddress.Address2,
			AddresS3: shippingAddress.Address3,
			Country:  "Indonesia",
			State:    shippingAddress.State,
			City:     shippingAddress.City,
			Zip:      shippingAddress.PostalCode,
			PhonE1:   customer.Data[0].PhonE1,
			PhonE2:   req.PicOrderContact,
		},
		Shipmthd:        req.ShippingMethodID,
		Stmtname:        req.BrandName,
		Slprsnid:        req.SalespersonID,
		Salsterr:        req.SalesTerritoryID,
		Comment1:        req.Comment1,
		Comment2:        req.Comment2,
		Prclevel:        req.PriceLevelID,
		Shipcomplete:    utils.ToString(shipComplete),
		GnlCustTypeId:   req.CustomerTypeID,
		GnlReferrerCode: req.ReferrerCode,
		GnlBusinessType: int32(req.BusinessTypeID),
		Staddrcd:        addressCodeCompanyAddress,
		Prbtadcd:        AddressCodeBillingAddress,
		Prstadcd:        addressCodeShippingAddress,
		GnlSocialSecNum: req.IDCardDocNumber,
		Pymtrmid:        req.PaymentTermID,
	})

	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// Update field upgrade status on customer
	if req.CustomerCode != "" {
		customerInternal.UpgradeStatus = statusx.ConvertStatusName(statusx.Approved)
		customerInternal.UpdatedAt = time.Now()
		customerInternal.ProspectiveCustomerID = prospectiveCustomer.ID

		var eligibleMembership *configurationService.GetConfigAppDetailResponse
		eligibleMembership, err = s.opt.Client.ConfigurationServiceGrpc.GetConfigAppDetail(ctx, &configurationService.GetConfigAppDetailRequest{
			Attribute: "eligible_membership_business_type",
		})
		if err != nil {
			err = edenlabs.ErrorInvalid("config_app")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		// Check the customer type is eligible for membership campaign
		if strings.Contains(eligibleMembership.Data.Value, req.CustomerTypeID) {
			customerInternal.MembershipLevelID = 1
			customerInternal.MembershipCheckpointID = 1
		}

		err = s.RepositoryCustomer.Update(ctx, customerInternal, "UpgradeStatus", "UpdatedAt", "MembershipLevelID", "MembershipCheckpointID", "ProspectiveCustomerID")
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	}

	_, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &auditService.CreateLogRequest{
		Log: &auditService.Log{
			UserId:      userID,
			ReferenceId: utils.ToString(prospectiveCustomer.ID),
			Type:        "prospective_customer",
			Function:    "upgrade",
			CreatedAt:   timestamppb.New(time.Now()),
		},
	})

	res = &dto.ProspectiveCustomerResponse{
		ID:                        prospectiveCustomer.ID,
		Code:                      prospectiveCustomer.Code,
		CustomerClass:             customerClassResponse,
		Salesperson:               salespersonResponse,
		Archetype:                 archetypeResponse,
		CustomerType:              customerTypeResponse,
		SalesTerritory:            salesTerritoryResponse,
		PriceLevel:                priceLevelResponse,
		Site:                      siteResponse,
		ShippingMethod:            shippingMethodResponse,
		Customer:                  customerResponse,
		BusinessName:              prospectiveCustomer.BusinessName,
		RegStatus:                 prospectiveCustomer.RegStatus,
		RegStatusConvert:          statusx.ConvertStatusValue(prospectiveCustomer.RegStatus),
		CreatedAt:                 timex.ToLocTime(ctx, prospectiveCustomer.CreatedAt),
		UpdatedAt:                 timex.ToLocTime(ctx, prospectiveCustomer.UpdatedAt),
		ProcessedAt:               timex.ToLocTime(ctx, prospectiveCustomer.ProcessedAt),
		ProcessedBy:               createdByResponse,
		DeclineType:               prospectiveCustomer.DeclineType,
		DeclineNote:               prospectiveCustomer.DeclineNote,
		BrandName:                 prospectiveCustomer.BrandName,
		Application:               applicationResponse,
		OutletImage:               utils.StringToStringArray(prospectiveCustomer.OutletImage),
		TimeConsent:               timeConcentResponse,
		ReferenceInfo:             referenceInfoResponse,
		ReferrerCode:              prospectiveCustomer.ReferrerCode,
		OwnerName:                 prospectiveCustomer.OwnerName,
		OwnerRole:                 prospectiveCustomer.OwnerRole,
		Email:                     prospectiveCustomer.Email,
		BusinessType:              businessTypeResponse,
		PicOrderName:              prospectiveCustomer.PicOrderName,
		PicOrderContact:           prospectiveCustomer.PicOrderContact,
		PicFinanceContact:         prospectiveCustomer.PicFinanceContact,
		PicFinanceName:            prospectiveCustomer.PicFinanceName,
		IDCardDocName:             prospectiveCustomer.IDCardDocName,
		IDCardDocNumber:           prospectiveCustomer.IDCardDocNumber,
		IDCardDocURL:              prospectiveCustomer.IDCardDocURL,
		TaxpayerDocName:           prospectiveCustomer.TaxpayerDocName,
		TaxpayerDocNumber:         prospectiveCustomer.TaxpayerDocNumber,
		TaxpayerDocURL:            prospectiveCustomer.TaxpayerDocURL,
		CompanyContractDocName:    prospectiveCustomer.CompanyContractDocName,
		CompanyContractDocURL:     prospectiveCustomer.CompanyContractDocURL,
		NotarialDeedDocName:       prospectiveCustomer.NotarialDeedDocName,
		NotarialDeedDocURL:        prospectiveCustomer.NotarialDeedDocURL,
		TaxableEntrepeneurDocName: prospectiveCustomer.TaxableEntrepeneurDocName,
		TaxableEntrepeneurDocURL:  prospectiveCustomer.TaxableEntrepeneurDocURL,
		CompanyCertificateRegName: prospectiveCustomer.CompanyCertificateRegName,
		CompanyCertificateRegURL:  prospectiveCustomer.CompanyCertificateRegURL,
		BusinessLicenseDocName:    prospectiveCustomer.BusinessLicenseDocName,
		BusinessLicenseDocURL:     prospectiveCustomer.BusinessLicenseDocURL,
		PaymentTerm:               paymentTermResponse,
		ExchangeInvoice:           prospectiveCustomer.ExchangeInvoice,
		ExchangeInvoiceTime:       prospectiveCustomer.ExchangeInvoiceTime,
		FinanceEmail:              prospectiveCustomer.FinanceEmail,
		InvoiceTerm:               invoiceTermResponse,
		Comment1:                  prospectiveCustomer.Comment1,
		Comment2:                  prospectiveCustomer.Comment2,
		OwnerContact:              prospectiveCustomer.OwnerContact,
	}

	return
}
