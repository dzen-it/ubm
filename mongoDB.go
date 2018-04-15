package ubm

import (
	"time"

	"github.com/globalsign/mgo"
	"gopkg.in/mgo.v2/bson"
)

type mongoDB struct {
	session *mgo.Session
	dbName  string
}

// NewMongoDB returns implemetation of DB interface for MongoDB
func NewMongoDB(addr string, dbName string) (DB, error) {
	var err error

	db := new(mongoDB)
	db.dbName = dbName
	db.session, err = mgo.Dial(addr)

	return db, err
}

func (m mongoDB) AddAction(id interface{}, actionName string) error {
	session := m.session.Clone()
	defer session.Close()

	now := time.Now()

	_, err := session.DB(m.dbName).C(collectionActionName).Upsert(
		bson.M{
			"id": id,
		},
		bson.M{
			"$set": bson.M{
				"actions." + actionName + ".last_call": now,
				"last_call":                            now,
				"last_action":                          actionName,
			},
			"$inc": bson.M{
				"actions." + actionName + ".count": 1,
			},
		},
	)

	return err
}

func (m mongoDB) GetAction(id interface{}, actionName string) (Action, error) {
	session := m.session.Clone()
	defer session.Close()

	a := actionCollection{}

	err := session.DB(m.dbName).C(collectionActionName).Find(bson.M{
		"id": id,
		"actions." + actionName: bson.M{
			"$exists": true,
		},
	}).Select(bson.M{
		"actions." + actionName: 1,
	}).One(&a)

	if err == mgo.ErrNotFound {
		err = ErrActionNotFound
	}

	return a.Actions[actionName], err
}

func (m mongoDB) GetLastAction(id interface{}) (LastAction, error) {
	session := m.session.Clone()
	defer session.Close()

	a := LastAction{}

	err := session.DB(m.dbName).C(collectionActionName).Find(bson.M{
		"id": id,
	}).Select(bson.M{
		"last_action": 1,
		"last_call":   1,
	}).One(&a)

	if err == mgo.ErrNotFound {
		err = ErrUserNotFound
	}

	return a, err
}
