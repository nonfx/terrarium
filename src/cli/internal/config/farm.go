package config

import "github.com/cldcvr/terrarium/src/pkg/confighelper"

// FarmDefault link of the farm repo to use by default
func FarmDefault() string {
	return confighelper.MustGetString("farm.default")
}
