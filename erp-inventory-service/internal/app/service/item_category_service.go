package service

import (
	"context"
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
	configurationService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IItemCategoryService interface {
	Get(ctx context.Context, offset, limit, status int, search, orderBy string, regionID string) (res []dto.ItemCategoryResponse, total int64, err error)
	GetByID(ctx context.Context, id int64) (res dto.ItemCategoryResponse, err error)
	UpdateImage(ctx context.Context, req dto.ItemCategoryImageRequestUpdate, id int64) (res dto.ItemCategoryResponse, err error)
	Archive(ctx context.Context, id int64) (res dto.ItemCategoryResponse, err error)
	Unarchive(ctx context.Context, id int64) (res dto.ItemCategoryResponse, err error)
	Create(ctx context.Context, req dto.ItemCategoryRequestCreate) (res dto.ItemCategoryResponse, err error)
	Update(ctx context.Context, req dto.ItemCategoryRequestUpdate, id int64) (res dto.ItemCategoryResponse, err error)
}

type ItemCategoryService struct {
	opt                         opt.Options
	RepositoryItemCategory      repository.IItemCategoryRepository
	RepositoryItemCategoryImage repository.IItemCategoryImageRepository
}

func NewItemCategoryService() IItemCategoryService {
	return &ItemCategoryService{
		opt:                         global.Setup.Common,
		RepositoryItemCategory:      repository.NewItemCategoryRepository(),
		RepositoryItemCategoryImage: repository.NewItemCategoryImageRepository(),
	}
}

func (s *ItemCategoryService) Get(ctx context.Context, offset, limit, status int, search, orderBy string, regionID string) (res []dto.ItemCategoryResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemCategoryService.Get")
	defer span.End()

	var (
		itemCategoryCategorys []*model.ItemCategory
		detailRegion          *bridgeService.GetAdmDivisionGPResponse
	)

	itemCategoryCategorys, total, err = s.RepositoryItemCategory.Get(ctx, offset, limit, status, search, orderBy, regionID)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, itemCategory := range itemCategoryCategorys {
		var region string
		if itemCategory.Regions != "" {
			listIDArrStr := strings.Split(itemCategory.Regions, ",")
			for _, v := range listIDArrStr {
				detailRegion, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPDetail(ctx, &bridgeService.GetAdmDivisionGPDetailRequest{
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
				region += detailRegion.Data[0].Region + ","
			}

			region = strings.TrimSuffix(region, ",")
		}

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

		res = append(res, dto.ItemCategoryResponse{
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

	return
}

func (s *ItemCategoryService) GetByID(ctx context.Context, id int64) (res dto.ItemCategoryResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemCategoryService.GetByID")
	defer span.End()

	var (
		itemCategory *model.ItemCategory
		detailRegion *bridgeService.GetAdmDivisionGPResponse
		regionList   []*dto.RegionResponse
	)
	itemCategory, err = s.RepositoryItemCategory.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("item_category_id")
		return
	}

	var region string
	if itemCategory.Regions != "" {
		listIDArrStr := strings.Split(itemCategory.Regions, ",")
		for _, v := range listIDArrStr {
			detailRegion, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPDetail(ctx, &bridgeService.GetAdmDivisionGPDetailRequest{
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
			region += detailRegion.Data[0].Region + ","

			regionList = append(regionList, &dto.RegionResponse{
				ID:          detailRegion.Data[0].Region,
				Code:        detailRegion.Data[0].Region,
				Description: detailRegion.Data[0].Region,
			})
		}
	}

	region = strings.TrimSuffix(region, ",")

	res = dto.ItemCategoryResponse{
		ID:            itemCategory.ID,
		Code:          itemCategory.Code,
		Name:          itemCategory.Name,
		RegionID:      itemCategory.Regions,
		Region:        region,
		Status:        itemCategory.Status,
		StatusConvert: statusx.ConvertStatusValue(itemCategory.Status),
		CreatedAt:     itemCategory.CreatedAt,
		UpdatedAt:     itemCategory.UpdatedAt,
		Regions:       regionList,
	}

	var itemCategoryImage *model.ItemCategoryImage
	itemCategoryImage, err = s.RepositoryItemCategoryImage.GetByItemCategoryID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if itemCategoryImage.ID != 0 {
		res.ItemCategoryImage = &dto.ItemCategoryImageResponse{
			ID:             itemCategoryImage.ID,
			ItemCategoryID: itemCategoryImage.ItemCategoryID,
			ImageUrl:       itemCategoryImage.ImageUrl,
			CreatedAt:      itemCategoryImage.CreatedAt,
		}
	}

	return
}

func (s *ItemCategoryService) UpdateImage(ctx context.Context, req dto.ItemCategoryImageRequestUpdate, id int64) (res dto.ItemCategoryResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemCategoryService.UpdateImage")
	defer span.End()

	// validate itemCategory id
	_, err = s.RepositoryItemCategory.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	err = s.RepositoryItemCategoryImage.DeleteByItemID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	itemCategoryImage := &model.ItemCategoryImage{
		ItemCategoryID: id,
		ImageUrl:       req.ImageUrl,
		CreatedAt:      time.Now(),
	}

	err = s.RepositoryItemCategoryImage.Create(ctx, itemCategoryImage)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.ItemCategoryResponse{
		ID: id,
		ItemCategoryImage: &dto.ItemCategoryImageResponse{
			ID:             itemCategoryImage.ID,
			ItemCategoryID: itemCategoryImage.ItemCategoryID,
			ImageUrl:       itemCategoryImage.ImageUrl,
			CreatedAt:      itemCategoryImage.CreatedAt,
		},
	}

	return
}

func (s *ItemCategoryService) Archive(ctx context.Context, id int64) (res dto.ItemCategoryResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemCategoryService.Archive")
	defer span.End()

	var ItemCategoryOld *model.ItemCategory
	ItemCategoryOld, err = s.RepositoryItemCategory.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("item_category_id")
		return
	}

	if ItemCategoryOld.Status != statusx.ConvertStatusName("Active") {
		err = edenlabs.ErrorMustActive("status")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	ItemCategory := &model.ItemCategory{
		ID:        ItemCategoryOld.ID,
		Status:    statusx.ConvertStatusName("Archived"),
		UpdatedAt: time.Now(),
	}

	err = s.RepositoryItemCategory.Update(ctx, ItemCategory, "Status", "UpdatedAt")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	userID := ctx.Value(constants.KeyUserID).(int64)

	_, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &auditService.CreateLogRequest{
		Log: &auditService.Log{
			UserId:      userID,
			ReferenceId: strconv.Itoa(int(ItemCategory.ID)),
			Type:        "item_category",
			Function:    "archive",
			CreatedAt:   timestamppb.New(time.Now()),
		},
	})

	res = dto.ItemCategoryResponse{
		ID:            ItemCategoryOld.ID,
		Code:          ItemCategory.Code,
		Name:          ItemCategoryOld.Name,
		RegionID:      ItemCategoryOld.Regions,
		Status:        statusx.ConvertStatusName("Archived"),
		StatusConvert: statusx.Archived,
	}

	return
}

func (s *ItemCategoryService) Unarchive(ctx context.Context, id int64) (res dto.ItemCategoryResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemCategoryService.Archive")
	defer span.End()

	var ItemCategoryOld *model.ItemCategory
	ItemCategoryOld, err = s.RepositoryItemCategory.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("item_category_id")
		return
	}

	if ItemCategoryOld.Status != statusx.ConvertStatusName("Archived") {
		err = edenlabs.ErrorMustArchived("status")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	ItemCategory := &model.ItemCategory{
		ID:        ItemCategoryOld.ID,
		Status:    statusx.ConvertStatusName("Active"),
		UpdatedAt: time.Now(),
	}

	err = s.RepositoryItemCategory.Update(ctx, ItemCategory, "Status", "UpdatedAt")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	userID := ctx.Value(constants.KeyUserID).(int64)

	_, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &auditService.CreateLogRequest{
		Log: &auditService.Log{
			UserId:      userID,
			ReferenceId: strconv.Itoa(int(ItemCategory.ID)),
			Type:        "item_category",
			Function:    "unarchive",
			CreatedAt:   timestamppb.New(time.Now()),
		},
	})

	res = dto.ItemCategoryResponse{
		ID:            ItemCategoryOld.ID,
		Code:          ItemCategory.Code,
		Name:          ItemCategoryOld.Name,
		RegionID:      ItemCategoryOld.Regions,
		Status:        statusx.ConvertStatusName("Active"),
		StatusConvert: statusx.Active,
	}

	return
}

func (s *ItemCategoryService) Create(ctx context.Context, req dto.ItemCategoryRequestCreate) (res dto.ItemCategoryResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemCategoryService.UpdateImage")
	defer span.End()

	var (
		id           int64
		listRegionID string
		isExist      bool
	)

	// Validate for characters length
	if len(req.Name) < 1 || len(req.Name) > 20 {
		err = edenlabs.ErrorMustEqualOrLess("name", "20 characters")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	isExist = s.RepositoryItemCategory.IsExistNameItemCategory(ctx, 0, req.Name)
	// validate the name cannot be the same with existing
	if isExist {
		err = edenlabs.ErrorExists("name")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, v := range req.RegionID {
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

	var codeGenerator *configurationService.GetGenerateCodeResponse
	codeGenerator, err = s.opt.Client.ConfigurationServiceGrpc.GetGenerateCode(ctx, &configurationService.GetGenerateCodeRequest{
		Format: "ICT",
		Domain: "item_category",
		Length: 6,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("configuration", "generate_code")
		return
	}

	listRegionID = utils.ArrayStringToString(req.RegionID)

	itemCategory := &model.ItemCategory{
		Regions:   listRegionID,
		Code:      codeGenerator.Data.Code,
		Name:      req.Name,
		Status:    statusx.ConvertStatusName("Active"),
		CreatedAt: time.Now(),
	}

	id, err = s.RepositoryItemCategory.Create(ctx, itemCategory)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	itemCategoryImage := &model.ItemCategoryImage{
		ItemCategoryID: id,
		ImageUrl:       req.ImageUrl,
		CreatedAt:      time.Now(),
	}

	err = s.RepositoryItemCategoryImage.Create(ctx, itemCategoryImage)
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
			Type:        "item_category",
			Function:    "create",
			CreatedAt:   timestamppb.New(time.Now()),
		},
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.ItemCategoryResponse{
		ID:       id,
		Code:     itemCategory.Code,
		RegionID: listRegionID,
		Name:     req.Name,
		Status:   itemCategory.Status,
		ItemCategoryImage: &dto.ItemCategoryImageResponse{
			ID:             itemCategoryImage.ID,
			ItemCategoryID: itemCategoryImage.ItemCategoryID,
			ImageUrl:       itemCategoryImage.ImageUrl,
			CreatedAt:      itemCategoryImage.CreatedAt,
		},
	}

	return
}

func (s *ItemCategoryService) Update(ctx context.Context, req dto.ItemCategoryRequestUpdate, id int64) (res dto.ItemCategoryResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemCategoryService.Update")
	defer span.End()

	var (
		detailItemCategory *model.ItemCategory
		listRegionID       string
		isExist            bool
	)

	// Validate for characters length of name
	if len(req.Name) < 1 || len(req.Name) > 20 {
		err = edenlabs.ErrorMustEqualOrLess("name", "20 characters")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	isExist = s.RepositoryItemCategory.IsExistNameItemCategory(ctx, id, req.Name)
	// validate the name cannot be the same with existing
	if isExist {
		err = edenlabs.ErrorExists("name")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, v := range req.RegionID {
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

	listRegionID = utils.ArrayStringToString(req.RegionID)

	detailItemCategory, err = s.RepositoryItemCategory.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("item_category_id")
		return
	}

	detailItemCategory.Regions = listRegionID
	detailItemCategory.Name = req.Name
	detailItemCategory.UpdatedAt = time.Now()

	err = s.RepositoryItemCategory.Update(ctx, detailItemCategory, "Regions", "Name", "UpdatedAt")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	itemCategoryImage := dto.ItemCategoryImageRequestUpdate{
		ImageUrl: req.ImageUrl,
	}

	res, err = s.UpdateImage(ctx, itemCategoryImage, id)
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
			Type:        "item_category",
			Function:    "update",
			CreatedAt:   timestamppb.New(time.Now()),
		},
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.ItemCategoryResponse{
		ID:       id,
		Code:     detailItemCategory.Code,
		RegionID: listRegionID,
		Name:     req.Name,
		Status:   detailItemCategory.Status,
		ItemCategoryImage: &dto.ItemCategoryImageResponse{
			ID:             res.ItemCategoryImage.ID,
			ItemCategoryID: res.ItemCategoryImage.ItemCategoryID,
			ImageUrl:       res.ItemCategoryImage.ImageUrl,
			CreatedAt:      res.ItemCategoryImage.CreatedAt,
		},
	}

	return
}
