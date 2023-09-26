// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package auth

import (
	"errors"
	"strings"
	"time"

	"git.edenfarm.id/project-version2/datamodel/model"

	"fmt"

	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/cuxs/dbredis"
	"git.edenfarm.id/cuxs/orm"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// SessionData struktur data current user logged in.
type SessionData struct {
	Token string       `json:"token"`
	Staff *model.Staff `json:"staff"`
}

// Login mendapatkan session data dari model user.
// user diasumsikan sudah valid untuk login.
// jadi disini tidak ada validasi untuk login, hanya
// untuk mendapatkan session data
func Login(user *model.User, dashboardToken ...string) (sd *SessionData, e error) {
	//application menu
	if sd, e = StartSession(user.ID); e == nil {
		// update last login dari user tersebut
		user.LastLoginAt = time.Now()
		user.ForceLogout = 2
		if len(dashboardToken) > 0 {
			user.DashboardNotifToken = dashboardToken[0]
		}
		user.Save("LastLoginAt", "ForceLogout", "PickingNotifToken", "DashboardNotifToken")

		return sd, nil
	}

	return nil, e
}

// StartSession mendapatkan data user entity dengan token
// untuk menandakan session user yang sedang login.
func StartSession(userID int64, token ...string) (sd *SessionData, e error) {
	sd = new(SessionData)
	// buat token baru atau menggunakan yang sebelumnya
	if len(token) == 0 {
		sd.Token = cuxs.JwtToken("id", userID)
	} else {
		sd.Token = token[0]
	}

	// membaca data user terlebih dahulu
	sd.Staff = &model.Staff{User: &model.User{ID: userID}}
	if e = sd.Staff.Read("User"); e == nil {
		sd.Staff.StaffID = userID + 56
		sd.Staff.User.Read("ID")
		sd.Staff.Role.Read("ID")
	}

	return sd, e
}

// UserSession mendapatkan session data dari user yang mengirimkan request.
//func UserSession(ctx *cuxs.Context) (*SessionData, error) {
func UserSession(ctx *cuxs.Context) (*SessionData, error) {
	if u := ctx.Get("user"); u != nil {
		c := u.(*jwt.Token).Claims.(jwt.MapClaims)
		var userID int64

		// id adalah user id
		if c["id"] != nil {
			userID = int64(c["id"].(float64))
		}

		// memakai token sebelumnya
		token := ctx.Get("user").(*jwt.Token).Raw

		return StartSession(userID, token)
	}

	return nil, errors.New(strings.Title("invalid jwt token"))
}

// isAuthorized cek apakah userID mempunyai privileges untuk mengakses
// module berdasarkan nama yang diberikan.
func isAuthorized(userID int64, module string) bool {
	var c int64

	key := fmt.Sprintf("SELECT COUNT(id) FROM user_permission WHERE user_id = %d AND permission_value = '%s'", userID, module)

	if !dbredis.Redis.CheckExistKey(key) {
		o := orm.NewOrm()
		o.Using("read_only")
		o.Raw(key).QueryRow(&c)
		dbredis.Redis.SetCache(key, c, 24*time.Hour)
	} else {
		dbredis.Redis.GetCache(key, &c)
	}

	return c != 0
}

// login to Field Purchaser Apps
func LoginFieldPurchaser(r SignInFieldPurchaserRequest) (sd *SessionData, e error) {
	user := r.User

	// check if role is Sourcing Manager or Sourcing Admin
	if !r.IsRoleValid {
		return nil, echo.NewHTTPError(401, "You're not registered in Sourcing Division, please contact IT Support or HRD")
	}

	//application menu
	if sd, e = StartSession(user.ID); e != nil {
		return nil, e
	}

	// update last login dari user tersebut
	user.LastLoginAt = time.Now()
	user.ForceLogout = 2
	user.Save("LastLoginAt", "ForceLogout", "PurchaserNotifToken")

	return sd, nil
}
