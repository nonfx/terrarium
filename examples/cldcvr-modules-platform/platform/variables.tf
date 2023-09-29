#Common Variables
variable "common_name_prefix" {
  description = "The prefix to use for all resources in this example"
  default     = "terrarium"
}

variable "common_tags" {
  description = "The tags to apply to all resources in this example"
  type        = map(string)
  default = {
    Example   = "terrarium-demo"
    Terraform = "true"
  }
}

variable "environment" {
  description = "The environment to deploy to"
  default     = "dev"
}

#VPC Variables
variable "vpc_cidr" {
  description = "The CIDR block for the VPC. Default value is a valid CIDR, but not acceptable by AWS and should be overridden"
  default     = "10.5.0.0/16"
}

#EC2 Variables
variable "ec2_config" {
  description = "The configuration for the EC2 instance"
  default = {
    "default" : {
      ami : "ami-0e18308d78c527c8a"
      instance_type : "t2.micro"
      monitoring_enabled : true
    }
  }
}

#DB Config
variable "db_config" {
  description = "The configuration for the DB instance"
  default = {
    "default" : {
      instance_class : "db.t3.micro"
      allocated_storage : 20
      storage_type : "gp2"
      publicly_accessible : false
      multi_az : false
      db_subnet_group_name : "default"
      storage_encrypted : false
    }
  }
}

#Redis Config
variable "redis_config" {
  description = "The configuration for the Redis instance"
  default = {
    "default" : {
      instance_class : "cache.t3.micro"
      engine_version : "5.0.6"
      port : 6379
      parameter_group_name : "default.redis5.0"
      maintenance_window : "sun:05:00-sun:06:00"
      snapshot_window : "05:00-06:00"
      snapshot_retention_limit : 1
    }
  }
}
