---
title: "Installation"
slug: "installation"
excerpt: "Get started with Terrarium by following these simple installation instructions."
hidden: false
category: 64fad6b9f42eab14641a53a6
---

# Installation Instructions

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

4. Verify the Installation

   To check if Terrarium is installed correctly, open your terminal or command prompt and run:

   ```bash
   terrarium version
   ```
