package service

import (
	"fmt"
	"strconv"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/validation"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/crm_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/settlement_service"
	"github.com/labstack/echo/v4"

	"git.edenfarm.id/edenlabs/edenlabs/jwt"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/repository"
)

func ServicePersonal() IPersonalService {
	m := new(PersonalService)
	m.opt = global.Setup.Common
	m.RepositoryUserCustomer = repository.NewUserCustomerRepository()
	m.RepositoryWhiteListLogin = repository.NewWhiteListLoginRepository()
	m.RepositoryOTPOutgoing = repository.NewOtpOutgoingRepository()
	return m
}

type IPersonalService interface {
	SaveRegistration(ctx echo.Context, req *dto.SaveRegistrationRequest) (token string, e error)
}

type PersonalService struct {
	opt                      opt.Options
	RepositoryUserCustomer   repository.IUserCustomerRepository
	RepositoryWhiteListLogin repository.IWhiteListLoginRepository
	RepositoryOTPOutgoing    repository.IOtpOutgoingRepository
}

func NewPersonalService() IPersonalService {
	return &PersonalService{
		opt:                      global.Setup.Common,
		RepositoryUserCustomer:   repository.NewUserCustomerRepository(),
		RepositoryWhiteListLogin: repository.NewWhiteListLoginRepository(),
		RepositoryOTPOutgoing:    repository.NewOtpOutgoingRepository(),
	}
}

func (s *PersonalService) SaveRegistration(ctx echo.Context, req *dto.SaveRegistrationRequest) (token string, e error) {
	c := ctx.Request().Context()
	c, span := s.opt.Trace.Start(c, "PersonalService.DeleteAccount")
	var err error
	var otp string
	fmt.Print(span)
	OTPValidDuration := s.opt.Env.GetString("OTP.VALID_DURATION") //env.GetString("OTP_VALID_DURATION", "600")

	appVersionInt, err := strconv.ParseInt(req.AppVersion, 10, 64)

	if (req.Platform == "orca") && (appVersionInt < 4001086 || req.OTP == "") {
		// o.Failure("message", "Update aplikasimu, minimum versi 4.1.86")
		// return o
	}

	if (req.Platform != "orca") && (appVersionInt < 2001128 || req.OTP == "") {
		// o.Failure("message", "Update aplikasimu, minimum versi 2.1.128")
		// return o
	}

	_, err = s.RepositoryWhiteListLogin.GetDetail(ctx.Request().Context(), 0, req.CustomerPhoneNumber, req.OTP)
	if err != nil {
		db := s.opt.Database.Read
		if err = db.Raw("SELECT otp FROM otp_outgoing oo WHERE oo.phone_number = ? "+
			"AND oo.created_at BETWEEN NOW()-INTERVAL ? second "+
			"AND NOW() "+
			"AND oo.otp_status = 1 "+
			"AND oo.usage_type = 2 ORDER BY oo.created_at DESC LIMIT 1", req.CustomerPhoneNumber, OTPValidDuration).QueryRow(&otp); err != nil {
			// o.Failure("message", "Kode OTP sudah tidak berlaku. Silakan coba kirim ulang.")
			// return o
		}
		if req.OTP == "" {
			// o.Failure("message", "Kode OTP sudah tidak berlaku. Silakan coba kirim ulang.")
			// return o
		}
		if req.OTP != otp {
			// o.Failure("message", "Kode OTP tidak valid. Silakan masukkan kode OTP yang benar.")
			// return o
		}
	}

	// if req.CodeUserCustomer, err = util.CheckTable("user_merchant"); err != nil {
	// 	// o.Failure("code.invalid", util.ErrorInvalidData("code"))
	// }

	// if c.CodeMerchant, err = util.CheckTable("merchant"); err != nil {
	// 	// o.Failure("code.invalid", util.ErrorInvalidData("code"))
	// }

	if req.CustomerName != "" {
		if !validation.CharacterOnly(req.CustomerName) {
			// o.Failure("customer_name", util.ErrorValidNameFormat())
		}
	} else {
		// o.Failure("customer_name", "Nama harus diisi.")
	}
	// check number only
	if !validation.NumberOnly(req.CustomerPhoneNumber) {
		// o.Failure("merchant_phone_number", util.ErrorValidPhoneFormat())
	}
	if !validation.NumberOnly(req.CustomerAltPhoneNumber) {
		// o.Failure("merchant_alt_phone_number", util.ErrorValidPhoneFormat())
	}
	//check first phone number
	if string(req.CustomerPhoneNumber[0]) == "0" {
		// o.Failure("merchant_phone_number", "Nomor telepon tidak valid.")
	}
	if len(req.CustomerPhoneNumber) < 8 || len(req.CustomerPhoneNumber) > 12 {
		// o.Failure("merchant_phone_number", "Nomor telepon harus 8 sampai 12 karakter.")
	}
	if string(req.CustomerAltPhoneNumber[0]) == "0" {
		// o.Failure("merchant_alt_phone_number", "Nomor telepon tidak valid.")
	}
	if len(req.CustomerAltPhoneNumber) < 8 || len(req.CustomerAltPhoneNumber) > 12 {
		// o.Failure("merchant_alt_phone_number", "Nomor telepon harus 8 sampai 12 karakter.")
	}
	// check valid email
	if req.CustomerEmail != "" {
		if !validation.EmailOnly(req.CustomerEmail) {
			// o.Failure("merchant_email", util.ErrorValidEmailFormat())
		}
	} else {
		// o.Failure("merchant_email", "Email harus diisi.")
	}
	if req.CustomerPhoneNumber == "" {
		// o.Failure("phone_number.required", util.ErrorInputRequired("phone number"))
	}

	if req.CustomerBirthDate != "" {
		if req.BirthDateAt, err = time.Parse("2006-01-02", req.CustomerBirthDate); err != nil {
			// o.Failure("birth_date.invalid", "Invalid birth_date")
		}
	}
	if req.ReferenceInfo == "" {
		// o.Failure("reference_info.required", util.ErrorInputRequired("reference info"))
	}

	//is phone number have a merchant
	// if customer, err := s.opt.Client.BridgeServiceGrpc.GetCustomerDetail(ctx.Request().Context(), &bridge_service.GetCustomerDetailRequest{
	// 	PhoneNumber: req.CustomerPhoneNumber,
	// }); err == nil {
	// 	//if len data lebih dari 0{}
	// 	if customer.Data.Id == 0 {
	// 		// o.Failure("merchant_phone_number", util.ErrorRegisteredPhoneNumber())
	// 	}
	// }

	if customer, err := s.opt.Client.BridgeServiceGrpc.GetCustomerGPList(ctx.Request().Context(), &bridge_service.GetCustomerGPListRequest{
		Phone:  req.CustomerPhoneNumber,
		Limit:  1,
		Offset: 0,
	}); err == nil {
		//if len data lebih dari 0{}
		if len(customer.Data) > 0 {
			err = edenlabs.ErrorValidation("merchant_phone_number", "Nomor ponsel sudah terdaftar. Silakan masukan nomor lain yang belum terdaftar.")
			// o.Failure("merchant_phone_number", util.ErrorRegisteredPhoneNumber())
			e = err
			return
		}
	}

	//check email
	if email, err := s.opt.Client.CrmServiceGrpc.GetCustomerDetail(ctx.Request().Context(), &crm_service.GetCustomerDetailRequest{
		Email: req.CustomerEmail,
	}); err == nil {
		//if len data lebih dari 0{}
		if email.Data.Id != 0 {
			err = edenlabs.ErrorValidation("merchant_email", "Email sudah terdaftar. Silakan masukan email yang belum terdaftar.")
			// o.Failure("merchant_phone_number", util.ErrorRegisteredPhoneNumber())
			e = err
			return
		}
	}

	//check referral_code
	if req.ReferrerCode != "" {
		refCode, err := s.opt.Client.CrmServiceGrpc.GetCustomerDetail(ctx.Request().Context(), &crm_service.GetCustomerDetailRequest{
			ReferrerCode: req.ReferrerCode,
		})
		if err != nil {
			//span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorValidation("referral_code", "Code referral tidak ditemukan.")
			e = err
			return
		}
		if refCode.Data.Id == 0 {
			err = edenlabs.ErrorValidation("referral_code", "Code referral tidak ditemukan.")
			e = err
			return
		}
		req.ReferrerID = refCode.Data.Id
	}

	if req.SubDistrictID == "" {
		// o.Failure("sub_district.invalid", util.ErrorInvalidData("sub district"))
	} else {
		// idSubdistrict, _ := common.Decrypt(c.SubDistrictID)
		// c.SubDistrict = &model.SubDistrict{ID: idSubdistrict}
		// c.SubDistrict.Read("ID")
		// c.SubDistrict.Area.Read("ID")

		// c.WarehouseCoverage = &model.WarehouseCoverage{SubDistrict: c.SubDistrict, MainWarehouse: 1}

		// c.PriceSet = &model.PriceSet{ID: c.SubDistrict.Area.ID}

		// if err = c.PriceSet.Read("ID"); err != nil {
		// 	o.Failure("area", util.ErrorInvalidData("area"))
		// }

		// if err = c.WarehouseCoverage.Read("SubDistrict", "MainWarehouse"); err != nil {
		// 	o.Failure("warehousecoverage", util.ErrorInvalidData("Warehouse Coverage"))
		// } else {
		// 	if err = c.WarehouseCoverage.Warehouse.Read("ID"); err != nil {
		// 		o.Failure("warehouse.invalid", util.ErrorMustBeSame("warehouse sub district", "sub district"))
		// 	}
		// }

		// if c.CodeBranch, err = util.CheckTable("branch"); err != nil {
		// 	o.Failure("code.invalid", util.ErrorInvalidData("code"))
		// }
	}

	// //read sales person config app to assign default sales person on branch
	// salesPersonConfig := &model.ConfigApp{Attribute: "salesperson_self_registration"}
	// salesPersonConfig.Read("Attribute")
	// salesPersonValue, _ := strconv.ParseInt(salesPersonConfig.Value, 10, 64)
	// c.SalesPerson = &model.Staff{ID: salesPersonValue}
	// c.SalesPerson.Read("ID")

	// // Set coordinate validation
	// if c.Latitude != nil && c.Longitude != nil {
	// 	c.IsValidCoordinate = 1
	// }

	// req.CodeUserCustomer, _ = util.GenerateCode("", "user_merchant")
	var codeGenerator *configuration_service.GetGenerateCodeResponse
	codeGenerator, err = s.opt.Client.ConfigurationServiceGrpc.GetGenerateCode(ctx.Request().Context(), &configuration_service.GetGenerateCodeRequest{
		Format: "UCS",
		Domain: "user_customer",
		Length: 4,
	})
	req.CodeUserCustomer = codeGenerator.Data.Code
	um := &model.UserCustomer{
		Code:          req.CodeUserCustomer,
		Verification:  2,
		Status:        1,
		ForceLogout:   2,
		FirebaseToken: req.FirebaseToken,
	}

	// //default data for registration
	// req.InvoiceTerm = &model.InvoiceTerm{Code: "SIT0001"}
	// req.InvoiceTerm.Read("Code")
	// req.PaymentTerm = &model.SalesTerm{Code: "SPT0011"}
	// req.PaymentTerm.Read("Code")
	// req.PaymentMethod = &model.PaymentMethod{Code: "PYM0002"}
	// req.PaymentMethod.Read("Code")
	// req.BusinessType = &model.BusinessType{Code: "BTY0009"}
	// req.BusinessType.Read("Code")
	// req.PaymentGroup = &model.PaymentGroup{Code: "PYG0001"}
	// req.PaymentGroup.Read("Code")
	if e = s.RepositoryUserCustomer.Create(ctx.Request().Context(), um); e == nil {
		codeGenerator, err = s.opt.Client.ConfigurationServiceGrpc.GetGenerateCustomerCode(ctx.Request().Context(), &configuration_service.GetGenerateCodeRequest{
			Format: "EC",
			Domain: "customer",
			Length: 8,
		})
		req.CodeCustomer = codeGenerator.Data.Code
		codeGenerator, err = s.opt.Client.ConfigurationServiceGrpc.GetGenerateCode(ctx.Request().Context(), &configuration_service.GetGenerateCodeRequest{
			Format: req.CodeCustomer + "-",
			Domain: "address",
			Length: 3,
		})
		req.CodeAddress = codeGenerator.Data.Code
		codeGenerator, err = s.opt.Client.ConfigurationServiceGrpc.GetGenerateReferralCode(ctx.Request().Context(), &configuration_service.GetGenerateCodeRequest{
			Format: "CRF",
			Domain: "referral",
			Length: 8,
		})
		// req.CodeReferral = util.GenerateReferralCode(8)
		req.CodeReferral = codeGenerator.Data.Code

		bcaVA, err := s.opt.Client.SettlementGrpc.GenerateFixedVaXendit(ctx.Request().Context(), &settlement_service.GenerateFixedVaXenditRequest{
			ExternalId: "BCA_FVA-" + req.CodeCustomer,
			BankCode:   "BCA",
		})
		permataVA, err := s.opt.Client.SettlementGrpc.GenerateFixedVaXendit(ctx.Request().Context(), &settlement_service.GenerateFixedVaXenditRequest{
			ExternalId: "PERMATA_FVA-" + req.CodeCustomer,
			BankCode:   "PERMATA",
		})
		fmt.Print(bcaVA, permataVA, err)
		m := &bridge_service.CreateCustomerGPRequest{
			Custnmbr:     req.CodeCustomer,
			Custname:     req.CustomerName,
			Address1:     req.Address1,
			Address2:     req.Address2,
			Address3:     req.Address3,
			Adrscode:     req.CodeAddress,
			City:         "DUMMY CITY",
			Cntcprsn:     req.CustomerName,
			Cprcstnm:     "",
			Custpriority: "1",
			Phone1:       req.CustomerPhoneNumber,
			Phone2:       req.CustomerAltPhoneNumber,
			// Shrtname:        "Shirt Name",
			State: "JKT",
			// Stmtname:        "STMT Name",
			GnlReferrerCode: req.ReferrerCode,
			GnlReferralCode: req.CodeReferral,
			GnlCustTypeId:   "BTY0009",
			GnlBusinessType: "2",
			Userdef1:        bcaVA.Data.AccountNumber,
			Userdef2:        permataVA.Data.AccountNumber,
		}
		if _, e = s.opt.Client.BridgeServiceGrpc.CreateCustomerGP(ctx.Request().Context(), m); e == nil {
			// var priceSet *model.PriceSet

			// orm := orm.NewOrm()
			// orm.Using("read_only")
			// orm.Raw("SELECT ap.id," +
			// 	" ap.area_id," +
			// 	" ap.default_price_set" +
			// 	" FROM area_policy ap" +
			// 	" JOIN area a ON ap.area_id = a.id" +
			// 	" WHERE a.aux_data = 2").QueryRows(&req.PriceSetArea)

			// // transaction to merchant price set
			// // get default price set id for branch
			// for _, v := range req.PriceSetArea {
			// 	mps := &model.MerchantPriceSet{
			// 		PriceSet: v.DefaultPriceSet,
			// 		Area:     v.Area,
			// 		Merchant: m,
			// 	}
			// 	if _, e = o.Insert(mps); e != nil {
			// 		o.Rollback()
			// 	}

			// 	//area suitable with selected sub district
			// 	if v.Area.ID == req.SubDistrict.Area.ID {
			// 		priceSet = v.DefaultPriceSet
			// 	}
			// }

			admDivision, err := s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPList(ctx.Request().Context(), &bridge_service.GetAdmDivisionGPListRequest{
				Limit:       1,
				Offset:      0,
				SubDistrict: req.SubDistrictID,
			})
			if err != nil {
				err = edenlabs.ErrorRpcNotFound("bridge", "adm division")
				e = err
				return
			}
			if len(admDivision.Data) == 0 {
				err = edenlabs.ErrorValidation("adm_division", "sub district not found.")
				e = err
				return
			}
			b := &bridge_service.UpdateAddressRequest{
				Custnmbr: req.CodeCustomer,
				Adrscode: req.CodeAddress,
				Cntcprsn: m.Cntcprsn,
				AddresS1: req.Address1,
				AddresS2: req.Address2,
				AddresS3: req.Address3,
				City:     admDivision.Data[0].City,
				State:    admDivision.Data[0].State,
				Zip:      m.Zip,
				CCode:    "ID",
				// Country:          admDivision.Data[0].,
				GnL_Address_Note: req.AddressNote,
				// GnL_Longitude:    strconv.FormatFloat(float64(*req.Longitude), 'f', 1, 64),
				// GnL_Latitude:     strconv.FormatFloat(float64(*req.Latitude), 'f', 1, 64),
				ShipToName:              req.CustomerName,
				GnL_Archetype_ID:        "ARC0001",
				PhonE1:                  m.Phone1,
				PhonE2:                  m.Phone2,
				PhonE3:                  m.Phone3,
				GnL_Administrative_Code: admDivision.Data[0].Code,
				UserdeF1:                bcaVA.Data.AccountNumber,
				UserdeF2:                permataVA.Data.AccountNumber,

				// TypeAddress:             "ship_to",
			}
			if req.Latitude != nil {
				b.GnL_Latitude = strconv.FormatFloat(float64(*req.Latitude), 'f', 1, 64)
			}
			if req.Longitude != nil {
				b.GnL_Longitude = strconv.FormatFloat(float64(*req.Longitude), 'f', 1, 64)
			}

			if _, e = s.opt.Client.BridgeServiceGrpc.UpdateAddress(ctx.Request().Context(), b); e == nil {

				if _, e = s.opt.Client.BridgeServiceGrpc.SetDefaultAddress(ctx.Request().Context(), &bridge_service.SetDefaultAddressRequest{
					Adrscode: b.Adrscode,
					Custnmbr: b.Custnmbr,
				}); e == nil {

				}
				// o.Commit()
				// e = log.AuditLogByMerchant(m, m.ID, "agent", "create", "")
				// xendit.BCAXenditFixedVA(m)
				// xendit.PermataXenditFixedVA(m)
				jwtInit := jwt.NewJWT([]byte(s.opt.Config.Jwt.Key))
				uc := jwt.UserMobile{
					PhoneNo:   req.CustomerPhoneNumber,
					ExpiresAt: time.Now().Add(time.Hour * 4).Unix(),
					// Timezone:  req.Timezone,
					//StandardClaims: jwt.StandardClaims{},
				}

				jwtGenerate, err := jwtInit.Create(uc)
				if err != nil {
					s.opt.Logger.AddMessage(log.ErrorLevel, err)
					e = err
					return
				}
				reqCreateRequestCRM := &crm_service.CreateCustomerRequest{
					Data: &crm_service.Customer{
						CustomerIdGp:            m.Custnmbr,
						ProspectiveCustomerId:   0,
						MembershipLevelId:       0,
						MembershipCheckpointId:  0,
						TotalPoint:              0,
						ProfileCode:             m.Custnmbr,
						Email:                   req.CustomerEmail,
						ReferenceInfo:           req.ReferenceInfo,
						UpgradeStatus:           0,
						KtpPhotosUrl:            "",
						CustomerPhotosUrl:       "",
						CustomerSelfieUrl:       "",
						MembershipRewardId:      0,
						MembershipRewardAmmount: 0,
						ReferralCode:            req.CodeReferral,
						ReferrerId:              req.ReferrerID,
						ReferrerCode:            req.ReferrerCode,
						Gender:                  int32(req.CustomerGender),
						BirthDate:               req.CustomerBirthDate,
					},
				}
				fmt.Print(reqCreateRequestCRM)
				customerResponse, err := s.opt.Client.CrmServiceGrpc.CreateCustomer(ctx.Request().Context(), reqCreateRequestCRM)
				//jwtGenerate = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwaG9uZV9ubyI6IjgwMDAwMTExMTAxIiwiZXhwaXJlX2F0IjoxNjc4MTc0MTA2LCJ0aW1lem9uZSI6IkFzaWEvSmFrYXJ0YSJ9.4LD2NcAoCr2Itdvvm-HGz7IFIusmhFzEO-lV1SJy1hY"
				um.LastLoginAt = time.Now()
				um.LoginToken = jwtGenerate
				um.ForceLogout = 2
				um.Verification = 2
				// um.CustomerIDGP = m.Custnmbr
				um.CustomerID = customerResponse.Data.Id

				err = s.RepositoryUserCustomer.Update(ctx.Request().Context(), um)
				if err != nil {
					span.RecordError(err)
					s.opt.Logger.AddMessage(log.ErrorLevel, err)
					e = err
					return
				}
				token = um.LoginToken

			} else {
				// o.Rollback()
			}

		}
	} else {
		// o.Rollback()
	}

	return
}
