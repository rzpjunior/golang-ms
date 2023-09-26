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

type IItemClassService interface {
	GetGP(ctx context.Context, req *pb.GetItemClassGPListRequest) (res *pb.GetItemClassGPResponse, err error)
	GetDetailGP(ctx context.Context, req *pb.GetItemClassGPDetailRequest) (res *pb.GetItemClassGPResponse, err error)
}

type ItemClassService struct {
	opt opt.Options
}

func NewItemClassService() IItemClassService {
	return &ItemClassService{
		opt: global.Setup.Common,
	}
}

func (s *ItemClassService) GetGP(ctx context.Context, req *pb.GetItemClassGPListRequest) (res *pb.GetItemClassGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemClassService.GetGP")
	defer span.End()

	params := map[string]string{
		"interid":    global.EnvDatabaseGP,
		"PageNumber": strconv.Itoa(int(req.Offset)),
		"PageSize":   strconv.Itoa(int(req.Limit)),
	}

	if req.Search != "" {
		req.Search = url.PathEscape(req.Search)
		params["itmclsdc"] = req.Search
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "ItemClass/getall", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *ItemClassService) GetDetailGP(ctx context.Context, req *pb.GetItemClassGPDetailRequest) (res *pb.GetItemClassGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemClassService.GetDetailGP")
	defer span.End()

	req.Id = url.PathEscape(req.Id)
	params := map[string]string{
		"interid": global.EnvDatabaseGP,
		"id":      req.Id,
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "ItemClass/getbyid", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}
