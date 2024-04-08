package rpc_test

import (
	"golsp/rpc"
	"testing"
)

type EncodingExample struct {
	Testing bool
}

func TestEncodeMessage(t *testing.T) {
	type args struct {
		msg any
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestEncodeMessage 1",
			args: args{msg: EncodingExample{Testing: true}},
			want: "Content-Length: 16\r\n\r\n{\"Testing\":true}",
		},
		{
			name: "TestEncodeMessage 2",
			args: args{msg: EncodingExample{Testing: false}},
			want: "Content-Length: 17\r\n\r\n{\"Testing\":false}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rpc.EncodeMessage(tt.args.msg); got != tt.want {
				t.Errorf("EncodeMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
