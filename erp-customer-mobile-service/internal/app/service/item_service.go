package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/utils"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/catalog_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
)

type IItemService interface {
	GetPrivate(ctx context.Context, req dto.RequestGetPrivateItemList) (res []dto.ItemResponse, err error)
	GetPublic(ctx context.Context, req dto.RequestGetItemList) (res []dto.ItemResponse, err error)
	GetPrivateDetail(ctx context.Context, req dto.ItemDetailPrivateRequest) (res dto.ItemResponse, err error)
	GetPublicDetail(ctx context.Context, req dto.ItemDetailRequest) (res dto.ItemResponse, err error)
	GetPrivateItemByListID(ctx context.Context, req *dto.RequestGetPrivateItemByListID) (res []*dto.ItemResponse, err error)
	GetLastFinTrans(ctx context.Context, req dto.RequestGetFinishedItems) (res []dto.LastFinTransItemResponse, err error)
}

type ItemService struct {
	opt opt.Options
	//RepositoryOTPOutgoing repository.IOtpOutgoingRepository
}

func NewItemService() IItemService {
	return &ItemService{
		opt: global.Setup.Common,
		//RepositoryOTPOutgoing: repository.NewOtpOutgoingRepository(),
	}
}

func (s *ItemService) GetPrivate(ctx context.Context, req dto.RequestGetPrivateItemList) (res []dto.ItemResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemService.Get")
	defer span.End()

	layout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Asia/Jakarta")
	currentTime, _ := time.ParseInLocation(layout, time.Now().Format(layout), loc)

	// check Address
	// address, err := s.opt.Client.BridgeServiceGrpc.GetAddressDetail(ctx, &bridge_service.GetAddressDetailRequest{
	// 	Id: req.Data.AddressID,
	// })
	// if err != nil {
	// 	span.RecordError(err)
	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 	err = edenlabs.ErrorRpcNotFound("bridge", "address")
	// 	return
	// }
	// if req.Session.Customer.AddressID != address.Data.Id {
	// 	//error
	// }

	//check Address
	// admDivision, err := s.opt.Client.BridgeServiceGrpc.GetAdmDivisionDetail(ctx, &bridge_service.GetAdmDivisionDetailRequest{
	// 	Id: req.Session.Customer.AdmDivisionID,
	// })
	// if err != nil {
	// 	span.RecordError(err)
	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 	err = edenlabs.ErrorRpcNotFound("bridge", "admDivision")
	// 	return
	// }

	glossary, err := s.opt.Client.ConfigurationServiceGrpc.GetGlossaryDetail(ctx, &configuration_service.GetGlossaryDetailRequest{
		Table:     "sales_order",
		Attribute: "order_channel",
		// ValueItemName: req.Platform,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("configuration", "glossary")
		return
	}
	itemCategoryID, _ := strconv.Atoi(req.Data.ItemCategoryID)

	Items, err := s.opt.Client.CatalogServiceGrpc.GetItemList(ctx, &catalog_service.GetItemListRequest{
		// admDivisionId: admDivision.Data.Id,
		Limit:          int32(req.Limit),
		Offset:         int32(req.Offset),
		Search:         req.Data.Search,
		ItemCategoryId: int64(itemCategoryID),
		// RegionId: admDivision.Data.RegionId,
		// OrderChannel:  glossary.Data.ValueInt,
		// Type:   int32(req.Data.Type),
		Status: 1,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
	}

	for _, Item := range Items.Data {
		var uomTemp *bridge_service.GetUomDetailResponse
		uomTemp, err = s.opt.Client.BridgeServiceGrpc.GetUomDetail(ctx, &bridge_service.GetUomDetailRequest{
			Id: utils.ToInt64(Item.UomId),
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		// var itemImageTemp *catalog_service.GetItemImageDetailResponse
		// itemImageTemp, err = s.opt.Client.CatalogServiceGrpc.GetItemImageDetail(ctx, &catalog_service.GetItemImageDetailRequest{
		// 	ItemId:    Item.Id,
		// 	MainImage: 1,
		// })
		// if err != nil {
		// 	span.RecordError(err)
		// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
		// 	return
		// }

		// default image
		imageUrl := "https://sgp1.digitaloceanspaces.com/image-erp-dev-eden/item/2342554562.jpg"

		if strings.Contains(Item.Description, req.Data.Search) {
			res = append(res, dto.ItemResponse{
				ID:               strconv.Itoa(int(Item.Id)),
				ItemName:         Item.Description,
				ItemUomName:      uomTemp.Data.Description,
				Description:      Item.Note,
				UnitPrice:        "5000",
				OrderMinQty:      "1",
				DecimalEnabled:   "1",
				ImageUrl:         imageUrl,
				ItemCategoryName: Item.ItemCategoryName,
			})
		}
	}

	fmt.Println(currentTime, glossary)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *ItemService) GetPublic(ctx context.Context, req dto.RequestGetItemList) (res []dto.ItemResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemService.Get")
	defer span.End()

	layout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Asia/Jakarta")
	currentTime, _ := time.ParseInLocation(layout, time.Now().Format(layout), loc)

	glossary, err := s.opt.Client.ConfigurationServiceGrpc.GetGlossaryDetail(ctx, &configuration_service.GetGlossaryDetailRequest{
		Table:     "sales_order",
		Attribute: "order_channel",
		// ValueItemName: req.Platform,
	})

	// admDivision, err := s.opt.Client.BridgeServiceGrpc.GetAdmDivisionDetail(ctx, &bridge_service.GetAdmDivisionDetailRequest{
	// 	Id: req.Data.AdmDivisionID,
	// })
	// if err != nil {
	// 	span.RecordError(err)
	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 	err = edenlabs.ErrorRpcNotFound("bridge", "admDivision")
	// 	return
	// }

	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("configuration", "glossary")
		return
	}
	itemCategoryID, _ := strconv.Atoi(req.Data.ItemCategoryID)

	Items, err := s.opt.Client.CatalogServiceGrpc.GetItemList(ctx, &catalog_service.GetItemListRequest{
		// admDivisionId: admDivision.Data.Id,
		Limit:          int32(req.Limit),
		Offset:         int32(req.Offset),
		Search:         req.Data.Search,
		ItemCategoryId: int64(itemCategoryID),
		// RegionId: admDivision.Data.RegionId,
		// OrderChannel:  glossary.Data.ValueInt,
		// Type:   int32(req.Data.Type),
		Status: 1,
	})

	for _, Item := range Items.Data {
		uomTemp, err := s.opt.Client.BridgeServiceGrpc.GetUomDetail(ctx, &bridge_service.GetUomDetailRequest{
			Id: utils.ToInt64(Item.UomId),
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "uom")
		}

		// itemImageTemp, err := s.opt.Client.CatalogServiceGrpc.GetItemImageDetail(ctx, &catalog_service.GetItemImageDetailRequest{
		// 	ItemId:    Item.Id,
		// 	MainImage: 1,
		// })
		// if err != nil {
		// 	span.RecordError(err)
		// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
		// 	err = edenlabs.ErrorRpcNotFound("catalog", "itemimage")
		// }

		// default image
		imageUrl := "https://sgp1.digitaloceanspaces.com/image-erp-dev-eden/item/2342554562.jpg"

		if strings.Contains(Item.Description, req.Data.Search) {
			res = append(res, dto.ItemResponse{
				ID:               strconv.Itoa(int(Item.Id)),
				ItemName:         Item.Description,
				ItemUomName:      uomTemp.Data.Description,
				Description:      Item.Note,
				UnitPrice:        "5000",
				OrderMinQty:      "1",
				DecimalEnabled:   "1",
				ImageUrl:         imageUrl,
				ItemCategoryName: Item.ItemCategoryName,
			})
		}

	}

	fmt.Println(currentTime, glossary)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *ItemService) GetPrivateDetail(ctx context.Context, req dto.ItemDetailPrivateRequest) (res dto.ItemResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemService.Get")
	defer span.End()

	layout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Asia/Jakarta")
	currentTime, _ := time.ParseInLocation(layout, time.Now().Format(layout), loc)
	var itemCategoryNameTemp []string
	var itemImageUrlTemp []string

	// check Address
	// address, err := s.opt.Client.BridgeServiceGrpc.GetAddressDetail(ctx, &bridge_service.GetAddressDetailRequest{
	// 	Id: req.Data.AddressID,
	// })
	// if err != nil {
	// 	span.RecordError(err)
	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 	err = edenlabs.ErrorRpcNotFound("bridge", "address")
	// 	return
	// }
	// if req.Session.Customer.AddressID != address.Data.Id {
	// 	//error
	// }

	//check Address
	// admDivision, err := s.opt.Client.BridgeServiceGrpc.GetAdmDivisionDetail(ctx, &bridge_service.GetAdmDivisionDetailRequest{
	// 	Id: req.Session.Customer.AdmDivisionID,
	// })
	// if err != nil {
	// 	span.RecordError(err)
	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 	err = edenlabs.ErrorRpcNotFound("bridge", "admDivision")
	// 	return
	// }

	glossary, err := s.opt.Client.ConfigurationServiceGrpc.GetGlossaryDetail(ctx, &configuration_service.GetGlossaryDetailRequest{
		Table:     "sales_order",
		Attribute: "order_channel",
		// ValueItemName: req.Platform,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("configuration", "glossary")
		return
	}

	ItemDetail, err := s.opt.Client.CatalogServiceGrpc.GetItemDetailByInternalId(ctx, &catalog_service.GetItemDetailByInternalIdRequest{
		Id: req.Data.ItemID,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("catalog", "item")
		return
	}

	for _, ItemCategoryID := range ItemDetail.Data.ItemCategoryId {
		var itemCategoryTemp *catalog_service.GetItemCategoryDetailResponse
		itemCategoryTemp, err = s.opt.Client.CatalogServiceGrpc.GetItemCategoryDetail(ctx, &catalog_service.GetItemCategoryDetailRequest{
			Id: ItemCategoryID,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("catalog", "item category")
			return
		}

		itemCategoryNameTemp = append(itemCategoryNameTemp, itemCategoryTemp.Data.Name)
	}

	for _, ItemImage := range ItemDetail.Data.ItemImage {
		itemImageUrlTemp = append(itemImageUrlTemp, ItemImage.ImageUrl)
	}

	// dummy data for item category and item image
	itemCategoryNameTemp = append(itemCategoryNameTemp, "Temp Category")
	itemImageUrlTemp = append(itemImageUrlTemp, "https://hips.hearstapps.com/hmg-prod.s3.amazonaws.com/images/cute-cat-photos-1593441022.jpg")

	uomTemp, err := s.opt.Client.BridgeServiceGrpc.GetUomDetail(ctx, &bridge_service.GetUomDetailRequest{
		Id: utils.ToInt64(ItemDetail.Data.UomId),
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "uom")
		return
	}

	res = dto.ItemResponse{
		ID:                  strconv.Itoa(int(ItemDetail.Data.Id)),
		Code:                ItemDetail.Data.Code,
		ItemName:            ItemDetail.Data.Description,
		ItemUomName:         uomTemp.Data.Description,
		Description:         ItemDetail.Data.Note,
		UnitPrice:           "5000",
		OrderMinQty:         "1",
		DecimalEnabled:      "1",
		ItemCategoryNameArr: itemCategoryNameTemp,
		ImagesUrlArr:        itemImageUrlTemp,
	}

	fmt.Println(currentTime, glossary)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *ItemService) GetPublicDetail(ctx context.Context, req dto.ItemDetailRequest) (res dto.ItemResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemService.Get")
	defer span.End()

	layout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Asia/Jakarta")
	currentTime, _ := time.ParseInLocation(layout, time.Now().Format(layout), loc)
	var itemCategoryNameTemp []string
	var itemImageUrlTemp []string

	// check Address
	// address, err := s.opt.Client.BridgeServiceGrpc.GetAddressDetail(ctx, &bridge_service.GetAddressDetailRequest{
	// 	Id: req.Data.AddressID,
	// })
	// if err != nil {
	// 	span.RecordError(err)
	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 	err = edenlabs.ErrorRpcNotFound("bridge", "address")
	// 	return
	// }
	// if req.Session.Customer.AddressID != address.Data.Id {
	// 	//error
	// }

	//check Address
	// admDivision, err := s.opt.Client.BridgeServiceGrpc.GetAdmDivisionDetail(ctx, &bridge_service.GetAdmDivisionDetailRequest{
	// 	Id: req.Session.Customer.AdmDivisionID,
	// })
	// if err != nil {
	// 	span.RecordError(err)
	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 	err = edenlabs.ErrorRpcNotFound("bridge", "admDivision")
	// 	return
	// }

	glossary, err := s.opt.Client.ConfigurationServiceGrpc.GetGlossaryDetail(ctx, &configuration_service.GetGlossaryDetailRequest{
		Table:     "sales_order",
		Attribute: "order_channel",
		// ValueItemName: req.Platform,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("configuration", "glossary")
		return
	}

	ItemDetail, err := s.opt.Client.CatalogServiceGrpc.GetItemDetailByInternalId(ctx, &catalog_service.GetItemDetailByInternalIdRequest{
		Id: utils.ToString(req.Data.ItemID),
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("catalog", "item")
		return
	}

	for _, ItemCategoryID := range ItemDetail.Data.ItemCategoryId {
		var itemCategoryTemp *catalog_service.GetItemCategoryDetailResponse
		itemCategoryTemp, err = s.opt.Client.CatalogServiceGrpc.GetItemCategoryDetail(ctx, &catalog_service.GetItemCategoryDetailRequest{
			Id: ItemCategoryID,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("catalog", "item category")
			return
		}

		itemCategoryNameTemp = append(itemCategoryNameTemp, itemCategoryTemp.Data.Name)
	}

	for _, ItemImage := range ItemDetail.Data.ItemImage {
		itemImageUrlTemp = append(itemImageUrlTemp, ItemImage.ImageUrl)
	}

	// dummy data for item category and item image
	itemCategoryNameTemp = append(itemCategoryNameTemp, "Temp Category")
	itemImageUrlTemp = append(itemImageUrlTemp, "https://hips.hearstapps.com/hmg-prod.s3.amazonaws.com/images/cute-cat-photos-1593441022.jpg")

	uomTemp, err := s.opt.Client.BridgeServiceGrpc.GetUomDetail(ctx, &bridge_service.GetUomDetailRequest{
		Id: utils.ToInt64(ItemDetail.Data.UomId),
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "uom")
	}

	res = dto.ItemResponse{
		ID:                  strconv.Itoa(int(ItemDetail.Data.Id)),
		ItemName:            ItemDetail.Data.Description,
		ItemUomName:         uomTemp.Data.Description,
		Description:         ItemDetail.Data.Note,
		UnitPrice:           "5000",
		OrderMinQty:         "1",
		DecimalEnabled:      "1",
		ItemCategoryNameArr: itemCategoryNameTemp,
		ImagesUrlArr:        itemImageUrlTemp,
	}

	fmt.Println(currentTime, glossary)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *ItemService) GetPrivateItemByListID(ctx context.Context, req *dto.RequestGetPrivateItemByListID) (res []*dto.ItemResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemService.Get")
	defer span.End()

	// check Address
	_, err = s.opt.Client.BridgeServiceGrpc.GetAddressDetail(ctx, &bridge_service.GetAddressDetailRequest{
		Id: utils.ToInt64(req.Data.AddressID),
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "address")
		return
	}

	var orderChannel *configuration_service.GetGlossaryDetailResponse
	orderChannel, err = s.opt.Client.ConfigurationServiceGrpc.GetGlossaryDetail(ctx, &configuration_service.GetGlossaryDetailRequest{
		Table:     "sales_order",
		Attribute: "order_channel",
		ValueName: req.Platform,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "item")
		return
	}

	req.Data.ItemsID = strings.ReplaceAll(req.Data.ItemsID, " ", "")

	for i, item := range utils.StringToStringArray(req.Data.ItemsID) {
		if utils.ToInt64(item) == 0 {
			err = edenlabs.ErrorValidation(fmt.Sprintf("item_id%d", i+1), "Item tidak valid")
			return
		}

		var ItemDetail *catalog_service.GetItemDetailByInternalIdResponse
		ItemDetail, err = s.opt.Client.CatalogServiceGrpc.GetItemDetailByInternalId(ctx, &catalog_service.GetItemDetailByInternalIdRequest{
			Id: item,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("catalog", "item")
			return
		}

		if strings.Contains(ItemDetail.Data.OrderChannelRestriction, orderChannel.Data.Note) {
			continue
		}

		var uom *bridge_service.GetUomDetailResponse
		uom, err = s.opt.Client.BridgeServiceGrpc.GetUomDetail(ctx, &bridge_service.GetUomDetailRequest{
			Id: utils.ToInt64(ItemDetail.Data.UomId),
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "uom")
			return
		}

		// default image
		imageUrl := "https://sgp1.digitaloceanspaces.com/image-erp-dev-eden/item/2342554562.jpg"
		if len(ItemDetail.Data.ItemImage) > 0 {
			imageUrl = ItemDetail.Data.ItemImage[0].ImageUrl
		}
		res = append(res, &dto.ItemResponse{
			ID:             strconv.Itoa(int(ItemDetail.Data.Id)),
			ItemName:       ItemDetail.Data.Description,
			ItemUomName:    uom.Data.Description,
			Description:    ItemDetail.Data.Note,
			UnitPrice:      "5000",
			DecimalEnabled: "1",
			ImageUrl:       imageUrl,
			OrderMinQty:    "1",
		})
	}

	return
}

// function to get items of last finished transactions
func (s *ItemService) GetLastFinTrans(ctx context.Context, req dto.RequestGetFinishedItems) (res []dto.LastFinTransItemResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemService.Get")
	defer span.End()

	var (
		addressID int64
		address   *bridge_service.GetAddressDetailResponse
		soReq     *bridge_service.GetSalesOrderListRequest
		soResp    *bridge_service.GetSalesOrderListResponse
		soiReq    *bridge_service.GetSalesOrderItemListRequest
		soiResp   *bridge_service.GetSalesOrderItemListResponse
		itemImage *catalog_service.GetItemImageDetailResponse
		uom       *bridge_service.GetUomDetailResponse
		itemDet   *bridge_service.GetItemDetailResponse
	)

	// get address from gp
	addressID, err = strconv.ParseInt(req.Data.AddressID, 10, 64)
	address, err = s.opt.Client.BridgeServiceGrpc.GetAddressDetail(ctx, &bridge_service.GetAddressDetailRequest{
		Id: addressID,
	})

	// get last finished sales order of address from gp
	soReq = &bridge_service.GetSalesOrderListRequest{
		AddressId: address.Data.Id,
		Status:    2,
		OrderBy:   "id desc",
		Limit:     1,
	}
	if soResp, err = s.opt.Client.BridgeServiceGrpc.GetSalesOrderList(ctx, soReq); err != nil {
		return nil, edenlabs.ErrorInvalid("sales order")
	}

	// get items of last sales order
	soiReq = &bridge_service.GetSalesOrderItemListRequest{
		SalesOrderId: soResp.Data[0].Id,
	}
	if soiResp, err = s.opt.Client.BridgeServiceGrpc.GetSalesOrderItemList(ctx, soiReq); err != nil {
		return nil, edenlabs.ErrorInvalid("sales order item")
	}

	// loop through sales order item data
	for i, item := range soiResp.Data {
		// get item data
		if itemDet, err = s.opt.Client.BridgeServiceGrpc.GetItemDetail(ctx, &bridge_service.GetItemDetailRequest{
			Id: utils.ToInt64(item.ItemId),
		}); err != nil {
			return nil, edenlabs.ErrorInvalid("item detail")
		}

		// get item image data
		if itemImage, err = s.opt.Client.CatalogServiceGrpc.GetItemImageDetail(ctx, &catalog_service.GetItemImageDetailRequest{
			ItemId:    utils.ToInt64(item.ItemId),
			MainImage: 1,
		}); err != nil {
			return nil, edenlabs.ErrorInvalid("item image")
		}

		// get uom data
		if uom, err = s.opt.Client.BridgeServiceGrpc.GetUomDetail(ctx, &bridge_service.GetUomDetailRequest{
			Id: itemDet.Data.UomId,
		}); err != nil {
			return nil, edenlabs.ErrorInvalid("uom")
		}

		// combine data into response
		res = append(res, dto.LastFinTransItemResponse{
			ItemId:      utils.ToString(314 + i),
			ItemName:    itemDet.Data.Description,
			ItemUomName: uom.Data.Description,
			UnitPrice:   fmt.Sprintf("%f", item.UnitPrice),
			ShadowPrice: fmt.Sprintf("%f", item.ShadowPrice),
			OrderQty:    "1",
			Subtotal:    fmt.Sprintf("%f", item.Subtotal),
			Weight:      fmt.Sprintf("%f", item.Weight),
			Note:        item.Note,
			ImageUrl:    itemImage.Data.ImageUrl,
		})
	}

	return
}
