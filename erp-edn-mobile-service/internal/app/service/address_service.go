package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/utils"
	"git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/global"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/internal/app/dto"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

func NewServiceAddress() IAddressService {
	m := new(AddressService)
	m.opt = global.Setup.Common
	return m
}

type IAddressService interface {
	GetAddresss(ctx context.Context, req dto.AddressListRequest) (res []*dto.Address, total int64, err error)
	GetAddressDetailById(ctx context.Context, req dto.AddressDetailRequest) (res *dto.AddressResponse, err error)
	GetListGp(ctx context.Context, req dto.GetAddressGPListRequest) (res []*dto.AddressGP, total int64, err error)
	GetDetailGp(ctx context.Context, id string) (res *dto.AddressGP, err error)
}

type AddressService struct {
	opt opt.Options
}

func (s *AddressService) GetAddresss(ctx context.Context, req dto.AddressListRequest) (res []*dto.Address, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AddressService.GetAddresss")
	defer span.End()

	// get adddress from bridge
	var addrResponse *bridgeService.GetAddressGPResponse
	addrResponse, err = s.opt.Client.BridgeServiceGrpc.GetAddressGPList(ctx, &bridgeService.GetAddressGPListRequest{
		Limit:    req.Limit,
		Offset:   req.Offset,
		Status:   utils.ToString(req.Status),
		Adrscode: req.Search,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "address")
		return
	}

	datas := []*dto.Address{}
	for _, addr := range addrResponse.Data {
		mainbranch := 0
		if addr.TypeAddress == "ship_to" {
			mainbranch = 1
		}
		datas = append(datas, &dto.Address{
			Code:            addr.Adrscode,
			Name:            addr.ShipToName,
			PicName:         addr.Custname,
			PhoneNumber:     addr.PhonE1,
			AltPhoneNumber:  addr.PhonE2,
			AddressName:     addr.ShipToName,
			ShippingAddress: addr.AddresS1 + " " + addr.AddresS2 + " " + addr.AddresS3,
			Latitude:        &addr.GnL_Latitude,
			Longitude:       &addr.GnL_Longitude,
			Note:            "",
			MainBranch:      int8(mainbranch),
			Status:          int8(addr.Inactive),
			// CreatedAt:       time.Time{},
			// CreatedBy:       0,
			// LastUpdatedAt:   time.Time{},
			// LastUpdatedBy:   0,
			City:          addr.AdministrativeDiv.GnlCity,
			State:         addr.State,
			ZipCode:       addr.Zip,
			CountryCode:   addr.CCode,
			Country:       addr.Country,
			UpsZone:       "",
			CustomerID:    addr.Custnmbr,
			CustomerName:  addr.Custname,
			RegionID:      addr.AdministrativeDiv.GnlRegion,
			ArchetypeID:   addr.GnL_Archetype_ID,
			SiteID:        addr.Locncode,
			AdmDivisionID: addr.AdministrativeDiv.GnlAdministrativeCode,
			ContactPerson: addr.Cntcprsn,
			StatusConvert: "",
		})
	}
	res = datas
	total = int64(addrResponse.TotalRecords)

	return
}

func (s *AddressService) GetAddressDetailById(ctx context.Context, req dto.AddressDetailRequest) (res *dto.AddressResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AddressService.GetAddresss")
	defer span.End()

	// get address from bridge
	var addrResponse *bridgeService.GetAddressGPResponse
	addrResponse, err = s.opt.Client.BridgeServiceGrpc.GetAddressGPDetail(ctx, &bridgeService.GetAddressGPDetailRequest{
		Id: req.Id,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "address")
		return
	}

	res = &dto.AddressResponse{
		// ID:               addrResponse.Data[0].Id,
		Code:          addrResponse.Data[0].Adrscode,
		CustomerName:  addrResponse.Data[0].Custname,
		ArchetypeID:   addrResponse.Data[0].GnL_Archetype_ID,
		AdmDivisionID: addrResponse.Data[0].GnL_Administrative_Code,
		SiteID:        addrResponse.Data[0].Locncode,
		// SalespersonID:    addrResponse.Data[0].SalespersonId,
		// TerritoryID:      addrResponse.Data[0].TerritoryId,
		AddressCode:   addrResponse.Data[0].Adrscode,
		AddressName:   addrResponse.Data[0].ShipToName,
		ContactPerson: addrResponse.Data[0].Cntcprsn,
		City:          addrResponse.Data[0].City,
		State:         addrResponse.Data[0].State,
		ZipCode:       addrResponse.Data[0].Zip,
		CountryCode:   addrResponse.Data[0].CCode,
		Country:       addrResponse.Data[0].Country,
		Latitude:      &addrResponse.Data[0].GnL_Latitude,
		Longitude:     &addrResponse.Data[0].GnL_Longitude,
		// UpsZone:          addrResponse.Data[0].UpsZone,
		ShippingMethod: addrResponse.Data[0].Shipmthd,
		// TaxScheduleID:    addrResponse.Data[0].TaxScheduleId,
		// PrintPhoneNumber: int8(addrResponse.Data[0].PrintPhoneNumber),
		Phone1: addrResponse.Data[0].PhonE1,
		Phone2: addrResponse.Data[0].PhonE2,
		Phone3: addrResponse.Data[0].PhonE3,
		// FaxNumber:        addrResponse.Data[0].FaxNumber,
		ShippingAddress: addrResponse.Data[0].AddresS1 + " " + addrResponse.Data[0].AddresS2 + " " + addrResponse.Data[0].AddresS3,
		// BcaVa:           addrResponse.Data[0].BcaVa,
		// OtherVa:         addrResponse.Data[0].OtherVa,
		// Note:            addrResponse.Data[0].Note,
		Status: int8(addrResponse.Data[0].Inactive),
		// CreatedAt:       addrResponse.Data[0].CreatedAt.AsTime(),
		// UpdatedAt:       addrResponse.Data[0].UpdatedAt.AsTime(),
	}

	return
}

func (s *AddressService) GetListGp(ctx context.Context, req dto.GetAddressGPListRequest) (res []*dto.AddressGP, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AddressService.GetListGp")
	defer span.End()

	var address *bridgeService.GetAddressGPResponse

	if address, err = s.opt.Client.BridgeServiceGrpc.GetAddressGPList(ctx, &bridgeService.GetAddressGPListRequest{
		Limit:  int32(req.Limit),
		Offset: int32(req.Offset),
	}); err != nil || !address.Succeeded {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "address")
		return
	}

	for _, addr := range address.Data {
		res = append(res, &dto.AddressGP{
			CCode:                   addr.CCode,
			Custnmbr:                addr.Custnmbr,
			Custname:                addr.Custname,
			Cntcprsn:                addr.Cntcprsn,
			Adrscode:                addr.Adrscode,
			Shipmthd:                addr.Shipmthd,
			Taxschid:                addr.Taxschid,
			AddresS1:                addr.AddresS1,
			AddresS2:                addr.AddresS2,
			AddresS3:                addr.AddresS3,
			Country:                 addr.Country,
			City:                    addr.City,
			State:                   addr.State,
			Zip:                     addr.Zip,
			PhonE1:                  addr.PhonE1,
			PhonE2:                  addr.PhonE2,
			PhonE3:                  addr.PhonE3,
			Slprsnid:                addr.Slprsnid,
			UserdeF1:                addr.UserdeF1,
			UserdeF2:                addr.UserdeF2,
			Salsterr:                addr.Salsterr,
			Locncode:                addr.Locncode,
			ShipToName:              addr.ShipToName,
			GnL_Administrative_Code: addr.GnL_Administrative_Code,
			GnL_Archetype_ID:        addr.GnL_Archetype_ID,
			GnL_Longitude:           addr.GnL_Longitude,
			GnL_Latitude:            addr.GnL_Latitude,
			GnL_Address_Note:        addr.GnL_Address_Note,
			Crusrid:                 addr.Crusrid,
			Creatddt:                addr.Creatddt,
			Mdfusrid:                addr.Mdfusrid,
			Modifdt:                 addr.Modifdt,
			TypeAddress:             addr.TypeAddress,
			Inactive:                addr.Inactive,
		})
	}

	total = int64(len(address.Data))

	return
}

func (s *AddressService) GetDetailGp(ctx context.Context, id string) (res *dto.AddressGP, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AddressService.GetDetailGp")
	defer span.End()

	var (
		address      *bridgeService.GetAddressGPResponse
		customer     *bridgeService.GetCustomerGPResponse
		paymentTerm  *bridgeService.GetPaymentTermGPResponse
		customerType *bridgeService.GetCustomerTypeDetailResponse
		region       *bridgeService.GetRegionDetailResponse
		site         *bridgeService.GetSiteDetailResponse
	)

	if address, err = s.opt.Client.BridgeServiceGrpc.GetAddressGPDetail(ctx, &bridgeService.GetAddressGPDetailRequest{
		Id: id,
	}); err != nil || !address.Succeeded {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "address")
		return
	}

	if customer, err = s.opt.Client.BridgeServiceGrpc.GetCustomerGPDetail(ctx, &bridgeService.GetCustomerGPDetailRequest{
		Id: address.Data[0].Custnmbr,
	}); err != nil || !customer.Succeeded {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "customer")
		return
	}

	if customerType, err = s.opt.Client.BridgeServiceGrpc.GetCustomerTypeDetail(ctx, &bridgeService.GetCustomerTypeDetailRequest{
		Id: 1,
	}); err != nil || !customer.Succeeded {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "customer")
		return
	}

	pymtrmids := []*dto.CustomerGPpymtrmid{}
	if customer.Data[0].Pymtrmid != nil {
		for _, pymtrmid := range customer.Data[0].Pymtrmid {
			pymtrmids = append(pymtrmids, &dto.CustomerGPpymtrmid{
				Pymtrmid:              pymtrmid.Pymtrmid,
				CalculateDateFromDays: pymtrmid.CalculateDateFromDays,
			})
		}
	}

	if paymentTerm, err = s.opt.Client.BridgeServiceGrpc.GetPaymentTermGPDetail(ctx, &bridgeService.GetPaymentTermGPDetailRequest{
		Id: "2% 10/Net 30",
	}); err != nil || !paymentTerm.Succeeded {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "payment term")
		return
	}

	region, err = s.opt.Client.BridgeServiceGrpc.GetRegionDetail(ctx, &bridgeService.GetRegionDetailRequest{
		Id: 1,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "region")
		return
	}

	site, err = s.opt.Client.BridgeServiceGrpc.GetSiteDetail(ctx, &bridgeService.GetSiteDetailRequest{
		Id: 1,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "site")
		return
	}

	res = &dto.AddressGP{
		CCode:                   address.Data[0].CCode,
		Custnmbr:                address.Data[0].Custnmbr,
		Custname:                address.Data[0].Custname,
		Cntcprsn:                address.Data[0].Cntcprsn,
		Adrscode:                address.Data[0].Adrscode,
		Shipmthd:                address.Data[0].Shipmthd,
		Taxschid:                address.Data[0].Taxschid,
		AddresS1:                address.Data[0].AddresS1,
		AddresS2:                address.Data[0].AddresS2,
		AddresS3:                address.Data[0].AddresS3,
		Country:                 address.Data[0].Country,
		City:                    address.Data[0].City,
		State:                   address.Data[0].State,
		Zip:                     address.Data[0].Zip,
		PhonE1:                  address.Data[0].PhonE1,
		PhonE2:                  address.Data[0].PhonE2,
		PhonE3:                  address.Data[0].PhonE3,
		Slprsnid:                address.Data[0].Slprsnid,
		UserdeF1:                address.Data[0].UserdeF1,
		UserdeF2:                address.Data[0].UserdeF2,
		Salsterr:                address.Data[0].Salsterr,
		Locncode:                address.Data[0].Locncode,
		ShipToName:              address.Data[0].ShipToName,
		GnL_Administrative_Code: address.Data[0].GnL_Administrative_Code,
		GnL_Archetype_ID:        address.Data[0].GnL_Archetype_ID,
		GnL_Longitude:           address.Data[0].GnL_Longitude,
		GnL_Latitude:            address.Data[0].GnL_Latitude,
		GnL_Address_Note:        address.Data[0].GnL_Address_Note,
		Crusrid:                 address.Data[0].Crusrid,
		Creatddt:                address.Data[0].Creatddt,
		Mdfusrid:                address.Data[0].Mdfusrid,
		Modifdt:                 address.Data[0].Modifdt,
		TypeAddress:             address.Data[0].TypeAddress,
		Inactive:                address.Data[0].Inactive,
		Customer: &dto.CustomerGP{
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
			PaymentTermGP: &dto.SalesPaymentTermGPResponse{
				Pymtrmid:              paymentTerm.Data[0].Pymtrmid,
				Duetype:               paymentTerm.Data[0].Duetype,
				Duedesc:               paymentTerm.Data[0].Duedesc,
				Duedtds:               paymentTerm.Data[0].Discdtds,
				CalculateDateFrom:     paymentTerm.Data[0].CalculateDateFrom,
				CalculateDateFromDays: paymentTerm.Data[0].CalculateDateFromDays,
				Disctype:              paymentTerm.Data[0].Disctype,
				Discdtds:              paymentTerm.Data[0].Discdtds,
				Dsclctyp:              paymentTerm.Data[0].Dsclctyp,
				Dscdlram:              paymentTerm.Data[0].Dscdlram,
				Dscpctam:              paymentTerm.Data[0].Dscpctam,
				Salpurch:              paymentTerm.Data[0].Salpurch,
				Discntcb:              paymentTerm.Data[0].Discntcb,
				Freight:               paymentTerm.Data[0].Freight,
				Misc:                  paymentTerm.Data[0].Misc,
				Tax:                   paymentTerm.Data[0].Tax,
			},
			CustomerType: &dto.CustomerTypeResponse{
				ID:           customerType.Data.Id,
				Code:         customerType.Data.Code,
				Description:  customerType.Data.Description,
				GroupType:    customerType.Data.GroupType,
				Abbreviation: customerType.Data.Abbreviation,
				Status:       int8(customerType.Data.Status),
			},
			Region: &dto.RegionResponse{
				// ID:          region.Data.Id,
				Code:        region.Data.Code,
				Description: region.Data.Description,
				Status:      int8(region.Data.Status),
			},
			Site: &dto.SiteResponse{
				// ID:          site.Data.Id,
				Code:        site.Data.Code,
				Description: site.Data.Description,
				Name:        site.Data.Description,
				Status:      int8(site.Data.Status),
			},
		},
	}

	return
}
