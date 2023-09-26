package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/edenlabs/edenlabs/timex"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/repository"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"github.com/sirupsen/logrus"
)

type IReceivingService interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (res []dto.ReceivingResponse, total int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (res dto.ReceivingResponse, err error)
	Create(ctx context.Context, req *dto.CreateReceivingRequest) (res dto.ReceivingResponse, err error)
	Confirm(ctx context.Context, req *dto.ConfirmReceivingRequest) (res dto.ReceivingResponse, err error)

	// gp integrated
	GetGP(ctx context.Context, req *pb.GetGoodsReceiptGPListRequest) (res *pb.GetGoodsReceiptGPResponse, err error)
	GetDetailGP(ctx context.Context, req *pb.GetGoodsReceiptGPDetailRequest) (res *pb.GetGoodsReceiptGPResponse, err error)
	CreateGP(ctx context.Context, req *pb.CreateGoodsReceiptGPRequest) (res *pb.CreateTransferRequestGPResponse, err error)
	UpdateGP(ctx context.Context, req *pb.UpdateGoodsReceiptGPRequest) (res *pb.CreateTransferRequestGPResponse, err error)
}

type ReceivingService struct {
	opt                     opt.Options
	RepositoryReceiving     repository.IReceivingRepository
	RepositoryReceivingItem repository.IReceivingItemRepository
	RepositoryPurchaseOrder repository.IPurchaseOrderRepository
	RepositoryItemTransfer  repository.IItemTransferRepository
	RepositorySite          repository.ISiteRepository
	RepositoryItem          repository.IItemRepository
	RepositoryUom           repository.IUomRepository
}

func NewReceivingService() IReceivingService {
	return &ReceivingService{
		opt:                     global.Setup.Common,
		RepositoryReceiving:     repository.NewReceivingRepository(),
		RepositoryReceivingItem: repository.NewReceivingItemRepository(),
		RepositoryPurchaseOrder: repository.NewPurchaseOrderRepository(),
		RepositoryItemTransfer:  repository.NewItemTransferRepository(),
		RepositorySite:          repository.NewSiteRepository(),
		RepositoryItem:          repository.NewItemRepository(),
		RepositoryUom:           repository.NewUomRepository(),
	}
}

func (s *ReceivingService) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (res []dto.ReceivingResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ReceivingService.Get")
	defer span.End()

	var (
		gr           []*model.Receiving
		site         *model.Site
		po           *model.PurchaseOrder
		itemTransfer *model.ItemTransfer
	)
	gr, total, err = s.RepositoryReceiving.Get(ctx, offset, limit, status, search, orderBy)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, it := range gr {
		site, err = s.RepositorySite.GetDetail(ctx, it.SiteId, "")
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		po, err = s.RepositoryPurchaseOrder.GetDetail(ctx, it.PurchaseOrderId, "")
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		itemTransfer, err = s.RepositoryItemTransfer.GetDetail(ctx, it.ItemTransferId, "")
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		res = append(res, dto.ReceivingResponse{
			ID:                  it.ID,
			SiteId:              it.SiteId,
			PurchaseOrderId:     it.PurchaseOrderId,
			ItemTransferId:      it.ItemTransferId,
			InboundType:         it.InboundType,
			ValidSupplierReturn: it.ValidSupplierReturn,
			CreatedAt:           it.CreatedAt,
			CreatedBy:           it.CreatedBy,
			ConfirmedAt:         it.ConfirmedAt,
			ConfirmedBy:         it.ConfirmedBy,
			Code:                it.Code,
			AtaDate:             it.AtaDate,
			AtaTime:             it.AtaTime,
			StockType:           it.StockType,
			TotalWeight:         it.TotalWeight,
			Note:                it.Note,
			Status:              it.Status,
			Locked:              it.Locked,
			LockedBy:            it.LockedBy,
			UpdatedAt:           it.UpdatedAt,
			UpdatedBy:           it.UpdatedBy,
			PurchaseOrder: &dto.PurchaseOrderResponse{
				ID:                     po.ID,
				Code:                   po.Code,
				VendorID:               po.VendorID,
				SiteID:                 po.SiteID,
				TermPaymentPurID:       po.TermPaymentPurID,
				VendorClassificationID: po.VendorClassificationID,
				PurchasePlanID:         po.PurchasePlanID,
				ConsolidatedShipmentID: po.ConsolidatedShipmentID,
				Status:                 po.Status,
				RecognitionDate:        po.RecognitionDate,
				EtaDate:                po.EtaDate,
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
				UpdatedAt:              po.UpdatedAt,
				UpdatedBy:              po.UpdatedBy,
				CreatedFrom:            po.CreatedFrom,
				HasFinishedGr:          po.HasFinishedGr,
				CreatedAt:              po.CreatedAt,
				CreatedBy:              po.CreatedBy,
				CommittedAt:            po.CommittedAt,
				CommittedBy:            po.CommittedBy,
				AssignedTo:             po.AssignedTo,
				AssignedBy:             po.AssignedBy,
				AssignedAt:             po.AssignedAt,
				Locked:                 po.Locked,
				LockedBy:               po.LockedBy,
			},
			Site: &dto.SiteResponse{
				ID:          site.ID,
				Code:        site.Code,
				Description: site.Description,
				Status:      site.Status,
				CreatedAt:   site.CreatedAt,
				UpdatedAt:   site.UpdatedAt,
			},
			ItemTransfer: &dto.ItemTransferResponse{
				ID:                 itemTransfer.ID,
				Code:               itemTransfer.Code,
				RecognitionDate:    itemTransfer.RecognitionDate,
				RequestDate:        itemTransfer.RequestDate,
				EtaDate:            itemTransfer.EtaDate,
				EtaTime:            itemTransfer.EtaTime,
				AtaDate:            itemTransfer.AtaDate,
				AtaTime:            itemTransfer.AtaTime,
				AdditionalCost:     itemTransfer.AdditionalCost,
				AdditionalCostNote: itemTransfer.AdditionalCostNote,
				SiteOriginID:       itemTransfer.SiteOriginID,
				SiteDestinationID:  itemTransfer.SiteDestinationID,
				StockType:          itemTransfer.StockType,
				Note:               itemTransfer.Note,
				TotalCost:          itemTransfer.TotalCost,
				TotalCharge:        itemTransfer.TotalCharge,
				TotalWeight:        itemTransfer.TotalWeight,
				Locked:             itemTransfer.Locked,
				LockedBy:           itemTransfer.LockedBy,
				TotalSku:           itemTransfer.TotalSku,
				Status:             itemTransfer.Status,
				UpdatedAt:          itemTransfer.UpdatedAt,
			},
		})
	}

	return
}

func (s *ReceivingService) GetDetail(ctx context.Context, id int64, code string) (res dto.ReceivingResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ReceivingService.GetDetail")
	defer span.End()

	var (
		it              *model.Receiving
		iri             []*model.ReceivingItem
		site            *model.Site
		po              *model.PurchaseOrder
		itemModel       *model.Item
		uomModel        *model.Uom
		itemTransfer    *model.ItemTransfer
		poRes           *dto.PurchaseOrderResponse
		itemTransferRes *dto.ItemTransferResponse
		iriRes          []*dto.ReceivingItemResponse
	)

	it, err = s.RepositoryReceiving.GetDetail(ctx, id, "")
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	site, err = s.RepositorySite.GetDetail(ctx, it.SiteId, "")
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if it.InboundType == 1 {
		po, err = s.RepositoryPurchaseOrder.GetDetail(ctx, it.PurchaseOrderId, "")
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		poRes = &dto.PurchaseOrderResponse{
			ID:                     po.ID,
			Code:                   po.Code,
			VendorID:               po.VendorID,
			SiteID:                 po.SiteID,
			TermPaymentPurID:       po.TermPaymentPurID,
			VendorClassificationID: po.VendorClassificationID,
			PurchasePlanID:         po.PurchasePlanID,
			ConsolidatedShipmentID: po.ConsolidatedShipmentID,
			Status:                 po.Status,
			RecognitionDate:        po.RecognitionDate,
			EtaDate:                po.EtaDate,
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
			UpdatedAt:              po.UpdatedAt,
			UpdatedBy:              po.UpdatedBy,
			CreatedFrom:            po.CreatedFrom,
			HasFinishedGr:          po.HasFinishedGr,
			CreatedAt:              po.CreatedAt,
			CreatedBy:              po.CreatedBy,
			CommittedAt:            po.CommittedAt,
			CommittedBy:            po.CommittedBy,
			AssignedTo:             po.AssignedTo,
			AssignedBy:             po.AssignedBy,
			AssignedAt:             po.AssignedAt,
			Locked:                 po.Locked,
			LockedBy:               po.LockedBy,
		}
	} else if it.InboundType == 2 {
		itemTransfer, err = s.RepositoryItemTransfer.GetDetail(ctx, it.ItemTransferId, "")
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		itemTransferRes = &dto.ItemTransferResponse{
			ID:                 itemTransfer.ID,
			Code:               itemTransfer.Code,
			RecognitionDate:    itemTransfer.RecognitionDate,
			RequestDate:        itemTransfer.RequestDate,
			EtaDate:            itemTransfer.EtaDate,
			EtaTime:            itemTransfer.EtaTime,
			AtaDate:            itemTransfer.AtaDate,
			AtaTime:            itemTransfer.AtaTime,
			AdditionalCost:     itemTransfer.AdditionalCost,
			AdditionalCostNote: itemTransfer.AdditionalCostNote,
			SiteOriginID:       itemTransfer.SiteOriginID,
			SiteDestinationID:  itemTransfer.SiteDestinationID,
			StockType:          itemTransfer.StockType,
			Note:               itemTransfer.Note,
			TotalCost:          itemTransfer.TotalCost,
			TotalCharge:        itemTransfer.TotalCharge,
			TotalWeight:        itemTransfer.TotalWeight,
			Locked:             itemTransfer.Locked,
			LockedBy:           itemTransfer.LockedBy,
			TotalSku:           itemTransfer.TotalSku,
			Status:             itemTransfer.Status,
			UpdatedAt:          itemTransfer.UpdatedAt,
		}
	}

	// receiving items
	iri, err = s.RepositoryReceivingItem.GetByReceivingId(ctx, it.ID)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, item := range iri {
		itemModel, err = s.RepositoryItem.GetDetail(ctx, item.ID, "")
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		uomModel, err = s.RepositoryUom.GetDetail(ctx, itemModel.UomID, "")
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		iriRes = append(iriRes, &dto.ReceivingItemResponse{
			ID:                  item.ID,
			PurchaseOrderItemID: item.PurchaseOrderItemID,
			ItemTransferItemID:  item.ItemTransferItemID,
			DeliverQty:          item.DeliverQty,
			RejectQty:           item.RejectQty,
			ReceiveQty:          item.ReceiveQty,
			Weight:              item.Weight,
			Note:                item.Note,
			RejectReason:        item.RejectReason,
			Item: &dto.ItemResponse{
				ID:                      itemModel.ID,
				Code:                    itemModel.Code,
				UomID:                   itemModel.UomID,
				ClassID:                 itemModel.ClassID,
				ItemCategoryID:          itemModel.ItemCategoryID,
				Description:             itemModel.Description,
				UnitWeightConversion:    itemModel.UnitWeightConversion,
				OrderMinQty:             itemModel.OrderMinQty,
				OrderMaxQty:             itemModel.OrderMaxQty,
				ItemType:                itemModel.ItemType,
				Packability:             itemModel.Packability,
				Capitalize:              itemModel.Capitalize,
				ExcludeArchetype:        itemModel.ExcludeArchetype,
				MaxDayDeliveryDate:      itemModel.MaxDayDeliveryDate,
				FragileGoods:            itemModel.FragileGoods,
				Taxable:                 itemModel.Taxable,
				OrderChannelRestriction: itemModel.OrderChannelRestriction,
				Note:                    itemModel.Note,
				Status:                  itemModel.Status,
				StatusConvert:           statusx.ConvertStatusValue(itemModel.Status),
				CreatedAt:               timex.ToLocTime(ctx, itemModel.CreatedAt),
				UpdatedAt:               timex.ToLocTime(ctx, itemModel.UpdatedAt),
				Uom: &dto.UomResponse{
					ID:             uomModel.ID,
					Code:           uomModel.Code,
					Description:    uomModel.Description,
					Status:         uomModel.Status,
					DecimalEnabled: uomModel.DecimalEnabled,
					CreatedAt:      uomModel.CreatedAt,
					UpdatedAt:      uomModel.UpdatedAt,
				},
			},
		})
	}

	res = dto.ReceivingResponse{
		ID:                  it.ID,
		SiteId:              it.SiteId,
		PurchaseOrderId:     it.PurchaseOrderId,
		ItemTransferId:      it.ItemTransferId,
		InboundType:         it.InboundType,
		ValidSupplierReturn: it.ValidSupplierReturn,
		CreatedAt:           it.CreatedAt,
		CreatedBy:           it.CreatedBy,
		ConfirmedAt:         it.ConfirmedAt,
		ConfirmedBy:         it.ConfirmedBy,
		Code:                it.Code,
		AtaDate:             it.AtaDate,
		AtaTime:             it.AtaTime,
		StockType:           it.StockType,
		TotalWeight:         it.TotalWeight,
		Note:                it.Note,
		Status:              it.Status,
		Locked:              it.Locked,
		LockedBy:            it.LockedBy,
		UpdatedAt:           it.UpdatedAt,
		UpdatedBy:           it.UpdatedBy,
		PurchaseOrder:       poRes,
		Site: &dto.SiteResponse{
			ID:          site.ID,
			Code:        site.Code,
			Description: site.Description,
			Status:      site.Status,
			CreatedAt:   site.CreatedAt,
			UpdatedAt:   site.UpdatedAt,
		},
		ItemTransfer:   itemTransferRes,
		ReceivingItems: iriRes,
	}

	return
}

func (s *ReceivingService) Create(ctx context.Context, req *dto.CreateReceivingRequest) (res dto.ReceivingResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ReceivingService.Create")
	defer span.End()

	var (
		r        model.Receiving
		items    []*model.ReceivingItem
		resItems []*dto.ReceivingItemResponse
		result   *model.Receiving
	)

	result, err = s.RepositoryReceiving.CreateWithItem(ctx, &r, items)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.ReceivingResponse{
		ID:                  result.ID,
		SiteId:              result.SiteId,
		PurchaseOrderId:     result.PurchaseOrderId,
		ItemTransferId:      result.ItemTransferId,
		InboundType:         result.InboundType,
		ValidSupplierReturn: result.ValidSupplierReturn,
		CreatedAt:           result.CreatedAt,
		CreatedBy:           result.CreatedBy,
		ConfirmedAt:         result.ConfirmedAt,
		ConfirmedBy:         result.ConfirmedBy,
		Code:                result.Code,
		StockType:           result.StockType,
		Note:                result.Note,
		Status:              result.Status,
		UpdatedAt:           result.UpdatedAt,
		UpdatedBy:           result.UpdatedBy,
		ReceivingItems:      resItems,
	}

	jsonRes, _ := json.Marshal(req)
	fmt.Println(string(jsonRes))
	return
}

func (s *ReceivingService) Confirm(ctx context.Context, req *dto.ConfirmReceivingRequest) (res dto.ReceivingResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ReceivingService.Confirm")
	defer span.End()

	var (
		r        model.Receiving
		items    []*model.ReceivingItem
		resItems []*dto.ReceivingItemResponse
		result   *model.Receiving
	)

	result, err = s.RepositoryReceiving.Confirm(ctx, &r, items)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.ReceivingResponse{
		ID:                  result.ID,
		SiteId:              result.SiteId,
		PurchaseOrderId:     result.PurchaseOrderId,
		ItemTransferId:      result.ItemTransferId,
		InboundType:         result.InboundType,
		ValidSupplierReturn: result.ValidSupplierReturn,
		CreatedAt:           result.CreatedAt,
		CreatedBy:           result.CreatedBy,
		ConfirmedAt:         result.ConfirmedAt,
		ConfirmedBy:         result.ConfirmedBy,
		Code:                result.Code,
		StockType:           result.StockType,
		Note:                result.Note,
		Status:              result.Status,
		UpdatedAt:           result.UpdatedAt,
		UpdatedBy:           result.UpdatedBy,
		ReceivingItems:      resItems,
	}

	jsonRes, _ := json.Marshal(req)
	fmt.Println(string(jsonRes))
	return
}

func (s *ReceivingService) GetGP(ctx context.Context, req *pb.GetGoodsReceiptGPListRequest) (res *pb.GetGoodsReceiptGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ReceivingService.GetGP")
	defer span.End()

	params := map[string]string{
		"interid":    global.EnvDatabaseGP,
		"PageNumber": strconv.Itoa(int(req.Offset)),
		"PageSize":   strconv.Itoa(int(req.Limit)),
	}

	if req.Poprctnm != "" {
		params["poprctnm"] = req.Poprctnm
	}

	if req.Doctype != "" {
		params["doctype"] = req.Doctype
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "goodsreceipt/list", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *ReceivingService) GetDetailGP(ctx context.Context, req *pb.GetGoodsReceiptGPDetailRequest) (res *pb.GetGoodsReceiptGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ReceivingService.GetDetailGP")
	defer span.End()

	params := map[string]string{
		"interid":  global.EnvDatabaseGP,
		"poprctnm": req.Poprctnm,
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "goodsreceipt/detail", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *ReceivingService) CreateGP(ctx context.Context, req *pb.CreateGoodsReceiptGPRequest) (res *pb.CreateTransferRequestGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ReceivingService.Create")
	defer span.End()

	req.Interid = global.EnvDatabaseGP

	err = global.HttpRestApiToMicrosoftGP("POST", "GoodsReceipt/create", req, &res, nil)
	if err != nil {
		logrus.Error(err.Error())
		err = errors.New("Connection to the server could not be established")
		return
	}

	return
}

func (s *ReceivingService) UpdateGP(ctx context.Context, req *pb.UpdateGoodsReceiptGPRequest) (res *pb.CreateTransferRequestGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ReceivingService.Update")
	defer span.End()

	req.Interid = global.EnvDatabaseGP

	err = global.HttpRestApiToMicrosoftGP("POST", "GoodsReceipt/update", req, &res, nil)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = errors.New("Connection to the server could not be established")
		return
	}

	return
}
