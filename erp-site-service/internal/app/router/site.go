package router

import "git.edenfarm.id/project-version3/erp-services/erp-site-service/internal/app/handler"

func init() {
	handlers["site"] = &handler.SiteHandler{}
}
