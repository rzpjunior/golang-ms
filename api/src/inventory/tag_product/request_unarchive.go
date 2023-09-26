package tag_product

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type unarchiveRequest struct {
	ID      int64             `json:"-"`
	Session *auth.SessionData `json:"-"`

	TagProduct *model.TagProduct
}

func (c *unarchiveRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	if c.TagProduct, err = repository.ValidTagProduct(c.ID); err == nil {
		if c.TagProduct.Status != 2 {
			o.Failure("id.active", util.ErrorArchived("status"))
		}
	} else {
		o.Failure("id.invalid", util.ErrorInvalidData("tag product"))
	}

	return o
}

func (r *unarchiveRequest) Messages() map[string]string {
	return map[string]string{}
}
