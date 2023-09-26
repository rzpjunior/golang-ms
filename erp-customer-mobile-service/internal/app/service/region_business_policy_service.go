package service

import (
	"context"
	"strconv"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

type IRegionBusinessPolicyService interface {
	Get(ctx context.Context, req dto.RegionBusinessPolRequest) (res []dto.RegionBusinessPolResponse, total int64, err error)
}

type RegionBusinessPolicyService struct {
	opt opt.Options
}

func NewRegionBusinessPolicyService() IRegionBusinessPolicyService {
	return &RegionBusinessPolicyService{
		opt: global.Setup.Common,
	}
}

func (s *RegionBusinessPolicyService) Get(ctx context.Context, req dto.RegionBusinessPolRequest) (res []dto.RegionBusinessPolResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "RegionBusinessPolicyService.Get")
	defer span.End()

	regionID, _ := strconv.Atoi(req.Data.RegionID)
	Region, err := s.opt.Client.BridgeServiceGrpc.GetRegionDetail(ctx, &bridge_service.GetRegionDetailRequest{
		Id: int64(regionID),
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	if Region.Data.Code == "" {
		//throw error
	}

	// customerTypeID, _ := strconv.Atoi(req.Data.CustomerTypeID)
	// customerType, err := s.opt.Client.BridgeServiceGrpc.GetRegionDetail(ctx, &bridge_service.GetRegionDetailRequest{
	// 	Id: int64(regionID),
	// })
	// if err != nil {
	// 	span.RecordError(err)
	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 	return
	// }

	if Region.Data.Code == "" {
		//throw error
	}
	// Region, err := s.opt.Client.BridgeServiceGrpc.GetRegionDetail(ctx, &bridge_service.GetRegionDetailRequest{
	// 	Id: int64(regionID),
	// })
	// if err != nil {
	// 	span.RecordError(err)
	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 	return
	// }
	// if Region.Data.Code == "" {
	// 	//throw error
	// }
	//fmt.Println(customerTypeID)
	// RegionBusinessPolicy, err := s.opt.Client.ConfigurationServiceGrpc.GetRegionBusinessPolicyList(ctx, &configuration_service.GetRegionBusinessPolicyListRequest{
	// 	Limit:    int32(limit),
	// 	Offset:   int32(offset),
	// 	Search:   search,
	// 	RegionId: int64(region_id),
	// })

	// for _, RegionBusinessPolicy := range RegionBusinessPolicy.Data {
	// 	res = append(res, dto.RegionBusinessPolicyResponse{
	// 		ID:                 RegionBusinessPolicy.Id,
	// 		OrderTimeLimit:     RegionBusinessPolicy.OrderTimeLimit,
	// 		MaxDayDeliveryDate: int(RegionBusinessPolicy.MaxDayDeliveryDate),
	// 		WeeklyDayOff:       int(RegionBusinessPolicy.WeeklyDayOff),
	// 		Region: &dto.RegionResponse{
	// 			ID:          Region.Data.Id,
	// 			Code:        Region.Data.Code,
	// 			Description: Region.Data.Description,
	// 		},
	// 	})
	// }
	total = int64(len(res))

	return
}
