package sales_assignment

import (
	"reflect"
	"strings"

	"git.edenfarm.id/cuxs/orm"
)

// getBranchBySalesGroup : query to get data of template assignment task
func getBranchBySalesGroup(cond map[string]interface{}) (m []*templateBranchBySalesGroup, e error) {
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

	// get data from branch
	q := "SELECT 'Existing Customer' AS customer_type, b.name outlet_name, sg.id sales_group_id, sg.name sales_group_name, b.id branch_id, b.code branch_code, sd.name sub_district_name, d.name district_name, s.id salesperson_id," +
		" s.code staff_code, s.name staff_name" +
		" FROM branch b" +
		" JOIN staff s ON s.id = b.salesperson_id" +
		" JOIN sales_group sg ON sg.id = s.sales_group_id" +
		" JOIN sub_district sd ON sd.id = b.sub_district_id" +
		" JOIN district d ON d.id = sd.district_id" +
		" WHERE b.status = 1 AND s.status = 1 AND " + where + ";"

	_, e = o.Raw(q, values).QueryRows(&m)

	// get data from customer acquisition
	var m1 []*templateBranchBySalesGroup
	q = "SELECT 'Customer Acquisition' AS customer_type, ca.name outlet_name, sg.id sales_group_id, sg.name sales_group_name, ca.id branch_id, s.id salesperson_id," +
		" s.code staff_code, s.name staff_name" +
		" FROM customer_acquisition ca" +
		" JOIN staff s ON s.id = ca.salesperson_id" +
		" JOIN sales_group sg ON sg.id = s.sales_group_id" +
		" WHERE ca.status = 2 AND s.status = 1 AND " + where + ";"

	_, e = o.Raw(q, values).QueryRows(&m1)
	m = append(m, m1...)

	// return error some thing went wrong
	return m, nil

}
