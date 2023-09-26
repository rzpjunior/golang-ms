package repository

import (
	"git.edenfarm.id/project-version2/api/src/sales/sales_assignment/unit_test/entity"
	"github.com/stretchr/testify/mock"
)

type SubmissionListRepositoryMock struct {
	Mock mock.Mock
}

func (repository *SubmissionListRepositoryMock) CheckSubmission(task int, finishDate string) *entity.Submission {
	arguments := repository.Mock.Called(task, finishDate)

	if arguments.Get(0) == nil {
		return nil
	} else {
		submission := arguments.Get(0).(entity.Submission)
		return &submission
	}
}
