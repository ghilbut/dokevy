################################################################
##
##  cert-manager
##

##--------------------------------------------------------------
##  kubernetes

resource kubernetes_namespace cert-manager {
  metadata {
    name = "cert-manager"
  }
}

resource kubernetes_secret letsencrypt {
  metadata {
    name = "aws-credential-secret"
    namespace = kubernetes_namespace.cert-manager.metadata[0].name
  }

  data = {
    secret-access-key = var.aws_iam_letsencrypt.secret_key
  }
}

##--------------------------------------------------------------
##  argo application

resource null_resource cert-manager {
  triggers = {
    manifest = data.template_file.cert-manager.rendered
  }

  provisioner local-exec {
    command = self.triggers.manifest
  }

  #provisioner local-exec {
  #  when    = destroy
  #  command = ""
  #}
}

data template_file cert-manager {
  template = <<-EOT
    kubectl \
        --context ${var.k8s_context} \
      apply \
        --validate=true \
        --wait=true \
        -f - <<EOF
    ---
    apiVersion: argoproj.io/v1alpha1
    kind: Application
    metadata:
      name: cert-manager
      namespace: ${kubernetes_namespace.argo.metadata[0].name}
      labels:
        argo.${var.business_domain}/category: network
        argo.${var.business_domain}/organization: plarform
    spec:
      project: default
      source:
        repoURL: ${var.argo_apps_repository}
        targetRevision: ${var.argo_apps_revision}
        path: k8s/helm/network/cert-manager
        helm:
          parameters:
          - name:  aws.region
            value: ${var.aws_region}
          - name:  aws.access_key
            value: ${var.aws_iam_letsencrypt.access_key}
          valueFiles:
          - values.yaml
          version: v3
      destination:
        server: https://kubernetes.default.svc
        namespace: ${kubernetes_namespace.cert-manager.metadata[0].name}
      syncPolicy:
        automated:
          prune: true
          selfHeal: true
        syncOptions:
        - Validate=true
    EOF
  EOT
}