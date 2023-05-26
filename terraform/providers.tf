terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "4.67.0"
    }
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "3.56.0"
    }
    google = {
      source  = "hashicorp/google"
      version = "4.65.1"
    }
    random = {
      source  = "hashicorp/random"
      version = "3.5.1"
    }
    time = {
      source  = "hashicorp/time"
      version = "0.9.1"
    }
    tls = {
      source  = "hashicorp/tls"
      version = "4.0.4"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "2.20.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "2.9.0"
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

provider "random" {
  # Configuration options
}

provider "time" {
  # Configuration options
}

provider "tls" {
  # Configuration options
}

provider "kubernetes" {
  # Configuration options
}

provider "helm" {
  # Configuration options
}
