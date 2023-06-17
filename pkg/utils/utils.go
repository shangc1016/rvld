package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"runtime/debug"
)

func Fatal(v any) {
	fmt.Printf("rvld: \033[0;1;31mfatal:\033[0m %v\n", v)
	// 打印出错的调用栈
	debug.PrintStack()
	os.Exit(1)
}

// 保证一定不能出错
func MustNo(err error) {
	if err != nil {
		Fatal(err.Error())
	}
}

func Assert(condition bool) {
	if !condition {
		Fatal("assert failed")
	}
}

// template
func Read[T any](data []byte) (val T) {
	reader := bytes.NewReader(data)
	// 读文件的偏移怎么确定？
	err := binary.Read(reader, binary.LittleEndian, &val)
	MustNo(err)
	// 返回类型处已经指定了返回值
	return
}
