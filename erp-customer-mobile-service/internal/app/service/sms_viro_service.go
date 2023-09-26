package service

import (
	"context"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/repository"
)

type ISMSViroService interface {
	Update(ctx context.Context, req dto.UpdateRequestSMSViro) (res dto.UpdateRequestSMSViro, err error)
}

type SMSViroService struct {
	opt                   opt.Options
	RepositoryOTPOutgoing repository.IOtpOutgoingRepository
}

func NewSMSViroService() ISMSViroService {
	return &SMSViroService{
		opt:                   global.Setup.Common,
		RepositoryOTPOutgoing: repository.NewOtpOutgoingRepository(),
	}
}

func (s *SMSViroService) Update(ctx context.Context, req dto.UpdateRequestSMSViro) (res dto.UpdateRequestSMSViro, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SMSViroService.Update")
	defer span.End()
	otpOutgoing := &model.OtpOutgoing{
		VendorMessageID: res.Results[0].MessageId,
	}
	otpOut, err := s.RepositoryOTPOutgoing.GetDetail(ctx, otpOutgoing)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	otpOut.DeliveryStatus = req.Results[0].Status.GroupID
	otpOut.UpdatedAt = time.Now()
	_, err = s.RepositoryOTPOutgoing.Create(ctx, otpOut)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	return
}
