package branch

type templateBranchBySalesperson struct {
	AreaId         int64  `orm:"column(area_id)" json:"area_id"`
	AreaName       string `orm:"column(area_name)" json:"area_name"`
	SalesGroupID   int64  `orm:"column(sales_group_id)" json:"sales_group_id"`
	SalesGroupName string `orm:"column(sales_group_name)" json:"sales_group_name"`
	SalespersonID  int64  `orm:"column(salesperson_id)" json:"salesperson_id"`
	StaffCode      string `orm:"column(staff_code)" json:"staff_code"`
	StaffName      string `orm:"column(staff_name)" json:"staff_name"`
	BranchID       int64  `orm:"column(branch_id)" json:"branch_id"`
	BranchCode     string `orm:"column(branch_code)" json:"branch_code"`
	BranchName     string `orm:"column(branch_name)" json:"branch_name"`
}
