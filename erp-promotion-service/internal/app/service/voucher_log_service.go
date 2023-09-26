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

type IVoucherLogService interface {
	Get(ctx context.Context, req *dto.VoucherLogRequestGet) (res []*dto.VoucherLogResponse, total int64, err error)
	Create(ctx context.Context, req *dto.VoucherLogRequestCreate) (err error)
	Cancel(ctx context.Context, req *dto.VoucherLogRequestCancel) (err error)
}

type VoucherLogService struct {
	opt                  opt.Options
	RepositoryVoucherLog repository.IVoucherLogRepository
}

func NewVoucherLogService() IVoucherLogService {
	return &VoucherLogService{
		opt:                  global.Setup.Common,
		RepositoryVoucherLog: repository.NewVoucherLogRepository(),
	}
}

func (s *VoucherLogService) Get(ctx context.Context, req *dto.VoucherLogRequestGet) (res []*dto.VoucherLogResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "VoucherLogService.Get")
	defer span.End()

	var VoucherLogs []*model.VoucherLog
	VoucherLogs, total, err = s.RepositoryVoucherLog.Get(ctx, req)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, VoucherLog := range VoucherLogs {
		res = append(res, &dto.VoucherLogResponse{
			ID:                    VoucherLog.ID,
			VoucherID:             VoucherLog.VoucherID,
			CreatedAt:             VoucherLog.CreatedAt,
			SalesOrderIDGP:        VoucherLog.SalesOrderIDGP,
			CustomerID:            VoucherLog.CustomerID,
			AddressIDGP:           VoucherLog.AddressIDGP,
			Status:                VoucherLog.Status,
			VoucherDiscountAmount: VoucherLog.VoucherDiscountAmount,
		})
	}

	return
}

func (s *VoucherLogService) Create(ctx context.Context, req *dto.VoucherLogRequestCreate) (err error) {
	ctx, span := s.opt.Trace.Start(ctx, "VoucherLogService.UpdateImage")
	defer span.End()

	voucherLog := &model.VoucherLog{
		VoucherID:             req.VoucherID,
		CustomerID:            req.CustomerID,
		AddressIDGP:           req.AddressIDGP,
		SalesOrderIDGP:        req.SalesOrderIDGP,
		VoucherDiscountAmount: req.VoucherDiscountAmount,
		Status:                statusx.ConvertStatusName("Active"),
		CreatedAt:             time.Now(),
	}

	err = s.RepositoryVoucherLog.Create(ctx, voucherLog)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *VoucherLogService) Cancel(ctx context.Context, req *dto.VoucherLogRequestCancel) (err error) {
	ctx, span := s.opt.Trace.Start(ctx, "VoucherLogService.Update")
	defer span.End()

	var listVoucherLogs []*model.VoucherLog

	param := &dto.VoucherLogRequestGet{
		VoucherID:      req.VoucherID,
		CustomerID:     req.CustomerID,
		AddressIDGP:    req.AddressIDGP,
		SalesOrderIDGP: req.SalesOrderIDGP,
		Status:         statusx.ConvertStatusName(statusx.Active),
		OrderBy:        "-id",
		Offset:         0,
		Limit:          1,
		Code:           req.Code,
	}

	listVoucherLogs, _, err = s.RepositoryVoucherLog.Get(ctx, param)

	for _, v := range listVoucherLogs {
		voucherLog := &model.VoucherLog{
			ID:     v.ID,
			Status: statusx.ConvertStatusName(statusx.Cancelled),
		}
		err = s.RepositoryVoucherLog.Update(ctx, voucherLog, "Status")
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	}

	return
}
