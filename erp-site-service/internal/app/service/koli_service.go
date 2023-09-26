package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-site-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/internal/app/repository"
)

type IKoliService interface {
	Get(ctx context.Context, req *dto.KoliGetRequest) (res []*dto.KoliResponse, total int64, err error)
	GetDetail(ctx context.Context, id int64) (res *dto.KoliResponse, err error)
}

type KoliService struct {
	opt            opt.Options
	RepositoryKoli repository.IKoliRepository
}

func NewKoliService() IKoliService {
	return &KoliService{
		opt:            global.Setup.Common,
		RepositoryKoli: repository.NewKoliRepository(),
	}
}

func (s *KoliService) Get(ctx context.Context, req *dto.KoliGetRequest) (res []*dto.KoliResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "KoliService.Get")
	defer span.End()

	var kolis []*model.Koli
	kolis, total, err = s.RepositoryKoli.Get(ctx, req)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, koli := range kolis {
		res = append(res, &dto.KoliResponse{
			Id:     koli.Id,
			Code:   koli.Code,
			Value:  koli.Value,
			Name:   koli.Name,
			Note:   koli.Note,
			Status: koli.Status,
		})
	}

	return
}

func (s *KoliService) GetDetail(ctx context.Context, id int64) (res *dto.KoliResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "KoliService.GetDetail")
	defer span.End()

	var koli *model.Koli
	koli, err = s.RepositoryKoli.GetByID(ctx, id)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &dto.KoliResponse{
		Id:     koli.Id,
		Code:   koli.Code,
		Value:  koli.Value,
		Name:   koli.Name,
		Note:   koli.Note,
		Status: koli.Status,
	}

	return
}
