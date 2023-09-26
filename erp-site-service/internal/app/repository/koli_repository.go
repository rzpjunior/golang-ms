package repository

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/internal/app/model"
)

type IKoliRepository interface {
	Get(ctx context.Context, req *dto.KoliGetRequest) (kolis []*model.Koli, count int64, err error)
	GetByID(ctx context.Context, id int64) (koli *model.Koli, err error)
}

type KoliRepository struct {
	opt opt.Options
}

func NewKoliRepository() IKoliRepository {
	return &KoliRepository{
		opt: global.Setup.Common,
	}
}

func (r *KoliRepository) Get(ctx context.Context, req *dto.KoliGetRequest) (kolis []*model.Koli, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "KoliRepository.Get")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.Koli))

	cond := orm.NewCondition()

	if req.Status != 0 {
		cond = cond.And("status", req.Status)
	}

	qs = qs.SetCond(cond)

	if req.OrderBy != "" {
		qs = qs.OrderBy(req.OrderBy)
	}

	count, err = qs.Offset(req.Offset).Limit(req.Limit).AllWithCtx(ctx, &kolis)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *KoliRepository) GetByID(ctx context.Context, id int64) (koli *model.Koli, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "KoliRepository.GetByID")
	defer span.End()

	koli = &model.Koli{
		Id: id,
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, koli, "id")
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}
