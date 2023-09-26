package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/site_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LoginHelper
func (h *SiteGrpcHandler) LoginHelper(ctx context.Context, req *pb.LoginHelperRequest) (res *pb.LoginHelperResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.LoginHelper")
	defer span.End()

	res, err = h.ServicesPickingOrder.Login(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

// GetPickingOrderHeader
func (h *SiteGrpcHandler) GetPickingOrderHeader(ctx context.Context, req *pb.GetPickingOrderHeaderRequest) (res *pb.GetPickingOrderHeaderResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetPickingOrderHeader")
	defer span.End()

	res, err = h.ServicesPickingOrder.GetGrpc(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *SiteGrpcHandler) GetPickingOrderDetail(ctx context.Context, req *pb.GetPickingOrderDetailRequest) (res *pb.GetPickingOrderDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetPickingOrderDetail")
	defer span.End()

	res, err = h.ServicesPickingOrder.GetDetailGrpc(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *SiteGrpcHandler) GetAggregatedProductSalesOrder(ctx context.Context, req *pb.GetAggregatedProductSalesOrderRequest) (res *pb.GetAggregatedProductSalesOrderResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetAggregatedProductSalesOrder")
	defer span.End()

	res, err = h.ServicesPickingOrder.GetDetailAggregatedProduct(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *SiteGrpcHandler) StartPickingOrder(ctx context.Context, req *pb.StartPickingOrderRequest) (res *pb.SuccessResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.StartPickingOrder")
	defer span.End()

	res, err = h.ServicesPickingOrder.StartPickingOrder(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *SiteGrpcHandler) SubmitPicking(ctx context.Context, req *pb.SubmitPickingRequest) (res *pb.SuccessResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.SubmitPicking")
	defer span.End()

	res, err = h.ServicesPickingOrder.SubmitPicking(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *SiteGrpcHandler) GetSalesOrderPicking(ctx context.Context, req *pb.GetSalesOrderPickingRequest) (res *pb.GetSalesOrderPickingResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetSalesOrderPicking")
	defer span.End()

	res, err = h.ServicesPickingOrder.GetSalesOrderPicking(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *SiteGrpcHandler) GetSalesOrderPickingDetail(ctx context.Context, req *pb.GetSalesOrderPickingDetailRequest) (res *pb.GetSalesOrderPickingDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetSalesOrderPickingDetail")
	defer span.End()

	res, err = h.ServicesPickingOrder.GetSalesOrderPickingDetail(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *SiteGrpcHandler) SubmitSalesOrder(ctx context.Context, req *pb.SubmitSalesOrderRequest) (res *pb.SuccessResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.SubmitSalesOrder")
	defer span.End()

	res, err = h.ServicesPickingOrder.SubmitSalesOrder(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *SiteGrpcHandler) GetSalesOrderToCheck(ctx context.Context, req *pb.GetSalesOrderToCheckRequest) (res *pb.GetSalesOrderToCheckResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetSalesOrderToCheck")
	defer span.End()

	res, err = h.ServicesPickingOrder.GetSalesOrderToCheck(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *SiteGrpcHandler) SPVGetSalesOrderToCheckDetail(ctx context.Context, req *pb.GetSalesOrderToCheckDetailRequest) (res *pb.GetSalesOrderToCheckDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.SPVGetSalesOrderToCheckDetail")
	defer span.End()

	res, err = h.ServicesPickingOrder.SPVGetSalesOrderToCheckDetail(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *SiteGrpcHandler) SPVRejectSalesOrder(ctx context.Context, req *pb.SPVRejectSalesOrderRequest) (res *pb.SuccessResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.SPVRejectSalesOrder")
	defer span.End()

	res, err = h.ServicesPickingOrder.SPVRejectSalesOrder(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *SiteGrpcHandler) SPVAcceptSalesOrder(ctx context.Context, req *pb.SPVAcceptSalesOrderRequest) (res *pb.SuccessResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.SPVAcceptSalesOrder")
	defer span.End()

	res, err = h.ServicesPickingOrder.SPVAcceptSalesOrder(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *SiteGrpcHandler) CheckerGetSalesOrderToCheckDetail(ctx context.Context, req *pb.GetSalesOrderToCheckDetailRequest) (res *pb.GetSalesOrderToCheckDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.CheckerGetSalesOrderToCheckDetail")
	defer span.End()

	res, err = h.ServicesPickingOrder.CheckerGetSalesOrderToCheckDetail(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *SiteGrpcHandler) CheckerStartChecking(ctx context.Context, req *pb.CheckerStartCheckingRequest) (res *pb.SuccessResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.CheckerStartChecking")
	defer span.End()

	res, err = h.ServicesPickingOrder.CheckerStartChecking(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *SiteGrpcHandler) CheckerSubmitChecking(ctx context.Context, req *pb.CheckerSubmitCheckingRequest) (res *pb.SuccessResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.CheckerSubmitChecking")
	defer span.End()

	res, err = h.ServicesPickingOrder.CheckerSubmitChecking(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *SiteGrpcHandler) CheckerRejectSalesOrder(ctx context.Context, req *pb.CheckerRejectSalesOrderRequest) (res *pb.SuccessResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.CheckerRejectSalesOrder")
	defer span.End()

	res, err = h.ServicesPickingOrder.CheckerRejectSalesOrder(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *SiteGrpcHandler) CheckerGetDeliveryKoli(ctx context.Context, req *pb.CheckerGetDeliveryKoliRequest) (res *pb.CheckerGetDeliveryKoliResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.CheckerGetDeliveryKoli")
	defer span.End()

	res, err = h.ServicesPickingOrder.CheckerGetDeliveryKoli(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *SiteGrpcHandler) CheckerAcceptSalesOrder(ctx context.Context, req *pb.CheckerAcceptSalesOrderRequest) (res *pb.CheckerAcceptSalesOrderResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.CheckerAcceptSalesOrder")
	defer span.End()

	res, err = h.ServicesPickingOrder.CheckerAcceptSalesOrder(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *SiteGrpcHandler) CheckerHistory(ctx context.Context, req *pb.CheckerHistoryRequest) (res *pb.CheckerHistoryResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.CheckerHistory")
	defer span.End()

	res, err = h.ServicesPickingOrder.CheckerHistory(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *SiteGrpcHandler) CheckerHistoryDetail(ctx context.Context, req *pb.CheckerHistoryDetailRequest) (res *pb.CheckerHistoryDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.CheckerHistoryDetail")
	defer span.End()

	res, err = h.ServicesPickingOrder.CheckerHistoryDetail(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *SiteGrpcHandler) PickerWidget(ctx context.Context, req *pb.PickerWidgetRequest) (res *pb.PickerWidgetResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.PickerWidget")
	defer span.End()

	res, err = h.ServicesPickingOrder.PickerWidget(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *SiteGrpcHandler) SPVWidget(ctx context.Context, req *pb.SPVWidgetRequest) (res *pb.SPVWidgetResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.SPVWidget")
	defer span.End()

	res, err = h.ServicesPickingOrder.SPVWidget(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *SiteGrpcHandler) CheckerWidget(ctx context.Context, req *pb.CheckerWidgetRequest) (res *pb.CheckerWidgetResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.CheckerWidget")
	defer span.End()

	res, err = h.ServicesPickingOrder.CheckerWidget(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *SiteGrpcHandler) History(ctx context.Context, req *pb.HistoryRequest) (res *pb.HistoryResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.History")
	defer span.End()

	res, err = h.ServicesPickingOrder.History(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *SiteGrpcHandler) HistoryDetail(ctx context.Context, req *pb.HistoryDetailRequest) (res *pb.HistoryDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.HistoryDetail")
	defer span.End()

	res, err = h.ServicesPickingOrder.HistoryDetail(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *SiteGrpcHandler) SPVWrtMonitoring(ctx context.Context, req *pb.GetWrtMonitoringListRequest) (res *pb.GetWrtMonitoringListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.SPVWrtMonitoring")
	defer span.End()

	res, err = h.ServicesPickingOrder.SPVWrtMonitoring(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *SiteGrpcHandler) SPVWrtMonitoringDetail(ctx context.Context, req *pb.GetWrtMonitoringDetailRequest) (res *pb.GetWrtMonitoringDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.SPVWrtMonitoringDetail")
	defer span.End()

	res, err = h.ServicesPickingOrder.SPVWrtMonitoringDetail(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}
