package handler

import (
	context "context"
	"strings"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/utils"
	dto "git.edenfarm.id/project-version3/erp-services/erp-inventory-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-inventory-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/catalog_service"
	catalogService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/catalog_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *CatalogGrpcHandler) GetItemList(ctx context.Context, req *catalogService.GetItemListRequest) (res *catalogService.GetItemListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetItemList")
	defer span.End()

	var items []*dto.ItemGPResponse
	items, _, err = h.ServicesItem.GetListItemComplex(ctx, &dto.ItemRequestGet{
		Offset:           int(req.Offset),
		Limit:            int(req.Limit),
		Status:           int(req.Status),
		Search:           req.Search,
		OrderBy:          req.OrderBy,
		UomID:            req.UomId,
		ItemCategoryID:   req.ItemCategoryId,
		SiteIDGP:         req.LocationCode,
		CustomerTypeIDGP: req.CustomerTypeIdGp,
		RegionIDGP:       req.RegionIdGp,
		Salability:       utils.ToString(req.Salability),
		OrderChannel:     req.OrderChannel,
		ArchetypeIDGP:    req.ArchetypeIdGp,
		PriceLevel:       req.PriceLevel,
	})
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*catalogService.Item
	for _, item := range items {

		itemResponse := &catalogService.Item{
			Id:                      item.ID,
			Code:                    item.Code,
			UomId:                   item.Uom.ID,
			ClassId:                 item.ClassID,
			Description:             item.Description,
			UnitWeightConversion:    item.UnitWeightConversion,
			OrderMinQty:             item.OrderMinQty,
			OrderMaxQty:             item.OrderMaxQty,
			ItemType:                item.ItemType,
			Capitalize:              item.Capitalize,
			ExcludeArchetype:        item.ExcludeArchetypeName,
			MaxDayDeliveryDate:      int32(item.MaxDayDeliveryDate),
			FragileGoods:            item.Fragile,
			Taxable:                 item.Taxable,
			OrderChannelRestriction: item.OrderChannelRestrictionName,
			Note:                    item.Note,
			Status:                  int32(item.Status),
			ItemCategoryName:        item.ItemCategoryName,
			UomName:                 item.Uom.Name,
			ClassName:               item.Class.Name,
			Price:                   item.Price,
			DecimalEnabled:          item.DecimalEnabled,
		}

		for _, itemPrice := range item.ItemPrice {
			itemResponse.ItemPrice = append(itemResponse.ItemPrice, &catalogService.Item_PriceLevel{
				RegionId:   itemPrice.Region,
				CustTypeId: itemPrice.CustomerType,
				Pricelevel: itemPrice.PriceLevel,
				Price:      itemPrice.Price,
			})
		}
		for _, itemPriceTiering := range item.ItemPriceTiering {
			itemResponse.PriceTiering = append(itemResponse.PriceTiering, &catalogService.Item_PriceTiering{
				Docnumbr:          itemPriceTiering.Docnumbr,
				GnlRegion:         itemPriceTiering.GnlRegion,
				EffectiveDate:     itemPriceTiering.EffectiveDate,
				GnlMinQty:         itemPriceTiering.GnlMinQty,
				GnlDiscountAmount: itemPriceTiering.GnlDiscountAmount,
				GnlQuotaUser:      itemPriceTiering.GnlQuotaUser,
			})
		}

		for _, itemSite := range item.ItemSite {
			itemResponse.ItemSite = append(itemResponse.ItemSite, &catalogService.Item_Site{
				RegionId:            itemSite.Region,
				SiteId:              itemSite.Location,
				GnlCbSalability:     itemSite.GnlCbSalability,
				GnlCbSalabilityDesc: itemSite.GnlCbSalabilityDesc,
				TotalStock:          itemSite.TotalStock,
			})
		}

		data = append(data, itemResponse)

	}

	res = &catalogService.GetItemListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *CatalogGrpcHandler) GetItemDetail(ctx context.Context, req *catalogService.GetItemDetailRequest) (res *catalogService.GetItemDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetItemDetail")
	defer span.End()

	var item dto.ItemGPResponse

	item, err = h.ServicesItem.GetByID(ctx, req.Id)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	data := &catalogService.Item{
		Id:                      item.ID,
		Code:                    item.Code,
		UomId:                   item.Uom.ID,
		ClassId:                 item.ClassID,
		Description:             item.Description,
		UnitWeightConversion:    item.UnitWeightConversion,
		OrderMinQty:             item.OrderMinQty,
		OrderMaxQty:             item.OrderMaxQty,
		ItemType:                item.ItemType,
		Capitalize:              item.Capitalize,
		ExcludeArchetype:        item.ExcludeArchetypeName,
		MaxDayDeliveryDate:      int32(item.MaxDayDeliveryDate),
		FragileGoods:            item.Fragile,
		Taxable:                 item.Taxable,
		OrderChannelRestriction: item.OrderChannelRestrictionName,
		Note:                    item.Note,
		Status:                  int32(item.Status),
		UomName:                 item.Uom.Name,
		ClassName:               item.Class.Name,
		ItemCategoryName:        item.ItemCategoryName,
		Price:                   item.Price,
	}
	for _, itemImage := range item.ItemImages {
		data.ItemImage = append(data.ItemImage, &catalogService.ItemImage{
			ItemId:    itemImage.ItemID,
			ImageUrl:  itemImage.ImageUrl,
			MainImage: int32(itemImage.MainImage),
		})
	}

	res = &catalogService.GetItemDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *CatalogGrpcHandler) GetItemDetailByInternalId(ctx context.Context, req *catalogService.GetItemDetailByInternalIdRequest) (res *catalogService.GetItemDetailByInternalIdResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetItemDetailByInternalId")
	defer span.End()

	var item dto.ItemGPResponse

	item, err = h.ServicesItem.GetDetailByInternalID(ctx, utils.ToInt64(req.Id), req.ItemIdGp)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	data := &catalogService.Item{
		Id:                      item.ID,
		Code:                    item.Code,
		UomId:                   item.Uom.ID,
		ClassId:                 item.ClassID,
		Description:             item.Description,
		UnitWeightConversion:    item.UnitWeightConversion,
		OrderMinQty:             item.OrderMinQty,
		OrderMaxQty:             item.OrderMaxQty,
		ItemType:                item.ItemType,
		Capitalize:              item.Capitalize,
		ExcludeArchetype:        item.ExcludeArchetypeName,
		MaxDayDeliveryDate:      int32(item.MaxDayDeliveryDate),
		Packability:             item.Packability,
		FragileGoods:            item.Fragile,
		Taxable:                 item.Taxable,
		OrderChannelRestriction: item.OrderChannelRestrictionName,
		Note:                    item.Note,
		Status:                  int32(item.Status),
		UomName:                 item.Uom.Name,
		ClassName:               item.Class.Name,
		ItemCategoryName:        item.ItemCategoryName,
		Price:                   item.Price,
	}
	for _, itemImage := range item.ItemImages {
		data.ItemImage = append(data.ItemImage, &catalogService.ItemImage{
			ItemId:    itemImage.ItemID,
			ImageUrl:  itemImage.ImageUrl,
			MainImage: int32(itemImage.MainImage),
		})
	}

	res = &catalogService.GetItemDetailByInternalIdResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *CatalogGrpcHandler) GetItemDetailMasterComplexByInternalID(ctx context.Context, req *catalogService.GetItemDetailByInternalIdRequest) (res *catalogService.GetItemDetailByInternalIdResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetItemDetailByInternalId")
	defer span.End()

	var item dto.ItemGPResponse

	item, err = h.ServicesItem.GetItemDetailMasterComplexByInternalID(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	data := &catalogService.Item{
		Id:                      item.ID,
		Code:                    item.Code,
		UomId:                   item.Uom.ID,
		ClassId:                 item.ClassID,
		Description:             item.Description,
		UnitWeightConversion:    item.UnitWeightConversion,
		OrderMinQty:             item.OrderMinQty,
		OrderMaxQty:             item.OrderMaxQty,
		ItemType:                item.ItemType,
		Capitalize:              item.Capitalize,
		ExcludeArchetype:        item.ExcludeArchetypeName,
		MaxDayDeliveryDate:      int32(item.MaxDayDeliveryDate),
		FragileGoods:            item.Fragile,
		Taxable:                 item.Taxable,
		OrderChannelRestriction: item.OrderChannelRestrictionName,
		Note:                    item.Note,
		Status:                  int32(item.Status),
		UomName:                 item.Uom.Name,
		ClassName:               item.Class.Name,
		ItemCategoryName:        item.ItemCategoryName,
		Price:                   item.Price,
	}
	for _, itemImage := range item.ItemImages {
		data.ItemImage = append(data.ItemImage, &catalogService.ItemImage{
			ItemId:    itemImage.ItemID,
			ImageUrl:  itemImage.ImageUrl,
			MainImage: int32(itemImage.MainImage),
		})
	}
	for _, itemPrice := range item.ItemPrice {
		data.ItemPrice = append(data.ItemPrice, &catalogService.Item_PriceLevel{
			RegionId:   itemPrice.Region,
			CustTypeId: itemPrice.CustomerType,
			Pricelevel: itemPrice.PriceLevel,
			Price:      itemPrice.Price,
		})
	}
	for _, itemSite := range item.ItemSite {
		data.ItemSite = append(data.ItemSite, &catalogService.Item_Site{
			RegionId:            itemSite.Region,
			SiteId:              itemSite.Location,
			GnlCbSalability:     itemSite.GnlCbSalability,
			GnlCbSalabilityDesc: itemSite.GnlCbSalabilityDesc,
			TotalStock:          itemSite.TotalStock,
		})
	}

	for _, itemPriceTiering := range item.ItemPriceTiering {
		data.PriceTiering = append(data.PriceTiering, &catalogService.Item_PriceTiering{
			Docnumbr:          itemPriceTiering.Docnumbr,
			GnlRegion:         itemPriceTiering.GnlRegion,
			EffectiveDate:     itemPriceTiering.EffectiveDate,
			GnlMinQty:         itemPriceTiering.GnlMinQty,
			GnlDiscountAmount: itemPriceTiering.GnlDiscountAmount,
			GnlQuotaUser:      itemPriceTiering.GnlQuotaUser,
		})
	}

	res = &catalogService.GetItemDetailByInternalIdResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

// get item list of internal
func (h *CatalogGrpcHandler) GetItemListInternal(ctx context.Context, req *catalogService.GetItemListRequest) (res *catalogService.GetItemListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetItemListInternal")
	defer span.End()

	var itemIntern []*model.Item

	itemIntern, err = h.ServicesItem.GetItemListInternal(ctx, &dto.ItemRequestGet{
		Offset:         int(req.Offset),
		Limit:          int(req.Limit),
		Search:         req.Search,
		ItemCategoryID: utils.ToInt64(req.ItemCategoryId),
		OrderChannel:   req.OrderChannel,
		ArchetypeIDGP:  req.ArchetypeIdGp,
		OrderBy:        req.OrderBy,
		ID:             req.ItemId,
	})
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = new(catalogService.GetItemListResponse)
	for _, item := range itemIntern {
		itemResponse := &catalogService.Item{
			Id:             item.ID,
			ItemIdGp:       item.ItemIDGP,
			Note:           item.Note,
			ItemCategoryId: utils.ArrayStringToInt64Array(strings.Split(item.ItemCategoryID, ",")),
		}

		for _, itemImage := range item.ItemImage {
			itemResponse.ItemImage = append(itemResponse.ItemImage, &catalog_service.ItemImage{
				Id:        itemImage.ID,
				ItemId:    itemImage.ItemID,
				ImageUrl:  itemImage.ImageUrl,
				MainImage: int32(itemImage.MainImage),
			})
		}

		res.Data = append(res.Data, itemResponse)
	}

	return
}

func (h *CatalogGrpcHandler) GetItemDetailInternal(ctx context.Context, req *catalogService.GetItemDetailRequest) (res *catalogService.GetItemDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetItemListInternal")
	defer span.End()

	var (
		itemIntern *model.Item
		itemImages []*catalogService.ItemImage
	)

	itemIntern, err = h.ServicesItem.GetItemDetailInternal(ctx, &dto.ItemRequestGet{
		ID:            req.Id,
		OrderChannel:  req.OrderChannel,
		ArchetypeIDGP: req.ArchetypeIdGp,
	})
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	for _, itemImage := range itemIntern.ItemImage {
		image := &catalogService.ItemImage{
			Id:        itemImage.ID,
			ItemId:    itemImage.ItemID,
			ImageUrl:  itemImage.ImageUrl,
			MainImage: int32(itemImage.MainImage),
		}

		itemImages = append(itemImages, image)
	}

	res = new(catalogService.GetItemDetailResponse)
	res.Data = &catalogService.Item{
		Id:                  itemIntern.ID,
		ItemIdGp:            itemIntern.ItemIDGP,
		ItemCategoryNameArr: itemIntern.ItemCategoryNameArr,
		ItemImage:           itemImages,
	}

	return
}
