package testutils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	OldToken       = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiIiLCJlbWFpbCI6InRlc3QiLCJwcm92aWRlclRva2VuIjoiZXlKaGJHY2lPaUpCTWpVMlIwTk5TMWNpTENKbGJtTWlPaUpCTVRJNFIwTk5JaXdpYVhZaU9pSm9XVU4xWmxsZmNrMWhheTB5UTBzMklpd2lkR0ZuSWpvaU1uRjFhMHhCVVZoWU1uTXlVRU4yY1hGSGNESmZRU0lzSW5wcGNDSTZJa1JGUmlKOS52T3JLcVdiMmZqWTF5VThGenJHTkRRLlc0N1E3MkhxMzUydy1obnguRXY4VjBBSS5ONVUyS1RrY3J4RWV1ZDVBRjNSbXJ3IiwiZ2NUb2tlbiI6InRlc3QiLCJleHAiOjE2MDQxMjIwMDAsImlhdCI6MTYwNDAzNTYwMH0.OqhloDOpk-QTPdi8yBeb3XJ1-RcQmotqX5i21roNYfY"
	invalidSigning = "eyJhbGciOiJIUzM4NCIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiIiLCJlbWFpbCI6InRlc3QiLCJwcm92aWRlclRva2VuIjoiZXlKaGJHY2lPaUpCTWpVMlIwTk5TMWNpTENKbGJtTWlPaUpCTVRJNFIwTk5JaXdpYVhZaU9pSTNPRWhKYVRoSVltcGFWWEowY1RkbElpd2lkR0ZuSWpvaVZrcDVVMWs1WXpaQ1ZGQXpSbVJxWm0xWlFuWnNVU0lzSW5wcGNDSTZJa1JGUmlKOS5VejR0Y2wyU05NX1NCQXIta0hVOGxnLmVWM09xT3NiZ2ppQVlFYVEuN2tCYllTTS5yS1ZILWZqWUtxdkdiWktUOFpUNmJ3IiwiZ2NUb2tlbiI6InRlc3QiLCJleHAiOjE2MDQzODE1NDIsImlhdCI6MTYwNDI5NTE0Mn0.ugIfb_H7J1JOpsZlw8VAKNntIVEpUlSQsApbrVwOnJnXvQ6tfg1CS4Qv2-KMK9H9"
)

func TestCreateToken(t *testing.T) {
	type args struct {
		key           string
		userID        string
		email         string
		gcToken       string
		providerToken string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Fail",
			args: args{
				email:   "test",
				gcToken: "test",
			},
			wantErr: true,
		},
		{
			name: "Success",
			args: args{
				key:     "mysupersecretkey",
				email:   "test",
				gcToken: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateToken(tt.args.key, tt.args.userID, tt.args.email, tt.args.gcToken, tt.args.providerToken, time.Now())
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Equal(t, 488, len(got))
			}
		})
	}
}

func TestCreateRefreshToken(t *testing.T) {
	type args struct {
		userID      string
		jwtKey      string
		currentTime time.Time
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				jwtKey:      "mysupersecretkey",
				userID:      "u",
				currentTime: time.Now(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CreateRefreshToken(tt.args.userID, tt.args.jwtKey, tt.args.currentTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateRefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
