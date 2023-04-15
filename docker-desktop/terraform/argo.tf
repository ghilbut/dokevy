################################################################
##
##  Argo CD
##

resource helm_release argo {
  lifecycle { ignore_changes = [values] }

  name      = "argo"
  chart     = "../helm/argo/"
  dependency_update = true
  namespace = kubernetes_namespace.system["argo"].metadata.0.name

  values = [
    "${file("../helm/argo/values.yaml")}",
  ]
}

resource kubernetes_secret argo_repository_credentials {
  metadata {
    name = "github-repo-creds-ultary-https"
    namespace = kubernetes_namespace.system["argo"].metadata.0.name
    labels = {
      "argocd.argoproj.io/secret-type" = "repo-creds"
    }
  }
  data = {
    type     = "git"
    url      = "https://github.com/ultary"
    password = var.github_password
    username = var.github_username
  }
}
