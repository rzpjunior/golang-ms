package middleware

import (
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/global"
)

// Middleware defines object for order api custom middleware
type Middleware struct {
	Option opt.Options
}

func NewMiddleware() *Middleware {
	return &Middleware{
		Option: global.Setup.Common,
	}
}
