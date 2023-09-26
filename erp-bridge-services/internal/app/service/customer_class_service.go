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

type ICustomerClassService interface {
	Get(ctx context.Context, req *pb.GetCustomerClassListRequest) (res *pb.GetCustomerClassResponse, err error)
	GetDetail(ctx context.Context, req *pb.GetCustomerClassDetailRequest) (res *pb.GetCustomerClassResponse, err error)
}

type CustomerClassService struct {
	opt opt.Options
}

func NewCustomerClassService() ICustomerClassService {
	return &CustomerClassService{
		opt: global.Setup.Common,
	}
}

func (s *CustomerClassService) Get(ctx context.Context, req *pb.GetCustomerClassListRequest) (res *pb.GetCustomerClassResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CustomerClassService.Get")
	defer span.End()

	params := map[string]string{
		"interid":    global.EnvDatabaseGP,
		"PageNumber": strconv.Itoa(int(req.Offset)),
		"PageSize":   strconv.Itoa(int(req.Limit)),
	}

	if req.Classid != "" {
		req.Classid = url.PathEscape(req.Classid)
		params["classid"] = req.Classid
	}

	if req.Clasdscr != "" {
		req.Clasdscr = url.PathEscape(req.Clasdscr)
		params["clasdscr"] = req.Clasdscr
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "customer/class/getall", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *CustomerClassService) GetDetail(ctx context.Context, req *pb.GetCustomerClassDetailRequest) (res *pb.GetCustomerClassResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CustomerClassService.GetDetail")
	defer span.End()

	req.Id = url.PathEscape(req.Id)
	params := map[string]string{
		"interid": global.EnvDatabaseGP,
		"id":      req.Id,
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "customer/class/getbyid", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if len(res.Data) == 0 {
		err = edenlabs.ErrorNotFound("customer_class")
	}

	return
}
