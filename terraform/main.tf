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

module "cdn" {
  source = "github.com/cldcvr/vanguard-demo//cdn/infra/aws/eks"
}

module "codepipes-iac" {
  source = "github.com/cldcvr/vanguard-iac//modules/vanguard-infra"
}
