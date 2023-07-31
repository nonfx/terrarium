# Public modules

module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "18.31.2"
}

module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "4.0.2"
}

module "security-group" {
  source  = "terraform-aws-modules/security-group/aws"
  version = "5.1.0"
}

module "rds" {
  source  = "terraform-aws-modules/rds/aws"
  version = "5.1.1"
}

module "kms" {
  source  = "terraform-aws-modules/kms/aws"
  version = "1.5.0"
}

module "cloudwatch-kms-key" {
  source  = "dod-iac/cloudwatch-kms-key/aws"
  version = "1.0.1"
}

#  Custom tf templates for mappings discovery

module "tr-hide-banking-demo" {
  source = "github.com/cldcvr/codepipes-tutorials//tfs/aws-ecr-apprunner-vpc?ref=terrarium-sources"
}

module "tr-hide-voting-demo" {
  source = "github.com/cldcvr/codepipes-tutorials//voting/infra/aws/eks?ref=terrarium-sources"
}

module "tr-hide-serverless-sample" {
  source = "github.com/cldcvr/codepipes-tutorials//serverless-sample/infra/aws?ref=terrarium-sources"
}

module "tr-hide-wpdemo-eks" {
  source = "github.com/cldcvr/codepipes-tutorials//wpdemo/infra/aws/eks?ref=terrarium-sources"
}
