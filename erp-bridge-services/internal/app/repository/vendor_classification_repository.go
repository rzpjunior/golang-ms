package repository

import (
	"context"
	"fmt"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
)

type IVendorClassificationRepository interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (vendorClassifications []*model.VendorClassification, count int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (vendorClassification *model.VendorClassification, err error)
}

type VendorClassificationRepository struct {
	opt opt.Options
}

func NewVendorClassificationRepository() IVendorClassificationRepository {
	return &VendorClassificationRepository{
		opt: global.Setup.Common,
	}
}

func (r *VendorClassificationRepository) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (vendorClassifications []*model.VendorClassification, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "VendorClassificationRepository.Get")
	defer span.End()

	// RETURN DUMMIES
	count = 10
	for i := 1; i <= int(count); i++ {
		vendorClassifications = append(vendorClassifications, r.MockDatas(int64(i)))
	}
	return
}

func (r *VendorClassificationRepository) GetDetail(ctx context.Context, id int64, code string) (vendorClassification *model.VendorClassification, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "VendorClassificationRepository.GetDetail")
	defer span.End()

	// RETURN DUMMIES
	vendorClassification = r.MockDatas(id)
	return
}

func (r *VendorClassificationRepository) MockDatas(id int64) (mockDatas *model.VendorClassification) {
	if id == 0 {
		id = 1
	}
	mockDatas = &model.VendorClassification{
		ID:            id,
		CommodityCode: fmt.Sprintf("VCC%d", id),
		CommodityName: "Dummy CommodityName",
		BadgeCode:     fmt.Sprintf("VCB%d", id),
		BadgeName:     "Dummy BadgeName",
		TypeCode:      fmt.Sprintf("VCT%d", id),
		TypeName:      "Dummy TypeName",
	}

	return mockDatas
}
