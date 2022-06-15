package observer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/csothen/lift/internal/config"
	"github.com/csothen/lift/internal/models"
	"github.com/csothen/lift/internal/sdk/jenkins"
	"github.com/csothen/lift/internal/sdk/sonarqube"
	"github.com/csothen/lift/internal/services"
	"github.com/csothen/lift/internal/terraform"
	"github.com/csothen/lift/internal/utils"
)

type Worker struct {
	s   *services.Service
	cfg *config.Config
}

const (
	delay time.Duration = time.Minute * 1
	limit time.Duration = time.Minute * 15
)

func NewWorker(s *services.Service, cfg *config.Config) *Worker {
	return &Worker{s, cfg}
}

func (w *Worker) Start() error {
	tfw := terraform.NewWorker(w.cfg.TerraformExecPath)

	for {
		log.Println("Started iterating over deployments...")
		ctx := context.Background()
		deployments := w.s.ReadAllDeployments(ctx)
		now := time.Now()
		for _, d := range deployments {
			// check whether we already exceeded the time limit
			exceededLimit := now.Sub(d.CreatedAt) > limit

			/**
			For each instance we will check whether it is up and running
			If it is not and the limit was exceeded we delete the whole deployment
			Otherwise we skip it
			If it is running we will create a user for the end user to use to login
			we persist that information and notify the user
			*/
			for _, i := range d.Instances {
				if i.State != models.Pending {
					continue
				}

				ok := w.checkConnection(i, d.Type)
				if !ok {
					err := w.handleFailedConnection(d.Canonical, exceededLimit, tfw)
					if err != nil {
						log.Println(err)
					}
					continue
				}

				err := w.handleSuccessfulConnection(i, *d)
				if err != nil {
					log.Println(err)
				}
			}
		}
		time.Sleep(delay)
	}
}

func (w *Worker) checkConnection(i models.Instance, st models.Type) bool {
	log.Printf("Checking connection on URL %s\n", i.URL)
	switch st {
	case models.SonarqubeService:
		status, err := sonarqube.Status(i.URL)
		return err == nil && status == sonarqube.Up
	case models.JenkinsService:
		status, err := jenkins.Status(i.URL)
		return err == nil && status == jenkins.Up
	default:
		return false
	}
}

func (w *Worker) handleFailedConnection(canonical string, exceededLimit bool, tfw *terraform.Worker) error {
	log.Println("Handling failed connection")
	if !exceededLimit {
		return nil
	}

	go func() {
		_, err := w.s.DeleteDeployment(context.Background(), canonical)
		if err != nil {
			log.Println(err)
		}
	}()
	return nil
}

func (w *Worker) handleSuccessfulConnection(i models.Instance, d models.Deployment) error {
	log.Println("Handling successful connection")

	// set the instance's admin credentials
	adminUsername, adminPassword := utils.GetAdminCredentials(d.Type, i.URL)

	i.AdminCredential = models.Credential{
		Username: adminUsername,
		Password: adminPassword,
	}

	userCreds, err := w.createUserCredentials(i, d)
	if err != nil {
		return fmt.Errorf("could not create user credentials: %w", err)
	}

	// update the instance to hold credentials and new state
	i, err = w.updateInstance(i, d.Canonical, userCreds)
	if err != nil {
		return fmt.Errorf("could not update instance: %w", err)
	}

	// obfuscate the admin credentials
	i.AdminCredential = models.Credential{}

	// notify the user that the service requested is available
	data, err := json.Marshal(i)
	if err != nil {
		return fmt.Errorf("could not parse data into json: %w", err)
	}

	r := strings.NewReader(string(data))
	_, err = http.Post(fmt.Sprintf("%s/%s", d.CallbackURL, d.Canonical), "application/json", r)
	if err != nil {
		return fmt.Errorf("could not call user's service: %w", err)
	}
	return nil
}

func (w *Worker) createUserCredentials(i models.Instance, d models.Deployment) (*models.Credential, error) {
	host := i.URL
	username, password := "user", utils.GeneratePassword(int64(rand.Int()))
	adminUsername, adminPassword := i.AdminCredential.Username, i.AdminCredential.Password

	cred := models.Credential{
		Username: username,
		Password: password,
	}

	switch d.Type {
	case models.SonarqubeService:
		err := sonarqube.CreateUser(host, username, password, adminUsername, adminPassword)
		if err != nil {
			return nil, fmt.Errorf("failed to create sonarqube user: %w", err)
		}
		return &cred, nil
	case models.JenkinsService:
		// TODO: Configure Jenkins so that it actually requires credentials
		// as it is it does not
		return &models.Credential{}, nil
	default:
		return nil, fmt.Errorf("service type %s not supported", d.Type.String())
	}
}

func (w *Worker) updateInstance(i models.Instance, dcan string, userCreds *models.Credential) (models.Instance, error) {
	i.UserCredential = *userCreds
	i.State = models.Running

	err := w.s.UpdateInstance(context.Background(), dcan, &i)
	if err != nil {
		return i, fmt.Errorf("could not persist user credential: %w", err)
	}
	return i, nil
}
