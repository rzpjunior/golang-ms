package service

import (
	"context"
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

type IVehicleProfileService interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, courierVendorID int64) (res []dto.VehicleProfileResponse, total int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (res dto.VehicleProfileResponse, err error)
	GetGP(ctx context.Context, req *pb.GetVehicleProfileGPListRequest) (res *pb.GetVehicleProfileGPResponse, err error)
	GetDetailGP(ctx context.Context, req *pb.GetVehicleProfileGPDetailRequest) (res *pb.GetVehicleProfileGPResponse, err error)
}

type VehicleProfileService struct {
	opt                      opt.Options
	RepositoryVehicleProfile repository.IVehicleProfileRepository
}

func NewVehicleProfileService() IVehicleProfileService {
	return &VehicleProfileService{
		opt:                      global.Setup.Common,
		RepositoryVehicleProfile: repository.NewVehicleProfileRepository(),
	}
}

func (s *VehicleProfileService) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, courierVendorID int64) (res []dto.VehicleProfileResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "VehicleProfileService.Get")
	defer span.End()

	var vehicleProfiles []*model.VehicleProfile
	vehicleProfiles, total, err = s.RepositoryVehicleProfile.Get(ctx, offset, limit, status, search, orderBy, courierVendorID)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, vProfile := range vehicleProfiles {
		res = append(res, dto.VehicleProfileResponse{
			ID:                  vProfile.ID,
			Code:                vProfile.Code,
			Name:                vProfile.Name,
			MaxKoli:             vProfile.MaxKoli,
			MaxWeight:           vProfile.MaxWeight,
			MaxFragile:          vProfile.MaxFragile,
			SpeedFactor:         vProfile.SpeedFactor,
			RoutingProfile:      vProfile.RoutingProfile,
			Status:              vProfile.Status,
			StatusConvert:       statusx.ConvertStatusValue(vProfile.Status),
			Skills:              vProfile.Skills,
			InitialCost:         vProfile.InitialCost,
			SubsequentCost:      vProfile.SubsequentCost,
			MaxAvailableVehicle: vProfile.MaxAvailableVehicle,
			CourierVendorID:     vProfile.CourierVendorID,
		})
	}

	return
}

func (s *VehicleProfileService) GetDetail(ctx context.Context, id int64, code string) (res dto.VehicleProfileResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "VehicleProfileService.GetDetail")
	defer span.End()

	var vehicleProfile *model.VehicleProfile
	vehicleProfile, err = s.RepositoryVehicleProfile.GetDetail(ctx, id, code)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.VehicleProfileResponse{
		ID:                  vehicleProfile.ID,
		Code:                vehicleProfile.Code,
		Name:                vehicleProfile.Name,
		MaxKoli:             vehicleProfile.MaxKoli,
		MaxWeight:           vehicleProfile.MaxWeight,
		MaxFragile:          vehicleProfile.MaxFragile,
		SpeedFactor:         vehicleProfile.SpeedFactor,
		RoutingProfile:      vehicleProfile.RoutingProfile,
		Status:              vehicleProfile.Status,
		StatusConvert:       statusx.ConvertStatusValue(vehicleProfile.Status),
		Skills:              vehicleProfile.Skills,
		InitialCost:         vehicleProfile.InitialCost,
		SubsequentCost:      vehicleProfile.SubsequentCost,
		MaxAvailableVehicle: vehicleProfile.MaxAvailableVehicle,
		CourierVendorID:     vehicleProfile.CourierVendorID,
	}

	return
}

func (s *VehicleProfileService) GetGP(ctx context.Context, req *pb.GetVehicleProfileGPListRequest) (res *pb.GetVehicleProfileGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "VehicleProfileService.GetGP")
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

	if req.Locncode != "" {
		params["locncode"] = req.Locncode
	}

	if req.Inactive != 0{
		params["inactive"] = strconv.Itoa(int(req.Inactive))
	}

	if req.GnlDescription100 != ""{
		params["gnl_description100"] = req.GnlDescription100
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "VehicleProfile/getall", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *VehicleProfileService) GetDetailGP(ctx context.Context, req *pb.GetVehicleProfileGPDetailRequest) (res *pb.GetVehicleProfileGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "VehicleProfileService.GetDetailGP")
	defer span.End()

	params := map[string]string{
		"interid": global.EnvDatabaseGP,
		"id":      req.Id,
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "VehicleProfile/getbyid", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}
