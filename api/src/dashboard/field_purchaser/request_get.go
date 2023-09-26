package field_purchaser

import (
	"strings"

	"git.edenfarm.id/cuxs/orm"
)

// getPurchasePlanSummary : query to get total of purchase plan with condition
func GetPurchasePlanSummary(cond map[string]interface{}) (ppsum *getPurchasePlanSummary, total int64, e error) {
	var where string
	var values []interface{}

	o := orm.NewOrm()
	o.Using("read_only")

	for k, v := range cond {
		where = where + " " + k + "? and"
		values = append(values, v)
	}
	where = where + " pp.status = 1 "
	where = strings.TrimSuffix(where, " and")

	q := "select count(pp.id) as total_purchase_plan_active, count(pp.assigned_to) as total_assigned_purchase_plan " +
		"from purchase_plan pp " +
		"JOIN supplier_organization so ON so.id = pp.supplier_organization_id " +
		"where " + where

	e = o.Raw(q, values).QueryRow(&ppsum)

	return
}
