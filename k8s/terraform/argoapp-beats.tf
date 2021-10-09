resource kubernetes_namespace beats {
  metadata {
    name = "elastic-beats"
  }
}

resource null_resource beats {
  triggers = {
    manifest = data.template_file.beats.rendered
  }

  provisioner local-exec {
    command = self.triggers.manifest
  }
}

data template_file beats {
  template = <<-EOT
    kubectl \
      --context docker-desktop \
      apply --validate=true \
            --wait=true \
            -f - <<EOF
    ---
    apiVersion: argoproj.io/v1alpha1
    kind: Application
    metadata:
      name: elastic-beats
      namespace: ${helm_release.argo.namespace}
      labels:
        argo.${var.business_domain}/category: observer
        argo.${var.business_domain}/organization: system
    spec:
      project: default
      source:
        repoURL: ${var.argo_apps_repository}
        targetRevision: ${var.argo_apps_revision}
        path: k8s/helm/observer/elastic-beats/
        helm:
          valueFiles:
            - values.yaml
          version: v3
      destination:
        server: https://kubernetes.default.svc
        namespace: ${kubernetes_namespace.beats.metadata[0].name}
      syncPolicy:
        automated:
          prune: true
          selfHeal: true
        syncOptions:
          - Validate=true
    EOF
  EOT
}