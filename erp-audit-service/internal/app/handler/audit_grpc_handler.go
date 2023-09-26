package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-audit-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-audit-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-audit-service/internal/app/service"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/audit_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuditGrpcHandler struct {
	Option       global.HandlerOptions
	ServiceAudit service.IAuditService
}

func (h *AuditGrpcHandler) CreateLog(ctx context.Context, req *pb.CreateLogRequest) (res *pb.CreateLogResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.CreateLog")
	defer span.End()

	auditLog, err := h.ServiceAudit.Create(ctx, dto.AuditRequestCreate{
		UserID:      req.Log.UserId,
		UserIdGp:    req.Log.UserIdGp,
		ReferenceID: req.Log.ReferenceId,
		Type:        req.Log.Type,
		Function:    req.Log.Function,
		Note:        req.Log.Note,
	})
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &pb.CreateLogResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &pb.Log{
			Id:        auditLog.ID,
			CreatedAt: timestamppb.New(auditLog.CreatedAt),
		},
	}
	return
}
