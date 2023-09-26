package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *BridgeGrpcHandler) GetAddressList(ctx context.Context, req *bridgeService.GetAddressListRequest) (res *bridgeService.GetAddressListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetAddressList")
	defer span.End()

	var addresses []dto.AddressResponse
	addresses, _, err = h.ServicesAddress.GetListGRPC(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*bridgeService.Address
	for _, address := range addresses {
		data = append(data, &bridgeService.Address{
			Id:               address.ID,
			Code:             address.Code,
			CustomerName:     address.CustomerName,
			ArchetypeId:      address.ArchetypeID,
			AdmDivisionId:    address.AdmDivisionID,
			SiteId:           address.SiteID,
			SalespersonId:    address.SalespersonID,
			TerritoryId:      address.TerritoryID,
			AddressCode:      address.AddressCode,
			AddressName:      address.AddressName,
			ContactPerson:    address.ContactPerson,
			City:             address.City,
			State:            address.State,
			ZipCode:          address.ZipCode,
			CountryCode:      address.CountryCode,
			Country:          address.Country,
			Latitude:         &address.Latitude,
			Longitude:        &address.Longitude,
			UpsZone:          address.UpsZone,
			ShippingMethod:   address.ShippingMethod,
			TaxScheduleId:    address.TaxScheduleID,
			PrintPhoneNumber: int32(address.PrintPhoneNumber),
			Phone_1:          address.Phone1,
			Phone_2:          address.Phone2,
			Phone_3:          address.Phone3,
			FaxNumber:        address.FaxNumber,
			ShippingAddress:  address.ShippingAddress,
			BcaVa:            address.BcaVa,
			OtherVa:          address.OtherVa,
			Note:             address.Note,
			DistrictId:       address.DistrictId,
			Status:           int32(address.Status),
			CreatedAt:        timestamppb.New(address.CreatedAt),
			UpdatedAt:        timestamppb.New(address.UpdatedAt),
		})
	}

	res = &bridgeService.GetAddressListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *BridgeGrpcHandler) GetAddressDetail(ctx context.Context, req *bridgeService.GetAddressDetailRequest) (res *bridgeService.GetAddressDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetAddressDetail")
	defer span.End()

	var address dto.AddressResponse
	address, err = h.ServicesAddress.GetDetailGRPC(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.GetAddressDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &bridgeService.Address{
			Id:               address.ID,
			Code:             address.Code,
			CustomerName:     address.CustomerName,
			ArchetypeId:      address.ArchetypeID,
			AdmDivisionId:    address.AdmDivisionID,
			SiteId:           address.SiteID,
			SalespersonId:    address.SalespersonID,
			TerritoryId:      address.TerritoryID,
			AddressCode:      address.AddressCode,
			AddressName:      address.AddressName,
			ContactPerson:    address.ContactPerson,
			City:             address.City,
			State:            address.State,
			ZipCode:          address.ZipCode,
			CountryCode:      address.CountryCode,
			Country:          address.Country,
			Latitude:         &address.Latitude,
			Longitude:        &address.Longitude,
			UpsZone:          address.UpsZone,
			ShippingMethod:   address.ShippingMethod,
			TaxScheduleId:    address.TaxScheduleID,
			PrintPhoneNumber: int32(address.PrintPhoneNumber),
			Phone_1:          address.Phone1,
			Phone_2:          address.Phone2,
			Phone_3:          address.Phone3,
			FaxNumber:        address.FaxNumber,
			ShippingAddress:  address.ShippingAddress,
			BcaVa:            address.BcaVa,
			OtherVa:          address.OtherVa,
			Note:             address.Note,
			DistrictId:       address.DistrictId,
			Status:           int32(address.Status),
			CreatedAt:        timestamppb.New(address.CreatedAt),
			UpdatedAt:        timestamppb.New(address.UpdatedAt),
		},
	}
	return
}

func (h *BridgeGrpcHandler) GetAddressListWithExcludedIds(ctx context.Context, req *bridgeService.GetAddressListWithExcludedIdsRequest) (res *bridgeService.GetAddressListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetAddressListWithExcludedIds")
	defer span.End()

	var addresses []dto.AddressResponse
	addresses, _, err = h.ServicesAddress.GetWithExcludedIds(ctx, int(req.Offset), int(req.Limit), int(req.Status), req.Search, req.OrderBy, req.ArchetypeId, req.AdmDivisionId, req.SiteId, req.SalespersonId, req.TerritoryId, req.TaxScheduleId, req.ExcludedIds)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*bridgeService.Address
	for _, address := range addresses {
		data = append(data, &bridgeService.Address{
			Id:               address.ID,
			Code:             address.Code,
			CustomerName:     address.CustomerName,
			ArchetypeId:      address.ArchetypeID,
			AdmDivisionId:    address.AdmDivisionID,
			SiteId:           address.SiteID,
			SalespersonId:    address.SalespersonID,
			TerritoryId:      address.TerritoryID,
			AddressCode:      address.AddressCode,
			AddressName:      address.AddressName,
			ContactPerson:    address.ContactPerson,
			City:             address.City,
			State:            address.State,
			ZipCode:          address.ZipCode,
			CountryCode:      address.CountryCode,
			Country:          address.Country,
			Latitude:         &address.Latitude,
			Longitude:        &address.Longitude,
			UpsZone:          address.UpsZone,
			ShippingMethod:   address.ShippingMethod,
			TaxScheduleId:    address.TaxScheduleID,
			PrintPhoneNumber: int32(address.PrintPhoneNumber),
			Phone_1:          address.Phone1,
			Phone_2:          address.Phone2,
			Phone_3:          address.Phone3,
			FaxNumber:        address.FaxNumber,
			ShippingAddress:  address.ShippingAddress,
			BcaVa:            address.BcaVa,
			OtherVa:          address.OtherVa,
			Note:             address.Note,
			DistrictId:       address.DistrictId,
			Status:           int32(address.Status),
			CreatedAt:        timestamppb.New(address.CreatedAt),
			UpdatedAt:        timestamppb.New(address.UpdatedAt),
		})
	}

	res = &bridgeService.GetAddressListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *BridgeGrpcHandler) DeleteAddress(ctx context.Context, req *bridgeService.DeleteAddressRequest) (res *bridgeService.DeleteAddressResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.DeleteAddress")
	defer span.End()

	param := &dto.DeleteAddressRequest{
		AdrsCode: req.Adrscode,
		InterID:  global.EnvDatabaseGP,
	}

	_, err = h.ServicesAddress.Delete(ctx, param)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.DeleteAddressResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
	}
	return
}

func (h *BridgeGrpcHandler) GetAddressGPList(ctx context.Context, req *bridgeService.GetAddressGPListRequest) (res *bridgeService.GetAddressGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetAddressGPList")
	defer span.End()

	res, err = h.ServicesAddress.GetGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	return
}

func (h *BridgeGrpcHandler) GetAddressGPDetail(ctx context.Context, req *bridgeService.GetAddressGPDetailRequest) (res *bridgeService.GetAddressGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetAddressGPDetail")
	defer span.End()

	res, err = h.ServicesAddress.GetDetailGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	return
}

func (h *BridgeGrpcHandler) CreateAddress(ctx context.Context, req *bridgeService.CreateAddressRequest) (res *bridgeService.CreateAddressResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.CreateAddress")
	defer span.End()

	param := &dto.AddressRequestCreate{
		InterID:               global.EnvDatabaseGP,
		CustNmbr:              req.Custnmbr,
		AdrsCode:              req.Adrscode,
		CntcPrsn:              req.Cntcprsn,
		Address1:              req.AddresS1,
		Address2:              req.AddresS2,
		Address3:              req.AddresS3,
		City:                  req.City,
		State:                 req.State,
		Zip:                   req.Zip,
		CCode:                 req.CCode,
		Country:               req.Country,
		GnlAdministrativeCode: req.GnL_Administrative_Code,
		GnlArchetypeID:        req.GnL_Archetype_ID,
		Upszone:               req.Upzone,
		Shipmthd:              req.Shipmthd,
		Taxschid:              req.Taxschid,
		Locncode:              req.Locncode,
		Slprsnid:              req.Slprsnid,
		Salsterr:              req.Salsterr,
		GnlLongitude:          req.GnL_Longitude,
		GnlLatitude:           req.GnL_Latitude,
		Userdef1:              req.UserdeF1,
		Userdef2:              req.UserdeF2,
		ShipToName:            req.ShipToName,
		Phone1:                req.PhonE1,
		Phone2:                req.PhonE2,
		Phone3:                req.PhonE3,
		Fax:                   req.Fax,
		GnlAddressNote:        req.GnL_Address_Note,
		Param:                 req.Param,
	}

	_, err = h.ServicesAddress.Create(ctx, param)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.CreateAddressResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
	}

	return
}

func (h *BridgeGrpcHandler) UpdateAddress(ctx context.Context, req *bridgeService.UpdateAddressRequest) (res *bridgeService.UpdateAddressResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.UpdateAddress")
	defer span.End()

	// param := &dto.UpdateAddressRequest{
	// 	InterID:               global.EnvDatabaseGP,
	// 	CustNmbr:              req.Custnmbr,
	// 	AdrsCode:              req.Adrscode,
	// 	CntcPrsn:              req.Cntcprsn,
	// 	Address1:              req.AddresS1,
	// 	Address2:              req.AddresS2,
	// 	Address3:              req.AddresS3,
	// 	City:                  req.City,
	// 	State:                 req.State,
	// 	Zip:                   req.Zip,
	// 	CCode:                 req.CCode,
	// 	Country:               req.Country,
	// 	GnlAdministrativeCode: req.GnL_Administrative_Code,
	// 	GnlArchetypeID:        req.GnL_Archetype_ID,
	// 	Upszone:               req.Upzone,
	// 	Shipmthd:              req.Shipmthd,
	// 	Taxschid:              req.Taxschid,
	// 	Locncode:              req.Locncode,
	// 	Slprsnid:              req.Slprsnid,
	// 	Salsterr:              req.Salsterr,
	// 	GnlLongitude:          req.GnL_Longitude,
	// 	GnlLatitude:           req.GnL_Latitude,
	// 	Userdef1:              req.UserdeF1,
	// 	Userdef2:              req.UserdeF2,
	// 	ShipToName:            req.ShipToName,
	// 	Phone1:                req.PhonE1,
	// 	Phone2:                req.PhonE2,
	// 	Phone3:                req.PhonE3,
	// 	Fax:                   req.Fax,
	// 	GnlAddressNote:        req.GnL_Address_Note,
	// 	Param:                 req.Param,
	// 	Inactive:              req.Inactive,
	// }

	_, err = h.ServicesAddress.Update(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.UpdateAddressResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
	}

	return
}

func (h *BridgeGrpcHandler) SetDefaultAddress(ctx context.Context, req *bridgeService.SetDefaultAddressRequest) (res *bridgeService.SetDefaultAddressResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.SetDefaultAddress")
	defer span.End()

	param := &dto.SetDefaultAddressRequest{
		AdrsCode: req.Adrscode,
		CustNmbr: req.Custnmbr,
		InterID:  global.EnvDatabaseGP,
	}

	_, err = h.ServicesAddress.SetDefault(ctx, param)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.SetDefaultAddressResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
	}
	return
}
