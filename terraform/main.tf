terraform {
  required_version = "~> 1.3.6"

  backend kubernetes {
    secret_suffix  = "ultary"
    labels         = {}
    namespace      = "default"
    config_path    = "~/.kube/config"
    config_context = "docker-desktop"
  }

  required_providers {
    postgresql = {
      source  = "cyrilgdn/postgresql"
      version = "~> 1.17.1"
    }
  }
}

provider kubernetes {
  experiments {
    manifest_resource = true
  }

  config_path    = var.k8s_context.path
  config_context = var.k8s_context.name
}

provider helm {
  kubernetes {
    config_path    = var.k8s_context.path
    config_context = var.k8s_context.name
  }
}
