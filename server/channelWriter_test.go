package server

import (
	"reflect"
	"testing"
)

func Test_newChannelWriter(t *testing.T) {
	tests := []struct {
		name string
		want *channelWriter
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newChannelWriter(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newChannelWriter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_channelWriter_addOutputChannel(t *testing.T) {
	type args struct {
		out chan string
	}
	tests := []struct {
		name string
		cw   *channelWriter
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cw.addOutputChannel(tt.args.out)
		})
	}
}

func Test_channelWriter_start(t *testing.T) {
	tests := []struct {
		name string
		cw   *channelWriter
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cw.start()
		})
	}
}
