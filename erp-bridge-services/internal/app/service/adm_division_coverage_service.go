package service

import (
	"context"
	"net/url"
	"strconv"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

type IAdmDivisionCoverageService interface {
	GetGP(ctx context.Context, req *pb.GetAdmDivisionCoverageGPListRequest) (res *pb.GetAdmDivisionCoverageGPResponse, err error)
	GetDetailGP(ctx context.Context, req *pb.GetAdmDivisionCoverageGPDetailRequest) (res *pb.GetAdmDivisionCoverageGPResponse, err error)
}

type AdmDivisionCoverageService struct {
	opt opt.Options
}

func NewAdmDivisionCoverageService() IAdmDivisionCoverageService {
	return &AdmDivisionCoverageService{
		opt: global.Setup.Common,
	}
}

func (s *AdmDivisionCoverageService) GetGP(ctx context.Context, req *pb.GetAdmDivisionCoverageGPListRequest) (res *pb.GetAdmDivisionCoverageGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AdmDivisionCoverageService.GetGP")
	defer span.End()

	params := map[string]string{
		"interid":    global.EnvDatabaseGP,
		"PageNumber": strconv.Itoa(int(req.Offset)),
		"PageSize":   strconv.Itoa(int(req.Limit)),
	}

	if req.GnlAdministrativeCode != "" {
		req.GnlAdministrativeCode = url.PathEscape(req.GnlAdministrativeCode)
		params["gnl_administrative_code"] = req.GnlAdministrativeCode
	}

	if req.GnlProvince != "" {
		req.GnlProvince = url.PathEscape(req.GnlProvince)
		params["gnl_province"] = req.GnlProvince
	}

	if req.GnlCity != "" {
		req.GnlCity = url.PathEscape(req.GnlCity)
		params["gnl_city"] = req.GnlCity
	}

	if req.GnlDistrict != "" {
		req.GnlDistrict = url.PathEscape(req.GnlDistrict)
		params["gnl_district"] = req.GnlDistrict
	}

	if req.GnlSubdistrict != "" {
		req.GnlSubdistrict = url.PathEscape(req.GnlSubdistrict)
		params["gnl_subdistrict"] = req.GnlSubdistrict
	}

	if req.Locncode != "" {
		req.Locncode = url.PathEscape(req.Locncode)
		params["locncode"] = req.Locncode
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "AdmDivisionCoverage/getall", nil, &res, params)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *AdmDivisionCoverageService) GetDetailGP(ctx context.Context, req *pb.GetAdmDivisionCoverageGPDetailRequest) (res *pb.GetAdmDivisionCoverageGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AdmDivisionCoverageService.GetDetailGP")
	defer span.End()

	req.Id = url.PathEscape(req.Id)
	params := map[string]string{
		"interid": global.EnvDatabaseGP,
		"id":      req.Id,
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "AdmDivisionCoverage/getbyid", nil, &res, params)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}
