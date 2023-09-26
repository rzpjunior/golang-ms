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
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/catalog_service"
	catalogService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/catalog_service"
	configurationService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/crm_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/customer_mobile_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/notification_service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IPushNotificationService interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, regionID string, scheduledAtFrom time.Time, scheduledAtTo time.Time) (res []*dto.PushNotificationResponse, total int64, err error)
	GetByID(ctx context.Context, id int64) (res dto.PushNotificationResponse, err error)
	Create(ctx context.Context, req dto.PushNotificationRequestCreate) (res dto.PushNotificationResponse, err error)
	Update(ctx context.Context, req dto.PushNotificationRequestUpdate, id int64) (res dto.PushNotificationResponse, err error)
	Cancel(ctx context.Context, req dto.PushNotificationRequestCancel, id int64) (res dto.PushNotificationResponse, err error)
	UpdateOpened(ctx context.Context, req *dto.PushNotificationRequestUpdateOpened) (err error)
	GetDetailMobile(ctx context.Context, id int64) (res *dto.PushNotificationResponse, err error)
}

type PushNotificationService struct {
	opt                        opt.Options
	RepositoryPushNotification repository.IPushNotificationRepository
	RepositoryItemSection      repository.IItemSectionRepository
}

func NewPushNotificationService() IPushNotificationService {
	return &PushNotificationService{
		opt:                        global.Setup.Common,
		RepositoryPushNotification: repository.NewPushNotificationRepository(),
		RepositoryItemSection:      repository.NewItemSectionRepository(),
	}
}

func (s *PushNotificationService) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, regionID string, scheduledAtFrom time.Time, scheduledAtTo time.Time) (res []*dto.PushNotificationResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PushNotificationService.Get")
	defer span.End()

	var pushNotifications []*model.PushNotification
	pushNotifications, total, err = s.RepositoryPushNotification.Get(ctx, offset, limit, status, search, orderBy, regionID, scheduledAtFrom, scheduledAtTo)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	currentTime := time.Now()
	for _, pushNotification := range pushNotifications {
		if pushNotification.Status == statusx.ConvertStatusName(statusx.Draft) {
			if currentTime.After(pushNotification.ScheduledAt) {
				// update status from draft to active
				pushNotification.Status = statusx.ConvertStatusName(statusx.Active)
				err = s.RepositoryPushNotification.Update(ctx, &model.PushNotification{ID: pushNotification.ID, Status: pushNotification.Status, UpdatedAt: time.Now()}, "Status", "UpdatedAt")
				if err != nil {
					span.RecordError(err)
					s.opt.Logger.AddMessage(log.ErrorLevel, err)
					return
				}
			}
		}

		if pushNotification.Status == statusx.ConvertStatusName(statusx.Active) {
			if currentTime.After(pushNotification.ScheduledAt) {
				// update status from active to finish
				pushNotification.Status = statusx.ConvertStatusName(statusx.Finished)
				err = s.RepositoryPushNotification.Update(ctx, &model.PushNotification{ID: pushNotification.ID, Status: pushNotification.Status, UpdatedAt: time.Now()}, "Status", "UpdatedAt")
				if err != nil {
					span.RecordError(err)
					s.opt.Logger.AddMessage(log.ErrorLevel, err)
					return
				}

			}
		}

		var redirect *dto.RedirectResponse
		switch pushNotification.RedirectTo {
		case 1:
			var item *catalog_service.GetItemDetailByInternalIdResponse
			item, err = s.opt.Client.CatalogServiceGrpc.GetItemDetailByInternalId(ctx, &catalog_service.GetItemDetailByInternalIdRequest{
				Id: utils.ToString(pushNotification.RedirectValue),
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("catalog", "item")
				return
			}
			redirect = &dto.RedirectResponse{
				To: pushNotification.RedirectTo,
				Value: dto.RedirectToItemResponse{
					ID:   item.Data.Id,
					Code: item.Data.Code,
					Name: item.Data.Description,
				},
				Name:      "Item",
				ValueName: "Item",
			}
		case 2:
			var itemCategory *catalogService.GetItemCategoryDetailResponse
			itemCategory, err = s.opt.Client.CatalogServiceGrpc.GetItemCategoryDetail(ctx, &catalogService.GetItemCategoryDetailRequest{
				Id: utils.ToInt64(pushNotification.RedirectValue),
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "item_category")
				return
			}
			redirect = &dto.RedirectResponse{
				To: pushNotification.RedirectTo,
				Value: dto.RedirectToItemCategoryResponse{
					ID:   itemCategory.Data.Id,
					Name: itemCategory.Data.Name,
					Code: itemCategory.Data.Code,
				},
				Name:      "Item Category",
				ValueName: "Item Category",
			}
		case 3:
			redirect = &dto.RedirectResponse{
				To:        pushNotification.RedirectTo,
				Value:     pushNotification.RedirectValue,
				Name:      "Cart",
				ValueName: "Cart",
			}
		case 4:
			redirect = &dto.RedirectResponse{
				To:        pushNotification.RedirectTo,
				Value:     pushNotification.RedirectValue,
				Name:      "URL",
				ValueName: "URL",
			}
		case 5:
			redirect = &dto.RedirectResponse{
				To:        pushNotification.RedirectTo,
				Value:     pushNotification.RedirectValue,
				Name:      "Promo",
				ValueName: "Promo",
			}
		case 6:
			redirect = &dto.RedirectResponse{
				To:        pushNotification.RedirectTo,
				Value:     pushNotification.RedirectValue,
				Name:      "Home",
				ValueName: "Home",
			}
		default:
			redirect = &dto.RedirectResponse{}
		}

		// get region name
		var regionNames []string
		for _, regionID := range utils.StringToStringArray(pushNotification.Regions) {
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
		for _, archetypeID := range utils.StringToStringArray(pushNotification.Archetypes) {
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

		res = append(res, &dto.PushNotificationResponse{
			ID:             pushNotification.ID,
			Code:           pushNotification.Code,
			CampaginName:   pushNotification.CampaignName,
			Regions:        utils.StringToStringArray(pushNotification.Regions),
			RegionNames:    regionNames,
			Archetypes:     utils.StringToStringArray(pushNotification.Archetypes),
			ArchetypeNames: archetypeNames,
			RedirectTo:     pushNotification.RedirectTo,
			RedirectValue:  pushNotification.RedirectValue,
			Redirect:       redirect,
			Title:          pushNotification.Title,
			Message:        pushNotification.Message,
			PushNow:        pushNotification.PushNow,
			ScheduledAt:    pushNotification.ScheduledAt,
			SuccessSent:    pushNotification.SuccessSent,
			FailedSent:     pushNotification.FailedSent,
			Opened:         pushNotification.Opened,
			CreatedAt:      timex.ToLocTime(ctx, pushNotification.CreatedAt),
			CreatedBy:      pushNotification.CreatedBy,
			UpdatedAt:      timex.ToLocTime(ctx, pushNotification.UpdatedAt),
			Status:         pushNotification.Status,
			StatusConvert:  statusx.ConvertStatusValue(pushNotification.Status),
		})
	}

	return
}

func (s *PushNotificationService) GetByID(ctx context.Context, id int64) (res dto.PushNotificationResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PushNotificationService.GetByID")
	defer span.End()

	var (
		pushNotification *model.PushNotification
		regionList       []*dto.RegionResponse
		archetypeList    []*dto.ArchetypeResponse
	)
	pushNotification, err = s.RepositoryPushNotification.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	currentTime := time.Now()
	if pushNotification.Status == statusx.ConvertStatusName(statusx.Draft) {
		if currentTime.After(pushNotification.ScheduledAt) {
			// update status from draft to active
			pushNotification.Status = statusx.ConvertStatusName(statusx.Active)
			err = s.RepositoryPushNotification.Update(ctx, &model.PushNotification{ID: pushNotification.ID, Status: pushNotification.Status, UpdatedAt: time.Now()}, "Status", "UpdatedAt")
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}
		}
	}

	if pushNotification.Status == statusx.ConvertStatusName(statusx.Active) {
		if currentTime.After(pushNotification.ScheduledAt) {
			// update status from active to finish
			pushNotification.Status = statusx.ConvertStatusName(statusx.Finished)
			err = s.RepositoryPushNotification.Update(ctx, &model.PushNotification{ID: pushNotification.ID, Status: pushNotification.Status, UpdatedAt: time.Now()}, "Status", "UpdatedAt")
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}

		}
	}

	var redirect *dto.RedirectResponse
	switch pushNotification.RedirectTo {
	case 1:
		var item *catalog_service.GetItemDetailByInternalIdResponse
		item, err = s.opt.Client.CatalogServiceGrpc.GetItemDetailByInternalId(ctx, &catalog_service.GetItemDetailByInternalIdRequest{
			Id: utils.ToString(pushNotification.RedirectValue),
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("catalog", "item")
			return
		}

		redirect = &dto.RedirectResponse{
			To: pushNotification.RedirectTo,
			Value: dto.RedirectToItemResponse{
				ID:   item.Data.Id,
				Code: item.Data.Code,
				Name: item.Data.Description,
			},
			Name:      "Item",
			ValueName: "Item",
		}
	case 2:
		var itemCategory *catalogService.GetItemCategoryDetailResponse
		itemCategory, err = s.opt.Client.CatalogServiceGrpc.GetItemCategoryDetail(ctx, &catalogService.GetItemCategoryDetailRequest{
			Id: utils.ToInt64(pushNotification.RedirectValue),
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "item_category")
			return
		}
		redirect = &dto.RedirectResponse{
			To: pushNotification.RedirectTo,
			Value: dto.RedirectToItemCategoryResponse{
				ID:   itemCategory.Data.Id,
				Name: itemCategory.Data.Name,
				Code: itemCategory.Data.Code,
			},
			Name:      "Item Category",
			ValueName: "Item Category",
		}
	case 3:
		redirect = &dto.RedirectResponse{
			To:        pushNotification.RedirectTo,
			Value:     pushNotification.RedirectValue,
			Name:      "Cart",
			ValueName: "Cart",
		}
	case 4:
		redirect = &dto.RedirectResponse{
			To:        pushNotification.RedirectTo,
			Value:     pushNotification.RedirectValue,
			Name:      "URL",
			ValueName: "URL",
		}
	case 5:
		redirect = &dto.RedirectResponse{
			To:        pushNotification.RedirectTo,
			Value:     pushNotification.RedirectValue,
			Name:      "Promo",
			ValueName: "Promo",
		}
	case 6:
		redirect = &dto.RedirectResponse{
			To:        pushNotification.RedirectTo,
			Value:     pushNotification.RedirectValue,
			Name:      "Home",
			ValueName: "Home",
		}
	default:
		redirect = &dto.RedirectResponse{}
	}

	// get region name
	var regionNames []string
	for _, regionID := range utils.StringToStringArray(pushNotification.Regions) {
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

		regionList = append(regionList, &dto.RegionResponse{
			ID:          region.Data[0].Region,
			Code:        region.Data[0].Region,
			Description: region.Data[0].Region,
		})
		regionNames = append(regionNames, region.Data[0].Region)
	}

	// get archetype name
	var archetypeNames []string
	for _, archetypeID := range utils.StringToStringArray(pushNotification.Archetypes) {
		var (
			archetype                           *bridge_service.GetArchetypeGPResponse
			statusArchetype, statusCustomerType int8
		)
		archetype, err = s.opt.Client.BridgeServiceGrpc.GetArchetypeGPDetail(ctx, &bridge_service.GetArchetypeGPDetailRequest{
			Id: archetypeID,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("archetype_id")
			return
		}

		var customerType *bridge_service.GetCustomerTypeGPResponse
		customerType, err = s.opt.Client.BridgeServiceGrpc.GetCustomerTypeGPDetail(ctx, &bridge_service.GetCustomerTypeGPDetailRequest{
			Id: archetype.Data[0].GnlCustTypeId,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		if archetype.Data[0].Inactive == 0 {
			statusArchetype = statusx.ConvertStatusName(statusx.Active)
		} else {
			statusArchetype = statusx.ConvertStatusName(statusx.Archived)
		}

		if customerType.Data[0].Inactive == 0 {
			statusCustomerType = statusx.ConvertStatusName(statusx.Active)
		} else {
			statusCustomerType = statusx.ConvertStatusName(statusx.Archived)
		}

		archetypeList = append(archetypeList, &dto.ArchetypeResponse{
			ID:             archetype.Data[0].GnlArchetypeId,
			Code:           archetype.Data[0].GnlArchetypeId,
			Description:    archetype.Data[0].GnlArchetypedescription,
			CustomerTypeID: archetype.Data[0].GnlCustTypeId,
			Status:         statusArchetype,
			ConvertStatus:  statusx.ConvertStatusValue(statusArchetype),
			CustomerType: &dto.CustomerTypeResponse{
				ID:            customerType.Data[0].GnL_Cust_Type_ID,
				Code:          customerType.Data[0].GnL_Cust_Type_ID,
				Description:   customerType.Data[0].GnL_CustType_Description,
				Status:        statusCustomerType,
				ConvertStatus: statusx.ConvertStatusValue(statusArchetype),
				CustomerGroup: customerType.Data[0].GnL_Cust_GroupDesc,
			},
		})
		archetypeNames = append(archetypeNames, archetype.Data[0].GnlArchetypedescription)
	}

	res = dto.PushNotificationResponse{
		ID:             pushNotification.ID,
		Code:           pushNotification.Code,
		CampaginName:   pushNotification.CampaignName,
		Regions:        utils.StringToStringArray(pushNotification.Regions),
		RegionNames:    regionNames,
		Archetypes:     utils.StringToStringArray(pushNotification.Archetypes),
		ArchetypeNames: archetypeNames,
		RedirectTo:     pushNotification.RedirectTo,
		RedirectValue:  pushNotification.RedirectValue,
		Redirect:       redirect,
		Title:          pushNotification.Title,
		Message:        pushNotification.Message,
		PushNow:        pushNotification.PushNow,
		ScheduledAt:    pushNotification.ScheduledAt,
		SuccessSent:    pushNotification.SuccessSent,
		FailedSent:     pushNotification.FailedSent,
		Opened:         pushNotification.Opened,
		CreatedAt:      timex.ToLocTime(ctx, pushNotification.CreatedAt),
		CreatedBy:      pushNotification.CreatedBy,
		UpdatedAt:      timex.ToLocTime(ctx, pushNotification.UpdatedAt),
		Status:         pushNotification.Status,
		StatusConvert:  statusx.ConvertStatusValue(pushNotification.Status),
		Region:         regionList,
		Archetype:      archetypeList,
	}

	return
}

func (s *PushNotificationService) Create(ctx context.Context, req dto.PushNotificationRequestCreate) (res dto.PushNotificationResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PushNotificationService.Create")
	defer span.End()

	// validate time
	if req.PushNow == 2 {
		currentTime := time.Now()

		if req.ScheduledAt.IsZero() {
			err = edenlabs.ErrorInvalid("scheduled_at")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		if req.ScheduledAt.Before(currentTime) {
			err = edenlabs.ErrorMustLater("scheduled_at", "current_time")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	} else {
		req.ScheduledAt = time.Now()
	}

	var redirectToName string
	switch req.RedirectTo {
	case 1:
		// item
		if req.RedirectValue == "" {
			err = edenlabs.ErrorRequired("redirect_value")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		_, err = s.opt.Client.BridgeServiceGrpc.GetItemDetail(ctx, &bridgeService.GetItemDetailRequest{
			Id: utils.ToInt64(req.RedirectValue),
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("redirect_value")
			return
		}
		redirectToName = "Item"

	case 2:
		// item category
		if req.RedirectValue == "" {
			err = edenlabs.ErrorRequired("redirect_value")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		if _, err = strconv.Atoi(req.RedirectValue); err != nil {
			err = edenlabs.ErrorInvalid("redirect_value")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		_, err = s.opt.Client.CatalogServiceGrpc.GetItemCategoryDetail(ctx, &catalogService.GetItemCategoryDetailRequest{
			Id: utils.ToInt64(req.RedirectValue),
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("redirect_value")
			return
		}
		redirectToName = "Item Category"
	case 3:
		// cart
		redirectToName = "Cart"
	case 4:
		// url
		if req.RedirectValue == "" {
			err = edenlabs.ErrorRequired("redirect_value")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		_, err = url.ParseRequestURI(req.RedirectValue)
		if err != nil {
			err = edenlabs.ErrorInvalid("redirect_value")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	case 5:
		// promo
		redirectToName = "Promo"
	case 6:
		// home
		redirectToName = "Home"
	default:
		err = edenlabs.ErrorRequired("redirect_to")
		return
	}

	// get region name
	var regionNames []string
	for _, regionID := range req.Regions {
		var region *bridgeService.GetAdmDivisionGPResponse
		region, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPDetail(ctx, &bridgeService.GetAdmDivisionGPDetailRequest{
			Region: regionID,
			Limit:  1,
			Offset: 0,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("region_id")
			return
		}
		regionNames = append(regionNames, region.Data[0].Region)
	}

	// get archetype name
	var archetypeNames []string
	for _, archetypeID := range req.Archetypes {
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

	var codeGenerator *configurationService.GetGenerateCodeResponse
	codeGenerator, err = s.opt.Client.ConfigurationServiceGrpc.GetGenerateCode(ctx, &configurationService.GetGenerateCodeRequest{
		Format: "PNT",
		Domain: "push_notification",
		Length: 6,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("configuration", "generate_code")
		return
	}
	code := codeGenerator.Data.Code

	pushNotification := &model.PushNotification{
		Code:          code,
		CampaignName:  req.CampaignName,
		Regions:       utils.ArrayStringToString(req.Regions),
		Archetypes:    utils.ArrayStringToString(req.Archetypes),
		RedirectTo:    int8(req.RedirectTo),
		RedirectValue: req.RedirectValue,
		Title:         req.Title,
		Message:       req.Message,
		PushNow:       req.PushNow,
		ScheduledAt:   req.ScheduledAt,
		SuccessSent:   0,
		FailedSent:    0,
		CreatedAt:     time.Now(),
		CreatedBy:     ctx.Value(constants.KeyUserID).(int64),
	}

	if req.PushNow == 1 {
		pushNotification.Status = statusx.ConvertStatusName(statusx.Active)
	} else {
		pushNotification.Status = statusx.ConvertStatusName(statusx.Draft)
	}

	err = s.RepositoryPushNotification.Create(ctx, pushNotification)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	userID := ctx.Value(constants.KeyUserID).(int64)
	_, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &auditService.CreateLogRequest{
		Log: &auditService.Log{
			UserId:      userID,
			ReferenceId: strconv.Itoa(int(pushNotification.ID)),
			Type:        "push_notification",
			Function:    "create",
			CreatedAt:   timestamppb.New(time.Now()),
		},
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if req.PushNow == 1 {
		// Set Default Customer
		var customerList *bridgeService.GetCustomerGPResponse
		customerList, err = s.opt.Client.BridgeServiceGrpc.GetCustomerGPList(ctx, &bridgeService.GetCustomerGPListRequest{
			Limit:  1000,
			Offset: 0,
			Phone:  "200001000011",
		})

		var userCustomers []*notification_service.UserCustomer
		for _, v := range customerList.Data {
			var customerInternal *crm_service.GetCustomerDetailResponse
			customerInternal, err = s.opt.Client.CrmServiceGrpc.GetCustomerDetail(ctx, &crm_service.GetCustomerDetailRequest{
				CustomerIdGp: v.Custnmbr,
			})
			if err != nil {
				continue
			}

			var userCustomer *customer_mobile_service.GetUserCustomerDetailResponse
			userCustomer, err = s.opt.Client.CustomerMobileServiceGrpc.GetUserCustomerDetail(ctx, &customer_mobile_service.GetUserCustomerDetailRequest{
				CustomerId: customerInternal.Data.Id,
			})
			if err != nil {
				continue
			}

			userCustomers = append(userCustomers, &notification_service.UserCustomer{
				CustomerId:     userCustomer.Data.CustomerId,
				UserCustomerId: userCustomer.Data.Id,
				FirebaseToken:  userCustomer.Data.FirebaseToken,
			})

		}

		var notificationStatus *notification_service.SendNotificationCampaignResponse
		notificationStatus, err = s.opt.Client.NotificationServiceGrpc.SendNotificationCampaign(ctx, &notification_service.SendNotificationCampaignRequest{
			NotificationCampaignId:   utils.ToString(pushNotification.ID),
			NotificationCampaignCode: pushNotification.Code,
			NotificationCampaignName: pushNotification.CampaignName,
			Title:                    pushNotification.Title,
			Message:                  pushNotification.Message,
			RedirectTo:               int64(pushNotification.RedirectTo),
			RedirectToName:           redirectToName,
			RedirectValue:            pushNotification.RedirectValue,
			UserCustomers:            userCustomers,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		pushNotification.Status = statusx.ConvertStatusName(statusx.Finished)
		pushNotification.SuccessSent = int(notificationStatus.Data.SuccessSent)
		pushNotification.FailedSent = int(notificationStatus.Data.FailedSent)

		err = s.RepositoryPushNotification.Update(ctx, pushNotification, "Status", "SuccessSent", "FailedSent")
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	}

	res = dto.PushNotificationResponse{
		ID:            pushNotification.ID,
		Code:          pushNotification.Code,
		CampaginName:  pushNotification.CampaignName,
		RedirectTo:    pushNotification.RedirectTo,
		RedirectValue: pushNotification.RedirectValue,
		Title:         pushNotification.Title,
		Message:       pushNotification.Message,
		PushNow:       pushNotification.PushNow,
		ScheduledAt:   pushNotification.ScheduledAt,
		SuccessSent:   pushNotification.SuccessSent,
		FailedSent:    pushNotification.FailedSent,
		Opened:        pushNotification.Opened,
		CreatedAt:     timex.ToLocTime(ctx, pushNotification.CreatedAt),
		CreatedBy:     pushNotification.CreatedBy,
		UpdatedAt:     timex.ToLocTime(ctx, pushNotification.UpdatedAt),
		Status:        pushNotification.Status,
		StatusConvert: statusx.ConvertStatusValue(pushNotification.Status),
	}

	return
}

func (s *PushNotificationService) Update(ctx context.Context, req dto.PushNotificationRequestUpdate, id int64) (res dto.PushNotificationResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PushNotificationService.Update")
	defer span.End()

	// validate is exist
	var (
		pushNotificationOld *model.PushNotification
		regionList          []*dto.RegionResponse
		archetypeList       []*dto.ArchetypeResponse
	)
	pushNotificationOld, err = s.RepositoryPushNotification.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// validate time
	if req.PushNow == 2 {
		currentTime := time.Now()

		if req.ScheduledAt.IsZero() {
			err = edenlabs.ErrorInvalid("scheduled_at")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		if req.ScheduledAt.Before(currentTime) {
			err = edenlabs.ErrorMustLater("scheduled_at", "current_time")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	}

	if pushNotificationOld.PushNow == 1 {
		err = edenlabs.ErrorMustDraft("status")
		return
	}

	if pushNotificationOld.Status != statusx.ConvertStatusName(statusx.Draft) {
		err = edenlabs.ErrorMustDraft("status")
		return
	}

	switch req.RedirectTo {
	case 1:
		// item
		if req.RedirectValue == "" {
			err = edenlabs.ErrorRequired("redirect_value")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		_, err = s.opt.Client.BridgeServiceGrpc.GetItemDetail(ctx, &bridgeService.GetItemDetailRequest{
			Id: utils.ToInt64(req.RedirectValue),
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("redirect_value")
			return
		}

	case 2:
		// item category
		if req.RedirectValue == "" {
			err = edenlabs.ErrorRequired("redirect_value")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		if _, err = strconv.Atoi(req.RedirectValue); err != nil {
			err = edenlabs.ErrorInvalid("redirect_value")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		_, err = s.opt.Client.CatalogServiceGrpc.GetItemCategoryDetail(ctx, &catalogService.GetItemCategoryDetailRequest{
			Id: utils.ToInt64(req.RedirectValue),
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("redirect_value")
			return
		}
	case 3:
		// cart
	case 4:
		// url
		if req.RedirectValue == "" {
			err = edenlabs.ErrorRequired("redirect_value")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		_, err = url.ParseRequestURI(req.RedirectValue)
		if err != nil {
			err = edenlabs.ErrorInvalid("redirect_value")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	case 5:
		// promo
	case 6:
		// home
	default:
		err = edenlabs.ErrorRequired("redirect_to")
		return
	}

	// get region name
	var regionNames []string
	for _, regionID := range req.Regions {
		var region *bridgeService.GetAdmDivisionGPResponse
		region, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPDetail(ctx, &bridgeService.GetAdmDivisionGPDetailRequest{
			Region: regionID,
			Limit:  1,
			Offset: 0,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("region_id")
			return
		}

		regionList = append(regionList, &dto.RegionResponse{
			ID:          utils.ToString(region.Data[0].Region),
			Code:        region.Data[0].Region,
			Description: region.Data[0].Region,
		})
		regionNames = append(regionNames, region.Data[0].Region)
	}

	// get archetype name
	var archetypeNames []string
	for _, archetypeID := range req.Archetypes {
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

		var customer *bridge_service.GetCustomerTypeGPResponse
		customer, err = s.opt.Client.BridgeServiceGrpc.GetCustomerTypeGPDetail(ctx, &bridge_service.GetCustomerTypeGPDetailRequest{
			Id: archetype.Data[0].GnlCustTypeId,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		archetypeList = append(archetypeList, &dto.ArchetypeResponse{
			ID:             archetype.Data[0].GnlArchetypeId,
			Code:           archetype.Data[0].GnlArchetypeId,
			Description:    archetype.Data[0].GnlArchetypedescription,
			CustomerTypeID: archetype.Data[0].GnlCustTypeId,
			CustomerType: &dto.CustomerTypeResponse{
				ID:          customer.Data[0].GnL_Cust_Type_ID,
				Code:        customer.Data[0].GnL_Cust_Type_ID,
				Description: customer.Data[0].GnL_CustType_Description,
			},
		})
		archetypeNames = append(archetypeNames, archetype.Data[0].GnlArchetypedescription)
	}

	pushNotification := &model.PushNotification{
		ID:            id,
		CampaignName:  req.CampaignName,
		Regions:       utils.ArrayStringToString(req.Regions),
		Archetypes:    utils.ArrayStringToString(req.Archetypes),
		RedirectTo:    int8(req.RedirectTo),
		RedirectValue: req.RedirectValue,
		Title:         req.Title,
		Message:       req.Message,
		PushNow:       req.PushNow,
		ScheduledAt:   req.ScheduledAt,
		UpdatedAt:     time.Now(),
	}
	err = s.RepositoryPushNotification.Update(ctx, pushNotification, "CampaignName", "Regions", "Archetypes", "RedirectTo", "RedirectValue", "Title", "Message", "PushNow", "ScheduledAt", "UpdatedAt")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	userID := ctx.Value(constants.KeyUserID).(int64)
	_, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &auditService.CreateLogRequest{
		Log: &auditService.Log{
			UserId:      userID,
			ReferenceId: strconv.Itoa(int(pushNotification.ID)),
			Type:        "push_notification",
			Function:    "update",
			CreatedAt:   timestamppb.New(time.Now()),
		},
	})

	res = dto.PushNotificationResponse{
		ID:             pushNotificationOld.ID,
		Code:           pushNotificationOld.Code,
		CampaginName:   pushNotification.CampaignName,
		Regions:        req.Regions,
		RegionNames:    regionNames,
		Archetypes:     req.Archetypes,
		ArchetypeNames: archetypeNames,
		RedirectTo:     pushNotification.RedirectTo,
		RedirectValue:  pushNotification.RedirectValue,
		Title:          pushNotification.Title,
		Message:        pushNotification.Message,
		PushNow:        pushNotification.PushNow,
		ScheduledAt:    pushNotification.ScheduledAt,
		SuccessSent:    pushNotificationOld.SuccessSent,
		FailedSent:     pushNotificationOld.FailedSent,
		Opened:         pushNotificationOld.Opened,
		CreatedAt:      timex.ToLocTime(ctx, pushNotificationOld.CreatedAt),
		CreatedBy:      pushNotification.CreatedBy,
		UpdatedAt:      timex.ToLocTime(ctx, pushNotification.UpdatedAt),
		Status:         pushNotificationOld.Status,
		StatusConvert:  statusx.ConvertStatusValue(pushNotificationOld.Status),
		Region:         regionList,
		Archetype:      archetypeList,
	}

	return
}

func (s *PushNotificationService) Cancel(ctx context.Context, req dto.PushNotificationRequestCancel, id int64) (res dto.PushNotificationResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PushNotificationService.Cancel")
	defer span.End()

	// validate pushNotification is exist
	var (
		pushNotificationOld *model.PushNotification
		regionList          []*dto.RegionResponse
		archetypeList       []*dto.ArchetypeResponse
	)
	pushNotificationOld, err = s.RepositoryPushNotification.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if pushNotificationOld.PushNow == 1 {
		err = edenlabs.ErrorMustDraft("status")
		return
	}

	if len(req.Note) > 100 {
		err = edenlabs.ErrorMustEqualOrLess("note", "100 characters")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// validate pushNotification time
	currentTime := time.Now()
	if currentTime.After(pushNotificationOld.ScheduledAt) {
		err = edenlabs.ErrorMustDraft("status")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if pushNotificationOld.Status != statusx.ConvertStatusName(statusx.Draft) {
		err = edenlabs.ErrorMustDraft("status")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// get region name
	var regionNames []string
	for _, regionID := range utils.StringToStringArray(pushNotificationOld.Regions) {
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

		regionList = append(regionList, &dto.RegionResponse{
			ID:          region.Data[0].Region,
			Code:        region.Data[0].Region,
			Description: region.Data[0].Region,
		})
		regionNames = append(regionNames, region.Data[0].Region)
	}

	// get archetype name
	var archetypeNames []string
	for _, archetypeID := range utils.StringToStringArray(pushNotificationOld.Archetypes) {
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

		var customer *bridge_service.GetCustomerTypeGPResponse
		customer, err = s.opt.Client.BridgeServiceGrpc.GetCustomerTypeGPDetail(ctx, &bridge_service.GetCustomerTypeGPDetailRequest{
			Id: archetype.Data[0].GnlCustTypeId,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		archetypeList = append(archetypeList, &dto.ArchetypeResponse{
			ID:             archetype.Data[0].GnlArchetypeId,
			Code:           archetype.Data[0].GnlArchetypeId,
			Description:    archetype.Data[0].GnlArchetypedescription,
			CustomerTypeID: archetype.Data[0].GnlCustTypeId,
			CustomerType: &dto.CustomerTypeResponse{
				ID:          customer.Data[0].GnL_Cust_Type_ID,
				Code:        customer.Data[0].GnL_Cust_Type_ID,
				Description: customer.Data[0].GnL_CustType_Description,
			},
		})
		archetypeNames = append(archetypeNames, archetype.Data[0].GnlArchetypedescription)
	}

	err = s.RepositoryPushNotification.Update(ctx, &model.PushNotification{
		ID:        id,
		Status:    statusx.ConvertStatusName(statusx.Cancelled),
		UpdatedAt: time.Now(),
	}, "Status", "UpdatedAt")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	userID := ctx.Value(constants.KeyUserID).(int64)
	_, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &auditService.CreateLogRequest{
		Log: &auditService.Log{
			UserId:      userID,
			ReferenceId: strconv.Itoa(int(pushNotificationOld.ID)),
			Type:        "push_notification",
			Function:    "cancel",
			CreatedAt:   timestamppb.New(time.Now()),
			Note:        req.Note,
		},
	})

	res = dto.PushNotificationResponse{
		ID:             pushNotificationOld.ID,
		Code:           pushNotificationOld.Code,
		CampaginName:   pushNotificationOld.CampaignName,
		RedirectTo:     pushNotificationOld.RedirectTo,
		RedirectValue:  pushNotificationOld.RedirectValue,
		Archetypes:     utils.StringToStringArray(pushNotificationOld.Archetypes),
		Regions:        utils.StringToStringArray(pushNotificationOld.Regions),
		Title:          pushNotificationOld.Title,
		Message:        pushNotificationOld.Message,
		PushNow:        pushNotificationOld.PushNow,
		ScheduledAt:    pushNotificationOld.ScheduledAt,
		SuccessSent:    pushNotificationOld.SuccessSent,
		FailedSent:     pushNotificationOld.FailedSent,
		Opened:         pushNotificationOld.Opened,
		CreatedAt:      timex.ToLocTime(ctx, pushNotificationOld.CreatedAt),
		CreatedBy:      pushNotificationOld.CreatedBy,
		UpdatedAt:      timex.ToLocTime(ctx, pushNotificationOld.UpdatedAt),
		Status:         statusx.ConvertStatusName(statusx.Cancelled),
		StatusConvert:  statusx.ConvertStatusValue(statusx.ConvertStatusName(statusx.Cancelled)),
		RegionNames:    regionNames,
		Region:         regionList,
		ArchetypeNames: archetypeNames,
		Archetype:      archetypeList,
	}

	return
}

func (s *PushNotificationService) UpdateOpened(ctx context.Context, req *dto.PushNotificationRequestUpdateOpened) (err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PushNotificationService.UpdateOpened")
	defer span.End()

	pushNotification := &model.PushNotification{
		ID:     req.ID,
		Opened: req.Opened,
	}

	err = s.RepositoryPushNotification.Update(ctx, pushNotification, "Opened")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	_, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &auditService.CreateLogRequest{
		Log: &auditService.Log{
			UserId:      22,
			ReferenceId: strconv.Itoa(int(pushNotification.ID)),
			Type:        "push_notification",
			Function:    "update",
			CreatedAt:   timestamppb.New(time.Now()),
		},
	})

	return
}

func (s *PushNotificationService) GetDetailMobile(ctx context.Context, id int64) (res *dto.PushNotificationResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PushNotificationService.GetDetailMobile")
	defer span.End()

	var pushNotification *model.PushNotification

	pushNotification, err = s.RepositoryPushNotification.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &dto.PushNotificationResponse{
		ID:            pushNotification.ID,
		Code:          pushNotification.Code,
		CampaginName:  pushNotification.CampaignName,
		Regions:       utils.StringToStringArray(pushNotification.Regions),
		Archetypes:    utils.StringToStringArray(pushNotification.Archetypes),
		RedirectTo:    pushNotification.RedirectTo,
		RedirectValue: pushNotification.RedirectValue,
		Title:         pushNotification.Title,
		Message:       pushNotification.Message,
		PushNow:       pushNotification.PushNow,
		ScheduledAt:   pushNotification.ScheduledAt,
		SuccessSent:   pushNotification.SuccessSent,
		FailedSent:    pushNotification.FailedSent,
		Opened:        pushNotification.Opened,
		CreatedAt:     timex.ToLocTime(ctx, pushNotification.CreatedAt),
		CreatedBy:     pushNotification.CreatedBy,
		UpdatedAt:     timex.ToLocTime(ctx, pushNotification.UpdatedAt),
		Status:        pushNotification.Status,
		StatusConvert: statusx.ConvertStatusValue(pushNotification.Status),
	}
	return
}
