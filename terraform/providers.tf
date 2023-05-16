terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
      version = "4.67.0"
    }
    azurerm = {
      source = "hashicorp/azurerm"
      version = "3.56.0"
    }
    google = {
      source = "hashicorp/google"
      version = "4.65.1"
    }
  }
}

provider "aws" {
  # Configuration options
}

provider "azurerm" {
  # Configuration options
}

provider "google" {
  # Configuration options
}
