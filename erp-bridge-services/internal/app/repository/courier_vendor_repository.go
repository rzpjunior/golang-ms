package repository

import (
	"context"
	"fmt"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
)

type ICourierVendorRepository interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, siteID int64) (courierVendors []*model.CourierVendor, total int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (courier *model.CourierVendor, err error)
}

type CourierVendorRepository struct {
	opt opt.Options
}

func NewCourierVendorRepository() ICourierVendorRepository {
	return &CourierVendorRepository{
		opt: global.Setup.Common,
	}
}

func (r *CourierVendorRepository) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, siteID int64) (courierVendors []*model.CourierVendor, total int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "CourierRepository.Get")
	defer span.End()

	// DUMMY RETURN
	dummies := r.MockDatas(10)
	return dummies, int64(len(dummies)), nil

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.Courier))

	cond := orm.NewCondition()

	if search != "" {
		cond = cond.And("name__icontains", search)
	}

	if status != 0 {
		cond = cond.And("status", status)
	}

	if siteID != 0 {
		cond = cond.And("site_id", siteID)
	}

	qs = qs.SetCond(cond)

	if orderBy != "" {
		qs = qs.OrderBy(orderBy)
	}

	total, err = qs.Offset(offset).Limit(limit).AllWithCtx(ctx, &courierVendors)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *CourierVendorRepository) GetDetail(ctx context.Context, id int64, code string) (courierVendor *model.CourierVendor, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ClassRepository.GetDetail")
	defer span.End()

	// DUMMY RETURN
	return r.MockDatas(1)[0], nil

	courierVendor = &model.CourierVendor{}

	var cols []string

	if id != 0 {
		cols = append(cols, "id")
		courierVendor.ID = id
	}

	if code != "" {
		cols = append(cols, "code")
		courierVendor.Code = code
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, courierVendor, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *CourierVendorRepository) MockDatas(total int) (mockDatas []*model.CourierVendor) {
	for i := 1; i <= total; i++ {
		mockDatas = append(mockDatas,
			&model.CourierVendor{
				ID:     int64(i),
				Code:   fmt.Sprintf("CRV%d", i),
				Name:   fmt.Sprintf("Dummy CourierVendor %d", i),
				Note:   "note",
				Status: 1,
				SiteID: 1,
			})
	}
	return
}
