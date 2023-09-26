// Copyright 2022 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"time"

	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	loc, err := time.LoadLocation("")
	if err != nil {
		logrus.Error(err.Error())
		os.Exit(1)
	}
	time.Local = loc

	if err := cmd.NewRootCommand().Execute(); err != nil {
		logrus.Error(err.Error())
		os.Exit(1)
	}
}
