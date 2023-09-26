package service

import (
	"git.edenfarm.id/project-version2/api/src/sales/sales_assignment/unit_test/entity"
	"git.edenfarm.id/project-version2/api/src/sales/sales_assignment/unit_test/repository"
)

type SubmissionListServie struct {
	Repository repository.SubmissionListRepository
}

func (service SubmissionListServie) CheckGetData(task int, finishDate string) *entity.Submission {
	submissionList := service.Repository.CheckSubmission(task, finishDate)

	if submissionList.Validation == "" {
		return nil
	}

	return submissionList
}

type SubmissionDetailService struct {
	Repository repository.SubmissionDetailRepository
}

// CheckGetData : simple service function to get model based on ID
func (s SubmissionDetailService) CheckGetData(id int8) (*entity.SubmissionDetail, error) {
	submissionDetail, e := s.Repository.CheckSubmissionDetail(id)

	if e != nil {
		return nil, e
	}

	return submissionDetail, nil
}
