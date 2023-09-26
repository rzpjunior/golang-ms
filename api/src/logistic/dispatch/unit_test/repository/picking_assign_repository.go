package repository

import "git.edenfarm.id/project-version2/api/src/logistic/dispatch/unit_test/entity"

type DispatchRepository interface {
	CheckCourierVendorCombination(id string) *entity.Dispacth
	CheckListDispatch(id string) *entity.Dispacth
}
