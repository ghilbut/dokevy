variable k8s_pv_root {
  type = string
  default = "~/work/workdata/docker-desktop/ghilbut/k8s-pv"
}

variable domain_root {
  type = string
  default = "kubelik.io"
}


################################################################
##
##  operator
##

##--------------------------------------------------------------
##  argo project

variable argo_admin_password {
  type = string
}

variable argo_apps_repository {
  type = string
  default = "git@github.com:ghilbut/docker-desktop.git"
}

variable argo_apps_revision {
  type = string
  default = "main"
}

variable github_ssh_privatekey_path {
  type = string
  default = "~/.ssh/github_rsa"
}

variable github_personal_token {
  type = string
}


################################################################
##
##  network
##

##--------------------------------------------------------------
##  cert-manager

##--------------------------------------------------------------
##  ingress-nginx


################################################################
##
##  data
##

##--------------------------------------------------------------
##  cassandra

variable cassandra_data_size {
  type = string
  default = "8Gi"
}

variable cassandra_persistence_enabled {
  type = bool
  default = false
}

variable cassandra_replica_count {
  type = number
  default = 1
}

##--------------------------------------------------------------
##  elasticssearch

variable elasticsearch_data_size {
  type = string
  default = "16Gi"
}

variable elasticsearch_persistence_enabled {
  type = bool
  default = false
}

variable elasticsearch_replica_count {
  type = number
  default = 1
}

##--------------------------------------------------------------
##  kafka

variable kafka_data_size {
  type = string
  default = "4Gi"
}

variable kafka_persistence_enabled {
  type = bool
  default = false
}

variable kafka_replica_count {
  type = number
  default = 1
}

variable zookeeper_data_size {
  type = string
  default = "1Gi"
}

variable zookeeper_log_size {
  type = string
  default = "4Gi"
}

variable zookeeper_persistence_enabled {
  type = bool
  default = false
}

variable zookeeper_replica_count {
  type = number
  default = 1
}

##--------------------------------------------------------------
##  mariadb

variable mariadb_data_size {
  type = string
  default = "8Gi"
}

variable mariadb_persistence_enabled {
  type = bool
  default = false
}
