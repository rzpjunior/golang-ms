package picking

import (
	"encoding/json"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"git.edenfarm.id/cuxs/mongodb"
	"go.mongodb.org/mongo-driver/bson"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

func getPickingOrder(rq *orm.RequestQuery, cond map[string]interface{}) (m []*templatePickingOrder, err error) {

	// get data requested
	var where string
	var values []interface{}

	o := orm.NewOrm()
	o.Using("read_only")

	for k, v := range cond {
		if reflect.TypeOf(v).Kind().String() == "slice" {
			where = where + " " + k + "? and ? and"
		} else {
			where = where + " " + k + "? and"
		}

		values = append(values, v)
	}
	where = strings.TrimSuffix(where, " and")
	if where != "" {
		where = "AND " + where
	}

	q := "SELECT w.name warehouse , m.name merchant_name , bt.name business_type , so.code order_code, adm.province_name province ," +
		" adm.city_name city , adm.district_name district , adm.sub_district_name sub_district , adm.postal_code postal_code ," +
		" wrt.name wrt , so.total_weight order_weight , so.delivery_date delivery_date , tps.name payment_term ," +
		" os.name sales_order_type , upper( gg.value_name) AS sales_order_status , so.shipping_address shipping_address , s.name picker ," +
		" cv.name vendor, poa.planning_vendor planning, pl.code pl_code" +
		" FROM sales_order so" +
		" LEFT JOIN picking_order_assign poa ON poa.sales_order_id = so.id" +
		" LEFT JOIN courier_vendor cv ON poa.courier_vendor_id = cv.id" +
		" LEFT JOIN staff s ON s.id = poa.staff_id" +
		" JOIN order_type_sls os ON os.id = so.order_type_sls_id" +
		" JOIN adm_division adm ON adm.sub_district_id = so.sub_district_id" +
		" JOIN wrt ON wrt.id = so.wrt_id" +
		" JOIN warehouse w ON w.id = so.warehouse_id" +
		" JOIN term_payment_sls tps ON tps.id = so.term_payment_sls_id" +
		" JOIN branch b on b.id = so.branch_id" +
		" JOIN merchant m on m.id = b.merchant_id" +
		" JOIN business_type bt on bt.id = m.business_type_id" +
		" JOIN glossary gg ON gg.value_int = so.status AND gg.attribute = 'status' AND gg.table = 'sales_order'" +
		" LEFT JOIN picking_list pl ON pl.id = poa.picking_list_id AND pl.status IN (1,3) " +
		"WHERE(so.order_type_sls_id != 10 AND so.status in(1,9,12) " + where + ") AND (poa.status = 1 OR poa.id IS NULL)"

	_, err = o.Raw(q, values).QueryRows(&m)

	return
}

type ResponsePrint struct {
	LinkPrint string  `json:"link_print"`
	TotalKoli float64 `json:"total_koli"`
}

func UpdatePrintLabel(soID int64) (e error) {
	o := orm.NewOrm()

	if _, e = o.Raw("update delivery_koli_increment set print_label = ? where sales_order_id = ?", 1, soID).Exec(); e != nil {
		return e
	}
	return e
}

func getListSOMonitoring(r wrtMonitoringRequest) (m []*ItemSOMonitoring, err error) {

	o := orm.NewOrm()
	o.Using("read_only")

	var result []*ItemSOMonitoring
	var qMarkIdStaff string
	var qIdStaff []string
	var query string

	if len(r.HelperID) != 0 {
		for _, v := range r.HelperID {
			qMarkIdStaff += "?,"
			qIdStaff = append(qIdStaff, v)
		}
		qMarkIdStaff = qMarkIdStaff[:len(qMarkIdStaff)-1]
		if r.HelperType == "picker" {
			query = "AND poa.staff_id IN (" + qMarkIdStaff + ") "
		} else {
			query = "AND poa.checked_by IN (" + qMarkIdStaff + ") "
		}
	}

	q := "SELECT so.code " +
		", poa.status " +
		", g.value_name 'picking_status', m.name 'merchant' " +
		", IF(FIND_IN_SET(1,m.tag_customer) > 0 , 'NC', IF(FIND_IN_SET(8,m.tag_customer) > 0, 'PC', NULL)) 'customer_tag' " +
		", SUM(dk.quantity) 'total_koli'" +
		", IF(s2.id IS NOT NULL, s2.name, s.name)'helper_name' " +
		", IF(s2.id IS NOT NULL, s2.code, s.code)'helper_code' " +
		", so.wrt_id " +
		"FROM picking_order_assign poa " +
		"JOIN picking_order po ON po.id = poa.picking_order_id " +
		"JOIN sales_order so ON so.id = poa.sales_order_id " +
		"JOIN branch b ON b.id = so.branch_id " +
		"JOIN merchant m ON m.id = b.merchant_id " +
		"LEFT JOIN delivery_koli dk ON dk.sales_order_id = poa.sales_order_id " +
		"LEFT JOIN staff s ON s.id = poa.staff_id " +
		"LEFT JOIN staff s2 ON s2.id = poa.checked_by " +
		"LEFT JOIN glossary g ON g.value_int = poa.status AND g.`table` = 'picking_order' AND g.`attribute` = 'doc_status_picking' " +
		"WHERE po.warehouse_id = ? " + query +
		"AND po.recognition_date IN(CURDATE(), CURDATE() + INTERVAL 1 DAY) " +
		"AND so.wrt_id = ? " +
		"GROUP BY poa.sales_order_id"

	if _, err = o.Raw(q, r.Session.Staff.Warehouse.ID, r.Staff, r.WRT.ID).QueryRows(&result); err != nil {
		return
	}

	return result, err
}

func getWRTMonitorings(r wrtMonitoringRequest) (m []*ItemWRTMonitoring, err error) {

	o := orm.NewOrm()
	o.Using("read_only")

	type qtyStore struct {
		WRTID           string `orm:"-" json:"id"`
		WRTName         string `json:"-" json:"wrt"`
		TotalSalesOrder int64  `json:"-"`
		TotalOnProgress int64  `orm:"-" json:"total_on_progress"`
		TotalPicked     int64  `orm:"-" json:"total_picked"`
		TotalChecking   int64  `orm:"-" json:"total_checking"`
		TotalFinished   int64  `orm:"-" json:"total_finished"`
	}

	var resTemp []*ItemWRTMonitoring
	var result []*ItemWRTMonitoring
	duplicateWRT := make(map[string]qtyStore)
	var query string
	var qMarkIdStaff string
	var qIdStaff []int64

	if len(r.HelperID) != 0 {
		for _, v := range r.Staff {
			qMarkIdStaff += "?,"
			qIdStaff = append(qIdStaff, v.ID)
		}

		qMarkIdStaff = qMarkIdStaff[:len(qMarkIdStaff)-1]
	}

	if r.HelperType == "picker" {
		if len(r.HelperID) != 0 {
			query = "AND poa.staff_id IN (" + qMarkIdStaff + ") "
		}
	} else {
		if len(r.HelperID) != 0 {
			query = "AND poa.checked_by IN (" + qMarkIdStaff + ") "
		}
	}

	q := "SELECT w.id, w.name, poa.status, poa.staff_id, poa.checked_by, s.name 'lead_picker', s2.name 'checker' " +
		"FROM picking_order_assign poa " +
		"JOIN sales_order so ON so.id = poa.sales_order_id JOIN branch b ON b.id = so.branch_id " +
		"JOIN picking_order po ON po.id = poa.picking_order_id " +
		"JOIN wrt w ON w.id = so.wrt_id " +
		"LEFT JOIN staff s ON s.id = poa.staff_id " +
		"LEFT JOIN staff s2 ON s2.id = poa.checked_by " +
		"WHERE po.recognition_date IN (CURDATE(), CURDATE() + INTERVAL 1 DAY) AND po.warehouse_id = ? " + query

	if _, err = o.Raw(q, r.Session.Staff.Warehouse.ID, qIdStaff).QueryRows(&resTemp); err != nil {
		return
	}

	for _, v := range resTemp {
		var qtyTemp qtyStore
		if val, ok := duplicateWRT[v.WRTID]; ok {
			if v.PickingOrderStatus == 2 {
				val.TotalFinished += 1
			}
			if v.PickingOrderStatus == 3 || v.PickingOrderStatus == 8 {
				val.TotalOnProgress += 1
			}
			if v.PickingOrderStatus == 5 {
				val.TotalPicked += 1
			}
			if v.PickingOrderStatus == 6 {
				val.TotalChecking += 1
			}
			val.TotalSalesOrder += 1
			val.WRTName = v.WRTName
			val.WRTID = v.WRTID
			duplicateWRT[v.WRTID] = val
		} else {
			if v.PickingOrderStatus == 2 {
				qtyTemp.TotalFinished += 1
			}
			if v.PickingOrderStatus == 3 || v.PickingOrderStatus == 8 {
				qtyTemp.TotalOnProgress += 1
			}
			if v.PickingOrderStatus == 5 {
				qtyTemp.TotalPicked += 1
			}
			if v.PickingOrderStatus == 6 {
				qtyTemp.TotalChecking += 1
			}
			qtyTemp.TotalSalesOrder += 1
			qtyTemp.WRTName = v.WRTName
			qtyTemp.WRTID = v.WRTID
			duplicateWRT[v.WRTID] = qtyTemp
		}
	}

	for _, v := range duplicateWRT {
		temp := new(ItemWRTMonitoring)
		temp.WRTID = common.Encrypt(v.WRTID)
		temp.WRTName = v.WRTName
		temp.TotalSalesOrder = v.TotalSalesOrder
		temp.TotalOnProgress = v.TotalOnProgress
		temp.TotalPicked = v.TotalPicked
		temp.TotalChecking = v.TotalChecking
		temp.TotalFinished = v.TotalFinished
		result = append(result, temp)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].WRTName < result[j].WRTName
	})

	return result, err
}

func getListHelper(r listHelperRequest) (m []*model.Staff, err error) {
	// get data requested
	o := orm.NewOrm()
	o.Using("read_only")

	var result []*model.Staff

	q := "SELECT s.id, s.code, s.name FROM staff s JOIN role r ON r.id = s.role_id  WHERE s.warehouse_id = ? AND r.name = ?"

	if _, err = o.Raw(q, r.Session.Staff.Warehouse.ID, r.HelperRole).QueryRows(&result); err != nil {
		return
	}

	return result, err
}

func getProfile(r profileRequestPicking) (m *MobileProfile, err error) {

	// get data requested
	o := orm.NewOrm()
	o.Using("read_only")

	// for leadpicker
	leadPickerRoleCode, _ := repository.GetConfigApp("attribute", "lead_picker_role_id")
	if r.Session.Staff.Role.Code == leadPickerRoleCode.Value {
		var status []int
		res := new(MobileProfile)

		q := "SELECT poa.status " +
			"FROM picking_order_assign poa JOIN picking_order po ON po.id = poa.picking_order_id " +
			"WHERE po.recognition_date IN(CURDATE(), CURDATE() + INTERVAL 1 DAY) AND poa.staff_id = ? AND poa.status IN (1,3,4,5)"

		if _, err = o.Raw(q, r.Session.Staff.ID).QueryRows(&status); err != nil {
			return
		}

		for _, v := range status {
			switch v {
			case 1:
				res.New += 1
			case 3:
				res.OnProgress += 1
			case 8:
				res.OnProgress += 1
			case 5:
				res.Picked += 1
			case 4:
				res.NeedApproval += 1
			}
			res.TotalAssign += 1
		}
		m = res

		if m.TotalAssign == 0 {
			m.OnProgressTask = 0
			m.FinishedTask = 0
			m.NeedApprovalTask = 0
			return

		}

		m.OnProgressTask = common.Rounder(m.OnProgress/m.TotalAssign*100, 0.5, 2)
		m.FinishedTask = common.Rounder(m.Picked/m.TotalAssign*100, 0.5, 2)
		m.NeedApprovalTask = common.Rounder(m.NeedApproval/m.TotalAssign*100, 0.5, 2)

		return

	}
	if r.Session.Staff.Role.Code == "ROL0049" { //for checker
		var checkerTemp []CheckerProfileTemp
		res := new(MobileProfile)

		q := "SELECT poa.status, poa.checked_by " +
			"FROM picking_order_assign poa JOIN picking_order po ON po.id = poa.picking_order_id " +
			"WHERE po.recognition_date IN(CURDATE(), CURDATE() + INTERVAL 1 DAY) AND po.warehouse_id = ? AND poa.status IN (2,5,6)"

		if _, err = o.Raw(q, r.Session.Staff.Warehouse.ID).QueryRows(&checkerTemp); err != nil {
			return
		}

		if len(checkerTemp) == 0 {
			return res, err
		}

		for _, v := range checkerTemp {
			switch v.Status {
			case 2:
				if r.Session.Staff.ID == v.CheckedBy {
					res.Finished += 1
				}
			case 5:
				res.Picked += 1
			case 6:
				if r.Session.Staff.ID == v.CheckedBy {
					res.Checking += 1
				}
			}
			res.TotalAssign += 1
		}
		m = res

		if m.Picked+m.Checking+m.Finished == 0 {
			m.OnProgressTask = 0
			m.FinishedTask = 0
			return
		}

		m.OnProgressTask = common.Rounder(m.Checking/(m.Picked+m.Checking+m.Finished)*100, 0.5, 2)
		m.FinishedTask = common.Rounder(m.Finished/(m.Picked+m.Checking+m.Finished)*100, 0.5, 2)

		return

	}
	if r.Session.Staff.Role.Code == "ROL0022" { //for packer
		res := new(MobileProfile)
		m = res
		return
	} else { //for SPV
		var pickingTemp []MobileProfile
		res := new(MobileProfile)

		q := "SELECT poa.sales_order_id, poa.status FROM picking_order_assign poa JOIN picking_order po ON po.id = poa.picking_order_id WHERE po.warehouse_id = ? AND po.recognition_date IN(CURDATE(), CURDATE() + INTERVAL 1 DAY)"

		if _, err = o.Raw(q, r.Session.Staff.Warehouse.ID).QueryRows(&pickingTemp); err != nil {
			return
		}

		for _, v := range pickingTemp {
			switch v.PickingOrderStatus {
			case 1:
				res.New += 1
			case 2:
				res.Finished += 1
			case 3:
				res.OnProgress += 1
			case 8:
				res.OnProgress += 1
			case 4:
				res.NeedApproval += 1
			}
			res.TotalAssign += 1
		}

		m = res

		if m.TotalAssign == 0 {
			m.OnProgressTask = 0
			m.FinishedTask = 0
			m.NeedApprovalTask = 0
			return
		}

		m.OnProgressTask = common.Rounder(m.OnProgress/m.TotalAssign*100, 0.5, 2)
		m.FinishedTask = common.Rounder(m.Finished/m.TotalAssign*100, 0.5, 2)
		m.NeedApprovalTask = common.Rounder(m.NeedApproval/m.TotalAssign*100, 0.5, 2)

		return
	}

}

func getListProduct(r listProductRequestPicking) (m *ItemPickingList, err error) {

	// get data requested
	o := orm.NewOrm()
	o.Using("read_only")

	type qtyStore struct {
		PickingRouting           int8    `json:"picking_routing"`
		PickingListStatus        int8    `json:"picking_list_status"`
		OrderQty                 float64 `json:"order_qty"`
		PickQty                  float64 `json:"pick_qty"`
		TotalSalesOrder          int8    `json:"-"`
		TotalCancelledSalesOrder int8    `json:"-"`
		IsSKUReviewable          int8    `json:"-"`
		RejectedByChecker        int8    `json:"-"`
	}

	var (
		where                    string
		resTemp                  []*ItemPickingListProduct
		resultPickingListProduct []*ItemPickingListProduct
		pickerStr                string
		pickerArr                []string
		pickers                  []*model.Staff
		pickingRoutingFlag       int8
		filter, exclude          map[string]interface{}
	)

	duplicateProducts := make(map[int64]qtyStore)

	if r.Query != "" {
		where = " and (m.name like ? or so.code like ? or p.name like ? or p.code like ?) "
	}

	q := "SELECT pl.status AS picking_list_status,p.id, p.name, u.name , poi.order_qty, poi.pick_qty, poi.unfullfill_note, poi.flag_saved_pick, poa.status, so.status 'status_sales_order', poi.picking_flag 'poi_status' " +
		"FROM picking_order_item poi " +
		"JOIN picking_order_assign poa ON poa.id = poi.picking_order_assign_id AND poa.status IN(1,3,7,8) " +
		"JOIN sales_order so on poa.sales_order_id = so.id AND so.status in(1,3,4,9,13,12) " +
		"JOIN branch b ON b.id = so.branch_id " +
		"JOIN merchant m ON m.id = b.merchant_id " +
		"JOIN picking_list pl ON pl.id = poa.picking_list_id AND pl.id = ? " +
		"JOIN product p ON p.id = poi.product_id " +
		"JOIN uom u ON u.id = p.uom_id " + where +
		" ORDER BY p.name"

	if r.Query != "" {
		if _, err = o.Raw(q, r.PickingList.ID, "%"+r.Query+"%", "%"+r.Query+"%", "%"+r.Query+"%", "%"+r.Query+"%").QueryRows(&resTemp); err != nil {
			return
		}
	} else {
		if _, err = o.Raw(q, r.PickingList.ID).QueryRows(&resTemp); err != nil {
			return
		}
	}

	// looping for classified product based on unique product and summarize the quantity
	for _, v := range resTemp {
		var qtyTemp qtyStore
		if val, ok := duplicateProducts[v.ProductID]; ok {
			val.PickingListStatus = v.PickingListStatus
			val.OrderQty += v.TotalOrder
			val.PickQty += v.PickQty
			val.TotalSalesOrder += 1
			val.IsSKUReviewable, _ = GetSavedPickForProducts(val.IsSKUReviewable, v.FlagSavedPick)
			val.TotalCancelledSalesOrder, _ = GetFlagDisableForProducts(val.TotalCancelledSalesOrder, v.SalesOrderStatus)

			if duplicateProducts[v.ProductID].RejectedByChecker == 0 {
				// 4 is picking order item status code for rejected by checker
				if v.PickingOrderItemStatus == 4 {
					val.RejectedByChecker = 1
				}
			} else {
				val.RejectedByChecker = 1
			}

			duplicateProducts[v.ProductID] = val
		} else {
			qtyTemp.PickingListStatus = v.PickingListStatus
			qtyTemp.OrderQty = v.TotalOrder
			qtyTemp.PickQty = v.PickQty
			qtyTemp.TotalSalesOrder += 1
			qtyTemp.IsSKUReviewable, _ = GetSavedPickForProducts(0, v.FlagSavedPick)
			qtyTemp.TotalCancelledSalesOrder, _ = GetFlagDisableForProducts(qtyTemp.TotalCancelledSalesOrder, v.SalesOrderStatus)

			if duplicateProducts[v.ProductID].RejectedByChecker == 0 {
				// 8 is picking order item status code for rejected by checker
				if v.PickingOrderItemStatus == 4 {
					qtyTemp.RejectedByChecker = 1
				}
			} else {
				qtyTemp.RejectedByChecker = 1
			}

			duplicateProducts[v.ProductID] = qtyTemp
		}
	}

	/**
	looping for add all of unique product
	to resultPickingListProduct slice with summarize quantity too
	*/
	var p *model.Product
	for k, v := range duplicateProducts {
		temp := new(ItemPickingListProduct)
		if v.TotalSalesOrder == v.TotalCancelledSalesOrder {
			temp.FlagDisableSku = 1
		}
		if v.TotalSalesOrder == v.IsSKUReviewable {
			temp.FlagSavedPick = 1
		}
		if p, err = repository.ValidProduct(k); err != nil {
			return
		}
		o.LoadRelated(p, "ProductImage", 1)
		temp.Product = p
		temp.TotalOrder = v.OrderQty
		temp.PickQty = v.PickQty
		temp.PickingListStatus = v.PickingListStatus

		temp.FlagRejectedByChecker = v.RejectedByChecker
		resultPickingListProduct = append(resultPickingListProduct, temp)
	}

	for _, v := range resultPickingListProduct {
		v.Product.Uom.Read("ID")
	}

	// sort product by product name descending
	sort.Slice(resultPickingListProduct, func(i, j int) bool {
		return resultPickingListProduct[i].Product.Name < resultPickingListProduct[j].Product.Name
	})

	filter = map[string]interface{}{"picking_list_id": r.PickingList.ID, "status_step__in": []int64{2, 3}}
	_, total, err := repository.CheckPickingRoutingStepData(filter, exclude)
	if err != nil {
		return
	}
	if total == 0 {
		pickingRoutingFlag = 3
	} else {
		pickingRoutingFlag = 1
	}

	for _, v := range resultPickingListProduct {
		v.PickingRouting = pickingRoutingFlag
	}

	// prepare the pickers data
	o.Raw("select sub_picker_id from picking_order_assign poa where poa.picking_list_id = ? LIMIT 1", r.PickingList.ID).QueryRow(&pickerStr)
	pickerArr = strings.Split(pickerStr, ",")

	for _, v2 := range pickerArr {
		pickerID, err := strconv.Atoi(v2)
		if err != nil {
			continue
		}
		picker := &model.Staff{ID: int64(pickerID)}
		picker.Read("id")

		pickers = append(pickers, picker)
	}

	result := &ItemPickingList{
		ItemPickingListProduct: resultPickingListProduct,
		Pickers:                pickers,
	}

	return result, err
}

func getListPicking(r listRequestPicking, cond map[string]interface{}) (m []*GroupPickingList, err error) {

	var (
		where                string
		values               []interface{}
		qMarkPickingList     string
		qFilterPickingList   []int
		qMarkPickingStatus   string
		qFilterPickingStatus []int
		filter, exclude      map[string]interface{}
		total                int64
	)

	// get data requested
	o := orm.NewOrm()
	o.Using("read_only")
	var queryString string

	for k, v := range cond {
		if reflect.TypeOf(v).Kind().String() == "slice" {
			where = where + " " + k + " >= ? and " + k + " <= ? and"
			values = append(values, v)
		} else if k == "query" {
			where = where + " (m.name like '%" + r.Query + "%' or so.code like '%" + r.Query + "%' or p.name like '%" + r.Query + "%' or pl.code like '%" + r.Query + "%' or pl.note like '%" + r.Query + "%') and"
			queryString = " and (m.name like ? or so.code like ? or p.name like ? or pl.code like ? or pl.note like ?)"
		} else if k == "filter_picking_list" {
			if r.FilterPickingList == "all" {
				qMarkPickingList = "?,?,?"
				qFilterPickingList = append(qFilterPickingList, 1, 3, 4)
			} else if r.FilterPickingList == "new" {
				qMarkPickingList = "?"
				qFilterPickingList = append(qFilterPickingList, 1)
			} else if r.FilterPickingList == "onProgress" {
				qMarkPickingList = "?"
				qFilterPickingList = append(qFilterPickingList, 3)
			} else {
				qMarkPickingList = "?"
				qFilterPickingList = append(qFilterPickingList, 4)
			}
		} else if k == "filter_picking_status" {
			if r.FilterPickingStatus == "all" {
				qMarkPickingStatus = "?,?,?,?"
				qFilterPickingStatus = append(qFilterPickingStatus, 1, 3, 4, 8)
			} else if r.FilterPickingStatus == "new" {
				qMarkPickingStatus = "?"
				qFilterPickingStatus = append(qFilterPickingStatus, 1)
			} else if r.FilterPickingStatus == "onProgress" {
				qMarkPickingStatus = "?"
				qFilterPickingStatus = append(qFilterPickingStatus, 3)
			} else if r.FilterPickingStatus == "needApproval" {
				qMarkPickingStatus = "?"
				qFilterPickingStatus = append(qFilterPickingStatus, 4)
			} else {
				qMarkPickingStatus = "?"
				qFilterPickingStatus = append(qFilterPickingStatus, 8)
			}
		} else {
			where = where + " " + k + "? and"
			values = append(values, v)
		}
	}
	where = strings.TrimSuffix(where, " and")
	if where != "" {
		where = "AND " + where + queryString
	}

	var result []*GroupPickingList
	duplicatePickingList := make(map[int64]*GroupPickingList)

	q := "SELECT pl.delivery_date, pl.status, pl.note, pl.id , pl.code ,m.tag_customer,poa.sales_order_id " +
		"FROM picking_list pl " +
		"JOIN picking_order_assign poa ON poa.picking_list_id = pl.id AND poa.status IN(" + qMarkPickingStatus + ") AND poa.staff_id = ? " +
		"JOIN sales_order so ON so.id = poa.sales_order_id " +
		"JOIN branch b ON b.id = so.branch_id " +
		"JOIN merchant m ON m.id = b.merchant_id " +
		"JOIN picking_order_item poi on poi.picking_order_assign_id = poa.id " +
		"JOIN product p on poi.product_id = p.id " +
		"WHERE pl.status IN (" + qMarkPickingList + ") " + where +
		" GROUP BY poa.sales_order_id"

	if _, ok := cond["query"]; ok {
		if _, err = o.Raw(q, qFilterPickingStatus, r.Staff.ID, qFilterPickingList, values, "%"+r.Query+"%", "%"+r.Query+"%", "%"+r.Query+"%", "%"+r.Query+"%", "%"+r.Query+"%").QueryRows(&m); err != nil {
			return
		}
	} else {
		if _, err = o.Raw(q, qFilterPickingStatus, r.Staff.ID, qFilterPickingList, values).QueryRows(&m); err != nil {
			return
		}
	}

	for _, v := range m {
		var qtyTemp *GroupPickingList
		qtyTemp = new(GroupPickingList)

		if val, ok := duplicatePickingList[v.PickingListID]; ok {
			val.TotalSalesOrder += 1

			if v.TagCustomer == 1 {
				val.NewCustomer += 1
			} else if v.TagCustomer == 8 {
				val.PowerCustomer += 1
			} else {
				val.Other += 1
			}
			duplicatePickingList[v.PickingListID] = val

		} else {
			qtyTemp.TotalSalesOrder += 1
			qtyTemp.PickingListCode = v.PickingListCode
			qtyTemp.PickingListNote = v.PickingListNote
			qtyTemp.PickingListID = v.PickingListID
			qtyTemp.Status = v.Status
			qtyTemp.DeliveryDate = v.DeliveryDate
			if v.TagCustomer == 1 {
				qtyTemp.NewCustomer += 1
			} else if v.TagCustomer == 8 {
				qtyTemp.PowerCustomer += 1
			} else {
				qtyTemp.Other += 1
			}
			duplicatePickingList[v.PickingListID] = qtyTemp
		}
	}

	for _, v := range duplicatePickingList {
		var pickerStr string
		var pickerArr []string
		o.Raw("select sub_picker_id from picking_order_assign poa where poa.picking_list_id = ? LIMIT 1", v.PickingListID).QueryRow(&pickerStr)
		pickerArr = strings.Split(pickerStr, ",")

		for _, v2 := range pickerArr {
			pickerID, err := strconv.Atoi(v2)
			if err != nil {
				continue
			}
			picker := &model.Staff{ID: int64(pickerID)}
			picker.Read("id")

			v.Pickers = append(v.Pickers, picker)
		}

		filter = map[string]interface{}{"picking_list_id": v.PickingListID, "status_step__in": []int64{2, 3}}
		_, total, err = repository.CheckPickingRoutingStepData(filter, exclude)
		if err != nil {
			return
		}

		if total == 0 {
			v.PickingRouting = 3
		} else {
			v.PickingRouting = 1
		}

		statusName, _ := repository.GetGlossaryMultipleValue("table", "picking_list", "attribute", "status", "value_int", v.Status)
		v.StatusConvert = statusName.ValueName
		iDEncrypt := common.Encrypt(v.PickingListID)
		plID, _ := strconv.ParseInt(iDEncrypt, 10, 64)
		v.PickingListID = plID

		result = append(result, v)
	}

	// sort product by product name descending
	sort.Slice(result, func(i, j int) bool {
		return result[i].PickingListID < result[j].PickingListID
	})

	return result, err
}

func getPickingListGenCode(r generateCodePickingRequest) (m map[string]ListPl, err error) {

	// get data requested
	o := orm.NewOrm()
	o.Using("read_only")

	plMap := make(map[int64]GenerateCodePickingList)

	var filterDistrict string
	if len(r.ArrDistrictID) != 0 {
		filterDistrict = "AND ad.district_id IN (" + r.QueryStringDistrict + ") "
	}

	var filterSubDistrict string
	if len(r.ArrSubDistrictID) != 0 {
		filterDistrict = "AND ad.sub_district_id IN (" + r.QueryStringSubDistrict + ") "
	}

	var filterWrt string
	if len(r.ArrWrtID) != 0 {
		filterWrt = "AND w.id IN (" + r.QueryStringWrt + ") "
	}

	var filterBusinessType string
	if len(r.ArrBusinessTypeID) != 0 {
		filterBusinessType = "AND bt.id IN (" + r.QueryStringBusinessType + ") "
	}

	var filterSalesOrderType string
	if len(r.ArrSalesOrderTypeID) != 0 {
		filterSalesOrderType = "AND so.order_type_sls_id IN (" + r.QueryStringSalesOrderType + ") "
	}

	var filterCity string
	if len(r.ArrCity) != 0 {
		filterCity = "AND ad.city_id = ? "
	}

	var res []GenerateCodePickingList

	q := "SELECT so2.id, so2.code 'so_code', so2.so_total, w.name 'wrt' , p.name 'product_name', soi2.weight 'weight_item', soi2.order_qty 'order_item' " +
		"FROM sales_order_item soi2 " +
		"JOIN(SELECT so.id, so.wrt_id, so.status , so.code , so.branch_id , SUM(soi.weight) 'so_total' " +
		"FROM sales_order so " +
		"JOIN sales_order_item soi ON soi.sales_order_id = so.id " +
		"LEFT JOIN picking_order_assign poa ON poa.sales_order_id = so.id " +
		"LEFT JOIN archetype a ON a.id = so.archetype_id " +
		"LEFT JOIN business_type bt ON bt.id = a.business_type_id " +
		"WHERE so.delivery_date = ? " + filterBusinessType + filterSalesOrderType +
		"and so.status IN (1,9,12) and so.order_type_sls_id != 10 and so.id not in (SELECT id from sales_order so2 WHERE so2.status =1 and so2.term_payment_sls_id = 11) " +
		"and so.warehouse_id = ? and poa.id is NULL " +
		"GROUP BY so.id) so2 on so2.id = soi2.sales_order_id " +
		"JOIN product p ON p.id = soi2.product_id " +
		"JOIN wrt w ON w.id = so2.wrt_id " + filterWrt + " " +
		"JOIN branch b ON b.id = so2.branch_id " +
		"JOIN adm_division ad ON ad.sub_district_id = b.sub_district_id " + filterDistrict + filterSubDistrict + filterCity + "ORDER BY ad.district_id, w.name, so2.code"

	if _, err = o.Raw(q, r.DeliveryDate, r.ArrBusinessTypeID, r.ArrSalesOrderTypeID, r.Warehouse, r.ArrWrtID, r.ArrDistrictID, r.ArrSubDistrictID, r.ArrCity).QueryRows(&res); err != nil {
		return
	}

	// group by sales order
	for _, v := range res {
		plMap[v.SalesOrderID] = v
	}

	if len(plMap) == 0 {
		return nil, err
	}

	if m, err = GeneratePickingList(plMap, r.Warehouse, r.LimitSalesOrder, r.LimitWeight); err != nil {
		return nil, err
	}

	return
}
func getSalesOrderGroupByPickingList(r groupingSalesOrderRequest) (m []*GroupingSalesOrder, err error) {

	// get data requested
	o := orm.NewOrm()
	o.Using("read_only")
	mdb := mongodb.NewMongo()

	q := "SELECT soi.note sales_order_item_note,so.status 'status_sales_order',w.name 'wrt', poa.id 'picking_order_assign_id',p.name 'product', u.name 'uom', poa.sales_order_id ,poa.status, so.code , m.name 'merchant', m.tag_customer , poi.order_qty , poi.pick_qty , poi.unfullfill_note, poi.picking_flag `picking_flag` " +
		"FROM picking_order_item poi " +
		"JOIN picking_order_assign poa ON poa.id = poi.picking_order_assign_id AND poa.status IN(1,3,7,8) " +
		"JOIN sales_order so ON so.id = poa.sales_order_id AND so.status in (1,9,12,3,4,13) " +
		"LEFT JOIN sales_order_item soi ON soi.sales_order_id = so.id and soi.product_id = ? " +
		"JOIN wrt w ON so.wrt_id = w.id " +
		"JOIN branch b ON b.id = so.branch_id " +
		"JOIN merchant m ON m.id = b.merchant_id " +
		"JOIN picking_list pl ON pl.id = poa.picking_list_id AND pl.id = ? " +
		"JOIN product p ON p.id = poi.product_id AND p.id = ? " +
		"JOIN uom u ON u.id = p.uom_id " +
		"ORDER BY poa.sales_order_id"

	if _, err = o.Raw(q, r.Product.ID, r.PickingList.ID, r.Product.ID).QueryRows(&m); err != nil {
		return
	}

	for _, v := range m {
		v.SalesOrderID = common.Encrypt(v.SalesOrder)
		v.PickingOrderAssignID = common.Encrypt(v.PickingOrderAssign)
		strSepComma := strings.Split(v.TagCustomerDB, ",")
		v.TagCustomer, _ = util.GetCustomerTag(strSepComma)
	}

	for _, v := range m {
		filter := bson.D{
			{"sales_order_id", v.SalesOrder},
			{"product_id", r.Product.ID},
			{"pack_type", bson.M{"$ne": -1}},
			{"status", 1},
		}

		var res []byte
		if res, err = mdb.GetAllDataWithFilter("Packing_Sales_Order", filter); err != nil {
			mdb.DisconnectMongoClient()
			return nil, err
		}

		if len(res) == 0 {
			continue
		}
		// region convert byte data to json data
		if err = json.Unmarshal(res, &r.PackRecommendation); err != nil {
			mdb.DisconnectMongoClient()
			return nil, err
		}
		// endregion
		v.PackRecommendation = r.PackRecommendation

		for _, v := range v.PackRecommendation {
			pcIDStr := common.Encrypt(v.PackingOrderID)
			pcID, _ := strconv.Atoi(pcIDStr)
			v.PackingOrderID = int64(pcID)

			pIDStr := common.Encrypt(v.ProductID)
			pID, _ := strconv.Atoi(pIDStr)
			v.ProductID = int64(pID)

			soIDStr := common.Encrypt(v.SalesOrderID)
			soID, _ := strconv.Atoi(soIDStr)
			v.SalesOrderID = int64(soID)
		}

	}

	mdb.DisconnectMongoClient()
	return
}

/*
method for generate picking list
if weight > 80 build picking list code directly
if not, do the operation
*/

var plMapFinal = make(map[string]ListPl)
var plMapReturn = make(map[string]ListPl)

func GeneratePickingList(plMap map[int64]GenerateCodePickingList, w *model.Warehouse, limitSalesOrder int, limitMaxWeight float64) (m map[string]ListPl, err error) {
	plMapCall := make(map[int64]GenerateCodePickingList)
	var plTemp ListPl
	var plTempLess ListPl

	/*
		initialize map not in first build api
	*/
	// map generator
	if len(plMapFinal) == 0 {
		plMapFinal = make(map[string]ListPl)
	}

	for _, v := range plMap {
		if v.TotalWeight < limitMaxWeight {

			// case if the next summarize > limitMaxWeight
			if plTempLess.TotalWeight+v.TotalWeight <= limitMaxWeight {

				plTempLess.TotalWeight += v.TotalWeight
				plTempLess.SalesOrderID = append(plTempLess.SalesOrderID, v.SalesOrderID)

				if len(plTempLess.SalesOrderID) == limitSalesOrder {
					plMapFinal, _ = InsertOnePickingList(plTempLess.TotalWeight, plTempLess.SalesOrderID, w, plMapFinal)

					plTempLess.TotalWeight = 0
					plTempLess.SalesOrderID = nil
					continue

				}
			} else {
				plMapCall[v.SalesOrderID] = v
			}

		} else {
			// bikin PL tersendiri untuk SO yang melebihi max weight
			plTemp.SalesOrderID = append(plTemp.SalesOrderID, v.SalesOrderID)
			plMapFinal, _ = InsertOnePickingList(v.TotalWeight, plTemp.SalesOrderID, w, plMapFinal)

			plTemp.TotalWeight = 0
			plTemp.SalesOrderID = nil

			continue
		}
	}
	/*
		so in here whatever the value in tempLess, it will become one picking order
	*/
	if len(plTempLess.SalesOrderID) != 0 {
		plMapFinal, _ = InsertOnePickingList(plTempLess.TotalWeight, plTempLess.SalesOrderID, w, plMapFinal)

		plTempLess.TotalWeight = 0
		plTempLess.SalesOrderID = nil
	}

	/*
		if there is value in this map,
		it means there is sales order summarize that still grater than 80
	*/
	if len(plMapCall) != 0 {
		_, err = GeneratePickingList(plMapCall, w, limitSalesOrder, limitMaxWeight)
	}

	if len(plMapCall) == 0 {
		plMapReturn = plMapFinal
	}

	plMapFinal = nil
	return plMapReturn, nil
}

func InsertOnePickingList(weight float64, listSOID []int64, w *model.Warehouse, plFinal map[string]ListPl) (map[string]ListPl, error) {
	code, _ := util.GenerateDocCode("PL", w.Code, "picking_list")
	plFinal[code] = ListPl{
		SalesOrderID: listSOID,
		TotalWeight:  weight,
	}

	return plFinal, nil
}

func GetFlagDisableForProducts(totalCancelledSalesOrder, currentStatusSalesOrder int8) (int8, error) {
	var e error

	if currentStatusSalesOrder == 3 || currentStatusSalesOrder == 4 {
		return totalCancelledSalesOrder + 1, e
	}

	return totalCancelledSalesOrder, e

}

func GetSavedPickForProducts(totalSavedPickQty, flagSavedQty int8) (int8, error) {
	var e error

	if flagSavedQty == 1 {
		return totalSavedPickQty + 1, e
	}

	return totalSavedPickQty, e

}
