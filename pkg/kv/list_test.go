package kv_test

import (
	"reflect"
	"testing"

	"github.com/rrodolfo-vmw/rvault/pkg/kv"

	vapi "github.com/hashicorp/vault/api"
	"github.com/spf13/viper"
)

func TestRList(t *testing.T) {
	cluster := createTestVault(t)
	defer cluster.Cleanup()
	client := cluster.Cores[0].Client
	type args struct {
		c            *vapi.Client
		engine       string
		kvVersion    string
		path         string
		concurrency  uint32
		includePaths []string
		excludePaths []string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "Smoke Test V2",
			args: args{
				c:            client,
				engine:       engineV2,
				kvVersion:    "2",
				path:         "/",
				includePaths: []string{"*"},
			},
			want: []string{
				"/france/paris/key",
				"/spain/admin",
				"/spain/malaga/random",
				"/uk/london/mi5",
			},
			wantErr: false,
		},
		{
			name: "Smoke Test V1",
			args: args{
				c:            client,
				engine:       engine,
				kvVersion:    "1",
				path:         "/",
				includePaths: []string{"*"},
			},
			want: []string{
				"/france/paris/key",
				"/spain/admin",
				"/spain/malaga/random",
				"/uk/london/mi5",
			},
			wantErr: false,
		},
		{
			name: "Inclusion and Exclusion Paths",
			args: args{
				c:            client,
				engine:       engine,
				kvVersion:    "1",
				path:         "/",
				includePaths: []string{"/spain/*", "/uk/*"},
				excludePaths: []string{"*/admin"},
			},
			want: []string{
				"/spain/malaga/random",
				"/uk/london/mi5",
			},
			wantErr: false,
		},
		{
			name: "Secret not found",
			args: args{
				c:            client,
				engine:       engineV2,
				kvVersion:    "2",
				path:         "/france/fakesecret",
				includePaths: []string{"*"},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Unknown KV version",
			args: args{
				c:            client,
				engine:       engine,
				kvVersion:    "99",
				path:         "/",
				includePaths: []string{"*"},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := kv.RList(tt.args.c, tt.args.engine, tt.args.kvVersion, tt.args.path, tt.args.includePaths, tt.args.excludePaths,
				tt.args.concurrency)
			viper.Reset()
			if (err != nil) != tt.wantErr {
				t.Errorf("RList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RList() got = %v, want %v", got, tt.want)
			}
		})
	}
}
