// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PackingRecommendation struct {
	ID                      primitive.ObjectID `orm:"column(_id)" json:"_id" bson:"_id,omitempty" `
	PackingOrderID          int64              `orm:"-" json:"packing_order_id" bson:"packing_order_id,omitempty"`
	ProductID               int64              `orm:"-" json:"product_id" bson:"product_id,omitempty"`
	PackType                float64            `orm:"-" json:"pack_type" bson:"pack_type,omitempty"`
	TotalProgressPercentage float64            `orm:"-" json:"total_progress_pct" bson:"total_progress_pct,omitempty"`

	ProductPack []*PackAdjustment `orm:"-" json:"product_pack"`
	Product     *Product          `orm:"-" json:"product"`
}

type PackAdjustment struct {
	PackType          float64 `orm:"-" json:"pack_type,omitempty" bson:"pack_type,omitempty"`
	ExpectedTotalPack float64 `orm:"-" json:"expected_total_pack" bson:"expected_total_pack"`
	ActualTotalPack   float64 `orm:"-" json:"actual_total_pack" bson:"actual_total_pack"`
}

type ProductPercentage struct {
	ExpectedTotalPack  float64 `orm:"-" json:"expected_total_pack" bson:"expected_total_pack"`
	ActualTotalPack    float64 `orm:"-" json:"actual_total_pack" bson:"actual_total_pack"`
	ProgressPercentage float64 `orm:"-" json:"progress_pct" bson:"progress_pct"`
}

type ResponseData struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	PackingOrderID    int64              `bson:"packing_order_id" json:"packing_order_id"`
	ProductID         int64              `bson:"product_id" json:"product_id"`
	PackType          float64            `bson:"pack_type" json:"pack_type"`
	ExpectedTotalPack float64            `bson:"expected_total_pack" json:"expected_total_pack"`
	WeightScale       float64            `json:"weight_scale"`
	ActualTotalPack   float64            `bson:"actual_total_pack" json:"actual_total_pack"`
	WeightPack        float64            `bson:"weight_pack" json:"weight_pack"`
	Status            int                `bson:"status" json:"status"`
	CodePrint         string             `bson:"code_print" json:"code_print"`
	Product           *Product           `json:"product,omitempty"`
	PackingOrder      *PackingOrder      `json:"packing_order,omitempty"`
}

type BarcodeModel struct {
	ID             int64   `json:"id" bson:"id,omitempty"`
	PackingOrderID int64   `bson:"packing_order_id" json:"packing_order_id"`
	Code           string  `json:"code" bson:"code,omitempty"`
	ProductID      int64   `bson:"product_id" json:"product_id"`
	PackType       float64 `bson:"pack_type" json:"pack_type"`
	WeightScale    float64 `bson:"weight_scale" json:"weight_scale"`
	Status         int     `bson:"status" json:"status"`
	DeltaPrint     int     `bson:"delta_print" json:"delta_print"`
	CreatedAt      string  `bson:"created_at" json:"created_at"`
	CreatedBy      int64   `bson:"created_by" json:"created_by"`
	DeletedAt      string  `bson:"deleted_at" json:"deleted_at"`
	DeletedBy      int64   `bson:"deleted_by" json:"deleted_by"`
	CreatedObj     *Staff  `json:"created_by_obj,omitempty"`
	DeletedObj     *Staff  `json:"deleted_by_obj,omitempty"`

	PackingOrder *PackingOrder `json:"packing_order"`
	Product      *Product      `json:"product"`
}

type PackRecommendation struct {
	PackingOrderID    int64   `bson:"packing_order_id" json:"packing_order_id"`
	SalesOrderID      int64   `bson:"sales_order_id" json:"sales_order_id"`
	ProductID         int64   `bson:"product_id" json:"product_id"`
	PackType          float64 `orm:"-" json:"pack_type" bson:"pack_type"`
	ExpectedTotalPack float64 `orm:"-" bson:"expected_total_pack" json:"expected_total_pack"`
	Status            int     `bson:"status" json:"status"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *PackingRecommendation) MarshalJSON() ([]byte, error) {
	type Alias PackingRecommendation

	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	})
}
