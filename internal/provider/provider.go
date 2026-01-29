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

func (p *helmfileProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	tflog.Debug(ctx, "[HELMFILE] in Schema request")
	resp.Schema = provider_helmfile.HelmfileProviderSchema(ctx)
	// FIXME: It should be generated from the schema definition
	resp.Schema.Description = "Helmfile provider configures the Helmfile CLI tool for managing Helm charts deployments."
}

func (p *helmfileProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
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
	providerInstance, diags = provider_helmfile.NewHelmfileProvider(config)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	if performInit {
		output, err := providerInstance.Executor.Init(ctx, nil)
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Error during init: %s", err.Error()), output)
			return
		}
		tflog.Info(ctx, "Helmfile init output:\n"+output)
	}
	resp.DataSourceData = providerInstance
	resp.ResourceData = providerInstance
}

func (p *helmfileProvider) Metadata(_ context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "helmfile"
	resp.Version = p.version
}

func (p *helmfileProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *helmfileProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resource_helmfile_release.NewHelmfileReleaseResource,
	}
}
