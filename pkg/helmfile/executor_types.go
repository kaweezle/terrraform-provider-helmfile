// Copyright Antoine Martin 2026
// SPDX-License-Identifier: MIT

package helmfile

var (
	_ CommonOptionsProvider       = (*CommonOptions)(nil)
	_ BaseGlobalOptionsProvider   = (*BaseGlobalOptions)(nil)
	_ GlobalOptionsProvider       = (*GlobalOptions)(nil)
	_ BaseResourceOptionsProvider = (*BaseResourceOptions)(nil)
	_ ResourceOptionsProvider     = (*ResourceOptions)(nil)
	_ OptionsProvider             = (*Options)(nil)
)

// CommonOptions implements CommonOptionsProvider.
type CommonOptions struct {
	kubeconfig  string
	environment string
	envVars     map[string]string
	logLevel    string
}

// Kubeconfig returns the path to the kubeconfig file.
func (o *CommonOptions) Kubeconfig() string {
	return o.kubeconfig
}

// Environment returns the helmfile environment name.
func (o *CommonOptions) Environment() string {
	return o.environment
}

// EnvVars returns environment variables to set.
func (o *CommonOptions) EnvVars() map[string]string {
	return o.envVars
}

// LogLevel returns the log level.
func (o *CommonOptions) LogLevel() string {
	return o.logLevel
}

// CopyFrom copies common options from a provider.
func (o *CommonOptions) CopyFrom(from CommonOptionsProvider) {
	if from == nil {
		return
	}
	o.kubeconfig = from.Kubeconfig()
	o.environment = from.Environment()
	o.envVars = from.EnvVars()
	o.logLevel = from.LogLevel()
}

// NewCommonOptions creates a CommonOptions from a provider.
func NewCommonOptions(from CommonOptionsProvider) CommonOptions {
	opts := CommonOptions{}
	opts.CopyFrom(from)
	return opts
}

// WithKubeconfig sets the kubeconfig path.
func (o *CommonOptions) WithKubeconfig(kubeconfig string) *CommonOptions {
	o.kubeconfig = kubeconfig
	return o
}

// WithEnvironment sets the environment name.
func (o *CommonOptions) WithEnvironment(environment string) *CommonOptions {
	o.environment = environment
	return o
}

// WithEnvVars sets the environment variables.
func (o *CommonOptions) WithEnvVars(envVars map[string]string) *CommonOptions {
	o.envVars = envVars
	return o
}

// WithLogLevel sets the log level.
func (o *CommonOptions) WithLogLevel(logLevel string) *CommonOptions {
	o.logLevel = logLevel
	return o
}

// BaseGlobalOptions implements BaseGlobalOptionsProvider.
type BaseGlobalOptions struct {
	defaultArgs                string
	helmBinary                 string
	kustomizeBinary            string
	stripArgsValuesOnExitError bool
	disableForceUpdate         bool
	enforcePluginVerification  bool
	helmOCIPlainHTTP           bool
	skipDeps                   bool
	skipRefresh                bool
}

// DefaultArgs returns default args to pass to helmfile.
func (o *BaseGlobalOptions) DefaultArgs() string {
	return o.defaultArgs
}

// HelmBinary returns the path to the helm binary.
func (o *BaseGlobalOptions) HelmBinary() string {
	return o.helmBinary
}

// KustomizeBinary returns the path to the kustomize binary.
func (o *BaseGlobalOptions) KustomizeBinary() string {
	return o.kustomizeBinary
}

// StripArgsValuesOnExitError returns whether to strip args values on exit error.
func (o *BaseGlobalOptions) StripArgsValuesOnExitError() bool {
	return o.stripArgsValuesOnExitError
}

// DisableForceUpdate returns whether to disable force update.
func (o *BaseGlobalOptions) DisableForceUpdate() bool {
	return o.disableForceUpdate
}

// EnforcePluginVerification returns whether to enforce plugin verification.
func (o *BaseGlobalOptions) EnforcePluginVerification() bool {
	return o.enforcePluginVerification
}

// HelmOCIPlainHTTP returns whether to use plain HTTP for OCI.
func (o *BaseGlobalOptions) HelmOCIPlainHTTP() bool {
	return o.helmOCIPlainHTTP
}

// SkipDeps returns whether to skip dependency updates.
func (o *BaseGlobalOptions) SkipDeps() bool {
	return o.skipDeps
}

// SkipRefresh returns whether to skip refresh.
func (o *BaseGlobalOptions) SkipRefresh() bool {
	return o.skipRefresh
}

// CopyFrom copies base global options from a provider.
func (o *BaseGlobalOptions) CopyFrom(from BaseGlobalOptionsProvider) {
	if from == nil {
		return
	}
	o.defaultArgs = from.DefaultArgs()
	o.helmBinary = from.HelmBinary()
	o.kustomizeBinary = from.KustomizeBinary()
	o.stripArgsValuesOnExitError = from.StripArgsValuesOnExitError()
	o.disableForceUpdate = from.DisableForceUpdate()
	o.enforcePluginVerification = from.EnforcePluginVerification()
	o.helmOCIPlainHTTP = from.HelmOCIPlainHTTP()
	o.skipDeps = from.SkipDeps()
	o.skipRefresh = from.SkipRefresh()
}

// NewBaseGlobalOptions creates BaseGlobalOptions from a provider.
func NewBaseGlobalOptions(from BaseGlobalOptionsProvider) BaseGlobalOptions {
	opts := BaseGlobalOptions{}
	opts.CopyFrom(from)
	return opts
}

// WithDefaultArgs sets the default args.
func (o *BaseGlobalOptions) WithDefaultArgs(defaultArgs string) *BaseGlobalOptions {
	o.defaultArgs = defaultArgs
	return o
}

// WithHelmBinary sets the helm binary path.
func (o *BaseGlobalOptions) WithHelmBinary(helmBinary string) *BaseGlobalOptions {
	o.helmBinary = helmBinary
	return o
}

// WithKustomizeBinary sets the kustomize binary path.
func (o *BaseGlobalOptions) WithKustomizeBinary(kustomizeBinary string) *BaseGlobalOptions {
	o.kustomizeBinary = kustomizeBinary
	return o
}

// WithStripArgsValuesOnExitError sets whether to strip args values on exit error.
func (o *BaseGlobalOptions) WithStripArgsValuesOnExitError(strip bool) *BaseGlobalOptions {
	o.stripArgsValuesOnExitError = strip
	return o
}

// WithDisableForceUpdate sets whether to disable force update.
func (o *BaseGlobalOptions) WithDisableForceUpdate(disable bool) *BaseGlobalOptions {
	o.disableForceUpdate = disable
	return o
}

// WithEnforcePluginVerification sets whether to enforce plugin verification.
func (o *BaseGlobalOptions) WithEnforcePluginVerification(enforce bool) *BaseGlobalOptions {
	o.enforcePluginVerification = enforce
	return o
}

// WithHelmOCIPlainHTTP sets whether to use plain HTTP for OCI.
func (o *BaseGlobalOptions) WithHelmOCIPlainHTTP(plain bool) *BaseGlobalOptions {
	o.helmOCIPlainHTTP = plain
	return o
}

// WithSkipDeps sets whether to skip dependency updates.
func (o *BaseGlobalOptions) WithSkipDeps(skip bool) *BaseGlobalOptions {
	o.skipDeps = skip
	return o
}

// WithSkipRefresh sets whether to skip refresh.
func (o *BaseGlobalOptions) WithSkipRefresh(skip bool) *BaseGlobalOptions {
	o.skipRefresh = skip
	return o
}

// GlobalOptions implements GlobalOptionsProvider by embedding
// BaseGlobalOptions and CommonOptions.
type GlobalOptions struct {
	CommonOptions
	BaseGlobalOptions
}

// CopyFrom copies global options from a provider.
func (o *GlobalOptions) CopyFrom(from GlobalOptionsProvider) {
	if from == nil {
		return
	}
	o.BaseGlobalOptions.CopyFrom(from)
	o.CommonOptions.CopyFrom(from)
}

// NewGlobalOptions creates GlobalOptions from a provider.
func NewGlobalOptions(from GlobalOptionsProvider) GlobalOptions {
	opts := GlobalOptions{}
	opts.CopyFrom(from)
	return opts
}

// BaseResourceOptions implements BaseResourceOptionsProvider.
type BaseResourceOptions struct {
	args             string
	fileOrDir        string
	kubeContext      string
	namespace        string
	chart            string
	selectors        []string
	stateValuesSet   map[string]any
	stateValuesFiles []string
}

// Args returns additional args to pass to helmfile.
func (o *BaseResourceOptions) Args() string {
	return o.args
}

// FileOrDir returns the helmfile file or directory path.
func (o *BaseResourceOptions) FileOrDir() string {
	return o.fileOrDir
}

// KubeContext returns the kubernetes context to use.
func (o *BaseResourceOptions) KubeContext() string {
	return o.kubeContext
}

// Namespace returns the kubernetes namespace.
func (o *BaseResourceOptions) Namespace() string {
	return o.namespace
}

// Chart returns the chart name or path.
func (o *BaseResourceOptions) Chart() string {
	return o.chart
}

// Selectors returns the release selectors.
func (o *BaseResourceOptions) Selectors() []string {
	return o.selectors
}

// StateValuesSet returns state values to set.
func (o *BaseResourceOptions) StateValuesSet() map[string]any {
	return o.stateValuesSet
}

// StateValuesFiles returns state values files.
func (o *BaseResourceOptions) StateValuesFiles() []string {
	return o.stateValuesFiles
}

// CopyFrom copies base resource options from a provider.
func (o *BaseResourceOptions) CopyFrom(from BaseResourceOptionsProvider) {
	if from == nil {
		return
	}
	o.args = from.Args()
	o.fileOrDir = from.FileOrDir()
	o.kubeContext = from.KubeContext()
	o.namespace = from.Namespace()
	o.chart = from.Chart()
	o.selectors = from.Selectors()
	o.stateValuesSet = from.StateValuesSet()
	o.stateValuesFiles = from.StateValuesFiles()
}

// WithArgs sets the args.
func (o *BaseResourceOptions) WithArgs(args string) *BaseResourceOptions {
	o.args = args
	return o
}

// WithFileOrDir sets the file or directory path.
func (o *BaseResourceOptions) WithFileOrDir(fileOrDir string) *BaseResourceOptions {
	o.fileOrDir = fileOrDir
	return o
}

// WithKubeContext sets the kubernetes context.
func (o *BaseResourceOptions) WithKubeContext(kubeContext string) *BaseResourceOptions {
	o.kubeContext = kubeContext
	return o
}

// WithNamespace sets the namespace.
func (o *BaseResourceOptions) WithNamespace(namespace string) *BaseResourceOptions {
	o.namespace = namespace
	return o
}

// WithChart sets the chart name or path.
func (o *BaseResourceOptions) WithChart(chart string) *BaseResourceOptions {
	o.chart = chart
	return o
}

// WithSelectors sets the release selectors.
func (o *BaseResourceOptions) WithSelectors(selectors []string) *BaseResourceOptions {
	o.selectors = selectors
	return o
}

// WithStateValuesSet sets the state values.
func (o *BaseResourceOptions) WithStateValuesSet(
	stateValuesSet map[string]any,
) *BaseResourceOptions {
	o.stateValuesSet = stateValuesSet
	return o
}

// WithStateValuesFiles sets the state values files.
func (o *BaseResourceOptions) WithStateValuesFiles(stateValuesFiles []string) *BaseResourceOptions {
	o.stateValuesFiles = stateValuesFiles
	return o
}

// NewBaseResourceOptions creates BaseResourceOptions from a provider.
func NewBaseResourceOptions(from BaseResourceOptionsProvider) BaseResourceOptions {
	opts := BaseResourceOptions{}
	opts.CopyFrom(from)
	return opts
}

// ResourceOptions implements ResourceOptionsProvider by embedding
// BaseResourceOptions and CommonOptions.
type ResourceOptions struct {
	CommonOptions
	BaseResourceOptions
}

// CopyFrom copies resource options from a provider.
func (o *ResourceOptions) CopyFrom(from ResourceOptionsProvider) {
	if from == nil {
		return
	}
	o.BaseResourceOptions.CopyFrom(from)
	o.CommonOptions.CopyFrom(from)
}

// NewResourceOptions creates ResourceOptions from a provider.
func NewResourceOptions(from ResourceOptionsProvider) ResourceOptions {
	opts := ResourceOptions{}
	opts.CopyFrom(from)
	return opts
}

// CopyFrom copies resource options from a provider.
type Options struct {
	CommonOptions
	BaseGlobalOptions
	BaseResourceOptions
}

// CopyFrom copies options from a provider.
func (o *Options) CopyFrom(from OptionsProvider) {
	if from == nil {
		return
	}
	o.CommonOptions.CopyFrom(from)
	o.BaseGlobalOptions.CopyFrom(from)
	o.BaseResourceOptions.CopyFrom(from)
}

// NewOptions creates Options from a provider.
func NewOptions(from OptionsProvider) Options {
	opts := Options{}
	opts.CopyFrom(from)
	return opts
}
