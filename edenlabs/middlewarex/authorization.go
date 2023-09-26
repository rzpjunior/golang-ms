package middlewarex

import (
	"context"
	"net/http"
	"regexp"
	"strings"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/constants"
	jwtx "git.edenfarm.id/edenlabs/edenlabs/jwt"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func Authorization(ctx echo.Context, lgr *log.Logger, signKey string, permissions []string) (err error) {
	authorization := ctx.Request().Header.Get("Authorization")
	if authorization != "" {
		var match bool
		match, err = regexp.MatchString("^Bearer .+", authorization)
		if err != nil || !match {
			lgr.AddMessage(log.DebugLevel, echo.ErrUnauthorized).Print()
			return echo.ErrUnauthorized
		}

		j := jwtx.NewJWT([]byte(signKey))

		tokenStr := strings.Split(authorization, " ")

		var token *jwt.Token
		token, err = j.Parse(tokenStr[1])
		if err != nil {
			lgr.AddMessage(log.DebugLevel, echo.ErrUnauthorized).Print()
			return echo.ErrUnauthorized
		}

		var claims *jwtx.UserClaim
		var ok bool
		claims, ok = token.Claims.(*jwtx.UserClaim)
		if !ok {
			lgr.AddMessage(log.DebugLevel, echo.ErrUnauthorized).Print()
			return echo.ErrUnauthorized
		}

		expiresAt := claims.ExpiresAt
		if expiresAt <= time.Now().Unix() {
			lgr.AddMessage(log.DebugLevel, echo.ErrBadRequest).Print()
			ctx.JSON(http.StatusUnauthorized, edenlabs.FormatResponse{
				Code:    http.StatusUnauthorized,
				Status:  "failure",
				Message: "Your token is expired",
			})
			return echo.ErrUnauthorized
		}

		if len(permissions) != 0 {
			isExisted := false
			for _, permission := range permissions {
				for _, userPermission := range claims.Permissions {
					if permission == userPermission {
						isExisted = true
					}
				}
			}
			if !isExisted {
				lgr.AddMessage(log.DebugLevel, echo.ErrBadRequest).Print()
				ctx.JSON(http.StatusForbidden, edenlabs.FormatResponse{
					Code:    http.StatusForbidden,
					Status:  "failure",
					Message: "You don't have permission for this",
				})
				return echo.ErrForbidden
			}
		}
		ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), constants.KeyToken, token)))
		ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), constants.KeyUserID, claims.UserID)))
		ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), constants.KeyTimezone, claims.Timezone)))
		return
	}
	return echo.ErrForbidden
}

func AuthorizationCourierApp(ctx echo.Context, lgr *log.Logger, signKey string) (err error) {
	authorization := ctx.Request().Header.Get("Authorization")
	if authorization == "" {
		return echo.ErrForbidden
	}

	match, err := regexp.MatchString("^Bearer .+", authorization)
	if err != nil || !match {
		lgr.AddMessage(log.DebugLevel, echo.ErrUnauthorized).Print()
		return echo.ErrUnauthorized
	}

	j := jwtx.NewJWT([]byte(signKey))

	tokenStr := strings.Split(authorization, " ")

	token, err := j.ParseCourier(tokenStr[1])
	if err != nil {
		lgr.AddMessage(log.DebugLevel, echo.ErrUnauthorized).Print()
		return echo.ErrUnauthorized
	}

	claims, ok := token.Claims.(*jwtx.UserCourierClaim)
	if !ok {
		lgr.AddMessage(log.DebugLevel, echo.ErrUnauthorized).Print()
		return echo.ErrUnauthorized
	}

	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), constants.KeyToken, token)))
	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), constants.KeyCourierID, claims.CourierID)))
	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), constants.KeySiteID, claims.SiteID)))
	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), constants.KeyTimezone, claims.Timezone)))

	return
}

func AuthorizationHelperMobileApp(ctx echo.Context, lgr *log.Logger, signKey string) (err error) {
	authorization := ctx.Request().Header.Get("Authorization")
	if authorization == "" {
		return echo.ErrForbidden
	}

	match, err := regexp.MatchString("^Bearer .+", authorization)
	if err != nil || !match {
		lgr.AddMessage(log.DebugLevel, echo.ErrUnauthorized).Print()
		return echo.ErrUnauthorized
	}

	j := jwtx.NewJWT([]byte(signKey))

	tokenStr := strings.Split(authorization, " ")

	token, err := j.ParseHelperMobile(tokenStr[1])
	if err != nil {
		lgr.AddMessage(log.DebugLevel, echo.ErrUnauthorized).Print()
		return echo.ErrUnauthorized
	}

	claims, ok := token.Claims.(*jwtx.UserHelperMobileClaim)
	if !ok {
		lgr.AddMessage(log.DebugLevel, echo.ErrUnauthorized).Print()
		return echo.ErrUnauthorized
	}

	platform := claims.Platform
	if platform != "helper-mobile" {
		lgr.AddMessage(log.DebugLevel, echo.ErrBadRequest).Print()
		ctx.JSON(http.StatusUnauthorized, edenlabs.FormatResponse{
			Code:    http.StatusForbidden,
			Status:  "failure",
			Message: "You don't have permission for this",
		})
		return echo.ErrUnauthorized
	}

	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), constants.KeyToken, token)))
	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), constants.KeyUserID, claims.UserID)))
	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), constants.KeySiteID, claims.SiteId)))
	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), constants.KeyTimezone, claims.Timezone)))

	return
}
