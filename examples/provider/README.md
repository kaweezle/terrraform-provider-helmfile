# Helmfile Provider Examples

This directory contains example Terraform configurations for the Helmfile
provider.

## Files

- `provider.tf` - Provider configuration examples

## Basic Usage

```hcl
terraform {
  required_providers {
    helmfile = {
      source = "kaweezle/helmfile"
    }
  }
}

provider "helmfile" {
  perform_init = true
}
```

## Configuration Options

### Minimal Configuration

The provider can be used with minimal configuration, relying on defaults:

```hcl
provider "helmfile" {}
```

This will:

- Use `helm` and `kustomize` from system PATH
- Not perform initialization automatically
- Use default error message handling
- Not force repository updates

### Development Configuration

For development environments, you might want more verbose output:

```hcl
provider "helmfile" {
  perform_init = true
  default_args = ["--debug"]

  additional_plugins = [
    "https://github.com/databus23/helm-diff",
    "https://github.com/jkroepke/helm-secrets"
  ]
}
```

### Production Configuration

For production, focus on security and reliability:

```hcl
provider "helmfile" {
  perform_init                    = true
  strip_args_values_on_exit_error = true
  enforce_plugin_verification     = true
  skip_refresh                    = false  # Always refresh repos

  additional_plugins = [
    "https://github.com/databus23/helm-diff"
  ]
}
```

### CI/CD Configuration

For CI/CD pipelines, optimize for speed and security:

```hcl
provider "helmfile" {
  perform_init                    = true
  skip_refresh                    = true  # Repos already updated in pipeline
  skip_deps                       = true  # Dependencies handled separately
  strip_args_values_on_exit_error = true

  default_args = ["--no-color"]  # Better for CI logs
}
```

### Custom Binary Paths

If you have Helm or Kustomize installed in non-standard locations:

```hcl
provider "helmfile" {
  helm_binary_path      = "/opt/helm/bin/helm"
  kustomize_binary_path = "/opt/kustomize/bin/kustomize"
}
```

### OCI Registry Support

For working with OCI-based Helm registries:

```hcl
provider "helmfile" {
  helm_oci_plain_http = false  # Use HTTPS (recommended)
  # helm_oci_plain_http = true  # Use only for local testing
}
```

## Environment Variables

The provider respects the following environment variables:

- `KUBECONFIG` - Path to Kubernetes configuration file
- `HELM_CACHE_HOME` - Helm cache directory
- `HELM_CONFIG_HOME` - Helm configuration directory
- `HELM_DATA_HOME` - Helm data directory

## Prerequisites

Before using the provider, ensure you have:

1. **Helm** installed (version 3.x or later)

   ```bash
   helm version
   ```

2. **Helmfile** installed (version 0.144.0 or later)

   ```bash
   helmfile version
   ```

3. **helm-diff plugin** installed

   ```bash
   helm plugin install https://github.com/databus23/helm-diff
   ```

4. **kubectl** configured with access to your cluster
   ```bash
   kubectl cluster-info
   ```

## Troubleshooting

### Plugin Installation Failures

If plugins fail to install automatically, install them manually:

```bash
helm plugin install https://github.com/databus23/helm-diff
```

Then set `perform_init = false` in the provider configuration.

### Binary Not Found

If the provider cannot find Helm or Kustomize binaries, specify their paths
explicitly:

```hcl
provider "helmfile" {
  helm_binary_path      = "/usr/local/bin/helm"
  kustomize_binary_path = "/usr/local/bin/kustomize"
}
```

### Authentication Issues

Ensure your kubeconfig is properly configured and you have the necessary
permissions:

```bash
kubectl auth can-i create deployments --all-namespaces
```

## Additional Resources

- [Provider Documentation](../../docs/index.md)
- [Helmfile Documentation](https://helmfile.readthedocs.io/)
- [Helm Documentation](https://helm.sh/docs/)
- [Terraform Provider Development](https://developer.hashicorp.com/terraform/plugin)
