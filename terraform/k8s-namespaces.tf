################################################################
##
##  Kubernetes namespaces
##

data kubernetes_namespace kube_system {
  metadata {
    name = "kube-system"
  }
}

locals {
  namespace_system_labels = {
    metallb-system = {
      "pod-security.kubernetes.io/enforce" = "privileged"
      "pod-security.kubernetes.io/audit"   = "privileged"
      "pod-security.kubernetes.io/warn"    = "privileged"
    }
  }
}

resource kubernetes_namespace system {
  for_each = toset([
    "argo",
    "cert-manager",
    "external-dns",
    "ingress-nginx",
    "mariadb",
    "metallb-system",
    "postgres",
    "prometheus",
  ])

  metadata {
    name = each.key
    labels = lookup(local.namespace_system_labels, each.key, {})
  }
}
