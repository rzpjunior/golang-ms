package repository

import (
	"context"
	"fmt"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
)

type ISalesPaymentRepository interface {
	Get(ctx context.Context, salesInvoiceID int64) (salesOrders []*model.SalesPayment, count int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (salesOrder *model.SalesPayment, err error)
}

type SalesPaymentRepository struct {
	opt opt.Options
}

func NewSalesPaymentRepository() ISalesPaymentRepository {
	return &SalesPaymentRepository{
		opt: global.Setup.Common,
	}
}

func (r *SalesPaymentRepository) Get(ctx context.Context, salesInvoiceID int64) (salesOrders []*model.SalesPayment, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SalesPaymentRepository.Get")
	defer span.End()

	// RETURN DUMMIES
	dummies := r.MockDatas(10)
	return dummies, int64(len(dummies)), nil

	return
}

func (r *SalesPaymentRepository) GetDetail(ctx context.Context, id int64, code string) (salesOrder *model.SalesPayment, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SalesPaymentRepository.GetDetail")
	defer span.End()

	// RETURN DUMMIES
	return r.MockDatas(1)[0], nil

	salesOrder = &model.SalesPayment{}

	var cols []string

	if id != 0 {
		cols = append(cols, "id")
		salesOrder.ID = id
	}

	if code != "" {
		cols = append(cols, "code")
		salesOrder.Code = code
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, salesOrder, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *SalesPaymentRepository) MockDatas(total int) (mockDatas []*model.SalesPayment) {
	for i := 1; i <= total; i++ {
		mockDatas = append(mockDatas,
			&model.SalesPayment{
				ID:              int64(i),
				Code:            fmt.Sprintf("dummy%d", 1),
				Status:          statusx.ConvertStatusName(statusx.Finished),
				Amount:          20000,
				RecognitionDate: time.Now(),
				ReceivedDate:    time.Now(),
			})
	}
	return mockDatas
}
