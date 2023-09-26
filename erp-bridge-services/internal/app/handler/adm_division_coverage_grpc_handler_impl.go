package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetAdmDivisionCoverageGPList
func (h *BridgeGrpcHandler) GetAdmDivisionCoverageGPList(ctx context.Context, req *bridgeService.GetAdmDivisionCoverageGPListRequest) (res *bridgeService.GetAdmDivisionCoverageGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetAdmDivisionCoverageGPList")
	defer span.End()

	res, err = h.ServiceAdmDivisionCoverage.GetGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	return
}

// GetAdmDivisionCoverageGPDetail
func (h *BridgeGrpcHandler) GetAdmDivisionCoverageGPDetail(ctx context.Context, req *bridgeService.GetAdmDivisionCoverageGPDetailRequest) (res *bridgeService.GetAdmDivisionCoverageGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetAdmDivisionCoverageGPDetail")
	defer span.End()

	res, err = h.ServiceAdmDivisionCoverage.GetDetailGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	return
}
