package ubm

type clientInMemory struct{}

// NewClientInMemory returns the Client working without backend
func NewClientInMemory() Client {
	return newClient(clientInMemory{})
}

func (c clientInMemory) AddAction(id string, action string) error {
	return nil
}

func (c clientInMemory) GetAction(id string, action string) (Action, error) {
	a := Action{}
	return a, nil
}

func (c clientInMemory) GetLastAction(id string) (LastAction, error) {
	la := LastAction{}
	return la, nil
}
