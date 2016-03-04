package models

import (
	"time"

	"github.com/RangelReale/osin"
)

// MgoClient represents a bson struct of osin.DefaultClient
type MgoClient struct {
	ID          string      `bson:"_id"`
	RedirectURI string      `bson:"redirect_uri"`
	Secret      string      `bson:"secret"`
	UserData    interface{} `bson:"user_data"`
}

// MapToOsinClient creates a new osin.Client with properties from MgoClient
func (c *MgoClient) MapToOsinClient() osin.Client {
	return &osin.DefaultClient{
		Id:          c.ID,
		RedirectUri: c.RedirectURI,
		Secret:      c.Secret,
		UserData:    c.UserData,
	}
}

// MgoAuthorizeData represents a bson struct of osin.AuthorizeData
type MgoAuthorizeData struct {
	Client      MgoClient   `bson:"client"`
	Code        string      `bson:"code"`
	ExpiresIn   int32       `bson:"expires_in"`
	Scope       string      `bson:"scope"`
	RedirectURI string      `bson:"redirect_uri"`
	State       string      `bson:"state"`
	CreatedAt   time.Time   `bson:"created_at"`
	UserData    interface{} `bson:"user_data"`
}

// MapToOsinAuthorizeData creates a new osin.AuthorizeData with properties from MgoAuthorizeData
func (a *MgoAuthorizeData) MapToOsinAuthorizeData() *osin.AuthorizeData {
	return &osin.AuthorizeData{
		Client:      a.Client.MapToOsinClient(),
		Code:        a.Code,
		CreatedAt:   a.CreatedAt,
		ExpiresIn:   a.ExpiresIn,
		RedirectUri: a.RedirectURI,
		Scope:       a.Scope,
		State:       a.State,
		UserData:    a.UserData,
	}
}

// MgoAccessData represents a bson struct of osin.AccessData
type MgoAccessData struct {
	Client        MgoClient         `bson:"client"`
	AuthorizeData *MgoAuthorizeData `bson:"authorize_data"`
	AccessData    *MgoAccessData    `bson:"previous_access_data"`
	AccessToken   string            `bson:"access_token"`
	RefreshToken  string            `bson:"refresh_token"`
	ExpiresIn     int32             `bson:"expires_in"`
	Scope         string            `bson:"scope"`
	RedirectURI   string            `bson:"redirect_uri"`
	CreatedAt     time.Time         `bson:"created_at"`
	UserData      interface{}       `bson:"user_data"`
}

// MapToOsinAccessData creates a new osin.AccessData with properties from MgoAccessData
func (a *MgoAccessData) MapToOsinAccessData() *osin.AccessData {
	accessData := &osin.AccessData{
		AccessToken:  a.AccessToken,
		Client:       a.Client.MapToOsinClient(),
		CreatedAt:    a.CreatedAt,
		ExpiresIn:    a.ExpiresIn,
		RedirectUri:  a.RedirectURI,
		RefreshToken: a.RefreshToken,
		Scope:        a.Scope,
		UserData:     a.UserData,
	}

	if a.AccessData != nil {
		accessData.AccessData = a.AccessData.MapToOsinAccessData()
	}
	if a.AuthorizeData != nil {
		accessData.AuthorizeData = a.AuthorizeData.MapToOsinAuthorizeData()
	}

	return accessData
}

// NewMgoAccessData creates a new MgoAccessData with properties from osin.AccessData
func NewMgoAccessData(data *osin.AccessData) *MgoAccessData {
	if data == nil {
		return new(MgoAccessData)
	}

	MgoAccessData := &MgoAccessData{
		AccessToken:  data.AccessToken,
		Client:       *NewMgoClient(data.Client),
		CreatedAt:    data.CreatedAt,
		ExpiresIn:    data.ExpiresIn,
		RedirectURI:  data.RedirectUri,
		RefreshToken: data.RefreshToken,
		Scope:        data.Scope,
		UserData:     data.UserData,
	}

	if data.AccessData != nil {
		MgoAccessData.AccessData = NewMgoAccessData(data.AccessData)
	}
	if data.AuthorizeData != nil {
		MgoAccessData.AuthorizeData = NewMgoAuthorizeData(data.AuthorizeData)
	}

	return MgoAccessData
}

// NewMgoAuthorizeData creates a new MgoAuthorizeData with properties from osin.AuthorizeData
func NewMgoAuthorizeData(data *osin.AuthorizeData) *MgoAuthorizeData {
	if data == nil {
		return new(MgoAuthorizeData)
	}

	return &MgoAuthorizeData{
		Client:      *NewMgoClient(data.Client),
		Code:        data.Code,
		ExpiresIn:   data.ExpiresIn,
		Scope:       data.Scope,
		RedirectURI: data.RedirectUri,
		State:       data.State,
		CreatedAt:   data.CreatedAt,
		UserData:    data.UserData,
	}
}

// NewMgoClient creates a new MgoClient with properties from osin.Client
func NewMgoClient(client osin.Client) *MgoClient {
	if client == nil {
		return new(MgoClient)
	}

	return &MgoClient{
		ID:          client.GetId(),
		RedirectURI: client.GetRedirectUri(),
		Secret:      client.GetSecret(),
		UserData:    client.GetUserData(),
	}
}
