package service

import (
	"errors"

	"git.edenfarm.id/project-version2/api/src/sales/sales_assignment/unit_test/entity"
	"git.edenfarm.id/project-version2/api/src/sales/sales_assignment/unit_test/repository"
)

type TaskAssignmentService struct {
	Repository repository.TaskAssignmentRepository
}

func (service TaskAssignmentService) CheckCountData(id string) (*entity.TaskAssignment, error) {
	taskAssignment := service.Repository.CheckCountData(id)
	if taskAssignment.CountData == 0 {
		return nil, errors.New("There is no data for this sales group")
	} else {
		return taskAssignment, nil
	}
}
