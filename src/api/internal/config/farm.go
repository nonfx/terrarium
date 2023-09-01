package config

import "github.com/cldcvr/terrarium/src/pkg/confighelper"

// FarmDefault link to the farm repo being utilized
func FarmDefault() string {
	return confighelper.MustGetString("farm.default")
}
