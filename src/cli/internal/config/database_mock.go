// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

//go:build mock
// +build mock

package config

import (
	"github.com/cldcvr/terrarium/src/pkg/db/mocks"
	"github.com/spf13/viper"
)

func SetDBMocks(m *mocks.DB) {
	EnableDBMocks()
	mockdb = m
}

func EnableDBMocks() {
	viper.Set("db.type", "mock")
}
