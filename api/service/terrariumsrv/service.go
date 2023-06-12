package terrariumsrv

import (
	"github.com/cldcvr/terrarium/api/db"
)

type Service struct {
	db db.DB
}

func New(db db.DB) *Service {
	return &Service{db: db}
}
