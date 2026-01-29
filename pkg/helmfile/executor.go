// Copyright Antoine Martin 2026
// SPDX-License-Identifier: MIT

package helmfile

import (
	"context"

	"github.com/helmfile/helmfile/pkg/app"
)

// HelmfileExecutor defines the interface for executing helmfile commands
type HelmfileExecutor interface {
	// Init runs helmfile init to install dependencies
	Init(ctx context.Context, options app.ConfigProvider) (string, error)

	// Apply runs helmfile apply/sync to deploy releases
	Apply(
		ctx context.Context,
		options app.ConfigProvider,
		applyOptions app.ApplyConfigProvider,
	) (string, error)

	// Diff runs helmfile diff to show changes
	Diff(
		ctx context.Context,
		options app.ConfigProvider,
		diffOptions app.DiffConfigProvider,
	) (string, error)

	// Template runs helmfile template to render manifests
	Template(
		ctx context.Context,
		options app.ConfigProvider,
		templateOptions app.TemplateConfigProvider,
	) (string, error)

	// Destroy runs helmfile destroy to delete releases
	Destroy(
		ctx context.Context,
		options app.ConfigProvider,
		destroyOptions app.DestroyConfigProvider,
	) (string, error)

	// Build runs helmfile build to validate configuration
	Build(ctx context.Context, options app.ConfigProvider) (string, error)
	// Version returns the helmfile version
	Version(ctx context.Context) (string, error)
}
