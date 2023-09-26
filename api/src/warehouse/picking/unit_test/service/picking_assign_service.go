package service

import (
	"errors"
	"git.edenfarm.id/project-version2/api/src/warehouse/picking/unit_test/entity"
	"git.edenfarm.id/project-version2/api/src/warehouse/picking/unit_test/repository"
)

type PickingAssignService struct {
	Repository repository.PickingAssignRepository
}

func (service PickingAssignService) CheckStatusApproveById(id string) (*entity.PickingAssign, error) {
	pickingAssign := service.Repository.CheckStatusById(id)
	if pickingAssign.Status != 2 {
		return nil, errors.New("status should be Finished")
	} else {
		return pickingAssign, nil
	}
}

func (service PickingAssignService) CheckStatusRejectById(id string) (*entity.PickingAssign, error) {
	pickingAssign := service.Repository.CheckStatusById(id)
	if pickingAssign.Status != 3 {
		return nil, errors.New("status should be On Progress")
	} else {
		return pickingAssign, nil
	}
}

func (service PickingAssignService) CheckStatusPrintLabelById(id string) (*entity.PickingAssign, error) {
	pickingAssign := service.Repository.CheckStatusPrintById(id)
	if pickingAssign.StatusPrint != 0 {
		return nil, errors.New("status should be label state (0)")
	} else {
		return pickingAssign, nil
	}
}

func (service PickingAssignService) CheckStatusPrintInvoiceById(id string) (*entity.PickingAssign, error) {
	pickingAssign := service.Repository.CheckStatusPrintById(id)
	if pickingAssign.StatusPrint != 1 {
		return nil, errors.New("status should be invoice state (1)")
	} else {
		return pickingAssign, nil
	}
}

func (service PickingAssignService) CheckPickingPrintState(state string) (*entity.PickingAssign, error) {
	pickingAssign := service.Repository.CheckPrintState(state)
	if pickingAssign.PrintState != "label_picking" {
		return nil, errors.New("print state should be print label")
	} else {
		return pickingAssign, nil
	}
}

func (service PickingAssignService) CheckSJPrintState(state string) (*entity.PickingAssign, error) {
	pickingAssign := service.Repository.CheckPrintState(state)
	if pickingAssign.PrintState != "sj" {
		return nil, errors.New("print state should be print sj")
	} else {
		return pickingAssign, nil
	}
}

func (service PickingAssignService) CheckTolerableByID(id string) (*entity.PickingAssign, error) {
	pickingAssign := service.Repository.CheckTolerable(id)
	if pickingAssign.Tolerable != 0.01 {
		return nil, errors.New("tolerable percentage should be 1%")
	} else {
		return pickingAssign, nil
	}
}

func (service PickingAssignService) CheckRefreshedSalesOrderByID(id string) (*entity.PickingAssign, error) {
	pickingAssign := service.Repository.CheckRefreshedItemSalesOrder(id)

	var items []entity.Item
	var item entity.Item

	item.ProductCode = "PRD0001"
	item.ProductCode = "PRD0002"
	items = append(items, item)
	pickingAssign.ProductItem = items

	for _, v := range pickingAssign.ProductItem {
		for _, v2 := range items {
			if v.ProductCode == v2.ProductCode {
				v.FlagOrder = 2
			} else {
				v.FlagOrder = 4
			}
		}
	}

	if pickingAssign.LockedSalesOrder != 1 {
		return nil, errors.New("Sales Order should be locked")
	} else {
		return pickingAssign, nil
	}

	return pickingAssign, nil
}
