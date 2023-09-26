package sales_assignment

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type cancelRequest struct {
	ID      int64             `json:"-" valid:"required"`
	Session *auth.SessionData `json:"-"`

	SalesAssignment     *model.SalesAssignment
	SalesAssignmentItem []*model.SalesAssignmentItem
}

func (c *cancelRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error
	q := orm.NewOrm()
	q.Using("read_only")

	if c.SalesAssignment, err = repository.ValidSalesAssignment(c.ID); err != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("sales assignment"))
		return o
	}

	if c.SalesAssignment.Status != 1 {
		o.Failure("id.inactive", util.ErrorActive("sales assignment"))
		return o
	}

	// get sales assignment item related to this sales assignment
	q.Raw("SELECT * FROM sales_assignment_item sai WHERE sai.sales_assignment_id = ? AND sai.status = 1", c.SalesAssignment.ID).QueryRows(&c.SalesAssignmentItem)

	return o
}

func (c *cancelRequest) Messages() map[string]string {
	return map[string]string{}
}
