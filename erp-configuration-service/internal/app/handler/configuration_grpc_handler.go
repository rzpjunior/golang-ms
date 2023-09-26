package handler

import (
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/internal/app/service"
)

type ConfigurationGrpcHandler struct {
	Option               global.HandlerOptions
	ServiceCodeGenerator service.ICodeGeneratorService
	ServiceGlossary      service.IGlossaryService
	ServiceWrt           service.IWrtService
	ServiceConfigApp     service.IApplicationConfigService
	ServiceRegionPolicy  service.IRegionPolicyService
	ServiceDayOff        service.IDayOffService
}
