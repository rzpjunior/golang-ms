package repository

import (
	"git.edenfarm.id/project-version2/api/src/warehouse/picking/unit_test/entity"
	"github.com/stretchr/testify/mock"
)

type PickingAssignRepositoryMock struct {
	Mock mock.Mock
}

func (repository *PickingAssignRepositoryMock) CheckStatusById(id string) *entity.PickingAssign {
	arguments := repository.Mock.Called(id)

	if arguments.Get(0) == nil {
		return nil
	} else {
		pickingAssign := arguments.Get(0).(entity.PickingAssign)
		return &pickingAssign
	}
}

func (repository *PickingAssignRepositoryMock) CheckStatusPrintById(id string) *entity.PickingAssign {
	arguments := repository.Mock.Called(id)

	if arguments.Get(0) == nil {
		return nil
	} else {
		pickingAssign := arguments.Get(0).(entity.PickingAssign)
		return &pickingAssign
	}
}

func (repository *PickingAssignRepositoryMock) CheckPrintState(state string) *entity.PickingAssign {
	arguments := repository.Mock.Called(state)

	if arguments.Get(0) == nil {
		return nil
	} else {
		pickingAssign := arguments.Get(0).(entity.PickingAssign)
		return &pickingAssign
	}
}

func (repository *PickingAssignRepositoryMock) CheckTolerable(id string) *entity.PickingAssign {
	arguments := repository.Mock.Called(id)

	if arguments.Get(0) == nil {
		return nil
	} else {
		pickingAssign := arguments.Get(0).(entity.PickingAssign)
		return &pickingAssign
	}
}

func (repository *PickingAssignRepositoryMock) CheckRefreshedItemSalesOrder(id string) *entity.PickingAssign {
	arguments := repository.Mock.Called(id)

	if arguments.Get(0) == nil {
		return nil
	} else {
		pickingAssign := arguments.Get(0).(entity.PickingAssign)
		return &pickingAssign
	}
}
