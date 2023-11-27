output "tr_component_eks_cluster_arn" {
    value = {for k, v in module.tr_component_eks : k => v.cluster_arn }
}

output "tr_component_eks_cluster_endpoint" {
    value = {for k, v in module.tr_component_eks : k => v.cluster_endpoint }
}


output "tr_component_eks_cluster_id" {
    value = {for k, v in module.tr_component_eks : k => v.cluster_id }
}


output "tr_component_eks_cluster_status" {
    value = {for k, v in module.tr_component_eks : k => v.cluster_status }
}


output "tr_component_eks_cluster_version" {
    value = {for k, v in module.tr_component_eks : k => v.cluster_version }
}

output "tr_component_eks_oidc_provider" {
    value = {for k, v in module.tr_component_eks : k => v.oidc_provider }
}


output "tr_component_eks_oidc_provider_arn" {
    value = {for k, v in module.tr_component_eks : k => v.oidc_provider_arn }
}

