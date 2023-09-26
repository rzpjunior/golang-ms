package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/constants"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/edenlabs/edenlabs/utils"
	"git.edenfarm.id/project-version3/erp-services/erp-inventory-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-inventory-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-inventory-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-inventory-service/internal/app/repository"
	auditService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/audit_service"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	catalogService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/catalog_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IItemService interface {
	Get(ctx context.Context, req dto.ItemRequestGet) (res []dto.ItemGPResponse, total int64, err error)
	GetByID(ctx context.Context, id string) (res dto.ItemGPResponse, err error)
	UpdateImage(ctx context.Context, req dto.ItemImageRequestUpdate, id int64) (res dto.ItemResponse, err error)
	UpdatePackable(ctx context.Context, req dto.ItemRequestPackable, id int64) (err error)
	UpdateFragile(ctx context.Context, req dto.ItemRequestFragile, id int64) (err error)
	Update(ctx context.Context, req dto.ItemRequestUpdate, id int64) (res dto.ItemResponse, err error)
	GetDetailByInternalID(ctx context.Context, id int64, itemIdGP string) (res dto.ItemGPResponse, err error)
	GetItemDetailMasterComplexByInternalID(ctx context.Context, req *catalogService.GetItemDetailByInternalIdRequest) (res dto.ItemGPResponse, err error)
	GetListItemComplex(ctx context.Context, req *dto.ItemRequestGet) (res []*dto.ItemGPResponse, total int64, err error)
	GetItemListInternal(ctx context.Context, req *dto.ItemRequestGet) (res []*model.Item, err error)
	GetItemDetailInternal(ctx context.Context, req *dto.ItemRequestGet) (res *model.Item, err error)
}

type ItemService struct {
	opt                         opt.Options
	RepositoryItemImage         repository.IItemImageRepository
	RepositoryItemCategory      repository.IItemCategoryRepository
	RepositoryItem              repository.IItemRepository
	RepositoryItemCategoryImage repository.IItemCategoryImageRepository
}

func NewItemService() IItemService {
	return &ItemService{
		opt:                         global.Setup.Common,
		RepositoryItemImage:         repository.NewItemImageRepository(),
		RepositoryItemCategory:      repository.NewItemCategoryRepository(),
		RepositoryItem:              repository.NewItemRepository(),
		RepositoryItemCategoryImage: repository.NewItemCategoryImageRepository(),
	}
}

func (s *ItemService) Get(ctx context.Context, req dto.ItemRequestGet) (res []dto.ItemGPResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemService.Get")
	defer span.End()

	var (
		items  *bridgeService.GetItemGPResponse
		status string
	)

	// Convert status to GP
	if req.Status != 0 {
		if req.Status == 1 {
			status = "0"
		} else {
			status = "1"
		}
	}

	if items, err = s.opt.Client.BridgeServiceGrpc.GetItemGPList(ctx, &bridgeService.GetItemGPListRequest{
		Limit:       int32(req.Limit),
		Offset:      int32(req.Offset),
		Description: req.Search,
		ClassId:     req.ClassID,
		UomId:       req.UomID,
		Inactive:    status,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "item")
		return
	}

	for _, item := range items.Data {
		var (
			detailItem                                                       *model.Item
			orderChannelRestrictionStr, excludeArchetypeStr, itemCategoryStr string
			itemCategoryList                                                 []*dto.ItemCategoryResponse
			excludeArchetypeList                                             []*dto.ArchetypeResponse
			orderChannelRestrictionList                                      []*dto.OrderChannelRestrictionResponse
			uom                                                              *bridgeService.GetUomGPResponse
			status                                                           int8
			itemClass                                                        *bridgeService.GetItemClassGPResponse
		)

		// sync GP
		detailItem = &model.Item{
			ItemIDGP: item.Itemnmbr,
			Name:     item.Itemdesc,
		}
		// try get or create twice before sending errors, as database will reject if there's two same item_id
		// meaning that someone is writing into database at the same time
		if err = s.RepositoryItem.SyncGP(ctx, detailItem); err != nil {
			if err = s.RepositoryItem.SyncGP(ctx, detailItem); err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorInvalid("item sync")
				return
			}
		}

		if detailItem.ItemCategoryID != "" {
			for _, v := range utils.StringToInt64Array(detailItem.ItemCategoryID) {
				var itemCategory *model.ItemCategory
				itemCategory, err = s.RepositoryItemCategory.GetByID(ctx, v)
				if err != nil {
					span.RecordError(err)
					s.opt.Logger.AddMessage(log.ErrorLevel, err)
					return
				}
				ItemCategory := &dto.ItemCategoryResponse{
					ID:       itemCategory.ID,
					RegionID: itemCategory.Regions,
					Name:     itemCategory.Name,
					Status:   itemCategory.Status,
				}
				itemCategoryList = append(itemCategoryList, ItemCategory)
				itemCategoryStr += itemCategory.Name + ", "
			}

			itemCategoryStr = strings.TrimSuffix(itemCategoryStr, ", ")
		}

		if detailItem.ExcludeArchetype != "" {
			// Handling space characters
			detailItem.ExcludeArchetype = strings.ReplaceAll(detailItem.ExcludeArchetype, " ", "")
			for _, v := range utils.StringToStringArray(detailItem.ExcludeArchetype) {
				var archetype *bridgeService.GetArchetypeGPResponse
				archetype, err = s.opt.Client.BridgeServiceGrpc.GetArchetypeGPDetail(ctx, &bridgeService.GetArchetypeGPDetailRequest{
					Id: v,
				})
				if err != nil {
					span.RecordError(err)
					s.opt.Logger.AddMessage(log.ErrorLevel, err)
					err = edenlabs.ErrorInvalid("archetype_id")
					return
				}
				Archetype := &dto.ArchetypeResponse{
					ID:          archetype.Data[0].GnlArchetypeId,
					Description: archetype.Data[0].GnlArchetypedescription,
				}
				excludeArchetypeList = append(excludeArchetypeList, Archetype)
				excludeArchetypeStr += archetype.Data[0].GnlArchetypedescription + ", "
			}
			excludeArchetypeStr = strings.TrimSuffix(excludeArchetypeStr, ", ")
		}

		if detailItem.OrderChannelRestriction != "" {
			for _, v := range utils.StringToInt64Array(detailItem.OrderChannelRestriction) {
				var glossary *configuration_service.GetGlossaryDetailResponse
				glossary, err = s.opt.Client.ConfigurationServiceGrpc.GetGlossaryDetail(ctx, &configuration_service.GetGlossaryDetailRequest{
					ValueInt:  int32(v),
					Table:     "sales_order",
					Attribute: "order_channel",
				})
				if err != nil {
					span.RecordError(err)
					s.opt.Logger.AddMessage(log.ErrorLevel, err)
					return
				}
				OrderChannelRestriction := &dto.OrderChannelRestrictionResponse{
					ID:        int64(glossary.Data.Id),
					Table:     glossary.Data.Table,
					Attribute: glossary.Data.Attribute,
					ValueInt:  int8(glossary.Data.ValueInt),
					ValueName: glossary.Data.ValueName,
					Note:      glossary.Data.Note,
				}
				orderChannelRestrictionList = append(orderChannelRestrictionList, OrderChannelRestriction)
				orderChannelRestrictionStr += glossary.Data.Note + ", "
			}

			orderChannelRestrictionStr = strings.TrimSuffix(orderChannelRestrictionStr, ", ")
		}

		if uom, err = s.opt.Client.BridgeServiceGrpc.GetUomGPDetail(ctx, &bridgeService.GetUomGPDetailRequest{
			Id: item.Uomschdl,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "uom")
			return
		}

		if itemClass, err = s.opt.Client.BridgeServiceGrpc.GetItemClassGPDetail(ctx, &bridgeService.GetItemClassGPDetailRequest{
			Id: item.Itmclscd,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "item class")
			return
		}

		if item.Inactive == 0 {
			status = statusx.ConvertStatusName(statusx.Active)
		} else {
			status = statusx.ConvertStatusName(statusx.Archived)
		}
		var fragile bool
		if detailItem.FragileGoods == "fragile" || detailItem.FragileGoods == "1" {
			fragile = true
		} else {
			fragile = false
		}

		var packable bool
		if detailItem.Packability == "packable" || detailItem.Packability == "1" {
			packable = true
		} else {
			packable = false
		}

		if req.ItemCategoryID != 0 {
			itemCategoryListID := utils.StringToInt64Array(detailItem.ItemCategoryID)
			for _, v := range itemCategoryListID {
				if req.ItemCategoryID == v {
					res = append(res, dto.ItemGPResponse{
						ID:   detailItem.ID,
						Code: item.Itemnmbr,
						Uom: &dto.UomGPResponse{
							ID:   uom.Data[0].Uomschdl,
							Name: uom.Data[0].Umschdsc,
						},
						ClassID:          itemClass.Data[0].Itmclscd,
						ItemCategory:     itemCategoryList,
						ItemCategoryName: itemCategoryStr,
						Description:      item.Itmgedsc,
						// UnitWeightConversion:        item.UnitWeightConversion,
						// OrderMinQty:                 item.OrderMinQty,
						// OrderMaxQty:                 item.OrderMaxQty,
						ItemType: item.ItemTypeDesc,
						Packable: packable,
						Fragile:  fragile,
						// Capitalize:                  item.Capitalize,
						Note:                 detailItem.Note,
						ExcludeArchetypeName: excludeArchetypeStr,
						MaxDayDeliveryDate:   detailItem.MaxDayDeliveryDate,
						// Taxable:                     item.Taxable,
						OrderChannelRestrictionName: orderChannelRestrictionStr,
						Status:                      status,
						ExcludeArchetypes:           excludeArchetypeList,
						OrderChannelRestrictions:    orderChannelRestrictionList,
						Class: &dto.ItemClassResponse{
							ID:   itemClass.Data[0].Itmclscd,
							Name: itemClass.Data[0].Itmclsdc,
						},
					})
				}
			}
		} else {
			res = append(res, dto.ItemGPResponse{
				ID:   detailItem.ID,
				Code: item.Itemnmbr,
				Uom: &dto.UomGPResponse{
					ID:   uom.Data[0].Uomschdl,
					Name: uom.Data[0].Umschdsc,
				},
				ClassID:              itemClass.Data[0].Itmclscd,
				ItemCategory:         itemCategoryList,
				ItemCategoryName:     itemCategoryStr,
				Description:          item.Itemdesc,
				UnitWeightConversion: item.Itemshwt,
				OrderMinQty:          item.Minorqty,
				OrderMaxQty:          item.Maxordqty,
				ItemType:             item.ItemTypeDesc,
				Packable:             true,
				Fragile:              true,
				Capitalize:           item.GnlCbCapitalitemDesc,
				Note:                 detailItem.Note,
				ExcludeArchetypeName: excludeArchetypeStr,
				MaxDayDeliveryDate:   detailItem.MaxDayDeliveryDate,
				// Taxable:                     item.Taxable,
				OrderChannelRestrictionName: orderChannelRestrictionStr,
				Status:                      status,
				ExcludeArchetypes:           excludeArchetypeList,
				OrderChannelRestrictions:    orderChannelRestrictionList,
				Class: &dto.ItemClassResponse{
					ID:   itemClass.Data[0].Itmclscd,
					Name: itemClass.Data[0].Itmclsdc,
				},
			})
		}

	}

	total = int64(items.TotalRecords)

	return
}

func (s *ItemService) GetByID(ctx context.Context, id string) (res dto.ItemGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemService.GetByID")
	defer span.End()

	var (
		item                        *bridgeService.GetItemGPResponse
		detailItem                  *model.Item
		itemCategoryList            []*dto.ItemCategoryResponse
		excludeArchetypeList        []*dto.ArchetypeResponse
		orderChannelRestrictionList []*dto.OrderChannelRestrictionResponse
		uom                         *bridgeService.GetUomGPResponse
		status                      int8
		itemClass                   *bridgeService.GetItemClassGPResponse
	)

	fmt.Println("ID Load", id)
	if item, err = s.opt.Client.BridgeServiceGrpc.GetItemGPDetail(ctx, &bridgeService.GetItemGPDetailRequest{
		Id: id,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "item")
		return
	}

	// sync GP
	detailItem = &model.Item{
		ItemIDGP: item.Data[0].Itemnmbr,
	}
	// try get or create twice before sending errors, as database will reject if there's two same item_id
	// meaning that someone is writing into database at the same time
	if err = s.RepositoryItem.SyncGP(ctx, detailItem); err != nil {
		if err = s.RepositoryItem.SyncGP(ctx, detailItem); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("item sync")
			return
		}
	}

	var itemImages []*model.ItemImage
	itemImages, _ = s.RepositoryItemImage.GetByItemID(ctx, detailItem.ID)

	var orderChannelRestrictionStr, excludeArchetypeStr, itemCategoryStr string

	if detailItem.ItemCategoryID != "" {
		for _, v := range utils.StringToInt64Array(detailItem.ItemCategoryID) {
			var itemCategory *model.ItemCategory
			itemCategory, err = s.RepositoryItemCategory.GetByID(ctx, v)
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}

			itemCategoryStr += itemCategory.Name + ", "

			var region string
			listIDArrStr := strings.Split(itemCategory.Regions, ",")
			listID := utils.ArrayStringToInt64Array(listIDArrStr)

			for _, v := range listID {
				var detailRegion *bridgeService.GetRegionDetailResponse
				detailRegion, err = s.opt.Client.BridgeServiceGrpc.GetRegionDetail(ctx, &bridgeService.GetRegionDetailRequest{
					Id: v,
				})
				if err != nil {
					span.RecordError(err)
					s.opt.Logger.AddMessage(log.ErrorLevel, err)
					err = edenlabs.ErrorInvalid("region_id")
					return
				}
				region += detailRegion.Data.Description + ","
			}

			region = strings.TrimSuffix(region, ",")

			var itemCategoryImage *model.ItemCategoryImage
			itemCategoryImage, err = s.RepositoryItemCategoryImage.GetByItemCategoryID(ctx, itemCategory.ID)
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}

			dtoItemCategoryImage := &dto.ItemCategoryImageResponse{
				ID:             itemCategoryImage.ID,
				ItemCategoryID: itemCategoryImage.ItemCategoryID,
				ImageUrl:       itemCategoryImage.ImageUrl,
				CreatedAt:      itemCategoryImage.CreatedAt,
			}

			itemCategoryList = append(itemCategoryList, &dto.ItemCategoryResponse{
				ID:                itemCategory.ID,
				Code:              itemCategory.Code,
				Name:              itemCategory.Name,
				RegionID:          itemCategory.Regions,
				Region:            region,
				Status:            itemCategory.Status,
				ItemCategoryImage: dtoItemCategoryImage,
				StatusConvert:     statusx.ConvertStatusValue(itemCategory.Status),
				CreatedAt:         itemCategory.CreatedAt,
				UpdatedAt:         itemCategory.UpdatedAt,
			})
		}

		itemCategoryStr = strings.TrimSuffix(itemCategoryStr, ", ")

	}

	if detailItem.ExcludeArchetype != "" {
		// Handling space characters
		detailItem.ExcludeArchetype = strings.ReplaceAll(detailItem.ExcludeArchetype, " ", "")
		for _, v := range utils.StringToStringArray(detailItem.ExcludeArchetype) {
			var archetype *bridgeService.GetArchetypeGPResponse
			archetype, err = s.opt.Client.BridgeServiceGrpc.GetArchetypeGPDetail(ctx, &bridgeService.GetArchetypeGPDetailRequest{
				Id: v,
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}
			Archetype := &dto.ArchetypeResponse{
				ID:          archetype.Data[0].GnlArchetypeId,
				Description: archetype.Data[0].GnlArchetypedescription,
			}
			excludeArchetypeList = append(excludeArchetypeList, Archetype)
			excludeArchetypeStr += archetype.Data[0].GnlArchetypedescription + ", "
		}
		excludeArchetypeStr = strings.TrimSuffix(excludeArchetypeStr, ", ")
	}

	if detailItem.OrderChannelRestriction != "" {
		for _, v := range utils.StringToInt64Array(detailItem.OrderChannelRestriction) {
			var glossary *configuration_service.GetGlossaryDetailResponse
			glossary, err = s.opt.Client.ConfigurationServiceGrpc.GetGlossaryDetail(ctx, &configuration_service.GetGlossaryDetailRequest{
				ValueInt:  int32(v),
				Table:     "sales_order",
				Attribute: "order_channel",
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}
			OrderChannelRestriction := &dto.OrderChannelRestrictionResponse{
				ID:        int64(glossary.Data.Id),
				Table:     glossary.Data.Table,
				Attribute: glossary.Data.Attribute,
				ValueInt:  int8(glossary.Data.ValueInt),
				ValueName: glossary.Data.ValueName,
				Note:      glossary.Data.Note,
			}
			orderChannelRestrictionList = append(orderChannelRestrictionList, OrderChannelRestriction)
			orderChannelRestrictionStr += glossary.Data.Note + ", "
		}

		orderChannelRestrictionStr = strings.TrimSuffix(orderChannelRestrictionStr, ", ")
	}

	if uom, err = s.opt.Client.BridgeServiceGrpc.GetUomGPDetail(ctx, &bridgeService.GetUomGPDetailRequest{
		Id: item.Data[0].Uomschdl,
	}); err != nil || !uom.Succeeded {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "uom")
		return
	}

	if itemClass, err = s.opt.Client.BridgeServiceGrpc.GetItemClassGPDetail(ctx, &bridgeService.GetItemClassGPDetailRequest{
		Id: item.Data[0].Itmclscd,
	}); err != nil || len(itemClass.Data) == 0 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "class")
		return
	}
	var itemImagesResponse []*dto.ItemImageResponse
	for _, itemImage := range itemImages {
		itemImagesResponse = append(itemImagesResponse, &dto.ItemImageResponse{
			ID:        itemImage.ID,
			ItemID:    itemImage.ItemID,
			ImageUrl:  itemImage.ImageUrl,
			MainImage: itemImage.MainImage,
			CreatedAt: itemImage.CreatedAt,
		})
	}

	var fragile bool
	if detailItem.FragileGoods == "fragile" {
		fragile = true
	} else {
		fragile = false
	}

	var packable bool
	if detailItem.Packability == "packable" {
		packable = true
	} else {
		packable = false
	}

	if item.Data[0].Inactive == 0 {
		status = statusx.ConvertStatusName(statusx.Active)
	} else {
		status = statusx.ConvertStatusName(statusx.Archived)
	}

	res = dto.ItemGPResponse{
		ID:   detailItem.ID,
		Code: item.Data[0].Itemnmbr,
		Uom: &dto.UomGPResponse{
			ID:   uom.Data[0].Uomschdl,
			Name: uom.Data[0].Umschdsc,
		},
		ClassID:              itemClass.Data[0].Itmclscd,
		ItemCategory:         itemCategoryList,
		ItemCategoryName:     itemCategoryStr,
		Description:          item.Data[0].Itemdesc,
		UnitWeightConversion: item.Data[0].Itemshwt,
		OrderMinQty:          item.Data[0].Minorqty,
		OrderMaxQty:          item.Data[0].Maxordqty,
		ItemType:             item.Data[0].ItemTypeDesc,
		Packable:             packable,
		Fragile:              fragile,
		Capitalize:           item.Data[0].GnlCbCapitalitemDesc,
		Note:                 detailItem.Note,
		ExcludeArchetypeName: excludeArchetypeStr,
		MaxDayDeliveryDate:   detailItem.MaxDayDeliveryDate,
		// Taxable:                     item.Data.Taxable,
		OrderChannelRestrictionName: orderChannelRestrictionStr,
		Status:                      status,
		ItemImages:                  itemImagesResponse,
		Class: &dto.ItemClassResponse{
			ID:   itemClass.Data[0].Itmclscd,
			Name: itemClass.Data[0].Itmclsdc,
		},
		OrderChannelRestrictions: orderChannelRestrictionList,
		ExcludeArchetypes:        excludeArchetypeList,
	}

	return
}

func (s *ItemService) UpdateImage(ctx context.Context, req dto.ItemImageRequestUpdate, id int64) (res dto.ItemResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemService.UpdateImage")
	defer span.End()

	// validate item id
	_, err = s.opt.Client.BridgeServiceGrpc.GetItemDetail(ctx, &bridgeService.GetItemDetailRequest{
		Id: id,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "item")
		return
	}

	err = s.RepositoryItemImage.DeleteByItemID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	var itemImagesResponse []*dto.ItemImageResponse
	for i, image := range req.Images {
		var isMainImage int8
		if i == 0 {
			isMainImage = 1
		} else {
			isMainImage = 2
		}

		itemImage := &model.ItemImage{
			ItemID:    id,
			ImageUrl:  image.ImageUrl,
			MainImage: isMainImage,
			CreatedAt: time.Now(),
		}

		err = s.RepositoryItemImage.Create(ctx, itemImage)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		itemImagesResponse = append(itemImagesResponse, &dto.ItemImageResponse{
			ID:        itemImage.ID,
			ItemID:    itemImage.ItemID,
			ImageUrl:  itemImage.ImageUrl,
			MainImage: itemImage.MainImage,
			CreatedAt: itemImage.CreatedAt,
		})
	}

	res = dto.ItemResponse{
		ID:         id,
		ItemImages: itemImagesResponse,
	}

	return
}

func (s *ItemService) UpdatePackable(ctx context.Context, req dto.ItemRequestPackable, id int64) (err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemService.UpdatePackable")
	defer span.End()

	// validate item id
	var item *bridgeService.GetItemDetailResponse
	item, err = s.opt.Client.BridgeServiceGrpc.GetItemDetail(ctx, &bridgeService.GetItemDetailRequest{
		Id: id,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "item")
		return
	}

	var detailItem *model.Item
	detailItem, err = s.RepositoryItem.GetDetail(ctx, item.Data.Id, "", 0, "", 0)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("item")
		return
	}

	var packability string
	if req.Packable {
		packability = "1"
	} else {
		packability = "0"
	}

	detailItem.Packability = packability

	err = s.RepositoryItem.Update(ctx, detailItem, "Packability")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *ItemService) UpdateFragile(ctx context.Context, req dto.ItemRequestFragile, id int64) (err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemService.UpdateFragile")
	defer span.End()

	// validate item id
	var item *bridgeService.GetItemDetailResponse
	item, err = s.opt.Client.BridgeServiceGrpc.GetItemDetail(ctx, &bridgeService.GetItemDetailRequest{
		Id: id,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "item")
		return
	}

	var detailItem *model.Item
	detailItem, err = s.RepositoryItem.GetDetail(ctx, item.Data.Id, "", 0, "", 0)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("item")
		return
	}

	var fragile string
	if req.Fragile {
		fragile = "fragile"
	} else {
		fragile = "non fragile"
	}

	detailItem.FragileGoods = fragile

	err = s.RepositoryItem.Update(ctx, detailItem, "FragileGoods")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *ItemService) Update(ctx context.Context, req dto.ItemRequestUpdate, id int64) (res dto.ItemResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemService.Update")
	defer span.End()

	var (
		itemBridge                  *bridgeService.GetItemGPResponse
		detailItem                  *model.Item
		itemCategoryList            []*dto.ItemCategoryResponse
		excludeArchetypeList        []*dto.ArchetypeResponse
		orderChannelRestrictionList []*dto.OrderChannelRestrictionResponse
	)
	detailItem, err = s.RepositoryItem.GetDetail(ctx, id, "", 0, "", 0)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("item")
		return
	}

	// validate item id
	itemBridge, err = s.opt.Client.BridgeServiceGrpc.GetItemGPDetail(ctx, &bridgeService.GetItemGPDetailRequest{
		Id: detailItem.ItemIDGP,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "item")
		return
	}

	var orderChannelRestrictionStr, excludeArchetypeStr, itemCategoryStr string

	if len(req.ItemCategory) != 0 {
		for _, v := range req.ItemCategory {
			var itemCategory *model.ItemCategory
			itemCategory, err = s.RepositoryItemCategory.GetByID(ctx, v)
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorInvalid("item_category")
				return
			}
			ItemCategory := &dto.ItemCategoryResponse{
				ID:       itemCategory.ID,
				RegionID: itemCategory.Regions,
				Name:     itemCategory.Name,
				Status:   itemCategory.Status,
			}
			itemCategoryList = append(itemCategoryList, ItemCategory)
			itemCategoryStr += itemCategory.Name + ", "
		}

		itemCategoryStr = strings.TrimSuffix(itemCategoryStr, ", ")
	}

	if len(req.ExcludeArchetype) != 0 {
		for _, v := range req.ExcludeArchetype {
			var archetype *bridgeService.GetArchetypeGPResponse
			archetype, err = s.opt.Client.BridgeServiceGrpc.GetArchetypeGPDetail(ctx, &bridgeService.GetArchetypeGPDetailRequest{
				Id: v,
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}
			Archetype := &dto.ArchetypeResponse{
				ID:          archetype.Data[0].GnlArchetypeId,
				Description: archetype.Data[0].GnlArchetypedescription,
			}
			excludeArchetypeList = append(excludeArchetypeList, Archetype)
			excludeArchetypeStr += archetype.Data[0].GnlArchetypedescription + ", "
		}
		excludeArchetypeStr = strings.TrimSuffix(excludeArchetypeStr, ", ")
	}

	if len(req.OrderChannelRestriction) != 0 {
		for _, v := range req.OrderChannelRestriction {
			var glossary *configuration_service.GetGlossaryDetailResponse
			glossary, err = s.opt.Client.ConfigurationServiceGrpc.GetGlossaryDetail(ctx, &configuration_service.GetGlossaryDetailRequest{
				ValueInt:  int32(v),
				Table:     "sales_order",
				Attribute: "order_channel",
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("configuration", "glossary")
				return
			}
			OrderChannelRestriction := &dto.OrderChannelRestrictionResponse{
				ID:        int64(glossary.Data.Id),
				Table:     glossary.Data.Table,
				Attribute: glossary.Data.Attribute,
				ValueInt:  int8(glossary.Data.ValueInt),
				ValueName: glossary.Data.ValueName,
				Note:      glossary.Data.Note,
			}
			orderChannelRestrictionList = append(orderChannelRestrictionList, OrderChannelRestriction)
			orderChannelRestrictionStr += glossary.Data.Note + ", "
		}

		orderChannelRestrictionStr = strings.TrimSuffix(orderChannelRestrictionStr, ", ")
	}

	// Validation note length
	if len(req.Note) > 500 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorMustEqualOrLess("note", "500 characters")
		return
	}

	detailItem.ItemCategoryID = utils.ArrayInt64ToString(req.ItemCategory)
	detailItem.MaxDayDeliveryDate = req.MaxDayDeliveryDate
	detailItem.ExcludeArchetype = utils.ArrayStringToString(req.ExcludeArchetype)
	detailItem.OrderChannelRestriction = utils.ArrayInt64ToString(req.OrderChannelRestriction)
	detailItem.Note = req.Note

	err = s.RepositoryItem.Update(ctx, detailItem, "ItemCategoryID", "MaxDayDeliveryDate", "ExcludeArchetype", "OrderChannelRestriction", "Note")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// Temporary takeout handling err
	_ = s.RepositoryItemImage.DeleteByItemID(ctx, id)
	// if err != nil {
	// 	span.RecordError(err)
	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 	return
	// }

	var itemImagesResponse []*dto.ItemImageResponse
	for i, image := range req.Images {
		var isMainImage int8
		if i == 0 {
			isMainImage = 1
		} else {
			isMainImage = 2
		}

		itemImage := &model.ItemImage{
			ItemID:    id,
			ImageUrl:  image,
			MainImage: isMainImage,
			CreatedAt: time.Now(),
		}

		err = s.RepositoryItemImage.Create(ctx, itemImage)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		itemImagesResponse = append(itemImagesResponse, &dto.ItemImageResponse{
			ID:        itemImage.ID,
			ItemID:    itemImage.ItemID,
			ImageUrl:  itemImage.ImageUrl,
			MainImage: itemImage.MainImage,
			CreatedAt: itemImage.CreatedAt,
		})
	}

	userID := ctx.Value(constants.KeyUserID).(int64)

	_, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &auditService.CreateLogRequest{
		Log: &auditService.Log{
			UserId:      userID,
			ReferenceId: strconv.Itoa(int(id)),
			Type:        "item",
			Function:    "update",
			CreatedAt:   timestamppb.New(time.Now()),
		},
	})

	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	var fragile bool
	if detailItem.FragileGoods == "fragile" || detailItem.FragileGoods == "1" {
		fragile = true
	} else {
		fragile = false
	}

	var packable bool
	if detailItem.Packability == "packable" || detailItem.Packability == "1" {
		packable = true
	} else {
		packable = false
	}

	res = dto.ItemResponse{
		ID: detailItem.ID,
		Uom: &dto.UomGPResponse{
			ID:   itemBridge.Data[0].Uomschdl,
			Name: "kg",
		},
		ItemCategory:                itemCategoryList,
		Code:                        itemBridge.Data[0].Itemnmbr,
		Description:                 itemBridge.Data[0].Itemdesc,
		UnitWeightConversion:        itemBridge.Data[0].Itemshwt,
		OrderMinQty:                 itemBridge.Data[0].Minorqty,
		OrderMaxQty:                 itemBridge.Data[0].Maxordqty,
		ItemType:                    itemBridge.Data[0].ItemTypeDesc,
		Packable:                    packable,
		Fragile:                     fragile,
		Capitalize:                  itemBridge.Data[0].GnlCbCapitalitemDesc,
		Note:                        detailItem.Note,
		ExcludeArchetypeName:        excludeArchetypeStr,
		ExcludeArchetypes:           excludeArchetypeList,
		MaxDayDeliveryDate:          int8(detailItem.MaxDayDeliveryDate),
		OrderChannelRestrictionName: orderChannelRestrictionStr,
		OrderChannelRestrictions:    orderChannelRestrictionList,
		ItemImages:                  itemImagesResponse,
		ItemCategoryName:            itemCategoryStr,
	}

	return
}

func (s *ItemService) GetDetailByInternalID(ctx context.Context, id int64, itemIdGP string) (res dto.ItemGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemService.GetDetailByInternalID")
	defer span.End()

	var (
		item                        *bridgeService.GetItemGPResponse
		detailItem                  *model.Item
		itemCategoryList            []*dto.ItemCategoryResponse
		excludeArchetypeList        []*dto.ArchetypeResponse
		orderChannelRestrictionList []*dto.OrderChannelRestrictionResponse
		uom                         *bridgeService.GetUomGPResponse
		status                      int8
		itemClass                   *bridgeService.GetItemClassGPResponse
	)

	detailItem, err = s.RepositoryItem.GetDetail(ctx, id, itemIdGP, 0, "", 0)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("item")
		return
	}

	if item, err = s.opt.Client.BridgeServiceGrpc.GetItemGPDetail(ctx, &bridgeService.GetItemGPDetailRequest{
		Id: detailItem.ItemIDGP,
	}); err != nil || !item.Succeeded {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "item")
		return
	}

	var itemImages []*model.ItemImage
	itemImages, _ = s.RepositoryItemImage.GetByItemID(ctx, detailItem.ID)

	var orderChannelRestrictionStr, excludeArchetypeStr, itemCategoryStr string

	if detailItem.ItemCategoryID != "" {
		for _, v := range utils.StringToInt64Array(detailItem.ItemCategoryID) {
			var itemCategory *model.ItemCategory
			itemCategory, err = s.RepositoryItemCategory.GetByID(ctx, v)
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}

			itemCategoryStr += itemCategory.Name + ", "

			var region string
			listIDArrStr := strings.Split(itemCategory.Regions, ",")
			listID := utils.ArrayStringToInt64Array(listIDArrStr)

			for _, v := range listID {
				var detailRegion *bridgeService.GetRegionDetailResponse
				detailRegion, err = s.opt.Client.BridgeServiceGrpc.GetRegionDetail(ctx, &bridgeService.GetRegionDetailRequest{
					Id: v,
				})
				if err != nil {
					span.RecordError(err)
					s.opt.Logger.AddMessage(log.ErrorLevel, err)
					err = edenlabs.ErrorInvalid("region_id")
					return
				}
				region += detailRegion.Data.Description + ","
			}

			region = strings.TrimSuffix(region, ",")

			var itemCategoryImage *model.ItemCategoryImage
			itemCategoryImage, err = s.RepositoryItemCategoryImage.GetByItemCategoryID(ctx, itemCategory.ID)
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}

			dtoItemCategoryImage := &dto.ItemCategoryImageResponse{
				ID:             itemCategoryImage.ID,
				ItemCategoryID: itemCategoryImage.ItemCategoryID,
				ImageUrl:       itemCategoryImage.ImageUrl,
				CreatedAt:      itemCategoryImage.CreatedAt,
			}

			itemCategoryList = append(itemCategoryList, &dto.ItemCategoryResponse{
				ID:                itemCategory.ID,
				Code:              itemCategory.Code,
				Name:              itemCategory.Name,
				RegionID:          itemCategory.Regions,
				Region:            region,
				Status:            itemCategory.Status,
				ItemCategoryImage: dtoItemCategoryImage,
				StatusConvert:     statusx.ConvertStatusValue(itemCategory.Status),
				CreatedAt:         itemCategory.CreatedAt,
				UpdatedAt:         itemCategory.UpdatedAt,
			})
		}

		itemCategoryStr = strings.TrimSuffix(itemCategoryStr, ", ")

	}

	if detailItem.ExcludeArchetype != "" {
		// Handling space characters
		detailItem.ExcludeArchetype = strings.ReplaceAll(detailItem.ExcludeArchetype, " ", "")
		for _, v := range utils.StringToStringArray(detailItem.ExcludeArchetype) {
			var archetype *bridgeService.GetArchetypeGPResponse
			archetype, err = s.opt.Client.BridgeServiceGrpc.GetArchetypeGPDetail(ctx, &bridgeService.GetArchetypeGPDetailRequest{
				Id: v,
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}
			Archetype := &dto.ArchetypeResponse{
				ID:          archetype.Data[0].GnlArchetypeId,
				Description: archetype.Data[0].GnlArchetypedescription,
				Code:        archetype.Data[0].GnlArchetypeId,
				CustomerType: &dto.CustomerTypeResponse{
					ID:          archetype.Data[0].GnlCustTypeId,
					Description: archetype.Data[0].GnlCusttypeDescription,
				},
			}
			excludeArchetypeList = append(excludeArchetypeList, Archetype)
			excludeArchetypeStr += archetype.Data[0].GnlArchetypedescription + ", "
		}
		excludeArchetypeStr = strings.TrimSuffix(excludeArchetypeStr, ", ")
	}

	if detailItem.OrderChannelRestriction != "" {
		for _, v := range utils.StringToInt64Array(detailItem.OrderChannelRestriction) {
			var glossary *configuration_service.GetGlossaryDetailResponse
			glossary, err = s.opt.Client.ConfigurationServiceGrpc.GetGlossaryDetail(ctx, &configuration_service.GetGlossaryDetailRequest{
				ValueInt:  int32(v),
				Table:     "sales_order",
				Attribute: "order_channel",
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}
			OrderChannelRestriction := &dto.OrderChannelRestrictionResponse{
				ID:        int64(glossary.Data.Id),
				Table:     glossary.Data.Table,
				Attribute: glossary.Data.Attribute,
				ValueInt:  int8(glossary.Data.ValueInt),
				ValueName: glossary.Data.ValueName,
				Note:      glossary.Data.Note,
			}
			orderChannelRestrictionList = append(orderChannelRestrictionList, OrderChannelRestriction)
			orderChannelRestrictionStr += glossary.Data.Note + ", "
		}

		orderChannelRestrictionStr = strings.TrimSuffix(orderChannelRestrictionStr, ", ")
	}

	if uom, err = s.opt.Client.BridgeServiceGrpc.GetUomGPDetail(ctx, &bridgeService.GetUomGPDetailRequest{
		Id: item.Data[0].Uomschdl,
	}); err != nil || !uom.Succeeded || len(uom.Data) == 0 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "uom")
		return
	}

	if itemClass, err = s.opt.Client.BridgeServiceGrpc.GetItemClassGPDetail(ctx, &bridgeService.GetItemClassGPDetailRequest{
		Id: item.Data[0].Itmclscd,
	}); err != nil || len(itemClass.Data) == 0 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "item class")
		return
	}

	var itemImagesResponse []*dto.ItemImageResponse
	for _, itemImage := range itemImages {
		itemImagesResponse = append(itemImagesResponse, &dto.ItemImageResponse{
			ID:        itemImage.ID,
			ItemID:    itemImage.ItemID,
			ImageUrl:  itemImage.ImageUrl,
			MainImage: itemImage.MainImage,
			CreatedAt: itemImage.CreatedAt,
		})
	}

	var fragile bool
	if detailItem.FragileGoods == "fragile" || detailItem.FragileGoods == "1" {
		fragile = true
	} else {
		fragile = false
	}

	var packable bool
	if detailItem.Packability == "packable" || detailItem.Packability == "1" {
		packable = true
	} else {
		packable = false
	}

	if item.Data[0].Inactive == 0 {
		status = statusx.ConvertStatusName(statusx.Active)
	} else {
		status = statusx.ConvertStatusName(statusx.Archived)
	}

	res = dto.ItemGPResponse{
		ID:   detailItem.ID,
		Code: item.Data[0].Itemnmbr,
		Uom: &dto.UomGPResponse{
			ID:   uom.Data[0].Uomschdl,
			Name: uom.Data[0].Umschdsc,
		},
		ClassID:              itemClass.Data[0].Itmclscd,
		ItemCategory:         itemCategoryList,
		ItemCategoryName:     itemCategoryStr,
		Description:          item.Data[0].Itemdesc,
		UnitWeightConversion: item.Data[0].Itemshwt,
		OrderMinQty:          item.Data[0].Minorqty,
		OrderMaxQty:          item.Data[0].Maxordqty,
		ItemType:             item.Data[0].ItemTypeDesc,
		Packable:             packable,
		Fragile:              fragile,
		Packability:          detailItem.Packability,
		Capitalize:           item.Data[0].GnlCbCapitalitemDesc,
		Note:                 detailItem.Note,
		ExcludeArchetypeName: excludeArchetypeStr,
		MaxDayDeliveryDate:   detailItem.MaxDayDeliveryDate,
		// Taxable:                     item.Data[0].Tax,
		OrderChannelRestrictionName: orderChannelRestrictionStr,
		Status:                      status,
		ItemImages:                  itemImagesResponse,
		Class: &dto.ItemClassResponse{
			ID:   itemClass.Data[0].Itmclscd,
			Name: itemClass.Data[0].Itmclsdc,
		},
		OrderChannelRestrictions: orderChannelRestrictionList,
		ExcludeArchetypes:        excludeArchetypeList,
	}

	return
}

func (s *ItemService) GetItemDetailMasterComplexByInternalID(ctx context.Context, req *catalogService.GetItemDetailByInternalIdRequest) (res dto.ItemGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemService.GetDetailByInternalID")
	defer span.End()

	var (
		item                     *bridgeService.GetItemMasterComplexGPListResponse
		detailItem               *model.Item
		uom                      *bridgeService.GetUomGPResponse
		status                   int8
		itemClass                *bridgeService.GetItemClassGPResponse
		itemPriceTieringResponse []*dto.ItemPriceTieringResponse
		itemPriceResponse        []*dto.ItemPriceResponse
		itemSite                 []*dto.ItemSiteResponse
		itemPrice                float64
	)
	id := utils.ToInt64(req.Id)
	detailItem, err = s.RepositoryItem.GetDetail(ctx, id, req.ItemIdGp, 0, "", 0)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("item")
		return
	}

	// check availability item by archetype and order channel
	if req.ArchetypeIdGp != "" && strings.Contains(detailItem.ExcludeArchetype, req.ArchetypeIdGp) {
		err = edenlabs.ErrorValidation("item", "The item not available for archetype")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if req.OrderChannel != 0 && strings.Contains(detailItem.OrderChannelRestriction, utils.ToString(req.OrderChannel)) {
		err = edenlabs.ErrorValidation("item", "The item not available for order channel")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	//ke bridge buat ambil detail
	if item, err = s.opt.Client.BridgeServiceGrpc.GetItemMasterComplexGP(ctx, &bridgeService.GetItemMasterComplexGPListRequest{
		Limit:         1,
		Offset:        0,
		ItemNumber:    detailItem.ItemIDGP,
		GnlRegion:     req.RegionIdGp,
		Inactive:      "0",
		Locncode:      req.LocationCode,
		GnlSalability: utils.ToString(req.Salability),
		Prclevel:      req.PriceLevel,
	}); err != nil || !item.Succeeded {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "item")
		return
	}
	if len(item.Data) == 0 {
		err = edenlabs.ErrorRpcNotFound("bridge", "item")
		return
	}
	var itemImages []*model.ItemImage
	itemImages, _ = s.RepositoryItemImage.GetByItemID(ctx, detailItem.ID)

	var itemCategoryStr string

	if detailItem.ItemCategoryID != "" {
		for _, v := range utils.StringToInt64Array(detailItem.ItemCategoryID) {
			var itemCategory *model.ItemCategory
			itemCategory, err = s.RepositoryItemCategory.GetByID(ctx, v)
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}
			itemCategoryStr += itemCategory.Name + ", "
		}
		itemCategoryStr = strings.TrimSuffix(itemCategoryStr, ", ")
	}

	if uom, err = s.opt.Client.BridgeServiceGrpc.GetUomGPDetail(ctx, &bridgeService.GetUomGPDetailRequest{
		Id: item.Data[0].Uomschdl,
	}); err != nil || !uom.Succeeded || len(uom.Data) == 0 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "uom")
		return
	}

	if itemClass, err = s.opt.Client.BridgeServiceGrpc.GetItemClassGPDetail(ctx, &bridgeService.GetItemClassGPDetailRequest{
		Id: item.Data[0].Itmclscd,
	}); err != nil || len(itemClass.Data) == 0 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "item class")
		return
	}

	var itemImagesResponse []*dto.ItemImageResponse
	for _, itemImage := range itemImages {
		itemImagesResponse = append(itemImagesResponse, &dto.ItemImageResponse{
			ID:        itemImage.ID,
			ItemID:    itemImage.ItemID,
			ImageUrl:  itemImage.ImageUrl,
			MainImage: itemImage.MainImage,
			CreatedAt: itemImage.CreatedAt,
		})
	}

	var fragile bool
	if detailItem.FragileGoods == "fragile" {
		fragile = true
	} else {
		fragile = false
	}

	var packable bool
	if detailItem.Packability == "packable" {
		packable = true
	} else {
		packable = false
	}

	if item.Data[0].Inactive == 0 {
		status = statusx.ConvertStatusName(statusx.Active)
	} else {
		status = statusx.ConvertStatusName(statusx.Archived)
	}

	if len(item.Data[0].PriceLevel) > 0 {
		for _, v := range item.Data[0].PriceLevel {
			itemPriceResponse = append(itemPriceResponse, &dto.ItemPriceResponse{
				Region:       v.GnlRegion,
				CustomerType: v.GnlCustTypeId,
				PriceLevel:   v.Prclevel,
				Price:        v.Price,
			})
		}
		itemPrice = item.Data[0].PriceLevel[0].Price
	}

	if len(item.Data[0].PriceTiering) > 0 {
		for _, v := range item.Data[0].PriceTiering {
			itemPriceTieringResponse = append(itemPriceTieringResponse, &dto.ItemPriceTieringResponse{
				Docnumbr:          v.Docnumbr,
				GnlRegion:         v.GnlRegion,
				EffectiveDate:     v.EffectiveDate,
				GnlMinQty:         v.GnlMinQty,
				GnlDiscountAmount: v.GnlDiscountAmount,
				GnlQuotaUser:      v.GnlQuotaUser,
			})
		}
	}

	if len(item.Data[0].Site) > 0 {
		for _, v := range item.Data[0].Site {
			itemSite = append(itemSite, &dto.ItemSiteResponse{
				Region:              v.GnlRegion,
				Location:            v.Locncode,
				GnlCbSalability:     v.GnlCbSalability,
				GnlCbSalabilityDesc: v.GnlCbSalabilityDesc,
				TotalStock:          v.TotalStock,
			})
		}
	}
	res = dto.ItemGPResponse{
		ID:   detailItem.ID,
		Code: item.Data[0].Itemnmbr,
		Uom: &dto.UomGPResponse{
			ID:   uom.Data[0].Uomschdl,
			Name: uom.Data[0].Umschdsc,
		},
		ClassID:              itemClass.Data[0].Itmclscd,
		ItemCategoryName:     itemCategoryStr,
		Description:          item.Data[0].Itemdesc,
		UnitWeightConversion: item.Data[0].Itemshwt,
		OrderMinQty:          item.Data[0].Minorqty,
		OrderMaxQty:          item.Data[0].Maxordqty,
		// ItemType:             item.Data[0].ItemTypeDesc,
		Packable: fragile,
		Fragile:  packable,
		// Capitalize:           item.Data[0].GnlCbCapitalitemDesc,
		Note:               detailItem.Note,
		MaxDayDeliveryDate: detailItem.MaxDayDeliveryDate,
		// Taxable:                     item.Data[0].Tax,
		Status:     status,
		ItemImages: itemImagesResponse,
		Class: &dto.ItemClassResponse{
			ID:   itemClass.Data[0].Itmclscd,
			Name: itemClass.Data[0].Itmclsdc,
		},
		Price:            itemPrice,
		ItemPrice:        itemPriceResponse,
		ItemPriceTiering: itemPriceTieringResponse,
		ItemSite:         itemSite,
	}

	return
}

func (s *ItemService) GetListItemComplex(ctx context.Context, req *dto.ItemRequestGet) (res []*dto.ItemGPResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemService.Get")
	defer span.End()

	var (
		statusGP   string
		itemList   []*model.Item
		itemDetail *model.Item
	)

	// Convert status to GP
	if req.Status != 0 {
		if req.Status == 1 {
			statusGP = "0"
		} else {
			statusGP = "1"
		}
	}

	if req.Search == "" {
		tempOffset := req.Offset
		req.Offset = (req.Offset * req.Limit)
		itemList, total, err = s.RepositoryItem.Get(ctx, req)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		req.Offset = tempOffset

		for _, item := range itemList {
			var (
				orderChannelRestrictionStr, excludeArchetypeStr, itemCategoryStr string
				itemCategoryList                                                 []*dto.ItemCategoryResponse
				excludeArchetypeList                                             []*dto.ArchetypeResponse
				orderChannelRestrictionList                                      []*dto.OrderChannelRestrictionResponse
				uom                                                              *bridgeService.GetUomGPResponse
				itemClass                                                        *bridgeService.GetItemClassGPResponse
				status                                                           int8
				itemPriceTieringResponse                                         []*dto.ItemPriceTieringResponse
				itemPriceResponse                                                []*dto.ItemPriceResponse
				itemSiteResponse                                                 []*dto.ItemSiteResponse
				isEnableDecimal                                                  bool
				itemGP                                                           *bridgeService.GetItemMasterComplexGPListResponse
				itemPrice                                                        float64
			)
			if itemGP, err = s.opt.Client.BridgeServiceGrpc.GetItemMasterComplexGP(ctx, &bridgeService.GetItemMasterComplexGPListRequest{
				Limit:  int32(req.Limit),
				Offset: 0,
				// Description:   req.Search,
				Inactive:      statusGP,
				Locncode:      req.SiteIDGP,
				GnlRegion:     req.RegionIDGP,
				GnlSalability: req.Salability,
				ItemNumber:    item.ItemIDGP,
				Prclevel:      req.PriceLevel,
			}); err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "item")
				return
			}

			// Handling the item id not in gp
			if itemGP.Data == nil {
				continue
			}

			if uom, err = s.opt.Client.BridgeServiceGrpc.GetUomGPDetail(ctx, &bridgeService.GetUomGPDetailRequest{
				Id: itemGP.Data[0].Uomschdl,
			}); err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "uom")
				return
			}

			if uom.Data[0].Umdpqtys == 3 {
				isEnableDecimal = true
			}

			if itemClass, err = s.opt.Client.BridgeServiceGrpc.GetItemClassGPDetail(ctx, &bridgeService.GetItemClassGPDetailRequest{
				Id: itemGP.Data[0].Itmclscd,
			}); err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "item class")
				return
			}

			if item.ItemCategoryID != "" {
				for _, v := range utils.StringToInt64Array(item.ItemCategoryID) {
					var itemCategory *model.ItemCategory
					itemCategory, err = s.RepositoryItemCategory.GetByID(ctx, v)
					if err != nil {
						span.RecordError(err)
						s.opt.Logger.AddMessage(log.ErrorLevel, err)
						return
					}
					itemCategoryStr += itemCategory.Name + ", "
				}
				itemCategoryStr = strings.TrimSuffix(itemCategoryStr, ", ")
			}

			if len(itemGP.Data[0].PriceLevel) > 0 {
				for _, v := range itemGP.Data[0].PriceLevel {
					itemPriceResponse = append(itemPriceResponse, &dto.ItemPriceResponse{
						Region:       v.GnlRegion,
						CustomerType: v.GnlCustTypeId,
						PriceLevel:   v.Prclevel,
						Price:        v.Price,
					})
				}
				itemPrice = itemGP.Data[0].PriceLevel[0].Price
			}

			if len(itemGP.Data[0].PriceTiering) > 0 {
				for _, v := range itemGP.Data[0].PriceTiering {
					itemPriceTieringResponse = append(itemPriceTieringResponse, &dto.ItemPriceTieringResponse{
						Docnumbr:          v.Docnumbr,
						GnlRegion:         v.GnlRegion,
						EffectiveDate:     v.EffectiveDate,
						GnlMinQty:         v.GnlMinQty,
						GnlDiscountAmount: v.GnlDiscountAmount,
						GnlQuotaUser:      v.GnlQuotaUser,
					})
				}
			}

			if len(itemGP.Data[0].Site) > 0 {
				for _, v := range itemGP.Data[0].Site {
					itemSiteResponse = append(itemSiteResponse, &dto.ItemSiteResponse{
						Region:              v.GnlRegion,
						Location:            v.Locncode,
						GnlCbSalability:     v.GnlCbSalability,
						GnlCbSalabilityDesc: v.GnlCbSalabilityDesc,
						TotalStock:          v.TotalStock,
					})
				}
			}

			// Convert status GP to internal
			if itemGP.Data[0].Inactive == 0 {
				status = statusx.ConvertStatusName(statusx.Active)
			} else {
				status = statusx.ConvertStatusName(statusx.Archived)
			}

			res = append(res, &dto.ItemGPResponse{
				ID:   item.ID,
				Code: itemGP.Data[0].Itemnmbr,
				Uom: &dto.UomGPResponse{
					ID:   uom.Data[0].Uomschdl,
					Name: uom.Data[0].Umschdsc,
				},
				ClassID:              itemClass.Data[0].Itmclscd,
				ItemCategory:         itemCategoryList,
				ItemCategoryName:     itemCategoryStr,
				Description:          itemGP.Data[0].Itemdesc,
				UnitWeightConversion: itemGP.Data[0].Itemshwt,
				OrderMinQty:          itemGP.Data[0].Minorqty,
				OrderMaxQty:          itemGP.Data[0].Maxordqty,
				// ItemType: item.ItemTypeDesc,
				Packable: true,
				Fragile:  true,
				// Capitalize:                  item.Capitalize,
				Note:                 item.Note,
				ExcludeArchetypeName: excludeArchetypeStr,
				MaxDayDeliveryDate:   item.MaxDayDeliveryDate,
				// Taxable:                     item.Taxable,
				OrderChannelRestrictionName: orderChannelRestrictionStr,
				Status:                      status,
				ExcludeArchetypes:           excludeArchetypeList,
				OrderChannelRestrictions:    orderChannelRestrictionList,
				Class: &dto.ItemClassResponse{
					ID:   itemClass.Data[0].Itmclscd,
					Name: itemClass.Data[0].Itmclsdc,
				},
				Price:            itemPrice,
				ItemPrice:        itemPriceResponse,
				ItemPriceTiering: itemPriceTieringResponse,
				DecimalEnabled:   isEnableDecimal,
				ItemSite:         itemSiteResponse,
			})

		}
	} else {
		var (
			excludeArchetypeList        []*dto.ArchetypeResponse
			orderChannelRestrictionList []*dto.OrderChannelRestrictionResponse
			itemCategoryStr             string
			itemCategoryList            []*dto.ItemCategoryResponse
			uom                         *bridgeService.GetUomGPResponse
			itemClass                   *bridgeService.GetItemClassGPResponse
			status                      int8
			itemPriceTieringResponse    []*dto.ItemPriceTieringResponse
			itemPriceResponse           []*dto.ItemPriceResponse
			isEnableDecimal             bool
			// itemGP                                                           *bridgeService.GetItemMasterComplexGPListResponse
		)

		itemGP, _ := s.opt.Client.BridgeServiceGrpc.GetItemMasterComplexGP(ctx, &bridgeService.GetItemMasterComplexGPListRequest{
			Limit:         int32(req.Limit),
			Offset:        int32(req.Offset),
			Description:   req.Search,
			Inactive:      statusGP,
			Locncode:      req.SiteIDGP,
			GnlRegion:     req.RegionIDGP,
			GnlSalability: req.Salability,
			Prclevel:      req.PriceLevel,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "item")
			return
		}

		for _, v := range itemGP.Data {
			var itemPrice float64
			itemDetail, err = s.RepositoryItem.GetDetail(ctx, 0, v.Itemnmbr, req.ItemCategoryID, "", 0)
			if err != nil {
				err = nil
				continue
			}

			uom, _ = s.opt.Client.BridgeServiceGrpc.GetUomGPDetail(ctx, &bridgeService.GetUomGPDetailRequest{
				Id: v.Uomschdl,
			})
			// if err != nil {
			// 	span.RecordError(err)
			// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
			// 	err = edenlabs.ErrorRpcNotFound("bridge", "uom")
			// 	return
			// }

			if uom.Data[0].Umdpqtys == 3 {
				isEnableDecimal = true
			}

			if itemClass, err = s.opt.Client.BridgeServiceGrpc.GetItemClassGPDetail(ctx, &bridgeService.GetItemClassGPDetailRequest{
				Id: v.Itmclscd,
			}); err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "item class")
				return
			}

			if itemDetail.ItemCategoryID != "" {
				for _, v := range utils.StringToInt64Array(itemDetail.ItemCategoryID) {
					var itemCategory *model.ItemCategory
					itemCategory, err = s.RepositoryItemCategory.GetByID(ctx, v)
					if err != nil {
						span.RecordError(err)
						s.opt.Logger.AddMessage(log.ErrorLevel, err)
						return
					}
					itemCategoryStr += itemCategory.Name + ", "
				}
				itemCategoryStr = strings.TrimSuffix(itemCategoryStr, ", ")
			}

			if len(v.PriceLevel) > 0 {
				for _, v := range v.PriceLevel {
					itemPriceResponse = append(itemPriceResponse, &dto.ItemPriceResponse{
						Region:       v.GnlRegion,
						CustomerType: v.GnlCustTypeId,
						PriceLevel:   v.Prclevel,
						Price:        v.Price,
					})
				}
				itemPrice = v.PriceLevel[0].Price
			}

			if len(v.PriceTiering) > 0 {
				for _, v := range v.PriceTiering {
					itemPriceTieringResponse = append(itemPriceTieringResponse, &dto.ItemPriceTieringResponse{
						Docnumbr:          v.Docnumbr,
						GnlRegion:         v.GnlRegion,
						EffectiveDate:     v.EffectiveDate,
						GnlMinQty:         v.GnlMinQty,
						GnlDiscountAmount: v.GnlDiscountAmount,
						GnlQuotaUser:      v.GnlQuotaUser,
					})
				}
			}

			// Convert status GP to internal
			if itemGP.Data[0].Inactive == 0 {
				status = statusx.ConvertStatusName(statusx.Active)
			} else {
				status = statusx.ConvertStatusName(statusx.Archived)
			}
			res = append(res, &dto.ItemGPResponse{
				ID:   itemDetail.ID,
				Code: v.Itemnmbr,
				Uom: &dto.UomGPResponse{
					ID:   uom.Data[0].Uomschdl,
					Name: uom.Data[0].Umschdsc,
				},
				ClassID:              itemClass.Data[0].Itmclscd,
				ItemCategory:         itemCategoryList,
				ItemCategoryName:     itemCategoryStr,
				Description:          v.Itemdesc,
				UnitWeightConversion: v.Itemshwt,
				OrderMinQty:          v.Minorqty,
				OrderMaxQty:          v.Maxordqty,
				// ItemType: item.ItemTypeDesc,
				Packable: true,
				Fragile:  true,
				// Capitalize:                  item.Capitalize,
				Note:                 itemDetail.Note,
				ExcludeArchetypeName: "",
				MaxDayDeliveryDate:   itemDetail.MaxDayDeliveryDate,
				// Taxable:                     item.Taxable,
				OrderChannelRestrictionName: "",
				Status:                      status,
				ExcludeArchetypes:           excludeArchetypeList,
				OrderChannelRestrictions:    orderChannelRestrictionList,
				Class: &dto.ItemClassResponse{
					ID:   itemClass.Data[0].Itmclscd,
					Name: itemClass.Data[0].Itmclsdc,
				},
				Price:            itemPrice,
				ItemPrice:        itemPriceResponse,
				ItemPriceTiering: itemPriceTieringResponse,
				DecimalEnabled:   isEnableDecimal,
			})

		}

	}
	total = int64(len(res))

	return
}

func (s *ItemService) GetItemListInternal(ctx context.Context, req *dto.ItemRequestGet) (res []*model.Item, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemService.Get")
	defer span.End()

	res, _, err = s.RepositoryItem.Get(ctx, req)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, item := range res {
		item.ItemImage, err = s.RepositoryItemImage.GetByItemID(ctx, item.ID)
	}

	return
}

func (s *ItemService) GetItemDetailInternal(ctx context.Context, req *dto.ItemRequestGet) (res *model.Item, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemService.Get")
	defer span.End()

	var (
		items             []*model.Item
		itemCategoryNames []string
	)

	items, _, err = s.RepositoryItem.Get(ctx, req)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	items[0].ItemImage, err = s.RepositoryItemImage.GetByItemID(ctx, items[0].ID)

	itemCategoryNames = strings.Split(items[0].ItemCategoryID, ",")
	for i, itemCategoryId := range itemCategoryNames {
		itemCategoryIdInt := utils.ToInt64(itemCategoryId)
		itemCategory, _ := s.RepositoryItemCategory.GetByID(ctx, itemCategoryIdInt)
		itemCategoryNames[i] = itemCategory.Name
	}

	items[0].ItemCategoryNameArr = itemCategoryNames

	return items[0], err
}
