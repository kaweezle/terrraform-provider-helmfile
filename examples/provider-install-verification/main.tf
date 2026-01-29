terraform {
  required_providers {
    helmfile = {
      source = "kaweezle/helmfile"
    }
  }
}

provider "helmfile" {
  helm_binary_path = "/usr/bin/helm"
}

data "helmfile_coffee" "example" {
  helmfile_path = "${path.module}/example-helmfile.yaml"
}
