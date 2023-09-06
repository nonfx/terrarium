// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package testutils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/Netflix/go-expect"
)

func GetTempFileName(dir string, prefix string, suffix string) (string, error) {
	tmpFile, err := os.CreateTemp(dir, prefix)
	if err != nil {
		return "", err
	}
	tmpFile.Close()
	os.Remove(tmpFile.Name())
	return fmt.Sprintf("%s.%s", tmpFile.Name(), suffix), nil
}

func SendInterrupt(c *expect.Console) {
	time.Sleep(time.Millisecond * 100)
	c.SendLine(string(terminal.KeyInterrupt))
	c.SendLine(string(terminal.KeyInterrupt))
}

func GenerateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}
	return string(ret), nil
}
