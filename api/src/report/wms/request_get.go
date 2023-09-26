package wms

import (
	"reflect"
	"strconv"
	"strings"
	"time"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// getStockLog : query to get data of report stock log
func getStockLog(cond map[string]interface{}) (sl []*reportStockLog, e error) {
	var where string
	var values []interface{}

	o := orm.NewOrm()
	o.Using("read_only")

	for k, v := range cond {
		if reflect.TypeOf(v).Kind().String() == "slice" {
			where = where + " " + k + " >= ? and " + k + " <= ? and"
		} else {
			where = where + " " + k + "? and"
		}

		values = append(values, v)
	}
	where = strings.TrimSuffix(where, " and")

	q := "SELECT sl.created_at AS Timestamp , sl.type AS Log_Type ,sl.ref_type AS Ref_Type , " +
		"CASE WHEN sl.ref_type = 1 THEN do.code " +
		"WHEN sl.ref_type = 2 THEN dr.code " +
		"WHEN sl.ref_type = 3 THEN gr.code " +
		"WHEN sl.ref_type = 4 THEN gt.code " +
		"WHEN sl.ref_type = 5 THEN so.code " +
		"WHEN sl.ref_type = 6 THEN we.code " +
		"WHEN sl.ref_type = 7 THEN ts.code END AS Reference_Code , " +
		"p.code AS Product_Code , p.name AS Product_Name , u.name AS UOM , " +
		"sl.initial_stock AS Initial_Stock , sl.quantity AS Quantity , sl.final_stock AS Final_Stock , " +
		"w.name AS Warehouse , a.name AS Area , sl.status AS Status , " +
		"CASE WHEN sl.ref_type = 1 THEN do.note " +
		"WHEN sl.ref_type = 2 THEN dr.note " +
		"WHEN sl.ref_type = 3 THEN gr.note " +
		"WHEN sl.ref_type = 4 THEN gt.note " +
		"WHEN sl.ref_type = 5 THEN so.note " +
		"WHEN sl.ref_type = 6 THEN we.note END AS Doc_Note , " +
		"sl.item_note AS Note " +
		"FROM stock_log sl " +
		"LEFT JOIN delivery_order do ON do.id = sl.ref_id and sl.ref_type = 1 " +
		"LEFT JOIN delivery_return dr ON dr.id = sl.ref_id and sl.ref_type = 2 " +
		"LEFT JOIN goods_receipt gr ON gr.id = sl.ref_id and sl.ref_type = 3 " +
		"LEFT JOIN goods_transfer gt ON gt.id = sl.ref_id and sl.ref_type = 4 " +
		"LEFT JOIN stock_opname so ON so.id = sl.ref_id and sl.ref_type = 5 " +
		"LEFT JOIN waste_entry we ON we.id = sl.ref_id and sl.ref_type = 6 " +
		"LEFT JOIN transfer_sku ts ON ts.id = sl.ref_id and sl.ref_type = 7 " +
		"JOIN product p ON p.id = sl.product_id JOIN uom u ON u.id = p.uom_id " +
		"JOIN warehouse w ON w.id = sl.warehouse_id JOIN area a ON a.id = w.area_id " +
		"WHERE " + where

	_, e = o.Raw(q, values).QueryRows(&sl)

	return
}

// getWasteLog : query to get data of report waste log
func getWasteLog(cond map[string]interface{}) (sl []*reportWasteLog, e error) {
	var where string
	var values []interface{}

	o := orm.NewOrm()
	o.Using("read_only")

	for k, v := range cond {
		if reflect.TypeOf(v).Kind().String() == "slice" {
			where = where + " " + k + " >= ? and " + k + " <= ? and"
		} else {
			where = where + " " + k + "? and"
		}

		values = append(values, v)
	}
	where = strings.TrimSuffix(where, " and")

	q := "SELECT wl.created_at AS Timestamp, " +
		"wl.type AS Log_Type, " +
		"wl.ref_type as Ref_Type, " +
		"CASE " +
		" WHEN wl.ref_type = 1 THEN dr.code " +
		" WHEN wl.ref_type = 2 THEN we.code " +
		" WHEN wl.ref_type = 3 THEN wd.code " +
		" WHEN wl.ref_type = 4 THEN ts.code " +
		" WHEN wl.ref_type = 5 THEN gt.code " +
		" WHEN wl.ref_type = 6 THEN gr.code " +
		" WHEN wl.ref_type = 7 THEN do.code " +
		" WHEN wl.ref_type = 8 THEN op.code " +
		"END AS Reference_Code, " +
		"CASE " +
		" WHEN wl.ref_type = 8 THEN glso.value_name " +
		" ELSE gl.value_name " +
		"END AS Waste_Reason, " +
		"gt.code AS Good_Transfer_Code, " +
		"gr.code AS Good_Receipt_Code, " +
		"po.code AS Purchase_Order_Code, " +
		"sl.name AS Suplier_Name, " +
		"slty.name AS Suplier_Type, " +
		"wo.name AS Warehouse_Origin, " +
		"p.code AS Product_Code, " +
		"p.name AS Product_Name, " +
		"u.name AS UOM, " +
		"wl.quantity AS Quantity, " +
		"wl.final_stock AS Final_Stock, " +
		"w.name AS Warehouse, " +
		"a.name AS Area, " +
		"CASE " +
		" WHEN wl.ref_type = 1 THEN dr.note " +
		" WHEN wl.ref_type = 2 THEN we.note " +
		" WHEN wl.ref_type = 3 THEN wd.note " +
		" WHEN wl.ref_type = 4 THEN ts.note " +
		" WHEN wl.ref_type = 5 THEN gt.note " +
		" WHEN wl.ref_type = 6 THEN gr.note " +
		" WHEN wl.ref_type = 7 THEN do.note " +
		" WHEN wl.ref_type = 8 THEN op.note " +
		"END AS Doc_Note, " +
		"wl.item_note AS Item_Note " +
		"FROM waste_log wl " +
		"LEFT JOIN delivery_return dr ON dr.id = wl.ref_id AND wl.ref_type = 1  " +
		"LEFT JOIN waste_entry we ON we.id = wl.ref_id AND wl.ref_type = 2 " +
		"LEFT JOIN waste_disposal wd ON wd.id = wl.ref_id AND wl.ref_type = 3 " +
		"LEFT JOIN transfer_sku ts ON ts.id = wl.ref_id AND wl.ref_type = 4 " +
		"LEFT JOIN goods_transfer gt ON gt.id = wl.ref_id AND wl.ref_type = 5 " +
		"LEFT JOIN goods_receipt gr ON gr.id = wl.ref_id AND wl.ref_type = 6 " +
		"LEFT JOIN delivery_order `do` ON `do`.id = wl.ref_id AND wl.ref_type = 7 " +
		"LEFT JOIN stock_opname op ON op.id = wl.ref_id AND wl.ref_type = 8 " +
		"LEFT JOIN purchase_order po ON po.id = gr.purchase_order_id  " +
		"LEFT JOIN supplier sl ON sl.id = po.supplier_id  " +
		"LEFT JOIN supplier_type slty ON slty.id = sl.supplier_type_id  " +
		"LEFT JOIN warehouse wo ON wo.id = gt.origin_id  " +
		"LEFT JOIN glossary gl ON gl.value_int = wl.waste_reason AND gl.attribute = 'waste_reason' AND gl.`table` = 'all' " +
		"LEFT JOIN glossary glso ON glso.value_int = wl.waste_reason AND glso.attribute = 'opname_reason' AND glso.`table` = 'stock_opname' " +
		"JOIN product p ON p.id = wl.product_id " +
		"JOIN uom u ON u.id = p.uom_id " +
		"JOIN warehouse w ON w.id = wl.warehouse_id " +
		"JOIN area a ON a.id = w.area_id " +
		"WHERE " + where + " ORDER BY wl.id"

	_, e = o.Raw(q, values).QueryRows(&sl)

	for _, v := range sl {
		logTypeGlossary, _ := repository.GetGlossaryMultipleValue("table", "waste_log", "attribute", "type", "value_int", v.LogType)
		v.LogType = strings.ToUpper(logTypeGlossary.ValueName)

		refTypeGlossary, _ := repository.GetGlossaryMultipleValue("table", "waste_log", "attribute", "ref_type", "value_int", v.RefType)
		v.RefType = refTypeGlossary.ValueName

		if v.RefType == "transfer_sku" {
			transferSku := &model.TransferSku{Code: v.ReferenceCode}
			transferSku.Read("code")

			if transferSku.GoodsReceipt != nil {
				transferSku.GoodsReceipt.Read("id")
				v.GoodReceiptCode = transferSku.GoodsReceipt.Code
			}
			if transferSku.GoodsTransfer != nil {
				transferSku.GoodsTransfer.Read("id")
				v.GoodTransferCode = transferSku.GoodsTransfer.Code

				warehouse := &model.Warehouse{ID: transferSku.GoodsTransfer.Origin.ID}
				warehouse.Read("id")
				v.WarehouseOrigin = warehouse.Name
			}
			if transferSku.PurchaseOrder != nil {
				transferSku.PurchaseOrder.Read("id")
				transferSku.PurchaseOrder.Supplier.Read("id")
				transferSku.PurchaseOrder.Supplier.SupplierType.Read("id")

				v.PurchaseOrderCode = transferSku.PurchaseOrder.Code
				v.SuplierName = transferSku.PurchaseOrder.Supplier.Name
				v.SuplierType = transferSku.PurchaseOrder.Supplier.SupplierType.Name
			}
		}

		if v.RefType == "goods_receipt" {
			goodsReceipt := &model.GoodsReceipt{Code: v.ReferenceCode}
			goodsReceipt.Read("code")

			if goodsReceipt.GoodsTransfer != nil {
				goodsReceipt.GoodsTransfer.Read("id")
				v.GoodTransferCode = goodsReceipt.GoodsTransfer.Code

				warehouse := &model.Warehouse{ID: goodsReceipt.GoodsTransfer.Origin.ID}
				warehouse.Read("id")
				v.WarehouseOrigin = warehouse.Name
			}
		}

	}

	return
}

// getStock : query to get data of stock
func getStock(cond map[string]interface{}) (s []*reportStock, e error) {
	var where string
	var values []interface{}

	o := orm.NewOrm()
	o.Using("read_only")

	for k, v := range cond {
		if reflect.TypeOf(v).Kind().String() == "slice" {
			where = where + " " + k + " >= ? and " + k + " <= ? and"
		} else {
			where = where + " " + k + "? and"
		}

		values = append(values, v)
	}
	where = strings.TrimSuffix(where, " and")
	if where != "" {
		where = "AND " + where
	}

	q := "SELECT p.code AS product_code , p.name AS product_name , w.name AS warehouse_name , stock.available_stock AS available_stock , " +
		"stock.waste_stock AS waste_stock , stock.safety_stock AS safety_stock , stock.commited_in_stock AS commited_in_stock , " +
		"stock.commited_out_stock AS commited_out_stock , stock.expected_qty AS expected_qty , stock.received_qty AS received_qty , " +
		"stock.intransit_qty AS intransit_qty , stock.intransit_waste_qty AS intransit_waste_qty , " +
		"case when stock.salable = 1 then 'Yes' else 'No' end AS salable , case when stock.purchasable = 1 then 'Yes' else 'No' end AS purchasable , " +
		"case when stock.status = 1 then 'Yes' else 'No' end AS status " +
		"FROM stock JOIN product p ON p.id = stock.product_id " +
		"JOIN warehouse w ON w.id = stock.warehouse_id " +
		"WHERE stock.status=1 " + where

	_, e = o.Raw(q, values).QueryRows(&s)

	return
}

// getGoodsReceiptItem : query to get data of goods receipt item
func getGoodsReceiptItem(cond map[string]interface{}) (gri []*reportGoodsReceiptItem, e error) {
	var where string
	var values []interface{}

	o := orm.NewOrm()
	o.Using("read_only")

	for k, v := range cond {
		if reflect.TypeOf(v).Kind().String() == "slice" {
			where = where + " " + k + " >= ? and " + k + " <= ? and"
		} else {
			where = where + " " + k + "? and"
		}

		values = append(values, v)
	}
	where = strings.TrimSuffix(where, " and")

	q := "SELECT gri2.* , " +
		"CASE WHEN tsi2.transfer_sku_id IS NULL THEN gri2.request_order_qty " +
		"ELSE IF(tsi2.product_id != tsi2.transfer_product_id, 0, gri2.request_order_qty) END AS Ordered_Qty , " +
		"CASE WHEN gri2.gt_deliver_qty IS NULL THEN NULL WHEN tsi2.transfer_sku_id IS NULL THEN gri2.gt_deliver_qty " +
		"ELSE IF(tsi2.product_id != tsi2.transfer_product_id, 0, gri2.gt_deliver_qty) END AS Transfer_qty , " +
		"CASE WHEN tsi2.transfer_sku_id IS NULL THEN gri2.gr_deliver_qty ELSE IF(tsi2.product_id != tsi2.transfer_product_id, 0, gri2.gr_deliver_qty) END AS Delivered_Qty , " +
		"CASE WHEN tsi2.transfer_sku_id IS NULL THEN gri2.gr_reject_qty ELSE IF(tsi2.product_id != tsi2.transfer_product_id, 0, gri2.gr_reject_qty) END AS Reject_Qty , " +
		"CASE WHEN tsi2.transfer_sku_id IS NULL THEN gri2.gr_received_qty ELSE IF(tsi2.product_id != tsi2.transfer_product_id, 0, gri2.gr_received_qty) END AS Received_Qty , " +
		"CASE WHEN tsi2.transfer_sku_id IS NULL THEN gri2.sr_return_qty ELSE IF(tsi2.product_id != tsi2.transfer_product_id, 0, gri2.sr_return_qty) END AS Return_Qty , " +
		"tsi2.product_down AS Product_DownGrade ,tsi2.transfer_qty AS After_Sortir_Good_Qty , tsi2.waste_qty AS After_Sortir_Waste_Qty , tsi2.downgrade_qty AS After_Sortir_DownGrade_Qty " +
		"FROM( SELECT " +
		"IF(po.code IS NULL, gt.code, po.code) AS Inbound_Code , s.code AS Supplier_Code, s.name AS Supplier_Name , upper(" +
		"IF(g2.value_name is null, g3.value_name, g2.value_name)) AS Inbound_Status , gr.code AS GR_Code, upper( g.value_name) AS GR_Status , sr.code AS SR_Code, upper( g4.value_name ) AS SR_Status , dn.code AS DN_Code, upper( g5.value_name ) AS DN_Status , " +
		"w2.name AS Warehouse_Origin , IF(w.name IS NULL, w3.name, w.name) AS Warehouse_Destination , " +
		"IF(po.eta_date IS NULL, gt.eta_date, po.eta_date) AS Estimation_Arrival_Date , " +
		"IF(po.eta_time IS NULL, gt.eta_time , po.eta_time) AS Estimation_Arrival_Time , gr.ata_date AS Actual_Arrival_Date , gr.ata_time AS Actual_Arrival_Time , gr.note AS GR_Note , p.code AS Product_Code , p.name AS Product_Name , u.name AS UOM, gri.note AS GR_Item_Note , " +
		"IF(poi.order_qty IS NULL, gti.request_qty, poi.order_qty) AS request_order_qty , gti.deliver_qty AS gt_deliver_qty , gri.deliver_qty AS gr_deliver_qty , gri.reject_qty AS gr_reject_qty , gri.receive_qty AS gr_received_qty , sri.return_good_qty AS sr_return_qty , ts.code AS TS_Code , " +
		"upper( g6.value_name ) AS TS_Status ,gri.product_id 'gri_product_id' ,ts.ts_id " +
		"FROM goods_receipt_item gri JOIN goods_receipt gr ON gr.id = gri.goods_receipt_id " +
		"LEFT JOIN purchase_order po ON po.id = gr.purchase_order_id " +
		"LEFT JOIN purchase_order_item poi ON poi.purchase_order_id = po.id AND poi.product_id = gri.product_id " +
		"LEFT JOIN goods_transfer gt ON gt.id = gr.goods_transfer_id " +
		"LEFT JOIN goods_transfer_item gti ON gti.goods_transfer_id = gt.id AND gti.product_id = gri.product_id " +
		"LEFT JOIN supplier s ON s.id = po.supplier_id JOIN product p ON p.id = gri.product_id " +
		"JOIN uom u ON p.uom_id = u.id " +
		"LEFT JOIN supplier_return sr ON sr.goods_receipt_id = gr.id AND sr.status IN (1,2) " +
		"LEFT JOIN supplier_return_item sri ON sri.supplier_return_id = sr.id AND sri.product_id = p.id " +
		"LEFT JOIN debit_note dn ON dn.supplier_return_id = sr.id AND dn.status IN (1,2) " +
		"LEFT JOIN warehouse w ON w.id = po.warehouse_id LEFT JOIN warehouse w2 ON w2.id = gt.origin_id " +
		"LEFT JOIN warehouse w3 ON w3.id = gt.destination_id " +
		"LEFT JOIN (SELECT ts.id 'ts_id', ts.code, ts.goods_receipt_id, ts.status ,tsi.* " +
		"FROM transfer_sku ts JOIN transfer_sku_item tsi ON tsi.transfer_sku_id = ts.id " +
		"GROUP BY ts.id, tsi.product_id, tsi.transfer_product_id ) ts ON ts.goods_receipt_id = gr.id AND ts.product_id = gri.product_id " +
		"LEFT JOIN glossary g ON g.value_int = gr.status AND g.`table` = 'goods_receipt' AND g.attribute = 'status' " +
		"LEFT JOIN glossary g2 ON g2.value_int = gt.status AND g2.`table` = 'goods_transfer' AND g2.attribute = 'status' " +
		"LEFT JOIN glossary g3 ON g3.value_int = po.status AND g3.`table` = 'purchase_order' AND g3.attribute = 'status' " +
		"LEFT JOIN glossary g4 ON g4.value_int = sr.status AND g4.`table` = 'supplier_return' AND g4.attribute = 'status' " +
		"LEFT JOIN glossary g5 ON g5.value_int = dn.status AND g5.`table` = 'debit_note' AND g5.attribute = 'status' " +
		"LEFT JOIN glossary g6 ON g6.value_int = ts.status AND g6.`table` = 'transfer_sku' AND g6.attribute = 'status' " +
		"WHERE  " + where +
		" GROUP BY gr.id, gri.product_id ) gri2 LEFT JOIN (SELECT tsi.transfer_sku_id , tsi.product_id , tsi.transfer_product_id ,IF(tsi.product_id != tsi.transfer_product_id, 0,tsi.transfer_qty) 'transfer_qty' , tsi.waste_qty ,IF(tsi.product_id = tsi.transfer_product_id, 0,tsi.transfer_qty) 'downgrade_qty' , IF(tsi.product_id = tsi.transfer_product_id , NULL, p2.name) 'product_down' FROM transfer_sku_item tsi " +
		"JOIN product p2 ON p2.id = tsi.transfer_product_id) tsi2 ON tsi2.transfer_sku_id = gri2.ts_id AND tsi2.product_id = gri2.gri_product_id;"

	_, e = o.Raw(q, values).QueryRows(&gri)

	return
}

// getDeliveryReturnItem : query to get data of delivery return item
func getDeliveryReturnItem(cond map[string]interface{}) (dri []*reportDeliveryReturnItem, e error) {
	var where string
	var values []interface{}

	o := orm.NewOrm()
	o.Using("read_only")

	for k, v := range cond {
		if reflect.TypeOf(v).Kind().String() == "slice" {
			where = where + " " + k + " >= ? and " + k + " <= ? and"
		} else {
			where = where + " " + k + "? and"
		}

		values = append(values, v)
	}
	where = strings.TrimSuffix(where, " and")

	q := "SELECT dr.recognition_date AS Return_Date, p.code AS Product_Code, p.name AS Product_Name, u.name AS Unit," +
		" dri.return_good_qty AS Good_Stock_Return_Qty, dri.return_waste_qty AS Waste_Return_Qty," +
		" (dri.return_good_qty+dri.return_waste_qty) AS Total_Return_Qty, soi.unit_price AS Product_Price, a.name AS Area," +
		" w.name AS Warehouse, so.code AS Order_Code, do.code AS Delivery_Code, do.recognition_date AS Delivery_Date," +
		" b.code AS Customer_Code, b.name AS Customer_Name, dr.note AS Delivery_Return_Note, dri.note AS Delivery_Return_Item_Note " +
		"FROM delivery_return dr " +
		"JOIN delivery_return_item dri ON dri.delivery_return_id = dr.id " +
		"JOIN product p ON p.id = dri.product_id " +
		"JOIN uom u ON u.id = p.uom_id " +
		"JOIN warehouse w ON w.id = dr.warehouse_id " +
		"JOIN area a ON a.id = w.area_id " +
		"JOIN delivery_order do ON do.id = dr.delivery_order_id " +
		"JOIN delivery_order_item doi ON doi.id = dri.delivery_order_item_id " +
		"JOIN sales_order so ON so.id = do.sales_order_id " +
		"JOIN sales_order_item soi ON soi.sales_order_id = so.id AND soi.id = doi.sales_order_item_id " +
		"JOIN branch b ON b.id = so.branch_id " +
		"WHERE " + where

	_, e = o.Raw(q, values).QueryRows(&dri)

	return
}

// getProducts : query to get data of products
func getProducts(cond map[string]interface{}) (p []*reportProducts, e error) {
	var where string
	var values []interface{}

	o := orm.NewOrm()
	o.Using("read_only")

	for k, v := range cond {
		if reflect.TypeOf(v).Kind().String() == "slice" {
			where = where + " " + k + " >= ? and " + k + " <= ? and"
		} else {
			if strings.Contains(k, "like") {
				str, ok := v.(string)
				if ok {
					str = "%" + str + "%"
				}
				v = str
			} else if strings.Contains(k, "=") {
				where = where + " " + k + "? and"
			} else {
				where = where + " find_in_set (? , " + k + ") and"
			}
		}

		values = append(values, v)
	}
	where = strings.TrimSuffix(where, " and")
	if where != "" {
		where = "WHERE " + where
	}

	q := "SELECT p.code AS product_code, p.name AS product_name ,c.code 'Code category - C2' " +
		",c.name 'Category - C2' ,c2.code 'Code category - C1' ,c2.name 'Category - C1' " +
		",c3.code 'Code category - C0' ,c3.name 'Category - C0' , u.name AS UOM, p.order_min_qty AS Minimal_Order_Qty" +
		",p.unit_weight AS total_weight , p.note AS product_note, p.description AS product_description" +
		", tp.name AS product_tag , g.value_name AS product_status, sal.name AS warehouse_salability " +
		", pur.name AS warehouse_purchasability, sto.name AS warehouse_storability , p.spare_percentage " +
		"FROM eden_v2.product p " +
		"LEFT JOIN( SELECT p.id AS id, p.tag_product AS tag_product, group_concat( tp.name SEPARATOR ',') AS name " +
		"FROM( eden_v2.product p JOIN eden_v2.tag_product tp) " +
		"WHERE find_in_set( tp.value, p.tag_product ) GROUP BY p.id ) tp ON tp.id = p.id " +
		"JOIN eden_v2.category c on c.id = p.category_id " +
		"LEFT JOIN eden_v2.category c2 on c2.id = c.parent_id " +
		"LEFT JOIN eden_v2.category c3 on c3.id = c2.grandparent_id " +
		"LEFT JOIN ( SELECT p.id AS id, p.warehouse_sal AS warehouse_sal, group_concat( w.name SEPARATOR ',' ) AS name " +
		"FROM ( eden_v2.product p JOIN eden_v2.warehouse w ) " +
		"WHERE find_in_set( w.id, p.warehouse_sal ) GROUP BY p.id ) sal ON sal.id = p.id " +
		"LEFT JOIN ( SELECT p.id AS id, p.warehouse_pur AS warehouse_pur, group_concat( w.name SEPARATOR ',' ) AS name " +
		"FROM ( eden_v2.product p JOIN eden_v2.warehouse w ) " +
		"WHERE find_in_set( w.id, p.warehouse_pur ) GROUP BY p.id ) pur ON pur.id = p.id " +
		"LEFT JOIN ( SELECT p.id AS id, p.warehouse_sto AS warehouse_sto, group_concat( w.name SEPARATOR ',' ) AS name " +
		"FROM ( eden_v2.product p JOIN eden_v2.warehouse w ) " +
		"WHERE find_in_set( w.id, p.warehouse_sto ) GROUP BY p.id ) sto ON sto.id = p.id " +
		"JOIN eden_v2.uom u ON u.id = p.uom_id " +
		"LEFT JOIN eden_v2.glossary g ON g.value_int = p.status AND g.attribute = 'status' " + where

	_, e = o.Raw(q, values).QueryRows(&p)

	return
}

// getDeliveryOrder : query to get data of delivery order
func getDeliveryOrder(rq *orm.RequestQuery, cond map[string]interface{}) (m []*reportDeliveryOrder, err error) {

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

	q := "SELECT w.name warehouse , m.name merchant_name , bt.name business_type , so.code order_code , " +
		"do.code delivery_code , upper(g.value_name) delivery_status , adm.province_name province , " +
		"adm.city_name city , adm.district_name district , adm.sub_district_name sub_district , adm.postal_code postal_code , " +
		"wrt.name wrt , so.total_weight order_weight , so.delivery_date delivery_date , tps.name payment_term , " +
		"os.name sales_order_type , upper( gg.value_name) AS sales_order_status , so.shipping_address shipping_address , " +
		"m.tag_customer, a.name area_name " +
		"FROM sales_order so " +
		"LEFT JOIN delivery_order do ON do.sales_order_id = so.id " +
		"JOIN order_type_sls os ON os.id = so.order_type_sls_id " +
		"JOIN adm_division adm ON adm.sub_district_id = so.sub_district_id " +
		"JOIN wrt ON wrt.id = so.wrt_id " +
		"JOIN warehouse w ON w.id = so.warehouse_id " +
		"JOIN term_payment_sls tps ON tps.id = so.term_payment_sls_id " +
		"LEFT JOIN glossary g ON g.value_int = do.status AND g.table='delivery_order' AND g.attribute= 'status' " +
		"JOIN branch b on b.id = so.branch_id " +
		"JOIN(SELECT m.id, m.name , business_type_id, GROUP_CONCAT(tc.name) 'tag_customer' " +
		"FROM merchant m " +
		"LEFT JOIN tag_customer tc ON FIND_IN_SET(tc.id, m.tag_customer) " +
		"GROUP BY m.id) m ON m.id = b.merchant_id " +
		"JOIN business_type bt on bt.id = m.business_type_id " +
		"JOIN area a ON a.id = so.area_id " +
		"LEFT JOIN glossary gg ON gg.value_int = so.status AND gg.attribute = 'doc_status' " +
		"WHERE so.status != 4 " + where +
		" ORDER BY so.delivery_date desc"

	_, err = o.Raw(q, values).QueryRows(&m)

	return
}

// getItemRecap : query to get data of item recap
func getItemRecap(cond map[string]interface{}) (ir []*reportItemRecap, e error) {
	var where string
	var values []interface{}

	o := orm.NewOrm()
	o.Using("read_only")

	for k, v := range cond {
		if reflect.TypeOf(v).Kind().String() == "slice" {
			where = where + " " + k + " >= ? and " + k + " <= ? and"
		} else {
			where = where + " " + k + "? and"
		}

		values = append(values, v)
	}
	where = strings.TrimSuffix(where, " and")
	if where != "" {
		where = "AND " + where
	}

	q := "select so.delivery_date order_delivery_date , a.name area, w.name warehouse , p.code product_code," +
		"p.name product_name, iu.name uom , sum(IF( ots.`name` NOT IN ('zero waste'),soi.order_qty, 0 )) total_qty, sum(IF( ots.`name` = 'zero waste',soi.order_qty, 0 )) total_quantity_zero_waste, sum(soi.weight) total_weight ," +
		"c.code 'Code category - C2' ,c.name 'Category - C2' ,c2.code 'Code category - C1' ," +
		"c2.name 'Category - C1' ,c3.code 'Code category - C0' ,c3.name 'Category - C0' " +
		"from sales_order so " +
		"join sales_order_item soi on soi.sales_order_id = so.id " +
		"join area a on a.id = so.area_id " +
		"join warehouse w on w.id =so.warehouse_id " +
		"join product p on p.id = soi.product_id " +
		"join uom iu on iu.id = p.uom_id " +
		"JOIN category c on c.id = p.category_id " +
		"LEFT JOIN category c2 on c2.id = c.parent_id " +
		"LEFT JOIN category c3 on c3.id = c2.grandparent_id " +
		"join order_type_sls ots on so.order_type_sls_id = ots.id " +
		"where((so.payment_group_sls_id = 1 and so.status not in (1,3,4)) " +
		"or ( so.payment_group_sls_id != 1 and so.status not in (3,4))) " +
		"and ots.value != 'draft' " + where + " group by w.id, p.id"

	_, e = o.Raw(q, values).QueryRows(&ir)

	return
}

// getMovementStock : query to get data of movement stock
func getMovementStock(cond map[string]interface{}) (ms []*reportMovementStock, e error) {

	var whID = strconv.Itoa(int(cond["w.id"].(int64)))
	var dateValue = cond["recognition_date"].(string)

	dateParsing, e := time.Parse("2006-01-02", dateValue)

	dateParsingMinOneDay := dateParsing.Add(-1).Format("2006-01-02")

	o := orm.NewOrm()
	o.Using("read_only")

	q := "select p.code 'product_code', p.name 'product_name', c.name 'category', u.name 'uom', coalesce(stp.final_stock,0) 'stock', coalesce(po.sum,0) 'plan_inbound'," +
		"coalesce(gr.sum,0) 'actual_inbound', coalesce(we.sum,0) 'waste', coalesce(so.sum,0) 'plan_delivery', coalesce(do.sum,0) 'actual_delivery', coalesce(gt1.sum,0) 'stock_transfer_in', coalesce(gt2.sum,0) 'stock_transfer_out', coalesce(dr.sum,0) 'goods_return', (coalesce(stp.final_stock,0) + coalesce(gr.sum,0) - coalesce(we.sum,0) - coalesce(do.sum,0) + coalesce(gt1.sum,0) - coalesce(gt2.sum,0) + coalesce(dr.sum,0)) as 'stock_akhir', s.available_stock 'actual_stock', (coalesce(s.available_stock,0) - (coalesce(stp.final_stock,0) + coalesce(gr.sum,0) - coalesce(we.sum,0) - coalesce(do.sum,0) + coalesce(gt1.sum,0) - coalesce(gt2.sum,0) + coalesce(dr.sum,0))) as 'selisih_stock' " +
		"from (select id, name, category_id, uom_id, code from product where status = 1) p " +
		"join category c on c.id = p.category_id " +
		"join uom u on u.id = p.uom_id " +
		"left join (select id, warehouse_id, available_stock, product_id from stock s where s.warehouse_id = " + whID + " ) s on s.product_id = p.id " +
		"left join (SELECT t.product_id, t.final_stock FROM (SELECT soi.product_id,  soi.final_stock, IF(@prev <> soi.product_id, @rn:=1,@rn), @prev:= soi.product_id, @rn:=@rn+1 AS rn FROM stock_opname_item soi JOIN stock_opname so on so.id = soi.stock_opname_id, (SELECT @rn:=0) rn, (SELECT @prev:='') prev WHERE so.status = 2 AND so.recognition_date = '" + dateParsingMinOneDay + "' and so.warehouse_id = " + whID + " ORDER BY soi.product_id ASC, soi.id DESC) t group by t.product_id having max(rn) ) stp on stp.product_id = p.id " +
		"left join (select sum(poi.order_qty) 'sum', poi.product_id 'product_id' from purchase_order_item poi " +
		"join purchase_order po on po.id = poi.purchase_order_id and po.status in (1,2) and po.recognition_date = '" + dateValue + "' and po.warehouse_id = " + whID + " group by poi.product_id) po on po.product_id = p.id " +
		"left join(select sum(gri.receive_qty) 'sum', gri.product_id 'product_id' from goods_receipt_item gri " +
		"join goods_receipt gr on gr.id = gri.goods_receipt_id and gr.status in (2) and gr.ata_date = '" + dateValue + "' and gr.warehouse_id = " + whID + " group by gri.product_id) gr on gr.product_id = p.id " +
		"left join (select sum(wei.waste_qty) 'sum', wei.product_id 'product_id' from waste_entry_item wei " +
		"join waste_entry we on we.id = wei.waste_entry_id and we.status in (2) and we.recognition_date = '" + dateValue + "' and we.warehouse_id = " + whID + " group by wei.product_id) we on we.product_id = p.id " +
		"left join (select sum(soi.order_qty) 'sum', soi.product_id 'product_id' from sales_order_item soi " +
		"join sales_order so on so.id = soi.sales_order_id and so.status in (1,9,12) and so.recognition_date = '" + dateValue + "' and so.warehouse_id = " + whID + " group by soi.product_id) so on so.product_id = p.id " +
		"left join (select sum(doi.receive_qty) 'sum', doi.product_id 'product_id' from delivery_order_item doi " +
		"join delivery_order do on do.id = doi.delivery_order_id and do.status in (2) and do.recognition_date = '" + dateValue + "' and do.warehouse_id = " + whID + " group by doi.product_id) do on do.product_id = p.id " +
		"left join (select sum(gti.receive_qty) 'sum', gti.product_id 'product_id' from goods_transfer_item gti " +
		"join goods_transfer gt on gt.id = gti.goods_transfer_id and gt.status in (2) and gt.destination_id = " + whID + " and gt.recognition_date = '" + dateValue + "' group by gti.product_id) gt1 on gt1.product_id = p.id " +
		"left join (select sum(gti2.deliver_qty) 'sum', gti2.product_id 'product_id' from goods_transfer_item gti2 " +
		"join goods_transfer gt2 on gt2.id = gti2.goods_transfer_id and gt2.status in (1) and gt2.origin_id = " + whID + " and gt2.recognition_date = '" + dateValue + "' group by gti2.product_id) gt2 on gt2.product_id = p.id " +
		"left join (select sum(dri.return_good_qty) 'sum', dri.product_id 'product_id' from delivery_return_item dri " +
		"join delivery_return dr on dr.id = dri.delivery_return_id and dr.status in (2) and dr.recognition_date = '" + dateValue + "' and dr.warehouse_id = " + whID + " group by dri.product_id) dr on dr.product_id = p.id "

	_, e = o.Raw(q).QueryRows(&ms)

	return
}

// getPicking : query to get data of report picking
func getPicking(cond map[string]interface{}) (pc []*reportPicking, e error) {
	var where string
	var values []interface{}

	o := orm.NewOrm()
	o.Using("read_only")

	// condition for getting values regard on time format
	// because there is 6 time stamp in
	var timeFormat = "%H:%i:%s"
	for i := 0; i < 6; i++ {
		values = append(values, timeFormat)
	}

	for k, v := range cond {
		if reflect.TypeOf(v).Kind().String() == "slice" {
			where = where + " " + k + "? and ? and"
		} else {
			where = where + " " + k + "? and"
		}

		values = append(values, v)
	}
	where = strings.TrimSuffix(where, " and")

	q1 := "SELECT so.delivery_date , TIME_FORMAT(poa.assign_timestamp, ?) 'timestamp assign' , pl.code 'pl_code', " +
		"so.code 'so code', ots.name 'order_type' , tps.name 'payment_term' , m.name 'merchant' , bt.name 'business_type', " +
		"COUNT(poi.id)'total item' , so.total_weight 'sales_order_weight' ,  sum(poi.check_qty) 'total_weight', " +
		"poa.total_koli ,so.shipping_address , wrt.name 'wrt', w.name 'warehouse' , s.name 'Picker', " +
		"TIME_FORMAT(poa.checkin_timestamp, ?) 'Time Start Picked' ," +
		"TIME_FORMAT(poa.checkout_timestamp , ?) 'Time Finish Picked' ,s2.name 'Checker' ," +
		"TIME_FORMAT(poa.checker_in_timestamp , ?) 'Time Checkin' ," +
		"TIME_FORMAT(poa.checker_out_timestamp , ?) 'Time Checkout' ,cv.name 'vendor' ," +
		"poa.planning_vendor 'planning' ,c.name 'courier' ,TIME_FORMAT(poa.dispatch_timestamp , ?) 'Dispatch Time' , " +
		"g2.value_name 'status picking assigned' , g.value_name 'status so' " +
		"FROM sales_order so " +
		"LEFT JOIN picking_order_assign poa ON so.id = poa.sales_order_id " +
		"LEFT JOIN picking_list pl on poa.picking_list_id = pl.id " +
		"LEFT JOIN picking_order po ON po.id = poa.picking_order_id " +
		"LEFT JOIN picking_order_item poi ON poi.picking_order_assign_id = poa.id " +
		"JOIN wrt ON wrt.id = so.wrt_id " +
		"JOIN branch b ON b.id = so.branch_id " +
		"JOIN merchant m ON m.id = b.merchant_id " +
		"JOIN warehouse w ON w.id = so.warehouse_id " +
		"JOIN order_type_sls ots ON so.order_type_sls_id=ots.id " +
		"JOIN term_payment_sls tps ON so.term_payment_sls_id=tps.id " +
		"JOIN business_type bt ON m.business_type_id=bt.id " +
		"LEFT JOIN staff s ON s.id = poa.staff_id " +
		"LEFT JOIN staff s2 ON s2.id = poa.checked_by " +
		"LEFT JOIN courier c ON c.id = poa.courier_id " +
		"LEFT JOIN courier_vendor cv ON cv.id = poa.courier_vendor_id " +
		"LEFT JOIN glossary g ON g.`table` = 'sales_order' and g.`attribute` = 'status' and g.value_int = so.status " +
		"LEFT JOIN glossary g2 ON g2.`table` = 'picking_order' and g2.`attribute` = 'doc_status_picking' and g2.value_int = poa.status " +
		"WHERE " + where +
		" GROUP BY so.code ORDER BY so.delivery_date"

	_, e = o.Raw(q1, values).QueryRows(&pc)

	return
}

// getPickingOrderItem : query to get data of report picking order item
func getPickingOrderItem(cond map[string]interface{}) (poi []*reportPickingOrderItem, e error) {
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

	q1 := "SELECT so.delivery_date , " +
		"pl.code 'pl_code' , " +
		"so.code 'so code' , " +
		"m.name 'merchant' , " +
		"p.code 'product_code' , " +
		"p.name 'product_name' , " +
		"u.name 'uom' , " +
		"poi.order_qty 'order_qty' , " +
		"poi.pick_qty 'qty_picker' , " +
		"poi.check_qty 'qty_checker' , " +
		"wrt.name 'wrt' , " +
		"w.name 'warehouse' , " +
		"poi.unfullfill_note " +
		"FROM sales_order so " +
		"JOIN picking_order_assign poa ON so.id = poa.sales_order_id " +
		"JOIN picking_list pl ON pl.id = poa.picking_list_id " +
		"JOIN picking_order po ON po.id = poa.picking_order_id " +
		"JOIN picking_order_item poi ON poi.picking_order_assign_id = poa.id " +
		"JOIN product p ON p.id = poi.product_id " +
		"JOIN uom u ON u.id = p.uom_id " +
		"JOIN wrt ON wrt.id = so.wrt_id " +
		"JOIN branch b ON b.id = so.branch_id " +
		"JOIN merchant m ON m.id = b.merchant_id " +
		"JOIN warehouse w ON w.id = so.warehouse_id " +
		"WHERE  " + where +
		" ORDER BY so.delivery_date, so.code"

	_, e = o.Raw(q1, values).QueryRows(&poi)

	return
}

// getGoodsTransferItem : query to get data of report goods transfer item
func getGoodsTransferItem(cond map[string]interface{}) (gti []*reportGoodsTransferItem, e error) {
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

	q1 := "SELECT gt.recognition_date 'Timestamp', gt.code 'Goods Transfer Code', p.code 'Product_Code', " +
		"p.name 'Product_Name' , u.name 'UOM', w.name 'Warehouse Origin', w2.name 'Warehouse Destination', " +
		"gti.request_qty 'Request Qty' , gti.deliver_qty 'Transfer Qty', gti.receive_qty 'Received Qty', " +
		"g.value_name 'Status', gt.note 'Doc_Note', gti.note 'Note' " +
		"FROM goods_transfer gt " +
		"JOIN goods_transfer_item gti ON gti.goods_transfer_id = gt.id " +
		"JOIN product p ON p.id = gti.product_id " +
		"JOIN uom u ON u.id = p.uom_id " +
		"JOIN warehouse w ON w.id = gt.origin_id " +
		"JOIN warehouse w2 ON w2.id = gt.destination_id " +
		"JOIN glossary g ON g.value_int = gt.status AND g.`table` = 'goods_transfer' " +
		"AND g.`attribute` = 'status' " +
		"WHERE " + where

	_, e = o.Raw(q1, values).QueryRows(&gti)

	return
}

// getRoutingReport : query to get data of report picking routing
func getPickingRoutingReport(cond map[string]interface{}) (rpr []*reportPickingRouting, e error) {
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

	q1 := "SELECT DISTINCT (so.code) as sales_order_code,pl.code as picking_list_code, poa.staff_id as lead_picker_id, prs.staff_id as picker_id, " +
		"p.name as product_name, u.name as UOM, poi.order_qty, " +
		"b.name as rack_name, prs.step_type, prs.`sequence`, " +
		"prs.expected_walking_duration, " +
		"prs.expected_service_duration, " +
		"prs.walking_start_time, prs.walking_finish_time, " +
		"prs.picking_start_time, prs.picking_finish_time, " +
		"prs.status_step as status " +
		"from picking_routing_step prs " +
		"join bin b on b.id = prs.bin_id " +
		"join picking_order_item poi on poi.id  = prs.picking_order_item_id " +
		"join picking_order_assign poa on poa.id  = poi.picking_order_assign_id " +
		"join picking_list pl on pl.id = prs.picking_list_id " +
		"join product p on p.id = poi.product_id " +
		"join uom u on u.id = p.uom_id " +
		"join picking_order po on po.id = poa.picking_order_id " +
		"join sales_order_item soi on soi.sales_order_id = poa.sales_order_id " +
		"join sales_order so on so.id = soi.sales_order_id " +
		"WHERE " + where

	_, e = o.Raw(q1, values).QueryRows(&rpr)

	mapStaff := map[int64]string{}
	mapStepType := map[int64]string{}
	mapStatus := map[int64]string{}

	filter := map[string]interface{}{"table": "picking_routing_step", "attribute": "step_type"}
	exclude := map[string]interface{}{}
	data, _, _ := repository.GetGlossariesByFilter(filter, exclude)
	for _, v := range data {
		mapStepType[int64(v.ValueInt)] = v.ValueName
	}

	filter = map[string]interface{}{"table": "picking_routing_step", "attribute": "status_step"}
	data, _, _ = repository.GetGlossariesByFilter(filter, exclude)
	for _, v := range data {
		mapStatus[int64(v.ValueInt)] = v.ValueName
	}

	for _, v := range rpr {
		v.StepTypeStr = mapStepType[v.StepType]

		if _, ok := mapStaff[v.LeadPickerID]; !ok {
			staff := &model.Staff{ID: v.LeadPickerID}
			staff.Read("id")

			mapStaff[v.LeadPickerID] = staff.Name
		}

		if _, ok := mapStaff[v.PickerID]; !ok {
			staff := &model.Staff{ID: v.PickerID}
			staff.Read("id")

			mapStaff[v.PickerID] = staff.Name
		}

		v.LeadPicker = mapStaff[v.LeadPickerID]
		v.Picker = mapStaff[v.PickerID]

		if !v.WalkingFinishTime.IsZero() {
			walking := v.WalkingFinishTime.Sub(v.WalkingStartTime)
			walkingfloat := walking.Seconds()
			v.ActualWalkingDuration = int64(walkingfloat)
		}

		if !v.PickingFinishTime.IsZero() {
			picking := v.PickingFinishTime.Sub(v.PickingStartTime)
			pickingFloat := picking.Seconds()
			v.ActualPickingDuration = int64(pickingFloat)
		}

		v.StatusStr = mapStatus[v.Status]
	}

	return
}

func getTransferSkuItem(cond map[string]interface{}) (gri []*reportTransferSkuItem, e error) {
	var where string
	var values []interface{}

	o := orm.NewOrm()
	o.Using("read_only")

	for k, v := range cond {
		if reflect.TypeOf(v).Kind().String() == "slice" {
			where = where + " " + k + " >= ? and " + k + " <= ? and"
		} else {
			where = where + " " + k + "? and"
		}

		values = append(values, v)
	}

	where = strings.TrimSuffix(where, " and")
	q := "SELECT ts.code 'TS_Code', g2.value_name 'TS_Status' , ts.recognition_date 'Recognition_Date' , " +
		"IF(po.id IS NULL, gt.code, po.code) 'Inbound_Code' , " +
		"s.code 'Supplier_Code', s.name 'Supplier_Name' , gr.code 'GR_Code', g.value_name 'GR_Status' , w.name 'Warehouse_Origin' , " +
		"IF(w3.id IS NULL, w2.name, w3.name) 'Warehouse_Destination' , p.code 'Product_Code', p.name 'Product_Name', u.name 'UOM' , " +
		"IF(tsi.product_id != tsi.transfer_product_id, NULL,st.available_stock) 'Goods_Stock' , " +
		"IF(tsi.product_id != tsi.transfer_product_id, 0, gri.receive_qty) 'Received_Qty' ," +
		"IF(tsi.product_id != tsi.transfer_product_id, 0,tsi.transfer_qty) 'After_Sortir_Qty_Good' , tsi.waste_qty 'After_Sortir_Qty_Waste' , tsi.discrepancy 'Discrepancy' ," +
		"IF(tsi.product_id = tsi.transfer_product_id, 0,tsi.transfer_qty) 'After_Sortir_Qty_Down_Grade' , IF(tsi.product_id = tsi.transfer_product_id , NULL, p2.code)'Product_Code_Downgrade' , " +
		"IF(tsi.product_id = tsi.transfer_product_id , NULL, p2.name) 'Product_Name_Downgrade' , u2.name 'UOM_Downgrade' " +
		"FROM transfer_sku ts " +
		"JOIN transfer_sku_item tsi ON tsi.transfer_sku_id = ts.id " +
		"JOIN glossary g2 ON g2.value_int = ts.status AND g2.`table` = 'transfer_sku' AND g2.attribute = 'status' " +
		"LEFT JOIN goods_transfer gt ON gt.id = ts.goods_transfer_id " +
		"LEFT JOIN purchase_order po ON po.id = ts.purchase_order_id " +
		"LEFT JOIN supplier s ON s.id = po.supplier_id " +
		"LEFT JOIN goods_receipt gr ON gr.id = ts.goods_receipt_id " +
		"LEFT JOIN glossary g ON g.value_int = gr.status AND g.`table` = 'goods_receipt' AND g.attribute = 'status' " +
		"LEFT JOIN goods_receipt_item gri ON gri.goods_receipt_id = gr.id AND gri.product_id = tsi.product_id " +
		"LEFT JOIN warehouse w ON w.id = gt.origin_id LEFT JOIN warehouse w2 ON w2.id = gt.destination_id " +
		"LEFT JOIN warehouse w3 ON w3.id = po.warehouse_id JOIN product p ON p.id = tsi.product_id " +
		"JOIN product p2 ON p2.id = tsi.transfer_product_id JOIN uom u ON u.id = p.uom_id " +
		"JOIN uom u2 ON u2.id = p2.uom_id JOIN stock st ON st.product_id = tsi.product_id AND st.warehouse_id = ts.warehouse_id " +
		"WHERE " + where

	_, e = o.Raw(q, values).QueryRows(&gri)

	return
}

func getPackingRecommendation(cond map[string]interface{}) (p []*model.PackingOrder, e error) {
	var where string
	var values []interface{}

	o := orm.NewOrm()
	o.Using("read_only")

	for k, v := range cond {
		if reflect.TypeOf(v).Kind().String() == "slice" {
			where = where + " " + k + " >= ? and " + k + " <= ? and"
		} else {
			where = where + " " + k + "? and"
		}

		values = append(values, v)
	}

	where = strings.TrimSuffix(where, " and")
	q := "select po.id from packing_order po where " + where

	_, e = o.Raw(q, values).QueryRows(&p)

	return
}
