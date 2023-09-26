package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
	"git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/internal/app/dto"
)

type IConfigService interface {
	GetConfigApp(ctx context.Context, application int32, field string, attribute string, value string) (res []dto.ConfigAppResponse, total int64, err error)
	GetGlossary(ctx context.Context, table string, attribute string, valueInt int, valueName string) (res []dto.GlossaryResponse, total int64, err error)
}

type ConfigService struct {
	opt opt.Options
}

func NewConfigService() IConfigService {
	return &ConfigService{
		opt: global.Setup.Common,
	}
}

func (s *ConfigService) GetConfigApp(ctx context.Context, application int32, field string, attribute string, value string) (res []dto.ConfigAppResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ConfigService.GetConfigApp")
	defer span.End()

	appConfig, err := s.opt.Client.ConfigurationServiceGrpc.GetConfigAppList(ctx, &configuration_service.GetConfigAppListRequest{
		Application: application,
		Field:       field,
		Attribute:   attribute,
		Value:       value,
	})

	for _, applicationConfig := range appConfig.Data {
		res = append(res, dto.ConfigAppResponse{
			ID:          int64(applicationConfig.Id),
			Application: int8(applicationConfig.Application),
			Field:       applicationConfig.Field,
			Attribute:   applicationConfig.Attribute,
			Value:       applicationConfig.Value,
		})
	}
	total = int64(len(appConfig.Data))

	return
}

func (s *ConfigService) GetGlossary(ctx context.Context, table string, attribute string, valueInt int, valueName string) (res []dto.GlossaryResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ConfigService.GetGlossary")
	defer span.End()

	glossary, err := s.opt.Client.ConfigurationServiceGrpc.GetGlossaryList(ctx, &configuration_service.GetGlossaryListRequest{
		Table:     table,
		Attribute: attribute,
		ValueInt:  int32(valueInt),
		ValueName: valueName,
	})

	for _, glossary := range glossary.Data {
		res = append(res, dto.GlossaryResponse{
			ID:        int64(glossary.Id),
			Table:     glossary.Table,
			Attribute: glossary.Attribute,
			ValueInt:  int8(glossary.ValueInt),
			ValueName: glossary.ValueName,
		})
	}

	total = int64(len(glossary.Data))

	return
}
