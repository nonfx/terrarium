variable "display_name" {
  type        = string
  description = "The folder’s display name"
  default     = ""
}

variable "parent" {
  type        = string
  description = "The resource name of the parent Folder or Organization."
  default     = ""
}