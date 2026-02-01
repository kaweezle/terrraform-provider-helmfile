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

// NewHelmfileLibraryExecutor creates a new HelmfileLibraryExecutor instance.
//
// It initializes the executor with global options and a list of additional Helm plugins
// that should be installed and verified during helmfile operations.
//
// Parameters:
//   - options: Global configuration options for helmfile operations
//   - additionalPlugins: List of Helm plugins to install and manage
//
// Returns:
//   - *HelmfileLibraryExecutor: A new executor instance ready to perform helmfile operations
func NewHelmfileLibraryExecutor(
	options *GlobalOptions,
	additionalPlugins []HelmPlugin,
) *HelmfileLibraryExecutor {
	return &HelmfileLibraryExecutor{
		globalOptions:     options,
		additionalPlugins: additionalPlugins,
	}
}

// InstallAdditionalPlugins installs and verifies additional Helm plugins.
//
// This method checks each configured plugin's version and installs or updates
// it as necessary. It ensures all required plugins are available before helmfile
// operations are executed.
//
// Parameters:
//   - ctx: Context for cancellation and deadline control
//   - logger: Structured logger for operation tracking
//
// Returns:
//   - error: If plugin installation or verification fails
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

// copyOptions merges options from an OptionsProvider into GlobalOptions.
//
// This helper function copies configuration values from the provider interface
// into the helmfile GlobalOptions structure. It handles conditional copying,
// only overwriting values if they are set in the provider.
//
// Parameters:
//   - globalOptions: Target GlobalOptions to populate
//   - options: Source provider containing configuration values
//
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

// createGlobalOptionsFromConfigProvider creates a helmfile GlobalImpl from configuration.
//
// This method constructs a complete GlobalImpl instance by merging the executor's
// stored global options with resource-specific options. It handles logger configuration
// and state value sets.
//
// Parameters:
//   - options: Resource-specific options to merge with global settings
//   - logger: Logger instance to attach to the configuration
//
// Returns:
//   - *config.GlobalImpl: Fully configured global options for helmfile operations
//   - error: If configuration creation or validation fails
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

// MergeEnvVars merges global and resource-specific environment variables.
//
// This method combines environment variables from global configuration with
// resource-specific overrides. Resource-level variables take precedence over
// global variables when there are conflicts.
//
// Parameters:
//   - options: Provider containing resource-specific environment variables
//
// Returns:
//   - map[string]string: Merged environment variables map
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

// Execute runs a custom helmfile operation with the given options.
//
// This is the core execution method that sets up the environment, captures output,
// manages environment variables, and executes the provided function. It handles
// stdout redirection, logging configuration, and cleanup.
//
// Parameters:
//   - ctx: Context for cancellation and deadline control
//   - options: Configuration options for the operation
//   - fn: Function to execute with the configured GlobalImpl
//
// Returns:
//   - string: Captured stdout from the operation
//   - string: Captured logs from the operation
//   - error: If the operation fails
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

// Init runs helmfile init to install dependencies.
//
// This method initializes helmfile by installing chart dependencies and any
// configured additional plugins. It should be called before performing other
// helmfile operations if dependencies need to be resolved.
//
// Parameters:
//   - ctx: Context for cancellation and deadline control
//   - options: Global options provider (can be nil to use executor defaults)
//
// Returns:
//   - string: Captured stdout from the init operation
//   - string: Captured logs from the init operation
//   - error: If initialization fails
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

// Apply runs helmfile apply/sync to deploy releases.
//
// This method applies the helmfile configuration, deploying or updating Helm releases
// as specified. It is equivalent to running 'helmfile apply' on the command line.
//
// Parameters:
//   - ctx: Context for cancellation and deadline control
//   - options: Configuration options for the operation
//   - applyOptions: Specific options for the apply operation (wait, timeout, etc.)
//
// Returns:
//   - string: Captured stdout from the apply operation
//   - string: Captured logs from the apply operation
//   - error: If the apply operation fails
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

// Diff runs helmfile diff to show changes.
//
// This method compares the current state with the desired state defined in the
// helmfile configuration. It shows what changes would be applied without actually
// applying them.
//
// Parameters:
//   - ctx: Context for cancellation and deadline control
//   - options: Configuration options for the operation
//   - diffOptions: Specific options for the diff operation
//
// Returns:
//   - string: Captured stdout showing the diff output
//   - string: Captured logs from the diff operation
//   - error: If the diff operation fails
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

// Template runs helmfile template to render manifests.
//
// This method renders the Kubernetes manifests from Helm charts defined in the
// helmfile configuration without installing them. Useful for validation and
// inspection of generated resources.
//
// Parameters:
//   - ctx: Context for cancellation and deadline control
//   - options: Configuration options for the operation
//   - templateOptions: Specific options for the template operation
//
// Returns:
//   - string: Captured stdout containing rendered manifests
//   - string: Captured logs from the template operation
//   - error: If the template operation fails
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

// Destroy runs helmfile destroy to delete releases.
//
// This method removes all Helm releases defined in the helmfile configuration
// from the Kubernetes cluster. It is equivalent to running 'helmfile destroy'
// on the command line.
//
// Parameters:
//   - ctx: Context for cancellation and deadline control
//   - options: Configuration options for the operation
//   - destroyOptions: Specific options for the destroy operation (cascade, wait, etc.)
//
// Returns:
//   - string: Captured stdout from the destroy operation
//   - string: Captured logs from the destroy operation
//   - error: If the destroy operation fails
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

// Build runs helmfile build to validate configuration.
//
// This method validates the helmfile configuration and outputs the processed
// state. It can optionally embed values into the output. This is primarily
// used for validation and generating checksums for change detection.
//
// Parameters:
//   - ctx: Context for cancellation and deadline control
//   - options: Configuration options for the operation
//   - embedValues: Whether to embed values in the build output
//
// Returns:
//   - string: Captured stdout containing the build output
//   - string: Captured logs from the build operation
//   - error: If the build operation fails
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

// List runs helmfile list to list releases.
//
// This method lists all releases defined in the helmfile configuration,
// returning them in JSON format. It can optionally skip chart information
// for faster execution.
//
// Parameters:
//   - ctx: Context for cancellation and deadline control
//   - options: Configuration options for the operation
//   - skipCharts: Whether to skip chart information in the output
//
// Returns:
//   - string: Captured stdout containing JSON list of releases
//   - string: Captured logs from the list operation
//   - error: If the list operation fails
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

// Version returns the helmfile version.
//
// This method returns the version string of the helmfile library being used.
// It does not require any configuration options.
//
// Parameters:
//   - _: Context (not used in this implementation)
//
// Returns:
//   - string: Version string
//   - string: Empty logs string
//   - error: Always nil (this operation cannot fail)
func (p *HelmfileLibraryExecutor) Version(_ context.Context) (string, string, error) {
	return version.Version(), "", nil
}
