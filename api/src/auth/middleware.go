// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package auth

import (
	"net/http"
	"strconv"
	"time"

	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/cuxs/dbredis"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"github.com/labstack/echo/v4"
)

// Authorized is middleware that will check if user has authorize in endpoint
// note: pada handler UrlMapping ditambahkan parameter auth.Authorized() dan application modulnya
// e.g:  r.Post("",h.auth, auth.Authorized("dashboard"))
func Authorized(permission string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return cuxs.Authorized()(func(context echo.Context) error {
			var e error
			var sd *SessionData
			var maintenance int

			if sd, e = UserSession(cuxs.NewContext(context)); e == nil {
				o := orm.NewOrm()
				o.Using("read_only")

				o.Raw("SELECT value from config_app where attribute = 'maintenance_dino'").QueryRow(&maintenance)
				if maintenance == 1 {
					return echo.NewHTTPError(503, "server maintenance")
				}

				if permission != "" && permission != "privilege" {
					sd.Staff.Role.Read("ID")
					sd.Staff.Role.Division.Read("ID")
					if isAuthorized(sd.Staff.User.ID, permission) && sd.Staff.User.Email != "superadmin" {
						return next(context)
					} else if sd.Staff.User.Email == "superadmin" {
						return next(context)
					} else if sd.Staff.Role.Code == "ROL0023" || sd.Staff.Role.Code == "ROL0022" {
						return echo.ErrForbidden
					}
				} else if permission == "privilege" {
					return next(context)
				}

			}
			return echo.ErrForbidden
		})
	}
}

func AuthorizedMobileUniversal() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return cuxs.Authorized()(func(context echo.Context) error {
			var e error
			var sd *SessionData
			o := orm.NewOrm()
			o.Using("read_only")

			if sd, e = UserSession(cuxs.NewContext(context)); e == nil {
				sd.Staff.Role.Read("ID")
				sd.Staff.Role.Division.Read("ID")
				sd.Staff.User.Read("ID")

				leadPickerRoleCode, _ := repository.GetConfigApp("attribute", "lead_picker_role_id")
				pickerRoleCode, _ := repository.GetConfigApp("attribute", "picker_role_id")
				if sd.Staff.Role.Division.Code == "DIV0009" {
					var count int
					o.Raw("select count(*) from user_permission up where up.user_id = ? and up.permission_value = 'pco_app'", sd.Staff.User.ID).QueryRow(&count)
					if count != 0 {
						return next(context)
					} else if sd.Staff.Role.Code == leadPickerRoleCode.Value || sd.Staff.Role.Code == pickerRoleCode.Value || sd.Staff.Role.Code == "ROL0022" || sd.Staff.Role.Code == "ROL0049" {
						if sd.Staff.User.ForceLogout == 1 {
							return echo.ErrForbidden
						} else {
							return next(context)
						}
					} else {
						return echo.ErrForbidden
					}
				} else {
					return echo.ErrForbidden

				}

			}
			return echo.ErrForbidden
		})
	}
}

// PACKING ORDER MOBILE
func AuthorizedMobile() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return cuxs.Authorized()(func(context echo.Context) error {
			var e error
			var sd *SessionData
			o := orm.NewOrm()
			o.Using("read_only")

			if sd, e = UserSession(cuxs.NewContext(context)); e == nil {
				sd.Staff.Role.Read("ID")
				sd.Staff.Role.Division.Read("ID")
				sd.Staff.User.Read("ID")

				leadPickerRoleCode, _ := repository.GetConfigApp("attribute", "lead_picker_role_id")
				if sd.Staff.Role.Division.Code == "DIV0009" {
					var count int
					o.Raw("select count(*) from user_permission up where up.user_id = ? and up.permission_value = 'pco_app'", sd.Staff.User.ID).QueryRow(&count)
					if count != 0 {
						return next(context)
					} else if sd.Staff.Role.Code == leadPickerRoleCode.Value || sd.Staff.Role.Code == "ROL0022" || sd.Staff.Role.Code == "ROL0049" {
						if sd.Staff.User.ForceLogout == 1 {
							return echo.ErrForbidden
						} else {
							return next(context)
						}
					} else {
						return echo.ErrForbidden
					}
				} else {
					return echo.ErrForbidden

				}

			}
			return echo.ErrForbidden
		})
	}
}

func AuthorizedPickerMobile() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return cuxs.Authorized()(func(context echo.Context) error {
			var e error
			var sd *SessionData
			o := orm.NewOrm()
			o.Using("read_only")

			if sd, e = UserSession(cuxs.NewContext(context)); e == nil {
				sd.Staff.Role.Read("ID")
				sd.Staff.Role.Division.Read("ID")
				sd.Staff.User.Read("ID")

				pickerRoleCode, _ := repository.GetConfigApp("attribute", "picker_role_id")
				if sd.Staff.Role.Division.Code == "DIV0009" {
					var count int
					o.Raw("select count(*) from user_permission up where up.user_id = ? and up.permission_value = 'pco_app'", sd.Staff.User.ID).QueryRow(&count)
					if count != 0 {
						return next(context)
					} else if sd.Staff.Role.Code == pickerRoleCode.Value {
						if sd.Staff.User.ForceLogout == 1 {
							return echo.ErrForbidden
						} else {
							return next(context)
						}
					} else {
						return echo.ErrForbidden
					}
				} else {
					return echo.ErrForbidden

				}

			}
			return echo.ErrForbidden
		})
	}
}

// FIELD PURCHASER MOBILE
func AuthorizedFieldPurchaserMobile() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return cuxs.Authorized()(func(context echo.Context) error {
			var e error
			var sd *SessionData
			o := orm.NewOrm()
			o.Using("read_only")

			if sd, e = UserSession(cuxs.NewContext(context)); e != nil {
				return echo.NewHTTPError(403, "Invalid JWT Token")
			}

			if e = sd.Staff.Role.Read("ID"); e != nil {
				return echo.NewHTTPError(403, "Invalid Role ID")
			}
			if e = sd.Staff.Role.Division.Read("ID"); e != nil {
				return echo.NewHTTPError(403, "Invalid Division ID")
			}
			if e = sd.Staff.User.Read("ID"); e != nil {
				return echo.NewHTTPError(403, "Invalid User ID")
			}

			if sd.Staff.Role.Division.Name != "Sourcing" {
				return echo.NewHTTPError(403, "You're not registered in Sourcing Division, please contact IT Support or HRD")
			}

			isRoleValid, err := repository.IsRoleFieldPurchaser(strconv.FormatInt(sd.Staff.Role.ID, 10))
			if err != nil {
				return echo.NewHTTPError(403, "Invalid Role ID")
			}

			if !isRoleValid {
				return echo.NewHTTPError(403, "You're not registered as Field Purchaser or Purchasing Manager, please contact IT support or HRD")
			}

			if sd.Staff.User.ForceLogout == 1 {
				return echo.NewHTTPError(403, "You've logged out, please login again")
			}

			return next(context)
		})
	}
}

func AuthorizedFarmStack() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(context echo.Context) error {
			var e error
			if e = validKey(context); e != nil {
				e = echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
			}
			if e == nil {
				return next(context)
			}

			return e
		}
	}
}

func validKey(c echo.Context) error {
	var e error
	key := c.Request().Header.Get("KEY")
	SecretKey := c.Request().Header.Get("SECRET-KEY")
	if key != "6oxTe4CJkDFIoyYi7b3qEB6kUpcDtD2uoKSejrc4JN4Q2IWjkPyCwyvGlPRxzbAb" || SecretKey != "0JBqcF24PtwJnHDqhrDR" {
		e = echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	return e
}

// untuk melakukan pengecekan apakah key yang dikirimkan untuk melimit
// akses ada di redis atau tidak.
func isLimit(key string, period time.Duration) bool {
	var c int64

	//to get retry count value from database
	configApp, _ := repository.GetConfigApp("attribute", "dino_failed_login_retry_count")
	maxRetry, _ := strconv.Atoi(configApp.Value)

	//tidak ada di redis
	if !dbredis.Redis.CheckExistKey(key) {
		dbredis.Redis.SetCache(key, c+1, period)
		return true
	}

	dbredis.Redis.GetCache(key, &c)
	if c != int64(maxRetry) {
		dbredis.Redis.SetCache(key, c+1, period)
		return true
	}
	//ada di redis
	return false
}

// untuk mengecek apakah counter di redis sudah memenuhi max retry
// tanpa mengubah value di redis
func isMaxLimit(key string, period time.Duration) bool {
	var c int64

	//to get retry count value from database
	configApp, _ := repository.GetConfigApp("attribute", "dino_failed_login_retry_count")
	maxRetry, _ := strconv.Atoi(configApp.Value)

	dbredis.Redis.GetCache(key, &c)

	//sudah max jadi ditolak untuk coba lagi
	if c == int64(maxRetry) {
		return false
	}

	return true
}

// resetCache : handler to reset cache
func resetCache(key string) (e error) {
	e = dbredis.Redis.DeleteCacheWhereLike(key)
	return e
}
