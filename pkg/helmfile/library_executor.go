// Copyright Antoine Martin 2026
// SPDX-License-Identifier: MIT

// cSpell: words helmexec cliv3 cliv4

package helmfile

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"maps"
	"os"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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

// HelmfileLibraryExecutor is the Helmfile provider implementation.
type HelmfileLibraryExecutor struct {
	globalOptions     *GlobalOptions
	additionalPlugins []HelmPlugin
}

// NewHelmfileProvider creates a new HelmfileProvider instance from the given HelmfileModel.
func NewHelmfileLibraryExecutor(
	options *GlobalOptions,
	additionalPlugins []HelmPlugin,
) *HelmfileLibraryExecutor {
	return &HelmfileLibraryExecutor{
		globalOptions:     options,
		additionalPlugins: additionalPlugins,
	}
}

func (p *HelmfileLibraryExecutor) InstallAdditionalPlugins(
	ctx context.Context,
	logger *zap.SugaredLogger,
) error {
	logger.Debug("Installing additional plugins...")
	runner := &helmexec.ShellRunner{
		Logger:                     logger,
		Ctx:                        ctx,
		StripArgsValuesOnExitError: p.globalOptions.StripArgsValuesOnExitError(),
	}
	helm, err := helmexec.New(
		p.globalOptions.HelmBinary(),
		helmexec.HelmExecOptions{},
		logger,
		"",
		"",
		runner,
	)
	if err != nil {
		return fmt.Errorf("error creating helm executor: %w", err)
	}

	// Use version-specific cli based on detected Helm version
	var pluginsDir string
	if helm.IsHelm3() {
		pluginsDir = cliv3.New().PluginsDirectory
	} else {
		pluginsDir = cliv4.New().PluginsDirectory
	}

	for _, p := range p.additionalPlugins {
		logger.Debugf(
			"Checking installation of plugin %s/%s/%s in directory %s",
			p.Name,
			p.Repo,
			p.Version,
			pluginsDir,
		)
		pluginVersion, err := helmexec.GetPluginVersion(p.Name, pluginsDir)
		if err != nil {
			if !strings.Contains(err.Error(), "not installed") {
				return fmt.Errorf("error getting plugin version: %w", err)
			}

			logger.Debugf(
				"Plugin %s/%s/%s is not installed, installing...",
				p.Name,
				p.Repo,
				p.Version,
			)
			err = helm.AddPlugin(p.Name, p.Repo, p.Version)
			if err != nil {
				return fmt.Errorf("error adding plugin: %w", err)
			}
			pluginVersion, err = helmexec.GetPluginVersion(p.Name, pluginsDir)
			if err != nil {
				logger.Errorf(
					"Error while getting version of just installed plugin: %s/%s/%s: %w",
					p.Name,
					p.Repo,
					p.Version,
					err,
				)
				return fmt.Errorf(
					"error while getting version of just installed plugin: %s/%s/%s: %w",
					p.Name,
					p.Repo,
					p.Version,
					err,
				)
			}
		}
		requiredVersion, err := semver.NewVersion(p.Version)
		if err != nil {
			return fmt.Errorf("error parsing required plugin version %s: %w", p.Version, err)
		}
		if pluginVersion.LessThan(requiredVersion) {
			err = helm.UpdatePlugin(p.Name)
			if err != nil {
				return fmt.Errorf("error updating plugin: %w", err)
			}
		}
	}
	return nil
}

//nolint:gocyclo // Structure is heavy
func copyOptions(globalOptions *config.GlobalOptions, options OptionsProvider) {
	if options == nil {
		return
	}
	if options.DefaultArgs() != "" {
		globalOptions.Args = options.DefaultArgs()
	}
	if options.HelmBinary() != "" {
		globalOptions.HelmBinary = options.HelmBinary()
	}
	if options.KustomizeBinary() != "" {
		globalOptions.KustomizeBinary = options.KustomizeBinary()
	}

	globalOptions.StripArgsValuesOnExitError = globalOptions.StripArgsValuesOnExitError ||
		options.StripArgsValuesOnExitError()
	globalOptions.DisableForceUpdate = globalOptions.DisableForceUpdate ||
		options.DisableForceUpdate()
	globalOptions.EnforcePluginVerification = globalOptions.EnforcePluginVerification ||
		options.EnforcePluginVerification()
	globalOptions.HelmOCIPlainHTTP = globalOptions.HelmOCIPlainHTTP ||
		options.HelmOCIPlainHTTP()
	globalOptions.SkipDeps = globalOptions.SkipDeps || options.SkipDeps()
	globalOptions.SkipRefresh = globalOptions.SkipRefresh || options.SkipRefresh()

	if options.Args() != "" {
		globalOptions.Args = options.Args()
	}
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
	if options.Environment() != "" {
		globalOptions.Environment = options.Environment()
	}
}

func (p *HelmfileLibraryExecutor) createGlobalOptionsFromConfigProvider(
	options OptionsProvider,
	logger *zap.SugaredLogger,
) (*config.GlobalImpl, error) {
	globalOptions := config.GlobalOptions{
		Args:                       p.globalOptions.DefaultArgs(),
		HelmBinary:                 p.globalOptions.HelmBinary(),
		KustomizeBinary:            p.globalOptions.KustomizeBinary(),
		StripArgsValuesOnExitError: p.globalOptions.StripArgsValuesOnExitError(),
		DisableForceUpdate:         p.globalOptions.DisableForceUpdate(),
		EnforcePluginVerification:  p.globalOptions.EnforcePluginVerification(),
		HelmOCIPlainHTTP:           p.globalOptions.HelmOCIPlainHTTP(),
		SkipDeps:                   p.globalOptions.SkipDeps(),
		SkipRefresh:                p.globalOptions.SkipRefresh(),
		Kubeconfig:                 p.globalOptions.Kubeconfig(),
		Environment:                p.globalOptions.Environment(),
		LogLevel:                   p.globalOptions.LogLevel(),
		EnableLiveOutput:           false,
		Color:                      false,
		NoColor:                    true,
	}
	copyOptions(&globalOptions, options)
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

func (p *HelmfileLibraryExecutor) MergeEnvVars(options CommonOptionsProvider) map[string]string {
	mergedEnvVars := make(map[string]string)
	// Start with global env vars
	maps.Copy(mergedEnvVars, p.globalOptions.envVars)
	// Override with options env vars
	if options != nil {
		maps.Copy(mergedEnvVars, options.EnvVars())
	}
	return mergedEnvVars
}

// Execute runs a custom function with the given options.
func (p *HelmfileLibraryExecutor) Execute(
	ctx context.Context,
	options OptionsProvider,
	fn func(context.Context, *config.GlobalImpl) error,
) (string, string, error) {
	capture := NewOutputCapture()
	logContext := tflog.SetField(tflog.NewSubsystem(ctx, "helmfile"), "executor", "HELMFILE")
	logger := CreateCaptureLogger(logContext, capture)
	globalOptions, err := p.createGlobalOptionsFromConfigProvider(options, logger)
	if err != nil {
		return "", "", fmt.Errorf("error performing init: %w", err)
	}

	cleanup := SetEnvVars(logger.Desugar(), p.MergeEnvVars(options))
	defer cleanup()

	// Redirect also stdout
	old := os.Stdout // keep backup of the real stdout
	r, w, err := os.Pipe()
	if err != nil {
		return "", "", fmt.Errorf("error creating pipe: %w", err)
	}
	os.Stdout = w

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		_, err = io.Copy(&buf, r)
		if err != nil {
			logger.Warnf("error copying stdout: %v", err)
		}
		outC <- buf.String()
	}()

	err = fn(ctx, globalOptions)

	if closeErr := w.Close(); closeErr != nil {
		logger.Warnf("error closing pipe writer: %v", closeErr)
	}
	os.Stdout = old // restoring the real stdout
	out := <-outC
	return out, capture.String(), err
}

func (p *HelmfileLibraryExecutor) Init(
	ctx context.Context,
	options GlobalOptionsProvider,
) (string, string, error) {
	var globalOptions OptionsProvider
	if options == nil {
		globalOptions = nil
	} else {
		// Create a simple OptionsProvider that just wraps the BaseGlobalOptionsProvider
		// and our stored global options
		globalOptions = struct {
			BaseGlobalOptionsProvider
			CommonOptionsProvider
			BaseResourceOptionsProvider
		}{
			BaseGlobalOptionsProvider:   options,
			CommonOptionsProvider:       &p.globalOptions.CommonOptions,
			BaseResourceOptionsProvider: &BaseResourceOptions{},
		}
	}
	return p.Execute(ctx, globalOptions, func(ctx context.Context, gi *config.GlobalImpl) error {
		initOptions := &config.InitOptions{
			Force: true,
		}
		initImpl := config.NewInitImpl(gi, initOptions)

		helmfileApp := app.New(initImpl)
		err := helmfileApp.Init(initImpl)
		if err == nil {
			err = p.InstallAdditionalPlugins(ctx, gi.Logger())
		}
		return err
	})
}

func (p *HelmfileLibraryExecutor) Apply(
	ctx context.Context,
	options OptionsProvider,
	applyOptions *config.ApplyOptions,
) (string, string, error) {
	return p.Execute(ctx, options, func(_ context.Context, gi *config.GlobalImpl) error {
		applyOptionsImpl := config.NewApplyImpl(gi, applyOptions)

		helmfileApp := app.New(applyOptionsImpl)

		// Run apply operation
		return helmfileApp.Apply(applyOptionsImpl)
	})
}

func (p *HelmfileLibraryExecutor) Diff(
	ctx context.Context,
	options OptionsProvider,
	diffOptions *config.DiffOptions,
) (string, string, error) {
	return p.Execute(ctx, options, func(_ context.Context, gi *config.GlobalImpl) error {
		diffOptionsImpl := config.NewDiffImpl(gi, diffOptions)

		helmfileApp := app.New(diffOptionsImpl)

		// Run apply operation
		return helmfileApp.Diff(diffOptionsImpl)
	})
}

func (p *HelmfileLibraryExecutor) Template(
	ctx context.Context,
	options OptionsProvider,
	templateOptions *config.TemplateOptions,
) (string, string, error) {
	return p.Execute(ctx, options, func(_ context.Context, gi *config.GlobalImpl) error {
		templateOptionsImpl := config.NewTemplateImpl(gi, templateOptions)
		helmfileApp := app.New(templateOptionsImpl)
		// Run apply operation
		return helmfileApp.Template(templateOptionsImpl)
	})
}

func (p *HelmfileLibraryExecutor) Destroy(
	ctx context.Context,
	options OptionsProvider,
	destroyOptions *config.DestroyOptions,
) (string, string, error) {
	return p.Execute(ctx, options, func(_ context.Context, gi *config.GlobalImpl) error {
		destroyOptionsImpl := config.NewDestroyImpl(gi, destroyOptions)
		helmfileApp := app.New(destroyOptionsImpl)
		// Run apply operation
		return helmfileApp.Destroy(destroyOptionsImpl)
	})
}

func (p *HelmfileLibraryExecutor) Build(
	ctx context.Context,
	options OptionsProvider,
	embedValues bool,
) (string, string, error) {
	return p.Execute(ctx, options, func(_ context.Context, gi *config.GlobalImpl) error {
		buildImpl := config.NewBuildImpl(gi, &config.BuildOptions{
			EmbedValues: embedValues,
		})
		helmfileApp := app.New(buildImpl)
		return helmfileApp.PrintState(buildImpl)
	})
}

func (p *HelmfileLibraryExecutor) List(
	ctx context.Context,
	options OptionsProvider,
	skipCharts bool,
) (string, string, error) {
	return p.Execute(ctx, options, func(_ context.Context, gi *config.GlobalImpl) error {
		listImpl := config.NewListImpl(gi, &config.ListOptions{
			SkipCharts: skipCharts,
			Output:     "json",
		})
		helmfileApp := app.New(listImpl)
		return helmfileApp.ListReleases(listImpl)
	})
}

func (p *HelmfileLibraryExecutor) Version(_ context.Context) (string, string, error) {
	return version.Version(), "", nil
}
