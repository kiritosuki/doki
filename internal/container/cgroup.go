package container

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// CgroupManager 用于管理 cgroup
type CgroupManager struct {
	Path string
}

// ResourceConfig 资源字段对象
type ResourceConfig struct {
	Memory string
	Cpus   float64
	Cpuset string
}

const base = "/sys/fs/cgroup/doki"

// NewCgroupManager 新建 CgroupManager 对象
func NewCgroupManager(containerId string) (*CgroupManager, error) {
	path := base + "/container-" + containerId
	// 0开头表示八进制 755是权限
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return nil, err
	}
	return &CgroupManager{Path: path}, nil
}

// EnableControllers 启动 controller
func EnableControllers() error {
	err := os.MkdirAll(base, 0755)
	if err != nil {
		return err
	}
	// 启用 controller
	data := []byte("+memory +cpu +cpuset")
	return os.WriteFile(base+"/cgroup.subtree_control", data, 0644)
}

// Enabled 判断是否开启了资源管理
func (r *ResourceConfig) Enabled() bool {
	return r.Memory != "" || r.Cpus > 0 || r.Cpuset != ""
}

// ParseMemory 格式化 memory 字符串
func (r *ResourceConfig) ParseMemory(memory string) (string, error) {
	if memory == "" {
		return "max", nil
	}
	// 转换小写
	memory = strings.ToLower(memory)
	var unit int64 = 1
	switch {
	case strings.HasSuffix(memory, "k"):
		unit = 1024
	case strings.HasSuffix(memory, "m"):
		unit = 1024 * 1024
	case strings.HasSuffix(memory, "g"):
		unit = 1024 * 1024 * 1024
	}
	memory = memory[:(len(memory) - 1)]
	// memory: 原数字	 10: 进制	 64: 位数
	// 比 Atoi 更高级的转换函数
	value, err := strconv.ParseInt(memory, 10, 64)
	if err != nil || value <= 0 {
		return "", fmt.Errorf("invalid memory limit : %s", memory)
	}
	// 参数: 数字 进制
	return strconv.FormatInt(value*unit, 10), nil
}
