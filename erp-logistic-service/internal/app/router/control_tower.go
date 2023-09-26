package router

import "git.edenfarm.id/project-version3/erp-services/erp-logistic-service/internal/app/handler"

func init() {
	handlers["control_tower"] = &handler.ControlTowerHandler{}
}
