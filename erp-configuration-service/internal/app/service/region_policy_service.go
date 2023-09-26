package service

import (
	"context"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-configuration-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/internal/app/repository"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

type IRegionPolicyService interface {
	Get(ctx context.Context, offset int, limit int, search string, orderBy string, regionID string) (res []dto.RegionPolicyResponse, total int64, err error)
	GetDetail(ctx context.Context, id int64, code string, regionId string) (res dto.RegionPolicyResponse, err error)
	Update(ctx context.Context, req dto.RegionPolicyRequestUpdate, id int64) (res dto.RegionPolicyResponse, err error)
}

type RegionPolicyService struct {
	opt                    opt.Options
	RepositoryRegionPolicy repository.IRegionPolicyRepository
}

func NewRegionPolicyService() IRegionPolicyService {
	return &RegionPolicyService{
		opt:                    global.Setup.Common,
		RepositoryRegionPolicy: repository.NewRegionPolicyRepository(),
	}
}

func (s *RegionPolicyService) Get(ctx context.Context, offset int, limit int, search string, orderBy string, regionID string) (res []dto.RegionPolicyResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "RegionPolicyService.Get")
	defer span.End()

	var RegionPolicys []*model.RegionPolicy
	RegionPolicys, total, err = s.RepositoryRegionPolicy.Get(ctx, offset, limit, search, orderBy, regionID)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, regionPolicy := range RegionPolicys {
		var region *bridgeService.GetAdmDivisionGPResponse
		region, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPDetail(ctx, &bridgeService.GetAdmDivisionGPDetailRequest{
			Region: regionPolicy.RegionIDGP,
			Limit:  1,
			Offset: 0,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "region")
			return
		}
		regionResponse := &dto.RegionResponse{
			ID:          region.Data[0].Region,
			Code:        region.Data[0].Region,
			Description: region.Data[0].Region,
		}

		var priceLevel *bridgeService.GetSalesPriceLevelResponse
		priceLevel, err = s.opt.Client.BridgeServiceGrpc.GetSalesPriceLevelDetail(ctx, &bridgeService.GetSalesPriceLevelDetailRequest{
			Id: regionPolicy.DefaultPriceLevel,
		})
		if err != nil || len(priceLevel.Data) == 0 {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "price level")
			return
		}

		priceLevelResponse := &dto.PriceLevelResponse{
			ID:             priceLevel.Data[0].Prclevel,
			Description:    priceLevel.Data[0].Prclevel,
			RegionID:       priceLevel.Data[0].GnlRegion,
			CustomerTypeID: priceLevel.Data[0].GnlCustTypeId,
		}

		res = append(res, dto.RegionPolicyResponse{
			ID:                 regionPolicy.ID,
			OrderTimeLimit:     regionPolicy.OrderTimeLimit,
			MaxDayDeliveryDate: regionPolicy.MaxDayDeliveryDate,
			WeeklyDayOff:       regionPolicy.WeeklyDayOff,
			Region:             regionResponse,
			CSPhoneNumber:      regionPolicy.CSPhoneNumber,
			DefaultPriceLevel:  priceLevelResponse,
		})
	}

	return
}

func (s *RegionPolicyService) GetDetail(ctx context.Context, id int64, code string, regionId string) (res dto.RegionPolicyResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "RegionPolicyService.GetDetail")
	defer span.End()

	var regionPolicy *model.RegionPolicy
	regionPolicy, err = s.RepositoryRegionPolicy.GetDetail(ctx, id, code, regionId)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	var region *bridgeService.GetAdmDivisionGPResponse
	region, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPDetail(ctx, &bridgeService.GetAdmDivisionGPDetailRequest{
		Region: regionPolicy.RegionIDGP,
		Limit:  1,
		Offset: 0,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "region")
		return
	}

	regionResponse := &dto.RegionResponse{
		ID:          region.Data[0].Region,
		Code:        region.Data[0].Region,
		Description: region.Data[0].Region,
	}

	var priceLevel *bridgeService.GetSalesPriceLevelResponse
	priceLevel, err = s.opt.Client.BridgeServiceGrpc.GetSalesPriceLevelDetail(ctx, &bridgeService.GetSalesPriceLevelDetailRequest{
		Id: regionPolicy.DefaultPriceLevel,
	})
	if err != nil || len(priceLevel.Data) == 0 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "price level")
		return
	}

	priceLevelResponse := &dto.PriceLevelResponse{
		ID:             priceLevel.Data[0].Prclevel,
		Description:    priceLevel.Data[0].Prclevel,
		RegionID:       priceLevel.Data[0].GnlRegion,
		CustomerTypeID: priceLevel.Data[0].GnlCustTypeId,
	}

	res = dto.RegionPolicyResponse{
		ID:                 regionPolicy.ID,
		OrderTimeLimit:     regionPolicy.OrderTimeLimit,
		MaxDayDeliveryDate: regionPolicy.MaxDayDeliveryDate,
		WeeklyDayOff:       regionPolicy.WeeklyDayOff,
		Region:             regionResponse,
		CSPhoneNumber:      regionPolicy.CSPhoneNumber,
		DefaultPriceLevel:  priceLevelResponse,
	}

	return
}
func (s *RegionPolicyService) Update(ctx context.Context, req dto.RegionPolicyRequestUpdate, id int64) (res dto.RegionPolicyResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "RegionPolicyService.Update")
	defer span.End()

	var priceLevel *bridgeService.GetSalesPriceLevelResponse
	priceLevel, err = s.opt.Client.BridgeServiceGrpc.GetSalesPriceLevelDetail(ctx, &bridgeService.GetSalesPriceLevelDetailRequest{
		Id: req.DefaultPriceLevel,
	})
	if err != nil || len(priceLevel.Data) == 0 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("default_price_level")
		return
	}

	// validate data is exist
	_, err = s.RepositoryRegionPolicy.GetDetail(ctx, id, "", "")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// validate format time
	if _, err = time.Parse("15:04", req.OrderTimeLimit); err != nil {
		err = edenlabs.ErrorInvalid("order_time_limit")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// Validate weekly day off
	if req.WeeklyDayOff < 1 || req.WeeklyDayOff > 7 {
		err = edenlabs.ErrorMustEqualOrLess("weekly_day_off", "7")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	RegionPolicy := &model.RegionPolicy{
		ID:                 id,
		OrderTimeLimit:     req.OrderTimeLimit,
		MaxDayDeliveryDate: req.MaxDayDeliveryDate,
		WeeklyDayOff:       req.WeeklyDayOff,
		DefaultPriceLevel:  req.DefaultPriceLevel,
	}

	err = s.RepositoryRegionPolicy.Update(ctx, RegionPolicy, "OrderTimeLimit", "MaxDayDeliveryDate", "WeeklyDayOff", "default_price_level")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.RegionPolicyResponse{
		ID:                 id,
		OrderTimeLimit:     req.OrderTimeLimit,
		MaxDayDeliveryDate: req.MaxDayDeliveryDate,
		WeeklyDayOff:       req.WeeklyDayOff,
		DefaultPriceLevel:  res.DefaultPriceLevel,
	}

	return
}
