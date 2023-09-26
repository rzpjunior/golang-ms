package service

import (
	"context"
	"strconv"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

type IHelperService interface {
	GetGP(ctx context.Context, req *pb.GetHelperGPListRequest) (res *pb.GetHelperGPResponse, err error)
	GetDetailGP(ctx context.Context, req *pb.GetHelperGPDetailRequest) (res *pb.GetHelperGPResponse, err error)
	Login(ctx context.Context, req *pb.LoginHelperRequest) (res *pb.LoginHelperResponse, err error)
}

type HelperService struct {
	opt opt.Options
}

func NewHelperService() IHelperService {
	return &HelperService{
		opt: global.Setup.Common,
	}
}

func (s *HelperService) GetGP(ctx context.Context, req *pb.GetHelperGPListRequest) (res *pb.GetHelperGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "HelperService.GetGP")
	defer span.End()

	// mandatory params
	params := map[string]string{
		"interid":    global.EnvDatabaseGP,
		"PageNumber": strconv.Itoa(int(req.Offset)),
		"PageSize":   strconv.Itoa(int(req.Limit)),
		"inactive":   strconv.Itoa(int(req.Inactive)),
	}

	if req.GnlHelperId != "" {
		params["gnl_helper_id"] = req.GnlHelperId
	}

	if req.GnlHelperType != "" {
		params["gnl_helper_type"] = req.GnlHelperType
	}

	if req.GnlHelperName != "" {
		params["gnl_helper_name"] = req.GnlHelperName
	}

	if req.Employid != "" {
		params["employid"] = req.Employid
	}

	if req.Userid != "" {
		params["userid"] = req.Userid
	}

	if req.Locncode != "" {
		params["locncode"] = req.Locncode
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "helper/getall", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *HelperService) GetDetailGP(ctx context.Context, req *pb.GetHelperGPDetailRequest) (res *pb.GetHelperGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "HelperService.GetDetailGP")
	defer span.End()

	params := map[string]string{
		"interid": global.EnvDatabaseGP,
		"id":      req.Id,
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "helper/getbyid", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *HelperService) Login(ctx context.Context, req *pb.LoginHelperRequest) (res *pb.LoginHelperResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DeliveryOrderService.Create")
	defer span.End()

	req.Interid = global.EnvDatabaseGP
	req.Type = 2

	err = global.HttpRestApiToMicrosoftGP("POST", "User/authentication", req, &res, nil)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}
