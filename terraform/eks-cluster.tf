#
# EKS Cluster Resources
#  * IAM Role to allow EKS service to manage other AWS services
#  * EC2 Security Group to allow networking traffic with EKS cluster
#  * EKS Cluster

# resource "aws_cloudwatch_log_group" "demo" {
#   # The log group name format is /aws/eks/<cluster-name>/cluster
#   # Reference: https://docs.aws.amazon.com/eks/latest/userguide/control-plane-logs.html
#   name              = "/aws/eks/${aws_eks_cluster.demo.name}/cluster"
#   retention_in_days = 7
#   # ... potentially other configuration ...
# }

resource "aws_security_group" "demo-cluster" {
  depends_on = [aws_vpc.demo]
  name        = "eks-demo-cluster-sg-${random_string.cluster.id}"
  description = "Cluster communication with worker nodes"
  vpc_id      = aws_vpc.demo.id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "terraform-eks-demo"
  }
}

resource "aws_security_group_rule" "demo-cluster-ingress-workstation-https" {
  depends_on        = [aws_security_group.demo-cluster]
  cidr_blocks       = [local.workstation-external-cidr]
  description       = "Allow workstation to communicate with the cluster API Server"
  from_port         = 443
  protocol          = "tcp"
  security_group_id = aws_security_group.demo-cluster.id
  to_port           = 443
  type              = "ingress"
}

resource "aws_eks_cluster" "demo" {
  enabled_cluster_log_types = ["api", "audit","authenticator","controllerManager","scheduler"]
  name     = "${var.cluster-name}-${random_string.cluster.id}"
  role_arn = aws_iam_role.cluster.arn

  vpc_config {
    security_group_ids = [aws_security_group.demo-cluster.id]
    subnet_ids         = aws_subnet.demo[*].id
  }

  depends_on = [
    aws_internet_gateway.demo,
    aws_security_group.demo-cluster,
    aws_iam_role_policy_attachment.cluster-AmazonEKSClusterPolicy,
    aws_iam_role_policy_attachment.cluster-AmazonEKSVPCResourceController,
    aws_iam_role_policy_attachment.cluster-AmazonVPCFullAccess,
  ]
}
