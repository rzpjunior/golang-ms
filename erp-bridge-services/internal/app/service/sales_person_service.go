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

type ISalesPersonService interface {
	GetGP(ctx context.Context, req *pb.GetSalesPersonGPListRequest) (res *pb.GetSalesPersonGPResponse, err error)
	GetDetailGP(ctx context.Context, req *pb.GetSalesPersonGPDetailRequest) (res *pb.GetSalesPersonGPResponse, err error)
}

type SalesPersonService struct {
	opt opt.Options
}

func NewSalesPersonService() ISalesPersonService {
	return &SalesPersonService{
		opt: global.Setup.Common,
	}
}

func (s *SalesPersonService) GetGP(ctx context.Context, req *pb.GetSalesPersonGPListRequest) (res *pb.GetSalesPersonGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesPersonService.GetGP")
	defer span.End()

	params := map[string]string{
		"interid":    global.EnvDatabaseGP,
		"PageNumber": strconv.Itoa(int(req.Offset)),
		"PageSize":   strconv.Itoa(int(req.Limit)),
	}

	if req.SalesTerritoryId != "" {
		req.SalesTerritoryId = url.PathEscape(req.SalesTerritoryId)
		params["salsterr"] = req.SalesTerritoryId
	}

	if req.Status != "" {
		params["inactive"] = req.Status
	}

	if req.Search != "" {
		req.Search = url.PathEscape(req.Search)
		params["sprsnsln"] = req.Search
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "salesperson/getall", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *SalesPersonService) GetDetailGP(ctx context.Context, req *pb.GetSalesPersonGPDetailRequest) (res *pb.GetSalesPersonGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesPersonService.GetDetailGP")
	defer span.End()

	req.Id = url.PathEscape(req.Id)

	params := map[string]string{
		"interid": global.EnvDatabaseGP,
		"id":      req.Id,
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "salesperson/getbyid", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}
