package service

import (
	"context"
	"strconv"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
)

type IRegionPolicyService interface {
	Get(ctx context.Context, offset int, limit int, search string, adm_division_id string, Type int) (res dto.RegionPolicy, total int64, err error)
}

type RegionPolicyService struct {
	opt opt.Options
}

func NewRegionPolicyService() IRegionPolicyService {
	return &RegionPolicyService{
		opt: global.Setup.Common,
	}
}

func (s *RegionPolicyService) Get(ctx context.Context, offset int, limit int, search string, adm_division_id string, Type int) (res dto.RegionPolicy, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "RegionPolicyService.Get")
	defer span.End()

	// admDivId := int64(adm_division_id)
	// admDivision, err := s.opt.Client.BridgeServiceGrpc.GetAdmDivisionDetail(ctx, &bridge_service.GetAdmDivisionDetailRequest{
	// 	Id: admDivId,
	// })
	admDivision, err := s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPList(ctx, &bridge_service.GetAdmDivisionGPListRequest{
		AdmDivisionCode: adm_division_id,
		Limit:           1,
		Offset:          0,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	if len(admDivision.Data) == 0 {
		//throw error
	}
	// Region, err := s.opt.Client.BridgeServiceGrpc.GetRegionDetail(ctx, &bridge_service.GetRegionDetailRequest{
	// 	Id: int64(admDivision.Data.RegionId),
	// })
	// if err != nil {
	// 	span.RecordError(err)
	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 	return
	// }
	// if Region.Data.Code == "" {
	// 	//throw error
	// }

	RegionPolicy, err := s.opt.Client.ConfigurationServiceGrpc.GetRegionPolicyList(ctx, &configuration_service.GetRegionPolicyListRequest{
		Limit:  1,
		Offset: 0,
		Search: admDivision.Data[0].Region,
		// RegionId: ,
	})

	res1 := dto.RegionPolicyResponse{
		ID:                 strconv.Itoa(int(RegionPolicy.Data[0].Id)),
		OrderTimeLimit:     RegionPolicy.Data[0].OrderTimeLimit,
		MaxDayDeliveryDate: strconv.Itoa(int(RegionPolicy.Data[0].MaxDayDeliveryDate)),
		WeeklyDayOff:       strconv.Itoa(int(RegionPolicy.Data[0].WeeklyDayOff)),
		CSPhoneNumber:      RegionPolicy.Data[0].CsPhoneNumber,
		Region: &dto.RegionResponse{
			ID:          admDivision.Data[0].Region,
			Code:        admDivision.Data[0].Region,
			Description: admDivision.Data[0].Region,
		},
	}
	res.RegPol = res1
	return
}
