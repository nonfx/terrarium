output "vpc_id" {
  value = {
    for vpc_name, vpc in module.vpc :
    vpc_name => vpc.vpc_id
  }
}

output "vpc_cidr_blocks" {
  value = {
    for vpc_name, vpc in module.vpc :
    vpc_name => vpc.vpc_cidr_block
  }
}

output "public_subnet_ids" {
  value = {
    for vpc_name, vpc in module.vpc :
    vpc_name => vpc.public_subnets
  }
}

output "private_subnet_ids" {
  value = {
    for vpc_name, vpc in module.vpc :
    vpc_name => vpc.private_subnets
  }
}

output "availability_zones" {
  value = {
    for vpc_name, vpc in module.vpc :
    vpc_name => vpc.azs
  }
}

output "database_subnets" {
  value = {
    for vpc_name, vpc in module.vpc :
    vpc_name => vpc.database_subnets
  }
}

output "database_subnet_group" {
  value = {
    for vpc_name, vpc in module.vpc :
    vpc_name => vpc.database_subnet_group
  }
}
