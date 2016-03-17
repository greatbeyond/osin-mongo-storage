package mgostore

import (
	"github.com/ashkang/osin-mongo-storage/mgostore/models"

	"github.com/RangelReale/osin"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// collection names for the entities
const (
	CLIENT_COL    = "clients"
	AUTHORIZE_COL = "authorizations"
	ACCESS_COL    = "accesses"
)

const REFRESHTOKEN = "refresh_token"

type MongoStorage struct {
	dbName  string
	session *mgo.Session
}

func New(session *mgo.Session, dbName string) *MongoStorage {
	storage := &MongoStorage{dbName, session}
	index := mgo.Index{
		Key:        []string{REFRESHTOKEN},
		Unique:     false, // refreshtoken is sometimes empty
		DropDups:   false,
		Background: true,
		Sparse:     true,
	}
	accesses := storage.session.DB(dbName).C(ACCESS_COL)
	err := accesses.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
	return storage
}

// SetClient stores a new osin.Client or updates an existing
func (store *MongoStorage) SetClient(id string, client osin.Client) error {
	session := store.session.Copy()
	defer session.Close()
	clients := session.DB(store.dbName).C(CLIENT_COL)
	mgoClient := models.NewMgoClient(client)
	_, err := clients.UpsertId(id, mgoClient)
	return err
}

// LoadAccesses returns all accesses that matches given query
func (store *MongoStorage) LoadAccesses(query bson.M) ([]*osin.AccessData, error) {
	session := store.session.Copy()
	defer session.Close()
	accesses := session.DB(store.dbName).C(ACCESS_COL)
	var mgoAccessData []*models.MgoAccessData
	err := accesses.Find(query).All(&mgoAccessData)

	accessData := make([]*osin.AccessData, len(mgoAccessData))
	for k, v := range mgoAccessData {
		accessData[k] = v.MapToOsinAccessData()
	}

	return accessData, err
}

// Clone implements the osin.Storage interface
func (store *MongoStorage) Clone() osin.Storage {
	return store
}

// Close implements the osin.Storage interface
func (store *MongoStorage) Close() {

}

// GetClient implements the osin.Storage interface
func (store *MongoStorage) GetClient(id string) (osin.Client, error) {
	session := store.session.Copy()
	defer session.Close()
	clients := session.DB(store.dbName).C(CLIENT_COL)
	mgoClient := models.NewMgoClient(nil)
	err := clients.FindId(id).One(mgoClient)
	return mgoClient.MapToOsinClient(), err
}

// SaveAuthorize implements the osin.Storage interface
func (store *MongoStorage) SaveAuthorize(data *osin.AuthorizeData) error {
	session := store.session.Copy()
	defer session.Close()
	authorizations := session.DB(store.dbName).C(AUTHORIZE_COL)
	mgoAuthorizeData := models.NewMgoAuthorizeData(data)
	_, err := authorizations.UpsertId(data.Code, mgoAuthorizeData)
	return err
}

// LoadAuthorize implements the osin.Storage interface
func (store *MongoStorage) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
	session := store.session.Copy()
	defer session.Close()
	authorizations := session.DB(store.dbName).C(AUTHORIZE_COL)
	mgoAuthorizeData := models.NewMgoAuthorizeData(nil)
	err := authorizations.FindId(code).One(mgoAuthorizeData)
	return mgoAuthorizeData.MapToOsinAuthorizeData(), err
}

// RemoveAuthorize implements the osin.Storage interface
func (store *MongoStorage) RemoveAuthorize(code string) error {
	session := store.session.Copy()
	defer session.Close()
	authorizations := session.DB(store.dbName).C(AUTHORIZE_COL)
	return authorizations.RemoveId(code)
}

// SaveAccess implements the osin.Storage interface
func (store *MongoStorage) SaveAccess(data *osin.AccessData) error {
	session := store.session.Copy()
	defer session.Close()
	accesses := session.DB(store.dbName).C(ACCESS_COL)
	mgoAccessData := models.NewMgoAccessData(data)
	_, err := accesses.UpsertId(data.AccessToken, mgoAccessData)
	return err
}

// LoadAccess implements the osin.Storage interface
func (store *MongoStorage) LoadAccess(token string) (*osin.AccessData, error) {
	session := store.session.Copy()
	defer session.Close()
	accesses := session.DB(store.dbName).C(ACCESS_COL)
	mgoAccessData := models.NewMgoAccessData(nil)
	err := accesses.FindId(token).One(mgoAccessData)
	return mgoAccessData.MapToOsinAccessData(), err
}

// RemoveAccess implements the osin.Storage interface
func (store *MongoStorage) RemoveAccess(token string) error {
	session := store.session.Copy()
	defer session.Close()
	accesses := session.DB(store.dbName).C(ACCESS_COL)
	return accesses.RemoveId(token)
}

// LoadRefresh implements the osin.Storage interface
func (store *MongoStorage) LoadRefresh(token string) (*osin.AccessData, error) {
	session := store.session.Copy()
	defer session.Close()
	accesses := session.DB(store.dbName).C(ACCESS_COL)
	mgoAccessData := models.NewMgoAccessData(nil)
	err := accesses.Find(bson.M{REFRESHTOKEN: token}).One(mgoAccessData)
	return mgoAccessData.MapToOsinAccessData(), err
}

// RemoveRefresh implements the osin.Storage interface
func (store *MongoStorage) RemoveRefresh(token string) error {
	session := store.session.Copy()
	defer session.Close()
	accesses := session.DB(store.dbName).C(ACCESS_COL)
	return accesses.Update(bson.M{REFRESHTOKEN: token}, bson.M{
		"$unset": bson.M{
			REFRESHTOKEN: 1,
		}})
}
