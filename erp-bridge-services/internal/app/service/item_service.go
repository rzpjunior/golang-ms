package service

import (
	"context"
	"net/url"
	"strconv"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/edenlabs/edenlabs/timex"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/repository"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

type IItemService interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, uomID int64, classID int64, itemCategoryID int64) (res []dto.ItemResponse, total int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (res dto.ItemResponse, err error)
	GetGP(ctx context.Context, req *pb.GetItemGPListRequest) (res *pb.GetItemGPResponse, err error)
	GetItemMasterComplexGP(ctx context.Context, req *pb.GetItemMasterComplexGPListRequest) (res *pb.GetItemMasterComplexGPListResponse, err error)
	GetDetailGP(ctx context.Context, req *pb.GetItemGPDetailRequest) (res *pb.GetItemGPResponse, err error)
}

type ItemService struct {
	opt            opt.Options
	RepositoryItem repository.IItemRepository
}

func NewItemService() IItemService {
	return &ItemService{
		opt:            global.Setup.Common,
		RepositoryItem: repository.NewItemRepository(),
	}
}

func (s *ItemService) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, uomID int64, classID int64, itemCategoryID int64) (res []dto.ItemResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemService.Get")
	defer span.End()

	var item []*model.Item
	item, total, err = s.RepositoryItem.Get(ctx, offset, limit, status, search, orderBy, uomID, classID, itemCategoryID)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, item := range item {
		res = append(res, dto.ItemResponse{
			ID:                      item.ID,
			Code:                    item.Code,
			UomID:                   item.UomID,
			ClassID:                 item.ClassID,
			ItemCategoryID:          item.ItemCategoryID,
			Description:             item.Description,
			UnitWeightConversion:    item.UnitWeightConversion,
			OrderMinQty:             item.OrderMinQty,
			OrderMaxQty:             item.OrderMaxQty,
			ItemType:                item.ItemType,
			Packability:             item.Packability,
			Capitalize:              item.Capitalize,
			Note:                    item.Note,
			ExcludeArchetype:        item.ExcludeArchetype,
			MaxDayDeliveryDate:      item.MaxDayDeliveryDate,
			FragileGoods:            item.FragileGoods,
			Taxable:                 item.Taxable,
			OrderChannelRestriction: item.OrderChannelRestriction,
			Status:                  item.Status,
			StatusConvert:           statusx.ConvertStatusValue(item.Status),
			CreatedAt:               timex.ToLocTime(ctx, item.CreatedAt),
			UpdatedAt:               timex.ToLocTime(ctx, item.UpdatedAt),
		})
	}

	return
}

func (s *ItemService) GetDetail(ctx context.Context, id int64, code string) (res dto.ItemResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemService.GetDetail")
	defer span.End()

	var item *model.Item
	item, err = s.RepositoryItem.GetDetail(ctx, id, code)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.ItemResponse{
		ID:                      item.ID,
		Code:                    item.Code,
		UomID:                   item.UomID,
		ClassID:                 item.ClassID,
		ItemCategoryID:          item.ItemCategoryID,
		Description:             item.Description,
		UnitWeightConversion:    item.UnitWeightConversion,
		OrderMinQty:             item.OrderMinQty,
		OrderMaxQty:             item.OrderMaxQty,
		ItemType:                item.ItemType,
		Packability:             item.Packability,
		Capitalize:              item.Capitalize,
		ExcludeArchetype:        item.ExcludeArchetype,
		MaxDayDeliveryDate:      item.MaxDayDeliveryDate,
		FragileGoods:            item.FragileGoods,
		Taxable:                 item.Taxable,
		OrderChannelRestriction: item.OrderChannelRestriction,
		Note:                    item.Note,
		Status:                  item.Status,
		StatusConvert:           statusx.ConvertStatusValue(item.Status),
		CreatedAt:               timex.ToLocTime(ctx, item.CreatedAt),
		UpdatedAt:               timex.ToLocTime(ctx, item.UpdatedAt),
	}

	return
}

func (s *ItemService) UpdateItemPackable(ctx context.Context, id int64, code string) (res dto.ItemResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemService.UpdateItemPackable")
	defer span.End()

	res = dto.ItemResponse{}
	return
}

func (s *ItemService) UpdateItemFragile(ctx context.Context, id int64, code string) (res dto.ItemResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemService.UpdateItemFragile")
	defer span.End()

	res = dto.ItemResponse{}
	return
}

func (s *ItemService) GetGP(ctx context.Context, req *pb.GetItemGPListRequest) (res *pb.GetItemGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemService.GetGP")
	defer span.End()

	if req.ItemNumber != "" {

	}

	params := map[string]string{
		"interid":    global.EnvDatabaseGP,
		"PageNumber": strconv.Itoa(int(req.Offset)),
		"PageSize":   strconv.Itoa(int(req.Limit)),
	}

	if req.ItemNumber != "" {
		params["itemnmbr"] = req.ItemNumber
	}

	if req.Description != "" {
		req.Description = url.PathEscape(req.Description)
		params["itemdesc"] = req.Description
	}

	if req.ClassId != "" {
		req.ClassId = url.PathEscape(req.ClassId)
		params["itmclscd"] = req.ClassId
	}

	if req.UomId != "" {
		req.UomId = url.PathEscape(req.UomId)
		params["uomschdl"] = req.UomId
	}

	if req.Inactive != "" {
		params["inactive"] = req.Inactive
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "item/getall", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *ItemService) GetDetailGP(ctx context.Context, req *pb.GetItemGPDetailRequest) (res *pb.GetItemGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemService.GetDetailGP")
	defer span.End()

	req.Id = url.PathEscape(req.Id)
	params := map[string]string{
		"interid": global.EnvDatabaseGP,
		"id":      req.Id,
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "item/getbyid", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if len(res.Data) == 0 || len(res.Data) > 1 {
		err = edenlabs.ErrorNotFound("item")
	}

	return
}

func (s *ItemService) GetItemMasterComplexGP(ctx context.Context, req *pb.GetItemMasterComplexGPListRequest) (res *pb.GetItemMasterComplexGPListResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemService.GetItemMasterComplexGP")
	defer span.End()

	params := map[string]string{
		"interid":    global.EnvDatabaseGP,
		"PageNumber": strconv.Itoa(int(req.Offset)),
		"PageSize":   strconv.Itoa(int(req.Limit)),
	}

	if req.ItemNumber != "" {
		params["itemnmbr"] = req.ItemNumber
	}
	if req.Description != "" {
		req.Description = url.PathEscape(req.Description)
		params["itemdesc"] = req.Description
	}
	if req.ClassId != "" {
		req.ClassId = url.PathEscape(req.ClassId)
		params["itmclscd"] = req.ClassId
	}
	if req.UomId != "" {
		req.UomId = url.PathEscape(req.UomId)
		params["uomschdl"] = req.UomId
	}
	if req.Inactive != "" {
		params["inactive"] = req.Inactive
	}
	if req.Locncode != "" {
		params["locncode"] = req.Locncode
	}
	if req.GnlRegion != "" {
		req.GnlRegion = url.PathEscape(req.GnlRegion)
		params["gnl_region"] = req.GnlRegion
	}
	if req.GnlSalability != "" {
		params["gnl_salability"] = req.GnlSalability
	}
	if req.GnlStorability != "" {
		params["gnl_storability"] = req.GnlStorability
	}
	if req.GnlCustTypeId != "" {
		params["gnl_cust_type_id"] = req.GnlCustTypeId
	}
	if req.Prclevel != "" {
		params["prclevel"] = req.Prclevel
	}
	if req.OrderBy != "" {
		params["order_by"] = req.OrderBy
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "item/complex", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}
