//go:build !mock
// +build !mock

package commander

func GetCommander() Commander {
	return &osExec{}
}
