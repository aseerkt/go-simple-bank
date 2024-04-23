data "http" "lb_controll_iam_policy" {
  url = "https://raw.githubusercontent.com/kubernetes-sigs/aws-load-balancer-controller/${local.albc_version}/docs/install/iam_policy.json"

  request_headers = {
    Accept = "application/json"
  }
}

resource "aws_iam_policy" "albc" {
  name   = "${local.name}AWSLoadBalancerControllerIAMPolicy"
  policy = data.http.lb_controll_iam_policy.response_body
}


data "aws_iam_policy_document" "eks_lb_trust" {
  version = "2012-10-17"
  statement {
    actions = ["sts:AssumeRoleWithWebIdentity"]
    effect  = "Allow"

    condition {
      test     = "StringEquals"
      variable = "${local.cluster_oidc_provder_url}:sub"
      values   = ["system:serviceaccount:kube-system:aws-load-balancer-controller"]
    }

    condition {
      test     = "StringEquals"
      variable = "${local.cluster_oidc_provder_url}:aud"
      values   = ["sts.amazonaws.com"]
    }

    principals {
      identifiers = [local.cluster_oidc_provder_arn]
      type        = "Federated"
    }
  }
}

resource "aws_iam_role" "eks_albc" {
  name               = "${local.name}AmazonEKSLoadBalancerControllerRole"
  assume_role_policy = data.aws_iam_policy_document.eks_lb_trust.json
}

resource "aws_iam_role_policy_attachment" "lbc_iam_policy" {
  policy_arn = aws_iam_policy.albc.arn
  role       = aws_iam_role.eks_albc.name
}

resource "kubernetes_service_account" "lbc_sa" {
  metadata {
    name      = local.albc_sa_name
    namespace = "kube-system"
    annotations = {
      "eks.amazonaws.com/role-arn" = aws_iam_role.eks_albc.arn
    }
    labels = {
      "app.kubernetes.io/component" = "controller"
      "app.kubernetes.io/name"      = local.albc_sa_name
    }
  }
}

resource "helm_release" "albc" {
  name       = "aws-load-balancer-controller"
  repository = "https://aws.github.io/eks-charts"
  chart      = "aws-load-balancer-controller"

  namespace = "kube-system"

  dynamic "set" {
    for_each = {
      "clusterName"           = local.name
      "serviceAccount.create" = false
      "serviceAccount.name"   = local.albc_sa_name
      "region"                = data.aws_region.current.name
      "vpcId"                 = local.vpc_id
    }
    content {
      name  = set.key
      value = set.value
    }
  }

  depends_on = [
    kubernetes_service_account.lbc_sa,
    aws_iam_role_policy_attachment.lbc_iam_policy
  ]
}
