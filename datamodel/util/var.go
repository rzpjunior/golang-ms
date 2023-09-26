package util

import (
	"strings"

	"git.edenfarm.id/cuxs/orm"
)

func ConvertStatusMaster(status int8) (st string) {
	if status == 1 {
		st = "active"
	} else if status == 2 {
		st = "archived"
	} else if status == 3 {
		st = "deleted"
	}
	return st
}

func ConvertStatusPicking(status int8) (st string) {
	if status == 1 {
		st = "New"
	} else if status == 2 {
		st = "Finished"
	} else if status == 3 {
		st = "On Progress"
	} else if status == 4 {
		st = "Need Approval"
	} else if status == 5 {
		st = "Picked"
	} else if status == 6 {
		st = "Checking"
	} else if status == 7 {
		st = "Cancelled"
	} else if status == 8 {
		st = "Rejected"
	}
	return st
}

func ConvertStatusPickingList(status int8) (st string) {
	if status == 1 {
		st = "New"
	} else if status == 2 {
		st = "Finished"
	} else if status == 3 {
		st = "On Progress"
	} else if status == 4 {
		st = "Rejected"
	}
	return st
}

func ConvertPurchaseOrderItemTaxStatus(status int8) (st string) {
	switch status {
	case 1:
		st = "Yes"
	default:
		st = "No"
	}

	return st
}

func ConvertPurchaseInvoiceItemTaxStatus(status int8) (st string) {
	switch status {
	case 1:
		st = "Yes"
	default:
		st = "No"
	}

	return st
}

func ConvertStatusDoc(status int8) (st string) {
	if status == 1 {
		st = "active"
	} else if status == 2 {
		st = "finished"
	} else if status == 3 {
		st = "cancelled"
	} else if status == 4 {
		st = "deleted"
	} else if status == 5 {
		st = "draft"
	} else if status == 6 {
		st = "partial"
	} else if status == 7 {
		st = "on_delivery"
	} else if status == 8 {
		st = "delivered"
	} else if status == 9 {
		st = "invoiced_not_delivered"
	} else if status == 10 {
		st = "invoiced_on_delivery"
	} else if status == 11 {
		st = "invoiced_delivered"
	} else if status == 12 {
		st = "paid_not_delivered"
	} else if status == 13 {
		st = "paid_on_delivery"
	} else if status == 14 {
		st = "new"
	} else if status == 15 {
		st = "registered"
	} else if status == 16 {
		st = "declined"
	}
	return st
}

// GetOrderChannel : function to get order channel name from glossary
func GetOrderChannel(value ...string) (name string, e error) {
	orm := orm.NewOrm()
	orm.Using("read_only")
	var qMark string
	for _, _ = range value {
		qMark = qMark + "?,"
	}
	qMark = strings.TrimSuffix(qMark, ",")
	e = orm.Raw("select group_concat(g.value_name order by g.id) from glossary g where g.attribute = 'order_channel' and value_int in ("+qMark+")", value).QueryRow(&name)

	return
}

func ConvertStatusFarmingProject(status int8) (st string) {
	switch status {
	case 1:
		st = "leads"
	case 2:
		st = "active"
	case 3:
		st = "deactived"
	case 4:
		st = "cancelled"
	}
	return st
}

func ConvertRejectReasonDoc(rejectReason int8) (st string) {
	if rejectReason == 1 {
		st = "Lost"
	} else if rejectReason == 2 {
		st = "Damaged"
	}
	return st
}
