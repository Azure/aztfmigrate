// Package configschema is an adoption of a subset of the github.com/hashicorp/terraform/internal/configs/configschema@92574a7811111f6afef4a16c9b72b8bd53e882e1.
// It only focus on the implied type (and its dependencies) of the schema `Block` type. But instead of the `Block` defined internally by terraform core, it target
// to the github.com/hashicorp/terraform-json.SchemaBlock.
package jsonschema
