variable "name" {
  type = string
  description = "Name cluster name"
  default = ""
}

resource "random_id" "random" {
  byte_length = 8
}

output "cluster_id" {
  description = "cluster id"
  value       = random_id.random.id
}
