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
      --context docker-desktop \
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
        argo.local.in/category: data
        argo.local.in/organization: platform
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
                  enabled: ${var.zookeeper_persistence_enabled}
                  dataDirSize: ${var.zookeeper_data_size}
                  dataLogDirSize: ${var.zookeeper_log_size}
              cp-kafka:
                persistence:
                  enabled: ${var.kafka_persistence_enabled}
                  size: ${var.kafka_data_size}
            ingress:
              hosts:
                - kafka.${var.domain_root}
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
    volume_name = kubernetes_persistent_volume.zookeeper[count.index].metadata.0.name
    storage_class_name = "local-storage"
  }
  wait_until_bound = true
}

resource kubernetes_persistent_volume zookeeper {
  count = length(null_resource.zookeeper_path)

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
        path = null_resource.zookeeper_path[count.index].triggers.path
      }
    }
    storage_class_name = "local-storage"
    volume_mode = "Filesystem"
  }
}

resource null_resource zookeeper_path {
  count = var.zookeeper_persistence_enabled ? var.zookeeper_replica_count : 0

  depends_on = [
    kubernetes_namespace.kafka,
  ]

  triggers = {
    path = pathexpand("${var.k8s_pv_root}/zookeeper-${count.index}/data")
  }

  provisioner local-exec {
    command = "mkdir -p ${self.triggers.path}"
  }

  provisioner local-exec {
    when    = destroy
    command = "rm -rf ${self.triggers.path}"
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
    volume_name = kubernetes_persistent_volume.zookeeper_log[count.index].metadata.0.name
    storage_class_name = "local-storage"
  }
  wait_until_bound = true
}

resource kubernetes_persistent_volume zookeeper_log {
  count = length(null_resource.zookeeper_log_path)

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
        path = null_resource.zookeeper_log_path[count.index].triggers.path
      }
    }
    storage_class_name = "local-storage"
    volume_mode = "Filesystem"
  }
}

resource null_resource zookeeper_log_path {
  count = var.zookeeper_persistence_enabled ? var.zookeeper_replica_count ? 0

  depends_on = [
    kubernetes_namespace.kafka,
  ]

  triggers = {
    path = pathexpand("${var.k8s_pv_root}/zookeeper-${count.index}/log")
  }

  provisioner local-exec {
    command = "mkdir -p ${self.triggers.path}"
  }

  provisioner local-exec {
    when    = destroy
    command = "rm -rf ${self.triggers.path}"
  }
}

##--------------------------------------------------------------
##  kafka pvc

resource kubernetes_persistent_volume_claim kafka {
  count = length(kubernetes_persistent_volume.kafka)

  metadata {
    # name: volumeclaimtemplates-name-statefulset-name-replica-index
    name = "datadir-0-kafka-cp-kafka-0"
    namespace = kubernetes_namespace.kafka.metadata.0.name
  }
  spec {
    access_modes = ["ReadWriteOnce"]
    resources {
      requests = {
        storage = var.kafka_data_size
      }
    }
    volume_name = kubernetes_persistent_volume.kafka[count.index].metadata.0.name
    storage_class_name = "local-storage"
  }
  wait_until_bound = true
}

resource kubernetes_persistent_volume kafka {
  count = length(null_resource.kafka_path)

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
        path = null_resource.kafka_path[count.index].triggers.path
      }
    }
    storage_class_name = "local-storage"
    volume_mode = "Filesystem"
  }
}

resource null_resource kafka_path {
  count = var.kafka_persistence_enabled ? var.kafka_replica_count : 0

  depends_on = [
    kubernetes_namespace.kafka,
  ]

  triggers = {
    path = pathexpand("${var.k8s_pv_root}/kafka-${count.index}/data-0")
  }

  provisioner local-exec {
    command = "mkdir -p ${self.triggers.path}"
  }

  provisioner local-exec {
    when    = destroy
    command = "rm -rf ${self.triggers.path}"
  }
}