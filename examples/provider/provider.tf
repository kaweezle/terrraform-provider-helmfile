# Basic provider configuration
terraform {
  required_providers {
    helmfile = {
      source = "kaweezle/helmfile"
    }
  }
}

provider "helmfile" {
  # Whether to run 'helmfile init' before applying any helmfile operations
  perform_init = true

  # Optional: Path to the helm binary
  # helm_binary_path = "/usr/local/bin/helm"

  # Optional: Path to the kustomize binary
  # kustomize_binary_path = "/usr/local/bin/kustomize"

  # Optional: Default arguments to pass to every helmfile command
  # default_args = ["--debug"]

  # Optional: Strip secret values from error messages (default: true)
  # strip_args_values_on_exit_error = true

  # Optional: Do not force helm repos to update
  # disable_force_update = false

  # Optional: Enforce helm plugin verification when installing missing plugins
  # enforce_plugin_verification = false

  # Optional: Allow using plain HTTP when pulling charts from OCI registries
  # helm_oci_plain_http = false

  # Optional: Skip running "helm repo update" and "helm dependency build"
  # skip_deps = false

  # Optional: Skip running 'helmfile repos' before applying operations
  # skip_refresh = false

  # Optional: Additional helm plugins to install
  # additional_plugins = [
  #   {
  #     name    = "x"
  #     version = "0.8.0"
  #     repo    = "https://github.com/mumoshu/helm-x"
  #   },
  # ]
}
