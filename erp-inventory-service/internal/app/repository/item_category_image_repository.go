package repository

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-inventory-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-inventory-service/internal/app/model"
)

type IItemCategoryImageRepository interface {
	GetByItemCategoryID(ctx context.Context, itemCategoryID int64) (itemCategoryImage *model.ItemCategoryImage, err error)
	Create(ctx context.Context, itemCategoryImage *model.ItemCategoryImage) (err error)
	DeleteByItemID(ctx context.Context, itemCategoryID int64) (err error)
}

type ItemCategoryImageRepository struct {
	opt opt.Options
}

func NewItemCategoryImageRepository() IItemCategoryImageRepository {
	return &ItemCategoryImageRepository{
		opt: global.Setup.Common,
	}
}

func (r *ItemCategoryImageRepository) GetByItemCategoryID(ctx context.Context, itemCategoryID int64) (itemCategoryImage *model.ItemCategoryImage, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ItemCategoryImageRepository.Get")
	defer span.End()

	db := r.opt.Database.Read

	itemCategoryImage = &model.ItemCategoryImage{
		ItemCategoryID: itemCategoryID,
	}

	err = db.ReadWithCtx(ctx, itemCategoryImage, "item_category_id")
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *ItemCategoryImageRepository) Create(ctx context.Context, itemCategoryImage *model.ItemCategoryImage) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ItemCategoryImageRepository.Create")
	defer span.End()

	db := r.opt.Database.Write

	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.InsertWithCtx(ctx, itemCategoryImage)
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

func (r *ItemCategoryImageRepository) DeleteByItemID(ctx context.Context, itemCategoryID int64) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ItemCategoryImageRepository.Delete")
	defer span.End()

	db := r.opt.Database.Write

	_, err = db.QueryTable(new(model.ItemCategoryImage)).Filter("item_category_id", itemCategoryID).DeleteWithCtx(ctx)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}
