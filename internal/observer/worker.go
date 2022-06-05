package observer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/csothen/tmdei-project/internal/config"
	"github.com/csothen/tmdei-project/internal/models"
	"github.com/csothen/tmdei-project/internal/services"
	"github.com/csothen/tmdei-project/internal/terraform"
	"github.com/csothen/tmdei-project/internal/utils"
)

type Worker struct {
	s   *services.Service
	cfg *config.Config
}

const (
	delay time.Duration = time.Minute * 1
	limit time.Duration = time.Minute * 30

	sonarqubeHealthcheckURL string = "http://%s/health"
)

func NewWorker(s *services.Service, cfg *config.Config) *Worker {
	return &Worker{s, cfg}
}

func (w *Worker) Start() error {
	tfw := terraform.NewWorker(w.cfg.TerraformExecPath)

	for {
		ctx := context.Background()
		deployments := w.s.ReadAllDeployments(ctx)
		now := time.Now()
		for _, d := range deployments {
			if d.State != models.Pending {
				continue
			}

			ok, err := checkConnection(d)
			if !ok || err != nil {
				if err != nil {
					log.Println(err)
				}
				// In case the deployment has exceeded the limit
				// of time it is allowed to take to be up without working
				// we destroy it
				// TODO: Add the creation date to the domain model
				if now.Sub(d.CreatedAt) > limit {
					deploymentPath, err := utils.BuildDeploymentFolderPath(d.Canonical)
					if err != nil {
						log.Println(err)
						continue
					}

					err = tfw.Teardown(deploymentPath)
					if err != nil {
						log.Println(err)
						continue
					}

					d.State = models.Stopped
					err = w.s.DeleteDeployment(ctx, d.Canonical)
					if err != nil {
						log.Println(err)
						continue
					}
				}
			}

			d.State = models.Running
			err = notifyUser(d)
			if err != nil {
				log.Println(err)
			}
			w.s.UpdateDeployment(ctx, d.Canonical, d)
		}

		time.Sleep(delay)
	}
}

func checkConnection(d *models.Deployment) (bool, error) {
	switch d.Type {
	case models.SonarqubeService:
		res, err := http.Get(fmt.Sprintf(sonarqubeHealthcheckURL, d.URL))
		if err != nil || res.StatusCode != 200 {
			return false, nil
		}
		return true, nil
	default:
		return false, fmt.Errorf("service %s not supported", d.Type.String())
	}
}

func notifyUser(d *models.Deployment) error {
	d.AdminCredential = models.Credential{}
	data, err := json.Marshal(d)
	if err != nil {
		return fmt.Errorf("could not parse data into json: %w", err)
	}

	r := strings.NewReader(string(data))
	_, err = http.Post(d.CallbackURL, "application/json", r)
	if err != nil {
		return fmt.Errorf("could not call user's service: %w", err)
	}
	return nil
}
