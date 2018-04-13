package ubm

type commander interface {
	AddAction(id string, action string) error
	GetAction(id string, action string) (Action, error)
	GetLastAction(id string) (LastAction, error)
}

// Client implements the wrapper over the concrete transport API
type Client struct {
	client commander
}

// AddAction registers an action for this id
func (c Client) AddAction(id string, action string) error {
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
