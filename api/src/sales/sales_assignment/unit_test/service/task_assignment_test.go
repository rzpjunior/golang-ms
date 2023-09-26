package service

import (
	"git.edenfarm.id/project-version2/api/src/sales/sales_assignment/unit_test/entity"
	"git.edenfarm.id/project-version2/api/src/sales/sales_assignment/unit_test/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

var taskAssignmentRepository = &repository.TaskAssignmentRepositoryMock{Mock: mock.Mock{}}
var taskAssignmentService = TaskAssignmentService{Repository: taskAssignmentRepository}

func TestAssignmentCount(t *testing.T) {
	taskAssignment := entity.TaskAssignment{
		Id:        "7",
		CountData: 1,
	}

	taskAssignmentRepository.Mock.On("CheckCountData", "7").Return(taskAssignment)

	result, _ := taskAssignmentService.CheckCountData("7")
	assert.Equal(t, taskAssignment.Id, result.Id)
	assert.Equal(t, taskAssignment.CountData, result.CountData)
}
