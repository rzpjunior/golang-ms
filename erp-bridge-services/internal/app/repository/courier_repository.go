package repository

import (
	"context"
	"fmt"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/generator"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
)

type ICourierRepository interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, vehicleProfileID int64, emergencyMode int64) (couriers []*model.Courier, total int64, err error)
	GetDetail(ctx context.Context, id int64, code string, userID int64) (courier *model.Courier, err error)
	Create(ctx context.Context, model *model.Courier) (err error)
	Update(ctx context.Context, model *model.Courier, columns ...string) (err error)
}

type CourierRepository struct {
	opt opt.Options
}

func NewCourierRepository() ICourierRepository {
	return &CourierRepository{
		opt: global.Setup.Common,
	}
}

func (r *CourierRepository) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, vehicleProfileID int64, emergencyMode int64) (couriers []*model.Courier, total int64, err error) {
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

	if vehicleProfileID != 0 {
		cond = cond.And("vehicle_profile_id", vehicleProfileID)
	}

	if emergencyMode != 0 {
		cond = cond.And("emergency_mode", emergencyMode)
	}

	qs = qs.SetCond(cond)

	if orderBy != "" {
		qs = qs.OrderBy(orderBy)
	}

	total, err = qs.Offset(offset).Limit(limit).AllWithCtx(ctx, &couriers)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *CourierRepository) GetDetail(ctx context.Context, id int64, code string, userID int64) (courier *model.Courier, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ClassRepository.GetDetail")
	defer span.End()

	// DUMMY RETURN
	return r.MockDatas(1)[0], nil

	courier = &model.Courier{}

	var cols []string

	if id != 0 {
		cols = append(cols, "id")
		courier.ID = id
	}

	if code != "" {
		cols = append(cols, "code")
		courier.Code = code
	}

	if userID != 0 {
		cols = append(cols, "user_id")
		courier.UserID = userID
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, courier, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *CourierRepository) Create(ctx context.Context, model *model.Courier) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "CourierRepository.Create")
	defer span.End()

	// RETURN DUMMY
	return nil

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

func (r *CourierRepository) Update(ctx context.Context, model *model.Courier, columns ...string) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "CourierRepository.Update")
	defer span.End()

	// RETURN DUMMY
	return nil

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

func (r *CourierRepository) MockDatas(total int) (mockDatas []*model.Courier) {
	for i := 1; i <= total; i++ {
		mockDatas = append(mockDatas,
			&model.Courier{
				ID:                int64(i),
				Code:              fmt.Sprintf("CRR%d", i),
				Name:              fmt.Sprintf("Dummy Courier %d", i),
				PhoneNumber:       "8987654321",
				LicensePlate:      "du 33 y",
				EmergencyMode:     1,
				LastEmergencyTime: generator.DummyTime(),
				Status:            1,
				RoleID:            1,
				UserID:            1,
				VehicleProfileID:  1,
			})
	}
	return
}
