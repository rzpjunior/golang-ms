package repository

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/internal/app/model"
)

type IHelperTokenRepository interface {
	GetByHelperId(ctx context.Context, helperId string) (res *model.HelperToken, err error)
	Create(ctx context.Context, model *model.HelperToken) (err error)
	Update(ctx context.Context, token *model.HelperToken, columns ...string) (err error)
}

type HelperTokenRepository struct {
	opt opt.Options
}

func NewHelperTokenRepository() IHelperTokenRepository {
	return &HelperTokenRepository{
		opt: global.Setup.Common,
	}
}

func (r *HelperTokenRepository) GetByHelperId(ctx context.Context, helperId string) (res *model.HelperToken, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "KoliRepository.GetByID")
	defer span.End()

	res = &model.HelperToken{
		HelperIdGp: helperId,
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, res, "helper_id_gp")
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *HelperTokenRepository) Create(ctx context.Context, model *model.HelperToken) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "HelperTokenRepository.CreateUpdate")
	defer span.End()

	db := r.opt.Database.Write

	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.InsertWithCtx(ctx, model)
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

func (r *HelperTokenRepository) Update(ctx context.Context, token *model.HelperToken, columns ...string) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "HelperTokenRepository.Update")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.UpdateWithCtx(ctx, token, columns...)
	if err != nil {
		span.RecordError(err)
		err = tx.Rollback()
		if err != nil {
			span.RecordError(err)
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}
