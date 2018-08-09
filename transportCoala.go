package ubm

import (
	"fmt"

	"github.com/coalalib/coalago"
)

const (
	defaultCoalaPort     = 5683
	coalaPathActions     = "/ubm/actions"
	coalaPathActionsLast = "/ubm/actions/last"
)

type transportCoala struct {
	addr   string
	server *coalago.Server
}

func (t transportCoala) Serve() error {
	return t.server.Listen(t.addr)
}

// NewServerCoala returns implementation of Listener interface for Coala protocol
func NewServerCoala(port int, db DB) Server {
	h := coalaBaseHandler{
		db: db,
	}

	server := coalago.NewServer()

	server.AddPOSTResource(coalaPathActions, h.handlerPostAction)
	server.AddGETResource(coalaPathActions, h.handlerGetAction)
	server.AddGETResource(coalaPathActionsLast, h.handlerGetLastAction)

	t := transportCoala{
		server: server,
		addr:   fmt.Sprintf(":%d", port),
	}

	return t
}

type coalaBaseHandler struct {
	db DB
}

func (c *coalaBaseHandler) handlerPostAction(message *coalago.CoAPMessage) *coalago.CoAPResourceHandlerResult {
	id := message.GetURIQuery("id")
	action := message.GetURIQuery("action")

	if len(id) == 0 || len(action) == 0 {
		return coalago.NewResponse(coalago.NewStringPayload("invalid queries"), coalago.CoapCodeBadRequest)
	}

	if err := c.db.AddAction(id, action); err != nil {
		return coalago.NewResponse(coalago.NewStringPayload("internal error"), coalago.CoapCodeInternalServerError)
	}

	return coalago.NewResponse(coalago.NewEmptyPayload(), coalago.CoapCodeChanged)
}

func (c *coalaBaseHandler) handlerGetAction(message *coalago.CoAPMessage) *coalago.CoAPResourceHandlerResult {
	id := message.GetURIQuery("id")
	action := message.GetURIQuery("action")

	if len(id) == 0 || len(action) == 0 {
		return coalago.NewResponse(coalago.NewStringPayload("invalid queries"), coalago.CoapCodeBadRequest)
	}

	a, err := c.db.GetAction(id, action)
	if err != nil {
		if err == ErrActionNotFound {
			return coalago.NewResponse(coalago.NewStringPayload(err.Error()), coalago.CoapCodeBadRequest)
		}
		return coalago.NewResponse(coalago.NewStringPayload("internal error: "+err.Error()), coalago.CoapCodeInternalServerError)
	}

	return coalago.NewResponse(coalago.NewJSONPayload(a), coalago.CoapCodeContent)
}

func (c *coalaBaseHandler) handlerGetLastAction(message *coalago.CoAPMessage) *coalago.CoAPResourceHandlerResult {
	id := message.GetURIQuery("id")

	if len(id) == 0 {
		return coalago.NewResponse(coalago.NewStringPayload("invalid queries"), coalago.CoapCodeBadRequest)
	}

	la, err := c.db.GetLastAction(id)
	if err != nil {
		if err == ErrUserNotFound {
			return coalago.NewResponse(coalago.NewStringPayload(err.Error()), coalago.CoapCodeBadRequest)
		}
		return coalago.NewResponse(coalago.NewStringPayload("internal error"), coalago.CoapCodeInternalServerError)
	}

	return coalago.NewResponse(coalago.NewJSONPayload(la), coalago.CoapCodeContent)
}
