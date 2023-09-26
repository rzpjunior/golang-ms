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

type IDeliveryKoliService interface {
	Get(ctx context.Context, req *dto.DeliveryKoliGetRequest) (res []*dto.DeliveryKoliResponse, total int64, err error)
}

type DeliveryKoliService struct {
	opt                    opt.Options
	RepositoryDeliveryKoli repository.IDeliveryKoliRepository
}

func NewDeliveryKoliService() IDeliveryKoliService {
	return &DeliveryKoliService{
		opt:                    global.Setup.Common,
		RepositoryDeliveryKoli: repository.NewDeliveryKoliRepository(),
	}
}

func (s *DeliveryKoliService) Get(ctx context.Context, req *dto.DeliveryKoliGetRequest) (res []*dto.DeliveryKoliResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DeliveryKoliService.Get")
	defer span.End()

	var deliveryKoli []*model.DeliveryKoli
	deliveryKoli, total, err = s.RepositoryDeliveryKoli.Get(ctx, req)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, v := range deliveryKoli {
		res = append(res, &dto.DeliveryKoliResponse{
			Id:        v.Id,
			SopNumber: v.SopNumber,
			KoliId:    v.KoliId,
			Quantity:  v.Quantity,
			Note:      v.Note,
		})
	}

	return
}
