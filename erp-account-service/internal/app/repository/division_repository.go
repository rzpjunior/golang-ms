package repository

import (
	"context"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/constants"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/utils"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/global"
	"google.golang.org/protobuf/types/known/timestamppb"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/model"
	auditService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/audit_service"
)

func RepositoryDivision() IDivisionRepository {
	m := new(DivisionRepository)
	m.opt = global.Setup.Common
	return m
}

type IDivisionRepository interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (divisions []*model.Division, count int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (division *model.Division, err error)
	GetByName(ctx context.Context, name string) (division *model.Division, err error)
	Create(ctx context.Context, division *model.Division) (err error)
	Update(ctx context.Context, division *model.Division, columns ...string) (err error)
	Archive(ctx context.Context, id int64) (err error)
}

type DivisionRepository struct {
	opt opt.Options
}

func (r *DivisionRepository) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (divisions []*model.Division, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "DivisionRepository.Get")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.Division))

	cond := orm.NewCondition()

	if search != "" {
		condGroup := orm.NewCondition()
		condGroup = condGroup.And("name__icontains", search).Or("code__icontains", search)
		cond = cond.AndCond(condGroup)
	}

	if status != 0 {
		cond = cond.And("status", status)
	}

	qs = qs.SetCond(cond)

	if orderBy != "" {
		qs = qs.OrderBy(orderBy)
	}

	count, err = qs.Offset(offset).Limit(limit).AllWithCtx(ctx, &divisions)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *DivisionRepository) GetDetail(ctx context.Context, id int64, code string) (division *model.Division, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "DivisionRepository.GetByID")
	defer span.End()
	var cols []string
	division = &model.Division{}

	if id != 0 {
		division.ID = id
		cols = append(cols, "id")
	}

	if code != "" {
		division.Code = code
		cols = append(cols, "code")
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, division, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *DivisionRepository) GetByName(ctx context.Context, name string) (division *model.Division, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "DivisionRepository.GetByName")
	defer span.End()

	division = &model.Division{
		Name: name,
	}
	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, division, "name")
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *DivisionRepository) Create(ctx context.Context, division *model.Division) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "DivisionRepository.Create")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.InsertWithCtx(ctx, division)
	if err != nil {
		span.RecordError(err)
		tx.Rollback()
		return
	}

	// insert mongo for audit log

	userID := ctx.Value(constants.KeyUserID).(int64)

	_, err = r.opt.Client.AuditServiceGrpc.CreateLog(ctx, &auditService.CreateLogRequest{
		Log: &auditService.Log{
			UserId:      userID,
			ReferenceId: utils.ToString(division.ID),
			Type:        "division",
			Function:    "Create",
			CreatedAt:   timestamppb.New(time.Now()),
		},
	})

	if err != nil {
		r.opt.Logger.AddMessage(log.ErrorLevel, err)
		span.RecordError(err)
		tx.Rollback()
		err = edenlabs.ErrorRpcNotFound("audit", "audit_log")
		return
	}

	err = tx.Commit()
	if err != nil {
		span.RecordError(err)
		tx.Rollback()
		return
	}

	return
}

func (r *DivisionRepository) Update(ctx context.Context, division *model.Division, columns ...string) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "DivisionRepository.Update")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.UpdateWithCtx(ctx, division, columns...)
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

func (r *DivisionRepository) Archive(ctx context.Context, id int64) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "DivisionRepository.Archive")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	division := &model.Division{
		ID:        id,
		Status:    2,
		UpdatedAt: time.Now(),
	}

	_, err = tx.UpdateWithCtx(ctx, division, "Status", "UpdatedAt")
	if err != nil {
		span.RecordError(err)
		tx.Rollback()
		return
	}

	value := &orm.Params{
		"Status":    2,
		"UpdatedAt": time.Now(),
	}

	qsRole := tx.QueryTable(new(model.Role)).Filter("division_id", id)

	var roles []*model.Role
	qsRole.AllWithCtx(ctx, &roles)

	for _, role := range roles {
		_, err = tx.QueryTable(new(model.UserRole)).Filter("role_id", role.ID).UpdateWithCtx(ctx, *value)
		if err != nil {
			span.RecordError(err)
			return
		}
	}

	_, err = qsRole.UpdateWithCtx(ctx, *value)
	if err != nil {
		span.RecordError(err)
		return
	}

	err = tx.Commit()
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}
