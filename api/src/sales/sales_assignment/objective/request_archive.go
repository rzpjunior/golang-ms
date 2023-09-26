package objective

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

type archiveRequest struct {
	ID      int64             `json:"-" valid:"required"`
	Session *auth.SessionData `json:"-"`

	SalesAssignmentObjective *model.SalesAssignmentObjective
}

func (c *archiveRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	if c.SalesAssignmentObjective, err = repository.ValidSalesAssignmentObjective(c.ID); err != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("id"))
	}

	if c.SalesAssignmentObjective.Status != 1 {
		o.Failure("id.invalid", util.ErrorActive("id"))
	}

	return o
}

func (c *archiveRequest) Messages() map[string]string {
	return map[string]string{
		"id.required": util.ErrorInputRequired("id"),
	}
}
