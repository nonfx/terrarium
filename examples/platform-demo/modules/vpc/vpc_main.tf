
data "aws_availability_zones" "available" {}

module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "5.0.0"

  for_each = var.tr_component_vpc

  name = "${var.extract_resource_name}-vpc"
  cidr = coalesce(each.value.vpc_cidr_block, "10.0.0.0/16")

  manage_default_network_acl    = true
  manage_default_route_table    = true
  manage_default_security_group = true

  # Fetch the availability zones based on the number provided in coalesce(each.value.number_of_azs, 2)
  azs = slice(data.aws_availability_zones.available.names, 0, coalesce(each.value.number_of_azs, 2))

  # Generate private subnets based on the VPC CIDR block and the number of availability zones
  private_subnets = [for i in range(coalesce(each.value.number_of_azs, 2)) : cidrsubnet(coalesce(each.value.vpc_cidr_block, "10.0.0.0/16"), 8, i)]

  # Generate public subnets based on the VPC CIDR block, the number of availability zones, and the offset of coalesce(each.value.number_of_azs, 2)
  public_subnets = [for i in range(coalesce(each.value.number_of_azs, 2)) : cidrsubnet(coalesce(each.value.vpc_cidr_block, "10.0.0.0/16"), 8, i + coalesce(each.value.number_of_azs, 2))]

  database_subnets = [for i in range(coalesce(each.value.number_of_azs, 2)) : cidrsubnet(coalesce(each.value.vpc_cidr_block, "10.0.0.0/16"), 8, 4 + i + coalesce(each.value.number_of_azs, 2))]

  create_database_subnet_group           = true
  create_database_subnet_route_table     = true

  enable_nat_gateway = true
  enable_vpn_gateway = true

  # Enable VPC flow logs with role and groups if the environment is production, otherwise disable them
  enable_flow_log                      = var.environment == "production" ? true : false
  create_flow_log_cloudwatch_iam_role  = var.environment == "production" ? true : false
  create_flow_log_cloudwatch_log_group = var.environment == "production" ? true : false

  tags = merge(
    {
      "Name" = format("%s", "${var.extract_resource_name}-vpc")
    },
    {
      environment = var.environment
    },
    var.tags,
  )
}
