resource kubernetes_namespace argo {
  metadata {
    name = "argo"
  }
}

resource kubernetes_secret argo {
  metadata {
    name = "repo-github-ssh-auth"
    namespace = kubernetes_namespace.argo.metadata[0].name
  }

  data = {
    ssh-privatekey = file(var.github_ssh_privatekey_path)
  }

  type = "kubernetes_secret"
}

data external argo_admin {
  program = [
    "${path.module}/scripts/argo.sh",
    var.argo_admin_password,
    "2021-05-24T01:00:00Z",
  ]
}

##--------------------------------------------------------------
##  helm v3

resource helm_release argo {
  depends_on = [
    kubernetes_secret.argo,
  ]

  lifecycle {
    ignore_changes = [
      set_sensitive,
    ]
  }

  name      = "argo"
  chart     = "../helm/operator/argo/"
  dependency_update = true
  namespace = kubernetes_namespace.argo.metadata[0].name

  values = [
    "${file("../helm/operator/argo/values.yaml")}",
  ]

  set {
    name  = "cd.server.ingress.hosts[0]"
    value = "argo.${var.inhouse_domain}"
  }

  set_sensitive {
    name  = "cd.configs.secret.argocdServerAdminPassword"
    value = data.external.argo_admin.result.encpw
  }

  set {
    name  = "cd.configs.secret.argocdServerAdminPasswordMtime"
    value = data.external.argo_admin.result.mtime
  }
}

##--------------------------------------------------------------
##  argo application

resource null_resource argo {
  triggers = {
    manifest = data.template_file.argo.rendered
  }

  provisioner local-exec {
    command = self.triggers.manifest
  }
}

data template_file argo {
  template = <<-EOT
    kubectl \
      --context ${var.k8s_context} \
      apply --validate=true \
            --wait=true \
            -f - <<EOF
    ---
    apiVersion: argoproj.io/v1alpha1
    kind: Application
    metadata:
      name: argo
      namespace: ${helm_release.argo.namespace}
      labels:
        argo.${var.business_domain}/category: operator
        argo.${var.business_domain}/organization: plarform
    spec:
      project: default
      source:
        repoURL: ${var.argo_apps_repository}
        targetRevision: ${var.argo_apps_revision}
        path: k8s/helm/operator/argo/
        helm:
          values: |
            cd:
              server:
                ingress:
                  hosts:
                    - argo.${var.inhouse_domain}
                  tls:
                    - hosts:
                        - argo.${var.inhouse_domain}
                      secretName: argo-tls
            ---
          valueFiles:
            - values.yaml
          version: v3
      destination:
        server: https://kubernetes.default.svc
        namespace: ${kubernetes_namespace.argo.metadata[0].name}
      syncPolicy:
        syncOptions:
          - Validate=true
    EOF
  EOT
}