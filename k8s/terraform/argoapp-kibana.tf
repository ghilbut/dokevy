resource kubernetes_namespace kibana {
  metadata {
    name = "kibana"
  }
}

resource null_resource kibana {
  triggers = {
    manifest = data.template_file.kibana.rendered
  }

  provisioner local-exec {
    command = self.triggers.manifest
  }
}

data template_file kibana {
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
      name: kibana
      namespace: ${helm_release.argo.namespace}
      labels:
        argo.${var.business_domain}/category: observer
        argo.${var.business_domain}/organization: platform
    spec:
      project: default
      source:
        ## https://github.com/elastic/helm-charts
        repoURL: https://helm.elastic.co/
        chart: kibana
        targetRevision: 7.14.0
        helm:
          values: |
            elasticsearchHosts: "http://elasticsearch.elasticsearch.svc:9200"
            podAnnotations:
              #co.elastic.logs/enabled: "true"
              co.elastic.logs/enabled: "false"
              co.elastic.logs/format: json
              co.elastic.logs/json.add_error_key: "true"
              co.elastic.logs/json.keys_under_root: "true"
              co.elastic.logs/json.message_key: message
              co.elastic.logs/fileset.stdout: access
              co.elastic.logs/fileset.stderr: error
              #co.elastic.logs/exclude_lines: "request ok"
            #kibanaConfig:
            #  kibana.yml: |
            #    xpack.security.enabled: true
            #    xpack.monitoring.enabled: true
            ingress:
              enabled: true
              annotations:
                #kubernetes.io/ingress.class: nginx
                kubernetes.io/ingress.class: istio
              hosts:
                - host: kibana.${var.inhouse_domain}
                  paths:
                    - path: /*
            ---
          valueFiles:
            - values.yaml
          version: v3
      destination:
        server: https://kubernetes.default.svc
        namespace: ${kubernetes_namespace.kibana.metadata[0].name}
      syncPolicy:
        automated:
          prune: true
          selfHeal: true
        syncOptions:
          - Validate=true
    EOF
  EOT
}