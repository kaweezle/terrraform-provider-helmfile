// Copyright Antoine Martin 2026
// SPDX-License-Identifier: MIT

package helmfile

import (
	"context"

	"github.com/helmfile/helmfile/pkg/config"
)

// CommonOptionsProvider defines common configuration options shared across all helmfile operations.
// These options include authentication (kubeconfig), environment selection, environment variables,
// and logging settings.
//
// This interface is embedded in both GlobalOptionsProvider and ResourceOptionsProvider to ensure
// consistent configuration across different levels of the provider.
type CommonOptionsProvider interface {
	Kubeconfig() string
	Environment() string
	EnvVars() map[string]string
	LogLevel() string
}

// BaseGlobalOptionsProvider defines global helmfile configuration options that apply to the
// entire helmfile execution context.
//
// These options control:
//   - Helm and Kustomize binary paths and default arguments
//   - Plugin behavior (verification, force updates)
//   - OCI registry settings (plain HTTP)
//   - Dependency and refresh handling
//
// Global options are typically configured at the provider level and inherited by all resources.
type BaseGlobalOptionsProvider interface {
	DefaultArgs() string
	HelmBinary() string
	KustomizeBinary() string
	StripArgsValuesOnExitError() bool
	DisableForceUpdate() bool
	EnforcePluginVerification() bool
	HelmOCIPlainHTTP() bool
	SkipDeps() bool
	SkipRefresh() bool
}

// GlobalOptionsProvider combines BaseGlobalOptionsProvider and CommonOptionsProvider to provide
// all configuration options needed for provider-level helmfile operations.
//
// This interface is typically used for initialization operations that don't target specific
// releases or resources.
type GlobalOptionsProvider interface {
	BaseGlobalOptionsProvider
	CommonOptionsProvider
}

// BaseResourceOptionsProvider defines resource-specific helmfile configuration options.
//
// These options control:
//   - Target helmfile file or directory
//   - Kubernetes context and namespace overrides
//   - Release filtering (selectors, chart)
//   - State value overrides (inline values and value files)
//   - Additional command-line arguments
//
// Resource options are typically configured per helmfile_release resource and override
// global settings.
type BaseResourceOptionsProvider interface {
	Args() string
	FileOrDir() string
	KubeContext() string
	Namespace() string
	Chart() string
	Selectors() []string
	StateValuesSet() map[string]any
	StateValuesFiles() []string
}

// ResourceOptionsProvider combines BaseResourceOptionsProvider and CommonOptionsProvider to
// provide all configuration options needed for resource-level helmfile operations.
//
// This interface is used when executing helmfile commands that target specific releases
// or resources, combining resource-specific settings with common authentication and
// environment configuration.
type ResourceOptionsProvider interface {
	BaseResourceOptionsProvider
	CommonOptionsProvider
}

// OptionsProvider combines all three option provider interfaces to provide complete
// configuration for helmfile operations.
//
// This interface includes:
//   - BaseGlobalOptionsProvider: Global helmfile settings
//   - BaseResourceOptionsProvider: Resource-specific settings
//   - CommonOptionsProvider: Shared authentication and environment settings
//
// It is used by executor methods that need access to both global and resource-specific
// configuration, such as Apply, Diff, Template, Destroy, Build, and List.
type OptionsProvider interface {
	BaseGlobalOptionsProvider
	BaseResourceOptionsProvider
	CommonOptionsProvider
}

// HelmfileExecutor defines the interface for executing helmfile commands.
type HelmfileExecutor interface {
	// Init runs helmfile init to install dependencies
	Init(ctx context.Context, options GlobalOptionsProvider) (string, string, error)

	// Apply runs helmfile apply/sync to deploy releases
	Apply(
		ctx context.Context,
		options OptionsProvider,
		applyOptions *config.ApplyOptions,
	) (string, string, error)

	// Diff runs helmfile diff to show changes
	Diff(
		ctx context.Context,
		options OptionsProvider,
		diffOptions *config.DiffOptions,
	) (string, string, error)

	// Template runs helmfile template to render manifests
	Template(
		ctx context.Context,
		options OptionsProvider,
		templateOptions *config.TemplateOptions,
	) (string, string, error)

	// Destroy runs helmfile destroy to delete releases
	Destroy(
		ctx context.Context,
		options OptionsProvider,
		destroyOptions *config.DestroyOptions,
	) (string, string, error)

	// Build runs helmfile build to validate configuration
	Build(ctx context.Context, options OptionsProvider, embedValues bool) (string, string, error)

	// List runs helmfile list to list releases
	List(ctx context.Context, options OptionsProvider, skipCharts bool) (string, string, error)

	// Version returns the helmfile version
	Version(ctx context.Context) (string, string, error)
}
