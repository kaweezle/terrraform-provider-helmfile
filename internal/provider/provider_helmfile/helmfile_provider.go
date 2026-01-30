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

// NewGlobalOptionsFromModel creates GlobalOptions from HelmfileModel.
func NewGlobalOptionsFromModel(
	ctx context.Context,
	model *HelmfileModel,
) (helmfile.GlobalOptions, diag.Diagnostics) {
	stringArgs := make([]string, 0)
	if diags := model.DefaultArgs.ElementsAs(ctx, &stringArgs, false); diags.HasError() {
		return helmfile.GlobalOptions{}, diags
	}
	envVars := make(map[string]string)
	if model.EnvVars.IsNull() == false {
		for k, v := range model.EnvVars.Elements() {
			envVars[k] = v.String()
		}
	}

	globalOptions := helmfile.GlobalOptions{
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

// NewHelmfileProvider creates a new HelmfileProvider instance from the given HelmfileModel.
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
