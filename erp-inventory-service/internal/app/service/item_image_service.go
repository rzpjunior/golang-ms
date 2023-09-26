package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-inventory-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-inventory-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-inventory-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-inventory-service/internal/app/repository"
)

type IItemImageService interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, itemID int64, mainImage int8) (res []dto.ItemImageResponse, total int64, err error)
	GetDetail(ctx context.Context, id int64, itemID int64, mainImage int8) (res dto.ItemImageResponse, err error)
}

type ItemImageService struct {
	opt                 opt.Options
	RepositoryItemImage repository.IItemImageRepository
}

func NewItemImageService() IItemImageService {
	return &ItemImageService{
		opt:                 global.Setup.Common,
		RepositoryItemImage: repository.NewItemImageRepository(),
	}
}

func (s *ItemImageService) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, itemID int64, mainImage int8) (res []dto.ItemImageResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemImageService.Get")
	defer span.End()

	return
}

func (s *ItemImageService) GetDetail(ctx context.Context, id int64, itemID int64, mainImage int8) (res dto.ItemImageResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemImageService.GetDetail")
	defer span.End()

	var itemImage *model.ItemImage
	itemImage, err = s.RepositoryItemImage.GetDetail(ctx, id, itemID, mainImage)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.ItemImageResponse{
		ID:        itemImage.ID,
		ItemID:    itemImage.ItemID,
		ImageUrl:  itemImage.ImageUrl,
		MainImage: itemImage.MainImage,
	}

	return
}
