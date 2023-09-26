package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-inventory-service/internal/app/dto"
	catalogService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/catalog_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *CatalogGrpcHandler) GetItemImageList(ctx context.Context, req *catalogService.GetItemImageListRequest) (res *catalogService.GetItemImageListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetItemImageList")
	defer span.End()

	var itemImages []dto.ItemImageResponse
	itemImages, _, err = h.ServicesItemImage.Get(ctx, int(req.Offset), int(req.Limit), int(req.Status), req.Search, req.OrderBy, req.ItemId, int8(req.MainImage))
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*catalogService.ItemImage
	for _, itemImage := range itemImages {
		data = append(data, &catalogService.ItemImage{
			Id:        itemImage.ID,
			ItemId:    itemImage.ItemID,
			ImageUrl:  itemImage.ImageUrl,
			MainImage: int32(itemImage.MainImage),
			CreatedAt: timestamppb.New(itemImage.CreatedAt),
		})
	}

	res = &catalogService.GetItemImageListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *CatalogGrpcHandler) GetItemImageDetail(ctx context.Context, req *catalogService.GetItemImageDetailRequest) (res *catalogService.GetItemImageDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetItemImageDetail")
	defer span.End()

	var itemImage dto.ItemImageResponse
	itemImage, err = h.ServicesItemImage.GetDetail(ctx, req.Id, req.ItemId, int8(req.MainImage))
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &catalogService.GetItemImageDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &catalogService.ItemImage{
			Id:        itemImage.ID,
			ItemId:    itemImage.ItemID,
			ImageUrl:  itemImage.ImageUrl,
			MainImage: int32(itemImage.MainImage),
			CreatedAt: timestamppb.New(itemImage.CreatedAt),
		},
	}
	return
}
