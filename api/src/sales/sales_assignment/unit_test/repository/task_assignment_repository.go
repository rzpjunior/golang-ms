package repository

import "git.edenfarm.id/project-version2/api/src/sales/sales_assignment/unit_test/entity"

type TaskAssignmentRepository interface {
	CheckCountData(id string) *entity.TaskAssignment
}
