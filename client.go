package ubm

import (
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

type commander interface {
	AddAction(id string, action string) error
	GetAction(id string, action string) (Action, error)
	GetLastAction(id string) (LastAction, error)
}

// Client implements the wrapper over the concrete transport API
type Client struct {
	client                 commander
	globalDisabledTriggers *sync.Map
	globalEnabledTriggers  *sync.Map
	globalTriggersCache    *cache.Cache
}

func newClient(cmndr commander) Client {
	c := Client{
		client:                 cmndr,
		globalDisabledTriggers: new(sync.Map),
		globalEnabledTriggers:  new(sync.Map),
		globalTriggersCache:    cache.New(cache.NoExpiration, time.Second*10),
	}

	return c
}

// AddAction registers an action for this id
func (c Client) AddAction(id string, action string) error {
	c.updateGlobalTrigger(id, action)
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
func (c Client) SetTrigger(enableAction string, disableAction string, lifetime time.Duration) {
	c.globalDisabledTriggers.Store(disableAction, enableAction)
	c.globalEnabledTriggers.Store(enableAction, lifetime)
}

// TriggerStatus checks the action "Is it blocked by a trigger?"
// If blocked, then returns True.
// To unlock the required to be added the disableAction.
func (c *Client) TriggerStatus(id string, action string) bool {
	_, ok := c.globalTriggersCache.Get(id + action)
	return ok
}

func (c *Client) updateGlobalTrigger(id, action string) {
	e, ok := c.globalDisabledTriggers.Load(action)
	if ok {
		enaction := e.(string)
		isRestored, ok := c.globalTriggersCache.Get(id + enaction)
		if !ok {
			return
		}

		if isRestored.(bool) {
			return
		}

		if lifetime, ok := c.globalEnabledTriggers.Load(enaction); ok {
			c.globalTriggersCache.Set(id+enaction, true, lifetime.(time.Duration))
		}

		return
	}

	_, ok = c.globalEnabledTriggers.Load(action)
	if ok {
		c.globalTriggersCache.Set(id+action, false, cache.NoExpiration)
		return
	}
}
