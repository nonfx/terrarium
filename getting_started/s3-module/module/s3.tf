module "tr_component_s3" {
  source  = "../"
  #version = "~> 4.0"

  bucket_name = "himanshu"
  name = "tag_name"
  environment ="env_name"
}