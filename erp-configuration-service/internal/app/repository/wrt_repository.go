package repository

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/internal/app/model"
)

type IWrtRepository interface {
	Get(ctx context.Context, offset, limit, typeWrt int, regionID int64, search string) (Wrtes []*model.Wrt, count int64, err error)
	GetDetail(ctx context.Context, id int64, wrtGpID string) (Wrt *model.Wrt, err error)
	Update(ctx context.Context, Wrt *model.Wrt, columns ...string) (err error)
	SyncGP(ctx context.Context, Wrt *model.Wrt) (err error)
}

type WrtRepository struct {
	opt opt.Options
}

func NewWrtRepository() IWrtRepository {
	return &WrtRepository{
		opt: global.Setup.Common,
	}
}

func (r *WrtRepository) Get(ctx context.Context, offset, limit, wrtType int, regionID int64, search string) (Wrtes []*model.Wrt, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "WrtRepository.Get")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.Wrt))

	cond := orm.NewCondition()

	if regionID != 0 {
		cond = cond.And("region_id", regionID)
	}

	if wrtType != 0 {
		cond = cond.And("type", wrtType)
	}

	if search != "" {
		condGroup := orm.NewCondition()
		condGroup = condGroup.And("code__icontains", search)
		cond = cond.AndCond(condGroup)
	}

	qs = qs.SetCond(cond)

	count, err = qs.Offset(offset).Limit(limit).AllWithCtx(ctx, &Wrtes)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *WrtRepository) GetDetail(ctx context.Context, id int64, wrtGpID string) (Wrt *model.Wrt, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "WrtRepository.GetByCode")
	defer span.End()

	var cols []string

	Wrt = &model.Wrt{
		ID:    id,
		WrtID: wrtGpID,
	}

	if id != 0 {
		cols = append(cols, "id")
	}

	if wrtGpID != "" {
		cols = append(cols, "WrtId")
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, Wrt, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *WrtRepository) Update(ctx context.Context, Wrt *model.Wrt, columns ...string) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "WrtRepository.Update")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.UpdateWithCtx(ctx, Wrt, columns...)

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

func (r *WrtRepository) SyncGP(ctx context.Context, Wrt *model.Wrt) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ItemRepository.SyncGP")
	defer span.End()

	db := r.opt.Database.Write

	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, _, err = tx.ReadOrCreateWithCtx(ctx, Wrt, "WrtID", []string{}...)
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
