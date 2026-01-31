// Copyright Antoine Martin 2026
// SPDX-License-Identifier: MIT

package helmfile

import (
	"context"

	"github.com/helmfile/helmfile/pkg/config"
)

type CommonOptionsProvider interface {
	Kubeconfig() string
	Environment() string
	EnvVars() map[string]string
	LogLevel() string
}

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

type GlobalOptionsProvider interface {
	BaseGlobalOptionsProvider
	CommonOptionsProvider
}

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

type ResourceOptionsProvider interface {
	BaseResourceOptionsProvider
	CommonOptionsProvider
}
type OptionsProvider interface {
	BaseGlobalOptionsProvider
	BaseResourceOptionsProvider
	CommonOptionsProvider
}

// HelmfileExecutor defines the interface for executing helmfile commands.
type HelmfileExecutor interface {
	// Init runs helmfile init to install dependencies
	Init(ctx context.Context, options GlobalOptionsProvider) (string, error)

	// Apply runs helmfile apply/sync to deploy releases
	Apply(
		ctx context.Context,
		options OptionsProvider,
		applyOptions *config.ApplyOptions,
	) (string, error)

	// Diff runs helmfile diff to show changes
	Diff(
		ctx context.Context,
		options OptionsProvider,
		diffOptions *config.DiffOptions,
	) (string, error)

	// Template runs helmfile template to render manifests
	Template(
		ctx context.Context,
		options OptionsProvider,
		templateOptions *config.TemplateOptions,
	) (string, error)

	// Destroy runs helmfile destroy to delete releases
	Destroy(
		ctx context.Context,
		options OptionsProvider,
		destroyOptions *config.DestroyOptions,
	) (string, error)

	// Build runs helmfile build to validate configuration
	Build(ctx context.Context, options OptionsProvider, embedValues bool) (string, error)

	// Version returns the helmfile version
	Version(ctx context.Context) (string, error)
}
