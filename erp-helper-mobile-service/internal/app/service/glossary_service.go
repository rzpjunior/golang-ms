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

type IGlossaryService interface {
	GetGlossary(ctx context.Context, req dto.GetGlossaryRequest) (res []dto.GlossaryResponse, total int64, err error)
}

type GlossaryService struct {
	opt opt.Options
}

func NewServiceGlossary() IGlossaryService {
	return &GlossaryService{
		opt: global.Setup.Common,
	}
}

func (s *GlossaryService) GetGlossary(ctx context.Context, req dto.GetGlossaryRequest) (res []dto.GlossaryResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "GlossaryService.GetGlossary")
	defer span.End()

	var glossary *configurationService.GetGlossaryListResponse

	if glossary, err = s.opt.Client.ConfigurationServiceGrpc.GetGlossaryList(ctx, &configurationService.GetGlossaryListRequest{
		Table:     req.Table,
		Attribute: req.Attribute,
		ValueInt:  int32(req.ValueInt),
		ValueName: req.ValueName,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("configuration", "glossary")
		return
	}

	for _, glossary := range glossary.Data {
		res = append(res, dto.GlossaryResponse{
			ID:        int64(glossary.Id),
			Table:     glossary.Table,
			Attribute: glossary.Attribute,
			ValueInt:  int8(glossary.ValueInt),
			ValueName: glossary.ValueName,
			Note:      glossary.Note,
		})
	}

	total = int64(len(glossary.Data))

	return
}
