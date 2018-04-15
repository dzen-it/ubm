package ubm

type withoutDB struct {
}

//NewWithoutDB create an empty instance that does not write and does not read anything
func NewWithoutDB() DB {
	return withoutDB{}
}

func (w withoutDB) AddAction(id interface{}, actionName string) error {
	return nil
}

func (w withoutDB) GetAction(id interface{}, actionName string) (a Action, err error) {
	return a, nil
}

func (w withoutDB) GetLastAction(id interface{}) (a LastAction, err error) {
	return a, nil
}
