package repository

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/internal/app/model"
)

type IConsolidatedShipmentSignatureRepository interface {
	GetByConsolidatedShipmentID(ctx context.Context, consolidatedShipmentID int64) (consolidatedShipmentSignatures []*model.ConsolidatedShipmentSignature, count int64, err error)
	GetByID(ctx context.Context, id int64) (consolidatedShipmentSignature *model.ConsolidatedShipmentSignature, err error)
	Create(ctx context.Context, consolidatedShipmentSignature *model.ConsolidatedShipmentSignature) (err error)
	CheckAlreadySigned(ctx context.Context, consolidatedShipmentID int64, jobFunction string) (isAlready bool, err error)
}

type ConsolidatedShipmentSignatureRepository struct {
	opt opt.Options
}

func NewConsolidatedShipmentSignatureRepository() IConsolidatedShipmentSignatureRepository {
	return &ConsolidatedShipmentSignatureRepository{
		opt: global.Setup.Common,
	}
}

func (r *ConsolidatedShipmentSignatureRepository) GetByConsolidatedShipmentID(ctx context.Context, consolidatedShipmentID int64) (consolidatedShipmentSignatures []*model.ConsolidatedShipmentSignature, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ConsolidatedShipmentSignatureRepository.GetByConsolidatedShipmentID")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.ConsolidatedShipmentSignature))

	cond := orm.NewCondition()

	if consolidatedShipmentID != 0 {
		cond = cond.And("consolidated_shipment_id", consolidatedShipmentID)
	}

	qs = qs.SetCond(cond)

	count, err = qs.AllWithCtx(ctx, &consolidatedShipmentSignatures)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *ConsolidatedShipmentSignatureRepository) GetByID(ctx context.Context, id int64) (consolidatedShipmentSignature *model.ConsolidatedShipmentSignature, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ConsolidatedShipmentSignatureRepository.GetByID")
	defer span.End()

	consolidatedShipmentSignature = &model.ConsolidatedShipmentSignature{
		ID: id,
	}
	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, consolidatedShipmentSignature, "id")
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *ConsolidatedShipmentSignatureRepository) Create(ctx context.Context, consolidatedShipmentSignature *model.ConsolidatedShipmentSignature) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ConsolidatedShipmentSignatureRepository.Create")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.InsertWithCtx(ctx, consolidatedShipmentSignature)
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

func (r *ConsolidatedShipmentSignatureRepository) CheckAlreadySigned(ctx context.Context, consolidatedShipmentID int64, jobFunction string) (isAlready bool, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ConsolidatedShipmentSignatureRepository.Create")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.ConsolidatedShipmentSignature))

	cond := orm.NewCondition()

	qs = qs.SetCond(cond)

	var count int64
	count, err = qs.Filter("consolidated_shipment_id", consolidatedShipmentID).Filter("job_function", jobFunction).CountWithCtx(ctx)
	if err != nil || count != 0 {
		span.RecordError(err)
		isAlready = false
		return
	}

	isAlready = true
	return
}
