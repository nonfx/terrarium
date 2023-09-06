// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package confighelper

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T) {
	defaults := map[string]interface{}{
		"db": map[string]interface{}{
			"host":      "localhost",
			"port":      5432,
			"user":      "postgres",
			"ssl_mode":  true,
			"test_var2": true,
			"test_var3": false,
		},
	}

	prefix := "TEST"

	viper.Reset()
	LoadDefaults(defaults, prefix)

	type testCase struct {
		key     string
		wantVal interface{}
		wantErr bool
	}

	t.Setenv("TEST_DB_PORT", "2222")
	t.Setenv("TEST_DB_USER", "testuser")
	t.Setenv("TEST_DB_TEST_VAR2", "false")
	t.Setenv("TEST_DB_TEST_VAR3", "true")

	t.Run("LoadDefaults", func(t *testing.T) {
		for _, tt := range []testCase{
			{key: "db.host", wantVal: "localhost"},
			{key: "db.port", wantVal: 2222},
			{key: "db.ssl_mode", wantVal: true},
			{key: "db.password", wantErr: true},
		} {
			t.Run(tt.key, func(t *testing.T) {
				assert.Equal(t, tt.wantErr, !viper.IsSet(tt.key), "check if variable is set")
				if !tt.wantErr {
					val := viper.Get(tt.key)
					assert.Equal(t, tt.wantVal, val, "assert value of config")
				}
			})
		}
	})

	t.Run("MustGetAny", func(t *testing.T) {
		for _, tt := range []testCase{
			{key: "db.host", wantVal: "localhost"},
			{key: "db.port", wantVal: 2222},
			{key: "db.ssl_mode", wantVal: true},
			{key: "db.password", wantErr: true},
		} {
			t.Run(tt.key, func(t *testing.T) {
				AssertPanic(t, tt.wantErr, func() {
					val := MustGetAny(tt.key)
					assert.Equal(t, tt.wantVal, val, "assert value of config")
				})
			})
		}
	})

	t.Run("MustGetString", func(t *testing.T) {
		for _, tt := range []testCase{
			{key: "db.host", wantVal: "localhost"},
			{key: "db.port", wantVal: "2222"},
			{key: "db.password", wantErr: true},
		} {
			t.Run(tt.key, func(t *testing.T) {
				AssertPanic(t, tt.wantErr, func() {
					val := MustGetString(tt.key)
					assert.Equal(t, tt.wantVal, val, "assert value of config")
				})
			})
		}
	})

	t.Run("MustGetBool", func(t *testing.T) {
		for _, tt := range []testCase{
			{key: "db.host", wantVal: false}, // no error, no default value, viper returns empty value for unmatched type
			{key: "db.ssl_mode", wantVal: true},
			{key: "db.test_var2", wantVal: false},
			{key: "db.test_var3", wantVal: true},
			{key: "db.password", wantErr: true},
		} {
			t.Run(tt.key, func(t *testing.T) {
				AssertPanic(t, tt.wantErr, func() {
					val := MustGetBool(tt.key)
					assert.Equal(t, tt.wantVal, val, "assert value of config")
				})
			})
		}
	})

	t.Run("MustGetInt", func(t *testing.T) {
		for _, tt := range []testCase{
			{key: "db.host", wantVal: 0}, // no error, no default value, viper returns empty value for unmatched type
			{key: "db.port", wantVal: 2222},
		} {
			t.Run(tt.key, func(t *testing.T) {
				AssertPanic(t, tt.wantErr, func() {
					val := MustGetInt(tt.key)
					assert.Equal(t, tt.wantVal, val, "assert value of config")
				})
			})
		}
	})

	t.Run("MustGetInt64", func(t *testing.T) {
		for _, tt := range []testCase{
			{key: "db.host", wantVal: int64(0)}, // no error, no default value, viper returns empty value for unmatched type
			{key: "db.port", wantVal: int64(2222)},
		} {
			t.Run(tt.key, func(t *testing.T) {
				AssertPanic(t, tt.wantErr, func() {
					val := MustGetInt64(tt.key)
					assert.Equal(t, tt.wantVal, val, "assert value of config")
				})
			})
		}
	})

	t.Run("MustGetFloat64", func(t *testing.T) {
		for _, tt := range []testCase{
			{key: "db.host", wantVal: float64(0)}, // no error, no default value, viper returns empty value for unmatched type
			{key: "db.port", wantVal: float64(2222)},
		} {
			t.Run(tt.key, func(t *testing.T) {
				AssertPanic(t, tt.wantErr, func() {
					val := MustGetFloat64(tt.key)
					assert.Equal(t, tt.wantVal, val, "assert value of config")
				})
			})
		}
	})
}

func AssertPanic(t *testing.T, expectPanic bool, f assert.PanicTestFunc, msgAndArgs ...interface{}) bool {
	panicAssert := assert.NotPanics
	if expectPanic {
		panicAssert = assert.Panics
	}

	return panicAssert(t, f, msgAndArgs...)
}
