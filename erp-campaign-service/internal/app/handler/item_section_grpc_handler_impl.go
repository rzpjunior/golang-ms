package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/dto"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/campaign_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *CampaignGrpcHandler) GetItemSectionList(ctx context.Context, req *pb.GetItemSectionListRequest) (res *pb.GetItemSectionListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetItemSectionList")
	defer span.End()

	param := &dto.ItemSectionRequestGet{
		Offset:        int64(req.Offset),
		Limit:         int64(req.Limit),
		RegionID:      req.RegionId,
		ArchetypeID:   req.ArchetypeId,
		Status:        req.Status,
		Search:        req.Search,
		OrderBy:       req.OrderBy,
		CurrentTime:   req.CurrentTime.AsTime(),
		ItemSectionID: req.ItemSectionId,
		Type:          int8(req.Type),
	}

	var itemSectionList []*dto.ItemSectionResponse

	itemSectionList, _, err = h.ServiceItemSection.GetListMobile(ctx, param)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	var data []*pb.ItemSection
	for _, v := range itemSectionList {
		data = append(data, &pb.ItemSection{
			Id:               v.ID,
			Code:             v.Code,
			Name:             v.Name,
			Regions:          v.Regions,
			Archetypes:       v.Archetypes,
			BackgroundImages: v.BackgroundImage,
			Sequence:         int32(v.Sequence),
			StartAt:          timestamppb.New(v.StartAt),
			FinishAt:         timestamppb.New(v.FinishAt),
			Status:           int32(v.Status),
			CreatedAt:        timestamppb.New(v.CreatedAt),
			UpdatedAt:        timestamppb.New(v.UpdatedAt),
			Type:             int32(v.Type),
			Note:             v.Note,
			Items:            v.ItemID,
		})
	}

	res = &pb.GetItemSectionListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *CampaignGrpcHandler) GetItemSectionDetail(ctx context.Context, req *pb.GetItemSectionDetailRequest) (res *pb.GetItemSectionDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetItemSectionDetail")
	defer span.End()

	var itemSection *dto.ItemSectionResponse

	itemSection, err = h.ServiceItemSection.GetDetailMobile(ctx, req.Id)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data *pb.ItemSection

	data = &pb.ItemSection{
		Id:               itemSection.ID,
		Code:             itemSection.Code,
		Name:             itemSection.Name,
		Regions:          itemSection.Regions,
		Archetypes:       itemSection.Archetypes,
		BackgroundImages: itemSection.BackgroundImage,
		Sequence:         int32(itemSection.Sequence),
		StartAt:          timestamppb.New(itemSection.StartAt),
		FinishAt:         timestamppb.New(itemSection.FinishAt),
		Status:           int32(itemSection.Status),
		CreatedAt:        timestamppb.New(itemSection.CreatedAt),
		UpdatedAt:        timestamppb.New(itemSection.UpdatedAt),
		Type:             int32(itemSection.Type),
		Note:             itemSection.Note,
	}
	for _, item := range itemSection.Items {
		data.Items = append(data.Items, item.ID)
	}

	res = &pb.GetItemSectionDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}
