package ubm

import "time"

type Action struct {
	LastCall time.Time `bson:"last_call" json:"last_call"`
	Count    int64     `bson:"count" json:"count"`
}

type LastAction struct {
	Name     string    `bson:"last_action" json:"last_action"`
	LastCall time.Time `bson:"last_call" json:"last_call"`
}

type actionCollection struct {
	ID      string            `bson:"id"`
	Actions map[string]Action `bson:"actions"`
	LastAction
}
