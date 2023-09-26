package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/timex"
	"git.edenfarm.id/edenlabs/edenlabs/utils"

	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/repository"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"github.com/sirupsen/logrus"
)

type ICustomerService interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, customerTypeId int64) (res []dto.CustomerResponse, total int64, err error)
	GetDetail(ctx context.Context, id int64, code string, phoneNumber string) (res dto.CustomerResponse, err error)
	UpdateCustomer(ctx context.Context, req *bridgeService.UpdateCustomerRequest) (res dto.CustomerResponse, err error)
	GetGP(ctx context.Context, req *pb.GetCustomerGPListRequest) (res *pb.GetCustomerGPResponse, err error)
	GetDetailGP(ctx context.Context, req *pb.GetCustomerGPDetailRequest) (res *pb.GetCustomerGPResponse, err error)
	CreateCustomerGP(ctx context.Context, req *bridgeService.CreateCustomerGPRequest) (res dto.CreateCustomerGPResponse, err error)
	UpdateCustomerGP(ctx context.Context, req *bridgeService.UpdateCustomerGPRequest) (res *bridgeService.GetCustomerGPResponse, err error)
	UpdateFixedVa(ctx context.Context, req *bridgeService.UpdateFixedVaRequest) (res *bridgeService.UpdateFixedVaResponse, err error)
}

type CustomerService struct {
	opt                opt.Options
	RepositoryCustomer repository.ICustomerRepository
}

func NewCustomerService() ICustomerService {
	return &CustomerService{
		opt:                global.Setup.Common,
		RepositoryCustomer: repository.NewCustomerRepository(),
	}
}

func (s *CustomerService) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, customerTypeId int64) (res []dto.CustomerResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CustomerService.Get")
	defer span.End()

	// customer
	var customers []*model.Customer
	customers, total, err = s.RepositoryCustomer.Get(ctx, offset, limit, status, search, orderBy, customerTypeId)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, cust := range customers {
		res = append(res, dto.CustomerResponse{
			ID:                         cust.ID,
			Code:                       cust.Code,
			Name:                       cust.Name,
			Gender:                     cust.Gender,
			BirthDate:                  cust.BirthDate,
			PicName:                    cust.PicName,
			PhoneNumber:                cust.PhoneNumber,
			AltPhoneNumber:             cust.AltPhoneNumber,
			Email:                      cust.Email,
			Password:                   cust.Password,
			BillingAddress:             cust.BillingAddress,
			Note:                       cust.Note,
			ReferenceInfo:              cust.ReferenceInfo,
			TagCustomer:                cust.TagCustomer,
			TagCustomerName:            cust.TagCustomerName,
			Status:                     cust.Status,
			Suspended:                  cust.Suspended,
			UpgradeStatus:              cust.UpgradeStatus,
			CustomerGroup:              cust.CustomerGroup,
			ReferralCode:               cust.ReferralCode,
			ReferrerCode:               cust.ReferrerCode,
			TotalPoint:                 cust.TotalPoint,
			CustomerTypeCreditLimit:    cust.CustomerTypeCreditLimit,
			EarnedPoint:                cust.EarnedPoint,
			RedeemedPoint:              cust.RedeemedPoint,
			CustomCreditLimit:          cust.CustomCreditLimit,
			CreditLimitAmount:          cust.CreditLimitAmount,
			RemainingCreditLimitAmount: cust.RemainingCreditLimitAmount,
			ProfileCode:                cust.ProfileCode,
			AverageSales:               cust.AverageSales,
			RemainingOutstanding:       cust.RemainingOutstanding,
			OverdueDebt:                cust.OverdueDebt,
			KTPPhotosUrl:               cust.KTPPhotosUrl,
			KTPPhotosUrlArr:            cust.KTPPhotosUrlArr,
			MerchantPhotosUrl:          cust.MerchantPhotosUrl,
			MerchantPhotosUrlArr:       cust.MerchantPhotosUrlArr,
			MembershipLevelID:          cust.MembershipLevelID,
			MembershipRewardID:         cust.MembershipRewardID,
			MembershipCheckpointID:     cust.MembershipCheckpointID,
			MembershipRewardAmount:     cust.MembershipRewardAmount,
			SalesPaymentTermID:         cust.SalesPaymentTermID,
			CreatedAt:                  timex.ToLocTime(ctx, cust.CreatedAt),
			CreatedBy:                  cust.CreatedBy,
			LastUpdatedAt:              cust.LastUpdatedAt,
			LastUpdatedBy:              cust.LastUpdatedBy,
			BirthDateString:            cust.BirthDateString,
			CustomerTypeId:             cust.CustomerTypeID,
		})
	}

	return
}

func (s *CustomerService) GetDetail(ctx context.Context, id int64, code string, phoneNumber string) (res dto.CustomerResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CustomerService.GetDetail")
	defer span.End()

	var customer *model.Customer
	customer, err = s.RepositoryCustomer.GetDetail(ctx, id, code, phoneNumber)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// customer response
	res = dto.CustomerResponse{
		ID:                         customer.ID,
		Code:                       customer.Code,
		Name:                       customer.Name,
		Gender:                     customer.Gender,
		BirthDate:                  customer.BirthDate,
		PicName:                    customer.PicName,
		PhoneNumber:                customer.PhoneNumber,
		AltPhoneNumber:             customer.AltPhoneNumber,
		Email:                      customer.Email,
		Password:                   customer.Password,
		BillingAddress:             customer.BillingAddress,
		Note:                       customer.Note,
		ReferenceInfo:              customer.ReferenceInfo,
		TagCustomer:                customer.TagCustomer,
		TagCustomerName:            customer.TagCustomerName,
		Status:                     customer.Status,
		Suspended:                  customer.Suspended,
		UpgradeStatus:              customer.UpgradeStatus,
		CustomerGroup:              customer.CustomerGroup,
		ReferralCode:               customer.ReferralCode,
		ReferrerCode:               customer.ReferrerCode,
		TotalPoint:                 customer.TotalPoint,
		CustomerTypeCreditLimit:    customer.CustomerTypeCreditLimit,
		EarnedPoint:                customer.EarnedPoint,
		RedeemedPoint:              customer.RedeemedPoint,
		CustomCreditLimit:          customer.CustomCreditLimit,
		CreditLimitAmount:          customer.CreditLimitAmount,
		RemainingCreditLimitAmount: customer.RemainingCreditLimitAmount,
		ProfileCode:                customer.ProfileCode,
		AverageSales:               customer.AverageSales,
		RemainingOutstanding:       customer.RemainingOutstanding,
		OverdueDebt:                customer.OverdueDebt,
		KTPPhotosUrl:               customer.KTPPhotosUrl,
		KTPPhotosUrlArr:            customer.KTPPhotosUrlArr,
		MerchantPhotosUrl:          customer.MerchantPhotosUrl,
		MerchantPhotosUrlArr:       customer.MerchantPhotosUrlArr,
		MembershipLevelID:          customer.MembershipLevelID,
		MembershipRewardID:         customer.MembershipRewardID,
		MembershipCheckpointID:     customer.MembershipCheckpointID,
		MembershipRewardAmount:     customer.MembershipRewardAmount,
		SalesPaymentTermID:         customer.SalesPaymentTermID,
		CreatedAt:                  timex.ToLocTime(ctx, customer.CreatedAt),
		CreatedBy:                  customer.CreatedBy,
		LastUpdatedAt:              customer.LastUpdatedAt,
		LastUpdatedBy:              customer.LastUpdatedBy,
		BirthDateString:            customer.BirthDateString,
		CustomerTypeId:             customer.CustomerTypeID,
	}

	return
}

func (s *CustomerService) UpdateCustomer(ctx context.Context, req *bridgeService.UpdateCustomerRequest) (res dto.CustomerResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CustomerService.UpdateCustomer")
	defer span.End()

	//var customer *model.Customer
	//send to gp to update customer
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	reqRest := &dto.CustomerResponse{
		ID:                         req.Data.Id,
		Code:                       req.Data.Code,
		Name:                       req.Data.Name,
		Gender:                     int8(req.Data.Gender),
		BirthDate:                  req.Data.BirthDate.AsTime(),
		PicName:                    req.Data.PicName,
		PhoneNumber:                req.Data.PhoneNumber,
		AltPhoneNumber:             req.Data.AltPhoneNumber,
		Email:                      req.Data.Email,
		Password:                   req.Data.Password,
		BillingAddress:             req.Data.BillingAddress,
		Note:                       req.Data.Note,
		ReferenceInfo:              req.Data.ReferenceInfo,
		TagCustomer:                req.Data.TagCustomer,
		TagCustomerName:            req.Data.TagCustomerName,
		Status:                     int8(req.Data.Status),
		Suspended:                  int8(req.Data.Suspended),
		UpgradeStatus:              int8(req.Data.UpgradeStatus),
		CustomerGroup:              int8(req.Data.CustomerGroup),
		ReferralCode:               req.Data.ReferralCode,
		ReferrerCode:               req.Data.ReferrerCode,
		TotalPoint:                 req.Data.TotalPoint,
		CustomerTypeCreditLimit:    int8(req.Data.CustomerTypeCreditLimit),
		EarnedPoint:                req.Data.EarnedPoint,
		RedeemedPoint:              req.Data.RedeemedPoint,
		CustomCreditLimit:          int8(req.Data.CustomCreditLimit),
		CreditLimitAmount:          req.Data.CreditLimitAmount,
		RemainingCreditLimitAmount: req.Data.RemainingCreditLimitAmount,
		ProfileCode:                req.Data.ProfileCode,
		AverageSales:               req.Data.AverageSales,
		RemainingOutstanding:       req.Data.RemainingOutstanding,
		OverdueDebt:                req.Data.OverdueDebt,
		KTPPhotosUrl:               req.Data.KTPPhotosUrl,
		KTPPhotosUrlArr:            req.Data.KTPPhotosUrlArr,
		MerchantPhotosUrl:          req.Data.MerchantPhotosUrl,
		MerchantPhotosUrlArr:       req.Data.MerchantPhotosUrlArr,
		MembershipLevelID:          req.Data.MembershipLevelID,
		MembershipRewardID:         req.Data.MembershipRewardID,
		MembershipCheckpointID:     req.Data.MembershipCheckpointID,
		MembershipRewardAmount:     req.Data.MembershipRewardAmount,
		//SalesPaymentTermID:         req.Data.SalesPaymentTermID,
		CreatedAt:       req.Data.CreatedAt.AsTime(),
		CreatedBy:       req.Data.CreatedBy,
		LastUpdatedAt:   req.Data.LastUpdatedAt.AsTime(),
		LastUpdatedBy:   req.Data.LastUpdatedBy,
		BirthDateString: req.Data.BirthDateString,
		CustomerTypeId:  req.Data.CustomerTypeId,
	}
	payload, _ := json.Marshal(reqRest)

	fmt.Println("payload customer : ", payload)
	// resp := &bridgeService.GetCustomerDetailResponse{}
	// err = global.HttpRestApiToMicrosoftGP("PUT", "Update/customer", reqRest, &resp)
	// if err != nil {
	// 	logrus.Error(err.Error())
	// 	os.Exit(1)
	// 	return
	// }

	// if resp.Code != 200 {
	// 	logrus.Error("Error Login: " + resp.Message)
	// 	os.Exit(1)
	// 	return
	// }
	//tes
	res = dto.CustomerResponse{
		ID:                         req.Data.Id,
		Code:                       req.Data.Code,
		Name:                       req.Data.Name,
		Gender:                     int8(req.Data.Gender),
		BirthDate:                  req.Data.BirthDate.AsTime(),
		PicName:                    req.Data.PicName,
		PhoneNumber:                req.Data.PhoneNumber,
		AltPhoneNumber:             req.Data.AltPhoneNumber,
		Email:                      req.Data.Email,
		Password:                   req.Data.Password,
		BillingAddress:             req.Data.BillingAddress,
		Note:                       req.Data.Note,
		ReferenceInfo:              req.Data.ReferenceInfo,
		TagCustomer:                req.Data.TagCustomer,
		TagCustomerName:            req.Data.TagCustomerName,
		Status:                     int8(req.Data.Status),
		Suspended:                  int8(req.Data.Suspended),
		UpgradeStatus:              int8(req.Data.UpgradeStatus),
		CustomerGroup:              int8(req.Data.CustomerGroup),
		ReferralCode:               req.Data.ReferralCode,
		ReferrerCode:               req.Data.ReferrerCode,
		TotalPoint:                 req.Data.TotalPoint,
		CustomerTypeCreditLimit:    int8(req.Data.CustomerTypeCreditLimit),
		EarnedPoint:                req.Data.EarnedPoint,
		RedeemedPoint:              req.Data.RedeemedPoint,
		CustomCreditLimit:          int8(req.Data.CustomCreditLimit),
		CreditLimitAmount:          req.Data.CreditLimitAmount,
		RemainingCreditLimitAmount: req.Data.RemainingCreditLimitAmount,
		ProfileCode:                req.Data.ProfileCode,
		AverageSales:               req.Data.AverageSales,
		RemainingOutstanding:       req.Data.RemainingOutstanding,
		OverdueDebt:                req.Data.OverdueDebt,
		KTPPhotosUrl:               req.Data.KTPPhotosUrl,
		KTPPhotosUrlArr:            req.Data.KTPPhotosUrlArr,
		MerchantPhotosUrl:          req.Data.MerchantPhotosUrl,
		MerchantPhotosUrlArr:       req.Data.MerchantPhotosUrlArr,
		MembershipLevelID:          req.Data.MembershipLevelID,
		MembershipRewardID:         req.Data.MembershipRewardID,
		MembershipCheckpointID:     req.Data.MembershipCheckpointID,
		MembershipRewardAmount:     req.Data.MembershipRewardAmount,
		//SalesPaymentTermID:         req.Data.SalesPaymentTermID,
		CreatedAt:       req.Data.CreatedAt.AsTime(),
		CreatedBy:       req.Data.CreatedBy,
		LastUpdatedAt:   req.Data.LastUpdatedAt.AsTime(),
		LastUpdatedBy:   req.Data.LastUpdatedBy,
		BirthDateString: req.Data.BirthDateString,
		CustomerTypeId:  req.Data.CustomerTypeId,
	}

	return
}

func (s *CustomerService) GetGP(ctx context.Context, req *pb.GetCustomerGPListRequest) (res *pb.GetCustomerGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CustomerService.GetGP")
	defer span.End()

	params := map[string]string{
		"interid":    global.EnvDatabaseGP,
		"PageNumber": strconv.Itoa(int(req.Offset)),
		"PageSize":   strconv.Itoa(int(req.Limit)),
	}
	if req.Phone != "" {
		params["phone1"] = req.Phone
	}
	if req.Id != "" {
		req.Id = url.PathEscape(req.Id)
		params["custnmbr"] = req.Id
	}
	if req.AdressId != "" {
		req.AdressId = url.PathEscape(req.AdressId)
		params["adrscode"] = req.AdressId
	}
	if req.Name != "" {
		req.Name = url.PathEscape(req.Name)
		params["custname"] = req.Name
	}
	if req.Prstadcd != "" {
		params["prstadcd"] = req.Prstadcd
	}
	if req.Salsterr != "" {
		req.Id = url.PathEscape(req.Id)
		params["salsterr"] = req.Salsterr
	}
	if req.SalesPersonId != "" {
		req.Id = url.PathEscape(req.Id)
		params["slprsnid"] = req.SalesPersonId
	}
	if req.CustomerTypeId != "" {
		req.Id = url.PathEscape(req.Id)
		params["gnl_cust_type_id"] = req.CustomerTypeId
	}
	if req.ReferrerCode != "" {
		params["gnl_referrer_code"] = req.ReferrerCode
	}
	if req.ReferralCode != "" {
		params["gnl_referral_code"] = req.ReferralCode
	}
	if req.Inactive != "" {
		params["inactive"] = req.Inactive
	}
	if req.Orderby != "" {
		params["orderby"] = req.Orderby
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "customer/list", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *CustomerService) GetDetailGP(ctx context.Context, req *pb.GetCustomerGPDetailRequest) (res *pb.GetCustomerGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CustomerService.GetDetailGP")
	defer span.End()

	req.Id = url.PathEscape(req.Id)
	params := map[string]string{
		"interid": global.EnvDatabaseGP,
		"id":      req.Id,
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "customer/detail", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if len(res.Data) == 0 {
		err = edenlabs.ErrorNotFound("customer")
	}

	return
}

func (s *CustomerService) CreateCustomerGP(ctx context.Context, req *bridgeService.CreateCustomerGPRequest) (res dto.CreateCustomerGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CustomerService.CreateCustomer")
	defer span.End()
	reqGP := &dto.CreateCustomerGPRequest{
		InterID:           global.EnvDatabaseGP,
		CustName:          req.Custname,
		CustClas:          req.Custclas,
		CustPriority:      req.Custpriority,
		CprCstNm:          req.Cprcstnm,
		StmtName:          req.Stmtname,
		ShrtName:          req.Shrtname,
		Address:           &dto.CreateOrUpdateAddressGpRequest{},
		UPSZone:           req.Upszone,
		ShipMthd:          req.Shipmthd,
		TaxSchID:          req.Taxschid,
		PrbtAdcd:          req.Prbtadcd,
		PrstAdcd:          req.Prstadcd,
		StAdrcd:           req.Staddrcd,
		SlprsnID:          req.Slprsnid,
		PymtrmID:          req.Pymtrmid,
		Inactive:          req.Inactive,
		Hold:              req.Hold,
		Salsterr:          req.Salsterr,
		UserDef1:          req.Userdef1,
		UserDef2:          req.Userdef2,
		DeclID:            req.Declid,
		Comment1:          req.Comment1,
		Comment2:          req.Comment2,
		CustDisc:          req.Custdisc,
		DisGrper:          req.Disgrper,
		DueGrper:          req.Duegrper,
		PrcLevel:          req.Prclevel,
		GnlCustTypeID:     req.GnlCustTypeId,
		GnlReferrerCode:   req.GnlReferrerCode,
		GnlBusinessType:   req.GnlBusinessType,
		GnlSocialSecNum:   req.GnlSocialSecNum,
		ShipComplete:      1, // Default to direct invoice
		CreditLimitType:   req.Crlmttyp,
		CreditLimitDesc:   req.CrlmttypDesc,
		CreditLimitAmount: req.Crlmtamt,
	}
	req.Interid = global.EnvDatabaseGP

	err = global.HttpRestApiToMicrosoftGP("POST", "customer/create", reqGP, &res, nil)
	if err != nil {
		logrus.Error(err.Error())
		err = errors.New("Connection to the server could not be established")
		return
	}

	if res.Code != 200 {
		logrus.Error("Error Login: " + res.Message)
		err = errors.New("Connection to the server could not be established")
		return
	}

	return
}

func (s *CustomerService) UpdateCustomerGP(ctx context.Context, req *bridgeService.UpdateCustomerGPRequest) (res *bridgeService.GetCustomerGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CustomerService.UpdateCustomerGP")
	defer span.End()

	reqGP := &dto.UpdateCustomerGPRequest{
		InterID:      global.EnvDatabaseGP,
		CustNmbr:     req.Custnmbr,
		CustName:     req.Custname,
		CustClas:     req.Custclas,
		CustPriority: utils.ToString(req.Custpriority),
		CprCstNm:     req.Cprcstnm,
		StmtName:     req.Stmtname,
		ShrtName:     req.Shrtname,
		Address: &dto.CreateOrUpdateAddressGpRequest{
			AdrsCode: req.Address.Adrscode,
			CntcPrsn: req.Address.Cntcprsn,
			Address1: req.Address.AddresS1,
			Address2: req.Address.AddresS2,
			Address3: req.Address.AddresS3,
			City:     req.Address.City,
			State:    req.Address.State,
			Zip:      req.Address.Zip,
			CCode:    req.Address.CCode,
			Country:  req.Address.Country,
			Phone1:   req.Address.PhonE1,
			Phone2:   req.Address.PhonE2,
			Phone3:   req.Address.PhonE3,
			Fax:      req.Address.Fax,
		},
		UPSZone:         req.Upszone,
		ShipMthd:        req.Shipmthd,
		TaxSchID:        req.Taxschid,
		PrbtAdcd:        req.Prbtadcd,
		PrstAdcd:        req.Prstadcd,
		StAdrcd:         req.Staddrcd,
		SlprsnID:        req.Slprsnid,
		PymtrmID:        req.Pymtrmid,
		Inactive:        req.Inactive,
		Hold:            req.Hold,
		Salsterr:        req.Salsterr,
		UserDef1:        req.Userdef1,
		UserDef2:        req.Userdef2,
		DeclID:          req.Declid,
		Comment1:        req.Comment1,
		Comment2:        req.Comment2,
		CustDisc:        req.Custdisc,
		DisGrper:        req.Disgrper,
		DueGrper:        req.Duegrper,
		PrcLevel:        req.Prclevel,
		GnlCustTypeID:   req.GnlCustTypeId,
		GnlReferrerCode: req.GnlReferrerCode,
		GnlBusinessType: req.GnlBusinessType,
		GnlSocialSecNum: req.GnlSocialSecNum,
		ShipComplete:    req.Shipcomplete,
	}

	err = global.HttpRestApiToMicrosoftGP("PUT", "customer/update", reqGP, &res, nil)
	if err != nil {
		logrus.Error(err.Error())
		err = errors.New("Connection to the server could not be established")
		return
	}

	return
}

func (s *CustomerService) UpdateFixedVa(ctx context.Context, req *bridgeService.UpdateFixedVaRequest) (res *bridgeService.UpdateFixedVaResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CustomerService.UpdateFixedVa")
	defer span.End()

	reqGP := &dto.UpdateFixedVaRequest{
		InterID:  global.EnvDatabaseGP,
		CustNmbr: req.CustomerIdGp,
		UserDef1: req.FixedVaBca,
		UserDef2: req.FixedVaPermata,
	}

	err = global.HttpRestApiToMicrosoftGP("PUT", "customer/update", reqGP, &res, nil)

	if err != nil {
		logrus.Error(err.Error())
		err = errors.New("Connection to the server could not be established")
		return
	}

	return
}
