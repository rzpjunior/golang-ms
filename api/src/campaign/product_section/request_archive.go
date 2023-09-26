// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package product_section

import (
	"time"

	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// archiveRequest : struct to hold request data
type archiveRequest struct {
	ID   int64  `json:"-" valid:"required"`
	Note string `json:"note" valid:"required"`

	ProductSection *model.ProductSection `json:"-"`

	Session *auth.SessionData
}

// Validate : function to validate request data
func (r *archiveRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	if r.ProductSection, err = repository.ValidProductSection(r.ID); err != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("product section"))
	}

	currentTime := time.Now()
	if currentTime.After(r.ProductSection.EndAt) {
		o.Failure("id.invalid", util.ErrorActive("product section"))
	}

	// get value_int status archive in glossary
	statusGlossary, err := repository.GetGlossaryMultipleValue("table", "product_section", "attribute", "status", "value_name", "Archived")
	if err != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("product section"))
	}
	r.ProductSection.Status = statusGlossary.ValueInt

	return o
}

// Messages : function to return error validation messages
func (r *archiveRequest) Messages() map[string]string {
	return map[string]string{
		"id.required":   util.ErrorInputRequired("id"),
		"note.required": util.ErrorInputRequired("note"),
	}
}
