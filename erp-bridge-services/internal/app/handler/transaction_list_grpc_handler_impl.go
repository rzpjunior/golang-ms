package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *BridgeGrpcHandler) GetTransactionListGPList(ctx context.Context, req *bridgeService.GetTransactionListGPListRequest) (res *bridgeService.GetTransactionListGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetTransactionListGPList")
	defer span.End()

	res, err = h.ServicesTransactionList.GetGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	return
}

func (h *BridgeGrpcHandler) GetTransactionListGPDetail(ctx context.Context, req *bridgeService.GetTransactionListGPDetailRequest) (res *bridgeService.GetTransactionListGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetTransactionListGPDetail")
	defer span.End()

	res, err = h.ServicesTransactionList.GetDetailGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	return
}
