package router

import "git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/handler"

func init() {
	handlers["sales/performance"] = &handler.SalesPerformanceHandler{}
}
