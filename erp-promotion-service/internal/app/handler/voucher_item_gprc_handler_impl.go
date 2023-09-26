package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-promotion-service/internal/app/dto"
	promotionService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/promotion_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *PromotionGrpcHandler) GetVoucherItemList(ctx context.Context, req *promotionService.GetVoucherItemListRequest) (res *promotionService.GetVoucherItemListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetVoucherItemList")
	defer span.End()

	param := &dto.VoucherItemRequestGet{
		VoucherID: req.VoucherId,
	}

	var VoucherItems []*dto.VoucherItemResponse
	VoucherItems, _, err = h.ServicesVoucherItem.Get(ctx, param)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*promotionService.VoucherItem
	for _, voucherItem := range VoucherItems {

		data = append(data, &promotionService.VoucherItem{
			Id:         voucherItem.ID,
			VoucherId:  voucherItem.VoucherID,
			MinQtyDisc: voucherItem.MinQtyDisc,
			CreatedAt:  timestamppb.New(voucherItem.CreatedAt),
			ItemId:     voucherItem.ItemID,
		})
	}

	res = &promotionService.GetVoucherItemListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}
