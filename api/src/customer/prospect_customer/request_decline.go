package prospect_customer

import (
	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type declineRequest struct {
	ID                int64  `json:"id" valid:"required"`
	DeclineTypeID     string `json:"decline_type_id" valid:"required"`
	DeclineNote       string `json:"decline_note"`
	DeclineType       int64  `json:"-"`
	DeclineTypeString string `json:"-"`

	ProspectiveCustomer *model.ProspectCustomer `json:"-"`
	Session             *auth.SessionData       `json:"-"`
	Staff               *model.Staff            `json:"-"`
}

func (c *declineRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	r := orm.NewOrm()
	r.Using("read_only")
	var e error
	if c.ID != 0 {
		if c.ProspectiveCustomer, e = repository.ValidProspectiveCustomer(c.ID); e != nil {
			o.Failure("id.invalid", "Invalid data")
		} else {
			if c.ProspectiveCustomer.RegStatus != 1 {
				o.Failure("id.invalid", util.ErrorActive("prospect customer"))
			}
		}
	}

	if c.DeclineType, e = common.Decrypt(c.DeclineTypeID); e != nil {
		o.Failure("decline_type_id.invalid", util.ErrorInvalidData("decline type"))
	} else {
		g := &model.Glossary{Table: "prospect_customer", Attribute: "decline_type", ValueInt: int8(c.DeclineType)}
		if e = g.Read("Table", "Attribute", "ValueInt"); e != nil {
			o.Failure("decline_type_id.invalid", util.ErrorInvalidData("decline type"))
		}
	}

	if len(c.DeclineNote) > 250 {
		o.Failure("decline_note", util.ErrorCharLength("decline note", 250))
	}

	return o
}

func (c *declineRequest) Messages() map[string]string {
	return map[string]string{
		"decline_type_id.required": util.ErrorInputRequired("decline type"),
	}
}
