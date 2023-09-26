package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/edenlabs/edenlabs/timex"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/repository"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"github.com/sirupsen/logrus"
)

type IAddressService interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, archetypeID int64, admDivisionID int64, siteID int64, salespersonID int64, territoryID int64, taxScheduleID int64) (res []dto.AddressResponse, total int64, err error)
	GetWithExcludedIds(ctx context.Context, offset int, limit int, status int, search string, orderBy string, archetypeID int64, admDivisionID int64, siteID int64, salespersonID int64, territoryID int64, taxScheduleID int64, excludedIds []int64) (res []dto.AddressResponse, total int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (res dto.AddressResponse, err error)
	GetListGRPC(ctx context.Context, req *bridgeService.GetAddressListRequest) (res []dto.AddressResponse, total int64, err error)
	GetDetailGRPC(ctx context.Context, req *bridgeService.GetAddressDetailRequest) (res dto.AddressResponse, err error)
	Delete(ctx context.Context, req *dto.DeleteAddressRequest) (eres dto.CommonGPResponse, err error)
	GetGP(ctx context.Context, req *pb.GetAddressGPListRequest) (res *pb.GetAddressGPResponse, err error)
	GetDetailGP(ctx context.Context, req *pb.GetAddressGPDetailRequest) (res *pb.GetAddressGPResponse, err error)
	Create(ctx context.Context, req *dto.AddressRequestCreate) (res dto.CommonGPResponse, err error)
	Update(ctx context.Context, req *bridgeService.UpdateAddressRequest) (res dto.CommonGPResponse, err error)
	SetDefault(ctx context.Context, req *dto.SetDefaultAddressRequest) (res dto.CommonGPResponse, err error)
}

type AddressService struct {
	opt               opt.Options
	RepositoryAddress repository.IAddressRepository
}

func NewAddressService() IAddressService {
	return &AddressService{
		opt:               global.Setup.Common,
		RepositoryAddress: repository.NewAddressRepository(),
	}
}

func (s *AddressService) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, archetypeID int64, admDivisionID int64, siteID int64, salespersonID int64, territoryID int64, taxScheduleID int64) (res []dto.AddressResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AddressService.Get")
	defer span.End()

	var addresses []*model.Address
	addresses, total, err = s.RepositoryAddress.Get(ctx, offset, limit, status, search, orderBy, archetypeID, admDivisionID, siteID, salespersonID, territoryID, taxScheduleID)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, address := range addresses {
		res = append(res, dto.AddressResponse{
			ID:               address.ID,
			Code:             address.Code,
			CustomerName:     address.CustomerName,
			ArchetypeID:      address.ArchetypeID,
			AdmDivisionID:    address.AdmDivisionID,
			SiteID:           address.SiteID,
			SalespersonID:    address.SalespersonID,
			TerritoryID:      address.TerritoryID,
			AddressCode:      address.AddressCode,
			AddressName:      address.AddressName,
			ContactPerson:    address.ContactPerson,
			City:             address.City,
			State:            address.State,
			ZipCode:          address.ZipCode,
			CountryCode:      address.CountryCode,
			Country:          address.Country,
			Latitude:         address.Latitude,
			Longitude:        address.Longitude,
			UpsZone:          address.UpsZone,
			ShippingMethod:   address.ShippingMethod,
			TaxScheduleID:    address.TaxScheduleID,
			PrintPhoneNumber: address.PrintPhoneNumber,
			Phone1:           address.Phone1,
			Phone2:           address.Phone2,
			Phone3:           address.Phone3,
			FaxNumber:        address.FaxNumber,
			ShippingAddress:  address.ShippingAddress,
			BcaVa:            address.BcaVa,
			OtherVa:          address.OtherVa,
			Note:             address.Note,
			DistrictId:       address.DistrictId,
			Status:           address.Status,
			StatusConvert:    statusx.ConvertStatusValue(address.Status),
			CreatedAt:        timex.ToLocTime(ctx, address.CreatedAt),
			UpdatedAt:        timex.ToLocTime(ctx, address.UpdatedAt),
		})
	}

	return
}

func (s *AddressService) GetDetail(ctx context.Context, id int64, code string) (res dto.AddressResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AddressService.GetDetail")
	defer span.End()

	var address *model.Address
	address, err = s.RepositoryAddress.GetDetail(ctx, id, code)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.AddressResponse{
		ID:               address.ID,
		Code:             address.Code,
		CustomerName:     address.CustomerName,
		ArchetypeID:      address.ArchetypeID,
		AdmDivisionID:    address.AdmDivisionID,
		SiteID:           address.SiteID,
		SalespersonID:    address.SalespersonID,
		TerritoryID:      address.TerritoryID,
		AddressCode:      address.AddressCode,
		AddressName:      address.AddressName,
		ContactPerson:    address.ContactPerson,
		City:             address.City,
		State:            address.State,
		ZipCode:          address.ZipCode,
		CountryCode:      address.CountryCode,
		Country:          address.Country,
		Latitude:         address.Latitude,
		Longitude:        address.Longitude,
		UpsZone:          address.UpsZone,
		ShippingMethod:   address.ShippingMethod,
		TaxScheduleID:    address.TaxScheduleID,
		PrintPhoneNumber: address.PrintPhoneNumber,
		Phone1:           address.Phone1,
		Phone2:           address.Phone2,
		Phone3:           address.Phone3,
		FaxNumber:        address.FaxNumber,
		ShippingAddress:  address.ShippingAddress,
		BcaVa:            address.BcaVa,
		OtherVa:          address.OtherVa,
		Note:             address.Note,
		DistrictId:       address.DistrictId,
		Status:           address.Status,
		StatusConvert:    statusx.ConvertStatusValue(address.Status),
		CreatedAt:        timex.ToLocTime(ctx, address.CreatedAt),
		UpdatedAt:        timex.ToLocTime(ctx, address.UpdatedAt),
	}

	return
}

func (s *AddressService) GetWithExcludedIds(ctx context.Context, offset int, limit int, status int, search string, orderBy string, archetypeID int64, admDivisionID int64, siteID int64, salespersonID int64, territoryID int64, taxScheduleID int64, excludedIds []int64) (res []dto.AddressResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AddressService.Get")
	defer span.End()

	var addresses []*model.Address
	addresses, total, err = s.RepositoryAddress.GetWithExcludedIds(ctx, offset, limit, status, search, orderBy, archetypeID, admDivisionID, siteID, salespersonID, territoryID, taxScheduleID, excludedIds)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, address := range addresses {
		res = append(res, dto.AddressResponse{
			ID:               address.ID,
			Code:             address.Code,
			CustomerName:     address.CustomerName,
			ArchetypeID:      address.ArchetypeID,
			AdmDivisionID:    address.AdmDivisionID,
			SiteID:           address.SiteID,
			SalespersonID:    address.SalespersonID,
			TerritoryID:      address.TerritoryID,
			AddressCode:      address.AddressCode,
			AddressName:      address.AddressName,
			ContactPerson:    address.ContactPerson,
			City:             address.City,
			State:            address.State,
			ZipCode:          address.ZipCode,
			CountryCode:      address.CountryCode,
			Country:          address.Country,
			Latitude:         address.Latitude,
			Longitude:        address.Longitude,
			UpsZone:          address.UpsZone,
			ShippingMethod:   address.ShippingMethod,
			TaxScheduleID:    address.TaxScheduleID,
			PrintPhoneNumber: address.PrintPhoneNumber,
			Phone1:           address.Phone1,
			Phone2:           address.Phone2,
			Phone3:           address.Phone3,
			FaxNumber:        address.FaxNumber,
			ShippingAddress:  address.ShippingAddress,
			BcaVa:            address.BcaVa,
			OtherVa:          address.OtherVa,
			Note:             address.Note,
			Status:           address.Status,
			DistrictId:       address.DistrictId,
			StatusConvert:    statusx.ConvertStatusValue(address.Status),
			CreatedAt:        timex.ToLocTime(ctx, address.CreatedAt),
			UpdatedAt:        timex.ToLocTime(ctx, address.UpdatedAt),
		})
	}

	return
}

func (s *AddressService) GetListGRPC(ctx context.Context, req *bridgeService.GetAddressListRequest) (res []dto.AddressResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AddressService.Get")
	defer span.End()

	var addresses []*model.Address
	addresses, total, err = s.RepositoryAddress.GetListGRPC(ctx, req)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, address := range addresses {
		res = append(res, dto.AddressResponse{
			ID:               address.ID,
			Code:             address.Code,
			CustomerName:     address.CustomerName,
			ArchetypeID:      address.ArchetypeID,
			AdmDivisionID:    address.AdmDivisionID,
			SiteID:           address.SiteID,
			SalespersonID:    address.SalespersonID,
			TerritoryID:      address.TerritoryID,
			AddressCode:      address.AddressCode,
			AddressName:      address.AddressName,
			ContactPerson:    address.ContactPerson,
			City:             address.City,
			State:            address.State,
			ZipCode:          address.ZipCode,
			CountryCode:      address.CountryCode,
			Country:          address.Country,
			Latitude:         address.Latitude,
			Longitude:        address.Longitude,
			UpsZone:          address.UpsZone,
			ShippingMethod:   address.ShippingMethod,
			TaxScheduleID:    address.TaxScheduleID,
			PrintPhoneNumber: address.PrintPhoneNumber,
			Phone1:           address.Phone1,
			Phone2:           address.Phone2,
			Phone3:           address.Phone3,
			FaxNumber:        address.FaxNumber,
			ShippingAddress:  address.ShippingAddress,
			BcaVa:            address.BcaVa,
			OtherVa:          address.OtherVa,
			Note:             address.Note,
			DistrictId:       address.DistrictId,
			Status:           address.Status,
			StatusConvert:    statusx.ConvertStatusValue(address.Status),
			CreatedAt:        timex.ToLocTime(ctx, address.CreatedAt),
			UpdatedAt:        timex.ToLocTime(ctx, address.UpdatedAt),
		})
	}

	return
}

func (s *AddressService) GetDetailGRPC(ctx context.Context, req *bridgeService.GetAddressDetailRequest) (res dto.AddressResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AddressService.GetDetail")
	defer span.End()

	var address *model.Address
	address, err = s.RepositoryAddress.GetDetailGRPC(ctx, req)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.AddressResponse{
		ID:               address.ID,
		Code:             address.Code,
		CustomerName:     address.CustomerName,
		ArchetypeID:      address.ArchetypeID,
		AdmDivisionID:    address.AdmDivisionID,
		SiteID:           address.SiteID,
		SalespersonID:    address.SalespersonID,
		TerritoryID:      address.TerritoryID,
		AddressCode:      address.AddressCode,
		AddressName:      address.AddressName,
		ContactPerson:    address.ContactPerson,
		City:             address.City,
		State:            address.State,
		ZipCode:          address.ZipCode,
		CountryCode:      address.CountryCode,
		Country:          address.Country,
		Latitude:         address.Latitude,
		Longitude:        address.Longitude,
		UpsZone:          address.UpsZone,
		ShippingMethod:   address.ShippingMethod,
		TaxScheduleID:    address.TaxScheduleID,
		PrintPhoneNumber: address.PrintPhoneNumber,
		Phone1:           address.Phone1,
		Phone2:           address.Phone2,
		Phone3:           address.Phone3,
		FaxNumber:        address.FaxNumber,
		ShippingAddress:  address.ShippingAddress,
		BcaVa:            address.BcaVa,
		OtherVa:          address.OtherVa,
		Note:             address.Note,
		DistrictId:       address.DistrictId,
		Status:           address.Status,
		StatusConvert:    statusx.ConvertStatusValue(address.Status),
		CreatedAt:        timex.ToLocTime(ctx, address.CreatedAt),
		UpdatedAt:        timex.ToLocTime(ctx, address.UpdatedAt),
	}

	return
}

func (s *AddressService) DeleteAddress(ctx context.Context, req *bridgeService.DeleteAddressRequest) (res dto.AddressResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AddressService.GetDetail")
	defer span.End()

	// var address *model.Address
	//address, err = s.RepositoryAddress.GetDetailGRPC(ctx, req)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	reqRest := &dto.AddressResponse{
		Code: req.Adrscode,
	}
	payload, _ := json.Marshal(reqRest)
	fmt.Println("payload address : ", payload)
	// resp := &bridgeService.GetCustomerDetailResponse{}
	// err = global.HttpRestApiToMicrosoftGP("PUT", "Delete/address", reqRest, &resp)
	// if err != nil {
	// 	logrus.Error(err.Error())
	// 	os.Exit(1)
	// 	return
	// }

	res = dto.AddressResponse{
		Code: req.Adrscode,
	}

	return
}

func (s *AddressService) GetGP(ctx context.Context, req *pb.GetAddressGPListRequest) (res *pb.GetAddressGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AddressService.GetGP")
	defer span.End()

	params := map[string]string{
		"interid":    global.EnvDatabaseGP,
		"PageNumber": strconv.Itoa(int(req.Offset)),
		"PageSize":   strconv.Itoa(int(req.Limit)),
	}

	if req.CustomerNumber != "" {
		req.CustomerNumber = url.PathEscape(req.CustomerNumber)
		params["custnmbr"] = req.CustomerNumber
	}

	if req.CustomerName != "" {
		req.CustomerName = url.PathEscape(req.CustomerName)
		params["custname"] = req.CustomerName
	}

	if req.Adrscode != "" {
		params["adrscode"] = req.Adrscode
	}

	if req.ExcludeType != "" {
		params["typeEx"] = req.ExcludeType
	}

	if req.Status != "" {
		params["Inactive"] = req.Status
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "address/getall", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	return
}

func (s *AddressService) GetDetailGP(ctx context.Context, req *pb.GetAddressGPDetailRequest) (res *pb.GetAddressGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AddressService.GetDetailGP")
	defer span.End()

	req.Id = url.PathEscape(req.Id)
	params := map[string]string{
		"interid": global.EnvDatabaseGP,
		"id":      req.Id,
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "address/getbyid", nil, &res, params)

	if err != nil || len(res.Data) == 0 {
		err = edenlabs.ErrorNotFound("address")
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	return
}

//push bridge
func (s *AddressService) Create(ctx context.Context, req *dto.AddressRequestCreate) (res dto.CommonGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AddressService.Create")
	defer span.End()

	err = global.HttpRestApiToMicrosoftGP("POST", "customer/address/create", req, &res, nil)
	if err != nil {
		logrus.Error(err.Error())
		err = errors.New("Connection to the server could not be established")
		return
	}

	if res.Code != 200 {
		logrus.Error("Error Login: " + res.Message)
		err = errors.New("Connection to the server could not be established")
		return
	}

	return
}

func (s *AddressService) Update(ctx context.Context, req *bridgeService.UpdateAddressRequest) (res dto.CommonGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AddressService.Update")
	defer span.End()

	req.Interid = global.EnvDatabaseGP

	err = global.HttpRestApiToMicrosoftGP("PUT", "customer/address/update", req, &res, nil)
	if err != nil {
		logrus.Error(err.Error())
		err = errors.New("Connection to the server could not be established")
		return
	}

	if res.Code != 200 {
		logrus.Error("Error Login: " + res.Message)
		err = errors.New("Connection to the server could not be established")
		return
	}

	return
}

func (s *AddressService) SetDefault(ctx context.Context, req *dto.SetDefaultAddressRequest) (res dto.CommonGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AddressService.SetDefault")
	defer span.End()

	err = global.HttpRestApiToMicrosoftGP("PUT", "customer/setdefaultaddress", req, &res, nil)
	if err != nil {
		logrus.Error(err.Error())
		err = errors.New("Connection to the server could not be established")
		return
	}

	if res.Code != 200 {
		logrus.Error("Error Login: " + res.Message)
		err = errors.New("Connection to the server could not be established")
		return
	}

	return
}

func (s *AddressService) Delete(ctx context.Context, req *dto.DeleteAddressRequest) (res dto.CommonGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AddressService.Delete")
	defer span.End()

	// params := map[string]string{
	// 	"interid": global.EnvDatabaseGP,
	// 	"id":      req.Id,
	// }

	// err = global.HttpRestApiToMicrosoftGP("GET", "address/getbyid", nil, &res, params)

	// if err != nil {
	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 	return
	// }

	return
}
