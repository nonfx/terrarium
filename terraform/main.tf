# Public modules

module "eks" {
  source = "terraform-aws-modules/eks/aws"
  version = "18.31.2"
}

module "vpc" {
  source = "terraform-aws-modules/vpc/aws"
  version = "4.0.2"
}

module "security-group" {
  source = "terraform-aws-modules/security-group/aws"
  version = "5.1.0"
}

#  Custom tf templates for mappings discovery

module "voting-demo" {
  source = "github.com/cldcvr/codepipes-tutorials//voting/infra/aws/eks?ref=terrarium-sources"
}

module "serverless-sample" {
  source = "github.com/cldcvr/codepipes-tutorials//serverless-sample/infra/aws?ref=terrarium-sources"
}

module "wpdemo-eks" {
  source = "github.com/cldcvr/codepipes-tutorials//wpdemo/infra/aws/eks?ref=terrarium-sources"
}

module "wpdemo-ec2" {
  source = "github.com/cldcvr/codepipes-tutorials//wpdemo/infra/aws/ec2?ref=terrarium-sources"
}

# # Private repos
# module "cdn" {
#   source = "github.com/cldcvr/vanguard-demo//cdn/infra/aws/eks"
# }

# module "codepipes-iac" {
#   source = "github.com/cldcvr/vanguard-iac//modules/vanguard-infra"
# }
