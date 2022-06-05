package terraform

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-exec/tfexec"
)

type Worker struct {
	execPath string
}

func NewWorker(execPath string) *Worker {
	return &Worker{execPath}
}

func (w *Worker) Deploy(dir string) error {
	ctx := context.Background()

	tf, err := tfexec.NewTerraform(dir, w.execPath)
	if err != nil {
		return fmt.Errorf("could not create terraform: %w", err)
	}

	err = tf.Init(ctx)
	if err != nil {
		return fmt.Errorf("could not init terraform: %w", err)
	}

	changed, err := tf.Plan(ctx)
	if err != nil {
		return fmt.Errorf("could not plan deployment: %w", err)
	}

	// in case there are no changes we don't apply the deployment
	if !changed {
		return nil
	}

	err = tf.Apply(ctx)
	if err != nil {
		return fmt.Errorf("could not apply deployment plan: %w", err)
	}
	return nil
}

func (w *Worker) Outputs(dir string) (map[string]string, error) {
	ctx := context.Background()

	tf, err := tfexec.NewTerraform(dir, w.execPath)
	if err != nil {
		return nil, fmt.Errorf("could not create terraform: %w", err)
	}

	outputs, err := tf.Output(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get the outputs: %w", err)
	}

	om := make(map[string]string)
	for key, output := range outputs {
		om[key] = string(output.Value)
	}
	return om, nil
}

func (w *Worker) Teardown(dir string) error {
	ctx := context.Background()

	tf, err := tfexec.NewTerraform(dir, w.execPath)
	if err != nil {
		return fmt.Errorf("could not create terraform: %w", err)
	}

	err = tf.Destroy(ctx)
	if err != nil {
		return fmt.Errorf("could not destroy deployment: %w", err)
	}
	return nil
}
