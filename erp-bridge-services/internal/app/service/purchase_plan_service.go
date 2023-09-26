package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/repository"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

type IPurchasePlanService interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (res []dto.PurchasePlanResponse, total int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (res dto.PurchasePlanResponse, err error)
	GetGP(ctx context.Context, req *pb.GetPurchasePlanGPListRequest) (res *pb.GetPurchasePlanGPResponse, err error)
	GetDetailGP(ctx context.Context, req *pb.GetPurchasePlanGPDetailRequest) (res *pb.GetPurchasePlanGPResponse, err error)
	AssignPurchasePlanGP(ctx context.Context, req *dto.AssignPurchasePlanGPRequest) (res *dto.AssignPurchasePlanGPResponse, err error)
}

type PurchasePlanService struct {
	opt                    opt.Options
	RepositoryPurchasePlan repository.IPurchasePlanRepository
}

func NewPurchasePlanService() IPurchasePlanService {
	return &PurchasePlanService{
		opt:                    global.Setup.Common,
		RepositoryPurchasePlan: repository.NewPurchasePlanRepository(),
	}
}

func (s *PurchasePlanService) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (res []dto.PurchasePlanResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PurchasePlanService.Get")
	defer span.End()

	var purchasePlans []*model.PurchasePlan
	purchasePlans, total, err = s.RepositoryPurchasePlan.Get(ctx, offset, limit, status, search, orderBy)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, purchasePlan := range purchasePlans {
		res = append(res, dto.PurchasePlanResponse{
			ID:                   purchasePlan.ID,
			Code:                 purchasePlan.Code,
			VendorOrganizationID: purchasePlan.VendorOrganizationID,
			SiteID:               purchasePlan.SiteID,
			RecognitionDate:      purchasePlan.RecognitionDate,
			EtaDate:              purchasePlan.EtaDate,
			EtaTime:              purchasePlan.EtaTime,
			TotalPrice:           purchasePlan.TotalPrice,
			TotalWeight:          purchasePlan.TotalWeight,
			TotalPurchasePlanQty: purchasePlan.TotalPurchasePlanQty,
			TotalPurchaseQty:     purchasePlan.TotalPurchaseQty,
			Note:                 purchasePlan.Note,
			Status:               purchasePlan.Status,
			AssignedTo:           purchasePlan.AssignedTo,
			AssignedBy:           purchasePlan.AssignedBy,
			AssignedAt:           purchasePlan.AssignedAt,
			CreatedAt:            purchasePlan.CreatedAt,
			CreatedBy:            purchasePlan.CreatedBy,
		})
	}

	return
}

func (s *PurchasePlanService) GetDetail(ctx context.Context, id int64, code string) (res dto.PurchasePlanResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PurchasePlanService.GetDetail")
	defer span.End()

	var purchasePlan *model.PurchasePlan
	purchasePlan, err = s.RepositoryPurchasePlan.GetDetail(ctx, id, code)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.PurchasePlanResponse{
		ID:                   purchasePlan.ID,
		Code:                 purchasePlan.Code,
		VendorOrganizationID: purchasePlan.VendorOrganizationID,
		SiteID:               purchasePlan.SiteID,
		RecognitionDate:      purchasePlan.RecognitionDate,
		EtaDate:              purchasePlan.EtaDate,
		EtaTime:              purchasePlan.EtaTime,
		TotalPrice:           purchasePlan.TotalPrice,
		TotalWeight:          purchasePlan.TotalWeight,
		TotalPurchasePlanQty: purchasePlan.TotalPurchasePlanQty,
		TotalPurchaseQty:     purchasePlan.TotalPurchaseQty,
		Note:                 purchasePlan.Note,
		Status:               purchasePlan.Status,
		AssignedTo:           purchasePlan.AssignedTo,
		AssignedBy:           purchasePlan.AssignedBy,
		AssignedAt:           purchasePlan.AssignedAt,
		CreatedAt:            purchasePlan.CreatedAt,
		CreatedBy:            purchasePlan.CreatedBy,
	}

	return
}

func (s *PurchasePlanService) GetGP(ctx context.Context, req *pb.GetPurchasePlanGPListRequest) (res *pb.GetPurchasePlanGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PurchasePlanService.GetGP")
	defer span.End()

	params := map[string]string{
		"interid":    global.EnvDatabaseGP,
		"PageNumber": strconv.Itoa(int(req.Offset)),
		"PageSize":   strconv.Itoa(int(req.Limit)),
	}

	if req.PrpPurchaseplanNo != "" {
		params["prp_purchaseplan_no"] = req.PrpPurchaseplanNo
	}
	if req.Locncode != "" {
		params["locncode"] = req.Locncode
	}
	if req.PrpPurchaseplanDateFrom != "" {
		params["prp_purchaseplan_date_from"] = req.PrpPurchaseplanDateFrom
	}

	if req.PrpPurchaseplanDateTo != "" {
		params["prp_purchaseplan_date_to"] = req.PrpPurchaseplanDateTo
	}

	if req.Status != 0 {
		params["status"] = strconv.Itoa(int(req.Status))
	}

	if req.FieldPurchaser != "" {
		params["prp_purchaseplan_user"] = req.FieldPurchaser
	}

	// filter LIKE (both) for search
	if req.PrpVendorOrgDesc != "" {
		params["prp_purchaseplan_no"] = req.PrpPurchaseplanNo
		params["prp_vendor_org_desc"] = req.PrpVendorOrgDesc
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "PurchasePlanList/list", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *PurchasePlanService) GetDetailGP(ctx context.Context, req *pb.GetPurchasePlanGPDetailRequest) (res *pb.GetPurchasePlanGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PurchasePlanService.GetDetailGP")
	defer span.End()

	params := map[string]string{
		"interid": global.EnvDatabaseGP,
		"id":      req.Id,
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "purchaseplanlist/detail", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *PurchasePlanService) AssignPurchasePlanGP(ctx context.Context, req *dto.AssignPurchasePlanGPRequest) (res *dto.AssignPurchasePlanGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PurchaseOrderService.Update")
	defer span.End()

	req.Interid = global.EnvDatabaseGP
	fmt.Println("aku")

	err = global.HttpRestApiToMicrosoftGP("POST", "PurchasePlanList/assign", req, &res, nil)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = errors.New("Connection to the server could not be established")
		return
	}

	if res.Code != 200 {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = errors.New("Connection to the server could not be established")
		return
	}

	return
}
