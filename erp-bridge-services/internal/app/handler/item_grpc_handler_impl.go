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

func (h *BridgeGrpcHandler) GetItemList(ctx context.Context, req *bridgeService.GetItemListRequest) (res *bridgeService.GetItemListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetItemList")
	defer span.End()

	var items []dto.ItemResponse
	items, _, err = h.ServicesItem.Get(ctx, int(req.Offset), int(req.Limit), int(req.Status), req.Search, req.OrderBy, req.UomId, req.ClassId, req.ItemCategoryId)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*bridgeService.Item
	for _, item := range items {
		data = append(data, &bridgeService.Item{
			Id:                      item.ID,
			Code:                    item.Code,
			UomId:                   item.UomID,
			ClassId:                 item.ClassID,
			ItemCategoryId:          item.ItemCategoryID,
			Description:             item.Description,
			UnitWeightConversion:    item.UnitWeightConversion,
			OrderMinQty:             item.OrderMinQty,
			OrderMaxQty:             item.OrderMaxQty,
			ItemType:                item.ItemType,
			Packability:             item.Packability,
			Capitalize:              item.Capitalize,
			ExcludeArchetype:        item.ExcludeArchetype,
			MaxDayDeliveryDate:      int32(item.MaxDayDeliveryDate),
			FragileGoods:            item.FragileGoods,
			Taxable:                 item.Taxable,
			OrderChannelRestriction: item.OrderChannelRestriction,
			Note:                    item.Note,
			Status:                  int32(item.Status),
			CreatedAt:               timestamppb.New(item.CreatedAt),
			UpdatedAt:               timestamppb.New(item.UpdatedAt),
		})
	}

	res = &bridgeService.GetItemListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *BridgeGrpcHandler) GetItemDetail(ctx context.Context, req *bridgeService.GetItemDetailRequest) (res *bridgeService.GetItemDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetItemDetail")
	defer span.End()

	var item dto.ItemResponse
	item, err = h.ServicesItem.GetDetail(ctx, req.Id, req.Code)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.GetItemDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &bridgeService.Item{
			Id:                      item.ID,
			Code:                    item.Code,
			UomId:                   item.UomID,
			ClassId:                 item.ClassID,
			ItemCategoryId:          item.ItemCategoryID,
			Description:             item.Description,
			UnitWeightConversion:    item.UnitWeightConversion,
			OrderMinQty:             item.OrderMinQty,
			OrderMaxQty:             item.OrderMaxQty,
			ItemType:                item.ItemType,
			Packability:             item.Packability,
			Capitalize:              item.Capitalize,
			ExcludeArchetype:        item.ExcludeArchetype,
			MaxDayDeliveryDate:      int32(item.MaxDayDeliveryDate),
			FragileGoods:            item.FragileGoods,
			Taxable:                 item.Taxable,
			OrderChannelRestriction: item.OrderChannelRestriction,
			Note:                    item.Note,
			Status:                  int32(item.Status),
			CreatedAt:               timestamppb.New(item.CreatedAt),
			UpdatedAt:               timestamppb.New(item.UpdatedAt),
		},
	}
	return
}

func (h *BridgeGrpcHandler) UpdateItemPackable(ctx context.Context, req *bridgeService.UpdateItemPackableRequest) (res *bridgeService.UpdateItemPackableResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.UpdateItemPackable")
	defer span.End()

	res = &bridgeService.UpdateItemPackableResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
	}
	return
}

func (h *BridgeGrpcHandler) UpdateItemFragile(ctx context.Context, req *bridgeService.UpdateItemFragileRequest) (res *bridgeService.UpdateItemFragileResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.UpdateItemFragile")
	defer span.End()

	res = &bridgeService.UpdateItemFragileResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
	}
	return
}

func (h *BridgeGrpcHandler) GetItemGPList(ctx context.Context, req *bridgeService.GetItemGPListRequest) (res *bridgeService.GetItemGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetItemGPList")
	defer span.End()

	res, err = h.ServicesItem.GetGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *BridgeGrpcHandler) GetItemGPDetail(ctx context.Context, req *bridgeService.GetItemGPDetailRequest) (res *bridgeService.GetItemGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetItemGPDetail")
	defer span.End()

	res, err = h.ServicesItem.GetDetailGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *BridgeGrpcHandler) GetItemMasterComplexGP(ctx context.Context, req *bridgeService.GetItemMasterComplexGPListRequest) (res *bridgeService.GetItemMasterComplexGPListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetItemGPList")
	defer span.End()

	res, err = h.ServicesItem.GetItemMasterComplexGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}
