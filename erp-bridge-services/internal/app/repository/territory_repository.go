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

type ITerritoryRepository interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, regionID int64, salespersonID int64, CustomerTypeID int64, subDistrictID int64) (territories []*model.Territory, count int64, err error)
	GetDetail(ctx context.Context, id int64, code string, salespersonID int64) (territory *model.Territory, err error)
}

type TerritoryRepository struct {
	opt opt.Options
}

func NewTerritoryRepository() ITerritoryRepository {
	return &TerritoryRepository{
		opt: global.Setup.Common,
	}
}

func (r *TerritoryRepository) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, regionID int64, salespersonID int64, CustomerTypeID int64, subDistrictID int64) (territories []*model.Territory, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "TerritoryRepository.Get")
	defer span.End()

	dummies := r.MockDatas(10)
	return dummies, int64(len(dummies)), nil

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.Territory))

	cond := orm.NewCondition()

	if search != "" {
		condGroup := orm.NewCondition()
		condGroup = condGroup.And("description__icontains", search).Or("code__icontains", search)
		cond = cond.AndCond(condGroup)
	}

	if status != 0 {
		cond.And("status", status)
	}

	if regionID != 0 {
		cond.And("region_id", regionID)
	}

	if salespersonID != 0 {
		cond.And("salesperson_id", salespersonID)
	}

	if CustomerTypeID != 0 {
		cond.And("customer_type_id", CustomerTypeID)
	}

	if subDistrictID != 0 {
		cond.And("sub_district_id", subDistrictID)
	}
	qs = qs.SetCond(cond)

	if orderBy != "" {
		qs = qs.OrderBy(orderBy)
	}

	count, err = qs.Offset(offset).Limit(limit).AllWithCtx(ctx, &territories)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *TerritoryRepository) GetDetail(ctx context.Context, id int64, code string, salespersonID int64) (territory *model.Territory, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "TerritoryRepository.GetDetail")
	defer span.End()

	return r.MockDatas(1)[0], nil

	territory = &model.Territory{}

	var cols []string

	if id != 0 {
		cols = append(cols, "id")
		territory.ID = id
	}

	if code != "" {
		cols = append(cols, "code")
		territory.Code = code
	}

	if salespersonID != 0 {
		cols = append(cols, "salesperson_id")
		territory.SalespersonID = salespersonID
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, territory, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *TerritoryRepository) MockDatas(total int) (mockDatas []*model.Territory) {
	for i := 1; i <= total; i++ {
		mockDatas = append(mockDatas,
			&model.Territory{
				ID:             int64(i),
				Code:           fmt.Sprintf("TER%d", i),
				Description:    fmt.Sprintf("Dummy Territory %d", i),
				RegionID:       1,
				SalespersonID:  1,
				CustomerTypeID: 1,
				SubDistrictID:  1,
				CreatedAt:      generator.DummyTime(),
				UpdatedAt:      generator.DummyTime(),
			})
	}
	return
}
