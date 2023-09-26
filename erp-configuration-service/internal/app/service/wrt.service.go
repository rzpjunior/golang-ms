package service

import (
	"context"
	"strconv"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/constants"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-configuration-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/internal/app/repository"
	auditService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/audit_service"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type IWrtService interface {
	Get(ctx context.Context, offset, limit, wrtType int, regionID string, search string, status int8) (res []dto.WrtResponse, total int64, err error)
	GetDetail(ctx context.Context, id int64) (res dto.WrtResponse, err error)
	Update(ctx context.Context, req dto.WrtRequestUpdate, id int64) (res dto.WrtResponse, err error)
	GetIDDetail(ctx context.Context, id int64) (res dto.WrtResponse, err error)
}

type WrtService struct {
	opt           opt.Options
	RepositoryWrt repository.IWrtRepository
}

func NewWrtService() IWrtService {
	return &WrtService{
		opt:           global.Setup.Common,
		RepositoryWrt: repository.NewWrtRepository(),
	}
}

func (s *WrtService) Get(ctx context.Context, offset, limit, wrtType int, regionID string, search string, status int8) (res []dto.WrtResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "WrtService.Get")
	defer span.End()

	var (
		regionResponse *dto.RegionResponse
		wrtListGP      *bridgeService.GetWrtGPResponse
		wrtDetail      *model.Wrt
		wrtResponse    dto.WrtResponse
		region         *bridgeService.GetAdmDivisionGPResponse
	)

	wrtListGP, err = s.opt.Client.BridgeServiceGrpc.GetWrtGPList(ctx, &bridgeService.GetWrtGPListRequest{
		Offset:    int32(offset),
		Limit:     int32(limit),
		Search:    search,
		GnlRegion: regionID,
	})
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	layout := "2006-01-02"

	for _, wrtGP := range wrtListGP.Data {

		createAt := time.Now().Format(layout)
		// sync GP
		WrtGP := &model.Wrt{
			WrtID: wrtGP.GnL_WRT_ID,
			Type:  1,
			Note:  "From GP " + createAt,
		}
		if err = s.RepositoryWrt.SyncGP(ctx, WrtGP); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("item sync")
			return
		}

		if *wrtGP.Inactive == 1 {
			continue
		}
		wrtDetail, err = s.RepositoryWrt.GetDetail(ctx, 0, wrtGP.GnL_WRT_ID)

		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		switch wrtGP.GnL_Region {
		case "GJ":
			wrtGP.GnL_Region = "Greater Jakarta"
		case "EJ":
			wrtGP.GnL_Region = "East Java"
		case "CJ":
			wrtGP.GnL_Region = "Central Java"
		case "NS":
			wrtGP.GnL_Region = "North Sumatera"
		case "WJ":
			wrtGP.GnL_Region = "West Java"
		case "HO":
			wrtGP.GnL_Region = "Greater Jakarta"
		case "Head Office":
			wrtGP.GnL_Region = "Greater Jakarta"
		case "North Sumatra":
			wrtGP.GnL_Region = "North Sumatera"
		}

		region, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPDetail(ctx, &bridgeService.GetAdmDivisionGPDetailRequest{
			Region: wrtGP.GnL_Region,
			Limit:  1,
			Offset: 0,
		})

		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("region_id")
			return
		}
		regionResponse = &dto.RegionResponse{
			ID:          region.Data[0].Region,
			Code:        region.Data[0].Region,
			Description: region.Data[0].Region,
		}

		wrtResponse = dto.WrtResponse{
			ID:       wrtDetail.ID,
			Name:     wrtGP.Strttime[0:5] + "-" + wrtGP.Endtime[0:5],
			Type:     wrtDetail.Type,
			Note:     wrtDetail.Note,
			Region:   regionResponse,
			RegionID: regionResponse.ID,
			Code:     wrtGP.GnL_WRT_ID,
		}

		// filter wrt type
		if wrtType != 0 && wrtDetail.Type == int8(wrtType) {
			res = append(res, wrtResponse)
		} else if wrtType == 0 {
			res = append(res, wrtResponse)
		}
	}
	total = int64(len(res))

	return
}

func (s *WrtService) GetDetail(ctx context.Context, id int64) (res dto.WrtResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "WrtService.GetDetail")
	defer span.End()

	var (
		wrt            *model.Wrt
		wrtGP          *bridgeService.GetWrtGPResponse
		regionResponse *dto.RegionResponse
	)

	wrt, err = s.RepositoryWrt.GetDetail(ctx, id, "")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("id")
		return
	}

	wrtGP, err = s.opt.Client.BridgeServiceGrpc.GetWrtGPDetail(ctx, &bridgeService.GetWrtGPDetailRequest{
		Id: wrt.WrtID,
	})

	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("wrt_id")
		return
	}

	switch wrtGP.Data[0].GnL_Region {
	case "GJ":
		wrtGP.Data[0].GnL_Region = "Greater Jakarta"
	case "EJ":
		wrtGP.Data[0].GnL_Region = "East Java"
	case "CJ":
		wrtGP.Data[0].GnL_Region = "Central Java"
	case "NS":
		wrtGP.Data[0].GnL_Region = "North Sumatra"
	case "WJ":
		wrtGP.Data[0].GnL_Region = "West Java"
	case "HO":
		wrtGP.Data[0].GnL_Region = "Greater Jakarta"
	case "Head Office":
		wrtGP.Data[0].GnL_Region = "Greater Jakarta"
	case "North Sumatra":
		wrtGP.Data[0].GnL_Region = "North Sumatera"
	}

	region, err := s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPDetail(ctx, &bridgeService.GetAdmDivisionGPDetailRequest{
		Region: wrtGP.Data[0].GnL_Region,
		Limit:  1,
		Offset: 0,
	})

	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("region_id")
		return
	}
	regionResponse = &dto.RegionResponse{
		ID:          region.Data[0].Region,
		Code:        region.Data[0].Region,
		Description: region.Data[0].Region,
	}

	res = dto.WrtResponse{
		ID:       wrt.ID,
		Code:     wrt.WrtID,
		Name:     wrtGP.Data[0].Strttime[0:5] + "-" + wrtGP.Data[0].Endtime[0:5],
		Type:     wrt.Type,
		Note:     wrt.Note,
		RegionID: region.Data[0].Region,
		Region:   regionResponse,
	}

	return
}

func (s *WrtService) GetIDDetail(ctx context.Context, id int64) (res dto.WrtResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "WrtService.GetDetail")
	defer span.End()

	var (
		wrt *model.Wrt
		// wrtGP *bridgeService.GetWrtGPResponse
		// regionResponse *dto.RegionResponse
	)

	wrt, err = s.RepositoryWrt.GetDetail(ctx, id, "")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("id")
		return
	}

	res = dto.WrtResponse{
		ID:   wrt.ID,
		Code: wrt.WrtID,
		// Name: wrtGP.Data[0].Strttime[0:5] + "-" + wrtGP.Data[0].Endtime[0:5],
		Type: wrt.Type,
		Note: wrt.Note,
		// RegionID: region.Data[0].Region,
		// Region:   regionResponse,
	}

	return
}
func (s *WrtService) Update(ctx context.Context, req dto.WrtRequestUpdate, id int64) (res dto.WrtResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "WrtService.Update")
	defer span.End()

	// Validate the value delivery(1) or self pickup(2)
	if req.Type != 1 && req.Type != 2 {
		err = edenlabs.ErrorInvalid("type")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	var WrtOld *model.Wrt
	WrtOld, err = s.RepositoryWrt.GetDetail(ctx, id, "")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	wrtGP, err := s.opt.Client.BridgeServiceGrpc.GetWrtGPDetail(ctx, &bridgeService.GetWrtGPDetailRequest{
		Id: WrtOld.WrtID,
	})

	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("wrt_id")
		return
	}

	Wrt := &model.Wrt{
		ID:   WrtOld.ID,
		Type: req.Type,
		Note: req.Note,
	}

	err = s.RepositoryWrt.Update(ctx, Wrt, "Type", "Note")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	userID := ctx.Value(constants.KeyUserID).(int64)

	_, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &auditService.CreateLogRequest{
		Log: &auditService.Log{
			UserId:      userID,
			ReferenceId: strconv.Itoa(int(Wrt.ID)),
			Type:        "wrt",
			Function:    "Update",
			CreatedAt:   timestamppb.New(time.Now()),
		},
	})

	res = dto.WrtResponse{
		ID:   WrtOld.ID,
		Name: wrtGP.Data[0].Strttime[0:5] + "-" + wrtGP.Data[0].Endtime[0:5],
		Type: Wrt.Type,
		Note: WrtOld.Note,
	}

	return
}
