package repository

import (
	"context"
	"fmt"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

type IAddressRepository interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, archetypeID int64, admDivisionID int64, siteID int64, salespersonID int64, territoryID int64, taxScheduleID int64) (addresses []*model.Address, count int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (address *model.Address, err error)
	GetWithExcludedIds(ctx context.Context, offset int, limit int, status int, search string, orderBy string, archetypeID int64, admDivisionID int64, siteID int64, salespersonID int64, territoryID int64, taxScheduleID int64, excludedIds []int64) (addresses []*model.Address, count int64, err error)
	GetListGRPC(ctx context.Context, req *bridgeService.GetAddressListRequest) (addresses []*model.Address, count int64, err error)
	GetDetailGRPC(ctx context.Context, req *bridgeService.GetAddressDetailRequest) (address *model.Address, err error)
}

type AddressRepository struct {
	opt opt.Options
}

func NewAddressRepository() IAddressRepository {
	return &AddressRepository{
		opt: global.Setup.Common,
	}
}

func (r *AddressRepository) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, archetypeID int64, admDivisionID int64, siteID int64, salespersonID int64, territoryID int64, taxScheduleID int64) (addresses []*model.Address, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "AddressRepository.Get")
	defer span.End()

	dummies := r.MockDatas(10)
	return dummies, int64(len(dummies)), nil

	// TODO: connect to Mysql or GP
	// db := r.opt.Database.Read

	// qs := db.QueryTable(new(model.Address))

	// cond := orm.NewCondition()

	// if search != "" {
	// 	condGroup := orm.NewCondition()
	// 	condGroup = condGroup.And("code", search)
	// 	cond = cond.AndCond(condGroup)
	// }

	// if status != 0 {
	// 	cond = cond.And("status", status)
	// }

	// if archetypeID != 0 {
	// 	cond = cond.And("archetype_id", archetypeID)
	// }

	// if admDivisionID != 0 {
	// 	cond = cond.And("admDivision_id", admDivisionID)
	// }

	// if siteID != 0 {
	// 	cond = cond.And("site_id", siteID)
	// }

	// if salespersonID != 0 {
	// 	cond = cond.And("salesperson_id", salespersonID)
	// }

	// if territoryID != 0 {
	// 	cond = cond.And("territory_id", territoryID)
	// }

	// if taxScheduleID != 0 {
	// 	cond = cond.And("tax_schedule_id", taxScheduleID)
	// }

	// qs = qs.SetCond(cond)

	// if orderBy != "" {
	// 	qs = qs.OrderBy(orderBy)
	// }

	// count, err = qs.Offset(offset).Limit(limit).AllWithCtx(ctx, &addresses)
	// if err != nil {
	// 	span.RecordError(err)
	// 	return
	// }

	return
}

func (r *AddressRepository) GetDetail(ctx context.Context, id int64, code string) (address *model.Address, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "AddressRepository.GetDetail")
	defer span.End()

	return r.MockDatas(1)[0], nil

	// TODO: connect to Mysql or GP
	// var cols []string

	// if id != 0 {
	// 	cols = append(cols, "id")
	// 	address.ID = id
	// }

	// if code != "" {
	// 	cols = append(cols, "code")
	// 	address.Code = code
	// }

	// db := r.opt.Database.Read
	// err = db.ReadWithCtx(ctx, address, cols...)
	// if err != nil {
	// 	span.RecordError(err)
	// 	return
	// }

	return
}

func (r *AddressRepository) GetListGRPC(ctx context.Context, req *bridgeService.GetAddressListRequest) (addresses []*model.Address, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "AddressRepository.Get")
	defer span.End()

	dummies := r.MockDatas(10)
	return dummies, int64(len(dummies)), nil

	// TODO: connect to Mysql or GP
	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.Address))

	cond := orm.NewCondition()

	if req.Search != "" {
		condGroup := orm.NewCondition()
		condGroup = condGroup.And("code", req.Search)
		cond = cond.AndCond(condGroup)
	}

	if req.Status != 0 {
		cond = cond.And("status", req.Status)
	}

	if req.ArchetypeId != 0 {
		cond = cond.And("archetype_id", req.ArchetypeId)
	}

	if req.AdmDivisionId != 0 {
		cond = cond.And("admDivision_id", req.AdmDivisionId)
	}

	if req.SiteId != 0 {
		cond = cond.And("site_id", req.SiteId)
	}

	if req.SalespersonId != 0 {
		cond = cond.And("salesperson_id", req.SalespersonId)
	}

	if req.TerritoryId != 0 {
		cond = cond.And("territory_id", req.TerritoryId)
	}

	if req.TaxScheduleId != 0 {
		cond = cond.And("tax_schedule_id", req.TaxScheduleId)
	}

	if req.CustomerId != 0 {
		cond = cond.And("customer_id", req.CustomerId)
	}
	qs = qs.SetCond(cond)

	if req.OrderBy != "" {
		qs = qs.OrderBy(req.OrderBy)
	}

	count, err = qs.Offset(req.Offset).Limit(req.Limit).AllWithCtx(ctx, &addresses)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *AddressRepository) GetDetailGRPC(ctx context.Context, req *bridgeService.GetAddressDetailRequest) (address *model.Address, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "AddressRepository.GetDetail")
	defer span.End()

	return r.MockDatas(1)[0], nil

	// TODO: connect to Mysql or GP
	var cols []string

	if req.Id != 0 {
		cols = append(cols, "id")
		address.ID = req.Id
	}

	if req.Code != "" {
		cols = append(cols, "code")
		address.Code = req.Code
	}

	if req.CustomerId != 0 {
		cols = append(cols, "code")
		address.CustomerID = req.CustomerId
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, address, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *AddressRepository) GetWithExcludedIds(ctx context.Context, offset int, limit int, status int, search string, orderBy string, archetypeID int64, admDivisionID int64, siteID int64, salespersonID int64, territoryID int64, taxScheduleID int64, excludedIds []int64) (addresses []*model.Address, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "AddressRepository.Get")
	defer span.End()

	dummies := r.MockDatas(10)
	return dummies, int64(len(dummies)), nil

	// TODO: connect to Mysql or GP
}

func (r *AddressRepository) MockDatas(total int) (mockDatas []*model.Address) {
	for i := 1; i <= total; i++ {
		mockDatas = append(mockDatas,
			&model.Address{
				ID:               int64(i),
				Code:             fmt.Sprintf("ADR%d", i),
				CustomerName:     fmt.Sprintf("Dummy CustomerName %d", i),
				ArchetypeID:      1,
				AdmDivisionID:    1,
				SiteID:           1,
				SalespersonID:    1,
				TerritoryID:      1,
				AddressCode:      "12345",
				AddressName:      "Dummy AddressName",
				ContactPerson:    "Dummy ContactPerson",
				City:             "Dummy City",
				State:            "Dummy State",
				ZipCode:          "Dummy ZipCode",
				CountryCode:      "Dummy CountryCode",
				Country:          "Dummy Country",
				Latitude:         12.3456789,
				Longitude:        12.3456789,
				UpsZone:          "Dummy UpsZone",
				ShippingMethod:   "Dummy ShippingMethod",
				TaxScheduleID:    1,
				PrintPhoneNumber: 1,
				ShippingAddress:  "Dummy ShippingAddress",
				DistrictId:       1,
				Status:           1,
				BcaVa:            "Dummy BcaVa",
				OtherVa:          "Dummy OtherVa",
			})
	}

	return mockDatas
}
