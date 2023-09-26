package opt

import (
	"git.edenfarm.id/edenlabs/edenlabs/config"
	"git.edenfarm.id/edenlabs/edenlabs/env"
	"git.edenfarm.id/edenlabs/edenlabs/jwt"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/mongox"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/edenlabs/edenlabs/producer"
	"git.edenfarm.id/edenlabs/edenlabs/pubsubx"
	"git.edenfarm.id/edenlabs/edenlabs/redisx"
	"git.edenfarm.id/edenlabs/edenlabs/s3x"
	"git.edenfarm.id/project-version3/erp-pkg/erp-client-grpc/client"
	"go.opentelemetry.io/otel/trace"
)

type Options struct {
	Env      env.Provider
	Config   config.Config
	Database Database
	Mongox   mongox.Client
	Redisx   redisx.Client
	Logger   *log.Logger
	Trace    trace.Tracer
	S3x      s3x.Client
	Jwt      jwt.JWT
	Client   client.Clients
	Producer *producer.Producers
	PubSub   pubsubx.IPubSub
}

type Database struct {
	Write orm.Ormer
	Read  orm.Ormer
}
