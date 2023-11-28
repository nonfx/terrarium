locals {
  project_org_id    = var.folder_id != "" ? null : var.org_id
  project_folder_id = var.folder_id != "" ? var.folder_id : null
}

resource "google_project" "project" {
  name                = var.name
  project_id          = var.project_id
  org_id              = local.project_org_id
  folder_id           = local.project_folder_id
  billing_account     = var.billing_account
  auto_create_network = var.auto_create_network
  labels              = var.labels
  skip_delete         = var.skip_delete
}

resource "google_project_service" "enable_service" {
  for_each                   = toset(var.service)
  project                    = var.project_id
  service                    = each.key
  disable_dependent_services = var.disable_dependent_services
}

output "number" {
  value       = google_project.project.number
  description = "The numeric identifier of the project."
}

variable "name" {
  description = "The name for the project"
  type        = string
  default     = ""
}

variable "project_id" {
  description = "The ID to give the project. If not provided, the `name` will be used."
  type        = string
  default     = ""
}

variable "org_id" {
  description = "The organization ID."
  type        = string
  default     = ""
}

variable "folder_id" {
  description = "The ID of a folder to host this project"
  type        = string
  default     = ""
}

variable "billing_account" {
  description = "The ID of the billing account to associate this project with"
  type        = string
  default     = ""
}

variable "skip_delete" {
  description = "If true, the Terraform resource can be deleted without deleting the Project via the Google API."
  type        = bool
  default     = false
}

variable "labels" {
  description = "Map of labels for project"
  type        = map(string)
  default     = {}
}

variable "auto_create_network" {
  description = "Create the default network"
  type        = bool
  default     = true
}

variable "service" {
  description = "The service to enable."
  type        = list(string)
  default     = []
}

variable "disable_dependent_services" {
  description = "The service to enable."
  type        = bool
  default     = false
}

terraform {
  required_version = ">= 0.13"
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 4.0"
    }
  }
}
