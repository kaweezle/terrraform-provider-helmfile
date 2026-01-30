terraform {
  required_providers {
    helmfile = {
      source = "kaweezle/helmfile"
    }
  }
}

provider "helmfile" {
  helm_binary_path = "/usr/bin/helm"
  perform_init     = true
  additional_plugins = [
    {
      name    = "x"
      version = "0.8.0"
      repo    = "https://github.com/mumoshu/helm-x"
    }
  ]
  log_level = "debug"

}

resource "helmfile_release" "example" {
  name         = "prometheus"
  file_or_path = "${path.module}/helmfile.yaml"
  kubeconfig   = "/root/kubeconfig"
}
