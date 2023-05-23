variable k8s_context {
  type    = object({ path = string, name = string })
  default = {
    path = "~/.kube/config",
    name = "docker-desktop"
  }
}

variable github_password {
  type = string
}

variable github_username {
  type = string
}

variable aws_letsencrypt_access_key {
  type = string
}

variable aws_letsencrypt_secret_key {
  type = string
}

variable pv_root {
  type = string
  default = "/Users/jhkim/work/workdata/docker-desktop/k8s-pv"
}

variable mariadb_admin_password {
  type = string
}

variable postgres_admin_password {
  type = string
}

variable postgres_harbor_password {
  type = string
}
