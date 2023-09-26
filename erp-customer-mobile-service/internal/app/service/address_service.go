package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/utils"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
)

type IAddressService interface {
	Get(ctx context.Context, req dto.RequestGetAddressList) (res []dto.ListAddressList, err error)
	Create(ctx context.Context, req *dto.CreateAddressRequest) (err error)
	Update(ctx context.Context, req *dto.UpdateAddressRequest) (err error)
	SetDefault(ctx context.Context, req *dto.SetDefaultAddressRequest) (err error)
	Delete(ctx context.Context, req *dto.DeleteAddressRequest) (err error)
}

type AddressService struct {
	opt opt.Options
	//RepositoryOTPOutgoing repository.IOtpOutgoingRepository
}

func NewAddressService() IAddressService {
	return &AddressService{
		opt: global.Setup.Common,
		//RepositoryOTPOutgoing: repository.NewOtpOutgoingRepository(),
	}
}

func (s *AddressService) Get(ctx context.Context, req dto.RequestGetAddressList) (res []dto.ListAddressList, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AddressService.Get")
	defer span.End()

	// Get Address List by Customer code with status active
	addressList, err := s.opt.Client.BridgeServiceGrpc.GetAddressGPList(ctx, &bridge_service.GetAddressGPListRequest{
		CustomerNumber: req.Session.Customer.Code,
		Status:         "0",
		Limit:          100,
		Offset:         1,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "Address")
		return
	}

	for _, address := range addressList.Data {
		admDivision, err := s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPList(ctx, &bridge_service.GetAdmDivisionGPListRequest{
			AdmDivisionCode: address.AdministrativeDiv.GnlAdministrativeCode,
			Limit:           10,
			Offset:          1,
		})
		if err != nil || len(admDivision.Data) == 0 {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return res, err
		}

		res = append(res, dto.ListAddressList{
			ArchetypeID:   address.GnL_Archetype_ID,
			AddressID:     address.Adrscode,
			AddressName:   address.ShipToName,
			PicName:       address.Cntcprsn,
			PhoneNumber:   address.PhonE1,
			Address1:      address.AddresS1,
			Address2:      address.AddresS2,
			Address3:      address.AddresS3,
			AddressType:   address.TypeAddress,
			AddressNote:   address.GnL_Address_Note,
			Latitude:      utils.ToFloat(address.GnL_Latitude),
			Longitude:     utils.ToFloat(address.GnL_Longitude),
			AdmDivisionId: address.AdministrativeDiv.GnlAdministrativeCode,
			RegionID:      admDivision.Data[0].Region,
			Province:      admDivision.Data[0].State,
			City:          admDivision.Data[0].City,
			District:      admDivision.Data[0].District,
			SubDistrict:   admDivision.Data[0].Subdistrict,
		})
	}

	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *AddressService) Create(ctx context.Context, req *dto.CreateAddressRequest) (err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AddressService.Create")
	defer span.End()

	var admDivisionDetail *bridge_service.AdmDivisionGP
	admDivision, err := s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPList(ctx, &bridge_service.GetAdmDivisionGPListRequest{
		SubDistrict: req.Data.SubDistrict,
		Limit:       10,
		Offset:      1,
	})
	if err != nil || len(admDivision.Data) == 0 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("bridge", "sub_district")
		return err
	}

	// Get Address list for get existing archetype
	address, err := s.opt.Client.BridgeServiceGrpc.GetAddressGPList(ctx, &bridge_service.GetAddressGPListRequest{
		CustomerNumber: req.Session.Customer.Code,
		Status:         "0",
		Limit:          1,
		Offset:         0,
	})
	if err != nil || len(address.Data) == 0 {
		err = edenlabs.ErrorValidation("address", "address id tidak valid")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// Validate for max length address
	if len(req.Data.Address1) > 60 || len(req.Data.Address2) > 60 || len(req.Data.Address3) > 60 {
		err = edenlabs.ErrorValidation("address", "Jumlah karakter address harus sama atau kurang dari 60 karakter")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return err
	}

	codeGenerator, err := s.opt.Client.ConfigurationServiceGrpc.GetGenerateCode(ctx, &configuration_service.GetGenerateCodeRequest{
		Format: req.Session.Customer.Code + "-",
		Domain: "address",
		Length: 3,
	})

	admDivisionDetail = admDivision.Data[0]
	_, err = s.opt.Client.BridgeServiceGrpc.CreateAddress(ctx, &bridge_service.CreateAddressRequest{
		Custnmbr:                req.Session.Customer.Code,
		Custname:                req.Session.Customer.Name,
		Cntcprsn:                req.Data.PICName,
		ShipToName:              req.Data.AddressName,
		Adrscode:                codeGenerator.Data.Code,
		AddresS1:                req.Data.Address1,
		AddresS2:                req.Data.Address2,
		AddresS3:                req.Data.Address3,
		GnL_Address_Note:        req.Data.AddressNote,
		Country:                 "Indonesia",
		City:                    admDivisionDetail.City,
		State:                   admDivisionDetail.State,
		PhonE1:                  req.Data.PhoneNumber,
		Inactive:                "0",
		GnL_Administrative_Code: admDivision.Data[0].Code,
		Locncode:                "WAREHOUSE",
		TypeAddress:             "other",
		GnL_Latitude:            req.Data.Latitude,
		GnL_Longitude:           req.Data.Longitude,
		GnL_Archetype_ID:        address.Data[0].GnL_Archetype_ID,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return err
	}

	return

}

func (s *AddressService) Update(ctx context.Context, req *dto.UpdateAddressRequest) (err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AddressService.Update")
	defer span.End()

	var admDivisionDetail *bridge_service.AdmDivisionGP
	admDivision, err := s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPList(ctx, &bridge_service.GetAdmDivisionGPListRequest{
		SubDistrict: req.Data.SubDistrict,
		Limit:       10,
		Offset:      1,
	})
	if err != nil || len(admDivision.Data) == 0 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("bridge", "sub_district")
		return err
	}

	// check Address
	address, err := s.opt.Client.BridgeServiceGrpc.GetAddressGPDetail(ctx, &bridge_service.GetAddressGPDetailRequest{
		Id: req.Data.AddressID,
	})
	if err != nil {
		err = edenlabs.ErrorValidation("address", "address id tidak valid")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// Validate the address is match with customer session
	if address.Data[0].Custnmbr != req.Session.Customer.Code {
		err = edenlabs.ErrorValidation("address", "address id tidak valid")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// Validate for max length address
	if len(req.Data.Address1) > 60 || len(req.Data.Address2) > 60 || len(req.Data.Address3) > 60 {
		err = edenlabs.ErrorValidation("address", "Jumlah karakter address harus sama atau kurang dari 60 karakter")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return err
	}

	admDivisionDetail = admDivision.Data[0]

	_, err = s.opt.Client.BridgeServiceGrpc.UpdateAddress(ctx, &bridge_service.UpdateAddressRequest{
		Custnmbr:                req.Session.Customer.Code,
		Custname:                req.Session.Customer.Name,
		Adrscode:                req.Data.AddressID,
		Cntcprsn:                req.Data.PICName,
		ShipToName:              req.Data.AddressName,
		AddresS1:                req.Data.Address1,
		AddresS2:                req.Data.Address2,
		AddresS3:                req.Data.Address3,
		GnL_Address_Note:        req.Data.AddressNote,
		Country:                 "Indonesia",
		City:                    admDivisionDetail.City,
		State:                   admDivisionDetail.State,
		PhonE1:                  req.Data.PhoneNumber,
		Inactive:                0,
		GnL_Administrative_Code: admDivision.Data[0].Code,
		Locncode:                "WAREHOUSE",
		TypeAddress:             address.Data[0].TypeAddress,
		GnL_Latitude:            req.Data.Latitude,
		GnL_Longitude:           req.Data.Longitude,
		GnL_Archetype_ID:        address.Data[0].GnL_Archetype_ID,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return err
	}

	return

}

func (s *AddressService) SetDefault(ctx context.Context, req *dto.SetDefaultAddressRequest) (err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AddressService.SetDefault")
	defer span.End()

	// check Address
	address, err := s.opt.Client.BridgeServiceGrpc.GetAddressGPDetail(ctx, &bridge_service.GetAddressGPDetailRequest{
		Id: req.Data.AddressID,
	})
	if err != nil {
		err = edenlabs.ErrorValidation("address_id", "address id tidak valid")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// Validate the address is match with customer session
	if address.Data[0].Custnmbr != req.Session.Customer.Code {
		err = edenlabs.ErrorValidation("address", "address id tidak valid")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	_, err = s.opt.Client.BridgeServiceGrpc.SetDefaultAddress(ctx, &bridge_service.SetDefaultAddressRequest{
		Adrscode: req.Data.AddressID,
		Custnmbr: req.Session.Customer.Code,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return err
	}
	return

}

func (s *AddressService) Delete(ctx context.Context, req *dto.DeleteAddressRequest) (err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AddressService.Delete")
	defer span.End()

	//check Address
	address, err := s.opt.Client.BridgeServiceGrpc.GetAddressGPDetail(ctx, &bridge_service.GetAddressGPDetailRequest{
		Id: req.Data.AddressID,
	})
	if err != nil {
		err = edenlabs.ErrorValidation("address_id", "address id tidak valid")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// Validate the address is match with customer session
	if address.Data[0].Custnmbr != req.Session.Customer.Code {
		err = edenlabs.ErrorValidation("address", "address id tidak valid")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if address.Data[0].TypeAddress == "ship_to" {
		err = edenlabs.ErrorValidation("default_address", "Default Address Tidak Bisa dihapus, Silahkan ubah default address terlebih dahulu")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return err
	}

	_, err = s.opt.Client.BridgeServiceGrpc.UpdateAddress(ctx, &bridge_service.UpdateAddressRequest{
		Custnmbr:                req.Session.Customer.Code,
		Custname:                req.Session.Customer.Name,
		Adrscode:                address.Data[0].Adrscode,
		Cntcprsn:                address.Data[0].Cntcprsn,
		ShipToName:              address.Data[0].ShipToName,
		AddresS1:                address.Data[0].AddresS1,
		AddresS2:                address.Data[0].AddresS2,
		AddresS3:                address.Data[0].AddresS3,
		GnL_Address_Note:        address.Data[0].GnL_Address_Note,
		Country:                 address.Data[0].Country,
		City:                    address.Data[0].City,
		State:                   address.Data[0].State,
		PhonE1:                  address.Data[0].PhonE1,
		Inactive:                1,
		GnL_Administrative_Code: address.Data[0].AdministrativeDiv.GnlAdministrativeCode,
		Locncode:                address.Data[0].Locncode,
		TypeAddress:             address.Data[0].TypeAddress,
		GnL_Latitude:            address.Data[0].GnL_Latitude,
		GnL_Longitude:           address.Data[0].GnL_Longitude,
		GnL_Archetype_ID:        address.Data[0].GnL_Archetype_ID,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return err
	}

	return

}
