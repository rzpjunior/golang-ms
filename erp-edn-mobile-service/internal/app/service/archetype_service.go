package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/global"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/internal/app/dto"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

func NewServiceArchetype() IArchetypeService {
	m := new(ArchetypeService)
	m.opt = global.Setup.Common
	return m
}

type IArchetypeService interface {
	GetGP(ctx context.Context, req dto.GetArchetypeGPListRequest) (res []*dto.ArchetypeGP, total int64, err error)
	GetDetaiGPlById(ctx context.Context, id string) (res *dto.ArchetypeGP, err error)
}

type ArchetypeService struct {
	opt opt.Options
}

func (s *ArchetypeService) GetGP(ctx context.Context, req dto.GetArchetypeGPListRequest) (res []*dto.ArchetypeGP, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ArchetypeService.GetGP")
	defer span.End()

	// get archetype from bridge
	var archeypeList *bridgeService.GetArchetypeGPResponse
	archeypeList, err = s.opt.Client.BridgeServiceGrpc.GetArchetypeGPList(ctx, &bridgeService.GetArchetypeGPListRequest{
		Limit:                   req.Limit,
		Offset:                  req.Offset,
		GnlArchetypeId:          req.GnlArchetypeId,
		GnlArchetypedescription: req.GnlArchetypedescription,
		GnlCustTypeId:           req.GnlCustTypeId,
		Inactive:                req.Inactive,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "archetype")
		return
	}

	datas := []*dto.ArchetypeGP{}
	for _, archetype := range archeypeList.Data {
		datas = append(datas, &dto.ArchetypeGP{
			GnlArchetypeId:          archetype.GnlArchetypeId,
			GnlArchetypedescription: archetype.GnlArchetypedescription,
			GnlCustTypeId:           archetype.GnlCustTypeId,
			GnlCusttypeDescription:  archetype.GnlCusttypeDescription,
			Inactive:                archetype.Inactive,
			InactiveDesc:            archetype.InactiveDesc,
		})
	}

	total = int64(len(datas))
	res = datas

	return
}

func (s *ArchetypeService) GetDetaiGPlById(ctx context.Context, id string) (res *dto.ArchetypeGP, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ArchetypeService.GetDetaiGPlById")
	defer span.End()

	// get archetype from bridge
	var archetype *bridgeService.GetArchetypeGPResponse
	archetype, err = s.opt.Client.BridgeServiceGrpc.GetArchetypeGPDetail(ctx, &bridgeService.GetArchetypeGPDetailRequest{
		Id: id,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "archetype")
		return
	}

	if len(archetype.Data) > 0 {
		res = &dto.ArchetypeGP{
			GnlArchetypeId:          archetype.Data[0].GnlArchetypeId,
			GnlArchetypedescription: archetype.Data[0].GnlArchetypedescription,
			GnlCustTypeId:           archetype.Data[0].GnlCustTypeId,
			GnlCusttypeDescription:  archetype.Data[0].GnlCusttypeDescription,
			Inactive:                archetype.Data[0].Inactive,
			InactiveDesc:            archetype.Data[0].InactiveDesc,
		}
	}

	return
}
