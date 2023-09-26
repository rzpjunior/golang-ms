package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/edenlabs/edenlabs/utils"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

type IArchetypeService interface {
	Get(ctx context.Context, req *dto.ArchetypeGetListRequest) (res []*dto.ArchetypeResponse, total int64, err error)
	GetDetail(ctx context.Context, id string) (res *dto.ArchetypeResponse, err error)
}

type ArchetypeService struct {
	opt opt.Options
}

func NewArchetypeService() IArchetypeService {
	return &ArchetypeService{
		opt: global.Setup.Common,
	}
}

func (s *ArchetypeService) Get(ctx context.Context, req *dto.ArchetypeGetListRequest) (res []*dto.ArchetypeResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "Archetype.Get")
	defer span.End()

	var (
		archetypes                          *bridge_service.GetArchetypeGPResponse
		statusArchetype, statusCustomerType int8
		statusGP                            string
	)

	if req.Status != 0 {
		switch req.Status {
		case 1:
			statusGP = "0"
		case 7:
			statusGP = "1"
		default:
			statusGP = utils.ToString(req.Status)
		}
	}

	archetypes, err = s.opt.Client.BridgeServiceGrpc.GetArchetypeGPList(ctx, &bridge_service.GetArchetypeGPListRequest{
		Limit:                   int32(req.Limit),
		Offset:                  int32(req.Offset) - 1,
		GnlArchetypedescription: req.Search,
		Inactive:                statusGP,
		GnlCustTypeId:           req.CustomerTypeID,
	})
	if err != nil || len(archetypes.Data) == 0 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, v := range archetypes.Data {
		var customerType *bridge_service.GetCustomerTypeGPResponse
		customerType, err = s.opt.Client.BridgeServiceGrpc.GetCustomerTypeGPDetail(ctx, &bridge_service.GetCustomerTypeGPDetailRequest{
			Id: v.GnlCustTypeId,
		})
		if err != nil {
			// span.RecordError(err)
			// s.opt.Logger.AddMessage(log.ErrorLevel, err)
			// return
			continue
		}

		if v.Inactive == 0 {
			statusArchetype = statusx.ConvertStatusName(statusx.Active)
		} else {
			statusArchetype = statusx.ConvertStatusName(statusx.Archived)
		}

		if customerType.Data[0].Inactive == 0 {
			statusCustomerType = statusx.ConvertStatusName(statusx.Active)
		} else {
			statusCustomerType = statusx.ConvertStatusName(statusx.Archived)
		}

		// Manual filter customer type
		res = append(res, &dto.ArchetypeResponse{
			ID:             v.GnlArchetypeId,
			Code:           v.GnlArchetypeId,
			Description:    v.GnlArchetypedescription,
			Status:         statusArchetype,
			ConvertStatus:  statusx.ConvertStatusValue(statusArchetype),
			CustomerTypeID: v.GnlCustTypeId,
			CustomerType: &dto.CustomerTypeResponse{
				ID:            customerType.Data[0].GnL_Cust_Type_ID,
				Code:          customerType.Data[0].GnL_Cust_Type_ID,
				Description:   customerType.Data[0].GnL_CustType_Description,
				Status:        statusCustomerType,
				ConvertStatus: statusx.ConvertStatusValue(statusArchetype),
				CustomerGroup: customerType.Data[0].GnL_Cust_GroupDesc,
			},
		})

	}

	total = int64(archetypes.TotalRecords)

	return
}

func (s *ArchetypeService) GetDetail(ctx context.Context, id string) (res *dto.ArchetypeResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "Archetype.GetDetail")
	defer span.End()

	var (
		archetype                           *bridge_service.GetArchetypeGPResponse
		statusArchetype, statusCustomerType int8
	)

	archetype, err = s.opt.Client.BridgeServiceGrpc.GetArchetypeGPDetail(ctx, &bridge_service.GetArchetypeGPDetailRequest{
		Id: id,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	var customerType *bridge_service.GetCustomerTypeGPResponse
	customerType, err = s.opt.Client.BridgeServiceGrpc.GetCustomerTypeGPDetail(ctx, &bridge_service.GetCustomerTypeGPDetailRequest{
		Id: archetype.Data[0].GnlCustTypeId,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if archetype.Data[0].Inactive == 0 {
		statusArchetype = statusx.ConvertStatusName(statusx.Active)
	} else {
		statusArchetype = statusx.ConvertStatusName(statusx.Archived)
	}

	if customerType.Data[0].Inactive == 0 {
		statusCustomerType = statusx.ConvertStatusName(statusx.Active)
	} else {
		statusCustomerType = statusx.ConvertStatusName(statusx.Archived)
	}

	res = &dto.ArchetypeResponse{
		ID:             archetype.Data[0].GnlArchetypeId,
		Code:           archetype.Data[0].GnlArchetypeId,
		Description:    archetype.Data[0].GnlArchetypedescription,
		CustomerTypeID: archetype.Data[0].GnlCustTypeId,
		Status:         statusArchetype,
		ConvertStatus:  statusx.ConvertStatusValue(statusArchetype),
		CustomerType: &dto.CustomerTypeResponse{
			ID:            customerType.Data[0].GnL_Cust_Type_ID,
			Code:          customerType.Data[0].GnL_Cust_Type_ID,
			Description:   customerType.Data[0].GnL_CustType_Description,
			Status:        statusCustomerType,
			ConvertStatus: statusx.ConvertStatusValue(statusCustomerType),
			CustomerGroup: customerType.Data[0].GnL_Cust_GroupDesc,
		},
	}

	return
}
