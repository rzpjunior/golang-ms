package dto

type HelperAppLoginRequest struct {
	Email         string `json:"email" valid:"required"`
	Password      string `json:"password" valid:"required"`
	FirebaseToken string `json:"firebase_token" valid:"required"`
	Timezone      string `json:"timezone"`
}

type HelperAppLoginResponse struct {
	Code          int32          `json:"code"`
	Message       string         `json:"message"`
	User          *HelperAppUser `json:"user"`
	Token         string         `json:"token"`
	FirebaseToken string         `json:"firebase_token"`
}

type HelperAppUser struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	SiteId   string `json:"site_id"`
	SiteName string `json:"site_name"`
	RoleName string `json:"role_name"`
}

type HelperAppGetPickingOrderRequest struct {
	Limit        int
	Offset       int
	LocationCode string
	SopNumber    string
	CustomerName string
	DocNumber    string
	ItemNumber   string
	HelperId     string
	Status       int8

	// jwt
	PickerId string
}

type HelperAppGetPickingOrderResponse struct {
	PickingOrder []*HelperAppGetPickingOrderPickingOrder `json:"picking_order"`
}

type HelperAppGetPickingOrderPickingOrder struct {
	Id              string `json:"id"`
	DocDate         string `json:"doc_date"`
	PickerId        string `json:"picker_id"`
	Status          int8   `json:"status"`
	TotalSalesOrder int64  `json:"total_sales_order"`
	Note            string `json:"note"`
}

type HelperAppGetPickingOrderProductsRequest struct {
	DocNumber      string
	ItemNameSearch string
}

type HelperAppGetPickingOrderProductsResponse struct {
	Id       string                        `json:"id"`
	Status   int8                          `json:"status"`
	Products []*HelperAppAggregatedProduct `json:"product"`
}

type HelperAppDocNumberRequest struct {
	DocNumber string `json:"doc_number"`
}

type HelperAppSopNumberRequest struct {
	SopNumber string `json:"sop_number"`

	SpvId string `json:"-"`
}

type HelperAppAggregatedProduct struct {
	ItemNumber      string  `json:"item_number"`
	ItemName        string  `json:"item_name"`
	Picture         string  `json:"picture"`
	UomDescription  string  `json:"uom_description"`
	TotalOrderQty   float64 `json:"total_order_qty"`
	TotalPickedQty  float64 `json:"total_picked_qty"`
	TotalSalesOrder int32   `json:"total_sales_order"`
	Status          int8    `json:"status"`
}

type HelperAppGetPickingOrderProductsSalesOrder struct {
	ItemId         string                          `json:"item_id"`
	ItemName       string                          `json:"item_name"`
	UomDescription string                          `json:"uom_description"`
	Picture        string                          `json:"picture"`
	SalesOrders    []*HelperAppAggregatedProductSO `json:"sales_order"`
}

type HelperAppAggregatedProductSO struct {
	Id            int64   `json:"id"`
	SopNumber     string  `json:"sop_number"`
	MerchantName  string  `json:"merchant_name"`
	Wrt           string  `json:"wrt"`
	OrderQty      float64 `json:"order_qty"`
	PickedQty     float64 `json:"picked_qty"`
	UnfulfillNote string  `json:"unfulfill_note"`
	Status        int8    `json:"status"`
}

type HelperAppSubmitPickingRequest struct {
	Request []*HelperAppSubmitPickingModel `json:"request"`
}

type HelperAppSubmitPickingModel struct {
	Id            int64   `json:"id"`
	PickQty       float64 `json:"pick_qty"`
	UnfulfillNote string  `json:"unfulfill_note"`
}

type HelperAppGetSalesOrderPickingResponse struct {
	SalesOrder []*HelperAppSalesOrderPickingModel `json:"sales_order"`
}

type HelperAppSalesOrderPickingModel struct {
	SopNumber        string  `json:"sop_number"`
	MerchantName     string  `json:"merchant_name"`
	SopNote          string  `json:"sop_note"`
	TotalKoli        float64 `json:"total_koli"`
	Status           int8    `json:"status"`
	ReadyToPack      bool    `json:"ready_to_pack"`
	ContainUnfulfill bool    `json:"contain_unfulfill"`
	CountPrintDo     int32   `json:"count_print_do"`
	CountPrintSi     int32   `json:"count_print_si"`
}

type HelperAppGetSalesOrderPickingDetailResponse struct {
	SopNumber           string                                     `json:"sop_number"`
	MerchantName        string                                     `json:"merchant_name"`
	Wrt                 string                                     `json:"wrt"`
	DeliveryDate        string                                     `json:"delivery_date"`
	TotalKoli           float64                                    `json:"total_koli"`
	TotalItemOnProgress int64                                      `json:"total_item_on_progress"`
	TotalItem           int64                                      `json:"total_item"`
	SopNote             string                                     `json:"sop_note"`
	Status              int64                                      `json:"status"`
	Item                []*HelperAppGetSalesOrderPickingDetailItem `json:"item"`
	HelperName          string                                     `json:"helper_name"`
	HelperId            string                                     `json:"helper_id"`
}

type HelperAppGetSalesOrderPickingDetailItem struct {
	Id                   int64   `json:"id"`
	PickingOrderAssignId int64   `json:"picking_order_assign_id"`
	ItemNumber           string  `json:"item_number"`
	ItemName             string  `json:"item_name"`
	Picture              string  `json:"picture"`
	OrderQty             float64 `json:"order_qty"`
	PickQty              float64 `json:"pick_qty"`
	CheckQty             float64 `json:"check_qty"`
	ExcessQty            float64 `json:"excess_qty"`
	UnfulfillNote        string  `json:"unfulfill_note"`
	Uom                  string  `json:"uom"`
	Status               int8    `json:"status"`
}
type HelperAppSubmitSalesOrderRequest struct {
	SopNumber string                           `json:"sop_number"`
	Koli      []*HelperAppSubmitSalesOrderKoli `json:"koli"`

	PickerId string `json:"-"`
}

type HelperAppSubmitSalesOrderKoli struct {
	Id       int64   `json:"id"`
	Quantity float64 `json:"quantity"`
}

type HelperAppGetSalesOrderToCheckRequest struct {
	Offset       int
	Limit        int
	SiteId       string
	SopNumber    string
	CustomerName string
	Statuses     []int
	WrtIDs       []string
}

type HelperAppGetSalesOrderToCheckResponse struct {
	SalesOrder []*HelperAppGetSalesOrderToCheckSOModel `json:"sales_order"`
}

type HelperAppGetSalesOrderToCheckSOModel struct {
	SopNumber           string  `json:"sop_number"`
	MerchantName        string  `json:"merchant_name"`
	DeliveryDate        string  `json:"delivery_date"`
	Wrt                 string  `json:"wrt"`
	SopNote             string  `json:"sop_note"`
	TotalItemOnProgress int64   `json:"total_item_on_progress"`
	TotalItem           int64   `json:"total_item"`
	TotalKoli           float64 `json:"total_koli"`
	CheckerName         string  `json:"checker_name"`
	PickerName          string  `json:"picker_name"`
	Status              int8    `json:"status"`
	CountPrintDo        int32   `json:"count_print_do"`
	CountPrintSi        int32   `json:"count_print_si"`
}

type HelperAppGetSalesOrderToCheckDetailResponse struct {
	SopNumber           string                                     `json:"sop_number"`
	MerchantName        string                                     `json:"merchant_name"`
	DeliveryDate        string                                     `json:"delivery_date"`
	Wrt                 string                                     `json:"wrt"`
	SopNote             string                                     `json:"sop_note"`
	TotalItemOnProgress int64                                      `json:"total_item_on_progress"`
	TotalItem           int64                                      `json:"total_item"`
	TotalKoli           float64                                    `json:"total_koli"`
	PickerName          string                                     `json:"picker_name"`
	Item                []*HelperAppGetSalesOrderPickingDetailItem `json:"item"`
	Status              int8                                       `json:"status"`
}

type HelperAppSuccessResponse struct {
	Success bool `json:"success"`
}

type HelperAppCheckerStartCheckingRequest struct {
	SopNumber string `json:"sop_number"`

	CheckerId string `json:"-"`
}

type HelperAppCheckerSubmitCheckingRequest struct {
	SopNumber string                                 `json:"sop_number"`
	Request   []*HelperAppCheckerSubmitCheckingModel `json:"request"`

	CheckerId string `json:"-"`
}
type HelperAppCheckerSubmitCheckingModel struct {
	ItemNumber    string  `json:"item_number"`
	CheckQuantity float64 `json:"check_qty"`
}

type HelperAppCheckerRejectSalesOrderRequest struct {
	SopNumber        string   `json:"sop_number"`
	ItemNumberReject []string `json:"item_number_reject"`

	CheckerId string `json:"-"`
}

type HelperAppCheckerGetDeliveryKoliResponse struct {
	DeliveryKoli []*DeliveryKoliResponse `json:"data"`
}

type HelperAppCheckerAcceptSalesOrderRequest struct {
	SopNumber string                           `json:"sop_number"`
	Koli      []*HelperAppSubmitSalesOrderKoli `json:"koli"`

	CheckerId string `json:"-"`
}

type HelperAppCheckerAcceptSalesOrderResponse struct {
	Success       bool   `json:"success"`
	DeliveryOrder string `json:"delivery_order"`
	SalesInvoice  string `json:"sales_invoice"`
}

type HelperAppCheckerHistoryRequest struct {
	Offset       int
	Limit        int
	SopNumber    string `json:"sop_number"`
	CustomerName string `json:"merchant_name"`
	WrtIdGP      string `json:"wrt_id_gp"`
	CheckerId    string `json:"-"`
}

type HelperAppCheckerHistoryResponse struct {
	SalesOrder []*HelperAppGetSalesOrderToCheckSOModel `json:"sales_order"`
}

type HelperAppCheckerHistoryDetailRequest struct {
	SopNumber string `json:"sop_number"`
	CheckerId string `json:"-"`
}

type HelperAppCheckerHistoryDetailResponse struct {
	SopNumber           string                                     `json:"sop_number"`
	MerchantName        string                                     `json:"merchant_name"`
	DeliveryDate        string                                     `json:"delivery_date"`
	Wrt                 string                                     `json:"wrt"`
	SopNote             string                                     `json:"sop_note"`
	TotalItemOnProgress int64                                      `json:"total_item_on_progress"`
	TotalItem           int64                                      `json:"total_item"`
	TotalKoli           float64                                    `json:"total_koli"`
	PickerName          string                                     `json:"picker_name"`
	Item                []*HelperAppGetSalesOrderPickingDetailItem `json:"item"`
	Status              int8                                       `json:"status"`
	CountPrintDo        int32                                      `json:"count_print_do"`
	CountPrintSi        int32                                      `json:"count_print_si"`
}

type HelperAppPickerWidgetRequest struct {
	HelperId string
}

type HelperAppPickerWidgetResponse struct {
	TotalSalesOrder             int64   `json:"total_sales_order"`
	TotalNew                    int64   `json:"total_new"`
	TotalOnProgress             int64   `json:"total_on_progress"`
	TotalOnProgressPercentage   float64 `json:"total_on_progress_percentage"`
	TotalPicked                 int64   `json:"total_picked"`
	TotalPickedPercentage       float64 `json:"total_picked_percentage"`
	TotalNeedApproval           int64   `json:"total_need_approval"`
	TotalNeedApprovalPercentage float64 `json:"total_need_approval_percentage"`
}

type HelperAppSPVWidgetRequest struct {
	SiteIdGp string
}

type HelperAppSPVWidgetResponse struct {
	TotalSalesOrder             int64   `json:"total_sales_order"`
	TotalNew                    int64   `json:"total_new"`
	TotalOnProgress             int64   `json:"total_on_progress"`
	TotalOnProgressPercentage   float64 `json:"total_on_progress_percentage"`
	TotalNeedApproval           int64   `json:"total_need_approval"`
	TotalNeedApprovalPercentage float64 `json:"total_need_approval_percentage"`
	TotalFinished               int64   `json:"total_finished"`
	TotalFinishedPercentage     float64 `json:"total_finished_percentage"`
}

type HelperAppCheckerWidgetRequest struct {
	CheckerId string
}

type HelperAppCheckerWidgetResponse struct {
	TotalSalesOrder         int64   `json:"total_sales_order"`
	TotalPicked             int64   `json:"total_picked"`
	TotalChecking           int64   `json:"total_checking"`
	TotalCheckingPercentage float64 `json:"total_checking_percentage"`
	TotalFinished           int64   `json:"total_finished"`
	TotalFinishedPercentage float64 `json:"total_finished_percentage"`
}

type HelperAppCheckerGetSalesOrderDetailRequest struct {
	SopNumber string
	CheckerId string
}

type HelperAppHistoryRequest struct {
	Limit        int
	Offset       int
	SopNumber    string
	CustomerName string

	// jwt
	PickerId string
}

type HelperAppHistoryResponse struct {
	SalesOrder []*HelperAppSalesOrderPickingModel `json:"sales_order"`
}

type HelperAppHistoryDetailRequest struct {
	SopNumber string
}

type HelperAppHistoryDetailResponse struct {
	SopNumber           string                                     `json:"sop_number"`
	MerchantName        string                                     `json:"merchant_name"`
	Wrt                 string                                     `json:"wrt"`
	DeliveryDate        string                                     `json:"delivery_date"`
	TotalKoli           float64                                    `json:"total_koli"`
	TotalItemOnProgress int64                                      `json:"total_item_on_progress"`
	TotalItem           int64                                      `json:"total_item"`
	SopNote             string                                     `json:"sop_note"`
	Status              int64                                      `json:"status"`
	Item                []*HelperAppGetSalesOrderPickingDetailItem `json:"item"`
}

type HelperAppWrtMonitoringDetailRequest struct {
	Type     int64    `json:"type"` // 1 = picker , 2 = checker
	HelperId []string `json:"helper_id"`
	WrtId    string   `json:"wrt_id"`

	// jwt token
	SiteId string `json:"-"`
}

type HelperAppWrtMonitoring struct {
	WrtId                string  `json:"wrt_id"`
	WrtDescription       string  `json:"wrt_description"`
	CountSalesOrder      int64   `json:"count_so"`
	OnProgress           int64   `json:"total_on_progress"`
	OnProgressPercentage float64 `json:"on_progress_percentage"`
	Finished             int64   `json:"total_finished"`
	FinishedPercentage   float64 `json:"finished_percentage"`
}

type HelperAppGetWrtMonitoringRequest struct {
	Type     int64    `json:"type"` // 1 = picker , 2 = checker
	HelperId []string `json:"helper_id"`

	// jwt token
	SiteId string `json:"-"`
}
type HelperAppGetWrtMonitoringResponse struct {
	Data []*HelperAppWrtMonitoring `json:"data"`
}

type HelperAppWrtMonitoringDetail struct {
	SopNumber    string  `json:"sop_number"`
	MerchantName string  `json:"merchant_name"`
	TotalKoli    float64 `json:"total_koli"`
	HelperCode   string  `json:"helper_code"`
	HelperName   string  `json:"helper_name"`
	Status       int8    `json:"status"`
}

type HelperAppGetWrtMonitoringDetailResponse struct {
	SalesOrder []*HelperAppWrtMonitoringDetail `json:"sales_order"`
}
