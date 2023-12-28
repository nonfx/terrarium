variable "bucket_name" {
  description = "Name of the bucket"
  type        = string
  default = "by_bucket"
}

variable "index_document" {
  description = "Index document for the website"
  type        = string
  default = "index.html"
}

variable "error_document" {
  description = "Error document for the website"
  type        = string
  default = "error.html"
}

resource "random_pet" "random" {
  length  = 2
}

output "website_endpoint" {
  description = "The website endpoint URL"
  value       = "http://${random_pet.random.id}.mocksite.com"
}
