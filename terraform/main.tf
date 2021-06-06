terraform {
  backend local {
    path = "terraform.tfstate"
  }
}


provider kubernetes {
  config_path    = var.kubeconfig_path
  config_context = var.kubeconfig_context
}

provider helm {
  kubernetes {
    config_path    = var.kubeconfig_path
    config_context = var.kubeconfig_context
  }
}

provider aws {
  region  = "ap-northeast-2"
  profile = "ghilbut"
}


locals {
  tags = {
    managed_by   = "terraform"
    organization = "ghilbut.com"
    owner        = "ghilbut@gmail.com"
    purpose      = "k8s over docker-desktop for ghilbut.com"
  }
}
