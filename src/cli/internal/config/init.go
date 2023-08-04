package config

import (
	"embed"

	"github.com/charmbracelet/log"
	"github.com/cldcvr/terrarium/src/pkg/confighelper"
	"gopkg.in/yaml.v3"
)

const ENV_PREFIX = "tr"

//go:embed defaults.yaml
var defaultsYamlFile embed.FS

func LoadDefaults() {
	defaultsYaml, err := defaultsYamlFile.ReadFile("defaults.yaml")
	if err != nil {
		panic(err)
	}

	defaultMap := map[string]interface{}{}
	err = yaml.Unmarshal(defaultsYaml, &defaultMap)
	if err != nil {
		panic(err)
	}

	confighelper.LoadDefaults(defaultMap, ENV_PREFIX)

	LoggerConfig(log.Default())
}
