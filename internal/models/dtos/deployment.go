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
	Canonical string `json:"canonical"`
	State     string `json:"state"`
	Type      string `json:"type"`
}
