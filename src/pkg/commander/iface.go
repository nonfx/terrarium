package commander

import "os/exec"

type Commander interface {
	Run(*exec.Cmd) error
}
