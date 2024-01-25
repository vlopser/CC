package main

import "testing"

func Test_checkQueue(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "testOK"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkQueue()
		})
	}
}
