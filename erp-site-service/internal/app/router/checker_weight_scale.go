package router

import "git.edenfarm.id/project-version3/erp-services/erp-site-service/internal/app/handler"

func init() {
	handlers["checker_weight_scale"] = &handler.CheckerWeightScaleHandler{}
}
