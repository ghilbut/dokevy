terraform {
  required_version = "~> 1.0.0"

  backend kubernetes {
    secret_suffix = "ghilbut"
    labels = {}
    namespace = "default"
    config_context = "docker-desktop"
  }
}


provider kubernetes {
  config_path    = var.k8s_context.path
  config_context = var.k8s_context.name
}

provider helm {
  kubernetes {
    config_path    = var.k8s_context.path
    config_context = var.k8s_context.name
  }
}
