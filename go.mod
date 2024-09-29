module github.com/Azure/azapi2azurerm

go 1.22.0

toolchain go1.23.1

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.9.2
	github.com/gertd/go-pluralize v0.2.1
	github.com/hashicorp/go-version v1.7.0
	github.com/hashicorp/hc-install v0.9.0
	github.com/hashicorp/hcl/v2 v2.22.0
	github.com/hashicorp/terraform-exec v0.21.0
	github.com/hashicorp/terraform-json v0.22.1
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.34.0
	github.com/magodo/aztft v0.3.1-0.20240802001407-2834614c6d9c
	github.com/magodo/tfadd v0.10.1-0.20240902124619-bd18a56f410d
	github.com/mitchellh/cli v1.1.5
	github.com/zclconf/go-cty v1.15.0
)

require (
	dario.cat/mergo v1.0.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.6.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/alertsmanagement/armalertsmanagement v0.8.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/apimanagement/armapimanagement v1.1.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appcontainers/armappcontainers/v3 v3.0.0-beta.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/applicationinsights/armapplicationinsights v1.1.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appplatform/armappplatform v1.1.0-beta.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appservice/armappservice v1.0.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/automation/armautomation v0.8.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/botservice/armbotservice v1.1.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cdn/armcdn v1.1.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5 v5.7.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/costmanagement/armcostmanagement/v2 v2.0.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/datafactory/armdatafactory/v7 v7.0.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dataprotection/armdataprotection v1.0.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/datashare/armdatashare v1.1.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/desktopvirtualization/armdesktopvirtualization v1.0.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/devtestlabs/armdevtestlabs v1.1.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/digitaltwins/armdigitaltwins v1.1.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/domainservices/armdomainservices v1.1.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/frontdoor/armfrontdoor v1.2.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/hdinsight/armhdinsight v1.1.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/iothub/armiothub v1.1.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault v1.2.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/kusto/armkusto v1.3.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/logic/armlogic v1.1.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/machinelearning/armmachinelearning/v3 v3.1.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor v0.7.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/netapp/armnetapp v1.0.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork v1.1.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/operationalinsights/armoperationalinsights v1.1.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/paloaltonetworksngfw/armpanngfw v1.0.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/recoveryservices/armrecoveryservicesbackup v1.0.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/recoveryservices/armrecoveryservicessiterecovery v1.1.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armdeploymentscripts v1.0.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/securityinsights/armsecurityinsights/v2 v2.0.0-beta.3 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage v1.3.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storagecache/armstoragecache v1.0.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storagemover/armstoragemover v1.1.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storagepool/armstoragepool v1.1.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/streamanalytics/armstreamanalytics v1.1.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/synapse/armsynapse v0.7.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/timeseriesinsights/armtimeseriesinsights v1.1.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/workloads/armworkloads v1.1.0 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver/v3 v3.3.0 // indirect
	github.com/Masterminds/sprig/v3 v3.3.0 // indirect
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/apparentlymart/go-textseg/v15 v15.0.0 // indirect
	github.com/armon/go-radix v1.0.0 // indirect
	github.com/bgentry/speakeasy v0.2.0 // indirect
	github.com/fatih/color v1.17.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cty v1.4.1-0.20200414143053-d3edf31b6320 // indirect
	github.com/hashicorp/go-hclog v1.5.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/hashicorp/logutils v1.0.0 // indirect
	github.com/hashicorp/terraform-plugin-go v0.23.0 // indirect
	github.com/hashicorp/terraform-plugin-log v0.9.0 // indirect
	github.com/huandu/xstrings v1.5.0 // indirect
	github.com/magodo/armid v0.0.0-20230511151020-27880e5961c3 // indirect
	github.com/magodo/tfpluginschema v0.0.0-20240902090353-0525d7d8c1c2 // indirect
	github.com/magodo/tfstate v0.0.0-20240829105815-03d52976fa13 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/posener/complete v1.2.3 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/spf13/cast v1.7.0 // indirect
	github.com/vmihailenco/msgpack v4.0.4+incompatible // indirect
	github.com/vmihailenco/msgpack/v5 v5.4.1 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	golang.org/x/crypto v0.27.0 // indirect
	golang.org/x/mod v0.21.0 // indirect
	golang.org/x/net v0.29.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	golang.org/x/text v0.18.0 // indirect
	golang.org/x/tools v0.25.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/protobuf v1.34.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)
