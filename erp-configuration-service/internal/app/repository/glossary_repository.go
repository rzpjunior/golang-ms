package repository

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/internal/app/model"
)

type IGlossaryRepository interface {
	Get(ctx context.Context, offset int, limit int, table string, attribute string, valueInt int, valueName string) (glossarys []*model.Glossary, count int64, err error)
	GetDetail(ctx context.Context, table string, attribute string, valueInt int, valueName string) (glossarys *model.Glossary, err error)
}

type GlossaryRepository struct {
	opt opt.Options
}

func NewGlossaryRepository() IGlossaryRepository {
	return &GlossaryRepository{
		opt: global.Setup.Common,
	}
}

func (r *GlossaryRepository) Get(ctx context.Context, offset int, limit int, table string, attribute string, valueInt int, valueName string) (glossarys []*model.Glossary, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "GlossaryRepository.Get")
	defer span.End()

	db := r.opt.Database.Read
	qs := db.QueryTable(new(model.Glossary))

	cond := orm.NewCondition()

	if table != "" {
		cond = cond.And("table", table)
	}
	if attribute != "" {
		cond = cond.And("attribute", attribute)
	}
	if valueInt != 0 {
		cond = cond.And("value_int", valueInt)
	}
	if valueName != "" {
		cond = cond.And("value_name", valueName)
	}

	qs = qs.SetCond(cond)

	if offset != 0 || limit != 0 {
		qs = qs.Offset(offset).Limit(limit)
	}

	count, err = qs.AllWithCtx(ctx, &glossarys)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *GlossaryRepository) GetDetail(ctx context.Context, table string, attribute string, valueInt int, valueName string) (glossary *model.Glossary, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "GlossaryRepository.Get")
	defer span.End()

	db := r.opt.Database.Read

	cols := []string{}

	glossary = &model.Glossary{}

	if table != "" {
		glossary.Table = table
		cols = append(cols, "table")
	}

	if attribute != "" {
		glossary.Attribute = attribute
		cols = append(cols, "attribute")
	}

	if valueInt != 0 {
		glossary.ValueInt = int8(valueInt)
		cols = append(cols, "value_int")
	}

	if valueName != "" {
		glossary.ValueName = valueName
		cols = append(cols, "value_name")
	}

	err = db.ReadWithCtx(ctx, glossary, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}
