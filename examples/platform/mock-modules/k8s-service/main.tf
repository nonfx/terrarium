variable "name" {
  type = string
  description = "Name of the service"
  default = ""
}

variable "port" {
  type = number
  description = "Port to bind the service at"
  default = 80
}

variable "cluster_id" {
  type = string
  description = "ID of the cluster"
  default = ""
}

variable "is_public" {
  type = bool
  description = "Should this service be exposed to the public network"
  default = false
}

resource "random_pet" "random" {
  length  = 2
}

output "host" {
  description = "The service host"
  value       = "${random_pet.random.id}.mocksite.com"
}
