package utils

import (
	"context"
	"github.com/sirupsen/logrus"
	"runtime"
	"sync"
	"sync/atomic"
)

func ExecAndCountFuncCtx(ctx context.Context, funcs ...func() error) int32 {
	wg := sync.WaitGroup{}
	errNum := int32(0)
	wg.Add(len(funcs))

	for _, fc := range funcs {
		go func(f func() error) {
			defer func() {
				wg.Done()
				if panicErr := recover(); panicErr != nil {
					const size = 4096
					buf := make([]byte, size)
					buf = buf[:runtime.Stack(buf, false)]
					logrus.Errorf("panic: %v\n %s\n", panicErr, buf)
				}
			}()

			if err := f(); err != nil {
				atomic.AddInt32(&errNum, 1)
			}
		}(fc)
	}

	wg.Wait()
	return errNum
}

func MaxLoopController(ctx context.Context, maxLoop int64, logic func() (bool, error)) {
	currLoop := int64(0)
	for {
		logrus.Infof("looping on %d", currLoop)
		if currLoop > maxLoop {
			logrus.Errorf("循环超过上线：次数=%d", currLoop)
			break
		}
		currLoop += 1
		isEnd, err := logic()
		if isEnd || err != nil {
			return
		}
	}
}
