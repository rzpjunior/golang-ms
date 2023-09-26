package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-configuration-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/internal/app/repository"
)

type IGlossaryService interface {
	Get(ctx context.Context, offset int, limit int, table string, attribute string, valueInt int, valueName string) (res []dto.GlossaryResponse, total int64, err error)
	GetDetail(ctx context.Context, table string, attribute string, valueInt int, valueName string) (res dto.GlossaryResponse, err error)
}

type GlossaryService struct {
	opt                opt.Options
	RepositoryGlossary repository.IGlossaryRepository
}

func NewGlossaryService() IGlossaryService {
	return &GlossaryService{
		opt:                global.Setup.Common,
		RepositoryGlossary: repository.NewGlossaryRepository(),
	}
}

func (s *GlossaryService) Get(ctx context.Context, offset int, limit int, table string, attribute string, valueInt int, valueName string) (res []dto.GlossaryResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "GlossaryService.GetList")
	defer span.End()

	var glossarys []*model.Glossary
	glossarys, total, err = s.RepositoryGlossary.Get(ctx, offset, limit, table, attribute, valueInt, valueName)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, glossary := range glossarys {
		res = append(res, dto.GlossaryResponse{
			ID:        glossary.ID,
			Table:     glossary.Table,
			Attribute: glossary.Attribute,
			ValueInt:  glossary.ValueInt,
			ValueName: glossary.ValueName,
			Note:      glossary.Note,
		})
	}

	return
}

func (s *GlossaryService) GetDetail(ctx context.Context, table string, attribute string, valueInt int, valueName string) (res dto.GlossaryResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "GlossaryService.GetDetail")
	defer span.End()

	var glossary *model.Glossary
	glossary, err = s.RepositoryGlossary.GetDetail(ctx, table, attribute, valueInt, valueName)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.GlossaryResponse{
		ID:        glossary.ID,
		Table:     glossary.Table,
		Attribute: glossary.Attribute,
		ValueInt:  glossary.ValueInt,
		ValueName: glossary.ValueName,
		Note:      glossary.Note,
	}
	return
}
