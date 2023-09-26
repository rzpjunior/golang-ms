package repository

import (
	"context"
	"fmt"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/generator"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
)

type IPurchasePlanRepository interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (purchasePlans []*model.PurchasePlan, count int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (purchasePlan *model.PurchasePlan, err error)
}

type PurchasePlanRepository struct {
	opt opt.Options
}

func NewPurchasePlanRepository() IPurchasePlanRepository {
	return &PurchasePlanRepository{
		opt: global.Setup.Common,
	}
}

func (r *PurchasePlanRepository) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (purchasePlans []*model.PurchasePlan, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PurchasePlanRepository.Get")
	defer span.End()

	// RETURN DUMMIES
	count = 10
	for i := 1; i <= int(count); i++ {
		purchasePlans = append(purchasePlans, r.MockDatas(int64(i)))
	}
	return
}

func (r *PurchasePlanRepository) GetDetail(ctx context.Context, id int64, code string) (purchasePlan *model.PurchasePlan, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PurchasePlanRepository.GetDetail")
	defer span.End()

	// RETURN DUMMIES
	purchasePlan = r.MockDatas(id)
	return
}

func (r *PurchasePlanRepository) MockDatas(id int64) (mockDatas *model.PurchasePlan) {
	if id == 0 {
		id = 1
	}

	mockDatas = &model.PurchasePlan{
		ID:                   id,
		Code:                 fmt.Sprintf("PP%d", id),
		VendorOrganizationID: generator.DummyInt64(1, 10),
		SiteID:               generator.DummyInt64(1, 10),
		RecognitionDate:      generator.DummyTime(),
		EtaDate:              generator.DummyTime(),
		EtaTime:              "10:00",
		TotalPrice:           generator.DummyFloat64(10000, 100000),
		TotalWeight:          0.25 * float64(generator.DummyInt(1, 50)),
		TotalPurchasePlanQty: 0.25 * float64(generator.DummyInt(1, 50)),
		TotalPurchaseQty:     0.25 * float64(generator.DummyInt(1, 50)),
		Note:                 "Dummy Note",
		Status:               1,
		CreatedAt:            time.Now(),
		CreatedBy:            1,
	}

	if id%2 != 0 {
		mockDatas.AssignedTo = generator.DummyInt64(1, 5)
		mockDatas.AssignedBy = generator.DummyInt64(1, 5)
		mockDatas.AssignedAt = time.Now()
	}

	return mockDatas
}
