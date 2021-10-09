resource kubernetes_namespace kiali {
  metadata {
    name = "kiali"
  }
}

resource null_resource kiali {
  count = var.ingress_type == "istio" ? 1 : 0

  triggers = {
    manifest = data.template_file.kiali.rendered
  }

  provisioner local-exec {
    command = self.triggers.manifest
  }
}

data template_file kiali {
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
      name: kiali
      namespace: ${helm_release.argo.namespace}
      labels:
        argo.${var.business_domain}/category: network
        argo.${var.business_domain}/organization: system
    spec:
      project: default
      source:
        repoURL: ${var.argo_apps_repository}
        targetRevision: ${var.argo_apps_revision}
        path: k8s/helm/network/kiali/
        helm:
          values: |
            server:
              server:
                istio_namespace: ${kubernetes_namespace.istio.metadata[0].name}
                web_fqdn: kiali.${var.inhouse_domain}
            ---
          valueFiles:
            - values.yaml
          version: v3
      destination:
        server: https://kubernetes.default.svc
        namespace: ${kubernetes_namespace.kiali.metadata[0].name}
      syncPolicy:
        automated:
          prune: true
          selfHeal: true
        syncOptions:
          - Validate=true
    EOF
  EOT
}