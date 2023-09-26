package repository

import (
	"context"
	"fmt"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/edenlabs/edenlabs/timex"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/sales_service"
	"git.edenfarm.id/project-version3/erp-services/erp-sales-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-sales-service/internal/app/model"
)

type ISalesOrderRepository interface {
	GetListGRPC(ctx context.Context, req *pb.GetSalesOrderListRequest) (SalesOrderes []*model.SalesOrder, count int64, err error)
	GetDetailGRPC(ctx context.Context, req *pb.GetSalesOrderDetailRequest) (SalesOrder *model.SalesOrder, err error)
	GetListItemGRPC(ctx context.Context, req *pb.GetSalesOrderItemListRequest) (SalesOrderItems []*model.SalesOrderItem, count int64, err error)
	GetDetailItemGRPC(ctx context.Context, req *pb.GetSalesOrderItemDetailRequest) (SalesOrderItem *model.SalesOrderItem, err error)
	CreateSalesOrder(ctx context.Context, req *pb.CreateSalesOrderRequest) (SalesOrder *model.SalesOrder, err error)
	UpdateSalesOrder(ctx context.Context, req *model.SalesOrder, columns ...string) (SalesOrder *model.SalesOrder, err error)
	GetListFeedbackGRPC(ctx context.Context, req *pb.GetSalesOrderFeedbackListRequest) (SalesOrders []*model.SalesOrderFeedback, count int64, err error)
	GetListUnreviewedGRPC(ctx context.Context, so []*pb.SalesOrder) (SalesOrderFeedbacks []*model.SalesOrderFeedback, count int64, err error)
	CreateSalesOrderFeedback(ctx context.Context, req *pb.CreateSalesOrderFeedbackRequest) (SalesOrderFeedbacks *model.SalesOrderFeedback, err error)
	GetSalesOrderListCronJob(ctx context.Context, req *pb.GetSalesOrderListCronjobRequest) (salesOrders []*model.SalesOrder, err error)
	UpdateSalesOrderRemindPayment(ctx context.Context, req *pb.UpdateSalesOrderRemindPaymentRequest) (res *pb.UpdateSalesOrderRemindPaymentResponse, err error)
	CreateSalesOrderPayment(ctx context.Context, req []*model.SalesOrderPayment) (res []*model.SalesOrderPayment, err error)
}

type SalesOrderRepository struct {
	opt opt.Options
}

func NewSalesOrderRepository() ISalesOrderRepository {
	return &SalesOrderRepository{
		opt: global.Setup.Common,
	}
}

func (r *SalesOrderRepository) GetListGRPC(ctx context.Context, req *pb.GetSalesOrderListRequest) (SalesOrders []*model.SalesOrder, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SalesOrderRepository.Get")
	defer span.End()

	// dummies := r.MockDatas(10)
	// return dummies, int64(len(dummies)), nil

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.SalesOrder))

	cond := orm.NewCondition()

	if req.Search != "" {
		condGroup := orm.NewCondition()
		condGroup = condGroup.And("sales_order_number__icontains", req.Search)
		cond = cond.AndCond(condGroup)
	}

	if len(req.Status) != 0 {
		cond = cond.And("status__in", req.Status)
	}

	if req.CustomerIdGp != "" {
		cond = cond.And("customer_id_gp", req.CustomerIdGp)
	}

	if req.AddressIdGp != "" {
		cond = cond.And("address_id_gp", req.AddressIdGp)
	}

	if req.SiteIdGp != "" {
		cond = cond.And("site_id_gp", req.SiteIdGp)
	}

	if timex.IsValid(req.OrderDateFrom.AsTime()) {
		cond = cond.And("requests_delivery_date__gte", timex.ToStartTime(req.OrderDateFrom.AsTime()))
	}

	if timex.IsValid(req.OrderDateTo.AsTime()) {
		cond = cond.And("requests_delivery_date__lte", timex.ToLastTime(req.OrderDateTo.AsTime()))
	}

	if req.PaymentTermIdGp != "" {
		cond = cond.And("term_payment_sls_id_gp", req.PaymentTermIdGp)
	}

	qs = qs.SetCond(cond)

	if req.OrderBy != "" {
		qs = qs.OrderBy("-id")
	}
	qs = qs.OrderBy("-id")
	count, err = qs.Offset(req.Offset).Limit(req.Limit).AllWithCtx(ctx, &SalesOrders)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}
func (r *SalesOrderRepository) GetDetailGRPC(ctx context.Context, req *pb.GetSalesOrderDetailRequest) (SalesOrder *model.SalesOrder, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SalesOrderRepository.GetDetailGRPC")
	defer span.End()

	// dummies := r.MockDatas(10)
	// return dummies, int64(len(dummies)), nil

	var cols []string
	SalesOrder = &model.SalesOrder{}

	if req.Id != 0 {
		cols = append(cols, "id")
		SalesOrder.ID = int64(req.Id)
	}

	if req.Code != "" {
		cols = append(cols, "sales_order_number")
		SalesOrder.SalesOrderNumber = req.Code
	}

	if req.SalesOrderNumberGp != "" {
		cols = append(cols, "sales_order_number_gp")
		SalesOrder.SalesOrderNumberGP = req.SalesOrderNumberGp
	}

	if req.CustomerIdGp != "" {
		cols = append(cols, "customer_id_gp")
		SalesOrder.CustomerIDGP = req.CustomerIdGp
	}

	if req.PaymentReminder != 0 {
		cols = append(cols, "payment_reminder")
		SalesOrder.PaymentReminder = int8(req.PaymentReminder)
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, SalesOrder, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *SalesOrderRepository) GetListItemGRPC(ctx context.Context, req *pb.GetSalesOrderItemListRequest) (SalesOrderItems []*model.SalesOrderItem, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SalesOrderRepository.Get")
	defer span.End()

	// dummies := r.MockDatas(10)
	// return dummies, int64(len(dummies)), nil

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.SalesOrderItem))

	cond := orm.NewCondition()

	if req.Search != "" {
		condGroup := orm.NewCondition()
		condGroup = condGroup.And("sales_order_number__icontains", req.Search)
		cond = cond.AndCond(condGroup)
	}

	if req.ItemId != 0 {
		cond = cond.And("id", req.ItemId)
	}

	if req.SalesOrderId != 0 {
		cond = cond.And("sales_order_id", req.SalesOrderId)
	}

	qs = qs.SetCond(cond)

	if req.OrderBy != "" {
		qs = qs.OrderBy(req.OrderBy)
	}

	count, err = qs.Offset(req.Offset).Limit(req.Limit).AllWithCtx(ctx, &SalesOrderItems)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}
func (r *SalesOrderRepository) GetDetailItemGRPC(ctx context.Context, req *pb.GetSalesOrderItemDetailRequest) (SalesOrderItem *model.SalesOrderItem, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SalesOrderRepository.GetDetailItemGRPC")
	defer span.End()

	// dummies := r.MockDatas(10)
	// return dummies, int64(len(dummies)), nil

	var cols []string
	SalesOrderItem = &model.SalesOrderItem{}

	if req.Id != 0 {
		cols = append(cols, "id")
		SalesOrderItem.ID = int64(req.Id)
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, SalesOrderItem, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *SalesOrderRepository) MockDatas(total int) (mockDatas []*model.SalesOrder) {
	for i := 1; i <= total; i++ {
		mockDatas = append(mockDatas,
			&model.SalesOrder{
				ID:               int64(i),
				SalesOrderNumber: fmt.Sprintf("DummySalesOrder%d", i),
				Status:           1,
				//CreatedAt: generator.DummyTime(),
			})
	}
	return
}

func (r *SalesOrderRepository) CreateSalesOrder(ctx context.Context, req *pb.CreateSalesOrderRequest) (SalesOrder *model.SalesOrder, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SalesOrderRepository.CreateSalesOrder")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	var (
		soi []*model.SalesOrderItem
		sov []*model.SalesOrderVoucher
	)
	SalesOrder = &model.SalesOrder{
		ID:                  req.Data.Id,
		SalesOrderNumber:    req.Data.SalesOrderNumber,
		AddressIDGP:         req.Data.AddressIdGp,
		CustomerIDGP:        req.Data.CustomerIdGp,
		WrtIDGP:             req.Data.WrtIdGp,
		TermPaymentSlsIDGP:  req.Data.TermPaymentSlsIdGp,
		SiteIDGP:            req.Data.SiteIdGp,
		SubDistrictIDGP:     req.Data.SubDistrictIdGp,
		RegionIDGP:          req.Data.RegionIdGp,
		PaymentGroupSlsID:   int32(req.Data.PaymentGroupSlsId),
		ArchetypeIDGP:       req.Data.ArchetypeIdGp,
		RecognitionDate:     time.Now(),
		RequestsShipDate:    req.Data.RequestsShipDate.AsTime(),
		BillingAddress:      req.Data.BillingAddress,
		ShippingAddress:     req.Data.ShippingAddress,
		DeliveryFee:         req.Data.DeliveryFee,
		VouDiscAmount:       req.Data.VouDiscAmount,
		CustomerPointLogID:  req.Data.CustomerPointLogId,
		TotalPrice:          req.Data.TotalPrice,
		TotalCharge:         req.Data.TotalCharge,
		TotalWeight:         req.Data.TotalWeight,
		Note:                req.Data.Note,
		ShippingAddressNote: req.Data.ShippingAddressNote,
		Status:              1,
		CreatedAt:           time.Now(),
		CreatedBy:           1,
		EdenPointCampaignID: req.Data.EdenPointCampaignId,
		IntegrationCode:     req.Data.IntegrationCode,
		PaymentReminder:     2,
		CancelType:          int8(req.Data.CancelType),
		PriceLevelIDGP:      req.Data.PriceLevelIdGp,
		ShippingMethodIDGP:  req.Data.ShippingMethodIdGp,
		CustomerNameGP:      req.Data.CustomerNameGp,
	}

	SalesOrder.ID, err = tx.InsertWithCtx(ctx, SalesOrder)
	if err != nil {
		span.RecordError(err)
		tx.Rollback()
		return
	}
	if req.Data.SalesOrderItem != nil {
		for _, v := range req.Data.SalesOrderItem {
			soi = append(soi, &model.SalesOrderItem{
				SalesOrderID:     SalesOrder.ID,
				ItemIDGP:         v.ItemIdGp,
				OrderQty:         v.OrderQty,
				UnitPrice:        v.UnitPrice,
				Subtotal:         v.Subtotal,
				Weight:           v.Weight,
				UomIDGP:          v.UomIdGp,
				PriceTieringIDGP: v.PriceTieringIdGp,
			})
		}
		_, err = tx.InsertMultiWithCtx(ctx, 100, soi)
		if err != nil {
			span.RecordError(err)
			tx.Rollback()
			return
		}
	}

	if req.Data.SalesOrderVoucher != nil {
		for _, v := range req.Data.SalesOrderVoucher {
			sov = append(sov, &model.SalesOrderVoucher{
				SalesOrderID: SalesOrder.ID,
				VoucherIDGP:  v.VoucherIdGp,
				DiscAmount:   v.DiscAmount,
				CreatedAt:    time.Now(),
				VoucherType:  int8(v.VoucherType),
			})
		}

		_, err = tx.InsertMultiWithCtx(ctx, 100, sov)
		if err != nil {
			span.RecordError(err)
			tx.Rollback()
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

func (r *SalesOrderRepository) UpdateSalesOrder(ctx context.Context, req *model.SalesOrder, columns ...string) (SalesOrder *model.SalesOrder, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SalesOrderRepository.UpdateSalesOrderHeader")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.UpdateWithCtx(ctx, req, columns...)
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

func (r *SalesOrderRepository) GetListFeedbackGRPC(ctx context.Context, req *pb.GetSalesOrderFeedbackListRequest) (SalesOrders []*model.SalesOrderFeedback, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SalesOrderRepository.GetListFeedbackGRPC")
	defer span.End()

	// dummies := r.MockDatas(10)
	// return dummies, int64(len(dummies)), nil

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.SalesOrderFeedback))

	cond := orm.NewCondition()

	if req.CustomerId != 0 {
		cond = cond.And("customer_id", req.CustomerId)
	}

	if req.SalesOrderId != 0 {
		cond = cond.And("sales_order_id", req.SalesOrderId)
	}

	if req.SalesOrderCode != "" {
		// cond = cond.And("site_id", req.SiteId)
	}

	// if req.SalespersonId != 0 {
	// 	cond = cond.And("sales_person_id", req.SalespersonId)
	// }

	// if timex.IsValid(req.OrderDateFrom.AsTime()) {
	// 	cond = cond.And("requests_delivery_date__gte", timex.ToStartTime(req.OrderDateFrom.AsTime()))
	// }
	// if timex.IsValid(req.OrderDateTo.AsTime()) {
	// 	cond = cond.And("requests_delivery_date__lte", timex.ToLastTime(req.OrderDateTo.AsTime()))
	// }

	qs = qs.SetCond(cond)

	qs = qs.OrderBy("-id")
	_, err = qs.AllWithCtx(ctx, &SalesOrders)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *SalesOrderRepository) GetListUnreviewedGRPC(ctx context.Context, so []*pb.SalesOrder) (SalesOrderFeedbacks []*model.SalesOrderFeedback, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SalesOrderRepository.GetListUnreviewedGRPC")
	defer span.End()

	db := r.opt.Database.Read
	var tempSalesOrderFeedback *model.SalesOrderFeedback
	// var qMark string
	// for range so {
	// 	qMark += "?,"
	// }
	// qMark = qMark[:len(qMark)-1]

	for _, v := range so {
		e := db.Raw("Select * from sales_order_feedback sof where sales_order_id=?", v.Id).QueryRow(&tempSalesOrderFeedback)
		if e != nil {
			layout := "2006-01-02"
			SalesOrderFeedbacks = append(SalesOrderFeedbacks, &model.SalesOrderFeedback{
				SalesOrderCode: v.SalesOrderNumber,
				DeliveryDate:   v.RequestsShipDate.AsTime().Format(layout),
				TotalCharge:    v.TotalCharge,
				SalesOrder:     v.Id,
			})
		}
	}

	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

// func (r *SalesOrderRepository) GetDetailFeedbackGRPC(ctx context.Context, req *pb.GetSalesOrderDetailRequest) (SalesOrder *model.SalesOrder, err error) {
// 	ctx, span := r.opt.Trace.Start(ctx, "SalesOrderRepository.GetDetailGRPC")
// 	defer span.End()

// 	// dummies := r.MockDatas(10)
// 	// return dummies, int64(len(dummies)), nil

// 	var cols []string
// 	SalesOrder = &model.SalesOrder{}

// 	if req.Id != 0 {
// 		cols = append(cols, "id")
// 		SalesOrder.ID = int64(req.Id)
// 	}

// 	if req.Code != "" {
// 		cols = append(cols, "sales_order_number")
// 		SalesOrder.SalesOrderNumber = req.Code
// 	}
// 	db := r.opt.Database.Read
// 	err = db.ReadWithCtx(ctx, SalesOrder, cols...)
// 	if err != nil {
// 		span.RecordError(err)
// 		return
// 	}

// 	return
// }

func (r *SalesOrderRepository) CreateSalesOrderFeedback(ctx context.Context, req *pb.CreateSalesOrderFeedbackRequest) (SalesOrderFeedbacks *model.SalesOrderFeedback, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SalesOrderRepository.GetListUnreviewedGRPC")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}
	sof := &model.SalesOrderFeedback{
		SalesOrderCode: req.Data.SalesOrderCode,
		DeliveryDate:   req.Data.DeliveryDate,
		RatingScore:    int8(req.Data.RatingScore),
		Description:    req.Data.Description,
		ToBeContacted:  int8(req.Data.ToBeContacted), //to be contacted
		CreatedAt:      time.Now(),
		SalesOrder:     req.Data.SalesOrderId, //sales order id
		Customer:       req.Data.CustomerId,   //customer id
		Tags:           req.Data.Tags,         //tags string
	}

	_, err = tx.InsertWithCtx(ctx, sof)
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

func (r *SalesOrderRepository) GetSalesOrderListCronJob(ctx context.Context, req *pb.GetSalesOrderListCronjobRequest) (salesOrders []*model.SalesOrder, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "GetSalesOrderListCronJob.Get")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.SalesOrder))

	cond := orm.NewCondition()
	if len(req.Status) != 0 {
		cond = cond.And("status__in", req.Status)
	}

	if req.RegionIdGp != "" {
		cond = cond.And("region_id_gp", req.RegionIdGp)
	}

	if req.RequestsDeliveryDate != "" {
		cond = cond.And("requests_delivery_date", req.RequestsDeliveryDate)
	}

	if req.PaymentReminder != 0 {
		cond = cond.And("payment_reminder", req.PaymentReminder)
	}

	qs = qs.SetCond(cond)

	qs = qs.OrderBy("-id")
	_, err = qs.AllWithCtx(ctx, &salesOrders)

	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *SalesOrderRepository) UpdateSalesOrderRemindPayment(ctx context.Context, req *pb.UpdateSalesOrderRemindPaymentRequest) (res *pb.UpdateSalesOrderRemindPaymentResponse, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SalesOrderRepository.UpdateSalesOrderRemindPayment")
	defer span.End()

	updateParam := make(orm.Params, 1)
	updateParam["payment_reminder"] = 1

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)

	if err != nil {
		return
	}
	_, err = tx.QueryTable("sales_order").Filter("id__in", req.SalesOrderId).UpdateWithCtx(ctx, updateParam)

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

func (r *SalesOrderRepository) CreateSalesOrderPayment(ctx context.Context, req []*model.SalesOrderPayment) (res []*model.SalesOrderPayment, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SalesOrderRepository.CreateSalesOrderPayment")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.InsertMultiWithCtx(ctx, 100, req)
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
