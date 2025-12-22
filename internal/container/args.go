package container

// InitArgs 是用匿名管道传输 init 函数参数时的类型
type InitArgs struct {
	Command   string   `json:"command"`
	Args      []string `json:"args"`
	EnableTty bool     `json:"enableTty"`
}
