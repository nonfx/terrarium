// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package localstate

import (
	"fmt"
	"log"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const (
	stateFileDir  = ".terrarium"
	stateFileName = ".terrarium_state"
	stateFileExt  = "yaml"
)

var (
	ls *localState
)

type localState struct {
	v             *viper.Viper
	stateFileName string
}

func init() {
	ls = new(localState)
	ls.stateFileName = fmt.Sprintf("%s/%s/%s.%s", findHomeDir(), stateFileDir, stateFileName, stateFileExt)
	reset()
}

// SetStateFileName sets the name of the underlying file used to store the state and causes
// the state to reload
func SetStateFileName(fn string) {
	ls.stateFileName = fn
	reset()
}

// Clear Removes the underlying file used to store the state and causes a reload (i.e. state will
// be empty)
func Clear() {
	os.Remove(ls.stateFileName)
	reset()
}

func reset() *localState {
	ls.v = viper.New()
	ls.v.SetConfigFile(ls.stateFileName)

	if err := ls.v.ReadInConfig(); err != nil {
		ls.v.WriteConfigAs(ls.stateFileName)
		ls.v.ReadInConfig()
	}
	return ls
}

func findHomeDir() string {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}
	return home
}

// Get return a value by key from the local state
func Get(key string) string {
	return ls.v.GetString(key)
}

// Set sets a local state value by key
func Set(key string, value string) {
	ls.v.Set(key, value)
	err := ls.v.WriteConfig()
	if err != nil {
		log.Fatal(err)
	}
}

// Unset remove a key from the local state
func Unset(key string) {
	configMap := ls.v.AllSettings()
	delete(configMap, key)

	// Unfortunately Viper doesn't have an Unset method yet. So this will reload all
	// the settings except the one to be deleted.
	ls.v = viper.New()
	for k, v := range configMap {
		ls.v.Set(k, v)
	}
	ls.v.WriteConfigAs(ls.stateFileName)
}

func List() map[string]interface{} {
	return ls.v.AllSettings()
}
