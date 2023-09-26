package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

type IAdmDivisionService interface {
	Get(ctx context.Context, offset int, limit int, search string, sub_district_id int, Type int, region_id int, province_id int, city_id int, district_id int) (res []dto.AdmDivisionResponse, total int64, err error)
	GetGP(ctx context.Context, req *dto.GetAdmDivisionRequest) (res []*dto.AdmDivisionGPResponse, total int64, err error)
	Search(ctx context.Context, req *dto.SearchAdmDivisionRequest) (res []*dto.AdmDivisionGPResponse, total int64, err error)
}

type AdmDivisionService struct {
	opt opt.Options
}

func NewAdmDivisionService() IAdmDivisionService {
	return &AdmDivisionService{
		opt: global.Setup.Common,
	}
}

func (s *AdmDivisionService) Get(ctx context.Context, offset int, limit int, search string, sub_district_id int, Type int, region_id int, province_id int, city_id int, district_id int) (res []dto.AdmDivisionResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AdmDivisionService.Get")
	defer span.End()

	admDivision, err := s.opt.Client.BridgeServiceGrpc.GetAdmDivisionList(ctx, &bridge_service.GetAdmDivisionListRequest{
		SubDistrictId: int64(sub_district_id),
		RegionId:      int64(region_id),
		Search:        search,
		ProvinceId:    int64(province_id),
		CityId:        int64(city_id),
		DistrictId:    int64(district_id),
	})

	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, admDiv := range admDivision.Data {
		region, err := s.opt.Client.BridgeServiceGrpc.GetRegionDetail(ctx, &bridge_service.GetRegionDetailRequest{
			Id: admDiv.RegionId,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return res, 0, err
		}

		res = append(res, dto.AdmDivisionResponse{
			ID:              strconv.Itoa(int(admDiv.Id)),
			Code:            admDiv.Code,
			RegionID:        strconv.Itoa(int(admDiv.RegionId)),
			RegionName:      region.Data.Description,
			Province:        "Dummy Province",
			City:            admDiv.City,
			District:        admDiv.District,
			SubDistrictID:   strconv.Itoa(int(admDiv.SubDistrictId)),
			SubDistrictName: "Dummy Sub District",
			DistrictID:      strconv.Itoa(int(0)),
			DistrictName:    "Dummy Distric",
			CityID:          strconv.Itoa(int(0)),
			CityName:        "Dummy City",
			ProvinceID:      strconv.Itoa(int(0)),
			ProvinceName:    "Dummy Province",
			ContryID:        strconv.Itoa(int(0)),
			CountryName:     "Dummy Country",
			PostalCode:      admDiv.PostalCode,
			Status:          strconv.Itoa(int(admDiv.Status)),
			StatusConvert:   "",
			CreatedAt:       admDiv.CreatedAt.AsTime(),
			UpdatedAt:       admDiv.UpdatedAt.AsTime(),
		})

	}
	total = int64(len(res))

	return
}

func (s *AdmDivisionService) GetGP(ctx context.Context, req *dto.GetAdmDivisionRequest) (res []*dto.AdmDivisionGPResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AdmDivisionService.GetGP")
	defer span.End()

	var typeAdmDivision, id string

	req.Data.Province = strings.Title(req.Data.Province)
	req.Data.City = strings.Title(req.Data.City)
	req.Data.District = strings.Title(req.Data.District)

	typeAdmDivision = "state"

	if req.Data.Province != "" {
		typeAdmDivision = "state"
		id = req.Data.Province
	}

	if req.Data.City != "" {
		typeAdmDivision = "city"
		id = req.Data.City
	}

	if req.Data.District != "" {
		typeAdmDivision = "district"
		id = req.Data.District
	}

	admDivision, err := s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPDetail(ctx, &bridge_service.GetAdmDivisionGPDetailRequest{
		Type: typeAdmDivision,
		Id:   id,
	})

	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, admDiv := range admDivision.Data {
		res = append(res, &dto.AdmDivisionGPResponse{
			Code:        admDiv.Code,
			Region:      admDiv.Region,
			Province:    admDiv.State,
			City:        admDiv.City,
			District:    admDiv.District,
			SubDistrict: admDiv.Subdistrict,
		})
	}

	total = int64(len(res))

	return
}

func (s *AdmDivisionService) Search(ctx context.Context, req *dto.SearchAdmDivisionRequest) (res []*dto.AdmDivisionGPResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AdmDivisionService.Search")
	defer span.End()

	admDivision, err := s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPList(ctx, &bridge_service.GetAdmDivisionGPListRequest{
		SubDistrict: req.Data.SubDistrict,
		Limit:       15,
		Offset:      1,
	})

	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, admDiv := range admDivision.Data {
		res = append(res, &dto.AdmDivisionGPResponse{
			Code:          admDiv.Code,
			Region:        admDiv.Region,
			Province:      admDiv.State,
			City:          admDiv.City,
			District:      admDiv.District,
			SubDistrict:   admDiv.Subdistrict,
			ConcatAddress: fmt.Sprintf("%s,%s, %s", admDiv.Subdistrict, admDiv.District, admDiv.City),
		})
	}

	total = int64(len(res))

	return
}
