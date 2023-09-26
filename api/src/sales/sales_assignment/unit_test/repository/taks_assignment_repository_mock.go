package repository

import (
	"git.edenfarm.id/project-version2/api/src/sales/sales_assignment/unit_test/entity"
	"github.com/stretchr/testify/mock"
)

type TaskAssignmentRepositoryMock struct {
	Mock mock.Mock
}

func (repository *TaskAssignmentRepositoryMock) CheckCountData(id string) *entity.TaskAssignment {
	arguments := repository.Mock.Called(id)

	if arguments.Get(0) == nil {
		return nil
	} else {
		taskAssignment := arguments.Get(0).(entity.TaskAssignment)
		return &taskAssignment
	}
}
