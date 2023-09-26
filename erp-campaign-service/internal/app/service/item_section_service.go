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
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/repository"
	auditService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/audit_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/catalog_service"
	configurationService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IItemSectionService interface {
	Get(ctx context.Context, req *dto.ItemSectionRequestGet) (res []*dto.ItemSectionResponse, total int64, err error)
	GetByID(ctx context.Context, id int64) (res dto.ItemSectionResponse, err error)
	Create(ctx context.Context, req dto.ItemSectionRequestCreate) (res dto.ItemSectionResponse, err error)
	Update(ctx context.Context, req dto.ItemSectionRequestUpdate, id int64) (res dto.ItemSectionResponse, err error)
	Archive(ctx context.Context, id int64, req dto.ItemSectionRequestArchive) (res dto.ItemSectionResponse, err error)
	GetListMobile(ctx context.Context, req *dto.ItemSectionRequestGet) (res []*dto.ItemSectionResponse, total int64, err error)
	GetDetailMobile(ctx context.Context, id int64) (res *dto.ItemSectionResponse, err error)
}

type ItemSectionService struct {
	opt                   opt.Options
	RepositoryItemSection repository.IItemSectionRepository
}

func NewItemSectionService() IItemSectionService {
	return &ItemSectionService{
		opt:                   global.Setup.Common,
		RepositoryItemSection: repository.NewItemSectionRepository(),
	}
}

func (s *ItemSectionService) Get(ctx context.Context, req *dto.ItemSectionRequestGet) (res []*dto.ItemSectionResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemSectionService.Get")
	defer span.End()

	var itemSections []*model.ItemSection
	itemSections, total, err = s.RepositoryItemSection.Get(ctx, req)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	currentTime := time.Now()
	for _, itemSection := range itemSections {
		if itemSection.Status == statusx.ConvertStatusName("Draft") {
			if currentTime.After(itemSection.StartAt) {
				// update status from draft to active
				err = s.RepositoryItemSection.Update(ctx, &model.ItemSection{ID: itemSection.ID, Status: 1, UpdatedAt: time.Now()}, "Status", "UpdatedAt")
				if err != nil {
					span.RecordError(err)
					s.opt.Logger.AddMessage(log.ErrorLevel, err)
					return
				}
				itemSection.Status = 1
			}
		}

		if itemSection.Status == 1 {
			if currentTime.After(itemSection.FinishAt) {
				// update status from active to finish
				err = s.RepositoryItemSection.Update(ctx, &model.ItemSection{ID: itemSection.ID, Status: statusx.ConvertStatusName("Finished"), UpdatedAt: time.Now()}, "Status", "UpdatedAt")
				if err != nil {
					span.RecordError(err)
					s.opt.Logger.AddMessage(log.ErrorLevel, err)
					return
				}
				itemSection.Status = statusx.ConvertStatusName("Draft")
			}
		}

		// get region name
		var regionNames []string
		for _, regionID := range utils.StringToStringArray(itemSection.Regions) {
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
		for _, archetypeID := range utils.StringToStringArray(itemSection.Archetypes) {
			var archetype *bridgeService.GetArchetypeGPResponse
			archetype, err = s.opt.Client.BridgeServiceGrpc.GetArchetypeGPDetail(ctx, &bridgeService.GetArchetypeGPDetailRequest{
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

		res = append(res, &dto.ItemSectionResponse{
			ID:              itemSection.ID,
			Code:            itemSection.Code,
			Name:            itemSection.Name,
			BackgroundImage: itemSection.BackgroundImage,
			StartAt:         itemSection.StartAt,
			FinishAt:        itemSection.FinishAt,
			Regions:         utils.StringToStringArray(itemSection.Regions),
			RegionNames:     regionNames,
			Archetypes:      utils.StringToStringArray(itemSection.Archetypes),
			ArchetypeNames:  archetypeNames,
			// Items:           items,
			Sequence:  itemSection.Sequence,
			Note:      itemSection.Note,
			Status:    itemSection.Status,
			CreatedAt: timex.ToLocTime(ctx, itemSection.CreatedAt),
			UpdatedAt: timex.ToLocTime(ctx, itemSection.UpdatedAt),
			Type:      itemSection.Type,
			ItemID:    utils.StringToInt64Array(itemSection.Items),
		})
	}

	return
}

func (s *ItemSectionService) GetByID(ctx context.Context, id int64) (res dto.ItemSectionResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemSectionService.GetByID")
	defer span.End()

	var (
		itemSection                 *model.ItemSection
		regions                     []*dto.RegionResponse
		archetypes                  []*dto.ArchetypeResponse
		items                       []*dto.ItemResponse
		regionNames, archetypeNames []string
	)
	itemSection, err = s.RepositoryItemSection.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	currentTime := time.Now()

	if itemSection.Status == statusx.ConvertStatusName("Draft") {
		if currentTime.After(itemSection.StartAt) {
			// update status from draft to active
			err = s.RepositoryItemSection.Update(ctx, &model.ItemSection{ID: itemSection.ID, Status: 1, UpdatedAt: time.Now()}, "Status", "UpdatedAt")
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}
			itemSection.Status = 1
		}
	}

	if itemSection.Status == 1 {
		if currentTime.After(itemSection.FinishAt) {
			// update status from active to finish
			err = s.RepositoryItemSection.Update(ctx, &model.ItemSection{ID: itemSection.ID, Status: statusx.ConvertStatusName("Finished"), UpdatedAt: time.Now()}, "Status", "UpdatedAt")
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}
			itemSection.Status = statusx.ConvertStatusName("Finished")
		}
	}

	itemsID := utils.StringToInt64Array(itemSection.Items)
	for _, itemID := range itemsID {

		var item *catalog_service.GetItemDetailByInternalIdResponse
		item, err = s.opt.Client.CatalogServiceGrpc.GetItemDetailByInternalId(ctx, &catalog_service.GetItemDetailByInternalIdRequest{
			Id: utils.ToString(itemID),
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("catalog", "item")
			return
		}
		items = append(items, &dto.ItemResponse{
			ID:          itemID,
			Code:        item.Data.Code,
			Description: item.Data.Description,
			Status:      int8(item.Data.Status),
		})
	}

	// get region name
	for _, regionID := range utils.StringToStringArray(itemSection.Regions) {
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

		regions = append(regions, &dto.RegionResponse{
			ID:          utils.ToString(region.Data[0].Region),
			Code:        region.Data[0].Region,
			Description: region.Data[0].Region,
		})
		regionNames = append(regionNames, region.Data[0].Region)
	}

	// get archetype name
	for _, archetypeID := range utils.StringToStringArray(itemSection.Archetypes) {
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

		archetypes = append(archetypes, &dto.ArchetypeResponse{
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

	res = dto.ItemSectionResponse{
		ID:              itemSection.ID,
		Name:            itemSection.Name,
		Code:            itemSection.Code,
		BackgroundImage: itemSection.BackgroundImage,
		StartAt:         itemSection.StartAt,
		FinishAt:        itemSection.FinishAt,
		Regions:         utils.StringToStringArray(itemSection.Regions),
		RegionNames:     regionNames,
		Archetypes:      utils.StringToStringArray(itemSection.Archetypes),
		ArchetypeNames:  archetypeNames,
		Items:           items,
		Sequence:        itemSection.Sequence,
		Note:            itemSection.Note,
		Status:          itemSection.Status,
		CreatedAt:       timex.ToLocTime(ctx, itemSection.CreatedAt),
		UpdatedAt:       timex.ToLocTime(ctx, itemSection.UpdatedAt),
		Type:            itemSection.Type,
		Region:          regions,
		Archetype:       archetypes,
	}

	return
}
func (s *ItemSectionService) Create(ctx context.Context, req dto.ItemSectionRequestCreate) (res dto.ItemSectionResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemSectionService.Create")
	defer span.End()

	// validate sequence limitation
	if req.Sequence < 1 || req.Sequence > 5 {
		err = edenlabs.ErrorValidation("sequence", "The sequence is invalid value")
		return
	}

	// validate item
	if len(req.Items) == 0 {
		err = edenlabs.ErrorValidation("items", "The items must be filled in")
		return
	}

	// validate start_at not greater than time.now
	if req.StartAt.Before(time.Now()) {
		err = edenlabs.ErrorMustGreater("start_at", "time now")
		return
	}

	// validate end_at not later than time.now and start_at
	if req.FinishAt.Before(time.Now()) || req.FinishAt.Before(req.StartAt) || req.FinishAt.Equal(req.StartAt) {
		err = edenlabs.ErrorMustGreater("finish_at", "time now or start at")
		return
	}

	// only check if is prod recommendation is false
	if req.Type != 2 {
		if req.BackgroundImage == "" {
			err = edenlabs.ErrorRequired("background_image")
			return
		}

		// validate background image is url
		_, err = url.ParseRequestURI(req.BackgroundImage)
		if err != nil {
			err = edenlabs.ErrorValidation("background_image", "The background image is invalid value")
			return
		}

		req.Type = 1
	} else {
		req.BackgroundImage = ""
	}
	// Validation Region
	for _, regionID := range req.Regions {
		_, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPDetail(ctx, &bridgeService.GetAdmDivisionGPDetailRequest{
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
	}

	// Validation Archetype
	for _, archetypeID := range req.Archetypes {
		_, err = s.opt.Client.BridgeServiceGrpc.GetArchetypeGPDetail(ctx, &bridge_service.GetArchetypeGPDetailRequest{
			Id: archetypeID,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("archetype_id")
			return
		}
	}

	// only check if it is product recommendation
	if req.Type == 2 {
		// check if there are already data in between the date range
		var isExist bool
		if isExist, err = s.RepositoryItemSection.CheckIsIntersect(ctx, req.Type, req.StartAt.Format("2006-01-02 15:04:05"), req.FinishAt.Format("2006-01-02 15:04:05"), 0); err != nil || isExist {
			err = edenlabs.ErrorValidation("start_at", "There are already active or draft data between "+req.StartAt.Format("2006-01-02 15:04:05")+" and "+req.FinishAt.Format("2006-01-02 15:04:05"))
			return
		}
	}

	var codeGenerator *configurationService.GetGenerateCodeResponse
	codeGenerator, err = s.opt.Client.ConfigurationServiceGrpc.GetGenerateCode(ctx, &configurationService.GetGenerateCodeRequest{
		Format: "ISC",
		Domain: "item_section",
		Length: 6,
	})

	itemSection := &model.ItemSection{
		Name:            req.Name,
		Code:            codeGenerator.Data.Code,
		BackgroundImage: req.BackgroundImage,
		StartAt:         req.StartAt,
		FinishAt:        req.FinishAt,
		Regions:         utils.ArrayStringToString(req.Regions),
		Archetypes:      utils.ArrayStringToString(req.Archetypes),
		Items:           utils.ArrayInt64ToString(req.Items),
		Sequence:        req.Sequence,
		Note:            req.Note,
		CreatedAt:       time.Now(),
		Status:          statusx.ConvertStatusName("Draft"),
		Type:            req.Type,
	}

	span.AddEvent("creating itemSection")
	itemSectionID, err := s.RepositoryItemSection.Create(ctx, itemSection)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	userID := ctx.Value(constants.KeyUserID).(int64)

	_, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &auditService.CreateLogRequest{
		Log: &auditService.Log{
			UserId:      userID,
			ReferenceId: strconv.Itoa(int(itemSectionID)),
			Type:        "item_section",
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

	res = dto.ItemSectionResponse{
		ID:              itemSection.ID,
		Code:            itemSection.Code,
		Name:            itemSection.Name,
		BackgroundImage: itemSection.BackgroundImage,
		StartAt:         itemSection.StartAt,
		FinishAt:        itemSection.FinishAt,
		Regions:         utils.StringToStringArray(itemSection.Regions),
		Archetypes:      utils.StringToStringArray(itemSection.Archetypes),
		Sequence:        itemSection.Sequence,
		Note:            itemSection.Note,
		Status:          itemSection.Status,
		CreatedAt:       timex.ToLocTime(ctx, itemSection.CreatedAt),
		UpdatedAt:       timex.ToLocTime(ctx, itemSection.UpdatedAt),
		Type:            itemSection.Type,
	}

	return
}

func (s *ItemSectionService) Update(ctx context.Context, req dto.ItemSectionRequestUpdate, id int64) (res dto.ItemSectionResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemSectionService.Update")
	defer span.End()

	// validate data is exist
	var (
		itemSectionOld              *model.ItemSection
		regions                     []*dto.RegionResponse
		archetypes                  []*dto.ArchetypeResponse
		items                       []*dto.ItemResponse
		regionNames, archetypeNames []string
	)

	itemSectionOld, err = s.RepositoryItemSection.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// validate product section must be draft
	if itemSectionOld.Status != statusx.ConvertStatusName("Draft") {
		err = edenlabs.ErrorValidation("status", "The item section must be draft")
		return
	}

	// validate sequence limitation
	if req.Sequence < 1 || req.Sequence > 5 {
		err = edenlabs.ErrorValidation("sequence", "The sequence is invalid value")
		return
	}

	// validate start_at not greater than time.now
	if req.StartAt.Before(time.Now()) {
		err = edenlabs.ErrorMustGreater("start_at", "time now")
		return
	}

	// validate end_at not later than time.now and start_at
	if req.FinishAt.Before(time.Now()) || req.FinishAt.Before(req.StartAt) || req.FinishAt.Equal(req.StartAt) {
		err = edenlabs.ErrorMustGreater("finish_at", "time now or start at")
		return
	}

	// only check if is prod recommendation is false
	if itemSectionOld.Type != 2 {
		if req.BackgroundImage == "" {
			err = edenlabs.ErrorRequired("background_image")
			return
		}

		// validate background image is url
		_, err = url.ParseRequestURI(req.BackgroundImage)
		if err != nil {
			err = edenlabs.ErrorValidation("background_image", "The background image is invalid value")
			return
		}

		itemSectionOld.Type = 1
	} else {
		req.BackgroundImage = ""
	}

	// only check if it is product recommendation
	if itemSectionOld.Type == 2 {
		// check if there are already data in between the date range
		var isExist bool
		if isExist, err = s.RepositoryItemSection.CheckIsIntersect(ctx, itemSectionOld.Type, req.StartAt.Format("2006-01-02 15:04:05"), req.FinishAt.Format("2006-01-02 15:04:05"), itemSectionOld.ID); err != nil || isExist {
			err = edenlabs.ErrorValidation("start_at", "There are already active or draft data between "+req.StartAt.Format("2006-01-02 15:04:05")+" and "+req.FinishAt.Format("2006-01-02 15:04:05"))
			return
		}
	}

	for _, itemID := range req.Items {
		var item *catalog_service.GetItemDetailByInternalIdResponse
		item, err = s.opt.Client.CatalogServiceGrpc.GetItemDetailByInternalId(ctx, &catalog_service.GetItemDetailByInternalIdRequest{
			Id: utils.ToString(itemID),
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("catalog", "item")
			return
		}
		items = append(items, &dto.ItemResponse{
			ID:          itemID,
			Code:        item.Data.Code,
			Description: item.Data.Description,
			Status:      int8(item.Data.Status),
		})
	}

	// get region name
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
			err = edenlabs.ErrorRpcNotFound("bridge", "region")
			return
		}

		regions = append(regions, &dto.RegionResponse{
			ID:          utils.ToString(region.Data[0].Region),
			Code:        region.Data[0].Region,
			Description: region.Data[0].Region,
		})
		regionNames = append(regionNames, region.Data[0].Region)
	}

	// get archetype name
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

		archetypes = append(archetypes, &dto.ArchetypeResponse{
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

	userID := ctx.Value(constants.KeyUserID).(int64)

	_, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &auditService.CreateLogRequest{
		Log: &auditService.Log{
			UserId:      userID,
			ReferenceId: strconv.Itoa(int(id)),
			Type:        "item_section",
			Function:    "update",
			CreatedAt:   timestamppb.New(time.Now()),
			Note:        req.Note,
		},
	})

	itemSection := &model.ItemSection{
		ID:              id,
		Name:            req.Name,
		BackgroundImage: req.BackgroundImage,
		StartAt:         req.StartAt,
		FinishAt:        req.FinishAt,
		Regions:         utils.ArrayStringToString(req.Regions),
		Archetypes:      utils.ArrayStringToString(req.Archetypes),
		Items:           utils.ArrayInt64ToString(req.Items),
		Sequence:        req.Sequence,
		Note:            req.Note,
		UpdatedAt:       time.Now(),
	}

	span.AddEvent("updating itemSection")
	err = s.RepositoryItemSection.Update(ctx, itemSection, "Name", "BackgroundImage", "StartAt", "FinishAt", "Regions", "Archetypes", "Items", "Sequence", "Note", "UpdatedAt")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.ItemSectionResponse{
		ID:              itemSectionOld.ID,
		Code:            itemSectionOld.Code,
		Name:            itemSection.Name,
		BackgroundImage: itemSection.BackgroundImage,
		StartAt:         itemSection.StartAt,
		FinishAt:        itemSection.FinishAt,
		Regions:         utils.StringToStringArray(itemSection.Regions),
		Archetypes:      utils.StringToStringArray(itemSection.Archetypes),
		Sequence:        itemSection.Sequence,
		Note:            itemSection.Note,
		Status:          itemSection.Status,
		CreatedAt:       timex.ToLocTime(ctx, itemSectionOld.CreatedAt),
		UpdatedAt:       timex.ToLocTime(ctx, itemSection.UpdatedAt),
		Type:            itemSectionOld.Type,
		Items:           items,
		Region:          regions,
		Archetype:       archetypes,
	}

	return
}

func (s *ItemSectionService) Archive(ctx context.Context, id int64, req dto.ItemSectionRequestArchive) (res dto.ItemSectionResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemSectionService.Archive")
	defer span.End()

	// validate data is exist
	var itemSectionOld *model.ItemSection
	itemSectionOld, err = s.RepositoryItemSection.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if itemSectionOld.Status == statusx.ConvertStatusName("Archived") {
		err = edenlabs.ErrorValidation("status", "The status has been archived")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	currentTime := time.Now()
	if currentTime.After(itemSectionOld.FinishAt) || (itemSectionOld.Status != statusx.ConvertStatusName("Active") && itemSectionOld.Status != statusx.ConvertStatusName("Draft")) {
		err = edenlabs.ErrorValidation("status", "The status must be active or draft")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if len(req.Note) > 100 {
		err = edenlabs.ErrorMustEqualOrLess("note", "100 characters")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	itemSection := &model.ItemSection{
		ID:        id,
		Status:    statusx.ConvertStatusName("Archived"),
		Note:      req.Note,
		UpdatedAt: time.Now(),
	}

	err = s.RepositoryItemSection.Update(ctx, itemSection, "Status", "Note", "UpdatedAt")
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
			Type:        "item_section",
			Function:    "archive",
			CreatedAt:   timestamppb.New(time.Now()),
			Note:        req.Note,
		},
	})

	res = dto.ItemSectionResponse{
		ID:        itemSection.ID,
		Note:      itemSection.Note,
		UpdatedAt: timex.ToLocTime(ctx, itemSection.UpdatedAt),
		Status:    itemSection.Status,
	}

	return
}

func (s *ItemSectionService) GetListMobile(ctx context.Context, req *dto.ItemSectionRequestGet) (res []*dto.ItemSectionResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemSectionService.GetListMobile")
	defer span.End()

	var itemSections []*model.ItemSection
	itemSections, total, err = s.RepositoryItemSection.Get(ctx, req)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, itemSection := range itemSections {
		res = append(res, &dto.ItemSectionResponse{
			ID:              itemSection.ID,
			Code:            itemSection.Code,
			Name:            itemSection.Name,
			BackgroundImage: itemSection.BackgroundImage,
			StartAt:         itemSection.StartAt,
			FinishAt:        itemSection.FinishAt,
			Regions:         utils.StringToStringArray(itemSection.Regions),
			Archetypes:      utils.StringToStringArray(itemSection.Archetypes),
			Sequence:        itemSection.Sequence,
			Note:            itemSection.Note,
			Status:          itemSection.Status,
			CreatedAt:       timex.ToLocTime(ctx, itemSection.CreatedAt),
			UpdatedAt:       timex.ToLocTime(ctx, itemSection.UpdatedAt),
			Type:            itemSection.Type,
			ItemID:          utils.StringToInt64Array(itemSection.Items),
		})
	}

	return
}

func (s *ItemSectionService) GetDetailMobile(ctx context.Context, id int64) (res *dto.ItemSectionResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemSectionService.GetDetailMobile")
	defer span.End()

	var itemSection *model.ItemSection
	itemSection, err = s.RepositoryItemSection.GetByID(ctx, id)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &dto.ItemSectionResponse{
		ID:              itemSection.ID,
		Code:            itemSection.Code,
		Name:            itemSection.Name,
		BackgroundImage: itemSection.BackgroundImage,
		StartAt:         itemSection.StartAt,
		FinishAt:        itemSection.FinishAt,
		Regions:         utils.StringToStringArray(itemSection.Regions),
		Archetypes:      utils.StringToStringArray(itemSection.Archetypes),
		Sequence:        itemSection.Sequence,
		Note:            itemSection.Note,
		Status:          itemSection.Status,
		CreatedAt:       timex.ToLocTime(ctx, itemSection.CreatedAt),
		UpdatedAt:       timex.ToLocTime(ctx, itemSection.UpdatedAt),
		Type:            itemSection.Type,
		ItemID:          utils.StringToInt64Array(itemSection.Items),
	}

	return
}
