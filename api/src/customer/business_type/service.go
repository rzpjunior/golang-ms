// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package business_type

// import (
// 	"git.edenfarm.id/project-version2/datamodel/model"
// 	"git.edenfarm.id/project-version2/api/util"
// )

// // Save : function to save data requested into database
// func Save(r createRequest) (u *model.PaymentMethod, e error) {
// 	r.Code, e = util.GenerateCode(r.Code, "payment_method")
// 	if e == nil {
// 		u = &model.PaymentMethod{
// 			Code:   r.Code,
// 			Name:   r.Name,
// 			Note:   r.Note,
// 			Status: int8(1),
// 		}

// 		e = u.Save()
// 	}

// 	return u, e
// }
