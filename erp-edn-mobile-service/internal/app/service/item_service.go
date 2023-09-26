package service

import (
	"context"
	"fmt"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/global"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	catalogService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/catalog_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
)

func NewServiceItem() IItemService {
	m := new(ItemService)
	m.opt = global.Setup.Common
	return m
}

type IItemService interface {
	Get(ctx context.Context, req dto.ItemListRequest) (res []*dto.ItemResponse, err error)
	GetProductGT(ctx context.Context, req dto.ItemListRequest) (res []*dto.ItemResponse, err error)
	GetByID(ctx context.Context, req dto.ItemDetailRequest) (res *dto.ItemResponse, err error)
}

type ItemService struct {
	opt opt.Options
}

func (s *ItemService) Get(ctx context.Context, req dto.ItemListRequest) (res []*dto.ItemResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemService.Get")
	defer span.End()

	var regionID, custTypeID, locationCode, priceLevel string
	var customer *bridge_service.GetCustomerGPResponse
	var address *bridge_service.GetAddressGPResponse
	tempAddress := &bridge_service.AddressGP{}

	if req.CustomerID != "" {
		customer, err = s.opt.Client.BridgeServiceGrpc.GetCustomerGPList(ctx, &bridge_service.GetCustomerGPListRequest{
			// Id: customerDetail.Data.CustomerIdGp,
			Limit:          1,
			Offset:         0,
			Id:             req.CustomerID,
			CustomerTypeId: "BTY0015",
		})
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		if len(customer.Data) == 0 {
			err = edenlabs.ErrorRpcNotFound("bridge", "customer")
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		custTypeID = customer.Data[0].CustomerType[0].GnL_Cust_Type_ID
		priceLevel = customer.Data[0].Prclevel
		address, err = s.opt.Client.BridgeServiceGrpc.GetAddressGPList(ctx, &bridge_service.GetAddressGPListRequest{
			// Id: customer.Data[0].Adrscode[0].Adrscode,
			Limit:          100,
			Offset:         0,
			CustomerNumber: customer.Data[0].Custnmbr,
			// Adrscode:       req.AddressID,
		})

		if err != nil {
			err = edenlabs.ErrorRpcNotFound("bridge", "address")
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		if len(address.Data) == 0 {
			err = edenlabs.ErrorRpcNotFound("bridge", "address")
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		for _, v := range address.Data {
			if v.TypeAddress == "ship_to" {
				tempAddress = v
			}
		}
		if tempAddress == nil {
			//span.RecordError(err)
			fmt.Println(address.Data)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return nil, err
		}
		regionID = address.Data[0].AdministrativeDiv.GnlRegion
		locationCode = address.Data[0].Locncode
	}

	orderChannel, err := s.opt.Client.ConfigurationServiceGrpc.GetGlossaryDetail(ctx, &configuration_service.GetGlossaryDetailRequest{
		Table:     "sales_order",
		Attribute: "order_channel",
		ValueName: "edn",
	})
	// get adddress from bridge
	var itemResponse *catalogService.GetItemListResponse
	itemResponse, err = s.opt.Client.CatalogServiceGrpc.GetItemList(ctx, &catalogService.GetItemListRequest{
		Limit:            req.Limit,
		Offset:           req.Offset,
		Status:           req.Status,
		Search:           req.Search,
		OrderBy:          "item_id_gp",
		RegionIdGp:       regionID,
		CustomerTypeIdGp: custTypeID,
		LocationCode:     locationCode,
		PriceLevel:       priceLevel,
		Salability:       1,
		OrderChannel:     orderChannel.Data.ValueInt,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("catalog", "item")
		return
	}

	for _, item := range itemResponse.Data {

		var unitPrice, stock float64
		// var regionID,siteID,priceLevel string
		// unitPrice = 10000
		if len(item.ItemPrice) > 0 {
			unitPrice = item.ItemPrice[0].Price
			priceLevel = item.ItemPrice[0].Pricelevel
			custTypeID = item.ItemPrice[0].CustTypeId
		} else {
			unitPrice = 0
			priceLevel = ""
			custTypeID = ""
		}
		if len(item.ItemSite) > 0 {
			regionID = item.ItemSite[0].RegionId
			locationCode = item.ItemSite[0].SiteId
			stock = item.ItemSite[0].TotalStock
		} else {
			regionID = ""
			locationCode = ""
			stock = 0
		}
		// unitPrice
		res = append(res, &dto.ItemResponse{
			ID:   item.Code,
			Name: item.Description,
			// Code: item.Code,
			Uom: &dto.UomResponse{
				ID:          item.UomId,
				Code:        item.UomId,
				Description: item.UomName,
			},
			Class: &dto.ClassResponse{
				ID: item.ClassId,
			},
			// Description:          item.Description,
			UnitWeightConversion: item.UnitWeightConversion,
			OrderMinQty:          item.OrderMinQty,
			OrderMaxQty:          item.OrderMaxQty,
			ItemType:             item.ItemType,
			Capitalize:           item.Capitalize,
			MaxDayDeliveryDate:   int8(item.MaxDayDeliveryDate),
			Taxable:              item.Taxable,
			Note:                 item.Note,
			Status:               int8(item.Status),
			UnitPrice:            unitPrice,
			RegionID:             regionID,
			SiteID:               locationCode,
			PriceLevel:           priceLevel,
			CustomerTypeID:       custTypeID,
			Stock:                stock,
		})
	}

	return
}

func (s *ItemService) GetProductGT(ctx context.Context, req dto.ItemListRequest) (res []*dto.ItemResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemService.Get")
	defer span.End()

	// get adddress from bridge
	var (
		itemResponse *bridge_service.GetItemMasterComplexGPListResponse
		inactive     string
	)
	if req.Status == 1 {
		inactive = "0"
	} else {
		inactive = "1"
	}
	itemResponse, err = s.opt.Client.BridgeServiceGrpc.GetItemMasterComplexGP(ctx, &bridge_service.GetItemMasterComplexGPListRequest{
		Limit:          req.Limit,
		Offset:         req.Offset,
		Locncode:       req.SiteID,
		Description:    req.Search,
		GnlStorability: "1",
		Inactive:       inactive,
		OrderBy:        req.OrderBy,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("catalog", "item")
		return
	}

	for _, item := range itemResponse.Data {

		var unitPrice float64
		unitPrice = 10000

		res = append(res, &dto.ItemResponse{
			ID:   item.Itemnmbr,
			Name: item.Itemdesc,
			// Code: item.Code,
			Uom: &dto.UomResponse{
				ID:          item.Uomschdl,
				Description: item.Uomschdl,
			},
			Class: &dto.ClassResponse{
				ID: item.Itmclscd,
			},
			// Description:          item.Description,
			// UnitWeightConversion: item.UnitWeightConversion,
			// OrderMinQty:          item.OrderMinQty,
			// OrderMaxQty:          item.OrderMaxQty,
			// ItemType:             item.ItemType,
			// Capitalize:           item.Capitalize,
			// MaxDayDeliveryDate:   int8(item.MaxDayDeliveryDate),
			// Taxable:              item.Taxable,
			// Note:                 item.Note,
			Status:    int8(req.Status),
			UnitPrice: unitPrice,
		})
	}

	return
}

func (s *ItemService) GetByID(ctx context.Context, req dto.ItemDetailRequest) (res *dto.ItemResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemService.GetByID")
	defer span.End()

	var unitPrice float64
	unitPrice = 10000

	// get item from bridge
	var item *catalogService.GetItemDetailResponse
	item, err = s.opt.Client.CatalogServiceGrpc.GetItemDetail(ctx, &catalogService.GetItemDetailRequest{
		Id: fmt.Sprintf("%d", req.Id),
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("catalog", "item")
		return
	}

	res = &dto.ItemResponse{
		// ID:   item.Data.Id,
		ID:   item.Data.Code,
		Name: item.Data.Description,
		Uom: &dto.UomResponse{
			ID: item.Data.UomId,
		},
		Class: &dto.ClassResponse{
			ID: item.Data.ClassId,
		},
		// Description:          item.Data.Description,
		UnitWeightConversion: item.Data.UnitWeightConversion,
		OrderMinQty:          item.Data.OrderMinQty,
		OrderMaxQty:          item.Data.OrderMaxQty,
		ItemType:             item.Data.ItemType,
		Capitalize:           item.Data.Capitalize,
		MaxDayDeliveryDate:   int8(item.Data.MaxDayDeliveryDate),
		Taxable:              item.Data.Taxable,
		Note:                 item.Data.Note,
		Status:               int8(item.Data.Status),
		UnitPrice:            unitPrice,
	}

	return
}
