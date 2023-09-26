package repository

import (
	"git.edenfarm.id/project-version2/api/src/warehouse/picking/unit_test/entity"
	"github.com/stretchr/testify/mock"
)

type PickingDashboardRepositoryMock struct {
	Mock mock.Mock
}

func (repository *PickingDashboardRepositoryMock) CheckSalesOrderByWrt(wrtId string) *entity.PickingDashboard {
	arguments := repository.Mock.Called(wrtId)

	if arguments.Get(0) == nil {
		return nil
	} else {
		pickingDashboard := arguments.Get(0).(entity.PickingDashboard)
		return &pickingDashboard
	}
}
