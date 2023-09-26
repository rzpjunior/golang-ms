package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/repository"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
	"github.com/labstack/echo/v4"
)

func (m *Middleware) Authorized(permissions string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) (err error) {
			var e error
			var valueMaintenance int
			var minVersion int64
			var platform string
			var appVersion string
			var um *model.UserCustomer
			var m = global.Setup.Common
			fmt.Println(um, platform, appVersion, valueMaintenance, minVersion)
			if um, platform, appVersion, e = validTokenLogin(ctx, permissions); e == nil {
				valueMaintenance, err := m.Client.ConfigurationServiceGrpc.GetConfigAppDetail(ctx.Request().Context(), &configuration_service.GetConfigAppDetailRequest{
					Attribute: "maintenance_" + platform,
				})
				if err != nil {
					return err
				}
				vm, err := strconv.Atoi(valueMaintenance.Data.Value)
				if vm == 1 && "maintenance_"+platform == "maintenance_orca" {
					e = echo.NewHTTPError(http.StatusServiceUnavailable, "Server Maintenance")
				}
				if vm == 1 && "maintenance_"+platform == "maintenance_mantis" {
					e = echo.NewHTTPError(http.StatusServiceUnavailable, "Server Maintenance")
				}
				if permissions != "public" {
					if um.ForceLogout != 2 {
						e = echo.NewHTTPError(http.StatusRequestTimeout, "Force Logout")
					}
				}
				appVersionInt, _ := strconv.ParseInt(appVersion, 10, 64)

				//check minimum version

				minVersion, err := m.Client.ConfigurationServiceGrpc.GetConfigAppDetail(ctx.Request().Context(), &configuration_service.GetConfigAppDetailRequest{
					Attribute: platform + "_min_version",
				})
				mv, err := strconv.Atoi(minVersion.Data.Value)

				if appVersionInt < int64(mv) {
					e = echo.NewHTTPError(http.StatusUpgradeRequired, "Force Update App")
				}
			}
			if e == nil {
				return next(ctx)
			}
			return e
		}
	}
}

func validTokenLogin(c echo.Context, state string) (*model.UserCustomer, string, string, error) {
	var e error
	var s = repository.NewUserCustomerRepository()
	var m = global.Setup.Common
	var uc *model.UserCustomer
	auth := c.Request().Header.Get("Authorization")
	platform := c.Request().Header.Get("Platform")
	appVersion := c.Request().Header.Get("AppVersion")
	if state != "public" && state != "private" {
		e = echo.NewHTTPError(http.StatusRequestTimeout, "invalid content")
	} else {
		if state == "private" {
			if auth != "" && len(auth) > 6 {
				bearer := auth[:strings.IndexByte(auth, ' ')]
				if bearer != "Bearer" {
					e = echo.NewHTTPError(http.StatusRequestTimeout, "invalid or expired jwt token")
				} else {
					tokenLogin := strings.Replace(auth, "Bearer ", "", 1) // get token from Authorization
					userCustomer := &model.UserCustomer{LoginToken: tokenLogin, Status: 1}
					userCustomer, err := s.GetDetail(c.Request().Context(), userCustomer)
					if err != nil {
						e = echo.NewHTTPError(http.StatusRequestTimeout, "invalid or expired jwt token")
						m.Logger.AddMessage(log.ErrorLevel, e)
						return nil, "", "", err

					}
					_, err = m.Client.BridgeServiceGrpc.GetCustomerDetail(c.Request().Context(), &bridge_service.GetCustomerDetailRequest{
						Id:     userCustomer.CustomerID,
						Status: 1,
					})
					if err != nil {
						e = echo.NewHTTPError(http.StatusRequestTimeout, "invalid or expired jwt token")
						m.Logger.AddMessage(log.ErrorLevel, e)
						return nil, "", "", err
					}
					uc = &model.UserCustomer{
						ID:            userCustomer.ID,
						Code:          userCustomer.Code,
						CustomerID:    userCustomer.CustomerID,
						FirebaseToken: userCustomer.FirebaseToken,
						LoginToken:    userCustomer.LoginToken,
						ForceLogout:   userCustomer.ForceLogout,
						CustomerIDGP:  userCustomer.CustomerIDGP,
					}
				}
			} else {
				e = echo.NewHTTPError(http.StatusRequestTimeout, "invalid or expired jwt token")
			}

			if appVersion == "" {
				e = echo.NewHTTPError(http.StatusForbidden, "invalid or expired jwt token")
			}
		}
	}

	return uc, platform, appVersion, e
}
