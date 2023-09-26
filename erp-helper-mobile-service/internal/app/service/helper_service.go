package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-helper-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-helper-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

type IHelperService interface {
	Get(ctx context.Context, req *dto.GetHelperRequest) (res []dto.GetHelperResponse, err error)
}

type HelperService struct {
	opt opt.Options
}

func NewServiceHelper() IHelperService {
	return &HelperService{
		opt: global.Setup.Common,
	}
}

func (s *HelperService) Get(ctx context.Context, req *dto.GetHelperRequest) (res []dto.GetHelperResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "HelperService.Get")
	defer span.End()

	var helper *bridge_service.GetHelperGPResponse

	if helper, err = s.opt.Client.BridgeServiceGrpc.GetHelperGPList(ctx, &bridge_service.GetHelperGPListRequest{
		Limit:         int32(req.Limit),
		Offset:        int32(req.Offset),
		Locncode:      req.SiteId,
		GnlHelperName: req.Name,
		GnlHelperType: req.Type,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "helper")
		return
	}

	for _, v := range helper.Data {
		res = append(res, dto.GetHelperResponse{
			Id:   v.GnlHelperId,
			Name: v.GnlHelperName,
		})
	}

	return
}
