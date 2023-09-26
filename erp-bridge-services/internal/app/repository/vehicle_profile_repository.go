package repository

import (
	"context"
	"fmt"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
)

type IVehicleProfileRepository interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, courierVendorID int64) (vehicleProfiles []*model.VehicleProfile, total int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (vehicleProfile *model.VehicleProfile, err error)
}

type VehicleProfileRepository struct {
	opt opt.Options
}

func NewVehicleProfileRepository() IVehicleProfileRepository {
	return &VehicleProfileRepository{
		opt: global.Setup.Common,
	}
}

func (r *VehicleProfileRepository) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, courierVendorID int64) (vehicleProfiles []*model.VehicleProfile, total int64, err error) {
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

	if courierVendorID != 0 {
		cond = cond.And("courier_vendor_id", courierVendorID)
	}

	qs = qs.SetCond(cond)

	if orderBy != "" {
		qs = qs.OrderBy(orderBy)
	}

	total, err = qs.Offset(offset).Limit(limit).AllWithCtx(ctx, &vehicleProfiles)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *VehicleProfileRepository) GetDetail(ctx context.Context, id int64, code string) (vehicleProfile *model.VehicleProfile, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "VehicleProfileRepository.GetDetail")
	defer span.End()

	// DUMMY RETURN
	return r.MockDatas(1)[0], nil

	vehicleProfile = &model.VehicleProfile{}

	var cols []string

	if id != 0 {
		cols = append(cols, "id")
		vehicleProfile.ID = id
	}

	if code != "" {
		cols = append(cols, "code")
		vehicleProfile.Code = code
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, vehicleProfile, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *VehicleProfileRepository) MockDatas(total int) (mockDatas []*model.VehicleProfile) {
	for i := 1; i <= total; i++ {
		mockDatas = append(mockDatas,
			&model.VehicleProfile{
				ID:                  int64(i),
				Code:                fmt.Sprintf("VHV%d", i),
				Name:                fmt.Sprintf("Dummy VehicleProvile %d", i),
				MaxKoli:             5,
				MaxWeight:           5,
				MaxFragile:          5,
				SpeedFactor:         1,
				RoutingProfile:      1,
				Skills:              "1,2",
				InitialCost:         10000,
				SubsequentCost:      5000,
				MaxAvailableVehicle: 10,
				Status:              1,
				CourierVendorID:     1,
			})
	}
	return mockDatas
}
