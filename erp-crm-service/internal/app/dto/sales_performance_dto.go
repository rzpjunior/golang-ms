package dto

type SalesPerformanceResponse struct {
	PlanVisit                int                  `json:"plan_visit"`
	PlanFollowUp             int                  `json:"plan_follow_up"`
	VisitActual              int                  `json:"visit_actual"`
	FollowUpActual           int                  `json:"follow_up_actual"`
	VisitPercentage          float64              `json:"visit_percentage"`
	FollowUpPercentage       float64              `json:"follow_up_percentage"`
	EffectiveCall            int                  `json:"effective_call"`
	EffectiveCallPercentage  float64              `json:"effective_call_percentage"`
	RevenueEffectiveCall     float64              `json:"revenue_effective_call"`
	RevenueTotal             float64              `json:"revenue_total"`
	TotalCustomerAcquisition int                  `json:"total_customer_acquisition"`
	Salesperson              *SalespersonResponse `json:"salesperson"`
}

type SalesPerformanceDetailResponse struct {
	VisitTracker              *PerformanceTrackerResponse          `json:"visit_tracker"`
	FollowUpTracker           *PerformanceTrackerResponse          `json:"follow_up_tracker"`
	EffectiveCallPercentage   float64                              `json:"effective_call_percentage"`
	SalesAssignmentSubmission []*SalesAssignmentSubmissionResponse `json:"sales_assignment_submissions"`
	CustomerAcquisition       []*CustomerAcquisitionResponse       `json:"customer_acquisitions"`
}

type PerformanceTrackerResponse struct {
	TotalPlan       int `json:"total_plan"`
	TotalFinished   int `json:"total_finished"`
	TotalFailed     int `json:"total_failed"`
	TotalCancelled  int `json:"total_cancelled"`
	TotalOutOfRoute int `json:"total_out_of_route"`
}
