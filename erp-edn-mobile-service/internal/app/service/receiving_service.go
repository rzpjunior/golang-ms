package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/constants"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/global"
	"google.golang.org/protobuf/types/known/timestamppb"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/internal/app/dto"
	auditService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/audit_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

func NewServiceReceiving() IReceivingService {
	m := new(ReceivingService)
	m.opt = global.Setup.Common
	return m
}

type IReceivingService interface {
	Get(ctx context.Context, req dto.ReceivingListRequest) (res []*dto.ReceivingResponse, err error)
	GetById(ctx context.Context, req dto.ReceivingDetailRequest) (res *dto.ReceivingResponse, err error)
	Create(ctx context.Context, req *dto.CreateReceivingRequest) (res *dto.ReceivingResponse, err error)
	Confirm(ctx context.Context, req *dto.ConfirmReceivingRequest) (res *dto.ReceivingResponse, err error)

	// gp integrated
	GetListGP(ctx context.Context, req dto.GetGoodsReceiptGPListRequest) (res []*dto.ReceivingResponse, total int64, err error)
	GetDetailGP(ctx context.Context, id string) (res *bridgeService.GoodsReceiptGP, err error)
	CreateGP(ctx context.Context, req *dto.CreateGoodsReceiptGPRequest) (res *bridgeService.CreateTransferRequestGPResponse, err error)
	UpdateGP(ctx context.Context, req *dto.UpdateGoodsReceiptGPRequest) (res *bridgeService.CreateTransferRequestGPResponse, err error)
}

type ReceivingService struct {
	opt opt.Options
}

func (s *ReceivingService) Get(ctx context.Context, req dto.ReceivingListRequest) (res []*dto.ReceivingResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ReceivingService.Get")
	defer span.End()

	// get Item Transfer from bridge
	var grRes *bridgeService.GetReceivingListResponse
	grRes, err = s.opt.Client.BridgeServiceGrpc.GetReceivingList(ctx, &bridgeService.GetReceivingListRequest{
		Limit:   req.Limit,
		Offset:  req.Offset,
		Status:  req.Status,
		Search:  req.Search,
		OrderBy: req.OrderBy,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "purchase order")
		return
	}

	datas := []*dto.ReceivingResponse{}
	for _, it := range grRes.Data {
		datas = append(datas, &dto.ReceivingResponse{
			// ID:                  it.Id,
			Code: it.Code,
			// SiteId:              it.SiteId,
			PurchaseOrderId:     it.PurchaseOrderId,
			ItemTransferId:      it.ItemTransferId,
			InboundType:         int8(it.InboundType),
			ValidSupplierReturn: int8(it.ValidSupplierReturn),
			AtaDate:             it.AtaDate.AsTime(),
			AtaTime:             it.AtaTime,
			StockType:           int8(it.StockType),
			TotalWeight:         it.TotalWeight,
			Note:                it.Note,
			Status:              int8(it.Status),
			Locked:              int8(it.Locked),
			CreatedAt:           it.CreatedAt.AsTime(),
			CreatedBy:           it.CreatedBy,
			ConfirmedAt:         it.ConfirmedAt.AsTime(),
			ConfirmedBy:         it.ConfirmedBy,
			LockedBy:            it.LockedBy,
			UpdatedAt:           it.UpdatedAt.AsTime(),
			UpdatedBy:           it.UpdatedBy,
		})
	}
	res = datas

	return
}

func (s *ReceivingService) GetById(ctx context.Context, req dto.ReceivingDetailRequest) (res *dto.ReceivingResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ReceivingService.GetById")
	defer span.End()

	// get Purchase Order from bridge
	var (
		itRes            *bridgeService.GetReceivingDetailResponse
		iti              []*dto.ReceivingItemResponse
		poRes            *bridgeService.GetPurchaseOrderDetailResponse
		poiRes           *bridgeService.GetPurchaseOrderItemDetailResponse
		vendorRes        *bridgeService.GetVendorDetailResponse
		itemRes          *bridgeService.GetItemDetailResponse
		uomRes           *bridgeService.GetUomDetailResponse
		itemTransRes     *bridgeService.GetItemTransferDetailResponse
		itemTransItemRes *bridgeService.GetItemTransferItemDetailResponse
		siteRes          *bridgeService.GetSiteDetailResponse
		regionRes        *bridgeService.GetRegionDetailResponse
		po               *dto.PurchaseOrderResponse
		poi              *dto.PurchaseOrderItemResponse
		it               *dto.ItemTransferResponse
		itItem           *dto.ItemTransferItemResponse
		site             *dto.SiteResponse
	)
	itRes, err = s.opt.Client.BridgeServiceGrpc.GetReceivingDetail(ctx, &bridgeService.GetReceivingDetailRequest{
		Id: int64(req.Id),
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "receiving")
		return
	}

	// mapping po
	if itRes.Data.PurchaseOrderId != 0 {
		poRes, err = s.opt.Client.BridgeServiceGrpc.GetPurchaseOrderDetail(ctx, &bridgeService.GetPurchaseOrderDetailRequest{Id: itRes.Data.PurchaseOrderId})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "purchase order")
			return
		}

		vendorRes, err = s.opt.Client.BridgeServiceGrpc.GetVendorDetail(ctx, &bridgeService.GetVendorDetailRequest{Id: poRes.Data.VendorId})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "vendor")
			return
		}

		po = &dto.PurchaseOrderResponse{
			// ID:                     poRes.Data.Id,
			Code: poRes.Data.Code,
			// VendorID:               poRes.Data.VendorId,
			SiteID:                 poRes.Data.SiteId,
			TermPaymentPurID:       poRes.Data.TermPaymentPurId,
			VendorClassificationID: poRes.Data.VendorClassificationId,
			PurchasePlanID:         poRes.Data.PurchasePlanId,
			ConsolidatedShipmentID: poRes.Data.ConsolidatedShipmentId,
			Status:                 poRes.Data.Status,
			RecognitionDate:        poRes.Data.RecognitionDate.AsTime(),
			EtaDate:                poRes.Data.EtaDate.AsTime(),
			SiteAddress:            poRes.Data.SiteAddress,
			EtaTime:                poRes.Data.EtaTime,
			TaxPct:                 poRes.Data.TaxPct,
			DeliveryFee:            poRes.Data.DeliveryFee,
			TotalPrice:             poRes.Data.TotalPrice,
			TaxAmount:              poRes.Data.TaxAmount,
			TotalCharge:            poRes.Data.TotalCharge,
			TotalInvoice:           poRes.Data.TotalInvoice,
			TotalWeight:            poRes.Data.TotalWeight,
			Note:                   poRes.Data.Note,
			DeltaPrint:             poRes.Data.DeltaPrint,
			Latitude:               poRes.Data.Latitude,
			Longitude:              poRes.Data.Longitude,
			CreatedFrom:            poRes.Data.CreatedFrom,
			HasFinishedGr:          int8(poRes.Data.HasFinishedGr),
			CommittedAt:            poRes.Data.CommittedAt.AsTime(),
			CommittedBy:            poRes.Data.CommittedBy,
			AssignedTo:             poRes.Data.AssignedTo,
			AssignedBy:             poRes.Data.AssignedBy,
			AssignedAt:             poRes.Data.AssignedAt.AsTime(),
			Locked:                 poRes.Data.Locked,
			LockedBy:               poRes.Data.LockedBy,
			CreatedBy:              poRes.Data.CreatedBy,
			CreatedAt:              poRes.Data.CreatedAt.AsTime(),
			UpdatedAt:              poRes.Data.UpdatedAt.AsTime(),
			UpdatedBy:              poRes.Data.UpdatedBy,
			Vendor: &dto.VendorResponse{
				// ID:                     vendorRes.Data.Id,
				Code: vendorRes.Data.Code,
				// VendorOrganizationID:   vendorRes.Data.VendorOrganizationId,
				VendorClassificationID: vendorRes.Data.VendorClassificationId,
				SubDistrictID:          vendorRes.Data.SubDistrictId,
				PaymentTermID:          vendorRes.Data.PaymentTermId,
				PicName:                vendorRes.Data.PicName,
				Email:                  vendorRes.Data.Email,
				PhoneNumber:            vendorRes.Data.PhoneNumber,
				Rejectable:             vendorRes.Data.Rejectable,
				Returnable:             vendorRes.Data.Rejectable,
				Address:                vendorRes.Data.Address,
				Note:                   vendorRes.Data.Note,
				Status:                 vendorRes.Data.Status,
				Latitude:               vendorRes.Data.Latitude,
				Longitude:              vendorRes.Data.Longitude,
				CreatedAt:              vendorRes.Data.CreatedAt.AsTime(),
			},
		}
	}

	// mapping item transfer
	if itRes.Data.ItemTransferId != 0 {
		itemTransRes, err = s.opt.Client.BridgeServiceGrpc.GetItemTransferDetail(ctx, &bridgeService.GetItemTransferDetailRequest{Id: itRes.Data.ItemTransferId})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "vendor")
			return
		}

		it = &dto.ItemTransferResponse{
			// ID:                 itemTransRes.Data.Id,
			Code:               itemTransRes.Data.Code,
			RequestDate:        itemTransRes.Data.RequestDate.AsTime(),
			RecognitionDate:    itemTransRes.Data.RecognitionDate.AsTime(),
			EtaDate:            itemTransRes.Data.EtaDate.AsTime(),
			EtaTime:            itemTransRes.Data.EtaTime,
			AtaDate:            itemTransRes.Data.AtaDate.AsTime(),
			AtaTime:            itemTransRes.Data.AtaTime,
			AdditionalCost:     itemTransRes.Data.AdditionalCost,
			AdditionalCostNote: itemTransRes.Data.AdditionalCostNote,
			StockType:          int8(itemTransRes.Data.StockType),
			TotalCost:          itemTransRes.Data.TotalCost,
			TotalCharge:        itemTransRes.Data.TotalCharge,
			TotalSku:           itemTransRes.Data.TotalSku,
			TotalWeight:        itemTransRes.Data.TotalWeight,
			Note:               itemTransRes.Data.Note,
			// Status:             int8(itemTransRes.Data.Status),
			Locked:    int8(itemTransRes.Data.Locked),
			LockedBy:  itemTransRes.Data.LockedBy,
			UpdatedAt: itemTransRes.Data.UpdatedAt.AsTime(),
			UpdatedBy: itemTransRes.Data.UpdatedBy,
		}
	}

	// mapping site
	if itRes.Data.SiteId != 0 {
		siteRes, err = s.opt.Client.BridgeServiceGrpc.GetSiteDetail(ctx, &bridgeService.GetSiteDetailRequest{Id: itRes.Data.SiteId})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "site")
			return
		}

		regionRes, err = s.opt.Client.BridgeServiceGrpc.GetRegionDetail(ctx, &bridgeService.GetRegionDetailRequest{Id: 1})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "site")
			return
		}

		site = &dto.SiteResponse{
			// ID:            siteRes.Data.Id,
			Code:          siteRes.Data.Code,
			Description:   siteRes.Data.Description,
			Name:          siteRes.Data.Description,
			Status:        int8(siteRes.Data.Status),
			StatusConvert: statusx.ConvertStatusValue(int8(siteRes.Data.Status)),
			CreatedAt:     siteRes.Data.CreatedAt.AsTime(),
			UpdatedAt:     siteRes.Data.UpdatedAt.AsTime(),
			Region: &dto.RegionResponse{
				// ID:            regionRes.Data.Id,
				Code:          regionRes.Data.Code,
				Description:   regionRes.Data.Description,
				Status:        int8(regionRes.Data.Status),
				StatusConvert: statusx.ConvertStatusValue(int8(regionRes.Data.Status)),
				CreatedAt:     regionRes.Data.CreatedAt.AsTime(),
				UpdatedAt:     regionRes.Data.UpdatedAt.AsTime(),
			},
		}
	}

	// mapping receiving items
	for _, item := range itRes.Data.ReceivingItems {
		var itemId int64
		if item.PurchaseOrderItem != nil {
			itemId = item.PurchaseOrderItem.ItemId
		}
		if item.ItemTransferItem != nil {
			itemId = item.PurchaseOrderItem.ItemId
		}
		itemRes, err = s.opt.Client.BridgeServiceGrpc.GetItemDetail(ctx, &bridgeService.GetItemDetailRequest{Id: itemId})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "item")
			return
		}
		uomRes, err = s.opt.Client.BridgeServiceGrpc.GetUomDetail(ctx, &bridgeService.GetUomDetailRequest{Id: itemRes.Data.UomId})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "uom")
			return
		}

		if item.PurchaseOrderItemId != 0 {
			poiRes, err = s.opt.Client.BridgeServiceGrpc.GetPurchaseOrderItemDetail(ctx, &bridgeService.GetPurchaseOrderItemDetailRequest{Id: item.PurchaseOrderItemId})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "item")
				return
			}
			poi = &dto.PurchaseOrderItemResponse{
				// ID:                 poiRes.Data.Id,
				// PurchaseOrderID:    poiRes.Data.PurchaseOrderId,
				PurchasePlanItemID: poiRes.Data.PurchasePlanItemId,
				// ItemID:             poiRes.Data.ItemId,
				OrderQty:      poiRes.Data.OrderQty,
				UnitPrice:     poiRes.Data.UnitPrice,
				TaxableItem:   poiRes.Data.TaxableItem,
				IncludeTax:    poiRes.Data.IncludeTax,
				TaxPercentage: poiRes.Data.TaxPercentage,
				TaxAmount:     poiRes.Data.TaxAmount,
				UnitPriceTax:  poiRes.Data.UnitPriceTax,
				Subtotal:      poiRes.Data.Subtotal,
				Weight:        poiRes.Data.Weight,
				Note:          poiRes.Data.Note,
				PurchaseQty:   poiRes.Data.PurchaseQty,
			}
		}

		if item.ItemTransferItemId != 0 {
			itemTransItemRes, err = s.opt.Client.BridgeServiceGrpc.GetItemTransferItemDetail(ctx, &bridgeService.GetItemTransferItemDetailRequest{Id: item.ItemTransferItemId})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "item")
				return
			}
			itItem = &dto.ItemTransferItemResponse{
				ID: itemTransItemRes.Data.Id,
				// ItemTransferID: itemTransItemRes.Data.ItemTransferId,
				// ItemID:         itemTransItemRes.Data.ItemId,
				// DeliverQty:     itemTransItemRes.Data.DeliverQty,
				ReceiveQty:  itemTransItemRes.Data.ReceiveQty,
				RequestQty:  itemTransItemRes.Data.RequestQty,
				ReceiveNote: itemTransItemRes.Data.ReceiveNote,
				UnitCost:    itemTransItemRes.Data.UnitCost,
				Subtotal:    itemTransItemRes.Data.Subtotal,
				Weight:      itemTransItemRes.Data.Weight,
				Note:        itemTransItemRes.Data.Note,
			}
		}

		var rejectReasonConvert string
		if item.RejectReason == 1 {
			rejectReasonConvert = "Lost"
		} else if item.RejectReason == 2 {
			rejectReasonConvert = "Damaged"
		}

		iti = append(iti, &dto.ReceivingItemResponse{
			// ID:                  item.Id,
			PurchaseOrderItemID: item.PurchaseOrderItemId,
			ItemTransferItemID:  item.ItemTransferItemId,
			DeliverQty:          item.DeliverQty,
			// ReceiveQty:          item.ReceiveQty,
			Weight:              item.Weight,
			Note:                item.Note,
			RejectQty:           item.RejectQty,
			RejectReason:        int8(item.RejectReason),
			RejectReasonConvert: rejectReasonConvert,
			IsDisabled:          int8(item.IsDisabled),
			Item: &dto.ItemResponse{
				// ID:   itemRes.Data.Id,
				Code: itemRes.Data.Code,
				Uom: &dto.UomResponse{
					ID:             "1",
					Code:           uomRes.Data.Code,
					Description:    uomRes.Data.Description,
					Status:         int8(uomRes.Data.Status),
					StatusConvert:  statusx.ConvertStatusValue(int8(uomRes.Data.Status)),
					DecimalEnabled: int8(uomRes.Data.DecimalEnabled),
				},
				Class: &dto.ClassResponse{
					ID: "1",
				},
				Description:          itemRes.Data.Description,
				UnitWeightConversion: itemRes.Data.UnitWeightConversion,
				OrderMinQty:          itemRes.Data.OrderMinQty,
				OrderMaxQty:          itemRes.Data.OrderMaxQty,
				ItemType:             itemRes.Data.ItemType,
				Capitalize:           itemRes.Data.Capitalize,
				MaxDayDeliveryDate:   int8(itemRes.Data.MaxDayDeliveryDate),
				Taxable:              itemRes.Data.Taxable,
				Note:                 itemRes.Data.Note,
				Status:               int8(itemRes.Data.Status),
				UnitPrice:            10000,
			},
			PurchaseOrderItem: poi,
			ItemTransferItem:  itItem,
		})
	}

	res = &dto.ReceivingResponse{
		// ID:                  itRes.Data.Id,
		Code: itRes.Data.Code,
		// SiteId:              itRes.Data.SiteId,
		PurchaseOrderId:     itRes.Data.PurchaseOrderId,
		ItemTransferId:      itRes.Data.ItemTransferId,
		InboundType:         int8(itRes.Data.InboundType),
		ValidSupplierReturn: int8(itRes.Data.ValidSupplierReturn),
		AtaDate:             itRes.Data.AtaDate.AsTime(),
		AtaTime:             itRes.Data.AtaTime,
		StockType:           int8(itRes.Data.StockType),
		TotalWeight:         itRes.Data.TotalWeight,
		Note:                itRes.Data.Note,
		Status:              int8(itRes.Data.Status),
		Locked:              int8(itRes.Data.Locked),
		CreatedAt:           itRes.Data.CreatedAt.AsTime(),
		CreatedBy:           itRes.Data.CreatedBy,
		ConfirmedAt:         itRes.Data.ConfirmedAt.AsTime(),
		ConfirmedBy:         itRes.Data.ConfirmedBy,
		LockedBy:            itRes.Data.LockedBy,
		UpdatedAt:           itRes.Data.UpdatedAt.AsTime(),
		UpdatedBy:           itRes.Data.UpdatedBy,
		ReceivingItems:      iti,
		PurchaseOrder:       po,
		ItemTransfer:        it,
		Site:                site,
	}

	return
}

func (s *ReceivingService) Create(ctx context.Context, req *dto.CreateReceivingRequest) (res *dto.ReceivingResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ReceivingService.Create")
	defer span.End()

	var (
		r        bridge_service.CreateReceivingRequest
		items    []*bridge_service.CreateReceivingItemRequest
		resItems []*dto.ReceivingItemResponse
		result   *bridge_service.GetReceivingDetailResponse
	)

	// TODO: check stock opname

	if len(req.Note) > 250 {
		err = edenlabs.ErrorInvalid("note")
		return
	}

	for _, v := range req.ReceivingItem {
		if len(v.Note) > 100 {
			err = edenlabs.ErrorMustEqualOrGreater("note", "100")
			return
		}

		items = append(items, &bridge_service.CreateReceivingItemRequest{
			ItemId: "1",
			Note:   v.Note,
		})
		resItems = append(resItems, &dto.ReceivingItemResponse{
			// ID:   int64(i + 1),
			Note: v.Note,
		})
	}

	r = bridgeService.CreateReceivingRequest{
		Note:          req.Note,
		ReceivingItem: items,
	}
	result, err = s.opt.Client.BridgeServiceGrpc.CreateReceiving(ctx, &r)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, v := range result.Data.ReceivingItems {
		resItems = append(resItems, &dto.ReceivingItemResponse{
			// ID:           v.Id,
			DeliverQty: v.DeliverQty,
			// ReceiveQty:   v.ReceiveQty,
			Weight:       v.Weight,
			Note:         v.Note,
			RejectQty:    v.RejectQty,
			RejectReason: int8(v.RejectReason),
			IsDisabled:   int8(v.IsDisabled),
		})
	}

	res = &dto.ReceivingResponse{
		// ID:                  result.Data.Id,
		Code: result.Data.Code,
		// SiteId:              result.Data.SiteId,
		PurchaseOrderId:     result.Data.PurchaseOrderId,
		ItemTransferId:      result.Data.ItemTransferId,
		InboundType:         int8(result.Data.InboundType),
		ValidSupplierReturn: int8(result.Data.ValidSupplierReturn),
		AtaDate:             result.Data.AtaDate.AsTime(),
		AtaTime:             result.Data.AtaTime,
		StockType:           int8(result.Data.StockType),
		TotalWeight:         result.Data.TotalWeight,
		Note:                result.Data.Note,
		Status:              int8(result.Data.Status),
		Locked:              int8(result.Data.Locked),
		CreatedAt:           result.Data.CreatedAt.AsTime(),
		CreatedBy:           result.Data.CreatedBy,
		ConfirmedAt:         result.Data.ConfirmedAt.AsTime(),
		ConfirmedBy:         result.Data.ConfirmedBy,
		LockedBy:            result.Data.LockedBy,
		UpdatedAt:           result.Data.UpdatedAt.AsTime(),
		UpdatedBy:           result.Data.UpdatedBy,
		ReceivingItems:      resItems,
	}

	jsonRes, _ := json.Marshal(req)
	fmt.Println(string(jsonRes))
	return
}

func (s *ReceivingService) Confirm(ctx context.Context, req *dto.ConfirmReceivingRequest) (res *dto.ReceivingResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ReceivingService.Confirm")
	defer span.End()

	var (
		r        bridge_service.ConfirmReceivingRequest
		resItems []*dto.ReceivingItemResponse
		result   *bridge_service.GetReceivingDetailResponse
	)

	r = bridgeService.ConfirmReceivingRequest{}
	result, err = s.opt.Client.BridgeServiceGrpc.ConfirmReceiving(ctx, &r)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, v := range result.Data.ReceivingItems {
		resItems = append(resItems, &dto.ReceivingItemResponse{
			// ID:           v.Id,
			DeliverQty: v.DeliverQty,
			// ReceiveQty:   v.ReceiveQty,
			Weight:       v.Weight,
			Note:         v.Note,
			RejectQty:    v.RejectQty,
			RejectReason: int8(v.RejectReason),
			IsDisabled:   int8(v.IsDisabled),
		})
	}

	res = &dto.ReceivingResponse{
		// ID:                  result.Data.Id,
		Code: result.Data.Code,
		// SiteId:              result.Data.SiteId,
		PurchaseOrderId:     result.Data.PurchaseOrderId,
		ItemTransferId:      result.Data.ItemTransferId,
		InboundType:         int8(result.Data.InboundType),
		ValidSupplierReturn: int8(result.Data.ValidSupplierReturn),
		AtaDate:             result.Data.AtaDate.AsTime(),
		AtaTime:             result.Data.AtaTime,
		StockType:           int8(result.Data.StockType),
		TotalWeight:         result.Data.TotalWeight,
		Note:                result.Data.Note,
		Status:              int8(result.Data.Status),
		Locked:              int8(result.Data.Locked),
		CreatedAt:           result.Data.CreatedAt.AsTime(),
		CreatedBy:           result.Data.CreatedBy,
		ConfirmedAt:         result.Data.ConfirmedAt.AsTime(),
		ConfirmedBy:         result.Data.ConfirmedBy,
		LockedBy:            result.Data.LockedBy,
		UpdatedAt:           result.Data.UpdatedAt.AsTime(),
		UpdatedBy:           result.Data.UpdatedBy,
		ReceivingItems:      resItems,
	}

	jsonRes, _ := json.Marshal(req)
	fmt.Println(string(jsonRes))
	return
}

func (s *ReceivingService) GetListGP(ctx context.Context, req dto.GetGoodsReceiptGPListRequest) (res []*dto.ReceivingResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ReceivingService.GetListGP")
	defer span.End()

	var (
		po                                                             *bridgeService.GetPurchaseOrderGPResponse
		statusStr                                                      string
		status                                                         int32
		ataDate, ataTime, etaDatePO, etaTimePO, etaDateITT, etaTimeITT time.Time

		gr                *bridgeService.GetGoodsReceiptGPResponse
		site              *bridgeService.GetSiteGPResponse
		itt               *bridgeService.GetInTransitTransferGPResponse
		purchaseOrder     *dto.PurchaseOrderResponse
		inTransitTransfer *dto.ItemTransferResponse
		vendorResponse    *dto.VendorResponse
		vendorDetail      *bridgeService.GetVendorGPResponse
		receivingItems    []*dto.ReceivingItemResponse
	)

	if gr, err = s.opt.Client.BridgeServiceGrpc.GetGoodsReceiptGPList(ctx, &bridgeService.GetGoodsReceiptGPListRequest{
		Limit:    int32(req.Limit),
		Offset:   int32(req.Offset),
		Poprctnm: req.Poprctnm,
		Doctype:  req.Doctype,
	}); err != nil || !gr.Succeeded {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "purchase order")
		return
	}

	/*
		type 1: GR that coming from PO(Shipment)
		type 8: GR that coming from GT(In Transit Transfer)
	*/

	for _, v := range gr.Data {
		switch v.Poptype {
		case 1:
			// get Site
			if site, err = s.opt.Client.BridgeServiceGrpc.GetSiteGPDetail(ctx, &bridgeService.GetSiteGPDetailRequest{
				Id: v.Polist[0].Locncode,
			}); err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "site")
				return
			}

			// get PO
			if po, err = s.opt.Client.BridgeServiceGrpc.GetPurchaseOrderGPList(ctx, &bridgeService.GetPurchaseOrderGPListRequest{
				Ponumber: v.Polist[0].Ponumber,
			}); err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "purchase order")
				return
			}

			// get vendor
			if v.Vendorid == "" {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorMustExistInActive("vendor", "goods receipt/receiving")
				return
			}
			if vendorDetail, err = s.opt.Client.BridgeServiceGrpc.GetVendorGPDetail(ctx, &bridgeService.GetVendorGPDetailRequest{
				Id: v.Vendorid,
			}); err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "vendor")
				return
			}

			vendorResponse = &dto.VendorResponse{
				ID:   vendorDetail.Data[0].VENDORID,
				Name: vendorDetail.Data[0].VENDNAME,
			}

			etaDatePO, err = time.Parse("2006-01-02T15:04:05", po.Data[0].Reqdate)
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorInvalid("eta_date")
				return
			}

			etaTimePO, err = time.Parse("2006-01-02T15:04:05", po.Data[0].Reqdate)
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorInvalid("eta_time")
				return
			}

			recognitionDate, _ := time.Parse("2006-01-02T15:04:05", po.Data[0].Docdate)

			purchaseOrder = &dto.PurchaseOrderResponse{
				ID:              po.Data[0].Ponumber,
				EtaDate:         etaDatePO,
				EtaTime:         etaTimePO.Format("15:04"),
				VendorID:        vendorDetail.Data[0].VENDORID,
				RecognitionDate: recognitionDate,
			}
			for _, v1 := range v.Polist {
				for _, poDet := range po.Data[0].PoDetail {
					if poDet.Itemnmbr == v1.Itemnmbr {
						receivingItems = append(receivingItems, &dto.ReceivingItemResponse{
							ID:          v1.Itemnmbr,
							ItemName:    v1.Itemdesc,
							OrderQty:    poDet.Qtyorder,
							TransferQty: 0,
							RejectQty:   poDet.Qtycance,
							ShippedQty:  v1.Qtyshppd, // should add first by GP for this field
							Note:        v1.Commntid,
							Uom:         v1.Uofm,
						})
					}
				}
			}
		case 8:
			// get Site
			if site, err = s.opt.Client.BridgeServiceGrpc.GetSiteGPDetail(ctx, &bridgeService.GetSiteGPDetailRequest{
				Id: v.Intrxlist[0].Locncode,
			}); err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "site")
				return
			}

			// get PO
			if itt, err = s.opt.Client.BridgeServiceGrpc.GetInTransitTransferGPDetail(ctx, &bridgeService.GetInTransitTransferGPDetailRequest{
				Id: v.Intrxlist[0].Orddocid,
			}); err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "in transit transfer")
				return
			}

			// get vendor
			// if v.Vendorid == "" {
			// 	span.RecordError(err)
			// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
			// 	err = edenlabs.ErrorMustExistInActive("vendor", "goods receipt/receiving")
			// 	return
			// }
			// if vendorDetail, err = s.opt.Client.BridgeServiceGrpc.GetVendorGPDetail(ctx, &bridgeService.GetVendorGPDetailRequest{
			// 	Id: v.Vendorid,
			// }); err != nil {
			// 	span.RecordError(err)
			// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
			// 	err = edenlabs.ErrorRpcNotFound("bridge", "vendor")
			// 	return
			// }
			etaDateITT, err = time.Parse("2006-01-02", itt.Data[0].Etadte)
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorInvalid("eta_date")
				return
			}

			etaTimeITT, err = time.Parse("15:04:05", itt.Data[0].Etatime)
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorInvalid("eta_time")
				return
			}

			recognitionDate, _ := time.Parse("2006-01-02T15:04:05", itt.Data[0].Ordrdate)

			inTransitTransfer = &dto.ItemTransferResponse{
				ID:              itt.Data[0].Orddocid,
				EtaDate:         etaDateITT,
				EtaTime:         etaTimeITT.Format("15:04"),
				VendorID:        "-",
				RecognitionDate: recognitionDate,
			}
			for _, v1 := range v.Intrxlist {
				for _, ittDet := range itt.Data[0].Details {
					if ittDet.Itemnmbr == v1.Itemnmbr {
						receivingItems = append(receivingItems, &dto.ReceivingItemResponse{
							ID:          v1.Itemnmbr,
							ItemName:    v1.Itemdesc,
							OrderQty:    v1.Qtyorder,
							TransferQty: ittDet.Trnsfqty,
							RejectQty:   ittDet.Trnsfqty - v1.Qtyorder,
							ShippedQty:  v1.Qtyshppd, // should add first by GP for this field
							Note:        v1.Commntid,
							Uom:         v1.Uofm,
						})
					}
				}
			}
		}

		ataDate, err = time.Parse("2006-01-02T15:04:05", v.Receiptdate)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("eta_date")
			return
		}

		ataTime, err = time.Parse("15:04:05", v.PrpActualarrivalTime)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("eta_time")
			return
		}
		switch v.Status {
		case "post":
			status = 2
			statusStr = "Finished"
		default:
			status = 1
			statusStr = "New"
		}

		res = append(res, &dto.ReceivingResponse{
			ID:          v.Poprctnm,
			InboundType: int8(v.Poptype),
			// SiteId:         site,
			AtaDate:        ataDate,
			AtaTime:        ataTime.Format("15:04"),
			PurchaseOrder:  purchaseOrder,
			ItemTransfer:   inTransitTransfer,
			ReceivingItems: receivingItems,
			Status:         int8(status),
			StatusStr:      statusStr,
			Region:         v.PrpRegion,
			Note:           "Not available in GP",
			Site: &dto.SiteResponse{
				ID:          site.Data[0].Locncode,
				Name:        site.Data[0].Locndscr,
				Description: site.Data[0].Locndscr,
			},
			Vendor: vendorResponse,
		})
	}

	total = int64(len(res))

	return
}

func (s *ReceivingService) GetDetailGP(ctx context.Context, id string) (res *bridgeService.GoodsReceiptGP, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ReceivingService.GetDetailGP")
	defer span.End()

	var po *bridgeService.GetGoodsReceiptGPResponse

	if po, err = s.opt.Client.BridgeServiceGrpc.GetGoodsReceiptGPDetail(ctx, &bridgeService.GetGoodsReceiptGPDetailRequest{
		Poprctnm: id,
	}); err != nil || !po.Succeeded {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "purchase order")
		return
	}

	if len(po.Data) > 0 {
		res = po.Data[0]
	}

	return
}

func (s *ReceivingService) CreateGP(ctx context.Context, req *dto.CreateGoodsReceiptGPRequest) (res *bridgeService.CreateTransferRequestGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ReceivingService.CreateGP")
	defer span.End()

	var grDetail []*bridgeService.CreateGoodsReceiptGPRequest_GoodsReceiptDetail
	for _, detail := range req.GoodsReceiptDTL {
		grDetail = append(grDetail, &bridgeService.CreateGoodsReceiptGPRequest_GoodsReceiptDetail{
			Poprctnm: detail.Poprctnm,
			Ponumber: detail.Ponumber,
			Locncode: detail.Locncode,
			Uofm:     detail.Uofm,
			Itemnmbr: detail.Itemnmbr,
			Qtyshppd: int32(detail.Qtyshppd),
			Unitcost: detail.Unitcost,
			Extdcost: detail.Extdcost,
		})
	}
	res, err = s.opt.Client.BridgeServiceGrpc.CreateGoodsReceiptGP(ctx, &bridgeService.CreateGoodsReceiptGPRequest{
		Poprctnm:        req.Poprctnm,
		Vnddocnm:        req.Vnddocnm,
		Receiptdate:     req.Receiptdate,
		Vendorid:        req.Vendorid,
		Curncyid:        req.Curncyid,
		Subtotal:        req.Subtotal,
		GoodsReceiptDTL: grDetail,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "purchase order")
		return
	}

	userID := ctx.Value(constants.KeyUserID).(int64)
	_, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &auditService.CreateLogRequest{
		Log: &auditService.Log{
			UserId:      userID,
			ReferenceId: "0",
			Type:        "receiving",
			Function:    "create",
			CreatedAt:   timestamppb.New(time.Now()),
		},
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpc("audit")
		return
	}

	return
}

func (s *ReceivingService) UpdateGP(ctx context.Context, req *dto.UpdateGoodsReceiptGPRequest) (res *bridgeService.CreateTransferRequestGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ReceivingService.UpdateGP")
	defer span.End()

	var grDetail []*bridgeService.UpdateGoodsReceiptGPRequest_Detail
	for _, detail := range req.Details {
		grDetail = append(grDetail, &bridgeService.UpdateGoodsReceiptGPRequest_Detail{
			Rcptlnnm: int32(detail.Rcptlnnm),
			Ponumber: detail.Ponumber,
			Locncode: detail.Locncode,
			Uofm:     detail.Locncode,
			Itemnmbr: detail.Itemnmbr,
			Qtyshppd: int32(detail.Qtyshppd),
			Unitcost: detail.Unitcost,
			Extdcost: detail.Extdcost,
		})
	}
	res, err = s.opt.Client.BridgeServiceGrpc.UpdateGoodsReceiptGP(ctx, &bridgeService.UpdateGoodsReceiptGPRequest{
		Poprctnm:             req.Poprctnm,
		PrpRegion:            req.PrpRegion,
		Actlship:             req.Actlship,
		PrpActualarrivalTime: req.PrpActualarrivalTime,
		Note:                 req.Note,
		Details:              grDetail,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "purchase order")
		return
	}

	userID := ctx.Value(constants.KeyUserID).(int64)
	_, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &auditService.CreateLogRequest{
		Log: &auditService.Log{
			UserId:      userID,
			ReferenceId: "0",
			Type:        "receiving",
			Function:    "update",
			CreatedAt:   timestamppb.New(time.Now()),
		},
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpc("audit")
		return
	}

	return
}
