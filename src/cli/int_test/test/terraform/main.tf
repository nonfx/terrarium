module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "18.31.2"
}

module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "4.0.2"
}
