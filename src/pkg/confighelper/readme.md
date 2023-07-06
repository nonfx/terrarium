
# Generic Configuration Package

This package provides a generic configuration system using the [Viper](https://github.com/spf13/viper) library, allowing users to define their own configuration prefix and provide default configuration.

## Usage

### Step 1: Create a default YAML configuration file

Create a YAML file that contains the default configuration for your project. For example, let's call it `defaults.yaml`:

```yaml
db:
  host: "localhost"
  user: "postgres"
  # password: "" # no default, panic if not set
  name: "cc_terrarium"
  port: 5432
  ssl_mode: false
```

### Step 2: Initialize the configuration

Import the `config` package and initialize the configuration by providing the prefix and default YAML configuration data:

```go
package main

import (
  "fmt"
  "embed"
  "log"

  "github.com/<username>/<package>/config"
)

//go:embed defaults.yaml
var defaultsYamlFile embed.FS

func main() {
  defaultsYaml, err := defaultsYamlFile.ReadFile("defaults.yaml")
  if err != nil {
    log.Fatal(err)
  }

  defaultsMap := map[string]interface{}{}
  err = yaml.Unmarshal(defaultsYaml, &defaultsMap)
  if err != nil {
    log.Fatal(err)
  }

  cfg, err := config.LoadDefaults(defaultsMap, "TR")
  if err != nil {
    log.Fatal(err)
  }

  // Use the configuration...
  host := viper.GetString("db.host")
  fmt.Println("Host:", host)
}
```

### Step 3: Set configuration values using environment variables

To override the default configuration values, you can set environment variables using the defined prefix and underscore separator. For example, to override the `db.host` configuration value, you can set the environment variable `TR_DB_HOST`.

### Step 4: Use the configuration values

You can access the configuration values using the Viper library and dot notation. For example, to access the `db.host` value, use:

```go
host := viper.GetString("db.host")
```

Refer to the Viper library documentation for more information on accessing configuration values.

#### Panic if not set

Use the `Must...` helper functions to panic if the value is not set. This is specially useful to help assert that the required values are set by the user. Example:

```go
host := config.MustGetString("db.host")
```
