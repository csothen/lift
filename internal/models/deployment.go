package models

import (
	"time"

	"github.com/csothen/lift/internal/db"
	"gorm.io/gorm"
)

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
	Canonical   string     `json:"canonical"`
	Instances   []Instance `json:"instances"`
	Type        Type       `json:"type"`
	CallbackURL string     `json:"callback_url"`
	CreatedAt   time.Time  `json:"created_at"`
}

type Instance struct {
	URL             string     `json:"url"`
	State           State      `json:"state"`
	AdminCredential Credential `json:"admin_cred"`
	UserCredential  Credential `json:"user_cred"`
}

type Credential struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	AccessToken string `json:"access_token"`
}

func (d *Deployment) FromDB(dbd *db.Deployment) {
	d.Canonical = dbd.Canonical
	d.Type = Type(dbd.Type)
	d.CallbackURL = dbd.CallbackURL
	d.CreatedAt = dbd.CreatedAt

	d.Instances = make([]Instance, len(dbd.Instances))
	for i, dbi := range dbd.Instances {
		instance := &Instance{}
		instance.FromDB(&dbi)
		d.Instances[i] = *instance
	}
}

func (d *Deployment) ToDB() *db.Deployment {
	dbd := &db.Deployment{
		Canonical:   d.Canonical,
		Instances:   make([]db.Instance, len(d.Instances)),
		Type:        uint(d.Type),
		CallbackURL: d.CallbackURL,
		Model: gorm.Model{
			CreatedAt: d.CreatedAt,
		},
	}

	for i, instance := range d.Instances {
		dbd.Instances[i] = *instance.ToDB(d.Canonical)
	}
	return dbd
}

func (i *Instance) FromDB(dbi *db.Instance) {
	i.URL = dbi.URL
	i.State = State(dbi.State)

	aCred := Credential{}
	aCred.FromDB(&dbi.AdminCredential)
	i.AdminCredential = aCred

	uCred := Credential{}
	uCred.FromDB(&dbi.UserCredential)
	i.UserCredential = uCred
}

func (i *Instance) ToDB(dcan string) *db.Instance {
	return &db.Instance{
		DeploymentCanonical: dcan,
		URL:                 i.URL,
		State:               uint(i.State),
		AdminCredential:     *i.AdminCredential.ToDB(),
		UserCredential:      *i.UserCredential.ToDB(),
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
