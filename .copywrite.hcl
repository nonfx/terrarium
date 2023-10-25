schema_version = 1

project {
  license        = "Apache-2.0"
  copyright_year = 2023
  copyright_holder = "Ollion"

  # (OPTIONAL) A list of globs that should not have copyright/license headers.
  # Supports doublestar glob patterns for more flexibility in defining which
  # files or folders should be ignored
  header_ignore = [
    "**.md",
    "**/*.tf",
    "**/*.tfvars",
    "**/testdata/**",
    "**/*.pb.go",
    "**/mocks",
    "**/.terraform",
    ".git",
    "coverage",
    ".bin",
    "dist",
    "**/.goreleaser.yaml"
  ]
}
