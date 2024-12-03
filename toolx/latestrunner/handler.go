package latestrunner

import (
	"fmt"
)

/*
 * 工程中经常有如下场景：
 * 配置在一段时间内间隔较短, 每次配置后需要进行全局耗时操作, 但只需要最新的配置生效后的全局耗时操作
 * 例如：topo修改后, 要生成新的链路图, 但是topo修改频繁, 生成链路图耗时, 只需要最新的topo生效后生成链路图
 * 效果: 传入一个函数, 短时间内多次调用, 只有最后一次(若近期该函数未被调用, 则还有第一次)调用会被执行, 从 O(n) 降到 O(1)
 */

type runner struct {
	queue     chan struct{}
	closer    chan struct{}
	listening bool
	f         func()
}

// 暂时不支持返回值
type Task struct {
	Func func(...any)
	Args []any
}

// 创建一个 runner, 传入多个 Task, 每个 Task 会被执行一次
func New(tasks ...Task) *runner {
	if len(tasks) < 1 {
		return &runner{
			queue:     make(chan struct{}, 60),
			closer:    make(chan struct{}),
			listening: false,
			f:         nil,
		}
	}
	fs := func() {
		for _, task := range tasks {
			task.Func(task.Args...)
		}
	}
	return &runner{
		queue:     make(chan struct{}, 60),
		closer:    make(chan struct{}),
		listening: false,
		f:         fs,
	}
}

// 监听 runner, 若不传入 Task, 则只监听初始化时传入的 Task, 若传入 Task, 则追加监听传入的 Task
func (r *runner) Listen(tasks ...Task) (func(), error) {
	if r.listening {
		return nil, fmt.Errorf("runner is listening")
	}
	if len(tasks) > 0 {
		hook := r.f
		fs := func() {
			if hook != nil {
				hook()
			}
			for _, task := range tasks {
				task.Func(task.Args...)
			}
		}
		r.f = fs
	}
	r.listening = true
	defer func() {
		r.listening = false
	}()
	go func() {
		for {
			select {
			case <-r.closer:
				return
			case <-r.queue:
				r.empty()
				if r.f != nil {
					r.f()
				}
			}
		}
	}()
	return r.run, nil
}

func (r *runner) Close() {
	if !r.listening {
		return
	}
	r.closer <- struct{}{}
	r.empty()
}

func (r *runner) run() {
	select {
	case r.queue <- struct{}{}:
		return
	default:
		r.empty()
		r.queue <- struct{}{}
	}
}

func (r *runner) empty() {
	// 清空 channel
	full := true
	for full {
		select {
		case <-r.queue:
			continue
		default:
			full = false
		}
	}
}
