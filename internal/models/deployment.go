package models

import "github.com/csothen/tmdei-project/internal/db"

//go:generate go run github.com/dmarkham/enumer -type=Type -transform=snake -output=type_string.go -linecomment=true
//go:generate go run github.com/dmarkham/enumer -type=State -transform=snake -output=state_string.go -linecomment=true

type Type uint

type State uint

const (
	SonarqubeService Type = iota // sonarqube

	Running State = iota // running
	Pending              // pending
	Stopped              // stopped
)

type Deployment struct {
	Canonical       string     `json:"canonical"`
	State           State      `json:"state"`
	Type            Type       `json:"type"`
	URL             string     `json:"url"`
	AdminCredential Credential `json:"admin_cred"`
	UserCredential  Credential `json:"user_cred"`
	CallbackURL     string     `json:"callback_url"`
}

type Credential struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	AccessToken string `json:"access_token"`
}

func (d *Deployment) FromDB(dbd *db.Deployment) {
	d.Canonical = dbd.Canonical
	d.State = State(dbd.State)
	d.Type = Type(dbd.Type)
	d.URL = dbd.URL
	d.CallbackURL = dbd.CallbackURL

	var aCred Credential
	aCred.FromDB(&dbd.AdminCredential)
	d.AdminCredential = aCred

	var uCred Credential
	uCred.FromDB(&dbd.UserCredential)
	d.UserCredential = uCred
}

func (d *Deployment) ToDB() *db.Deployment {
	return &db.Deployment{
		Canonical:       d.Canonical,
		State:           uint(d.State),
		Type:            uint(d.Type),
		URL:             d.URL,
		AdminCredential: *d.AdminCredential.ToDB(),
		UserCredential:  *d.UserCredential.ToDB(),
		CallbackURL:     d.CallbackURL,
	}
}

func (c *Credential) FromDB(dbc *db.Credential) {
	c.Username = dbc.Username
	c.Password = dbc.Password
	c.AccessToken = dbc.AccessToken
}

func (c *Credential) ToDB() *db.Credential {
	return &db.Credential{
		Username:    c.Username,
		Password:    c.Password,
		AccessToken: c.AccessToken,
	}
}
