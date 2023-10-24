// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

//go:build mock
// +build mock

package update

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/cldcvr/terrarium/src/cli/internal/config"
	"github.com/cldcvr/terrarium/src/pkg/db/mocks"
	"github.com/cldcvr/terrarium/src/pkg/testutils/clitesting"
	"github.com/stretchr/testify/mock"
	"gopkg.in/h2non/gock.v1"
)

func TestCmd(t *testing.T) {
	config.LoadDefaults()
	clitest := clitesting.CLITest{
		CmdToTest: NewCmd,
	}
	mockDB := &mocks.DB{}
	mockDB.On("ExecuteSQLStatement", mock.Anything).Return(func(dump string) error {
		if strings.Contains(dump, "fail") {
			return fmt.Errorf("mock error")
		}
		return nil
	})
	mockDB.On("CreateRelease", mock.Anything).Return(nil, nil)
	config.SetDBMocks(mockDB)
	clitest.RunTests(t, []clitesting.CLITestCase{
		{
			Name: "success",
			GockSetup: func(ctx context.Context, t *testing.T) {
				var buf bytes.Buffer
				gzWriter := gzip.NewWriter(&buf)
				gzWriter.Write([]byte("some dummy SQL command;"))
				gzWriter.Close()
				b, _ := json.Marshal(map[string]interface{}{"tag_name": "mock_tag"})
				gock.New("https://api.github.com").
					Get("/repos/cldcvr/terrarium-farm/releases/latest").
					Reply(http.StatusOK).
					Body(bytes.NewReader(b))

				gock.New("https://github.com").
					Get("/cldcvr/terrarium-farm/releases/download/mock_tag/cc_terrarium_data.sql.gz").
					Reply(http.StatusOK).
					Body(bytes.NewReader(buf.Bytes())).
					SetHeader("Content-Encoding", "gzip")
			},
		},
		{
			Name: "download artifact non ok http response",
			GockSetup: func(ctx context.Context, t *testing.T) {
				b, _ := json.Marshal(map[string]interface{}{"tag_name": "mock_tag"})
				gock.New("https://api.github.com").
					Get("/repos/cldcvr/terrarium-farm/releases/latest").
					Reply(http.StatusOK).
					Body(bytes.NewReader(b))
				gock.New("https://github.com").
					Get("/cldcvr/terrarium-farm/releases/download/mock_tag/cc_terrarium_data.sql.gz").
					Reply(http.StatusInternalServerError).
					SetHeader("Content-Encoding", "gzip")
			},
			ExpError: "HTTP request failed with status code 500",
			WantErr:  true,
		},
		{
			Name: "invalid artifact failure",
			GockSetup: func(ctx context.Context, t *testing.T) {
				b, _ := json.Marshal(map[string]interface{}{"tag_name": "mock_tag"})
				gock.New("https://api.github.com").
					Get("/repos/cldcvr/terrarium-farm/releases/latest").
					Reply(http.StatusOK).
					Body(bytes.NewReader(b))
				gock.New("https://github.com").
					Get("/cldcvr/terrarium-farm/releases/download/mock_tag/cc_terrarium_data.sql.gz").
					Reply(http.StatusOK)
			},
			ExpError: "error creating gzip reader",
			WantErr:  true,
		},
		{
			Name: "failure executing SQL statement",
			GockSetup: func(ctx context.Context, t *testing.T) {
				var buf bytes.Buffer
				gzWriter := gzip.NewWriter(&buf)
				gzWriter.Write([]byte("failure"))
				gzWriter.Close()
				b, _ := json.Marshal(map[string]interface{}{"tag_name": "mock_tag"})
				gock.New("https://api.github.com").
					Get("/repos/cldcvr/terrarium-farm/releases/latest").
					Reply(http.StatusOK).
					Body(bytes.NewReader(b))
				gock.New("https://github.com").
					Get("/cldcvr/terrarium-farm/releases/download/mock_tag/cc_terrarium_data.sql.gz").
					Reply(http.StatusOK).
					Body(bytes.NewReader(buf.Bytes())).
					SetHeader("Content-Encoding", "gzip")
			},
			ExpError: "error executing dump file: mock error",
			WantErr:  true,
		},
	})
}

func setupGockSuccess(ctx context.Context, t *testing.T) {
	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)
	gzWriter.Write([]byte("some dummy SQL command;"))
	gzWriter.Close()

	gock.New("https://github.com").
		Get("/cldcvr/terrarium-farm/releases/download/latest/cc_terrarium_data.sql.gz").
		Reply(http.StatusOK).
		Body(bytes.NewReader(buf.Bytes())).
		SetHeader("Content-Encoding", "gzip")
}

func setupGockFailure(ctx context.Context, t *testing.T) {
	gock.New("https://github.com").
		Get("/cldcvr/terrarium-farm/releases/download/latest/cc_terrarium_data.sql.gz").
		Reply(http.StatusInternalServerError).
		SetHeader("Content-Encoding", "gzip")
	return
}
