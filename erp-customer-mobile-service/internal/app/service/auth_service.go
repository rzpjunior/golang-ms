package service

import (
	"bytes"
	"context"
	rand2 "crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/validation"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/campaign_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/crm_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/sales_service"
	"github.com/labstack/echo/v4"

	"git.edenfarm.id/edenlabs/edenlabs/jwt"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/repository"
)

func ServiceAuth() IAuthService {
	m := new(AuthService)
	m.opt = global.Setup.Common
	m.RepositoryUserCustomer = repository.NewUserCustomerRepository()
	m.RepositoryWhiteListLogin = repository.NewWhiteListLoginRepository()
	m.RepositoryOTPOutgoing = repository.NewOtpOutgoingRepository()
	return m
}

type IAuthService interface {
	Login(ctx context.Context, req dto.LoginRequest) (res dto.LoginResponse, err error)
	ResponseOTP(ctx context.Context, req dto.RequestGetOTP) (res dto.ResponseSendOtp, err error)
	SignOut(ctx echo.Context, customerId int64) (m *dto.SessionDataCustomer, e error)
	CheckPhoneNumber(ctx echo.Context, req *dto.CheckPhoneNumberRequest) (m *dto.CheckPhoneNumberRequest, e error)
	SendOTPWASociomile(ctx context.Context, phoneNumber string, application int, usageType int) error
	SendOTPSMSVIRO(ctx context.Context, phoneNumber string, application int, usageType int) error
	Session(ctx echo.Context) (m *dto.SessionDataCustomer, e error)
	VerifyOtp(ctx echo.Context, req dto.VerifyRegistRequest) error
	DeleteAccount(ctx echo.Context, customerID int64, customerIdGP string) (e error)
}

type AuthService struct {
	opt                      opt.Options
	RepositoryUserCustomer   repository.IUserCustomerRepository
	RepositoryWhiteListLogin repository.IWhiteListLoginRepository
	RepositoryOTPOutgoing    repository.IOtpOutgoingRepository
}

func NewAuthService() IAuthService {
	return &AuthService{
		opt:                      global.Setup.Common,
		RepositoryUserCustomer:   repository.NewUserCustomerRepository(),
		RepositoryWhiteListLogin: repository.NewWhiteListLoginRepository(),
		RepositoryOTPOutgoing:    repository.NewOtpOutgoingRepository(),
	}
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (res dto.LoginResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AuthService.Login")

	//var e error
	var otp string
	OTPValidDuration := s.opt.Env.GetString("OTP.VALID_DURATION") //env.GetString("OTP_VALID_DURATION", "600")
	//to get expired period to limit how many seconds value from database
	appConfig, err := s.opt.Client.ConfigurationServiceGrpc.GetConfigAppDetail(ctx, &configuration_service.GetConfigAppDetailRequest{
		Attribute: "cma_failed_login_retry_time",
	})
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("configuration", "config app")
		return
	}
	period, _ := strconv.Atoi(appConfig.Data.Value)
	expPeriod := time.Duration(period) * time.Second

	fmt.Println(">>>>>>>>>>>>>>>>>>", expPeriod)
	// //to get retry count value from database
	appConfig, err = s.opt.Client.ConfigurationServiceGrpc.GetConfigAppDetail(ctx, &configuration_service.GetConfigAppDetailRequest{
		Attribute: "cma_failed_login_retry_count",
	})
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("configuration", "config app")
		return
	}
	maxRetry, _ := strconv.Atoi(appConfig.Data.Value)

	fmt.Println(">>>>>>>>>>>>>>>>>>", maxRetry)

	redisKey := req.Data.PhoneNumber
	// check number only
	fmt.Println(">>>>>>>>>>>>>>>>>>", redisKey)

	if validation.NumberOnly(req.Data.PhoneNumber) == false {
		// o.Failure("phone_number", "Nomor telepon harus dalam format angka.")
		err = edenlabs.ErrorValidation("phone_number", "Nomor telepon harus dalam format angka.")
		// err = errors.New("Nomor telepon harus dalam format angka.")
	}
	// //check first char
	if string(req.Data.PhoneNumber[0]) == "0" {
		err = edenlabs.ErrorValidation("phone_number", "Nomor telepon tidak valid.")
		//o.Failure("phone_number", "")
		// err = errors.New("Nomor telepon tidak valid.")
	}
	if len(req.Data.PhoneNumber) < 8 || len(req.Data.PhoneNumber) > 12 {
		err = edenlabs.ErrorValidation("phone_number", "Nomor telepon harus 8 sampai 12 karakter.")

		// err = errors.New("Nomor telepon harus 8 sampai 12 karakter.")
		// o.Failure("phone_number", "Nomor telepon harus 8 sampai 12 karakter.")
	}

	// //check Merchant
	custGP, err := s.opt.Client.BridgeServiceGrpc.GetCustomerGPList(ctx, &bridge_service.GetCustomerGPListRequest{
		Phone:    req.Data.PhoneNumber,
		Limit:    1,
		Offset:   0,
		Inactive: "0",
	})
	// customer, err := s.opt.Client.BridgeServiceGrpc.GetCustomerDetail(ctx, &bridge_service.GetCustomerDetailRequest{
	// 	PhoneNumber: req.Data.PhoneNumber,
	// })
	if err != nil {
		if !global.IsLimit(redisKey, expPeriod, maxRetry) {
			// err = errors.New("Terlalu banyak percobaan untuk login dengan nomor ini - Silahkan coba lagi dalam " + strconv.Itoa(period) + " detik")
			err = edenlabs.ErrorValidation("otp", "Terlalu banyak percobaan untuk login dengan nomor ini - Silahkan coba lagi dalam "+strconv.Itoa(period)+" detik")
			return
		}
		err = edenlabs.ErrorValidation("phone_number.invalid", "Nomor yang anda masukkan tidak terdaftar.")

		// err = errors.New("Nomor yang anda masukkan tidak terdaftar.")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if len(custGP.Data) == 0 {
		if !global.IsLimit(redisKey, expPeriod, maxRetry) {
			err = edenlabs.ErrorValidation("otp", "Terlalu banyak percobaan untuk login dengan nomor ini - Silahkan coba lagi dalam "+strconv.Itoa(period)+" detik")
			// err = errors.New("Terlalu banyak percobaan untuk login dengan nomor ini - Silahkan coba lagi dalam " + strconv.Itoa(period) + " detik")
			return
		}
		err = edenlabs.ErrorValidation("phone_number.invalid", "Nomor yang anda masukkan tidak terdaftar.")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	customer, err := s.opt.Client.CrmServiceGrpc.GetCustomerDetail(ctx, &crm_service.GetCustomerDetailRequest{
		CustomerIdGp: custGP.Data[0].Custnmbr,
	})
	if err != nil {
		if !global.IsLimit(redisKey, expPeriod, maxRetry) {
			// err = errors.New("Terlalu banyak percobaan untuk login dengan nomor ini - Silahkan coba lagi dalam " + strconv.Itoa(period) + " detik")
			err = edenlabs.ErrorValidation("otp", "Terlalu banyak percobaan untuk login dengan nomor ini - Silahkan coba lagi dalam "+strconv.Itoa(period)+" detik")
			return
		}
		err = edenlabs.ErrorValidation("phone_number.invalid", "Nomor yang anda masukkan tidak terdaftar.")

		// err = errors.New("Nomor yang anda masukkan tidak terdaftar.")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	userCustomer := &model.UserCustomer{
		CustomerID: customer.Data.Id,
		// CustomerIDGP: strconv.Itoa(int(customer.Data.Id)),
	}
	userCustomer, err = s.RepositoryUserCustomer.GetDetail(ctx, userCustomer)
	if err != nil {
		err = edenlabs.ErrorValidation("phone_number.invalid", "Nomor yang anda masukkan tidak terdaftar.")
		// err = errors.New("Nomor yang anda masukkan tidak terdaftar.")
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if userCustomer.Status != 1 {
		// err = errors.New("Akun ini tidak aktif. Silahkan hubungi tim Customer Service kami.")
		err = edenlabs.ErrorValidation("phone_number.invalid", "Akun ini tidak aktif. Silahkan hubungi tim Customer Service kami.")
		// o.Failure("phone_number.invalid", "Akun ini tidak aktif. Silahkan hubungi tim Customer Service kami.")
	}

	_, err = s.RepositoryWhiteListLogin.GetDetail(ctx, 0, req.Data.PhoneNumber, req.Data.OTP)
	if err == nil {
		// if !global.IsMaxLimit(redisKey, maxRetry) {
		// 	err = edenlabs.ErrorValidation("otp", "Terlalu banyak percobaan untuk login dengan nomor ini - Silahkan coba lagi dalam "+strconv.Itoa(period)+" detik")
		// 	return
		// }
	} else if err != nil {
		//kalo ga sama dengan nil berarti datanya kosong,kalau datanya kosong di whitelistlogin,berarti bisa login.

		db := s.opt.Database.Read
		db.Raw("SELECT otp FROM otp_outgoing oo WHERE oo.phone_number = ? "+
			"AND oo.created_at BETWEEN NOW()-INTERVAL ? second "+
			"AND NOW() "+
			"AND oo.otp_status = 1 "+
			"AND oo.usage_type = 1 ORDER BY oo.created_at DESC LIMIT 1", req.Data.PhoneNumber, OTPValidDuration).QueryRow(&otp)

		if otp != req.Data.OTP {
			if !global.IsLimit(redisKey, expPeriod, maxRetry) {
				err = edenlabs.ErrorValidation("otp", "Terlalu banyak percobaan untuk login dengan nomor ini - Silahkan coba lagi dalam "+strconv.Itoa(period)+" detik")
				return
			}
			err = edenlabs.ErrorValidation("otp.invalid", "Kode OTP tidak valid. Silakan masukkan kode OTP yang benar.")
			// err = errors.New("Kode OTP tidak valid. Silakan masukkan kode OTP yang benar.")
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		if otp == "" {
			if !global.IsLimit(redisKey, expPeriod, maxRetry) {
				err = edenlabs.ErrorValidation("otp", "Terlalu banyak percobaan untuk login dengan nomor ini - Silahkan coba lagi dalam "+strconv.Itoa(period)+" detik")
				// err = errors.New("Terlalu banyak percobaan untuk login dengan nomor ini - Silahkan coba lagi dalam " + strconv.Itoa(period) + " detik")
				return
			}
			err = edenlabs.ErrorValidation("otp.invalid", "Kode OTP sudah tidak berlaku. Silakan coba kirim ulang.")
			// err = errors.New("Kode OTP sudah tidak berlaku. Silakan coba kirim ulang.")
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		if !global.IsMaxLimit(redisKey, maxRetry) {
			// err = errors.New("Terlalu banyak percobaan untuk login dengan nomor ini - Silahkan coba lagi dalam " + strconv.Itoa(period) + " detik")
			err = edenlabs.ErrorValidation("otp", "Terlalu banyak percobaan untuk login dengan nomor ini - Silahkan coba lagi dalam "+strconv.Itoa(period)+" detik")

			return
		}
	}

	jwtInit := jwt.NewJWT([]byte(s.opt.Config.Jwt.Key))
	uc := jwt.UserMobile{
		PhoneNo:   req.Data.PhoneNumber,
		ExpiresAt: time.Now().Add(time.Hour * 4).Unix(),
		Timezone:  req.Timezone,
		//StandardClaims: jwt.StandardClaims{},
	}

	jwtGenerate, err := jwtInit.Create(uc)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	// jwtGenerate = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwaG9uZV9ubyI6IjgwMDAwMTExMTAxIiwiZXhwaXJlX2F0IjoxNjc4MTc0MTA2LCJ0aW1lem9uZSI6IkFzaWEvSmFrYXJ0YSJ9.4LD2NcAoCr2Itdvvm-HGz7IFIusmhFzEO-lV1SJy1hY"
	userCustomer.LastLoginAt = time.Now()
	userCustomer.LoginToken = jwtGenerate
	userCustomer.ForceLogout = 2
	userCustomer.Verification = 2
	userCustomer.FirebaseToken = req.Data.FcmToken

	err = s.RepositoryUserCustomer.Update(ctx, userCustomer)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	db2 := s.opt.Database.Write
	tx, err := db2.BeginWithCtx(ctx)
	_, err = tx.Raw("UPDATE otp_outgoing oo set oo.otp_status=2 where phone_number =? and otp =? ", req.Data.PhoneNumber, OTPValidDuration).Exec()
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	tx.Commit()

	res.Token = jwtGenerate
	return
}
func UpdateOTP() {

}
func (s *AuthService) ResponseOTP(ctx context.Context, req dto.RequestGetOTP) (res dto.ResponseSendOtp, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AuthService.ResponseOTP")
	var e error
	var otpCount, otpMaxReq int
	var maxRequestOTPDuration float64
	var createdAt string

	// //to get expired period to limit how many seconds value from database
	// configApp, _ := repository.GetConfigApp("attribute", "otp_resend_duration")
	// if e != nil {
	// 	o.Failure("", "")
	// 	return o
	// }
	// period, _ := strconv.Atoi(configApp.Value)
	// expPeriod := time.Duration(period) * time.Second

	// //to get retry count value from database
	// configApp, _ = repository.GetConfigApp("attribute", "otp_max_resend_count")

	// maxRetry, _ := strconv.Atoi(configApp.Value)

	// redisKey := "ctr_otp_" + r.Data.IPAddress + "_" + r.Data.PhoneNumber

	// redisKey2 := "otp_" + r.Data.IPAddress + "_" + r.Data.PhoneNumber

	// otpLimit := isLimitReqOTP(redisKey, redisKey2, 30*time.Minute, expPeriod, maxRetry, r.Data.Type)

	// //sudah retry 5 kali dengan nomor yang sama dalam 30 menit
	// if otpLimit == 0 {
	// 	o.Failure("otp", "Terlalu banyak percobaan untuk nomor ini - Silahkan coba lagi dalam 30 menit")
	// 	return o

	// }

	// check number only
	if !validation.NumberOnly(req.Data.PhoneNumber) {
		// o.Failure("phone_number", "Nomor telepon harus dalam format angka.")
	}
	//check first char
	if string(req.Data.PhoneNumber[0]) == "0" {
		// o.Failure("phone_number", "Nomor telepon tidak valid.")
	}
	if len(req.Data.PhoneNumber) < 8 || len(req.Data.PhoneNumber) > 12 {
		// o.Failure("phone_number", "Nomor telepon harus 8 sampai 12 karakter.")
	}
	if req.Data.Type == "1" {
		//request otp tapi masih ada di redis dalam batasan 1 menit
		// if otpLimit == 2 {
		// 	o.Failure("otp", "Terlalu banyak percobaan untuk nomor ini - Silahkan coba lagi dalam "+strconv.Itoa(period)+" detik")
		// 	return o
		// }

		//check Customer
		customer, err := s.opt.Client.BridgeServiceGrpc.GetCustomerDetail(ctx, &bridge_service.GetCustomerDetailRequest{
			PhoneNumber: req.Data.PhoneNumber,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			// o.Failure("phone_number.invalid", "Nomor yang anda masukkan tidak terdaftar.")
			return res, err
		}
		if customer.Data.PhoneNumber != req.Data.PhoneNumber {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return res, err
		}
		if customer.Data.Id == 0 {
			// o.Failure("phone_number", "Nomor yang anda masukkan tidak terdaftar.")
		}

		userCustomer := &model.UserCustomer{CustomerID: customer.Data.Id, CustomerIDGP: strconv.Itoa(int(customer.Data.Id))}

		userCustomer, err = s.RepositoryUserCustomer.GetDetail(ctx, userCustomer)
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return res, err
		}
		if userCustomer.Status != 1 {
			// o.Failure("phone_number.invalid", "Akun ini tidak aktif. Silahkan hubungi tim Customer Service kami.")
		}
	} else if req.Data.Type == "2" {
		//ini untuk registrasi personal
		// check character only
		if req.Data.RegistrationRequest.CustomerName != "" {
			if !validation.CharacterOnly(req.Data.RegistrationRequest.CustomerName) {
				// o.Failure("customer_name", util.ErrorValidNameFormat())
			}
		} else {
			// o.Failure("customer_name", "Nama harus diisi.")
		}
		// check number only
		if !validation.NumberOnly(req.Data.RegistrationRequest.CustomerPhoneNumber) {
			// o.Failure("merchant_phone_number", util.ErrorValidPhoneFormat())
		}
		if !validation.NumberOnly(req.Data.RegistrationRequest.CustomerAltPhoneNumber) {
			// o.Failure("merchant_alt_phone_number", util.ErrorValidPhoneFormat())
		}
		//check first phone number
		if string(req.Data.RegistrationRequest.CustomerPhoneNumber[0]) == "0" {
			// o.Failure("merchant_phone_number", "Nomor telepon tidak valid.")
		}
		if len(req.Data.RegistrationRequest.CustomerPhoneNumber) < 8 || len(req.Data.RegistrationRequest.CustomerPhoneNumber) > 12 {
			// o.Failure("merchant_phone_number", "Nomor telepon harus 8 sampai 12 karakter.")
		}
		if string(req.Data.RegistrationRequest.CustomerAltPhoneNumber[0]) == "0" {
			// o.Failure("merchant_alt_phone_number", "Nomor telepon tidak valid.")
		}
		if len(req.Data.RegistrationRequest.CustomerAltPhoneNumber) < 8 || len(req.Data.RegistrationRequest.CustomerAltPhoneNumber) > 12 {
			// o.Failure("merchant_alt_phone_number", "Nomor telepon harus 8 sampai 12 karakter.")
		}
		// check valid email
		if req.Data.RegistrationRequest.CustomerEmail != "" {
			if !validation.EmailOnly(req.Data.RegistrationRequest.CustomerEmail) {
				// o.Failure("merchant_email", util.ErrorValidEmailFormat())
			}
		} else {
			// o.Failure("merchant_email", "Email harus diisi.")
		}
		if req.Data.RegistrationRequest.CustomerPhoneNumber == "" {
			// o.Failure("phone_number.required", util.ErrorInputRequired("phone number"))
		}

		if req.Data.RegistrationRequest.CustomerBirthDate != "" {
			if req.Data.RegistrationRequest.BirthDateAt, err = time.Parse("2006-01-02", req.Data.RegistrationRequest.CustomerBirthDate); err != nil {
				// o.Failure("birth_date.invalid", "Invalid birth_date")
			}
		}

		if req.Data.RegistrationRequest.ReferenceInfo == "" {
			// o.Failure("reference_info.required", util.ErrorInputRequired("reference info"))
		}

		//is phone number have a merchant
		if customer, err := s.opt.Client.BridgeServiceGrpc.GetCustomerGPList(ctx, &bridge_service.GetCustomerGPListRequest{
			Phone:  req.Data.PhoneNumber,
			Limit:  1,
			Offset: 0,
		}); err == nil {
			//if len data lebih dari 0{}
			if len(customer.Data) > 0 {
				err = edenlabs.ErrorValidation("merchant_phone_number", "Nomor ponsel sudah terdaftar. Silakan masukan nomor lain yang belum terdaftar.")
				// o.Failure("merchant_phone_number", util.ErrorRegisteredPhoneNumber())

				return res, err
			}
		}

		//check email
		if email, err := s.opt.Client.CrmServiceGrpc.GetCustomerDetail(ctx, &crm_service.GetCustomerDetailRequest{
			Email: req.Data.RegistrationRequest.CustomerEmail,
		}); err == nil {
			//if len data lebih dari 0{}
			if email.Data.Id != 0 {
				err = edenlabs.ErrorValidation("merchant_email", "Email sudah terdaftar. Silakan masukan email yang belum terdaftar.")
				// o.Failure("merchant_phone_number", util.ErrorRegisteredPhoneNumber())
				e = err
				return res, err
			}
		}

		//check referral_code
		if req.Data.RegistrationRequest.ReferrerCode != "" {
			refCode, err := s.opt.Client.CrmServiceGrpc.GetCustomerDetail(ctx, &crm_service.GetCustomerDetailRequest{
				ReferrerCode: req.Data.RegistrationRequest.ReferrerCode,
			})
			if err != nil {
				//span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorValidation("referral_code", "Code referral tidak ditemukan.")
				e = err
				return res, err
			}
			if refCode.Data.Id == 0 {
				err = edenlabs.ErrorValidation("referral_code", "Code referral tidak ditemukan.")
				e = err
				return res, err
			}
		}

		if req.Data.RegistrationRequest.AdmDivisionID == "" {
			// o.Failure("sub_district.invalid", util.ErrorInvalidData("sub district"))
		} else {
			// //save adm division here
			// req.Data.RegistrationRequest.SubDistrict = &model.SubDistrict{ID: area}
			// req.Data.RegistrationRequest.SubDistrict.Read("ID")
			// req.Data.RegistrationRequest.Area = &model.Area{ID: r.Data.RegistrationRequest.SubDistrict.Area.ID}

			//get region first then check price set based on region
			// req.Data.RegistrationRequest.PriceSet = &model.PriceSet{ID: req.Data.RegistrationRequest.SubDistrict.Area.ID}
			// req.Data.RegistrationRequest.WarehouseCoverage = &model.WarehouseCoverage{SubDistrict: req.Data.RegistrationRequest.SubDistrict, MainWarehouse: 1}
			// req.Data.RegistrationRequest.WarehouseCoverage.Read("SubDistrict", "MainWarehouse")

			// if e = req.Data.RegistrationRequest.PriceSet.Read("ID"); e != nil {
			// 	o.Failure("area", util.ErrorInvalidData("area"))
			// }

			// if e = req.Data.RegistrationRequest.WarehouseCoverage.Warehouse.Read("ID"); e != nil {
			// 	o.Failure("warehouse.invalid", util.ErrorInvalidData("warehouse coverage"))
			// } else {
			// 	if req.Data.RegistrationRequest.WarehouseCoverage.Warehouse.Status != 1 {
			// 		o.Failure("warehouse.inactive", util.ErrorActive("default warehouse"))
			// 	}
		}
	}

	_, e = s.RepositoryWhiteListLogin.GetDetail(ctx, 0, req.Data.PhoneNumber, "")
	if e == nil {
		return res, e
	}

	db := s.opt.Database.Read
	db.Raw("SELECT COUNT(otp) FROM otp_outgoing oo WHERE oo.phone_number = ? "+
		"AND oo.created_at BETWEEN NOW()-INTERVAL ? second AND NOW() AND oo.otp_status = 1 AND oo.usage_type = ?", req.Data.PhoneNumber, 1800, req.Data.Type).QueryRow(&otpCount)
	otpMaxReq = s.opt.Env.GetInt("OTP.MAX_REQUEST_OTP")
	if otpCount >= otpMaxReq {
		db.Raw("SELECT created_at FROM otp_outgoing WHERE phone_number = ? AND "+
			"otp_status = 1 AND usage_type = ? ORDER BY created_at DESC LIMIT 1", req.Data.PhoneNumber, req.Data.Type).QueryRow(&createdAt)
		if req.Data.CreatedAt, e = time.Parse("2006-01-02 15:04:05", createdAt); e != nil {
			// o.Failure("created_at.invalid", util.ErrorInvalidData("created at"))
		} else {
			wib, _ := time.LoadLocation("Asia/Jakarta")
			currentTime := time.Now().In(wib)
			t1, _ := time.Parse("15:04:05", req.Data.CreatedAt.Format("15:04:05"))
			t2, _ := time.Parse("15:04:05", currentTime.Format("15:04:05"))
			maxRequestOTPDuration = 1800
			maxOtptoMinutes := maxRequestOTPDuration / 60
			fmt.Print(maxOtptoMinutes)
			if t1.Sub(t2).Seconds() <= maxRequestOTPDuration {
				// o.Failure("otp.invalid", "Anda telah mencapai batas kirim OTP. Silakan lakukan login ulang "+strconv.FormatFloat(maxOtptoMinutes, 'f', -1, 64)+" menit lagi.")
			}
		}

	}

	return
}

func (s *AuthService) Session(ctx echo.Context) (m *dto.SessionDataCustomer, e error) {
	var userCustomer *model.UserCustomer
	m = new(dto.SessionDataCustomer)
	auth := ctx.Request().Header.Get("Authorization")
	if auth != "" && len(auth) > 6 {
		bearer := auth[:strings.IndexByte(auth, ' ')]
		if bearer != "Bearer" {
			e = echo.NewHTTPError(http.StatusRequestTimeout, "invalid or expired jwt token")
		} else {
			var (
				membershipLevelResponse      *dto.MembershipLevelResponse
				membershipCheckpointResponse *dto.MembershipCheckpointResponse
			)

			tokenLogin := strings.Replace(auth, "Bearer ", "", 1) // get token from Authorization
			userCustomer = &model.UserCustomer{LoginToken: tokenLogin, Status: 1}
			userCustomer, err := s.RepositoryUserCustomer.GetDetail(ctx.Request().Context(), userCustomer)
			if err != nil {
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}

			// Get customer detail From CRM Service
			customerDetail, err := s.opt.Client.CrmServiceGrpc.GetCustomerDetail(ctx.Request().Context(), &crm_service.GetCustomerDetailRequest{
				Id: userCustomer.CustomerID,
			})
			if err != nil {
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("crm", "customer")
				return
			}

			customer, err := s.opt.Client.BridgeServiceGrpc.GetCustomerGPList(ctx.Request().Context(), &bridge_service.GetCustomerGPListRequest{
				// Id: customerDetail.Data.CustomerIdGp,
				Limit:  1,
				Offset: 0,
				Id:     customerDetail.Data.CustomerIdGp,
			})
			if err != nil {
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}
			if len(customer.Data) == 0 {
				err = edenlabs.ErrorRpcNotFound("bridge", "customer")
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}
			// Get membership From campaign Service
			if customerDetail.Data.MembershipLevelId != 0 {
				var membershipLevel *campaign_service.GetMembershipLevelDetailResponse
				membershipLevel, err = s.opt.Client.CampaignServiceGrpc.GetMembershipLevelDetail(ctx.Request().Context(), &campaign_service.GetMembershipLevelDetailRequest{
					Id: customerDetail.Data.MembershipLevelId,
				})
				if err != nil {
					s.opt.Logger.AddMessage(log.ErrorLevel, err)
					err = edenlabs.ErrorRpcNotFound("campaign", "membership level")
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

				var membershipCheckpoint *campaign_service.GetMembershipCheckpointDetailResponse
				membershipCheckpoint, err = s.opt.Client.CampaignServiceGrpc.GetMembershipCheckpointDetail(ctx.Request().Context(), &campaign_service.GetMembershipCheckpointDetailRequest{
					Id: customerDetail.Data.MembershipCheckpointId,
				})
				if err != nil {
					s.opt.Logger.AddMessage(log.ErrorLevel, err)
					err = edenlabs.ErrorRpcNotFound("campaign", "membership checkpoint")
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

			m.Customer = &dto.SessionCustomer{
				ID:           strconv.Itoa(int(customerDetail.Data.Id)),
				Code:         customerDetail.Data.CustomerIdGp,
				ReferralCode: customerDetail.Data.ReferralCode,
				Name:         customer.Data[0].Custname,
				// Gender:                     strconv.Itoa(int(customer.Data.Gender)),
				// BirthDate:                  customer.Data.BirthDate.AsTime(),
				PicName:        customer.Data[0].Custname,
				PhoneNumber:    customer.Data[0].PhonE1,
				AltPhoneNumber: customer.Data[0].PhonE2,
				Email:          customerDetail.Data.Email,
				// Password:                   customer.Data.Password,
				BillingAddress: customer.Data[0].AddresS1 + " " + customer.Data[0].AddresS2 + " " + customer.Data[0].AddresS3,
				// Note:                       customer.Data.Note,
				ReferenceInfo: customerDetail.Data.ReferenceInfo,
				// TagCustomer:                customer.Data.TagCustomer,
				// Status:                     strconv.Itoa(int(customer.Data.Status)),
				// Suspended:                  strconv.Itoa(int(customer.Data.Suspended)),
				UpgradeStatus: strconv.Itoa(int(customerDetail.Data.UpgradeStatus)),
				// CustomerGroup:              strconv.Itoa(int(customer.Data.CustomerGroup)),
				// TagCustomerName:            customer.Data.TagCustomerName,
				ReferrerCode: customerDetail.Data.ReferrerCode,
				// CreatedAt:                  customer.Data.CreatedAt.AsTime(),
				// CreatedBy:                  strconv.Itoa(int(customer.Data.CreatedBy)),
				// LastUpdatedAt:              customer.Data.LastUpdatedAt.AsTime(),
				// LastUpdatedBy:              strconv.Itoa(int(customer.Data.LastUpdatedBy)),
				TotalPoint: strconv.Itoa(int(customerDetail.Data.TotalPoint)),
				// CustomerTypeCreditLimit:    strconv.Itoa(int(customer.Data.CustomerTypeCreditLimit)), //customer class
				// EarnedPoint:                strconv.Itoa(int(customer.Data.EarnedPoint)),
				// RedeemedPoint:              strconv.Itoa(int(customer.Data.RedeemedPoint)),
				// CustomCreditLimit:          strconv.Itoa(int(customer.Data.CustomCreditLimit)), //dari gp
				// CreditLimitAmount:          strconv.Itoa(int(customer.Data.CreditLimitAmount)),
				ProfileCode: customerDetail.Data.ProfileCode,
				// RemainingCreditLimitAmount: strconv.Itoa(int(customer.Data.RemainingCreditLimitAmount)),
				// AverageSales:               strconv.Itoa(int(customer.Data.AverageSales)),
				// RemainingOutstanding:       strconv.Itoa(int(customer.Data.RemainingOutstanding)), //dari gp
				// OverdueDebt:                strconv.Itoa(int(customer.Data.OverdueDebt)),          //dari gp
				// KTPPhotosUrl:               customer.Data.KTPPhotosUrl,
				// CustomerPhotosUrl:          customer.Data.MerchantPhotosUrl,
				// KTPPhotosUrlArr:            customer.Data.KTPPhotosUrlArr,
				// CustomerPhotosUrlArr:       customer.Data.MerchantPhotosUrlArr, //dikita
				MembershipLevelID:      strconv.Itoa(int(customerDetail.Data.MembershipLevelId)),
				MembershipCheckpointID: strconv.Itoa(int(customerDetail.Data.MembershipCheckpointId)),
				MembershipRewardID:     strconv.Itoa(int(customerDetail.Data.MembershipRewardId)),
				MembershipRewardAmount: strconv.Itoa(int(customerDetail.Data.MembershipRewardAmmount)),
				// TermPaymentSlsId:           strconv.Itoa(int(customer.Data.TermPaymentSlsId)),
				// BirthDateString:            customer.Data.BirthDateString,
				UserCustomer:   &model.UserCustomer{},
				InvoiceTerm:    &model.InvoiceTerm{},
				TermPaymentSls: "PBD", // change after get data from GP
				PaymentMethod:  &model.PaymentMethod{},
				//CustomerType:               &model.CustomerType{},
				FinanceArea:      &model.Region{},
				PaymentGroup:     &model.PaymentGroup{},
				ProspectCustomer: &model.ProspectCustomer{},
				CustomerPriceSet: []*model.CustomerPriceSet{},
				CustomerAccNum:   []*model.CustomerAccNum{},
				// CustomerType:         customer.Data.CustomerTypeId,
				//CustomerType:         customer.Data[0].CustomerType[0].GnL_Cust_Type_ID,
				ReferrerCustomer:     &model.Customer{},
				MembershipLevel:      membershipLevelResponse,
				MembershipCheckpoint: membershipCheckpointResponse,
				MembershipReward:     &model.MembershipReward{},
			}
			m.Customer.UserCustomer = userCustomer
			if len(customer.Data[0].CustomerType) == 0 {
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			} else {
				m.Customer.CustomerType = customer.Data[0].CustomerType[0].GnL_Cust_Type_ID
			}

			custType, err := s.opt.Client.BridgeServiceGrpc.GetCustomerTypeGPList(ctx.Request().Context(), &bridge_service.GetCustomerTypeGPListRequest{
				Id:     m.Customer.CustomerType,
				Limit:  1,
				Offset: 0,
			})

			if len(custType.Data) == 0 {
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}

			m.Customer.CustomerGroup = custType.Data[0].GnL_Cust_GroupDesc

			// address, err := s.opt.Client.BridgeServiceGrpc.GetAddressList(ctx.Request().Context(), &bridge_service.GetAddressListRequest{
			// 	CustomerId: customer.Data.Id,
			// })
			address, err := s.opt.Client.BridgeServiceGrpc.GetAddressGPList(ctx.Request().Context(), &bridge_service.GetAddressGPListRequest{
				// Id: customer.Data[0].Adrscode[0].Adrscode,
				Limit:          10,
				Offset:         1,
				CustomerNumber: customer.Data[0].Custnmbr,
			})

			if err != nil {
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}
			if len(address.Data) == 0 {
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}
			tempAddress := &bridge_service.AddressGP{}
			for _, v := range address.Data {
				if v.TypeAddress == "ship_to" {
					tempAddress = v
				}
			}
			if tempAddress == nil {
				//span.RecordError(err)
				fmt.Println(address.Data)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return nil, err
			}
			// admDivision, err := s.opt.Client.BridgeServiceGrpc.GetAdmDivisionDetail(ctx.Request().Context(), &bridge_service.GetAdmDivisionDetailRequest{
			// 	Id: address.Data[0].AdmDivisionId,
			// })
			admDivision, err := s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPList(ctx.Request().Context(), &bridge_service.GetAdmDivisionGPListRequest{
				Limit:  1,
				Offset: 0,
				// AdmDivisionCode: address.Data[0].GnL_Administrative_Code,
				AdmDivisionCode: tempAddress.AdministrativeDiv.GnlAdministrativeCode,
			})
			if err != nil {
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}
			m.Address = &dto.SessionAddress{
				ID:              tempAddress.Adrscode,
				Code:            tempAddress.Adrscode,
				Name:            tempAddress.ShipToName,
				PicName:         m.Customer.Name,
				PhoneNumber:     tempAddress.PhonE1,
				AltPhoneNumber:  tempAddress.PhonE2,
				AddressName:     tempAddress.ShipToName,
				ShippingAddress: tempAddress.AddresS1 + " " + tempAddress.AddresS2 + " " + tempAddress.AddresS3,
				Latitude:        tempAddress.GnL_Latitude,
				Longitude:       tempAddress.GnL_Longitude,
				Note:            tempAddress.GnL_Address_Note,
				MainBranch:      "",
				// Status:             strconv.Itoa(int(address.Data[0].Status)),
				// CreatedAt:          address.Data[0].CreatedAt.AsTime(),
				CreatedBy: "",
				// LastUpdatedAt:      address.Data[0].UpdatedAt.AsTime(),
				LastUpdatedBy:      "",
				PinpointValidation: "",
				AdmDivisionId:      tempAddress.AdministrativeDiv.GnlAdministrativeCode,
				RegionID:           admDivision.Data[0].Region,
				CustomerID:         customer.Data[0].Custnmbr,
				ArchetypeID:        tempAddress.GnL_Archetype_ID,
				PriceSetID:         "",
				SiteID:             tempAddress.Locncode,
				SalesPersonID:      "",
				Customer:           &model.Customer{},
				Region:             &model.Region{},
				Archetype:          &model.Archetype{},
				PriceSet:           &model.PriceSet{},
				Site:               &model.Site{},
				Salesperson:        &model.Staff{},
				AdmDivision:        &model.AdmDivision{},
				StatusConvert:      "",
				City:               admDivision.Data[0].City,
				SubDistrictID:      admDivision.Data[0].Subdistrict,
			}
			// 	 else {
			// 			if e = m.Merchant.PaymentGroup.Read("ID"); e != nil {
			// 				e = echo.NewHTTPError(http.StatusForbidden, "invalid or expired jwt token")
			// 			}

			// 			if e = m.Merchant.PaymentTerm.Read("ID"); e != nil {
			// 				e = echo.NewHTTPError(http.StatusForbidden, "invalid or expired jwt token")
			// 			}

			// 			m.Merchant.CustomerBusinessType = m.Merchant.BusinessType.ID
			// 			m.Merchant.userCustomer = userCustomer

			// 			// check membership level of merchant
			// 			if m.Merchant.MembershipLevelID == 0 {
			// 				o := orm.NewOrm()
			// 				o.Using("read_only")

			// 				type LevelCheckpoint struct {
			// 					LevelID      int64 `orm:"column(membership_level_id)"`
			// 					CheckpointID int64 `orm:"column(membership_checkpoint_id)"`
			// 				}

			// 				var (
			// 					levelCheckpoint LevelCheckpoint
			// 					q               string
			// 					isExist         bool
			// 				)

			// 				// check if merchant's business type is eligible for membership
			// 				q = "select exists(" +
			// 					"select m.id " +
			// 					"from merchant m " +
			// 					"join config_app ca on ca.`attribute` = 'eligible_membership_business_type' and find_in_set(m.business_type_id, ca.value) > 0 " +
			// 					"where m.id = ?)"
			// 				e = o.Raw(q, m.Merchant.ID).QueryRow(&isExist)

			// 				if isExist {
			// 					// get lowest membership level & checkpoint
			// 					q = "select ml.id membership_level_id, mc.id membership_checkpoint_id " +
			// 						"from membership_level ml " +
			// 						"join membership_checkpoint mc on ml.id = mc.membership_level_id " +
			// 						"where ml.status = 1 and mc.status = 1 " +
			// 						"order by mc.checkpoint asc " +
			// 						"limit 1"
			// 					if e = o.Raw(q).QueryRow(&levelCheckpoint); e != nil || levelCheckpoint.LevelID == 0 {
			// 						e = echo.NewHTTPError(http.StatusForbidden, "invalid or expired jwt token")
			// 					}

			// 					// update membership level & checkpoint of merchant
			// 					o.Using("default")
			// 					q = "update merchant set membership_level_id = ?, membership_checkpoint_id = ? where id = ?"
			// 					if _, e = o.Raw(q, levelCheckpoint.LevelID, levelCheckpoint.CheckpointID, m.Merchant.ID).Exec(); e != nil {
			// 						e = echo.NewHTTPError(http.StatusForbidden, "invalid or expired jwt token")
			// 					}
			// 					m.Merchant.MembershipLevelID = levelCheckpoint.LevelID
			// 				}
			// 			}

			// 			// get membership level of merchant
			// 			if m.Merchant.MembershipLevelID == 0 {
			// 				m.Merchant.MembershipLevel = nil
			// 			} else {
			// 				m.Merchant.MembershipLevel = &model.MembershipLevel{ID: m.Merchant.MembershipLevelID}
			// 				m.Merchant.MembershipLevel.Read("ID")
			// 			}
			// 		}
			// 	} else {
			// 		e = echo.NewHTTPError(http.StatusForbidden, "invalid or expired jwt token")
			// 	}

			// 	//get merchant account number
			// 	o := orm.NewOrm()
			// 	o.Using("read_only")
			// 	o.Raw(
			// 		"SELECT * FROM merchant_acc_num "+
			// 			"WHERE merchant_id = ?", m.Merchant.ID).QueryRows(&m.Merchant.MerchantAccNum)

			// 	for i := 0; i < len(m.Merchant.MerchantAccNum); i++ {
			// 		m.Merchant.MerchantAccNum[i].PaymentChannel.Read("ID")
			// 	}
			// }
			// } else {
			// 	e = echo.NewHTTPError(http.StatusForbidden, "invalid or expired jwt token")
		}
	}

	return m, e
}

func CustomerSession(ctx echo.Context) (m *dto.SessionDataCustomer, e error) {
	var s = global.Setup.Common
	var r = repository.NewUserCustomerRepository()

	var (
		membershipLevelResponse      *dto.MembershipLevelResponse
		membershipCheckpointResponse *dto.MembershipCheckpointResponse
	)

	m = new(dto.SessionDataCustomer)
	auth := ctx.Request().Header.Get("Authorization")
	if len(auth) < 6 {
		return nil, echo.NewHTTPError(http.StatusUnprocessableEntity, "Authorization is required")
	}

	tokenLogin := strings.Replace(auth, "Bearer ", "", 1)
	userCustomerPhoneNo, e := GetUserCustomerPhoneNo(tokenLogin)

	if e != nil {
		e = echo.NewHTTPError(http.StatusBadRequest, "token user not valid")
		return nil, e
	}
	customer, err := s.Client.BridgeServiceGrpc.GetCustomerGPList(ctx.Request().Context(), &bridge_service.GetCustomerGPListRequest{
		Phone:  userCustomerPhoneNo,
		Limit:  1,
		Offset: 0,
	})
	if err != nil {
		//span.RecordError(err)
		s.Logger.AddMessage(log.ErrorLevel, err)
		return nil, err
	}
	// Get customer detail From CRM Service
	customerDetail, err := s.Client.CrmServiceGrpc.GetCustomerDetail(ctx.Request().Context(), &crm_service.GetCustomerDetailRequest{
		CustomerIdGp: customer.Data[0].Custnmbr,
	})

	if err != nil {
		s.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("crm", "customer")
		return
	}

	var membershipLevelID, membershipCheckpointID int64
	// Get membership From campaign Service
	if customerDetail.Data.MembershipLevelId != 0 {
		var membershipLevel *campaign_service.GetMembershipLevelDetailResponse
		membershipLevel, err = s.Client.CampaignServiceGrpc.GetMembershipLevelDetail(ctx.Request().Context(), &campaign_service.GetMembershipLevelDetailRequest{
			Id: customerDetail.Data.MembershipLevelId,
		})
		if err != nil {
			s.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("campaign", "membership level")
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
		membershipLevelID = membershipLevel.Data.Id

		var membershipCheckpoint *campaign_service.GetMembershipCheckpointDetailResponse
		membershipCheckpoint, err = s.Client.CampaignServiceGrpc.GetMembershipCheckpointDetail(ctx.Request().Context(), &campaign_service.GetMembershipCheckpointDetailRequest{
			Id: customerDetail.Data.MembershipCheckpointId,
		})
		if err != nil {
			s.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("campaign", "membership checkpoint")
			return
		}
		membershipCheckpointResponse = &dto.MembershipCheckpointResponse{
			ID:                membershipCheckpoint.Data.Id,
			Checkpoint:        int8(membershipCheckpoint.Data.Checkpoint),
			TargetAmount:      membershipCheckpoint.Data.TargetAmount,
			Status:            int8(membershipCheckpoint.Data.Status),
			MembershipLevelID: membershipCheckpoint.Data.MembershipLevelId,
		}
		membershipCheckpointID = membershipCheckpoint.Data.Id
	}

	userCustomer := &model.UserCustomer{
		CustomerID: customerDetail.Data.Id,
		// CustomerIDGP: strconv.Itoa(int(customer.Data.Id)),
	}
	userCustomer, err = r.GetDetail(ctx.Request().Context(), userCustomer)
	if err != nil {
		s.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	if tokenLogin != userCustomer.LoginToken {
		// e = echo.NewHTTPError(http.StatusBadRequest, "token user not valid")
	}
	if len(customer.Data[0].CustomerType) == 0 {
		err = edenlabs.ErrorValidation("customer_type", "customer do not have customer type.")
		return
	}
	m.Customer = &dto.SessionCustomer{
		ID:                     strconv.Itoa(int(customerDetail.Data.Id)),
		Code:                   customer.Data[0].Custnmbr,
		Name:                   customer.Data[0].Custname,
		Email:                  customerDetail.Data.Email,
		ReferenceInfo:          customerDetail.Data.ReferenceInfo,
		UpgradeStatus:          strconv.Itoa(int(customerDetail.Data.UpgradeStatus)),
		TotalPoint:             strconv.Itoa(int(customerDetail.Data.TotalPoint)),
		ProfileCode:            customerDetail.Data.ProfileCode,
		MembershipLevelID:      strconv.Itoa(int(membershipLevelID)),
		MembershipCheckpointID: strconv.Itoa(int(membershipCheckpointID)),
		MembershipLevel:        membershipLevelResponse,
		MembershipCheckpoint:   membershipCheckpointResponse,
		MembershipRewardID:     strconv.Itoa(int(customerDetail.Data.MembershipRewardId)),
		MembershipRewardAmount: strconv.FormatFloat(float64(customerDetail.Data.MembershipRewardAmmount), 'f', 1, 64),
		CustomerType:           customer.Data[0].CustomerType[0].GnL_Cust_Type_ID,
		UserCustomer:           userCustomer,
		TermPaymentSls:         "COD",
	}

	// customerID, _ := strconv.Atoi(m.Customer.ID)
	address, err := s.Client.BridgeServiceGrpc.GetAddressGPList(ctx.Request().Context(), &bridge_service.GetAddressGPListRequest{
		Limit:          10,
		Offset:         1,
		CustomerNumber: customer.Data[0].Custnmbr,
	})
	if err != nil {
		//span.RecordError(err)
		s.Logger.AddMessage(log.ErrorLevel, err)
		return nil, err
	}
	tempAddress := &bridge_service.AddressGP{}
	for _, v := range address.Data {
		if v.TypeAddress == "ship_to" {
			tempAddress = v
		}
	}
	if tempAddress == nil {
		//span.RecordError(err)
		fmt.Println(address.Data)
		s.Logger.AddMessage(log.ErrorLevel, err)
		return nil, err
	}

	custType, err := s.Client.BridgeServiceGrpc.GetCustomerTypeGPList(ctx.Request().Context(), &bridge_service.GetCustomerTypeGPListRequest{
		Id:     m.Customer.CustomerType,
		Limit:  1,
		Offset: 0,
	})

	if len(custType.Data) == 0 {
		s.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	m.Customer.CustomerGroup = custType.Data[0].GnL_Cust_GroupDesc

	admDivision, err := s.Client.BridgeServiceGrpc.GetAdmDivisionGPList(ctx.Request().Context(), &bridge_service.GetAdmDivisionGPListRequest{
		Limit:  1,
		Offset: 0,
		// AdmDivisionCode: address.Data[0].GnL_Administrative_Code,
		AdmDivisionCode: tempAddress.AdministrativeDiv.GnlAdministrativeCode,
	})
	if err != nil {
		s.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	m.Address = &dto.SessionAddress{
		ID:              tempAddress.Adrscode,
		Code:            tempAddress.Adrscode,
		AdmDivisionId:   tempAddress.AdministrativeDiv.GnlAdministrativeCode,
		RegionID:        admDivision.Data[0].Region, //strconv.Itoa(int(admDivision.Data.RegionId)),
		SiteID:          tempAddress.Locncode,       // strconv.Itoa(int(address.Data.SiteId)),
		Name:            tempAddress.ShipToName,
		PicName:         "",
		PhoneNumber:     tempAddress.PhonE1,
		ShippingAddress: tempAddress.AddresS1 + " " + tempAddress.AddresS2 + " " + tempAddress.AddresS3,
		Note:            tempAddress.GnL_Address_Note,
		Latitude:        tempAddress.GnL_Latitude,
		Longitude:       tempAddress.GnL_Longitude,
		City:            tempAddress.City,
		SubDistrictID:   admDivision.Data[0].Subdistrict, // strconv.Itoa(int(admDivision.Data.SubDistrictId)),
		ArchetypeID:     tempAddress.GnL_Archetype_ID,
	}
	return
}
func GetUserCustomerPhoneNo(tokenLogin string) (phoneNo string, e error) {
	var m = global.Setup.Common
	jwtInit := jwt.NewJWT([]byte(m.Config.Jwt.Key))

	token, e := jwtInit.ParseMobile(tokenLogin)

	if e != nil && e.Error() != "Token is expired" {
		return "", e
	}

	claims := token.Claims.(*jwt.UserMobile)

	return claims.PhoneNo, e
}

func (s *AuthService) SignOut(ctx echo.Context, customerID int64) (m *dto.SessionDataCustomer, e error) {
	db2 := s.opt.Database.Write
	tx, err := db2.BeginWithCtx(ctx.Request().Context())
	_, err = tx.Raw(
		"UPDATE user_customer um "+
			"SET um.login_token = null, "+
			"um.firebase_token = null "+
			"WHERE um.customer_id = ?", customerID).Exec()
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	tx.Commit()
	return
}

func (s *AuthService) CheckPhoneNumber(ctx echo.Context, req *dto.CheckPhoneNumberRequest) (m *dto.CheckPhoneNumberRequest, e error) {
	c := ctx.Request().Context()

	c, span := s.opt.Trace.Start(c, "AuthService.Login")

	if validation.NumberOnly(req.Data.PhoneNumber) == false {
		// o.Failure("phone_number", "Nomor telepon harus dalam format angka.")
	}
	// //check first char
	if string(req.Data.PhoneNumber[0]) == "0" {
		//o.Failure("phone_number", "Nomor telepon tidak valid.")
	}
	if len(req.Data.PhoneNumber) < 8 || len(req.Data.PhoneNumber) > 12 {
		// o.Failure("phone_number", "Nomor telepon harus 8 sampai 12 karakter.")
	}

	// //check Merchant
	// customer, err := s.opt.Client.BridgeServiceGrpc.GetCustomerDetail(c, &bridge_service.GetCustomerDetailRequest{
	// 	PhoneNumber: req.Data.PhoneNumber,
	// })
	customer, err := s.opt.Client.BridgeServiceGrpc.GetCustomerGPList(c, &bridge_service.GetCustomerGPListRequest{
		Phone:    req.Data.PhoneNumber,
		Limit:    1,
		Offset:   0,
		Inactive: "0",
	})
	if err != nil {
		e = edenlabs.ErrorValidation("phone_number.invalid", "Nomor yang anda masukkan tidak terdaftar.")
		span.RecordError(e)
		s.opt.Logger.AddMessage(log.ErrorLevel, e)
		return
	}
	if len(customer.Data) == 0 {
		e = edenlabs.ErrorValidation("phone_number.invalid", "Nomor yang anda masukkan tidak terdaftar.")
		span.RecordError(e)
		s.opt.Logger.AddMessage(log.ErrorLevel, e)
		return
	}

	// _, err = s.opt.Client.CrmServiceGrpc.GetCustomerDetail(ctx.Request().Context(), &crm_service.GetCustomerDetailRequest{
	// 	CustomerIdGp: customer.Data[0].Custnmbr,
	// })
	// if err != nil {
	// 	var codeGenerator *configuration_service.GetGenerateCodeResponse
	// 	codeGenerator, err = s.opt.Client.ConfigurationServiceGrpc.GetGenerateCode(ctx.Request().Context(), &configuration_service.GetGenerateCodeRequest{
	// 		Format: "UCS",
	// 		Domain: "user_customer",
	// 		Length: 4,
	// 	})
	// 	CodeUserCustomer := codeGenerator.Data.Code
	// 	um := &model.UserCustomer{
	// 		Code:         CodeUserCustomer,
	// 		Verification: 2,
	// 		Status:       1,
	// 		ForceLogout:  2,
	// 		// FirebaseToken: req.FirebaseToken,
	// 	}
	// 	if e = s.RepositoryUserCustomer.Create(ctx.Request().Context(), um); e == nil {
	// 		// codeGenerator, err = s.opt.Client.ConfigurationServiceGrpc.GetGenerateCustomerCode(ctx.Request().Context(), &configuration_service.GetGenerateCodeRequest{
	// 		// 	Format: "EC",
	// 		// 	Domain: "customer",
	// 		// 	Length: 8,
	// 		// })
	// 		// codeGenerator, err = s.opt.Client.ConfigurationServiceGrpc.GetGenerateReferralCode(ctx.Request().Context(), &configuration_service.GetGenerateCodeRequest{
	// 		// 	Format: "CRF",
	// 		// 	Domain: "referral",
	// 		// 	Length: 8,
	// 		// })
	// 		reqCreateRequestCRM := &crm_service.CreateCustomerRequest{
	// 			Data: &crm_service.Customer{
	// 				CustomerIdGp:            customer.Data[0].Custnmbr,
	// 				ProspectiveCustomerId:   0,
	// 				MembershipLevelId:       0,
	// 				MembershipCheckpointId:  0,
	// 				TotalPoint:              0,
	// 				ProfileCode:             customer.Data[0].Custnmbr,
	// 				Email:                   "",
	// 				ReferenceInfo:           "",
	// 				UpgradeStatus:           0,
	// 				KtpPhotosUrl:            "",
	// 				CustomerPhotosUrl:       "",
	// 				CustomerSelfieUrl:       "",
	// 				MembershipRewardId:      0,
	// 				MembershipRewardAmmount: 0,
	// 				// ReferralCode:            req.CodeReferral,
	// 				// ReferrerId:              req.ReferrerID,
	// 				// ReferrerCode:            "",
	// 			},
	// 		}
	// 		fmt.Print(reqCreateRequestCRM)
	// 		customerResponse, err := s.opt.Client.CrmServiceGrpc.CreateCustomer(ctx.Request().Context(), reqCreateRequestCRM)
	// 		um.LastLoginAt = time.Now()
	// 		// um.LoginToken = jwtGenerate
	// 		um.ForceLogout = 2
	// 		um.Verification = 2
	// 		// um.CustomerIDGP = m.Custnmbr
	// 		um.CustomerID = customerResponse.Data.Id

	// 		err = s.RepositoryUserCustomer.Update(ctx.Request().Context(), um)
	// 		if err != nil {
	// 			span.RecordError(err)
	// 			s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 			return
	// 		}
	// 	}
	// }
	// if customer.Data[0].PhonE1 != req.Data.PhoneNumber {
	// }
	if customer.Data[0].Custnmbr == "" || customer.Data[0].Inactive == 1 {
		// o.Failure("phone_number", "Nomor yang anda masukkan tidak terdaftar.")
	}

	valueMaintenance, err := s.opt.Client.ConfigurationServiceGrpc.GetConfigAppDetail(ctx.Request().Context(), &configuration_service.GetConfigAppDetailRequest{
		Attribute: "maintenance_" + req.Platform,
	})
	if err != nil {

	}
	if valueMaintenance.Data.Value == "1" && "maintenance_"+req.Platform == "maintenance_orca" {
		// o.Failure("maintenance_state", "Sistem sedang dalam pemeliharaan.")
	}
	if valueMaintenance.Data.Value == "1" && "maintenance_"+req.Platform == "maintenance_mantis" {
		// o.Failure("maintenance_state", "Sistem sedang dalam pemeliharaan.")
	}

	m = &dto.CheckPhoneNumberRequest{
		Platform: req.Platform,
		Data: dto.CheckPhoneNumberLoginRequest{
			PhoneNumber: req.Data.PhoneNumber,
		},
	}

	return m, e
}

func (s *AuthService) SendOTPWASociomile(ctx context.Context, phoneNumber string, application int, usageType int) error {

	var (
		client    = &http.Client{}
		data      dto.SendOtpWASociomile
		message   string
		component dto.Component
		parameter dto.Parameter
	)
	var m = global.Setup.Common

	WhatsAppSociomileHost := m.Env.GetString("WHATSAPP_SOCIOMILE.HOST")
	WhatsAppSociomileApiKey := m.Env.GetString("WHATSAPP_SOCIOMILE.API_KEY")
	WhatsAppSociomileTemplateID := m.Env.GetString("WHATSAPP_SOCIOMILE.TEMPLATE_ID")

	ocl := m.Env.GetInt("OTP.CODE_LENGTH")
	otp := RandomInt(ocl, phoneNumber)
	messageConfig, err := m.Client.ConfigurationServiceGrpc.GetConfigAppDetail(ctx, &configuration_service.GetConfigAppDetailRequest{
		Attribute: "otp_msg",
	})
	message = messageConfig.Data.Value
	//	message = "Eden Farm - #OTP# adalah kode verifikasi masuk ke akun. PENTING: Mohon untuk TIDAK MENYEBARKAN kode ke orang lain atau pihak Eden Farm, demi keamanan akun."
	message = strings.ReplaceAll(message, "#OTP#", otp)

	parameter.Type = "text"
	parameter.Text = otp
	component.Type = "body"
	component.Parameters = append(component.Parameters, parameter)
	data.WAID = "62" + phoneNumber
	data.TemplateID = WhatsAppSociomileTemplateID
	data.Components = append(data.Components, component)

	jsonReq, _ := json.Marshal(data)

	request, err := http.NewRequest("POST", WhatsAppSociomileHost+"messages/send-template-message", bytes.NewBuffer(jsonReq))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("x-api-key", WhatsAppSociomileApiKey)

	response, err := client.Do(request)

	if err != nil {
		return err
	}

	req := UnmarshalBody(response)

	if req.Errors != nil {
		return errors.New(req.Errors[0].Title)
	}

	og := &model.OtpOutgoing{
		PhoneNumber:     phoneNumber,
		OTP:             otp,
		Application:     application,
		UsageType:       usageType,
		DeliveryStatus:  1,
		OtpStatus:       1,
		MessageType:     2,
		Message:         message,
		Vendor:          2,
		VendorMessageID: req.Message[0].ID,
		CreatedAt:       time.Now(),
	}
	_, err = s.RepositoryOTPOutgoing.Create(ctx, og)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return err
	}

	defer response.Body.Close()

	return err
}

func RandomInt(n int, phoneNumber string) string {
	var otp string
	var total int
	var s = global.Setup.Common

	const letters = "0123456789"
	ret := make([]byte, n)
	db := s.Database.Read
	for i := 0; i < n; i++ {
		num, err := rand2.Int(rand2.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return ""
		}
		ret[i] = letters[num.Int64()]
	}
	otp = string(ret)
	db.Raw("SELECT COUNT(phone_number) FROM otp_outgoing WHERE phone_number = ? AND otp = ?", phoneNumber, otp).QueryRow(&total)
	if total > 0 {
		RandomInt(n, phoneNumber)
	}

	return otp
}

type GetResponsePost struct {
	Message []Message `json:"messages"`
	Errors  []Error   `json:"errors"` // For Handling Errors WA OTP
}

type Error struct {
	Title string `json:"title"`
}

type Message struct {
	To        string `json:"to"`
	MessageID string `json:"messageId"`
	Status    status `json:"status"`
	ID        string `json:"id"` // For WA Message ID
}
type status struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func UnmarshalBody(r *http.Response) GetResponsePost {
	var body []byte
	var grp GetResponsePost
	var err error
	if r != nil {
		body, err = ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(body, &grp)
		if err != nil {
			panic(err)
		}
	}

	return grp
}

type sendOtpSMSVIRO struct {
	From              string        `json:"from"`
	Destinations      []destination `json:"destinations"`
	Text              string        `json:"text"`
	NotifyURL         string        `json:"notifyUrl"`
	NotifyContentType string        `json:"notifyContentType"`
}

type destination struct {
	To string `json:"to"`
}
type message struct {
	Messages []sendOtpSMSVIRO `json:"messages"`
}

func (s *AuthService) SendOTPSMSVIRO(ctx context.Context, phoneNumber string, application int, usageType int) error {
	//phoneNumber = "85695601769"
	var client = &http.Client{}
	var data message
	var sendOtpSMSVIRO sendOtpSMSVIRO
	var message string
	var desti destination
	//var env env.Provider

	var m = global.Setup.Common

	SMSViroNotifyUrl := m.Env.GetString("SMS_VIRO.NOTIFY_URL")
	SecretKeySMSViro := m.Env.GetString("SMS_VIRO.SECRET_KEY")

	ocl := m.Env.GetInt("OTP.CODE_LENGTH")
	otp := RandomInt(ocl, phoneNumber)
	messageConfig, err := m.Client.ConfigurationServiceGrpc.GetConfigAppDetail(ctx, &configuration_service.GetConfigAppDetailRequest{
		Attribute: "otp_msg",
	})
	message = messageConfig.Data.Value
	message = strings.ReplaceAll(message, "#OTP#", otp)

	sendOtpSMSVIRO.NotifyURL = SMSViroNotifyUrl
	sendOtpSMSVIRO.From = "Eden Farm"
	desti.To = "62" + phoneNumber
	sendOtpSMSVIRO.Text = message
	sendOtpSMSVIRO.NotifyContentType = "application/json"
	sendOtpSMSVIRO.Destinations = append(sendOtpSMSVIRO.Destinations, desti)
	data.Messages = append(data.Messages, sendOtpSMSVIRO)
	jsonReq, _ := json.Marshal(data)

	request, err := http.NewRequest("POST", "https://api.smsviro.com/restapi/sms/1/text/advanced", bytes.NewBuffer(jsonReq))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", SecretKeySMSViro)

	response, err := client.Do(request)
	req := UnmarshalBody(response)

	og := &model.OtpOutgoing{
		PhoneNumber:     phoneNumber,
		OTP:             otp,
		Application:     application,
		UsageType:       usageType,
		DeliveryStatus:  1,
		OtpStatus:       1,
		MessageType:     1,
		Message:         message,
		Vendor:          1,
		VendorMessageID: req.Message[0].MessageID,
		CreatedAt:       time.Now(),
	}

	_, err = s.RepositoryOTPOutgoing.Create(ctx, og)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return err
	}

	defer response.Body.Close()

	return nil
}

// Validate implement validation.Requests interfaces.
func (s *AuthService) VerifyOtp(ctx echo.Context, req dto.VerifyRegistRequest) error {
	var err error
	var otp string

	//to get expired period to limit how many seconds value from database
	appConfig, err := s.opt.Client.ConfigurationServiceGrpc.GetConfigAppDetail(ctx.Request().Context(), &configuration_service.GetConfigAppDetailRequest{
		Attribute: "cma_failed_login_retry_time",
	})
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("configuration", "config app")
		return err
	}
	period, _ := strconv.Atoi(appConfig.Data.Value)
	expPeriod := time.Duration(period) * time.Second

	appConfig, err = s.opt.Client.ConfigurationServiceGrpc.GetConfigAppDetail(ctx.Request().Context(), &configuration_service.GetConfigAppDetailRequest{
		Attribute: "otp_max_attempt_count",
	})
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("configuration", "config app")
		return err
	}
	// //to get retry count value from database
	maxRetry, _ := strconv.Atoi(appConfig.Data.Value)

	redisKey := "verify_otp_" + "_" + req.Data.PhoneNumber
	OTPValidDuration := s.opt.Env.GetString("OTP.VALID_DURATION") //env.GetString("OTP_VALID_DURATION", "600")

	_, err = s.RepositoryWhiteListLogin.GetDetail(ctx.Request().Context(), 0, req.Data.PhoneNumber, req.Data.OTP)
	if err == nil {
		return err
	}
	db := s.opt.Database.Read
	err = db.Raw("SELECT otp FROM otp_outgoing oo WHERE oo.phone_number = ? "+
		"AND oo.created_at BETWEEN NOW()-INTERVAL ? second "+
		"AND NOW() "+
		"AND oo.otp_status = 1 "+
		"AND oo.usage_type = 2 ORDER BY oo.created_at DESC LIMIT 1", req.Data.PhoneNumber, OTPValidDuration).QueryRow(&otp)
	if err != nil {
		if !global.IsLimit(redisKey, expPeriod, maxRetry) {
			err = edenlabs.ErrorValidation("otp.invalid", "Terlalu banyak percobaan untuk verifikasi otp dengan nomor ini  - Silahkan coba lagi dalam "+strconv.Itoa(period)+" detik")
			return err
		}
		err = edenlabs.ErrorValidation("otp.invalid", "Kode OTP tidak valid. Silakan masukkan kode OTP yang benar.")
		return err
	}
	if otp != "" {
		if otp != req.Data.OTP {
			if !global.IsLimit(redisKey, expPeriod, maxRetry) {
				err = edenlabs.ErrorValidation("otp.invalid", "Terlalu banyak percobaan untuk verifikasi otp dengan nomor ini  - Silahkan coba lagi dalam "+strconv.Itoa(period)+" detik")
				return err
			}
			err = edenlabs.ErrorValidation("otp.invalid", "Kode OTP tidak valid. Silakan masukkan kode OTP yang benar.")
			// o.Failure("otp.invalid", "Kode OTP tidak valid. Silakan masukkan kode OTP yang benar.")
		}
	} else {
		if !global.IsLimit(redisKey, expPeriod, maxRetry) {
			err = edenlabs.ErrorValidation("otp.invalid", "Terlalu banyak percobaan untuk verifikasi otp dengan nomor ini  - Silahkan coba lagi dalam "+strconv.Itoa(period)+" detik")
			return err
		}
		err = edenlabs.ErrorValidation("otp.invalid", "Kode OTP sudah tidak berlaku. Silakan coba kirim ulang.")
		// o.Failure("otp.invalid", "Kode OTP sudah tidak berlaku. Silakan coba kirim ulang.")
	}
	if !global.IsMaxLimit(redisKey, maxRetry) {
		err = edenlabs.ErrorValidation("otp.invalid", "Terlalu banyak percobaan untuk verifikasi otp dengan nomor ini  - Silahkan coba lagi dalam "+strconv.Itoa(period)+" detik")
		return err
	}

	return err
}

func (s *AuthService) DeleteAccount(ctx echo.Context, customerID int64, customerIdGP string) (e error) {
	c := ctx.Request().Context()
	c, span := s.opt.Trace.Start(c, "AuthService.DeleteAccount")

	// //check sales order
	salesOrder, err := s.opt.Client.SalesServiceGrpc.GetSalesOrderList(c, &sales_service.GetSalesOrderListRequest{
		CustomerId:   customerID,
		Status:       []int32{1, 5, 6, 7, 8, 9, 10, 11},
		Limit:        1,
		Offset:       0,
		CustomerCode: customerIdGP,
	})
	if err == nil {
		//err = errors.New("Kamu masih memiliki transaksi yang belum selesai.")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	if salesOrder.Data != nil {
		err = errors.New("Kamu masih memiliki transaksi yang belum selesai.")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	userCustomer := &model.UserCustomer{
		CustomerID:   customerID,
		CustomerIDGP: customerIdGP,
	}
	userCustomer, err = s.RepositoryUserCustomer.GetDetail(c, userCustomer)
	if err != nil {
		err = errors.New("Nomor yang anda masukkan tidak terdaftar.")
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	userCustomer.Status = 3

	err = s.RepositoryUserCustomer.Update(c, userCustomer)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	//update status merchant here
	customer, err := s.opt.Client.BridgeServiceGrpc.GetCustomerGPDetail(ctx.Request().Context(), &bridge_service.GetCustomerGPDetailRequest{
		Id: customerIdGP,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	// customer.Data.Status = 3

	// customers, err := s.opt.Client.BridgeServiceGrpc.UpdateCustomer(ctx.Request().Context(), &bridge_service.UpdateCustomerRequest{
	// 	Data: customer.Data,
	// })
	// if customers.Code != 200 {

	// }
	//update branch here
	address, err := s.opt.Client.BridgeServiceGrpc.GetAddressGPDetail(ctx.Request().Context(), &bridge_service.GetAddressGPDetailRequest{
		Id: customer.Data[0].Adrscode[0].Adrscode,
	})
	address.Data[0].Inactive = 1
	for _, v := range address.Data {
		_, err = s.opt.Client.BridgeServiceGrpc.UpdateAddress(ctx.Request().Context(), &bridge_service.UpdateAddressRequest{
			Custnmbr:                v.Custnmbr,
			Custname:                v.Custname,
			Adrscode:                v.Adrscode,
			Slprsnid:                v.Slprsnid,
			Shipmthd:                v.Shipmthd,
			Taxschid:                v.Taxschid,
			Cntcprsn:                v.Cntcprsn,
			AddresS1:                v.AddresS1,
			AddresS2:                v.AddresS2,
			AddresS3:                v.AddresS3,
			Country:                 v.Country,
			City:                    v.City,
			State:                   v.State,
			Zip:                     v.Zip,
			PhonE1:                  v.PhonE1,
			PhonE2:                  v.PhonE2,
			PhonE3:                  v.PhonE3,
			CCode:                   v.CCode,
			Locncode:                v.Locncode,
			Salsterr:                v.Salsterr,
			UserdeF1:                v.UserdeF1,
			UserdeF2:                v.UserdeF2,
			ShipToName:              v.ShipToName,
			GnL_Administrative_Code: v.GnL_Administrative_Code,
			GnL_Archetype_ID:        v.GnL_Archetype_ID,
			GnL_Longitude:           v.GnL_Longitude,
			GnL_Latitude:            v.GnL_Latitude,
			GnL_Address_Note:        v.GnL_Address_Note,
			Inactive:                1,
			Crusrid:                 v.Crusrid,
			Creatddt:                v.Creatddt,
			Mdfusrid:                v.Mdfusrid,
			Modifdt:                 v.Modifdt,
			TypeAddress:             v.TypeAddress,
			Param:                   "",
			Fax:                     "",
			Upzone:                  "",
		})
	}
	//update prospect customer here
	prosCust, err := s.opt.Client.CrmServiceGrpc.GetProspectiveCustomerList(ctx.Request().Context(), &crm_service.GetProspectiveCustomerListRequest{
		CustomerId: customer.Data[0].Custnmbr,
	})
	for _, v := range prosCust.Data {
		_, err = s.opt.Client.CrmServiceGrpc.DeleteProspectiveCustomer(ctx.Request().Context(), &crm_service.DeleteProspectiveCustomerRequest{
			Id: v.Id,
		})
	}

	return
}
