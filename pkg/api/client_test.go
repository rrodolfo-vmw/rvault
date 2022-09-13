package api_test

import (
	"testing"

	"github.com/rrodolfo-vmw/rvault/pkg/api"

	"github.com/spf13/viper"
)

var smokeTestToken = "mytoken"
var smokeTestAddress = "127.0.0.1:8200"

func TestNewClient(t *testing.T) {
	tests := []struct {
		name        string
		token       string
		address     string
		wantToken   string
		wantAddress string
		wantErr     bool
	}{
		{
			name:        "Smoke Test",
			token:       smokeTestToken,
			address:     smokeTestAddress,
			wantToken:   smokeTestToken,
			wantAddress: "http://" + smokeTestAddress,
			wantErr:     false,
		},
		{
			name:        "Fail on Missing Address",
			token:       smokeTestToken,
			address:     "",
			wantToken:   "",
			wantAddress: "",
			wantErr:     true,
		},
		{
			name:        "Fail on Missing Token",
			token:       "",
			address:     smokeTestAddress,
			wantToken:   "",
			wantAddress: "",
			wantErr:     true,
		},
		{
			name:        "Fail on invalid URL",
			token:       smokeTestToken,
			address:     ":invalid/",
			wantToken:   "",
			wantAddress: "",
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set("global.address", tt.address)
			viper.Set("global.token", tt.token)
			got, err := api.NewClient()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != nil {
				gotToken := got.Token()
				gotAddress := got.Address()
				if tt.wantToken != gotToken {
					t.Errorf("got = %s, want %s", gotToken, tt.wantToken)
				}
				if tt.wantAddress != gotAddress {
					t.Errorf("got = %s, want %s", gotAddress, tt.wantAddress)
				}
			}

		})
	}
}
