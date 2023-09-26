package service

import (
	"context"
	"strconv"
	"strings"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/utils"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/campaign_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/catalog_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IItemSectionService interface {
	GetPrivate(ctx context.Context, req dto.RequestGetPrivateItemSection) (res []*dto.ItemSectionResponse, err error)
	GetPublic(ctx context.Context, req dto.RequestGetItemSection) (res []*dto.ItemSectionResponse, err error)
	GetPublicDetail(ctx context.Context, req dto.RequestGetDetailItemSection) (res *dto.ItemSectionResponse, err error)
	GetPrivateDetail(ctx context.Context, req dto.RequestGetPrivateItemSectionDetail) (res *dto.ItemSectionResponse, err error)
}

type ItemSectionService struct {
	opt opt.Options
	//RepositoryOTPOutgoing repository.IOtpOutgoingRepository
}

func NewItemSectionService() IItemSectionService {
	return &ItemSectionService{
		opt: global.Setup.Common,
		//RepositoryOTPOutgoing: repository.NewOtpOutgoingRepository(),
	}
}

func (s *ItemSectionService) GetPrivate(ctx context.Context, req dto.RequestGetPrivateItemSection) (res []*dto.ItemSectionResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemSectionService.GetPrivate")
	defer span.End()

	//check Address
	address, err := s.opt.Client.BridgeServiceGrpc.GetAddressGPDetail(ctx, &bridge_service.GetAddressGPDetailRequest{
		Id: req.Data.AddressID,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("address_id", "address id tidak valid")
		return
	}

	//check Address
	admDivision, err := s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPList(ctx, &bridge_service.GetAdmDivisionGPListRequest{
		AdmDivisionCode: address.Data[0].AdministrativeDiv.GnlAdministrativeCode,
		Limit:           1,
		Offset:          0,
	})
	if err != nil || len(admDivision.Data) == 0 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("adm_division_id", "adm division id tidak valid")
		return
	}

	Type, _ := strconv.Atoi(req.Data.Type)

	//get ItemSection based on region ID
	itemSections, err := s.opt.Client.CampaignServiceGrpc.GetItemSectionList(ctx, &campaign_service.GetItemSectionListRequest{
		RegionId:    admDivision.Data[0].Region,
		Type:        int32(Type),
		Status:      1,
		CurrentTime: timestamppb.Now(),
		ArchetypeId: address.Data[0].GnL_Archetype_ID,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	glossary, err := s.opt.Client.ConfigurationServiceGrpc.GetGlossaryDetail(ctx, &configuration_service.GetGlossaryDetailRequest{
		Table:     "sales_order",
		Attribute: "order_channel",
		ValueName: req.Platform,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, itemSection := range itemSections.Data {
		var items []*dto.ItemResponse
		for _, item := range itemSection.Items {
			var detailItem *catalog_service.GetItemDetailByInternalIdResponse
			detailItem, err = s.opt.Client.CatalogServiceGrpc.GetItemDetailByInternalId(ctx, &catalog_service.GetItemDetailByInternalIdRequest{
				Id: utils.ToString(item),
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorValidation("item_id", "item id tidak valid")
				return
			}

			if strings.Contains(detailItem.Data.ExcludeArchetype, utils.ToString(glossary.Data.ValueInt)) {
				continue
			}

			items = append(items, &dto.ItemResponse{
				ID:             utils.ToString(detailItem.Data.Id),
				Code:           detailItem.Data.Code,
				ItemName:       detailItem.Data.Description,
				ItemUomName:    detailItem.Data.UomName,
				UnitPrice:      "5000",
				Description:    detailItem.Data.Note,
				OrderMinQty:    "0.01",
				DecimalEnabled: "1",
			})
		}
		res = append(res, &dto.ItemSectionResponse{
			ID:              utils.ToString(itemSection.Id),
			Code:            itemSection.Code,
			Name:            itemSection.Name,
			Region:          utils.ArrayStringToString(itemSection.Regions),
			Archetype:       utils.ArrayStringToString(itemSection.Archetypes),
			BackgroundImage: itemSection.BackgroundImages,
			StartAt:         itemSection.StartAt.AsTime(),
			EndAt:           itemSection.FinishAt.AsTime(),
			Sequence:        utils.ToString(itemSection.Sequence),
			Type:            utils.ToString(itemSection.Type),
			CreatedAt:       itemSection.CreatedAt.AsTime(),
			UpdatedAt:       itemSection.UpdatedAt.AsTime(),
			Item:            items,
		})
	}

	return
}

func (s *ItemSectionService) GetPublic(ctx context.Context, req dto.RequestGetItemSection) (res []*dto.ItemSectionResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemSectionService.Get")
	defer span.End()

	//check Address
	admDivision, err := s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPList(ctx, &bridge_service.GetAdmDivisionGPListRequest{
		AdmDivisionCode: req.Data.AdmDivisionID,
		Limit:           1,
		Offset:          0,
	})
	if err != nil || len(admDivision.Data) == 0 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("adm_division_id", "adm division id tidak valid")
		return
	}

	glossary, err := s.opt.Client.ConfigurationServiceGrpc.GetGlossaryDetail(ctx, &configuration_service.GetGlossaryDetailRequest{
		Table:     "sales_order",
		Attribute: "order_channel",
		ValueName: req.Platform,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	Type, _ := strconv.Atoi(req.Data.Type)
	//get ItemSection based on region ID
	itemSections, err := s.opt.Client.CampaignServiceGrpc.GetItemSectionList(ctx, &campaign_service.GetItemSectionListRequest{
		RegionId:    utils.ToString(admDivision.Data[0].Region),
		Status:      1,
		Type:        int32(Type),
		CurrentTime: timestamppb.Now(),
		ArchetypeId: "ARC0001",
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, itemSection := range itemSections.Data {
		var items []*dto.ItemResponse
		for _, item := range itemSection.Items {
			var detailItem *catalog_service.GetItemDetailByInternalIdResponse
			detailItem, err = s.opt.Client.CatalogServiceGrpc.GetItemDetailByInternalId(ctx, &catalog_service.GetItemDetailByInternalIdRequest{
				Id: utils.ToString(item),
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorValidation("item_id", "item id tidak valid")
				return
			}

			if strings.Contains(detailItem.Data.ExcludeArchetype, utils.ToString(glossary.Data.ValueInt)) {
				continue
			}

			items = append(items, &dto.ItemResponse{
				ID:             utils.ToString(detailItem.Data.Id),
				Code:           detailItem.Data.Code,
				ItemName:       detailItem.Data.Description,
				ItemUomName:    detailItem.Data.UomName,
				UnitPrice:      "5000",
				Description:    detailItem.Data.Note,
				OrderMinQty:    "0.01",
				DecimalEnabled: "1",
			})
		}
		res = append(res, &dto.ItemSectionResponse{
			ID:              utils.ToString(itemSection.Id),
			Code:            itemSection.Code,
			Name:            itemSection.Name,
			Region:          utils.ArrayStringToString(itemSection.Regions),
			Archetype:       utils.ArrayStringToString(itemSection.Archetypes),
			BackgroundImage: itemSection.BackgroundImages,
			StartAt:         itemSection.StartAt.AsTime(),
			EndAt:           itemSection.FinishAt.AsTime(),
			Sequence:        utils.ToString(itemSection.Sequence),
			Type:            utils.ToString(itemSection.Type),
			CreatedAt:       itemSection.CreatedAt.AsTime(),
			UpdatedAt:       itemSection.UpdatedAt.AsTime(),
			Item:            items,
		})
	}

	return
}

func (s *ItemSectionService) GetPublicDetail(ctx context.Context, req dto.RequestGetDetailItemSection) (res *dto.ItemSectionResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemSectionService.Get")
	defer span.End()

	//check Adm
	admDivision, err := s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPList(ctx, &bridge_service.GetAdmDivisionGPListRequest{
		AdmDivisionCode: req.Data.AdmDivisionID,
		Limit:           1,
		Offset:          0,
	})
	if err != nil || len(admDivision.Data) == 0 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("adm_division_id", "adm division id tidak valid")
		return
	}

	glossary, err := s.opt.Client.ConfigurationServiceGrpc.GetGlossaryDetail(ctx, &configuration_service.GetGlossaryDetailRequest{
		Table:     "sales_order",
		Attribute: "order_channel",
		ValueName: req.Platform,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	itemSectionID, _ := strconv.Atoi(req.Data.ItemSectionID)
	//get ItemSection based on region ID
	itemSection, err := s.opt.Client.CampaignServiceGrpc.GetItemSectionList(ctx, &campaign_service.GetItemSectionListRequest{
		ItemSectionId: int64(itemSectionID),
		RegionId:      admDivision.Data[0].Region,
		CurrentTime:   timestamppb.Now(),
		Status:        1,
	})
	if err != nil || len(itemSection.Data) == 0 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	var items []*dto.ItemResponse
	for _, item := range itemSection.Data[0].Items {
		var detailItem *catalog_service.GetItemDetailByInternalIdResponse
		detailItem, err = s.opt.Client.CatalogServiceGrpc.GetItemDetailByInternalId(ctx, &catalog_service.GetItemDetailByInternalIdRequest{
			Id: utils.ToString(item),
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorValidation("item_id", "item id tidak valid")
			return
		}

		if strings.Contains(detailItem.Data.ExcludeArchetype, utils.ToString(glossary.Data.ValueInt)) {
			continue
		}

		items = append(items, &dto.ItemResponse{
			ID:             utils.ToString(detailItem.Data.Id),
			Code:           detailItem.Data.Code,
			ItemName:       detailItem.Data.Description,
			ItemUomName:    detailItem.Data.UomName,
			UnitPrice:      "5000",
			Description:    detailItem.Data.Note,
			OrderMinQty:    "0.01",
			DecimalEnabled: "1",
		})
	}

	res = &dto.ItemSectionResponse{
		ID:              utils.ToString(itemSection.Data[0].Id),
		Code:            itemSection.Data[0].Code,
		Name:            itemSection.Data[0].Name,
		Region:          utils.ArrayStringToString(itemSection.Data[0].Regions),
		Archetype:       utils.ArrayStringToString(itemSection.Data[0].Archetypes),
		BackgroundImage: itemSection.Data[0].BackgroundImages,
		StartAt:         itemSection.Data[0].StartAt.AsTime(),
		EndAt:           itemSection.Data[0].FinishAt.AsTime(),
		Sequence:        utils.ToString(itemSection.Data[0].Sequence),
		Type:            utils.ToString(itemSection.Data[0].Type),
		CreatedAt:       itemSection.Data[0].CreatedAt.AsTime(),
		UpdatedAt:       itemSection.Data[0].UpdatedAt.AsTime(),
		Item:            items,
	}

	return
}

func (s *ItemSectionService) GetPrivateDetail(ctx context.Context, req dto.RequestGetPrivateItemSectionDetail) (res *dto.ItemSectionResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemSectionService.Get")
	defer span.End()

	//check Address
	address, err := s.opt.Client.BridgeServiceGrpc.GetAddressGPDetail(ctx, &bridge_service.GetAddressGPDetailRequest{
		Id: req.Data.AddressID,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("address_id", "address id tidak valid")
		return
	}

	//check Address
	admDivision, err := s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPList(ctx, &bridge_service.GetAdmDivisionGPListRequest{
		AdmDivisionCode: address.Data[0].AdministrativeDiv.GnlAdministrativeCode,
		Limit:           1,
		Offset:          0,
	})
	if err != nil || len(admDivision.Data) == 0 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("adm_division_id", "adm division id tidak valid")
		return
	}

	glossary, err := s.opt.Client.ConfigurationServiceGrpc.GetGlossaryDetail(ctx, &configuration_service.GetGlossaryDetailRequest{
		Table:     "sales_order",
		Attribute: "order_channel",
		ValueName: req.Platform,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	itemSectionID, _ := strconv.Atoi(req.Data.ItemSectionID)
	//get ItemSection based on region ID
	itemSection, err := s.opt.Client.CampaignServiceGrpc.GetItemSectionList(ctx, &campaign_service.GetItemSectionListRequest{
		ItemSectionId: int64(itemSectionID),
		RegionId:      admDivision.Data[0].Region,
		ArchetypeId:   address.Data[0].GnL_Archetype_ID,
		Status:        1,
	})
	if err != nil || len(itemSection.Data) == 0 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	var items []*dto.ItemResponse
	for _, item := range itemSection.Data[0].Items {
		var detailItem *catalog_service.GetItemDetailByInternalIdResponse
		detailItem, err = s.opt.Client.CatalogServiceGrpc.GetItemDetailByInternalId(ctx, &catalog_service.GetItemDetailByInternalIdRequest{
			Id: utils.ToString(item),
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorValidation("item_id", "item id tidak valid")
			return
		}

		if strings.Contains(detailItem.Data.ExcludeArchetype, utils.ToString(glossary.Data.ValueInt)) {
			continue
		}

		items = append(items, &dto.ItemResponse{
			ID:             utils.ToString(detailItem.Data.Id),
			Code:           detailItem.Data.Code,
			ItemName:       detailItem.Data.Description,
			ItemUomName:    detailItem.Data.UomName,
			UnitPrice:      "5000",
			Description:    detailItem.Data.Note,
			OrderMinQty:    "0.01",
			DecimalEnabled: "1",
		})
	}

	res = &dto.ItemSectionResponse{
		ID:              utils.ToString(itemSection.Data[0].Id),
		Code:            itemSection.Data[0].Code,
		Name:            itemSection.Data[0].Name,
		Region:          utils.ArrayStringToString(itemSection.Data[0].Regions),
		Archetype:       utils.ArrayStringToString(itemSection.Data[0].Archetypes),
		BackgroundImage: itemSection.Data[0].BackgroundImages,
		StartAt:         itemSection.Data[0].StartAt.AsTime(),
		EndAt:           itemSection.Data[0].FinishAt.AsTime(),
		Sequence:        utils.ToString(itemSection.Data[0].Sequence),
		Type:            utils.ToString(itemSection.Data[0].Type),
		CreatedAt:       itemSection.Data[0].CreatedAt.AsTime(),
		UpdatedAt:       itemSection.Data[0].UpdatedAt.AsTime(),
		Item:            items,
	}

	return
}
