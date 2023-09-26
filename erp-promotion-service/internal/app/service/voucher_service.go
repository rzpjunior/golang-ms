package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/constants"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/edenlabs/edenlabs/utils"
	"git.edenfarm.id/project-version3/erp-services/erp-promotion-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-promotion-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-promotion-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-promotion-service/internal/app/repository"
	"google.golang.org/protobuf/types/known/timestamppb"

	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/account_service"
	auditService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/audit_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	campaignService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/campaign_service"
	catalogService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/catalog_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
	configurationService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
	crmService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/crm_service"
)

type IVoucherService interface {
	Get(ctx context.Context, req *dto.VoucherRequestGet) (res []*dto.VoucherResponse, total int64, err error)
	GetDetail(ctx context.Context, id int64) (res *dto.VoucherResponse, err error)
	Create(ctx context.Context, req *dto.VoucherRequestCreate) (res *dto.VoucherResponse, err error)
	Archive(ctx context.Context, id int64) (res *dto.VoucherResponse, err error)
	CreateBulky(ctx context.Context, req *dto.VoucherRequestBulky) (err error)
	GetMobileVoucherList(ctx context.Context, req *dto.VoucherRequestGetMobileVoucherList) (res []*dto.VoucherResponse, count int64, err error)
	GetMobileVoucherDetail(ctx context.Context, req *dto.VoucherRequestGetMobileVoucherDetail) (res *dto.VoucherResponse, err error)
	Update(ctx context.Context, req *dto.VoucherRequestUpdate) (err error)
}

type VoucherService struct {
	opt                   opt.Options
	RepositoryVoucher     repository.IVoucherRepository
	RepositoryVoucherItem repository.IVoucherItemRepository
}

func NewVoucherService() IVoucherService {
	return &VoucherService{
		opt:                   global.Setup.Common,
		RepositoryVoucher:     repository.NewVoucherRepository(),
		RepositoryVoucherItem: repository.NewVoucherItemRepository(),
	}
}

func (s *VoucherService) Get(ctx context.Context, req *dto.VoucherRequestGet) (res []*dto.VoucherResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "VoucherService.Get")
	defer span.End()

	// Validation param region
	if req.RegionID != "" {
		_, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPDetail(ctx, &bridgeService.GetAdmDivisionGPDetailRequest{
			Region: req.RegionID,
			Limit:  1,
			Offset: 0,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("region_id")
			return
		}
	}

	// validation param archetype
	if req.ArchetypeID != "" {
		var archetypeGP *bridgeService.GetArchetypeGPResponse
		archetypeGP, err = s.opt.Client.BridgeServiceGrpc.GetArchetypeGPDetail(ctx, &bridge_service.GetArchetypeGPDetailRequest{
			Id: req.ArchetypeID,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("archetype_id")
			return
		}
		req.CustomerTypeID = archetypeGP.Data[0].GnlCustTypeId
	}

	var vouchers []*model.Voucher
	vouchers, total, err = s.RepositoryVoucher.Get(ctx, req)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, voucher := range vouchers {
		var (
			customerResponse                    *dto.CustomerResponse
			regionResponse                      *dto.RegionResponse
			archetypeResponse                   *dto.ArchetypeResponse
			membershipLevelResponse             *dto.MembershipLevelResponse
			membershipCheckpointResponse        *dto.MembershipCheckpointResponse
			divisionResponse                    *dto.DivisionResponse
			archetype                           *bridge_service.GetArchetypeGPResponse
			customerType                        *bridge_service.GetCustomerTypeGPResponse
			customerTypeResponse                *dto.CustomerTypeResponse
			statusArchetype, statusCustomerType int8
		)

		if voucher.RegionIDGP == "" {
			regionResponse = &dto.RegionResponse{
				Description: "All Region",
			}
		} else {
			// Get Region from bridge
			var region *bridgeService.GetAdmDivisionGPResponse
			region, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPDetail(ctx, &bridgeService.GetAdmDivisionGPDetailRequest{
				Region: voucher.RegionIDGP,
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
		}

		if voucher.CustomerTypeIDGP == "" {
			if voucher.ArchetypeIDGP != "" {
				err = edenlabs.ErrorInvalid("archetype_id")
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}
			customerTypeResponse = &dto.CustomerTypeResponse{
				Description: "All Customer Type",
			}
			archetypeResponse = &dto.ArchetypeResponse{
				Description:  fmt.Sprintf("All Archetype"),
				CustomerType: customerTypeResponse,
			}
		} else {
			// Get customer type from bridge
			customerType, err = s.opt.Client.BridgeServiceGrpc.GetCustomerTypeGPDetail(ctx, &bridge_service.GetCustomerTypeGPDetailRequest{
				Id: voucher.CustomerTypeIDGP,
			})
			if err != nil || len(customerType.Data) == 0 {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorInvalid("customer_type_id")
				return
			}
			if customerType.Data[0].Inactive == 0 {
				statusCustomerType = statusx.ConvertStatusName(statusx.Active)
			} else {
				statusCustomerType = statusx.ConvertStatusName(statusx.Archived)
			}
			customerTypeResponse = &dto.CustomerTypeResponse{
				ID:            customerType.Data[0].GnL_Cust_Type_ID,
				Code:          customerType.Data[0].GnL_Cust_Type_ID,
				Description:   customerType.Data[0].GnL_CustType_Description,
				CustomerGroup: customerType.Data[0].GnL_Cust_GroupDesc,
				Status:        statusCustomerType,
				ConvertStatus: statusx.ConvertStatusValue(statusCustomerType),
			}
		}

		if voucher.ArchetypeIDGP == "" && voucher.CustomerTypeIDGP != "" {
			archetypeResponse = &dto.ArchetypeResponse{
				Description:  "All Archetype",
				CustomerType: customerTypeResponse,
			}
		} else if voucher.CustomerTypeIDGP != "" && voucher.ArchetypeIDGP != "" {
			// Get Archetype from bridge
			archetype, err = s.opt.Client.BridgeServiceGrpc.GetArchetypeGPDetail(ctx, &bridge_service.GetArchetypeGPDetailRequest{
				Id: voucher.ArchetypeIDGP,
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorInvalid("archetype_id")
				return
			}
			if archetype.Data[0].Inactive == 0 {
				statusArchetype = statusx.ConvertStatusName(statusx.Active)
			} else {
				statusArchetype = statusx.ConvertStatusName(statusx.Archived)
			}
			archetypeResponse = &dto.ArchetypeResponse{
				ID:             archetype.Data[0].GnlArchetypeId,
				Code:           archetype.Data[0].GnlArchetypeId,
				Description:    archetype.Data[0].GnlArchetypedescription,
				CustomerTypeID: archetype.Data[0].GnlCustTypeId,
				CustomerType:   customerTypeResponse,
				Status:         statusArchetype,
				ConvertStatus:  statusx.ConvertStatusValue(statusArchetype),
			}
		}

		if voucher.DivisionID != 0 {
			// Get division from bridge
			var division *account_service.GetDivisionDetailResponse
			division, err = s.opt.Client.AccountServiceGrpc.GetDivisionDetail(ctx, &account_service.GetDivisionDetailRequest{
				Id: voucher.DivisionID,
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "division")
				return
			}

			divisionResponse = &dto.DivisionResponse{
				ID:            division.Data.Id,
				Code:          division.Data.Code,
				Description:   division.Data.Name,
				Status:        int8(division.Data.Status),
				StatusConvert: statusx.ConvertStatusValue(int8(division.Data.Status)),
			}
		}

		if voucher.CustomerID != 0 {
			// Get customer from crm
			var customer *crmService.GetCustomerDetailResponse
			customer, err = s.opt.Client.CrmServiceGrpc.GetCustomerDetail(ctx, &crmService.GetCustomerDetailRequest{
				Id: voucher.CustomerID,
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("crm", "customer")
				return
			}
			var customergp *bridgeService.GetCustomerGPResponse
			customergp, err = s.opt.Client.BridgeServiceGrpc.GetCustomerGPList(ctx, &bridgeService.GetCustomerGPListRequest{
				Id:     customer.Data.CustomerIdGp,
				Limit:  1,
				Offset: 0,
			})
			if err != nil || len(customergp.Data) == 0 {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "customer")
				return
			}
			customerResponse = &dto.CustomerResponse{
				ID:   customer.Data.Id,
				Code: customer.Data.CustomerIdGp,
				Name: customergp.Data[0].Custname,
			}
		}

		// Get membership level from campaign service
		if voucher.MembershipLevelID != 0 {
			var membershipLevel *campaignService.GetMembershipLevelDetailResponse
			membershipLevel, err = s.opt.Client.CampaignServiceGrpc.GetMembershipLevelDetail(ctx, &campaignService.GetMembershipLevelDetailRequest{
				Id: voucher.MembershipLevelID,
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("campaign", "membership_level")
				return
			}
			membershipLevelResponse = &dto.MembershipLevelResponse{
				ID:       membershipLevel.Data.Id,
				Code:     membershipLevel.Data.Code,
				Level:    int8(membershipLevel.Data.Level),
				Name:     membershipLevel.Data.Name,
				ImageUrl: membershipLevel.Data.ImageUrl,
				Status:   int8(membershipLevel.Data.Status),
			}
		}

		// Get membership level from campaign service
		if voucher.MembershipCheckPointID != 0 {
			var membershipCheckpoint *campaignService.GetMembershipCheckpointDetailResponse
			membershipCheckpoint, err = s.opt.Client.CampaignServiceGrpc.GetMembershipCheckpointDetail(ctx, &campaignService.GetMembershipCheckpointDetailRequest{
				Id: voucher.MembershipCheckPointID,
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("campaign", "membership_checkpoint")
				return
			}
			membershipCheckpointResponse = &dto.MembershipCheckpointResponse{
				ID:                membershipCheckpoint.Data.Id,
				Checkpoint:        int8(membershipCheckpoint.Data.Checkpoint),
				TargetAmount:      membershipCheckpoint.Data.TargetAmount,
				Status:            int8(membershipCheckpoint.Data.Status),
				MembershipLevelID: membershipCheckpoint.Data.MembershipLevelId,
			}
		}

		res = append(res, &dto.VoucherResponse{
			ID:                   voucher.ID,
			Code:                 voucher.Code,
			RedeemCode:           voucher.RedeemCode,
			Name:                 voucher.Name,
			Type:                 voucher.Type,
			StartTime:            voucher.StartTime,
			EndTime:              voucher.EndTime,
			OverallQuota:         voucher.OverallQuota,
			UserQuota:            voucher.UserQuota,
			RemOverallQuota:      voucher.RemOverallQuota,
			MinOrder:             voucher.MinOrder,
			DiscAmount:           voucher.DiscAmount,
			TermConditions:       voucher.TermConditions,
			ImageUrl:             voucher.ImageUrl,
			VoidReason:           voucher.VoidReason,
			Note:                 voucher.Note,
			Status:               voucher.Status,
			StatusConvert:        statusx.ConvertStatusValue(voucher.Status),
			VoucherItem:          voucher.VoucherItem,
			CreatedAt:            voucher.CreatedAt,
			Region:               regionResponse,
			Archetype:            archetypeResponse,
			Customer:             customerResponse,
			MembershipLevel:      membershipLevelResponse,
			MembershipCheckpoint: membershipCheckpointResponse,
			Division:             divisionResponse,
		})
	}

	return
}

func (s *VoucherService) GetDetail(ctx context.Context, id int64) (res *dto.VoucherResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "VoucherService.GetDetail")
	defer span.End()

	var (
		voucher                             *model.Voucher
		customerResponse                    *dto.CustomerResponse
		regionResponse                      *dto.RegionResponse
		archetypeResponse                   *dto.ArchetypeResponse
		membershipLevelResponse             *dto.MembershipLevelResponse
		membershipCheckpointResponse        *dto.MembershipCheckpointResponse
		voucherItems                        []*model.VoucherItem
		divisionResponse                    *dto.DivisionResponse
		voucherItemsResponse                []*dto.VoucherItemResponse
		archetype                           *bridgeService.GetArchetypeGPResponse
		customerType                        *bridgeService.GetCustomerTypeGPResponse
		customerTypeResponse                *dto.CustomerTypeResponse
		statusArchetype, statusCustomerType int8
	)

	voucher, err = s.RepositoryVoucher.GetDetail(ctx, id)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if voucher.RegionIDGP == "" {
		regionResponse = &dto.RegionResponse{
			Description: "All Region",
		}
	} else {
		// Get Region from bridge
		var region *bridgeService.GetAdmDivisionGPResponse
		region, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPDetail(ctx, &bridgeService.GetAdmDivisionGPDetailRequest{
			Region: voucher.RegionIDGP,
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
	}

	if voucher.CustomerTypeIDGP == "" {
		if voucher.ArchetypeIDGP != "" {
			err = edenlabs.ErrorInvalid("archetype_id")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		customerTypeResponse = &dto.CustomerTypeResponse{
			Description: "All Customer Type",
		}
		archetypeResponse = &dto.ArchetypeResponse{
			Description:  "All Archetype",
			CustomerType: customerTypeResponse,
		}
	} else {
		// Get customer type from bridge
		customerType, err = s.opt.Client.BridgeServiceGrpc.GetCustomerTypeGPDetail(ctx, &bridge_service.GetCustomerTypeGPDetailRequest{
			Id: voucher.CustomerTypeIDGP,
		})
		if err != nil || len(customerType.Data) == 0 {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("customer_type_id")
			return
		}
		if customerType.Data[0].Inactive == 0 {
			statusCustomerType = statusx.ConvertStatusName(statusx.Active)
		} else {
			statusCustomerType = statusx.ConvertStatusName(statusx.Archived)
		}
		customerTypeResponse = &dto.CustomerTypeResponse{
			ID:            customerType.Data[0].GnL_Cust_Type_ID,
			Code:          customerType.Data[0].GnL_Cust_Type_ID,
			Description:   customerType.Data[0].GnL_CustType_Description,
			CustomerGroup: customerType.Data[0].GnL_Cust_GroupDesc,
			Status:        statusCustomerType,
			ConvertStatus: statusx.ConvertStatusValue(statusCustomerType),
		}
	}

	if voucher.ArchetypeIDGP == "" && voucher.CustomerTypeIDGP != "" {
		archetypeResponse = &dto.ArchetypeResponse{
			Description:  "All Archetype",
			CustomerType: customerTypeResponse,
		}
	} else if voucher.CustomerTypeIDGP != "" && voucher.ArchetypeIDGP != "" {
		// Get Archetype from bridge
		archetype, err = s.opt.Client.BridgeServiceGrpc.GetArchetypeGPDetail(ctx, &bridge_service.GetArchetypeGPDetailRequest{
			Id: voucher.ArchetypeIDGP,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("archetype_id")
			return
		}
		if archetype.Data[0].Inactive == 0 {
			statusArchetype = statusx.ConvertStatusName(statusx.Active)
		} else {
			statusArchetype = statusx.ConvertStatusName(statusx.Archived)
		}
		archetypeResponse = &dto.ArchetypeResponse{
			ID:             archetype.Data[0].GnlArchetypeId,
			Code:           archetype.Data[0].GnlArchetypeId,
			Description:    archetype.Data[0].GnlArchetypedescription,
			CustomerTypeID: archetype.Data[0].GnlCustTypeId,
			CustomerType:   customerTypeResponse,
			Status:         statusArchetype,
			ConvertStatus:  statusx.ConvertStatusValue(statusArchetype),
		}
	}

	if voucher.DivisionID != 0 {
		// Get division from bridge
		var division *account_service.GetDivisionDetailResponse
		division, err = s.opt.Client.AccountServiceGrpc.GetDivisionDetail(ctx, &account_service.GetDivisionDetailRequest{
			Id: voucher.DivisionID,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "division")
			return
		}

		divisionResponse = &dto.DivisionResponse{
			ID:            division.Data.Id,
			Code:          division.Data.Code,
			Description:   division.Data.Name,
			Status:        int8(division.Data.Status),
			StatusConvert: statusx.ConvertStatusValue(int8(division.Data.Status)),
		}
	}

	// Get Customer from bridge
	if voucher.CustomerID != 0 {
		// Get customer from crm
		var customer *crmService.GetCustomerDetailResponse
		customer, err = s.opt.Client.CrmServiceGrpc.GetCustomerDetail(ctx, &crmService.GetCustomerDetailRequest{
			Id: voucher.CustomerID,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("crm", "customer")
			return
		}
		var customergp *bridgeService.GetCustomerGPResponse
		customergp, err = s.opt.Client.BridgeServiceGrpc.GetCustomerGPList(ctx, &bridgeService.GetCustomerGPListRequest{
			Id:     customer.Data.CustomerIdGp,
			Limit:  1,
			Offset: 0,
		})
		if err != nil || len(customergp.Data) == 0 {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "customer")
			return
		}
		customerResponse = &dto.CustomerResponse{
			ID:   customer.Data.Id,
			Code: customer.Data.CustomerIdGp,
			Name: customergp.Data[0].Custname,
		}
	}

	// Get membership level from campaign service
	if voucher.MembershipLevelID != 0 {
		var membershipLevel *campaignService.GetMembershipLevelDetailResponse
		membershipLevel, err = s.opt.Client.CampaignServiceGrpc.GetMembershipLevelDetail(ctx, &campaignService.GetMembershipLevelDetailRequest{
			Id: voucher.MembershipLevelID,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("campaign", "membership_level")
			return
		}
		membershipLevelResponse = &dto.MembershipLevelResponse{
			ID:       membershipLevel.Data.Id,
			Code:     membershipLevel.Data.Code,
			Level:    int8(membershipLevel.Data.Level),
			Name:     membershipLevel.Data.Name,
			ImageUrl: membershipLevel.Data.ImageUrl,
			Status:   int8(membershipLevel.Data.Status),
		}
	}

	// Get membership level from campaign service
	if voucher.MembershipCheckPointID != 0 {
		var membershipCheckpoint *campaignService.GetMembershipCheckpointDetailResponse
		membershipCheckpoint, err = s.opt.Client.CampaignServiceGrpc.GetMembershipCheckpointDetail(ctx, &campaignService.GetMembershipCheckpointDetailRequest{
			Id: voucher.MembershipCheckPointID,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("campaign", "membership_checkpoint")
			return
		}
		membershipCheckpointResponse = &dto.MembershipCheckpointResponse{
			ID:                membershipCheckpoint.Data.Id,
			Checkpoint:        int8(membershipCheckpoint.Data.Checkpoint),
			TargetAmount:      membershipCheckpoint.Data.TargetAmount,
			Status:            int8(membershipCheckpoint.Data.Status),
			MembershipLevelID: membershipCheckpoint.Data.MembershipLevelId,
		}
	}

	if voucher.VoucherItem == 1 {
		filter := &dto.VoucherItemRequestGet{
			VoucherID: voucher.ID,
		}

		voucherItems, _, err = s.RepositoryVoucherItem.Get(ctx, filter)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("voucher_item")
			return
		}

		for _, v := range voucherItems {
			var item *catalogService.GetItemDetailByInternalIdResponse
			item, err = s.opt.Client.CatalogServiceGrpc.GetItemDetailByInternalId(ctx, &catalogService.GetItemDetailByInternalIdRequest{
				Id: utils.ToString(v.ItemID),
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("catalog", "item")
				return
			}

			var uom *bridgeService.GetUomGPResponse
			uom, err = s.opt.Client.BridgeServiceGrpc.GetUomGPDetail(ctx, &bridgeService.GetUomGPDetailRequest{
				Id: item.Data.UomId,
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("catalog", "uom")
				return
			}

			detailUom := &dto.UomResponse{
				ID:          uom.Data[0].Uomschdl,
				Code:        uom.Data[0].Uomschdl,
				Description: uom.Data[0].Umschdsc,
			}

			detailItem := &dto.ItemResponse{
				ID:                   item.Data.Id,
				Code:                 item.Data.Code,
				Description:          item.Data.Description,
				UnitWeightConversion: item.Data.UnitWeightConversion,
				Uom:                  detailUom,
			}

			voucherItemsResponse = append(voucherItemsResponse, &dto.VoucherItemResponse{
				ID:         v.ID,
				VoucherID:  v.VoucherID,
				ItemID:     v.ItemID,
				MinQtyDisc: v.MinQtyDisc,
				Item:       detailItem,
				CreatedAt:  v.CreatedAt,
			})
		}
	}

	res = &dto.VoucherResponse{
		ID:                   voucher.ID,
		Code:                 voucher.Code,
		RedeemCode:           voucher.RedeemCode,
		Name:                 voucher.Name,
		Type:                 voucher.Type,
		StartTime:            voucher.StartTime,
		EndTime:              voucher.EndTime,
		OverallQuota:         voucher.OverallQuota,
		UserQuota:            voucher.UserQuota,
		RemOverallQuota:      voucher.RemOverallQuota,
		MinOrder:             voucher.MinOrder,
		DiscAmount:           voucher.DiscAmount,
		TermConditions:       voucher.TermConditions,
		ImageUrl:             voucher.ImageUrl,
		VoidReason:           voucher.VoidReason,
		Note:                 voucher.Note,
		Status:               voucher.Status,
		StatusConvert:        statusx.ConvertStatusValue(voucher.Status),
		VoucherItem:          voucher.VoucherItem,
		CreatedAt:            voucher.CreatedAt,
		Region:               regionResponse,
		Archetype:            archetypeResponse,
		Customer:             customerResponse,
		MembershipLevel:      membershipLevelResponse,
		MembershipCheckpoint: membershipCheckpointResponse,
		Division:             divisionResponse,
		VoucherItems:         voucherItemsResponse,
	}

	return
}

func (s *VoucherService) Create(ctx context.Context, req *dto.VoucherRequestCreate) (res *dto.VoucherResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "VoucherService.Create")
	defer span.End()

	var (
		isExist                               bool
		voucher                               *model.Voucher
		customerResponse                      *dto.CustomerResponse
		regionResponse                        *dto.RegionResponse
		archetypeResponse                     *dto.ArchetypeResponse
		membershipLevelResponse               *dto.MembershipLevelResponse
		membershipCheckpointResponse          *dto.MembershipCheckpointResponse
		voucherItems                          []*dto.VoucherItemResponse
		divisionResponse                      *dto.DivisionResponse
		voucherItem                           int8
		archetype                             *bridge_service.GetArchetypeGPResponse
		customerType                          *bridgeService.GetCustomerTypeGPResponse
		customerTypeResponse                  *dto.CustomerTypeResponse
		statusArchetype, statusCustomerType   int8
		customerIDGP, expenseAccountAttribute string
	)

	// Validate for characters length
	if len(req.RedeemCode) < 5 || len(req.RedeemCode) > 20 {
		err = edenlabs.ErrorMustContain("redeem_code", "5 - 20 characters")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// validate the name cannot be the same with existing
	isExist = s.RepositoryVoucher.IsRedeemCodeExist(ctx, req.RedeemCode)
	if isExist {
		err = edenlabs.ErrorExists("redeem_code")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// Validate the redeem code is exist in GP
	voucherGP, err := s.opt.Client.BridgeServiceGrpc.GetVoucherGPList(ctx, &bridgeService.GetVoucherGPListRequest{
		Limit:            1,
		Offset:           0,
		GnlVoucherStatus: "0",
		GnlVoucherCode:   req.RedeemCode,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("voucherGP", "bridge")
		return
	}
	// validate the name cannot be the same with existing
	if len(voucherGP.Data) != 0 {
		err = edenlabs.ErrorExists("redeem_code")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// validate start_at not greater than time.now
	if req.StartTime.Before(time.Now()) {
		err = edenlabs.ErrorMustGreater("start_time", "time now")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// validate end time greater than start time and time now
	if req.StartTime.Equal(req.EndTime) || req.EndTime.Before(req.StartTime) {
		err = edenlabs.ErrorMustLater("end_time", "start_time")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// validate min order not less than equal 0
	if utils.ToInt(req.MinOrder) <= 0 {
		err = edenlabs.ErrorMustGreater("min_order", "0")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// validate disc amount not less than 1
	if req.DiscAmount < 1 {
		err = edenlabs.ErrorMustEqualOrGreater("disc_amount", "1")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// validate user quota not less than 1
	if req.UserQuota < 1 {
		err = edenlabs.ErrorMustEqualOrGreater("user_quota", "1")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// validate over all quota not less than 1
	if req.OverallQuota < 1 {
		err = edenlabs.ErrorMustEqualOrGreater("overall_quota", "1")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// validate user quota didn't greater that over all quota
	if req.OverallQuota < req.UserQuota {
		err = edenlabs.ErrorMustEqualOrGreater("overall_quota", "user_quota")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// validate type voucher
	if req.Type != 1 && req.Type != 2 && req.Type != 4 {
		err = edenlabs.ErrorInvalid("type")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// Manipulate data for all region
	if req.RegionID == "all" {
		req.RegionID = ""
		regionResponse = &dto.RegionResponse{
			Description: "All Region",
		}
	} else {
		// Get Region from bridge
		var region *bridgeService.GetAdmDivisionGPResponse
		region, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPDetail(ctx, &bridgeService.GetAdmDivisionGPDetailRequest{
			Region: req.RegionID,
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
	}

	// Get customer type from bridge
	customerType, err = s.opt.Client.BridgeServiceGrpc.GetCustomerTypeGPDetail(ctx, &bridge_service.GetCustomerTypeGPDetailRequest{
		Id: req.CustomerTypeID,
	})
	if err != nil || len(customerType.Data) == 0 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("customer_type_id")
		return
	}

	// Convert status GP to internal status
	if customerType.Data[0].Inactive == 0 {
		statusCustomerType = statusx.ConvertStatusName(statusx.Active)
	} else {
		statusCustomerType = statusx.ConvertStatusName(statusx.Archived)
	}
	customerTypeResponse = &dto.CustomerTypeResponse{
		ID:            customerType.Data[0].GnL_Cust_Type_ID,
		Code:          customerType.Data[0].GnL_Cust_Type_ID,
		Description:   customerType.Data[0].GnL_CustType_Description,
		CustomerGroup: customerType.Data[0].GnL_Cust_GroupDesc,
		Status:        statusCustomerType,
		ConvertStatus: statusx.ConvertStatusValue(statusCustomerType),
	}

	// Manipulate data for all archetype
	if req.ArchetypeID == "all" && req.CustomerTypeID != "" {
		req.ArchetypeID = ""
		archetypeResponse = &dto.ArchetypeResponse{
			Description:  "All Archetype",
			CustomerType: customerTypeResponse,
		}
	} else if req.CustomerTypeID != "" && req.ArchetypeID != "" {
		// Get Archetype from bridge
		archetype, err = s.opt.Client.BridgeServiceGrpc.GetArchetypeGPDetail(ctx, &bridge_service.GetArchetypeGPDetailRequest{
			Id: req.ArchetypeID,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("archetype_id")
			return
		}
		// Convert status GP to internal status
		if archetype.Data[0].Inactive == 0 {
			statusArchetype = statusx.ConvertStatusName(statusx.Active)
		} else {
			statusArchetype = statusx.ConvertStatusName(statusx.Archived)
		}
		archetypeResponse = &dto.ArchetypeResponse{
			ID:             archetype.Data[0].GnlArchetypeId,
			Code:           archetype.Data[0].GnlArchetypeId,
			Description:    archetype.Data[0].GnlArchetypedescription,
			CustomerTypeID: archetype.Data[0].GnlCustTypeId,
			CustomerType:   customerTypeResponse,
			Status:         statusArchetype,
			ConvertStatus:  statusx.ConvertStatusValue(statusArchetype),
		}
		// Validate archetype is syncronized with customer type
		if (req.CustomerTypeID != archetype.Data[0].GnlCustTypeId) && req.ArchetypeID != "all" {
			err = edenlabs.ErrorInvalid("archetype_id")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	}

	if req.DivisionID != 0 {
		// Get division from account service
		var division *account_service.GetDivisionDetailResponse
		division, err = s.opt.Client.AccountServiceGrpc.GetDivisionDetail(ctx, &account_service.GetDivisionDetailRequest{
			Id: req.DivisionID,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("division_id")
			return
		}

		divisionResponse = &dto.DivisionResponse{
			ID:            division.Data.Id,
			Code:          division.Data.Code,
			Description:   division.Data.Name,
			Status:        int8(division.Data.Status),
			StatusConvert: statusx.ConvertStatusValue(int8(division.Data.Status)),
		}

		// Set attribute COA Config based on division
		switch req.DivisionID {
		case 5:
			expenseAccountAttribute = "expense_account_coa_culinary"
		case 10:
			expenseAccountAttribute = "expense_account_coa_partnership"
		case 17:
			expenseAccountAttribute = "expense_account_coa_wholesale"
		default:
			expenseAccountAttribute = "expense_account_coa_others"
		}
	} else {
		expenseAccountAttribute = "expense_account_coa_others"
	}

	// Get COA number
	configCOA, err := s.opt.Client.ConfigurationServiceGrpc.GetConfigAppDetail(ctx, &configuration_service.GetConfigAppDetailRequest{
		Attribute: expenseAccountAttribute,
	})

	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if req.CustomerID != 0 {
		// Get customer from crm
		var customer *crmService.GetCustomerDetailResponse
		customer, err = s.opt.Client.CrmServiceGrpc.GetCustomerDetail(ctx, &crmService.GetCustomerDetailRequest{
			Id: req.CustomerID,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("customer_id")
			return
		}

		var customergp *bridgeService.GetCustomerGPResponse
		customergp, err = s.opt.Client.BridgeServiceGrpc.GetCustomerGPList(ctx, &bridgeService.GetCustomerGPListRequest{
			Id:     customer.Data.CustomerIdGp,
			Limit:  1,
			Offset: 0,
		})
		if err != nil || len(customergp.Data) == 0 {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("customer_id_gp")
			return
		}
		customerResponse = &dto.CustomerResponse{
			ID:   customer.Data.Id,
			Code: customer.Data.CustomerIdGp,
			Name: customergp.Data[0].Custname,
		}

		customerIDGP = customer.Data.CustomerIdGp
	}

	// Validate membership checkpoint is according with membership level
	if req.MembershipLevelID != 0 && req.MembershipCheckPointID == 0 {
		err = edenlabs.ErrorValidation("membership_checkpoint_id", "Please select membership lapak")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// Get membership level from campaign service
	if req.MembershipLevelID != 0 {
		var membershipLevel *campaignService.GetMembershipLevelDetailResponse
		membershipLevel, err = s.opt.Client.CampaignServiceGrpc.GetMembershipLevelDetail(ctx, &campaignService.GetMembershipLevelDetailRequest{
			Id: req.MembershipLevelID,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("campaign", "membership_level_id")
			return
		}
		membershipLevelResponse = &dto.MembershipLevelResponse{
			ID:       membershipLevel.Data.Id,
			Code:     membershipLevel.Data.Code,
			Level:    int8(membershipLevel.Data.Level),
			Name:     membershipLevel.Data.Name,
			ImageUrl: membershipLevel.Data.ImageUrl,
			Status:   int8(membershipLevel.Data.Status),
		}

		// Get membership checkpoint from campaign service
		var membershipCheckpoint *campaignService.GetMembershipCheckpointDetailResponse
		membershipCheckpoint, err = s.opt.Client.CampaignServiceGrpc.GetMembershipCheckpointDetail(ctx, &campaignService.GetMembershipCheckpointDetailRequest{
			Id: req.MembershipCheckPointID,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("campaign", "membership_checkpoint_id")
			return
		}

		// Validate membership checkpoint is according with membership level
		if membershipCheckpoint.Data.MembershipLevelId != req.MembershipLevelID {
			err = edenlabs.ErrorInvalid("membership_checkpoint_id")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		membershipCheckpointResponse = &dto.MembershipCheckpointResponse{
			ID:                membershipCheckpoint.Data.Id,
			Checkpoint:        int8(membershipCheckpoint.Data.Checkpoint),
			TargetAmount:      membershipCheckpoint.Data.TargetAmount,
			Status:            int8(membershipCheckpoint.Data.Status),
			MembershipLevelID: membershipCheckpoint.Data.MembershipLevelId,
		}
	}

	// Validate voucher item
	if len(req.VoucherItem) > 0 {
		voucherItem = 1
		productList := make(map[int64]string)
		for i, v := range req.VoucherItem {
			if v.ItemID == 0 {
				err = edenlabs.ErrorValidation(fmt.Sprintf("item_id%d", i), "Please Enter Product")
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}

			if v.MinQtyDisc < 0.01 {
				err = edenlabs.ErrorValidation(fmt.Sprintf("min_qty_disc%d", i), "Minimal Qty must be greater than 0.")
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}

			if _, exist := productList[v.ItemID]; exist {
				err = edenlabs.ErrorValidation(fmt.Sprintf("item_id%d", i), "Product is duplicate. Please enter another product.")
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			} else {
				productList[v.ItemID] = "t"
			}

			var item *catalogService.GetItemDetailByInternalIdResponse
			item, err = s.opt.Client.CatalogServiceGrpc.GetItemDetailByInternalId(ctx, &catalogService.GetItemDetailByInternalIdRequest{
				Id: utils.ToString(v.ItemID),
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorInvalid("item_id")
				return
			}

			detailItem := &dto.ItemResponse{
				ID:                   item.Data.Id,
				Code:                 item.Data.Code,
				Description:          item.Data.Description,
				UnitWeightConversion: item.Data.UnitWeightConversion,
			}

			voucherItems = append(voucherItems, &dto.VoucherItemResponse{
				ItemID:     v.ItemID,
				MinQtyDisc: v.MinQtyDisc,
				Item:       detailItem,
			})

		}
	} else {
		voucherItem = 2
	}

	// Generate Code For Voucher
	var codeGenerator *configurationService.GetGenerateCodeResponse
	codeGenerator, err = s.opt.Client.ConfigurationServiceGrpc.GetGenerateCode(ctx, &configurationService.GetGenerateCodeRequest{
		Format: "VOU",
		Domain: "voucher",
		Length: 6,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("configuration", "code generator")
		return
	}

	voucher = &model.Voucher{
		RegionIDGP:             req.RegionID,
		CustomerTypeIDGP:       req.CustomerTypeID,
		ArchetypeIDGP:          req.ArchetypeID,
		CustomerID:             req.CustomerID,
		MembershipLevelID:      req.MembershipLevelID,
		MembershipCheckPointID: req.MembershipCheckPointID,
		DivisionID:             req.DivisionID,
		Code:                   codeGenerator.Data.Code,
		RedeemCode:             req.RedeemCode,
		Name:                   req.Name,
		Type:                   req.Type,
		StartTime:              req.StartTime,
		EndTime:                req.EndTime,
		OverallQuota:           req.OverallQuota,
		RemOverallQuota:        req.OverallQuota,
		UserQuota:              req.UserQuota,
		MinOrder:               utils.ToFloat(req.MinOrder),
		DiscAmount:             req.DiscAmount,
		TermConditions:         req.TermConditions,
		ImageUrl:               req.ImageUrl,
		VoidReason:             req.VoidReason,
		Note:                   req.Note,
		Status:                 statusx.ConvertStatusName("Active"),
		VoucherItem:            voucherItem,
		CreatedAt:              time.Now(),
	}

	// Create new voucher to internal DB
	err = s.RepositoryVoucher.Create(ctx, voucher)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// Record item to voucher item
	for _, v := range voucherItems {
		item := &model.VoucherItem{
			VoucherID:  voucher.ID,
			ItemID:     v.ItemID,
			MinQtyDisc: v.MinQtyDisc,
			CreatedAt:  time.Now(),
		}
		err = s.RepositoryVoucherItem.Create(ctx, item)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		v.VoucherID = voucher.ID
		v.CreatedAt = time.Now()
	}

	// Create Voucher to GP
	_, err = s.opt.Client.BridgeServiceGrpc.CreateVoucherGP(ctx, &bridgeService.CreateVoucherGPRequest{
		GnlVoucherId:      voucher.Code,
		GnlChannel:        2,
		GnlVoucherType:    int32(voucher.Type),
		GnlVoucherName:    voucher.Name,
		GnlExpenseAccount: configCOA.Data.Value,
		GnlVoucherCode:    voucher.RedeemCode,
		GnlMinimumOrder:   int32(voucher.MinOrder),
		GnlDiscountAmount: int32(voucher.DiscAmount),
		GnlVoucherStatus:  0,
		Restriction: &bridgeService.CreateVoucherGPRequest_Restriction{
			GnlRegion:      voucher.RegionIDGP,
			GnlCustTypeId:  voucher.CustomerTypeIDGP,
			GnlArchetypeId: voucher.ArchetypeIDGP,
		},
		AdvancedProperties: &bridgeService.CreateVoucherGPRequest_AdvancedProperties{
			Custnmbr:              customerIDGP,
			GnlStartPeriod:        voucher.StartTime.Format("2006-01-02"),
			GnlEndPeriod:          voucher.EndTime.Format("2006-01-02"),
			GnlTotalQuotaCount:    int32(voucher.OverallQuota),
			GnlRemainingOverallQu: int32(voucher.RemOverallQuota),
			GnlTotalQuotaCountPe:  int32(voucher.UserQuota),
		},
	})

	if err != nil {
		err = edenlabs.ErrorValidation("voucherGP", "Error Create Voucher to GP")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	userID := ctx.Value(constants.KeyUserID).(int64)

	// Insert to audit log
	_, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &auditService.CreateLogRequest{
		Log: &auditService.Log{
			UserId:      userID,
			ReferenceId: strconv.Itoa(int(voucher.ID)),
			Type:        "voucher",
			Function:    "create",
			CreatedAt:   timestamppb.New(time.Now()),
			Note:        voucher.Note,
		},
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &dto.VoucherResponse{
		ID:                   voucher.ID,
		RedeemCode:           voucher.RedeemCode,
		Name:                 voucher.Name,
		Type:                 voucher.Type,
		StartTime:            voucher.StartTime,
		EndTime:              voucher.EndTime,
		OverallQuota:         voucher.OverallQuota,
		UserQuota:            voucher.UserQuota,
		RemOverallQuota:      voucher.RemOverallQuota,
		MinOrder:             voucher.MinOrder,
		DiscAmount:           voucher.DiscAmount,
		TermConditions:       voucher.TermConditions,
		ImageUrl:             voucher.ImageUrl,
		VoidReason:           voucher.VoidReason,
		Note:                 voucher.Note,
		Status:               voucher.Status,
		StatusConvert:        statusx.ConvertStatusValue(voucher.Status),
		VoucherItem:          voucher.VoucherItem,
		CreatedAt:            voucher.CreatedAt,
		Region:               regionResponse,
		Archetype:            archetypeResponse,
		Customer:             customerResponse,
		MembershipLevel:      membershipLevelResponse,
		MembershipCheckpoint: membershipCheckpointResponse,
		VoucherItems:         voucherItems,
		Division:             divisionResponse,
	}

	return
}

func (s *VoucherService) Archive(ctx context.Context, id int64) (res *dto.VoucherResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "VoucherService.Archive")
	defer span.End()

	var voucher *model.Voucher
	voucher, err = s.RepositoryVoucher.GetDetail(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("id")
		return
	}

	// Validate status must be active
	if voucher.Status != statusx.ConvertStatusName("Active") {
		err = edenlabs.ErrorMustActive("status")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if voucher.RemOverallQuota == 0 {
		voucher.VoidReason = 2
	} else if voucher.EndTime.Before(time.Now()) {
		voucher.VoidReason = 1
	} else {
		voucher.VoidReason = 3
	}

	voucher.Status = statusx.ConvertStatusName("Archived")

	err = s.RepositoryVoucher.Update(ctx, voucher, "Status", "VoidReason")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	userID := ctx.Value(constants.KeyUserID).(int64)

	_, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &auditService.CreateLogRequest{
		Log: &auditService.Log{
			UserId:      userID,
			ReferenceId: strconv.Itoa(int(voucher.ID)),
			Type:        "voucher",
			Function:    "archive",
			CreatedAt:   timestamppb.New(time.Now()),
		},
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &dto.VoucherResponse{
		ID:              voucher.ID,
		RedeemCode:      voucher.RedeemCode,
		Name:            voucher.Name,
		Type:            voucher.Type,
		StartTime:       voucher.StartTime,
		EndTime:         voucher.EndTime,
		OverallQuota:    voucher.OverallQuota,
		UserQuota:       voucher.UserQuota,
		RemOverallQuota: voucher.RemOverallQuota,
		MinOrder:        voucher.MinOrder,
		DiscAmount:      voucher.DiscAmount,
		TermConditions:  voucher.TermConditions,
		ImageUrl:        voucher.ImageUrl,
		VoidReason:      voucher.VoidReason,
		Note:            voucher.Note,
		Status:          voucher.Status,
		StatusConvert:   statusx.ConvertStatusValue(voucher.Status),
		VoucherItem:     voucher.VoucherItem,
		CreatedAt:       voucher.CreatedAt,
	}

	return
}

func (s *VoucherService) CreateBulky(ctx context.Context, req *dto.VoucherRequestBulky) (err error) {
	ctx, span := s.opt.Trace.Start(ctx, "VoucherService.CreateBulky")
	defer span.End()

	layout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Asia/Jakarta")
	redeemCodeList := make(map[string]int)
	errorList := ""
	for i, v := range req.Data {

		var expenseAccountAttribute string

		// Validate for characters length
		if len(v.RedeemCode) < 5 || len(v.RedeemCode) > 20 {
			errorList += fmt.Sprintf("Redeem Code must contains 5 and 20 characters at row %d.| ", i+1)
		}

		// Validate voucher name is required
		if v.VoucherName == "" {
			errorList += fmt.Sprintf("Voucher Name is required at row %d.| ", i+1)
		}

		// Validate voucher type is required
		if v.VoucherType == 0 {
			errorList += fmt.Sprintf("Voucher Type is required at row %d.| ", i+1)
		}

		// Validate voucher type value (exclude voucher type 3 redeem point )
		if v.VoucherType != 1 && v.VoucherType != 2 && v.VoucherType != 4 {
			errorList += fmt.Sprintf("Voucher Type is invalid at row %d.| ", i+1)
		}

		// Validate redeem code is required and exist
		if v.RedeemCode != "" {
			// validate the name cannot be the same with existing
			if _, exist := redeemCodeList[v.RedeemCode]; exist {
				errorList += fmt.Sprintf("Redeem Code is duplicate at row %d with data on row %d.| ", i+1, redeemCodeList[v.RedeemCode])
			} else {
				redeemCodeList[v.RedeemCode] = i + 1
			}

			var isExist bool
			isExist = s.RepositoryVoucher.IsRedeemCodeExist(ctx, v.RedeemCode)
			if isExist {
				errorList += fmt.Sprintf("Redeem Code is exist at row %d.| ", i+1)
			}
			var voucherGP *bridgeService.GetVoucherGPResponse
			voucherGP, err = s.opt.Client.BridgeServiceGrpc.GetVoucherGPList(ctx, &bridgeService.GetVoucherGPListRequest{
				Limit:            1,
				Offset:           0,
				GnlVoucherStatus: "0",
				GnlVoucherCode:   v.RedeemCode,
			})

			if len(voucherGP.Data) != 0 {
				errorList += fmt.Sprintf("Redeem Code is exist at row %d.| ", i+1)
			}
		} else {
			errorList += fmt.Sprintf("Redeem Code is required at row %d.| ", i+1)
		}

		// Validate start time format and Required
		if v.StartTime != "" {
			if v.StartTimeActual, err = time.ParseInLocation(layout, v.StartTime, loc); err != nil {
				errorList += fmt.Sprintf("Start time is invalid at row %d.| ", i+1)
			} else {
				v.StartTime = v.StartTimeActual.Format(time.RFC3339)
				v.StartTimeActual, err = time.ParseInLocation(time.RFC3339, v.StartTime, loc)
			}
		} else {
			errorList += fmt.Sprintf("Start time is required at row %d.| ", i+1)
		}

		// Validate End time format and Required
		if v.EndTime != "" {
			if v.EndTimeActual, err = time.ParseInLocation(layout, v.EndTime, loc); err != nil {
				errorList += fmt.Sprintf("End time is invalid at row %d.| ", i+1)
			} else {
				v.EndTime = v.EndTimeActual.Format(time.RFC3339)
				v.EndTimeActual, err = time.ParseInLocation(time.RFC3339, v.EndTime, loc)
			}
		} else {
			errorList += fmt.Sprintf("End time is required at row %d.| ", i+1)
		}

		// Validate start time more than now
		if v.StartTimeActual.Before(time.Now()) {
			errorList += fmt.Sprintf("Start time must greater than time now at row %d.| ", i+1)
		}

		// Validate end time must greater than start time
		if v.StartTimeActual.Equal(v.EndTimeActual) || v.EndTimeActual.Before(v.StartTimeActual) {
			errorList += fmt.Sprintf("End time must greater than start time at row %d.| ", i+1)
		}

		// validate min order not less than equal 0
		if v.MinOrder <= 0 {
			errorList += fmt.Sprintf("Min Order must equal or greater than 0 at row %d.| ", i+1)
		}

		// validate disc amount more than 0
		if v.DiscAmount < 1 {
			errorList += fmt.Sprintf("Disc Amount must greater than 0 at row %d.| ", i+1)
		}

		// validate user quota more than 0
		if v.UserQuota < 1 {
			errorList += fmt.Sprintf("User Quota must greater than 0 at row %d.| ", i+1)
		}

		// validate overall quota more than 0
		if v.OverallQuota < 1 {
			errorList += fmt.Sprintf("Overall Quota must greater than 0 at row %d.| ", i+1)
		}

		// validate user quota cannot more than overall quota
		if v.OverallQuota < v.UserQuota {
			errorList += fmt.Sprintf("Overall Quota must equal or greater than user quota at row %d.| ", i+1)
		}

		// Get Region from bridge
		if v.RegionName != "" {
			v.RegionName = strings.Title(v.RegionName)
			// check if it's for all regions
			if v.RegionName == "All" {
				v.RegionName = ""
			} else {
				var region *bridgeService.GetAdmDivisionGPResponse
				region, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPDetail(ctx, &bridgeService.GetAdmDivisionGPDetailRequest{
					Region: v.RegionName,
					Limit:  1,
					Offset: 0,
				})
				if err != nil {
					errorList += fmt.Sprintf("Region Code is invalid at row %d.| ", i+1)
				} else {
					v.RegionID = utils.ToString(region.Data[0].Region)
				}
			}
		} else {
			errorList += fmt.Sprintf("Region Code is required at row %d.| ", i+1)
		}

		// Get Customer Type from bridge
		if v.CustomerTypeCode != "" {
			// Get Archetype from bridge
			var CustomerType *bridgeService.GetCustomerTypeGPResponse
			CustomerType, err = s.opt.Client.BridgeServiceGrpc.GetCustomerTypeGPDetail(ctx, &bridgeService.GetCustomerTypeGPDetailRequest{
				Id: v.CustomerTypeCode,
			})
			if err != nil {
				errorList += fmt.Sprintf("Customer Type Code is invalid at row %d.| ", i+1)
			} else {
				v.CustomerTypeCode = CustomerType.Data[0].GnL_Cust_Type_ID
			}

		} else {
			errorList += fmt.Sprintf("Customer Type Code is required at row %d.| ", i+1)
		}

		// Get Archetype from bridge
		if v.ArchetypeCode != "" {
			v.ArchetypeCode = strings.Title(v.ArchetypeCode)
			// check if it's for all archetypes
			if v.ArchetypeCode == "All" {
				v.ArchetypeCode = ""
			} else {
				// Get Archetype from bridge
				var archetype *bridgeService.GetArchetypeGPResponse
				archetype, err = s.opt.Client.BridgeServiceGrpc.GetArchetypeGPDetail(ctx, &bridgeService.GetArchetypeGPDetailRequest{
					Id: v.ArchetypeCode,
				})
				if err != nil {
					errorList += fmt.Sprintf("Archetype Code is invalid at row %d.| ", i+1)
				} else {
					v.ArchetypeID = archetype.Data[0].GnlArchetypeId
					// Check the customer type and archetype is syncron
					if v.CustomerTypeCode != "" && (archetype.Data[0].GnlCustTypeId != v.CustomerTypeCode) {
						errorList += fmt.Sprintf("Archetype and customer type not match at row %d.| ", i+1)
					}
				}
			}
		} else {
			errorList += fmt.Sprintf("Archetype Code is required at row %d.| ", i+1)
		}

		// Get division from bridge
		var division *account_service.GetDivisionDefaultByCustomerTypeResponse
		division, err = s.opt.Client.AccountServiceGrpc.GetDivisionDefaultByCustomerType(ctx, &account_service.GetDivisionDefaultByCustomerTypeRequest{
			CustomerTypeIdGp: v.CustomerTypeCode,
		})
		if err != nil {
			errorList += fmt.Sprintf("Customer Type Code is invalid at row %d.| ", i+1)
		} else {
			v.DivisionID = division.Data.Id
			// Set attribute COA Config based on division
			switch division.Data.Id {
			case 5:
				expenseAccountAttribute = "expense_account_coa_culinary"
			case 10:
				expenseAccountAttribute = "expense_account_coa_partnership"
			case 17:
				expenseAccountAttribute = "expense_account_coa_wholesale"
			default:
				expenseAccountAttribute = "expense_account_coa_others"
			}

			// Get COA number
			var configCOA *configurationService.GetConfigAppDetailResponse
			configCOA, err = s.opt.Client.ConfigurationServiceGrpc.GetConfigAppDetail(ctx, &configuration_service.GetConfigAppDetailRequest{
				Attribute: expenseAccountAttribute,
			})

			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}

			v.ExpenseAccount = configCOA.Data.Value
		}

		// Get Customer from bridge
		if v.CustomerCode != "" {
			var customer *crmService.GetCustomerDetailResponse
			customer, err = s.opt.Client.CrmServiceGrpc.GetCustomerDetail(ctx, &crmService.GetCustomerDetailRequest{
				CustomerIdGp: v.CustomerCode,
			})
			if err != nil {
				errorList += fmt.Sprintf("Customer Code is invalid at row %d.| ", i+1)
			} else {
				v.CustomerID = customer.Data.Id
			}
		}

		// Get membership level from campaign service
		if v.MembershipLevel != 0 {
			var membershipLevel *campaignService.GetMembershipLevelDetailResponse
			membershipLevel, err = s.opt.Client.CampaignServiceGrpc.GetMembershipLevelDetail(ctx, &campaignService.GetMembershipLevelDetailRequest{
				Level: int64(v.MembershipLevel),
			})
			if err != nil {
				errorList += fmt.Sprintf("Membership Level is invalid at row %d.| ", i+1)
			} else {
				v.MembershipLevelID = membershipLevel.Data.Id
			}
			// Get membership level from campaign service
			var membershipCheckpoint *campaignService.GetMembershipCheckpointDetailResponse
			membershipCheckpoint, err = s.opt.Client.CampaignServiceGrpc.GetMembershipCheckpointDetail(ctx, &campaignService.GetMembershipCheckpointDetailRequest{
				Checkpoint: int64(v.MembershipCheckpoint),
			})
			if err != nil {
				errorList += fmt.Sprintf("Membership Checkpoint is invalid at row %d.| ", i+1)
			} else {
				if membershipCheckpoint.Data.MembershipLevelId != v.MembershipLevelID {
					errorList += fmt.Sprintf("Membership Checkpoint is invalid at row %d.| ", i+1)
				}
				v.MembershipCheckpointID = membershipCheckpoint.Data.Id
			}
		}
	}

	if errorList != "" {
		errorList = strings.TrimSuffix(errorList, "| ")
		err = edenlabs.ErrorValidation("error_callback", errorList)
		return
	}

	for _, v := range req.Data {
		var codeGenerator *configurationService.GetGenerateCodeResponse
		codeGenerator, err = s.opt.Client.ConfigurationServiceGrpc.GetGenerateCode(ctx, &configurationService.GetGenerateCodeRequest{
			Format: "VOU",
			Domain: "voucher",
			Length: 6,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("configuration", "code generator")
			return
		}

		voucher := &model.Voucher{
			RegionIDGP:             v.RegionID,
			CustomerTypeIDGP:       v.CustomerTypeCode,
			ArchetypeIDGP:          v.ArchetypeID,
			CustomerID:             v.CustomerID,
			MembershipLevelID:      v.MembershipLevelID,
			MembershipCheckPointID: v.MembershipCheckpointID,
			DivisionID:             v.DivisionID,
			Code:                   codeGenerator.Data.Code,
			RedeemCode:             v.RedeemCode,
			Name:                   v.VoucherName,
			Type:                   v.VoucherType,
			StartTime:              v.StartTimeActual,
			EndTime:                v.EndTimeActual,
			OverallQuota:           v.OverallQuota,
			RemOverallQuota:        v.OverallQuota,
			UserQuota:              v.UserQuota,
			MinOrder:               v.MinOrder,
			DiscAmount:             v.DiscAmount,
			Status:                 statusx.ConvertStatusName("Active"),
			VoucherItem:            2,
			CreatedAt:              time.Now(),
			Note:                   v.Note,
		}

		err = s.RepositoryVoucher.Create(ctx, voucher)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		// Create Voucher to GP
		_, err = s.opt.Client.BridgeServiceGrpc.CreateVoucherGP(ctx, &bridgeService.CreateVoucherGPRequest{
			GnlVoucherId:      voucher.Code,
			GnlChannel:        2,
			GnlVoucherType:    int32(voucher.Type),
			GnlVoucherName:    voucher.Name,
			GnlExpenseAccount: v.ExpenseAccount,
			GnlVoucherCode:    voucher.RedeemCode,
			GnlMinimumOrder:   int32(voucher.MinOrder),
			GnlDiscountAmount: int32(voucher.DiscAmount),
			GnlVoucherStatus:  0,
			Restriction: &bridgeService.CreateVoucherGPRequest_Restriction{
				GnlRegion:      voucher.RegionIDGP,
				GnlCustTypeId:  voucher.CustomerTypeIDGP,
				GnlArchetypeId: voucher.ArchetypeIDGP,
			},
			AdvancedProperties: &bridgeService.CreateVoucherGPRequest_AdvancedProperties{
				Custnmbr:              v.CustomerCode,
				GnlStartPeriod:        voucher.StartTime.Format("2006-01-02"),
				GnlEndPeriod:          voucher.EndTime.Format("2006-01-02"),
				GnlTotalQuotaCount:    int32(voucher.OverallQuota),
				GnlRemainingOverallQu: int32(voucher.RemOverallQuota),
				GnlTotalQuotaCountPe:  int32(voucher.UserQuota),
			},
		})
		if err != nil {
			err = edenlabs.ErrorValidation("voucherGP", "Error Create Voucher to GP")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		userID := ctx.Value(constants.KeyUserID).(int64)

		_, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &auditService.CreateLogRequest{
			Log: &auditService.Log{
				UserId:      userID,
				ReferenceId: strconv.Itoa(int(voucher.ID)),
				Type:        "voucher",
				Function:    "create bulky",
				CreatedAt:   timestamppb.New(time.Now()),
				Note:        voucher.Note,
			},
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	}

	return
}

func (s *VoucherService) GetMobileVoucherList(ctx context.Context, req *dto.VoucherRequestGetMobileVoucherList) (res []*dto.VoucherResponse, count int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "VoucherService.GetDetail")
	defer span.End()

	var vouchers []*model.Voucher
	vouchers, count, err = s.RepositoryVoucher.GetMobileVoucherList(ctx, req)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, voucher := range vouchers {
		res = append(res, &dto.VoucherResponse{
			ID:                   voucher.ID,
			Code:                 voucher.Code,
			RedeemCode:           voucher.RedeemCode,
			Name:                 voucher.Name,
			Type:                 voucher.Type,
			StartTime:            voucher.StartTime,
			EndTime:              voucher.EndTime,
			OverallQuota:         voucher.OverallQuota,
			UserQuota:            voucher.UserQuota,
			RemOverallQuota:      voucher.RemOverallQuota,
			MinOrder:             voucher.MinOrder,
			DiscAmount:           voucher.DiscAmount,
			TermConditions:       voucher.TermConditions,
			ImageUrl:             voucher.ImageUrl,
			VoidReason:           voucher.VoidReason,
			Note:                 voucher.Note,
			Status:               voucher.Status,
			StatusConvert:        statusx.ConvertStatusValue(voucher.Status),
			VoucherItem:          voucher.VoucherItem,
			CreatedAt:            voucher.CreatedAt,
			RemUserQuota:         voucher.RemUserQuota,
			MembershipLevel:      &dto.MembershipLevelResponse{ID: voucher.MembershipLevelID},
			MembershipCheckpoint: &dto.MembershipCheckpointResponse{ID: voucher.MembershipCheckPointID},
			Region:               &dto.RegionResponse{ID: voucher.RegionIDGP},
			Archetype:            &dto.ArchetypeResponse{ID: voucher.ArchetypeIDGP},
			Customer:             &dto.CustomerResponse{ID: voucher.CustomerID},
		})
	}

	return
}

func (s *VoucherService) GetMobileVoucherDetail(ctx context.Context, req *dto.VoucherRequestGetMobileVoucherDetail) (res *dto.VoucherResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "VoucherService.GetDetail")
	defer span.End()

	var voucher *model.Voucher
	voucher, err = s.RepositoryVoucher.GetMobileVoucherDetail(ctx, req)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &dto.VoucherResponse{
		ID:                   voucher.ID,
		Code:                 voucher.Code,
		RedeemCode:           voucher.RedeemCode,
		Name:                 voucher.Name,
		Type:                 voucher.Type,
		StartTime:            voucher.StartTime,
		EndTime:              voucher.EndTime,
		OverallQuota:         voucher.OverallQuota,
		UserQuota:            voucher.UserQuota,
		RemOverallQuota:      voucher.RemOverallQuota,
		MinOrder:             voucher.MinOrder,
		DiscAmount:           voucher.DiscAmount,
		TermConditions:       voucher.TermConditions,
		ImageUrl:             voucher.ImageUrl,
		VoidReason:           voucher.VoidReason,
		Note:                 voucher.Note,
		Status:               voucher.Status,
		StatusConvert:        statusx.ConvertStatusValue(voucher.Status),
		VoucherItem:          voucher.VoucherItem,
		CreatedAt:            voucher.CreatedAt,
		RemUserQuota:         voucher.RemUserQuota,
		MembershipLevel:      &dto.MembershipLevelResponse{ID: voucher.MembershipLevelID},
		MembershipCheckpoint: &dto.MembershipCheckpointResponse{ID: voucher.MembershipCheckPointID},
		Region:               &dto.RegionResponse{ID: voucher.RegionIDGP},
		Archetype:            &dto.ArchetypeResponse{ID: voucher.ArchetypeIDGP, CustomerType: &dto.CustomerTypeResponse{ID: voucher.CustomerTypeIDGP}},
		Customer:             &dto.CustomerResponse{ID: voucher.CustomerID},
	}
	return
}

func (s *VoucherService) Update(ctx context.Context, req *dto.VoucherRequestUpdate) (err error) {
	ctx, span := s.opt.Trace.Start(ctx, "VoucherService.Archive")
	defer span.End()

	voucher := &model.Voucher{
		ID:              req.VoucherID,
		RemOverallQuota: req.RemOverallQuota,
	}

	err = s.RepositoryVoucher.Update(ctx, voucher, "RemOverallQuota")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}
