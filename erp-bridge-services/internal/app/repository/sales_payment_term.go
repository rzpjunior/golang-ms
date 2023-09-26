package repository

import (
	"context"
	"fmt"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
)

type ISalesPaymentTermRepository interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (paymentTerms []*model.SalesPaymentTerm, count int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (paymentTerm *model.SalesPaymentTerm, err error)
}

type SalesPaymentTermRepository struct {
	opt opt.Options
}

func NewSalesPaymentTermRepository() ISalesPaymentTermRepository {
	return &SalesPaymentTermRepository{
		opt: global.Setup.Common,
	}
}

func (r *SalesPaymentTermRepository) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (paymentTerms []*model.SalesPaymentTerm, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SalesPaymentTermRepository.Get")
	defer span.End()

	// RETURN DUMMIES
	dummies := r.MockDatas(10)
	return dummies, int64(len(dummies)), nil

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.SalesPaymentTerm))

	cond := orm.NewCondition()

	if search != "" {
		condGroup := orm.NewCondition()
		condGroup = condGroup.And("description__icontains", search).Or("code__icontains", search)
		cond = cond.AndCond(condGroup)
	}

	if status != 0 {
		cond = cond.And("status", status)
	}

	qs = qs.SetCond(cond)

	if orderBy != "" {
		qs = qs.OrderBy(orderBy)
	}

	count, err = qs.Offset(offset).Limit(limit).AllWithCtx(ctx, &paymentTerms)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *SalesPaymentTermRepository) GetDetail(ctx context.Context, id int64, code string) (paymentTerm *model.SalesPaymentTerm, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SalesPaymentTermRepository.GetDetail")
	defer span.End()

	// RETURN DUMMIES
	return r.MockDatas(1)[0], nil

	paymentTerm = &model.SalesPaymentTerm{}

	var cols []string

	if id != 0 {
		cols = append(cols, "id")
		paymentTerm.ID = id
	}

	if code != "" {
		cols = append(cols, "code")
		paymentTerm.Code = code
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, paymentTerm, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *SalesPaymentTermRepository) MockDatas(total int) (mockDatas []*model.SalesPaymentTerm) {
	for i := 1; i <= total; i++ {
		mockDatas = append(mockDatas,
			&model.SalesPaymentTerm{
				ID:          int64(i),
				Code:        fmt.Sprintf("REG%d", i),
				Description: fmt.Sprintf("Dummy PaymentTerm %d", i),
				Status:      1,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			})
	}

	return mockDatas
}
