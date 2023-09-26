package fridge

type reportSoldProductFridge struct {
	UOMName       string `orm:"column(uom_name)" json:"uom_name"`
	TotalWeight   string `orm:"column(total_weight)" json:"total_weight"`
	ProductName   string `orm:"column(product_name)" json:"product_name"`
	BranchName    string `orm:"column(branch_name)" json:"branch_name"`
	MerchantName  string `orm:"column(merchant_name)" json:"merchant_name"`
	SoldDate      string `orm:"column(last_seen_at)" json:"last_seen_at"`
	WarehouseName string `orm:"column(warehouse_name)" json:"warehouse_name"`
	AreaName      string `orm:"column(area_name)" json:"area_name"`
	CreatedAt     string `orm:"column(created_at)" json:"created_at"`
}

type reportAllProductFridge struct {
	UOMName         string `orm:"column(uom_name)" json:"uom_name"`
	TotalWeight     string `orm:"column(total_weight)" json:"total_weight"`
	ProductName     string `orm:"column(product_name)" json:"product_name"`
	BranchName      string `orm:"column(branch_name)" json:"branch_name"`
	MerchantName    string `orm:"column(merchant_name)" json:"merchant_name"`
	ProcessedDate   string `orm:"column(last_seen_at)" json:"last_seen_at"`
	WarehouseName   string `orm:"column(warehouse_name)" json:"warehouse_name"`
	AreaName        string `orm:"column(area_name)" json:"area_name"`
	CreatedAt       string `orm:"column(created_at)" json:"created_at"`
	FinishedAt      string `orm:"column(finished_at)" json:"finished_at"`
	BoxItemStatus   int    `orm:"column(box_item_status)" json:"box_item_status"`
	BoxFridgeStatus int    `orm:"column(box_fridge_status)" json:"box_fridge_status"`
	Status          string `orm:"column(status)" json:"status"`
	ImageURL        string `orm:"column(image_url)" json:"image_url"`
	ProductCode     string `orm:"column(product_code)" json:"product_code"`
	UnitPrice       string `orm:"column(unit_price)" json:"unit_price"`
	TotalPrice      string `orm:"column(total_price)" json:"total_price"`
}
