package tf

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-exec/tfexec"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/ms-henglu/azurerm-restapi-to-azurerm/types"
)

type Terraform struct {
	exec             *tfexec.Terraform
	LogEnabled       bool
	workingDirectory string
}

const planfile = "tfplan"

// The required terraform version that has the `terraform add` command.
var minRequiredTFVersion = version.Must(version.NewSemver("v1.1.0-alpha20210630"))
var maxRequiredTFVersion = version.Must(version.NewSemver("v1.1.0-alpha20211006"))

func NewTerraform(workingDirectory string, logEnabled bool) (*Terraform, error) {
	execPath, err := FindTerraform(context.TODO(), minRequiredTFVersion, maxRequiredTFVersion)
	if err != nil {
		return nil, fmt.Errorf("error finding a terraform exectuable: %w", err)
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
		t.exec.SetStdout(ioutil.Discard)
		t.exec.SetStderr(ioutil.Discard)
		t.exec.SetLogger(log.New(ioutil.Discard, "", 0))
	}
}

func (t *Terraform) Init() error {
	if _, err := os.Stat(path.Join(t.GetWorkingDirectory(), ".terraform")); os.IsNotExist(err) {
		err := t.exec.Init(context.Background(), tfexec.Upgrade(false))
		// ignore the error if can't find azurerm-restapi
		if err != nil && strings.Contains(err.Error(), "Azure/azurerm-restapi: provider registry registry.terraform.io does not have") {
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
		if resourceChange == nil || resourceChange.Change == nil || resourceChange.Type != "azurerm-restapi_resource" {
			continue
		}
		if resourceChange.Index == nil {
			resourceMap[resourceChange.Address] = &types.GenericResource{
				Label: resourceChange.Name,
				Instances: []types.Instance{
					{
						ResourceId: getResourceId(resourceChange.Change.Before),
						ApiVersion: getApiVersion(resourceChange.Change.Before),
					},
				},
			}
		} else {
			address := fmt.Sprintf("%s.%s", resourceChange.Type, resourceChange.Name)
			if _, ok := resourceMap[address]; ok {
				resourceMap[address].Instances = append(resourceMap[address].Instances, types.Instance{
					Index:      resourceChange.Index,
					ResourceId: getResourceId(resourceChange.Change.Before),
					ApiVersion: getApiVersion(resourceChange.Change.Before),
				})
			} else {
				resourceMap[address] = &types.GenericResource{
					Label:        resourceChange.Name,
					ResourceType: "",
					Instances: []types.Instance{
						{
							Index:      resourceChange.Index,
							ResourceId: getResourceId(resourceChange.Change.Before),
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
				prop := strings.TrimPrefix(output.OldName, fmt.Sprintf("jsondecode(%s.output).", resource.OldAddress(instance.Index)))
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

func (t *Terraform) ListGenericPatchResources(p *tfjson.Plan) []types.GenericPatchResource {
	resources := make([]types.GenericPatchResource, 0)
	if p == nil {
		return resources
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
	tpl, err := t.exec.Add(context.TODO(), address, tfexec.FromState(true))
	// remove comments
	tpl = tpl[strings.Index(tpl, `resource "`):]
	if err != nil {
		return "", fmt.Errorf("converting terraform state to config for resource %s: %w", address, err)
	}
	return tpl, nil
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
