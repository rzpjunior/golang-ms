package sales_assignment

type templateBranchBySalesGroup struct {
	CustomerType    string `orm:"column(customer_type)" json:"customer_type"`
	SalesGroupID    int64  `orm:"column(sales_group_id)" json:"sales_group_id"`
	SalespersonID   int64  `orm:"column(salesperson_id)" json:"salesperson_id"`
	BranchID        int64  `orm:"column(branch_id)" json:"branch_id"`
	SalesGroupName  string `orm:"column(sales_group_name)" json:"sales_group_name"`
	BranchCode      string `orm:"column(branch_code)" json:"branch_code"`
	OutletName      string `orm:"column(outlet_name)" json:"outlet_name"`
	SubDistrictName string `orm:"column(sub_district_name)" json:"sub_district_name"`
	DistrictName    string `orm:"column(district_name)" json:"district_name"`
	StaffCode       string `orm:"column(staff_code)" json:"staff_code"`
	StaffName       string `orm:"column(staff_name)" json:"staff_name"`
}
