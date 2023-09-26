package service

import (
	"errors"
	"testing"

	"git.edenfarm.id/project-version2/api/src/sales/sales_assignment/unit_test/entity"
	"git.edenfarm.id/project-version2/api/src/sales/sales_assignment/unit_test/repository"
	"git.edenfarm.id/project-version2/api/src/sales/sales_assignment/unit_test/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var submissionListRepo = &repository.SubmissionListRepositoryMock{Mock: mock.Mock{}}
var submissionService = SubmissionListServie{Repository: submissionListRepo}

func TestSubmissionList(t *testing.T) {
	submissionList := entity.Submission{
		Task:        1,
		Finish_date: "2022-03-01",
		Validation:  "this is task one",
	}

	submissionListRepo.Mock.On("CheckSubmission", 1, "2022-03-01").Return(submissionList)

	result := submissionService.CheckGetData(1, "2022-03-01")
	if !assert.Equal(t, result.Validation, submissionList.Validation) {
		t.Errorf("result %v", result.Validation)
	}
}

func TestSubmissionDetail(t *testing.T) {
	const trueID = int8(1)
	const falseID = int8(10)
	submissionDetail := entity.SubmissionDetail{
		Id:         1,
		Task:       1,
		Validation: "this is task for id one",
	}
	submissionDetailRes := entity.SubmissionDetail{
		Id:         1,
		Task:       1,
		Validation: "this is task for id one",
	}

	type test = struct {
		name    string
		mock    func(mock *mocks.SubmissionDetailRepository)
		wantRes *entity.SubmissionDetail
		wantErr bool
	}

	// first test
	tests := []test{
		{
			name: "given id will return detail submission",
			mock: func(mock *mocks.SubmissionDetailRepository) {
				mock.On("CheckSubmissionDetail", trueID).Return(&submissionDetail, nil)
			},
			wantRes: &submissionDetailRes,
			wantErr: false,
		}, {
			name: "raise error if something wrong",
			mock: func(mock *mocks.SubmissionDetailRepository) {
				mock.On("CheckSubmissionDetail", trueID).Return(nil, errors.New("err"))
			},
			wantRes: nil,
			wantErr: true,
		}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newMock := new(mocks.SubmissionDetailRepository)
			s := &SubmissionDetailService{Repository: newMock}

			tt.mock(newMock)
			got, e := s.CheckGetData(trueID)
			if (e != nil) != tt.wantErr {
				t.Errorf("false err=%v, want err=%v", e, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.wantRes, got) {
				t.Errorf("false res=%v, want res=%v", got, tt.wantRes)
			}
		})
	}

	// second test
	tests = []test{
		{
			name: "given false id will return error",
			mock: func(mock *mocks.SubmissionDetailRepository) {
				mock.On("CheckSubmissionDetail", falseID).Return(nil, errors.New("err"))
			},
			wantRes: nil,
			wantErr: true,
		}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newMock := new(mocks.SubmissionDetailRepository)
			s := &SubmissionDetailService{Repository: newMock}

			tt.mock(newMock)
			got, e := s.CheckGetData(falseID)
			if (e != nil) != tt.wantErr {
				t.Errorf("false err=%v, want err=%v", got, e)
				return
			}
		})
	}
}
