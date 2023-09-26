package handler

import (
	"encoding/json"
	"fmt"

	"git.edenfarm.id/edenlabs/edenlabs/log"

	"git.edenfarm.id/project-version3/erp-services/erp-boilerplate-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-boilerplate-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-boilerplate-service/internal/app/service"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
)

type BoilerplateConsumerHandler struct {
	Option         global.HandlerOptions
	ServicesPerson service.IPersonService
}

func (h BoilerplateConsumerHandler) IncomingMessage(message *message.Message) (err error) {
	ctx, span := h.Option.Common.Trace.Start(message.Context(), "Consumer.IncomingMessage")
	defer span.End()

	messagePrint := fmt.Sprintf("Received message: %s %s metadata: %v", message.UUID, string(message.Payload), message.Metadata)

	h.Option.Common.Logger.AddMessage(log.InfoLevel, messagePrint).Print()
	span.AddEvent(messagePrint)

	var event dto.PersonRequestUpdate
	if err := json.Unmarshal(message.Payload, &event); err != nil {
		return err
	}

	_, err = h.ServicesPerson.Update(ctx, event)
	if err != nil {
		err = errors.Wrap(err, "cannot create person")
		span.RecordError(err)
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}
