// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package confighelper

import (
	"strings"

	"github.com/rotisserie/eris"
	"github.com/spf13/viper"
)

// LoadDefaults loads the default configuration values provided as a map into Viper.
// It sets the environment prefix and key replacer for Viper to handle environment variables.
func LoadDefaults(defaults map[string]interface{}, prefix string) {
	viper.SetEnvPrefix(prefix)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetTypeByDefaultValue(true)

	for k, v := range defaults {
		viper.SetDefault(k, v)
	}
}

// MustGetString retrieves the value associated with the given key as a string from Viper.
// It uses the MustGet function to retrieve the value and panics if the value is not found.
func MustGetString(key string) string {
	return MustGet(key, viper.GetString)
}

// MustGetBool retrieves the value associated with the given key as a bool from Viper.
// It uses the MustGet function to retrieve the value and panics if the value is not found.
func MustGetBool(key string) bool {
	return MustGet(key, viper.GetBool)
}

// MustGetInt retrieves the value associated with the given key as an int from Viper.
// It uses the MustGet function to retrieve the value and panics if the value is not found.
func MustGetInt(key string) int {
	return MustGet(key, viper.GetInt)
}

// MustGetInt64 retrieves the value associated with the given key as an int64 from Viper.
// It uses the MustGet function to retrieve the value and panics if the value is not found.
func MustGetInt64(key string) int64 {
	return MustGet(key, viper.GetInt64)
}

// MustGetFloat64 retrieves the value associated with the given key as a float64 from Viper.
// It uses the MustGet function to retrieve the value and panics if the value is not found.
func MustGetFloat64(key string) float64 {
	return MustGet(key, viper.GetFloat64)
}

// MustGetAny retrieves the value associated with the given key as an interface{} from Viper.
// It uses the MustGet function to retrieve the value and panics if the value is not found.
// The retrieved value is returned as an interface{} type, allowing flexibility in handling different types of values.
func MustGetAny(key string) interface{} {
	return MustGet(key, viper.Get)
}

// MustGet retrieves the value associated with the given key using the provided getter function from Viper.
// It panics if the value is not found.
func MustGet[T comparable](key string, getter func(string) T) T {
	val, err := MustGetOrError[T](key, getter)
	if err != nil {
		panic(err)
	}

	return val
}

// MustGetOrError retrieves the value associated with the given key using the provided getter function from Viper.
// It returns an error if the value is not found.
func MustGetOrError[T comparable](key string, getter func(string) T) (val T, err error) {
	if !viper.IsSet(key) {
		err = eris.Errorf("missing configuration: %s", key)
		return
	}

	val = getter(key)

	return
}
