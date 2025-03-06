package hookx

import (
	"fmt"
	"testing"
)

func TestInit(t *testing.T) {
	myTestHook := func() {
		fmt.Println("Hook 1")
	}
	type args struct {
		hooks []*func()
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test Init",
			args: args{
				hooks: []*func(){
					&myTestHook,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			myTestHook := func() {
				fmt.Println("myTestHook")
			}
			Init(&myTestHook)
			Init(&myTestHook)
			Init(tt.args.hooks...)
			Init(tt.args.hooks...)
		})
	}
}
