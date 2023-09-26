package sales_person

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type archiveRequest struct {
	ID int64 `json:"-" valid:"required"`

	Staff *model.Staff

	Session *auth.SessionData `json:"-"`
}

func (c *archiveRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	c.Staff = &model.Staff{ID: c.ID}
	if err = c.Staff.Read("ID"); err == nil {
		if c.Staff.Status != 1 {
			o.Failure("status.active", util.ErrorActive("status"))
		}
	} else {
		o.Failure("staff.invalid", util.ErrorInvalidData("staff"))
	}

	var count int
	orSelect.Raw("select count(*) from branch b where b.status in (1,2) and b.salesperson_id = ?", c.ID).QueryRow(&count)
	if count > 0 {
		o.Failure("id.invalid", util.ErrorRelated("active/archive", "branch", "user"))
	}

	return o
}

func (c *archiveRequest) Messages() map[string]string {
	return map[string]string{}
}
