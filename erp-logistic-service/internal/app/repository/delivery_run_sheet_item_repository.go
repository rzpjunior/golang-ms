package repository

import (
	"context"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/internal/app/model"
)

type IDeliveryRunSheetItemRepository interface {
	Get(ctx context.Context, req dto.DeliveryRunSheetItemGetRequest) (deliveryRunSheetItems []*model.DeliveryRunSheetItem, count int64, err error)
	GetByID(ctx context.Context, id int64) (deliveryRunSheetItem *model.DeliveryRunSheetItem, err error)
	Create(ctx context.Context, model *model.DeliveryRunSheetItem) (err error)
	Update(ctx context.Context, deliveryRunSheetItem *model.DeliveryRunSheetItem, columns ...string) (err error)
	GetAllGroupedDeliveryRunSheetItem(ctx context.Context, req dto.GroupedDeliveryRunSheetItemGetRequest) (deliveryRunSheetItems []*model.DeliveryRunSheetItem, count int64, err error)
}

type DeliveryRunSheetItemRepository struct {
	opt opt.Options
}

func NewDeliveryRunSheetItemRepository() IDeliveryRunSheetItemRepository {
	return &DeliveryRunSheetItemRepository{
		opt: global.Setup.Common,
	}
}

func (r *DeliveryRunSheetItemRepository) Get(ctx context.Context, req dto.DeliveryRunSheetItemGetRequest) (deliveryRunSheetItems []*model.DeliveryRunSheetItem, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "DeliveryRunSheetItemRepository.Get")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.DeliveryRunSheetItem))

	cond := orm.NewCondition()
	if len(req.StepType) > 0 {
		cond = cond.And("step_type__in", req.StepType)
	}

	if len(req.Status) > 0 {
		cond = cond.And("status__in", req.Status)
	}

	if len(req.DeliveryRunSheetIDs) > 0 {
		cond = cond.And("delivery_run_sheet_id__in", req.DeliveryRunSheetIDs)
	}

	if len(req.CourierIDs) > 0 {
		cond = cond.And("courier_id__in", req.CourierIDs)
	}

	if len(req.ArrSalesOrderIDs) > 0 {
		cond = cond.And("sales_order_id__in", req.ArrSalesOrderIDs)
	}

	if req.SearchSalesOrderCode != "" {
		cond = cond.And("sales_order_id__icontains", req.SearchSalesOrderCode)
	}

	qs = qs.SetCond(cond)

	if req.OrderBy != "" {
		qs = qs.OrderBy(req.OrderBy)
	}

	if req.GroupBy != "" {
		qs = qs.OrderBy(req.GroupBy)
	}

	count, err = qs.Offset(req.Offset).Limit(req.Limit).AllWithCtx(ctx, &deliveryRunSheetItems)

	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *DeliveryRunSheetItemRepository) GetByID(ctx context.Context, id int64) (deliveryRunSheetItem *model.DeliveryRunSheetItem, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "DeliveryRunSheetItemRepository.GetByID")
	defer span.End()

	deliveryRunSheetItem = &model.DeliveryRunSheetItem{
		ID: id,
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, deliveryRunSheetItem, "id")
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *DeliveryRunSheetItemRepository) Create(ctx context.Context, model *model.DeliveryRunSheetItem) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "DeliveryRunSheetItemRepository.Create")
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

func (r *DeliveryRunSheetItemRepository) Update(ctx context.Context, model *model.DeliveryRunSheetItem, columns ...string) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "DeliveryRunSheetItemRepository.Update")
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

func (r *DeliveryRunSheetItemRepository) GetAllGroupedDeliveryRunSheetItem(ctx context.Context, req dto.GroupedDeliveryRunSheetItemGetRequest) (deliveryRunSheetItems []*model.DeliveryRunSheetItem, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ControlTowerRepository.GetAllGroupedDeliveryRunSheet")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.DeliveryRunSheetItem))

	cond := orm.NewCondition()

	if len(req.ArrSalesOrderIDs) > 0 {
		cond = cond.And("sales_order_id__in", req.ArrSalesOrderIDs)
	}

	if len(req.ArrCourierVendorsCourierIDs) > 0 {
		cond = cond.And("courier_id__in", req.ArrCourierVendorsCourierIDs)
	}

	if req.CourierID != "" {
		cond = cond.And("courier_id", req.CourierID)
	}

	if len(req.Status) > 0 {
		cond = cond.And("status__in", req.Status)
	}

	if req.SearchSalesOrderID != "" {
		cond = cond.And("sales_order_id", req.SearchSalesOrderID)
	}

	qs = qs.SetCond(cond)

	if req.OrderBy != "" {
		qs = qs.OrderBy(req.OrderBy)
	}

	qs = qs.GroupBy("delivery_run_sheet_id")

	count, err = qs.Offset(req.Offset).Limit(req.Limit).AllWithCtx(ctx, &deliveryRunSheetItems)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *DeliveryRunSheetItemRepository) MockDatas() (mockDatas []*model.DeliveryRunSheetItem) {
	sampleLatitude := -6.1879324
	sampleLongitude := 106.7376164
	mockDatas = append(mockDatas,
		&model.DeliveryRunSheetItem{
			ID:                          1,
			StepType:                    2,
			Latitude:                    &sampleLatitude,
			Longitude:                   &sampleLongitude,
			Status:                      1,
			Note:                        "DUMMY NOTE DRSI",
			RecipientName:               "DUMMY NAME",
			MoneyReceived:               99999,
			DeliveryEvidenceImageURL:    "DUMMY DELIVERY EVIDENCE",
			TransactionEvidenceImageURL: "DUMMY TRANSACTION EVIDENCE",
			ArrivalTime:                 time.Now(),
			UnpunctualReason:            0,
			UnpunctualDetail:            0,
			FarDeliveryReason:           "",
			CreatedAt:                   time.Now(),
			StartedAt:                   time.Now(),
			FinishedAt:                  time.Now().Add(10),
			DeliveryRunSheetID:          1,
			CourierID:                   "COU0001",
			SalesOrderID:                "SO0001",
		},
		&model.DeliveryRunSheetItem{
			ID:                          2,
			StepType:                    2,
			Latitude:                    &sampleLatitude,
			Longitude:                   &sampleLongitude,
			Status:                      1,
			Note:                        "dummy note drsi",
			RecipientName:               "dummy drsi",
			MoneyReceived:               99999,
			DeliveryEvidenceImageURL:    "dummy delivery evidence",
			TransactionEvidenceImageURL: "dummy transaction evidence",
			ArrivalTime:                 time.Now(),
			UnpunctualReason:            0,
			UnpunctualDetail:            0,
			FarDeliveryReason:           "dummy far delivery reason",
			CreatedAt:                   time.Now(),
			StartedAt:                   time.Now(),
			FinishedAt:                  time.Now().Add(10),
			DeliveryRunSheetID:          2,
			CourierID:                   "COU0001",
			SalesOrderID:                "SO0002",
		},
	)

	return mockDatas
}
