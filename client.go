package ubm

import "sync"

type commander interface {
	AddAction(id string, action string) error
	GetAction(id string, action string) (Action, error)
	GetLastAction(id string) (LastAction, error)
}

// Client implements the wrapper over the concrete transport API
type Client struct {
	client   commander
	triggers *sync.Map
}

func newClient(cmndr commander) Client {
	c := Client{
		client:   cmndr,
		triggers: new(sync.Map),
	}

	return c
}

// AddAction registers an action for this id
func (c Client) AddAction(id string, action string) error {
	if v, ok := c.triggers.Load(id); ok {
		trgr := v.(*trigger)
		trgr.ProcessDisableAction(action)
	}
	return c.client.AddAction(id, action)
}

// GetAction gets information about specific action for this id
func (c Client) GetAction(id string, action string) (Action, error) {
	return c.client.GetAction(id, action)
}

// GetLastAction gets information about any last action of this id
func (c Client) GetLastAction(id string) (LastAction, error) {
	return c.client.GetLastAction(id)
}

// SetTrigger sets the state tracking on the client by id,
// per enableAction action, which can be disable via disableAction
func (c Client) SetTrigger(id string, enableAction string, disableAction string) {
	t, _ := c.triggers.LoadOrStore(id, &trigger{})
	trgr := t.(*trigger)
	trgr.Set(enableAction, disableAction)
	c.triggers.Store(id, trgr)
}

// IsTriggerBlocked checks the action "Is it blocked by a trigger?"
// If blocked, then returns True.
// To unlock the required to be added the disableAction.
func (c *Client) IsTriggerBlocked(id string, action string) bool {
	t, ok := c.triggers.Load(id)
	if ok {
		trgr := t.(*trigger)
		_, ok = trgr.enableActionList.Load(action)
		return ok
	}

	return false
}
