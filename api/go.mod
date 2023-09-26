module git.edenfarm.id/project-version2/api

go 1.15

replace git.edenfarm.id/project-version2/datamodel v0.0.0 => ../datamodel

require (
	git.edenfarm.id/cuxs/common v1.2.2
	git.edenfarm.id/cuxs/cuxs v1.3.11
	git.edenfarm.id/cuxs/dbredis v1.0.6
	git.edenfarm.id/cuxs/env v0.0.0-20191008033713-5add4703b723
	git.edenfarm.id/cuxs/mongodb v1.0.1
	git.edenfarm.id/cuxs/orm v1.2.5
	git.edenfarm.id/cuxs/validation v0.0.0-20191008040451-f67bf0c4c24d
	git.edenfarm.id/project-version2/datamodel v0.0.0
	github.com/ahmetb/go-linq/v3 v3.2.0
	github.com/getsentry/sentry-go v0.11.0
	github.com/go-playground/assert v1.2.1
	github.com/go-playground/assert/v2 v2.0.1
	github.com/go-redis/redis/v8 v8.11.4
	github.com/go-redsync/redsync/v4 v4.5.0
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/labstack/echo/v4 v4.6.1
	github.com/mattn/go-sqlite3 v2.0.3+incompatible // indirect
	github.com/minio/minio-go/v7 v7.0.11
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/robfig/cron/v3 v3.0.1
	github.com/segmentio/kafka-go v0.4.23
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/tealeg/xlsx v1.0.5
	github.com/ulule/limiter/v3 v3.10.0
	github.com/xendit/xendit-go v0.8.0
	go.mongodb.org/mongo-driver v1.9.1
	golang.org/x/text v0.3.7
	gopkg.in/check.v1 v1.0.0-20200902074654-038fdea0a05b // indirect
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
)
