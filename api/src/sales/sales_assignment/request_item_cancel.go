package sales_assignment

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type cancelItemRequest struct {
	ID      int64             `json:"-" valid:"required"`
	Session *auth.SessionData `json:"-"`

	SalesAssignmentItem *model.SalesAssignmentItem
}

func (c *cancelItemRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	if c.SalesAssignmentItem, err = repository.ValidSalesAssignmentItem(c.ID); err != nil {
		o.Failure("sales_assignment_item.invalid", util.ErrorInvalidData("sales assignment"))
		return o
	}

	if c.SalesAssignmentItem.Status != 1 {
		o.Failure("sales_assignment_item.inactive", util.ErrorActive("sales assignment"))
		return o
	}

	return o
}

func (c *cancelItemRequest) Messages() map[string]string {
	return map[string]string{}
}
