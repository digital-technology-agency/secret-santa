package utils

import "testing"

var testKey = "SANTA_TEST_KEY"

func TestGetEnv(t *testing.T) {
	type args struct {
		key      string
		fallback string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "get env from os",
			args: args{
				key:      testKey,
				fallback: "Key not found",
			},
			want: "Key not found",
		},
		{
			name: "get env from os",
			args: args{
				key:      "GOPATH",
				fallback: "Key not found",
			},
			want: "/home/deemak/go/bin",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetEnv(tt.args.key, tt.args.fallback); got != tt.want {
				t.Errorf("GetEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}
