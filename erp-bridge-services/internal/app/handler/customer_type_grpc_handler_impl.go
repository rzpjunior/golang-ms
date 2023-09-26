package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *BridgeGrpcHandler) GetCustomerTypeGPList(ctx context.Context, req *bridgeService.GetCustomerTypeGPListRequest) (res *bridgeService.GetCustomerTypeGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetCustomerTypeGPList")
	defer span.End()

	res, err = h.ServicesCustomerType.GetGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	return
}

func (h *BridgeGrpcHandler) GetCustomerTypeGPDetail(ctx context.Context, req *bridgeService.GetCustomerTypeGPDetailRequest) (res *bridgeService.GetCustomerTypeGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetCustomerTypeGPDetail")
	defer span.End()

	res, err = h.ServicesCustomerType.GetDetailGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	return
}
