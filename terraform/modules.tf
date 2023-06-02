
# # load module from Terraform Module Registry
#module "eks" {
#  source  = "terraform-aws-modules/eks/aws"
#  version = "19.14.0"
#}
#
# # load module directly from GitHub
# module "rds" {
#  source = "github.com/terraform-aws-modules/terraform-aws-vpc.git"
# }

# # load module from local file system
# module "vpc" {
#  source = "./vpc"
# }

module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "4.0.2"
}

module "rds" {
  source  = "terraform-aws-modules/rds/aws"
  version = "5.9.0"
}

module "security-group" {
  source  = "terraform-aws-modules/security-group/aws"
  version = "4.17.2"
}

module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "19.14.0"
}

module "s3-bucket" {
  source  = "terraform-aws-modules/s3-bucket/aws"
  version = "3.11.0"
}
