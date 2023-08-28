package commander

import "os/exec"

//go:generate mockery --all

type Commander interface {
	Run(*exec.Cmd) error
}
