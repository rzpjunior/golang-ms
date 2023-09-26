package price_set

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/src/auth"
)

type unarchiveRequest struct {
	ID 		int64 `json:"-" valid:"required"`
	Session *auth.SessionData `json:"-"`
}

func (c *unarchiveRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	return o
}

func (c *unarchiveRequest) Messages() map[string]string {
	return map[string]string{}
}
