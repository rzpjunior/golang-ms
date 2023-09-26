package service

import (
	"errors"
	"git.edenfarm.id/project-version2/api/src/logistic/dispatch/unit_test/entity"
	"git.edenfarm.id/project-version2/api/src/logistic/dispatch/unit_test/repository"
)

type DispatchService struct {
	Repository repository.DispatchRepository
}

func (service DispatchService) CheckVendorCourierCombinationById(id string) (*entity.Dispacth, error) {
	dispatch := service.Repository.CheckCourierVendorCombination(id)
	if dispatch.CourierName != "Doni" {
		return nil, errors.New("Courier name should be Doni")
	} else if dispatch.VendorName != "Eden Farm" {
		return nil, errors.New("Vendor name should be Eden Farm")
	} else {
		return dispatch, nil
	}
}

func (service DispatchService) CheckListDispatchById(id string) (*entity.Dispacth, error) {
	dispatch := service.Repository.CheckListDispatch(id)
	if dispatch.DispatchStatus != 1 {
		return nil, errors.New("Dispatch status should be 1(New)")
	} else if dispatch.PickingStatus != 2 {
		return nil, errors.New("Picking status should be 2(Finished)")
	} else if dispatch.SalesOrderCode != "SO-BASG9867-12210019" {
		return nil, errors.New("Invalid sales order")
	} else {
		return dispatch, nil
	}
}
