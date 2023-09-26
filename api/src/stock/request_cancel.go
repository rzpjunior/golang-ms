// Copyright 2020 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package stock

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type cancelRequest struct {
	JobsID primitive.ObjectID `json:"jobs_id"`

	Warehouse         *model.Warehouse           `json:"-"`
	DeliveryOrder     *model.DeliveryOrder       `json:"-"`
	DeliveryOrderItem []*model.DeliveryOrderItem `json:"-"`
	Stock             []*model.Stock             `json:"-"`
	StockLog          []*model.StockLog          `json:"-"`
	WasteLog          []*model.WasteLog          `json:"-"`
	Session           *auth.SessionData          `json:"-"`
}

// cancelRequest : function to validate cancel delivery order based from request data
func (r *cancelRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	o1 := orm.NewOrm()
	o1.Using("read_only")
	var stockOpname int64
	var err error

	warehouseID := r.DeliveryOrder.Warehouse.ID

	r.Warehouse = &model.Warehouse{ID: warehouseID}
	err = r.Warehouse.Read("ID")
	if err != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("warehouse"))
	}

	if err = o1.Raw("SELECT count(id) from stock_opname where warehouse_id = ? AND status = 1", warehouseID).QueryRow(&stockOpname); err != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("stock opname"))
	}
	if stockOpname > 0 {
		o.Failure("id.invalid", util.ErrorRelated("active", "stock opname", r.DeliveryOrder.Warehouse.Name))
	}

	return o
}

func (r *cancelRequest) Messages() map[string]string {
	return map[string]string{}
}
