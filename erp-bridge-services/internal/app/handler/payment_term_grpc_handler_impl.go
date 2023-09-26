package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *BridgeGrpcHandler) GetPaymentTermGPList(ctx context.Context, req *bridgeService.GetPaymentTermGPListRequest) (res *bridgeService.GetPaymentTermGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetPaymentTermGPList")
	defer span.End()

	res, err = h.ServicesPaymentTerm.GetGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	return
}

func (h *BridgeGrpcHandler) GetPaymentTermGPDetail(ctx context.Context, req *bridgeService.GetPaymentTermGPDetailRequest) (res *bridgeService.GetPaymentTermGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetPaymentTermGPDetail")
	defer span.End()

	res, err = h.ServicesPaymentTerm.GetDetailGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	return
}
