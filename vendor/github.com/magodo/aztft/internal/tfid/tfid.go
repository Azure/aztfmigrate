package tfid

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
	"github.com/magodo/aztft/internal/resmap"
)

type builderFunc func(*client.ClientBuilder, armid.ResourceId, string) (string, error)

var dynamicBuilders = map[string]builderFunc{
	"azurerm_active_directory_domain_service":                        buildActiveDirectoryDomainService,
	"azurerm_storage_object_replication":                             buildStorageObjectReplication,
	"azurerm_storage_share":                                          buildStorageShare,
	"azurerm_storage_container":                                      buildStorageContainer,
	"azurerm_storage_queue":                                          buildStorageQueue,
	"azurerm_storage_table":                                          buildStorageTable,
	"azurerm_key_vault_key":                                          buildKeyVaultKey,
	"azurerm_key_vault_secret":                                       buildKeyVaultSecret,
	"azurerm_key_vault_certificate":                                  buildKeyVaultCertificate,
	"azurerm_key_vault_certificate_contacts":                         buildKeyVaultCertificateContacts,
	"azurerm_key_vault_certificate_issuer":                           buildKeyVaultCertificateIssuer,
	"azurerm_key_vault_managed_storage_account":                      buildKeyVaultStorageAccount,
	"azurerm_key_vault_managed_storage_account_sas_token_definition": buildKeyVaultStorageAccountSasTokenDefinition,
	"azurerm_storage_blob":                                           buildStorageBlob,
	"azurerm_storage_share_directory":                                buildStorageShareDirectory,
	"azurerm_storage_share_file":                                     buildStorageShareFile,
	"azurerm_storage_table_entity":                                   buildStorageTableEntity,
	"azurerm_storage_data_lake_gen2_filesystem":                      buildStorageDfs,
	"azurerm_storage_data_lake_gen2_path":                            buildStorageDfsPath,
	"azurerm_api_management_api":                                     buildApiManagementApi,
	"azurerm_automation_job_schedule":                                buildAutomationJobSchedule,
}

func NeedsAPI(rt string) bool {
	_, ok := dynamicBuilders[rt]
	return ok
}

func DynamicBuild(id armid.ResourceId, rt string, cred azcore.TokenCredential, clientOpt arm.ClientOptions) (string, error) {
	id = id.Clone()

	importSpec, err := GetImportSpec(id, rt)
	if err != nil {
		return "", fmt.Errorf("getting import spec for %s as %s: %v", id, rt, err)
	}

	builder, ok := dynamicBuilders[rt]
	if !ok {
		return "", fmt.Errorf("unknown resource type: %q", rt)
	}

	b := &client.ClientBuilder{
		Cred:      cred,
		ClientOpt: clientOpt,
	}

	return builder(b, id, importSpec)
}

func StaticBuild(id armid.ResourceId, rt string) (string, error) {
	id = id.Clone()

	importSpec, err := GetImportSpec(id, rt)
	if err != nil {
		return "", fmt.Errorf("getting import spec for %s as %s: %v", id, rt, err)
	}

	rid, ok := id.(*armid.ScopedResourceId)
	if !ok {
		return id.String(), nil
	}

	lastItem := func(l []string) string {
		if len(l) == 0 {
			return ""
		}
		return l[len(l)-1]
	}

	switch rt {
	case "azurerm_app_service_slot_virtual_network_swift_connection":
		rid.AttrTypes[2] = "config"
	case "azurerm_app_service_virtual_network_swift_connection":
		rid.AttrTypes[1] = "config"
	case "azurerm_synapse_workspace_sql_aad_admin":
		rid.AttrTypes[1] = "sqlAdministrators"
	case "azurerm_monitor_diagnostic_setting":
		// input: <target id>/providers/Microsoft.Insights/diagnosticSettings/setting1
		// tfid : <target id>|setting1
		id = id.ParentScope()
		return id.String() + "|" + rid.Names()[0], nil

	case "azurerm_synapse_role_assignment":
		pid := id.Parent()
		if err := pid.Normalize(importSpec); err != nil {
			return "", fmt.Errorf("normalizing id %q for %q with import spec %q: %v", pid.String(), rt, importSpec, err)
		}
		return pid.String() + "|" + id.Names()[1], nil
	case "azurerm_postgresql_active_directory_administrator":
		pid := id.Parent()
		if err := pid.Normalize(importSpec); err != nil {
			return "", fmt.Errorf("normalizing id %q for %q with import spec %q: %v", pid.String(), rt, importSpec, err)
		}
		return pid.String(), nil
	case "azurerm_servicebus_namespace_network_rule_set":
		pid := id.Parent()
		if err := pid.Normalize(importSpec); err != nil {
			return "", fmt.Errorf("normalizing id %q for %q with import spec %q: %v", pid.String(), rt, importSpec, err)
		}
		return pid.String(), nil
	case "azurerm_iotcentral_application_network_rule_set":
		pid := id.Parent()
		if err := pid.Normalize(importSpec); err != nil {
			return "", fmt.Errorf("normalizing id %q for %q with import spec %q: %v", pid.String(), rt, importSpec, err)
		}
		return pid.String(), nil
	case "azurerm_role_definition":
		scopeId := id.Parent()
		if scopeId == nil {
			scopeId = id.ParentScope()
		}
		scopePart := scopeId.String()
		routePart := strings.TrimPrefix(id.String(), scopePart)
		return routePart + "|" + scopePart, nil
	case "azurerm_network_manager_deployment":
		managerId := id.Parent().Parent()
		if err := managerId.Normalize(importSpec); err != nil {
			return "", fmt.Errorf("normalizing id %q for %q with import spec %q: %v", managerId.String(), rt, importSpec, err)
		}
		return managerId.String() + "/commit|" + id.Names()[1] + "|" + id.Names()[2], nil
	// Porperty-like resources
	case "azurerm_disk_pool_iscsi_target_lun":
		return buildIdForPropertyLikeResource(id.Parent(), lastItem(id.Names()), "azurerm_disk_pool_iscsi_target", "azurerm_managed_disk", "/lun|")
	case "azurerm_disk_pool_managed_disk_attachment":
		return buildIdForPropertyLikeResource(id.Parent(), lastItem(id.Names()), "azurerm_disk_pool", "azurerm_managed_disk", "/managedDisk|")
	case "azurerm_nat_gateway_public_ip_association":
		return buildIdForPropertyLikeResource(id.Parent(), lastItem(id.Names()), "azurerm_nat_gateway", "azurerm_public_ip", "|")
	case "azurerm_nat_gateway_public_ip_prefix_association":
		return buildIdForPropertyLikeResource(id.Parent(), lastItem(id.Names()), "azurerm_nat_gateway", "azurerm_public_ip_prefix", "|")
	case "azurerm_network_interface_application_gateway_backend_address_pool_association":
		return buildIdForPropertyLikeResource(id.Parent(), lastItem(id.Names()), "fake_azurerm_network_interface_ipconfig", "fake_azurerm_application_gateway_backend_address_pool", "|")
	case "azurerm_network_interface_application_security_group_association":
		return buildIdForPropertyLikeResource(id.Parent(), lastItem(id.Names()), "fake_azurerm_network_interface_ipconfig", "azurerm_application_security_group", "|")
	case "azurerm_network_interface_backend_address_pool_association":
		return buildIdForPropertyLikeResource(id.Parent(), lastItem(id.Names()), "fake_azurerm_network_interface_ipconfig", "azurerm_lb_backend_address_pool", "|")
	case "azurerm_network_interface_nat_rule_association":
		return buildIdForPropertyLikeResource(id.Parent(), lastItem(id.Names()), "fake_azurerm_network_interface_ipconfig", "azurerm_lb_nat_rule", "|")
	case "azurerm_network_interface_security_group_association":
		return buildIdForPropertyLikeResource(id.Parent(), lastItem(id.Names()), "azurerm_network_interface", "azurerm_network_security_group", "|")
	case "azurerm_virtual_desktop_workspace_application_group_association":
		return buildIdForPropertyLikeResource(id.Parent(), lastItem(id.Names()), "azurerm_virtual_desktop_workspace", "azurerm_virtual_desktop_application_group", "|")
	case "azurerm_subnet_nat_gateway_association",
		"azurerm_subnet_network_security_group_association",
		"azurerm_subnet_route_table_association":
		return id.Parent().String(), nil
	case "azurerm_iothub_endpoint_cosmosdb_account",
		"azurerm_iothub_endpoint_eventhub",
		"azurerm_iothub_endpoint_servicebus_queue",
		"azurerm_iothub_endpoint_servicebus_topic",
		"azurerm_iothub_endpoint_storage_container":
		id := id.(*armid.ScopedResourceId)
		id.AttrTypes[len(id.AttrTypes)-1] = "endpoints"
		return id.String(), nil
	case "azurerm_api_management_api_operation_policy",
		"azurerm_api_management_api_policy",
		"azurerm_api_management_policy",
		"azurerm_api_management_product_policy":
		pid := id.Parent()
		if err := pid.Normalize(importSpec); err != nil {
			return "", fmt.Errorf("normalizing id %q for %q with import spec %q: %v", pid.String(), rt, importSpec, err)
		}
		return pid.String(), nil
	case "azurerm_netapp_account_encryption":
		pid := id.Parent()
		if err := pid.Normalize(importSpec); err != nil {
			return "", fmt.Errorf("normalizing id %q for %q with import spec %q: %v", pid.String(), rt, importSpec, err)
		}
		return pid.String(), nil
	case "azurerm_storage_blob_inventory_policy":
		pid := id.Parent()
		if err := pid.Normalize(importSpec); err != nil {
			return "", fmt.Errorf("normalizing id %q for %q with import spec %q: %v", pid.String(), rt, importSpec, err)
		}
		return pid.String(), nil
	case "azurerm_container_app_environment_custom_domain":
		pid := id.Parent()
		if err := pid.Normalize(importSpec); err != nil {
			return "", fmt.Errorf("normalizing id %q for %q with import spec %q: %v", pid.String(), rt, importSpec, err)
		}
		return pid.String(), nil
	case "azurerm_role_management_policy":
		parentScopeId := id.ParentScope()
		return id.String() + "|" + parentScopeId.String(), nil
	}

	if importSpec != "" {
		if err := rid.Normalize(importSpec); err != nil {
			return "", fmt.Errorf("normalizing id %q for %q with import spec %q: %v", id.String(), rt, importSpec, err)
		}
	}
	return id.String(), nil
}

func GetImportSpec(id armid.ResourceId, rt string) (string, error) {
	resmap.Init()
	m := resmap.TF2ARMIdMap
	_ = m
	item, ok := resmap.TF2ARMIdMap[rt]
	if !ok {
		return "", fmt.Errorf("unknown resource type %q", rt)
	}

	if id.ParentScope() == nil {
		// For root scope resource id, the import spec is guaranteed to be only one.
		return item.ManagementPlane.ImportSpecs[0], nil
	}

	switch len(item.ManagementPlane.ImportSpecs) {
	case 0:
		// The ID is dynamically built (e.g. for property-like or some of the data plane only resources)
		return "", nil
	case 1:
		return item.ManagementPlane.ImportSpecs[0], nil
	default:
		// Needs to be matched with the scope. Or there might be zero import spec, as for the hypothetic resource ids.
		idscope := id.ParentScope().ScopeString()
		i := -1
		for idx, scope := range item.ManagementPlane.ParentScopes {
			if strings.EqualFold(scope, idscope) {
				i = idx
			}
		}
		if i == -1 {
			return "", fmt.Errorf("id %q doesn't correspond to resource type %q", id, rt)
		}
		return item.ManagementPlane.ImportSpecs[i], nil
	}
}

func buildIdForPropertyLikeResource(mainId armid.ResourceId, secondaryIdEnc string, mainRt, propRt, sep string) (string, error) {
	mainTFId, err := StaticBuild(mainId, mainRt)
	if err != nil {
		return "", fmt.Errorf("building resource id for %q: %v", mainId, err)
	}
	b, err := base64.StdEncoding.DecodeString(secondaryIdEnc)
	if err != nil {
		return "", fmt.Errorf("base64 decoding resource id %q: %v", secondaryIdEnc, err)
	}
	secondaryId, err := armid.ParseResourceId(string(b))
	if err != nil {
		return "", fmt.Errorf("parsing resource id %q: %v", string(b), err)
	}
	secondaryTFId, err := StaticBuild(secondaryId, propRt)
	if err != nil {
		return "", fmt.Errorf("building resource id for %q: %v", secondaryId, err)
	}
	return mainTFId + sep + secondaryTFId, nil
}
