package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/dto"
)

type IProspectCustomerService interface {
	Upgrade(ctx context.Context, req dto.RequestUpgradeBusiness) (err error)
}

type ProspectCustomerService struct {
	opt opt.Options
	//RepositoryOTPOutgoing repository.IOtpOutgoingRepository
}

func NewProspectCustomerService() ICustomerService {
	return &CustomerService{
		opt: global.Setup.Common,
		//RepositoryOTPOutgoing: repository.NewOtpOutgoingRepository(),
	}
}
func (s *ProspectCustomerService) Upgrade(ctx context.Context, req dto.RequestUpgradeBusiness) (err error) {
	if !req.Data.TNCDataIsRight {
		// o.Failure("tnc_data_is_right.invalid", "Mohon untuk konfirmasi data sudah benar")
		return
	}

	// glossaryRegChannel, e := repository.GetGlossaryMultipleValue("table", "prospect_customer", "attribute", "reg_channel", "value_name", c.Platform)
	// if e != nil {
	// 	o.Failure("reg_channel.invalid", util.ErrorInvalidData("reg channel"))
	// 	return o
	// }
	// c.Data.RegChannel = glossaryRegChannel.ValueInt

	return
}
