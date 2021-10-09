resource kubernetes_namespace kafka {
  metadata {
    name = "kafka"
  }
}

resource null_resource kafka {
  depends_on = [
    kubernetes_persistent_volume_claim.kafka,
    kubernetes_persistent_volume_claim.zookeeper,
    kubernetes_persistent_volume_claim.zookeeper_log,
  ]

  triggers = {
    manifest = data.template_file.kafka.rendered
  }

  provisioner local-exec {
    command = self.triggers.manifest
  }
}

data template_file kafka {
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
      name: kafka
      namespace: ${helm_release.argo.namespace}
      labels:
        argo.${var.business_domain}/category: data
        argo.${var.business_domain}/organization: platform
    spec:
      project: default
      source:
        repoURL: ${var.argo_apps_repository}
        targetRevision: ${var.argo_apps_revision}
        path: k8s/helm/data/kafka/
        helm:
          values: |
            kafka:
              cp-zookeeper:
                persistence:
                  enabled: ${var.kafka_persistence_enabled}
                  dataDirSize: ${var.zookeeper_data_size}
                  dataLogDirSize: ${var.zookeeper_log_size}
              cp-kafka:
                persistence:
                  enabled: ${var.kafka_persistence_enabled}
                  size: ${var.kafka_data_size}
            ingress:
              hosts:
                - kafka.${var.inhouse_domain}
            ---
          valueFiles:
            - values.yaml
          version: v3
      destination:
        server: https://kubernetes.default.svc
        namespace: ${kubernetes_namespace.kafka.metadata.0.name}
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
##  zookeeper pvc for data

resource kubernetes_persistent_volume_claim zookeeper {
  count = length(kubernetes_persistent_volume.zookeeper)

  metadata {
    # name: volumeclaimtemplates-name-statefulset-name-replica-index
    name = "datadir-kafka-cp-zookeeper-${count.index}"
    namespace = kubernetes_namespace.kafka.metadata.0.name
  }
  spec {
    access_modes = ["ReadWriteOnce"]
    resources {
      requests = {
        storage = var.zookeeper_data_size
      }
    }
    volume_name = kubernetes_persistent_volume.zookeeper.0.metadata.0.name
    storage_class_name = "local-storage"
  }
  wait_until_bound = true
}

resource kubernetes_persistent_volume zookeeper {
  count = var.kafka_persistence_enabled ? 1 : 0

  metadata {
    name = "zookeeper-${count.index}"
  }
  spec {
    access_modes = ["ReadWriteOnce"]
    capacity = {
      storage = var.zookeeper_data_size
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
        path = "${var.k8s_pv_root}/zookeeper-${count.index}/data"
      }
    }
    storage_class_name = "local-storage"
    volume_mode = "Filesystem"
  }
}

##--------------------------------------------------------------
##  zookeeper pvc for log

resource kubernetes_persistent_volume_claim zookeeper_log {
  count = length(kubernetes_persistent_volume.zookeeper_log)

  metadata {
    # name: volumeclaimtemplates-name-statefulset-name-replica-index
    name = "datalogdir-kafka-cp-zookeeper-${count.index}"
    namespace = kubernetes_namespace.kafka.metadata.0.name
  }
  spec {
    access_modes = ["ReadWriteOnce"]
    resources {
      requests = {
        storage = var.zookeeper_log_size
      }
    }
    volume_name = kubernetes_persistent_volume.zookeeper_log.0.metadata.0.name
    storage_class_name = "local-storage"
  }
  wait_until_bound = true
}

resource kubernetes_persistent_volume zookeeper_log {
  count = var.kafka_persistence_enabled ? 1 : 0

  metadata {
    name = "zookeeper-log-${count.index}"
  }
  spec {
    access_modes = ["ReadWriteOnce"]
    capacity = {
      storage = var.zookeeper_log_size
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
        path = "${var.k8s_pv_root}/zookeeper-${count.index}/log"
      }
    }
    storage_class_name = "local-storage"
    volume_mode = "Filesystem"
  }
}

##--------------------------------------------------------------
##  kafka pvc

resource kubernetes_persistent_volume_claim kafka {
  count = length(kubernetes_persistent_volume.kafka)

  metadata {
    # name: volumeclaimtemplates-name-statefulset-name-replica-index
    name = "datadir-0-kafka-cp-kafka-${count.index}"
    namespace = kubernetes_namespace.kafka.metadata.0.name
  }
  spec {
    access_modes = ["ReadWriteOnce"]
    resources {
      requests = {
        storage = var.kafka_data_size
      }
    }
    volume_name = kubernetes_persistent_volume.kafka.0.metadata.0.name
    storage_class_name = "local-storage"
  }
  wait_until_bound = true
}

resource kubernetes_persistent_volume kafka {
  count = var.kafka_persistence_enabled ? 1 : 0

  metadata {
    name = "kafka-${count.index}-0"
  }
  spec {
    access_modes = ["ReadWriteOnce"]
    capacity = {
      storage = var.kafka_data_size
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
        path = "${var.k8s_pv_root}/kafka-${count.index}/data-0"
      }
    }
    storage_class_name = "local-storage"
    volume_mode = "Filesystem"
  }
}