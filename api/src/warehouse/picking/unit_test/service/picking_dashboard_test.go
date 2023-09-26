package service

import (
	"git.edenfarm.id/project-version2/api/src/warehouse/picking/unit_test/entity"
	"git.edenfarm.id/project-version2/api/src/warehouse/picking/unit_test/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

var pickingDashboardRepository = &repository.PickingDashboardRepositoryMock{Mock: mock.Mock{}}
var pickingDashboardService = PickingDashboardService{Repository: pickingDashboardRepository}

func TestDashboardPicking(t *testing.T) {
	pickingDashboard := entity.PickingDashboard{
		WrtId:                     "1",
		NewPickingStatus:          10,
		FinishedPickingStatus:     10,
		OnProgressPickingStatus:   10,
		NeedApprovalPickingStatus: 10,
		PickedPickingStatus:       10,
		CheckingPickingStatus:     10,
		TotalSO:                   10,
	}

	pickingDashboardRepository.Mock.On("CheckSalesOrderByWrt", "1").Return(pickingDashboard)

	result, err := pickingDashboardService.CheckSalesOrderByWrtId("1")
	assert.Nil(t, err)
	assert.Equal(t, pickingDashboard.WrtId, result.WrtId)
	assert.Equal(t, pickingDashboard.NewPickingStatus, result.NewPickingStatus)
	assert.Equal(t, pickingDashboard.FinishedPickingStatus, result.FinishedPickingStatus)
	assert.Equal(t, pickingDashboard.OnProgressPickingStatus, result.OnProgressPickingStatus)
	assert.Equal(t, pickingDashboard.NeedApprovalPickingStatus, result.NeedApprovalPickingStatus)
	assert.Equal(t, pickingDashboard.PickedPickingStatus, result.PickedPickingStatus)
	assert.Equal(t, pickingDashboard.CheckingPickingStatus, result.CheckingPickingStatus)
	assert.Equal(t, pickingDashboard.TotalSO, result.TotalSO)
}
