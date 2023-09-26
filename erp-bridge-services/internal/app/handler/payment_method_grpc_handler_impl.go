package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *BridgeGrpcHandler) GetPaymentMethodGPList(ctx context.Context, req *bridgeService.GetPaymentMethodGPListRequest) (res *bridgeService.GetPaymentMethodGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetPaymentMethodGPList")
	defer span.End()

	res, err = h.ServicesPaymentMethod.GetGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	return
}

func (h *BridgeGrpcHandler) GetPaymentMethodGPDetail(ctx context.Context, req *bridgeService.GetPaymentMethodGPDetailRequest) (res *bridgeService.GetPaymentMethodGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetPaymentMethodGPDetail")
	defer span.End()

	res, err = h.ServicesPaymentMethod.GetDetailGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	return
}
