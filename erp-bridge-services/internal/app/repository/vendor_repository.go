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

type IVendorRepository interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (vendors []*model.Vendor, count int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (vendor *model.Vendor, err error)
}

type VendorRepository struct {
	opt opt.Options
}

func NewVendorRepository() IVendorRepository {
	return &VendorRepository{
		opt: global.Setup.Common,
	}
}

func (r *VendorRepository) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (vendors []*model.Vendor, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "VendorRepository.Get")
	defer span.End()

	// RETURN DUMMIES
	count = 10
	for i := 1; i <= int(count); i++ {
		vendors = append(vendors, r.MockDatas(int64(i)))
	}
	return
}

func (r *VendorRepository) GetDetail(ctx context.Context, id int64, code string) (vendor *model.Vendor, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "VendorRepository.GetDetail")
	defer span.End()

	// RETURN DUMMIES
	vendor = r.MockDatas(id)
	return
}

func (r *VendorRepository) MockDatas(id int64) (mockDatas *model.Vendor) {
	if id == 0 {
		id = 1
	}
	mockDatas = &model.Vendor{
		ID:                     id,
		Code:                   fmt.Sprintf("VEN%d", id),
		VendorOrganizationID:   generator.DummyInt64(1, 10),
		VendorClassificationID: generator.DummyInt64(1, 10),
		SubDistrictID:          generator.DummyInt64(1, 10),
		PicName:                "Dummy PicName",
		Email:                  "Dummy Email",
		PhoneNumber:            "Dummy PhoneNumber",
		PaymentTermID:          generator.DummyInt64(1, 10),
		Rejectable:             1,
		Returnable:             1,
		Address:                "Dummy Address",
		Note:                   "Dummy Note",
		Status:                 1,
		Latitude:               "",
		Longitude:              "",
		CreatedAt:              time.Now(),
		CreatedBy:              1,
	}

	return mockDatas
}
