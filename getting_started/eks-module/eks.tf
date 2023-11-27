module "tr_component_eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "~> 19.0"
  for_each = local.tr_component_eks
  cluster_name    = each.value.cluster_name
  cluster_version = each.value.cluster_version

  cluster_endpoint_public_access  = each.value.cluster_endpoint_public_access

  cluster_addons = var.cluster_addons

  vpc_id                   = var.vpc_id
  subnet_ids               = var.subnet_ids
  control_plane_subnet_ids = var.control_plane_subnet_ids

  # Self Managed Node Group(s)
  self_managed_node_group_defaults = var.self_managed_node_group_defaults

  self_managed_node_groups = var.self_managed_node_groups

  # EKS Managed Node Group(s)
  eks_managed_node_group_defaults = var.eks_managed_node_group_defaults

  eks_managed_node_groups = var.eks_managed_node_groups

  # Fargate Profile(s)
  fargate_profiles = var.fargate_profiles

  # aws-auth configmap
  manage_aws_auth_configmap = var.manage_aws_auth_configmap

  aws_auth_roles = var.aws_auth_roles

  aws_auth_users = var.aws_auth_users

  aws_auth_accounts = var.aws_auth_accounts

  tags = var.tags
}