
//--------------------------------------------------------------
//  kubernetes

resource kubernetes_namespace letsencrypt {
  metadata {
    name = "cert-manager"
  }
}

resource kubernetes_secret letsencrypt {
  metadata {
    name = "aws-credential-secret"
    namespace = kubernetes_namespace.letsencrypt.metadata[0].name
  }

  data = {
    secret_key = aws_iam_access_key.letsencrypt.secret
  }
}


//--------------------------------------------------------------
//  aws

resource aws_iam_user letsencrypt {
  name = "letsencrypt"
  path = "/"

  tags = local.tags
}

resource aws_iam_access_key letsencrypt {
  user = aws_iam_user.letsencrypt.name
}

resource aws_iam_user_policy letsencrypt {
  name = "letsencrypt"
  user = aws_iam_user.letsencrypt.name

  policy = <<-EOF
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "route53:GetChange",
          "Resource": "arn:aws:route53:::change/*"
        },
        {
          "Effect": "Allow",
          "Action": [
            "route53:ChangeResourceRecordSets",
            "route53:ListResourceRecordSets"
          ],
          "Resource": "arn:aws:route53:::hostedzone/*"
        },
        {
          "Effect": "Allow",
          "Action": "route53:ListHostedZonesByName",
          "Resource": "*"
        }
      ]
    }
    EOF
}
