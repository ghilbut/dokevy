resource kubernetes_namespace istio {
  metadata {
    name = "istio-system"
  }
}

resource null_resource istio {
  triggers = {
    manifest = data.template_file.istio.rendered
  }

  provisioner local-exec {
    command = self.triggers.manifest
  }
}

data template_file istio {
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
      name: istio
      namespace: ${helm_release.argo.namespace}
      labels:
        argo.${var.business_domain}/category: network
        argo.${var.business_domain}/organization: system
    spec:
      project: default
      source:
        repoURL: ${var.argo_apps_repository}
        targetRevision: ${var.argo_apps_revision}
        path: k8s/helm/network/istio/
        helm:
          valueFiles:
            - values.yaml
          version: v3
      destination:
        server: https://kubernetes.default.svc
        namespace: ${kubernetes_namespace.istio.metadata[0].name}
      syncPolicy:
        automated:
          prune: true
          selfHeal: true
        syncOptions:
          - Validate=true
      ignoreDifferences:
        - group: admissionregistration.k8s.io
          kind: MutatingWebhookConfiguration
          name: istio-sidecar-injector
          #namespace: ${kubernetes_namespace.istio.metadata[0].name}
          jsonPointers:
            - /webhooks/0/clientConfig/caBundle
            - /webhooks/1/clientConfig/caBundle
            - /webhooks/2/clientConfig/caBundle
            - /webhooks/3/clientConfig/caBundle
        - group: admissionregistration.k8s.io
          kind: ValidatingWebhookConfiguration
          name: istio-validator-istio-system
          #namespace: ${kubernetes_namespace.istio.metadata[0].name}
          jsonPointers:
            - /webhooks/0/clientConfig/caBundle
            - /webhooks/0/failurePolicy
            - /webhooks/1/clientConfig/caBundle
            - /webhooks/1/failurePolicy
    EOF
  EOT
}