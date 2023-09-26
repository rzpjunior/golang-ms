package repository

import (
	"context"
	"fmt"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

type IAdmDivisionRepository interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, regionID int64, subDistrictID int64) (admDivisions []*model.AdmDivision, count int64, err error)
	GetDetail(ctx context.Context, id int64, code string, regionID int64, subDistrictID int64) (admDivision *model.AdmDivision, err error)
	GetListGRPC(ctx context.Context, req *bridgeService.GetAdmDivisionListRequest) (admDivisions []*model.AdmDivision, count int64, err error)
	GetDetailGRPC(ctx context.Context, req *bridgeService.GetAdmDivisionDetailRequest) (admDivision *model.AdmDivision, err error)
}

type AdmDivisionRepository struct {
	opt opt.Options
}

func NewAdmDivisionRepository() IAdmDivisionRepository {
	return &AdmDivisionRepository{
		opt: global.Setup.Common,
	}
}

func (r *AdmDivisionRepository) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, regionID int64, subDistrictID int64) (admDivisions []*model.AdmDivision, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "AdmDivisionRepository.Get")
	defer span.End()

	// RETURN DUMMIES
	dummies := r.MockDatas(10)
	return dummies, int64(len(dummies)), nil

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.AdmDivision))

	cond := orm.NewCondition()

	if search != "" {
		condGroup := orm.NewCondition()
		condGroup = condGroup.And("description__icontains", search).Or("code__icontains", search)
		cond = cond.AndCond(condGroup)
	}

	if status != 0 {
		cond = cond.And("status", status)
	}

	if regionID != 0 {
		cond = cond.And("region_id", regionID)
	}

	if subDistrictID != 0 {
		cond = cond.And("sub_district_id", subDistrictID)
	}

	qs = qs.SetCond(cond)

	if orderBy != "" {
		qs = qs.OrderBy(orderBy)
	}

	count, err = qs.Offset(offset).Limit(limit).AllWithCtx(ctx, &admDivisions)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *AdmDivisionRepository) GetDetail(ctx context.Context, id int64, code string, regionID int64, subDistrictID int64) (admDivision *model.AdmDivision, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "AdmDivisionRepository.GetDetail")
	defer span.End()

	// RETURN DUMMIES
	return r.MockDatas(1)[0], nil

	admDivision = &model.AdmDivision{}

	var cols []string

	if id != 0 {
		cols = append(cols, "id")
		admDivision.ID = id
	}

	if code != "" {
		cols = append(cols, "code")
		admDivision.Code = code
	}

	if regionID != 0 {
		cols = append(cols, "region_id")
		admDivision.RegionID = regionID
	}

	if subDistrictID != 0 {
		cols = append(cols, "sub_district_id")
		admDivision.SubDistrictID = subDistrictID
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, admDivision, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}
func (r *AdmDivisionRepository) GetListGRPC(ctx context.Context, req *bridgeService.GetAdmDivisionListRequest) (admDivisions []*model.AdmDivision, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "AdmDivisionRepository.Get")
	defer span.End()

	// RETURN DUMMIES
	dummies := r.MockDatas(10)
	return dummies, int64(len(dummies)), nil

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.AdmDivision))

	cond := orm.NewCondition()

	if req.Search != "" {
		condGroup := orm.NewCondition()
		condGroup = condGroup.And("description__icontains", req.Search).Or("code__icontains", req.Search)
		cond = cond.AndCond(condGroup)
	}

	if req.Status != 0 {
		cond = cond.And("status", req.Status)
	}

	if req.RegionId != 0 {
		cond = cond.And("region_id", req.RegionId)
	}

	if req.SubDistrictId != 0 {
		cond = cond.And("sub_district_id", req.SubDistrictId)
	}
	if req.CityId != 0 {
		cond = cond.And("city_id", req.CityId)
	}

	if req.DistrictId != 0 {
		cond = cond.And("district_id", req.DistrictId)
	}

	if req.ProvinceId != 0 {
		cond = cond.And("province_id", req.ProvinceId)
	}

	qs = qs.SetCond(cond)

	if req.OrderBy != "" {
		qs = qs.OrderBy(req.OrderBy)
	}

	count, err = qs.Offset(req.Offset).Limit(req.Limit).AllWithCtx(ctx, &admDivisions)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *AdmDivisionRepository) GetDetailGRPC(ctx context.Context, req *bridgeService.GetAdmDivisionDetailRequest) (admDivision *model.AdmDivision, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "AdmDivisionRepository.GetDetail")
	defer span.End()

	// RETURN DUMMIES
	return r.MockDatas(1)[0], nil

	admDivision = &model.AdmDivision{}

	var cols []string

	if req.Id != 0 {
		cols = append(cols, "id")
		admDivision.ID = req.Id
	}

	if req.Code != "" {
		cols = append(cols, "code")
		admDivision.Code = req.Code
	}

	if req.RegionId != 0 {
		cols = append(cols, "region_id")
		admDivision.RegionID = req.RegionId
	}

	if req.SubDistrictId != 0 {
		cols = append(cols, "sub_district_id")
		admDivision.SubDistrictID = req.SubDistrictId
	}

	if req.DistrictId != 0 {
		cols = append(cols, "district_id")
		admDivision.DistrictID = req.DistrictId
	}
	if req.CityId != 0 {
		cols = append(cols, "city_id")
		admDivision.CityID = req.CityId
	}

	if req.ProvinceId != 0 {
		cols = append(cols, "province_id")
		admDivision.ProvinceID = req.ProvinceId
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, admDivision, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}
func (r *AdmDivisionRepository) MockDatas(total int) (mockDatas []*model.AdmDivision) {
	for i := 1; i <= total; i++ {
		mockDatas = append(mockDatas,
			&model.AdmDivision{
				ID:            int64(i),
				Code:          fmt.Sprintf("ADM%d", i),
				RegionID:      1,
				City:          fmt.Sprintf("Dummy City %d", i),
				District:      fmt.Sprintf("Dummy District %d", i),
				Province:      fmt.Sprintf("Dummy Province %d", i),
				SubDistrictID: 1,
				PostalCode:    "12345",
				Status:        1,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			})
	}

	return
}
