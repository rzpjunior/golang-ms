package handler

import (
	"git.edenfarm.id/project-version3/erp-services/erp-storage-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-storage-service/internal/app/service"
)

type StorageGrpcHandler struct {
	Option         global.HandlerOptions
	ServicesUpload service.IUploadService
}
