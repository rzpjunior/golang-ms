package handler

import (
	"strconv"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/service"
	"github.com/labstack/echo/v4"
	"github.com/ulule/limiter/v3"
)

type AuthHandler struct {
	Option global.HandlerOptions
}

// URLMapping declare endpoint with handler function.
func (h *AuthHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	cMiddleware := middleware.NewMiddleware()

	//r.POST("/login", h.Login)
	r.POST("", h.Login, cMiddleware.Authorized("public"))
	r.POST("/reqotp", h.requestOtp, cMiddleware.Authorized("public"))
	r.POST("/verifyOtpRegist", h.verifyOtpRegist, cMiddleware.Authorized("public"))
	r.POST("/checkNumberLogin", h.checkNumberLogin, cMiddleware.Authorized("public"))
	r.POST("/session", h.session, cMiddleware.Authorized("private"))
	r.POST("/signOut", h.deleteSignOut, cMiddleware.Authorized("private"))
	r.POST("/delete-account", h.deleteAccount, cMiddleware.Authorized("private"))
}

func (h AuthHandler) Login(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	timezone := ctx.Request().Header.Get("Timezone")

	var req dto.LoginRequest
	login := service.ServiceAuth()

	if timezone != "" {
		req.Timezone = timezone
	} else {
		req.Timezone = "Asia/Jakarta"
	}
	options := limiter.Options{
		TrustForwardHeader: true,
	}
	ipRequester := limiter.GetIP(ctx.Request(), options).String()

	req.Data.IPAddress = ipRequester
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = login.Login(ctx.Request().Context(), req)

	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h AuthHandler) requestOtp(c echo.Context) (e error) {
	var (
		req         dto.RequestGetOTP
		application int
		response    dto.ResponseSendOtp
	)
	login := service.ServiceAuth()

	ctx := c.(*edenlabs.Context)

	options := limiter.Options{
		TrustForwardHeader: true,
	}
	ipRequester := limiter.GetIP(ctx.Request(), options).String()

	req.Data.IPAddress = ipRequester

	if e = ctx.Bind(&req); e != nil {
		return ctx.Serve(e)
	}
	// if req.Data.WhiteListLogin.ID != 0 {
	// 	return ctx.Serve(e)
	// }

	if req.Platform == "orca" {
		application = 2
	} else if req.Platform == "mantis" {
		application = 3
	}

	if _, e := login.ResponseOTP(ctx.Request().Context(), req); e != nil {
		return ctx.Serve(e)
	}

	t, _ := strconv.Atoi(req.Data.Type)

	if req.Data.OtpType == "0" || req.Data.OtpType == "2" {
		response.OtpName = "WhatsApp"
		ctx.ResponseData = "response"
		e = login.SendOTPWASociomile(ctx.Request().Context(), req.Data.PhoneNumber, application, t)
	}

	// if otp type = 1 or invalid phone number for WA, direct to sms otp
	if req.Data.OtpType == "1" || e != nil {
		response.OtpName = "SMS"
		ctx.ResponseData = "response"
		e = login.SendOTPSMSVIRO(ctx.Request().Context(), req.Data.PhoneNumber, application, t)
	}
	return ctx.Serve(e)
}

func (h AuthHandler) verifyOtpRegist(c echo.Context) (e error) {
	ctx := c.(*edenlabs.Context)
	var req dto.VerifyRegistRequest
	login := service.ServiceAuth()

	options := limiter.Options{
		TrustForwardHeader: true,
	}
	ipRequester := limiter.GetIP(ctx.Request(), options).String()

	req.Data.IPAddress = ipRequester

	if e = ctx.Bind(&req); e == nil {
		if e := login.VerifyOtp(ctx, req); e != nil {
			return ctx.Serve(e)
		}
		ctx.Data("success")
	}
	return ctx.Serve(e)
}

func (h AuthHandler) checkNumberLogin(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)
	var req *dto.CheckPhoneNumberRequest
	login := service.ServiceAuth()

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = login.CheckPhoneNumber(c, req)

	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}
	return ctx.Serve(err)
}

func (h AuthHandler) session(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)
	var req dto.CheckSessionRequest
	var s *dto.SessionDataCustomer
	login := service.ServiceAuth()

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	if s, err = login.Session(ctx); err != nil {
		return ctx.Serve(err)
	}
	ctx.Data(s)

	return ctx.Serve(err)
}

func (h AuthHandler) deleteSignOut(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)
	var req dto.SignOutRequest
	login := service.ServiceAuth()

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	if req.Session, err = service.CustomerSession(ctx); err == nil {
		customerID, _ := strconv.Atoi(req.Session.Customer.ID)
		login.SignOut(ctx, int64(customerID))
	}

	return ctx.Serve(err)
}
func (h AuthHandler) deleteAccount(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)
	var req dto.RequestDeleteAccount
	login := service.ServiceAuth()

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}
	if req.Session, err = service.CustomerSession(ctx); err == nil {
		customerID, _ := strconv.Atoi(req.Session.Customer.ID)
		err = login.DeleteAccount(ctx, int64(customerID), req.Session.Customer.Code)
	}

	return ctx.Serve(err)
}
