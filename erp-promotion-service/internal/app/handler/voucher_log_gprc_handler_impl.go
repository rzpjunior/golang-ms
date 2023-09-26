package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-promotion-service/internal/app/dto"
	promotionService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/promotion_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *PromotionGrpcHandler) GetVoucherLogList(ctx context.Context, req *promotionService.GetVoucherLogListRequest) (res *promotionService.GetVoucherLogListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetVoucherLogList")
	defer span.End()

	param := &dto.VoucherLogRequestGet{
		VoucherID:      req.VoucherId,
		CustomerID:     req.CustomerId,
		AddressIDGP:    req.AddressIdGp,
		SalesOrderIDGP: req.SalesOrderIdGp,
		Status:         int8(req.Status),
		Offset:         req.Offset,
		Limit:          req.Limit,
	}

	var (
		voucherLogss []*dto.VoucherLogResponse
		total        int64
	)
	voucherLogss, total, err = h.ServicesVoucherLog.Get(ctx, param)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*promotionService.VoucherLog
	for _, voucherLog := range voucherLogss {
		data = append(data, &promotionService.VoucherLog{
			Id:                    voucherLog.ID,
			VoucherId:             voucherLog.VoucherID,
			CustomerId:            voucherLog.CustomerID,
			AddressIdGp:           voucherLog.AddressIDGP,
			SalesOrderIdGp:        voucherLog.SalesOrderIDGP,
			VoucherDiscountAmount: voucherLog.VoucherDiscountAmount,
			Status:                int32(voucherLog.Status),
			CreatedAt:             timestamppb.New(voucherLog.CreatedAt),
		})
	}

	res = &promotionService.GetVoucherLogListResponse{
		Code:         int32(codes.OK),
		Message:      codes.OK.String(),
		Data:         data,
		TotalRecords: total,
	}
	return
}

func (h *PromotionGrpcHandler) CreateVoucherLog(ctx context.Context, req *promotionService.CreateVoucherLogRequest) (res *promotionService.CreateVoucherLogResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetVoucherMobileDetail")
	defer span.End()

	param := &dto.VoucherLogRequestCreate{
		VoucherID:             req.VoucherId,
		CustomerID:            req.CustomerId,
		AddressIDGP:           req.AddressIdGp,
		SalesOrderIDGP:        req.SalesOrderIdGp,
		VoucherDiscountAmount: req.VoucherDiscountAmount,
	}

	err = h.ServicesVoucherLog.Create(ctx, param)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &promotionService.CreateVoucherLogResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
	}
	return
}

func (h *PromotionGrpcHandler) CancelVoucherLog(ctx context.Context, req *promotionService.CancelVoucherLogRequest) (res *promotionService.CancelVoucherLogResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetVoucherMobileDetail")
	defer span.End()

	param := &dto.VoucherLogRequestCancel{
		VoucherID:      req.VoucherId,
		CustomerID:     req.CustomerId,
		AddressIDGP:    req.AddressIdGp,
		SalesOrderIDGP: req.SalesOrderIdGp,
		Code:           req.Code,
	}

	err = h.ServicesVoucherLog.Cancel(ctx, param)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &promotionService.CancelVoucherLogResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
	}
	return
}
