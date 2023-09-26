package price_set

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type archiveRequest struct {
	ID      int64             `json:"-" valid:"required"`
	Session *auth.SessionData `json:"-"`

	PriceSet *model.PriceSet
}

func (c *archiveRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	if c.PriceSet, err = repository.ValidPriceSet(c.ID); err != nil {
		o.Failure("price_set_id.invalid", util.ErrorInvalidData("price_set"))
	}

	var count int

	orSelect.Raw("select count(*) from branch b where b.price_set_id = ? and b.status in (1,2)", c.PriceSet.ID).QueryRow(&count)

	if count > 0 {
		o.Failure("price_set_id.invalid", util.ErrorRelated("archived/active", "branch", "price set"))
	}

	return o
}

func (c *archiveRequest) Messages() map[string]string {
	return map[string]string{}
}
