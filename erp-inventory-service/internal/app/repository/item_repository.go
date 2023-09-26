package repository

import (
	"context"
	"strings"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/project-version3/erp-services/erp-inventory-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-inventory-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-inventory-service/internal/app/model"
)

type IItemRepository interface {
	GetDetail(ctx context.Context, id int64, code string, itemCategoryID int64, archetypeID string, orderchannel int64) (item *model.Item, err error)
	Update(ctx context.Context, Item *model.Item, columns ...string) (err error)
	SyncGP(ctx context.Context, Item *model.Item) (err error)
	Get(ctx context.Context, req *dto.ItemRequestGet) (items []*model.Item, count int64, err error)
}

type ItemRepository struct {
	opt opt.Options
}

func NewItemRepository() IItemRepository {
	return &ItemRepository{
		opt: global.Setup.Common,
	}
}

func (r *ItemRepository) GetDetail(ctx context.Context, id int64, code string, itemCategoryID int64, archetypeID string, orderchannel int64) (item *model.Item, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ItemRepository.GetDetail")
	defer span.End()
	var itemDetail []*model.Item

	db := r.opt.Database.Read

	cond := orm.NewCondition()

	qs := db.QueryTable(new(model.Item))

	if itemCategoryID != 0 {
		cond = cond.And("item_category_id", itemCategoryID)
	}

	if archetypeID != "" {
		cond1 := orm.NewCondition()
		cond1 = cond1.AndNot("exclude_archetype__icontains", archetypeID).Or("exclude_archetype__isnull", true).Or("exclude_archetype", "")
		cond = cond.AndCond(cond1)
	}

	if orderchannel != 0 {
		cond1 := orm.NewCondition()
		cond1 = cond1.AndNot("order_channel_restriction__icontains", orderchannel).Or("order_channel_restriction__isnull", true).Or("order_channel_restriction", "")
		cond = cond.AndCond(cond1)
	}
	if code != "" {
		cond = cond.And("item_id_gp", code)
	}
	if id != 0 {
		cond = cond.And("id", id)
	}

	qs = qs.SetCond(cond)

	err = qs.OneWithCtx(ctx, &itemDetail)
	if err != nil {
		span.RecordError(err)
		return
	}

	item = &model.Item{
		ID:                      itemDetail[0].ID,
		ItemIDGP:                itemDetail[0].ItemIDGP,
		ItemCategoryID:          itemDetail[0].ItemCategoryID,
		Note:                    itemDetail[0].Note,
		ExcludeArchetype:        itemDetail[0].ExcludeArchetype,
		MaxDayDeliveryDate:      itemDetail[0].MaxDayDeliveryDate,
		OrderChannelRestriction: itemDetail[0].OrderChannelRestriction,
		Packability:             itemDetail[0].Packability,
		FragileGoods:            itemDetail[0].FragileGoods,
	}
	return
}

func (r *ItemRepository) Update(ctx context.Context, Item *model.Item, columns ...string) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ItemRepository.Update")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.UpdateWithCtx(ctx, Item, columns...)

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

func (r *ItemRepository) SyncGP(ctx context.Context, item *model.Item) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ItemRepository.SyncGP")
	defer span.End()

	db := r.opt.Database.Write

	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, _, err = tx.ReadOrCreateWithCtx(ctx, item, "item_id_gp", []string{}...)
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

func (r *ItemRepository) Get(ctx context.Context, req *dto.ItemRequestGet) (items []*model.Item, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ItemRepository.Get")
	defer span.End()

	db := r.opt.Database.Read

	cond := orm.NewCondition()

	qs := db.QueryTable(new(model.Item))

	if req.Search != "" {
		cond = cond.And("name__icontains", req.Search)
	}

	if req.ID != "" {
		idArr := strings.Split(req.ID, ",")
		cond = cond.And("id__in", idArr)
	}

	if req.ItemCategoryID != 0 {
		cond = cond.And("item_category_id", req.ItemCategoryID)
	}

	if req.ArchetypeIDGP != "" {
		cond1 := orm.NewCondition()
		cond1 = cond1.AndNot("exclude_archetype__icontains", req.ArchetypeIDGP).Or("exclude_archetype__isnull", true).Or("exclude_archetype", "")
		cond = cond.AndCond(cond1)
	}

	if req.OrderChannel != 0 {
		cond1 := orm.NewCondition()
		cond1 = cond1.AndNot("order_channel_restriction__icontains", req.OrderChannel).Or("order_channel_restriction__isnull", true).Or("order_channel_restriction", "")
		cond = cond.AndCond(cond1)
	}

	if len(req.ItemIdGP) != 0 {
		cond = cond.And("item_id_gp__in", req.ItemIdGP)
	}

	qs = qs.SetCond(cond)

	if req.OrderBy != "" {
		qs = qs.OrderBy(req.OrderBy)
	} else {
		qs = qs.OrderBy("-item_id_gp")
	}

	_, err = qs.Offset(req.Offset).Limit(req.Limit).AllWithCtx(ctx, &items)
	if err != nil {
		span.RecordError(err)
		return
	}

	if err != nil {
		span.RecordError(err)
		return
	}

	return
}
