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

type IShippingMethodService interface {
	Get(ctx context.Context, req *pb.GetShippingMethodListRequest) (res *pb.GetShippingMethodResponse, err error)
	GetDetail(ctx context.Context, req *pb.GetShippingMethodDetailRequest) (res *pb.GetShippingMethodResponse, err error)
}

type ShippingMethodService struct {
	opt opt.Options
}

func NewShippingMethodService() IShippingMethodService {
	return &ShippingMethodService{
		opt: global.Setup.Common,
	}
}

func (s *ShippingMethodService) Get(ctx context.Context, req *pb.GetShippingMethodListRequest) (res *pb.GetShippingMethodResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ShippingMethodService.Get")
	defer span.End()

	params := map[string]string{
		"interid":    global.EnvDatabaseGP,
		"PageNumber": strconv.Itoa(int(req.Offset)),
		"PageSize":   strconv.Itoa(int(req.Limit)),
	}

	if req.Shiptype != "" {
		params["shiptype"] = req.Shiptype
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "ShippingMethod/GetAll", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *ShippingMethodService) GetDetail(ctx context.Context, req *pb.GetShippingMethodDetailRequest) (res *pb.GetShippingMethodResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ShippingMethodService.GetDetail")
	defer span.End()

	req.Id = url.PathEscape(req.Id)
	params := map[string]string{
		"interid": global.EnvDatabaseGP,
		"id":      req.Id,
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "ShippingMethod/getbyid", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if len(res.Data) == 0 {
		err = edenlabs.ErrorNotFound("shipping_method")
	}

	return
}
