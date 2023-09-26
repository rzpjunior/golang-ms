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

type ISalesTerritoryService interface {
	GetGP(ctx context.Context, req *pb.GetSalesTerritoryGPListRequest) (res *pb.GetSalesTerritoryGPResponse, err error)
	GetDetailGP(ctx context.Context, req *pb.GetSalesTerritoryGPDetailRequest) (res *pb.GetSalesTerritoryGPResponse, err error)
}

type SalesTerritoryService struct {
	opt opt.Options
}

func NewSalesTerritoryService() ISalesTerritoryService {
	return &SalesTerritoryService{
		opt: global.Setup.Common,
	}
}

func (s *SalesTerritoryService) GetGP(ctx context.Context, req *pb.GetSalesTerritoryGPListRequest) (res *pb.GetSalesTerritoryGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesTerritoryService.GetGP")
	defer span.End()

	params := map[string]string{
		"interid":    global.EnvDatabaseGP,
		"PageNumber": strconv.Itoa(int(req.Offset)),
		"PageSize":   strconv.Itoa(int(req.Limit)),
	}

	if req.Search != "" {
		req.Search = url.PathEscape(req.Search)
		params["slterdsc"] = req.Search
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "SalesTerritory/getall", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *SalesTerritoryService) GetDetailGP(ctx context.Context, req *pb.GetSalesTerritoryGPDetailRequest) (res *pb.GetSalesTerritoryGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesTerritoryService.GetDetailGP")
	defer span.End()

	req.Id = url.PathEscape(req.Id)
	params := map[string]string{
		"interid": global.EnvDatabaseGP,
		"id":      req.Id,
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "SalesTerritory/getbyid", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}
