package commander

import "os/exec"

type osExec struct{}

func (e *osExec) Run(cmd *exec.Cmd) error {
	return cmd.Run()
}
