package dto

import (
	"time"
)

type SalesAssignmentObjectiveResponse struct {
	ID            int64              `json:"id"`
	Code          string             `json:"code"`
	Name          string             `json:"name"`
	Objective     string             `json:"objective"`
	SurveyLink    string             `json:"survey_link"`
	Status        int8               `json:"status"`
	CreatedAt     time.Time          `json:"created_at"`
	CreatedBy     *CreatedByResponse `json:"created_by"`
	UpdatedAt     time.Time          `json:"updated_at"`
	StatusConvert string             `json:"status_convert"`
}

type CreatedByResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type SalesAssignmentObjectiveRequestCreate struct {
	Name       string `json:"name" valid:"required"`
	Objective  string `json:"objective" valid:"required"`
	SurveyLink string `json:"survey_link" valid:"required"`
}

type SalesAssignmentObjectiveRequestUpdate struct {
	Objective  string `json:"objective" valid:"required"`
	SurveyLink string `json:"survey_link" valid:"required"`
}
