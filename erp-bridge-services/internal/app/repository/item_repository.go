package repository

import (
	"context"
	"fmt"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
)

type IItemRepository interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, uomID int64, classID int64, itemCategoryID int64) (items []*model.Item, count int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (item *model.Item, err error)
}

type ItemRepository struct {
	opt opt.Options
}

func NewItemRepository() IItemRepository {
	return &ItemRepository{
		opt: global.Setup.Common,
	}
}

func (r *ItemRepository) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, uomID int64, classID int64, itemCategoryID int64) (items []*model.Item, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ItemRepository.Get")
	defer span.End()

	count = 10
	for i := 1; i <= int(count); i++ {
		items = append(items, r.MockDatas(int64(i)))
	}
	return

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.Item))

	cond := orm.NewCondition()

	if search != "" {
		cond = cond.And("description__icontains", search)
	}

	if status != 0 {
		cond = cond.And("status", status)
	}

	if uomID != 0 {
		cond = cond.And("uom_id", uomID)
	}

	if classID != 0 {
		cond = cond.And("class_id", classID)
	}

	if itemCategoryID != 0 {
		cond = cond.And("item_category_id", itemCategoryID)
	}

	qs = qs.SetCond(cond)

	if orderBy != "" {
		qs = qs.OrderBy(orderBy)
	}

	count, err = qs.Offset(offset).Limit(limit).AllWithCtx(ctx, &items)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *ItemRepository) GetDetail(ctx context.Context, id int64, code string) (item *model.Item, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ItemRepository.GetDetail")
	defer span.End()

	// RETURN DUMMIES
	item = r.MockDatas(id)
	return

	item = &model.Item{}

	var cols []string

	if id != 0 {
		cols = append(cols, "id")
		item.ID = id
	}

	if code != "" {
		cols = append(cols, "code")
		item.Code = code
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, item, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *ItemRepository) MockDatas(id int64) (mockDatas *model.Item) {
	if id == 0 {
		id = 1
	}
	mockDatas = &model.Item{
		ID:                      id,
		Code:                    fmt.Sprintf("ITM%d", id),
		UomID:                   1,
		ClassID:                 1,
		ItemCategoryID:          1,
		Description:             fmt.Sprintf("Dummy Item %d", id),
		UnitWeightConversion:    1,
		OrderMinQty:             1,
		OrderMaxQty:             1,
		ItemType:                "Dummy ItemType",
		Packability:             "Dummy Packability",
		Capitalize:              "Dummy Capitalize",
		Note:                    "Dummy Note",
		ExcludeArchetype:        "Dummy ExcludeArchetype",
		MaxDayDeliveryDate:      1,
		FragileGoods:            "Dummy FragileGoods",
		Taxable:                 "Dummy Taxable",
		OrderChannelRestriction: "Dummy OrderChannelRestriction",
		Status:                  1,
		CreatedAt:               time.Now(),
		UpdatedAt:               time.Now(),
	}

	return mockDatas
}
