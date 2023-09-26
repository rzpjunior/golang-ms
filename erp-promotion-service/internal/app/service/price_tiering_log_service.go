package service

import (
	"context"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/project-version3/erp-services/erp-promotion-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-promotion-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-promotion-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-promotion-service/internal/app/repository"
)

type IPriceTieringLogService interface {
	Get(ctx context.Context, req *dto.PriceTieringLogRequestGet) (res []*dto.PriceTieringLogResponse, totalQty int64, err error)
	Create(ctx context.Context, req *dto.PriceTieringLogRequestCreate) (err error)
	Cancel(ctx context.Context, req *dto.PriceTieringLogRequestCancel) (err error)
}

type PriceTieringLogService struct {
	opt                       opt.Options
	RepositoryPriceTieringLog repository.IPriceTieringLogRepository
}

func NewPriceTieringLogService() IPriceTieringLogService {
	return &PriceTieringLogService{
		opt:                       global.Setup.Common,
		RepositoryPriceTieringLog: repository.NewPriceTieringLogRepository(),
	}
}

func (s *PriceTieringLogService) Get(ctx context.Context, req *dto.PriceTieringLogRequestGet) (res []*dto.PriceTieringLogResponse, totalQty int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PriceTieringLogService.Get")
	defer span.End()

	var PriceTieringLogs []*model.PriceTieringLog
	PriceTieringLogs, err = s.RepositoryPriceTieringLog.Get(ctx, req)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, PriceTieringLog := range PriceTieringLogs {
		totalQty += int64(PriceTieringLog.DiscountQty)
		res = append(res, &dto.PriceTieringLogResponse{
			ID:               PriceTieringLog.ID,
			PriceTieringIDGP: PriceTieringLog.PriceTieringIDGP,
			SalesOrderIDGP:   PriceTieringLog.SalesOrderIDGP,
			CustomerID:       PriceTieringLog.CustomerID,
			AddressIDGP:      PriceTieringLog.AddressIDGP,
			ItemID:           PriceTieringLog.ItemID,
			DiscountQty:      PriceTieringLog.DiscountQty,
			DiscountAmount:   PriceTieringLog.DiscountAmount,
			CreatedAt:        PriceTieringLog.CreatedAt,
			Status:           PriceTieringLog.Status,
		})
	}

	return
}

func (s *PriceTieringLogService) Create(ctx context.Context, req *dto.PriceTieringLogRequestCreate) (err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PriceTieringLogService.UpdateImage")
	defer span.End()

	PriceTieringLog := &model.PriceTieringLog{
		PriceTieringIDGP: req.PriceTieringIDGP,
		SalesOrderIDGP:   req.SalesOrderIDGP,
		CustomerID:       req.CustomerID,
		AddressIDGP:      req.AddressIDGP,
		ItemID:           req.ItemID,
		DiscountQty:      req.DiscountQty,
		DiscountAmount:   req.DiscountAmount,
		Status:           statusx.ConvertStatusName("Active"),
		CreatedAt:        time.Now(),
	}

	err = s.RepositoryPriceTieringLog.Create(ctx, PriceTieringLog)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *PriceTieringLogService) Cancel(ctx context.Context, req *dto.PriceTieringLogRequestCancel) (err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PriceTieringLogService.Update")
	defer span.End()

	var listPriceTieringLogs []*model.PriceTieringLog

	param := &dto.PriceTieringLogRequestGet{
		PriceTieringIDGP: req.PriceTieringIDGP,
		CustomerID:       req.CustomerID,
		AddressIDGP:      req.AddressIDGP,
		SalesOrderIDGP:   req.SalesOrderIDGP,
		ItemID:           req.ItemID,
		Status:           statusx.ConvertStatusName(statusx.Active),
		OrderBy:          "-id",
		Offset:           0,
		Limit:            1,
	}

	listPriceTieringLogs, err = s.RepositoryPriceTieringLog.Get(ctx, param)

	for _, v := range listPriceTieringLogs {
		PriceTieringLog := &model.PriceTieringLog{
			ID:     v.ID,
			Status: statusx.ConvertStatusName(statusx.Cancelled),
		}
		err = s.RepositoryPriceTieringLog.Update(ctx, PriceTieringLog, "Status")
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	}

	return
}
