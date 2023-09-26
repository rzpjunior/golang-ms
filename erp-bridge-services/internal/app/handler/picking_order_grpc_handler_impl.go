package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *BridgeGrpcHandler) GetPickingOrderGPHeader(ctx context.Context, req *pb.GetPickingOrderGPHeaderRequest) (res *pb.GetPickingOrderGPHeaderResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetPickingOrderGPHeader")
	defer span.End()

	res, err = h.ServicePickingOrder.GetGrpc(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

// GetPickingOrderGPDetail
func (h *BridgeGrpcHandler) GetPickingOrderGPDetail(ctx context.Context, req *pb.GetPickingOrderGPDetailRequest) (res *pb.GetPickingOrderGPDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetPickingOrderGPDetail")
	defer span.End()

	res, err = h.ServicePickingOrder.GetDetailGrpc(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *BridgeGrpcHandler) SubmitPickingCheckingPickingOrder(ctx context.Context, req *pb.SubmitPickingCheckingRequest) (res *pb.SubmitPickingCheckingResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.SubmitPickingChecking")
	defer span.End()

	res, err = h.ServicePickingOrder.SubmitPickingChecking(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}
