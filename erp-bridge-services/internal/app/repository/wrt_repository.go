package repository

import (
	"context"
	"fmt"
	"strings"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
)

type IWrtRepository interface {
	Get(ctx context.Context, offset, limit int, regionID int64, search string) (Wrtes []*model.Wrt, count int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (Wrt *model.Wrt, err error)
}

type WrtRepository struct {
	opt opt.Options
}

func NewWrtRepository() IWrtRepository {
	return &WrtRepository{
		opt: global.Setup.Common,
	}
}

func (r *WrtRepository) Get(ctx context.Context, offset, limit int, regionID int64, search string) (Wrtes []*model.Wrt, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "WrtRepository.Get")
	defer span.End()

	// RETURN DUMMIES
	dummies := r.MockDatas(10)
	return dummies, int64(len(dummies)), nil

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.Wrt))

	cond := orm.NewCondition()

	if regionID != 0 {
		cond = cond.And("region_id", regionID)
	}

	if search != "" {
		condGroup := orm.NewCondition()
		// for special conditions if user filter by full name of time wrt on dashboard
		if strings.Contains(search, "-") {
			timeWrt := strings.Split(search, "-")
			condGroup = condGroup.And("start_time__icontains", timeWrt[0]).And("end_time__icontains", timeWrt[1])

		} else {
			condGroup = condGroup.Or("start_time__icontains", search).Or("end_time__icontains", search).Or("code__icontains", search)
		}
		cond = cond.AndCond(condGroup)
	}

	qs = qs.SetCond(cond)

	count, err = qs.Offset(offset).Limit(limit).AllWithCtx(ctx, &Wrtes)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *WrtRepository) GetDetail(ctx context.Context, id int64, code string) (Wrt *model.Wrt, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "WrtRepository.GetByCode")
	defer span.End()

	// RETURN DUMMIES
	return r.MockDatas(1)[0], nil

	Wrt = &model.Wrt{}

	var cols []string

	if id != 0 {
		cols = append(cols, "id")
		Wrt.ID = id
	}

	if code != "" {
		cols = append(cols, "code")
		Wrt.Code = code
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, Wrt, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *WrtRepository) MockDatas(total int) (mockDatas []*model.Wrt) {
	for i := 1; i <= total; i++ {
		mockDatas = append(mockDatas,
			&model.Wrt{
				ID:        int64(i),
				RegionID:  1,
				Code:      fmt.Sprintf("WRT%d", i),
				StartTime: "00.00",
				EndTime:   "23:59",
			})
	}
	return mockDatas
}
