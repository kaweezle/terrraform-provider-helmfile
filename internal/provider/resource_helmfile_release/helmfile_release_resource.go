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

// NewHelmfileReleaseResource creates a new helmfile_release resource instance.
//
// This is the factory function used by the Terraform framework to create
// new instances of the helmfile_release resource.
//
// Returns:
//   - resource.Resource: A new resource instance
func NewHelmfileReleaseResource() resource.Resource {
	return &HelmfileReleaseResource{}
}

// Metadata sets resource-level metadata.
//
// This method returns metadata about the resource that Terraform uses for
// resource identification and configuration. It sets the resource type name
// by appending "_release" to the provider type name.
//
// Parameters:
//   - _: Context (not used in this implementation)
//   - req: Metadata request from the framework
//   - resp: Metadata response to populate
func (r *HelmfileReleaseResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_release"
}

// NewOptionsFromModel creates helmfile Options from a Terraform model.
//
// This function converts the resource configuration model from Terraform into
// the internal Options structure used by the helmfile executor. It extracts
// and validates configuration values including args, environment variables,
// selectors, state values, and overrides.
//
// Parameters:
//   - ctx: Context for the operation
//   - model: Resource model from Terraform configuration
//
// Returns:
//   - *helmfile.Options: Converted options for helmfile operations
//   - diag.Diagnostics: Diagnostics if conversion encounters issues
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

// NewApplyOptionsFromModel creates helmfile ApplyOptions from a Terraform model.
//
// This function converts the resource configuration model into the helmfile
// ApplyOptions structure. It maps all apply-specific configuration values including
// cascade behavior, concurrency, wait settings, hooks, CRDs, cleanup, and more.
//
// Parameters:
//   - ctx: Context for the operation
//   - model: Resource model from Terraform configuration
//
// Returns:
//   - *config.ApplyOptions: Converted options for helmfile apply operation
//   - diag.Diagnostics: Diagnostics if conversion encounters issues
//
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
		diags := model.SuppressOutputLineRegex.ElementsAs(ctx, &stringArgs, false)
		if diags.HasError() {
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

// Schema defines the resource schema.
//
// This method returns the schema for the resource configuration. The schema
// is automatically generated from the provider-code-spec.json file.
//
// Parameters:
//   - ctx: Context for the operation
//   - _: Schema request (not used in this implementation)
//   - resp: Schema response to populate with the resource schema
func (r *HelmfileReleaseResource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = HelmfileReleaseResourceSchema(ctx)
	resp.Schema.Description = "Manages a Helm release defined in a Helmfile."
}

// Configure configures the resource with provider data.
//
// This method is called by the framework during resource initialization to
// provide the resource with access to the configured provider instance. It
// validates the provider type and stores it for use in CRUD operations.
//
// Parameters:
//   - ctx: Context for the operation
//   - req: Configuration request containing provider data
//   - resp: Configuration response to populate with diagnostics
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

// removeNondeterministicBuildLogLines filters build output for deterministic hashing.
//
// This function removes non-deterministic lines from helmfile build output to enable
// consistent hash calculation. It filters out comment lines and filepath lines that
// contain randomly generated temporary file names. This makes the helmfile-diff result
// deterministic by using the output of `helmfile build --embed-values` as the hash key.
//
// Parameters:
//   - s: Raw helmfile build output string
//
// Returns:
//   - string: Filtered output with non-deterministic lines removed
//   - error: If filtering fails
//
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

// Create implements the resource Create operation.
//
// This method creates a new helmfile release by applying the helmfile configuration.
// It extracts the configuration from the plan, executes helmfile apply, and saves
// the resulting state.
//
// Parameters:
//   - ctx: Context for the operation
//   - req: Create request containing the planned configuration
//   - resp: Create response to populate with the resulting state
//
//nolint:gocritic // Interface implementation
func (r *HelmfileReleaseResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	tflog.Debug(ctx, "=HLMFL=> Creating Helmfile release resource")
	var data HelmfileReleaseModel
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	executor := r.provider.Executor
	resp.Diagnostics.Append(applyHelmfileRelease(ctx, executor, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// applyHelmfileRelease executes the helmfile apply operation.
//
// This function converts the resource model to helmfile options and apply options,
// then executes the apply operation using the helmfile executor. It handles errors
// and logs the operation results.
//
// Parameters:
//   - ctx: Context for the operation
//   - executor: Helmfile executor to use for the apply operation
//   - data: Resource model containing configuration
//
// Returns:
//   - diag.Diagnostics: Diagnostics if the operation encounters issues
func applyHelmfileRelease(
	ctx context.Context,
	executor helmfile.HelmfileExecutor,
	data *HelmfileReleaseModel,
) diag.Diagnostics {
	options, diags := NewOptionsFromModel(ctx, data)
	if diags.HasError() {
		return diags
	}

	applyOptions, diags := NewApplyOptionsFromModel(ctx, data)
	if diags.HasError() {
		return diags
	}

	output, logs, err := executor.Apply(ctx, options, applyOptions)
	if err != nil {
		diags.AddError(
			"Error applying Helmfile releases",
			fmt.Sprintf(
				"An error occurred while applying Helmfile releases: %s\nLogs:\n%s\nOutput:\n%s",
				err.Error(),
				logs,
				output,
			),
		)
		return diags
	}
	tflog.Debug(ctx, "Helmfile releases applied successfully", map[string]any{
		"output": output,
		"logs":   logs,
	})
	return diag.Diagnostics{}
}

// Read implements the resource Read operation.
//
// This method reads the current state of the helmfile release. Currently, it
// simply preserves the existing state without querying the cluster, as helmfile
// does not provide a native read operation.
//
// Parameters:
//   - ctx: Context for the operation
//   - req: Read request containing the current state
//   - resp: Read response to populate with refreshed state
//
//nolint:gocritic // Interface implementation
func (r *HelmfileReleaseResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	tflog.Debug(ctx, "=HLMFL=> Reading Helmfile release resource")
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

	// Set the state
	diags = resp.State.Set(ctx, &data)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
}

// Update implements the resource Update operation.
//
// This method updates an existing helmfile release by applying the updated
// configuration. It follows the same logic as Create, using helmfile apply
// to bring the release to the desired state.
//
// Parameters:
//   - ctx: Context for the operation
//   - req: Update request containing the new planned configuration
//   - resp: Update response to populate with the resulting state
//
//nolint:gocritic // Interface implementation
func (r *HelmfileReleaseResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	tflog.Debug(ctx, "=HLMFL=> Updating Helmfile release resource")
	var data HelmfileReleaseModel
	diags := req.Plan.Get(ctx, &data)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	executor := r.provider.Executor
	resp.Diagnostics.Append(applyHelmfileRelease(ctx, executor, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state
	diags = resp.State.Set(ctx, &data)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
}

// Delete implements the resource Delete operation.
//
// This method deletes the helmfile release by executing helmfile destroy.
// It removes all Helm releases defined in the helmfile configuration from
// the Kubernetes cluster.
//
// Parameters:
//   - ctx: Context for the operation
//   - req: Delete request containing the current state
//   - resp: Delete response to populate with diagnostics
//
//nolint:gocritic // Interface implementation
func (r *HelmfileReleaseResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	tflog.Debug(ctx, "=HLMFL=> Deleting Helmfile release resource")
	// Implementation of Delete operation
	var data HelmfileReleaseModel
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	executor := r.provider.Executor
	resp.Diagnostics.Append(destroyHelmfileRelease(ctx, executor, &data)...)
}

// destroyHelmfileRelease executes the helmfile destroy operation.
//
// This function converts the resource model to helmfile options and destroy options,
// then executes the destroy operation using the helmfile executor. It handles errors
// and logs the operation results.
//
// Parameters:
//   - ctx: Context for the operation
//   - executor: Helmfile executor to use for the destroy operation
//   - data: Resource model containing configuration
//
// Returns:
//   - diag.Diagnostics: Diagnostics if the operation encounters issues
func destroyHelmfileRelease(
	ctx context.Context,
	executor helmfile.HelmfileExecutor,
	data *HelmfileReleaseModel,
) diag.Diagnostics {
	options, diags := NewOptionsFromModel(ctx, data)
	if diags.HasError() {
		return diags
	}

	destroyOptions := NewDestroyOptionsFromModel(data)

	output, logs, err := executor.Destroy(ctx, options, destroyOptions)
	if err != nil {
		diags.AddError(
			"Error destroying Helmfile releases",
			fmt.Sprintf(
				"An error occurred while destroying Helmfile releases: %s\nLogs:\n%s\nOutput:\n%s",
				err.Error(),
				logs,
				output,
			),
		)
		return diags
	}
	tflog.Debug(ctx, "Helmfile releases destroyed successfully", map[string]any{
		"output": output,
		"logs":   logs,
	})
	return diag.Diagnostics{}
}

// NewDestroyOptionsFromModel creates helmfile DestroyOptions from a Terraform model.
//
// This function converts the destroy configuration from the resource model into
// the helmfile DestroyOptions structure. It extracts cascade behavior, concurrency,
// wait settings, and timeout values.
//
// Parameters:
//   - model: Resource model from Terraform configuration
//
// Returns:
//   - *config.DestroyOptions: Converted options for helmfile destroy operation
func NewDestroyOptionsFromModel(
	model *HelmfileReleaseModel,
) *config.DestroyOptions {
	result := &config.DestroyOptions{}
	destroyOptions := model.Destroy
	if destroyOptions.IsNull() || destroyOptions.IsUnknown() {
		return result
	}
	result.Cascade = destroyOptions.Cascade.ValueString()
	result.Concurrency = int(destroyOptions.Concurrency.ValueInt64())
	result.DeleteWait = destroyOptions.Wait.ValueBool()
	result.DeleteTimeout = int(destroyOptions.Timeout.ValueInt64())
	return result
}

//nolint:gocritic // Interface implementation
func (r *HelmfileReleaseResource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	tflog.Debug(ctx, "=HLMFL=> Modifying plan for Helmfile release resource")
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

// updateReleaseState updates the Terraform state for a helmfile release by calculating
// the build digest (SHA256 checksum) and retrieving the list of managed releases.
//
// This function performs two main operations:
//  1. Calculates a SHA256 checksum of the helmfile build output to detect changes
//  2. Retrieves the current list of releases managed by helmfile
//
// The calculated digest and release list are stored in the provided HelmfileReleaseModel
// data structure, which is then persisted to Terraform state.
//
// Parameters:
//   - ctx: Context for logging and cancellation
//   - executor: HelmfileExecutor instance for running helmfile commands
//   - options: Options containing helmfile configuration and paths
//   - data: Pointer to the resource model that will be updated with state information
//
// Returns diagnostics containing any errors encountered during digest calculation
// or release list retrieval.
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

// releaseDigest calculates a SHA256 checksum of the helmfile build output to detect
// configuration changes.
//
// This function runs `helmfile build` to generate the complete helmfile configuration,
// filters out non-deterministic log lines (timestamps, version info), and calculates
// a SHA256 hash of the resulting output. This digest is used by Terraform to detect
// when the helmfile configuration has changed and requires an update.
//
// The digest is stored in the resource's Sha256Checksum attribute and is used by the
// BuildDigestModifier plan modifier to determine if an update is needed.
//
// Parameters:
//   - ctx: Context for logging and cancellation
//   - executor: HelmfileExecutor instance for running helmfile build
//   - options: Options containing helmfile configuration and paths
//
// Returns:
//   - string: Hexadecimal SHA256 checksum of the filtered build output
//   - diag.Diagnostics: Any errors encountered during build or processing
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

// jsonToReleasesListValue parses JSON output from `helmfile list` and converts it into
// a Terraform ListValue containing release information.
//
// This function unmarshals the JSON output from `helmfile list -o json` into a slice of
// app.HelmRelease structures, then converts each release into a ReleasesListValue with
// typed Terraform attributes. The resulting list is stored in the resource's releases_list
// computed attribute.
//
// Each release contains:
//   - name: Release name
//   - namespace: Kubernetes namespace
//   - chart: Helm chart reference
//   - version: Chart version
//   - enabled: Whether the release is enabled in helmfile
//   - installed: Whether the release is currently installed
//   - labels: Map of labels parsed from the comma-separated labels string
//
// Labels are parsed from the format "key1:value1,key2:value2" into a map.
//
// Parameters:
//   - ctx: Context for type operations
//   - output: JSON string output from `helmfile list -o json`
//
// Returns:
//   - *basetypes.ListValue: List of ReleasesListValue objects, or nil on error
//   - diag.Diagnostics: Any errors encountered during parsing or conversion
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
