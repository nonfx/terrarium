<a href="https://terrarium.cldcvr.com/">
  <img alt="terrarium" src="https://storage.googleapis.com/codepipes-assets/terrarium/assets/T8-logo.png" width="100" />
</a>

[![Documentation](https://img.shields.io/badge/Documentation-Design_Docs-blue?logo=read-the-docs)](https://terrarium.readme.io/docs/)
[![Build results](https://github.com/cldcvr/terrarium/actions/workflows/release.yaml/badge.svg)](https://github.com/cldcvr/terrarium/actions/workflows/release.yaml)
[![GitHub license](https://img.shields.io/github/license/cldcvr/terrarium)](https://github.com/cldcvr/terrarium/blob/main/LICENSE)
[![Discord](https://img.shields.io/discord/1130781563928444978)](https://discord.gg/gG3gDm9GmF)
[![Contribute](https://img.shields.io/badge/Contribute-How_to_Contribute-blue?logo=github)](./CONTRIBUTING.md)

---

# Terrarium

The Terrarium project aims to empower platform engineering professionals by providing a comprehensive suite of tools for creating reusable Terraform templates. In the evolving landscape of DevOps, Terrarium aligns with the principles of [Internal Developer Platforms (IDP)](https://internaldeveloperplatform.org/), ensuring streamlined infrastructure provisioning and management.

## Installation Instructions

### Prerequisite

<img alt="Terraform" src="https://www.datocms-assets.com/2885/1620155116-brandhcterraformverticalcolor.svg" width="150px">                     <img alt="Golang" src="https://raw.githubusercontent.com/github/explore/80688e429a7d4ef2fca1e82350fe8e3517d3494d/topics/yaml/yaml.png" width="100px">                    <img alt="Terraform" src="https://user-images.githubusercontent.com/25181517/192149581-88194d20-1a37-4be8-8801-5dc0017ffbbe.png" width="100px">

### Steps
1. Download [Terrarium](https://github.com/cldcvr/terrarium/releases) and extract the TAR archive.

   ```bash
   wget https://github.com/cldcvr/terrarium/releases/download/$VERSION/terrarium-$VERSION-linux-amd64.tar.gz
   ```

   Example:

```bash
  wget https://github.com/cldcvr/terrarium/releases/download/v0.4/terrarium-v0.4-macos-amd64.tar.gz
  tar -xzf terrarium-v0.4-macos-amd64.tar.gz
```

2. Move the `terrarium` binary to a directory in your system's PATH, like `/usr/local/bin/`.
   Add this to your shell:

   ```bash
   PATH="$PATH:/path/to/terrarium"
   ```

3. Alternatively, install using the source code.
   Clone this repo and execute:

   ```bash
   make install
   ```
> [!IMPORTANT]
> Make sure you have go version 1.20 or above if not try this link to download https://go.dev/dl/



4. Verify the Installation

   To check if Terrarium is installed correctly, open your terminal or command prompt and run:

   ```bash
   terrarium version
   ```

### Get Started with Terrarium

## Client Libraries

| Tool                                                                   | Description                                                          |
| ---------------------------------------------------------------------- | -------------------------------------------------------------------- |
| [VS Code Extension](https://github.com/cldcvr/terrarium-vscode-plugin) | Assists DevOps in writing Terraform code and managing dependencies.  |
| [CLI](./setup.md)                                                      | Command-line interface for working with Terrarium and app templates. |
| [Web UI](https://github.com/cldcvr/terrarium-frontend) (coming soon)   | User interface for visualizing Terrarium Platform Templates.         |
| [API](./src/api/) (internal)                                           | Internal query server for Terrarium Farm repo content.               |

## Concepts

- [App Manifest](./src/pkg/metadata/app/readme.md) - App Manifest provides a way for an applications to declare its infrastructure requirements using generic dependency interfaces. Such as, a working Terraform template can be generated at the time of deployment using the best practice defined in the Terrarium platform template.
- [Terrarium Dependency Interface](./src/pkg/metadata/dependency/readme.md) - The Terrarium Dependency Interface is an agreement that outlines how applications and Infrastructure as Code (IaC) interact. Dependencies are implemented in platforms and used in apps. A single dependency can be built into various platforms but only once per platform. However, an app can use the same dependency multiple times.
- [Terrarium Platform Framework](./examples/platform/readme.md) - The Terrarium Platform Framework helps make reusable templates with Terraform. A Terrarium Platform Template implements dependencies in an opinionated way, exposing only relevant controls to the app and generating a defined set of outputs for the app to use as environment variables.
- [Terrarium Farm](./examples/farm/readme.md) - The Terrarium Farm is a repository containing seed data like tf-modules, dependencies, taxonomy & platforms. The farm repository has workflows to ensure the sanity of the content as well as scan the content to extract key information. The Ollion maintained Farm repo is at [cldcvr/terrarium-farm](https://github.com/cldcvr/terrarium-farm).

## Flow

### Basic

```mermaid
flowchart LR
    PT["Terrarium Platform Template\neg: cc-aws/cc-gcp/cc-azure"]
    DEP["Dependencies\neg: pgsql, redis,\nweb-server, background-worker, static-server"]
    APP["Applications\neg: backend, frontend, worker"]
    TR["Terrarium CLI"]
    WT["Working Terraform Template"]
    ENV["Environment Variables"]

    DEP --"implemented in"--> PT
    DEP --"used in"--> APP

    PT --"one"--> TR
    APP --"many"--> TR
    TR --"generates"--> WT
    TR --"generates"--> ENV
```

### How it works

<div align="center">
<img src="./_docs/terrarium-ref-diag.png" alt="end-to-end" style="border-radius: 10px;">
</div>

### Repos & Deployment setup

```mermaid
flowchart TD
    subgraph "Git Repos"
    Farm["Terrarium Farm\nCC Owned"]
    CCP["Terrarium tempaltes\nCC Owned"]
    OSP["Terrarium tempaltes\nOther open source repos"]
    end

    subgraph "Public SaaS"
    DB["Public DB"]
    API
    UI["T8 Web UI Interface"]

    CCP --"Release pipeline"--> DB
    Farm --"Release pipeline"--> DB
    DB --> API
    API --> UI
    OSP --"link to\nbranch/tag/commit"--> UI
    end

    subgraph "User Local"
    CLI["T8 CLI"]
    VS["T8 VS Code Ext."]
    GC["git clone"]

    Farm --"Release pipeline"--> CLI
    OSP --"link to\nbranch/tag/commit"--> GC
    GC --> CLI
    CLI --> VS
    end
```

## Progress

- [x] Release VS-Code extension with basic auto-complete from curated modules in the farm-repo.
- [x] Document the Terrarium Platform Framework and app dependencies format with examples.
- [x] Implement CLI command to lint & parse Terrarium Platform Templates.
- [x] Simplify the installation of Postgres & T8-API in Docker for the VS Code extension by integrating it into the CLI.
- [x] Develop CLI command to compose working Terraform templates using T8-Platform Templates & App Dependencies.
- [x] Enhance VS Code plugin to support auto-complete with local modules.
- [x] Add dependency-interfaces content to the farm-repo.
- [ ] Develop CLI & VS Code plugin features to assist developers in declaring infrastructure dependencies for their apps.
- [ ] Create a UI to help developers select App Dependencies by showcasing platform and farm insights.
- [ ] Add taxonomy mappings to the farm-repo.
- [ ] Enhance VS Code plugin to automatically implement dependency-interfaces in a platform (best guess).

## Get Involved

Join our Discord community - [**Terrarium Community**](https://discord.gg/gG3gDm9GmF).

Terrarium is still in its early stages, and we welcome your contributions.

To file a bug, suggest an improvement, or request a new feature please open an issue, refer to our [contributing guide](./CONTRIBUTING.md)
