package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-inventory-service/internal/app/dto"
	catalogService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/catalog_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *CatalogGrpcHandler) GetItemCategoryList(ctx context.Context, req *catalogService.GetItemCategoryListRequest) (res *catalogService.GetItemCategoryListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetItemCategoryList")
	defer span.End()

	var ItemCategoryes []dto.ItemCategoryResponse
	ItemCategoryes, _, err = h.ServicesItemCategory.Get(ctx, int(req.Offset), int(req.Limit), int(req.Status), req.Search, req.OrderBy, req.RegionId)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*catalogService.ItemCategory
	for _, ItemCategory := range ItemCategoryes {

		data = append(data, &catalogService.ItemCategory{
			Id:       ItemCategory.ID,
			RegionId: ItemCategory.RegionID,
			Name:     ItemCategory.Name,
			Status:   int32(ItemCategory.Status),
			ImageUrl: ItemCategory.ItemCategoryImage.ImageUrl,
			Code:     ItemCategory.Code,
		})
	}

	res = &catalogService.GetItemCategoryListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *CatalogGrpcHandler) GetItemCategoryDetail(ctx context.Context, req *catalogService.GetItemCategoryDetailRequest) (res *catalogService.GetItemCategoryDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetItemCategoryDetail")
	defer span.End()

	var ItemCategory dto.ItemCategoryResponse

	ItemCategory, err = h.ServicesItemCategory.GetByID(ctx, req.Id)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data *catalogService.ItemCategory

	data = &catalogService.ItemCategory{
		Id:       ItemCategory.ID,
		RegionId: ItemCategory.RegionID,
		Name:     ItemCategory.Name,
		Status:   int32(ItemCategory.Status),
		ImageUrl: ItemCategory.ItemCategoryImage.ImageUrl,
		Code:     ItemCategory.Code,
	}

	res = &catalogService.GetItemCategoryDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}
