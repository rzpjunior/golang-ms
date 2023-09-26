package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *BridgeGrpcHandler) GetSalesTerritoryGPList(ctx context.Context, req *bridgeService.GetSalesTerritoryGPListRequest) (res *bridgeService.GetSalesTerritoryGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetSalesTerritoryGPList")
	defer span.End()

	res, err = h.ServicesSalesTerritory.GetGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	return
}

func (h *BridgeGrpcHandler) GetSalesTerritoryGPDetail(ctx context.Context, req *bridgeService.GetSalesTerritoryGPDetailRequest) (res *bridgeService.GetSalesTerritoryGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetSalesTerritoryGPDetail")
	defer span.End()

	res, err = h.ServicesSalesTerritory.GetDetailGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	return
}
