package container

import (
	"os"
	"syscall"

	log "github.com/sirupsen/logrus"
)

// RunContainerInitProcess 启动容器的 init 进程
// 容器所在进程创建之后 这是容器内部执行的第一个进程
// 使用 mount 挂载 proc 文件系统
// Linux 中的文件 有真实文件 有虚拟文件
// 真实目录挂载 也可以用 mount source为源路径 target为目标路径
// proc是虚拟文件系统 linux系统下的/proc目录不是真实存在的 是内核根据进程状态动态维护的一个虚拟目录
// 挂载虚拟目录时 source只需要写标识符即可 这里的标识符就是proc 挂载到容器内的/proc
func RunContainerInitProcess(command string, args []string) error {
	log.Infof("command: %s", command)
	// 这里当source和文件类型都为 “” 的时候 表示修改挂载配置
	// 把该namespace的/目录之下的传播机制设置为private
	// 在新版本的linux内核中 由于有mount传播机制 子进程的mount会被传播的父进程
	// 如果不加这一行 创建容器进程的namespace执行了mount之后 会影响到父进程(即宿主机)
	// 例如这里 宿主机的/proc也会被污染 当容器关停之后 宿主机的/proc就存放的是脏数据了
	// private 表示禁止传播 rec 表示递归对 / 目录下所有生效
	_ = syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
	// 这里是 mount 的flags
	defaultMountFlags := syscall.MS_NOEXEC | // 禁止执行二进制文件
		syscall.MS_NOSUID | // 禁止su来提升权限
		syscall.MS_NODEV // 禁止创建文件
	_ = syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	argv := []string{command}
	if err := syscall.Exec(command, argv, os.Environ()); err != nil {
		log.Errorf(err.Error())
	}
	return nil
}
