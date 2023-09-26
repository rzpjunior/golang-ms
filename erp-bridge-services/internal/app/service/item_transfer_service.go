package service

import (
	"context"
	"encoding/json"
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
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/repository"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

type IItemTransferService interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (res []dto.ItemTransferResponse, total int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (res dto.ItemTransferResponse, err error)
	Create(ctx context.Context, req *dto.CreateItemTransferRequest) (res dto.ItemTransferResponse, err error)
	Update(ctx context.Context, req *dto.UpdateItemTransferRequest) (res dto.ItemTransferResponse, err error)
	Commit(ctx context.Context, req *dto.CommitItemTransferRequest) (res dto.ItemTransferResponse, err error)
	GetInTransitTransferGP(ctx context.Context, req *pb.GetInTransitTransferGPListRequest) (res *pb.GetInTransitTransferGPResponse, err error)
	GetInTransitTransferDetailGP(ctx context.Context, req *pb.GetInTransitTransferGPDetailRequest) (res *pb.GetInTransitTransferGPResponse, err error)
	GetTransferRequestGP(ctx context.Context, req *pb.GetTransferRequestGPListRequest) (res *pb.GetTransferRequestGPResponse, err error)
	GetTransferRequestDetailGP(ctx context.Context, req *pb.GetTransferRequestGPDetailRequest) (res *pb.GetTransferRequestGPResponse, err error)
	CreateGP(ctx context.Context, req *dto.CreateTransferRequestGPRequest) (res *pb.CreateTransferRequestGPResponse, err error)
	UpdateGP(ctx context.Context, req *dto.UpdateTransferRequestGPRequest) (res *pb.CreateTransferRequestGPResponse, err error)
	UpdateITTGP(ctx context.Context, req *dto.UpdateInTransiteTransferGPRequest) (res *pb.UpdateInTransitTransferGPResponse, err error)
	CommitGP(ctx context.Context, req *dto.CommitTransferRequestGPRequest) (res *pb.CommitTransferRequestGPResponse, err error)
}

type ItemTransferService struct {
	opt                        opt.Options
	RepositoryItemTransfer     repository.IItemTransferRepository
	RepositoryItemTransferItem repository.IItemTransferItemRepository
	RepositorySite             repository.ISiteRepository
	RepositoryItem             repository.IItemRepository
	RepositoryUom              repository.IUomRepository
	RepositoryRegion           repository.IRegionRepository
	RepositoryReceiving        repository.IReceivingRepository
}

func NewItemTransferService() IItemTransferService {
	return &ItemTransferService{
		opt:                        global.Setup.Common,
		RepositoryItemTransfer:     repository.NewItemTransferRepository(),
		RepositoryItemTransferItem: repository.NewItemTransferItemRepository(),
		RepositorySite:             repository.NewSiteRepository(),
		RepositoryItem:             repository.NewItemRepository(),
		RepositoryUom:              repository.NewUomRepository(),
		RepositoryRegion:           repository.NewRegionRepository(),
		RepositoryReceiving:        repository.NewReceivingRepository(),
	}
}

func (s *ItemTransferService) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (res []dto.ItemTransferResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemTransferService.Get")
	defer span.End()

	var (
		itemTransfers                   []*model.ItemTransfer
		siteOrigin, siteDestination     *model.Site
		regionOrigin, regionDestination *model.Region
	)
	itemTransfers, total, err = s.RepositoryItemTransfer.Get(ctx, offset, limit, status, search, orderBy)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, it := range itemTransfers {
		var (
			receiving []*model.Receiving
			grs       []*dto.ReceivingListinDetailResponse
		)
		siteOrigin, err = s.RepositorySite.GetDetail(ctx, it.SiteOriginID, "")
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		regionOrigin, err = s.RepositoryRegion.GetDetail(ctx, siteOrigin.RegionId, "")
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		siteDestination, err = s.RepositorySite.GetDetail(ctx, it.SiteDestinationID, "")
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		regionDestination, err = s.RepositoryRegion.GetDetail(ctx, siteDestination.RegionId, "")
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		receiving, err = s.RepositoryReceiving.GetByInbound(ctx, 1, it.ID)
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

		res = append(res, dto.ItemTransferResponse{
			ID:                 it.ID,
			Code:               it.Code,
			RequestDate:        it.RequestDate,
			RecognitionDate:    it.RecognitionDate,
			EtaDate:            it.EtaDate,
			EtaTime:            it.EtaTime,
			AtaDate:            it.AtaDate,
			AtaTime:            it.AtaTime,
			AdditionalCost:     it.AdditionalCost,
			AdditionalCostNote: it.AdditionalCostNote,
			StockType:          it.StockType,
			TotalCost:          it.TotalCost,
			TotalCharge:        it.TotalCharge,
			TotalWeight:        it.TotalWeight,
			Note:               it.Note,
			Status:             it.Status,
			Locked:             it.Locked,
			LockedBy:           it.LockedBy,
			TotalSku:           it.TotalSku,
			UpdatedAt:          it.UpdatedAt,
			UpdatedBy:          it.UpdatedBy,
			HasFinishedGr:      it.HasFinishedGr,
			Receiving:          grs,

			SiteOrigin: &dto.SiteResponse{
				ID:            siteOrigin.ID,
				Code:          siteOrigin.Code,
				Description:   siteOrigin.Description,
				Status:        siteOrigin.Status,
				StatusConvert: statusx.ConvertStatusValue(siteOrigin.Status),
				CreatedAt:     siteOrigin.CreatedAt,
				UpdatedAt:     siteOrigin.UpdatedAt,
				Region: &dto.RegionResponse{
					ID:            regionOrigin.ID,
					Code:          regionOrigin.Code,
					Description:   regionOrigin.Description,
					Status:        regionOrigin.Status,
					StatusConvert: statusx.ConvertStatusValue(regionOrigin.Status),
					CreatedAt:     regionOrigin.CreatedAt,
					UpdatedAt:     regionOrigin.UpdatedAt,
				},
			},
			SiteDestination: &dto.SiteResponse{
				ID:            siteDestination.ID,
				Code:          siteDestination.Code,
				Description:   siteDestination.Description,
				Status:        siteDestination.Status,
				StatusConvert: statusx.ConvertStatusValue(siteDestination.Status),
				CreatedAt:     siteDestination.CreatedAt,
				UpdatedAt:     siteDestination.UpdatedAt,
				Region: &dto.RegionResponse{
					ID:            regionDestination.ID,
					Code:          regionDestination.Code,
					Description:   regionDestination.Description,
					Status:        regionDestination.Status,
					StatusConvert: statusx.ConvertStatusValue(regionDestination.Status),
					CreatedAt:     regionDestination.CreatedAt,
					UpdatedAt:     regionDestination.UpdatedAt,
				},
			},
		})
	}

	return
}

func (s *ItemTransferService) GetDetail(ctx context.Context, id int64, code string) (res dto.ItemTransferResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemTransferService.GetDetail")
	defer span.End()

	var (
		it                              *model.ItemTransfer
		iti                             []*model.ItemTransferItem
		itiRes                          []*dto.ItemTransferItemResponse
		receiving                       []*model.Receiving
		grs                             []*dto.ReceivingListinDetailResponse
		siteOrigin, siteDestination     *model.Site
		regionOrigin, regionDestination *model.Region
	)
	it, err = s.RepositoryItemTransfer.GetDetail(ctx, id, code)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	siteOrigin, err = s.RepositorySite.GetDetail(ctx, it.SiteOriginID, "")
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	regionOrigin, err = s.RepositoryRegion.GetDetail(ctx, siteOrigin.RegionId, "")
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	siteDestination, err = s.RepositorySite.GetDetail(ctx, it.SiteDestinationID, "")
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	regionDestination, err = s.RepositoryRegion.GetDetail(ctx, siteDestination.RegionId, "")
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	iti, err = s.RepositoryItemTransferItem.GetByItemTransferId(ctx, it.ID)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, item := range iti {
		var (
			itemModel *model.Item
			uomModel  *model.Uom
		)
		itemModel, err = s.RepositoryItem.GetDetail(ctx, item.ItemID, "")
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		uomModel, err = s.RepositoryUom.GetDetail(ctx, itemModel.UomID, "")
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		itiRes = append(itiRes, &dto.ItemTransferItemResponse{
			ID:          item.ID,
			DeliverQty:  item.DeliverQty,
			RequestQty:  item.RequestQty,
			ReceiveQty:  item.ReceiveQty,
			ReceiveNote: item.ReceiveNote,
			UnitCost:    item.UnitCost,
			Subtotal:    item.Subtotal,
			Weight:      item.Weight,
			Note:        item.Note,
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
				Note:                    itemModel.Note,
				ExcludeArchetype:        itemModel.ExcludeArchetype,
				MaxDayDeliveryDate:      itemModel.MaxDayDeliveryDate,
				FragileGoods:            itemModel.FragileGoods,
				Taxable:                 itemModel.Taxable,
				OrderChannelRestriction: itemModel.OrderChannelRestriction,
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

	receiving, err = s.RepositoryReceiving.GetByInbound(ctx, 1, it.ID)
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

	res = dto.ItemTransferResponse{
		ID:                 it.ID,
		Code:               it.Code,
		RequestDate:        it.RequestDate,
		RecognitionDate:    it.RecognitionDate,
		EtaDate:            it.EtaDate,
		EtaTime:            it.EtaTime,
		AtaDate:            it.AtaDate,
		AtaTime:            it.AtaTime,
		AdditionalCost:     it.AdditionalCost,
		AdditionalCostNote: it.AdditionalCostNote,
		StockType:          it.StockType,
		TotalCost:          it.TotalCost,
		TotalCharge:        it.TotalCharge,
		TotalWeight:        it.TotalWeight,
		Note:               it.Note,
		Status:             it.Status,
		Locked:             it.Locked,
		LockedBy:           it.LockedBy,
		TotalSku:           it.TotalSku,
		UpdatedAt:          it.UpdatedAt,
		UpdatedBy:          it.UpdatedBy,
		ItemTransferItem:   itiRes,
		HasFinishedGr:      it.HasFinishedGr,
		SiteOrigin: &dto.SiteResponse{
			ID:            siteOrigin.ID,
			Code:          siteOrigin.Code,
			Description:   siteOrigin.Description,
			Status:        siteOrigin.Status,
			StatusConvert: statusx.ConvertStatusValue(siteOrigin.Status),
			CreatedAt:     siteOrigin.CreatedAt,
			UpdatedAt:     siteOrigin.UpdatedAt,
			Region: &dto.RegionResponse{
				ID:            regionOrigin.ID,
				Code:          regionOrigin.Code,
				Description:   regionOrigin.Description,
				Status:        regionOrigin.Status,
				StatusConvert: statusx.ConvertStatusValue(regionOrigin.Status),
				CreatedAt:     regionOrigin.CreatedAt,
				UpdatedAt:     regionOrigin.UpdatedAt,
			},
		},
		SiteDestination: &dto.SiteResponse{
			ID:            siteDestination.ID,
			Code:          siteDestination.Code,
			Description:   siteDestination.Description,
			Status:        siteDestination.Status,
			StatusConvert: statusx.ConvertStatusValue(siteDestination.Status),
			CreatedAt:     siteDestination.CreatedAt,
			UpdatedAt:     siteDestination.UpdatedAt,
			Region: &dto.RegionResponse{
				ID:            regionDestination.ID,
				Code:          regionDestination.Code,
				Description:   regionDestination.Description,
				Status:        regionDestination.Status,
				StatusConvert: statusx.ConvertStatusValue(regionDestination.Status),
				CreatedAt:     regionDestination.CreatedAt,
				UpdatedAt:     regionDestination.UpdatedAt,
			},
		},
		Receiving: grs,
	}

	return
}

func (s *ItemTransferService) Create(ctx context.Context, req *dto.CreateItemTransferRequest) (res dto.ItemTransferResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemTransferService.Create")
	defer span.End()

	var (
		r           model.ItemTransfer
		items       []*model.ItemTransferItem
		resItems    []*dto.ItemTransferItemResponse
		result      *model.ItemTransfer
		requestDate time.Time
		productList = make(map[int64]bool)
	)

	// Time Validation
	if requestDate, err = time.Parse("2006-01-02", req.RequestDateStr); err != nil {
		err = edenlabs.ErrorInvalid("order_date")
		return
	}

	// cek site
	_, err = s.RepositorySite.GetDetail(ctx, req.SiteOriginID, "")
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("site_id")
		return
	}
	_, err = s.RepositorySite.GetDetail(ctx, req.SiteDestinationID, "")
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("site_id")
		return
	}

	// TODO: check stock opname

	if len(req.Note) > 250 {
		err = edenlabs.ErrorInvalid("note")
		return
	}

	if req.SiteOriginID == req.SiteDestinationID {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("site_id", "cannot be same with site destination")
		return
	}

	for i, v := range req.ItemTransferItems {
		var (
			uom  *model.Uom
			item *model.Item
		)

		// if _, exist := productList[v.ItemID]; exist {
		// 	err = edenlabs.ErrorDuplicate("product")
		// 	return
		// }

		if len(v.Note) > 100 {
			err = edenlabs.ErrorMustEqualOrGreater("note", "100")
			return
		}

		if v.ItemID == 0 {
			err = edenlabs.ErrorInvalid("item_id")
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
			if math.Mod(v.RequestQty, 1) != 0 {
				err = edenlabs.ErrorInvalid("request_qty")
				return
			}
		}

		productList[v.ItemID] = true
		items = append(items, &model.ItemTransferItem{
			ItemTransferID: 1,
			ItemID:         v.ItemID,
			RequestQty:     v.RequestQty,
			Note:           v.Note,
		})
		resItems = append(resItems, &dto.ItemTransferItemResponse{
			ID:             int64(i + 1),
			ItemTransferID: 1,
			ItemID:         v.ItemID,
			RequestQty:     v.RequestQty,
			Note:           v.Note,
		})
	}

	r = model.ItemTransfer{
		Code:              "IT0001",
		RequestDate:       requestDate,
		SiteOriginID:      req.SiteOriginID,
		SiteDestinationID: req.SiteDestinationID,
		StockType:         req.StockTypeID,
		Note:              req.Note,
		Status:            1,
	}

	result, err = s.RepositoryItemTransfer.CreateWithItem(ctx, &r, items)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.ItemTransferResponse{
		ID:                result.ID,
		Code:              result.Code,
		RequestDate:       result.RequestDate,
		SiteOriginID:      result.SiteOriginID,
		SiteDestinationID: result.SiteDestinationID,
		StockType:         result.StockType,
		Note:              result.Note,
		Status:            result.Status,
		UpdatedAt:         result.UpdatedAt,
		UpdatedBy:         result.UpdatedBy,
		ItemTransferItem:  resItems,
	}

	jsonRes, _ := json.Marshal(req)
	fmt.Println(string(jsonRes))
	return
}

func (s *ItemTransferService) Update(ctx context.Context, req *dto.UpdateItemTransferRequest) (res dto.ItemTransferResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemTransferService.Update")
	defer span.End()

	var (
		r           model.ItemTransfer
		items       []*model.ItemTransferItem
		resItems    []*dto.ItemTransferItemResponse
		result      *model.ItemTransfer
		requestDate time.Time
		productList = make(map[int64]bool)
	)

	// cek Id
	_, err = s.RepositoryItemTransfer.GetDetail(ctx, req.Id, "")
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("site_id")
		return
	}

	// Time Validation
	if requestDate, err = time.Parse("2006-01-02", req.RequestDateStr); err != nil {
		err = edenlabs.ErrorInvalid("request_date")
		return
	}

	// cek site
	_, err = s.RepositorySite.GetDetail(ctx, req.SiteOriginID, "")
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("site_id")
		return
	}
	_, err = s.RepositorySite.GetDetail(ctx, req.SiteDestinationID, "")
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("site_id")
		return
	}

	// TODO: check stock opname

	if len(req.Note) > 250 {
		err = edenlabs.ErrorInvalid("note")
		return
	}

	if req.SiteOriginID == req.SiteDestinationID {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("site_id", "cannot be same with site destination")
		return
	}

	for i, v := range req.ItemTransferItems {
		var (
			uom  *model.Uom
			item *model.Item
		)

		_, err = s.RepositoryItemTransferItem.GetDetail(ctx, v.Id)
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("site_id")
			return
		}

		if _, exist := productList[v.ItemID]; exist {
			err = edenlabs.ErrorDuplicate("product")
			return
		}

		if len(v.Note) > 100 {
			err = edenlabs.ErrorMustEqualOrGreater("note", "100")
			return
		}

		if v.ItemID == 0 {
			err = edenlabs.ErrorInvalid("item_id")
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
			if math.Mod(v.RequestQty, 1) != 0 {
				err = edenlabs.ErrorInvalid("request_qty")
				return
			}
		}

		productList[v.ItemID] = true
		items = append(items, &model.ItemTransferItem{
			ItemTransferID: 1,
			ItemID:         v.ItemID,
			RequestQty:     v.RequestQty,
			Note:           v.Note,
		})
		resItems = append(resItems, &dto.ItemTransferItemResponse{
			ID:             int64(i + 1),
			ItemTransferID: 1,
			ItemID:         v.ItemID,
			RequestQty:     v.RequestQty,
			Note:           v.Note,
		})
	}

	r = model.ItemTransfer{
		Code:              "IT0001",
		RequestDate:       requestDate,
		SiteOriginID:      req.SiteOriginID,
		SiteDestinationID: req.SiteDestinationID,
		Note:              req.Note,
		Status:            1,
	}

	result, err = s.RepositoryItemTransfer.UpdateWithItem(ctx, &r, items)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.ItemTransferResponse{
		ID:                result.ID,
		Code:              result.Code,
		RequestDate:       result.RequestDate,
		SiteOriginID:      result.SiteOriginID,
		SiteDestinationID: result.SiteDestinationID,
		StockType:         result.StockType,
		Note:              result.Note,
		Status:            result.Status,
		UpdatedAt:         result.UpdatedAt,
		UpdatedBy:         result.UpdatedBy,
		ItemTransferItem:  resItems,
	}

	jsonRes, _ := json.Marshal(req)
	fmt.Println(string(jsonRes))
	return
}

func (s *ItemTransferService) Commit(ctx context.Context, req *dto.CommitItemTransferRequest) (res dto.ItemTransferResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemTransferService.Create")
	defer span.End()

	var (
		r      model.ItemTransfer
		result *model.ItemTransfer
		// it                       *model.ItemTransfer
		items                    []*model.ItemTransferItem
		resItems                 []*dto.ItemTransferItemResponse
		recognitionDate, etaDate time.Time
		productList              = make(map[int64]bool)
	)

	// cek Id
	_, err = s.RepositoryItemTransfer.GetDetail(ctx, req.Id, "")
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("site_id")
		return
	}

	// if it.Status != 5 {
	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 	err = edenlabs.ErrorMustActive("id")
	// 	return
	// }

	// Time Validation
	if recognitionDate, err = time.Parse("2006-01-02", req.RecognitionDateStr); err != nil {
		err = edenlabs.ErrorInvalid("order_date")
		return
	}
	if etaDate, err = time.Parse("2006-01-02", req.EtaDateStr); err != nil {
		err = edenlabs.ErrorInvalid("order_date")
		return
	}

	if _, err = time.Parse("15:04", req.EtaTimeStr); err != nil {
		err = edenlabs.ErrorInvalid("order_date")
		return
	}

	if etaDate.Before(recognitionDate) {
		err = edenlabs.ErrorMustEqualOrGreater("eta_date", "recognition_date")
		return
	}

	if req.AdditionalCost < 0 {
		err = edenlabs.ErrorMustEqualOrGreater("biaya tambahan", "0")
		return
	}

	if req.AdditionalCost > 0 && req.AdditionalCostNote == "" {
		err = edenlabs.ErrorRequired("catatan tambahan")
		return
	}

	// TODO: check stock opname
	// TODO: check site restriction

	for i, v := range req.ItemTransferItems {
		var (
			uom  *model.Uom
			item *model.Item
		)

		// if _, exist := productList[v.ItemID]; exist {
		// 	err = edenlabs.ErrorDuplicate("product")
		// 	return
		// }

		_, err = s.RepositoryItemTransferItem.GetDetail(ctx, v.Id)
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("site_id")
			return
		}

		if len(v.Note) > 100 {
			err = edenlabs.ErrorMustEqualOrGreater("note", "100")
			return
		}

		if v.ItemID == 0 {
			err = edenlabs.ErrorInvalid("item_id")
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
			if v.RequestQty != float64((int64(v.RequestQty))) {
				err = edenlabs.ErrorInvalid("request_qty")
				return
			}
		}
		if v.RequestQty <= 0 {
			err = edenlabs.ErrorMustGreater("request_qty", "0")
			return
		}

		// TODO: CHECK STOCK DATA

		productList[v.ItemID] = true
		items = append(items, &model.ItemTransferItem{
			ItemTransferID: 1,
			ItemID:         v.ItemID,
			RequestQty:     v.RequestQty,
			Note:           v.Note,
		})
		resItems = append(resItems, &dto.ItemTransferItemResponse{
			ID:             int64(i + 1),
			ItemTransferID: 1,
			ItemID:         v.ItemID,
			RequestQty:     v.RequestQty,
			Note:           v.Note,
		})
	}

	r = model.ItemTransfer{
		Code:            "IT0001",
		RecognitionDate: recognitionDate,
		EtaDate:         etaDate,
		EtaTime:         req.EtaTimeStr,
		Status:          1,
	}

	result, err = s.RepositoryItemTransfer.CommitItemTransfer(ctx, &r, items)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.ItemTransferResponse{
		ID:                 result.ID,
		Code:               result.Code,
		RecognitionDate:    result.RecognitionDate,
		RequestDate:        result.RequestDate,
		EtaDate:            result.EtaDate,
		EtaTime:            result.EtaTime,
		AtaDate:            result.AtaDate,
		AtaTime:            result.AtaTime,
		AdditionalCost:     result.AdditionalCost,
		AdditionalCostNote: result.AdditionalCostNote,
		SiteOriginID:       result.SiteOriginID,
		SiteDestinationID:  result.SiteDestinationID,
		StockType:          result.StockType,
		Note:               result.Note,
		TotalCost:          result.TotalCost,
		TotalCharge:        result.TotalCharge,
		TotalWeight:        result.TotalWeight,
		Locked:             result.Locked,
		LockedBy:           result.LockedBy,
		TotalSku:           result.TotalSku,
		Status:             result.Status,
		UpdatedAt:          result.UpdatedAt,
		UpdatedBy:          result.UpdatedBy,
		ItemTransferItem:   resItems,
	}

	jsonRes, _ := json.Marshal(req)
	fmt.Println(string(jsonRes))
	return
}

func (s *ItemTransferService) GetInTransitTransferGP(ctx context.Context, req *pb.GetInTransitTransferGPListRequest) (res *pb.GetInTransitTransferGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "GetInTransitTransferGP.GetGP")
	defer span.End()

	params := map[string]string{
		"interid":    global.EnvDatabaseGP,
		"PageNumber": strconv.Itoa(int(req.Offset)),
		"PageSize":   strconv.Itoa(int(req.Limit)),
	}

	if req.Orddocid != "" {
		params["orddocid"] = req.Orddocid
	}
	if req.IvmTrType != "" {
		params["ivm_tr_type"] = req.IvmTrType
	}
	if req.Ordrdate != "" {
		params["ordrdate"] = req.Ordrdate
	}
	if req.Trnsfloc != "" {
		params["trnsfloc"] = req.Trnsfloc
	}
	if req.Locncode != "" {
		params["locncode"] = req.Locncode
	}
	if req.RequestDate != "" {
		params["request_date"] = req.RequestDate
	}
	if req.Etadte != "" {
		params["etadte"] = req.Etadte
	}
	if req.Status != 0 {
		params["status"] = strconv.Itoa(int(req.Status))
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "InTransitTransfer/list", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *ItemTransferService) GetInTransitTransferDetailGP(ctx context.Context, req *pb.GetInTransitTransferGPDetailRequest) (res *pb.GetInTransitTransferGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "GetInTransitTransferDetailGP.GetDetailGP")
	defer span.End()

	params := map[string]string{
		"interid": global.EnvDatabaseGP,
		"id":      req.Id,
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "InTransitTransfer/detail", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *ItemTransferService) GetTransferRequestGP(ctx context.Context, req *pb.GetTransferRequestGPListRequest) (res *pb.GetTransferRequestGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "GetTransferRequestGP.GetGP")
	defer span.End()

	var statusInt []string

	params := map[string]string{
		"interid":    global.EnvDatabaseGP,
		"PageNumber": strconv.Itoa(int(req.Offset)),
		"PageSize":   strconv.Itoa(int(req.Limit)),
	}

	if req.Docnumbr != "" {
		params["docnumbr"] = req.Docnumbr
	}
	if req.DocnumbrLike != "" {
		params["docnumbr_like"] = req.DocnumbrLike
	}
	if req.IvmTrType != "" {
		params["ivm_tr_type"] = req.IvmTrType
	}
	if req.RequestDateFrom != "" && req.RequestDateFrom != "0001-01-01" {
		params["request_date_from"] = req.RequestDateFrom
	}
	if req.RequestDateTo != "" && req.RequestDateTo != "0001-01-01" {
		params["request_date_to"] = req.RequestDateTo
	}
	if req.IvmLocncodeFrom != "" {
		params["ivm_locncode_from"] = req.IvmLocncodeFrom
	}
	if req.IvmLocncodeTo != "" {
		params["ivm_locncode_to"] = req.IvmLocncodeTo
	}
	if req.DocdateFrom != "" {
		params["docdate_from"] = req.DocdateFrom
	}
	if req.DocdateTo != "" {
		params["docdate_to"] = req.DocdateTo
	}
	if req.IvmReqEtaFrom != "" {
		params["ivm_req_eta_from"] = req.IvmReqEtaFrom
	}
	if req.IvmReqEtaTo != "" {
		params["ivm_req_eta_to"] = req.IvmReqEtaTo
	}
	// if req.Status != 0 {
	// 	params["status"] = strconv.Itoa(int(req.Status))
	// }
	// add status
	if len(req.Status) > 0 {
		for _, v := range req.Status {
			statusInt = append(statusInt, strconv.Itoa(int(v)))
		}
		params["status"] = strings.Join(statusInt, ",")
	}
	if req.Orderby != "" {
		params["orderby"] = req.Orderby
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "TransferRequest/GetAll", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *ItemTransferService) GetTransferRequestDetailGP(ctx context.Context, req *pb.GetTransferRequestGPDetailRequest) (res *pb.GetTransferRequestGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "GetTransferRequestDetailGP.GetDetailGP")
	defer span.End()

	params := map[string]string{
		"interid": global.EnvDatabaseGP,
		"id":      req.Id,
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "TransferRequest/GetByID", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *ItemTransferService) CreateGP(ctx context.Context, req *dto.CreateTransferRequestGPRequest) (res *pb.CreateTransferRequestGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CreateTransferRequestGP.CreateGP")
	defer span.End()

	req.Interid = global.EnvDatabaseGP
	err = global.HttpRestApiToMicrosoftGP("POST", "GoodsTransfer/create", req, &res, nil)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *ItemTransferService) UpdateGP(ctx context.Context, req *dto.UpdateTransferRequestGPRequest) (res *pb.CreateTransferRequestGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemTransferService.Update")
	defer span.End()

	req.Interid = global.EnvDatabaseGP
	err = global.HttpRestApiToMicrosoftGP("PUT", "GoodsTransfer/update", req, &res, nil)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *ItemTransferService) UpdateITTGP(ctx context.Context, req *dto.UpdateInTransiteTransferGPRequest) (res *pb.UpdateInTransitTransferGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "InTransitTransferService.Update")
	defer span.End()

	req.Interid = global.EnvDatabaseGP
	err = global.HttpRestApiToMicrosoftGP("PUT", "InTransitTransfer/update", req, &res, nil)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *ItemTransferService) CommitGP(ctx context.Context, req *dto.CommitTransferRequestGPRequest) (res *pb.CommitTransferRequestGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemTransferService.Commit")
	defer span.End()

	req.Interid = global.EnvDatabaseGP
	err = global.HttpRestApiToMicrosoftGP("PUT", "goodstransfer/commit", req, &res, nil)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}
