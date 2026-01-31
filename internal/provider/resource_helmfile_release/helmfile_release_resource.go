// Copyright Antoine Martin 2026
// SPDX-License-Identifier: MIT
// cSpell: words norbac basetypes

package resource_helmfile_release

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/helmfile/helmfile/pkg/app"
	"github.com/helmfile/helmfile/pkg/config"
	"github.com/kaweezle/terraform-provider-helmfile/internal/provider/provider_helmfile"
	"github.com/kaweezle/terraform-provider-helmfile/pkg/helmfile"
)

var (
	_ resource.ResourceWithConfigure  = (*HelmfileReleaseResource)(nil)
	_ resource.Resource               = (*HelmfileReleaseResource)(nil)
	_ resource.ResourceWithModifyPlan = (*HelmfileReleaseResource)(nil)
)

type HelmfileReleaseResource struct {
	provider *provider_helmfile.HelmfileProvider
}

func NewHelmfileReleaseResource() resource.Resource {
	return &HelmfileReleaseResource{}
}

func (r *HelmfileReleaseResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_release"
}

func NewOptionsFromModel(
	ctx context.Context,
	model *HelmfileReleaseModel,
) (*helmfile.Options, diag.Diagnostics) {
	stringArgs := make([]string, 0)
	if diags := model.Args.ElementsAs(ctx, &stringArgs, false); diags.HasError() {
		return nil, diags
	}
	envVars := make(map[string]string)
	if !model.EnvVars.IsNull() {
		for k, v := range model.EnvVars.Elements() {
			envVars[k] = v.String()
		}
	}

	options := &helmfile.Options{
		BaseResourceOptions: helmfile.BaseResourceOptions{},
		BaseGlobalOptions:   helmfile.BaseGlobalOptions{},
		CommonOptions:       helmfile.CommonOptions{},
	}

	// Set base global options
	options.BaseResourceOptions.
		WithArgs(strings.Join(stringArgs, " ")).
		WithChart(model.Chart.ValueString()).
		WithFileOrDir(model.FileOrPath.ValueString()).
		WithKubeContext(model.KubeContext.ValueString()).
		WithNamespace(model.Namespace.ValueString()).
		WithSelectors(func() []string {
			selectors := make([]string, 0)
			if !model.Selectors.IsNull() {
				for _, v := range model.Selectors.Elements() {
					selectors = append(selectors, v.String())
				}
			}
			return selectors
		}()).
		WithStateValuesSet(func() map[string]any {
			valuesSet := make(map[string]any)
			if !model.StateValuesSet.IsNull() {
				for k, v := range model.StateValuesSet.Elements() {
					valuesSet[k] = v.String()
				}
			}
			return valuesSet
		}()).
		WithStateValuesFiles(func() []string {
			valuesFiles := make([]string, 0)
			if !model.StateValuesFiles.IsNull() {
				for _, v := range model.StateValuesFiles.Elements() {
					valuesFiles = append(valuesFiles, v.String())
				}
			}
			return valuesFiles
		}())

	if !model.Overrides.IsNull() && !model.Overrides.IsUnknown() {
		options.BaseGlobalOptions.WithDisableForceUpdate(model.Overrides.DisableForceUpdate.ValueBool()).
			WithEnforcePluginVerification(model.Overrides.EnforcePluginVerification.ValueBool()).
			WithHelmBinary(model.Overrides.HelmBinaryPath.ValueString()).
			WithHelmOCIPlainHTTP(model.Overrides.HelmOciPlainHttp.ValueBool()).
			WithKustomizeBinary(model.Overrides.KustomizeBinaryPath.ValueString()).
			WithSkipDeps(model.Overrides.SkipDeps.ValueBool()).
			WithSkipRefresh(model.Overrides.SkipRefresh.ValueBool()).
			WithStripArgsValuesOnExitError(model.Overrides.StripArgsValuesOnExitError.ValueBool())
	}

	// Set common options
	options.CommonOptions.
		WithLogLevel(model.LogLevel.ValueString()).
		WithKubeconfig(model.Kubeconfig.ValueString()).
		WithEnvironment(model.Environment.ValueString()).
		WithEnvVars(envVars)

	return options, diag.Diagnostics{}
}

//nolint:gocyclo // Function is complex due to many fields to map
func NewApplyOptionsFromModel(
	ctx context.Context,
	model *HelmfileReleaseModel,
) (*config.ApplyOptions, diag.Diagnostics) {
	applyOptions := &config.ApplyOptions{
		Cascade:                  model.Cascade.ValueString(),
		Concurrency:              int(model.Concurrency.ValueInt64()),
		Context:                  int(model.Context.ValueInt64()),
		DetailedExitcode:         model.DetailedExitcode.ValueBool(),
		DiffArgs:                 model.DiffArgs.ValueString(),
		EnforceNeedsAreInstalled: model.EnforceNeedsAreInstalled.ValueBool(),
		HideNotes:                model.HideNotes.ValueBool(),
		IncludeNeeds:             model.IncludeNeeds.ValueBool(),
		IncludeTests:             model.IncludeTests.ValueBool(),
		IncludeTransitiveNeeds:   model.IncludeTransitiveNeeds.ValueBool(),
		NoHooks:                  model.NoHooks.ValueBool(),
		Output:                   model.Output.ValueString(),
		PostRenderer:             model.PostRenderer.ValueString(),
		ResetValues:              model.ResetValues.ValueBool(),
		ReuseValues:              model.ReuseValues.ValueBool(),
		ShowSecrets:              model.ShowSecrets.ValueBool(),
		SkipCleanup:              model.SkipCleanup.ValueBool(),
		SkipCRDs:                 model.SkipCrds.ValueBool(),
		SkipDiffOnInstall:        model.SkipDiffOnInstall.ValueBool(),
		SkipNeeds:                model.SkipNeeds.ValueBool(),
		SkipSchemaValidation:     model.SkipSchemaValidation.ValueBool(),
		StripTrailingCR:          model.StripTrailingCr.ValueBool(),
		SuppressDiff:             model.SuppressDiff.ValueBool(),
		SuppressSecrets:          model.SuppressSecrets.ValueBool(),
		SyncArgs:                 model.SyncArgs.ValueString(),
		SyncReleaseLabels:        model.SyncReleaseLabels.ValueBool(),
		TakeOwnership:            model.TakeOwnership.ValueBool(),
		Validate:                 model.Validate.ValueBool(),
		Wait:                     model.Wait.ValueBool(),
		WaitForJobs:              model.WaitForJobs.ValueBool(),
		WaitRetries:              int(model.WaitRetries.ValueInt64()),
	}

	if !model.PostRendererArgs.IsNull() && !model.PostRendererArgs.IsUnknown() {
		stringArgs := make([]string, 0)
		if diags := model.PostRendererArgs.ElementsAs(ctx, &stringArgs, false); diags.HasError() {
			return nil, diags
		}
		applyOptions.PostRendererArgs = stringArgs
	}

	if !model.Set.IsNull() && !model.Set.IsUnknown() {
		stringArgs := make([]string, 0)
		if diags := model.Set.ElementsAs(ctx, &stringArgs, false); diags.HasError() {
			return nil, diags
		}
		applyOptions.Set = stringArgs
	}
	if !model.Suppress.IsNull() && !model.Suppress.IsUnknown() {
		stringArgs := make([]string, 0)
		if diags := model.Suppress.ElementsAs(ctx, &stringArgs, false); diags.HasError() {
			return nil, diags
		}
		applyOptions.Suppress = stringArgs
	}

	if !model.SuppressOutputLineRegex.IsNull() && !model.SuppressOutputLineRegex.IsUnknown() {
		stringArgs := make([]string, 0)
		if diags := model.SuppressOutputLineRegex.ElementsAs(ctx, &stringArgs, false); diags.HasError() {
			return nil, diags
		}
		applyOptions.SuppressOutputLineRegex = stringArgs
	}
	if !model.Values.IsNull() && !model.Values.IsUnknown() {
		stringArgs := make([]string, 0)
		if diags := model.Values.ElementsAs(ctx, &stringArgs, false); diags.HasError() {
			return nil, diags
		}
		applyOptions.Values = stringArgs
	}

	// Set apply options based on model fields if needed
	return applyOptions, diag.Diagnostics{}
}

func (r *HelmfileReleaseResource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = HelmfileReleaseResourceSchema(ctx)
	resp.Schema.Description = "Manages a Helm release defined in a Helmfile."
}

func (r *HelmfileReleaseResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	tflog.Debug(ctx, "Configuring Helmfile release resource")
	if req.ProviderData == nil {
		return
	}
	provider, ok := req.ProviderData.(*provider_helmfile.HelmfileProvider)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf(
				"Expected *provider_helmfile.HelmfileProvider, got: %T. Please report this issue to the provider developers.",
				req.ProviderData,
			),
		)

		return
	}
	r.provider = provider
}

// This provider uses output of `helmfile build` to calculate the hash key of the helmfile-diff cache, which is used
// to make originally non-deterministic `helmfile-diff` result to be deterministic.
//
// In `removeNondeterministicBuildLogLines`, we remove some part of `helm build --embed-values` that is
// non-deterministic due to that the temporary helmfile.yaml generated by the provider has a random name.
func removeNondeterministicBuildLogLines(s string) (string, error) {
	buf := &bytes.Buffer{}
	w := bufio.NewWriter(buf)

	b := bufio.NewScanner(strings.NewReader(s))
	for b.Scan() {
		l := b.Text()
		if !strings.HasPrefix(l, "#") && !strings.HasPrefix(l, "filepath: ") {
			if _, err := w.WriteString(l); err != nil {
				return "", fmt.Errorf("writing line: %w", err)
			}
			if err := w.WriteByte('\n'); err != nil {
				return "", fmt.Errorf("writing newline: %w", err)
			}
		}
	}
	if err := w.Flush(); err != nil {
		return "", fmt.Errorf("filtering helmfile build output: %w", err)
	}

	return buf.String(), nil
}

//nolint:gocritic // Interface implementation
func (r *HelmfileReleaseResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	tflog.Debug(ctx, "################### Creating Helmfile release resource")
	var data HelmfileReleaseModel
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// TODO: Implement the logic to create the Helmfile release using r.provider
	// Steps (To be completed):
	// 1. DONE: Make a build --embed-values call to helmfile with the relevant helmfile and environment and
	//    create a sha256 hash of the output
	// 2. DONE: Make a list in order to list releases
	// 3. Make an apply call to helmfile with the relevant helmfile and environment to ensure the release is
	//    created/updated
	// 4. DONE: Store the hash in the state to detect future changes

	// executor := r.provider.Executor
	// options, diags := NewOptionsFromModel(ctx, &data)
	// if diags.HasError() {
	// 	resp.Diagnostics.Append(diags...)
	// 	return
	// }

	// applyOptions, diags := NewApplyOptionsFromModel(ctx, &data)
	// if diags.HasError() {
	// 	resp.Diagnostics.Append(diags...)
	// 	return
	// }

	// Set the state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func updateReleaseState(
	ctx context.Context,
	executor helmfile.HelmfileExecutor,
	options *helmfile.Options,
	data *HelmfileReleaseModel,
) diag.Diagnostics {
	var diags diag.Diagnostics
	var digest string
	tflog.Debug(ctx, "Calculating Helmfile build digest for release resource")
	digest, diags = releaseDigest(ctx, executor, options)
	if diags.HasError() {
		diags.Append(diags...)
		return diags
	}

	tflog.Debug(ctx, "Helmfile build digest calculated", map[string]any{
		"sha256_checksum": digest,
	})
	data.Sha256Checksum = types.StringValue(digest)

	tflog.Debug(ctx, "Getting the list of Helmfile releases")
	output, logs, err := executor.List(ctx, options, true)
	if err != nil {
		diags.AddError(
			"Error listing Helmfile releases",
			fmt.Sprintf(
				"An error occurred while listing Helmfile releases: %s\nLogs:\n%s",
				err.Error(),
				logs,
			),
		)
		return diags
	}

	// Create the list of objects from the JSON in output
	releaseList, diags := jsonToReleasesListValue(ctx, output)
	if diags.HasError() {
		diags.Append(diags...)
		return diags
	}
	data.ReleasesList = *releaseList
	return diags
}

// releaseDigest runs `helmfile build` and returns the sha256 checksum of its output.
func releaseDigest(
	ctx context.Context,
	executor helmfile.HelmfileExecutor,
	options *helmfile.Options,
) (string, diag.Diagnostics) {
	output, logs, err := executor.Build(ctx, options, true)
	var diags diag.Diagnostics
	if err != nil {
		diags.AddError(
			"Error building Helmfile releases",
			fmt.Sprintf(
				"An error occurred while building Helmfile releases: %s\nLogs:\n%s",
				err.Error(),
				logs,
			),
		)
		return "", diags
	}
	// Keep for debugging purposes. Don't show in normal runs as it may contain sensitive info.
	// tflog.Debug(ctx, "Build output", map[string]interface{}{
	// 	"output": output,
	// })

	var cleanOutput string
	cleanOutput, err = removeNondeterministicBuildLogLines(output)
	if err != nil {
		diags.AddError(
			"Error filtering Helmfile build output",
			fmt.Sprintf("An error occurred while filtering Helmfile build output: %s", err.Error()),
		)
		return "", diags
	}
	// Keep for debugging purposes. Don't show in normal runs as it may contain sensitive info.
	// tflog.Debug(ctx, "Clean output", map[string]interface{}{
	// 	"output": cleanOutput,
	// })

	hash := sha256.New()
	hash.Write([]byte(cleanOutput))
	return fmt.Sprintf("%x", hash.Sum(nil)), diags
}

// jsonToReleasesListValue converts the JSON output of `helmfile list -o json` into a ListValue of ReleasesListValue.
func jsonToReleasesListValue(
	ctx context.Context,
	output string,
) (*basetypes.ListValue, diag.Diagnostics) {
	var releases []app.HelmRelease
	var diags diag.Diagnostics
	if err := json.Unmarshal([]byte(output), &releases); err != nil {
		diags.AddError(
			"Error unmarshaling Helmfile releases list",
			fmt.Sprintf(
				"An error occurred while unmarshaling Helmfile releases list: %s",
				err.Error(),
			),
		)
		return nil, diags
	}

	if len(releases) == 0 {
		result := types.ListValueMust(
			NewReleasesListValueNull().Type(ctx),
			[]attr.Value{},
		)
		return &result, diags
	}

	// Now create the list of releases in the state
	values := make([]attr.Value, 0, len(releases))
	for _, release := range releases {
		// Populate tags map. Tags are in the frm "chart:prometheus,name:prom-norbac-ubuntu,namespace:prometheus"
		labelTags := strings.Split(release.Labels, ",")
		tags := make(map[string]attr.Value, len(labelTags))
		for _, tag := range labelTags {
			parts := strings.SplitN(tag, ":", 2)
			if len(parts) == 2 {
				tags[parts[0]] = types.StringValue(parts[1])
			} else {
				diags.AddError(
					"Error parsing release tag",
					fmt.Sprintf("An error occurred while parsing release tag: %s", tag),
				)
				return nil, diags
			}
		}

		values = append(values, ReleasesListValue{
			Name:      types.StringValue(release.Name),
			Namespace: types.StringValue(release.Namespace),
			Chart:     types.StringValue(release.Chart),
			Version:   types.StringValue(release.Version),
			Enabled:   types.BoolValue(release.Enabled),
			Installed: types.BoolValue(release.Installed),
			Labels:    types.MapValueMust(basetypes.StringType{}, tags),
			state:     attr.ValueStateKnown,
		})
	}
	result := types.ListValueMust(
		values[0].Type(ctx),
		values,
	)
	return &result, diags
}

//nolint:gocritic // Interface implementation
func (r *HelmfileReleaseResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	tflog.Debug(ctx, "################### Reading Helmfile release resource")
	var data HelmfileReleaseModel
	diags := req.State.Get(ctx, &data)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	// executor := r.provider.Executor
	// options, diags := NewOptionsFromModel(ctx, &data)
	// if diags.HasError() {
	// 	resp.Diagnostics.Append(diags...)
	// 	return
	// }

	tflog.Debug(ctx, "Helmfile release resource state updated during read", map[string]any{
		"sha256_checksum": data.Sha256Checksum.ValueString(),
	})

	// Set the state
	diags = resp.State.Set(ctx, &data)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
}

//nolint:gocritic // Interface implementation
func (r *HelmfileReleaseResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	tflog.Debug(ctx, "################### Updating Helmfile release resource")
	var data HelmfileReleaseModel
	diags := req.Plan.Get(ctx, &data)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// executor := r.provider.Executor
	// options, diags := NewOptionsFromModel(ctx, &data)
	// if diags.HasError() {
	// 	resp.Diagnostics.Append(diags...)
	// 	return
	// }

	// TODO: Implement the logic to create the Helmfile release using r.provider

	// Set the state
	diags = resp.State.Set(ctx, &data)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
}

//nolint:gocritic // Interface implementation
func (r *HelmfileReleaseResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	tflog.Debug(ctx, "################### Deleting Helmfile release resource")
	// Implementation of Delete operation
	var data HelmfileReleaseModel
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

//nolint:gocritic // Interface implementation
func (r *HelmfileReleaseResource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	tflog.Debug(ctx, "################### Modifying plan for Helmfile release resource")
	if req.Plan.Raw.IsNull() {
		// Resource is being destroyed, no need to modify the plan
		return
	}
	// Implementation of ModifyPlan operation
	var data HelmfileReleaseModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	executor := r.provider.Executor
	options, diags := NewOptionsFromModel(ctx, &data)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	diags = updateReleaseState(ctx, executor, options, &data)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.Plan.Set(ctx, &data)...)
}
