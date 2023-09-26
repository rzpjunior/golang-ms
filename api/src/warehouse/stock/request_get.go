package stock

type IntransitStock struct {
	ProductID    int64   `orm:"column(product_id)" json:"-"`
	WarehouseID  int64   `orm:"column(warehouse_id)" json:"-"`
	IntransitQty float64 `orm:"column(intransit_qty)" json:"-"`
}

type ReceivedStock struct {
	ProductID   int64   `orm:"column(product_id)" json:"-"`
	WarehouseID int64   `orm:"column(warehouse_id)" json:"-"`
	ReceivedQty float64 `orm:"column(received_qty)" json:"-"`
}

type IntransitWasteStock struct {
	ProductID         int64   `orm:"column(product_id)" json:"-"`
	WarehouseID       int64   `orm:"column(warehouse_id)" json:"-"`
	IntransitWasteQty float64 `orm:"column(intransit_Waste_qty)" json:"-"`
}

type ExpectedStockGTPO struct {
	ProductID   int64   `orm:"column(product_id)" json:"-"`
	WarehouseID int64   `orm:"column(warehouse_id)" json:"-"`
	DraftQty    float64 `orm:"column(draft_qty)" json:"-"`
}

type ExpectedStockGTPP struct {
	ProductID   int64   `orm:"column(product_id)" json:"-"`
	WarehouseID int64   `orm:"column(warehouse_id)" json:"-"`
	PurchaseQty float64 `orm:"column(purchase_qty)" json:"-"`
}
