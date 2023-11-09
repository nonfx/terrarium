variable "name" {
  description = "name of the bucket"
  type = string
}

output "uri" {
  value = base64encode(var.name)
}
