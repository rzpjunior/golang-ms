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

func (h *PromotionGrpcHandler) GetPriceTieringLogList(ctx context.Context, req *promotionService.GetPriceTieringLogListRequest) (res *promotionService.GetPriceTieringLogListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetPriceTieringLogList")
	defer span.End()

	param := &dto.PriceTieringLogRequestGet{
		PriceTieringIDGP: req.PriceTieringIdGp,
		CustomerID:       req.CustomerId,
		AddressIDGP:      req.AddressIdGp,
		SalesOrderIDGP:   req.SalesOrderIdGp,
		ItemID:           req.ItemId,
		Status:           int8(req.Status),
		Offset:           req.Offset,
		Limit:            req.Limit,
	}

	var (
		PriceTieringLogss []*dto.PriceTieringLogResponse
		total             int64
	)
	PriceTieringLogss, total, err = h.ServicesPriceTieringLog.Get(ctx, param)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*promotionService.PriceTieringLog
	for _, PriceTieringLog := range PriceTieringLogss {
		data = append(data, &promotionService.PriceTieringLog{
			Id:               PriceTieringLog.ID,
			PriceTieringIdGp: PriceTieringLog.PriceTieringIDGP,
			CustomerId:       PriceTieringLog.CustomerID,
			AddressIdGp:      PriceTieringLog.AddressIDGP,
			SalesOrderIdGp:   PriceTieringLog.SalesOrderIDGP,
			ItemId:           PriceTieringLog.ItemID,
			DiscountQty:      PriceTieringLog.DiscountQty,
			DiscountAmount:   PriceTieringLog.DiscountAmount,
			Status:           int32(PriceTieringLog.Status),
			CreatedAt:        timestamppb.New(PriceTieringLog.CreatedAt),
		})
	}

	res = &promotionService.GetPriceTieringLogListResponse{
		Code:     int32(codes.OK),
		Message:  codes.OK.String(),
		Data:     data,
		TotalQty: total,
	}
	return
}

func (h *PromotionGrpcHandler) CreatePriceTieringLog(ctx context.Context, req *promotionService.CreatePriceTieringLogRequest) (res *promotionService.CreatePriceTieringLogResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetVoucherMobileDetail")
	defer span.End()

	param := &dto.PriceTieringLogRequestCreate{
		PriceTieringIDGP: req.PriceTieringIdGp,
		CustomerID:       req.CustomerId,
		AddressIDGP:      req.AddressIdGp,
		SalesOrderIDGP:   req.SalesOrderIdGp,
		ItemID:           req.ItemId,
		DiscountQty:      req.DiscountQty,
		DiscountAmount:   req.DiscountAmount,
	}

	err = h.ServicesPriceTieringLog.Create(ctx, param)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &promotionService.CreatePriceTieringLogResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
	}
	return
}

func (h *PromotionGrpcHandler) CancelPriceTieringLog(ctx context.Context, req *promotionService.CancelPriceTieringLogRequest) (res *promotionService.CancelPriceTieringLogResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetVoucherMobileDetail")
	defer span.End()

	param := &dto.PriceTieringLogRequestCancel{
		PriceTieringIDGP: req.PriceTieringIdGp,
		CustomerID:       req.CustomerId,
		AddressIDGP:      req.AddressIdGp,
		SalesOrderIDGP:   req.SalesOrderIdGp,
		ItemID:           req.ItemId,
	}

	err = h.ServicesPriceTieringLog.Cancel(ctx, param)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &promotionService.CancelPriceTieringLogResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
	}
	return
}
