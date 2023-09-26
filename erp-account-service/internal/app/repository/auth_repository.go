package repository

import (
	"context"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/model"
)

func RepositoryAuth() IAuthRepository {
	m := new(AuthRepository)
	m.opt = global.Setup.Common
	return m
}

type IAuthRepository interface {
	GetByEmail(ctx context.Context, email, password string) (user *model.User, err error)
}

type AuthRepository struct {
	opt opt.Options
}

func (r *AuthRepository) GetByEmail(ctx context.Context, email, password string) (user *model.User, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "StaffRepository.GetByUserID")
	defer span.End()

	db := r.opt.Database.Read
	user = &model.User{
		Email:  email,
		Status: 1,
	}
	if err = db.ReadWithCtx(ctx, user, "email", "status"); err != nil {
		return
	}

	return
}
