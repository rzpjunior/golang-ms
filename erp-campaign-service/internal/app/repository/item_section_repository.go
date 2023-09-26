package repository

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/edenlabs/edenlabs/timex"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/model"
)

type IItemSectionRepository interface {
	Get(ctx context.Context, req *dto.ItemSectionRequestGet) (itemSections []*model.ItemSection, count int64, err error)
	GetByID(ctx context.Context, id int64) (itemSection *model.ItemSection, err error)
	Create(ctx context.Context, itemSection *model.ItemSection) (id int64, err error)
	Update(ctx context.Context, itemSection *model.ItemSection, columns ...string) (err error)
	CheckIsIntersect(ctx context.Context, sectionType int8, startDate, finishAt string, id int64) (isExist bool, e error)
}

type ItemSectionRepository struct {
	opt opt.Options
}

func NewItemSectionRepository() IItemSectionRepository {
	return &ItemSectionRepository{
		opt: global.Setup.Common,
	}
}

func (r *ItemSectionRepository) Get(ctx context.Context, req *dto.ItemSectionRequestGet) (itemSections []*model.ItemSection, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ItemSectionRepository.Get")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.ItemSection))

	cond := orm.NewCondition()

	if req.Search != "" {
		cond1 := orm.NewCondition()
		cond1 = cond1.And("name__icontains", req.Search).Or("code__icontains", req.Search)
		cond = cond.AndCond(cond1)
	}

	if req.RegionID != "" {
		cond1 := orm.NewCondition()
		cond1 = cond1.And("regions__icontains", ","+req.RegionID+",").Or("regions__istartswith", req.RegionID+",").Or("regions__iendswith", ","+req.RegionID).Or("regions", req.RegionID)
		cond = cond.AndCond(cond1)
	}

	if req.ArchetypeID != "" {
		cond1 := orm.NewCondition()
		cond1 = cond1.And("archetypes__icontains", ","+req.ArchetypeID+",").Or("archetypes__istartswith", req.ArchetypeID+",").Or("archetypes__iendswith", ","+req.ArchetypeID).Or("archetypes", req.ArchetypeID)
		cond = cond.AndCond(cond1)
	}

	if len(req.Status) != 0 {
		cond = cond.And("status__in", req.Status)
	}

	if req.Type != 0 {
		cond = cond.And("type", req.Type)
	}

	if timex.IsValid(req.CurrentTime) {
		cond1 := orm.NewCondition()
		cond1 = cond1.And("start_at__lte", req.CurrentTime).And("finish_at__gte", req.CurrentTime)

		cond = cond.AndCond(cond1)
	}

	if req.ItemSectionID != 0 {
		cond = cond.And("id", req.ItemSectionID)
	}

	qs = qs.SetCond(cond)

	if req.OrderBy != "" {
		qs = qs.OrderBy(req.OrderBy)
	}

	_, err = qs.Offset(req.Offset).Limit(req.Limit).AllWithCtx(ctx, &itemSections)

	count, err = qs.Count()

	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *ItemSectionRepository) GetByID(ctx context.Context, id int64) (itemSection *model.ItemSection, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ItemSectionRepository.GetByID")
	defer span.End()

	itemSection = &model.ItemSection{
		ID: id,
	}
	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, itemSection, "id")
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *ItemSectionRepository) Create(ctx context.Context, itemSection *model.ItemSection) (id int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ItemSectionRepository.Create")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	id, err = tx.InsertWithCtx(ctx, itemSection)

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

func (r *ItemSectionRepository) Update(ctx context.Context, itemSection *model.ItemSection, columns ...string) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ItemSectionRepository.Update")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.UpdateWithCtx(ctx, itemSection, columns...)

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

// CheckIsIntersect : to check if there already exist active product section at based on parameters
func (r *ItemSectionRepository) CheckIsIntersect(ctx context.Context, sectionType int8, startDate, finishAt string, id int64) (isExist bool, e error) {
	ctx, span := r.opt.Trace.Start(ctx, "ItemSectionRepository.CheckIsIntersect")
	defer span.End()

	db := r.opt.Database.Read

	q := "select exists(select id from item_section where status NOT IN (2,7) and type = ? and id != ? and (" +
		"(? BETWEEN start_at and finish_at) or (? BETWEEN start_at and finish_at) or (start_at BETWEEN ? and ?) or (finish_at BETWEEN ? and ?)" +
		"))"

	if e = db.Raw(q, sectionType, id, startDate, finishAt, startDate, finishAt, startDate, finishAt).QueryRow(&isExist); e != nil {
		return isExist, e
	}

	return isExist, nil
}
