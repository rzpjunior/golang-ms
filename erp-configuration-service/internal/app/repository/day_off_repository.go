package repository

import (
	"context"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/edenlabs/edenlabs/timex"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/internal/app/model"
)

type IDayOffRepository interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, startDate time.Time, endDate time.Time) (dayOffs []*model.DayOff, count int64, err error)
	GetByDate(ctx context.Context, startDate time.Time, endDate time.Time) (dayOffs []*model.DayOff, count int64, err error)
	GetByID(ctx context.Context, id int64) (dayOff *model.DayOff, err error)
	GetByOffDate(ctx context.Context, OffDate string) (dayOff *model.DayOff, err error)
	Create(ctx context.Context, dayOff *model.DayOff) (err error)
	Archive(ctx context.Context, dayOff *model.DayOff, columns ...string) (err error)
	UnArchive(ctx context.Context, dayOff *model.DayOff, columns ...string) (err error)
}

type DayOffRepository struct {
	opt opt.Options
}

func NewDayOffRepository() IDayOffRepository {
	return &DayOffRepository{
		opt: global.Setup.Common,
	}
}

func (r *DayOffRepository) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, startDate time.Time, endDate time.Time) (dayOffs []*model.DayOff, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "DayOffRepository.Get")
	defer span.End()

	db := r.opt.Database.Read
	qs := db.QueryTable(new(model.DayOff))

	cond := orm.NewCondition()

	if status != 0 {
		cond = cond.And("status", status)
	}

	if timex.IsValid(startDate) {
		cond = cond.And("off_date__gte", timex.ToStartTime(startDate))
	}

	if timex.IsValid(endDate) {
		cond = cond.And("off_date__lte", timex.ToLastTime(endDate))
	}

	qs = qs.SetCond(cond)

	if orderBy != "" {
		qs = qs.OrderBy(orderBy)
	}

	count, err = qs.Offset(offset).Limit(limit).AllWithCtx(ctx, &dayOffs)
	if err != nil {
		span.RecordError(err)
		return
	}

	return dayOffs, count, nil
}

func (r *DayOffRepository) GetByDate(ctx context.Context, startDate time.Time, endDate time.Time) (dayOffs []*model.DayOff, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "DayOffRepository.GetByDate")
	defer span.End()

	o := orm.NewOrm()

	o.Raw("SELECT * FROM day_off"+
		" WHERE status = 1 AND off_date between ? AND ?",
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02")).QueryRows(&dayOffs)

	count = int64(len(dayOffs))
	return dayOffs, count, nil
}

func (r *DayOffRepository) GetByID(ctx context.Context, id int64) (dayOff *model.DayOff, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "DayOffRepository.GetByID")
	defer span.End()

	dayOff = &model.DayOff{
		ID: id,
	}
	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, dayOff, "id")
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *DayOffRepository) GetByOffDate(ctx context.Context, offDate string) (dayOff *model.DayOff, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "DayOffRepository.GetByID")
	defer span.End()

	offDateTime, err := time.Parse("2006-01-02", offDate)
	if err != nil {
		span.RecordError(err)
		return
	}

	dayOff = &model.DayOff{
		OffDate: offDateTime,
		Status:  1,
	}
	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, dayOff, "off_date", "status")
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *DayOffRepository) Create(ctx context.Context, dayOff *model.DayOff) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "DayOffRepository.Create")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.InsertWithCtx(ctx, dayOff)
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

func (r *DayOffRepository) Archive(ctx context.Context, dayOff *model.DayOff, columns ...string) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "DayOffRepository.Archive")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.UpdateWithCtx(ctx, dayOff, columns...)
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

func (r *DayOffRepository) UnArchive(ctx context.Context, dayOff *model.DayOff, columns ...string) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "DayOffRepository.UnArchive")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.UpdateWithCtx(ctx, dayOff, columns...)
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
