# Basic helmfile_release resource example
resource "helmfile_release" "example_basic" {
  name         = "my-app"
  file_or_path = "./helmfile.yaml"
}

# Comprehensive helmfile_release resource with common options
resource "helmfile_release" "example_comprehensive" {
  # Required attributes
  name         = "nginx-ingress"
  file_or_path = "./helmfile.yaml"

  # Environment and context
  environment  = "production"
  namespace    = "ingress-system"
  kube_context = "my-cluster"
  kubeconfig   = "/path/to/kubeconfig"

  # Selectors
  selector = [
    "app=nginx",
    "tier=frontend"
  ]

  # Helm execution options
  args        = "--timeout 10m"
  cascade     = "background"
  concurrency = 4

  # Diff options
  context           = 3
  detailed_exitcode = true
  diff_args         = "--color"
  include_tests     = true
  no_hooks          = false
  output            = "simple"
  suppress_secrets  = true
  suppress_diff     = false

  # Values and settings
  values = [
    "./values/production.yaml",
    "./values/overrides.yaml"
  ]

  set = [
    "replicaCount=3",
    "image.tag=v1.2.3"
  ]

  # Upgrade options
  reuse_values  = false
  reset_values  = false
  wait          = true
  wait_for_jobs = true

  # Advanced options
  skip_crds                 = false
  skip_diff_on_install      = false
  skip_schema_validation    = false
  allow_no_matching_release = false

  # Debug and logging
  debug     = false
  log_level = "info"

  # State values (for templating within helmfile)
  state_values_file = ["./state-values.yaml"]
  state_values_set = [
    "env=production",
    "region=us-east-1"
  ]
}

# Example with post-renderer
resource "helmfile_release" "example_post_renderer" {
  name         = "kustomized-app"
  file_or_path = "./helmfile.yaml"

  post_renderer = "/path/to/kustomize-wrapper.sh"
  post_renderer_args = [
    "--enable-helm",
    "--namespace=default"
  ]
}

# Example with selectors and needs
resource "helmfile_release" "example_selectors" {
  name         = "microservices"
  file_or_path = "./helmfile.yaml"

  selector = [
    "tier=backend",
    "env=staging"
  ]

  include_needs            = true
  include_transitive_needs = true
  skip_needs               = false
}

# Example with suppression options
resource "helmfile_release" "example_suppression" {
  name         = "secure-app"
  file_or_path = "./helmfile.yaml"

  suppress_secrets = true
  suppress = [
    "Secret",
    "ConfigMap"
  ]
  suppress_output_line_regex = [
    ".*password.*",
    ".*token.*"
  ]
}
