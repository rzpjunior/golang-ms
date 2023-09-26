package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/dto"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/campaign_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *CampaignGrpcHandler) GetMembershipLevelList(ctx context.Context, req *pb.GetMembershipLevelListRequest) (res *pb.GetMembershipLevelListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetMembershipLevelList")
	defer span.End()

	param := &dto.MembershipLevelRequestGet{
		Limit:   int64(req.Limit),
		Offset:  int64(req.Offset),
		Status:  int8(req.Status),
		Search:  req.Search,
		OrderBy: req.OrderBy,
	}

	var membershipLevelList []*dto.MembershipLevelResponse
	membershipLevelList, _, err = h.ServiceMembership.GetMembershipLevelList(ctx, param)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	var data []*pb.MembershipLevel
	for _, v := range membershipLevelList {
		data = append(data, &pb.MembershipLevel{
			Id:       v.ID,
			Code:     v.Code,
			Level:    int32(v.Level),
			ImageUrl: v.ImageUrl,
			Name:     v.Name,
			Status:   int32(v.Status),
		})
	}

	res = &pb.GetMembershipLevelListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *CampaignGrpcHandler) GetMembershipLevelDetail(ctx context.Context, req *pb.GetMembershipLevelDetailRequest) (res *pb.GetMembershipLevelDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetMembershipLevelDetail")
	defer span.End()

	var membershipLevel dto.MembershipLevelResponse

	membershipLevel, err = h.ServiceMembership.GetMembeshipLevelDetail(ctx, req.Id, int8(req.Level))
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data *pb.MembershipLevel

	data = &pb.MembershipLevel{
		Id:       membershipLevel.ID,
		Code:     membershipLevel.Code,
		Level:    int32(membershipLevel.Level),
		ImageUrl: membershipLevel.ImageUrl,
		Name:     membershipLevel.Name,
		Status:   int32(membershipLevel.Status),
	}

	res = &pb.GetMembershipLevelDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *CampaignGrpcHandler) GetMembershipCheckpointList(ctx context.Context, req *pb.GetMembershipCheckpointListRequest) (res *pb.GetMembershipCheckpointListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetMembershipCheckpointList")
	defer span.End()

	param := &dto.MembershipCheckpointRequestGet{
		Offset:            int64(req.Offset),
		Limit:             int64(req.Limit),
		MembershipLevelID: req.MembershipLevelId,
		Status:            int8(req.Status),
		OrderBy:           req.OrderBy,
		ID:                req.Id,
	}

	var membershipCheckpointList []*dto.MembershipCheckpointResponse
	membershipCheckpointList, _, err = h.ServiceMembership.GetMembershipCheckpointList(ctx, param)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	var data []*pb.MembershipCheckpoint
	for _, v := range membershipCheckpointList {
		data = append(data, &pb.MembershipCheckpoint{
			Id:                v.ID,
			Checkpoint:        int32(v.Checkpoint),
			TargetAmount:      v.TargetAmount,
			Status:            int32(v.Status),
			MembershipLevelId: v.MembershipLevelID,
		})
	}

	res = &pb.GetMembershipCheckpointListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *CampaignGrpcHandler) GetMembershipCheckpointDetail(ctx context.Context, req *pb.GetMembershipCheckpointDetailRequest) (res *pb.GetMembershipCheckpointDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetMembershipCheckpointDetail")
	defer span.End()

	var membershipCheckpoint dto.MembershipCheckpointResponse

	membershipCheckpoint, err = h.ServiceMembership.GetMembershipCheckpointDetail(ctx, req.Id, int8(req.Checkpoint))
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data *pb.MembershipCheckpoint

	data = &pb.MembershipCheckpoint{
		Id:                membershipCheckpoint.ID,
		Checkpoint:        int32(membershipCheckpoint.Checkpoint),
		TargetAmount:      membershipCheckpoint.TargetAmount,
		Status:            int32(membershipCheckpoint.Status),
		MembershipLevelId: membershipCheckpoint.MembershipLevelID,
	}

	res = &pb.GetMembershipCheckpointDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *CampaignGrpcHandler) GetMembershipLevelAdvantageList(ctx context.Context, req *pb.GetMembershipLevelAdvantageListRequest) (res *pb.GetMembershipLevelAdvantageListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetMembershipLevelList")
	defer span.End()

	var membershipLevelAdvantageList []*dto.MembershipLevelAdvantage
	membershipLevelAdvantageList, _, err = h.ServiceMembership.GetMembershipLevelAdvantageList(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	var data []*pb.MembershipLevelAdvantage
	for _, v := range membershipLevelAdvantageList {
		data = append(data, &pb.MembershipLevelAdvantage{
			Id:                    v.ID,
			MembershipLevelId:     v.MembershipLevelID,
			MembershipAdvantageId: v.MembershipAdvantageID,
		})
	}

	res = &pb.GetMembershipLevelAdvantageListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *CampaignGrpcHandler) GetMembershipAdvantageDetail(ctx context.Context, req *pb.GetMembershipAdvantageDetailRequest) (res *pb.GetMembershipAdvantageDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetMembershipAdvantageDetail")
	defer span.End()

	var membershipAdvantage dto.MembershipAdvantage

	membershipAdvantage, err = h.ServiceMembership.GetMembershipAdvantageDetail(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data *pb.MembershipAdvantage

	data = &pb.MembershipAdvantage{
		Id:          membershipAdvantage.ID,
		Name:        membershipAdvantage.Name,
		LinkUrl:     membershipAdvantage.LinkUrl,
		Description: membershipAdvantage.Description,
		ImageUrl:    membershipAdvantage.ImageUrl,
		Status:      int32(membershipAdvantage.Status),
	}

	res = &pb.GetMembershipAdvantageDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *CampaignGrpcHandler) GetMembershipRewardList(ctx context.Context, req *pb.GetMembershipRewardListRequest) (res *pb.GetMembershipRewardListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetMembershipRewardList")
	defer span.End()

	var membershipRewardList []*dto.MembershipReward
	membershipRewardList, _, err = h.ServiceMembership.GetMembershipRewardList(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	var data []*pb.MembershipReward
	for _, v := range membershipRewardList {
		data = append(data, &pb.MembershipReward{
			Id:                 v.ID,
			OpenedImageUrl:     v.OpenedImageUrl,
			ClosedImageUrl:     v.ClosedImageUrl,
			BackgroundImageUrl: v.BackgroundImageUrl,
			RewardLevel:        int32(v.RewardLevel),
			MaxAmount:          v.MaxAmount,
			Status:             int32(v.Status),
			Description:        v.Description,
			IsPassed:           int32(v.IsPassed),
			CurrentPercentage:  v.CurrentPercentage,
			RemainingAmount:    v.RemainingAmount,
		})
	}

	res = &pb.GetMembershipRewardListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *CampaignGrpcHandler) GetMembershipRewardDetail(ctx context.Context, req *pb.GetMembershipRewardDetailRequest) (res *pb.GetMembershipRewardDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetMembershipRewardList")
	defer span.End()

	var membershipReward *dto.MembershipReward
	membershipReward, err = h.ServiceMembership.GetMembershipRewardDetail(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	data := &pb.MembershipReward{
		Id:                 membershipReward.ID,
		OpenedImageUrl:     membershipReward.OpenedImageUrl,
		ClosedImageUrl:     membershipReward.ClosedImageUrl,
		BackgroundImageUrl: membershipReward.BackgroundImageUrl,
		RewardLevel:        int32(membershipReward.RewardLevel),
		MaxAmount:          membershipReward.MaxAmount,
		Status:             int32(membershipReward.Status),
		Description:        membershipReward.Description,
		IsPassed:           int32(membershipReward.IsPassed),
		CurrentPercentage:  membershipReward.CurrentPercentage,
		RemainingAmount:    membershipReward.RemainingAmount,
	}

	res = &pb.GetMembershipRewardDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *CampaignGrpcHandler) GetCustomerMembershipDetail(ctx context.Context, req *pb.GetCustomerMembershipDetailRequest) (res *pb.GetCustomerMembershipDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.v")
	defer span.End()

	var customerMembership *dto.CustomerMembership
	customerMembership, err = h.ServiceMembership.GetCustomerMembership(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	data := &pb.CustomerMembership{
		MembershipLevel:      int32(customerMembership.MembershipLevel),
		MembershipLevelName:  customerMembership.MembershipLevelName,
		MembershipCheckpoint: int32(customerMembership.MembershipCheckpoint),
		CheckpointPercentage: customerMembership.CheckpointPercentage,
		CurrentAmount:        customerMembership.CurrentAmount,
		TargetAmount:         customerMembership.TargetAmount,
	}

	res = &pb.GetCustomerMembershipDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}
