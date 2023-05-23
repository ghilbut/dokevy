api_addr      = "http://127.0.0.1:8200"
cluster_addr  = "http://127.0.0.1:8201"
disable_clustering = true
disable_mlock = true
  
storage file {
  path    = "/vault/file"
}

log_level = "info"

listener tcp {
  address     = "0.0.0.0:8200"
  tls_disable = true
}

ui = true
