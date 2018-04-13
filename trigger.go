package ubm

import (
	"sync"
)

type trigger struct {
	enableActionList  sync.Map
	disableActionList sync.Map
}

func (trgr *trigger) Set(enableAction string, disableAction string) {
	trgr.enableActionList.Store(enableAction, nil)
	v, _ := trgr.disableActionList.LoadOrStore(disableAction, &sync.Map{})
	enableList := v.(*sync.Map)
	enableList.Store(enableAction, nil)

	trgr.enableActionList.Store(enableAction, disableAction)
}

func (trgr *trigger) ProcessDisableAction(disableAction string) {
	v, ok := trgr.disableActionList.Load(disableAction)
	if ok {
		enableList := v.(*sync.Map)
		enableList.Range(func(key, value interface{}) bool {
			trgr.enableActionList.Delete(key)
			return true
		})
		trgr.disableActionList.Delete(disableAction)
	}
}
