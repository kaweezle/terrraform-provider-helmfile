// Copyright Antoine Martin 2026
// SPDX-License-Identifier: MIT

// cSpell: words stretchr dupl

//nolint:dupl // Similar test code for different functions
package resource_helmfile_release

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGlobalOptionsFromModel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		model        *HelmfileReleaseModel
		validateFunc func(*testing.T, *HelmfileReleaseModel)
		name         string
		wantErr      bool
	}{
		{
			name: "basic configuration",
			model: &HelmfileReleaseModel{
				Args: types.ListValueMust(
					types.StringType,
					[]attr.Value{types.StringValue("--debug")},
				),
				LogLevel:    types.StringValue("info"),
				Kubeconfig:  types.StringValue("/path/to/kubeconfig"),
				Environment: types.StringValue("production"),
				EnvVars:     types.MapValueMust(types.StringType, map[string]attr.Value{}),
				Overrides:   OverridesValue{},
			},
			wantErr: false,
			validateFunc: func(t *testing.T, model *HelmfileReleaseModel) {
				t.Helper()
				ctx := context.Background()
				options, diags := NewOptionsFromModel(ctx, model)

				require.False(t, diags.HasError(), "Should not have diagnostics errors")
				require.NotNil(t, options, "GlobalOptions should not be nil")

				assert.Equal(t, "info", options.LogLevel())
				assert.Equal(t, "/path/to/kubeconfig", options.Kubeconfig())
				assert.Equal(t, "production", options.Environment())
			},
		},
		{
			name: "with environment variables",
			model: &HelmfileReleaseModel{
				Args:        types.ListValueMust(types.StringType, []attr.Value{}),
				LogLevel:    types.StringValue("debug"),
				Kubeconfig:  types.StringValue(""),
				Environment: types.StringValue(""),
				EnvVars: types.MapValueMust(types.StringType, map[string]attr.Value{
					"KEY1": types.StringValue("value1"),
					"KEY2": types.StringValue("value2"),
				}),
				Overrides: OverridesValue{},
			},
			wantErr: false,
			validateFunc: func(t *testing.T, model *HelmfileReleaseModel) {
				t.Helper()
				ctx := context.Background()
				options, diags := NewOptionsFromModel(ctx, model)

				require.False(t, diags.HasError())
				require.NotNil(t, options)

				assert.Equal(t, "debug", options.LogLevel())
				assert.Contains(t, options.EnvVars(), "KEY1")
				assert.Contains(t, options.EnvVars(), "KEY2")
			},
		},
		{
			name: "with null environment variables",
			model: &HelmfileReleaseModel{
				Args:        types.ListValueMust(types.StringType, []attr.Value{}),
				LogLevel:    types.StringValue("warn"),
				Kubeconfig:  types.StringValue(""),
				Environment: types.StringValue(""),
				EnvVars:     types.MapNull(types.StringType),
				Overrides:   OverridesValue{},
			},
			wantErr: false,
			validateFunc: func(t *testing.T, model *HelmfileReleaseModel) {
				t.Helper()
				ctx := context.Background()
				options, diags := NewOptionsFromModel(ctx, model)

				require.False(t, diags.HasError())
				require.NotNil(t, options)

				assert.Equal(t, "warn", options.LogLevel())
				assert.Empty(t, options.EnvVars())
			},
		},
		{
			name: "with multiple args",
			model: &HelmfileReleaseModel{
				Args: types.ListValueMust(types.StringType, []attr.Value{
					types.StringValue("--debug"),
					types.StringValue("--verbose"),
					types.StringValue("--flag=value"),
				}),
				LogLevel:    types.StringValue("info"),
				Kubeconfig:  types.StringValue(""),
				Environment: types.StringValue(""),
				EnvVars:     types.MapValueMust(types.StringType, map[string]attr.Value{}),
				Overrides:   OverridesValue{},
			},
			wantErr: false,
			validateFunc: func(t *testing.T, model *HelmfileReleaseModel) {
				t.Helper()
				ctx := context.Background()
				options, diags := NewOptionsFromModel(ctx, model)

				require.False(t, diags.HasError())
				require.NotNil(t, options)

				// Args should be joined with spaces
				assert.Contains(t, options.Args(), "--debug")
			},
		},
		{
			name: "with overrides - not null",
			model: &HelmfileReleaseModel{
				Args:        types.ListValueMust(types.StringType, []attr.Value{}),
				LogLevel:    types.StringValue("info"),
				Kubeconfig:  types.StringValue(""),
				Environment: types.StringValue(""),
				EnvVars:     types.MapValueMust(types.StringType, map[string]attr.Value{}),
				Overrides: OverridesValue{
					state:                      attr.ValueStateKnown,
					DisableForceUpdate:         types.BoolValue(true),
					EnforcePluginVerification:  types.BoolValue(true),
					HelmBinaryPath:             types.StringValue("/custom/helm"),
					HelmOciPlainHttp:           types.BoolValue(false),
					KustomizeBinaryPath:        types.StringValue("/custom/kustomize"),
					SkipDeps:                   types.BoolValue(true),
					SkipRefresh:                types.BoolValue(false),
					StripArgsValuesOnExitError: types.BoolValue(true),
				},
			},
			wantErr: false,
			validateFunc: func(t *testing.T, model *HelmfileReleaseModel) {
				t.Helper()
				ctx := context.Background()
				options, diags := NewOptionsFromModel(ctx, model)

				require.False(t, diags.HasError())
				require.NotNil(t, options)

				assert.True(t, options.DisableForceUpdate())
				assert.Equal(t, "/custom/helm", options.HelmBinary())
				assert.Equal(t, "/custom/kustomize", options.KustomizeBinary())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.validateFunc(t, tt.model)
		})
	}
}

func TestNewApplyOptionsFromModel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		model        *HelmfileReleaseModel
		validateFunc func(*testing.T, *HelmfileReleaseModel)
		name         string
		wantErr      bool
	}{
		{
			name: "basic apply options",
			model: &HelmfileReleaseModel{
				Cascade:                  types.StringValue("background"),
				Concurrency:              types.Int64Value(4),
				Context:                  types.Int64Value(3),
				DetailedExitcode:         types.BoolValue(true),
				DiffArgs:                 types.StringValue("--color=always"),
				EnforceNeedsAreInstalled: types.BoolValue(true),
				HideNotes:                types.BoolValue(false),
				IncludeNeeds:             types.BoolValue(true),
				IncludeTests:             types.BoolValue(false),
				IncludeTransitiveNeeds:   types.BoolValue(true),
				NoHooks:                  types.BoolValue(false),
				Output:                   types.StringValue("json"),
				PostRenderer:             types.StringValue("kustomize"),
				PostRendererArgs:         types.ListNull(types.StringType),
				ResetValues:              types.BoolValue(false),
				ReuseValues:              types.BoolValue(true),
				ShowSecrets:              types.BoolValue(false),
				SkipCleanup:              types.BoolValue(false),
				SkipCrds:                 types.BoolValue(false),
				SkipDiffOnInstall:        types.BoolValue(false),
				SkipNeeds:                types.BoolValue(false),
				SkipSchemaValidation:     types.BoolValue(false),
				StripTrailingCr:          types.BoolValue(true),
				SuppressDiff:             types.BoolValue(false),
				SuppressSecrets:          types.BoolValue(true),
				SyncArgs:                 types.StringValue("--force"),
				SyncReleaseLabels:        types.BoolValue(true),
				TakeOwnership:            types.BoolValue(false),
				Validate:                 types.BoolValue(true),
				Wait:                     types.BoolValue(true),
				WaitForJobs:              types.BoolValue(true),
				WaitRetries:              types.Int64Value(5),
				Set:                      types.ListNull(types.StringType),
				Suppress:                 types.ListNull(types.StringType),
				SuppressOutputLineRegex:  types.ListNull(types.StringType),
				Values:                   types.ListNull(types.StringType),
			},
			wantErr: false,
			validateFunc: func(t *testing.T, model *HelmfileReleaseModel) {
				t.Helper()
				ctx := context.Background()
				applyOpts, diags := NewApplyOptionsFromModel(ctx, model)

				require.False(t, diags.HasError(), "Should not have diagnostics errors")
				require.NotNil(t, applyOpts, "ApplyOptions should not be nil")

				assert.Equal(t, "background", applyOpts.Cascade)
				assert.Equal(t, 4, applyOpts.Concurrency)
				assert.Equal(t, 3, applyOpts.Context)
				assert.True(t, applyOpts.DetailedExitcode)
				assert.Equal(t, "--color=always", applyOpts.DiffArgs)
				assert.True(t, applyOpts.EnforceNeedsAreInstalled)
				assert.False(t, applyOpts.HideNotes)
				assert.True(t, applyOpts.IncludeNeeds)
				assert.False(t, applyOpts.IncludeTests)
				assert.True(t, applyOpts.IncludeTransitiveNeeds)
				assert.False(t, applyOpts.NoHooks)
				assert.Equal(t, "json", applyOpts.Output)
				assert.Equal(t, "kustomize", applyOpts.PostRenderer)
				assert.False(t, applyOpts.ResetValues)
				assert.True(t, applyOpts.ReuseValues)
				assert.False(t, applyOpts.ShowSecrets)
				assert.False(t, applyOpts.SkipCleanup)
				assert.False(t, applyOpts.SkipCRDs)
				assert.False(t, applyOpts.SkipDiffOnInstall)
				assert.False(t, applyOpts.SkipNeeds)
				assert.False(t, applyOpts.SkipSchemaValidation)
				assert.True(t, applyOpts.StripTrailingCR)
				assert.False(t, applyOpts.SuppressDiff)
				assert.True(t, applyOpts.SuppressSecrets)
				assert.Equal(t, "--force", applyOpts.SyncArgs)
				assert.True(t, applyOpts.SyncReleaseLabels)
				assert.False(t, applyOpts.TakeOwnership)
				assert.True(t, applyOpts.Validate)
				assert.True(t, applyOpts.Wait)
				assert.True(t, applyOpts.WaitForJobs)
				assert.Equal(t, 5, applyOpts.WaitRetries)
			},
		},
		{
			name: "with post renderer args",
			model: &HelmfileReleaseModel{
				Cascade:                  types.StringValue(""),
				Concurrency:              types.Int64Value(0),
				Context:                  types.Int64Value(0),
				DetailedExitcode:         types.BoolValue(false),
				DiffArgs:                 types.StringValue(""),
				EnforceNeedsAreInstalled: types.BoolValue(false),
				HideNotes:                types.BoolValue(false),
				IncludeNeeds:             types.BoolValue(false),
				IncludeTests:             types.BoolValue(false),
				IncludeTransitiveNeeds:   types.BoolValue(false),
				NoHooks:                  types.BoolValue(false),
				Output:                   types.StringValue(""),
				PostRenderer:             types.StringValue(""),
				PostRendererArgs: types.ListValueMust(types.StringType, []attr.Value{
					types.StringValue("--arg1"),
					types.StringValue("--arg2=value"),
				}),
				ResetValues:             types.BoolValue(false),
				ReuseValues:             types.BoolValue(false),
				ShowSecrets:             types.BoolValue(false),
				SkipCleanup:             types.BoolValue(false),
				SkipCrds:                types.BoolValue(false),
				SkipDiffOnInstall:       types.BoolValue(false),
				SkipNeeds:               types.BoolValue(false),
				SkipSchemaValidation:    types.BoolValue(false),
				StripTrailingCr:         types.BoolValue(false),
				SuppressDiff:            types.BoolValue(false),
				SuppressSecrets:         types.BoolValue(false),
				SyncArgs:                types.StringValue(""),
				SyncReleaseLabels:       types.BoolValue(false),
				TakeOwnership:           types.BoolValue(false),
				Validate:                types.BoolValue(false),
				Wait:                    types.BoolValue(false),
				WaitForJobs:             types.BoolValue(false),
				WaitRetries:             types.Int64Value(0),
				Set:                     types.ListNull(types.StringType),
				Suppress:                types.ListNull(types.StringType),
				SuppressOutputLineRegex: types.ListNull(types.StringType),
				Values:                  types.ListNull(types.StringType),
			},
			wantErr: false,
			validateFunc: func(t *testing.T, model *HelmfileReleaseModel) {
				t.Helper()
				ctx := context.Background()
				applyOpts, diags := NewApplyOptionsFromModel(ctx, model)

				require.False(t, diags.HasError())
				require.NotNil(t, applyOpts)

				require.Len(t, applyOpts.PostRendererArgs, 2)
				assert.Equal(t, "--arg1", applyOpts.PostRendererArgs[0])
				assert.Equal(t, "--arg2=value", applyOpts.PostRendererArgs[1])
			},
		},
		{
			name: "with set values",
			model: &HelmfileReleaseModel{
				Cascade:                  types.StringValue(""),
				Concurrency:              types.Int64Value(0),
				Context:                  types.Int64Value(0),
				DetailedExitcode:         types.BoolValue(false),
				DiffArgs:                 types.StringValue(""),
				EnforceNeedsAreInstalled: types.BoolValue(false),
				HideNotes:                types.BoolValue(false),
				IncludeNeeds:             types.BoolValue(false),
				IncludeTests:             types.BoolValue(false),
				IncludeTransitiveNeeds:   types.BoolValue(false),
				NoHooks:                  types.BoolValue(false),
				Output:                   types.StringValue(""),
				PostRenderer:             types.StringValue(""),
				PostRendererArgs:         types.ListNull(types.StringType),
				ResetValues:              types.BoolValue(false),
				ReuseValues:              types.BoolValue(false),
				ShowSecrets:              types.BoolValue(false),
				SkipCleanup:              types.BoolValue(false),
				SkipCrds:                 types.BoolValue(false),
				SkipDiffOnInstall:        types.BoolValue(false),
				SkipNeeds:                types.BoolValue(false),
				SkipSchemaValidation:     types.BoolValue(false),
				StripTrailingCr:          types.BoolValue(false),
				SuppressDiff:             types.BoolValue(false),
				SuppressSecrets:          types.BoolValue(false),
				SyncArgs:                 types.StringValue(""),
				SyncReleaseLabels:        types.BoolValue(false),
				TakeOwnership:            types.BoolValue(false),
				Validate:                 types.BoolValue(false),
				Wait:                     types.BoolValue(false),
				WaitForJobs:              types.BoolValue(false),
				WaitRetries:              types.Int64Value(0),
				Set: types.ListValueMust(types.StringType, []attr.Value{
					types.StringValue("key1=value1"),
					types.StringValue("key2=value2"),
				}),
				Suppress:                types.ListNull(types.StringType),
				SuppressOutputLineRegex: types.ListNull(types.StringType),
				Values:                  types.ListNull(types.StringType),
			},
			wantErr: false,
			validateFunc: func(t *testing.T, model *HelmfileReleaseModel) {
				t.Helper()
				ctx := context.Background()
				applyOpts, diags := NewApplyOptionsFromModel(ctx, model)

				require.False(t, diags.HasError())
				require.NotNil(t, applyOpts)

				require.Len(t, applyOpts.Set, 2)
				assert.Equal(t, "key1=value1", applyOpts.Set[0])
				assert.Equal(t, "key2=value2", applyOpts.Set[1])
			},
		},
		{
			name: "with suppress and suppress output line regex",
			model: &HelmfileReleaseModel{
				Cascade:                  types.StringValue(""),
				Concurrency:              types.Int64Value(0),
				Context:                  types.Int64Value(0),
				DetailedExitcode:         types.BoolValue(false),
				DiffArgs:                 types.StringValue(""),
				EnforceNeedsAreInstalled: types.BoolValue(false),
				HideNotes:                types.BoolValue(false),
				IncludeNeeds:             types.BoolValue(false),
				IncludeTests:             types.BoolValue(false),
				IncludeTransitiveNeeds:   types.BoolValue(false),
				NoHooks:                  types.BoolValue(false),
				Output:                   types.StringValue(""),
				PostRenderer:             types.StringValue(""),
				PostRendererArgs:         types.ListNull(types.StringType),
				ResetValues:              types.BoolValue(false),
				ReuseValues:              types.BoolValue(false),
				ShowSecrets:              types.BoolValue(false),
				SkipCleanup:              types.BoolValue(false),
				SkipCrds:                 types.BoolValue(false),
				SkipDiffOnInstall:        types.BoolValue(false),
				SkipNeeds:                types.BoolValue(false),
				SkipSchemaValidation:     types.BoolValue(false),
				StripTrailingCr:          types.BoolValue(false),
				SuppressDiff:             types.BoolValue(false),
				SuppressSecrets:          types.BoolValue(false),
				SyncArgs:                 types.StringValue(""),
				SyncReleaseLabels:        types.BoolValue(false),
				TakeOwnership:            types.BoolValue(false),
				Validate:                 types.BoolValue(false),
				Wait:                     types.BoolValue(false),
				WaitForJobs:              types.BoolValue(false),
				WaitRetries:              types.Int64Value(0),
				Set:                      types.ListNull(types.StringType),
				Suppress: types.ListValueMust(types.StringType, []attr.Value{
					types.StringValue("secret1"),
					types.StringValue("secret2"),
				}),
				SuppressOutputLineRegex: types.ListValueMust(types.StringType, []attr.Value{
					types.StringValue(".*password.*"),
					types.StringValue(".*token.*"),
				}),
				Values: types.ListNull(types.StringType),
			},
			wantErr: false,
			validateFunc: func(t *testing.T, model *HelmfileReleaseModel) {
				t.Helper()
				ctx := context.Background()
				applyOpts, diags := NewApplyOptionsFromModel(ctx, model)

				require.False(t, diags.HasError())
				require.NotNil(t, applyOpts)

				require.Len(t, applyOpts.Suppress, 2)
				assert.Equal(t, "secret1", applyOpts.Suppress[0])
				assert.Equal(t, "secret2", applyOpts.Suppress[1])

				require.Len(t, applyOpts.SuppressOutputLineRegex, 2)
				assert.Equal(t, ".*password.*", applyOpts.SuppressOutputLineRegex[0])
				assert.Equal(t, ".*token.*", applyOpts.SuppressOutputLineRegex[1])
			},
		},
		{
			name: "with values files",
			model: &HelmfileReleaseModel{
				Cascade:                  types.StringValue(""),
				Concurrency:              types.Int64Value(0),
				Context:                  types.Int64Value(0),
				DetailedExitcode:         types.BoolValue(false),
				DiffArgs:                 types.StringValue(""),
				EnforceNeedsAreInstalled: types.BoolValue(false),
				HideNotes:                types.BoolValue(false),
				IncludeNeeds:             types.BoolValue(false),
				IncludeTests:             types.BoolValue(false),
				IncludeTransitiveNeeds:   types.BoolValue(false),
				NoHooks:                  types.BoolValue(false),
				Output:                   types.StringValue(""),
				PostRenderer:             types.StringValue(""),
				PostRendererArgs:         types.ListNull(types.StringType),
				ResetValues:              types.BoolValue(false),
				ReuseValues:              types.BoolValue(false),
				ShowSecrets:              types.BoolValue(false),
				SkipCleanup:              types.BoolValue(false),
				SkipCrds:                 types.BoolValue(false),
				SkipDiffOnInstall:        types.BoolValue(false),
				SkipNeeds:                types.BoolValue(false),
				SkipSchemaValidation:     types.BoolValue(false),
				StripTrailingCr:          types.BoolValue(false),
				SuppressDiff:             types.BoolValue(false),
				SuppressSecrets:          types.BoolValue(false),
				SyncArgs:                 types.StringValue(""),
				SyncReleaseLabels:        types.BoolValue(false),
				TakeOwnership:            types.BoolValue(false),
				Validate:                 types.BoolValue(false),
				Wait:                     types.BoolValue(false),
				WaitForJobs:              types.BoolValue(false),
				WaitRetries:              types.Int64Value(0),
				Set:                      types.ListNull(types.StringType),
				Suppress:                 types.ListNull(types.StringType),
				SuppressOutputLineRegex:  types.ListNull(types.StringType),
				Values: types.ListValueMust(types.StringType, []attr.Value{
					types.StringValue("values1.yaml"),
					types.StringValue("values2.yaml"),
				}),
			},
			wantErr: false,
			validateFunc: func(t *testing.T, model *HelmfileReleaseModel) {
				t.Helper()
				ctx := context.Background()
				applyOpts, diags := NewApplyOptionsFromModel(ctx, model)

				require.False(t, diags.HasError())
				require.NotNil(t, applyOpts)

				require.Len(t, applyOpts.Values, 2)
				assert.Equal(t, "values1.yaml", applyOpts.Values[0])
				assert.Equal(t, "values2.yaml", applyOpts.Values[1])
			},
		},
		{
			name: "with all list options as null",
			model: &HelmfileReleaseModel{
				Cascade:                  types.StringValue(""),
				Concurrency:              types.Int64Value(1),
				Context:                  types.Int64Value(0),
				DetailedExitcode:         types.BoolValue(false),
				DiffArgs:                 types.StringValue(""),
				EnforceNeedsAreInstalled: types.BoolValue(false),
				HideNotes:                types.BoolValue(false),
				IncludeNeeds:             types.BoolValue(false),
				IncludeTests:             types.BoolValue(false),
				IncludeTransitiveNeeds:   types.BoolValue(false),
				NoHooks:                  types.BoolValue(false),
				Output:                   types.StringValue(""),
				PostRenderer:             types.StringValue(""),
				PostRendererArgs:         types.ListNull(types.StringType),
				ResetValues:              types.BoolValue(false),
				ReuseValues:              types.BoolValue(false),
				ShowSecrets:              types.BoolValue(false),
				SkipCleanup:              types.BoolValue(false),
				SkipCrds:                 types.BoolValue(false),
				SkipDiffOnInstall:        types.BoolValue(false),
				SkipNeeds:                types.BoolValue(false),
				SkipSchemaValidation:     types.BoolValue(false),
				StripTrailingCr:          types.BoolValue(false),
				SuppressDiff:             types.BoolValue(false),
				SuppressSecrets:          types.BoolValue(false),
				SyncArgs:                 types.StringValue(""),
				SyncReleaseLabels:        types.BoolValue(false),
				TakeOwnership:            types.BoolValue(false),
				Validate:                 types.BoolValue(false),
				Wait:                     types.BoolValue(false),
				WaitForJobs:              types.BoolValue(false),
				WaitRetries:              types.Int64Value(0),
				Set:                      types.ListNull(types.StringType),
				Suppress:                 types.ListNull(types.StringType),
				SuppressOutputLineRegex:  types.ListNull(types.StringType),
				Values:                   types.ListNull(types.StringType),
			},
			wantErr: false,
			validateFunc: func(t *testing.T, model *HelmfileReleaseModel) {
				t.Helper()
				ctx := context.Background()
				applyOpts, diags := NewApplyOptionsFromModel(ctx, model)

				require.False(t, diags.HasError())
				require.NotNil(t, applyOpts)

				assert.Nil(t, applyOpts.PostRendererArgs)
				assert.Nil(t, applyOpts.Set)
				assert.Nil(t, applyOpts.Suppress)
				assert.Nil(t, applyOpts.SuppressOutputLineRegex)
				assert.Nil(t, applyOpts.Values)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.validateFunc(t, tt.model)
		})
	}
}
