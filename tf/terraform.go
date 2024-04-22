package tf

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"

	"github.com/Azure/azapi2azurerm/types"
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

func (t *Terraform) Plan() (*tfjson.Plan, error) {
	_, err := t.exec.Plan(context.TODO(), tfexec.Out(planfile))
	if err != nil {
		return nil, err
	}
	t.SetLogEnabled(false)
	p, err := t.exec.ShowPlanFile(context.TODO(), planfile)
	t.SetLogEnabled(true)
	return p, err
}

func (t *Terraform) ListGenericResources(p *tfjson.Plan) []types.GenericResource {
	resources := make([]types.GenericResource, 0)
	if p == nil {
		return resources
	}
	resourceMap := make(map[string]*types.GenericResource)
	for _, resourceChange := range p.ResourceChanges {
		if resourceChange == nil || resourceChange.Change == nil || resourceChange.Type != "azapi_resource" {
			continue
		}
		if resourceChange.Index == nil {
			resourceMap[resourceChange.Address] = &types.GenericResource{
				Label: resourceChange.Name,
				Instances: []types.Instance{
					{
						ResourceId: getId(resourceChange.Change.Before),
						ApiVersion: getApiVersion(resourceChange.Change.Before),
					},
				},
			}
		} else {
			address := fmt.Sprintf("%s.%s", resourceChange.Type, resourceChange.Name)
			if _, ok := resourceMap[address]; ok {
				resourceMap[address].Instances = append(resourceMap[address].Instances, types.Instance{
					Index:      resourceChange.Index,
					ResourceId: getId(resourceChange.Change.Before),
					ApiVersion: getApiVersion(resourceChange.Change.Before),
				})
			} else {
				resourceMap[address] = &types.GenericResource{
					Label:        resourceChange.Name,
					ResourceType: "",
					Instances: []types.Instance{
						{
							Index:      resourceChange.Index,
							ResourceId: getId(resourceChange.Change.Before),
							ApiVersion: getApiVersion(resourceChange.Change.Before),
						},
					},
				}
			}
		}
	}
	for _, v := range resourceMap {
		resources = append(resources, *v)
	}
	refValueMap := getRefValueMap(p)
	for index, resource := range resources {
		resources[index].References = getReferencesForAddress(resource.OldAddress(nil), p, refValueMap)
		resources[index].InputProperties = getInputProperties(resource.OldAddress(nil), p)
	}
	for i, resource := range resources {
		outputPropSet := make(map[string]bool)
		for j, instance := range resource.Instances {
			resources[i].Instances[j].Outputs = getOutputsForAddress(resource.OldAddress(instance.Index), refValueMap)
			for _, output := range resources[i].Instances[j].Outputs {
				prop := strings.TrimPrefix(output.OldName, fmt.Sprintf("%s.output.", resource.OldAddress(instance.Index)))
				prop = strings.TrimPrefix(prop, fmt.Sprintf("jsondecode(%s.output).", resource.OldAddress(instance.Index)))
				if strings.HasPrefix(prop, "identity.userAssignedIdentities") {
					prop = "identity.userAssignedIdentities"
				}
				outputPropSet[prop] = true
			}
		}
		resources[i].OutputProperties = make([]string, 0)
		for key := range outputPropSet {
			resources[i].OutputProperties = append(resources[i].OutputProperties, key)
		}
	}
	return resources
}

func (t *Terraform) ListGenericUpdateResources(p *tfjson.Plan) []types.GenericUpdateResource {
	resources := make([]types.GenericUpdateResource, 0)
	if p == nil {
		return resources
	}
	idMap := make(map[string]*tfjson.ResourceChange)
	for _, resourceChange := range p.ResourceChanges {
		if resourceChange == nil || resourceChange.Change == nil || strings.Contains(resourceChange.Type, "azapi") {
			continue
		}
		idMap[getId(resourceChange.Change.Before)] = resourceChange
	}

	for _, resourceChange := range p.ResourceChanges {
		if resourceChange == nil || resourceChange.Change == nil || resourceChange.Type != "azapi_update_resource" {
			continue
		}
		resourceId := getId(resourceChange.Change.Before)
		if idMap[resourceId] == nil {
			log.Printf("[WARN] resource azapi_update_resource.%s's target is not in the same terraform working directory", resourceChange.Name)
			continue
		}
		rc := idMap[resourceId]
		resources = append(resources, types.GenericUpdateResource{
			OldLabel:     resourceChange.Name,
			Label:        rc.Name,
			Id:           resourceId,
			ApiVersion:   getApiVersion(resourceChange.Change.Before),
			ResourceType: rc.Type,
			Change:       rc.Change,
		})
	}
	refValueMap := getRefValueMap(p)
	for index, resource := range resources {
		resources[index].References = getReferencesForAddress(resource.OldAddress(), p, refValueMap)
		resources[index].InputProperties = getInputProperties(resource.OldAddress(), p)
	}
	for i, resource := range resources {
		resources[i].Outputs = getOutputsForAddress(resource.OldAddress(), refValueMap)
		resources[i].OutputProperties = make([]string, 0)
		for _, output := range resources[i].Outputs {
			resources[i].OutputProperties = append(resources[i].OutputProperties, strings.TrimPrefix(output.OldName, fmt.Sprintf("jsondecode(%s.output).", resource.OldAddress())))
		}
	}
	return resources
}

func (t *Terraform) ImportAdd(address string, id string) (string, error) {
	_ = t.Init()
	err := t.exec.Import(context.TODO(), address, id)
	if err != nil {
		log.Fatal(err)
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

func (t *Terraform) RefreshState(resources []string) error {
	// TODO: replace refresh command with apply -refresh-only
	opts := make([]tfexec.RefreshCmdOption, 0)
	for _, res := range resources {
		opts = append(opts, tfexec.Target(res))
	}
	return t.exec.Refresh(context.TODO(), opts...)
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
