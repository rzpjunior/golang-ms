package service

import (
	"context"
	"fmt"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

type IArchetypeService interface {
	Get(ctx context.Context, offset int, limit int, search string, customer_type_id int, Type int) (res []dto.ArchetypeResponse, total int64, err error)
}

type ArchetypeService struct {
	opt opt.Options
}

func NewArchetypeService() IArchetypeService {
	return &ArchetypeService{
		opt: global.Setup.Common,
	}
}

func (s *ArchetypeService) Get(ctx context.Context, offset int, limit int, search string, customer_type_id int, Type int) (res []dto.ArchetypeResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ArchetypeService.Get")
	defer span.End()

	// cusTypeId := int64(customer_type_id)
	// Archetype, err := s.opt.Client.BridgeServiceGrpc.GetArchetypeList(ctx, &bridge_service.GetArchetypeListRequest{
	// 	Limit:          int32(limit),
	// 	Offset:         int32(offset),
	// 	Search:         search,
	// 	CustomerTypeId: cusTypeId,
	// })

	ArchetypeGP, _ := s.opt.Client.BridgeServiceGrpc.GetArchetypeGPList(ctx, &bridge_service.GetArchetypeGPListRequest{
		Limit:  int32(1000),
		Offset: int32(1),
	})

	// fmt.Println(Archetype, ArchetypeGP)

	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// for _, archetype := range Archetype.Data {
	// 	customerType, err := s.opt.Client.BridgeServiceGrpc.GetCustomerTypeDetail(ctx, &bridge_service.GetCustomerTypeDetailRequest{
	// 		Id: cusTypeId})

	// 	if err != nil {
	// 		span.RecordError(err)
	// 		s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 		return res, 0, err
	// 	}
	// 	customerTypeGP, _ := s.opt.Client.BridgeServiceGrpc.GetCustomerTypeGPDetail(ctx, &bridge_service.GetCustomerTypeGPDetailRequest{
	// 		// Id: archetype.CustomerTypeId,
	// 	})
	// 	fmt.Println(customerTypeGP)

	// 	content := archetype.CustomerTypeId
	// 	var docRequired string
	// 	if content == 2 || content == 4 {
	// 		docRequired = "1"
	// 	} else {
	// 		docRequired = "2"
	// 	}
	// 	res = append(res, dto.ArchetypeResponse{
	// 		ID:               strconv.Itoa(int(archetype.Id)),
	// 		Code:             archetype.Code,
	// 		CustomerTypeID:   strconv.Itoa(int(cusTypeId)),
	// 		Description:      archetype.Description,
	// 		Status:           strconv.Itoa(int(archetype.Status)),
	// 		StatusConvert:    "",
	// 		CreatedAt:        archetype.CreatedAt.AsTime(),
	// 		UpdatedAt:        archetype.UpdatedAt.AsTime(),
	// 		CustomerGroup:    "",
	// 		Name:             "Dummy Archetype Name",
	// 		NameID:           "Dummy Archetype Name ID",
	// 		Abbreviation:     "Dummy Abbreviation",
	// 		Note:             "Dummy Note",
	// 		AuxData:          "",
	// 		DocRequired:      docRequired,
	// 		DocImageRequired: []string{},
	// 		CustomerType: &model.CustomerType{
	// 			ID:               customerType.Data.Id,
	// 			Code:             customerType.Data.Code,
	// 			Name:             customerType.Data.Description,
	// 			Note:             "",
	// 			AuxData:          0,
	// 			Status:           int8(customerType.Data.Status),
	// 			DocImageRequired: "",
	// 		},
	// 	})

	// }
	// total = int64(len(res))

	for _, archetype := range ArchetypeGP.Data {
		customerTypeGP, _ := s.opt.Client.BridgeServiceGrpc.GetCustomerTypeGPDetail(ctx, &bridge_service.GetCustomerTypeGPDetailRequest{
			Id: archetype.GnL_Cust_Type_ID,
		})
		fmt.Println(customerTypeGP)

		content := archetype.GnL_Cust_Type_ID
		var docRequired string
		if content == "BTY0002" || content == "BTY0004" {
			docRequired = "1"
		} else {
			docRequired = "2"
		}
		res = append(res, dto.ArchetypeResponse{
			ID:             archetype.GnL_Archetype_ID,
			Code:           archetype.GnL_Archetype_ID,
			CustomerTypeID: archetype.GnL_Cust_Type_ID,
			Description:    archetype.GnL_ArchetypeDescription,
			Status:         archetype.InactivE_DESC,
			StatusConvert:  "",
			// CreatedAt:        archetype.CreatedAt.AsTime(),
			// UpdatedAt:        archetype.UpdatedAt.AsTime(),
			CustomerGroup: "",
			Name:          archetype.GnL_ArchetypeDescription,
			NameID:        archetype.GnL_Archetype_ID,
			// Abbreviation:     "Dummy Abbreviation",
			// Note:             "Dummy Note",
			AuxData:          "",
			DocRequired:      docRequired,
			DocImageRequired: []string{},
			CustomerType: &model.CustomerType{
				// ID:               customerType.Data.Id,
				Code:    customerTypeGP.Data[0].GnL_Cust_Type_ID,
				Name:    customerTypeGP.Data[0].GnL_CustType_Description,
				Note:    "",
				AuxData: 0,
				// Status:           customerTypeGP.Data[0].,
				DocImageRequired: "",
			},
		})

	}
	return
}
