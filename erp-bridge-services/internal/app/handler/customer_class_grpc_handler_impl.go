package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *BridgeGrpcHandler) GetCustomerClassList(ctx context.Context, req *bridgeService.GetCustomerClassListRequest) (res *bridgeService.GetCustomerClassResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetCustomerClassGPList")
	defer span.End()

	res, err = h.ServicesCustomerClass.Get(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	return
}

func (h *BridgeGrpcHandler) GetCustomerClassDetail(ctx context.Context, req *bridgeService.GetCustomerClassDetailRequest) (res *bridgeService.GetCustomerClassResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetCustomerClassGPDetail")
	defer span.End()

	res, err = h.ServicesCustomerClass.GetDetail(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	return
}
