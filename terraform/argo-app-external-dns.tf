################################################################
##
##  Argo CD Application - External-DNS
##

resource kubernetes_secret external_dns {
  metadata {
    name = "external-dns-files"
    namespace = kubernetes_namespace.system["external-dns"].metadata.0.name
  }
  data = {
    credentials = <<-EOF
      [default]
      aws_access_key_id=${var.aws_letsencrypt_access_key}
      aws_secret_access_key=${var.aws_letsencrypt_secret_key}
      EOF
  }
}
