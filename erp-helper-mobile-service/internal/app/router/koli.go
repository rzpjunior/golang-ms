package router

import "git.edenfarm.id/project-version3/erp-services/erp-helper-mobile-service/internal/app/handler"

func init() {
	handlers["koli"] = &handler.KoliHandler{}
}