package repository

import "git.edenfarm.id/project-version2/api/src/warehouse/picking/unit_test/entity"

type PickingDashboardRepository interface {
	CheckSalesOrderByWrt(wrtId string) *entity.PickingDashboard
}
