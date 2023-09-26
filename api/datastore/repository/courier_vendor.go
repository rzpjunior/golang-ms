package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"
)

func ValidCourierVendor(id int64) (courierVendor *model.CourierVendor, e error) {
	courierVendor = &model.CourierVendor{ID: id}
	e = courierVendor.Read("ID")

	return
}
