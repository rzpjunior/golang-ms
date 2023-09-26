package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *BridgeGrpcHandler) GetClassList(ctx context.Context, req *bridgeService.GetClassListRequest) (res *bridgeService.GetClassListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetClassList")
	defer span.End()

	var classes []dto.ClassResponse
	classes, _, err = h.ServicesClass.Get(ctx, int(req.Offset), int(req.Limit), int(req.Status), req.Search, req.OrderBy)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*bridgeService.Class
	for _, class := range classes {
		data = append(data, &bridgeService.Class{
			Id:          class.ID,
			Code:        class.Code,
			Description: class.Description,
			Status:      int32(class.Status),
			CreatedAt:   timestamppb.New(class.CreatedAt),
			UpdatedAt:   timestamppb.New(class.UpdatedAt),
		})
	}

	res = &bridgeService.GetClassListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *BridgeGrpcHandler) GetClassDetail(ctx context.Context, req *bridgeService.GetClassDetailRequest) (res *bridgeService.GetClassDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetClassDetail")
	defer span.End()

	var class dto.ClassResponse
	class, err = h.ServicesClass.GetDetail(ctx, req.Id, req.Code)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.GetClassDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &bridgeService.Class{
			Id:          class.ID,
			Code:        class.Code,
			Description: class.Description,
			Status:      int32(class.Status),
			CreatedAt:   timestamppb.New(class.CreatedAt),
			UpdatedAt:   timestamppb.New(class.UpdatedAt),
		},
	}
	return
}
