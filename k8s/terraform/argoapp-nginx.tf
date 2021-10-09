resource kubernetes_namespace nginx {
  metadata {
    name = "ingress-nginx"
  }
}

resource null_resource nginx {
  triggers = {
    manifest = data.template_file.nginx.rendered
  }

  provisioner local-exec {
    command = self.triggers.manifest
  }
}

data template_file nginx {
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
      name: ingress-nginx
      namespace: ${helm_release.argo.namespace}
      labels:
        argo.local.in/category: network
        argo.local.in/organization: system
    spec:
      project: default
      source:
        ## https://github.com/kubernetes/ingress-nginx
        repoURL: https://kubernetes.github.io/ingress-nginx/
        chart: ingress-nginx
        targetRevision: 4.0.1
        helm:
          valueFiles:
            - values.yaml
          version: v3
      destination:
        server: https://kubernetes.default.svc
        namespace: ${kubernetes_namespace.nginx.metadata[0].name}
      syncPolicy:
        automated:
          prune: true
          selfHeal: true
        syncOptions:
          - Validate=true
    EOF
  EOT
}