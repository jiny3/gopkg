package latestrunner_test

import (
	"testing"
	"time"

	"github.com/jiny3/gopkg/toolx/latestrunner"
)

func ExampleRunner(tasks ...latestrunner.Task) int {
	counter := 0
	counterTask := latestrunner.Task{
		Func: func(args ...any) {
			counter++
		},
	}
	runner := latestrunner.New(counterTask)

	runFunc, err := runner.Listen(tasks...)
	if err != nil {
		panic(err)
	} else {
		defer runner.Close()
	}

	for i := 0; i < 10; i++ {
		runFunc()
		time.Sleep(1 * time.Millisecond)
	}
	time.Sleep(40 * time.Millisecond)
	return counter
}

func TestExampleRunner(t *testing.T) {
	type args struct {
		tasks []latestrunner.Task
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "test",
			args: args{
				tasks: []latestrunner.Task{
					{
						Func: func(args ...any) {
							time.Sleep(30 * time.Millisecond)
						},
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
