package repository

import "git.edenfarm.id/project-version2/api/src/sales/sales_group/unit_test/entity"

type SalesGroupRepository interface {
	CheckSalesGroup(id string) *entity.SalesGroup
	CheckDuplicateNameSalesGroup(id string) *entity.SalesGroup
}
