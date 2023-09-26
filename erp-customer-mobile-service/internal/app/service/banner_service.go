package service

import (
	"context"
	"strconv"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/utils"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/campaign_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/catalog_service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IBannerService interface {
	GetPrivate(ctx context.Context, req dto.RequestGetPrivateBanner) (res []dto.ResponseBanner, err error)
	GetPublic(ctx context.Context, req dto.RequestGetBanner) (res []dto.ResponseBanner, err error)
}

type BannerService struct {
	opt opt.Options
	//RepositoryOTPOutgoing repository.IOtpOutgoingRepository
}

func NewBannerService() IBannerService {
	return &BannerService{
		opt: global.Setup.Common,
		//RepositoryOTPOutgoing: repository.NewOtpOutgoingRepository(),
	}
}

func (s *BannerService) GetPrivate(ctx context.Context, req dto.RequestGetPrivateBanner) (res []dto.ResponseBanner, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "BannerService.Get")
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

	//check Admin Division
	admDivision, err := s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPList(ctx, &bridge_service.GetAdmDivisionGPListRequest{
		AdmDivisionCode: address.Data[0].AdministrativeDiv.GnlAdministrativeCode,
		Limit:           10,
		Offset:          1,
	})
	if err != nil || len(admDivision.Data) == 0 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("adm_division_id", "adm division id tidak valid")
		return
	}

	//get Banner based on region ID
	banners, err := s.opt.Client.CampaignServiceGrpc.GetBannerList(ctx, &campaign_service.GetBannerListRequest{
		RegionId:    admDivision.Data[0].Region,
		ArchetypeId: address.Data[0].GnL_Archetype_ID,
		CurrentTime: timestamppb.Now(),
		Status:      1,
		Limit:       5,
	})

	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, banner := range banners.Data {
		var (
			itemCategory         *catalog_service.GetItemCategoryDetailResponse
			itemCategoryResponse *dto.ItemCategoryResponse
			item                 *catalog_service.GetItemDetailByInternalIdResponse
			itemResponse         *dto.ItemResponse
			itemSection          *campaign_service.GetItemSectionDetailResponse
			itemSectionResponse  *dto.ItemSectionResponse
		)
		if banner.RedirectTo == 2 {
			itemCategory, err = s.opt.Client.CatalogServiceGrpc.GetItemCategoryDetail(ctx, &catalog_service.GetItemCategoryDetailRequest{
				Id: utils.ToInt64(banner.RedirectValue),
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}

			itemCategoryResponse = &dto.ItemCategoryResponse{
				ID:       utils.ToString(itemCategory.Data.Id),
				Region:   itemCategory.Data.RegionId,
				Name:     itemCategory.Data.Name,
				ImageUrl: itemCategory.Data.ImageUrl,
				Status:   utils.ToString(itemCategory.Data.Status),
			}

		} else if banner.RedirectTo == 3 {
			item, err = s.opt.Client.CatalogServiceGrpc.GetItemDetailByInternalId(ctx, &catalog_service.GetItemDetailByInternalIdRequest{
				Id: banner.RedirectValue,
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("catalog", "item")
				return
			}

			itemResponse = &dto.ItemResponse{
				ID:                  strconv.Itoa(int(item.Data.Id)),
				Code:                item.Data.Code,
				ItemName:            item.Data.Description,
				ItemUomName:         item.Data.UomName,
				Description:         item.Data.Note,
				UnitPrice:           "5000",
				OrderMinQty:         "1",
				DecimalEnabled:      "1",
				ItemCategoryNameArr: utils.StringToStringArray(item.Data.ItemCategoryName),
			}
		} else if banner.RedirectTo == 4 {
			itemSectionResponse = &dto.ItemSectionResponse{
				ID: "0",
			}
		} else if banner.RedirectTo == 6 {
			itemSection, err = s.opt.Client.CampaignServiceGrpc.GetItemSectionDetail(ctx, &campaign_service.GetItemSectionDetailRequest{
				Id: utils.ToInt64(banner.RedirectValue),
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}
			itemSectionResponse = &dto.ItemSectionResponse{
				ID:              utils.ToString(itemSection.Data.Id),
				Code:            itemSection.Data.Code,
				Name:            itemSection.Data.Name,
				Region:          utils.ArrayStringToString(itemSection.Data.Regions),
				Archetype:       utils.ArrayStringToString(itemSection.Data.Archetypes),
				BackgroundImage: itemSection.Data.BackgroundImages,
				StartAt:         itemSection.Data.StartAt.AsTime(),
				EndAt:           itemSection.Data.FinishAt.AsTime(),
				Sequence:        utils.ToString(itemSection.Data.Sequence),
				Type:            utils.ToString(itemSection.Data.Type),
				CreatedAt:       itemSection.Data.CreatedAt.AsTime(),
				UpdatedAt:       itemSection.Data.UpdatedAt.AsTime(),
			}
		}

		res = append(res, dto.ResponseBanner{
			ID:             strconv.Itoa(int(banner.Id)),
			Code:           banner.Code,
			Name:           banner.Name,
			ImageUrl:       banner.ImageUrl,
			StartDate:      banner.StartAt.AsTime(),
			EndDate:        banner.FinishAt.AsTime(),
			NavigationType: strconv.Itoa(int(banner.RedirectTo)),
			NavigationUrl:  banner.RedirectValue,
			Region:         utils.ArrayStringToString(banner.Regions),
			Archetype:      utils.ArrayStringToString(banner.Archetypes),
			Queue:          strconv.Itoa(int(banner.Queue)),
			Note:           banner.Note,
			Status:         strconv.Itoa(int(banner.Status)),
			CreatedAt:      banner.CreatedAt.AsTime(),
			ItemCategory:   itemCategoryResponse,
			ItemSection:    itemSectionResponse,
			Item:           itemResponse,
		})
	}

	return
}
func (s *BannerService) GetPublic(ctx context.Context, req dto.RequestGetBanner) (res []dto.ResponseBanner, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "BannerService.Get")
	defer span.End()

	//check Admin Division
	admDivision, err := s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPList(ctx, &bridge_service.GetAdmDivisionGPListRequest{
		AdmDivisionCode: req.Data.AdmDivisionId,
		Limit:           10,
		Offset:          1,
	})
	if err != nil || len(admDivision.Data) == 0 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("adm_division_id", "adm division id tidak valid")
		return
	}

	//get Banner based on region ID
	banners, err := s.opt.Client.CampaignServiceGrpc.GetBannerList(ctx, &campaign_service.GetBannerListRequest{
		RegionId:    admDivision.Data[0].Region,
		CurrentTime: timestamppb.Now(),
		Status:      1,
		Limit:       5,
	})
	for _, banner := range banners.Data {
		var (
			itemCategory         *catalog_service.GetItemCategoryDetailResponse
			itemCategoryResponse *dto.ItemCategoryResponse
			item                 *catalog_service.GetItemDetailByInternalIdResponse
			itemResponse         *dto.ItemResponse
			itemSection          *campaign_service.GetItemSectionDetailResponse
			itemSectionResponse  *dto.ItemSectionResponse
		)
		if banner.RedirectTo == 2 {
			itemCategory, err = s.opt.Client.CatalogServiceGrpc.GetItemCategoryDetail(ctx, &catalog_service.GetItemCategoryDetailRequest{
				Id: utils.ToInt64(banner.RedirectValue),
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}

			itemCategoryResponse = &dto.ItemCategoryResponse{
				ID:       utils.ToString(itemCategory.Data.Id),
				Region:   itemCategory.Data.RegionId,
				Name:     itemCategory.Data.Name,
				ImageUrl: itemCategory.Data.ImageUrl,
				Status:   utils.ToString(itemCategory.Data.Status),
			}

		} else if banner.RedirectTo == 3 {
			item, err = s.opt.Client.CatalogServiceGrpc.GetItemDetailByInternalId(ctx, &catalog_service.GetItemDetailByInternalIdRequest{
				Id: banner.RedirectValue,
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("catalog", "item")
				return
			}

			itemResponse = &dto.ItemResponse{
				ID:                  strconv.Itoa(int(item.Data.Id)),
				Code:                item.Data.Code,
				ItemName:            item.Data.Description,
				ItemUomName:         item.Data.UomName,
				Description:         item.Data.Note,
				UnitPrice:           "5000",
				OrderMinQty:         "1",
				DecimalEnabled:      "1",
				ItemCategoryNameArr: utils.StringToStringArray(item.Data.ItemCategoryName),
			}
		} else if banner.RedirectTo == 4 {
			itemSectionResponse = &dto.ItemSectionResponse{
				ID: "0",
			}
		} else if banner.RedirectTo == 6 {
			itemSection, err = s.opt.Client.CampaignServiceGrpc.GetItemSectionDetail(ctx, &campaign_service.GetItemSectionDetailRequest{
				Id: utils.ToInt64(banner.RedirectValue),
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}
			itemSectionResponse = &dto.ItemSectionResponse{
				ID:              utils.ToString(itemSection.Data.Id),
				Code:            itemSection.Data.Code,
				Name:            itemSection.Data.Name,
				Region:          utils.ArrayStringToString(itemSection.Data.Regions),
				Archetype:       utils.ArrayStringToString(itemSection.Data.Archetypes),
				BackgroundImage: itemSection.Data.BackgroundImages,
				StartAt:         itemSection.Data.StartAt.AsTime(),
				EndAt:           itemSection.Data.FinishAt.AsTime(),
				Sequence:        utils.ToString(itemSection.Data.Sequence),
				Type:            utils.ToString(itemSection.Data.Type),
				CreatedAt:       itemSection.Data.CreatedAt.AsTime(),
				UpdatedAt:       itemSection.Data.UpdatedAt.AsTime(),
			}
		}
		res = append(res, dto.ResponseBanner{
			ID:             strconv.Itoa(int(banner.Id)),
			Code:           banner.Code,
			Name:           banner.Name,
			ImageUrl:       banner.ImageUrl,
			StartDate:      banner.StartAt.AsTime(),
			EndDate:        banner.FinishAt.AsTime(),
			NavigationType: strconv.Itoa(int(banner.RedirectTo)),
			NavigationUrl:  banner.RedirectValue,
			Region:         utils.ArrayStringToString(banner.Regions),
			Archetype:      utils.ArrayStringToString(banner.Archetypes),
			Queue:          strconv.Itoa(int(banner.Queue)),
			Note:           banner.Note,
			Status:         strconv.Itoa(int(banner.Status)),
			CreatedAt:      banner.CreatedAt.AsTime(),
			ItemCategory:   itemCategoryResponse,
			ItemSection:    itemSectionResponse,
			Item:           itemResponse,
		})
	}

	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}
