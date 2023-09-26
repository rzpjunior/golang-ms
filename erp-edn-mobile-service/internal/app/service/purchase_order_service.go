package service

import (
	"context"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/constants"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/global"
	auditService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/audit_service"
	"google.golang.org/protobuf/types/known/timestamppb"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

func NewServicePurchaseOrder() IPurchaseOrderService {
	m := new(PurchaseOrderService)
	m.opt = global.Setup.Common
	return m
}

type IPurchaseOrderService interface {
	Create(ctx context.Context, req dto.CreatePurchaseOrderRequest) (res *dto.PurchaseOrderResponse, err error)
	Get(ctx context.Context, req dto.PurchaseOrderListRequest) (res []*dto.PurchaseOrderResponse, err error)
	GetById(ctx context.Context, req dto.PurchaseOrderDetailRequest) (res *dto.PurchaseOrderResponse, err error)
	Commit(ctx context.Context, id int64) (err error)
	Cancel(ctx context.Context, req dto.CancelPurchaseOrderRequest) (err error)
	Update(ctx context.Context, req dto.UpdatePurchaseOrderRequest) (res *dto.PurchaseOrderResponse, err error)
	UpdateProduct(ctx context.Context, req dto.UpdateProductPurchaseOrderRequest) (res *dto.PurchaseOrderResponse, err error)

	// gp
	CreateGP(ctx context.Context, req dto.CreatePurchaseOrderGPRequest) (res *bridgeService.CreatePurchaseOrderGPResponse, err error)
	UpdateGP(ctx context.Context, req dto.CreatePurchaseOrderGPRequest) (res *bridgeService.CreatePurchaseOrderGPResponse, err error)
	GetListGP(ctx context.Context, req dto.GetPurchaseOrderGPListRequest) (res []*dto.PurchaseOrderResponse, total int64, err error)
	GetDetailGP(ctx context.Context, id string) (res *dto.PurchaseOrderResponse, err error)
	CommitPurchaseOrderGP(ctx context.Context, payload *dto.CommitPurchaseOrderGPRequest) (res *bridgeService.CreateTransferRequestGPResponse, err error)
}

type PurchaseOrderService struct {
	opt opt.Options
}

func (s *PurchaseOrderService) Create(ctx context.Context, req dto.CreatePurchaseOrderRequest) (res *dto.PurchaseOrderResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PurchaseOrderService.Create")
	defer span.End()

	// create Purchase Order from bridge
	var (
		poRes *bridgeService.GetPurchaseOrderDetailResponse
		items []*bridgeService.CreatePurchaseOrderItemRequest
		poi   []*dto.PurchaseOrderItemResponse
	)
	for _, item := range req.PurchaseOrderItems {
		items = append(items, &bridgeService.CreatePurchaseOrderItemRequest{
			ItemId:        item.ItemID,
			OrderQty:      item.OrderQty,
			UnitPrice:     item.UnitPrice,
			Note:          item.Note,
			PurchaseQty:   item.PurchaseQty,
			IncludeTax:    int32(item.IncludeTax),
			TaxPercentage: item.TaxPercentage,
		})
	}
	poRes, err = s.opt.Client.BridgeServiceGrpc.CreatePurchaseOrder(ctx, &bridgeService.CreatePurchaseOrderRequest{
		VendorId:    req.VendorID,
		SiteId:      req.SiteID,
		OrderDate:   req.OrderDate,
		StrEtaDate:  req.StrEtaDate,
		EtaTime:     req.EtaTime,
		DeliveryFee: req.DeliveryFee,
		Note:        req.Note,
		TaxPct:      req.TaxPct,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		Items:       items,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "purchase order")
		return
	}

	for _, item := range poRes.Data.PurchaseOrderItems {
		poi = append(poi, &dto.PurchaseOrderItemResponse{
			// ID:                 item.Id,
			// PurchaseOrderID:    item.PurchaseOrderId,
			PurchasePlanItemID: item.PurchasePlanItemId,
			// ItemID:             item.ItemId,
			OrderQty:      item.OrderQty,
			UnitPrice:     item.UnitPrice,
			TaxableItem:   item.TaxableItem,
			IncludeTax:    item.IncludeTax,
			TaxPercentage: item.TaxPercentage,
			TaxAmount:     item.TaxAmount,
			UnitPriceTax:  item.UnitPriceTax,
			Subtotal:      item.Subtotal,
			Weight:        item.Weight,
			Note:          item.Note,
			PurchaseQty:   item.PurchaseQty,
		})
	}

	res = &dto.PurchaseOrderResponse{
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
		PurchaseOrderItems:     poi,
	}

	return
}

func (s *PurchaseOrderService) Get(ctx context.Context, req dto.PurchaseOrderListRequest) (res []*dto.PurchaseOrderResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PurchaseOrderService.Get")
	defer span.End()

	// get Purchase Orders from bridge
	var poRes *bridgeService.GetPurchaseOrderListResponse
	poRes, err = s.opt.Client.BridgeServiceGrpc.GetPurchaseOrderList(ctx, &bridgeService.GetPurchaseOrderListRequest{
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

	datas := []*dto.PurchaseOrderResponse{}
	for _, po := range poRes.Data {
		var (
			vendor *bridge_service.GetVendorDetailResponse
			site   *bridge_service.GetSiteDetailResponse
			grs    []*dto.ReceivingListinDetailResponse
		)
		vendor, err = s.opt.Client.BridgeServiceGrpc.GetVendorDetail(ctx, &bridgeService.GetVendorDetailRequest{
			Id: po.VendorId,
		})
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		site, err = s.opt.Client.BridgeServiceGrpc.GetSiteDetail(ctx, &bridgeService.GetSiteDetailRequest{
			Id: po.SiteId,
		})
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		for _, gr := range po.Receiving {
			grs = append(grs, &dto.ReceivingListinDetailResponse{
				ID:     gr.Id,
				Code:   gr.Code,
				Status: int8(gr.Status),
			})
		}
		datas = append(datas, &dto.PurchaseOrderResponse{
			// ID:                     po.Id,
			Code: po.Code,
			// VendorID:               po.VendorId,
			SiteID:                 po.SiteId,
			TermPaymentPurID:       po.TermPaymentPurId,
			VendorClassificationID: po.VendorClassificationId,
			PurchasePlanID:         po.PurchasePlanId,
			ConsolidatedShipmentID: po.ConsolidatedShipmentId,
			Status:                 po.Status,
			RecognitionDate:        po.RecognitionDate.AsTime(),
			EtaDate:                po.EtaDate.AsTime(),
			SiteAddress:            po.SiteAddress,
			EtaTime:                po.EtaTime,
			TaxPct:                 po.TaxPct,
			DeliveryFee:            po.DeliveryFee,
			TotalPrice:             po.TotalPrice,
			TaxAmount:              po.TaxAmount,
			TotalCharge:            po.TotalCharge,
			TotalInvoice:           po.TotalInvoice,
			TotalWeight:            po.TotalWeight,
			Note:                   po.Note,
			DeltaPrint:             po.DeltaPrint,
			Latitude:               po.Latitude,
			Longitude:              po.Longitude,
			CreatedFrom:            po.CreatedFrom,
			HasFinishedGr:          int8(po.HasFinishedGr),
			CommittedAt:            po.CommittedAt.AsTime(),
			CommittedBy:            po.CommittedBy,
			AssignedTo:             po.AssignedTo,
			AssignedBy:             po.AssignedBy,
			AssignedAt:             po.AssignedAt.AsTime(),
			Locked:                 po.Locked,
			LockedBy:               po.LockedBy,
			CreatedBy:              po.CreatedBy,
			CreatedAt:              po.CreatedAt.AsTime(),
			UpdatedAt:              po.UpdatedAt.AsTime(),
			UpdatedBy:              po.UpdatedBy,
			Receiving:              grs,
			Vendor: &dto.VendorResponse{
				// ID:                     vendor.Data.Id,
				Code: vendor.Data.Code,
				// VendorOrganizationID:   vendor.Data.VendorOrganizationId,
				VendorClassificationID: vendor.Data.VendorClassificationId,
				SubDistrictID:          vendor.Data.SubDistrictId,
				PaymentTermID:          vendor.Data.PaymentTermId,
				PicName:                vendor.Data.PicName,
				Email:                  vendor.Data.Email,
				PhoneNumber:            vendor.Data.PhoneNumber,
				Rejectable:             vendor.Data.Rejectable,
				Returnable:             vendor.Data.Rejectable,
				Address:                vendor.Data.Address,
				Note:                   vendor.Data.Note,
				Status:                 vendor.Data.Status,
				Latitude:               vendor.Data.Latitude,
				Longitude:              vendor.Data.Longitude,
				CreatedAt:              vendor.Data.CreatedAt.AsTime(),
			},
			Site: &dto.SiteResponse{
				// ID:            site.Data.Id,
				Code:          site.Data.Code,
				Description:   site.Data.Description,
				Name:          site.Data.Description,
				Status:        int8(site.Data.Status),
				StatusConvert: statusx.ConvertStatusValue(int8(site.Data.Status)),
				Region: &dto.RegionResponse{
					// ID:            site.Data.Region.Id,
					Code:          site.Data.Region.Code,
					Description:   site.Data.Region.Description,
					Status:        int8(site.Data.Region.Status),
					StatusConvert: statusx.ConvertStatusValue(int8(site.Data.Region.Status)),
					CreatedAt:     site.Data.Region.CreatedAt.AsTime(),
					UpdatedAt:     site.Data.Region.UpdatedAt.AsTime(),
				},
				CreatedAt: site.Data.CreatedAt.AsTime(),
				UpdatedAt: site.Data.UpdatedAt.AsTime(),
			},
		})
	}
	res = datas

	return
}

func (s *PurchaseOrderService) GetById(ctx context.Context, req dto.PurchaseOrderDetailRequest) (res *dto.PurchaseOrderResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PurchaseOrderService.GetById")
	defer span.End()

	// get Purchase Order from bridge
	var (
		poRes     *bridgeService.GetPurchaseOrderDetailResponse
		poi       []*dto.PurchaseOrderItemResponse
		grs       []*dto.ReceivingListinDetailResponse
		unitPrice = float64(10000)
	)
	poRes, err = s.opt.Client.BridgeServiceGrpc.GetPurchaseOrderDetail(ctx, &bridgeService.GetPurchaseOrderDetailRequest{
		Id: int64(req.Id),
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "purchase order")
		return
	}

	for _, item := range poRes.Data.PurchaseOrderItems {
		poi = append(poi, &dto.PurchaseOrderItemResponse{
			// ID:                 item.Id,
			// PurchaseOrderID:    item.PurchaseOrderId,
			PurchasePlanItemID: item.PurchasePlanItemId,
			// ItemID:             item.ItemId,
			OrderQty:      item.OrderQty,
			UnitPrice:     item.UnitPrice,
			TaxableItem:   item.TaxableItem,
			IncludeTax:    item.IncludeTax,
			TaxPercentage: item.TaxPercentage,
			TaxAmount:     item.TaxAmount,
			UnitPriceTax:  item.UnitPriceTax,
			Subtotal:      item.Subtotal,
			Weight:        item.Weight,
			Note:          item.Note,
			PurchaseQty:   item.PurchaseQty,
			Item: &dto.ItemResponse{
				// ID:   item.Item.Id,
				Code: item.Item.Code,
				Uom: &dto.UomResponse{
					ID:          "1",
					Code:        item.Item.Uom.Code,
					Description: item.Item.Uom.Description,
				},
				Description:          item.Item.Description,
				UnitWeightConversion: item.Item.UnitWeightConversion,
				OrderMinQty:          item.Item.OrderMinQty,
				OrderMaxQty:          item.Item.OrderMaxQty,
				ItemType:             item.Item.ItemType,
				Capitalize:           item.Item.Capitalize,
				MaxDayDeliveryDate:   int8(item.Item.MaxDayDeliveryDate),
				Taxable:              item.Item.Taxable,
				Note:                 item.Item.Note,
				Status:               int8(item.Item.Status),
				UnitPrice:            unitPrice,
			},
		})
	}

	for _, gr := range poRes.Data.Receiving {
		grs = append(grs, &dto.ReceivingListinDetailResponse{
			ID:     gr.Id,
			Code:   gr.Code,
			Status: int8(gr.Status),
		})
	}

	res = &dto.PurchaseOrderResponse{
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
		PurchaseOrderItems:     poi,
		Receiving:              grs,
	}

	return
}

func (s *PurchaseOrderService) Commit(ctx context.Context, id int64) (err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PurchaseOrderService.GetById")
	defer span.End()

	// get Purchase Order from bridge
	var poRes *bridgeService.GetPurchaseOrderDetailResponse
	poRes, err = s.opt.Client.BridgeServiceGrpc.GetPurchaseOrderDetail(ctx, &bridgeService.GetPurchaseOrderDetailRequest{
		Id: id,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "purchase order")
		return
	}

	if poRes.Data.Status != 5 {
		err = edenlabs.ErrorMustActive("status")
		return
	}

	// commit Purchase Order from bridge
	_, err = s.opt.Client.BridgeServiceGrpc.CommitPurchaseOrder(ctx, &bridgeService.CommitPurchaseOrderRequest{
		Id: id,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "purchase order")
		return
	}

	return
}

func (s *PurchaseOrderService) Cancel(ctx context.Context, req dto.CancelPurchaseOrderRequest) (err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PurchaseOrderService.GetById")
	defer span.End()

	// get Purchase Order from bridge
	var poRes *bridgeService.GetPurchaseOrderDetailResponse
	poRes, err = s.opt.Client.BridgeServiceGrpc.GetPurchaseOrderDetail(ctx, &bridgeService.GetPurchaseOrderDetailRequest{
		Id: req.Id,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "purchase order")
		return
	}

	if poRes.Data.Status != 5 {
		err = edenlabs.ErrorMustActive("status")
		return
	}

	if poRes.Data.PurchasePlanId != 0 {
		_, err = s.opt.Client.BridgeServiceGrpc.GetPurchasePlanDetail(ctx, &bridgeService.GetPurchasePlanDetailRequest{
			Id: req.Id,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "purchase plan")
			return
		}
	}

	if poRes.Data.DeltaPrint > 0 {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("delta_print")
		return
	}

	// cancel Purchase Order from bridge
	_, err = s.opt.Client.BridgeServiceGrpc.CancelPurchaseOrder(ctx, &bridgeService.CancelPurchaseOrderRequest{
		Id:   req.Id,
		Note: req.Note,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "cancel purchase order")
		return
	}

	return
}

func (s *PurchaseOrderService) Update(ctx context.Context, req dto.UpdatePurchaseOrderRequest) (res *dto.PurchaseOrderResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PurchaseOrderService.Update")
	defer span.End()

	// update Purchase Order from bridge
	var (
		poRes *bridgeService.GetPurchaseOrderDetailResponse
		items []*bridgeService.UpdatePurchaseOrderItemRequest
		poi   []*dto.PurchaseOrderItemResponse
	)
	for _, item := range req.PurchaseOrderItems {
		items = append(items, &bridgeService.UpdatePurchaseOrderItemRequest{
			ItemId:        item.ItemID,
			OrderQty:      item.OrderQty,
			UnitPrice:     item.UnitPrice,
			Note:          item.Note,
			PurchaseQty:   item.PurchaseQty,
			IncludeTax:    int32(item.IncludeTax),
			TaxPercentage: item.TaxPercentage,
		})
	}
	poRes, err = s.opt.Client.BridgeServiceGrpc.UpdatePurchaseOrder(ctx, &bridgeService.UpdatePurchaseOrderRequest{
		VendorId:    req.VendorID,
		SiteId:      req.SiteID,
		OrderDate:   req.OrderDate,
		StrEtaDate:  req.StrEtaDate,
		EtaTime:     req.EtaTime,
		DeliveryFee: req.DeliveryFee,
		Note:        req.Note,
		TaxPct:      req.TaxPct,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		Items:       items,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "purchase order")
		return
	}

	for _, item := range poRes.Data.PurchaseOrderItems {
		poi = append(poi, &dto.PurchaseOrderItemResponse{
			// ID:                 item.Id,
			// PurchaseOrderID:    item.PurchaseOrderId,
			PurchasePlanItemID: item.PurchasePlanItemId,
			// ItemID:             item.ItemId,
			OrderQty:      item.OrderQty,
			UnitPrice:     item.UnitPrice,
			TaxableItem:   item.TaxableItem,
			IncludeTax:    item.IncludeTax,
			TaxPercentage: item.TaxPercentage,
			TaxAmount:     item.TaxAmount,
			UnitPriceTax:  item.UnitPriceTax,
			Subtotal:      item.Subtotal,
			Weight:        item.Weight,
			Note:          item.Note,
			PurchaseQty:   item.PurchaseQty,
		})
	}

	res = &dto.PurchaseOrderResponse{
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
		PurchaseOrderItems:     poi,
	}

	return
}

func (s *PurchaseOrderService) UpdateProduct(ctx context.Context, req dto.UpdateProductPurchaseOrderRequest) (res *dto.PurchaseOrderResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PurchaseOrderService.UpdateProduct")
	defer span.End()

	// update Purchase Order from bridge
	var (
		poRes *bridgeService.GetPurchaseOrderDetailResponse
		items []*bridgeService.UpdatePurchaseOrderItemRequest
		poi   []*dto.PurchaseOrderItemResponse
	)
	for _, item := range req.PurchaseOrderItems {
		items = append(items, &bridgeService.UpdatePurchaseOrderItemRequest{
			ItemId:        item.ItemID,
			OrderQty:      item.OrderQty,
			UnitPrice:     item.UnitPrice,
			Note:          item.Note,
			PurchaseQty:   item.PurchaseQty,
			IncludeTax:    int32(item.IncludeTax),
			TaxPercentage: item.TaxPercentage,
		})
	}
	poRes, err = s.opt.Client.BridgeServiceGrpc.UpdateProductPurchaseOrder(ctx, &bridgeService.UpdateProductPurchaseOrderRequest{
		DeliveryFee: req.DeliveryFee,
		TaxPct:      req.TaxPct,
		Items:       items,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "purchase order")
		return
	}

	for _, item := range poRes.Data.PurchaseOrderItems {
		poi = append(poi, &dto.PurchaseOrderItemResponse{
			// ID:                 item.Id,
			// PurchaseOrderID:    item.PurchaseOrderId,
			PurchasePlanItemID: item.PurchasePlanItemId,
			// ItemID:             item.ItemId,
			OrderQty:      item.OrderQty,
			UnitPrice:     item.UnitPrice,
			TaxableItem:   item.TaxableItem,
			IncludeTax:    item.IncludeTax,
			TaxPercentage: item.TaxPercentage,
			TaxAmount:     item.TaxAmount,
			UnitPriceTax:  item.UnitPriceTax,
			Subtotal:      item.Subtotal,
			Weight:        item.Weight,
			Note:          item.Note,
			PurchaseQty:   item.PurchaseQty,
		})
	}

	res = &dto.PurchaseOrderResponse{
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
		PurchaseOrderItems:     poi,
	}

	return
}

func (s *PurchaseOrderService) CreateGP(ctx context.Context, req dto.CreatePurchaseOrderGPRequest) (res *bridgeService.CreatePurchaseOrderGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PurchaseOrderService.CreateGP")
	defer span.End()

	// create Purchase Order from bridge
	var items []*bridgeService.CreatePurchaseOrderGPRequest_Podtl
	itemList := make(map[string]bool)
	for _, item := range req.Detail {
		if _, exist := itemList[item.Itemnmbr]; exist {
			err = edenlabs.ErrorDuplicate("itmnmbr")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		items = append(items, &bridgeService.CreatePurchaseOrderGPRequest_Podtl{
			Ord:      item.Ord,
			Itemnmbr: item.Itemnmbr,
			Uofm:     item.Uofm,
			Qtyorder: item.Qtyorder,
			Qtycance: item.Qtycance,
			Unitcost: item.Unitcost,
			Notetext: item.Notetext,
		})
		itemList[item.Itemnmbr] = true
	}
	res, err = s.opt.Client.BridgeServiceGrpc.CreatePurchaseOrderGP(ctx, &bridgeService.CreatePurchaseOrderGPRequest{
		Potype:                  req.Potype,
		Ponumber:                req.Ponumber,
		Docdate:                 req.Docdate,
		Buyerid:                 req.Buyerid,
		Vendorid:                req.Vendorid,
		Curncyid:                req.Curncyid,
		Deprtmnt:                req.Deprtmnt,
		Locncode:                req.Locncode,
		Taxschid:                req.Taxschid,
		Subtotal:                req.Subtotal,
		Trdisamt:                req.Trdisamt,
		Frtamnt:                 req.Frtamnt,
		Miscamnt:                req.Miscamnt,
		Taxamnt:                 req.Taxamnt,
		PrpPurchaseplanNo:       req.PrpPurchasePlanNo,
		PrpPaymentMethod:        req.PrpPaymentMethod,
		PrpRegion:               req.PrpRegion,
		PrpEstimatedarrivalDate: req.PrpEstimatedArrivalDate,
		Notetext:                req.NoteText,
		Duedate:                 "1900-01-01",
		Pymtrmid:                "",
		Detail:                  items,
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
			ReferenceId: res.Ponumber,
			Type:        "purchase_order",
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

func (s *PurchaseOrderService) UpdateGP(ctx context.Context, req dto.CreatePurchaseOrderGPRequest) (res *bridgeService.CreatePurchaseOrderGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PurchaseOrderService.UpdateGP")
	defer span.End()

	// update Purchase Order from bridge
	var items []*bridgeService.UpdatePurchaseOrderGPRequest_Podtl
	itemList := make(map[string]bool)
	for _, item := range req.Detail {
		if _, exist := itemList[item.Itemnmbr]; exist {
			err = edenlabs.ErrorDuplicate("itmnmbr")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		items = append(items, &bridgeService.UpdatePurchaseOrderGPRequest_Podtl{
			Ord:      item.Ord,
			Itemnmbr: item.Itemnmbr,
			Uofm:     item.Uofm,
			Qtyorder: item.Qtyorder,
			Qtycance: item.Qtycance,
			Unitcost: item.Unitcost,
			Notetext: item.Notetext,
		})
		itemList[item.Itemnmbr] = true
	}
	res, err = s.opt.Client.BridgeServiceGrpc.UpdatePurchaseOrderGP(ctx, &bridgeService.UpdatePurchaseOrderGPRequest{
		Potype:                  req.Potype,
		Ponumber:                req.Ponumber,
		Docdate:                 req.Docdate,
		Buyerid:                 req.Buyerid,
		Vendorid:                req.Vendorid,
		Curncyid:                req.Curncyid,
		Deprtmnt:                req.Deprtmnt,
		Locncode:                req.Locncode,
		Taxschid:                req.Taxschid,
		Subtotal:                req.Subtotal,
		Trdisamt:                req.Trdisamt,
		Frtamnt:                 req.Frtamnt,
		Miscamnt:                req.Miscamnt,
		Taxamnt:                 req.Taxamnt,
		PrpPurchaseplanNo:       req.PrpPurchasePlanNo,
		PrpPaymentMethod:        req.PrpPaymentMethod,
		PrpRegion:               req.PrpRegion,
		PrpEstimatedarrivalDate: req.PrpEstimatedArrivalDate,
		Notetext:                req.NoteText,
		Detail:                  items,
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
			ReferenceId: res.Ponumber,
			Type:        "purchase_order",
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

func (s *PurchaseOrderService) GetListGP(ctx context.Context, req dto.GetPurchaseOrderGPListRequest) (res []*dto.PurchaseOrderResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PurchaseOrderService.GetListGP")
	defer span.End()

	var (
		po                                *bridgeService.GetPurchaseOrderGPResponse
		vendorDetail                      *bridgeService.GetVendorGPResponse
		siteDetail                        *bridgeService.GetSiteGPResponse
		item                              *bridgeService.GetItemGPResponse
		pi                                *dto.PurchaseInvoiceDetailResponse
		receiving                         []*dto.ReceivingListinDetailResponse
		recognitionDate, etaDate, etaTime time.Time
		totalWeight                       float64
		layoutInput, statusStr            string
		statusFilter                      []int32
		status                            int32
		// layoutOutput    string
	)

	layoutInput = "2006-01-02T15:04:05" // The layout for the input string
	// layoutOutput = "2006-01-02"         // The layout for the desired output
	switch req.Postatus {
	case 1:
		statusFilter = []int32{1}
	case 4:
		statusFilter = []int32{2, 3, 4}
	default:
		statusFilter = []int32{1, 2, 3, 4}
	}

	// get PO
	if po, err = s.opt.Client.BridgeServiceGrpc.GetPurchaseOrderGPList(ctx, &bridgeService.GetPurchaseOrderGPListRequest{
		Limit:        int32(req.Limit),
		Offset:       int32(req.Offset),
		Orderby:      req.OrderBy,
		Ponumber:     req.Ponumber,
		Ponumberlike: req.PonumberLike,
		ReqdateFrom:  req.ReqDateFrom.Format("2006-01-02"),
		ReqdateTo:    req.ReqDateTo.Format("2006-01-02"),
		VendorId:     req.Vendorid,
		Locncode:     req.Locncode,
		Status:       statusFilter,
	}); err != nil || !po.Succeeded {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "purchase order")
		return
	}

	for _, v := range po.Data {
		/**
		it should be deleted later on, because the data showed
		if there is data dummy from gp, will not match with amount of perpage
		**/

		if v.Postatus == 0 || v.Vendorid == "" {
			continue
		}
		recognitionDate, _ = time.Parse(layoutInput, v.Docdate)

		// get vendor
		if vendorDetail, err = s.opt.Client.BridgeServiceGrpc.GetVendorGPDetail(ctx, &bridgeService.GetVendorGPDetailRequest{
			Id: v.Vendorid,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "vendor")
			return
		}

		// get site
		if siteDetail, err = s.opt.Client.BridgeServiceGrpc.GetSiteGPDetail(ctx, &bridgeService.GetSiteGPDetailRequest{
			Id: v.PrpLocncode,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "site")
			return
		}

		etaDate, err = time.Parse("2006-01-02T15:04:05", po.Data[0].Reqdate)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("eta_date")
			return
		}

		etaTime, err = time.Parse("2006-01-02T15:04:05", po.Data[0].Reqdate)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("eta_time")
			return
		}

		for _, detail := range v.PoDetail {
			// handling uofm non KG, it has to be converted to KG
			if detail.Uofm != "KG" {
				// get item
				if item, err = s.opt.Client.BridgeServiceGrpc.GetItemGPList(ctx, &bridgeService.GetItemGPListRequest{
					ItemNumber: detail.Itemnmbr,
				}); err != nil {
					span.RecordError(err)
					s.opt.Logger.AddMessage(log.ErrorLevel, err)
					err = edenlabs.ErrorRpcNotFound("bridge", "item")
					return
				}
				totalWeight += item.Data[0].Itemshwt * detail.Qtyorder
			} else {
				totalWeight += detail.Qtyorder
			}
		}

		if len(v.Poinvoice) > 0 {
			pi = &dto.PurchaseInvoiceDetailResponse{
				ID:        v.Poinvoice[0].Poprctnm,
				Status:    int8(v.Poinvoice[0].Status),
				StatusStr: v.Poinvoice[0].StatusDesc,
			}
		}

		if len(v.Poreceiving) > 0 {
			for _, v1 := range v.Poreceiving {
				receiving = append(receiving, &dto.ReceivingListinDetailResponse{
					ID:        v1.Poprctnm,
					Status:    int8(v1.Status),
					StatusStr: v1.StatusDesc,
				})
			}
		}
		switch v.Postatus {
		case 1:
			status = 1
			statusStr = "Draft"
		case 2, 3, 4:
			status = 4
			statusStr = "Active"
		}
		res = append(res, &dto.PurchaseOrderResponse{
			ID:              v.Ponumber,
			RecognitionDate: recognitionDate,
			Vendor: &dto.VendorResponse{
				ID:   vendorDetail.Data[0].VENDORID,
				Name: vendorDetail.Data[0].VENDNAME,
			},
			Site: &dto.SiteResponse{
				ID:          siteDetail.Data[0].Locncode,
				Name:        siteDetail.Data[0].Locndscr,
				Description: siteDetail.Data[0].Locndscr,
				Address:     siteDetail.Data[0].AddresS1 + siteDetail.Data[0].AddresS2 + siteDetail.Data[0].AddresS3,
			},
			EtaDate:         etaDate,
			EtaTime:         etaTime.Format("15:04"),
			TotalWeight:     totalWeight,
			TotalPrice:      v.Subtotal,
			TotalCharge:     v.Total,
			StatusGP:        v.Postatus,
			Status:          status,
			StatusStr:       statusStr,
			PurchaseInvoice: pi,
			Receiving:       receiving,
		})
	}

	total = int64(len(res))

	return
}

func (s *PurchaseOrderService) GetDetailGP(ctx context.Context, id string) (res *dto.PurchaseOrderResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PurchaseOrderService.GetDetailGP")
	defer span.End()

	var (
		po                                *bridgeService.GetPurchaseOrderGPResponse
		vendorDetail                      *bridgeService.GetVendorGPResponse
		siteDetail                        *bridgeService.GetSiteGPResponse
		purchaseOrderItems                []*dto.PurchaseOrderItemResponse
		admDivision                       *bridgeService.GetAdmDivisionGPResponse
		paymentTerm                       *bridgeService.GetPaymentTermGPResponse
		item                              *bridgeService.GetItemGPResponse
		pymntTermRes                      *dto.PaymentTermResponse
		pi                                *dto.PurchaseInvoiceDetailResponse
		receiving                         []*dto.ReceivingListinDetailResponse
		status                            int32
		recognitionDate, etaDate, etaTime time.Time
		totalWeight                       float64
		layoutInput, statusStr            string
		// layoutOutput    string
	)
	receiving = []*dto.ReceivingListinDetailResponse{}

	layoutInput = "2006-01-02T15:04:05" // The layout for the input string
	// layoutOutput = "2006-01-02"         // The layout for the desired output

	if po, err = s.opt.Client.BridgeServiceGrpc.GetPurchaseOrderGPDetail(ctx, &bridgeService.GetPurchaseOrderGPDetailRequest{
		Id: id,
	}); err != nil || !po.Succeeded {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "purchase order")
		return
	}

	if len(po.Data) > 0 {

		recognitionDate, _ = time.Parse(layoutInput, po.Data[0].Docdate)

		// get vendor
		if vendorDetail, err = s.opt.Client.BridgeServiceGrpc.GetVendorGPDetail(ctx, &bridgeService.GetVendorGPDetailRequest{
			Id: po.Data[0].Vendorid,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "vendor")
			return
		}

		// get site
		if siteDetail, err = s.opt.Client.BridgeServiceGrpc.GetSiteGPDetail(ctx, &bridgeService.GetSiteGPDetailRequest{
			Id: po.Data[0].PrpLocncode,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "site")
			return
		}
		// get region of the adm division
		if admDivision, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPList(ctx, &bridgeService.GetAdmDivisionGPListRequest{
			Limit:           1000,
			AdmDivisionCode: siteDetail.Data[0].GnlAdministrativeCode,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "adm division")
			return
		}

		if po.Data[0].Pymtrmid != "" {
			paymentTerm, err = s.opt.Client.BridgeServiceGrpc.GetPaymentTermGPDetail(ctx, &bridge_service.GetPaymentTermGPDetailRequest{
				Id: po.Data[0].Pymtrmid,
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "payment_method")
				return
			}
			pymntTermRes = &dto.PaymentTermResponse{
				ID: paymentTerm.Data[0].Pymtrmid,
			}
		}

		etaDate, err = time.Parse("2006-01-02T15:04:05", po.Data[0].Reqdate)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("eta_date")
			return
		}

		etaTime, err = time.Parse("2006-01-02T15:04:05", po.Data[0].Reqdate)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("eta_time")
			return
		}
		if len(po.Data[0].Poinvoice) > 0 {
			switch po.Data[0].Poinvoice[0].Status {
			case 1:
				status = 1
				statusStr = "Draft"
			case 2, 3, 4:
				status = 4
				statusStr = "Active"
			default:
				status = 4
				statusStr = "Active"
			}
			pi = &dto.PurchaseInvoiceDetailResponse{
				ID:        po.Data[0].Poinvoice[0].Poprctnm,
				Status:    int8(status),
				StatusStr: statusStr,
			}
		}

		if len(po.Data[0].Poreceiving) > 0 {
			for _, v := range po.Data[0].Poreceiving {
				switch v.Status {
				case 1:
					status = 1
					statusStr = "Draft"
				case 2, 3, 4:
					status = 4
					statusStr = "Active"
				default:
					status = 4
					statusStr = "Active"
				}
				receiving = append(receiving, &dto.ReceivingListinDetailResponse{
					ID:        v.Poprctnm,
					Status:    int8(status),
					StatusStr: statusStr,
				})
			}
		}

		for _, detail := range po.Data[0].PoDetail {
			// handling uofm non KG, it has to be converted to KG
			if detail.Uofm != "KG" {
				// get item
				if item, err = s.opt.Client.BridgeServiceGrpc.GetItemGPList(ctx, &bridgeService.GetItemGPListRequest{
					ItemNumber: detail.Itemnmbr,
				}); err != nil {
					span.RecordError(err)
					s.opt.Logger.AddMessage(log.ErrorLevel, err)
					err = edenlabs.ErrorRpcNotFound("bridge", "item")
					return
				}
				totalWeight += item.Data[0].Itemshwt * detail.Qtyorder
			} else {
				totalWeight += detail.Qtyorder
			}
			purchaseOrderItems = append(purchaseOrderItems, &dto.PurchaseOrderItemResponse{
				ItemID:      detail.Itemnmbr,
				ItemName:    detail.Itemdesc,
				PurchaseQty: detail.Qtyorder - detail.Qtycance,
				OrderQty:    detail.Qtyorder,
				UnitPrice:   detail.Unitcost,
				Uom:         detail.Uofm,
				// IncludeTax: ,
				// TaxPercentage: detail.,
				TaxAmount: detail.Taxamnt,
				Note:      detail.Commntid,
				// UnitPriceTax: detail.U,
				Subtotal: detail.Qtyorder * detail.Unitcost,
			})
			// totalPrice += detail.Extdcost
			// TotalCharge += detail.Extdcost // it should be calculate with another addition tax amount
		}
		switch po.Data[0].Postatus {
		case 1:
			status = 1
			statusStr = "Draft"
		case 2, 3, 4:
			status = 4
			statusStr = "Active"
		}
		res = &dto.PurchaseOrderResponse{
			ID:              po.Data[0].Ponumber,
			RecognitionDate: recognitionDate,
			Vendor: &dto.VendorResponse{
				ID:   vendorDetail.Data[0].VENDORID,
				Name: vendorDetail.Data[0].VENDNAME,
				VendorOrganization: &dto.VendorOrganizationResponse{
					ID:   vendorDetail.Data[0].Organization.PRP_Vendor_Org_ID,
					Name: vendorDetail.Data[0].Organization.PRP_Vendor_Org_Desc,
				},
			},
			PaymentTerm: pymntTermRes,
			Site: &dto.SiteResponse{
				ID:          siteDetail.Data[0].Locncode,
				Name:        siteDetail.Data[0].Locndscr,
				Description: siteDetail.Data[0].Locndscr,
				Address:     siteDetail.Data[0].AddresS1 + siteDetail.Data[0].AddresS2 + siteDetail.Data[0].AddresS3,
				Region: &dto.RegionResponse{
					Name: admDivision.Data[0].Region,
				},
			},
			EtaDate:            etaDate,
			EtaTime:            etaTime.Format("15:04"),
			TotalWeight:        totalWeight,
			TotalPrice:         po.Data[0].Subtotal,
			TotalCharge:        po.Data[0].Total,
			StatusGP:           po.Data[0].Postatus,
			Status:             status,
			StatusStr:          statusStr,
			PurchaseOrderItems: purchaseOrderItems,
			Receiving:          receiving,
			Note:               po.Data[0].Commntid,
			PurchaseInvoice:    pi,
		}
	}

	return
}

func (s *PurchaseOrderService) CommitPurchaseOrderGP(ctx context.Context, payload *dto.CommitPurchaseOrderGPRequest) (res *bridgeService.CreateTransferRequestGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PurchaseOrderService.CreateDetailGP")
	defer span.End()

	if res, err = s.opt.Client.BridgeServiceGrpc.CommitPurchaseOrderGP(ctx, &bridgeService.CommitPurchaseOrderGPRequest{
		Docnumber: payload.Docnumber,
		Docdate:   payload.Docdate,
	}); err != nil {
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
			Type:        "purchase_order",
			Function:    "commit",
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
