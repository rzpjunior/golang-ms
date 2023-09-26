package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/internal/app/dto"
)

type IAdmDivisionService interface {
	Get(ctx context.Context, offset int, limit int, search string, sub_district_id string, Type int) (res []dto.AdmDivisionResponse, total int64, err error)
	GetByID(ctx context.Context, id string) (res *dto.AdmDivisionResponse, err error)
}

type AdmDivisionService struct {
	opt opt.Options
}

func NewAdmDivisionService() IAdmDivisionService {
	return &AdmDivisionService{
		opt: global.Setup.Common,
	}
}

func (s *AdmDivisionService) Get(ctx context.Context, offset int, limit int, search string, sub_district_id string, Type int) (res []dto.AdmDivisionResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AdmDivisionService.Get")
	defer span.End()

	admDivisions, err := s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPDetail(ctx, &bridge_service.GetAdmDivisionGPDetailRequest{
		Type: "state",
	})

	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, admDivision := range admDivisions.Data {
		res = append(res, dto.AdmDivisionResponse{
			ID:          admDivision.Code,
			Code:        admDivision.Code,
			Region:      admDivision.Region,
			Province:    admDivision.State,
			City:        admDivision.City,
			District:    admDivision.District,
			SubDistrict: admDivision.Subdistrict,
		})

	}
	total = int64(len(res))

	return
}

func (s *AdmDivisionService) GetByID(ctx context.Context, id string) (res *dto.AdmDivisionResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AdmDivisionService.GetByID")
	defer span.End()

	var admDivision *bridgeService.GetAdmDivisionGPResponse
	admDivision, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPDetail(ctx, &bridge_service.GetAdmDivisionGPDetailRequest{
		Id:   id,
		Type: "state",
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "adm_division")
		return
	}

	res = &dto.AdmDivisionResponse{
		ID:          admDivision.Data[0].Code,
		Code:        admDivision.Data[0].Code,
		Region:      admDivision.Data[0].Region,
		Province:    admDivision.Data[0].State,
		City:        admDivision.Data[0].City,
		District:    admDivision.Data[0].District,
		SubDistrict: admDivision.Data[0].Subdistrict,
	}

	return
}
