// Copyright 2016 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package edenlabs

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/validation"
	"github.com/labstack/echo/v4"
)

type (
	// Custom echo binder
	Binder struct{}

	// Validator custom validation object.
	BinderValidator struct {
		once      sync.Once
		validator *validation.Validation
	}
)

// Set binding validator using cuxs validation.
var bindValidator = &BinderValidator{}

// ValidateStruct evaluate an object,
// will run validation request if the object
// is implementing validatonRequests.
func (v *BinderValidator) validate(obj interface{}) error {
	v.lazyinit()

	var o *validation.Output
	if vr, ok := obj.(validation.Request); ok {
		o = v.validator.Request(vr)
	} else {
		o = v.validator.Struct(obj)
	}

	if !o.Valid {
		return o
	}
	return nil
}

// lazyinit initialing validator instances for one of time only.
func (v *BinderValidator) lazyinit() {
	v.once.Do(func() {
		v.validator = validation.New()
	})
}

// Bind is decode request body and injecting into interfaces,
// We only accept json data type other type will return error bad requests.
// Also automaticly validate data with interfaces.
func (b Binder) Bind(i interface{}, ctx echo.Context) (err error) {
	bindValidator.lazyinit()
	req := ctx.Request()
	ctype := req.Header.Get(echo.HeaderContentType)
	if strings.HasPrefix(ctype, echo.MIMEApplicationJSON) {
		if err = json.NewDecoder(req.Body).Decode(i); err != nil && err != io.EOF {
			if ute, ok := err.(*json.UnmarshalTypeError); ok {
				return echo.NewHTTPError(http.StatusUnprocessableEntity, fmt.Sprintf("Invalid payload, please check your payload : expected=%v, got=%v, offset=%v", ute.Type, ute.Value, ute.Offset))
			} else if se, ok := err.(*json.SyntaxError); ok {
				return echo.NewHTTPError(http.StatusUnprocessableEntity, fmt.Sprintf("Syntax error: offset=%v, error=%v", se.Offset, se.Error()))
			} else if _, ok := err.(*time.ParseError); ok {
				return echo.NewHTTPError(http.StatusUnprocessableEntity, "Invalid datetime format value, please check your datetime value")
			} else {
				return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
			}
		}
		return bindValidator.validate(i)
	}
	return echo.ErrUnsupportedMediaType
}
