package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/edenlabs/edenlabs/timex"
	"git.edenfarm.id/edenlabs/edenlabs/utils"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/repository"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"github.com/sirupsen/logrus"
)

type IPurchaseOrderService interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (res []dto.PurchaseOrderResponse, total int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (res dto.PurchaseOrderResponse, err error)
	Create(ctx context.Context, req *dto.CreatePurchaseOrderRequest) (res dto.PurchaseOrderResponse, err error)
	CreateGP(ctx context.Context, req *dto.CreatePurchaseOrderGPRequest) (res dto.CommonPurchaseOrderGPResponse, err error)
	Commit(ctx context.Context, id int64) (err error)
	Cancel(ctx context.Context, id int64) (err error)
	Update(ctx context.Context, req *dto.UpdatePurchaseOrderRequest) (res dto.PurchaseOrderResponse, err error)
	UpdateGP(ctx context.Context, req *dto.UpdatePurchaseOrderGPRequest) (res dto.CommonPurchaseOrderGPResponse, err error)
	UpdateProduct(ctx context.Context, req *dto.UpdateProductPurchaseOrderRequest) (res dto.PurchaseOrderResponse, err error)
	GetGP(ctx context.Context, req *pb.GetPurchaseOrderGPListRequest) (res *pb.GetPurchaseOrderGPResponse, err error)
	GetDetailGP(ctx context.Context, req *pb.GetPurchaseOrderGPDetailRequest) (res *pb.GetPurchaseOrderGPResponse, err error)
	CommitPurchaseOrderGP(ctx context.Context, req *pb.CommitPurchaseOrderGPRequest) (res *pb.CreateTransferRequestGPResponse, err error)
	CancelPurchaseOrderGP(ctx context.Context, req *dto.CancelPurchaseOrderGPRequest) (res *dto.CancelPurchaseOrderGPResponse, err error)
	CreateConsolidatedShipmentGP(ctx context.Context, req *dto.CreateConsolidatedShipmentGPRequest) (res dto.CreateConsolidatedShipmentGPResponse, err error)
	UpdateConsolidatedShipmentGP(ctx context.Context, req *dto.UpdateConsolidatedShipmentGPRequest) (res dto.UpdateConsolidatedShipmentGPResponse, err error)
}

type PurchaseOrderService struct {
	opt                         opt.Options
	RepositoryPurchaseOrder     repository.IPurchaseOrderRepository
	RepositoryVendor            repository.IVendorRepository
	RepositorySite              repository.ISiteRepository
	RepositoryUom               repository.IUomRepository
	RepositoryItem              repository.IItemRepository
	RepositoryPurchasePlan      repository.IPurchasePlanRepository
	RepositoryPurchaseOrderItem repository.IPurchaseOrderItemRepository
	RepositoryReceiving         repository.IReceivingRepository
}

func NewPurchaseOrderService() IPurchaseOrderService {
	return &PurchaseOrderService{
		opt:                         global.Setup.Common,
		RepositoryPurchaseOrder:     repository.NewPurchaseOrderRepository(),
		RepositoryVendor:            repository.NewVendorRepository(),
		RepositorySite:              repository.NewSiteRepository(),
		RepositoryUom:               repository.NewUomRepository(),
		RepositoryItem:              repository.NewItemRepository(),
		RepositoryPurchasePlan:      repository.NewPurchasePlanRepository(),
		RepositoryPurchaseOrderItem: repository.NewPurchaseOrderItemRepository(),
		RepositoryReceiving:         repository.NewReceivingRepository(),
	}
}

func (s *PurchaseOrderService) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (res []dto.PurchaseOrderResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PurchaseOrderService.Get")
	defer span.End()

	var purchaseOrders []*model.PurchaseOrder
	purchaseOrders, total, err = s.RepositoryPurchaseOrder.Get(ctx, offset, limit, status, search, orderBy)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, purchaseOrder := range purchaseOrders {
		var (
			receiving []*model.Receiving
			grs       []*dto.ReceivingListinDetailResponse
		)
		receiving, err = s.RepositoryReceiving.GetByInbound(ctx, 1, purchaseOrder.ID)
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		for _, gr := range receiving {
			grs = append(grs, &dto.ReceivingListinDetailResponse{
				ID:     fmt.Sprintf("%d", gr.ID),
				Code:   gr.Code,
				Status: gr.Status,
			})
		}

		res = append(res, dto.PurchaseOrderResponse{
			ID:                     purchaseOrder.ID,
			Code:                   purchaseOrder.Code,
			VendorID:               purchaseOrder.VendorID,
			SiteID:                 purchaseOrder.SiteID,
			TermPaymentPurID:       purchaseOrder.TermPaymentPurID,
			VendorClassificationID: purchaseOrder.VendorClassificationID,
			PurchasePlanID:         purchaseOrder.PurchasePlanID,
			ConsolidatedShipmentID: purchaseOrder.ConsolidatedShipmentID,
			Status:                 purchaseOrder.Status,
			RecognitionDate:        purchaseOrder.RecognitionDate,
			EtaDate:                purchaseOrder.EtaDate,
			SiteAddress:            purchaseOrder.SiteAddress,
			EtaTime:                purchaseOrder.EtaTime,
			TaxPct:                 purchaseOrder.TaxPct,
			DeliveryFee:            purchaseOrder.DeliveryFee,
			TotalPrice:             purchaseOrder.TotalPrice,
			TaxAmount:              purchaseOrder.TaxAmount,
			TotalCharge:            purchaseOrder.TotalCharge,
			TotalInvoice:           purchaseOrder.TotalInvoice,
			TotalWeight:            purchaseOrder.TotalWeight,
			Note:                   purchaseOrder.Note,
			DeltaPrint:             purchaseOrder.DeltaPrint,
			Latitude:               purchaseOrder.Latitude,
			Longitude:              purchaseOrder.Longitude,
			UpdatedAt:              purchaseOrder.UpdatedAt,
			UpdatedBy:              purchaseOrder.UpdatedBy,
			CreatedFrom:            purchaseOrder.CreatedFrom,
			HasFinishedGr:          purchaseOrder.HasFinishedGr,
			CreatedAt:              purchaseOrder.CreatedAt,
			CreatedBy:              purchaseOrder.CreatedBy,
			CommittedAt:            purchaseOrder.CommittedAt,
			CommittedBy:            purchaseOrder.CommittedBy,
			AssignedTo:             purchaseOrder.AssignedTo,
			AssignedBy:             purchaseOrder.AssignedBy,
			AssignedAt:             purchaseOrder.AssignedAt,
			Locked:                 purchaseOrder.Locked,
			LockedBy:               purchaseOrder.LockedBy,
			Receiving:              grs,
		})
	}

	jsonRes, _ := json.Marshal(res)
	fmt.Println(string(jsonRes))

	return
}

func (s *PurchaseOrderService) GetDetail(ctx context.Context, id int64, code string) (res dto.PurchaseOrderResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PurchaseOrderService.GetDetail")
	defer span.End()

	var (
		purchaseOrder      *model.PurchaseOrder
		purchaseOrderItems []*model.PurchaseOrderItem
		receiving          []*model.Receiving
		grs                []*dto.ReceivingListinDetailResponse
		poi                []*dto.PurchaseOrderItemResponse
		itemRes            *model.Item
		uomRes             *model.Uom
	)
	purchaseOrder, err = s.RepositoryPurchaseOrder.GetDetail(ctx, id, code)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	purchaseOrderItems, err = s.RepositoryPurchaseOrderItem.GetByPurchaseOrderId(ctx, purchaseOrder.ID)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, item := range purchaseOrderItems {
		itemRes, err = s.RepositoryItem.GetDetail(ctx, item.ItemID, "")
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		uomRes, err = s.RepositoryUom.GetDetail(ctx, itemRes.UomID, "")
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		poi = append(poi, &dto.PurchaseOrderItemResponse{
			ID:                 item.ID,
			PurchaseOrderID:    item.PurchaseOrderID,
			PurchasePlanItemID: item.PurchasePlanItemID,
			ItemID:             item.ItemID,
			OrderQty:           item.OrderQty,
			UnitPrice:          item.UnitPrice,
			TaxableItem:        item.TaxableItem,
			IncludeTax:         item.IncludeTax,
			TaxPercentage:      item.TaxPercentage,
			TaxAmount:          item.TaxAmount,
			UnitPriceTax:       item.UnitPriceTax,
			Subtotal:           item.Subtotal,
			Weight:             item.Weight,
			Note:               item.Note,
			PurchaseQty:        item.PurchaseQty,
			Item: &dto.ItemResponse{
				ID:                      itemRes.ID,
				Code:                    itemRes.Code,
				UomID:                   itemRes.UomID,
				ClassID:                 itemRes.ClassID,
				ItemCategoryID:          itemRes.ItemCategoryID,
				Description:             itemRes.Description,
				UnitWeightConversion:    itemRes.UnitWeightConversion,
				OrderMinQty:             itemRes.OrderMinQty,
				OrderMaxQty:             itemRes.OrderMaxQty,
				ItemType:                itemRes.ItemType,
				Packability:             itemRes.Packability,
				Capitalize:              itemRes.Capitalize,
				Note:                    itemRes.Note,
				ExcludeArchetype:        itemRes.ExcludeArchetype,
				MaxDayDeliveryDate:      itemRes.MaxDayDeliveryDate,
				FragileGoods:            itemRes.FragileGoods,
				Taxable:                 itemRes.Taxable,
				OrderChannelRestriction: itemRes.OrderChannelRestriction,
				Status:                  itemRes.Status,
				StatusConvert:           statusx.ConvertStatusValue(itemRes.Status),
				CreatedAt:               timex.ToLocTime(ctx, itemRes.CreatedAt),
				UpdatedAt:               timex.ToLocTime(ctx, itemRes.UpdatedAt),
				Uom: &dto.UomResponse{
					ID:             uomRes.ID,
					Code:           uomRes.Code,
					Description:    uomRes.Description,
					Status:         uomRes.Status,
					DecimalEnabled: uomRes.DecimalEnabled,
					CreatedAt:      uomRes.CreatedAt,
					UpdatedAt:      uomRes.UpdatedAt,
				},
			},
		})
	}

	receiving, err = s.RepositoryReceiving.GetByInbound(ctx, 1, purchaseOrder.ID)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, gr := range receiving {
		grs = append(grs, &dto.ReceivingListinDetailResponse{
			ID:     fmt.Sprintf("%d", gr.ID),
			Code:   gr.Code,
			Status: gr.Status,
		})
	}

	res = dto.PurchaseOrderResponse{
		ID:                     purchaseOrder.ID,
		Code:                   purchaseOrder.Code,
		VendorID:               purchaseOrder.VendorID,
		SiteID:                 purchaseOrder.SiteID,
		TermPaymentPurID:       purchaseOrder.TermPaymentPurID,
		VendorClassificationID: purchaseOrder.VendorClassificationID,
		PurchasePlanID:         purchaseOrder.PurchasePlanID,
		ConsolidatedShipmentID: purchaseOrder.ConsolidatedShipmentID,
		Status:                 purchaseOrder.Status,
		RecognitionDate:        purchaseOrder.RecognitionDate,
		EtaDate:                purchaseOrder.EtaDate,
		SiteAddress:            purchaseOrder.SiteAddress,
		EtaTime:                purchaseOrder.EtaTime,
		TaxPct:                 purchaseOrder.TaxPct,
		DeliveryFee:            purchaseOrder.DeliveryFee,
		TotalPrice:             purchaseOrder.TotalPrice,
		TaxAmount:              purchaseOrder.TaxAmount,
		TotalCharge:            purchaseOrder.TotalCharge,
		TotalInvoice:           purchaseOrder.TotalInvoice,
		TotalWeight:            purchaseOrder.TotalWeight,
		Note:                   purchaseOrder.Note,
		DeltaPrint:             purchaseOrder.DeltaPrint,
		Latitude:               purchaseOrder.Latitude,
		Longitude:              purchaseOrder.Longitude,
		UpdatedAt:              purchaseOrder.UpdatedAt,
		UpdatedBy:              purchaseOrder.UpdatedBy,
		CreatedFrom:            purchaseOrder.CreatedFrom,
		HasFinishedGr:          purchaseOrder.HasFinishedGr,
		CreatedAt:              purchaseOrder.CreatedAt,
		CreatedBy:              purchaseOrder.CreatedBy,
		CommittedAt:            purchaseOrder.CommittedAt,
		CommittedBy:            purchaseOrder.CommittedBy,
		AssignedTo:             purchaseOrder.AssignedTo,
		AssignedBy:             purchaseOrder.AssignedBy,
		AssignedAt:             purchaseOrder.AssignedAt,
		Locked:                 purchaseOrder.Locked,
		LockedBy:               purchaseOrder.LockedBy,
		PurchaseOrderItems:     poi,
		Receiving:              grs,
	}

	jsonRes, _ := json.Marshal(res)
	fmt.Println(string(jsonRes))
	return
}

func (s *PurchaseOrderService) Create(ctx context.Context, req *dto.CreatePurchaseOrderRequest) (res dto.PurchaseOrderResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PurchaseOrderService.Create")
	defer span.End()

	var (
		r                                                    model.PurchaseOrder
		items                                                []*model.PurchaseOrderItem
		resItems                                             []*dto.PurchaseOrderItemResponse
		result                                               *model.PurchaseOrder
		vendor                                               *model.Vendor
		recognitionDate, etaDate                             time.Time
		totalCharge, totalPrice, totalTaxAmount, totalWeight float64
		productList                                          = make(map[int64]bool)
	)

	// cek vendor
	vendor, err = s.RepositoryVendor.GetDetail(ctx, req.VendorID, "")
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("vendor_id")
		return
	}

	// cek site
	_, err = s.RepositorySite.GetDetail(ctx, req.SiteID, "")
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("site_id")
		return
	}

	// Time Validation
	if recognitionDate, err = time.Parse("2006-01-02", req.OrderDate); err != nil {
		err = edenlabs.ErrorInvalid("order_date")
		return
	}

	if etaDate, err = time.Parse("2006-01-02", req.StrEtaDate); err != nil {
		err = edenlabs.ErrorInvalid("estimated arrival date")
		return
	}

	if len(req.Note) > 250 {
		err = edenlabs.ErrorInvalid("note")
		return
	}

	for i, v := range req.PurchaseOrderItems {
		var (
			uom                     *model.Uom
			item                    *model.Item
			unitPriceTax, taxAmount float64
			taxableItem             int8
		)
		if v.OrderQty <= 0 {
			err = edenlabs.ErrorMustGreater("order_qty", "0")
			return
		}

		if v.UnitPrice < 0 {
			err = edenlabs.ErrorMustEqualOrGreater("unit_price", "0")
			return
		}

		if v.TaxPercentage < 0 {
			err = edenlabs.ErrorMustGreater("tax_percantage", "0")
			return
		}

		if len(v.Note) > 100 {
			err = edenlabs.ErrorMustEqualOrGreater("note", "100")
			return
		}

		if v.ItemID == 0 {
			err = edenlabs.ErrorInvalid("product_id")
			return
		}

		item, err = s.RepositoryItem.GetDetail(ctx, v.ItemID, "")
		if err != nil {
			err = edenlabs.ErrorInvalid("product")
			return
		}

		uom, err = s.RepositoryUom.GetDetail(ctx, item.UomID, "")
		if err != nil {
			err = edenlabs.ErrorInvalid("product_uom")
			return
		}

		if uom.DecimalEnabled == 2 {
			if math.Mod(v.OrderQty, 1) != 0 {
				err = edenlabs.ErrorInvalid("order_qty")
				return
			}
		}

		unitPriceInput := v.UnitPrice

		if item.Taxable != "" {
			taxableItem = 1
		}
		taxAmount = math.Round((unitPriceInput * v.TaxPercentage / 100) * v.OrderQty)
		unitPriceTax = math.Round(unitPriceInput * (100 + v.TaxPercentage) / 100)

		isIncludeTax := v.IncludeTax == 1
		isNotTaxableItem := taxableItem != 1

		if isIncludeTax {
			unitPriceNonTax := math.Round(unitPriceInput * 100 / (100 + v.TaxPercentage))
			unitPriceTax := unitPriceInput
			taxAmount = math.Round((unitPriceTax - unitPriceNonTax) * v.OrderQty)
			v.UnitPrice = unitPriceNonTax
		}

		if isNotTaxableItem {
			taxAmount = 0
			unitPriceTax = 0
		}

		subtotal := v.OrderQty * v.UnitPrice

		// Summarize all the item tax amount
		totalTaxAmount += taxAmount
		totalPrice = totalPrice + subtotal
		totalWeight = totalWeight + (v.OrderQty * item.UnitWeightConversion)

		if _, exist := productList[v.ItemID]; exist {
			err = edenlabs.ErrorDuplicate("product")
			return
		}

		productList[v.ItemID] = true
		items = append(items, &model.PurchaseOrderItem{
			ItemID:        v.ItemID,
			OrderQty:      v.OrderQty,
			UnitPrice:     v.UnitPrice,
			TaxableItem:   int32(taxableItem),
			IncludeTax:    int32(v.IncludeTax),
			TaxPercentage: v.TaxPercentage,
			TaxAmount:     taxAmount,
			UnitPriceTax:  unitPriceTax,
			Subtotal:      subtotal,
			Weight:        v.OrderQty * item.UnitWeightConversion,
			Note:          v.Note,
			PurchaseQty:   v.PurchaseQty,
		})
		resItems = append(resItems, &dto.PurchaseOrderItemResponse{
			ID:            int64(i + 1),
			ItemID:        v.ItemID,
			OrderQty:      v.OrderQty,
			UnitPrice:     v.UnitPrice,
			TaxableItem:   int32(taxableItem),
			IncludeTax:    int32(v.IncludeTax),
			TaxPercentage: v.TaxPercentage,
			TaxAmount:     taxAmount,
			UnitPriceTax:  unitPriceTax,
			Subtotal:      subtotal,
			Weight:        v.OrderQty * item.UnitWeightConversion,
			Note:          v.Note,
			PurchaseQty:   v.PurchaseQty,
		})
	}

	totalCharge = totalPrice + req.DeliveryFee + (req.TaxPct * totalPrice / 100) + totalTaxAmount

	r = model.PurchaseOrder{
		Code:             r.Code,
		VendorID:         req.VendorID,
		SiteID:           req.SiteID,
		TermPaymentPurID: vendor.PaymentTermID,
		Status:           5,
		RecognitionDate:  recognitionDate,
		EtaDate:          etaDate,
		EtaTime:          req.EtaTime,
		TaxPct:           req.TaxPct,
		DeliveryFee:      req.DeliveryFee,
		TotalPrice:       totalPrice,
		TaxAmount:        totalTaxAmount,
		TotalCharge:      utils.Rounder(totalCharge, 0.5, 2),
		TotalInvoice:     0,
		TotalWeight:      totalWeight,
		Note:             req.Note,
		CreatedAt:        time.Now(),
		CreatedBy:        1,
		CreatedFrom:      0,
		HasFinishedGr:    2,
		Locked:           2,
		Latitude:         req.Latitude,
		Longitude:        req.Longitude,
		DeltaPrint:       0,
	}
	result, err = s.RepositoryPurchaseOrder.CreateWithItem(ctx, &r, items)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.PurchaseOrderResponse{
		ID:                     result.ID,
		Code:                   result.Code,
		VendorID:               result.VendorID,
		SiteID:                 result.SiteID,
		TermPaymentPurID:       result.TermPaymentPurID,
		VendorClassificationID: result.VendorClassificationID,
		PurchasePlanID:         result.PurchasePlanID,
		ConsolidatedShipmentID: result.ConsolidatedShipmentID,
		Status:                 result.Status,
		RecognitionDate:        result.RecognitionDate,
		EtaDate:                result.EtaDate,
		SiteAddress:            result.SiteAddress,
		EtaTime:                result.EtaTime,
		TaxPct:                 result.TaxPct,
		DeliveryFee:            result.DeliveryFee,
		TotalPrice:             result.TotalPrice,
		TaxAmount:              result.TaxAmount,
		TotalCharge:            result.TotalCharge,
		TotalInvoice:           result.TotalInvoice,
		TotalWeight:            result.TotalWeight,
		Note:                   result.Note,
		DeltaPrint:             result.DeltaPrint,
		Latitude:               result.Latitude,
		Longitude:              result.Longitude,
		UpdatedAt:              result.UpdatedAt,
		UpdatedBy:              result.UpdatedBy,
		CreatedFrom:            result.CreatedFrom,
		HasFinishedGr:          result.HasFinishedGr,
		CreatedAt:              result.CreatedAt,
		CreatedBy:              result.CreatedBy,
		CommittedAt:            result.CommittedAt,
		CommittedBy:            result.CommittedBy,
		AssignedTo:             result.AssignedTo,
		AssignedBy:             result.AssignedBy,
		AssignedAt:             result.AssignedAt,
		Locked:                 result.Locked,
		LockedBy:               result.LockedBy,
		PurchaseOrderItems:     resItems,
	}

	jsonRes, _ := json.Marshal(req)
	fmt.Println(string(jsonRes))
	return
}

func (s *PurchaseOrderService) Commit(ctx context.Context, id int64) (err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PurchaseOrderService.SubmitTaskVisitFU")
	defer span.End()

	// check existed po with id
	var purchaseOrder *model.PurchaseOrder
	purchaseOrder, err = s.RepositoryPurchaseOrder.GetDetail(ctx, id, "")
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if purchaseOrder.Status != 5 {
		err = edenlabs.ErrorMustActive("status")
		return
	}

	err = s.RepositoryPurchaseOrder.CommitPurchaseOrder(ctx, &model.PurchaseOrder{
		ID:          id,
		Status:      int32(statusx.ConvertStatusName(statusx.Active)),
		CommittedAt: time.Now(),
		CommittedBy: 1,
	})
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *PurchaseOrderService) Cancel(ctx context.Context, id int64) (err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PurchaseOrderService.SubmitTaskVisitFU")
	defer span.End()

	// check existed po with id
	var (
		purchaseOrder *model.PurchaseOrder
	)
	purchaseOrder, err = s.RepositoryPurchaseOrder.GetDetail(ctx, id, "")
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("id")
		return
	}

	if purchaseOrder.Status != 5 {
		err = edenlabs.ErrorMustActive("status")
		return
	}

	if purchaseOrder.PurchasePlanID != 0 {
		_, err = s.RepositoryPurchasePlan.GetDetail(ctx, purchaseOrder.PurchasePlanID, "")
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("purchase_plan_id")
			return
		}
	}

	if purchaseOrder.DeltaPrint > 0 {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("delta_print")
	}

	err = s.RepositoryPurchaseOrder.CancelPurchaseOrder(ctx, &model.PurchaseOrder{
		ID:     id,
		Status: int32(statusx.ConvertStatusName(statusx.Active)),
	})
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *PurchaseOrderService) Update(ctx context.Context, req *dto.UpdatePurchaseOrderRequest) (res dto.PurchaseOrderResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PurchaseOrderService.Update")
	defer span.End()

	var (
		r                                                    model.PurchaseOrder
		resPO                                                *model.PurchaseOrder
		items                                                []*model.PurchaseOrderItem
		resItems                                             []*dto.PurchaseOrderItemResponse
		result                                               *model.PurchaseOrder
		vendor                                               *model.Vendor
		recognitionDate, etaDate                             time.Time
		totalCharge, totalPrice, totalTaxAmount, totalWeight float64
		productList                                          = make(map[int64]bool)
	)

	// cek po
	resPO, err = s.RepositoryPurchaseOrder.GetDetail(ctx, req.Id, "")
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("id")
		return
	}

	if resPO.Status != 5 && resPO.Status != 1 {
		err = edenlabs.ErrorInvalid("id")
		return
	}

	// cek site
	_, err = s.RepositorySite.GetDetail(ctx, req.SiteID, "")
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("site_id")
		return
	}

	// cek vendor
	vendor, err = s.RepositoryVendor.GetDetail(ctx, req.VendorID, "")
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("vendor_id")
		return
	}

	// Time Validation
	if recognitionDate, err = time.Parse("2006-01-02", req.OrderDate); err != nil {
		err = edenlabs.ErrorInvalid("order_date")
		return
	}

	if etaDate, err = time.Parse("2006-01-02", req.StrEtaDate); err != nil {
		err = edenlabs.ErrorInvalid("estimated arrival date")
		return
	}

	if _, err = time.Parse("15:04", req.EtaTime); err != nil {
		err = edenlabs.ErrorInvalid("eta_time")
		return
	}

	if len(req.Note) > 250 {
		err = edenlabs.ErrorInvalid("note")
		return
	}

	for i, v := range req.PurchaseOrderItems {
		var (
			uom                     *model.Uom
			item                    *model.Item
			unitPriceTax, taxAmount float64
			taxableItem             int8
		)

		// cek poi
		_, err = s.RepositoryPurchaseOrderItem.GetDetail(ctx, v.Id, "")
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("purchase_order_item_id")
			return
		}

		if v.OrderQty <= 0 {
			err = edenlabs.ErrorMustGreater("order_qty", "0")
			return
		}

		if v.UnitPrice < 0 {
			err = edenlabs.ErrorMustEqualOrGreater("unit_price", "0")
			return
		}

		if v.TaxPercentage < 0 {
			err = edenlabs.ErrorMustGreater("tax_percantage", "0")
			return
		}

		if len(v.Note) > 100 {
			err = edenlabs.ErrorMustEqualOrGreater("note", "100")
			return
		}

		if v.ItemID == 0 {
			err = edenlabs.ErrorInvalid("product_id")
			return
		}

		item, err = s.RepositoryItem.GetDetail(ctx, v.ItemID, "")
		if err != nil {
			err = edenlabs.ErrorInvalid("product")
			return
		}

		uom, err = s.RepositoryUom.GetDetail(ctx, v.ItemID, "")
		if err != nil {
			err = edenlabs.ErrorInvalid("product_uom")
			return
		}

		if uom.DecimalEnabled == 2 {
			if math.Mod(v.OrderQty, 1) != 0 {
				err = edenlabs.ErrorInvalid("order_qty")
			}
		}

		unitPriceInput := v.UnitPrice

		if item.Taxable != "" {
			taxableItem = 1
		}
		taxAmount = math.Round((unitPriceInput * v.TaxPercentage / 100) * v.OrderQty)
		unitPriceTax = math.Round(unitPriceInput * (100 + v.TaxPercentage) / 100)

		isIncludeTax := v.IncludeTax == 1
		isNotTaxableItem := taxableItem != 1

		if isIncludeTax {
			unitPriceNonTax := math.Round(unitPriceInput * 100 / (100 + v.TaxPercentage))
			unitPriceTax := unitPriceInput
			taxAmount = math.Round((unitPriceTax - unitPriceNonTax) * v.OrderQty)
			v.UnitPrice = unitPriceNonTax
		}

		if isNotTaxableItem {
			taxAmount = 0
			unitPriceTax = 0
		}

		subtotal := v.OrderQty * v.UnitPrice

		// Summarize all the item tax amount
		totalTaxAmount += taxAmount
		totalPrice = totalPrice + subtotal
		totalWeight = totalWeight + (v.OrderQty * item.UnitWeightConversion)

		if _, exist := productList[v.ItemID]; exist {
			err = edenlabs.ErrorDuplicate("product")
			return
		}

		productList[v.ItemID] = true
		items = append(items, &model.PurchaseOrderItem{
			ItemID:        v.ItemID,
			OrderQty:      v.OrderQty,
			UnitPrice:     v.UnitPrice,
			TaxableItem:   int32(taxableItem),
			IncludeTax:    int32(v.IncludeTax),
			TaxPercentage: v.TaxPercentage,
			TaxAmount:     taxAmount,
			UnitPriceTax:  unitPriceTax,
			Subtotal:      subtotal,
			Weight:        v.OrderQty * item.UnitWeightConversion,
			Note:          v.Note,
			PurchaseQty:   v.PurchaseQty,
		})

		resItems = append(resItems, &dto.PurchaseOrderItemResponse{
			ID:            int64(i + 1),
			ItemID:        v.ItemID,
			OrderQty:      v.OrderQty,
			UnitPrice:     v.UnitPrice,
			TaxableItem:   int32(taxableItem),
			IncludeTax:    int32(v.IncludeTax),
			TaxPercentage: v.TaxPercentage,
			TaxAmount:     taxAmount,
			UnitPriceTax:  unitPriceTax,
			Subtotal:      subtotal,
			Weight:        v.OrderQty * item.UnitWeightConversion,
			Note:          v.Note,
			PurchaseQty:   v.PurchaseQty,
		})
	}

	totalCharge = totalPrice + req.DeliveryFee + (req.TaxPct * totalPrice / 100) + totalTaxAmount

	r = model.PurchaseOrder{
		Code:             r.Code,
		VendorID:         req.VendorID,
		SiteID:           req.SiteID,
		TermPaymentPurID: vendor.PaymentTermID,
		Status:           5,
		RecognitionDate:  recognitionDate,
		EtaDate:          etaDate,
		EtaTime:          req.EtaTime,
		TaxPct:           req.TaxPct,
		DeliveryFee:      req.DeliveryFee,
		TotalPrice:       totalPrice,
		TaxAmount:        totalTaxAmount,
		TotalCharge:      utils.Rounder(totalCharge, 0.5, 2),
		TotalInvoice:     0,
		TotalWeight:      totalWeight,
		Note:             req.Note,
		CreatedAt:        time.Now(),
		CreatedBy:        1,
		CreatedFrom:      0,
		HasFinishedGr:    2,
		Locked:           2,
		Latitude:         req.Latitude,
		Longitude:        req.Longitude,
		DeltaPrint:       0,
	}
	result, err = s.RepositoryPurchaseOrder.UpdateWithItem(ctx, &r, items)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.PurchaseOrderResponse{
		ID:                     result.ID,
		Code:                   result.Code,
		VendorID:               result.VendorID,
		SiteID:                 result.SiteID,
		TermPaymentPurID:       result.TermPaymentPurID,
		VendorClassificationID: result.VendorClassificationID,
		PurchasePlanID:         result.PurchasePlanID,
		ConsolidatedShipmentID: result.ConsolidatedShipmentID,
		Status:                 result.Status,
		RecognitionDate:        result.RecognitionDate,
		EtaDate:                result.EtaDate,
		SiteAddress:            result.SiteAddress,
		EtaTime:                result.EtaTime,
		TaxPct:                 result.TaxPct,
		DeliveryFee:            result.DeliveryFee,
		TotalPrice:             result.TotalPrice,
		TaxAmount:              result.TaxAmount,
		TotalCharge:            result.TotalCharge,
		TotalInvoice:           result.TotalInvoice,
		TotalWeight:            result.TotalWeight,
		Note:                   result.Note,
		DeltaPrint:             result.DeltaPrint,
		Latitude:               result.Latitude,
		Longitude:              result.Longitude,
		UpdatedAt:              result.UpdatedAt,
		UpdatedBy:              result.UpdatedBy,
		CreatedFrom:            result.CreatedFrom,
		HasFinishedGr:          result.HasFinishedGr,
		CreatedAt:              result.CreatedAt,
		CreatedBy:              result.CreatedBy,
		CommittedAt:            result.CommittedAt,
		CommittedBy:            result.CommittedBy,
		AssignedTo:             result.AssignedTo,
		AssignedBy:             result.AssignedBy,
		AssignedAt:             result.AssignedAt,
		Locked:                 result.Locked,
		LockedBy:               result.LockedBy,
		PurchaseOrderItems:     resItems,
	}

	jsonRes, _ := json.Marshal(req)
	fmt.Println(string(jsonRes))
	return
}

func (s *PurchaseOrderService) UpdateProduct(ctx context.Context, req *dto.UpdateProductPurchaseOrderRequest) (res dto.PurchaseOrderResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PurchaseOrderService.Update")
	defer span.End()

	var (
		r                                                    model.PurchaseOrder
		resPO                                                *model.PurchaseOrder
		items                                                []*model.PurchaseOrderItem
		resItems                                             []*dto.PurchaseOrderItemResponse
		result                                               *model.PurchaseOrder
		etaDate                                              time.Time
		totalCharge, totalPrice, totalTaxAmount, totalWeight float64
		productList                                          = make(map[int64]bool)
	)

	// cek po
	resPO, err = s.RepositoryPurchaseOrder.GetDetail(ctx, req.Id, "")
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("id")
		return
	}

	if resPO.Status != 5 && resPO.Status != 1 {
		err = edenlabs.ErrorInvalid("id")
		return
	}

	// TODO:
	// CHECK GR/ITEM TRANSFER ON GP
	// CHECK Debit Note on GP
	// CHECK Purchase Invoice on GP

	for i, v := range req.PurchaseOrderItems {
		var (
			uom                     *model.Uom
			item                    *model.Item
			unitPriceTax, taxAmount float64
			taxableItem             int8
		)

		// cek poi
		_, err = s.RepositoryPurchaseOrderItem.GetDetail(ctx, v.Id, "")
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("purchase_order_item_id")
			return
		}

		if v.OrderQty <= 0 {
			err = edenlabs.ErrorMustGreater("order_qty", "0")
			return
		}

		if v.UnitPrice < 0 {
			err = edenlabs.ErrorMustEqualOrGreater("unit_price", "0")
			return
		}

		if v.TaxPercentage < 0 {
			err = edenlabs.ErrorMustGreater("tax_percantage", "0")
			return
		}

		if len(v.Note) > 100 {
			err = edenlabs.ErrorMustEqualOrGreater("note", "100")
			return
		}

		if v.ItemID == 0 {
			err = edenlabs.ErrorInvalid("product_id")
			return
		}

		item, err = s.RepositoryItem.GetDetail(ctx, v.ItemID, "")
		if err != nil {
			err = edenlabs.ErrorInvalid("product")
			return
		}

		uom, err = s.RepositoryUom.GetDetail(ctx, v.ItemID, "")
		if err != nil {
			err = edenlabs.ErrorInvalid("product_uom")
			return
		}

		if uom.DecimalEnabled == 2 {
			if math.Mod(v.OrderQty, 1) != 0 {
				err = edenlabs.ErrorInvalid("order_qty")
			}
		}

		unitPriceInput := v.UnitPrice

		if item.Taxable != "" {
			taxableItem = 1
		}

		//TODO: CHECK GR

		taxAmount = math.Round((unitPriceInput * v.TaxPercentage / 100) * v.OrderQty)
		unitPriceTax = math.Round(unitPriceInput * (100 + v.TaxPercentage) / 100)

		isIncludeTax := v.IncludeTax == 1
		isNotTaxableItem := taxableItem != 1

		if isIncludeTax {
			unitPriceNonTax := math.Round(unitPriceInput * 100 / (100 + v.TaxPercentage))
			unitPriceTax := unitPriceInput
			taxAmount = math.Round((unitPriceTax - unitPriceNonTax) * v.OrderQty)
			v.UnitPrice = unitPriceNonTax
		}

		if isNotTaxableItem {
			taxAmount = 0
			unitPriceTax = 0
		}

		subtotal := v.OrderQty * v.UnitPrice

		// Summarize all the item tax amount
		totalTaxAmount += taxAmount
		totalPrice = totalPrice + subtotal
		totalWeight = totalWeight + (v.OrderQty * item.UnitWeightConversion)

		if _, exist := productList[v.ItemID]; exist {
			err = edenlabs.ErrorDuplicate("product")
			return
		}

		productList[v.ItemID] = true
		items = append(items, &model.PurchaseOrderItem{
			ItemID:        v.ItemID,
			OrderQty:      v.OrderQty,
			UnitPrice:     v.UnitPrice,
			TaxableItem:   int32(taxableItem),
			IncludeTax:    int32(v.IncludeTax),
			TaxPercentage: v.TaxPercentage,
			TaxAmount:     taxAmount,
			UnitPriceTax:  unitPriceTax,
			Subtotal:      subtotal,
			Weight:        v.OrderQty * item.UnitWeightConversion,
			Note:          v.Note,
			PurchaseQty:   v.PurchaseQty,
		})

		resItems = append(resItems, &dto.PurchaseOrderItemResponse{
			ID:            int64(i + 1),
			ItemID:        v.ItemID,
			OrderQty:      v.OrderQty,
			UnitPrice:     v.UnitPrice,
			TaxableItem:   int32(taxableItem),
			IncludeTax:    int32(v.IncludeTax),
			TaxPercentage: v.TaxPercentage,
			TaxAmount:     taxAmount,
			UnitPriceTax:  unitPriceTax,
			Subtotal:      subtotal,
			Weight:        v.OrderQty * item.UnitWeightConversion,
			Note:          v.Note,
			PurchaseQty:   v.PurchaseQty,
		})
	}

	totalCharge = totalPrice + req.DeliveryFee + (req.TaxPct * totalPrice / 100) + totalTaxAmount

	r = model.PurchaseOrder{
		Code:             r.Code,
		VendorID:         resPO.VendorID,
		SiteID:           resPO.SiteID,
		TermPaymentPurID: 1,
		Status:           5,
		RecognitionDate:  resPO.RecognitionDate,
		EtaDate:          etaDate,
		EtaTime:          resPO.EtaTime,
		TaxPct:           req.TaxPct,
		DeliveryFee:      req.DeliveryFee,
		TotalPrice:       totalPrice,
		TaxAmount:        totalTaxAmount,
		TotalCharge:      utils.Rounder(totalCharge, 0.5, 2),
		TotalInvoice:     0,
		TotalWeight:      totalWeight,
		Note:             resPO.Note,
		CreatedAt:        time.Now(),
		CreatedBy:        1,
		CreatedFrom:      0,
		HasFinishedGr:    2,
		Locked:           2,
		Latitude:         resPO.Latitude,
		Longitude:        resPO.Longitude,
		DeltaPrint:       0,
	}

	result, err = s.RepositoryPurchaseOrder.UpdateProduct(ctx, &r, items)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.PurchaseOrderResponse{
		ID:                     result.ID,
		Code:                   result.Code,
		VendorID:               result.VendorID,
		SiteID:                 result.SiteID,
		TermPaymentPurID:       result.TermPaymentPurID,
		VendorClassificationID: result.VendorClassificationID,
		PurchasePlanID:         result.PurchasePlanID,
		ConsolidatedShipmentID: result.ConsolidatedShipmentID,
		Status:                 result.Status,
		RecognitionDate:        result.RecognitionDate,
		EtaDate:                result.EtaDate,
		SiteAddress:            result.SiteAddress,
		EtaTime:                result.EtaTime,
		TaxPct:                 result.TaxPct,
		DeliveryFee:            result.DeliveryFee,
		TotalPrice:             result.TotalPrice,
		TaxAmount:              result.TaxAmount,
		TotalCharge:            result.TotalCharge,
		TotalInvoice:           result.TotalInvoice,
		TotalWeight:            result.TotalWeight,
		Note:                   result.Note,
		DeltaPrint:             result.DeltaPrint,
		Latitude:               result.Latitude,
		Longitude:              result.Longitude,
		UpdatedAt:              result.UpdatedAt,
		UpdatedBy:              result.UpdatedBy,
		CreatedFrom:            result.CreatedFrom,
		HasFinishedGr:          result.HasFinishedGr,
		CreatedAt:              result.CreatedAt,
		CreatedBy:              result.CreatedBy,
		CommittedAt:            result.CommittedAt,
		CommittedBy:            result.CommittedBy,
		AssignedTo:             result.AssignedTo,
		AssignedBy:             result.AssignedBy,
		AssignedAt:             result.AssignedAt,
		Locked:                 result.Locked,
		LockedBy:               result.LockedBy,
		PurchaseOrderItems:     resItems,
	}

	jsonRes, _ := json.Marshal(req)
	fmt.Println(string(jsonRes))
	return
}

func (s *PurchaseOrderService) GetGP(ctx context.Context, req *pb.GetPurchaseOrderGPListRequest) (res *pb.GetPurchaseOrderGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PurchaseOrderService.GetGP")
	defer span.End()
	var statusInt []string

	params := map[string]string{
		"interid":    global.EnvDatabaseGP,
		"PageNumber": strconv.Itoa(int(req.Offset)),
		"PageSize":   strconv.Itoa(int(req.Limit)),
	}

	if req.Ponumber != "" {
		params["PONUMBER"] = req.Ponumber
	}

	if req.Ponumberlike != "" {
		params["PONUMBER_LIKE"] = req.Ponumberlike
	}

	if req.PurchasePlanId != "" {
		params["prp_purchaseplan_no"] = req.PurchasePlanId
	}

	// add status
	if len(req.Status) > 0 {
		for _, v := range req.Status {
			statusInt = append(statusInt, strconv.Itoa(int(v)))
		}
		params["POSTATUS"] = strings.Join(statusInt, ",")
	}

	if req.ReqdateFrom != "" && req.ReqdateFrom != "0001-01-01" {
		params["REQDATE_FROM"] = req.ReqdateFrom
	}

	if req.ReqdateTo != "" && req.ReqdateTo != "0001-01-01" {
		params["REQDATE_TO"] = req.ReqdateTo
	}

	if req.IsNotConsolidated {
		params["is_consolidated"] = "no"
	}

	if req.IsPurchasePlan != "" {
		params["is_purchase_plan"] = "yes"
	}

	if req.Orderby != "" {
		params["orderby"] = "desc"
	}

	if req.PrpPurchaseplanUser != "" {
		params["prp_purchaseplan_user"] = req.PrpPurchaseplanUser
	}

	if req.Locncode != "" {
		params["LOCNCODE"] = req.Locncode
	}

	if req.PrpCsNo != "" {
		params["prp_cs_no"] = req.PrpCsNo
	}

	if req.VendorId != "" {
		params["VENDORID"] = req.VendorId
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "purchaseorder/getall", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *PurchaseOrderService) GetDetailGP(ctx context.Context, req *pb.GetPurchaseOrderGPDetailRequest) (res *pb.GetPurchaseOrderGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PurchaseOrderService.GetDetailGP")
	defer span.End()

	params := map[string]string{
		"interid": global.EnvDatabaseGP,
		"id":      req.Id,
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "purchaseorder/getbyid", nil, &res, params)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *PurchaseOrderService) CreateGP(ctx context.Context, req *dto.CreatePurchaseOrderGPRequest) (res dto.CommonPurchaseOrderGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PurchaseOrderService.Create")
	defer span.End()

	req.Interid = global.EnvDatabaseGP

	err = global.HttpRestApiToMicrosoftGP("POST", "purchaseorder/create", req, &res, nil)
	if err != nil {
		logrus.Error(err.Error())
		err = errors.New("Connection to the server could not be established")
		return
	}

	if res.Code != 200 {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = errors.New("Connection to the server could not be established")
		return
	}

	return
}

func (s *PurchaseOrderService) UpdateGP(ctx context.Context, req *dto.UpdatePurchaseOrderGPRequest) (res dto.CommonPurchaseOrderGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PurchaseOrderService.Update")
	defer span.End()

	req.Interid = global.EnvDatabaseGP

	err = global.HttpRestApiToMicrosoftGP("POST", "purchaseorder/update", req, &res, nil)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = errors.New("Connection to the server could not be established")
		return
	}

	if res.Code != 200 {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = errors.New("Connection to the server could not be established")
		return
	}

	return
}

func (s *PurchaseOrderService) CommitPurchaseOrderGP(ctx context.Context, req *pb.CommitPurchaseOrderGPRequest) (res *pb.CreateTransferRequestGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PurchaseOrderService.CommitPurchaseOrderGP")
	defer span.End()

	req.Interid = global.EnvDatabaseGP

	err = global.HttpRestApiToMicrosoftGP("PUT", "purchaseorder/commit", req, &res, nil)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *PurchaseOrderService) CancelPurchaseOrderGP(ctx context.Context, req *dto.CancelPurchaseOrderGPRequest) (res *dto.CancelPurchaseOrderGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PurchaseOrderService.Update")
	defer span.End()

	req.Interid = global.EnvDatabaseGP

	err = global.HttpRestApiToMicrosoftGP("POST", "purchaseorder/cancel", req, &res, nil)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = errors.New("Connection to the server could not be established")
		return
	}

	if res.Code != 200 {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = errors.New("Connection to the server could not be established")
		return
	}

	return
}

func (s *PurchaseOrderService) CreateConsolidatedShipmentGP(ctx context.Context, req *dto.CreateConsolidatedShipmentGPRequest) (res dto.CreateConsolidatedShipmentGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PurchaseOrderService.CreateConsolidatedShipmentGP")
	defer span.End()

	var (
		poList = make(map[string]bool)
	)

	req.Interid = global.EnvDatabaseGP

	if req.PrpDriverName == "" {
		edenlabs.ErrorRequired("prp_driver_name")
	}

	if req.PrVehicleNumber == "" {
		edenlabs.ErrorRequired("pr_vehicle_number")
	}

	if req.PhoneName == "" {
		edenlabs.ErrorRequired("phonname")
	}

	for _, v := range req.PurchaseOrders {

		if _, exist := poList[v.Ponumber]; exist {
			err = edenlabs.ErrorDuplicate("purchase_order")
			return
		}
	}

	err = global.HttpRestApiToMicrosoftGP("POST", "purchaseorder/update/CSDetail", req, &res, nil)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = errors.New("Connection to the server could not be established")
		return

	}

	if res.Code != 200 {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = errors.New("Connection to the server could not be established")
		return
	}

	return
}

func (s *PurchaseOrderService) UpdateConsolidatedShipmentGP(ctx context.Context, req *dto.UpdateConsolidatedShipmentGPRequest) (res dto.UpdateConsolidatedShipmentGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PurchaseOrderService.UpdateConsolidatedShipmentGP")
	defer span.End()

	var (
		poList = make(map[string]bool)
	)

	req.Interid = global.EnvDatabaseGP

	if req.PrpDriverName == "" {
		edenlabs.ErrorRequired("prp_driver_name")
	}

	if req.PrVehicleNumber == "" {
		edenlabs.ErrorRequired("pr_vehicle_number")
	}

	if req.PhoneName == "" {
		edenlabs.ErrorRequired("phonname")
	}

	for _, v := range req.PurchaseOrders {

		if _, exist := poList[v.Ponumber]; exist {
			err = edenlabs.ErrorDuplicate("purchase_order")
			return
		}
	}

	err = global.HttpRestApiToMicrosoftGP("POST", "purchaseorder/update/CSDetail", req, &res, nil)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = errors.New("Connection to the server could not be established")
		return

	}

	if res.Code != 200 {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = errors.New("Connection to the server could not be established")
		return
	}

	return
}
