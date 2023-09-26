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

func (h *BridgeGrpcHandler) GetSiteList(ctx context.Context, req *bridgeService.GetSiteListRequest) (res *bridgeService.GetSiteListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetSiteList")
	defer span.End()

	var sites []dto.SiteResponse
	sites, _, err = h.ServicesSite.Get(ctx, int(req.Offset), int(req.Limit), int(req.Status), req.Search, req.OrderBy)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*bridgeService.Site
	for _, site := range sites {
		data = append(data, &bridgeService.Site{
			Id:          site.ID,
			Code:        site.Code,
			Description: site.Description,
			Status:      int32(site.Status),
			CreatedAt:   timestamppb.New(site.CreatedAt),
			UpdatedAt:   timestamppb.New(site.UpdatedAt),
		})
	}

	res = &bridgeService.GetSiteListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *BridgeGrpcHandler) GetSiteInIdsList(ctx context.Context, req *bridgeService.GetSiteInIdsListRequest) (res *bridgeService.GetSiteListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetSiteList")
	defer span.End()

	var sites []dto.SiteResponse
	sites, _, err = h.ServicesSite.GetInIds(ctx, int(req.Offset), int(req.Limit), int(req.Status), req.Ids, req.OrderBy)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*bridgeService.Site
	for _, site := range sites {
		data = append(data, &bridgeService.Site{
			Id:          site.ID,
			Code:        site.Code,
			Description: site.Description,
			Status:      int32(site.Status),
			CreatedAt:   timestamppb.New(site.CreatedAt),
			UpdatedAt:   timestamppb.New(site.UpdatedAt),
		})
	}

	res = &bridgeService.GetSiteListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *BridgeGrpcHandler) GetSiteDetail(ctx context.Context, req *bridgeService.GetSiteDetailRequest) (res *bridgeService.GetSiteDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetSiteDetail")
	defer span.End()

	var site dto.SiteResponse
	site, err = h.ServicesSite.GetDetail(ctx, req.Id, req.Code)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.GetSiteDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &bridgeService.Site{
			Id:          site.ID,
			Code:        site.Code,
			Description: site.Description,
			Status:      int32(site.Status),
			CreatedAt:   timestamppb.New(site.CreatedAt),
			UpdatedAt:   timestamppb.New(site.UpdatedAt),
			Region: &bridgeService.Region{
				Id:          site.Region.ID,
				Code:        site.Region.Code,
				Description: site.Region.Description,
				Status:      int32(site.Region.Status),
				CreatedAt:   timestamppb.New(site.Region.CreatedAt),
				UpdatedAt:   timestamppb.New(site.Region.UpdatedAt),
			},
		},
	}
	return
}

func (h *BridgeGrpcHandler) GetSiteGPList(ctx context.Context, req *bridgeService.GetSiteGPListRequest) (res *bridgeService.GetSiteGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetSiteGPList")
	defer span.End()

	res, err = h.ServicesSite.GetGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *BridgeGrpcHandler) GetSiteGPDetail(ctx context.Context, req *bridgeService.GetSiteGPDetailRequest) (res *bridgeService.GetSiteGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetSiteGPDetail")
	defer span.End()

	res, err = h.ServicesSite.GetDetailGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}
