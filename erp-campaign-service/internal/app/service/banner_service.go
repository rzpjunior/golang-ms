package service

import (
	"context"
	"net/url"
	"strconv"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/constants"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/edenlabs/edenlabs/timex"
	"git.edenfarm.id/edenlabs/edenlabs/utils"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/repository"
	auditService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/audit_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	catalogService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/catalog_service"
	configurationService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IBannerService interface {
	Get(ctx context.Context, req *dto.BannerRequestGet) (res []*dto.BannerResponse, total int64, err error)
	GetByID(ctx context.Context, id int64) (res dto.BannerResponse, err error)
	Create(ctx context.Context, req dto.BannerRequestCreate) (res dto.BannerResponse, err error)
	Archive(ctx context.Context, id int64, req dto.BannerRequestArchive) (res dto.BannerResponse, err error)
	GetListMobile(ctx context.Context, req *dto.BannerRequestGet) (res []*dto.BannerResponse, total int64, err error)
}

type BannerService struct {
	opt                   opt.Options
	RepositoryBanner      repository.IBannerRepository
	RepositoryItemSection repository.IItemSectionRepository
}

func NewBannerService() IBannerService {
	return &BannerService{
		opt:                   global.Setup.Common,
		RepositoryBanner:      repository.NewBannerRepository(),
		RepositoryItemSection: repository.NewItemSectionRepository(),
	}
}

func (s *BannerService) Get(ctx context.Context, req *dto.BannerRequestGet) (res []*dto.BannerResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "BannerService.Get")
	defer span.End()

	var banners []*model.Banner
	banners, total, err = s.RepositoryBanner.Get(ctx, req)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	currentTime := time.Now()
	for _, banner := range banners {
		if banner.Status == statusx.ConvertStatusName("Draft") {
			if currentTime.After(banner.StartAt) {
				// update status from draft to active
				err = s.RepositoryBanner.Update(ctx, &model.Banner{ID: banner.ID, Status: statusx.ConvertStatusName("Active"), UpdatedAt: time.Now()}, "Status", "UpdatedAt")
				if err != nil {
					span.RecordError(err)
					s.opt.Logger.AddMessage(log.ErrorLevel, err)
					return
				}
				banner.Status = statusx.ConvertStatusName("Active")
			}
		}

		if banner.Status == statusx.ConvertStatusName("Active") {
			if currentTime.After(banner.FinishAt) {
				// update status from active to finish
				err = s.RepositoryBanner.Update(ctx, &model.Banner{ID: banner.ID, Status: statusx.ConvertStatusName("Finished"), UpdatedAt: time.Now()}, "Status", "UpdatedAt")
				if err != nil {
					span.RecordError(err)
					s.opt.Logger.AddMessage(log.ErrorLevel, err)
					return
				}
				banner.Status = statusx.ConvertStatusName("Finished")
			}
		}

		var redirect *dto.RedirectResponse
		switch banner.RedirectTo {
		case 1:
			redirect = &dto.RedirectResponse{
				To:        banner.RedirectTo,
				Value:     banner.RedirectValue,
				Name:      "url",
				ValueName: "url",
			}
		case 2:
			var itemCategory *catalogService.GetItemCategoryDetailResponse
			itemCategory, err = s.opt.Client.CatalogServiceGrpc.GetItemCategoryDetail(ctx, &catalogService.GetItemCategoryDetailRequest{
				Id: utils.ToInt64(banner.RedirectValue),
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("catalog", "item_category")
				return
			}
			redirect = &dto.RedirectResponse{
				To: banner.RedirectTo,
				Value: dto.RedirectToItemCategoryResponse{
					ID:   itemCategory.Data.Id,
					Name: itemCategory.Data.Name,
				},
				Name:      "item category",
				ValueName: "item category",
			}
		case 3:
			var item *catalogService.GetItemDetailByInternalIdResponse
			item, err = s.opt.Client.CatalogServiceGrpc.GetItemDetailByInternalId(ctx, &catalogService.GetItemDetailByInternalIdRequest{
				Id: utils.ToString(banner.RedirectValue),
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("catalog", "item")
				return
			}

			redirect = &dto.RedirectResponse{
				To: banner.RedirectTo,
				Value: dto.RedirectToItemResponse{
					ID:   item.Data.Id,
					Code: item.Data.Code,
					Name: item.Data.Description,
				},
				Name:      "item",
				ValueName: "item",
			}
		case 4:
			redirect = &dto.RedirectResponse{
				To:        banner.RedirectTo,
				Value:     banner.RedirectValue,
				Name:      "promo",
				ValueName: "promo",
			}
		case 5:
			redirect = &dto.RedirectResponse{
				To:        banner.RedirectTo,
				Value:     banner.RedirectValue,
				Name:      "no redirect",
				ValueName: "no redirect",
			}
		case 6:
			var itemSection *model.ItemSection
			itemSection, err = s.RepositoryItemSection.GetByID(ctx, utils.ToInt64(banner.RedirectValue))
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}

			redirect = &dto.RedirectResponse{
				To: banner.RedirectTo,
				Value: dto.RedirectToItemSectionResponse{
					ID:   itemSection.ID,
					Name: itemSection.Name,
				},
				Name:      "item section",
				ValueName: "item section",
			}
		case 7:
			redirect = &dto.RedirectResponse{
				To:        banner.RedirectTo,
				Name:      "Eden Rewards",
				ValueName: "Eden Rewards",
			}
		case 8:
			redirect = &dto.RedirectResponse{
				To:        banner.RedirectTo,
				Name:      "Referral",
				ValueName: "Referral",
			}
		case 9:
			redirect = &dto.RedirectResponse{
				To:        banner.RedirectTo,
				Name:      "Form Pendaftaran Akun Bisnis",
				ValueName: "Form Pendaftaran Akun Bisnis",
			}
		default:
			redirect = &dto.RedirectResponse{}
		}

		// get region name
		var regionNames []string
		for _, regionID := range utils.StringToStringArray(banner.Regions) {
			var region *bridgeService.GetAdmDivisionGPResponse
			region, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPDetail(ctx, &bridgeService.GetAdmDivisionGPDetailRequest{
				Region: regionID,
				Limit:  1,
				Offset: 0,
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "region")
				return
			}
			regionNames = append(regionNames, region.Data[0].Region)
		}

		// get archetype name
		var archetypeNames []string
		for _, archetypeID := range utils.StringToStringArray(banner.Archetypes) {
			var archetype *bridge_service.GetArchetypeGPResponse
			archetype, err = s.opt.Client.BridgeServiceGrpc.GetArchetypeGPDetail(ctx, &bridge_service.GetArchetypeGPDetailRequest{
				Id: archetypeID,
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorInvalid("archetype_id")
				return
			}
			archetypeNames = append(archetypeNames, archetype.Data[0].GnlArchetypedescription)
		}

		res = append(res, &dto.BannerResponse{
			ID:             banner.ID,
			Code:           banner.Code,
			Name:           banner.Name,
			Regions:        utils.StringToStringArray(banner.Regions),
			RegionNames:    regionNames,
			Archetypes:     utils.StringToStringArray(banner.Archetypes),
			ArchetypeNames: archetypeNames,
			Queue:          banner.Queue,
			ImageUrl:       banner.ImageUrl,
			StartAt:        timex.ToLocTime(ctx, banner.StartAt),
			FinishAt:       timex.ToLocTime(ctx, banner.FinishAt),
			Redirect:       redirect,
			Note:           banner.Note,
			CreatedAt:      timex.ToLocTime(ctx, banner.CreatedAt),
			UpdatedAt:      timex.ToLocTime(ctx, banner.UpdatedAt),
			Status:         banner.Status,
		})
	}

	return
}

func (s *BannerService) GetByID(ctx context.Context, id int64) (res dto.BannerResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "BannerService.GetByID")
	defer span.End()

	var banner *model.Banner
	banner, err = s.RepositoryBanner.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	currentTime := time.Now()
	if banner.Status == statusx.ConvertStatusName("Draft") {
		if currentTime.After(banner.StartAt) {
			// update status from draft to active
			err = s.RepositoryBanner.Update(ctx, &model.Banner{ID: banner.ID, Status: statusx.ConvertStatusName("Active"), UpdatedAt: time.Now()}, "Status", "UpdatedAt")
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}
			banner.Status = statusx.ConvertStatusName("Active")
		}
	}

	if banner.Status == statusx.ConvertStatusName("Active") {
		if currentTime.After(banner.FinishAt) {
			// update status from active to finish
			err = s.RepositoryBanner.Update(ctx, &model.Banner{ID: banner.ID, Status: statusx.ConvertStatusName("Finished"), UpdatedAt: time.Now()}, "Status", "UpdatedAt")
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}
			banner.Status = statusx.ConvertStatusName("Finished")
		}
	}

	var redirect *dto.RedirectResponse
	switch banner.RedirectTo {
	case 1:
		redirect = &dto.RedirectResponse{
			To:        banner.RedirectTo,
			Value:     banner.RedirectValue,
			Name:      "url",
			ValueName: "url",
		}
	case 2:
		var itemCategory *catalogService.GetItemCategoryDetailResponse
		itemCategory, err = s.opt.Client.CatalogServiceGrpc.GetItemCategoryDetail(ctx, &catalogService.GetItemCategoryDetailRequest{
			Id: utils.ToInt64(banner.RedirectValue),
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("catalog", "item_category")
			return
		}
		redirect = &dto.RedirectResponse{
			To: banner.RedirectTo,
			Value: dto.RedirectToItemCategoryResponse{
				ID:   itemCategory.Data.Id,
				Name: itemCategory.Data.Name,
			},
			Name:      "item category",
			ValueName: "item category",
		}
	case 3:
		var item *catalogService.GetItemDetailByInternalIdResponse
		item, err = s.opt.Client.CatalogServiceGrpc.GetItemDetailByInternalId(ctx, &catalogService.GetItemDetailByInternalIdRequest{
			Id: utils.ToString(banner.RedirectValue),
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("catalog", "item")
			return
		}

		redirect = &dto.RedirectResponse{
			To: banner.RedirectTo,
			Value: dto.RedirectToItemResponse{
				ID:   item.Data.Id,
				Code: item.Data.Code,
				Name: item.Data.Description,
			},
			Name:      "item",
			ValueName: "item",
		}
	case 4:
		redirect = &dto.RedirectResponse{
			To:        banner.RedirectTo,
			Value:     banner.RedirectValue,
			Name:      "promo",
			ValueName: "promo",
		}
	case 5:
		redirect = &dto.RedirectResponse{
			To:        banner.RedirectTo,
			Value:     banner.RedirectValue,
			Name:      "no redirect",
			ValueName: "no redirect",
		}
	case 6:
		var itemSection *model.ItemSection
		itemSection, err = s.RepositoryItemSection.GetByID(ctx, utils.ToInt64(banner.RedirectValue))
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		redirect = &dto.RedirectResponse{
			To: banner.RedirectTo,
			Value: dto.RedirectToItemSectionResponse{
				ID:   itemSection.ID,
				Name: itemSection.Name,
			},
			Name:      "item section",
			ValueName: "item section",
		}
	case 7:
		redirect = &dto.RedirectResponse{
			To:        banner.RedirectTo,
			Name:      "Eden Rewards",
			ValueName: "Eden Rewards",
		}
	case 8:
		redirect = &dto.RedirectResponse{
			To:        banner.RedirectTo,
			Name:      "Referral",
			ValueName: "Referral",
		}
	case 9:
		redirect = &dto.RedirectResponse{
			To:        banner.RedirectTo,
			Name:      "Form Pendaftaran Akun Bisnis",
			ValueName: "Form Pendaftaran Akun Bisnis",
		}
	default:
		redirect = &dto.RedirectResponse{}
	}

	// get region name
	var regionNames []string
	for _, regionID := range utils.StringToStringArray(banner.Regions) {
		var region *bridgeService.GetAdmDivisionGPResponse
		region, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPDetail(ctx, &bridgeService.GetAdmDivisionGPDetailRequest{
			Region: regionID,
			Limit:  1,
			Offset: 0,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "region")
			return
		}
		regionNames = append(regionNames, region.Data[0].Region)
	}

	// get archetype name
	var archetypeNames []string
	for _, archetypeID := range utils.StringToStringArray(banner.Archetypes) {
		var archetype *bridge_service.GetArchetypeGPResponse
		archetype, err = s.opt.Client.BridgeServiceGrpc.GetArchetypeGPDetail(ctx, &bridge_service.GetArchetypeGPDetailRequest{
			Id: archetypeID,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("archetype_id")
			return
		}
		archetypeNames = append(archetypeNames, archetype.Data[0].GnlArchetypedescription)
	}

	res = dto.BannerResponse{
		ID:             banner.ID,
		Name:           banner.Name,
		Regions:        utils.StringToStringArray(banner.Regions),
		RegionNames:    regionNames,
		Archetypes:     utils.StringToStringArray(banner.Archetypes),
		ArchetypeNames: archetypeNames,
		Code:           banner.Code,
		Queue:          banner.Queue,
		ImageUrl:       banner.ImageUrl,
		StartAt:        timex.ToLocTime(ctx, banner.StartAt),
		FinishAt:       timex.ToLocTime(ctx, banner.FinishAt),
		Redirect:       redirect,
		Note:           banner.Note,
		CreatedAt:      timex.ToLocTime(ctx, banner.CreatedAt),
		UpdatedAt:      timex.ToLocTime(ctx, banner.UpdatedAt),
		Status:         banner.Status,
	}

	return
}

func (s *BannerService) Create(ctx context.Context, req dto.BannerRequestCreate) (res dto.BannerResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "BannerService.Create")
	defer span.End()

	// validate time
	currentTime := time.Now()
	if !req.StartAt.IsZero() && req.StartAt.Before(currentTime) {
		err = edenlabs.ErrorValidation("start_at", "Start date must be later than current date")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if !req.StartAt.IsZero() && !req.FinishAt.IsZero() && (req.StartAt.After(req.FinishAt) || req.StartAt.Equal(req.FinishAt)) {
		err = edenlabs.ErrorValidation("finish_at", "End date must be later than start date")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// Validate region id
	for _, v := range req.Regions {
		_, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPDetail(ctx, &bridgeService.GetAdmDivisionGPDetailRequest{
			Region: v,
			Limit:  1,
			Offset: 0,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("region_id")
			return
		}
	}

	// Validate archetype id
	for _, v := range req.Archetypes {
		_, err = s.opt.Client.BridgeServiceGrpc.GetArchetypeGPDetail(ctx, &bridge_service.GetArchetypeGPDetailRequest{
			Id: v,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("archetype_id")
			return
		}
	}

	switch req.RedirectTo {
	case 1:
		// url
		if req.RedirectValue == "" {
			err = edenlabs.ErrorValidation("redirect_value", "The redirect value is required")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		_, err = url.ParseRequestURI(req.RedirectValue)
		if err != nil {
			err = edenlabs.ErrorValidation("redirect_value", "The redirect value is invalid value")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	case 2:
		// item category
		if req.RedirectValue == "" {
			err = edenlabs.ErrorValidation("redirect_value", "The redirect value is required")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		if _, err = strconv.Atoi(req.RedirectValue); err != nil {
			err = edenlabs.ErrorValidation("redirect_value", "The redirect value is invalid value")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		// validate the redirect value is existed
		_, err = s.opt.Client.CatalogServiceGrpc.GetItemCategoryDetail(ctx, &catalogService.GetItemCategoryDetailRequest{
			Id: utils.ToInt64(req.RedirectValue),
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorValidation("redirect_value", "The redirect value is invalid value")
			return
		}

	case 3:
		// item
		if req.RedirectValue == "" {
			err = edenlabs.ErrorValidation("redirect_value", "The redirect value is required")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		if _, err = strconv.Atoi(req.RedirectValue); err != nil {
			err = edenlabs.ErrorValidation("redirect_value", "The redirect value is invalid value")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		// validate the redirect value is existed
		_, err = s.opt.Client.BridgeServiceGrpc.GetItemDetail(ctx, &bridgeService.GetItemDetailRequest{
			Id: utils.ToInt64(req.RedirectValue),
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorValidation("redirect_value", "The redirect value is invalid value")
			return
		}
	case 4:
		// promo
	case 5:
		// no redirect
	case 6:
		// item section
		if req.RedirectValue == "" {
			err = edenlabs.ErrorValidation("redirect_value", "The redirect value is required")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		if _, err = strconv.Atoi(req.RedirectValue); err != nil {
			err = edenlabs.ErrorValidation("redirect_value", "The redirect value is invalid value")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	default:
	}

	var codeGenerator *configurationService.GetGenerateCodeResponse
	codeGenerator, err = s.opt.Client.ConfigurationServiceGrpc.GetGenerateCode(ctx, &configurationService.GetGenerateCodeRequest{
		Format: "BNR",
		Domain: "banner",
		Length: 6,
	})

	banner := &model.Banner{
		Name:          req.Name,
		Code:          codeGenerator.Data.Code,
		Regions:       utils.ArrayStringToString(req.Regions),
		Archetypes:    utils.ArrayStringToString(req.Archetypes),
		Queue:         req.Queue,
		RedirectTo:    int8(req.RedirectTo),
		RedirectValue: req.RedirectValue,
		StartAt:       req.StartAt,
		FinishAt:      req.FinishAt,
		ImageUrl:      req.ImageUrl,
		Note:          req.Note,
		CreatedAt:     time.Now(),
		Status:        statusx.ConvertStatusName("Draft"),
	}

	span.AddEvent("creating new banner")
	err = s.RepositoryBanner.Create(ctx, banner)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	span.AddEvent("banner is created", trace.WithAttributes(attribute.Int64("banner_id", banner.ID)))

	userID := ctx.Value(constants.KeyUserID).(int64)

	_, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &auditService.CreateLogRequest{
		Log: &auditService.Log{
			UserId:      userID,
			ReferenceId: strconv.Itoa(int(banner.ID)),
			Type:        "banner",
			Function:    "create",
			CreatedAt:   timestamppb.New(time.Now()),
			Note:        req.Note,
		},
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.BannerResponse{
		ID:         banner.ID,
		Code:       banner.Code,
		Name:       banner.Name,
		Regions:    utils.StringToStringArray(banner.Regions),
		Archetypes: utils.StringToStringArray(banner.Archetypes),
		Queue:      banner.Queue,
		ImageUrl:   banner.ImageUrl,
		StartAt:    timex.ToLocTime(ctx, banner.StartAt),
		FinishAt:   timex.ToLocTime(ctx, banner.FinishAt),
		Redirect: &dto.RedirectResponse{
			To:    int8(req.RedirectTo),
			Value: req.RedirectValue,
		},
		Note:      banner.Note,
		CreatedAt: timex.ToLocTime(ctx, banner.CreatedAt),
		UpdatedAt: timex.ToLocTime(ctx, banner.UpdatedAt),
		Status:    banner.Status,
	}

	return
}

func (s *BannerService) Archive(ctx context.Context, id int64, req dto.BannerRequestArchive) (res dto.BannerResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "BannerService.Delete")
	defer span.End()

	// validate banner is exist
	var bannerOld *model.Banner
	bannerOld, err = s.RepositoryBanner.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// validate banner time
	currentTime := time.Now()
	if currentTime.After(bannerOld.FinishAt) {
		err = edenlabs.ErrorValidation("status", "The banner must be active")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if bannerOld.Status == statusx.ConvertStatusName("Archived") {
		err = edenlabs.ErrorValidation("status", "The status has been archived")
		return
	}

	err = s.RepositoryBanner.Archive(ctx, id, req)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	userID := ctx.Value(constants.KeyUserID).(int64)

	_, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &auditService.CreateLogRequest{
		Log: &auditService.Log{
			UserId:      userID,
			ReferenceId: strconv.Itoa(int(id)),
			Type:        "banner",
			Function:    "archive",
			CreatedAt:   timestamppb.New(time.Now()),
			Note:        req.Note,
		},
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.BannerResponse{
		ID:         bannerOld.ID,
		Name:       bannerOld.Name,
		Regions:    utils.StringToStringArray(bannerOld.Regions),
		Archetypes: utils.StringToStringArray(bannerOld.Archetypes),
		Code:       bannerOld.Code,
		Queue:      bannerOld.Queue,
		ImageUrl:   bannerOld.ImageUrl,
		StartAt:    timex.ToLocTime(ctx, bannerOld.StartAt),
		FinishAt:   timex.ToLocTime(ctx, bannerOld.FinishAt),
		Note:       req.Note,
		CreatedAt:  timex.ToLocTime(ctx, bannerOld.CreatedAt),
		UpdatedAt:  timex.ToLocTime(ctx, time.Now()),
		Status:     statusx.ConvertStatusName("Archived"),
	}

	return
}

func (s *BannerService) GetListMobile(ctx context.Context, req *dto.BannerRequestGet) (res []*dto.BannerResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "BannerService.Get")
	defer span.End()

	var banners []*model.Banner
	banners, total, err = s.RepositoryBanner.Get(ctx, req)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	currentTime := time.Now()
	for _, banner := range banners {
		if banner.Status == statusx.ConvertStatusName("Draft") {
			if currentTime.After(banner.StartAt) {
				// update status from draft to active
				err = s.RepositoryBanner.Update(ctx, &model.Banner{ID: banner.ID, Status: statusx.ConvertStatusName("Active"), UpdatedAt: time.Now()}, "Status", "UpdatedAt")
				if err != nil {
					span.RecordError(err)
					s.opt.Logger.AddMessage(log.ErrorLevel, err)
					return
				}
				banner.Status = statusx.ConvertStatusName("Active")
			}
		}

		if banner.Status == statusx.ConvertStatusName("Active") {
			if currentTime.After(banner.FinishAt) {
				// update status from active to finish
				err = s.RepositoryBanner.Update(ctx, &model.Banner{ID: banner.ID, Status: statusx.ConvertStatusName("Finished"), UpdatedAt: time.Now()}, "Status", "UpdatedAt")
				if err != nil {
					span.RecordError(err)
					s.opt.Logger.AddMessage(log.ErrorLevel, err)
					return
				}
				banner.Status = statusx.ConvertStatusName("Finished")
			}
		}

		res = append(res, &dto.BannerResponse{
			ID:         banner.ID,
			Code:       banner.Code,
			Name:       banner.Name,
			Regions:    utils.StringToStringArray(banner.Regions),
			Archetypes: utils.StringToStringArray(banner.Archetypes),
			Queue:      banner.Queue,
			ImageUrl:   banner.ImageUrl,
			StartAt:    timex.ToLocTime(ctx, banner.StartAt),
			FinishAt:   timex.ToLocTime(ctx, banner.FinishAt),
			Redirect:   &dto.RedirectResponse{To: banner.RedirectTo, Value: banner.RedirectValue},
			Note:       banner.Note,
			CreatedAt:  timex.ToLocTime(ctx, banner.CreatedAt),
			UpdatedAt:  timex.ToLocTime(ctx, banner.UpdatedAt),
			Status:     banner.Status,
		})
	}

	return
}
