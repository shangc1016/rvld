package main

import (
	"os"

	"github.com/shangc1016/rvld/pkg/linker"
	"github.com/shangc1016/rvld/pkg/utils"
)

func main() {
	if len(os.Args) < 2 {
		utils.Fatal("wrong args")
	}

	// 从命令行的第一个参数解析对象文件名
	file := linker.MustNewFile(os.Args[1])

	// 按照对象文件的elf格式解析出文件的section header
	inputfile := linker.NewInputFile(file)
	utils.Assert(len(inputfile.ElfSections) == 10)

}
