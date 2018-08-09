package ubm

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/coalalib/coalago"
)

type clientCoala struct {
	host net.Addr
}

// NewClientCoala returns the Client for Coala protocol
func NewClientCoala(host string) Client {
	addr, err := net.ResolveUDPAddr("udp", host)
	if err != nil {
		panic(err)
	}

	c := newClient(clientCoala{
		host: addr,
	})

	return c
}

func (c clientCoala) AddAction(id string, action string) error {
	message := newCoalaMessageAddAction(id, action)

	client := coalago.NewClient()
	_, err := client.Send(message, c.host.String())
	return err
}

func (c clientCoala) GetAction(id string, action string) (Action, error) {
	a := Action{}
	message := newCoalaMessageGetAction(id, action)
	client := coalago.NewClient()
	resp, err := client.Send(message, c.host.String())
	if err != nil {
		return a, err
	}

	if err = errFromResponse(resp); err != nil {
		return a, err
	}

	err = json.Unmarshal(resp.Body, &a)
	return a, err
}

func (c clientCoala) GetLastAction(id string) (LastAction, error) {
	la := LastAction{}
	message := newCoalaMessageGetLastAction(id)
	client := coalago.NewClient()
	resp, err := client.Send(message, c.host.String())
	if err != nil {
		return la, err
	}

	if err = errFromResponse(resp); err != nil {
		return la, err
	}

	err = json.Unmarshal(resp.Body, &la)
	return la, err
}

func newCoalaMessageAddAction(id, action string) *coalago.CoAPMessage {
	message := coalago.NewCoAPMessage(coalago.NON, coalago.POST)
	message.SetURIPath(coalaPathActions)
	message.SetURIQuery("id", id)
	message.SetURIQuery("action", action)
	message.SetSchemeCOAPS()
	return message
}

func newCoalaMessageGetAction(id, action string) *coalago.CoAPMessage {
	message := coalago.NewCoAPMessage(coalago.CON, coalago.GET)
	message.SetURIPath(coalaPathActions)
	message.SetURIQuery("id", id)
	message.SetURIQuery("action", action)
	message.SetSchemeCOAPS()
	return message
}

func newCoalaMessageGetLastAction(id string) *coalago.CoAPMessage {
	message := coalago.NewCoAPMessage(coalago.CON, coalago.GET)
	message.SetURIPath(coalaPathActionsLast)
	message.SetURIQuery("id", id)
	message.SetSchemeCOAPS()
	return message
}

func errFromResponse(resp *coalago.Response) error {
	if resp.Code < coalago.CoapCodeBadRequest {
		return nil
	}

	return fmt.Errorf("%s", resp.Body)
}
