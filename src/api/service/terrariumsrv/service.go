// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package terrariumsrv

import (
	"github.com/cldcvr/terrarium/src/pkg/db"
)

type Service struct {
	db db.DB
}

func New(db db.DB) *Service {
	return &Service{db: db}
}
