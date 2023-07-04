[![PkgGoDev](https://pkg.go.dev/badge/github.com/magodo/tfstate)](https://pkg.go.dev/github.com/magodo/tfstate)

# tfstate

Helper types for Terraform state on top of https://github.com/hashicorp/terraform-json.

## Why

Currently, there is no good way to interact with Terraform state. This makes sense as Terraform core intentionally hide the details of the state as it is evolving quickly. If you look into the Terraform core codebase, you'll notce that [there are multiple formats defined for states along the lifetime of Terraform](https://github.com/hashicorp/terraform/tree/d3e7c5e8a9162617f9cc12ccd005347978825d02/internal/states/statefile).

Whilst for some reason, developers still want a way to inspect the Terraform state file via some means. Currently, the correct way to do so is via https://github.com/hashicorp/terraform-exec. Where it provides a method [`Show()`](https://pkg.go.dev/github.com/hashicorp/terraform-exec@v0.16.0/tfexec#Terraform.Show), that returns you a [`tfjson.State`](https://pkg.go.dev/github.com/hashicorp/terraform-json#State).

Everything works just fine, the only problem is for each resource instance inside `tfjson.State`, its main content [`AttributeValues`](https://pkg.go.dev/github.com/hashicorp/terraform-json#StateResource) is of type `map[string]interface{}`. This makes the user can hardly do some fancy inspection on the resource attributes, as they are not typed.

This package aims to fix this last gap by defining a thin wrapper `tfstate.State` around the `tfjson.State`, which has almost the same structure, except the `AttributeValues` is replaced with `Value`, which is of type `cty.Value`.

## Note

This package only works for the V4 format of state file, which is the used since Terraform v0.12.

## Example

See: https://github.com/magodo/tfstate/blob/main/state_example_test.go.