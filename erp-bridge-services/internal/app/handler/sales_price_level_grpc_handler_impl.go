package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *BridgeGrpcHandler) GetSalesPriceLevelList(ctx context.Context, req *bridgeService.GetSalesPriceLevelListRequest) (res *bridgeService.GetSalesPriceLevelResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetSalesPriceLevelGPList")
	defer span.End()

	res, err = h.ServicesSalesPriceLevel.Get(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	return
}

func (h *BridgeGrpcHandler) GetSalesPriceLevelDetail(ctx context.Context, req *bridgeService.GetSalesPriceLevelDetailRequest) (res *bridgeService.GetSalesPriceLevelResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetSalesPriceLevelGPDetail")
	defer span.End()

	res, err = h.ServicesSalesPriceLevel.GetDetail(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	return
}
