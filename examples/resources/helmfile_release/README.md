# Helmfile Release Examples

This directory contains example Terraform configurations for the
`helmfile_release` resource.

## Files

- `resource.tf` - Comprehensive examples showing different use cases
- `helmfile.yaml` - Sample helmfile configuration
- `values.yaml` - Sample values file for the helmfile

## Usage

1. Copy the example files to your Terraform project
2. Modify the `helmfile.yaml` to match your desired Helm releases
3. Update the `resource.tf` with your specific configuration
4. Run `terraform init` to initialize the provider
5. Run `terraform plan` to preview changes
6. Run `terraform apply` to deploy

## Examples Included

### Basic Example

The simplest possible configuration with just required fields.

### Comprehensive Example

Shows most common configuration options including:

- Environment and namespace configuration
- Selectors for targeting specific releases
- Diff and sync options
- Values files and inline values
- Wait and upgrade behavior

### Post-Renderer Example

Demonstrates using a post-renderer (like Kustomize) to transform manifests
before deployment.

### Selectors Example

Shows how to use selectors to target specific releases and include dependencies.

### Suppression Example

Demonstrates security-focused options for suppressing secrets and sensitive
output.

## Prerequisites

- Terraform >= 1.0
- Helm >= 3.0
- Helmfile >= 0.144.0
- helm-diff plugin

## Additional Resources

- [Helmfile Documentation](https://helmfile.readthedocs.io/)
- [Helm Documentation](https://helm.sh/docs/)
- [Provider Documentation](../../docs/index.md)
