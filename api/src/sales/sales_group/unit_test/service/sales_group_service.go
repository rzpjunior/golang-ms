package service

import (
	"errors"
	"git.edenfarm.id/project-version2/api/src/sales/sales_group/unit_test/entity"
	"git.edenfarm.id/project-version2/api/src/sales/sales_group/unit_test/repository"
)

type SalesGroupService struct {
	Repository repository.SalesGroupRepository
}

func (service SalesGroupService) CheckSalesGroupById(id string) (*entity.SalesGroup, error) {
	salesGroup := service.Repository.CheckSalesGroup(id)
	if salesGroup.SalesGroupStatus != 2 {
		return nil, errors.New("Sales Group is archived")
	} else {
		return salesGroup, nil
	}
}

func (service SalesGroupService) CheckDuplicateNameSalesGroupByName(name string) (*entity.SalesGroup, error) {
	salesGroup := service.Repository.CheckDuplicateNameSalesGroup(name)
	if salesGroup.SalesGroupName == "jakarta1" {
		return nil, errors.New("Sales Group Name is already registered")
	} else {
		return salesGroup, nil
	}
}
