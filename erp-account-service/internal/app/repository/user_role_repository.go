package repository

import (
	"context"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/global"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/model"
)

func RepositoryUserRole() IUserRoleRepository {
	m := new(UserRoleRepository)
	m.opt = global.Setup.Common
	return m
}

type IUserRoleRepository interface {
	GetByUserID(ctx context.Context, userID int64) (userRole []*model.UserRole, err error)
	GetActiveByUserID(ctx context.Context, userID int64) (userRole []*model.UserRole, err error)
	GetByRoleID(ctx context.Context, roleID int64) (userRole []*model.UserRole, err error)
	Create(ctx context.Context, userRole *model.UserRole) (err error)
	Update(ctx context.Context, userRole *model.UserRole, columns ...string) (err error)
	Delete(ctx context.Context, id int64) (err error)
	ArchiveByUserID(ctx context.Context, userID int64) (err error)
}

type UserRoleRepository struct {
	opt opt.Options
}

func (r *UserRoleRepository) GetByUserID(ctx context.Context, userID int64) (userRole []*model.UserRole, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "UserRoleRepository.GetByRoleID")
	defer span.End()

	db := r.opt.Database.Read
	_, err = db.QueryTable(new(model.UserRole)).Filter("user_id", userID).AllWithCtx(ctx, &userRole)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *UserRoleRepository) GetActiveByUserID(ctx context.Context, userID int64) (userRole []*model.UserRole, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "UserRoleRepository.GetActiveByUserID")
	defer span.End()

	db := r.opt.Database.Read
	_, err = db.QueryTable(new(model.UserRole)).Filter("user_id", userID).Filter("status", 1).AllWithCtx(ctx, &userRole)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *UserRoleRepository) GetByRoleID(ctx context.Context, roleID int64) (userRole []*model.UserRole, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "UserRoleRepository.GetByRoleID")
	defer span.End()

	db := r.opt.Database.Read
	_, err = db.QueryTable(new(model.UserRole)).Filter("role_id", roleID).Filter("status", 1).AllWithCtx(ctx, &userRole)
	if err != nil {
		span.RecordError(err)
		return
	}
	return
}

func (r *UserRoleRepository) Create(ctx context.Context, userRole *model.UserRole) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "UserRoleRepository.Create")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.InsertWithCtx(ctx, userRole)

	if err == nil {
		err = tx.Commit()
		if err != nil {
			span.RecordError(err)
			return
		}
	} else {
		span.RecordError(err)
		err = tx.Rollback()
		if err != nil {
			span.RecordError(err)
			return
		}
	}

	return
}

func (r *UserRoleRepository) Update(ctx context.Context, userRole *model.UserRole, columns ...string) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "UserRoleRepository.Update")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.UpdateWithCtx(ctx, userRole, columns...)

	if err == nil {
		err = tx.Commit()
		if err != nil {
			span.RecordError(err)
			return
		}
	} else {
		span.RecordError(err)
		err = tx.Rollback()
		if err != nil {
			span.RecordError(err)
			return
		}
	}

	return
}

func (r *UserRoleRepository) Delete(ctx context.Context, id int64) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "UserRoleRepository.Delete")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	userRole := &model.UserRole{
		ID: id,
	}

	_, err = tx.DeleteWithCtx(ctx, userRole)

	if err == nil {
		err = tx.Commit()
		if err != nil {
			span.RecordError(err)
			return
		}
	} else {
		span.RecordError(err)
		err = tx.Rollback()
		if err != nil {
			span.RecordError(err)
			return
		}
	}

	return
}

func (r *UserRoleRepository) ArchiveByUserID(ctx context.Context, userID int64) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "UserRoleRepository.Update")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	value := &orm.Params{
		"Status":    2,
		"UpdatedAt": time.Now(),
	}

	_, err = tx.QueryTable(new(model.UserRole)).Filter("user_id", userID).UpdateWithCtx(ctx, *value)

	if err == nil {
		err = tx.Commit()
		if err != nil {
			span.RecordError(err)
			return
		}
	} else {
		span.RecordError(err)
		err = tx.Rollback()
		if err != nil {
			span.RecordError(err)
			return
		}
	}

	return
}
