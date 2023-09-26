package repository

import (
	"context"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/edenlabs/edenlabs/timex"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/model"
)

type IBannerRepository interface {
	Get(ctx context.Context, req *dto.BannerRequestGet) (banners []*model.Banner, count int64, err error)
	GetByID(ctx context.Context, id int64) (banner *model.Banner, err error)
	Create(ctx context.Context, banner *model.Banner) (err error)
	Archive(ctx context.Context, id int64, req dto.BannerRequestArchive) (err error)
	Update(ctx context.Context, banner *model.Banner, columns ...string) (err error)
}

type BannerRepository struct {
	opt opt.Options
}

func NewBannerRepository() IBannerRepository {
	return &BannerRepository{
		opt: global.Setup.Common,
	}
}

func (r *BannerRepository) Get(ctx context.Context, req *dto.BannerRequestGet) (banners []*model.Banner, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "BannerRepository.Get")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.Banner))

	cond := orm.NewCondition()

	if req.Search != "" {
		cond1 := orm.NewCondition()
		cond1 = cond1.Or("name__icontains", req.Search).Or("code__icontains", req.Search)
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

	if timex.IsValid(req.CurrentTime) {
		cond1 := orm.NewCondition()
		cond1 = cond1.And("start_at__lte", req.CurrentTime).And("finish_at__gte", req.CurrentTime)

		cond = cond.AndCond(cond1)
	}

	qs = qs.SetCond(cond)

	if req.OrderBy != "" {
		qs = qs.OrderBy(req.OrderBy)
	}

	_, err = qs.Offset(req.Offset).Limit(req.Limit).AllWithCtx(ctx, &banners)
	if err != nil {
		span.RecordError(err)
		return
	}
	count, err = qs.Count()

	return
}

func (r *BannerRepository) GetByID(ctx context.Context, id int64) (banner *model.Banner, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "BannerRepository.GetByID")
	defer span.End()

	banner = &model.Banner{
		ID: id,
	}
	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, banner, "id")
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *BannerRepository) Create(ctx context.Context, banner *model.Banner) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "BannerRepository.Create")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.InsertWithCtx(ctx, banner)
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

func (r *BannerRepository) Update(ctx context.Context, banner *model.Banner, columns ...string) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "BannerRepository.Update")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.UpdateWithCtx(ctx, banner, columns...)

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

func (r *BannerRepository) Archive(ctx context.Context, id int64, req dto.BannerRequestArchive) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "BannerRepository.Archive")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	banner := &model.Banner{
		ID:        id,
		Status:    statusx.ConvertStatusName("Archived"),
		UpdatedAt: time.Now(),
		Note:      req.Note,
	}

	_, err = tx.UpdateWithCtx(ctx, banner, "Status", "Note", "UpdatedAt")
	if err != nil {
		span.RecordError(err)
		tx.Rollback()
		return
	}
	err = tx.Commit()

	return
}
