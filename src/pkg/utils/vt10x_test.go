// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/Netflix/go-expect"
	"github.com/stretchr/testify/assert"
)

var failNewConsole bool

func expectNewConsoleMock(opts ...expect.ConsoleOpt) (*expect.Console, error) {
	if failNewConsole {
		return nil, fmt.Errorf("error from NewConsole")
	} else {
		return expect.NewConsole(opts...)
	}
}

func TestNewVT10XConsole(t *testing.T) {
	expectNewConsole = expectNewConsoleMock
	defer func() { expectNewConsole = expect.NewConsole }()

	// This timeout will make it so the test fails within 5 seconds if something is wrong
	timeoutOpt := expect.WithDefaultTimeout(time.Duration(5) * time.Second)
	outputBuffer := new(bytes.Buffer)
	type args struct {
		opts []expect.ConsoleOpt
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid console",
			args: args{opts: []expect.ConsoleOpt{expect.WithStdout(outputBuffer)}},
		},
		{
			name:    "Fail - Console error",
			args:    args{opts: []expect.ConsoleOpt{}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			failNewConsole = tt.wantErr
			console, err := NewVT10XConsole(append(tt.args.opts, timeoutOpt)...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewVT10XConsole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			doneConsole := make(chan struct{})
			go func() {
				defer close(doneConsole)
				_, err := console.ExpectString("testing 1 2 3")
				assert.Nil(t, err)
			}()

			console.SendLine("testing 1 2 3")
			<-doneConsole
			console.Tty().Close()
			console.Close()
		})
	}
}
