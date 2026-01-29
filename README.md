# Terraform Provider Helmfile

A Terraform provider for managing Helmfile releases. This provider allows you to
manage Helm releases declaratively through Helmfile within your Terraform
infrastructure as code.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features

- **Declarative Helmfile Management**: Manage Helmfile releases as Terraform
  resources
- **Full Helmfile Support**: Access all helmfile commands and options through
  Terraform
- **State Management**: Terraform tracks the state of your Helmfile releases
- **Environment Support**: Target specific Helmfile environments
- **Selector Support**: Use label selectors to target specific releases
- **Flexible Configuration**: Customize helm binary paths, kustomize, and other
  settings

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21 (for building from source)
- [Helm](https://helm.sh/docs/intro/install/) >= 3.0
- [Helmfile](https://github.com/helmfile/helmfile) >= 0.150.0

## Installation

### Using Terraform Registry

Add the provider to your Terraform configuration:

```hcl
terraform {
  required_providers {
    helmfile = {
      source  = "kaweezle/helmfile"
      version = "~> 0.1.0"
    }
  }
}

provider "helmfile" {
  helm_binary_path = "/usr/local/bin/helm"
  perform_init     = true
}
```

### Building from Source

```bash
git clone https://github.com/kaweezle/terraform-provider-helmfile.git
cd terraform-provider-helmfile
go build -o terraform-provider-helmfile
```

## Usage

### Basic Example

```hcl
provider "helmfile" {
  helm_binary_path = "/usr/local/bin/helm"
  perform_init     = true
}

resource "helmfile_release" "nginx" {
  name         = "nginx"
  file_or_path = "./helmfile.yaml"
  environment  = "production"

  values = [
    "values.yaml",
    "values-prod.yaml"
  ]
}
```

### Advanced Example

```hcl
resource "helmfile_release" "app_stack" {
  name         = "app-stack"
  file_or_path = "./helmfiles"
  environment  = "staging"
  namespace    = "applications"
  kube_context = "staging-cluster"

  selector = [
    "app=myapp",
    "tier=backend"
  ]

  values = [
    "common-values.yaml",
    "staging-values.yaml"
  ]

  set = [
    "image.tag=v1.2.3",
    "replicas=3"
  ]

  concurrency          = 2
  suppress_secrets     = true
  skip_diff_on_install = false
  wait                 = true
  wait_for_jobs        = true
}
```

## Provider Configuration

The provider supports the following configuration options:

| Name                    | Type         | Description                                | Default     |
| ----------------------- | ------------ | ------------------------------------------ | ----------- |
| `perform_init`          | bool         | Run 'helmfile init' before operations      | `false`     |
| `helm_binary_path`      | string       | Path to helm binary                        | System PATH |
| `kustomize_binary_path` | string       | Path to kustomize binary                   | System PATH |
| `default_args`          | list(string) | Default args for all helmfile commands     | `[]`        |
| `skip_refresh`          | bool         | Skip 'helmfile repos' before operations    | `false`     |
| `skip_deps`             | bool         | Skip helm repo update and dependency build | `false`     |
| `disable_force_update`  | bool         | Don't force helm repos to update           | `false`     |
| `additional_plugins`    | list(string) | Additional helm plugins to install         | `[]`        |

## Resource Configuration

The `helmfile_release` resource supports numerous options. Key attributes
include:

- **Required**:
  - `name` - Release name
  - `file_or_path` - Path to helmfile.yaml or directory

- **Optional**:
  - `environment` - Helmfile environment
  - `namespace` - Kubernetes namespace
  - `kube_context` - kubectl context
  - `selector` - Label selectors
  - `values` - Value files
  - `set` - Set values
  - `wait`, `wait_for_jobs` - Wait configurations
  - `suppress_secrets` - Hide secrets in diff output
  - And many more (see [full documentation](docs/resources/release.md))

## Examples

Check out the [examples](examples/) directory for more usage patterns:

- [Complete Example](examples/complete/) - Full-featured example with multiple
  configurations
- [Provider Setup](examples/provider/) - Provider configuration examples
- [Resource Examples](examples/resources/helmfile_release/) - Various resource
  configurations

## Development

### Building

```bash
go build -o terraform-provider-helmfile
```

### Testing

```bash
go test -v ./...
```

### Running with Local Provider

Create a `~/.terraformrc` file:

```hcl
provider_installation {
  dev_overrides {
    "kaweezle/helmfile" = "/path/to/your/terraform-provider-helmfile"
  }
  direct {}
}
```

## Documentation

Full documentation is available in the [docs](docs/) directory:

- [Provider Documentation](docs/index.md)
- [Resource: helmfile_release](docs/resources/release.md)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major
changes, please open an issue first to discuss what you would like to change.

### Development Setup

1. Fork the repository
2. Clone your fork
3. Create a feature branch
4. Make your changes
5. Add tests
6. Run tests: `go test ./...`
7. Submit a pull request

## Troubleshooting

### Provider Not Found

Make sure the provider is properly installed:

```bash
terraform init
```

### Helmfile Command Fails

Enable debug mode for more verbose output:

```hcl
resource "helmfile_release" "example" {
  # ...
  debug     = true
  log_level = "debug"
}
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file
for details.

## Acknowledgments

- Built using the
  [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework)
- Integrates with [Helmfile](https://github.com/helmfile/helmfile)
- Inspired by the
  [terraform-provider-helmfile](https://github.com/mumoshu/terraform-provider-helmfile)
  by mumoshu

## Support

- 📖 [Documentation](docs/)
- 🐛
  [Issue Tracker](https://github.com/kaweezle/terraform-provider-helmfile/issues)
- 💬
  [Discussions](https://github.com/kaweezle/terraform-provider-helmfile/discussions)

## Project Status

This project is under active development. Contributions and feedback are
welcome!
