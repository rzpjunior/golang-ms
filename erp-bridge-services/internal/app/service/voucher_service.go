package service

import (
	"context"
	"net/url"
	"strconv"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

type IVoucherService interface {
	Get(ctx context.Context, req *pb.GetVoucherGPListRequest) (res *pb.GetVoucherGPResponse, err error)
	Create(ctx context.Context, req *dto.CreateVoucherGPRequest) (res *pb.CreateVoucherGPResponse, err error)
}

type VoucherService struct {
	opt opt.Options
}

func NewVoucherService() IVoucherService {
	return &VoucherService{
		opt: global.Setup.Common,
	}
}

func (s *VoucherService) Get(ctx context.Context, req *pb.GetVoucherGPListRequest) (res *pb.GetVoucherGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "VoucherService.Get")
	defer span.End()

	params := map[string]string{
		"interid":    global.EnvDatabaseGP,
		"PageNumber": strconv.Itoa(int(req.Offset)),
		"PageSize":   strconv.Itoa(int(req.Limit)),
	}

	if req.GnlVoucherId != "" {
		req.GnlVoucherId = url.PathEscape(req.GnlVoucherId)
		params["gnl_voucher_id"] = req.GnlVoucherId
	}

	if req.GnlVoucherCode != "" {
		req.GnlVoucherCode = url.PathEscape(req.GnlVoucherCode)
		params["gnl_voucher_code"] = req.GnlVoucherCode
	}

	if req.GnlVoucherStatus != "" {
		req.GnlVoucherStatus = url.PathEscape(req.GnlVoucherStatus)
		params["gnl_voucher_status"] = req.GnlVoucherStatus
	}

	if req.GnlStartPeriod != "" {
		req.GnlStartPeriod = url.PathEscape(req.GnlStartPeriod)
		params["gnl_start_period"] = req.GnlStartPeriod
	}

	if req.GnlEndPeriod != "" {
		req.GnlEndPeriod = url.PathEscape(req.GnlEndPeriod)
		params["gnl_end_period"] = req.GnlEndPeriod
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "Voucher/list", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *VoucherService) Create(ctx context.Context, req *dto.CreateVoucherGPRequest) (res *pb.CreateVoucherGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "VoucherService.GetDetail")
	defer span.End()

	req.InterID = global.EnvDatabaseGP

	err = global.HttpRestApiToMicrosoftGP("POST", "Voucher/create", req, &res, nil)
	if err != nil {

		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}
