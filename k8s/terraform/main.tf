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
  config_path    = "~/.kube/config"
  config_context = "docker-desktop"
}

provider helm {
  kubernetes {
    config_path    = "~/.kube/config"
    config_context = "docker-desktop"
  }
}
