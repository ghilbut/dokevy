################################################################
##
##  Argo CD Application - PostgreSQL
##

resource kubernetes_secret postgres_configs {
  metadata {
    name = "postgres-config-files"
    namespace = kubernetes_namespace.system["postgres"].metadata.0.name
  }
  data = {
    "pg_hba.conf"     = file("./files/postgres-14.5/pg_hba.conf")
    "postgresql.conf" = file("./files/postgres-14.5/postgresql.conf")
  }
}

resource kubernetes_secret postgres_envs {
  metadata {
    name = "postgres-env-vars"
    namespace = kubernetes_namespace.system["postgres"].metadata.0.name
  }
  data = {
    #POSTGRES_PASSWORD = random_password.postgres_admin.result
    POSTGRES_PASSWORD = var.postgres_admin_password
  }
}

resource kubernetes_persistent_volume_claim postgres {
  metadata {
    name = "postgres"
    namespace = kubernetes_namespace.system["postgres"].metadata.0.name
  }
  spec {
    access_modes = ["ReadWriteOnce"]
    resources {
      requests = {
        storage = kubernetes_persistent_volume.postgres.spec.0.capacity.storage
      }
    }
    storage_class_name = kubernetes_persistent_volume.postgres.spec.0.storage_class_name
    volume_name = kubernetes_persistent_volume.postgres.metadata.0.name
  }
}

resource kubernetes_persistent_volume postgres {
  metadata {
    name = "postgres"
  }
  spec {
    access_modes = ["ReadWriteOnce"]
    capacity     = {
      storage = "1Gi"
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
        path = join("/", [trimsuffix(var.pv_root, "/"), "postgres"])
      }
    }
    storage_class_name = "hostpath"
  }
}

# resource random_password postgres_admin {
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
# output postgres_admin_password {
#   value = random_password.postgres_admin.result
#   sensitive = true
# }

provider postgresql {
  host             = "127.0.0.1"
  port             = 5432
  username         = "postgres"
  password         = var.postgres_admin_password
  superuser        = true
  sslmode          = "disable"
  connect_timeout  = 10
  expected_version = "15.1"
}
