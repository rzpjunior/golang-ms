package repository

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/internal/app/model"
)

type IDeliveryRunReturnItemRepository interface {
	Get(ctx context.Context, req dto.DeliveryRunReturnItemGetRequest) (deliveryRunReturnItems []*model.DeliveryRunReturnItem, count int64, err error)
	GetByID(ctx context.Context, id int64, deliveryOrderItemId string) (res *model.DeliveryRunReturnItem, err error)
	Create(ctx context.Context, model *model.DeliveryRunReturnItem) (err error)
	Update(ctx context.Context, model *model.DeliveryRunReturnItem, columns ...string) (err error)
	Delete(ctx context.Context, id int64) (err error)
}

type DeliveryRunReturnItemRepository struct {
	opt opt.Options
}

func NewDeliveryRunReturnItemRepository() IDeliveryRunReturnItemRepository {
	return &DeliveryRunReturnItemRepository{
		opt: global.Setup.Common,
	}
}

func (r *DeliveryRunReturnItemRepository) Get(ctx context.Context, req dto.DeliveryRunReturnItemGetRequest) (deliveryRunReturnItems []*model.DeliveryRunReturnItem, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "DeliveryRunReturnItemRepository.Get")
	defer span.End()

	// // RETURN DUMMY
	// dummies := []*model.DeliveryRunReturnItem{
	// 	{
	// 		ID:                  1,
	// 		ReceiveQty:          999,
	// 		ReturnReason:        1,
	// 		ReturnEvidence:      "DUMMY EVIDENCE",
	// 		Subtotal:            9999,
	// 		DeliveryRunReturnID: 1,
	// 		DeliveryOrderItemID: "DOI0001",
	// 	},
	// 	{
	// 		ID:                  2,
	// 		ReceiveQty:          999,
	// 		ReturnReason:        2,
	// 		ReturnEvidence:      "dummy evidence",
	// 		Subtotal:            9999,
	// 		DeliveryRunReturnID: 1,
	// 		DeliveryOrderItemID: "DOI0002",
	// 	},
	// }
	// return dummies, 2, nil

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.DeliveryRunReturnItem))

	cond := orm.NewCondition()

	if len(req.ArrDeliveryRunReturnIDs) > 0 {
		cond = cond.And("delivery_run_return_id__in", req.ArrDeliveryRunReturnIDs)
	}

	if len(req.ArrDeliveryOrderItemIDs) > 0 {
		cond = cond.And("delivery_order_item_id__in", req.ArrDeliveryOrderItemIDs)
	}

	qs = qs.SetCond(cond)

	if req.OrderBy != "" {
		qs = qs.OrderBy(req.OrderBy)
	}

	count, err = qs.Offset(req.Offset).Limit(req.Limit).AllWithCtx(ctx, &deliveryRunReturnItems)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *DeliveryRunReturnItemRepository) GetByID(ctx context.Context, id int64, deliveryOrderItemId string) (res *model.DeliveryRunReturnItem, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "DeliveryRunReturnItemRepository.GetByID")
	defer span.End()

	// RETURN DUMMY
	// dummy := &model.DeliveryRunReturnItem{
	// 	ID:                  1,
	// 	ReceiveQty:          999,
	// 	ReturnReason:        1,
	// 	ReturnEvidence:      "DUMMY EVIDENCE",
	// 	Subtotal:            9999,
	// 	DeliveryRunReturnID: 1,
	// 	DeliveryOrderItemID: "DOI0001",
	// }
	// return dummy, nil

	res = &model.DeliveryRunReturnItem{}

	var cols []string

	if id != 0 {
		cols = append(cols, "id")
		res.ID = id
	}

	if deliveryOrderItemId != "" {
		cols = append(cols, "delivery_order_item_id")
		res.DeliveryOrderItemID = deliveryOrderItemId
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, res, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *DeliveryRunReturnItemRepository) Create(ctx context.Context, model *model.DeliveryRunReturnItem) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "DeliveryRunReturnItemRepository.Create")
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

func (r *DeliveryRunReturnItemRepository) Update(ctx context.Context, model *model.DeliveryRunReturnItem, columns ...string) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "DeliveryRunReturnItemRepository.Update")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.UpdateWithCtx(ctx, model, columns...)
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

func (r *DeliveryRunReturnItemRepository) Delete(ctx context.Context, id int64) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "DeliveryRunReturnItemRepository.Delete")
	defer span.End()

	// RETURN DUMMY
	// return nil

	db := r.opt.Database.Write

	qs := db.QueryTable(new(model.DeliveryRunReturnItem))

	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = qs.Filter("id", id).Delete()
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
