package repository

import (
	"context"
	"encoding/json"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/internal/app/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IPackingOrderItemRepository interface {
	Get(ctx context.Context, packingOrderID int64, itemID string) (packingOrderItems []*model.PackingOrderItem, count int64, err error)
	GetDetail(ctx context.Context, packingOrderID int64, packType float64, itemID string) (packingOrderItems *model.PackingOrderItem, err error)
	GetDetailPack(ctx context.Context, packingOrderID int64, packType float64, itemID string) (packingOrderItems *model.PackingOrderItem, err error)
	GetByPackingOrderID(ctx context.Context, packingOrderID int64) (packingOrderItems []*model.PackingOrderItem, count int64, err error)
	GetListBarcode(ctx context.Context, packingOrderIDs []int64, itemID string) (packingOrderItemBarcodes []*model.PackingOrderItemBarcode, err error)
	GetBarcode(ctx context.Context, packingOrderID int64, packType float64, itemID string) (packingOrderItemBarcode *model.PackingOrderItemBarcode, err error)
	Create(ctx context.Context, packingOrderItem *model.PackingOrderItem) (err error)
	CreateMany(ctx context.Context, packingOrderItems []*model.PackingOrderItem) (err error)
	CreateBarcode(ctx context.Context, packingOrderItemBarcode *model.PackingOrderItemBarcode) (err error)
	Update(ctx context.Context, packingOrderItem *model.PackingOrderItem, packingOrderID int64, packType float64, itemID string) (err error)
	UpdateBarcode(ctx context.Context, packingOrderItemBarcode *model.PackingOrderItemBarcode, packingOrderID int64, packType float64, itemID string, code string) (err error)
	DeleteMany(ctx context.Context, packingOrderID int64, filter interface{}) (int64, error)
}

type PackingOrderItemRepository struct {
	opt opt.Options
}

func NewPackingOrderItemRepository() IPackingOrderItemRepository {
	return &PackingOrderItemRepository{
		opt: global.Setup.Common,
	}
}

func (r *PackingOrderItemRepository) Get(ctx context.Context, packingOrderID int64, itemID string) (packingOrderItems []*model.PackingOrderItem, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PackingOrderItemRepository.Get")
	defer span.End()

	db := r.opt.Mongox

	opts := options.FindOptions{}

	filter := map[string]interface{}{
		"packing_order_id": packingOrderID,
		"status":           statusx.ConvertStatusName(statusx.Active),
	}

	if itemID != "" {
		filter["item_id_gp"] = itemID
	}

	var ret []byte
	ret, err = db.GetByFilter(ctx, "packing_order_item", &opts)
	if err != nil {
		span.RecordError(err)
		return
	}

	err = json.Unmarshal(ret, &packingOrderItems)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *PackingOrderItemRepository) GetDetail(ctx context.Context, packingOrderID int64, packType float64, itemID string) (packingOrderItems *model.PackingOrderItem, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PackingOrderItemRepository.GetDetail")
	defer span.End()

	db := r.opt.Mongox

	filter := map[string]interface{}{
		"packing_order_id": packingOrderID,
		"pack_type":        packType,
		"item_id_gp":       itemID,
	}

	var ret []byte
	ret, err = db.FindByFilter(ctx, "packing_order_item", filter)
	if err != nil {
		span.RecordError(err)
		return
	}

	err = json.Unmarshal(ret, &packingOrderItems)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *PackingOrderItemRepository) GetDetailPack(ctx context.Context, packingOrderID int64, packType float64, itemID string) (packingOrderItems *model.PackingOrderItem, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PackingOrderItemBarcodeRepository.GetDetail")
	defer span.End()

	db := r.opt.Mongox

	filter := map[string]interface{}{
		"packing_order_id": packingOrderID,
		// "pack_type":        packType,
		// "item_id_gp":       itemID,
	}

	if packType != 0 {
		filter["pack_type"] = packType
	}

	if itemID != "" {
		filter["item_id_gp"] = itemID
	}

	var ret []byte
	ret, err = db.FindByFilter(ctx, "packing_order_item", filter)
	if err != nil {
		span.RecordError(err)
		return
	}

	err = json.Unmarshal(ret, &packingOrderItems)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *PackingOrderItemRepository) GetByPackingOrderID(ctx context.Context, packingOrderID int64) (packingOrderItems []*model.PackingOrderItem, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PackingOrderItemRepository.GetByPackingOrderID")
	defer span.End()

	db := r.opt.Mongox

	opts := options.FindOptions{}
	opts.SetSort(map[string]interface{}{"pack_type": 1})
	filter := map[string]interface{}{
		"packing_order_id": packingOrderID,
		"status":           statusx.ConvertStatusName(statusx.Active),
	}

	var ret []byte
	ret, err = db.GetByFilter(ctx, "packing_order_item", filter, &opts)
	if err != nil {
		span.RecordError(err)
		return
	}

	err = json.Unmarshal(ret, &packingOrderItems)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *PackingOrderItemRepository) GetListBarcode(ctx context.Context, packingOrderIDs []int64, itemID string) (packingOrderItemBarcodes []*model.PackingOrderItemBarcode, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PackingOrderItemRepository.GetListBarcode")
	defer span.End()

	db := r.opt.Mongox

	filter := map[string]interface{}{
		"packing_order_id": bson.M{"$in": packingOrderIDs},
	}

	if itemID != "" {
		filter["item_id_gp"] = itemID
	}

	var ret []byte
	ret, err = db.GetByFilter(ctx, "packing_order_item_barcode", filter)
	if err != nil {
		span.RecordError(err)
		return
	}

	err = json.Unmarshal(ret, &packingOrderItemBarcodes)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *PackingOrderItemRepository) GetBarcode(ctx context.Context, packingOrderID int64, packType float64, itemID string) (packingOrderItemBarcode *model.PackingOrderItemBarcode, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PackingOrderItemRepository.GetBarcode")
	defer span.End()

	db := r.opt.Mongox

	opts := options.FindOneOptions{}
	opts.SetSort(map[string]interface{}{"code": -1})

	filter := map[string]interface{}{
		"packing_order_id": packingOrderID,
		"pack_type":        packType,
		"item_id_gp":       itemID,
		"status":           statusx.ConvertStatusName(statusx.Active),
	}

	var ret []byte
	ret, err = db.FindByFilter(ctx, "packing_order_item_barcode", filter, &opts)
	if err != nil {
		span.RecordError(err)
		return
	}

	err = json.Unmarshal(ret, &packingOrderItemBarcode)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *PackingOrderItemRepository) Create(ctx context.Context, packingOrderItem *model.PackingOrderItem) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PackingOrderItemRepository.Create")
	defer span.End()

	db := r.opt.Mongox
	_, err = db.Insert(ctx, "packing_order_item", packingOrderItem)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *PackingOrderItemRepository) CreateMany(ctx context.Context, packingOrderItems []*model.PackingOrderItem) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PackingOrderItemRepository.Create")
	defer span.End()

	db := r.opt.Mongox
	ipackingOrderItems := make([]interface{}, 0)
	for _, item := range packingOrderItems {
		ipackingOrderItems = append(ipackingOrderItems, item)
	}

	_, err = db.InsertBulk(ctx, "packing_order_item", ipackingOrderItems)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *PackingOrderItemRepository) CreateBarcode(ctx context.Context, packingOrderItemBarcode *model.PackingOrderItemBarcode) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PackingOrderItemRepository.CreateBarcode")
	defer span.End()

	db := r.opt.Mongox
	_, err = db.Insert(ctx, "packing_order_item_barcode", packingOrderItemBarcode)
	if err != nil {
		span.RecordError(err)
		return
	}
	return
}

func (r *PackingOrderItemRepository) Update(ctx context.Context, packingOrderItem *model.PackingOrderItem, packingOrderID int64, packType float64, itemID string) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PackingOrderItemRepository.Update")
	defer span.End()

	filter := map[string]interface{}{
		"packing_order_id": packingOrderID,
		"pack_type":        packType,
		"item_id_gp":       itemID,
	}

	db := r.opt.Mongox
	err = db.Update(ctx, "packing_order_item", filter, packingOrderItem)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *PackingOrderItemRepository) UpdateBarcode(ctx context.Context, packingOrderItemBarcode *model.PackingOrderItemBarcode, packingOrderID int64, packType float64, itemID string, code string) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PackingOrderItemRepository.UpdateBarcode")
	defer span.End()

	filter := map[string]interface{}{
		"packing_order_id": packingOrderID,
		"pack_type":        packType,
		"item_id_gp":       itemID,
		"status":           statusx.ConvertStatusName(statusx.Active),
	}

	if code != "" {
		filter["code"] = code
	}

	db := r.opt.Mongox
	err = db.Update(ctx, "packing_order_item_barcode", filter, packingOrderItemBarcode)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *PackingOrderItemRepository) DeleteMany(ctx context.Context, packingOrderID int64, filter interface{}) (deleted int64, err error) {

	ctx, span := r.opt.Trace.Start(ctx, "PackingOrderItemRepository.DeleteMany")
	defer span.End()

	filter = map[string]interface{}{
		"packing_order_id": packingOrderID,
	}

	db := r.opt.Mongox
	_, err = db.DeleteMany(ctx, "packing_order_item", filter)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}
