// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"git.edenfarm.id/cuxs/dbredis"

	"git.edenfarm.id/cuxs/common/log"
	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/cuxs/env"
	"git.edenfarm.id/project-version2/api/engine"
	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
)

// init preparing application instances.
func init() {
	log.DebugMode = cuxs.IsDebug()
	log.Log = log.New()

	if e := cuxs.DbSetup(); e != nil {
		panic(e)
	}
	if e := cuxs.DbSetupReadOnly(); e != nil {
		panic(e)
	}

	if e := cuxs.DbSetupScraping(); e != nil {
		panic(e)
	}

	if _, e := dbredis.RedisStart(); e != nil {
		panic(e)
	}
}

// main creating new instances application
// and serving application server.
func main() {
	// Registering Sentry
	sentryEnv := env.GetString("SENTRY_ENVIRONMENT", "")
	sentryClientOptions := sentry.ClientOptions{
		Dsn:              env.GetString("SENTRY_DSN", ""),
		Environment:      sentryEnv,
		TracesSampleRate: 0.1,
	}
	if sentryEnv == "production" {
		sentryClientOptions.Release = "api@" + env.GetString("COREAPI_RELEASE_TAG", "")
	} else {
		sentryClientOptions.Release = "api@staging"
	}
	if err := sentry.Init(sentryClientOptions); err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	}

	e := engine.Router()
	e.Use(sentryecho.New(sentryecho.Options{Repanic: true}))

	DbHostAuditor := env.GetString("DB_HOST_AUDITOR", "")
	DbNameAuditor := env.GetString("DB_NAME_AUDITOR", "")
	DbUserAuditor := env.GetString("DB_USER_AUDITOR", "")
	DbPasswordAuditor := env.GetString("DB_PASS_AUDITOR", "")

	if DbPasswordAuditor != "" && DbHostAuditor != "" && DbNameAuditor != "" && DbUserAuditor != "" {
		if e := cuxs.DbSetupAuditor(); e != nil {
			sentry.CaptureMessage("ERROR TO DB START AUDITOR : " + e.Error())
		}
	}
	// starting server
	cuxs.StartServer(e)
}
