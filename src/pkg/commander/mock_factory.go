//go:build mock
// +build mock

package commander

var m Commander

func SetCommander(mock Commander) {
	m = mock
}

func GetCommander() Commander {
	return m
}
