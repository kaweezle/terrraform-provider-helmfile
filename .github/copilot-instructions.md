# Terraform Provider Helmfile Development Guide for AI Agents

## Project Overview

Terraform Provider Helmfile enables declarative management of Helmfile releases
through Terraform. Built using the Terraform Plugin Framework, it wraps the
helmfile library to execute helmfile operations within Terraform's lifecycle.

### Core Architecture

**Provider Structure:**

- `internal/provider/` - Provider implementation using Plugin Framework
  - `provider.go` - Main provider definition and configuration
  - `provider_helmfile/` - Provider logic and helmfile executor integration
  - `resource_helmfile_release/` - Release resource (CRUD operations)
- `pkg/helmfile/` - Helmfile execution abstraction layer
  - `executor.go` - Interface defining helmfile operations (Apply, Diff,
    Template, Destroy, Build)
  - `output_capture.go` - Zap logger integration for capturing helmfile output

**Code Generation:**

- Schema is defined in `provider-code-spec.json` (JSON schema)
- Generated code in `*_gen.go` files (DO NOT EDIT)
- Run `make generate` to regenerate from spec using `tfplugingen-framework`

## Development Workflows

### Building and Testing

```bash
# Build the provider
make build

# Install locally for testing
make install

# Run unit tests
make test

# Run acceptance tests (requires TF_ACC=1)
make testacc

# Lint code
make lint

# Format code
make fmt

# Regenerate schemas from spec
make generate
```

### Local Development Setup

1. **Build and install**: `make install` places binary in Go bin path
2. **Test with Terraform**: Use `terraform init` with local provider in
   `~/.terraform.d/plugins/`
3. **Debug mode**: Run provider with `--debug` flag for attaching debugger

## Project-Specific Conventions

### Code Generation Pattern

**Critical Rule:** Schema changes go in `provider-code-spec.json`, not
`*_gen.go` files.

Workflow:

1. Edit `provider-code-spec.json` to add/modify resource attributes
2. Run `make generate` to regenerate `*_gen.go` files
3. Implement business logic in non-generated files (e.g.,
   `helmfile_release_resource.go`)

Example separation:

```go
// helmfile_release_resource_gen.go - Generated schema
func HelmfileReleaseResourceSchema(ctx context.Context) schema.Schema { ... }

// helmfile_release_resource.go - Your implementation
func (r *HelmfileReleaseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    // TODO: Implement using r.provider (HelmfileExecutor interface)
}
```

### Helmfile Integration

The provider uses helmfile as a library (not CLI):

```go
// Provider configures a HelmfileExecutor instance
type HelmfileProvider struct{}

// Resources receive the executor through Configure()
func (r *HelmfileReleaseResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
    r.provider, _ = req.ProviderData.(*provider_helmfile.HelmfileProvider)
}

// Use executor in CRUD operations
result, err := r.provider.Apply(ctx, options, applyOptions)
```

### Output Capture Pattern

Helmfile operations use a custom logger to capture output:

```go
capture := helmfile.NewOutputCapture()
logger := helmfile.CreateCaptureLogger(capture) // zap.SugaredLogger
// ... run helmfile operation ...
output := capture.String() // Get captured logs
```

This is critical for returning diagnostic information to Terraform users.

### Resource Implementation Pattern

Resources follow this structure:

1. **Metadata**: Set TypeName from provider prefix
2. **Schema**: Return generated schema from `*_gen.go`
3. **Configure**: Receive provider data (HelmfileExecutor)
4. **CRUD**: Implement using helmfile executor interface methods

Currently, CRUD operations have `// TODO` markers - implementation needed.

## Key Integration Points

### Terraform Plugin Framework

Uses Plugin Framework v2 (not older SDKv2):

- Schema via `schema.Schema` structs (generated)
- Configuration via `types.String`, `types.Bool`, etc.
- Diagnostics via `resp.Diagnostics.Add*()` methods

### Helmfile Library

Direct integration with `github.com/helmfile/helmfile/pkg/app`:

- `app.New(options)` creates helmfile app instance
- Operations: `Apply()`, `Diff()`, `Template()`, `Destroy()`
- Uses `app.ConfigProvider` interfaces for configuration

### External Dependencies

- **Helm binary**: Required for helmfile operations (path configurable via
  provider config)
- **Kustomize binary**: Optional (path configurable)
- **kubectl context**: Used if specified in resource config

## Common Operations

### Adding a New Resource Attribute

1. Edit `provider-code-spec.json`:

```json
{
  "name": "new_attribute",
  "string": {
    "description": "Description here",
    "optional_required": "optional"
  }
}
```

2. Regenerate: `make generate`
3. Update model struct in `*_resource.go` if needed
4. Implement handling in CRUD methods

### Implementing a CRUD Operation

Pattern for `Create()`:

```go
func (r *HelmfileReleaseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data HelmfileReleaseModel
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Build helmfile options from data
    options := buildOptionsFromModel(data)
    applyOptions := buildApplyOptions(data)

    // Execute helmfile apply
    result, err := r.provider.Apply(ctx, options, applyOptions)
    if err != nil {
        resp.Diagnostics.AddError("Failed to apply helmfile", result.Output)
        return
    }

    // Populate state from result
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
```

### Testing

Acceptance tests use `TF_ACC=1` environment variable:

```go
func TestAccHelmfileRelease_basic(t *testing.T) {
    resource.Test(t, resource.TestCase{
        ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            {
                Config: testAccHelmfileReleaseConfig_basic,
                Check: resource.ComposeAggregateTestCheckFunc(
                    resource.TestCheckResourceAttr("helmfile_release.test", "name", "test"),
                ),
            },
        },
    })
}
```

## CI/CD Pipeline

- **Test workflow**: Runs on PR and push (`.github/workflows/test.yml`)
  - Build verification
  - Linting with golangci-lint
  - Code generation check
  - Unit tests
- **Release workflow**: Triggered on `v*` tags (`.github/workflows/release.yml`)
  - Uses goreleaser
  - GPG signs artifacts
  - Publishes to Terraform Registry

## Documentation

- Provider docs in `docs/` directory
- Generated by `terraform-plugin-docs` tool
- Uses examples from `examples/` directory
- Resource schemas auto-generated from code

## Troubleshooting

### "Code generated by terraform-plugin-framework-generator DO NOT EDIT"

Don't edit `*_gen.go` files directly. Modify `provider-code-spec.json` and
regenerate.

### Provider not found in local testing

Run `make install` to place provider in `~/go/bin`, then configure Terraform
development overrides:

```hcl
# ~/.terraformrc
provider_installation {
  dev_overrides {
    "kaweezle/helmfile" = "/home/user/go/bin"
  }
  direct {}
}
```

### Helmfile operations failing

Check:

1. Helm binary path is correct in provider config
2. Kubeconfig is accessible
3. Output captured via `OutputCapture` for diagnostic messages

## Additional Resources

- [Terraform Plugin Framework Docs](https://developer.hashicorp.com/terraform/plugin/framework)
- [Helmfile Library Docs](https://github.com/helmfile/helmfile)
- [Provider Code Spec Schema](https://github.com/hashicorp/terraform-plugin-codegen-spec)
