package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PackingOrderPack struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	PackingOrderID     int64              `bson:"packing_order_id" json:"packing_order_id"`
	ItemID             int64              `bson:"item_id" json:"item_id"`
	ItemIDGP           string             `bson:"item_id_gp" json:"item_id_gp"`
	ItemName           string             `bson:"item_name" json:"item_name"`
	UomID              int64              `bson:"uom_id" json:"uom_id"`
	UomIDGP            string             `bson:"uom_id_gp" json:"uom_id_gp"`
	Uom                string             `bson:"uom" json:"uom"`
	OrderMinQty        float64            `bson:"order_min_qty" json:"order_min_qty"`
	WeightPack         float64            `bson:"weight_pack" json:"weight_pack"`
	ProgressPercentage float64            `bson:"progress_percentage" json:"progress_percentage"`
	ExcessPercentage   float64            `bson:"excess_percentage" json:"excess_percentage"`
	TotalOrderWeight   int64              `bson:"total_order_weight" json:"total_order_weight"`
	ItemPack           []*ItemPack        `bson:"item_pack" json:"item_pack"`
}

type ItemPack struct {
	PackType          float64 `bson:"pack_type" json:"pack_type" `
	ExpectedTotalPack float64 `bson:"expected_total_pack" json:"expected_total_pack" `
	ActualTotalPack   float64 `bson:"actual_total_pack" json:"actual_total_pack" `
	Status            int8    `bson:"status" json:"status"`
}

type ConfigItem struct {
	ItemName string `bson:"item_name" json:"item_name" `
	ItemID   int64  `bson:"item_id" json:"item_id"`
	Uom      string `bson:"uom" json:"uom"`
	UomIDGP  string `bson:"uom_id_gp" json:"uom_id_gp"`
}
