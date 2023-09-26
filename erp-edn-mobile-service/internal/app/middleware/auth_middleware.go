package middleware

import (
	"net/http"
	"strings"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/middlewarex"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/account_service"
	"github.com/labstack/echo/v4"
)

func (m *Middleware) Authorized(permissions ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) (err error) {
			err = middlewarex.Authorization(ctx, m.Option.Logger, m.Option.Config.Jwt.Key, permissions)
			if err != nil {
				return err
			}

			auth := ctx.Request().Header.Get("Authorization")
			tokenLogin := strings.Replace(auth, "Bearer ", "", 1) // get token from Authorization

			// update token
			var userResponse *account_service.GetUserDetailResponse
			userResponse, err = m.Option.Client.AccountServiceGrpc.GetUserByEdnAppLoginToken(ctx.Request().Context(), &account_service.GetUserByEdnAppLoginTokenRequest{
				EdnappLoginToken: tokenLogin,
			})
			if err != nil {
				m.Option.Logger.AddMessage(log.ErrorLevel, err)
				edenlabs.ErrorRpcNotFound("account", "GetUserByEdnAppLoginToken")
				err = echo.NewHTTPError(http.StatusUnauthorized, "Force Logout")
				return
			}

			if userResponse.Data.ForceLogout != 2 {
				err = echo.NewHTTPError(http.StatusUnauthorized, "Force Logout")
			}
			return next(ctx)
		}
	}
}
