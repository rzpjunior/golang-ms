package service

import (
	"context"
	"strings"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-configuration-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

type IAdmDivisionService interface {
	Get(ctx context.Context, req *dto.AdmDivisionGetRequest) (res []*dto.AdmDivisionResponse, total int64, err error)
}

type AdmDivisionService struct {
	opt opt.Options
}

func NewAdmDivisionService() IAdmDivisionService {
	return &AdmDivisionService{
		opt: global.Setup.Common,
	}
}

func (s *AdmDivisionService) Get(ctx context.Context, req *dto.AdmDivisionGetRequest) (res []*dto.AdmDivisionResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AdmDivisionService.Get")
	defer span.End()

	var typeAdmDivision string

	req.Region = strings.Title(req.Region)
	req.Province = strings.Title(req.Province)
	req.City = strings.Title(req.City)
	req.District = strings.Title(req.District)
	req.SubDistrict = strings.Title(req.SubDistrict)

	// Handling if all param is empty
	if req.Region == "" && req.Province == "" && req.City == "" && req.District == "" && req.SubDistrict == "" && req.RegionSearch == "" && req.ProvinceSearch == "" && req.CitySearch == "" && req.DistrictSearch == "" && req.SubDistrictSearch == "" {
		typeAdmDivision = "region"
	}

	admDivisions, err := s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPDetail(ctx, &bridge_service.GetAdmDivisionGPDetailRequest{
		Limit:           float64(req.Limit),
		Offset:          float64(req.Offset),
		Type:            typeAdmDivision,
		Region:          req.Region,
		RegionLike:      req.RegionSearch,
		State:           req.Province,
		StateLike:       req.ProvinceSearch,
		City:            req.City,
		CityLike:        req.CitySearch,
		District:        req.District,
		DistrictLike:    req.DistrictSearch,
		Subdistrict:     req.SubDistrict,
		SubdistrictLike: req.SubDistrictSearch,
	})

	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, admDivision := range admDivisions.Data {
		res = append(res, &dto.AdmDivisionResponse{
			Region:      admDivision.Region,
			Province:    admDivision.State,
			City:        admDivision.City,
			District:    admDivision.District,
			SubDistrict: admDivision.Subdistrict,
			PostalCode:  admDivision.Zipcode,
		})
	}

	total = int64(admDivisions.TotalRecords)

	return
}
