package repository

import (
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/global"
)

type ISalesAssignmentSubmissionRepository interface {
}

type SalesAssignmentSubmissionRepository struct {
	opt opt.Options
}

func NewSalesAssignmentSubmissionRepository() ISalesAssignmentSubmissionRepository {
	return &SalesAssignmentSubmissionRepository{
		opt: global.Setup.Common,
	}
}
