package picking

import "git.edenfarm.id/project-version2/datamodel/model"

type templatePickingOrder struct {
	ID              int64   `orm:"column(id);auto" json:"-"`
	Warehouse       string  `orm:"column(warehouse);null" json:"warehouse"`
	BusinessType    string  `orm:"column(business_type);null" json:"business_type"`
	OrderCode       string  `orm:"column(order_code);null" json:"order_code"`
	OrderType       string  `orm:"column(sales_order_type);null" json:"sales_order_type"`
	MerchantName    string  `orm:"column(merchant_name);null" json:"merchant_name"`
	OrderStatus     string  `orm:"column(sales_order_status);null" json:"sales_order_status"`
	DeliveryCode    string  `orm:"column(delivery_code);null" json:"delivery_code"`
	DeliveryStatus  string  `orm:"column(delivery_status);null" json:"delivery_status"`
	ShippingAddress string  `orm:"column(shipping_address);null" json:"shipping_address"`
	Province        string  `orm:"column(province);null" json:"province"`
	City            string  `orm:"column(city);null" json:"city"`
	District        string  `orm:"column(district);null" json:"district"`
	SubDistrict     string  `orm:"column(sub_district);null" json:"sub_district"`
	PostalCode      string  `orm:"column(postal_code);null" json:"postal_code"`
	Wrt             string  `orm:"column(wrt);null" json:"wrt"`
	OrderWeight     float64 `orm:"column(order_weight);null" json:"order_weight"`
	DeliveryDate    string  `orm:"column(delivery_date);null" json:"delivery_date"`
	PaymentTerm     string  `orm:"column(payment_term);null" json:"payment_term"`
	Picker          string  `orm:"column(picker);null" json:"picker"`
	Vendor          string  `orm:"column(vendor);null" json:"vendor"`
	Planning        string  `orm:"column(planning);null" json:"planning"`
	Courier         string  `orm:"column(courier);null" json:"courier"`
	PickingListCode string  `orm:"column(pl_code);null" json:"pl_code"`
}

type MobileProfile struct {
	TotalAssign        float64 `orm:"column(total_assign)" json:"total_assign"`
	TotalNeedApproval  float64 `orm:"column(total_need_approval)" json:"total_need_approval"`
	New                float64 `orm:"column(new)" json:"new"`
	Finished           float64 `orm:"column(finished)" json:"finished"`
	Picked             float64 `orm:"column(picked)" json:"picked"`
	Checking           float64 `orm:"column(checking)" json:"checking"`
	OnProgress         float64 `orm:"column(on_progress)" json:"on_progress"`
	NeedApproval       float64 `orm:"column(need_approval)" json:"need_approval"`
	OnProgressTask     float64 `orm:"-" json:"on_progress_task"`
	FinishedTask       float64 `orm:"-" json:"finished_task"`
	NeedApprovalTask   float64 `orm:"-" json:"need_approval_task"`
	PickingOrderStatus float64 `orm:"column(status)" json:"status"`
}

type ItemWRTMonitoring struct {
	WRTID              string  `orm:"column(id)" json:"id"`
	WRTName            string  `orm:"column(name)" json:"WRT"`
	PickingOrderID     int64   `orm:"column(id);auto" json:"-"`
	PickingOrderStatus float64 `orm:"column(status)" json:"status"`
	StaffId            float64 `orm:"column(staff_id)" json:"staff_id"`
	CheckedBy          float64 `orm:"column(checked_by)" json:"checked_by"`
	LeadPicker         string  `orm:"column(lead_picker)" json:"lead_picker"`
	Checker            string  `orm:"column(checker)" json:"checker"`

	TotalSalesOrder int64 `orm:"-" json:"total_so"`
	TotalOnProgress int64 `orm:"-" json:"total_on_progress"`
	TotalPicked     int64 `orm:"-" json:"total_picked"`
	TotalChecking   int64 `orm:"-" json:"total_checking"`
	TotalFinished   int64 `orm:"-" json:"total_finished"`
}

type ItemSOMonitoring struct {
	SalesOrder    string  `orm:"column(code)" json:"sales_order"`
	Status        int64   `orm:"column(status)" json:"status"`
	StatusConvert string  `orm:"column(picking_status)" json:"status_convert"`
	Merchant      string  `orm:"column(merchant)" json:"merchant"`
	CustomerTag   string  `orm:"column(customer_tag)" json:"customer_tag"`
	TotalKoli     float64 `orm:"column(total_koli)" json:"total_koli"`
	HelperName    string  `orm:"column(helper_name)" json:"helper_name"`
	HelperCode    string  `orm:"column(helper_code)" json:"helper_code"`
}

type CheckerProfileTemp struct {
	Status    int   `orm:"-" json:"status"`
	CheckedBy int64 `orm:"-" json:"checked_by"`
}

type ItemPickingList struct {
	ItemPickingListProduct []*ItemPickingListProduct `json:"item_picking_list_product"`
	Pickers                []*model.Staff            `json:"pickers"`
}
type ItemPickingListProduct struct {
	Product *model.Product `orm:"-" json:"product"`

	PickingRouting         int8    `json:"picking_routing"`
	PickingListStatus      int8    `orm:"column(picking_list_status)" json:"picking_list_status"`
	ProductID              int64   `orm:"column(id)" json:"-"`
	TotalOrder             float64 `orm:"column(order_qty)" json:"total_order"`
	PickQty                float64 `orm:"column(pick_qty)" json:"pick_qty"`
	FlagDisableSku         int8    `orm:"-" json:"flag_disable_sku"`
	FlagSavedPick          int8    `orm:"column(flag_saved_pick)" json:"flag_saved_pick"`
	FlagRejectedByChecker  int8    `orm:"-" json:"flag_rejected_by_checker"`
	PickingOrderItemStatus int8    `orm:"column(poi_status)" json:"poi_status"`
	SalesOrderStatus       int8    `orm:"column(status_sales_order)" json:"-"`
	PickersString          string  `orm:"column(sub_pickers)" json:"-"`
}
type GroupPickingList struct {
	PickingListID   int64  `orm:"column(id)" json:"picking_list_id"`
	PickingListCode string `orm:"column(code)" json:"picking_list_code"`
	PickingRouting  int    `json:"picking_routing"`
	PickingListNote string `orm:"column(note)" json:"note"`
	TagCustomer     int    `orm:"column(tag_customer)" json:"-"`
	Status          int    `orm:"column(status)" json:"status"`
	DeliveryDate    string `orm:"column(delivery_date)" json:"delivery_date"`
	StatusConvert   string `orm:"-" json:"status_convert"`

	Pickers []*model.Staff `json:"pickers"`

	NewCustomer     int `orm:"-" json:"new_customer"`
	PowerCustomer   int `orm:"-" json:"power_customer"`
	Other           int `orm:"-" json:"other"`
	TotalSalesOrder int `orm:"-" json:"total_sales_order"`
}

type GenerateCodePickingList struct {
	SalesOrderID   int64   `orm:"column(id)" json:"sales_order_id"`
	SalesOrderCode string  `orm:"column(so_code)" json:"sales_order_code"`
	TotalWeight    float64 `orm:"column(so_total)" json:"total_weight"`
	Wrt            string  `orm:"column(wrt)" json:"wrt"`
	ProductName    string  `orm:"column(product_name)" json:"product_name"`
	WeightItem     float64 `orm:"column(weight_item)" json:"weight_item"`
	OrderItem      float64 `orm:"column(order_item)" json:"order_item"`
}

type GroupingSalesOrder struct {
	ProductName          string                      `orm:"column(product)" json:"product_name"`
	SalesOrder           int64                       `orm:"column(sales_order_id)" json:"-"`
	SalesOrderItemNote   string                      `orm:"column(sales_order_item_note)" json:"sales_order_item_note"`
	PickingOrderAssign   int64                       `orm:"column(picking_order_assign_id)" json:"-"`
	Status               int                         `orm:"column(status)" json:"status"`
	StatusSalesOrder     int                         `orm:"column(status_sales_order)" json:"status_sales_order"`
	PickingFlag          int                         `orm:"column(picking_flag)" json:"picking_flag"`
	SalesOrderID         string                      `orm:"-" json:"sales_order_id"`
	PickingOrderAssignID string                      `orm:"-" json:"picking_order_assign_id"`
	SalesOrderCode       string                      `orm:"column(code)" json:"sales_order_code"`
	OrderQty             float64                     `orm:"column(order_qty)" json:"order_qty"`
	PickQty              float64                     `orm:"column(pick_qty)" json:"pick_qty"`
	UnfullfillNote       string                      `orm:"column(unfullfill_note)" json:"unfullfill_note"`
	Uom                  string                      `orm:"column(uom)" json:"uom"`
	MerchantName         string                      `orm:"column(merchant)" json:"merchant"`
	TagCustomer          string                      `orm:"-" json:"tag_customer"`
	TagCustomerDB        string                      `orm:"column(tag_customer)" json:"-"`
	Wrt                  string                      `orm:"column(wrt)" json:"wrt"`
	PackRecommendation   []*model.PackRecommendation `orm:"-" json:"pack_recommendation"`
}

type SalesOrderByPickingList struct {
	PickingListCode string              `orm:"-" json:"picking_list_code"`
	SalesOrders     []*model.SalesOrder `orm:"-" json:"sales_orders"`
}
