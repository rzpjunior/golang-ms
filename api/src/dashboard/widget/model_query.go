package widget

type getDashboardSalesOrderWithWRT struct {
	Wrt        string  `orm:"column(wrt)" json:"wrt"`
	TotalSO    float64 `orm:"column(total_so)" json:"total_so"`
	Active     int     `orm:"column(active)" json:"active"`
	Cancelled  int     `orm:"column(cancelled)" json:"cancelled"`
	OnDelivery int     `orm:"column(on_delivery)" json:"on_delivery"`
	Finished   int     `orm:"column(finished)" json:"finished"`
}

type getGrandTotalDashboardSalesOrderWithWRT struct {
	GrandTotalSO    float64                          `orm:"column(grand_total_so)" json:"grand_total_so"`
	GrandActive     int                              `orm:"column(grand_total_active)" json:"grand_total_active"`
	GrandCancelled  int                              `orm:"column(grand_total_cancelled)" json:"grand_total_cancelled"`
	GrandOnDelivery int                              `orm:"column(grand_total_on_delivery)" json:"grand_total_on_delivery"`
	GrandFinished   int                              `orm:"column(grand_total_finished)" json:"grand_total_finished"`
	TotalRow        int                              `orm:"-" json:"total_row"`
	DashboardSOWRT  []*getDashboardSalesOrderWithWRT `orm:"-" json:"data_dashboard"`
}

type getDashboardTotalPickingOrderWithWRT struct {
	NewPickingStatus          int `orm:"column(new)" json:"new"`
	OnProgressPickingStatus   int `orm:"column(on_progress)" json:"on_progress"`
	NeedApprovalPickingStatus int `orm:"column(need_approval)" json:"need_approval"`
	PickedPickingStatus       int `orm:"column(picked)" json:"picked"`
	CheckingPickingStatus     int `orm:"column(checking)" json:"checking"`
	FinishedPickingStatus     int `orm:"column(finished)" json:"finished"`
	TotalSO                   int `orm:"-" json:"total_so"`
}

type getIdlePickingObj struct {
	Staff           string `orm:"column(staff)" json:"staff"`
	DurationMinutes string `orm:"column(duration_minutes)" json:"-"`
	DurationIdle    string `orm:"column(duration_idle)" json:"duration_idle"`
	Warehouse       string `orm:"column(warehouse)" json:"warehouse"`
}
