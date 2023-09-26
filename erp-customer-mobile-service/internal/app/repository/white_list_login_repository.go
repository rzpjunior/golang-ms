package repository

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/model"
)

type IWhiteListLoginRepository interface {
	GetDetail(ctx context.Context, id int64, phoneNumber string, otp string) (whiteListLogin *model.WhiteListLogin, err error)
}

type WhiteListLoginRepository struct {
	opt opt.Options
}

func NewWhiteListLoginRepository() IWhiteListLoginRepository {
	return &WhiteListLoginRepository{
		opt: global.Setup.Common,
	}
}

func (r *WhiteListLoginRepository) GetDetail(ctx context.Context, id int64, phoneNumber string, otp string) (whiteListLogin *model.WhiteListLogin, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "WhiteListLoginRepository.GetDetail")
	defer span.End()

	var cols []string
	whiteListLogin = &model.WhiteListLogin{}
	if id != 0 {
		cols = append(cols, "id")
		whiteListLogin.ID = id
	}

	if phoneNumber != "" {
		cols = append(cols, "phone_number")
		whiteListLogin.PhoneNumber = phoneNumber
	}

	if otp != "" {
		cols = append(cols, "otp")
		whiteListLogin.OTP = otp
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, whiteListLogin, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}
