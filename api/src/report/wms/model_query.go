package wms

import "time"

type reportStockLog struct {
	TimeStamp     string  `orm:"column(Timestamp)" json:"timestamp"`
	LogType       string  `orm:"column(Log_Type)" json:"log_type"`
	RefType       string  `orm:"column(Ref_Type)" json:"ref_type"`
	InitialStock  float64 `orm:"column(Initial_Stock)" json:"initial_stock"`
	ReferenceCode string  `orm:"column(Reference_Code)" json:"reference_code"`
	ProductCode   string  `orm:"column(Product_Code)" json:"product_code"`
	ProductName   string  `orm:"column(Product_Name)" json:"product_name"`
	Uom           string  `orm:"column(UOM)" json:"uom"`
	Quantity      float64 `orm:"column(Quantity)" json:"quantity"`
	FinalStock    float64 `orm:"column(Final_Stock)" json:"final_stock"`
	Warehouse     string  `orm:"column(Warehouse)" json:"warehouse"`
	Area          string  `orm:"column(Area)" json:"area"`
	Status        string  `orm:"column(Status)" json:"status"`
	DocNote       string  `orm:"column(Doc_Note)" json:"doc_note"`
	Note          string  `orm:"column(Note)" json:"item_note"`
}

type reportWasteLog struct {
	TimeStamp         string  `orm:"column(Timestamp)" json:"timestamp"`
	LogType           string  `orm:"column(Log_Type)" json:"log_type"`
	RefType           string  `orm:"column(Ref_Type)" json:"ref_type"`
	ReferenceCode     string  `orm:"column(Reference_Code)" json:"reference_code"`
	WasteReason       string  `orm:"column(Waste_Reason)" json:"Waste_Reason"`
	GoodReceiptCode   string  `orm:"column(Good_Receipt_Code)" json:"Good_Receipt_Code"`
	GoodTransferCode  string  `orm:"column(Good_Transfer_Code)" json:"Good_Transfer_Code"`
	PurchaseOrderCode string  `orm:"column(Purchase_Order_Code)" json:"Purchase_Order_Code"`
	SuplierName       string  `orm:"column(Suplier_Name)" json:"Suplier_Name"`
	SuplierType       string  `orm:"column(Suplier_Type)" json:"Suplier_Type"`
	WarehouseOrigin   string  `orm:"column(Warehouse_Origin)" json:"Warehouse_Origin"`
	ProductCode       string  `orm:"column(Product_Code)" json:"product_code"`
	ProductName       string  `orm:"column(Product_Name)" json:"product_name"`
	Uom               string  `orm:"column(UOM)" json:"uom"`
	Quantity          float64 `orm:"column(Quantity)" json:"quantity"`
	FinalStock        float64 `orm:"column(Final_Stock)" json:"final_stock"`
	Warehouse         string  `orm:"column(Warehouse)" json:"warehouse"`
	Area              string  `orm:"column(Area)" json:"area"`
	DocNote           string  `orm:"column(Doc_Note)" json:"doc_note"`
	ItemNote          string  `orm:"column(Item_Note)" json:"item_note"`
}

type reportStock struct {
	ProductCode         string  `orm:"column(product_code)" json:"product_code"`
	ProductName         string  `orm:"column(product_name)" json:"product_name"`
	WarehouseName       string  `orm:"column(warehouse_name)" json:"warehouse_name"`
	AvailableStock      float64 `orm:"column(available_stock)" json:"available_stock"`
	WasteStock          float64 `orm:"column(waste_stock)" json:"waste_stock"`
	SafetyStock         float64 `orm:"column(safety_stock)" json:"safety_stock"`
	CommitedInStock     float64 `orm:"column(commited_in_stock)" json:"commited_in_stock"`
	CommitedOutStock    float64 `orm:"column(commited_out_stock)" json:"commited_out_stock"`
	ExpectedStock       float64 `orm:"column(expected_qty)" json:"expected_qty"`
	IntransitStock      float64 `orm:"column(intransit_qty)" json:"intransit_qty"`
	ReceivedStock       float64 `orm:"column(received_qty)" json:"received_qty"`
	IntransitWasteStock float64 `orm:"column(intransit_waste_qty)" json:"intransit_waste_qty"`
	Salable             string  `orm:"column(salable)" json:"salable"`
	Purchasable         string  `orm:"column(purchasable)" json:"purchasable"`
	Status              string  `orm:"column(status)" json:"status"`
}

type reportGoodsReceiptItem struct {
	InboundCode             string  `orm:"column(Inbound_Code)" json:"inbound_code"`
	SupplierCode            string  `orm:"column(Supplier_Code)" json:"supplier_code"`
	SupplierName            string  `orm:"column(Supplier_Name)" json:"supplier_name"`
	InboundStatus           string  `orm:"column(Inbound_Status)" json:"inbound_status"`
	GRCode                  string  `orm:"column(GR_Code)" json:"gr_code"`
	GRStatus                string  `orm:"column(GR_Status)" json:"goods_receipt_status"`
	SupplierReturnCode      string  `orm:"column(SR_Code)" json:"supplier_return_code"`
	SupplierReturnStatus    string  `orm:"column(SR_Status)" json:"supplier_return_status"`
	DebitNoteCode           string  `orm:"column(DN_Code)" json:"debit_note_code"`
	DebitNoteStatus         string  `orm:"column(DN_Status)" json:"debit_note_status"`
	WarehousePurchase       string  `orm:"column(Warehouse_Purchase)" json:"warehouse_purchase"`
	WarehouseOrigin         string  `orm:"column(Warehouse_Origin)" json:"warehouse_origin"`
	WarehouseDestination    string  `orm:"column(Warehouse_Destination)" json:"warehouse_destination"`
	EstimationArrivalDate   string  `orm:"column(Estimation_Arrival_Date)" json:"estimation_arrival_date"`
	EstimationArrivalTime   string  `orm:"column(Estimation_Arrival_Time)" json:"estimation_arrival_time"`
	ActualArrivalDate       string  `orm:"column(Actual_Arrival_Date)" json:"actual_arrival_date"`
	ActualArrivalTime       string  `orm:"column(Actual_Arrival_Time)" json:"actual_arrival_time"`
	GRNote                  string  `orm:"column(GR_Note)" json:"gr_note"`
	ProductCode             string  `orm:"column(Product_Code)" json:"product_code"`
	ProductName             string  `orm:"column(Product_Name)" json:"product_name"`
	UOM                     string  `orm:"column(UOM)" json:"uom"`
	GRItemNote              string  `orm:"column(GR_Item_Note)" json:"gr_item_note"`
	OrderedQty              float64 `orm:"column(Ordered_Qty)" json:"ordered_qty"`
	DeliveredQty            float64 `orm:"column(Delivered_Qty)" json:"delivered_qty"`
	RejectQty               float64 `orm:"column(Reject_Qty)" json:"reject_qty"`
	ReceivedQty             float64 `orm:"column(Received_Qty)" json:"received_qty"`
	SRIReturnQuantity       float64 `orm:"column(Return_Qty)" json:"return_qty"`
	AfterSortirGoodQty      float64 `orm:"column(After_Sortir_Good_Qty)" json:"good_qty"`
	AfterSortirWasteQty     float64 `orm:"column(After_Sortir_Waste_Qty)" json:"waste_qty"`
	AfterSortirDownGradeQty float64 `orm:"column(After_Sortir_DownGrade_Qty)" json:"down_grade_qty"`
	Product_DownGrade       string  `orm:"column(Product_DownGrade)" json:"product_downgrade"`
	TS_Code                 string  `orm:"column(TS_Code)" json:"ts_code"`
	TS_Status               string  `orm:"column(TS_Status)" json:"ts_status"`
}

type reportDeliveryReturnItem struct {
	ReturnDate             string  `orm:"column(Return_Date)" json:"return_date"`
	ProductCode            string  `orm:"column(Product_Code)" json:"product_code"`
	ProductName            string  `orm:"column(Product_Name)" json:"product_name"`
	Unit                   string  `orm:"column(Unit)" json:"unit"`
	GoodStockReturnQty     float64 `orm:"column(Good_Stock_Return_Qty)" json:"good_stock_return_qty"`
	WasteReturnQty         float64 `orm:"column(Waste_Return_Qty)" json:"waste_return_qty"`
	TotalReturnQty         float64 `orm:"column(Total_Return_Qty)" json:"total_return_qty"`
	ProductPrice           float64 `orm:"column(Product_Price)" json:"product_price"`
	Area                   string  `orm:"column(Area)" json:"area"`
	Warehouse              string  `orm:"column(Warehouse)" json:"warehouse"`
	OrderCode              string  `orm:"column(Order_Code)" json:"order_code"`
	DeliveryCode           string  `orm:"column(Delivery_Code)" json:"delivery_code"`
	DeliveryDate           string  `orm:"column(Delivery_Date)" json:"delivery_date"`
	CustomerCode           string  `orm:"column(Customer_Code)" json:"customer_code"`
	CustomerName           string  `orm:"column(Customer_Name)" json:"customer_name"`
	DeliveryReturnNote     string  `orm:"column(Delivery_Return_Note)" json:"delivery_return_note"`
	DeliveryReturnItemNote string  `orm:"column(Delivery_Return_Item_Note)" json:"delivery_return_item_note"`
}

type reportProducts struct {
	ProductCode             string  `orm:"column(product_code)" json:"product_code"`
	ProductName             string  `orm:"column(product_name)" json:"product_name"`
	Category                string  `orm:"column(Category - C2)" json:"category"`
	CategoryCode            string  `orm:"column(Code category - C2)" json:"category_code"`
	Parent                  string  `orm:"column(Category - C1)" json:"parent"`
	ParentCode              string  `orm:"column(Code category - C1)" json:"parent_code"`
	GrandParent             string  `orm:"column(Category - C0)" json:"grand_parent"`
	GrandParentCode         string  `orm:"column(Code category - C0)" json:"grand_parent_code"`
	UOM                     string  `orm:"column(UOM)" json:"uom"`
	TotalWeight             float64 `orm:"column(total_weight)" json:"total_weight"`
	MinimalOrderQty         float64 `orm:"column(Minimal_Order_Qty)" json:"minimal_order_qty"`
	ProductNote             string  `orm:"column(product_note)" json:"product_note"`
	ProductDescription      string  `orm:"column(product_description)" json:"product_description"`
	ProductTag              string  `orm:"column(product_tag)" json:"product_tag"`
	ProductStatus           string  `orm:"column(product_status)" json:"product_status"`
	WarehouseSalability     string  `orm:"column(warehouse_salability)" json:"warehouse_salability"`
	WarehousePurchasability string  `orm:"column(warehouse_purchasability)" json:"warehouse_purchasability"`
	WarehouseStorability    string  `orm:"column(warehouse_storability)" json:"warehouse_storability"`
	SparePercentage         float64 `orm:"column(spare_percentage)" json:"spare_percentage"`
}

type reportDeliveryOrder struct {
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
	TagCustomer     string  `orm:"column(tag_customer);null" json:"tag_customer"`
	AreaName        string  `orm:"column(area_name);null" json:"area_name"`
}

type reportItemRecap struct {
	DeliveryDate      string  `orm:"column(order_delivery_date);null" json:"order_delivery_date"`
	Area              string  `orm:"column(area);null" json:"area"`
	Warehouse         string  `orm:"column(warehouse);null" json:"warehouse"`
	ProductCode       string  `orm:"column(product_code);null" json:"product_code"`
	ProductName       string  `orm:"column(product_name);null" json:"product_name"`
	Uom               string  `orm:"column(uom);null" json:"uom"`
	Category          string  `orm:"column(Category - C2);null" json:"category"`
	CategoryCode      string  `orm:"column(Code category - C2);null" json:"category_code"`
	Parent            string  `orm:"column(Category - C1);null" json:"parent"`
	ParentCode        string  `orm:"column(Code category - C1);null" json:"parent_code"`
	GrandParent       string  `orm:"column(Category - C0);null" json:"grand_parent"`
	GrandParentCode   string  `orm:"column(Code category - C0);null" json:"grand_parent_code"`
	TotalQty          float64 `orm:"column(total_qty);null" json:"total_qty"`
	TotalQtyZeroWaste float64 `orm:"column(total_quantity_zero_waste);null" json:"total_quantity_zero_waste"`
	TotalWeight       float64 `orm:"column(total_weight);null" json:"total_weight"`
}

type reportMovementStock struct {
	ProductCode       string  `orm:"column(product_code);null" json:"product_code"`
	ProductName       string  `orm:"column(product_name);null" json:"product_name"`
	Uom               string  `orm:"column(uom);null" json:"uom"`
	Category          string  `orm:"column(category);null" json:"category"`
	Stock             float64 `orm:"column(stock);null" json:"stock"`
	PlanInbound       float64 `orm:"column(plan_inbound);null" json:"plan_inbound"`
	ActualInbound     float64 `orm:"column(actual_inbound);null" json:"actual_inbound"`
	Waste             float64 `orm:"column(waste);null" json:"waste"`
	PlanDelivery      float64 `orm:"column(plan_delivery);null" json:"plan_delivery"`
	ActualDelivery    float64 `orm:"column(actual_delivery);null" json:"actual_delivery"`
	StockTransferIn   float64 `orm:"column(stock_transfer_in);null" json:"stock_transfer_in"`
	StockTransferOut  float64 `orm:"column(stock_transfer_out);null" json:"stock_transfer_out"`
	GoodsReturn       float64 `orm:"column(goods_return);null" json:"goods_return"`
	FinalStock        float64 `orm:"column(stock_akhir);null" json:"stock_akhir"`
	ActualStock       float64 `orm:"column(actual_stock);null" json:"actual_stock"`
	StockDifferential float64 `orm:"column(selisih_stock);null" json:"selisih_stock"`
}

type reportPicking struct {
	DeliveryDate          string  `orm:"column(delivery_date);null" json:"delivery_date"`
	TimestampAssign       string  `orm:"column(timestamp assign);null" json:"timestamp assign"`
	PickingListCode       string  `orm:"column(pl_code);null" json:"picking_list_code"`
	SalesOrderCode        string  `orm:"column(so code);null" json:"so code"`
	SalesOrderType        string  `orm:"column(order_type);null" json:"order_type"`
	PaymentTerm           string  `orm:"column(payment_term);null" json:"payment_term"`
	Merchant              string  `orm:"column(merchant);null" json:"merchant"`
	BusinessType          string  `orm:"column(business_type);null" json:"business_type"`
	Item                  float64 `orm:"column(total item);null" json:"total item"`
	SalesOrderWeight      float64 `orm:"column(sales_order_weight);null" json:"sales_order_weight"`
	TotalWeight           float64 `orm:"column(total_weight);null" json:"total_weight"`
	TotalKoli             float64 `orm:"column(total_koli);null" json:"total_koli"`
	ShippingAddress       string  `orm:"column(shipping_address);null" json:"shipping_address"`
	Wrt                   string  `orm:"column(wrt);null" json:"wrt"`
	Warehouse             string  `orm:"column(warehouse);null" json:"warehouse"`
	Picker                string  `orm:"column(Picker);null" json:"Picker"`
	TimeStartPicked       string  `orm:"column(Time Start Picked);null" json:"Time Start Picked"`
	TimeFinishPicked      string  `orm:"column(Time Finish Picked);null" json:"Time Finish Picked"`
	Checker               string  `orm:"column(Checker);null" json:"Checker"`
	TimeCheckin           string  `orm:"column(Time Checkin);null" json:"Time Checkin"`
	TimeCheckout          string  `orm:"column(Time Checkout);null" json:"Time Checkout"`
	Vendor                string  `orm:"column(vendor);null" json:"vendor"`
	Planning              string  `orm:"column(planning);null" json:"planning"`
	Courier               string  `orm:"column(courier);null" json:"courier"`
	DispatchTime          string  `orm:"column(Dispatch Time);null" json:"Dispatch Time"`
	StatusPickingAssigned string  `orm:"column(status picking assigned);null" json:"status picking assigned"`
	StatusSalesOrder      string  `orm:"column(status so);null" json:"status so"`
}

type reportPickingOrderItem struct {
	DeliveryDate     string  `orm:"column(delivery_date);null" json:"delivery_date"`
	PickingListCode  string  `orm:"column(pl_code);null" json:"picking_list_code"`
	SalesOrderCode   string  `orm:"column(so code);null" json:"so code"`
	Merchant         string  `orm:"column(merchant);null" json:"merchant"`
	ProductCode      string  `orm:"column(product_code);null" json:"product_code"`
	ProductName      string  `orm:"column(product_name);null" json:"product_name"`
	Uom              string  `orm:"column(uom);null" json:"uom"`
	OrderQty         float64 `orm:"column(order_qty);null" json:"order_qty"`
	QtyPicker        float64 `orm:"column(qty_picker);null" json:"qty_picker"`
	QtyChecker       float64 `orm:"column(qty_checker);null" json:"qty_checker"`
	Wrt              string  `orm:"column(wrt);null" json:"wrt"`
	Warehouse        string  `orm:"column(warehouse);null" json:"warehouse"`
	UnfullfilledNote string  `orm:"column(unfullfill_note);null" json:"unfullfill_note"`
}

type reportGoodsTransferItem struct {
	Timestamp            string  `orm:"column(Timestamp);null" json:"timestamp"`
	GoodsTransferCode    string  `orm:"column(Goods Transfer Code);null" json:"goods_transfer_code"`
	ProductCode          string  `orm:"column(Product_Code);null" json:"product_code"`
	ProductName          string  `orm:"column(Product_Name);null" json:"product_name"`
	Uom                  string  `orm:"column(UOM);null" json:"uom"`
	WarehouseOrigin      string  `orm:"column(Warehouse Origin);null" json:"warehouse_origin"`
	WarehouseDestination string  `orm:"column(Warehouse Destination);null" json:"warehouse_destination"`
	RequestQty           float64 `orm:"column(Request Qty);null" json:"request_qty"`
	TransferQty          float64 `orm:"column(Transfer Qty);null" json:"transfer_qty"`
	ReceivedQty          float64 `orm:"column(Received Qty);null" json:"received_qty"`
	Status               string  `orm:"column(Status);null" json:"status"`
	DocNote              string  `orm:"column(Doc_Note);null" json:"doc_note"`
	Note                 string  `orm:"column(Note);null" json:"Note"`
}

type reportPickingRouting struct {
	SalesOrderCode          string    `orm:"column(sales_order_code);null" json:"sales_order_code"`
	PickingListCode         string    `orm:"column(picking_list_code);null" json:"picking_list_code"`
	LeadPickerID            int64     `orm:"column(lead_picker_id);null"`
	PickerID                int64     `orm:"column(picker_id);null"`
	Product                 string    `orm:"column(product_name);null" json:"product"`
	UOM                     string    `orm:"column(UOM);null" json:"UOM"`
	OrderQty                float64   `orm:"column(order_qty);null" json:"order_qty"`
	RackName                string    `orm:"column(rack_name);null" json:"rack_name"`
	StepType                int64     `orm:"column(step_type);null"`
	Sequence                int64     `orm:"column(sequence);null" json:"sequence"`
	ExpectedWalkingDuration int64     `orm:'column(expected_walking_duration);null" json:"expected_walking_duration"`
	ExpectedServiceDuration int64     `orm:'column(expected_service_duration);null" json:"expected_service_duration"`
	WalkingStartTime        time.Time `orm:'column(walking_start_time);null"`
	WalkingFinishTime       time.Time `orm:'column(walking_finish_time);null"`
	PickingStartTime        time.Time `orm:'column(picking_start_time);null"`
	PickingFinishTime       time.Time `orm:'column(picking_finish_time);null"`
	Status                  int64     `orm:"column(status);null"`

	StepTypeStr           string `json:"step_type"`
	LeadPicker            string `json:"lead_picker_name"`
	Picker                string `json:"picker_name"`
	ActualWalkingDuration int64  `json:"actual_walking_duration"`
	ActualPickingDuration int64  `json:"actual_picking_duration"`
	StatusStr             string `json:"status"`
}

type reportTransferSkuItem struct {
	TransferSkuCode            string  `orm:"column(TS_Code);null" json:"ts_code"`
	TransferSkuStatus          string  `orm:"column(TS_Status);null" json:"ts_status"`
	RecognitionDate            string  `orm:"column(Recognition_Date);null" json:"recognition_date"`
	InboundCode                string  `orm:"column(Inbound_Code);null" json:"inbound_code"`
	SupplierCode               string  `orm:"column(Supplier_Code);null" json:"supplier_code"`
	SupplierName               string  `orm:"column(Supplier_Name);null" json:"supplier_name"`
	GRCode                     string  `orm:"column(GR_Code);null" json:"gr_code"`
	GRStatus                   string  `orm:"column(GR_Status);null" json:"gr_status"`
	WarehouseOrigin            string  `orm:"column(Warehouse_Origin);null" json:"warehouse_origin"`
	WarehouseDestination       string  `orm:"column(Warehouse_Destination);null" json:"warehouse_destination"`
	ProductCode                string  `orm:"column(Product_Code);null" json:"product_code"`
	ProductName                string  `orm:"column(Product_Name);null" json:"product_name"`
	Uom                        string  `orm:"column(UOM);null" json:"uom"`
	GoodsStock                 float64 `orm:"column(Goods_Stock);null" json:"goods_stock"`
	ReceivedQty                float64 `orm:"column(Received_Qty);null" json:"received_qty"`
	AfterSortirQty             float64 `orm:"column(After_Sortir_Qty_Good);null" json:"after_sortir_qty"`
	AfterSortirQtyForWaste     float64 `orm:"column(After_Sortir_Qty_Waste);null" json:"after_sortir_qty(waste)"`
	Discrepancy                float64 `orm:"column(Discrepancy);null" json:"discrepancy"`
	AfterSortirQtyForDownGrade float64 `orm:"column(After_Sortir_Qty_Down_Grade);null" json:"after_sortir_qty(downgrade)"`
	ProductCodeDowngrade       string  `orm:"column(Product_Code_Downgrade);null" json:"product_code_downgrade"`
	ProductNameDowngrade       string  `orm:"column(Product_Name_Downgrade);null" json:"product_name_downgrade"`
	UOMDowngrade               string  `orm:"column(UOM_Downgrade);null" json:"uom_downgrade"`
}
