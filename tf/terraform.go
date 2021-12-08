package tf

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-exec/tfexec"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/ms-henglu/azurerm-restapi-to-azurerm/types"
)

type Terraform struct {
	exec tfexec.Terraform
}

const planfile = "tfplan"

// The required terraform version that has the `terraform add` command.
var minRequiredTFVersion = version.Must(version.NewSemver("v1.1.0-alpha20210630"))
var maxRequiredTFVersion = version.Must(version.NewSemver("v1.1.0-alpha20211006"))
var LogEnabled = false

func NewTerraform(workingDirectory string) (*Terraform, error) {
	execPath, err := FindTerraform(context.TODO(), minRequiredTFVersion, maxRequiredTFVersion)
	if err != nil {
		return nil, fmt.Errorf("error finding a terraform exectuable: %w", err)
	}
	if err != nil {
		return nil, err
	}
	tf, err := tfexec.NewTerraform(workingDirectory, execPath)
	if err != nil {
		return nil, err
	}

	t := &Terraform{
		exec: *tf,
	}
	t.SetLogEnabled(true)
	return t, nil
}

func (t *Terraform) SetLogEnabled(enabled bool) {
	if enabled && LogEnabled {
		t.exec.SetStdout(os.Stdout)
		t.exec.SetStderr(os.Stderr)
		t.exec.SetLogger(log.New(os.Stdout, "", 0))
	} else {
		t.exec.SetStdout(ioutil.Discard)
		t.exec.SetStderr(ioutil.Discard)
		t.exec.SetLogger(log.New(ioutil.Discard, "", 0))
	}
}

func (t *Terraform) Init() error {
	if _, err := os.Stat(".terraform"); os.IsNotExist(err) {
		err := t.exec.Init(context.Background(), tfexec.Upgrade(false))
		// ignore the error if can't find azurerm-restapi
		if err != nil && strings.Contains(err.Error(), "Azure/azurerm-restapi: provider registry registry.terraform.io does not have") {
			return nil
		}
		return err
	}
	log.Println("[INFO] skip running init command because .terraform folder exist")
	return nil
}

func (t *Terraform) Show() (*tfjson.State, error) {
	return t.exec.Show(context.TODO())
}

func (t *Terraform) Plan() (*tfjson.Plan, error) {
	ok, err := t.exec.Plan(context.TODO(), tfexec.Out(planfile))
	if err != nil {
		return nil, err
	}
	if !ok {
		// no changes
		return nil, nil
	}

	t.SetLogEnabled(false)
	p, err := t.exec.ShowPlanFile(context.TODO(), planfile)
	t.SetLogEnabled(true)
	return p, err
}

func (t *Terraform) ListGenericResources() ([]types.GenericResource, error) {
	_, err := t.exec.Plan(context.TODO(), tfexec.Out(planfile))
	if err != nil {
		return nil, err
	}
	t.SetLogEnabled(false)
	p, err := t.exec.ShowPlanFile(context.TODO(), planfile)
	t.SetLogEnabled(true)

	resources := make([]types.GenericResource, 0)
	if p == nil {
		return resources, nil
	}
	for _, resourceChange := range p.ResourceChanges {
		if resourceChange == nil || resourceChange.Change == nil || resourceChange.Type != "azurerm-restapi_resource" {
			continue
		}
		resources = append(resources, types.GenericResource{
			Label:        resourceChange.Name,
			Id:           getResourceId(resourceChange.Change.Before),
			ResourceType: "",
		})
	}
	refValueMap := getRefValueMap(p)
	for index, resource := range resources {
		resources[index].References = getReferencesForAddress(resource.OldAddress(), p, refValueMap)
	}
	for i, resource := range resources {
		resources[i].Outputs = getOutputsForAddress(resource.OldAddress(), refValueMap)
	}
	return resources, err
}

func (t *Terraform) ListGenericPatchResources() ([]types.GenericPatchResource, error) {
	_, err := t.exec.Plan(context.TODO(), tfexec.Out(planfile))
	if err != nil {
		return nil, err
	}
	t.SetLogEnabled(false)
	p, err := t.exec.ShowPlanFile(context.TODO(), planfile)
	t.SetLogEnabled(true)

	resources := make([]types.GenericPatchResource, 0)
	if p == nil {
		return resources, nil
	}
	idMap := make(map[string]*tfjson.ResourceChange)
	for _, resourceChange := range p.ResourceChanges {
		if resourceChange == nil || resourceChange.Change == nil || strings.Contains(resourceChange.Type, "azurerm-restapi") {
			continue
		}
		idMap[getId(resourceChange.Change.Before)] = resourceChange
	}

	for _, resourceChange := range p.ResourceChanges {
		if resourceChange == nil || resourceChange.Change == nil || resourceChange.Type != "azurerm-restapi_patch_resource" {
			continue
		}
		resourceId := getResourceId(resourceChange.Change.Before)
		if idMap[resourceId] == nil {
			log.Printf("[WARN] resource azurerm-restapi_patch_resource.%s's target is not in the same terraform working directory", resourceChange.Name)
			continue
		}
		rc := idMap[resourceId]
		resources = append(resources, types.GenericPatchResource{
			OldLabel:     resourceChange.Name,
			Label:        rc.Name,
			Id:           resourceId,
			ResourceType: rc.Type,
			Change:       rc.Change,
		})
	}
	refValueMap := getRefValueMap(p)
	for index, resource := range resources {
		resources[index].References = getReferencesForAddress(resource.OldAddress(), p, refValueMap)
	}
	for i, resource := range resources {
		resources[i].Outputs = getOutputsForAddress(resource.OldAddress(), refValueMap)
	}
	return resources, err
}

func (t *Terraform) Import(address string, id string) (string, error) {
	_ = t.Init()
	err := t.exec.Import(context.TODO(), address, id)
	if err != nil {
		log.Fatal(err)
	}
	tpl, err := t.exec.Add(context.TODO(), address, tfexec.FromState(true))
	// remove comments
	tpl = tpl[strings.Index(tpl, `resource "`):]
	if err != nil {
		return "", fmt.Errorf("converting terraform state to config for resource %s: %w", address, err)
	}
	return tpl, nil
}

func (t *Terraform) Apply() error {
	return t.exec.Apply(context.TODO())
}

func (t *Terraform) Destroy() error {
	return t.exec.Destroy(context.TODO())
}

func (t *Terraform) GetExec() tfexec.Terraform {
	return t.exec
}
