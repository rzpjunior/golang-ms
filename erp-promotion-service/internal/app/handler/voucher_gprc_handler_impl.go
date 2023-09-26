package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/utils"
	dto "git.edenfarm.id/project-version3/erp-services/erp-promotion-service/internal/app/dto"
	promotionService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/promotion_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *PromotionGrpcHandler) GetVoucherMobileList(ctx context.Context, req *promotionService.GetVoucherMobileListRequest) (res *promotionService.GetVoucherMobileListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetVoucherMobileList")
	defer span.End()

	param := &dto.VoucherRequestGetMobileVoucherList{
		RegionID:             req.RegionId,
		CustomerTypeID:       req.CustomerTypeId,
		ArchetypeID:          req.ArchetypeId,
		MembershipLevel:      int8(req.MembershipLevelId),
		MembershipCheckpoint: int8(req.MembershipCheckpointId),
		CustomerID:           req.CustomerId,
		CustomerLevel:        int8(req.CustomerLevelId),
		Offset:               req.Offset,
		Limit:                req.Limit,
		IsMembershipOnly:     req.IsMembershipOnly,
		Category:             int8(req.Category),
	}

	var Vouchers []*dto.VoucherResponse
	Vouchers, _, err = h.ServicesVoucher.GetMobileVoucherList(ctx, param)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*promotionService.Voucher
	for _, voucher := range Vouchers {
		data = append(data, &promotionService.Voucher{
			Id:                     voucher.ID,
			Code:                   voucher.Code,
			RedeemCode:             voucher.RedeemCode,
			Name:                   voucher.Name,
			ImageUrl:               voucher.ImageUrl,
			MinOrder:               voucher.MinOrder,
			DiscAmount:             voucher.DiscAmount,
			RemUserQuota:           voucher.RemUserQuota,
			RemOverallQuota:        voucher.RemOverallQuota,
			OverallQuota:           voucher.OverallQuota,
			UserQuota:              voucher.UserQuota,
			StartTime:              timestamppb.New(voucher.StartTime),
			EndTime:                timestamppb.New(voucher.EndTime),
			VoucherItem:            int32(voucher.VoucherItem),
			MembershipLevelId:      voucher.MembershipLevel.ID,
			MembershipCheckpointId: voucher.MembershipCheckpoint.ID,
			TermConditions:         voucher.TermConditions,
			Type:                   int32(voucher.Type),
			RegionId:               voucher.Region.ID,
			ArchetypeId:            voucher.Archetype.ID,
			CustomerId:             voucher.Customer.ID,
			Status:                 int32(voucher.Status),
			VoidReason:             int32(voucher.VoidReason),
		})
	}

	res = &promotionService.GetVoucherMobileListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *PromotionGrpcHandler) GetVoucherMobileDetail(ctx context.Context, req *promotionService.GetVoucherMobileDetailRequest) (res *promotionService.GetVoucherMobileDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetVoucherMobileDetail")
	defer span.End()

	param := &dto.VoucherRequestGetMobileVoucherDetail{
		CustomerID: req.CustomerId,
		RedeemCode: req.RedeemCode,
	}

	var voucher *dto.VoucherResponse
	voucher, err = h.ServicesVoucher.GetMobileVoucherDetail(ctx, param)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data *promotionService.Voucher

	data = &promotionService.Voucher{
		Id:                     voucher.ID,
		Code:                   voucher.Code,
		RedeemCode:             voucher.RedeemCode,
		Name:                   voucher.Name,
		ImageUrl:               voucher.ImageUrl,
		MinOrder:               voucher.MinOrder,
		DiscAmount:             voucher.DiscAmount,
		RemUserQuota:           voucher.RemUserQuota,
		RemOverallQuota:        voucher.RemOverallQuota,
		OverallQuota:           voucher.OverallQuota,
		UserQuota:              voucher.UserQuota,
		StartTime:              timestamppb.New(voucher.StartTime),
		EndTime:                timestamppb.New(voucher.EndTime),
		VoucherItem:            int32(voucher.VoucherItem),
		MembershipLevelId:      voucher.MembershipLevel.ID,
		MembershipCheckpointId: voucher.MembershipCheckpoint.ID,
		TermConditions:         voucher.TermConditions,
		Type:                   int32(voucher.Type),
		RegionId:               voucher.Region.ID,
		ArchetypeId:            voucher.Archetype.ID,
		CustomerId:             voucher.Customer.ID,
		Status:                 int32(voucher.Status),
		VoidReason:             int32(voucher.VoidReason),
		CustomerTypeId:         voucher.Archetype.CustomerType.ID,
	}

	res = &promotionService.GetVoucherMobileDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *PromotionGrpcHandler) CreateVoucher(ctx context.Context, req *promotionService.CreateVoucherRequest) (res *promotionService.CreateVoucherResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetVoucherMobileDetail")
	defer span.End()

	param := &dto.VoucherRequestCreate{
		RedeemCode:             req.RedeemCode,
		Name:                   req.Name,
		ImageUrl:               req.ImageUrl,
		MinOrder:               utils.ToString(req.MinOrder),
		DiscAmount:             req.DiscAmount,
		OverallQuota:           req.OverallQuota,
		UserQuota:              req.UserQuota,
		StartTime:              req.StartTime.AsTime(),
		EndTime:                req.EndTime.AsTime(),
		MembershipLevelID:      req.MembershipLevelId,
		MembershipCheckPointID: req.MembershipCheckpointId,
		TermConditions:         req.TermConditions,
		Type:                   int8(req.Type),
		RegionID:               req.RegionId,
		ArchetypeID:            req.ArchetypeId,
		CustomerID:             req.CustomerId,
		Status:                 int8(req.Status),
		VoidReason:             int8(req.VoidReason),
	}

	_, err = h.ServicesVoucher.Create(ctx, param)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &promotionService.CreateVoucherResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
	}
	return
}

func (h *PromotionGrpcHandler) UpdateVoucher(ctx context.Context, req *promotionService.UpdateVoucherRequest) (res *promotionService.UpdateVoucherResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetVoucherMobileDetail")
	defer span.End()

	param := &dto.VoucherRequestUpdate{
		VoucherID:       req.VoucherId,
		RemOverallQuota: req.RemOverallQuota,
	}

	err = h.ServicesVoucher.Update(ctx, param)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &promotionService.UpdateVoucherResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
	}
	return
}
