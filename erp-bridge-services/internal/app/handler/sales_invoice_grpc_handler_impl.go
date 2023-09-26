package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GetSalesInvoice
func (h *BridgeGrpcHandler) GetSalesInvoiceList(ctx context.Context, req *bridgeService.GetSalesInvoiceListRequest) (res *bridgeService.GetSalesInvoiceListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetSalesInvoiceList")
	defer span.End()

	var salesInvoices []dto.SalesInvoiceResponse
	salesInvoices, _, err = h.ServiceSalesInvoice.Get(ctx, int(req.Offset), int(req.Limit), int(req.Status), req.Search, req.OrderBy, req.AddressId, req.CustomerId, req.SalespersonId, req.OrderDateFrom.AsTime(), req.OrderDateTo.AsTime())
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var salesInvoicesGrpc []*bridgeService.SalesInvoice
	for _, salesInvoice := range salesInvoices {
		salesInvoicesGrpc = append(salesInvoicesGrpc, &bridgeService.SalesInvoice{
			Id:                salesInvoice.ID,
			Code:              salesInvoice.Code,
			CodeExt:           salesInvoice.CodeExt,
			Status:            int32(salesInvoice.Status),
			RecognitionDate:   timestamppb.New(salesInvoice.RecognitionDate),
			DueDate:           timestamppb.New(salesInvoice.DueDate),
			BillingAddress:    salesInvoice.BillingAddress,
			DeliveryFee:       salesInvoice.DeliveryFee,
			VouRedeemCode:     salesInvoice.VouRedeemCode,
			VouDiscAmount:     salesInvoice.VouDiscAmount,
			PointRedeemAmount: salesInvoice.PointRedeemAmount,
			Adjustment:        int32(salesInvoice.Adjustment),
			AdjAmount:         salesInvoice.AdjAmount,
			AdjNote:           salesInvoice.AdjNote,
			TotalPrice:        salesInvoice.TotalPrice,
			TotalCharge:       salesInvoice.TotalCharge,
			DeltaPrint:        salesInvoice.DeltaPrint,
			VoucherId:         salesInvoice.VoucherID,
			RemainingAmount:   salesInvoice.RemainingAmount,
			Note:              salesInvoice.Note,
		})
	}

	res = &bridgeService.GetSalesInvoiceListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    salesInvoicesGrpc,
	}
	return
}

func (h *BridgeGrpcHandler) GetSalesInvoiceDetail(ctx context.Context, req *bridgeService.GetSalesInvoiceDetailRequest) (res *bridgeService.GetSalesInvoiceDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetSalesInvoiceDetail")
	defer span.End()

	var salesInvoice dto.SalesInvoiceResponse
	salesInvoice, err = h.ServiceSalesInvoice.GetDetail(ctx, req.Id, req.Code)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.GetSalesInvoiceDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &bridgeService.SalesInvoice{
			Id:            salesInvoice.ID,
			Code:          salesInvoice.Code,
			Status:        int32(salesInvoice.Status),
			DeliveryFee:   salesInvoice.DeliveryFee,
			VouDiscAmount: salesInvoice.VouDiscAmount,
			TotalPrice:    salesInvoice.TotalPrice,
			TotalCharge:   salesInvoice.TotalCharge,
		},
	}
	return
}

func (h *BridgeGrpcHandler) GetSalesInvoiceGPList(ctx context.Context, req *bridgeService.GetSalesInvoiceGPListRequest) (res *bridgeService.GetSalesInvoiceGPListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetAddressGPList")
	defer span.End()

	res, err = h.ServiceSalesInvoice.GetGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	return
}

func (h *BridgeGrpcHandler) GetSalesInvoiceGPDetail(ctx context.Context, req *bridgeService.GetSalesInvoiceGPDetailRequest) (res *bridgeService.GetSalesInvoiceGPDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetAddressGPList")
	defer span.End()

	res, err = h.ServiceSalesInvoice.GetDetailGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	return
}

func (h *BridgeGrpcHandler) CreateSalesInvoiceGP(ctx context.Context, req *bridgeService.CreateSalesInvoiceGPRequest) (res *bridgeService.CreateSalesInvoiceGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.CreateSalesInvoiceGP")
	defer span.End()

	res, err = h.ServiceSalesInvoice.CreateGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}
