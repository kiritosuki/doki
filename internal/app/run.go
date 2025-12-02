package app

import (
	"os"

	"github.com/kiritosuki/doki/internal/container"
	log "github.com/sirupsen/logrus"
)

// Run 用来执行具体的 command
func Run(tty bool, command string) {
	// 这里NewParentProcess 相当于 fork 了一个进程 作为 init 进程
	// init 进程是所有进程的祖先进程 命名为parent
	parent := container.NewParentProcess(tty, command)
	// Start 方法表示 init 进程真正调用
	if err := parent.Start(); err != nil {
		log.Error(err)
	}
	// 容器运行时代码会阻塞在这里 因为祖先线程有义务回收所有子进程 避免僵尸进程
	_ = parent.Wait()
	// 当所有子进程都被回收之后 表示运行结束 容器推出
	// unix 规范 0表示正常退出 非0表示异常退出
	// 这里的逻辑是当所有子进程都结束之后 容器自动关闭
	// init 进程依赖于子进程的生命周期 故这里标记为 -1 异常退出 也是常见做法
	os.Exit(-1)
}
