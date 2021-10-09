resource kubernetes_namespace metallb {
  metadata {
    name = "metallb-system"
  }
}

resource null_resource metallb {
  triggers = {
    manifest = data.template_file.metallb.rendered
  }

  provisioner local-exec {
    command = self.triggers.manifest
  }
}

data template_file metallb {
  template = <<-EOT
    kubectl \
      --context ${var.k8s_context.name} \
      apply --validate=true \
            --wait=true \
            -f - <<EOF
    ---
    apiVersion: argoproj.io/v1alpha1
    kind: Application
    metadata:
      name: metallb
      namespace: ${helm_release.argo.namespace}
      labels:
        argo.${var.domain_root}/category: network
        argo.${var.domain_root}/organization: system
    spec:
      project: default
      source:
        ## https://github.com/bitnami/charts
        repoURL: https://charts.bitnami.com/bitnami/
        chart: metallb
        targetRevision: 2.5.4
        helm:
          values: |
            configInline:
              address-pools:
                - name: default
                  protocol: layer2
                  addresses: 192.168.255.10-192.168.255.250
            ---
          valueFiles:
            - values.yaml
          version: v3
      destination:
        server: https://kubernetes.default.svc
        namespace: ${kubernetes_namespace.metallb.metadata[0].name}
      syncPolicy:
        automated:
          prune: true
          selfHeal: true
        syncOptions:
          - Validate=true
    EOF
  EOT
}