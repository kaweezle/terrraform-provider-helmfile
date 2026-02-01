// Copyright Antoine Martin 2026
// SPDX-License-Identifier: MIT

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/kaweezle/terraform-provider-helmfile/internal/provider/provider_helmfile"
	"github.com/kaweezle/terraform-provider-helmfile/internal/provider/resource_helmfile_release"
)

var _ provider.Provider = (*helmfileProvider)(nil)

// New creates a new provider instance.
//
// This function is the factory method for creating the helmfile provider.
// It should return a function that creates a properly configured provider
// instance when called by the Terraform framework.
//
// Parameters:
//   - version: Version string of the provider (set at build time)
//
// Returns:
//   - func() provider.Provider: Factory function for provider instantiation
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &helmfileProvider{
			version: version,
		}
	}
}

type helmfileProvider struct {
	version string
}

// Schema defines the provider-level schema.
//
// This method returns the schema for provider configuration. The schema
// is automatically generated from the provider-code-spec.json file.
//
// Parameters:
//   - ctx: Context for the operation
//   - req: Schema request from the framework
//   - resp: Schema response to populate with the provider schema
func (p *helmfileProvider) Schema(
	ctx context.Context,
	_ provider.SchemaRequest,
	resp *provider.SchemaResponse,
) {
	tflog.Debug(ctx, "in Schema request")
	resp.Schema = provider_helmfile.HelmfileProviderSchema(ctx)
	// FIXME: It should be generated from the schema definition
	resp.Schema.Description = "Helmfile provider configures the Helmfile CLI tool for managing Helm charts deployments."
}

// Configure configures the provider with given configuration values.
//
// This method is called by Terraform during provider initialization to apply
// the user's configuration. It validates the configuration and creates the
// underlying helmfile provider instance that will be used by resources.
//
// Parameters:
//   - ctx: Context for the operation
//   - req: Configuration request containing user-provided values
//   - resp: Configuration response to populate with the configured provider instance
func (p *helmfileProvider) Configure(
	ctx context.Context,
	req provider.ConfigureRequest,
	resp *provider.ConfigureResponse,
) {
	tflog.Debug(ctx, "In provider configuration")
	var config provider_helmfile.HelmfileModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.HelmBinaryPath.IsNull() || config.HelmBinaryPath.IsUnknown() {
		config.HelmBinaryPath = types.StringValue("helm")
	}
	performInit := config.PerformInit.ValueBool()
	var providerInstance *provider_helmfile.HelmfileProvider
	providerInstance, diags = provider_helmfile.NewHelmfileProvider(ctx, &config)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	if performInit {
		output, logs, err := providerInstance.Executor.Init(ctx, nil)
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Error during init: %s", err.Error()), logs)
			return
		}
		tflog.Info(ctx, "Helmfile init output:\n"+output)
		tflog.Info(ctx, "Helmfile init logs:\n"+logs)
	}
	resp.DataSourceData = providerInstance
	resp.ResourceData = providerInstance
}

// Metadata sets provider-level metadata.
//
// This method returns metadata about the provider that Terraform uses for
// provider identification and configuration.
//
// Parameters:
//   - ctx: Context for the operation
//   - req: Metadata request from the framework
//   - resp: Metadata response to populate
func (p *helmfileProvider) Metadata(
	_ context.Context,
	_ provider.MetadataRequest,
	resp *provider.MetadataResponse,
) {
	resp.TypeName = "helmfile"
	resp.Version = p.version
}

// DataSources returns the provider's data source types.
//
// This method returns a list of data source factories for all data sources
// provided by this provider. Currently returns an empty list as no data
// sources are implemented.
//
// Parameters:
//   - ctx: Context for the operation
//
// Returns:
//   - []func() datasource.DataSource: Empty slice (no data sources)
func (p *helmfileProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

// Resources returns the provider's resource types.
//
// This method returns a list of resource factories for all resources
// provided by this provider. Currently provides the helmfile_release resource.
//
// Parameters:
//   - ctx: Context for the operation
//
// Returns:
//   - []func() resource.Resource: Slice of resource factory functions
func (p *helmfileProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resource_helmfile_release.NewHelmfileReleaseResource,
	}
}
