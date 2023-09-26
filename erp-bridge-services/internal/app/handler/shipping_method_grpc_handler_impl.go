package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *BridgeGrpcHandler) GetShippingMethodList(ctx context.Context, req *bridgeService.GetShippingMethodListRequest) (res *bridgeService.GetShippingMethodResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetShippingMethodGPList")
	defer span.End()

	res, err = h.ServicesShippingMethod.Get(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	return
}

func (h *BridgeGrpcHandler) GetShippingMethodDetail(ctx context.Context, req *bridgeService.GetShippingMethodDetailRequest) (res *bridgeService.GetShippingMethodResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetShippingMethodGPDetail")
	defer span.End()

	res, err = h.ServicesShippingMethod.GetDetail(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	return
}
