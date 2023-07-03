
module "vpc_endpoints2" {
  source = "terraform-aws-modules/vpc/aws//modules/vpc-endpoints"

  vpc_id             = "module.vpc.vpc_id"
  # security_group_ids = [data.aws_security_group.default.id]

  endpoints = {
    s3 = local.tr_if_s3 ? {
      service = "s3"
      tags    = { Name = "s3-vpc-endpoint" }
    }: null,
  #   dynamodb = {
  #     service         = "dynamodb"
  #     service_type    = "Gateway"
  #     route_table_ids = flatten([module.vpc.intra_route_table_ids, module.vpc.private_route_table_ids, module.vpc.public_route_table_ids])
  #     policy          = data.aws_iam_policy_document.dynamodb_endpoint_policy.json
  #     tags            = { Name = "dynamodb-vpc-endpoint" }
  #   },
  #   ssm = {
  #     service             = "ssm"
  #     private_dns_enabled = true
  #     subnet_ids          = module.vpc.private_subnets
  #     security_group_ids  = [aws_security_group.vpc_tls.id]
  #   },
  #   ssmmessages = {
  #     service             = "ssmmessages"
  #     private_dns_enabled = true
  #     subnet_ids          = module.vpc.private_subnets
  #     security_group_ids  = [aws_security_group.vpc_tls.id]
  #   },
  #   lambda = {
  #     service             = "lambda"
  #     private_dns_enabled = true
  #     subnet_ids          = module.vpc.private_subnets
  #   },
  #   ecs = {
  #     service             = "ecs"
  #     private_dns_enabled = true
  #     subnet_ids          = module.vpc.private_subnets
  #   },
  #   ecs_telemetry = {
  #     create              = false
  #     service             = "ecs-telemetry"
  #     private_dns_enabled = true
  #     subnet_ids          = module.vpc.private_subnets
  #   },
  #   ec2 = {
  #     service             = "ec2"
  #     private_dns_enabled = true
  #     subnet_ids          = module.vpc.private_subnets
  #     security_group_ids  = [aws_security_group.vpc_tls.id]
  #   },
  #   ec2messages = {
  #     service             = "ec2messages"
  #     private_dns_enabled = true
  #     subnet_ids          = module.vpc.private_subnets
  #     security_group_ids  = [aws_security_group.vpc_tls.id]
  #   },
  #   ecr_api = {
  #     service             = "ecr.api"
  #     private_dns_enabled = true
  #     subnet_ids          = module.vpc.private_subnets
  #     policy              = data.aws_iam_policy_document.generic_endpoint_policy.json
  #   },
  #   ecr_dkr = {
  #     service             = "ecr.dkr"
  #     private_dns_enabled = true
  #     subnet_ids          = module.vpc.private_subnets
  #     policy              = data.aws_iam_policy_document.generic_endpoint_policy.json
  #   },
  #   kms = {
  #     service             = "kms"
  #     private_dns_enabled = true
  #     subnet_ids          = module.vpc.private_subnets
  #     security_group_ids  = [aws_security_group.vpc_tls.id]
  #   },
  #   codedeploy = {
  #     service             = "codedeploy"
  #     private_dns_enabled = true
  #     subnet_ids          = module.vpc.private_subnets
  #   },
  #   codedeploy_commands_secure = {
  #     service             = "codedeploy-commands-secure"
  #     private_dns_enabled = true
  #     subnet_ids          = module.vpc.private_subnets
  #   },
  }

  # tags = merge(local.tags, {
  #   Project  = "Secret"
  #   Endpoint = "true"
  # })
}
