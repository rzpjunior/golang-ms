package field_purchaser

// getPurchasePlanSummary : struct to hold get purchase plan summary query
type getPurchasePlanSummary struct {
	TotalPurchasePlanActive   float64 `orm:"column(total_purchase_plan_active)" json:"total_purchase_plan_active"`
	TotalAssignedPurchasePlan float64 `orm:"column(total_assigned_purchase_plan)" json:"total_assigned_purchase_plan"`
}
