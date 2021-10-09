resource kubernetes_namespace mariadb {
  metadata {
    name = "mariadb"
  }
}

resource null_resource mariadb {
  triggers = {
    manifest = data.template_file.mariadb.rendered
  }

  provisioner local-exec {
    command = self.triggers.manifest
  }
}

data template_file mariadb {
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
      name: mariadb
      namespace: ${helm_release.argo.namespace}
      labels:
        argo.${var.business_domain}/category: data
        argo.${var.business_domain}/organization: platform
    spec:
      project: default
      source:
        ## https://github.com/bitnami/charts/tree/master/bitnami/mariadb
        repoURL: https://charts.bitnami.com/bitnami/
        chart: mariadb
        targetRevision: 9.6.0
        helm:
          values: |
            auth:
              rootPassword: rootpw
              database: ledger
              username: user
              password: userpw
              replicationPassword: replicatorpw
            primary:
              extraFlags: "--max_connections=5"
              extraEnvVars:
                - name: MARIADB_CHARACTER_SET
                  value: utf8mb4
                - name: MARIADB_COLLATE
                  value: utf8mb4_unicode_ci
                - name: TZ
                  value: Asia/Seoul
              persistence:
                enabled: ${var.mariadb_persistence_enabled}
                subPath: ${var.k8s_pv_root}/mariadb/primary
                size: ${var.mariadb_data_size}
              service:
                type: LoadBalancer
            secondary:
              replicaCount: 0
              extraFlags: "--max_connections=5"
              extraEnvVars:
                - name: MARIADB_CHARACTER_SET
                  value: utf8mb4
                - name: MARIADB_COLLATE
                  value: utf8mb4_unicode_ci
                - name: TZ
                  value: Asia/Seoul
              persistence:
                enabled: ${var.mariadb_persistence_enabled}
                subPath: ${var.k8s_pv_root}/mariadb/secondary
                size: ${var.mariadb_data_size}
              service:
                type: LoadBalancer
                port: 13306
            ---
          valueFiles:
            - values.yaml
          version: v3
      destination:
        server: https://kubernetes.default.svc
        namespace: ${kubernetes_namespace.mariadb.metadata[0].name}
      syncPolicy:
        automated:
          prune: true
          selfHeal: true
        syncOptions:
          - Validate=true
    EOF
  EOT
}

resource null_resource mariadb_pvs {
  count = var.mariadb_persistence_enabled ? 1 : 0

  depends_on = [
    kubernetes_namespace.mariadb,
  ]

  triggers = {
    primary = "${var.k8s_pv_root}/mariadb/primary"
    secondary = "${var.k8s_pv_root}/mariadb/secondary"
  }

  provisioner local-exec {
    command = <<-EOC
      mkdir -p \
        ${self.triggers.primary} \
        ${self.triggers.secondary}
      EOC
  }
}