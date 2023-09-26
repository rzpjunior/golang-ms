package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-helper-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-helper-mobile-service/internal/app/dto"
	configurationService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
)

type IConfigAppService interface {
	GetConfigApp(ctx context.Context, req dto.GetConfigAppRequest) (res []dto.ConfigAppResponse, total int64, err error)
}

type ConfigAppService struct {
	opt opt.Options
}

func NewServiceConfigApp() IConfigAppService {
	return &ConfigAppService{
		opt: global.Setup.Common,
	}
}

func (s *ConfigAppService) GetConfigApp(ctx context.Context, req dto.GetConfigAppRequest) (res []dto.ConfigAppResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ConfigAppService.GetConfigApp")
	defer span.End()

	var configApp *configurationService.GetConfigAppListResponse

	if configApp, err = s.opt.Client.ConfigurationServiceGrpc.GetConfigAppList(ctx, &configurationService.GetConfigAppListRequest{
		Id:          int32(req.Id),
		Application: int32(req.Application),
		Field:       req.Field,
		Attribute:   req.Attribute,
		Value:       req.Value,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("configuration", "config app")
		return
	}

	for _, v := range configApp.Data {
		res = append(res, dto.ConfigAppResponse{
			Id:          int64(v.Id),
			Application: int8(v.Application),
			Field:       v.Field,
			Attribute:   v.Attribute,
			Value:       v.Value,
		})
	}

	total = int64(len(configApp.Data))

	return
}
