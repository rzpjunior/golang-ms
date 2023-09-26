package report

import (
	"errors"
	"reflect"
	"strings"
	"time"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// getPaymentGateway : query to get data of report payment gateway
func getPaymentGateway(cond map[string]interface{}) (m []*reportPaymentGateway, e error) {
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

	// get data requested
	q := "SELECT m.code merchant_code, m.name merchant_name, b.code branch_code, b.name branch_name, so.code sales_order_code,(CASE " +
		"WHEN g.value_name = 'invoice_va' THEN 'Invoice VA' " +
		"WHEN g.value_name = 'fixed_va' THEN 'Fixed VA' " +
		"ELSE '' END) type, " +
		"pc.name channel, tx.amount total_amount, date_format(tx.transaction_date, '%d-%m-%Y') transaction_date, tx.transaction_time transaction_time, " +
		"tx.account_number " +
		"FROM txn_xendit tx " +
		"LEFT JOIN sales_order so ON so.id = tx.sales_order_id " +
		"LEFT JOIN branch b ON b.id = so.branch_id " +
		"JOIN merchant m ON m.id = tx.merchant_id " +
		"JOIN glossary g ON g.table = 'txn_xendit' AND g.attribute = 'type' AND g.value_int = tx.type " +
		"JOIN payment_channel pc ON pc.id = tx.payment_channel_id " +
		"WHERE tx.id > 0 and " + "(" + where + ") ;"

	_, e = o.Raw(q, values).QueryRows(&m)

	// return error some thing went wrong
	return m, nil

}

// getMainOutlet : query to get data of report main outlet
func getMainOutlet(cond map[string]interface{}) (m []*reportMainOutlet, e error) {
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

	// get data requested
	q := "SELECT  m.code code, m.name name, m.billing_address billing_address, a.name finance_area, m.pic_name pic_name, " +
		"m.phone_number phone_number, m.email email, m.total_point current_edenpoint, tps.name default_payment_term, tis.name default_invoice_term, pgs.name payment_group, " +
		"bt.name business_type, 'outlet' customer_group, " +
		"GROUP_CONCAT(tc.name ORDER BY tc.name SEPARATOR ', ') customer_tag, m.status status," +
		"DATE_FORMAT (m.created_at,'%d-%m-%Y %H:%i:%s') created_at, s.name created_by, DATE_FORMAT (m.last_updated_at ,'%d-%m-%Y %H:%i:%s') last_updated_at, s2.name last_updated_by, " +
		"m.business_type_credit_limit business_type_credit_limit, " +
		"m.credit_limit_amount credit_limit_amount, " +
		"m.suspended suspended," +
		"m.credit_limit_remaining " +
		"FROM merchant m " +
		"JOIN area a ON a.id = m.finance_area_id " +
		"JOIN term_payment_sls tps ON tps.id = m.term_payment_sls_id " +
		"JOIN term_invoice_sls tis ON tis.id = m.term_invoice_sls_id " +
		"JOIN payment_group_sls pgs ON pgs.id = m.payment_group_sls_id " +
		"JOIN business_type bt ON bt.id = m.business_type_id " +
		"LEFT JOIN tag_customer tc ON FIND_IN_SET(tc.id, m.tag_customer) != 0 " +
		"JOIN staff s ON s.id = m.created_by " +
		"LEFT JOIN staff s2 ON s2.id = m.last_updated_by " +
		"WHERE m.customer_group = 1 and " + "(" + where + ") " +
		"GROUP BY m.id"

	_, e = o.Raw(q, values).QueryRows(&m)

	// return error some thing went wrong
	return m, nil

}

// getOutlet : query to get data of report outlet
func getOutlet(cond map[string]interface{}) (m []*reportOutlet, e error) {
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

	// get data requested
	q := "SELECT m.code main_outlet_code, b.code outlet_code, b.name outlet_name, b.pic_name pic_name, b.phone_number phone_number, " +
		"ac.name archetype, w.name warehouse_default, ar.name area, ad.province_name province, ad.city_name city, ad.district_name district, " +
		"ad.sub_district_name sub_district, b.shipping_address shipping_address, ps.name price_set, s1.name salesperson, " +
		"b.status status, DATE_FORMAT (b.created_at,'%d-%m-%Y %H:%i:%s') created_at, s2.name created_by, " +
		"DATE_FORMAT (b.last_updated_at,'%d-%m-%Y %H:%i:%s') last_updated_at, s3.name last_updated_by " +
		"FROM branch b " +
		"JOIN merchant m ON m.id = b.merchant_id " +
		"JOIN archetype ac ON ac.id = b.archetype_id " +
		"JOIN warehouse w ON w.id = b.warehouse_id " +
		"JOIN area ar ON ar.id = b.area_id " +
		"JOIN adm_division ad ON ad.sub_district_id = b.sub_district_id " +
		"JOIN price_set ps ON ps.id = b.price_set_id " +
		"JOIN staff s1 ON s1.id = b.salesperson_id " +
		"JOIN staff s2 ON s2.id = b.created_by " +
		"LEFT JOIN staff s3 ON s3.id = b.last_updated_by " +
		"WHERE m.customer_group = 1 and " + "(" + where + ") ;"

	_, e = o.Raw(q, values).QueryRows(&m)

	// return error some thing went wrong
	return m, nil

}

// getAgent : query to get data of report agent
func getAgent(cond map[string]interface{}) (m []*reportAgent, e error) {
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

	// get data requested
	q := "SELECT m.code agent_code, b.name agent_name,m.name main_outlet_name, b.pic_name pic_name, b.phone_number phone_number, m.email email,  m.total_point current_edenpoint," +
		"tps.name AS default_payment_term, tis.name AS default_invoice_term, pgs.name AS payment_group, bt.name AS business_type, ac.name AS archetype, " +
		"GROUP_CONCAT(tc.name ORDER BY tc.name SEPARATOR ', ') AS customer_tag, w.name AS warehouse_default, ar1.name agent_area, ar2.name address_area, " +
		"ad.province_name AS province, ad.city_name AS city, ad.district_name AS district, ad.sub_district_name AS sub_district, " +
		"b.shipping_address AS shipping_address, ps.name AS price_set, s1.name AS salesperson, " +
		"m.status agent_status, b.status address_status, DATE_FORMAT (b.created_at,'%d-%m-%Y %H:%i:%s') AS created_at, s2.name AS created_by, " +
		"DATE_FORMAT (b.last_updated_at,'%d-%m-%Y %H:%i:%s') AS last_updated_at, s3.name AS last_updated_by, " +
		"m.business_type_credit_limit business_type_credit_limit, " +
		"m.credit_limit_amount credit_limit_amount, " +
		"m.suspended suspended, " +
		"m.credit_limit_remaining " +
		"FROM branch b " +
		"JOIN merchant m ON m.id = b.merchant_id " +
		"JOIN term_payment_sls tps ON tps.id = m.term_payment_sls_id " +
		"JOIN term_invoice_sls tis ON tis.id = m.term_invoice_sls_id " +
		"JOIN payment_group_sls pgs ON pgs.id = m.payment_group_sls_id " +
		"JOIN business_type bt ON bt.id = m.business_type_id " +
		"LEFT JOIN tag_customer tc ON FIND_IN_SET(tc.id, m.tag_customer) != 0 " +
		"JOIN archetype ac ON ac.id = b.archetype_id " +
		"JOIN warehouse w ON w.id = b.warehouse_id " +
		"JOIN area ar1 ON ar1.id = m.finance_area_id " +
		"JOIN area ar2 ON ar2.id = b.area_id " +
		"JOIN adm_division ad ON ad.sub_district_id = b.sub_district_id " +
		"JOIN price_set ps ON ps.id = b.price_set_id " +
		"JOIN staff s1 ON s1.id = b.salesperson_id " +
		"JOIN staff s2 ON s2.id = b.created_by " +
		"LEFT JOIN staff s3 ON s3.id = b.last_updated_by " +
		"WHERE m.customer_group = 2 and " + "(" + where + ") " +
		"GROUP BY b.id, m.id;"

	_, e = o.Raw(q, values).QueryRows(&m)

	// return error some thing went wrong
	return m, nil

}

// getVoucherLog : query to get data of report voucher log
func getVoucherLog(cond map[string]interface{}) (m []*reportVoucherLog, e error) {
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

	// get data requested
	q := "SELECT v.code AS voucher_code, v.redeem_code AS redeem_code, DATE_FORMAT (vl.timestamp,'%d-%m-%Y %H:%i:%s') AS redeem_date, g.value_name AS voucher_type, so.code AS order_code, " +
		"so.recognition_date AS order_date, g2.value_name AS order_status, vl.vou_disc_amount AS discount_amount, so.total_charge AS total_order, " +
		"m.code AS merchant_code, m.name AS merchant_name, b.code AS outlet_code, b.name AS outlet_name, " +
		"a.name AS area, bt.name AS business_type, at.name AS archetype " +
		"FROM voucher_log vl " +
		"JOIN voucher v ON v.id = vl.voucher_id " +
		"JOIN glossary g ON g.table = 'voucher' AND g.attribute = 'type' AND g.value_int = v.type " +
		"JOIN sales_order so ON so.id = vl.sales_order_id " +
		"JOIN glossary g2 ON g2.table = 'all' AND g2.attribute = 'doc_status' AND g2.value_int = so.status " +
		"JOIN branch b ON b.id = so.branch_id " +
		"JOIN merchant m ON m.id = b.merchant_id " +
		"JOIN area a ON a.id = so.area_id " +
		"JOIN archetype at ON at.id = so.archetype_id " +
		"JOIN business_type bt ON bt.id = at.business_type_id " +
		"WHERE vl.status = 1 and " + "(" + where + ") " +
		"ORDER BY vl.timestamp DESC;"

	_, e = o.Raw(q, values).QueryRows(&m)

	// return error some thing went wrong
	return m, nil

}

// getSubmission : query to get data of submission
func getSubmission(condSA, condCA, condOOR map[string]interface{}) (m []*reportSubmission, e error) {
	o := orm.NewOrm()
	o.Using("read_only")
	var m2, m3, m4, m5 []*reportSubmission

	var whereSA string
	var valuesSA []interface{}
	for k, v := range condSA {
		if reflect.TypeOf(v).Kind().String() == "slice" {
			whereSA = whereSA + " " + k + "? and ? and"
		} else {
			whereSA = whereSA + " " + k + "? and"
		}
		valuesSA = append(valuesSA, v)
	}
	whereSA = strings.TrimSuffix(whereSA, " and")

	// getting sai for existing customer
	// get data requested from sales_assignment_item
	q1 := "SELECT 'Existing Customer' AS customer_type_str, sai.out_of_route, sai.id, b.id outlet_id, sai.salesperson_id, sa.sales_group_id salesgroup_id, g.value_name task, sai.start_date start_date, sai.end_date end_date, sai.submit_date submission_date, sai.finish_date finish_date, sai.status, sai.objective_codes, s.name salesperson, " +
		"b.code outlet_code, b.name outlet_name, b.shipping_address shipping_address, b.phone_number phone_number, b.latitude latitude, " +
		"b.longitude longitude, sai.latitude task_latitude, sai.longitude task_longitude, sai.actual_distance actual_distance, sai.task_photo task_photo, g2.value_name result, '-' food_app, sg.name sales_group " +
		"FROM sales_assignment_item sai " +
		"JOIN glossary g ON g.value_int = sai.task AND g.attribute = 'task' AND g.table = 'sales_assignment_item' " +
		"JOIN staff s ON s.id = sai.salesperson_id " +
		"JOIN branch b ON b.id = sai.branch_id " +
		"LEFT JOIN glossary g2 ON g2.value_int = sai.answer_option_id AND g2.attribute = 'task_answer' AND g2.table = 'sales_assignment_item' " +
		"JOIN sales_assignment sa ON sa.id = sai.sales_assignment_id " +
		"JOIN sales_group sg ON sg.id = sa.sales_group_id " +
		"WHERE " + "(" + whereSA + ") AND sai.customer_type = 1 AND sai.out_of_route != 1"

	if _, e = o.Raw(q1, valuesSA).QueryRows(&m); e != nil {
		return nil, e
	}

	for _, v := range m {
		if v.Status == 1 {
			v.StatusStr = "active"
		} else if v.Status == 2 {
			v.StatusStr = "finished"

			if v.FinishDate != "" {
				var so []*model.SalesOrder
				finishDate, e := time.Parse("2006-01-02 15:04:05", v.FinishDate)
				if e != nil {
					return nil, e
				}

				if _, e = o.Raw("SELECT * FROM sales_order so WHERE so.branch_id = ? AND so.salesperson_id = ? AND so.sales_group_id = ? AND so.recognition_date = ? AND so.status NOT IN (3,4)",
					v.OutletId, v.SalespersonId, v.SalesGroupId, finishDate.Format("2006-01-02")).QueryRows(&so); e != nil && !errors.Is(e, orm.ErrNoRows) {
					return nil, e
				}
				if so != nil {
					v.EffectiveCall = true
					for _, rec := range so {
						v.RevenueEffectiveCall += rec.TotalCharge
					}

				}
			}

		} else if v.Status == 14 {
			v.StatusStr = "failed"
			var failedReason string
			if e = o.Raw("SELECT g.value_name FROM sales_failed_visit sfv LEFT JOIN glossary g ON g.value_int = sfv.failed_status WHERE g.`table`= 'sales_failed_visit' AND g.`attribute` = 'failed_status' AND sfv.sales_assignment_item_id = ?", v.Id).QueryRow(&failedReason); e != nil {
				return nil, e
			}
			v.Result = failedReason
		} else {
			v.StatusStr = "cancelled"
		}

		v.OutOfRouteStr = "No"

		if v.ObjectiveCodes != "" {
			var objectives []string
			objectiveArr := strings.Split(v.ObjectiveCodes, ",")
			qMark := ""
			for range objectiveArr {
				qMark = qMark + "?,"
			}
			qMark = strings.TrimSuffix(qMark, ",")
			query := "SELECT name FROM sales_assignment_objective WHERE code IN(" + qMark + ")"
			if _, e = o.Raw(query, objectiveArr).QueryRows(&objectives); e != nil {
				return nil, e
			}
			v.ObjectiveCodesStr = strings.Join(objectives, ",")
		}
	}

	// getting sai for customer acquisition
	// get data requested from sales_assignment_item
	q2 := "SELECT 'Customer Acquisition' AS customer_type_str, sai.id, ca.id outlet_id, sai.salesperson_id, sa.sales_group_id salesgroup_id, g.value_name task, sai.start_date start_date, sai.end_date end_date, sai.submit_date submission_date, sai.finish_date finish_date, sai.status, s.name salesperson, " +
		"ca.name outlet_name, ca.address_name shipping_address, ca.phone_number phone_number, ca.latitude latitude, " +
		"ca.longitude longitude, sai.latitude task_latitude, sai.longitude task_longitude, sai.actual_distance actual_distance, sai.task_photo task_photo, g2.value_name result, '-' food_app, sg.name sales_group " +
		"FROM sales_assignment_item sai " +
		"JOIN glossary g ON g.value_int = sai.task AND g.attribute = 'task' AND g.table = 'sales_assignment_item' " +
		"JOIN staff s ON s.id = sai.salesperson_id " +
		"JOIN customer_acquisition ca ON ca.id = sai.customer_acquisition_id " +
		"LEFT JOIN glossary g2 ON g2.value_int = sai.answer_option_id AND g2.attribute = 'task_answer' AND g2.table = 'sales_assignment_item' " +
		"JOIN sales_assignment sa ON sa.id = sai.sales_assignment_id " +
		"JOIN sales_group sg ON sg.id = sa.sales_group_id " +
		"WHERE " + "(" + whereSA + ") AND sai.customer_type = 2 AND sai.out_of_route != 1"

	if _, e = o.Raw(q2, valuesSA).QueryRows(&m2); e != nil {
		return nil, e
	}

	for _, v := range m2 {
		if v.Status == 1 {
			v.StatusStr = "active"
		} else if v.Status == 2 {
			v.StatusStr = "finished"
		} else if v.Status == 14 {
			v.StatusStr = "failed"
			var failedReason string
			if e = o.Raw("SELECT g.value_name FROM sales_failed_visit sfv LEFT JOIN glossary g ON g.value_int = sfv.failed_status WHERE g.`table`= 'sales_failed_visit' AND g.`attribute` = 'failed_status' AND sfv.sales_assignment_item_id = ?", v.Id).QueryRow(&failedReason); e != nil {
				return nil, e
			}
			v.Result = failedReason
		} else {
			v.StatusStr = "cancelled"
		}

		v.OutOfRouteStr = "No"

		if v.ObjectiveCodes != "" {
			var objectives []string
			objectiveArr := strings.Split(v.ObjectiveCodes, ",")
			qMark := ""
			for range objectiveArr {
				qMark = qMark + "?,"
			}
			qMark = strings.TrimSuffix(qMark, ",")
			query := "SELECT name FROM sales_assignment_objective WHERE code IN(" + qMark + ")"
			if _, e = o.Raw(query, objectiveArr).QueryRows(&objectives); e != nil {
				return nil, e
			}
			v.ObjectiveCodesStr = strings.Join(objectives, ",")
		}
	}
	m = append(m, m2...)

	// getting sai for existing customer Out of Route
	// get data requested from sales_assignment_item
	var whereOOR string
	var valuesOOR []interface{}
	for k, v := range condOOR {
		if reflect.TypeOf(v).Kind().String() == "slice" {
			whereOOR = whereOOR + " " + k + "? and ? and"
		} else {
			whereOOR = whereOOR + " " + k + "? and"
		}
		valuesOOR = append(valuesOOR, v)
	}
	whereOOR = strings.TrimSuffix(whereOOR, " and")

	q3 := "SELECT 'Existing Customer' AS customer_type_str, sai.out_of_route, sai.id, b.id outlet_id, sai.salesperson_id, s.sales_group_id salesgroup_id, g.value_name task, sai.start_date start_date, sai.end_date end_date, sai.submit_date submission_date, sai.finish_date finish_date, sai.status, sai.objective_codes, s.name salesperson, " +
		"b.code outlet_code, b.name outlet_name, b.shipping_address shipping_address, b.phone_number phone_number, b.latitude latitude, " +
		"b.longitude longitude, sai.latitude task_latitude, sai.longitude task_longitude, sai.actual_distance actual_distance, sai.task_photo task_photo, g2.value_name result, '-' food_app, sg.name sales_group " +
		"FROM sales_assignment_item sai " +
		"JOIN glossary g ON g.value_int = sai.task AND g.attribute = 'task' AND g.table = 'sales_assignment_item' " +
		"JOIN staff s ON s.id = sai.salesperson_id " +
		"JOIN sales_group sg ON sg.id = s.sales_group_id " +
		"JOIN branch b ON b.id = sai.branch_id " +
		"LEFT JOIN glossary g2 ON g2.value_int = sai.answer_option_id AND g2.attribute = 'task_answer' AND g2.table = 'sales_assignment_item' " +
		"WHERE " + "(" + whereOOR + ") AND sai.customer_type = 1 AND sai.out_of_route = 1"

	if _, e = o.Raw(q3, valuesOOR...).QueryRows(&m3); e != nil {
		return nil, e
	}

	for _, v := range m3 {
		if v.Status == 1 {
			v.StatusStr = "active"
		} else if v.Status == 2 {
			v.StatusStr = "finished"

			if v.FinishDate != "" {
				var so []*model.SalesOrder
				finishDate, e := time.Parse("2006-01-02 15:04:05", v.FinishDate)
				if e != nil {
					return nil, e
				}

				if _, e = o.Raw("SELECT * FROM sales_order so WHERE so.branch_id = ? AND so.salesperson_id = ? AND so.sales_group_id = ? AND so.recognition_date = ? AND so.status NOT IN (3,4)",
					v.OutletId, v.SalespersonId, v.SalesGroupId, finishDate.Format("2006-01-02")).QueryRows(&so); e != nil && !errors.Is(e, orm.ErrNoRows) {
					return nil, e
				}
				if so != nil {
					v.EffectiveCall = true
					for _, rec := range so {
						v.RevenueEffectiveCall += rec.TotalCharge
					}

				}
			}

		} else if v.Status == 14 {
			v.StatusStr = "failed"
			var failedReason string
			if e = o.Raw("SELECT g.value_name FROM sales_failed_visit sfv LEFT JOIN glossary g ON g.value_int = sfv.failed_status WHERE g.`table`= 'sales_failed_visit' AND g.`attribute` = 'failed_status' AND sfv.sales_assignment_item_id = ?", v.Id).QueryRow(&failedReason); e != nil {
				return nil, e
			}
			v.Result = failedReason
		} else {
			v.StatusStr = "cancelled"
		}

		v.OutOfRouteStr = "Yes"
	}
	m = append(m, m3...)

	// getting sai for customer acquisition Out of Route
	// get data requested from sales_assignment_item
	q4 := "SELECT 'Customer Acquisition' AS customer_type_str, sai.out_of_route, sai.id, ca.id outlet_id, sai.salesperson_id, s.sales_group_id salesgroup_id, g.value_name task, sai.start_date start_date, sai.end_date end_date, sai.submit_date submission_date, sai.finish_date finish_date, sai.status, sai.objective_codes, s.name salesperson, " +
		"ca.name outlet_name, ca.address_name shipping_address, ca.phone_number phone_number, ca.latitude latitude, " +
		"ca.longitude longitude, sai.latitude task_latitude, sai.longitude task_longitude, sai.actual_distance actual_distance, sai.task_photo task_photo, g2.value_name result, '-' food_app, sg.name sales_group " +
		"FROM sales_assignment_item sai " +
		"JOIN glossary g ON g.value_int = sai.task AND g.attribute = 'task' AND g.table = 'sales_assignment_item' " +
		"JOIN staff s ON s.id = sai.salesperson_id " +
		"JOIN sales_group sg ON sg.id = s.sales_group_id " +
		"JOIN customer_acquisition ca ON ca.id = sai.customer_acquisition_id " +
		"LEFT JOIN glossary g2 ON g2.value_int = sai.answer_option_id AND g2.attribute = 'task_answer' AND g2.table = 'sales_assignment_item' " +
		"WHERE " + "(" + whereOOR + ") AND sai.customer_type = 2 AND sai.out_of_route = 1"

	if _, e = o.Raw(q4, valuesOOR...).QueryRows(&m4); e != nil {
		return nil, e
	}

	for _, v := range m4 {
		if v.Status == 1 {
			v.StatusStr = "active"
		} else if v.Status == 2 {
			v.StatusStr = "finished"
		} else if v.Status == 14 {
			v.StatusStr = "failed"
			var failedReason string
			if e = o.Raw("SELECT g.value_name FROM sales_failed_visit sfv LEFT JOIN glossary g ON g.value_int = sfv.failed_status WHERE g.`table`= 'sales_failed_visit' AND g.`attribute` = 'failed_status' AND sfv.sales_assignment_item_id = ?", v.Id).QueryRow(&failedReason); e != nil {
				return nil, e
			}
			v.Result = failedReason
		} else {
			v.StatusStr = "cancelled"
		}

		v.OutOfRouteStr = "Yes"
	}
	m = append(m, m4...)

	var whereCA string
	var valuesCA []interface{}
	for k, v := range condCA {
		if reflect.TypeOf(v).Kind().String() == "slice" {
			whereCA = whereCA + " " + k + "? and ? and"
		} else {
			whereCA = whereCA + " " + k + "? and"
		}
		valuesCA = append(valuesCA, v)
	}
	whereCA = strings.TrimSuffix(whereCA, " and")
	// get data requested from customer_acquisition
	q5 := "SELECT 'Customer Acquisition' AS customer_type_str, g.value_name task, '-' start_date, '-' end_date, ca.submit_date submission_date, ca.finish_date finish_date, s.name salesperson, " +
		"'-' outlet_code, ca.name outlet_name, ca.address_name shipping_address, ca.phone_number phone_number, ca.latitude latitude, " +
		"ca.longitude longitude, ca.task_photo task_photo, '-' result, if(ca.food_app = 1, 'yes', 'no') food_app, sg.name sales_group " +
		"FROM customer_acquisition ca " +
		"JOIN glossary g ON g.value_int = ca.task AND g.attribute = 'task' AND g.table = 'sales_assignment_item' " +
		"JOIN staff s on s.id = ca.salesperson_id " +
		"JOIN sales_group sg ON sg.id = ca.sales_group_id " +
		"WHERE " + "(" + whereCA + ")"

	if _, e = o.Raw(q5, valuesCA).QueryRows(&m5); e != nil {
		return nil, e
	}

	for _, v := range m5 {
		v.OutOfRouteStr = "No"
	}

	m = append(m, m5...)
	return m, nil
}
