package ubm

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/coalalib/coalago"
	m "github.com/coalalib/coalago/message"
)

type clientCoala struct {
	coala *coalago.Coala
	host  net.Addr
}

// NewClientCoala returns the Client for Coala protocol
func NewClientCoala(host string) Client {
	addr, err := net.ResolveUDPAddr("udp", host)
	if err != nil {
		panic(err)
	}

	c := newClient(clientCoala{
		coala: coalago.NewListen(0),
		host:  addr,
	})

	return c
}

func (c clientCoala) AddAction(id string, action string) error {
	message := newCoalaMessageAddAction(id, action)
	_, err := c.coala.Send(message, c.host)
	return err
}

func (c clientCoala) GetAction(id string, action string) (Action, error) {
	a := Action{}
	message := newCoalaMessageGetAction(id, action)
	resp, err := c.coala.Send(message, c.host)
	if err != nil {
		return a, err
	}

	if err = errFromResponse(resp); err != nil {
		return a, err
	}

	err = json.Unmarshal(resp.Payload.Bytes(), &a)
	return a, err
}

func (c clientCoala) GetLastAction(id string) (LastAction, error) {
	la := LastAction{}
	message := newCoalaMessageGetLastAction(id)
	resp, err := c.coala.Send(message, c.host)
	if err != nil {
		return la, err
	}

	if err = errFromResponse(resp); err != nil {
		return la, err
	}

	err = json.Unmarshal(resp.Payload.Bytes(), &la)
	return la, err
}

func newCoalaMessageAddAction(id, action string) *m.CoAPMessage {
	message := m.NewCoAPMessage(m.NON, m.POST)
	message.SetURIPath(coalaPathActions)
	message.SetURIQuery("id", id)
	message.SetURIQuery("action", action)
	message.SetSchemeCOAPS()
	return message
}

func newCoalaMessageGetAction(id, action string) *m.CoAPMessage {
	message := m.NewCoAPMessage(m.CON, m.GET)
	message.SetURIPath(coalaPathActions)
	message.SetURIQuery("id", id)
	message.SetURIQuery("action", action)
	message.SetSchemeCOAPS()
	return message
}

func newCoalaMessageGetLastAction(id string) *m.CoAPMessage {
	message := m.NewCoAPMessage(m.CON, m.GET)
	message.SetURIPath(coalaPathActionsLast)
	message.SetURIQuery("id", id)
	message.SetSchemeCOAPS()
	return message
}

func errFromResponse(message *m.CoAPMessage) error {
	if message.Code < m.CoapCodeBadRequest {
		return nil
	}

	return fmt.Errorf("%s", message.Payload.String())
}
