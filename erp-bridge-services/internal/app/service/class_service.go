package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/edenlabs/edenlabs/timex"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/repository"
)

type IClassService interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (res []dto.ClassResponse, total int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (res dto.ClassResponse, err error)
}

type ClassService struct {
	opt             opt.Options
	RepositoryClass repository.IClassRepository
}

func NewClassService() IClassService {
	return &ClassService{
		opt:             global.Setup.Common,
		RepositoryClass: repository.NewClassRepository(),
	}
}

func (s *ClassService) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (res []dto.ClassResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ClassService.Get")
	defer span.End()

	var classes []*model.Class
	classes, total, err = s.RepositoryClass.Get(ctx, offset, limit, status, search, orderBy)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, class := range classes {
		res = append(res, dto.ClassResponse{
			ID:            class.ID,
			Code:          class.Code,
			Description:   class.Description,
			Status:        class.Status,
			StatusConvert: statusx.ConvertStatusValue(class.Status),
			CreatedAt:     timex.ToLocTime(ctx, class.CreatedAt),
			UpdatedAt:     timex.ToLocTime(ctx, class.UpdatedAt),
		})
	}

	return
}

func (s *ClassService) GetDetail(ctx context.Context, id int64, code string) (res dto.ClassResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ClassService.GetDetail")
	defer span.End()

	var class *model.Class
	class, err = s.RepositoryClass.GetDetail(ctx, id, code)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.ClassResponse{
		ID:            class.ID,
		Code:          class.Code,
		Description:   class.Description,
		Status:        class.Status,
		StatusConvert: statusx.ConvertStatusValue(class.Status),
		CreatedAt:     timex.ToLocTime(ctx, class.CreatedAt),
		UpdatedAt:     timex.ToLocTime(ctx, class.UpdatedAt),
	}

	return
}
