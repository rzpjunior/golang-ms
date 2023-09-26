package repository

import (
	"git.edenfarm.id/project-version2/api/src/sales/sales_assignment/unit_test/entity"
)

type SubmissionListRepository interface {
	CheckSubmission(task int, finishDate string) *entity.Submission
}

type SubmissionDetailRepository interface {
	CheckSubmissionDetail(id int8) (*entity.SubmissionDetail, error)
}
