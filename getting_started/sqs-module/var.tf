variable "tags" {
    default= {
    Environment = "dev"
  }
}

variable "kms_master_key_id"{
    default = "0d1ba9e8-9421-498a-9c8a-01e9772b2924"
} 
  
  
variable "kms_data_key_reuse_period_seconds" {
    default= 3600
}