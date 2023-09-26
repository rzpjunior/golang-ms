package service

import (
	"errors"
	"git.edenfarm.id/project-version2/api/src/warehouse/picking/unit_test/entity"
	"git.edenfarm.id/project-version2/api/src/warehouse/picking/unit_test/repository"
)

type PickingDashboardService struct {
	Repository repository.PickingDashboardRepository
}

func (service PickingDashboardService) CheckSalesOrderByWrtId(id string) (*entity.PickingDashboard, error) {
	pickingDashboard := service.Repository.CheckSalesOrderByWrt(id)
	if pickingDashboard.NewPickingStatus != 10 {
		return nil, errors.New("New status should be 10")
	}
	if pickingDashboard.FinishedPickingStatus != 10 {
		return nil, errors.New("Finished status should be 10")
	}
	if pickingDashboard.OnProgressPickingStatus != 10 {
		return nil, errors.New("On Progress status should be 10")
	}
	if pickingDashboard.NeedApprovalPickingStatus != 10 {
		return nil, errors.New("Need Approval status should be 10")
	}
	if pickingDashboard.PickedPickingStatus != 10 {
		return nil, errors.New("Picked status should be 10")
	}
	if pickingDashboard.CheckingPickingStatus != 10 {
		return nil, errors.New("Checking status should be 10")
	}
	if pickingDashboard.TotalSO != 10 {
		return nil, errors.New("Total Sales Order should be 10")
	} else {
		return pickingDashboard, nil
	}
}
