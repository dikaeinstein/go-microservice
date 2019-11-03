package data

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// MongoStore is a MongoDB data store which implements the Store interface
type MongoStore struct {
	session *mgo.Session
}

// NewMongoStore creates an instance of MongoStore with the given connection string
func NewMongoStore(connection string) (*MongoStore, error) {
	session, err := mgo.Dial(connection)
	if err != nil {
		return nil, err
	}

	return &MongoStore{session: session}, nil
}

// Search returns Kittens from the MongoDB instance which have the name name
func (m *MongoStore) Search(name string) ([]Kitten, error) {
	s := m.session.Clone()
	defer s.Close()

	var results []Kitten
	c := s.DB("kittenserver").C("kittens")

	err := c.Find(bson.M{"name": name}).All(&results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// DeleteAllKittens deletes all the kittens from the datastore
func (m *MongoStore) DeleteAllKittens() error {
	s := m.session.Clone()
	defer s.Close()

	return s.DB("kittenserver").C("kittens").DropCollection()
}

// InsertKittens inserts a slice of kittens into the datastore
func (m *MongoStore) InsertKittens(kittens ...Kitten) error {
	s := m.session.Clone()
	defer s.Close()

	documents := make([]interface{}, len(kittens))
	for i, v := range kittens {
		documents[i] = v
	}
	return s.DB("kittenserver").C("kittens").Insert(documents...)
}
