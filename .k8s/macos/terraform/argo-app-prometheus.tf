################################################################
##
##  Argo CD Application - Prometheus
##

resource kubernetes_persistent_volume prometheus_altermanager {
  metadata {
    name = "prometheus-alertmanager"
  }
  spec {
    access_modes = ["ReadWriteOnce"]
    capacity     = {
      storage = "2Gi"
    }
    node_affinity {
      required {
        node_selector_term {
          match_expressions {
            key      = "kubernetes.io/hostname"
            operator = "In"
            values = ["docker-desktop"]
          }
        }
      }
    }
    persistent_volume_source {
      local {
        path = join("/", [trimsuffix(var.pv_root, "/"), "prometheus", "alertmanager"])
      }
    }
    storage_class_name = "hostpath"
  }
}

resource kubernetes_persistent_volume prometheus_pushgateway {
  metadata {
    name = "prometheus-pushgateway"
  }
  spec {
    access_modes = ["ReadWriteOnce"]
    capacity     = {
      storage = "2Gi"
    }
    node_affinity {
      required {
        node_selector_term {
          match_expressions {
            key      = "kubernetes.io/hostname"
            operator = "In"
            values = ["docker-desktop"]
          }
        }
      }
    }
    persistent_volume_source {
      local {
        path = join("/", [trimsuffix(var.pv_root, "/"), "prometheus", "pushgateway"])
      }
    }
    storage_class_name = "hostpath"
  }
}

resource kubernetes_persistent_volume prometheus_server {
  metadata {
    name = "prometheus-server"
  }
  spec {
    access_modes = ["ReadWriteOnce"]
    capacity     = {
      storage = "8Gi"
    }
    node_affinity {
      required {
        node_selector_term {
          match_expressions {
            key      = "kubernetes.io/hostname"
            operator = "In"
            values = ["docker-desktop"]
          }
        }
      }
    }
    persistent_volume_reclaim_policy = "Retain"
    persistent_volume_source {
      local {
        path = join("/", [trimsuffix(var.pv_root, "/"), "prometheus", "server"])
      }
    }
    storage_class_name = "hostpath"
  }
}
