package router

import (
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/handler"
)

func init() {
	handlers["vehicle_profile"] = &handler.VehicleProfileHandler{}
}
