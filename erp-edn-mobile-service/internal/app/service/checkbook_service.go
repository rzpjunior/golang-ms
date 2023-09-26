package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
)

func NewServiceCheckbook() ICheckbookService {
	m := new(CheckbookService)
	m.opt = global.Setup.Common
	return m
}

type ICheckbookService interface {
	GetGP(ctx context.Context, req dto.CheckbookListRequest) (res []*dto.CheckbookGP, total int64, err error)
	// GetDetaiGPlById(ctx context.Context, id string) (res *dto.CheckbookGP, err error)
}

type CheckbookService struct {
	opt opt.Options
}

func (s *CheckbookService) GetGP(ctx context.Context, req dto.CheckbookListRequest) (res []*dto.CheckbookGP, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CheckbookService.GetGP")
	defer span.End()

	// get payment method from bridge
	// var pmRes *bridgeService.GetCheckbookGPResponse
	// pmRes, err = s.opt.Client.BridgeServiceGrpc.GetCheckbookGPList(ctx, &bridgeService.GetCheckbookGPListRequest{
	// 	Limit:  req.Limit,
	// 	Offset: req.Offset,
	// })
	// if err != nil {
	// 	span.RecordError(err)
	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 	err = edenlabs.ErrorRpcNotFound("bridge", "payment method")
	// 	return
	// }

	// datas := []*dto.CheckbookGP{}
	// for _, pm := range pmRes.Data {
	// 	datas = append(datas, &dto.CheckbookGP{
	// 		Pmntnmbr:             pm.Pmntnmbr,
	// 		Vendorid:             pm.Vendorid,
	// 		Docnumbr:             pm.Docnumbr,
	// 		Docdate:              pm.Docdate,
	// 		Chekbkid:             pm.Chekbkid,
	// 		Pyenttyp:             pm.Pyenttyp,
	// 		CheckbookDesc:    pm.CheckbookDesc,
	// 		PrpCheckbookId:   pm.PrpCheckbookId,
	// 		PrpCheckbookDesc: pm.PrpCheckbookDesc,
	// 		Inactive:             pm.Inactive,
	// 		InactiveDesc:         pm.InactiveDesc,
	// 	})
	// }
	// paymentChannel, err := s.opt.Client.SalesServiceGrpc.GetPaymentChannelList(ctx, &sales_service.GetPaymentChannelListRequest{
	// 	Status:     1,
	// 	PublishIva: 1,
	// 	// PublishFva:      1,
	// 	CheckbookId: 2,
	// })

	Checkbook, _ := s.opt.Client.ConfigurationServiceGrpc.GetConfigAppList(ctx, &configuration_service.GetConfigAppListRequest{
		Offset:    0,
		Limit:     100,
		Attribute: req.RegionID,
		Field:     "EDN App Checkbook ID",
	})

	for _, v := range Checkbook.Data {
		res = append(res, &dto.CheckbookGP{
			ID:            v.Value,
			CheckbookCode: v.Value,
			CheckbookDesc: v.Attribute,
		})
	}

	if req.RegionID != "" {
		Checkbook, _ = s.opt.Client.ConfigurationServiceGrpc.GetConfigAppList(ctx, &configuration_service.GetConfigAppListRequest{
			Offset:    0,
			Limit:     100,
			Attribute: "ALL",
			Field:     "EDN App Checkbook ID",
		})

		for _, v := range Checkbook.Data {
			res = append(res, &dto.CheckbookGP{
				ID:            v.Value,
				CheckbookCode: v.Value,
				CheckbookDesc: v.Attribute,
			})
		}

	}

	total = int64(len(res))

	return
}

// func (s *CheckbookService) GetDetaiGPlById(ctx context.Context, id string) (res *dto.CheckbookGP, err error) {
// 	ctx, span := s.opt.Trace.Start(ctx, "CheckbookService.GetDetaiGPlById")
// 	defer span.End()

// 	// get Checkbook from bridge
// 	Checkbook, err := s.opt.Client.SalesServiceGrpc.GetCheckbookList(ctx, &sales_service.GetCheckbookListRequest{
// 		Status: 1,
// 		Id:     id,
// 	})
// 	if err != nil {
// 		span.RecordError(err)
// 		s.opt.Logger.AddMessage(log.ErrorLevel, err)
// 		err = edenlabs.ErrorRpcNotFound("bridge", "payment method")
// 		return
// 	}

// 	if len(Checkbook.Data) > 0 {
// 		res = &dto.CheckbookGP{
// 			ID:                Checkbook.Data[0].Code,
// 			CheckbookDesc: Checkbook.Data[0].Name,
// 			CheckbookCode: Checkbook.Data[0].Code,
// 		}
// 	}

// 	return
// }
