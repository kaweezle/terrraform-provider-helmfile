// Copyright Antoine Martin 2026
// SPDX-License-Identifier: MIT

// cSpell: words helmexec cliv3 cliv4

package helmfile

import (
	"context"
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/helmfile/helmfile/pkg/app"
	"github.com/helmfile/helmfile/pkg/app/version"
	"github.com/helmfile/helmfile/pkg/config"
	"github.com/helmfile/helmfile/pkg/helmexec"
	"go.uber.org/zap"
	cliv3 "helm.sh/helm/v3/pkg/cli"
	cliv4 "helm.sh/helm/v4/pkg/cli"
)

var _ HelmfileExecutor = (*HelmfileLibraryExecutor)(nil)

type HelmPlugin struct {
	Name    string
	Version string
	Repo    string
}

// HelmfileLibraryExecutor is the Helmfile provider implementation
type HelmfileLibraryExecutor struct {
	globalOptions     *config.GlobalOptions
	additionalPlugins []HelmPlugin
}

// NewHelmfileProvider creates a new HelmfileProvider instance from the given HelmfileModel
func NewHelmfileLibraryExecutor(options *config.GlobalOptions, additionalPlugins []HelmPlugin) *HelmfileLibraryExecutor {
	return &HelmfileLibraryExecutor{
		globalOptions:     options,
		additionalPlugins: additionalPlugins,
	}
}

func (p *HelmfileLibraryExecutor) createGlobalOptionsFromConfigProvider(options app.ConfigProvider, logger *zap.SugaredLogger) (*config.GlobalImpl, error) {
	globalOptions := *p.globalOptions
	if options != nil {
		if options.Args() != "" {
			globalOptions.Args = options.Args()
		}
		if options.HelmBinary() != "" {
			globalOptions.HelmBinary = options.HelmBinary()
		}
		if options.KustomizeBinary() != "" {
			globalOptions.KustomizeBinary = options.KustomizeBinary()
		}

		globalOptions.StripArgsValuesOnExitError = globalOptions.StripArgsValuesOnExitError || options.StripArgsValuesOnExitError()
		globalOptions.DisableForceUpdate = globalOptions.DisableForceUpdate || options.DisableForceUpdate()
		globalOptions.EnforcePluginVerification = globalOptions.EnforcePluginVerification || options.EnforcePluginVerification()
		globalOptions.HelmOCIPlainHTTP = globalOptions.HelmOCIPlainHTTP || options.HelmOCIPlainHTTP()
		globalOptions.SkipDeps = globalOptions.SkipDeps || options.SkipDeps()
		globalOptions.SkipRefresh = globalOptions.SkipRefresh || options.SkipRefresh()

		if options.FileOrDir() != "" {
			globalOptions.File = options.FileOrDir()
		}
		if options.KubeContext() != "" {
			globalOptions.KubeContext = options.KubeContext()
		}
		if options.Namespace() != "" {
			globalOptions.Namespace = options.Namespace()
		}
		if options.Chart() != "" {
			globalOptions.Chart = options.Chart()
		}
		if len(options.Selectors()) > 0 {
			globalOptions.Selector = options.Selectors()
		}
		if len(options.StateValuesFiles()) > 0 {
			globalOptions.StateValuesFile = options.StateValuesFiles()
		}
		if options.Kubeconfig() != "" {
			globalOptions.Kubeconfig = options.Kubeconfig()
		}
		if options.Env() != "" {
			globalOptions.Environment = options.Env()
		}
	}
	globalOptions.SetLogger(logger)
	result := config.NewGlobalImpl(&globalOptions)
	if options != nil {
		stateValueSets := options.StateValuesSet()
		if stateValueSets != nil {
			result.SetSet(stateValueSets)
		}
	}
	err := config.NewCLIConfigImpl(result)
	if err != nil {
		return nil, fmt.Errorf("error while creating global options: %w", err)
	}
	return result, nil
}

func (p *HelmfileLibraryExecutor) InstallAdditionalPlugins(ctx context.Context, logger *zap.SugaredLogger) error {
	logger.Debug("Installing additional plugins...")
	runner := &helmexec.ShellRunner{
		Logger:                     logger,
		Ctx:                        ctx,
		StripArgsValuesOnExitError: p.globalOptions.StripArgsValuesOnExitError,
	}
	helm, err := helmexec.New(p.globalOptions.HelmBinary, helmexec.HelmExecOptions{}, logger, "", "", runner)
	if err != nil {
		return err
	}

	// Use version-specific cli based on detected Helm version
	var pluginsDir string
	if helm.IsHelm3() {
		pluginsDir = cliv3.New().PluginsDirectory
	} else {
		pluginsDir = cliv4.New().PluginsDirectory
	}

	for _, p := range p.additionalPlugins {
		logger.Debugf("Checking installation of plugin %s/%s/%s in directory %s", p.Name, p.Repo, p.Version, pluginsDir)
		pluginVersion, err := helmexec.GetPluginVersion(p.Name, pluginsDir)
		if err != nil {
			if !strings.Contains(err.Error(), "not installed") {
				return err
			}

			logger.Debugf("Plugin %s/%s/%s is not installed, installing...", p.Name, p.Repo, p.Version)
			err = helm.AddPlugin(p.Name, p.Repo, p.Version)
			if err != nil {
				return err
			}
			pluginVersion, err = helmexec.GetPluginVersion(p.Name, pluginsDir)
			if err != nil {
				logger.Errorf("Error while getting version of just installed plugin: %s/%s/%s: %w", p.Name, p.Repo, p.Version, err)
				return fmt.Errorf("error while getting version of just installed plugin: %s/%s/%s: %w", p.Name, p.Repo, p.Version, err)
			}
		}
		requiredVersion, _ := semver.NewVersion(p.Version)
		if pluginVersion.LessThan(requiredVersion) {
			err = helm.UpdatePlugin(p.Name)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *HelmfileLibraryExecutor) Init(ctx context.Context, options app.ConfigProvider) (string, error) {
	capture := NewOutputCapture()
	logger := CreateCaptureLogger(capture)

	globalOptions, err := p.createGlobalOptionsFromConfigProvider(options, logger)
	if err != nil {
		return "", fmt.Errorf("error performing init: %w", err)
	}
	helmfileApp := app.New(globalOptions)

	initOptions := &config.InitOptions{
		Force: true,
	}
	initImpl := config.NewInitImpl(globalOptions, initOptions)

	helmfileApp = app.New(initImpl)
	err = helmfileApp.Init(initImpl)
	if err == nil {
		err = p.InstallAdditionalPlugins(ctx, logger)
	}

	return capture.String(), err
}

func (p *HelmfileLibraryExecutor) Apply(ctx context.Context, options app.ConfigProvider, applyOptions app.ApplyConfigProvider) (string, error) {
	capture := NewOutputCapture()
	_ = CreateCaptureLogger(capture)

	helmfileApp := app.New(options)

	// Run apply operation
	err := helmfileApp.Apply(applyOptions)

	return capture.String(), err
}

func (p *HelmfileLibraryExecutor) Diff(ctx context.Context, options app.ConfigProvider, diffOptions app.DiffConfigProvider) (string, error) {
	capture := NewOutputCapture()
	_ = CreateCaptureLogger(capture)

	helmfileApp := app.New(options)

	// Run apply operation
	err := helmfileApp.Diff(diffOptions)

	return capture.String(), err
}

func (p *HelmfileLibraryExecutor) Template(ctx context.Context, options app.ConfigProvider, templateOptions app.TemplateConfigProvider) (string, error) {
	capture := NewOutputCapture()
	_ = CreateCaptureLogger(capture)

	helmfileApp := app.New(options)

	// Run apply operation
	err := helmfileApp.Template(templateOptions)

	return capture.String(), err
}

func (p *HelmfileLibraryExecutor) Destroy(ctx context.Context, options app.ConfigProvider, destroyOptions app.DestroyConfigProvider) (string, error) {
	capture := NewOutputCapture()
	_ = CreateCaptureLogger(capture)

	helmfileApp := app.New(options)

	// Run apply operation
	err := helmfileApp.Destroy(destroyOptions)

	// Get captured output and prepend debug info
	return capture.String(), err
}

func (p *HelmfileLibraryExecutor) Build(ctx context.Context, options app.ConfigProvider) (string, error) {
	capture := NewOutputCapture()
	logger := CreateCaptureLogger(capture)

	buildOptions := config.NewBuildOptions()

	globalOptions := &config.GlobalOptions{}
	globalOptions.SetLogger(logger)
	globalImpl := config.NewGlobalImpl(globalOptions)

	buildImpl := config.NewBuildImpl(globalImpl, buildOptions)

	helmfileApp := app.New(buildImpl)
	err := helmfileApp.PrintState(buildImpl)

	// Get captured output and prepend debug info
	return capture.String(), err
}

func (p *HelmfileLibraryExecutor) Version(ctx context.Context) (string, error) {
	return version.Version(), nil
}
