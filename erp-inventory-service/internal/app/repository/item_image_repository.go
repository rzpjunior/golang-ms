package repository

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-inventory-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-inventory-service/internal/app/model"
)

type IItemImageRepository interface {
	Get(ctx context.Context, limit int64, offset int64, status string, search string, orderBy string, itemID int64, mainImage int8) (itemImages []*model.ItemImage, err error)
	GetDetail(ctx context.Context, id int64, itemID int64, mainImage int8) (itemImage *model.ItemImage, err error)
	GetByItemID(ctx context.Context, itemID int64) (itemImages []*model.ItemImage, err error)
	Create(ctx context.Context, itemImage *model.ItemImage) (err error)
	DeleteByItemID(ctx context.Context, itemID int64) (err error)
}

type ItemImageRepository struct {
	opt opt.Options
}

func NewItemImageRepository() IItemImageRepository {
	return &ItemImageRepository{
		opt: global.Setup.Common,
	}
}
func (r *ItemImageRepository) Get(ctx context.Context, limit int64, offset int64, status string, search string, orderBy string, itemID int64, mainImage int8) (itemImages []*model.ItemImage, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ItemImageRepository.Get")
	defer span.End()

	return
}

func (r *ItemImageRepository) GetDetail(ctx context.Context, id int64, itemID int64, mainImage int8) (itemImage *model.ItemImage, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ItemImageRepository.GetDetail")
	defer span.End()

	var cols []string
	itemImage = &model.ItemImage{}

	if id != 0 {
		itemImage.ID = id
		cols = append(cols, "id")
	}

	if itemID != 0 {
		itemImage.ItemID = itemID
		cols = append(cols, "item_id")
	}

	if mainImage != 0 {
		itemImage.MainImage = mainImage
		cols = append(cols, "main_image")
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, itemImage, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *ItemImageRepository) GetByItemID(ctx context.Context, itemID int64) (itemImages []*model.ItemImage, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ItemImageRepository.Get")
	defer span.End()

	db := r.opt.Database.Read

	_, err = db.QueryTable(new(model.ItemImage)).Filter("item_id", itemID).OrderBy("main_image").AllWithCtx(ctx, &itemImages)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *ItemImageRepository) Create(ctx context.Context, itemImage *model.ItemImage) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ItemImageRepository.Create")
	defer span.End()

	db := r.opt.Database.Write

	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.InsertWithCtx(ctx, itemImage)
	if err != nil {
		span.RecordError(err)
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *ItemImageRepository) DeleteByItemID(ctx context.Context, itemID int64) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ItemImageRepository.Delete")
	defer span.End()

	db := r.opt.Database.Write

	_, err = db.QueryTable(new(model.ItemImage)).Filter("item_id", itemID).DeleteWithCtx(ctx)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}
