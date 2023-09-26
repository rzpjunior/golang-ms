package handler

import (
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/service"
)

type CrmGrpcHandler struct {
	Option                           global.HandlerOptions
	ServicesProspectiveCustomer      service.IProspectiveCustomerService
	ServicesSalesAssignment          service.ISalesAssignmentService
	ServicesSalesAssignmentItem      service.ISalesAssignmentItemService
	ServicesSalesAssignmentObjective service.ISalesAssignmentObjectiveService
	ServicesCustomerAcquisition      service.ICustomerAcquisitionService
	ServicesSalesFailedVisit         service.ISalesFailedVisitService
	ServicesSalesSubmission          service.ISalesAssignmentSubmissionService
	ServicesCustomer                 service.ICustomerService
}
