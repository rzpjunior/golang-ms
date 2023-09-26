package repository

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/model"
)

type IOtpOutgoingRepository interface {
	Create(ctx context.Context, reqOtpOutgoing *model.OtpOutgoing) (OtpOutgoing *model.OtpOutgoing, err error)
	GetDetail(ctx context.Context, reqOtpOutgoing *model.OtpOutgoing) (OtpOutgoing *model.OtpOutgoing, err error)
}

type OtpOutgoingRepository struct {
	opt opt.Options
}

func NewOtpOutgoingRepository() IOtpOutgoingRepository {
	return &OtpOutgoingRepository{
		opt: global.Setup.Common,
	}
}

func (r *OtpOutgoingRepository) Create(ctx context.Context, reqOtpOutgoing *model.OtpOutgoing) (OtpOutgoing *model.OtpOutgoing, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "OtpOutgoingRepository.Create")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.InsertWithCtx(ctx, reqOtpOutgoing)
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
func (r *OtpOutgoingRepository) GetDetail(ctx context.Context, reqOtpOutgoing *model.OtpOutgoing) (OtpOutgoing *model.OtpOutgoing, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "OtpOutgoingRepository.Create")
	defer span.End()

	var cols []string
	OtpOutgoing = &model.OtpOutgoing{}
	if reqOtpOutgoing.ID != 0 {
		cols = append(cols, "id")
		OtpOutgoing.ID = reqOtpOutgoing.ID
	}

	if reqOtpOutgoing.Application != 0 {
		cols = append(cols, "application")
		OtpOutgoing.Application = reqOtpOutgoing.Application
	}

	if reqOtpOutgoing.DeliveryStatus != 0 {
		cols = append(cols, "delivery_status")
		OtpOutgoing.DeliveryStatus = reqOtpOutgoing.DeliveryStatus
	}
	if reqOtpOutgoing.MessageType != 0 {
		cols = append(cols, "message_type")
		OtpOutgoing.MessageType = reqOtpOutgoing.MessageType
	}
	if reqOtpOutgoing.OtpStatus != 0 {
		cols = append(cols, "otp_status")
		OtpOutgoing.OtpStatus = reqOtpOutgoing.OtpStatus
	}
	if reqOtpOutgoing.UsageType != 0 {
		cols = append(cols, "usage_type")
		OtpOutgoing.UsageType = reqOtpOutgoing.UsageType
	}
	if reqOtpOutgoing.Vendor != 0 {
		cols = append(cols, "vendor")
		OtpOutgoing.Vendor = reqOtpOutgoing.Vendor
	}

	if reqOtpOutgoing.Message != "" {
		cols = append(cols, "message")
		OtpOutgoing.Message = reqOtpOutgoing.Message
	}
	if reqOtpOutgoing.OTP != "" {
		cols = append(cols, "otp")
		OtpOutgoing.OTP = reqOtpOutgoing.OTP
	}
	if reqOtpOutgoing.PhoneNumber != "" {
		cols = append(cols, "phone_number")
		OtpOutgoing.PhoneNumber = reqOtpOutgoing.PhoneNumber
	}
	if reqOtpOutgoing.VendorMessageID != "" {
		cols = append(cols, "vendor_message_id")
		OtpOutgoing.VendorMessageID = reqOtpOutgoing.VendorMessageID
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, OtpOutgoing, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}
