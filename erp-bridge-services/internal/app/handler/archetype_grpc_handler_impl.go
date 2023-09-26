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

// GetArchetypeList List
func (h *BridgeGrpcHandler) GetArchetypeList(ctx context.Context, req *bridgeService.GetArchetypeListRequest) (res *bridgeService.GetArchetypeListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetArchetypeList")
	defer span.End()

	var archetypes []dto.ArchetypeResponse
	archetypes, _, err = h.ServicesArchetype.Get(ctx, int(req.Offset), int(req.Limit), int(req.Status), req.Search, req.OrderBy, req.CustomerTypeId)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*bridgeService.Archetype
	for _, archetype := range archetypes {
		data = append(data, &bridgeService.Archetype{
			Id:             archetype.ID,
			Code:           archetype.Code,
			CustomerTypeId: archetype.CustomerTypeID,
			Description:    archetype.Description,
			Status:         int32(archetype.Status),
			CreatedAt:      timestamppb.New(archetype.CreatedAt),
			UpdatedAt:      timestamppb.New(archetype.UpdatedAt),
		})
	}

	res = &bridgeService.GetArchetypeListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

// GetArchetypeDetail
func (h *BridgeGrpcHandler) GetArchetypeDetail(ctx context.Context, req *bridgeService.GetArchetypeDetailRequest) (res *bridgeService.GetArchetypeDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetArchetypeDetail")
	defer span.End()

	var archetype dto.ArchetypeResponse
	archetype, err = h.ServicesArchetype.GetDetail(ctx, req.Id, req.Code)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.GetArchetypeDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &bridgeService.Archetype{
			Id:             archetype.ID,
			Code:           archetype.Code,
			CustomerTypeId: archetype.CustomerTypeID,
			Description:    archetype.Description,
			Status:         int32(archetype.Status),
			CreatedAt:      timestamppb.New(archetype.CreatedAt),
			UpdatedAt:      timestamppb.New(archetype.UpdatedAt),
		},
	}
	return
}

// GetArchetypeGPList endpoint
func (h *BridgeGrpcHandler) GetArchetypeGPList(ctx context.Context, req *bridgeService.GetArchetypeGPListRequest) (res *bridgeService.GetArchetypeGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetArchetypeGPList")
	defer span.End()

	res, err = h.ServicesArchetype.GetGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	return
}

func (h *BridgeGrpcHandler) GetArchetypeGPDetail(ctx context.Context, req *bridgeService.GetArchetypeGPDetailRequest) (res *bridgeService.GetArchetypeGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetArchetypeGPDetail")
	defer span.End()

	res, err = h.ServicesArchetype.GetDetailGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	return
}
