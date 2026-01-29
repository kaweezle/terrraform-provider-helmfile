// Copyright (c) Antoine Martin
// SPDX-License-Identifier: MIT

// cSpell: words helmexec cliv3 cliv4

package provider_helmfile

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/helmfile/helmfile/pkg/config"
	"github.com/kaweezle/terraform-provider-helmfile/pkg/helmfile"
)

// HelmfileProvider is the Helmfile provider implementation
type HelmfileProvider struct {
	Executor *helmfile.HelmfileLibraryExecutor
}

// NewGlobalOptionsFromModel creates GlobalOptions from HelmfileModel
func NewGlobalOptionsFromModel(model HelmfileModel) (*config.GlobalOptions, diag.Diagnostics) {
	stringArgs := make([]string, 0)
	if diag := model.DefaultArgs.ElementsAs(context.Background(), &stringArgs, false); diag.HasError() {
		return nil, diag
	}

	globalOptions := &config.GlobalOptions{
		Args:                       strings.Join(stringArgs, " "),
		DisableForceUpdate:         model.DisableForceUpdate.ValueBool(),
		EnforcePluginVerification:  model.EnforcePluginVerification.ValueBool(),
		HelmBinary:                 model.HelmBinaryPath.ValueString(),
		HelmOCIPlainHTTP:           model.HelmOciPlainHttp.ValueBool(),
		KustomizeBinary:            model.KustomizeBinaryPath.ValueString(),
		SkipDeps:                   model.SkipDeps.ValueBool(),
		SkipRefresh:                model.SkipRefresh.ValueBool(),
		StripArgsValuesOnExitError: model.StripArgsValuesOnExitError.ValueBool(),
		EnableLiveOutput:           false,
		Color:                      false,
		NoColor:                    true,
		Debug:                      model.Debug.ValueBool(),
		Quiet:                      true,
	}
	return globalOptions, diag.Diagnostics{}
}

func NewPluginsFromModel(model HelmfileModel) ([]helmfile.HelmPlugin, diag.Diagnostics) {
	modelPlugins := make([]AdditionalPluginsValue, 0)
	if diag := model.AdditionalPlugins.ElementsAs(context.Background(), &modelPlugins, false); diag.HasError() {
		return nil, diag
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

// NewHelmfileProvider creates a new HelmfileProvider instance from the given HelmfileModel
func NewHelmfileProvider(model HelmfileModel) (*HelmfileProvider, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	options, err := NewGlobalOptionsFromModel(model)
	diags.Append(err...)
	if err.HasError() {
		return nil, diags
	}
	additionalPlugins, err := NewPluginsFromModel(model)
	diags.Append(err...)
	if err.HasError() {
		return nil, diags
	}

	return &HelmfileProvider{
		Executor: helmfile.NewHelmfileLibraryExecutor(options, additionalPlugins),
	}, diags
}
