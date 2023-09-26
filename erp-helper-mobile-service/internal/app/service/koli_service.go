package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-helper-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-helper-mobile-service/internal/app/dto"
	siteService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/site_service"
)

type IKoliService interface {
	GetKoli(ctx context.Context, req dto.GetKoliRequest) (res []dto.KoliResponse, total int64, err error)
	GetKoliDetail(ctx context.Context, id int64) (res dto.KoliResponse, err error)
}

type KoliService struct {
	opt opt.Options
}

func NewServiceKoli() IKoliService {
	return &KoliService{
		opt: global.Setup.Common,
	}
}

func (s *KoliService) GetKoli(ctx context.Context, req dto.GetKoliRequest) (res []dto.KoliResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "KoliService.GetKoli")
	defer span.End()

	var koli *siteService.GetKoliListResponse

	if koli, err = s.opt.Client.SiteServiceGrpc.GetKoliList(ctx, &siteService.GetKoliListRequest{
		Limit:   int32(req.Limit),
		Offset:  int32(req.Offset),
		Status:  int32(req.Status),
		OrderBy: req.OrderBy,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("site", "koli")
		return
	}

	for _, koli := range koli.Data {
		res = append(res, dto.KoliResponse{
			Id:     koli.Id,
			Code:   koli.Code,
			Value:  koli.Value,
			Name:   koli.Name,
			Note:   koli.Note,
			Status: int8(koli.Status),
		})
	}

	total = int64(len(koli.Data))

	return
}

func (s *KoliService) GetKoliDetail(ctx context.Context, id int64) (res dto.KoliResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "KoliService.GetKoli")
	defer span.End()

	var koli *siteService.GetKoliDetailResponse

	if koli, err = s.opt.Client.SiteServiceGrpc.GetKoliDetail(ctx, &siteService.GetKoliDetailRequest{
		Id: id,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("site", "koli")
		return
	}

	res = dto.KoliResponse{
		Id:     koli.Data.Id,
		Code:   koli.Data.Code,
		Value:  koli.Data.Value,
		Name:   koli.Data.Name,
		Note:   koli.Data.Note,
		Status: int8(koli.Data.Status),
	}

	return
}
