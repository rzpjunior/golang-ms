// Copyright 2016 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package edenlabs

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"git.edenfarm.id/edenlabs/edenlabs/validation"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type UsualStruct struct {
	Name string        `json:"name" valid:"required"`
	Rs   RequestStruct `json:"request_struct" valid:"required"`
}

type RequestStruct struct {
	Username string `json:"username" valid:"required|email"`
	Password string `json:"password" valid:"required"`
}

func (r RequestStruct) Validates() *validation.Output {
	o := validation.Output{Valid: false}
	o.Failure("login", "wew")
	return &o
}

func (r RequestStruct) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	if r.Password != "s3cr3t" {
		o.Valid = false
		o.Failure("login", "password not match.")
	}
	return o
}

func (r RequestStruct) Messages() map[string]string {
	return map[string]string{
		"username.required": "username cannot be null",
		"username.email":    "username is invalid",
		"password.required": "password cannot be null",
		"login":             "Login failed, ",
	}
}

func Test_Binder(t *testing.T) {
	var obj UsualStruct

	ctx := requestWithBody("POST", "/", `{}`, "text/plain")
	assert.Error(t, ctx.Bind(&obj))

	ctx = requestWithBody("POST", "/", `{"hoge": !@#@#}`, "application/json")
	assert.Error(t, ctx.Bind(&obj))

	ctx = requestWithBody("POST", "/", `{"name": 123}`, "application/json")
	assert.Error(t, ctx.Bind(&obj))
}

func Test_ValidatorStruct(t *testing.T) {
	var obj UsualStruct
	ctx := requestWithBody("POST", "/", `{"hoge": 0}`, "application/json")
	assert.Error(t, ctx.Bind(&obj))
}

func Test_ValidatorRequest(t *testing.T) {
	var ors RequestStruct
	ctx := requestWithBody("POST", "/", `{"username": "x@x.com", "password": "n0ts3cr3t"}`, "application/json")
	assert.Error(t, ctx.Bind(&ors))
}

func Test_ValidValidatorRequest(t *testing.T) {
	var ors RequestStruct
	ctx := requestWithBody("POST", "/", `{"username": "x@x.com", "password": "s3cr3t"}`, "application/json")
	assert.NoError(t, ctx.Bind(&ors))
}

func Test_EmptyBodyShouldNotError(t *testing.T) {
	var ors RequestStruct
	ctx := requestWithBody("POST", "/", ``, "application/json")

	e := ctx.Bind(&ors)
	assert.IsType(t, &validation.Output{}, e)
}

func Test_BadRequest(t *testing.T) {
	var ors RequestStruct
	ctx := requestWithBody("POST", "/", `{"name": 123}`, "text/html")

	e := ctx.Bind(&ors)

	assert.NotNil(t, e)
}

func requestWithBody(method, path, body string, ctype string) echo.Context {
	e := echo.New()
	e.Binder = Binder{}
	e.HTTPErrorHandler = HTTPErrorHandler
	req, _ := http.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Add(echo.HeaderContentType, ctype)
	rec := httptest.NewRecorder()

	return e.NewContext(req, rec)
}
