package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *BridgeGrpcHandler) GetHelperGPList(ctx context.Context, req *bridgeService.GetHelperGPListRequest) (res *bridgeService.GetHelperGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetHelperList")
	defer span.End()

	res, err = h.ServicesHelper.GetGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	return
}

func (h *BridgeGrpcHandler) GetHelperGPDetail(ctx context.Context, req *bridgeService.GetHelperGPDetailRequest) (res *bridgeService.GetHelperGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetHelperDetail")
	defer span.End()

	res, err = h.ServicesHelper.GetDetailGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *BridgeGrpcHandler) LoginHelper(ctx context.Context, req *bridgeService.LoginHelperRequest) (res *bridgeService.LoginHelperResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.LoginHelper")
	defer span.End()

	res, err = h.ServicesHelper.Login(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}
