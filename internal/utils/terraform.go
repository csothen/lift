package utils

import (
	"context"
	"log"

	"github.com/csothen/lift/internal/config"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
)

// in case no terraform executable path was provided we will install it
// and pass the new executable path to the config
func InstallTerraform(cfg *config.Config) {
	if cfg.TerraformExecPath == "" {
		// we will install version 1.2.1 of Terraform
		installer := &releases.ExactVersion{
			Product: product.Terraform,
			Version: version.Must(version.NewVersion("1.2.1")),
		}

		execPath, err := installer.Install(context.Background())
		if err != nil {
			log.Fatal("error installing Terraform: %w", err)
		}
		cfg.TerraformExecPath = execPath
	}
}
