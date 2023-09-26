// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package auth_test

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"git.edenfarm.id/project-version2/api/test"

	"git.edenfarm.id/cuxs/common/tester"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	test.Setup()

	// run tests
	res := m.Run()

	// cleanup
	//test.DataCleanUp()

	os.Exit(res)
}

func TestAuthenticationHandler_URLMapping(t *testing.T) {
	var auth = []struct {
		req      tester.D
		expected int
	}{
		{tester.D{"email": "", "password": ""}, http.StatusUnprocessableEntity},
		{tester.D{"email": "sysadmin", "password": "xxxxx"}, http.StatusUnprocessableEntity},
		{tester.D{"email": "qasico", "password": "qasico123"}, http.StatusUnprocessableEntity},
		{tester.D{"email": "sysadmin@qasico.com", "password": "qasico123"}, http.StatusOK},
	}

	ng := tester.New()
	for _, tes := range auth {
		ng.POST("/v1/auth").
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}
}
