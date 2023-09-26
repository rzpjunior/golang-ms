package router

import "git.edenfarm.id/project-version3/erp-services/erp-sales-service/internal/app/handler"

func init() {
	handlers["sales_price_level"] = &handler.SalesPriceLevelHandler{}
}
