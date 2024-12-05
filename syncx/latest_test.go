package syncx_test

import (
	"testing"
	"time"

	"github.com/jiny3/gopkg/syncx"
)

func ExampleRunner(tasks ...func()) int {
	counter := 0
	counterTask := func() {
		counter++
	}
	lrunner := syncx.NewLatest(counterTask)

	runFunc, err := lrunner.Listen(tasks...)
	if err != nil {
		panic(err)
	}
	defer lrunner.Close()

	for i := 0; i < 10; i++ {
		runFunc()
		time.Sleep(1 * time.Millisecond)
	}
	time.Sleep(40 * time.Millisecond)
	return counter
}

func TestExampleRunner(t *testing.T) {
	type args struct {
		tasks []func()
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "test",
			args: args{
				tasks: []func(){
					func() {
						time.Sleep(30 * time.Millisecond)
					},
				},
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExampleRunner(tt.args.tasks...); got != tt.want {
				t.Errorf("ExampleRunner() = %v, want %v", got, tt.want)
			}
		})
	}
}
