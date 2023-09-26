package repository

import (
	"git.edenfarm.id/project-version2/api/src/sales/sales_group/unit_test/entity"
	"github.com/stretchr/testify/mock"
)

type SalesGroupRepositoryMock struct {
	Mock mock.Mock
}

func (repository *SalesGroupRepositoryMock) CheckSalesGroup(id string) *entity.SalesGroup {
	arguments := repository.Mock.Called(id)

	if arguments.Get(0) == nil {
		return nil
	} else {
		salesGroup := arguments.Get(0).(entity.SalesGroup)
		return &salesGroup
	}
}

func (repository *SalesGroupRepositoryMock) CheckDuplicateNameSalesGroup(name string) *entity.SalesGroup {
	arguments := repository.Mock.Called(name)

	if arguments.Get(0) == nil {
		return nil
	} else {
		salesGroup := arguments.Get(0).(entity.SalesGroup)
		return &salesGroup
	}
}
