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

type ISiteRepository interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (sites []*model.Site, count int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (site *model.Site, err error)
	GetInIds(ctx context.Context, offset int, limit int, status int, ids []int64, orderBy string) (sites []*model.Site, count int64, err error)
}

type SiteRepository struct {
	opt opt.Options
}

func NewSiteRepository() ISiteRepository {
	return &SiteRepository{
		opt: global.Setup.Common,
	}
}

func (r *SiteRepository) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (sites []*model.Site, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SiteRepository.Get")
	defer span.End()

	count = 10
	for i := 1; i <= int(count); i++ {
		sites = append(sites, r.MockDatas(int64(i)))
	}
	return

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.Site))

	cond := orm.NewCondition()

	if search != "" {
		condGroup := orm.NewCondition()
		condGroup = condGroup.And("description__icontains", search).Or("code__icontains", search)
		cond = cond.AndCond(condGroup)
	}
	if status != 0 {
		cond = cond.And("status", status)
	}

	qs = qs.SetCond(cond)

	if orderBy != "" {
		qs = qs.OrderBy(orderBy)
	}

	count, err = qs.Offset(offset).Limit(limit).AllWithCtx(ctx, &sites)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *SiteRepository) GetInIds(ctx context.Context, offset int, limit int, status int, ids []int64, orderBy string) (sites []*model.Site, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SiteRepository.Get")
	defer span.End()

	count = 10
	for i := 1; i <= int(count); i++ {
		sites = append(sites, r.MockDatas(int64(i)))
	}
	return

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.Site))

	cond := orm.NewCondition()

	if ids != nil {
		condGroup := orm.NewCondition()
		condGroup = condGroup.And("id__icontains", ids).Or("code__icontains", ids)
		cond = cond.AndCond(condGroup)
	}
	if status != 0 {
		cond = cond.And("status", status)
	}

	qs = qs.SetCond(cond)

	if orderBy != "" {
		qs = qs.OrderBy(orderBy)
	}

	count, err = qs.Offset(offset).Limit(limit).AllWithCtx(ctx, &sites)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *SiteRepository) GetDetail(ctx context.Context, id int64, code string) (site *model.Site, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PermissionRepository.GetDetail")
	defer span.End()

	site = r.MockDatas(id)
	return

	site = &model.Site{}

	var cols []string

	if id != 0 {
		cols = append(cols, "id")
		site.ID = id
	}

	if code != "" {
		cols = append(cols, "code")
		site.Code = code
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, site, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *SiteRepository) MockDatas(id int64) (mockDatas *model.Site) {
	if id == 0 {
		id = 1
	}
	mockDatas = &model.Site{
		ID:          id,
		Code:        fmt.Sprintf("STE%d", id),
		Description: fmt.Sprintf("Dummy Site %d", id),
		Status:      1,
		CreatedAt:   generator.DummyTime(),
		UpdatedAt:   generator.DummyTime(),
	}

	return
}
