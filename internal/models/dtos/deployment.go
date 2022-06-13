package dtos

type NewDeployments struct {
	Deployments []NewDeployment
	CallbackURL string
}

type NewDeployment struct {
	UseCase  string
	Services []NewService
}

type NewService struct {
	Service string
	Count   int
}

type CreatedDeployment struct {
	Canonical   string            `json:"canonical"`
	Instances   []CreatedInstance `json:"instances"`
	Type        string            `json:"type"`
	CallbackURL string            `json:"callbackURL"`
}

type CreatedInstance struct {
	State string `json:"state"`
}
