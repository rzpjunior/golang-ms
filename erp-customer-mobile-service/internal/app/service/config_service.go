package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
)

type IConfigService interface {
	GetAppConfig(ctx context.Context, application int32, field string, attribute string, value string) (res []dto.ApplicationConfigResponse, total int64, err error)
	// GetAppConfigByID(ctx context.Context, id int64) (res dto.ApplicationConfigResponse, err error)
	// UpdateAppConfig(ctx context.Context, req dto.ApplicationConfigRequestUpdate, id int64) (res dto.ApplicationConfigResponse, err error)
	GetGlossary(ctx context.Context, table string, attribute string, valueInt int, valueName string) (res []dto.GlossaryResponse, total int64, err error)
	GetDeliveryFee(ctx context.Context, req dto.RequestGetDeliveryFee) (res *dto.ResponseGetDeliveryFee, err error)
}

type ConfigService struct {
	opt opt.Options
}

func NewConfigService() IConfigService {
	return &ConfigService{
		opt: global.Setup.Common,
	}
}

func (s *ConfigService) GetAppConfig(ctx context.Context, application int32, field string, attribute string, value string) (res []dto.ApplicationConfigResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ConfigService.GetAppConfig")
	defer span.End()

	var applicationConfigs []*model.ApplicationConfig

	appConfig, err := s.opt.Client.ConfigurationServiceGrpc.GetConfigAppList(ctx, &configuration_service.GetConfigAppListRequest{
		Application: application,
		Field:       field,
		Attribute:   attribute,
		Value:       value,
	})

	for _, applicationConfig := range appConfig.Data {
		res = append(res, dto.ApplicationConfigResponse{
			ID:          strconv.Itoa(int((applicationConfig.Id))),
			Application: strconv.Itoa(int((applicationConfig.Application))),
			Field:       applicationConfig.Field,
			Attribute:   applicationConfig.Attribute,
			Value:       applicationConfig.Value,
		})
	}
	total = int64(len(applicationConfigs))

	return
}

func (s *ConfigService) GetGlossary(ctx context.Context, table string, attribute string, valueInt int, valueName string) (res []dto.GlossaryResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ConfigService.Get")
	defer span.End()

	glossary, err := s.opt.Client.ConfigurationServiceGrpc.GetGlossaryList(ctx, &configuration_service.GetGlossaryListRequest{
		Table:     table,
		Attribute: attribute,
		ValueInt:  int32(valueInt),
		ValueName: valueName,
	})

	for _, glossary := range glossary.Data {
		res = append(res, dto.GlossaryResponse{
			ID:        strconv.Itoa(int(glossary.Id)),
			Table:     glossary.Table,
			Attribute: glossary.Attribute,
			ValueInt:  strconv.Itoa(int(glossary.ValueInt)),
			ValueName: glossary.ValueName,
			Note:      "",
		})
	}

	total = int64(len(glossary.Data))

	return
}

func (s *ConfigService) GetDeliveryFee(ctx context.Context, req dto.RequestGetDeliveryFee) (res *dto.ResponseGetDeliveryFee, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ConfigService.Get")
	defer span.End()
	// regionID, _ := strconv.Atoi(req.Data.RegionID)
	// region, err := s.opt.Client.BridgeServiceGrpc.GetRegionDetail(ctx, &bridge_service.GetRegionDetailRequest{
	// 	Id: int64(regionID),
	// })
	region, err := s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPList(ctx, &bridge_service.GetAdmDivisionGPListRequest{
		Region: req.Data.RegionID,
		Limit:  1,
		Offset: 0,
	})
	fmt.Print(region)

	if err != nil {
		err = errors.New("Region invalid data.")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	// customerTypeID, _ := strconv.Atoi(req.Data.CustomerTypeID)

	// customerType, err := s.opt.Client.BridgeServiceGrpc.GetCustomerTypeDetail(ctx, &bridge_service.GetCustomerTypeDetailRequest{
	// 	Id: int64(customerTypeID),
	// })
	// if err != nil {
	// 	err = errors.New("Customer Type invalid data.")
	// 	span.RecordError(err)
	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 	return
	// }
	customerTypeGP, err := s.opt.Client.BridgeServiceGrpc.GetCustomerTypeGPDetail(ctx, &bridge_service.GetCustomerTypeGPDetailRequest{
		Id: req.Data.CustomerTypeID,
	})
	if err != nil {
		err = errors.New("Customer Type invalid data.")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// regionPolicy, err := s.opt.Client.ConfigurationServiceGrpc.GetRegionPolicyList(ctx, &configuration_service.GetRegionPolicyListRequest{
	// 	RegionId: int64(region.Data.Id),
	// })
	// if err != nil {
	// 	err = errors.New("Region Policy invalid data.")
	// 	span.RecordError(err)
	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 	return
	// }

	// deliveryFee, err := s.opt.Client.BridgeServiceGrpc.GetDeliveryFeeList(ctx, &bridge_service.GetDeliveryFeeListRequest{
	// 	RegionId: 1,
	// })
	deliveryFee, err := s.opt.Client.BridgeServiceGrpc.GetDeliveryFeeGPList(ctx, &bridge_service.GetDeliveryFeeGPListRequest{
		Limit:     100,
		Offset:    0,
		GnlRegion: region.Data[0].Region,
	})
	if err != nil {
		err = errors.New("Region Business Policy invalid data.")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	var exist bool
	for _, v := range deliveryFee.Data {
		if v.GnlCustTypeId == customerTypeGP.Data[0].GnL_Cust_Type_ID {
			res = &dto.ResponseGetDeliveryFee{
				// ID:          strconv.Itoa(int(v.)),
				MinOrder:    strconv.Itoa(int(v.Minorqty)),
				DeliveryFee: strconv.Itoa(int(v.GnlDeliveryFee)),
			}
			exist = true
			break
		}
	}
	if !exist {
		for _, v := range deliveryFee.Data {
			if v.GnlCustTypeId == "" {
				res = &dto.ResponseGetDeliveryFee{
					// ID:          strconv.Itoa(int(v.Id)),
					MinOrder:    strconv.Itoa(int(v.Minorqty)),
					DeliveryFee: strconv.Itoa(int(v.GnlDeliveryFee)),
				}
				break
			}
		}
	}
	// for _, glossary := range glossary.Data {
	// 	res = append(res, dto.GlossaryResponse{
	// 		ID:        int64(glossary.Id),
	// 		Table:     glossary.Table,
	// 		Attribute: glossary.Attribute,
	// 		ValueInt:  int8(glossary.ValueInt),
	// 		ValueName: glossary.ValueName,
	// 		Note:      "",
	// 	})
	// }

	//total = int64(len(res))

	return
}
