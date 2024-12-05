package syncx

import (
	"fmt"
)

/*
 * 工程中经常有如下场景：
 * 配置在一段时间内间隔较短, 每次配置后需要进行全局耗时操作, 但只需要最新的配置生效后的全局耗时操作
 * 例如：topo修改后, 要生成新的链路图, 但是topo修改频繁, 生成链路图耗时, 只需要最新的topo生效后生成链路图
 * 效果: 传入一个函数, 短时间内多次调用, 只有最后一次(若近期该函数未被调用, 则还有第一次)调用会被执行, 从 O(n) 降到 O(1)
 */

type latest struct {
	queue     chan struct{}
	closer    chan struct{}
	listening bool
	f         func()
}

// 创建一个 runner, 传入多个 func(), 顺序执行; 无入参, 无返回值, 推荐传入闭包函数
func NewLatest(fs ...func()) *latest {
	if len(fs) < 1 {
		return &latest{
			queue:     make(chan struct{}, 60),
			closer:    make(chan struct{}),
			listening: false,
			f:         nil,
		}
	}
	_fs := func() {
		for _, f := range fs {
			f()
		}
	}
	return &latest{
		queue:     make(chan struct{}, 60),
		closer:    make(chan struct{}),
		listening: false,
		f:         _fs,
	}
}

// 监听 runner, 若不传入 func(), 则只执行初始化时传入的 func(), 若传入 func(), 则追加执行传入的 func()
func (r *latest) Listen(fs ...func()) (func(), error) {
	if r.listening {
		return nil, fmt.Errorf("runner is listening")
	}
	if len(fs) > 0 {
		hook := r.f
		_fs := func() {
			if hook != nil {
				hook()
			}
			for _, f := range fs {
				f()
			}
		}
		r.f = _fs
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

func (r *latest) Close() {
	if !r.listening {
		return
	}
	r.closer <- struct{}{}
	r.empty()
}

func (r *latest) run() {
	select {
	case r.queue <- struct{}{}:
	default:
	}
}

func (r *latest) empty() {
	// 清空 channel
	for len(r.queue) > 0 {
		<-r.queue
	}
}
