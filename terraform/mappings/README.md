<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 0.12 |
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 0.13.1 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | 3.74 |
| <a name="requirement_helm"></a> [helm](#requirement\_helm) | 2.1.2 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | 3.74 |
| <a name="provider_helm"></a> [helm](#provider\_helm) | 2.1.2 |
| <a name="provider_http"></a> [http](#provider\_http) | n/a |
| <a name="provider_kubernetes"></a> [kubernetes](#provider\_kubernetes) | n/a |
| <a name="provider_null"></a> [null](#provider\_null) | n/a |
| <a name="provider_random"></a> [random](#provider\_random) | n/a |

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_container-insights"></a> [container-insights](#module\_container-insights) | Young-ook/eks/aws//modules/container-insights | 1.4.13 |

## Resources

| Name | Type |
|------|------|
| [aws_acm_certificate.main](https://registry.terraform.io/providers/hashicorp/aws/3.74/docs/resources/acm_certificate) | resource |
| [aws_acm_certificate_validation.main](https://registry.terraform.io/providers/hashicorp/aws/3.74/docs/resources/acm_certificate_validation) | resource |
| [aws_db_instance.default](https://registry.terraform.io/providers/hashicorp/aws/3.74/docs/resources/db_instance) | resource |
| [aws_db_subnet_group.database](https://registry.terraform.io/providers/hashicorp/aws/3.74/docs/resources/db_subnet_group) | resource |
| [aws_eks_cluster.demo](https://registry.terraform.io/providers/hashicorp/aws/3.74/docs/resources/eks_cluster) | resource |
| [aws_eks_node_group.demo](https://registry.terraform.io/providers/hashicorp/aws/3.74/docs/resources/eks_node_group) | resource |
| [aws_elasticache_cluster.demo](https://registry.terraform.io/providers/hashicorp/aws/3.74/docs/resources/elasticache_cluster) | resource |
| [aws_elasticache_subnet_group.redis](https://registry.terraform.io/providers/hashicorp/aws/3.74/docs/resources/elasticache_subnet_group) | resource |
| [aws_iam_instance_profile.node](https://registry.terraform.io/providers/hashicorp/aws/3.74/docs/resources/iam_instance_profile) | resource |
| [aws_iam_openid_connect_provider.cluster](https://registry.terraform.io/providers/hashicorp/aws/3.74/docs/resources/iam_openid_connect_provider) | resource |
| [aws_iam_policy.AWSLoadBalancerControllerIAMPolicy](https://registry.terraform.io/providers/hashicorp/aws/3.74/docs/resources/iam_policy) | resource |
| [aws_iam_policy.AWSVisualEditorPolicy](https://registry.terraform.io/providers/hashicorp/aws/3.74/docs/resources/iam_policy) | resource |
| [aws_iam_role.cluster](https://registry.terraform.io/providers/hashicorp/aws/3.74/docs/resources/iam_role) | resource |
| [aws_iam_role.node](https://registry.terraform.io/providers/hashicorp/aws/3.74/docs/resources/iam_role) | resource |
| [aws_iam_role_policy_attachment.cluster-AWSVisualEditorPolicy](https://registry.terraform.io/providers/hashicorp/aws/3.74/docs/resources/iam_role_policy_attachment) | resource |
| [aws_iam_role_policy_attachment.cluster-AmazonEKSClusterPolicy](https://registry.terraform.io/providers/hashicorp/aws/3.74/docs/resources/iam_role_policy_attachment) | resource |
| [aws_iam_role_policy_attachment.cluster-AmazonEKSServicePolicy](https://registry.terraform.io/providers/hashicorp/aws/3.74/docs/resources/iam_role_policy_attachment) | resource |
| [aws_iam_role_policy_attachment.cluster-AmazonEKSVPCResourceController](https://registry.terraform.io/providers/hashicorp/aws/3.74/docs/resources/iam_role_policy_attachment) | resource |
| [aws_iam_role_policy_attachment.cluster-AmazonVPCFullAccess](https://registry.terraform.io/providers/hashicorp/aws/3.74/docs/resources/iam_role_policy_attachment) | resource |
| [aws_iam_role_policy_attachment.node-AWSLoadBalancerControllerIAMPolicy](https://registry.terraform.io/providers/hashicorp/aws/3.74/docs/resources/iam_role_policy_attachment) | resource |
| [aws_iam_role_policy_attachment.node-AmazonEC2ContainerRegistryReadOnly](https://registry.terraform.io/providers/hashicorp/aws/3.74/docs/resources/iam_role_policy_attachment) | resource |
| [aws_iam_role_policy_attachment.node-AmazonEKSWorkerNodePolicy](https://registry.terraform.io/providers/hashicorp/aws/3.74/docs/resources/iam_role_policy_attachment) | resource |
| [aws_iam_role_policy_attachment.node-AmazonEKS_CNI_Policy](https://registry.terraform.io/providers/hashicorp/aws/3.74/docs/resources/iam_role_policy_attachment) | resource |
| [aws_iam_role_policy_attachment.node-AmazonVPCFullAccess](https://registry.terraform.io/providers/hashicorp/aws/3.74/docs/resources/iam_role_policy_attachment) | resource |
| [aws_internet_gateway.demo](https://registry.terraform.io/providers/hashicorp/aws/3.74/docs/resources/internet_gateway) | resource |
| [aws_route53_record.main](https://registry.terraform.io/providers/hashicorp/aws/3.74/docs/resources/route53_record) | resource |
| [aws_route_table.demo](https://registry.terraform.io/providers/hashicorp/aws/3.74/docs/resources/route_table) | resource |
| [aws_route_table_association.demo](https://registry.terraform.io/providers/hashicorp/aws/3.74/docs/resources/route_table_association) | resource |
| [aws_security_group.dbsg](https://registry.terraform.io/providers/hashicorp/aws/3.74/docs/resources/security_group) | resource |
| [aws_security_group.demo-cluster](https://registry.terraform.io/providers/hashicorp/aws/3.74/docs/resources/security_group) | resource |
| [aws_security_group.redissg](https://registry.terraform.io/providers/hashicorp/aws/3.74/docs/resources/security_group) | resource |
| [aws_security_group_rule.demo-cluster-ingress-workstation-https](https://registry.terraform.io/providers/hashicorp/aws/3.74/docs/resources/security_group_rule) | resource |
| [aws_subnet.dbsubnet](https://registry.terraform.io/providers/hashicorp/aws/3.74/docs/resources/subnet) | resource |
| [aws_subnet.demo](https://registry.terraform.io/providers/hashicorp/aws/3.74/docs/resources/subnet) | resource |
| [aws_subnet.redissubnet](https://registry.terraform.io/providers/hashicorp/aws/3.74/docs/resources/subnet) | resource |
| [aws_vpc.demo](https://registry.terraform.io/providers/hashicorp/aws/3.74/docs/resources/vpc) | resource |
| [helm_release.aws_alb_controller](https://registry.terraform.io/providers/hashicorp/helm/2.1.2/docs/resources/release) | resource |
| [kubernetes_config_map.aws_auth](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/config_map) | resource |
| [null_resource.delete_ingress](https://registry.terraform.io/providers/hashicorp/null/latest/docs/resources/resource) | resource |
| [null_resource.wait_for_cluster](https://registry.terraform.io/providers/hashicorp/null/latest/docs/resources/resource) | resource |
| [random_string.cluster](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/string) | resource |
| [random_string.role](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/string) | resource |
| [aws_availability_zones.available](https://registry.terraform.io/providers/hashicorp/aws/3.74/docs/data-sources/availability_zones) | data source |
| [aws_caller_identity.current](https://registry.terraform.io/providers/hashicorp/aws/3.74/docs/data-sources/caller_identity) | data source |
| [aws_eks_cluster_auth.cluster_auth](https://registry.terraform.io/providers/hashicorp/aws/3.74/docs/data-sources/eks_cluster_auth) | data source |
| [aws_route53_zone.main](https://registry.terraform.io/providers/hashicorp/aws/3.74/docs/data-sources/route53_zone) | data source |
| [http_http.workstation-external-ip](https://registry.terraform.io/providers/hashicorp/http/latest/docs/data-sources/http) | data source |
| [null_data_source.cluster](https://registry.terraform.io/providers/hashicorp/null/latest/docs/data-sources/data_source) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_aws_region"></a> [aws\_region](#input\_aws\_region) | n/a | `string` | `"us-east-2"` | no |
| <a name="input_certificate_enabled"></a> [certificate\_enabled](#input\_certificate\_enabled) | n/a | `bool` | `false` | no |
| <a name="input_cluster-name"></a> [cluster-name](#input\_cluster-name) | n/a | `string` | `"codepipes-demo"` | no |
| <a name="input_cluster_ipv4_cidr"></a> [cluster\_ipv4\_cidr](#input\_cluster\_ipv4\_cidr) | The IP address range of the kubernetes pods in this cluster. Default is an automatically assigned CIDR. | `string` | `"10.0.0.0/16"` | no |
| <a name="input_domain_name"></a> [domain\_name](#input\_domain\_name) | n/a | `string` | `null` | no |
| <a name="input_map_additional_iam_roles"></a> [map\_additional\_iam\_roles](#input\_map\_additional\_iam\_roles) | Additional IAM roles to add to `config-map-aws-auth` ConfigMap | <pre>list(object({<br>    rolearn  = string<br>    username = string<br>    groups   = list(string)<br>  }))</pre> | `[]` | no |
| <a name="input_node-group-name"></a> [node-group-name](#input\_node-group-name) | n/a | `string` | `"codepipes-cdn-node-group"` | no |
| <a name="input_role-eks-demo-node"></a> [role-eks-demo-node](#input\_role-eks-demo-node) | n/a | `string` | `"codepipes-cdn-eks-demo-node"` | no |
| <a name="input_vpc-eks-tag-name"></a> [vpc-eks-tag-name](#input\_vpc-eks-tag-name) | n/a | `string` | `"codepipes-cdn-eks-demo-tag-name"` | no |
| <a name="input_wait_for_cluster_cmd"></a> [wait\_for\_cluster\_cmd](#input\_wait\_for\_cluster\_cmd) | Custom local-exec command to execute for determining if the eks cluster is healthy. Cluster endpoint will be available as an environment variable called ENDPOINT | `string` | `" apk add curl; for i in `seq 1 60`; do curl -k $ENDPOINT/healthz >/dev/null && exit 0 || true; sleep 5; done; echo TIMEOUT && exit 1"` | no |
| <a name="input_wait_for_cluster_interpreter"></a> [wait\_for\_cluster\_interpreter](#input\_wait\_for\_cluster\_interpreter) | Custom local-exec command line interpreter for the command to determining if the eks cluster is healthy. | `list(string)` | <pre>[<br>  "/bin/sh",<br>  "-c"<br>]</pre> | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_certificate_arn"></a> [certificate\_arn](#output\_certificate\_arn) | The ARN of the certificate that is being validated. |
| <a name="output_cluster_endpoint"></a> [cluster\_endpoint](#output\_cluster\_endpoint) | Endpoint for your Kubernetes API server. |
| <a name="output_cluster_name"></a> [cluster\_name](#output\_cluster\_name) | Name of the cluster |
| <a name="output_cluster_region"></a> [cluster\_region](#output\_cluster\_region) | Cluster Region |
| <a name="output_config_map_aws_auth"></a> [config\_map\_aws\_auth](#output\_config\_map\_aws\_auth) | Generated AWS Auth Config Map |
| <a name="output_eks_arn"></a> [eks\_arn](#output\_eks\_arn) | ARN of the cluster role. |
| <a name="output_kubeconfig"></a> [kubeconfig](#output\_kubeconfig) | kubeconfig file |
| <a name="output_node_arn"></a> [node\_arn](#output\_node\_arn) | ARN of the node role. |
| <a name="output_rds_instance_address"></a> [rds\_instance\_address](#output\_rds\_instance\_address) | The hostname of the RDS instance. |
| <a name="output_rds_instance_endpoint"></a> [rds\_instance\_endpoint](#output\_rds\_instance\_endpoint) | The connection endpoint in address:port format. |
| <a name="output_rds_instance_id"></a> [rds\_instance\_id](#output\_rds\_instance\_id) | The RDS instance id. |
| <a name="output_redis_endpoint"></a> [redis\_endpoint](#output\_redis\_endpoint) | Elasticache redis connection endpoint in address:port format. |
| <a name="output_redis_hostname"></a> [redis\_hostname](#output\_redis\_hostname) | Elasticache redis address |
| <a name="output_redis_port"></a> [redis\_port](#output\_redis\_port) | Elasticache redis address |
| <a name="output_redis_security_group_id"></a> [redis\_security\_group\_id](#output\_redis\_security\_group\_id) | ID of the elasticache security group. |
| <a name="output_security_group_id"></a> [security\_group\_id](#output\_security\_group\_id) | ID of the db security group. |
| <a name="output_subnet_group_id"></a> [subnet\_group\_id](#output\_subnet\_group\_id) | The db subnet group name. |
<!-- END_TF_DOCS -->