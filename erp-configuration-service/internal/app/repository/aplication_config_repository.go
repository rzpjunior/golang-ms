package repository

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/internal/app/model"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
)

type IApplicationConfigRepository interface {
	Get(ctx context.Context, req *dto.ApplicationConfigRequestGet) (applicationConfigs []*model.ApplicationConfig, count int64, err error)
	GetByID(ctx context.Context, id int64) (applicationConfig *model.ApplicationConfig, err error)
	GetDetail(ctx context.Context, req *pb.GetConfigAppDetailRequest) (res *model.ApplicationConfig, err error)
	Create(ctx context.Context, applicationConfig *model.ApplicationConfig) (err error)
	Update(ctx context.Context, applicationConfig *model.ApplicationConfig, columns ...string) (err error)
}

type ApplicationConfigRepository struct {
	opt opt.Options
}

func NewApplicationConfigRepository() IApplicationConfigRepository {
	return &ApplicationConfigRepository{
		opt: global.Setup.Common,
	}
}

func (r *ApplicationConfigRepository) Get(ctx context.Context, req *dto.ApplicationConfigRequestGet) (applicationConfigs []*model.ApplicationConfig, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ApplicationConfigRepository.Get")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.ApplicationConfig))

	cond := orm.NewCondition()

	if req.Search != "" {
		cond = cond.And("field__icontains", req.Search)
	}

	if req.Application != 0 {
		cond = cond.And("application", req.Application)
	}

	if req.Attribute != "" {
		cond = cond.And("attribute__icontains", req.Attribute)
	}

	if req.Value != "" {
		cond = cond.And("value__icontains", req.Value)
	}

	qs = qs.SetCond(cond)

	if req.OrderBy != "" {
		qs = qs.OrderBy(req.OrderBy)
	}

	count, err = qs.Offset(req.Offset).Limit(req.Limit).AllWithCtx(ctx, &applicationConfigs)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *ApplicationConfigRepository) GetByID(ctx context.Context, id int64) (applicationConfig *model.ApplicationConfig, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ApplicationConfigRepository.GetByID")
	defer span.End()

	applicationConfig = &model.ApplicationConfig{
		ID: id,
	}
	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, applicationConfig, "id")
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *ApplicationConfigRepository) GetDetail(ctx context.Context, req *pb.GetConfigAppDetailRequest) (applicationConfig *model.ApplicationConfig, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ApplicationConfigRepository.GetByID")
	defer span.End()
	var cols []string
	applicationConfig = &model.ApplicationConfig{}

	if req.Id != 0 {
		cols = append(cols, "id")
		applicationConfig.ID = int64(req.Id)
	}

	if req.Application != 0 {
		cols = append(cols, "application")
		applicationConfig.Application = int8(req.Application)
	}

	if req.Attribute != "" {
		cols = append(cols, "attribute")
		applicationConfig.Attribute = req.Attribute
	}
	if req.Field != "" {
		cols = append(cols, "field")
		applicationConfig.Field = req.Field
	}
	if req.Value != "" {
		cols = append(cols, "value")
		applicationConfig.Value = req.Value
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, applicationConfig, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}
func (r *ApplicationConfigRepository) Create(ctx context.Context, applicationConfig *model.ApplicationConfig) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ApplicationConfigRepository.Create")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.InsertWithCtx(ctx, applicationConfig)
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

func (r *ApplicationConfigRepository) Update(ctx context.Context, applicationConfig *model.ApplicationConfig, columns ...string) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ApplicationConfigRepository.Update")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.UpdateWithCtx(ctx, applicationConfig, columns...)
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
