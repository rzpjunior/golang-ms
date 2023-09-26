package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-configuration-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/internal/app/repository"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
)

type IApplicationConfigService interface {
	Get(ctx context.Context, req *dto.ApplicationConfigRequestGet) (res []dto.ApplicationConfigResponse, total int64, err error)
	GetByID(ctx context.Context, id int64) (res dto.ApplicationConfigResponse, err error)
	GetDetail(ctx context.Context, req *pb.GetConfigAppDetailRequest) (res dto.ApplicationConfigResponse, err error)
	Update(ctx context.Context, req dto.ApplicationConfigRequestUpdate, id int64) (res dto.ApplicationConfigResponse, err error)
}

type ApplicationConfigService struct {
	opt                         opt.Options
	RepositoryApplicationConfig repository.IApplicationConfigRepository
}

func NewApplicationConfigService() IApplicationConfigService {
	return &ApplicationConfigService{
		opt:                         global.Setup.Common,
		RepositoryApplicationConfig: repository.NewApplicationConfigRepository(),
	}
}

func (s *ApplicationConfigService) Get(ctx context.Context, req *dto.ApplicationConfigRequestGet) (res []dto.ApplicationConfigResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ApplicationConfigService.Get")
	defer span.End()

	var applicationConfigs []*model.ApplicationConfig
	applicationConfigs, total, err = s.RepositoryApplicationConfig.Get(ctx, req)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, applicationConfig := range applicationConfigs {
		res = append(res, dto.ApplicationConfigResponse{
			ID:          applicationConfig.ID,
			Application: applicationConfig.Application,
			Field:       applicationConfig.Field,
			Attribute:   applicationConfig.Attribute,
			Value:       applicationConfig.Value,
		})
	}

	return
}

func (s *ApplicationConfigService) GetByID(ctx context.Context, id int64) (res dto.ApplicationConfigResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ApplicationConfigService.GetByID")
	defer span.End()

	var applicationConfig *model.ApplicationConfig
	applicationConfig, err = s.RepositoryApplicationConfig.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.ApplicationConfigResponse{
		ID:          applicationConfig.ID,
		Application: applicationConfig.Application,
		Field:       applicationConfig.Field,
		Attribute:   applicationConfig.Attribute,
		Value:       applicationConfig.Value,
	}

	return
}

func (s *ApplicationConfigService) GetDetail(ctx context.Context, req *pb.GetConfigAppDetailRequest) (res dto.ApplicationConfigResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ApplicationConfigService.GetByID")
	defer span.End()

	var applicationConfig *model.ApplicationConfig
	applicationConfig, err = s.RepositoryApplicationConfig.GetDetail(ctx, req)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.ApplicationConfigResponse{
		ID:          applicationConfig.ID,
		Application: applicationConfig.Application,
		Field:       applicationConfig.Field,
		Attribute:   applicationConfig.Attribute,
		Value:       applicationConfig.Value,
	}

	return
}

func (s *ApplicationConfigService) Update(ctx context.Context, req dto.ApplicationConfigRequestUpdate, id int64) (res dto.ApplicationConfigResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ApplicationConfigService.Update")
	defer span.End()

	applicationConfig := &model.ApplicationConfig{
		ID:          id,
		Application: req.Application,
		Field:       req.Field,
		Value:       req.Value,
	}

	// validate data is exist
	_, err = s.RepositoryApplicationConfig.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	err = s.RepositoryApplicationConfig.Update(ctx, applicationConfig, "Application", "Field", "Value")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.ApplicationConfigResponse{
		ID:          applicationConfig.ID,
		Application: req.Application,
		Field:       req.Field,
		Value:       req.Value,
	}

	return
}
