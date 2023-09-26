package repository

import (
	"context"
	"fmt"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/generator"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
)

type IVendorOrganizationRepository interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (vendorOrganizations []*model.VendorOrganization, count int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (vendorOrganization *model.VendorOrganization, err error)
}

type VendorOrganizationRepository struct {
	opt opt.Options
}

func NewVendorOrganizationRepository() IVendorOrganizationRepository {
	return &VendorOrganizationRepository{
		opt: global.Setup.Common,
	}
}

func (r *VendorOrganizationRepository) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (vendorOrganizations []*model.VendorOrganization, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "VendorOrganizationRepository.Get")
	defer span.End()

	// RETURN DUMMIES
	count = 10
	for i := 1; i <= int(count); i++ {
		vendorOrganizations = append(vendorOrganizations, r.MockDatas(int64(i)))
	}
	return
}

func (r *VendorOrganizationRepository) GetDetail(ctx context.Context, id int64, code string) (vendorOrganization *model.VendorOrganization, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "VendorOrganizationRepository.GetDetail")
	defer span.End()

	// RETURN DUMMIES
	vendorOrganization = r.MockDatas(id)
	return
}

func (r *VendorOrganizationRepository) MockDatas(id int64) (mockDatas *model.VendorOrganization) {
	if id == 0 {
		id = 1
	}
	mockDatas = &model.VendorOrganization{
		ID:                     id,
		Code:                   fmt.Sprintf("VENORG%d", id),
		VendorClassificationID: generator.DummyInt64(1, 10),
		SubDistrictID:          generator.DummyInt64(1, 10),
		PaymentTermID:          generator.DummyInt64(1, 10),
		Name:                   "Dummy VendorOrganization",
		Address:                "Dummy Address",
		Note:                   "Dummy Note",
		Status:                 1,
	}

	return mockDatas
}
