package tf

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"

	"github.com/hashicorp/terraform-exec/tfexec"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/magodo/tfadd/tfadd"
)

type Terraform struct {
	exec             *tfexec.Terraform
	LogEnabled       bool
	workingDirectory string
}

const planfile = "tfplan"

func NewTerraform(workingDirectory string, logEnabled bool) (*Terraform, error) {
	execPath, err := FindTerraform(context.Background())
	if err != nil {
		return nil, err
	}
	tf, err := tfexec.NewTerraform(workingDirectory, execPath)
	if err != nil {
		return nil, err
	}

	t := &Terraform{
		exec:             tf,
		workingDirectory: workingDirectory,
		LogEnabled:       logEnabled,
	}
	t.SetLogEnabled(true)
	return t, nil
}

func (t *Terraform) SetLogEnabled(enabled bool) {
	if enabled && t.LogEnabled {
		t.exec.SetStdout(os.Stdout)
		t.exec.SetStderr(os.Stderr)
		t.exec.SetLogger(log.New(os.Stdout, "", 0))
	} else {
		t.exec.SetStdout(io.Discard)
		t.exec.SetStderr(io.Discard)
		t.exec.SetLogger(log.New(io.Discard, "", 0))
	}
}

func (t *Terraform) Init() error {
	if _, err := os.Stat(path.Join(t.GetWorkingDirectory(), ".terraform")); os.IsNotExist(err) {
		err := t.exec.Init(context.Background(), tfexec.Upgrade(false))
		// ignore the error if can't find azapi
		if err != nil && strings.Contains(err.Error(), "Azure/azapi: provider registry registry.terraform.io does not have") {
			return nil
		}
		return err
	}
	return nil
}

func (t *Terraform) Show() (*tfjson.State, error) {
	return t.exec.Show(context.TODO())
}

func (t *Terraform) Plan(varFile *string) (*tfjson.Plan, error) {
	planOptions := []tfexec.PlanOption{tfexec.Out(planfile)}
	if varFile != nil && *varFile != "" {
		planOptions = append(planOptions, tfexec.VarFile(*varFile))
	}
	_, err := t.exec.Plan(context.TODO(), planOptions...)
	if err != nil {
		return nil, err
	}
	t.SetLogEnabled(false)
	p, err := t.exec.ShowPlanFile(context.TODO(), planfile)
	t.SetLogEnabled(true)
	return p, err
}

func (t *Terraform) ImportAdd(address string, id string) (string, error) {
	_ = t.Init()
	err := t.exec.Import(context.TODO(), address, id)
	if err != nil {
		return "", fmt.Errorf("importing resource %s: %w", address, err)
	}
	outputs, err := tfadd.StateForTargets(context.TODO(), t.exec, []string{address}, tfadd.Full(true))
	if err != nil {
		return "", fmt.Errorf("converting terraform state to config for resource %s: %w", address, err)
	}
	if len(outputs) == 0 {
		return "", fmt.Errorf("resource %s not found in state", address)
	}
	return string(outputs[0]), nil
}

func (t *Terraform) Import(address string, id string) error {
	_ = t.Init()
	return t.exec.Import(context.TODO(), address, id)
}

func (t *Terraform) Apply() error {
	return t.exec.Apply(context.TODO())
}

func (t *Terraform) Destroy() error {
	return t.exec.Destroy(context.TODO())
}

func (t *Terraform) GetExec() *tfexec.Terraform {
	return t.exec
}

func (t *Terraform) GetWorkingDirectory() string {
	return t.workingDirectory
}
