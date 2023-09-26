package service

import (
	"context"
	"strings"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/global"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	dto "git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

func NewServiceAdmDivision() IAdmDivisionService {
	m := new(AdmDivisionService)
	m.opt = global.Setup.Common
	return m
}

type IAdmDivisionService interface {
	GetAdmDivisions(ctx context.Context, req dto.AdmDivisionListRequest) (res []*dto.AdmDivisionResponse, err error)
	GetAdmDivisionDetailById(ctx context.Context, req dto.AdmDivisionDetailRequest) (res *dto.AdmDivisionResponse, err error)
	GetGP(ctx context.Context, req dto.AdmDivisionListRequest) (res []*dto.AdmDivisionGP, total int64, err error)
	GetDetaiGPlById(ctx context.Context, id, divType string, limit, offset int) (res []*dto.AdmDivisionGP, err error)
	GetCoverageGP(ctx context.Context, req dto.AdmDivisionCoverageListRequest) (res []*dto.AdmDivisionCoverageGP, total int64, err error)
	GetCoverageDetaiGPlById(ctx context.Context, id string) (res *dto.AdmDivisionCoverageGP, err error)
	GetGPAdmDiv(ctx context.Context, req *dto.GetAdmDivisionRequest) (res []*dto.AdmDivisionGPResponse, total int64, err error)
}

type AdmDivisionService struct {
	opt opt.Options
}

func (s *AdmDivisionService) GetGPAdmDiv(ctx context.Context, req *dto.GetAdmDivisionRequest) (res []*dto.AdmDivisionGPResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AdmDivisionService.GetGP")
	defer span.End()

	req.Region = strings.Title(req.Region)
	req.Province = strings.Title(req.Province)
	req.City = strings.Title(req.City)
	req.District = strings.Title(req.District)
	req.SubDistrict = strings.Title(req.SubDistrict)
	req.TypeAdm = strings.Title(req.TypeAdm)

	admDivision, err := s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPDetail(ctx, &bridge_service.GetAdmDivisionGPDetailRequest{
		Type:            req.TypeAdm,
		Limit:           100,
		RegionLike:      req.Region,
		StateLike:       req.Province,
		CityLike:        req.City,
		DistrictLike:    req.District,
		SubdistrictLike: req.SubDistrict,
		RegionOthers:    "no",
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
			ZipCode:     admDiv.Zipcode,
		})
	}

	total = int64(len(res))

	return
}

func (s *AdmDivisionService) GetAdmDivisions(ctx context.Context, req dto.AdmDivisionListRequest) (res []*dto.AdmDivisionResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AdmDivisionService.GetAdmDivisions")
	defer span.End()

	// get Adm Division from bridge
	var admRes *bridgeService.GetAdmDivisionListResponse
	admRes, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionList(ctx, &bridgeService.GetAdmDivisionListRequest{
		Limit:   req.Limit,
		Offset:  req.Offset,
		Status:  req.Status,
		Search:  req.Search,
		OrderBy: req.OrderBy,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "administrative division")
		return
	}

	datas := []*dto.AdmDivisionResponse{}
	for _, admDiv := range admRes.Data {
		datas = append(datas, &dto.AdmDivisionResponse{
			ID:            admDiv.Id,
			Code:          admDiv.Code,
			ProvinceID:    admDiv.ProvinceId,
			CityID:        admDiv.CityId,
			DistrictID:    admDiv.DistrictId,
			SubDistrictID: admDiv.SubDistrictId,
			RegionID:      admDiv.RegionId,
			PostalCode:    admDiv.PostalCode,
			Province:      admDiv.Province,
			City:          admDiv.City,
			District:      admDiv.District,
			Region:        admDiv.Region,
			Status:        int8(admDiv.Status),
			StatusConvert: statusx.ConvertStatusValue(int8(admDiv.Status)),
			CreatedAt:     admDiv.CreatedAt.AsTime(),
			UpdatedAt:     admDiv.UpdatedAt.AsTime(),
		})
	}
	res = datas

	return
}

func (s *AdmDivisionService) GetAdmDivisionDetailById(ctx context.Context, req dto.AdmDivisionDetailRequest) (res *dto.AdmDivisionResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "Adm DivisionService.GetAdmDivisionDetailById")
	defer span.End()

	// get Adm Division from bridge
	var admRes *bridgeService.GetAdmDivisionDetailResponse
	admRes, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionDetail(ctx, &bridgeService.GetAdmDivisionDetailRequest{
		Id: int64(req.Id),
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "administrative division")
		return
	}

	res = &dto.AdmDivisionResponse{
		ID:            admRes.Data.Id,
		Code:          admRes.Data.Code,
		ProvinceID:    admRes.Data.ProvinceId,
		CityID:        admRes.Data.CityId,
		DistrictID:    admRes.Data.DistrictId,
		SubDistrictID: admRes.Data.SubDistrictId,
		RegionID:      admRes.Data.RegionId,
		PostalCode:    admRes.Data.PostalCode,
		Province:      admRes.Data.Province,
		City:          admRes.Data.City,
		District:      admRes.Data.District,
		Region:        admRes.Data.Region,
		Status:        int8(admRes.Data.Status),
		StatusConvert: statusx.ConvertStatusValue(int8(admRes.Data.Status)),
		CreatedAt:     admRes.Data.CreatedAt.AsTime(),
		UpdatedAt:     admRes.Data.UpdatedAt.AsTime(),
	}

	return
}

func (s *AdmDivisionService) GetGP(ctx context.Context, req dto.AdmDivisionListRequest) (res []*dto.AdmDivisionGP, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AdmDivisionService.GetGP")
	defer span.End()

	// get adm division from bridge
	var admRes *bridgeService.GetAdmDivisionGPResponse
	admRes, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPList(ctx, &bridgeService.GetAdmDivisionGPListRequest{
		Limit:           req.Limit,
		Offset:          req.Offset,
		Region:          req.Region,
		AdmDivisionCode: req.Code,
		State:           req.State,
		City:            req.City,
		District:        req.District,
		SubDistrict:     req.Subdistrict,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "adm division")
		return
	}

	datas := []*dto.AdmDivisionGP{}
	for _, adm := range admRes.Data {
		datas = append(datas, &dto.AdmDivisionGP{
			Code:        adm.Code,
			Region:      adm.Region,
			State:       adm.State,
			City:        adm.City,
			District:    adm.District,
			SubDistrict: adm.Subdistrict,
		})
	}

	total = int64(len(datas))
	res = datas

	return
}

func (s *AdmDivisionService) GetDetaiGPlById(ctx context.Context, id, divType string, limit, offset int) (res []*dto.AdmDivisionGP, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AdmDivisionService.GetDetaiGPlById")
	defer span.End()

	// get adm division from bridge
	var adm *bridgeService.GetAdmDivisionGPResponse
	adm, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPDetail(ctx, &bridgeService.GetAdmDivisionGPDetailRequest{
		Id:     id,
		Type:   divType,
		Limit:  float64(limit),
		Offset: float64(offset),
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "adm division")
		return
	}

	for _, i := range adm.Data {
		res = append(res, &dto.AdmDivisionGP{
			Code:        i.Code,
			Region:      i.Region,
			State:       i.State,
			City:        i.City,
			District:    i.District,
			SubDistrict: i.Subdistrict,
		})
	}

	return
}

func (s *AdmDivisionService) GetCoverageGP(ctx context.Context, req dto.AdmDivisionCoverageListRequest) (res []*dto.AdmDivisionCoverageGP, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AdmDivisionService.GetCoverageGP")
	defer span.End()

	// get adm division from bridge
	var admRes *bridgeService.GetAdmDivisionCoverageGPResponse
	admRes, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionCoverageGPList(ctx, &bridgeService.GetAdmDivisionCoverageGPListRequest{
		Limit:                 req.Limit,
		Offset:                req.Offset,
		GnlAdministrativeCode: req.GnlAdministrativeCode,
		GnlProvince:           req.GnlProvince,
		GnlCity:               req.GnlCity,
		GnlDistrict:           req.GnlDistrict,
		GnlSubdistrict:        req.GnlSubdistrict,
		Locncode:              req.Locncode,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "adm division")
		return
	}

	datas := []*dto.AdmDivisionCoverageGP{}
	for _, adm := range admRes.Data {
		datas = append(datas, &dto.AdmDivisionCoverageGP{
			GnlAdministrativeCode: adm.GnlAdministrativeCode,
			GnlRegion:             adm.GnlRegion,
			GnlProvince:           adm.GnlProvince,
			GnlCity:               adm.GnlCity,
			GnlDistrict:           adm.GnlDistrict,
			GnlSubdistrict:        adm.GnlSubdistrict,
			Locncode:              adm.Locncode,
		})
	}

	total = int64(len(datas))
	res = datas

	return
}

func (s *AdmDivisionService) GetCoverageDetaiGPlById(ctx context.Context, id string) (res *dto.AdmDivisionCoverageGP, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AdmDivisionService.GetCoverageDetaiGPlById")
	defer span.End()

	// get adm division from bridge
	var adm *bridgeService.GetAdmDivisionCoverageGPResponse
	adm, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionCoverageGPDetail(ctx, &bridgeService.GetAdmDivisionCoverageGPDetailRequest{
		Id: id,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "adm division")
		return
	}

	res = &dto.AdmDivisionCoverageGP{
		GnlAdministrativeCode: adm.Data[0].GnlAdministrativeCode,
		GnlRegion:             adm.Data[0].GnlRegion,
		GnlProvince:           adm.Data[0].GnlProvince,
		GnlCity:               adm.Data[0].GnlCity,
		GnlDistrict:           adm.Data[0].GnlDistrict,
		GnlSubdistrict:        adm.Data[0].GnlSubdistrict,
		Locncode:              adm.Data[0].Locncode,
	}

	return
}
