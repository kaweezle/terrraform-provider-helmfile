// Copyright (c) Antoine Martin
// SPDX-License-Identifier: MIT

package provider_helmfile

import (
	"context"

	"github.com/helmfile/helmfile/pkg/app"
	"github.com/helmfile/helmfile/pkg/app/version"
	"github.com/helmfile/helmfile/pkg/config"
	"github.com/kaweezle/terraform-provider-helmfile/pkg/helmfile"
)

var _ helmfile.HelmfileExecutor = (*HelmfileProvider)(nil)

type HelmfileProvider struct{}

func NewHelmfileProvider(helmBinaryPath string, performInit bool) (*HelmfileProvider, error) {
	return &HelmfileProvider{}, nil
}

func (p *HelmfileProvider) Apply(ctx context.Context, options app.ConfigProvider, applyOptions app.ApplyConfigProvider) (*helmfile.Result, error) {
	capture := helmfile.NewOutputCapture()
	_ = helmfile.CreateCaptureLogger(capture)

	helmfileApp := app.New(options)

	// Run apply operation
	err := helmfileApp.Apply(applyOptions)

	// Get captured output and prepend debug info
	output := capture.String()

	if err != nil {
		return &helmfile.Result{
			Output:   output,
			ExitCode: 1,
			Error:    err,
		}, err
	}

	return &helmfile.Result{
		Output:   output,
		ExitCode: 0,
		Error:    nil,
	}, nil
}

func (p *HelmfileProvider) Diff(ctx context.Context, options app.ConfigProvider, diffOptions app.DiffConfigProvider) (*helmfile.Result, error) {
	capture := helmfile.NewOutputCapture()
	_ = helmfile.CreateCaptureLogger(capture)

	helmfileApp := app.New(options)

	// Run apply operation
	err := helmfileApp.Diff(diffOptions)

	// Get captured output and prepend debug info
	output := capture.String()

	if err != nil {
		return &helmfile.Result{
			Output:   output,
			ExitCode: 1,
			Error:    err,
		}, err
	}

	return &helmfile.Result{
		Output:   output,
		ExitCode: 0,
		Error:    nil,
	}, nil
}

func (p *HelmfileProvider) Template(ctx context.Context, options app.ConfigProvider, templateOptions app.TemplateConfigProvider) (*helmfile.Result, error) {
	capture := helmfile.NewOutputCapture()
	_ = helmfile.CreateCaptureLogger(capture)

	helmfileApp := app.New(options)

	// Run apply operation
	err := helmfileApp.Template(templateOptions)

	// Get captured output and prepend debug info
	output := capture.String()

	if err != nil {
		return &helmfile.Result{
			Output:   output,
			ExitCode: 1,
			Error:    err,
		}, err
	}

	return &helmfile.Result{
		Output:   output,
		ExitCode: 0,
		Error:    nil,
	}, nil
}

func (p *HelmfileProvider) Destroy(ctx context.Context, options app.ConfigProvider, destroyOptions app.DestroyConfigProvider) (*helmfile.Result, error) {
	capture := helmfile.NewOutputCapture()
	_ = helmfile.CreateCaptureLogger(capture)

	helmfileApp := app.New(options)

	// Run apply operation
	err := helmfileApp.Destroy(destroyOptions)

	// Get captured output and prepend debug info
	output := capture.String()

	if err != nil {
		return &helmfile.Result{
			Output:   output,
			ExitCode: 1,
			Error:    err,
		}, err
	}

	return &helmfile.Result{
		Output:   output,
		ExitCode: 0,
		Error:    nil,
	}, nil
}

func (p *HelmfileProvider) Build(ctx context.Context, options app.ConfigProvider) (*helmfile.Result, error) {
	capture := helmfile.NewOutputCapture()
	logger := helmfile.CreateCaptureLogger(capture)

	buildOptions := config.NewBuildOptions()

	globalOptions := &config.GlobalOptions{}
	globalOptions.SetLogger(logger)
	globalImpl := config.NewGlobalImpl(globalOptions)

	buildImpl := config.NewBuildImpl(globalImpl, buildOptions)

	helmfileApp := app.New(buildImpl)
	err := helmfileApp.PrintState(buildImpl)

	// Get captured output and prepend debug info
	output := capture.String()

	if err != nil {
		return &helmfile.Result{
			Output:   output,
			ExitCode: 1,
			Error:    err,
		}, err
	}

	return &helmfile.Result{
		Output:   output,
		ExitCode: 0,
		Error:    nil,
	}, nil
}

func (p *HelmfileProvider) Version(ctx context.Context) (string, error) {
	return version.Version(), nil
}
