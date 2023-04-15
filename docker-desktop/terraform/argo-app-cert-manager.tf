################################################################
##
##  Argo CD Application - Cert-Manager
##

resource kubernetes_manifest cluster_issuer {
  manifest = yamldecode(<<-EOY
    apiVersion: cert-manager.io/v1
    kind: ClusterIssuer
    metadata:
      name: letsencrypt
      annotations:
        argocd.argoproj.io/sync-options: SkipDryRunOnMissingResource=true
    spec:
      acme:
        email: ghilbut@gmail.com
        server: https://acme-v02.api.letsencrypt.org/directory
        privateKeySecretRef:
          name: letsencrypt
        solvers:
        - selector:
            dnsZones:
            - ghilbut.com
            - ghilbut.net
          dns01:
            route53:
              region: ap-northeast-2
              accessKeyID: ${var.aws_letsencrypt_access_key}
              secretAccessKeySecretRef:
                name: ${kubernetes_secret.letsencrypt.metadata.0.name}
                key:  secret-access-key
    EOY
  )
}

resource kubernetes_secret letsencrypt {
  metadata {
    name = "aws-iam-letsencrypt-credential-secret"
    namespace = kubernetes_namespace.system["cert-manager"].metadata.0.name
  }
  data = {
    secret-access-key = var.aws_letsencrypt_secret_key
  }
}
