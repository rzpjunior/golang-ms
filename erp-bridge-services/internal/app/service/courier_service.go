package service

import (
	"context"
	"errors"
	"net/url"
	"strconv"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/repository"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

type ICourierService interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, vehicleProfileID int64, emergencyMode int64) (res []dto.CourierResponse, total int64, err error)
	GetDetail(ctx context.Context, id int64, code string, userID int64) (res dto.CourierResponse, err error)
	GetGP(ctx context.Context, req *pb.GetCourierGPListRequest) (res *pb.GetCourierGPResponse, err error)
	GetDetailGP(ctx context.Context, req *pb.GetCourierGPDetailRequest) (res *pb.GetCourierGPResponse, err error)
	ActivateEmergency(ctx context.Context, id string) (res *pb.EmergencyCourierResponse, err error)
	DeactivateEmergency(ctx context.Context, id string) (res *pb.EmergencyCourierResponse, err error)
}

type CourierService struct {
	opt               opt.Options
	RepositoryCourier repository.ICourierRepository
}

func NewCourierService() ICourierService {
	return &CourierService{
		opt:               global.Setup.Common,
		RepositoryCourier: repository.NewCourierRepository(),
	}
}

func (s *CourierService) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, vehicleProfileID int64, emergencyMode int64) (res []dto.CourierResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CourierService.Get")
	defer span.End()

	var couriers []*model.Courier
	couriers, total, err = s.RepositoryCourier.Get(ctx, offset, limit, status, search, orderBy, vehicleProfileID, emergencyMode)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, courier := range couriers {
		res = append(res, dto.CourierResponse{
			ID:                courier.ID,
			RoleID:            courier.RoleID,
			UserID:            courier.UserID,
			Code:              courier.Code,
			Name:              courier.Name,
			PhoneNumber:       courier.PhoneNumber,
			VehicleProfileID:  courier.VehicleProfileID,
			LicensePlate:      courier.LicensePlate,
			EmergencyMode:     courier.EmergencyMode,
			LastEmergencyTime: courier.LastEmergencyTime,
			Status:            courier.Status,
			StatusConvert:     statusx.ConvertStatusValue(courier.Status),
		})
	}

	return
}

func (s *CourierService) GetDetail(ctx context.Context, id int64, code string, userID int64) (res dto.CourierResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CourierService.GetDetail")
	defer span.End()

	var courier *model.Courier
	courier, err = s.RepositoryCourier.GetDetail(ctx, id, code, userID)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.CourierResponse{
		ID:                courier.ID,
		RoleID:            courier.RoleID,
		UserID:            courier.UserID,
		Code:              courier.Code,
		Name:              courier.Name,
		PhoneNumber:       courier.PhoneNumber,
		VehicleProfileID:  courier.VehicleProfileID,
		LicensePlate:      courier.LicensePlate,
		EmergencyMode:     courier.EmergencyMode,
		LastEmergencyTime: courier.LastEmergencyTime,
		Status:            courier.Status,
		StatusConvert:     statusx.ConvertStatusValue(courier.Status),
	}

	return
}

func (s *CourierService) GetGP(ctx context.Context, req *pb.GetCourierGPListRequest) (res *pb.GetCourierGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CourierService.GetGP")
	defer span.End()

	params := map[string]string{
		"interid":    global.EnvDatabaseGP,
		"PageNumber": strconv.Itoa(int(req.Offset)),
		"PageSize":   strconv.Itoa(int(req.Limit)),
	}

	if req.GnlCourierVendorId != "" {
		params["gnl_courier_vendor_id"] = req.GnlCourierVendorId
	}

	if req.Orderby != "" {
		params["orderby"] = req.Orderby
	}

	if req.GnlVehicleProfileId != "" {
		params["gnl_vehicle_profile_id"] = req.GnlVehicleProfileId
	}

	if req.GnlCourierId != "" {
		params["gnl_courier_id"] = req.GnlCourierId
	}

	if req.GnlCourierName != "" {
		params["gnl_courier_name"] = url.PathEscape(req.GnlCourierName)
	}

	if req.Phonname != "" {
		params["phonname"] = req.Phonname
	}

	if req.Inactive == 0 || req.Inactive == 1 {
		params["inactive"] = strconv.Itoa(int(req.Inactive))
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "courier/getall", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *CourierService) GetDetailGP(ctx context.Context, req *pb.GetCourierGPDetailRequest) (res *pb.GetCourierGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CourierService.GetDetailGP")
	defer span.End()

	params := map[string]string{
		"interid": global.EnvDatabaseGP,
		"id":      req.Id,
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "courier/getbyid", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *CourierService) ActivateEmergency(ctx context.Context, id string) (res *pb.EmergencyCourierResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CourierService.ActivateEmergency")
	defer span.End()

	var request *dto.EmergencyModeRequest
	request = &dto.EmergencyModeRequest{
		InterID:          global.EnvDatabaseGP,
		GnlCourierID:     id,
		GnlEmergencyMode: 1,
	}

	err = global.HttpRestApiToMicrosoftGP("PUT", "Courier/update", request, &res, nil)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = errors.New("Fail update the emergency Mode")
		return
	}

	if res.Code != 200 {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = errors.New("Fail update the emergency Mode")
		return
	}

	return
}

func (s *CourierService) DeactivateEmergency(ctx context.Context, id string) (res *pb.EmergencyCourierResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CourierService.DeactivateEmergency")
	defer span.End()

	var request *dto.EmergencyModeRequest
	request = &dto.EmergencyModeRequest{
		InterID:          global.EnvDatabaseGP,
		GnlCourierID:     id,
		GnlEmergencyMode: 0,
	}

	err = global.HttpRestApiToMicrosoftGP("PUT", "Courier/update", request, &res, nil)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = errors.New("Fail update the emergency Mode")
		return
	}

	if res.Code != 200 {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = errors.New("Fail update the emergency Mode")
		return
	}

	return
}
