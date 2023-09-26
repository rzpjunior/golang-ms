package service

import (
	"context"
	"net/url"
	"strconv"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/edenlabs/edenlabs/timex"
	"git.edenfarm.id/edenlabs/edenlabs/utils"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/repository"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

type IAdmDivisionService interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, regionID int64, subDistrictID int64) (res []dto.AdmDivisionResponse, total int64, err error)
	GetDetail(ctx context.Context, id int64, code string, regionID int64, subDistrictID int64) (res dto.AdmDivisionResponse, err error)
	GetListGRPC(ctx context.Context, req *bridgeService.GetAdmDivisionListRequest) (res []dto.AdmDivisionResponse, total int64, err error)
	GetDetailGRPC(ctx context.Context, req *bridgeService.GetAdmDivisionDetailRequest) (res dto.AdmDivisionResponse, err error)
	GetGP(ctx context.Context, req *pb.GetAdmDivisionGPListRequest) (res *pb.GetAdmDivisionGPResponse, err error)
	GetDetailGP(ctx context.Context, req *pb.GetAdmDivisionGPDetailRequest) (res *pb.GetAdmDivisionGPResponse, err error)
}

type AdmDivisionService struct {
	opt                   opt.Options
	RepositoryAdmDivision repository.IAdmDivisionRepository
}

func NewAdmDivisionService() IAdmDivisionService {
	return &AdmDivisionService{
		opt:                   global.Setup.Common,
		RepositoryAdmDivision: repository.NewAdmDivisionRepository(),
	}
}

func (s *AdmDivisionService) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, regionID int64, subDistrictID int64) (res []dto.AdmDivisionResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AdmDivisionService.Get")
	defer span.End()

	var admDivisions []*model.AdmDivision
	admDivisions, total, err = s.RepositoryAdmDivision.Get(ctx, offset, limit, status, search, orderBy, regionID, subDistrictID)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, admDivision := range admDivisions {
		res = append(res, dto.AdmDivisionResponse{
			ID:            admDivision.ID,
			Code:          admDivision.Code,
			ProvinceID:    admDivision.ProvinceID,
			CityID:        admDivision.CityID,
			DistrictID:    admDivision.DistrictID,
			SubDistrictID: admDivision.SubDistrictID,
			RegionID:      admDivision.RegionID,
			PostalCode:    admDivision.PostalCode,
			Province:      admDivision.Province,
			City:          admDivision.City,
			District:      admDivision.District,
			Region:        admDivision.Region,
			Status:        admDivision.Status,
			StatusConvert: statusx.ConvertStatusValue(admDivision.Status),
			CreatedAt:     timex.ToLocTime(ctx, admDivision.CreatedAt),
			UpdatedAt:     timex.ToLocTime(ctx, admDivision.UpdatedAt),
		})
	}

	return
}
func (s *AdmDivisionService) GetListGRPC(ctx context.Context, req *bridgeService.GetAdmDivisionListRequest) (res []dto.AdmDivisionResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AdmDivisionService.Get")
	defer span.End()

	var admDivisions []*model.AdmDivision
	admDivisions, total, err = s.RepositoryAdmDivision.GetListGRPC(ctx, req)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, admDivision := range admDivisions {
		res = append(res, dto.AdmDivisionResponse{
			ID:            admDivision.ID,
			Code:          admDivision.Code,
			ProvinceID:    admDivision.ProvinceID,
			CityID:        admDivision.CityID,
			DistrictID:    admDivision.DistrictID,
			SubDistrictID: admDivision.SubDistrictID,
			RegionID:      admDivision.RegionID,
			PostalCode:    admDivision.PostalCode,
			Province:      admDivision.Province,
			City:          admDivision.City,
			District:      admDivision.District,
			Region:        admDivision.Region,
			Status:        admDivision.Status,
			StatusConvert: statusx.ConvertStatusValue(admDivision.Status),
			CreatedAt:     timex.ToLocTime(ctx, admDivision.CreatedAt),
			UpdatedAt:     timex.ToLocTime(ctx, admDivision.UpdatedAt),
		})
	}

	return
}
func (s *AdmDivisionService) GetDetailGRPC(ctx context.Context, req *bridgeService.GetAdmDivisionDetailRequest) (res dto.AdmDivisionResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AdmDivisionService.GetDetail")
	defer span.End()

	var admDivision *model.AdmDivision
	admDivision, err = s.RepositoryAdmDivision.GetDetailGRPC(ctx, req)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.AdmDivisionResponse{
		ID:            admDivision.ID,
		Code:          admDivision.Code,
		ProvinceID:    admDivision.ProvinceID,
		CityID:        admDivision.CityID,
		DistrictID:    admDivision.DistrictID,
		SubDistrictID: admDivision.SubDistrictID,
		RegionID:      admDivision.RegionID,
		PostalCode:    admDivision.PostalCode,
		Province:      admDivision.Province,
		City:          admDivision.City,
		District:      admDivision.District,
		Region:        admDivision.Region,
		Status:        admDivision.Status,
		StatusConvert: statusx.ConvertStatusValue(admDivision.Status),
		CreatedAt:     timex.ToLocTime(ctx, admDivision.CreatedAt),
		UpdatedAt:     timex.ToLocTime(ctx, admDivision.UpdatedAt),
	}

	return
}

func (s *AdmDivisionService) GetDetail(ctx context.Context, id int64, code string, regionID int64, subDistrictID int64) (res dto.AdmDivisionResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AdmDivisionService.GetDetail")
	defer span.End()

	var admDivision *model.AdmDivision
	admDivision, err = s.RepositoryAdmDivision.GetDetail(ctx, id, code, regionID, subDistrictID)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.AdmDivisionResponse{
		ID:            admDivision.ID,
		Code:          admDivision.Code,
		ProvinceID:    admDivision.ProvinceID,
		CityID:        admDivision.CityID,
		DistrictID:    admDivision.DistrictID,
		SubDistrictID: admDivision.SubDistrictID,
		RegionID:      admDivision.RegionID,
		PostalCode:    admDivision.PostalCode,
		Province:      admDivision.Province,
		City:          admDivision.City,
		District:      admDivision.District,
		Region:        admDivision.Region,
		Status:        admDivision.Status,
		StatusConvert: statusx.ConvertStatusValue(admDivision.Status),
		CreatedAt:     timex.ToLocTime(ctx, admDivision.CreatedAt),
		UpdatedAt:     timex.ToLocTime(ctx, admDivision.UpdatedAt),
	}

	return
}

func (s *AdmDivisionService) GetGP(ctx context.Context, req *pb.GetAdmDivisionGPListRequest) (res *pb.GetAdmDivisionGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AdmDivisionService.GetGP")
	defer span.End()

	params := map[string]string{
		"interid":    global.EnvDatabaseGP,
		"PageNumber": strconv.Itoa(int(req.Offset)),
		"PageSize":   strconv.Itoa(int(req.Limit)),
	}

	if req.AdmDivisionCode != "" {
		params["code"] = req.AdmDivisionCode
	}

	if req.Region != "" {
		req.Region = url.PathEscape(req.Region)
		params["Region"] = req.Region
	}

	if req.State != "" {
		req.State = url.PathEscape(req.State)
		params["State"] = req.State
	}

	if req.City != "" {
		req.City = url.PathEscape(req.City)
		params["City"] = req.City
	}

	if req.District != "" {
		req.District = url.PathEscape(req.District)
		params["District"] = req.District
	}

	if req.SubDistrict != "" {
		req.SubDistrict = url.PathEscape(req.SubDistrict)
		params["SubDistrict"] = req.SubDistrict
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "admdivision/getall", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *AdmDivisionService) GetDetailGP(ctx context.Context, req *pb.GetAdmDivisionGPDetailRequest) (res *pb.GetAdmDivisionGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AdmDivisionService.GetDetailGP")
	defer span.End()

	params := map[string]string{
		"interid": global.EnvDatabaseGP,
	}

	if req.Id != "" {
		req.Id = url.PathEscape(req.Id)
		params["id"] = req.Id
	}

	if req.Type != "" {
		req.Type = url.PathEscape(req.Type)
		params["type"] = req.Type
	}

	if req.Region != "" {
		req.Region = url.PathEscape(req.Region)
		params["region"] = req.Region
	}

	if req.State != "" {
		req.State = url.PathEscape(req.State)
		params["state"] = req.State
	}

	if req.City != "" {
		req.City = url.PathEscape(req.City)
		params["city"] = req.City
	}

	if req.District != "" {
		req.District = url.PathEscape(req.District)
		params["district"] = req.District
	}

	if req.Subdistrict != "" {
		req.Subdistrict = url.PathEscape(req.Subdistrict)
		params["subdistrict"] = req.Subdistrict
	}

	if req.RegionLike != "" {
		req.RegionLike = url.PathEscape(req.RegionLike)
		params["region_like"] = req.RegionLike
	}

	if req.StateLike != "" {
		req.StateLike = url.PathEscape(req.StateLike)
		params["state_like"] = req.StateLike
	}

	if req.CityLike != "" {
		req.CityLike = url.PathEscape(req.CityLike)
		params["city_like"] = req.CityLike
	}

	if req.DistrictLike != "" {
		req.DistrictLike = url.PathEscape(req.DistrictLike)
		params["district_like"] = req.DistrictLike
	}

	if req.SubdistrictLike != "" {
		req.SubdistrictLike = url.PathEscape(req.SubdistrictLike)
		params["subdistrict_like"] = req.SubdistrictLike
	}

	if req.RegionOthers != "" {
		req.RegionOthers = url.PathEscape(req.RegionOthers)
		params["region_others"] = req.RegionOthers
	}

	if req.Code != "" {
		req.Code = url.PathEscape(req.Code)
		params["code"] = req.Code
	}

	if req.Limit != 0 {
		params["PageSize"] = utils.ToString(req.Limit)
	}

	if req.Offset != 0 {
		params["PageNumber"] = utils.ToString(req.Offset)
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "admdivision/getbyid", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if len(res.Data) == 0 {
		err = edenlabs.ErrorNotFound("adm_division")
	}

	return
}
