locals {
  workers_role_arns = [aws_iam_role.cluster.arn, aws_iam_role.node.arn]

  # Add worker nodes role ARNs (could be from many un-managed worker groups) to the ConfigMap
  # Note that we don't need to do this for managed Node Groups since EKS adds their roles to the ConfigMap automatically
  map_worker_roles = [
    for role_arn in local.workers_role_arns : {
      rolearn : role_arn
      username : "${role_arn}-user"
      groups : [
        "system:bootstrappers",
        "system:masters"
      ]
    }
  ]
}


provider "kubernetes" {
    host                   = data.null_data_source.cluster.outputs["cluster_endpoint"]
    cluster_ca_certificate = base64decode(aws_eks_cluster.demo.certificate_authority.0.data)
    token                  = data.aws_eks_cluster_auth.cluster_auth.token
}

data "null_data_source" "cluster" {
   inputs = {
     cluster_endpoint = aws_eks_cluster.demo.endpoint
   }
}


resource "kubernetes_config_map" "aws_auth" {
  depends_on = [aws_eks_cluster.demo]

  metadata {
    name      = "aws-auth"
    namespace = "kube-system"
  }

  data = {
    mapRoles    = yamlencode(distinct(concat(local.map_worker_roles)))
  }
}

resource "null_resource" "wait_for_cluster" {
  depends_on = [
    aws_eks_cluster.demo,
    aws_security_group.demo-cluster
  ]

  provisioner "local-exec" {
    command     = var.wait_for_cluster_cmd
    interpreter = var.wait_for_cluster_interpreter
    environment = {
      ENDPOINT = aws_eks_cluster.demo.endpoint
    }
  }
}