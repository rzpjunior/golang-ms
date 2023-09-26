package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/constants"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/global"
	"google.golang.org/protobuf/types/known/timestamppb"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/account_service"
	auditService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/audit_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
)

func NewServiceItemTransfer() IItemTransferService {
	m := new(ItemTransferService)
	m.opt = global.Setup.Common
	return m
}

type IItemTransferService interface {
	Get(ctx context.Context, req dto.ItemTransferListRequest) (res []*dto.ItemTransferResponse, err error)
	GetById(ctx context.Context, req dto.ItemTransferDetailRequest) (res *dto.ItemTransferResponse, err error)
	Create(ctx context.Context, req *dto.CreateItemTransferRequest) (res *dto.ItemTransferResponse, err error)
	Update(ctx context.Context, req *dto.UpdateItemTransferRequest) (res *dto.ItemTransferResponse, err error)
	Commit(ctx context.Context, req *dto.CommitItemTransferRequest) (res *dto.ItemTransferResponse, err error)
	// gp
	GetInTransitTransferListGP(ctx context.Context, req *dto.GetInTransitTransferGPListRequest) (res []*bridgeService.InTransitTransferGP, total int64, err error)
	GetInTransitTransferDetailGP(ctx context.Context, id string) (res *dto.ItemTransferResponse, err error)
	GetTransferRequestListGP(ctx context.Context, req *dto.GetTransferRequestGPListRequest) (res []*dto.ItemTransferResponse, total int64, err error)
	GetTransferRequestDetailGP(ctx context.Context, id string) (res *dto.ItemTransferResponse, err error)
	CreateTransferRequestGP(ctx context.Context, payload *dto.CreateTransferRequestGPRequest) (res *bridgeService.CreateTransferRequestGPResponse, err error)
	UpdateTransferRequestGP(ctx context.Context, payload *dto.UpdateTransferRequestGPRequest) (res *bridgeService.CreateTransferRequestGPResponse, err error)
	UpdateInTransitTransferGP(ctx context.Context, payload *dto.UpdateInTransiteTransferGPRequest) (res *bridgeService.UpdateInTransitTransferGPResponse, err error)
	CommitTransferRequestGP(ctx context.Context, payload *dto.CommitTransferRequestGPRequest) (res *bridgeService.CommitTransferRequestGPResponse, err error)
}

type ItemTransferService struct {
	opt opt.Options
}

func (s *ItemTransferService) Get(ctx context.Context, req dto.ItemTransferListRequest) (res []*dto.ItemTransferResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemTransferService.Get")
	defer span.End()

	// get Item Transfer from bridge
	var itRes *bridgeService.GetItemTransferListResponse
	itRes, err = s.opt.Client.BridgeServiceGrpc.GetItemTransferList(ctx, &bridgeService.GetItemTransferListRequest{
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

	datas := []*dto.ItemTransferResponse{}
	for _, it := range itRes.Data {
		var grs []*dto.ReceivingListinDetailResponse
		for _, gr := range it.Receiving {
			grs = append(grs, &dto.ReceivingListinDetailResponse{
				ID:     gr.Id,
				Code:   gr.Code,
				Status: int8(gr.Status),
			})
		}
		datas = append(datas, &dto.ItemTransferResponse{
			// ID:                 it.Id,
			Code:               it.Code,
			RequestDate:        it.RequestDate.AsTime(),
			RecognitionDate:    it.RecognitionDate.AsTime(),
			EtaDate:            it.EtaDate.AsTime(),
			EtaTime:            it.EtaTime,
			AtaDate:            it.AtaDate.AsTime(),
			AtaTime:            it.AtaTime,
			AdditionalCost:     it.AdditionalCost,
			AdditionalCostNote: it.AdditionalCostNote,
			StockType:          int8(it.StockType),
			TotalCost:          it.TotalCost,
			TotalCharge:        it.TotalCharge,
			TotalSku:           it.TotalSku,
			TotalWeight:        it.TotalWeight,
			Note:               it.Note,
			// Status:             int8(it.Status),
			Locked:    int8(it.Locked),
			LockedBy:  it.LockedBy,
			UpdatedAt: it.UpdatedAt.AsTime(),
			UpdatedBy: it.UpdatedBy,
			Receiving: grs,
			SiteOrigin: &dto.SiteResponse{
				// ID:            it.SiteOrigin.Id,
				Code:          it.SiteOrigin.Code,
				Name:          it.SiteOrigin.Description,
				Description:   it.SiteOrigin.Description,
				Status:        int8(it.SiteOrigin.Status),
				StatusConvert: statusx.ConvertStatusValue(int8(it.SiteOrigin.Status)),
				CreatedAt:     it.SiteOrigin.CreatedAt.AsTime(),
				UpdatedAt:     it.SiteOrigin.UpdatedAt.AsTime(),
			},
			SiteDestination: &dto.SiteResponse{
				// ID:            it.SiteDestination.Id,
				Code:          it.SiteDestination.Code,
				Description:   it.SiteDestination.Description,
				Name:          it.SiteDestination.Description,
				Status:        int8(it.SiteDestination.Status),
				StatusConvert: statusx.ConvertStatusValue(int8(it.SiteDestination.Status)),
				CreatedAt:     it.SiteDestination.CreatedAt.AsTime(),
				UpdatedAt:     it.SiteDestination.UpdatedAt.AsTime(),
			},
		})
	}
	res = datas

	return
}

func (s *ItemTransferService) GetById(ctx context.Context, req dto.ItemTransferDetailRequest) (res *dto.ItemTransferResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemTransferService.GetById")
	defer span.End()

	// get Purchase Order from bridge
	var (
		itRes *bridgeService.GetItemTransferDetailResponse
		iti   []*dto.ItemTransferItemResponse
		grs   []*dto.ReceivingListinDetailResponse
	)
	itRes, err = s.opt.Client.BridgeServiceGrpc.GetItemTransferDetail(ctx, &bridgeService.GetItemTransferDetailRequest{
		Id: int64(req.Id),
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "purchase order")
		return
	}

	for _, item := range itRes.Data.ItemTransferItem {
		iti = append(iti, &dto.ItemTransferItemResponse{
			ID: item.Id,
			// ItemTransferID: item.ItemTransferId,
			// ItemID:         item.ItemId,
			// DeliverQty:     item.DeliverQty,
			ReceiveQty:  item.ReceiveQty,
			RequestQty:  item.RequestQty,
			ReceiveNote: item.ReceiveNote,
			UnitCost:    item.UnitCost,
			Subtotal:    item.Subtotal,
			Weight:      item.Weight,
			Note:        item.Note,
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
				UnitPrice:            float64(10000),
			},
		})
	}

	for _, gr := range itRes.Data.Receiving {
		grs = append(grs, &dto.ReceivingListinDetailResponse{
			ID:     gr.Id,
			Code:   gr.Code,
			Status: int8(gr.Status),
		})
	}

	res = &dto.ItemTransferResponse{
		// ID:                 itRes.Data.Id,
		Code:               itRes.Data.Code,
		RequestDate:        itRes.Data.RequestDate.AsTime(),
		RecognitionDate:    itRes.Data.RecognitionDate.AsTime(),
		EtaDate:            itRes.Data.EtaDate.AsTime(),
		EtaTime:            itRes.Data.EtaTime,
		AtaDate:            itRes.Data.AtaDate.AsTime(),
		AtaTime:            itRes.Data.AtaTime,
		AdditionalCost:     itRes.Data.AdditionalCost,
		AdditionalCostNote: itRes.Data.AdditionalCostNote,
		StockType:          int8(itRes.Data.StockType),
		TotalCost:          itRes.Data.TotalCost,
		TotalCharge:        itRes.Data.TotalCharge,
		TotalSku:           itRes.Data.TotalSku,
		TotalWeight:        itRes.Data.TotalWeight,
		Note:               itRes.Data.Note,
		// Status:             int8(itRes.Data.Status),
		Locked:           int8(itRes.Data.Locked),
		LockedBy:         itRes.Data.LockedBy,
		UpdatedAt:        itRes.Data.UpdatedAt.AsTime(),
		UpdatedBy:        itRes.Data.UpdatedBy,
		ItemTransferItem: iti,
		SiteOrigin: &dto.SiteResponse{
			// ID:            itRes.Data.SiteOrigin.Id,
			Code:          itRes.Data.SiteOrigin.Code,
			Name:          itRes.Data.SiteOrigin.Description,
			Description:   itRes.Data.SiteOrigin.Description,
			Status:        int8(itRes.Data.SiteOrigin.Status),
			StatusConvert: statusx.ConvertStatusValue(int8(itRes.Data.SiteOrigin.Status)),
			CreatedAt:     itRes.Data.SiteOrigin.CreatedAt.AsTime(),
			UpdatedAt:     itRes.Data.SiteOrigin.UpdatedAt.AsTime(),
			Region: &dto.RegionResponse{
				// ID:            itRes.Data.SiteOrigin.Region.Id,
				Code:          itRes.Data.SiteOrigin.Region.Code,
				Description:   itRes.Data.SiteOrigin.Region.Description,
				Status:        int8(itRes.Data.SiteOrigin.Region.Status),
				StatusConvert: statusx.ConvertStatusValue(int8(itRes.Data.SiteOrigin.Region.Status)),
				CreatedAt:     itRes.Data.SiteOrigin.Region.CreatedAt.AsTime(),
				UpdatedAt:     itRes.Data.SiteOrigin.Region.UpdatedAt.AsTime(),
			},
		},
		SiteDestination: &dto.SiteResponse{
			// ID:            itRes.Data.SiteDestination.Id,
			Code:          itRes.Data.SiteDestination.Code,
			Name:          itRes.Data.SiteDestination.Description,
			Description:   itRes.Data.SiteDestination.Description,
			Status:        int8(itRes.Data.SiteDestination.Status),
			StatusConvert: statusx.ConvertStatusValue(int8(itRes.Data.SiteDestination.Status)),
			CreatedAt:     itRes.Data.SiteDestination.CreatedAt.AsTime(),
			UpdatedAt:     itRes.Data.SiteDestination.UpdatedAt.AsTime(),
			Region: &dto.RegionResponse{
				// ID:            itRes.Data.SiteDestination.Region.Id,
				Code:          itRes.Data.SiteDestination.Region.Code,
				Description:   itRes.Data.SiteDestination.Region.Description,
				Status:        int8(itRes.Data.SiteDestination.Region.Status),
				StatusConvert: statusx.ConvertStatusValue(int8(itRes.Data.SiteDestination.Region.Status)),
				CreatedAt:     itRes.Data.SiteDestination.Region.CreatedAt.AsTime(),
				UpdatedAt:     itRes.Data.SiteDestination.Region.UpdatedAt.AsTime(),
			},
		},
		Receiving: grs,
	}

	return
}

func (s *ItemTransferService) Create(ctx context.Context, req *dto.CreateItemTransferRequest) (res *dto.ItemTransferResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemTransferService.Create")
	defer span.End()

	var (
		r               bridge_service.CreateItemTransferRequest
		items           []*bridge_service.CreateItemTransferItemRequest
		resItems        []*dto.ItemTransferItemResponse
		result          *bridge_service.GetItemTransferDetailResponse
		user            *account_service.GetUserDetailResponse
		site            *bridge_service.GetSiteListResponse
		productList     = make(map[int64]bool)
		siteAccess      []int64
		siteRestriction = make(map[int64]bool)
	)

	// Time Validation
	if _, err = time.Parse("2006-01-02", req.RequestDateStr); err != nil {
		err = edenlabs.ErrorInvalid("order_date")
		return
	}

	_, err = s.opt.Client.ConfigurationServiceGrpc.GetGlossaryDetail(ctx, &configuration_service.GetGlossaryDetailRequest{
		Table:     "all",
		Attribute: "stock_type",
		ValueInt:  int32(req.StockTypeID),
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// cek site
	_, err = s.opt.Client.BridgeServiceGrpc.GetSiteDetail(ctx, &bridge_service.GetSiteDetailRequest{
		Id: req.SiteOriginID,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	// cek site
	_, err = s.opt.Client.BridgeServiceGrpc.GetSiteDetail(ctx, &bridge_service.GetSiteDetailRequest{
		Id: req.SiteDestinationID,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// cek site
	user, err = s.opt.Client.AccountServiceGrpc.GetUserDetail(ctx, &account_service.GetUserDetailRequest{
		Id: ctx.Value(constants.KeyUserID).(int64),
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, site := range strings.Split(user.Data.SiteAccess, ",") {
		idSite, _ := strconv.ParseInt(site, 0, 64)
		siteAccess = append(siteAccess, idSite)
	}

	site, err = s.opt.Client.BridgeServiceGrpc.GetSiteInIdsList(ctx, &bridge_service.GetSiteInIdsListRequest{
		Limit:  1000,
		Offset: 0,
		Ids:    siteAccess,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, v := range site.Data {
		siteRestriction[v.Id] = true
	}

	if ok, _ := siteRestriction[req.SiteDestinationID]; !ok {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("site_destination_id")
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
			uom  *bridge_service.GetUomDetailResponse
			item *bridge_service.GetItemDetailResponse
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

		item, err = s.opt.Client.BridgeServiceGrpc.GetItemDetail(ctx, &bridge_service.GetItemDetailRequest{
			Id: v.ItemID,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		uom, err = s.opt.Client.BridgeServiceGrpc.GetUomDetail(ctx, &bridge_service.GetUomDetailRequest{
			Id: item.Data.UomId,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		if uom.Data.DecimalEnabled == 2 {
			if math.Mod(v.RequestQty, 1) != 0 {
				err = edenlabs.ErrorInvalid("order_qty")
				return
			}
		}

		productList[v.ItemID] = true

		items = append(items, &bridge_service.CreateItemTransferItemRequest{
			ItemId:     1,
			RequestQty: v.RequestQty,
			Note:       v.Note,
		})
		resItems = append(resItems, &dto.ItemTransferItemResponse{
			ID: int64(i + 1),
			// ItemTransferID: 1,
			// ItemID:         v.ItemID,
			RequestQty: v.RequestQty,
			Note:       v.Note,
		})
	}

	r = bridgeService.CreateItemTransferRequest{
		RequestDateStr:    req.RequestDateStr,
		SiteOriginId:      req.SiteOriginID,
		SiteDestinationId: req.SiteDestinationID,
		StockTypeId:       int32(req.StockTypeID),
		Note:              req.Note,
		ItemTransferItems: items,
	}
	result, err = s.opt.Client.BridgeServiceGrpc.CreateItemTransfer(ctx, &r)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &dto.ItemTransferResponse{
		// ID:                result.Data.Id,
		Code:              result.Data.Code,
		RequestDate:       result.Data.RequestDate.AsTime(),
		SiteOriginID:      result.Data.SiteOriginId,
		SiteDestinationID: result.Data.SiteDestinationId,
		StockType:         int8(result.Data.StockType),
		Note:              result.Data.Note,
		// Status:            int8(result.Data.Status),
		UpdatedAt:        result.Data.UpdatedAt.AsTime(),
		UpdatedBy:        result.Data.UpdatedBy,
		ItemTransferItem: resItems,
	}

	jsonRes, _ := json.Marshal(req)
	fmt.Println(string(jsonRes))
	return
}

func (s *ItemTransferService) Update(ctx context.Context, req *dto.UpdateItemTransferRequest) (res *dto.ItemTransferResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemTransferService.Update")
	defer span.End()

	var (
		r               bridge_service.UpdateItemTransferRequest
		items           []*bridge_service.UpdateItemTransferItemRequest
		resItems        []*dto.ItemTransferItemResponse
		result          *bridge_service.GetItemTransferDetailResponse
		user            *account_service.GetUserDetailResponse
		site            *bridge_service.GetSiteListResponse
		productList     = make(map[int64]bool)
		siteAccess      []int64
		siteRestriction = make(map[int64]bool)
	)

	// cek id
	_, err = s.opt.Client.BridgeServiceGrpc.GetItemTransferDetail(ctx, &bridge_service.GetItemTransferDetailRequest{
		Id: req.Id,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// Time Validation
	if _, err = time.Parse("2006-01-02", req.RequestDateStr); err != nil {
		err = edenlabs.ErrorInvalid("request_date")
		return
	}

	// cek site
	_, err = s.opt.Client.BridgeServiceGrpc.GetSiteDetail(ctx, &bridge_service.GetSiteDetailRequest{
		Id: req.SiteOriginID,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	// cek site
	_, err = s.opt.Client.BridgeServiceGrpc.GetSiteDetail(ctx, &bridge_service.GetSiteDetailRequest{
		Id: req.SiteDestinationID,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// cek site
	user, err = s.opt.Client.AccountServiceGrpc.GetUserDetail(ctx, &account_service.GetUserDetailRequest{
		Id: ctx.Value(constants.KeyUserID).(int64),
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, site := range strings.Split(user.Data.SiteAccess, ",") {
		idSite, _ := strconv.ParseInt(site, 0, 64)
		siteAccess = append(siteAccess, idSite)
	}

	site, err = s.opt.Client.BridgeServiceGrpc.GetSiteInIdsList(ctx, &bridge_service.GetSiteInIdsListRequest{
		Limit:  1000,
		Offset: 0,
		Ids:    siteAccess,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, v := range site.Data {
		siteRestriction[v.Id] = true
	}

	if ok, _ := siteRestriction[req.SiteDestinationID]; !ok {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("site_destination_id")
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
			uom  *bridge_service.GetUomDetailResponse
			item *bridge_service.GetItemDetailResponse
		)

		// cek id
		_, err = s.opt.Client.BridgeServiceGrpc.GetItemTransferItemDetail(ctx, &bridge_service.GetItemTransferItemDetailRequest{
			Id: v.Id,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
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

		item, err = s.opt.Client.BridgeServiceGrpc.GetItemDetail(ctx, &bridge_service.GetItemDetailRequest{
			Id: v.ItemID,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		uom, err = s.opt.Client.BridgeServiceGrpc.GetUomDetail(ctx, &bridge_service.GetUomDetailRequest{
			Id: item.Data.UomId,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		if uom.Data.DecimalEnabled == 2 {
			if math.Mod(v.RequestQty, 1) != 0 {
				err = edenlabs.ErrorInvalid("order_qty")
				return
			}
		}

		productList[v.ItemID] = true

		items = append(items, &bridge_service.UpdateItemTransferItemRequest{
			ItemId:     1,
			RequestQty: v.RequestQty,
			Note:       v.Note,
		})
		resItems = append(resItems, &dto.ItemTransferItemResponse{
			ID: int64(i + 1),
			// ItemTransferID: 1,
			// ItemID:         v.ItemID,
			RequestQty: v.RequestQty,
			Note:       v.Note,
		})
	}

	r = bridgeService.UpdateItemTransferRequest{
		Id:                 req.Id,
		RecognitionDateStr: req.RecognitionDateStr,
		EtaDateStr:         req.EtaDateStr,
		EtaTimeStr:         req.EtaTimeStr,
		AdditionalCost:     req.AdditionalCost,
		AdditionalCostNote: req.AdditionalCostNote,
		RequestDateStr:     req.RequestDateStr,
		SiteOriginId:       req.SiteOriginID,
		SiteDestinationId:  req.SiteDestinationID,
		Note:               req.Note,
		ItemTransferItems:  items,
	}
	result, err = s.opt.Client.BridgeServiceGrpc.UpdateItemTransfer(ctx, &r)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &dto.ItemTransferResponse{
		// ID:                 result.Data.Id,
		RecognitionDate:    result.Data.RecognitionDate.AsTime(),
		EtaDate:            result.Data.EtaDate.AsTime(),
		EtaTime:            result.Data.EtaTime,
		AtaDate:            result.Data.AtaDate.AsTime(),
		AtaTime:            result.Data.AtaTime,
		AdditionalCost:     result.Data.AdditionalCost,
		AdditionalCostNote: result.Data.AdditionalCostNote,
		TotalCost:          result.Data.TotalCost,
		TotalCharge:        result.Data.TotalCharge,
		TotalWeight:        result.Data.TotalWeight,
		Locked:             int8(result.Data.Locked),
		LockedBy:           result.Data.LockedBy,
		TotalSku:           result.Data.TotalSku,
		Code:               result.Data.Code,
		RequestDate:        result.Data.RequestDate.AsTime(),
		SiteOriginID:       result.Data.SiteOriginId,
		SiteDestinationID:  result.Data.SiteDestinationId,
		StockType:          int8(result.Data.StockType),
		Note:               result.Data.Note,
		// Status:             int8(result.Data.Status),
		UpdatedAt:        result.Data.UpdatedAt.AsTime(),
		UpdatedBy:        result.Data.UpdatedBy,
		ItemTransferItem: resItems,
	}

	jsonRes, _ := json.Marshal(req)
	fmt.Println(string(jsonRes))
	return
}

func (s *ItemTransferService) Commit(ctx context.Context, req *dto.CommitItemTransferRequest) (res *dto.ItemTransferResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemTransferService.Commit")
	defer span.End()

	var (
		r                        bridge_service.CommitItemTransferRequest
		items                    []*bridge_service.UpdateItemTransferItemRequest
		resItems                 []*dto.ItemTransferItemResponse
		result                   *bridge_service.GetItemTransferDetailResponse
		recognitionDate, etaDate time.Time
		productList              = make(map[int64]bool)
	)

	// cek id
	_, err = s.opt.Client.BridgeServiceGrpc.GetItemTransferDetail(ctx, &bridge_service.GetItemTransferDetailRequest{
		Id: req.Id,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// Time Validation
	if recognitionDate, err = time.Parse("2006-01-02", req.RecognitionDateStr); err != nil {
		err = edenlabs.ErrorInvalid("request_date")
		return
	}

	if etaDate, err = time.Parse("2006-01-02", req.EtaDateStr); err != nil {
		err = edenlabs.ErrorInvalid("request_date")
		return
	}

	if _, err = time.Parse("15:04", req.EtaTimeStr); err != nil {
		err = edenlabs.ErrorInvalid("request_date")
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

	for i, v := range req.ItemTransferItems {
		var (
			uom  *bridge_service.GetUomDetailResponse
			item *bridge_service.GetItemDetailResponse
		)

		// cek id
		_, err = s.opt.Client.BridgeServiceGrpc.GetItemTransferItemDetail(ctx, &bridge_service.GetItemTransferItemDetailRequest{
			Id: v.Id,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

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

		item, err = s.opt.Client.BridgeServiceGrpc.GetItemDetail(ctx, &bridge_service.GetItemDetailRequest{
			Id: v.ItemID,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		uom, err = s.opt.Client.BridgeServiceGrpc.GetUomDetail(ctx, &bridge_service.GetUomDetailRequest{
			Id: item.Data.UomId,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		if uom.Data.DecimalEnabled == 2 {
			if math.Mod(v.RequestQty, 1) != 0 {
				err = edenlabs.ErrorInvalid("order_qty")
				return
			}
		}

		productList[v.ItemID] = true

		items = append(items, &bridge_service.UpdateItemTransferItemRequest{
			Id:          v.Id,
			ItemId:      v.ItemID,
			TransferQty: v.TransferQty,
			ReceiveQty:  v.ReceiveQty,
			ReceiveNote: v.ReceiveNote,
			UnitCost:    v.UnitCost,
			RequestQty:  v.RequestQty,
			Note:        v.Note,
		})
		resItems = append(resItems, &dto.ItemTransferItemResponse{
			ID: int64(i + 1),
			// ItemTransferID: 1,
			// ItemID:         v.ItemID,
			RequestQty: v.RequestQty,
			Note:       v.Note,
		})
	}

	r = bridgeService.CommitItemTransferRequest{
		Id:                 req.Id,
		RecognitionDate:    req.RecognitionDateStr,
		EtaDate:            req.EtaDateStr,
		EtaTime:            req.EtaTimeStr,
		AdditionalCost:     req.AdditionalCost,
		AdditionalCostNote: req.AdditionalCostNote,
		ItemTransferItems:  items,
	}
	result, err = s.opt.Client.BridgeServiceGrpc.CommitItemTransfer(ctx, &r)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &dto.ItemTransferResponse{
		// ID:                 result.Data.Id,
		RecognitionDate:    result.Data.RecognitionDate.AsTime(),
		EtaDate:            result.Data.EtaDate.AsTime(),
		EtaTime:            result.Data.EtaTime,
		AtaDate:            result.Data.AtaDate.AsTime(),
		AtaTime:            result.Data.AtaTime,
		AdditionalCost:     result.Data.AdditionalCost,
		AdditionalCostNote: result.Data.AdditionalCostNote,
		TotalCost:          result.Data.TotalCost,
		TotalCharge:        result.Data.TotalCharge,
		TotalWeight:        result.Data.TotalWeight,
		Locked:             int8(result.Data.Locked),
		LockedBy:           result.Data.LockedBy,
		TotalSku:           result.Data.TotalSku,
		Code:               result.Data.Code,
		RequestDate:        result.Data.RequestDate.AsTime(),
		SiteOriginID:       result.Data.SiteOriginId,
		SiteDestinationID:  result.Data.SiteDestinationId,
		StockType:          int8(result.Data.StockType),
		Note:               result.Data.Note,
		// Status:             int8(result.Data.Status),
		UpdatedAt:        result.Data.UpdatedAt.AsTime(),
		UpdatedBy:        result.Data.UpdatedBy,
		ItemTransferItem: resItems,
	}

	jsonRes, _ := json.Marshal(req)
	fmt.Println(string(jsonRes))
	return
}

func (s *ItemTransferService) GetInTransitTransferListGP(ctx context.Context, req *dto.GetInTransitTransferGPListRequest) (res []*bridgeService.InTransitTransferGP, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "InTransitTransferService.GetInTransitTransferListGP")
	defer span.End()

	var po *bridgeService.GetInTransitTransferGPResponse

	if po, err = s.opt.Client.BridgeServiceGrpc.GetInTransitTransferGPList(ctx, &bridgeService.GetInTransitTransferGPListRequest{
		Limit:       int32(req.Limit),
		Offset:      int32(req.Offset),
		Orddocid:    req.Orddocid,
		IvmTrType:   req.IvmTrType,
		Ordrdate:    req.Ordrdate,
		Trnsfloc:    req.Trnsfloc,
		Locncode:    req.Locncode,
		RequestDate: req.RequestDate,
		Etadte:      req.Etadte,
		Status:      req.Status,
	}); err != nil || !po.Succeeded {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "in transit transfer")
		return
	}

	total = int64(len(po.Data))
	res = po.Data

	return
}

func (s *ItemTransferService) GetInTransitTransferDetailGP(ctx context.Context, id string) (res *dto.ItemTransferResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "InTransitTransferService.GetDetailGP")
	defer span.End()

	var (
		itt                    *bridgeService.GetInTransitTransferGPResponse
		etaDateITT, etaTimeITT time.Time
		statusStr              string
		itemTransferItem       []*dto.ItemTransferItemResponse
		status                 int32
	)

	// get ITT
	if itt, err = s.opt.Client.BridgeServiceGrpc.GetInTransitTransferGPDetail(ctx, &bridgeService.GetInTransitTransferGPDetailRequest{
		Id: id,
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

	recognitionDate, _ := time.Parse("2006-01-02", itt.Data[0].Ordrdate)
	switch itt.Data[0].Status {
	case 1:
		status = 1
		statusStr = "Draft"
	case 4:
		status = 4
		statusStr = "Active"
	default:
		status = 4
		statusStr = "Active"
	}

	for _, v := range itt.Data[0].Details {
		itemTransferItem = append(itemTransferItem, &dto.ItemTransferItemResponse{
			ItemID:        v.Itemnmbr,
			ItemName:      v.Dscriptn,
			Uom:           v.Uofm,
			TransferQty:   v.Trnsfqty,
			FullfilledQty: v.Qtyfulfi,
			ShippedQty:    v.Qtyshppd,
		})
	}

	res = &dto.ItemTransferResponse{
		ID:              itt.Data[0].Orddocid,
		EtaDate:         etaDateITT,
		EtaTime:         etaTimeITT.Format("15:04"),
		VendorID:        "-",
		RecognitionDate: recognitionDate,
		StatusGP:        itt.Data[0].Status,
		Status:          status,
		StatusStr:       statusStr,
		ReasonCodeStr:   "Not provivided by GP yet",
		ShippingMethod:  itt.Data[0].Shipmthd,
		SiteOrigin: &dto.SiteResponse{
			ID:          itt.Data[0].Itlocn[0].Locndscr, //Not provivided by GP yet
			Name:        itt.Data[0].Itlocn[0].Locndscr,
			Description: itt.Data[0].Itlocn[0].Locndscr,
		},
		SiteDestination: &dto.SiteResponse{
			ID:          itt.Data[0].Locncode[0].Locndscr, //Not provivided by GP yet
			Name:        itt.Data[0].Locncode[0].Locndscr,
			Description: itt.Data[0].Locncode[0].Locndscr,
		},
		ItemTransferItem: itemTransferItem,
	}

	return
}

func (s *ItemTransferService) GetTransferRequestListGP(ctx context.Context, req *dto.GetTransferRequestGPListRequest) (res []*dto.ItemTransferResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "TransferRequestService.GetTransferRequestListGP")
	defer span.End()

	var (
		tr                                *bridgeService.GetTransferRequestGPResponse
		recognitionDate, etaDate, etaTime time.Time
		layoutInput, statusStr            string
		statusFilter                      []int32
		status                            int32
	)

	layoutInput = "2006-01-02" // The layout for the input string
	// layoutOutput = "2006-01-02"         // The layout for the desired output

	switch req.Status {
	case 1:
		statusFilter = []int32{1}
	case 4:
		statusFilter = []int32{2, 3, 4}
	default:
		statusFilter = []int32{1, 2, 3, 4}
	}

	if tr, err = s.opt.Client.BridgeServiceGrpc.GetTransferRequestGPList(ctx, &bridgeService.GetTransferRequestGPListRequest{
		Limit:           int32(req.Limit),
		Offset:          int32(req.Offset),
		DocnumbrLike:    req.Docnumbr,
		RequestDateFrom: req.RequestDateFrom.Format("2006-01-02"),
		RequestDateTo:   req.RequestDateTo.Format("2006-01-02"),
		IvmLocncodeFrom: req.IvmLocncodeFrom,
		IvmLocncodeTo:   req.IvmLocncodeTo,
		Orderby:         req.OrderBy,
		Status:          statusFilter,
	}); err != nil || !tr.Succeeded {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "transfer request")
		return
	}

	for _, v := range tr.Data {
		/**
		it should be deleted later on, because the data showed
		if there is data dummy from gp, will not match with amount of perpage
		**/
		recognitionDate, _ = time.Parse(layoutInput, v.Docdate)

		etaDate, err = time.Parse("2006-01-02", v.RequestDate)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("eta_date")
			return
		}

		etaTime, err = time.Parse("15:04:05", v.IvmReqEta)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("eta_time")
			return
		}

		// for _, detail := range v.Details {
		// 	// handling uofm non KG, it has to be converted to KG
		// 	if detail.Uofm != "KG" {

		// 	} else {
		// 		totalWeight += detail.Qtyorder
		// 	}
		// }

		switch v.IvmStatus {
		case 1:
			status = 1
			statusStr = "Draft"
		case 2, 3, 4:
			status = 4
			statusStr = "Active"
		}
		res = append(res, &dto.ItemTransferResponse{
			ID:              v.Docnumbr,
			RecognitionDate: recognitionDate,
			ReasonCodeStr:   tr.Data[0].ReasonCode,
			EtaDate:         etaDate,
			EtaTime:         etaTime.Format("15:04"),
			SiteOrigin: &dto.SiteResponse{
				ID:          v.IvmLocncodeFrom[0].Locndscr, // id not showed by gp response
				Name:        v.IvmLocncodeFrom[0].Locndscr,
				Description: v.IvmLocncodeFrom[0].Locndscr,
				Address:     v.IvmLocncodeFrom[0].Address,
			},
			SiteDestination: &dto.SiteResponse{
				ID:          v.IvmLocncodeTo[0].Locndscr, // id not showed by gp response
				Name:        v.IvmLocncodeTo[0].Locndscr,
				Description: v.IvmLocncodeTo[0].Locndscr,
				Address:     v.IvmLocncodeTo[0].Address,
			},
			TotalWeight: v.TotalWeight,
			StatusGP:    v.IvmStatus,
			Status:      status,
			StatusStr:   statusStr,
			Type:        v.IvmTrType,
			TypeStr:     v.IvmTrTypeDesc,
		})
	}

	total = int64(len(res))

	return
}

func (s *ItemTransferService) GetTransferRequestDetailGP(ctx context.Context, id string) (res *dto.ItemTransferResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "TransferRequestService.GetDetailGP")
	defer span.End()

	var (
		tr                                *bridgeService.GetTransferRequestGPResponse
		itemTransferItem                  []*dto.ItemTransferItemResponse
		receiving                         []*dto.ReceivingListinDetailResponse
		recognitionDate, etaDate, etaTime time.Time
		layoutInput, statusStr            string
		status, statusITT, statusGR       int32
		itt                               *dto.InTransitTransferDetailResponse
	)

	layoutInput = "2006-01-02" // The layout for the input string
	// layoutOutput = "2006-01-02"         // The layout for the desired output

	receiving = []*dto.ReceivingListinDetailResponse{}

	if tr, err = s.opt.Client.BridgeServiceGrpc.GetTransferRequestGPList(ctx, &bridgeService.GetTransferRequestGPListRequest{
		Docnumbr: id,
		// Status:          req.Status,
	}); err != nil || !tr.Succeeded {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "transfer request")
		return
	}

	/**
	it should be deleted later on, because the data showed
	if there is data dummy from gp, will not match with amount of perpage
	**/
	recognitionDate, _ = time.Parse(layoutInput, tr.Data[0].Docdate)

	etaDate, err = time.Parse("2006-01-02", tr.Data[0].RequestDate)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("eta_date")
		return
	}

	etaTime, err = time.Parse("15:04:05", tr.Data[0].IvmReqEta)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("eta_time")
		return
	}

	if len(tr.Data[0].Intrxlist) > 0 {
		switch tr.Data[0].Intrxlist[0].Status {
		case 1:
			statusITT = 1
			statusStr = "Draft"
		case 2, 3, 4:
			statusITT = 4
			statusStr = "Active"
		default:
			statusITT = 4
			statusStr = "Active"
		}
		itt = &dto.InTransitTransferDetailResponse{
			ID:        tr.Data[0].Intrxlist[0].Orddocid,
			Status:    int8(statusITT),
			StatusStr: statusStr,
		}

		if len(tr.Data[0].Intrxlist[0].GoodsReceipt) > 0 {
			for _, v := range tr.Data[0].Intrxlist[0].GoodsReceipt {
				switch v.Status {
				case 1:
					statusGR = 1
					statusStr = "Active"
				case 2:
					statusGR = 2
					statusStr = "Finished"
				default:
					statusGR = 1
					statusStr = "Active"
				}
				receiving = append(receiving, &dto.ReceivingListinDetailResponse{
					ID:        v.Poprctnm,
					Status:    int8(statusGR),
					StatusStr: statusStr,
				})
			}
		}
	}

	for _, detail := range tr.Data[0].Details {
		// handling uofm non KG, it has to be converted to KG
		itemTransferItem = append(itemTransferItem, &dto.ItemTransferItemResponse{
			LnitmSeq:      detail.Lnitmseq,
			ItemID:        detail.Itemnmbr,
			ItemName:      detail.Itemdesc,
			Uom:           detail.Uofm,
			RequestQty:    detail.IvmQtyRequest,
			FullfilledQty: detail.IvmQtyFulfill,
			// Note:      detail.Commntid,
		})
	}

	switch tr.Data[0].IvmStatus {
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
	res = &dto.ItemTransferResponse{
		ID:              tr.Data[0].Docnumbr,
		RecognitionDate: recognitionDate,
		ReasonCodeStr:   tr.Data[0].ReasonCode,
		EtaDate:         etaDate,
		EtaTime:         etaTime.Format("15:04"),
		SiteOrigin: &dto.SiteResponse{
			ID:          tr.Data[0].IvmLocncodeFrom[0].Locndscr, // id not showed by gp response
			Name:        tr.Data[0].IvmLocncodeFrom[0].Locndscr,
			Description: tr.Data[0].IvmLocncodeFrom[0].Locndscr,
			Address:     tr.Data[0].IvmLocncodeFrom[0].Address,
		},
		SiteDestination: &dto.SiteResponse{
			ID:          tr.Data[0].IvmLocncodeTo[0].Locndscr, // id not showed by gp response
			Name:        tr.Data[0].IvmLocncodeTo[0].Locndscr,
			Description: tr.Data[0].IvmLocncodeTo[0].Locndscr,
			Address:     tr.Data[0].IvmLocncodeTo[0].Address,
		},
		TotalWeight:       tr.Data[0].TotalWeight,
		StatusGP:          tr.Data[0].IvmStatus,
		Status:            status,
		StatusStr:         statusStr,
		Type:              tr.Data[0].IvmTrType,
		TypeStr:           tr.Data[0].IvmTrTypeDesc,
		ItemTransferItem:  itemTransferItem,
		Receiving:         receiving,
		InTransitTransfer: itt,
	}

	return
}

// TR Create
func (s *ItemTransferService) CreateTransferRequestGP(ctx context.Context, payload *dto.CreateTransferRequestGPRequest) (res *bridgeService.CreateTransferRequestGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "TransferRequestService.CreateDetailGP")
	defer span.End()

	var details []*bridgeService.CreateTransferRequestGPRequest_Detail
	for _, v := range payload.Detail {
		details = append(details, &bridgeService.CreateTransferRequestGPRequest_Detail{
			Lnitmseq:      int32(v.Sequence),
			Itemnmbr:      v.ItemID,
			Uofm:          v.Uom,
			IvmQtyRequest: v.RequestQty,
			IvmQtyFulfill: v.FulfilledQty,
		})
	}

	// var codeGenerator *configuration_service.GetGenerateCodeResponse
	// if codeGenerator, err = s.opt.Client.ConfigurationServiceGrpc.GetGenerateCode(ctx, &configuration_service.GetGenerateCodeRequest{
	// 	Format: "TR",
	// 	Domain: "Site",
	// 	Length: 6,
	// }); err != nil {
	// 	span.RecordError(err)
	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 	err = edenlabs.ErrorRpcNotFound("configuration", "code generator")
	// 	return
	// }

	if res, err = s.opt.Client.BridgeServiceGrpc.CreateTransferRequestGP(ctx, &bridgeService.CreateTransferRequestGPRequest{
		Docnumbr:        "",
		Docdate:         payload.RecognitionDate,
		IvmTrType:       int32(payload.Type),
		RequestDate:     payload.EtaDate,
		IvmReqEta:       payload.EtaTime,
		IvmLocncodeFrom: payload.SiteFrom,
		IvmLocncodeTo:   payload.SiteTo,
		Reason_Code:     payload.ReasonCode,
		Detail:          details,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "transfer request")
		return
	}

	userID := ctx.Value(constants.KeyUserID).(int64)
	_, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &auditService.CreateLogRequest{
		Log: &auditService.Log{
			UserId:      userID,
			ReferenceId: res.Docnumbr,
			Type:        "item transfer",
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

func (s *ItemTransferService) UpdateTransferRequestGP(ctx context.Context, payload *dto.UpdateTransferRequestGPRequest) (res *bridgeService.CreateTransferRequestGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "TransferRequestService.UpdateTransferRequestGP")
	defer span.End()
	var (
		userID  int64
		tr      *bridgeService.GetTransferRequestGPResponse
		details []*bridgeService.UpdateTransferRequestGPRequest_Detail
	)

	if tr, err = s.opt.Client.BridgeServiceGrpc.GetTransferRequestGPDetail(ctx, &bridgeService.GetTransferRequestGPDetailRequest{
		Id: payload.Docnumbr,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "transfer request")
		return
	}

	if tr.Data[0].IvmStatus != 1 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorMustDraft("transfer request")
		return
	}

	for _, detail := range payload.Detail {
		details = append(details, &bridgeService.UpdateTransferRequestGPRequest_Detail{
			Lnitmseq:      int32(detail.Lnitmseq),
			Itemnmbr:      detail.Itemnmbr,
			IvmQtyRequest: detail.RequestQty,
		})
	}

	if res, err = s.opt.Client.BridgeServiceGrpc.UpdateTransferRequestGP(ctx, &bridgeService.UpdateTransferRequestGPRequest{
		Docnumbr:    payload.Docnumbr,
		Docdate:     payload.Docdate,
		RequestDate: payload.RequestDate,
		IvmReqEta:   payload.IvmReqEta,
		Reason_Code: payload.ReasonCode,
		Detail:      details,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "transfer request")
		return
	}

	userID = ctx.Value(constants.KeyUserID).(int64)
	_, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &auditService.CreateLogRequest{
		Log: &auditService.Log{
			UserId:      userID,
			ReferenceId: payload.Docnumbr,
			Type:        "item transfer",
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

func (s *ItemTransferService) UpdateInTransitTransferGP(ctx context.Context, payload *dto.UpdateInTransiteTransferGPRequest) (res *bridgeService.UpdateInTransitTransferGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "TransferRequestService.UpdateTransferRequestGP")
	defer span.End()
	var (
		userID      int64
		tr          *bridgeService.GetTransferRequestGPResponse
		itt         *bridgeService.GetInTransitTransferGPResponse
		details     []*bridgeService.UpdateInTransitTransferGPRequest_Detail
		trDetailMap map[string]float64
	)

	if itt, err = s.opt.Client.BridgeServiceGrpc.GetInTransitTransferGPList(ctx, &bridgeService.GetInTransitTransferGPListRequest{
		Orddocid: payload.Orddocid,
		// Status:          req.Status,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "in tarnsit transfer")
		return
	}

	if tr, err = s.opt.Client.BridgeServiceGrpc.GetTransferRequestGPList(ctx, &bridgeService.GetTransferRequestGPListRequest{
		Docnumbr: payload.IvmTrNumber,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "transfer request")
		return
	}

	if tr.Data[0].IvmStatus != 4 && itt.Data[0].Status != 1 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorMustActive("transfer request")
		return
	}

	trDetailMap = make(map[string]float64, 0)
	for _, v := range tr.Data[0].Details {
		trDetailMap[v.Itemnmbr] = v.IvmQtyRequest
	}

	for _, detail := range payload.Detail {
		val, ok := trDetailMap[detail.Itemnmbr]
		// If the key exists
		if ok {
			if val < detail.FulfillQty {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorMustGreater("fulfill qty", "request qty")
				return
			}
		}
		details = append(details, &bridgeService.UpdateInTransitTransferGPRequest_Detail{
			Lnitmseq:   int32(detail.Lnitmseq),
			Itemnmbr:   detail.Itemnmbr,
			ReasonCode: detail.ReasonCode,
			Qtyfulfi:   detail.FulfillQty,
		})
	}

	if res, err = s.opt.Client.BridgeServiceGrpc.UpdateInTransitTransferGP(ctx, &bridgeService.UpdateInTransitTransferGPRequest{
		Orddocid:    payload.Orddocid,
		IvmTrNumber: payload.IvmTrNumber,
		Ordrdate:    payload.Ordrdate,
		Etadte:      payload.Etadte,
		Eta:         payload.Etatime,
		Note:        payload.Note,
		Detail:      details,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "in transit transfer")
		return
	}

	userID = ctx.Value(constants.KeyUserID).(int64)
	_, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &auditService.CreateLogRequest{
		Log: &auditService.Log{
			UserId:      userID,
			ReferenceId: payload.Orddocid,
			Type:        "in transit transfer",
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

func (s *ItemTransferService) CommitTransferRequestGP(ctx context.Context, payload *dto.CommitTransferRequestGPRequest) (res *bridgeService.CommitTransferRequestGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "TransferRequestService.CommitTransferRequestGP")
	defer span.End()

	var details []*bridgeService.CommitTransferRequestGPRequest_Detail
	for _, detail := range payload.Detail {
		details = append(details, &bridgeService.CommitTransferRequestGPRequest_Detail{
			Lnitmseq:      int32(detail.Lnitmseq),
			IvmQtyFulfill: detail.FulfilledQty,
		})
	}

	if res, err = s.opt.Client.BridgeServiceGrpc.CommitTransferRequestGP(ctx, &bridgeService.CommitTransferRequestGPRequest{
		Docnumbr: payload.TRNumber,
		Detail:   details,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "transfer request")
		return
	}

	userID := ctx.Value(constants.KeyUserID).(int64)
	_, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &auditService.CreateLogRequest{
		Log: &auditService.Log{
			UserId:      userID,
			ReferenceId: payload.TRNumber,
			Type:        "item transfer",
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
