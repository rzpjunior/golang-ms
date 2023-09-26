package service

import (
	"context"
	"net/url"
	"strconv"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

type ISalesPriceLevelService interface {
	Get(ctx context.Context, req *pb.GetSalesPriceLevelListRequest) (res *pb.GetSalesPriceLevelResponse, err error)
	GetDetail(ctx context.Context, req *pb.GetSalesPriceLevelDetailRequest) (res *pb.GetSalesPriceLevelResponse, err error)
}

type SalesPriceLevelService struct {
	opt opt.Options
}

func NewSalesPriceLevelService() ISalesPriceLevelService {
	return &SalesPriceLevelService{
		opt: global.Setup.Common,
	}
}

func (s *SalesPriceLevelService) Get(ctx context.Context, req *pb.GetSalesPriceLevelListRequest) (res *pb.GetSalesPriceLevelResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesPriceLevelService.Get")
	defer span.End()

	params := map[string]string{
		"interid":    global.EnvDatabaseGP,
		"PageNumber": strconv.Itoa(int(req.Offset)),
		"PageSize":   strconv.Itoa(int(req.Limit)),
	}

	if req.GnlRegion != "" {
		req.GnlRegion = url.PathEscape(req.GnlRegion)
		params["gnl_region"] = req.GnlRegion
	}

	if req.GnlCustTypeId != "" {
		req.GnlCustTypeId = url.PathEscape(req.GnlCustTypeId)
		params["gnl_cust_type_id"] = req.GnlCustTypeId
	}

	if req.Prclevel != "" {
		req.Prclevel = url.PathEscape(req.Prclevel)
		params["prclevel"] = req.Prclevel
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "SalesPriceLevel/getall", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *SalesPriceLevelService) GetDetail(ctx context.Context, req *pb.GetSalesPriceLevelDetailRequest) (res *pb.GetSalesPriceLevelResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesPriceLevelService.GetDetail")
	defer span.End()

	req.Id = url.PathEscape(req.Id)
	params := map[string]string{
		"interid": global.EnvDatabaseGP,
		"id":      req.Id,
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "SalesPriceLevel/getbyid", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if len(res.Data) == 0 {
		err = edenlabs.ErrorNotFound("sales_price_level")
	}

	return
}
