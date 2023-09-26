package tag

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type archiveRequest struct {
	ID int64 `json:"-" valid:"required"`

	Session *auth.SessionData `json:"-"`
}

func (c *archiveRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	if customerTag, err := repository.ValidCustomerTag(c.ID); err == nil {
		if customerTag.Status != 1 {
			o.Failure("status.active", util.ErrorActive("status"))
		}
	} else {
		o.Failure("id.invalid", util.ErrorInvalidData("customer tag"))
	}

	if countMerchant, err := repository.CountMerchantTagCustomer(c.ID); err == nil {
		if countMerchant > 0 {
			o.Failure("id.invalid", util.ErrorRelated("active or archive ", "customer", "customer tag"))
		}
	}

	return o
}

func (c *archiveRequest) Messages() map[string]string {
	return map[string]string{}
}
