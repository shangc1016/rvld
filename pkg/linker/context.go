package linker

type ContextArgs struct {
	Output       string      // 保存链接器最后生成的文件名字
	Emulation    MachineType // 哪种架构
	LibraryPaths []string    // 库文件的路径，这是个数组
}

type Context struct {
	Args ContextArgs
}

func NewContext() *Context {
	// 链接器默认输出文件名字是a.out
	return &Context{Args: ContextArgs{
		Output:    "a.out",
		Emulation: MachineTypeNone}}
}
