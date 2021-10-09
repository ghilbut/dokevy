resource kubernetes_namespace elasticsearch {
  metadata {
    name = "elasticsearch"
  }
}

resource null_resource elasticsearch {
  triggers = {
    manifest = data.template_file.elasticsearch.rendered
  }

  provisioner local-exec {
    command = self.triggers.manifest
  }
}

data template_file elasticsearch {
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
      name: elasticsearch
      namespace: ${helm_release.argo.namespace}
      labels:
        argo.${var.business_domain}/category: data
        argo.${var.business_domain}/organization: platform
    spec:
      project: default
      source:
        ## https://github.com/elastic/helm-charts
        repoURL: https://helm.elastic.co/
        chart: elasticsearch
        targetRevision: 7.14.0
        helm:
          values: |
            replicas: 1
            minimumMasterNodes: 1
            #esConfig:
            #  elasticsearch.yml: |
            #    xpack.security.enabled: true
            #    xpack.monitoring.enabled: true
            volumeClaimTemplate:
              accessModes: [ "ReadWriteOnce" ]
              resources:
                requests:
                  storage: ${var.elasticsearch_data_size}
            persistence:
              enabled: ${var.elasticsearch_persistence_enabled}
            service:
              type: LoadBalancer
              annotations:
                metallb.universe.tf/allow-shared-ip: docker-desktop
            clusterHealthCheckParams: "wait_for_status=yellow&timeout=1s"
            fullnameOverride: elasticsearch
            ---
          valueFiles:
            - values.yaml
          version: v3
      destination:
        server: https://kubernetes.default.svc
        namespace: ${kubernetes_namespace.elasticsearch.metadata[0].name}
      syncPolicy:
        automated:
          prune: true
          selfHeal: true
        syncOptions:
          - Validate=true
    EOF
  EOT
}


##--------------------------------------------------------------
##  elasticsearch pvc

resource kubernetes_persistent_volume_claim elasticsearch {
  count = length(kubernetes_persistent_volume.elasticsearch)

  metadata {
    # name: volumeclaimtemplates-name-statefulset-name-replica-index
    name = "elasticsearch-elasticsearch-${count.index}"
    namespace = kubernetes_namespace.elasticsearch.metadata.0.name
  }
  spec {
    access_modes = ["ReadWriteOnce"]
    resources {
      requests = {
        storage = var.elasticsearch_data_size
      }
    }
    volume_name = kubernetes_persistent_volume.elasticsearch[count.index].metadata[count.index].name
    storage_class_name = "local-storage"
  }
  wait_until_bound = true
}

resource kubernetes_persistent_volume elasticsearch {
  count = var.elasticsearch_persistence_enabled ? var.elasticsearch_replica_count : 0

  metadata {
    name = "elasticsearch-${count.index}"
  }
  spec {
    access_modes = ["ReadWriteOnce"]
    capacity = {
      storage = var.elasticsearch_data_size
    }
    node_affinity {
      required {
        node_selector_term {
          match_expressions {
            key = "kubernetes.io/hostname"
            operator = "In"
            values = ["docker-desktop"]
          }
        }
      }
    }
    persistent_volume_reclaim_policy = "Recycle"
    persistent_volume_source {
      local {
        path = "${var.k8s_pv_root}/elasticsearch-${count.index}/data"
      }
    }
    storage_class_name = "local-storage"
    volume_mode = "Filesystem"
  }
}