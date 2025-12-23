package syncx_test

import (
	"context"
	"testing"
	"time"

	"github.com/jiny3/gopkg/syncx"
)

func TestListenLatest(t *testing.T) {
	testCounter := 0
	tests := []struct {
		name string
		fs   []func()
		want int
	}{
		{
			name: "Test case 1",
			fs: []func(){
				func() {
					testCounter++
					time.Sleep(600 * time.Millisecond)
				},
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			submit := syncx.ListenLatest(context.Background(), tt.fs...)
			for range 5 {
				submit()
				time.Sleep(100 * time.Millisecond)
			}
			time.Sleep(1 * time.Second)
			if testCounter != tt.want {
				t.Errorf("got %v, want %v", testCounter, tt.want)
			}
		})
	}
}
