package handler

import (
	context "context"
	"fmt"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/dto"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/campaign_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *CampaignGrpcHandler) GetBannerList(ctx context.Context, req *pb.GetBannerListRequest) (res *pb.GetBannerListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.Get")
	defer span.End()

	param := &dto.BannerRequestGet{
		Offset:      int64(req.Offset),
		Limit:       int64(req.Limit),
		RegionID:    req.RegionId,
		ArchetypeID: req.ArchetypeId,
		Status:      req.Status,
		Search:      req.Search,
		OrderBy:     req.OrderBy,
		CurrentTime: req.CurrentTime.AsTime(),
	}

	bannerList, _, err := h.ServiceBanner.GetListMobile(ctx, param)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*pb.Banner
	for _, v := range bannerList {
		var redirectValue string
		redirectValue = fmt.Sprintf("%v", v.Redirect.Value)
		data = append(data, &pb.Banner{
			Id:            v.ID,
			Code:          v.Code,
			Name:          v.Name,
			Regions:       v.Regions,
			Archetypes:    v.Archetypes,
			Queue:         int32(v.Queue),
			RedirectTo:    int32(v.Redirect.To),
			RedirectValue: redirectValue,
			ImageUrl:      v.ImageUrl,
			StartAt:       timestamppb.New(v.StartAt),
			FinishAt:      timestamppb.New(v.FinishAt),
			Note:          v.Note,
			Status:        int32(v.Status),
			CreatedAt:     timestamppb.New(v.CreatedAt),
			UpdatedAt:     timestamppb.New(v.UpdatedAt),
		})
	}

	res = &pb.GetBannerListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *CampaignGrpcHandler) GetBannerDetail(ctx context.Context, req *pb.GetBannerDetailRequest) (res *pb.GetBannerDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetBannerDetail")
	defer span.End()

	banner, err := h.ServiceBanner.GetByID(ctx, req.Id)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &pb.GetBannerDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &pb.Banner{
			Id:            banner.ID,
			Code:          banner.Code,
			Name:          banner.Name,
			Regions:       banner.Regions,
			Archetypes:    banner.Archetypes,
			Queue:         int32(banner.Queue),
			RedirectTo:    int32(banner.Redirect.To),
			RedirectValue: banner.Redirect.ValueName,
			ImageUrl:      banner.ImageUrl,
			StartAt:       timestamppb.New(banner.StartAt),
			FinishAt:      timestamppb.New(banner.FinishAt),
			Note:          banner.Note,
			Status:        int32(banner.Status),
			CreatedAt:     timestamppb.New(banner.CreatedAt),
			UpdatedAt:     timestamppb.New(banner.UpdatedAt),
		},
	}

	return
}
