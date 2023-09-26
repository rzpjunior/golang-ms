package service

import (
	"context"
	"strconv"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/global"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	dto "git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
)

func NewServiceWrt() IWrtService {
	m := new(WrtService)
	m.opt = global.Setup.Common
	return m
}

type IWrtService interface {
	Get(ctx context.Context, req dto.GetWrtListRequest) (res []*dto.WrtResponse, err error)
	GetDetailById(ctx context.Context, req dto.GetWrtDetailRequest) (res *dto.WrtResponse, err error)
	GetGP(ctx context.Context, req dto.GetWrtListRequest) (res []*dto.WrtGP, total int64, err error)
	GetDetaiGPlById(ctx context.Context, id string) (res *dto.WrtGP, err error)
}

type WrtService struct {
	opt opt.Options
}

func (s *WrtService) Get(ctx context.Context, req dto.GetWrtListRequest) (res []*dto.WrtResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "WrtService.Get")
	defer span.End()

	// get wrt from bridge
	Wrt, err := s.opt.Client.ConfigurationServiceGrpc.GetWrtList(ctx, &configuration_service.GetWrtListRequest{
		RegionId: req.RegionId,
		Search:   req.Search,
		// Type:     int32(Type),
		Limit:  req.Limit,
		Offset: req.Offset,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	Region, err := s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPList(ctx, &bridge_service.GetAdmDivisionGPListRequest{
		Region: req.RegionId,
		// Type: "region",
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, wrtList := range Wrt.Data {
		res = append(res, &dto.WrtResponse{
			ID:     strconv.Itoa(int(wrtList.Id)),
			Code:   wrtList.Code,
			Name:   wrtList.Name,
			Note:   wrtList.Note,
			Type:   strconv.Itoa(int(wrtList.Type)),
			Status: "1",
			Region: &dto.RegionResponse{
				// ID:          Region.Data[0].Region,
				Code:        Region.Data[0].Region,
				Description: Region.Data[0].Region,
			},
		})
	}

	return
}

func (s *WrtService) GetDetailById(ctx context.Context, req dto.GetWrtDetailRequest) (res *dto.WrtResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "WrtService.GetDetailById")
	defer span.End()

	// get Wrt from bridge
	var wrt *bridgeService.GetWrtDetailResponse
	wrt, err = s.opt.Client.BridgeServiceGrpc.GetWrtDetail(ctx, &bridgeService.GetWrtDetailRequest{
		Id: int64(req.Id),
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "wrt")
		return
	}

	res = &dto.WrtResponse{
		// ID:        wrt.Data.Id,
		Code: wrt.Data.Code,
		// RegionID:  wrt.Data.RegionId,
		// StartTime: wrt.Data.StartTime,
		// EndTime:   wrt.Data.EndTime,
		Region: &dto.RegionResponse{
			// ID:            wrt.Data.Region.Id,
			Code:          wrt.Data.Region.Code,
			Description:   wrt.Data.Region.Description,
			Status:        int8(wrt.Data.Region.Status),
			StatusConvert: statusx.ConvertStatusValue(int8(wrt.Data.Region.Status)),
			CreatedAt:     wrt.Data.Region.CreatedAt.AsTime(),
			UpdatedAt:     wrt.Data.Region.UpdatedAt.AsTime(),
		},
	}

	return
}

func (s *WrtService) GetGP(ctx context.Context, req dto.GetWrtListRequest) (res []*dto.WrtGP, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "WrtService.GetGP")
	defer span.End()

	// get wrt from bridge
	var wrtRes *bridgeService.GetWrtGPResponse
	wrtRes, err = s.opt.Client.BridgeServiceGrpc.GetWrtGPList(ctx, &bridgeService.GetWrtGPListRequest{
		Limit:  req.Limit,
		Offset: req.Offset,
		Search: req.Search,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "wrt")
		return
	}

	datas := []*dto.WrtGP{}
	for _, wrt := range wrtRes.Data {
		datas = append(datas, &dto.WrtGP{
			GnL_Region: wrt.GnL_Region,
			GnL_WRT_ID: wrt.GnL_WRT_ID,
			Strttime:   wrt.Strttime,
			Endtime:    wrt.Endtime,
			Inactive:   wrt.Inactive,
		})
	}

	total = int64(len(datas))
	res = datas

	return
}

func (s *WrtService) GetDetaiGPlById(ctx context.Context, id string) (res *dto.WrtGP, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "WrtService.GetDetaiGPlById")
	defer span.End()

	// get Wrt from bridge
	var wrt *bridgeService.GetWrtGPResponse
	wrt, err = s.opt.Client.BridgeServiceGrpc.GetWrtGPDetail(ctx, &bridgeService.GetWrtGPDetailRequest{
		Id: id,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "wrt")
		return
	}

	res = &dto.WrtGP{
		GnL_Region: wrt.Data[0].GnL_Region,
		GnL_WRT_ID: wrt.Data[0].GnL_WRT_ID,
		Strttime:   wrt.Data[0].Strttime,
		Endtime:    wrt.Data[0].Endtime,
		Inactive:   wrt.Data[0].Inactive,
	}

	return
}
