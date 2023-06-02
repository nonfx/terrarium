

module "container-insights" {
  source       = "Young-ook/eks/aws//modules/container-insights"
  version      = "1.4.13"
  cluster_name = aws_eks_cluster.demo.name
  oidc         = zipmap(
    ["url", "arn"],
    [local.oidc["url"], local.oidc["arn"]]
  )
  tags         = { env = "demo" }
  depends_on = [
    aws_eks_node_group.demo,
  ]
}

locals {
  oidc = {
    arn = aws_iam_openid_connect_provider.cluster.arn
    url = replace(aws_iam_openid_connect_provider.cluster.url, "https://", "")
  }
}

data "aws_eks_cluster_auth" "cluster_auth" {
  name = aws_eks_cluster.demo.name
  depends_on = [null_resource.wait_for_cluster]
}

provider "helm" {
  # Configuration options
  kubernetes {
    host                   = aws_eks_cluster.demo.endpoint
    cluster_ca_certificate = base64decode(aws_eks_cluster.demo.certificate_authority.0.data)
    token                  = data.aws_eks_cluster_auth.cluster_auth.token
  }
}

resource helm_release aws_alb_controller {
  name       = "aws-load-balancer-controller"

  repository = "https://aws.github.io/eks-charts"
  chart      = "aws-load-balancer-controller"
  namespace  = "kube-system"

  set {
    name  = "clusterName"
    value = aws_eks_cluster.demo.name
  }

  depends_on = [
    aws_eks_node_group.demo,
  ]
}

# resource "helm_release" "containerinsights" {
#   count            = var.enabled ? 1 : 0
#   name             = lookup(var.helm, "name", "eks-cw")
#   chart            = lookup(var.helm, "chart", "container-insights")
#   version          = lookup(var.helm, "version", null)
#   repository       = lookup(var.helm, "repository", join("/", [path.module, "charts"]))
#   namespace        = local.namespace
#   create_namespace = true
#   cleanup_on_fail  = lookup(var.helm, "cleanup_on_fail", true)

#   dynamic "set" {
#     for_each = {
#       "cluster.name"                                              = var.cluster_name
#       "cluster.region"                                            = data.aws_region.current.0.name
#       "serviceAccount.name"                                       = local.serviceaccount
#       "serviceAccount.annotations.eks\\.amazonaws\\.com/role-arn" = module.irsa[0].arn[0]
#     }
#     content {
#       name  = set.key
#       value = set.value
#     }
#   }