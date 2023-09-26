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

func (h *BridgeGrpcHandler) GetSalesPaymentList(ctx context.Context, req *bridgeService.GetSalesPaymentListRequest) (res *bridgeService.GetSalesPaymentListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetSalesPaymentList")
	defer span.End()

	var salesPayments []dto.SalesPaymentResponse
	salesPayments, _, err = h.ServiceSalesPayment.Get(ctx, req.SalesInvoiceId)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var salesPaymentsGrpc []*bridgeService.SalesPayment
	for _, salesPayment := range salesPayments {
		salesPaymentsGrpc = append(salesPaymentsGrpc, &bridgeService.SalesPayment{
			Id:              salesPayment.ID,
			Code:            salesPayment.Code,
			Status:          int32(salesPayment.Status),
			Amount:          salesPayment.Amount,
			RecognitionDate: timestamppb.New(salesPayment.RecognitionDate),
			ReceivedDate:    timestamppb.New(salesPayment.ReceivedDate),
		})
	}

	res = &bridgeService.GetSalesPaymentListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    salesPaymentsGrpc,
	}
	return
}

func (h *BridgeGrpcHandler) GetSalesPaymentDetail(ctx context.Context, req *bridgeService.GetSalesPaymentDetailRequest) (res *bridgeService.GetSalesPaymentDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetSalesPaymentDetail")
	defer span.End()

	var salesPayment dto.SalesPaymentResponse
	salesPayment, err = h.ServiceSalesPayment.GetDetail(ctx, req.Id, "")
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var salesPaymentGrpc *bridgeService.SalesPayment
	salesPaymentGrpc = &bridgeService.SalesPayment{
		Id:              salesPayment.ID,
		Code:            salesPayment.Code,
		Status:          int32(salesPayment.Status),
		Amount:          salesPayment.Amount,
		RecognitionDate: timestamppb.New(salesPayment.RecognitionDate),
		ReceivedDate:    timestamppb.New(salesPayment.ReceivedDate),
	}

	res = &bridgeService.GetSalesPaymentDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    salesPaymentGrpc,
	}
	return
}

func (h *BridgeGrpcHandler) CreateSalesPaymentGP(ctx context.Context, req *bridgeService.CreateSalesPaymentGPRequest) (res *bridgeService.CreateSalesInvoiceGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.CreateSalesPaymentGP")
	defer span.End()

	var resAPi *bridgeService.CreateSalesInvoiceGPResponse
	resAPi, err = h.ServiceSalesPayment.CreateGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.CreateSalesInvoiceGPResponse{
		Code:     resAPi.Code,
		Message:  resAPi.Message,
		Sopnumbe: resAPi.Sopnumbe,
	}
	return
}

func (h *BridgeGrpcHandler) CreateSalesPaymentGPnonPBD(ctx context.Context, req *bridgeService.CreateSalesPaymentGPnonPBDRequest) (res *bridgeService.CreateSalesPaymentGPnonPBDResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.CreateSalesPaymentGP")
	defer span.End()

	var resAPi *bridgeService.CreateSalesPaymentGPnonPBDResponse
	resAPi, err = h.ServiceSalesPayment.CreateGPnonPBD(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.CreateSalesPaymentGPnonPBDResponse{
		Code:          resAPi.Code,
		Message:       resAPi.Message,
		Paymentnumber: resAPi.Paymentnumber,
	}
	return
}

func (h *BridgeGrpcHandler) GetSalesPaymentGPList(ctx context.Context, req *bridgeService.GetSalesPaymentGPListRequest) (res *bridgeService.GetSalesPaymentGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetSalesPaymentGPList")
	defer span.End()

	res, err = h.ServiceSalesPayment.GetGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *BridgeGrpcHandler) GetSalesPaymentGPDetail(ctx context.Context, req *bridgeService.GetSalesPaymentGPDetailRequest) (res *bridgeService.GetSalesPaymentGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetSalesPaymentGPDetail")
	defer span.End()

	res, err = h.ServiceSalesPayment.GetDetailGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}
