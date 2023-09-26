// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package courier

import (
	"time"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// Save : function to save data requested into database
func Save(r createRequest) (response interface{}, e error) {

	o := orm.NewOrm()
	o.Begin()

	// deleted courier transaction row if already exist
	if r.CourierTransaction != nil {
		o.Delete(r.CourierTransaction)
	}

	if e == nil {
		u := &model.CourierTransaction{
			DeliveryOrder:  r.DeliveryOrder,
			Latitude:       r.Latitude,
			Longitude:      r.Longitude,
			Accuracy:       r.Accuracy,
			CourierName:    r.CourierName,
			CourierPhoneNo: r.CourierPhoneNo,
			CreatedAt:      time.Now(),
			Note:           r.Note,
		}

		if _, e = o.Insert(u); e != nil {
			o.Rollback()
			return nil, e
		} else {
			// update field hasDelivered on DO table on regarded DO
			r.DeliveryOrder.HasDelivered = 1
			if _, e = o.Update(r.DeliveryOrder, "HasDelivered"); e != nil {
				o.Rollback()
				return nil, e
			}
		}

	}

	o.Commit()

	response = "Data has been saved Successfully"
	return response, e
}
