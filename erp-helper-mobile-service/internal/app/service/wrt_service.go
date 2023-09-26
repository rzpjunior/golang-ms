package service

import (
	"context"
	"fmt"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-helper-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-helper-mobile-service/internal/app/dto"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

type IWrtService interface {
	Get(ctx context.Context, req dto.GetWrtRequest) (res []dto.WrtResponse, total int64, err error)
	GetDetail(ctx context.Context, id string) (res dto.WrtResponse, err error)
}

type WrtService struct {
	opt opt.Options
}

func NewServiceWrt() IWrtService {
	return &WrtService{
		opt: global.Setup.Common,
	}
}

func (s *WrtService) Get(ctx context.Context, req dto.GetWrtRequest) (res []dto.WrtResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "WrtService.Get")
	defer span.End()

	// get site's wrt
	// get site detail
	var site *bridgeService.GetSiteGPResponse
	if site, err = s.opt.Client.BridgeServiceGrpc.GetSiteGPDetail(ctx, &bridgeService.GetSiteGPDetailRequest{
		Id: req.SiteId,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "site")
		return
	}

	// get region of the adm division
	var admDivision *bridgeService.GetAdmDivisionGPResponse
	if admDivision, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPList(ctx, &bridgeService.GetAdmDivisionGPListRequest{
		Limit:           1000,
		AdmDivisionCode: site.Data[0].GnlAdministrativeCode,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "adm division")
		return
	}

	fmt.Println("Region", admDivision.Data[0].Region)
	// get wrt of the region
	var wrt *bridgeService.GetWrtGPResponse
	if wrt, err = s.opt.Client.BridgeServiceGrpc.GetWrtGPList(ctx, &bridgeService.GetWrtGPListRequest{
		Limit:     int32(req.Limit),
		Offset:    int32(req.Offset),
		GnlRegion: admDivision.Data[0].Region,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "wrt")
		return
	}

	for _, wrt := range wrt.Data {
		res = append(res, dto.WrtResponse{
			ID:        wrt.GnL_WRT_ID,
			RegionId:  wrt.GnL_Region,
			StartTime: wrt.Strttime,
			EndTime:   wrt.Endtime,
			Status:    int8(*wrt.Inactive),
		})
	}

	total = int64(len(wrt.Data))

	return
}

func (s *WrtService) GetDetail(ctx context.Context, id string) (res dto.WrtResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "WrtService.GetWrt")
	defer span.End()

	var wrt *bridgeService.GetWrtGPResponse

	if wrt, err = s.opt.Client.BridgeServiceGrpc.GetWrtGPDetail(ctx, &bridgeService.GetWrtGPDetailRequest{
		Id: id,
	}); err != nil || !wrt.Succeeded {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "wrt")
		return
	}

	res = dto.WrtResponse{
		ID:        wrt.Data[0].GnL_WRT_ID,
		RegionId:  wrt.Data[0].GnL_Region,
		StartTime: wrt.Data[0].Strttime,
		EndTime:   wrt.Data[0].Endtime,
		Status:    int8(*wrt.Data[0].Inactive),
	}

	return
}
