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

type IWRTService interface {
	Get(ctx context.Context, req dto.WrtRequest) (res []dto.WrtResponse, err error)
}

type WRTService struct {
	opt opt.Options
}

func NewWRTService() IWRTService {
	return &WRTService{
		opt: global.Setup.Common,
	}
}

func (s *WRTService) Get(ctx context.Context, req dto.WrtRequest) (res []dto.WrtResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "WRTService.Get")
	defer span.End()

	//cek data area id ada ga areanya,kalau ada ambil region policy.kalau uda ambil region policy nanti balikin yang ada di depan
	Region, err := s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPDetail(ctx, &bridge_service.GetAdmDivisionGPDetailRequest{
		Id:   "Greater Jakarta",
		Type: "region",
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// if Region.Data.Code == "" {
	// 	//throw error
	// }
	Type, _ := strconv.Atoi(req.Data.Type)
	Wrt, err := s.opt.Client.ConfigurationServiceGrpc.GetWrtList(ctx, &configuration_service.GetWrtListRequest{
		RegionId: "Greater Jakarta",
		Type:     int32(Type),
		Limit:    10,
		Offset:   1,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, wrtList := range Wrt.Data {
		res = append(res, dto.WrtResponse{
			ID:     strconv.Itoa(int(wrtList.Id)),
			Code:   wrtList.Code,
			Name:   wrtList.Name,
			Note:   wrtList.Note,
			Type:   strconv.Itoa(int(wrtList.Type)),
			Status: "1",
			Region: &dto.RegionResponse{
				ID:          Region.Data[0].Region,
				Code:        Region.Data[0].Region,
				Description: Region.Data[0].Region,
			},
		})
	}

	return
}
