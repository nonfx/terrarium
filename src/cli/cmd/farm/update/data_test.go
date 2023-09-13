// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

//go:build mock
// +build mock

package update

import "testing"

func Test_cleanup(t *testing.T) {
	tests := []struct {
		name string
		dump string
		want string
	}{
		{
			name: "success",
			dump: "INSERT INTO public.taxonomies valid sql statement; \nSET some statement; \nSELECT some statement;",
			want: "INSERT INTO taxonomies valid sql statement; \n\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cleanup(tt.dump); got != tt.want {
				t.Errorf("cleanup() = %v, want %v", got, tt.want)
			}
		})
	}
}
