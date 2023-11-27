# variable "cluster_endpoint_public_access"{
#    default = true 
# }  

variable  "cluster_addons"{
    default =  {
    coredns = {
      most_recent = true
    }
    kube-proxy = {
      most_recent = true
    }
    vpc-cni = {
      most_recent = true
    }
  }
} 

variable "vpc_id" {
    default = "vpc-1234556abcdef"
}

variable "subnet_ids" {
    default = ["subnet-abcde012", "subnet-bcde012a", "subnet-fghi345a"]
}
variable "control_plane_subnet_ids"{
  default = ["subnet-xyzde987", "subnet-slkjf456", "subnet-qeiru789"]
}

# Self Managed Node Group(s)
variable "self_managed_node_group_defaults" {
    default = {
    instance_type                          = "m6i.large"
    update_launch_template_default_version = true
    iam_role_additional_policies = {
      AmazonSSMManagedInstanceCore = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
    }
  }
}
variable   "self_managed_node_groups" {
    default = {
    one = {
      name         = "mixed-1"
      max_size     = 5
      desired_size = 2

      use_mixed_instances_policy = true
      mixed_instances_policy = {
        instances_distribution = {
          on_demand_base_capacity                  = 0
          on_demand_percentage_above_base_capacity = 10
          spot_allocation_strategy                 = "capacity-optimized"
        }

        override = [
          {
            instance_type     = "m5.large"
            weighted_capacity = "1"
          },
          {
            instance_type     = "m6i.large"
            weighted_capacity = "2"
          },
        ]
      }
    }
  }
}

# EKS Managed Node Group(s)
 variable "eks_managed_node_group_defaults"{ 
    default = {
    instance_types = ["m6i.large", "m5.large", "m5n.large", "m5zn.large"]
  }
 }
 variable "eks_managed_node_groups"{ 
    default = {
    blue = {}
    green = {
      min_size     = 1
      max_size     = 10
      desired_size = 1

      instance_types = ["t3.large"]
      capacity_type  = "SPOT"
    }
  }
 }

 # Fargate Profile(s)
  variable "fargate_profiles" {
    default = {
    default = {
      name = "default"
      selectors = [
        {
          namespace = "default"
        }
      ]
    }
  }
  }

# aws-auth configmap
 variable "manage_aws_auth_configmap"{
    default = true
 }
 variable "aws_auth_roles"{
    default = [
    {
      rolearn  = "arn:aws:iam::66666666666:role/role1"
      username = "role1"
      groups   = ["system:masters"]
    },
  ]
 }

variable "aws_auth_users" {
    default = [
    {
      userarn  = "arn:aws:iam::66666666666:user/user1"
      username = "user1"
      groups   = ["system:masters"]
    },
    {
      userarn  = "arn:aws:iam::66666666666:user/user2"
      username = "user2"
      groups   = ["system:masters"]
    },
  ]
}

 variable "aws_auth_accounts" {
    default= [
    "777777777777",
    "888888888888",
  ]
 }
variable "tags" {
    default= {
    Environment = "dev"
    Terraform   = "true"
  }
}
