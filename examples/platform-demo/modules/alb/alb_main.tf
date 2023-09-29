module "alb" {
  source  = "terraform-aws-modules/alb/aws"
  version = "~> 8.7"

  for_each = var.tf_component_lb

  name = "${var.extract_resource_name}-alb"

  create_lb = local.tr_web_service == true ? true : false

  load_balancer_type = each.value.load_balancer_type
  create_security_group = true

  vpc_id          = var.vpc_id
  subnets         = var.public_subnet_ids
  security_groups = var.security_group_ids

  access_logs = {
    bucket = module.s3_bucket[each.key].s3_bucket_id
  }

  security_group_rules = [
    {
      type        = "ingress"
      from_port   = 80
      to_port     = 80
      protocol    = "TCP"
      cidr_blocks = ["0.0.0.0/0"]
    },
    {
      type        = "ingress"
      from_port   = 443
      to_port     = 443
      protocol    = "TCP"
      cidr_blocks = ["0.0.0.0/0"]
    },
    {
      type        = "egress"
      from_port   = 0
      to_port     = 0
      protocol    = "-1"
      cidr_blocks = ["0.0.0.0/0"]
    }
  ]


  target_groups = [
    for service_key, service_value in var.tr_component_ecs_services : {
      name_prefix      = substr("${service_key}",0,6)
      backend_protocol = "HTTP"
      backend_port     = try(service_value.port, null)
      target_type      = "ip"
    }
    if try(service_value.port, null) != null
  ]

  https_listeners = [
    {
      port               = 443
      protocol           = "HTTPS"
      certificate_arn    = each.value.certificate_arn
      target_group_index = 0
    }
  ]

  http_tcp_listeners = [
    {
      port        = 80
      protocol    = "HTTP"
      action_type = "redirect"
      redirect = {
        port        = "443"
        protocol    = "HTTPS"
        status_code = "HTTP_301"
      }
    }
  ]

  tags = {
    Environment = "Test"
  }
}

resource "random_id" "bucket_suffix" {
  byte_length = 4
  keepers = {
    bucket_base_name = var.extract_resource_name
  }
}

module "s3_bucket" {
  source  = "terraform-aws-modules/s3-bucket/aws"
  version = "3.14.0"

  for_each = var.tf_component_lb

  bucket = "${var.extract_resource_name}-alb-logs-${random_id.bucket_suffix.hex}"
  acl    = "log-delivery-write"

  # Allow deletion of non-empty bucket
  force_destroy = var.environment == "production" || var.environment == "prod" ? false : true

  control_object_ownership = true
  object_ownership         = "ObjectWriter"

  attach_elb_log_delivery_policy = true # Required for ALB logs
  attach_lb_log_delivery_policy  = true # Required for ALB/NLB logs
}

