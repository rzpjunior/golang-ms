package repository

import "git.edenfarm.id/project-version2/api/src/warehouse/picking/unit_test/entity"

type PickingAssignRepository interface {
	CheckStatusById(id string) *entity.PickingAssign
	CheckStatusPrintById(id string) *entity.PickingAssign
	CheckPrintState(state string) *entity.PickingAssign
	CheckTolerable(id string) *entity.PickingAssign
	CheckRefreshedItemSalesOrder(id string) *entity.PickingAssign
}
