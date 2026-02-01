// Copyright Antoine Martin 2026
// SPDX-License-Identifier: MIT

// cSpell: words helmexec cliv3 cliv4

package provider_helmfile

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/kaweezle/terraform-provider-helmfile/pkg/helmfile"
)

// HelmfileProvider is the Helmfile provider implementation.
type HelmfileProvider struct {
	Executor *helmfile.HelmfileLibraryExecutor
}

// NewGlobalOptionsFromModel creates GlobalOptions from a Terraform model.
//
// This function converts the provider configuration model from Terraform into
// the internal GlobalOptions structure used by the helmfile executor. It extracts
// and validates configuration values from the model and populates both base global
// options and common options.
//
// Parameters:
//   - ctx: Context for the operation
//   - model: Provider model from Terraform configuration
//
// Returns:
//   - *helmfile.GlobalOptions: Converted global options
//   - diag.Diagnostics: Diagnostics if conversion encounters issues
func NewGlobalOptionsFromModel(
	ctx context.Context,
	model *HelmfileModel,
) (*helmfile.GlobalOptions, diag.Diagnostics) {
	stringArgs := make([]string, 0)
	if diags := model.DefaultArgs.ElementsAs(ctx, &stringArgs, false); diags.HasError() {
		return nil, diags
	}
	envVars := make(map[string]string)
	if !model.EnvVars.IsNull() {
		for k, v := range model.EnvVars.Elements() {
			envVars[k] = v.String()
		}
	}

	globalOptions := &helmfile.GlobalOptions{
		BaseGlobalOptions: helmfile.BaseGlobalOptions{},
		CommonOptions:     helmfile.CommonOptions{},
	}

	// Set base global options
	globalOptions.BaseGlobalOptions.
		WithDefaultArgs(strings.Join(stringArgs, " ")).
		WithDisableForceUpdate(model.DisableForceUpdate.ValueBool()).
		WithEnforcePluginVerification(model.EnforcePluginVerification.ValueBool()).
		WithHelmBinary(model.HelmBinaryPath.ValueString()).
		WithHelmOCIPlainHTTP(model.HelmOciPlainHttp.ValueBool()).
		WithKustomizeBinary(model.KustomizeBinaryPath.ValueString()).
		WithSkipDeps(model.SkipDeps.ValueBool()).
		WithSkipRefresh(model.SkipRefresh.ValueBool()).
		WithStripArgsValuesOnExitError(model.StripArgsValuesOnExitError.ValueBool())

	// Set common options
	globalOptions.CommonOptions.
		WithLogLevel(model.LogLevel.ValueString()).
		WithKubeconfig(model.Kubeconfig.ValueString()).
		WithEnvironment(model.Environment.ValueString()).
		WithEnvVars(envVars)

	return globalOptions, diag.Diagnostics{}
}

// NewPluginsFromModel converts plugin configuration to HelmPlugin structs.
//
// This function takes the plugin configuration from the Terraform provider
// model and converts it to the internal HelmPlugin representation. It validates
// and extracts plugin name, repository, and version information.
//
// Parameters:
//   - ctx: Context for the operation
//   - model: Provider model containing plugin configuration
//
// Returns:
//   - []helmfile.HelmPlugin: Converted plugin list
//   - diag.Diagnostics: Diagnostics if conversion encounters issues
func NewPluginsFromModel(
	ctx context.Context,
	model *HelmfileModel,
) ([]helmfile.HelmPlugin, diag.Diagnostics) {
	modelPlugins := make([]AdditionalPluginsValue, 0)
	if diags := model.AdditionalPlugins.ElementsAs(ctx, &modelPlugins, false); diags.HasError() {
		return nil, diags
	}
	result := make([]helmfile.HelmPlugin, 0, len(modelPlugins))
	for _, p := range modelPlugins {
		result = append(result, helmfile.HelmPlugin{
			Name:    p.Name.ValueString(),
			Repo:    p.Repo.ValueString(),
			Version: p.Version.ValueString(),
		})
	}
	return result, diag.Diagnostics{}
}

// NewHelmfileProvider creates a new HelmfileProvider instance from the given model.
//
// This function creates and initializes the provider instance with global options
// and additional plugins extracted from the Terraform configuration model. It creates
// the underlying helmfile executor that will be used by all resources.
//
// Parameters:
//   - ctx: Context for the operation
//   - model: Provider configuration model from Terraform
//
// Returns:
//   - *HelmfileProvider: Configured provider instance ready for use
//   - diag.Diagnostics: Diagnostics if provider creation encounters issues
func NewHelmfileProvider(
	ctx context.Context,
	model *HelmfileModel,
) (*HelmfileProvider, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	options, err := NewGlobalOptionsFromModel(ctx, model)
	diags.Append(err...)
	if err.HasError() {
		return nil, diags
	}
	additionalPlugins, err := NewPluginsFromModel(ctx, model)
	diags.Append(err...)
	if err.HasError() {
		return nil, diags
	}

	return &HelmfileProvider{
		Executor: helmfile.NewHelmfileLibraryExecutor(options, additionalPlugins),
	}, diags
}
