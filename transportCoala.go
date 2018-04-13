package ubm

import (
	"fmt"

	"github.com/coalalib/coalago"
	m "github.com/coalalib/coalago/message"
	"github.com/coalalib/coalago/resource"
)

const (
	defaultCoalaPort     = 5683
	coalaPathActions     = "/ubm/actions"
	coalaPathActionsLast = "/ubm/actions/last"
)

type transportCoala struct {
	coala *coalago.Coala
}

func (t transportCoala) Serve() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	select {}
}

// NewServerCoala returns implementation of Listener interface for Coala protocol
func NewServerCoala(port int, db DB) Server {
	h := coalaBaseHandler{
		db: db,
	}

	coala := coalago.NewListen(port)

	coala.AddPOSTResource(coalaPathActions, h.handlerPostAction)
	coala.AddGETResource(coalaPathActions, h.handlerGetAction)
	coala.AddGETResource(coalaPathActionsLast, h.handlerGetLastAction)

	t := transportCoala{
		coala: coala,
	}

	return t
}

type coalaBaseHandler struct {
	db DB
}

func (c *coalaBaseHandler) handlerPostAction(message *m.CoAPMessage) *resource.CoAPResourceHandlerResult {
	id := message.GetURIQuery("id")
	action := message.GetURIQuery("action")

	if len(id) == 0 || len(action) == 0 {
		return resource.NewResponse(m.NewStringPayload("invalid queries"), m.CoapCodeBadRequest)
	}

	if err := c.db.AddAction(id, action); err != nil {
		return resource.NewResponse(m.NewStringPayload("internal error"), m.CoapCodeInternalServerError)
	}

	return resource.NewResponse(m.NewEmptyPayload(), m.CoapCodeChanged)
}

func (c *coalaBaseHandler) handlerGetAction(message *m.CoAPMessage) *resource.CoAPResourceHandlerResult {
	id := message.GetURIQuery("id")
	action := message.GetURIQuery("action")

	if len(id) == 0 || len(action) == 0 {
		return resource.NewResponse(m.NewStringPayload("invalid queries"), m.CoapCodeBadRequest)
	}

	a, err := c.db.GetAction(id, action)
	if err != nil {
		if err == ErrActionNotFound {
			return resource.NewResponse(m.NewStringPayload(err.Error()), m.CoapCodeBadRequest)
		}
		return resource.NewResponse(m.NewStringPayload("internal error: "+err.Error()), m.CoapCodeInternalServerError)
	}

	return resource.NewResponse(m.NewJSONPayload(a), m.CoapCodeContent)
}

func (c *coalaBaseHandler) handlerGetLastAction(message *m.CoAPMessage) *resource.CoAPResourceHandlerResult {
	id := message.GetURIQuery("id")

	if len(id) == 0 {
		return resource.NewResponse(m.NewStringPayload("invalid queries"), m.CoapCodeBadRequest)
	}

	la, err := c.db.GetLastAction(id)
	if err != nil {
		if err == ErrUserNotFound {
			return resource.NewResponse(m.NewStringPayload(err.Error()), m.CoapCodeBadRequest)
		}
		return resource.NewResponse(m.NewStringPayload("internal error"), m.CoapCodeInternalServerError)
	}

	return resource.NewResponse(m.NewJSONPayload(la), m.CoapCodeContent)
}
