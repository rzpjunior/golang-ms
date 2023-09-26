package sales_person

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type unarchiveRequest struct {
	ID int64 `json:"-" valid:"required"`

	Staff *model.Staff

	Session *auth.SessionData `json:"-"`
}

func (c *unarchiveRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	c.Staff = &model.Staff{ID: c.ID}
	if err = c.Staff.Read("ID"); err == nil {
		if c.Staff.Status != 2 {
			o.Failure("status.archive", util.ErrorArchived("status"))
		}
	} else {
		o.Failure("staff.invalid", util.ErrorInvalidData("staff"))
	}

	return o
}

func (c *unarchiveRequest) Messages() map[string]string {
	return map[string]string{}
}
