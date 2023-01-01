################################################################
##
##  Argo CD Application - MariaDB
##

resource kubernetes_secret mariadb_envs {
  metadata {
    name = "mariadb-env-vars"
    namespace = kubernetes_namespace.system["mariadb"].metadata.0.name
  }
  data = {
    #MARIADB_ROOT_PASSWORD = random_password.mariadb_admin.result
    MARIADB_ROOT_PASSWORD = var.mariadb_admin_password
  }
}

resource kubernetes_persistent_volume_claim mariadb {
  metadata {
    name = "mariadb"
    namespace = kubernetes_namespace.system["mariadb"].metadata.0.name
  }
  spec {
    access_modes = ["ReadWriteOnce"]
    resources {
      requests = {
        storage = kubernetes_persistent_volume.mariadb.spec.0.capacity.storage
      }
    }
    storage_class_name = kubernetes_persistent_volume.mariadb.spec.0.storage_class_name
    volume_name = kubernetes_persistent_volume.mariadb.metadata.0.name
  }
}

resource kubernetes_persistent_volume mariadb {
  metadata {
    name = "mariadb"
  }
  spec {
    access_modes = ["ReadWriteOnce"]
    capacity     = {
      storage = "10Gi"
    }
    node_affinity {
      required {
        node_selector_term {
          match_expressions {
            key      = "kubernetes.io/hostname"
            operator = "In"
            values   = ["docker-desktop"]
          }
        }
      }
    }
    persistent_volume_reclaim_policy = "Retain"
    persistent_volume_source {
      local {
        path = join("/", [trimsuffix(var.pv_root, "/"), "mariadb"])
      }
    }
    storage_class_name = "hostpath"
  }
}

# resource random_password mariadb_admin {
#   length           = 12
#   lower            = true
#   min_lower        = 2
#   min_numeric      = 2
#   min_special      = 2
#   min_upper        = 2
#   numeric          = true
#   special          = true
#   override_special = "!@#%^&*()-"
#   upper            = true
# }
#
# output mariadb_admin_password {
#   value = random_password.mariadb_admin.result
#   sensitive = true
# }
