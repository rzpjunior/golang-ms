package repository

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/project-version3/erp-services/erp-inventory-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-inventory-service/internal/app/model"
)

type IItemCategoryRepository interface {
	Get(ctx context.Context, offset, limit, status int, search, orderBy string, regionID string) (items []*model.ItemCategory, count int64, err error)
	GetByID(ctx context.Context, id int64) (itemCategory *model.ItemCategory, err error)
	Update(ctx context.Context, ItemCategory *model.ItemCategory, columns ...string) (err error)
	Create(ctx context.Context, ItemCategory *model.ItemCategory) (res int64, err error)
	IsExistNameItemCategory(ctx context.Context, id int64, name string) (isExist bool)
}

type ItemCategoryRepository struct {
	opt opt.Options
}

func NewItemCategoryRepository() IItemCategoryRepository {
	return &ItemCategoryRepository{
		opt: global.Setup.Common,
	}
}

func (r *ItemCategoryRepository) Get(ctx context.Context, offset, limit, status int, search, orderBy string, regionID string) (items []*model.ItemCategory, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ItemCategoryRepository.Get")
	defer span.End()

	db := r.opt.Database.Read

	cond := orm.NewCondition()

	qs := db.QueryTable(new(model.ItemCategory))

	if search != "" {
		cond1 := orm.NewCondition()
		cond1 = cond1.And("name__icontains", search).Or("code__icontains", search)
		cond = cond.AndCond(cond1)
	}

	if status != 0 {
		cond = cond.And("status", status)
	}

	if regionID != "" {
		cond = cond.And("regions__icontains", regionID)
	}

	qs = qs.SetCond(cond)

	if orderBy != "" {
		qs = qs.OrderBy(orderBy)
	}

	_, err = qs.Offset(offset).Limit(limit).AllWithCtx(ctx, &items)
	if err != nil {
		span.RecordError(err)
		return
	}

	count, err = qs.Count()

	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *ItemCategoryRepository) GetByID(ctx context.Context, id int64) (itemCategory *model.ItemCategory, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ItemCategoryRepository.GetByID")
	defer span.End()

	itemCategory = &model.ItemCategory{
		ID: id,
	}
	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, itemCategory, "id")
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *ItemCategoryRepository) Update(ctx context.Context, ItemCategory *model.ItemCategory, columns ...string) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ItemCategoryRepository.Update")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.UpdateWithCtx(ctx, ItemCategory, columns...)

	if err == nil {
		err = tx.Commit()
		if err != nil {
			span.RecordError(err)
			return
		}
	} else {
		span.RecordError(err)
		err = tx.Rollback()
		if err != nil {
			span.RecordError(err)
			return
		}
	}
	return
}

func (r *ItemCategoryRepository) Create(ctx context.Context, ItemCategory *model.ItemCategory) (res int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ItemCategoryRepository.Create")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}
	res, err = tx.InsertWithCtx(ctx, ItemCategory)
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

func (r *ItemCategoryRepository) IsExistNameItemCategory(ctx context.Context, id int64, name string) (isExist bool) {
	ctx, span := r.opt.Trace.Start(ctx, "ItemCategoryRepository.GetByID")
	defer span.End()
	db := r.opt.Database.Read

	cond := orm.NewCondition()

	qs := db.QueryTable(new(model.ItemCategory))

	if id != 0 {
		cond = cond.AndNot("id", id)
	}

	if name != "" {
		cond = cond.And("name", name)
	}

	qs = qs.SetCond(cond)

	isExist = qs.Exist()

	return
}
