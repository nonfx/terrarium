---
title: "Introduction"
slug: "introduction"
excerpt: "This section provides an overview of Terrarium features and benefits."
hidden: false
category: 64e49e46c016e9000c549a3a
---

# Overview

## What is Terrarium?

Terrarium simplifies the process of creating contemporary, comprehensive infrastructure, offering exceptional assistance to platform engineers and developer app support.Watch the [**Terrarium deep dive**](https://storage.googleapis.com/codepipes-assets/terrarium/assets/deepdiveterrarium.mp4) video to learn more.
In the evolving landscape of DevOps, Terrarium aligns with the principles of **Internal Developer Platforms (IDP)**, ensuring streamlined infrastructure provisioning and management.

## Workflow

<div align="center">
  <img src="https://storage.googleapis.com/codepipes-assets/terrarium/assets/Architecture.jpeg" alt="Image" style="border-radius: 10px;">
</div>

## Concepts

- [Terrarium Farm](https://terrarium.readme.io/docs/terrarium-farm) - The Terrarium Farm is a repository containing curated collection of content that enables essential features such as generalization of infrastructure dependency interfaces, taxonomy management and auto-complete suggestions for terraform and app-dependencies. The Ollion maintained Farm repo is at [cldcvr/terrarium-farm](https://github.com/cldcvr/terrarium-farm).
- [Terrarium Dependency Interface](https://terrarium.readme.io/docs/dependency-interface) - The Terrarium Dependency Interface is a contract that defines the inputs and outputs of an infrastructure dependency, facilitating communication between applications and Infrastructure as Code (IaC).
- [Terrarium Platform Framework](https://terrarium.readme.io/docs/terrarium-platform-template-framework-doc) - The Terrarium Platform Framework facilitates the creation of reusable Infrastructure as Code configurations using Terraform. Templates built using this framework are referred to as Terrarium Platform templates.
- [Terrarium App Manifest](https://terrarium.readme.io/docs/app-manifest) - Terrarium App metadata provides a way for applications to declare their infrastructure dependencies.

## Client Libraries

- [CLI](./setup.md) - Our command-line interface (CLI) offers a seamless experience for working with Terrarium. It includes commands for Terrarium Auto-Complete, leveraging the Terrarium Platform framework, adding dependencies to your applications, and ultimately composing a fully functional Terraform template by combining the Platform template and the app-specific dependencies.
- [VS Code Extension](https://github.com/cldcvr/terrarium-vscode-plugin) - Our VS Code extension is designed to assist DevOps professionals in writing Terraform code and streamlining infrastructure dependency declaration for the app Developers.
- [Web UI](https://github.com/cldcvr/terrarium-frontend) (coming soon) - The UI component, currently in development, will provide comprehensive documentation and visual representations of the Terrarium Platform Templates, the implemented dependencies within the Platform, and platform insights. This will help DevOps professionals gain a clear overview of the platform template and assist developers in selecting appropriate dependencies.
