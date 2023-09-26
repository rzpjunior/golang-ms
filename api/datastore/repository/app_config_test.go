// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository_test

import (
	"testing"

	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"

	"git.edenfarm.id/cuxs/orm"
	"github.com/stretchr/testify/assert"
)

func TestGetAppConfig(t *testing.T) {
	_, e := repository.GetAppConfig("id", 1000)
	assert.Error(t, e, "Response should be error, because there are no data yet.")

	c := model.DummyAppConfig()
	cd, e := repository.GetAppConfig("id", c.ID)
	assert.NoError(t, e, "Data should be exists.")
	assert.Equal(t, c.ID, cd.ID, "ID Response should be a same.")
}

func TestGetAppConfigs(t *testing.T) {
	model.DummyAppConfig()
	qs := orm.RequestQuery{}
	_, _, e := repository.GetAppConfigs(&qs)
	assert.NoError(t, e, "Data should be exists.")
}
