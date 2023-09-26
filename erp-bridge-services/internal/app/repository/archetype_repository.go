package repository

import (
	"context"
	"fmt"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
)

type IArchetypeRepository interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, CustomerTypeID int64) (archetypes []*model.Archetype, count int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (archetype *model.Archetype, err error)
}

type ArchetypeRepository struct {
	opt opt.Options
}

func NewArchetypeRepository() IArchetypeRepository {
	return &ArchetypeRepository{
		opt: global.Setup.Common,
	}
}

func (r *ArchetypeRepository) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, CustomerTypeID int64) (archetypes []*model.Archetype, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ArchetypeRepository.Get")
	defer span.End()

	// RETURN DUMMIES
	dummies := r.MockDatas(10)
	return dummies, int64(len(dummies)), nil

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.Archetype))

	cond := orm.NewCondition()

	if search != "" {
		condGroup := orm.NewCondition()
		condGroup = condGroup.And("description__icontains", search).Or("code__icontains", search)
		cond = cond.AndCond(condGroup)
	}

	if status != 0 {
		cond = cond.And("status", status)
	}

	if CustomerTypeID != 0 {
		cond = cond.And("customer_type_id", CustomerTypeID)
	}

	qs = qs.SetCond(cond)

	if orderBy != "" {
		qs = qs.OrderBy(orderBy)
	}

	count, err = qs.Offset(offset).Limit(limit).AllWithCtx(ctx, &archetypes)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *ArchetypeRepository) GetDetail(ctx context.Context, id int64, code string) (archetype *model.Archetype, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ArchetypeRepository.GetDetail")
	defer span.End()

	// RETURN DUMMIES
	return r.MockDatas(1)[0], nil

	archetype = &model.Archetype{}

	var cols []string

	if id != 0 {
		cols = append(cols, "id")
		archetype.ID = id
	}

	if code != "" {
		cols = append(cols, "code")
		archetype.Code = code
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, archetype, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *ArchetypeRepository) MockDatas(total int) (mockDatas []*model.Archetype) {
	for i := 1; i <= total; i++ {
		mockDatas = append(mockDatas,
			&model.Archetype{
				ID:             int64(i),
				Code:           fmt.Sprintf("ARC%d", i),
				Description:    fmt.Sprintf("Dummy Archetype %d", i),
				CustomerTypeID: 1,
				Status:         1,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			})
	}

	return mockDatas
}
