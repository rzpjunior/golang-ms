package repository

import (
	"git.edenfarm.id/project-version2/api/src/logistic/dispatch/unit_test/entity"
	"github.com/stretchr/testify/mock"
)

type DispatchRepositoryMock struct {
	Mock mock.Mock
}

func (repository *DispatchRepositoryMock) CheckCourierVendorCombination(id string) *entity.Dispacth {
	arguments := repository.Mock.Called(id)

	if arguments.Get(0) == nil {
		return nil
	} else {
		dispatch := arguments.Get(0).(entity.Dispacth)
		return &dispatch
	}
}

func (repository *DispatchRepositoryMock) CheckListDispatch(id string) *entity.Dispacth {
	arguments := repository.Mock.Called(id)

	if arguments.Get(0) == nil {
		return nil
	} else {
		dispatch := arguments.Get(0).(entity.Dispacth)
		return &dispatch
	}
}
