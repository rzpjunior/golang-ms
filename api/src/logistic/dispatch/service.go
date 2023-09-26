package dispatch

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/log"
)

func Update(r updateRequest) (poa *model.PickingOrderAssign, e error) {
	o := orm.NewOrm()
	o.Begin()

	// this condition for function audit log belongs to
	var functionLog string
	if r.PickingOrderAssign.Courier != nil {
		if r.Courier == nil {
			r.PickingOrderAssign.Courier = nil
			functionLog = "delete_courier"
		} else {
			r.PickingOrderAssign.Courier = r.Courier

			functionLog = "add_courier"
		}
	} else {
		r.PickingOrderAssign.Courier = r.Courier
		functionLog = "add_courier"
	}
	r.PickingOrderAssign.Dispatcher = r.Session.Staff

	if _, e = o.Update(r.PickingOrderAssign, "Courier", "Dispatcher"); e == nil {
		e = log.AuditLogByUser(r.Session.Staff, r.PickingOrderAssign.ID, "picking_order_assign", functionLog, "")
	} else {
		o.Rollback()
		return nil, e
	}
	o.Commit()
	return r.PickingOrderAssign, e
}

func UpdateVendor(r updateVendorRequest) (poa *model.PickingOrderAssign, e error) {
	o := orm.NewOrm()
	o.Begin()

	r.PickingOrderAssign.CourierVendor = r.CourierVendor
	r.PickingOrderAssign.Courier = r.Courier
	r.PickingOrderAssign.Dispatcher = r.Session.Staff

	if _, e = o.Update(r.PickingOrderAssign, "CourierVendor", "Courier", "Dispatcher"); e == nil {
		e = log.AuditLogByUser(r.Session.Staff, r.PickingOrderAssign.ID, "picking_order_assign", "update vendor", "")
	} else {
		o.Rollback()
		return nil, e
	}
	o.Commit()
	return r.PickingOrderAssign, e
}

func ScanDispatch(r scanRequest) (poa *model.PickingOrderAssign, e error) {
	o := orm.NewOrm()
	o.Begin()

	switch r.TypeRequest {
	case "scan":
		r.PickingOrderAssign.TotalScanDispatch = r.PickingOrderAssign.TotalScanDispatch + 1
		if r.PickingOrderAssign.TotalKoli != float64(r.PickingOrderAssign.TotalScanDispatch) {
			r.PickingOrderAssign.DispatchStatus = 3
			if _, e = o.Update(r.PickingOrderAssign, "TotalScanDispatch", "DispatchStatus"); e != nil {
				o.Rollback()
				return nil, e
			}
		} else {
			r.PickingOrderAssign.DispatchStatus = 2
			if _, e = o.Update(r.PickingOrderAssign, "TotalScanDispatch", "DispatchStatus"); e != nil {
				o.Rollback()
				return nil, e
			}
		}
		r.DeliveryKoliIncrement.IsRead = 1
		if _, e = o.Update(r.DeliveryKoliIncrement, "IsRead"); e != nil {
			o.Rollback()
			return nil, e
		}

		e = log.AuditLogByUser(r.Session.Staff, r.PickingOrderAssign.ID, "dispatch", "scan", "")
	case "reset":
		r.PickingOrderAssign.TotalScanDispatch = 0
		r.PickingOrderAssign.DispatchStatus = 1
		if _, e = o.Update(r.PickingOrderAssign, "TotalScanDispatch", "DispatchStatus"); e != nil {
			o.Rollback()
			return nil, e
		}

		dki := orm.Params{
			"is_read": int8(0),
		}

		if _, e = o.QueryTable(new(model.DeliveryKoliIncrement)).Filter("sales_order_id", r.PickingOrderAssign.SalesOrder.ID).Update(dki); e != nil {
			o.Rollback()
			return nil, e
		}

		e = log.AuditLogByUser(r.Session.Staff, r.PickingOrderAssign.ID, "dispatch", "reset", "")
	}
	o.Commit()
	return r.PickingOrderAssign, e
}
