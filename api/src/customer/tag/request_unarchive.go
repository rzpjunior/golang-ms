package tag

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type unarchiveRequest struct {
	ID int64 `json:"-" valid:"required"`

	Items   []int64           `json:"-"`
	Session *auth.SessionData `json:"-"`
}

func (c *unarchiveRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	if customerTag, err := repository.ValidCustomerTag(c.ID); err == nil {
		if customerTag.Status != 2 {
			o.Failure("status.active", util.ErrorArchived("status"))
		}
	} else {
		o.Failure("id.invalid", util.ErrorInvalidData("customer tag"))
	}

	return o
}

func (c *unarchiveRequest) Messages() map[string]string {
	return map[string]string{}
}
