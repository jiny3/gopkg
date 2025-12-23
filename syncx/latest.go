package syncx

import (
	"context"
)

/*
 * 工程中经常有如下场景：
 * 配置在一段时间内间隔较短, 每次配置后需要进行全局耗时操作, 但只需要最新的配置生效后的全局耗时操作
 * 例如：topo修改后, 要生成新的链路图, 但是topo修改频繁, 生成链路图耗时, 只需要最新的topo生效后生成链路图
 * 效果: 传入一个函数, 短时间内多次调用, 只有最后一次(若近期该函数未被调用, 则还有第一次)调用会被执行, 从 O(n) 降到 O(1)
 */

// ListenLatest 传入多个 func(), 顺序执行，通过 submit() 提交执行请求
func ListenLatest(ctx context.Context, fs ...func()) (submit func()) {
	triggered := make(chan struct{}, 1)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-triggered:
				for _, f := range fs {
					f()
				}
			}
		}
	}()

	return func() {
		select {
		case triggered <- struct{}{}:
		default:
		}
	}
}
