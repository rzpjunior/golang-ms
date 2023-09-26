package service

import (
	"context"
	"strings"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/utils"
	"git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/global"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	dto "git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/account_service"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
)

func NewServiceSalesperson() ISalespersonService {
	m := new(SalespersonService)
	m.opt = global.Setup.Common
	return m
}

type ISalespersonService interface {
	Get(ctx context.Context, req dto.SalespersonListRequest) (res []*dto.SalespersonResponse, err error)
	GetDetailById(ctx context.Context, req dto.SalespersonDetailRequest) (res *dto.SalespersonResponse, err error)
	GetGP(ctx context.Context, req dto.SalespersonListRequest) (res []*dto.SalesPerson, total int64, err error)
	GetDetaiGPlById(ctx context.Context, id string) (res *dto.SalesPerson, err error)
}

type SalespersonService struct {
	opt opt.Options
}

func (s *SalespersonService) Get(ctx context.Context, req dto.SalespersonListRequest) (res []*dto.SalespersonResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalespersonService.Get")
	defer span.End()

	// get sites from bridge
	var spRes *bridgeService.GetSalespersonListResponse
	spRes, err = s.opt.Client.BridgeServiceGrpc.GetSalespersonList(ctx, &bridgeService.GetSalespersonListRequest{
		Limit:   req.Limit,
		Offset:  req.Offset,
		Status:  req.Status,
		Search:  req.Search,
		OrderBy: req.OrderBy,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "salesperson")
		return
	}

	datas := []*dto.SalespersonResponse{}
	for _, sp := range spRes.Data {
		datas = append(datas, &dto.SalespersonResponse{
			ID:            sp.Id,
			Code:          sp.Code,
			FirstName:     sp.Firstname,
			MiddleName:    sp.Middlename,
			LastName:      sp.Lastname,
			Status:        int8(sp.Status),
			StatusConvert: statusx.ConvertStatusValue(int8(sp.Status)),
			CreatedAt:     sp.CreatedAt.AsTime(),
			UpdatedAt:     sp.UpdatedAt.AsTime(),
		})
	}
	res = datas

	return
}

func (s *SalespersonService) GetDetailById(ctx context.Context, req dto.SalespersonDetailRequest) (res *dto.SalespersonResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalespersonService.GetDetailById")
	defer span.End()

	// get Salesperson from bridge
	var spRes *bridgeService.GetSalespersonDetailResponse
	spRes, err = s.opt.Client.BridgeServiceGrpc.GetSalespersonDetail(ctx, &bridgeService.GetSalespersonDetailRequest{
		Id: int64(req.Id),
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "salesperson")
		return
	}

	res = &dto.SalespersonResponse{
		ID:            spRes.Data.Id,
		Code:          spRes.Data.Code,
		FirstName:     spRes.Data.Firstname,
		MiddleName:    spRes.Data.Middlename,
		LastName:      spRes.Data.Lastname,
		Status:        int8(spRes.Data.Status),
		StatusConvert: statusx.ConvertStatusValue(int8(spRes.Data.Status)),
		CreatedAt:     spRes.Data.CreatedAt.AsTime(),
		UpdatedAt:     spRes.Data.UpdatedAt.AsTime(),
	}

	return
}

func (s *SalespersonService) GetGP(ctx context.Context, req dto.SalespersonListRequest) (res []*dto.SalesPerson, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalespersonService.GetGP")
	defer span.End()

	var status int32

	switch statusx.ConvertStatusValue(int8(req.Status)) {
	// status convert
	case statusx.Active:
		status = 1
	case statusx.Archived:
		status = 2
	}

	var arrSalesPerson []int64
	// arrSalesPerson = []int64{6, 8, 11}
	// fmt.Println(arrSalesPerson)
	adminLapak, _ := s.opt.Client.ConfigurationServiceGrpc.GetConfigAppList(ctx, &configuration_service.GetConfigAppListRequest{
		Offset:    0,
		Limit:     1,
		Attribute: "admin_edn_role_id",
		// Field:     "Admin Lapak EDN",
	})
	substrings := strings.Split(adminLapak.Data[0].Value, ",")

	for _, substring := range substrings {
		arrSalesPerson = append(arrSalesPerson, utils.ToInt64(substring))
	}

	var spAcc *account_service.GetUserListResponse
	spAcc, err = s.opt.Client.AccountServiceGrpc.GetUserList(ctx, &account_service.GetUserListRequest{
		Limit:  req.Limit,
		Offset: req.OffsetQuery,
		Search: req.Search,
		Status: status,
		// DivisionId: 8,
		ArrRoleId: arrSalesPerson,
	})
	// fmt.Println(spAcc)
	// get salesperson from bridge
	// var spRes *bridgeService.GetSalesPersonGPResponse
	// spRes, err = s.opt.Client.BridgeServiceGrpc.GetSalesPersonGPList(ctx, &bridgeService.GetSalesPersonGPListRequest{
	// 	Limit:  req.Limit,
	// 	Offset: req.Offset,
	// 	Search: req.Search,
	// })
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("account", "salesperson")
		return
	}

	datas := []*dto.SalesPerson{}
	for _, sp := range spAcc.Data {
		switch sp.Status {
		// status convert
		case 1:
			status = 1
		case 2:
			status = 7
		}

		datas = append(datas, &dto.SalesPerson{
			StaffID:      sp.EmployeeCode,
			EmployeeCode: sp.EmployeeCode,
			Name:         sp.Name,
			DisplayName:  sp.Nickname,
			Status:       int8(status),
		})
	}

	total = int64(len(datas))
	res = datas

	return
}

func (s *SalespersonService) GetDetaiGPlById(ctx context.Context, id string) (res *dto.SalesPerson, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalespersonService.GetDetaiGPlById")
	defer span.End()

	var sp *account_service.GetUserDetailResponse
	sp, err = s.opt.Client.AccountServiceGrpc.GetUserDetail(ctx, &account_service.GetUserDetailRequest{
		// Id: utils.toint64(id),
		EmployeeCode: id,
	})

	// // get sales person from bridge
	// var sp *bridgeService.GetSalesPersonGPResponse
	// sp, err = s.opt.Client.BridgeServiceGrpc.GetSalesPersonGPDetail(ctx, &bridgeService.GetSalesPersonGPDetailRequest{
	// 	Id: id,
	// })
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("account", "salesperson")
		return
	}

	res = &dto.SalesPerson{
		StaffID:      sp.Data.EmployeeCode,
		EmployeeCode: sp.Data.EmployeeCode,
		Name:         sp.Data.Name,
		DisplayName:  sp.Data.Nickname,
	}

	return
}
