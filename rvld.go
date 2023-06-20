package main

import (
	"fmt"
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

	// 继续解析objfile
	objfile := linker.NewObjFile(file)
	objfile.Parse()

	for _, sh := range objfile.ElfSections {
		fmt.Println(linker.ElfGetName(objfile.ShStrtab, sh.Name))
	}

	fmt.Println("section header lens: ", len(objfile.ElfSections))

	fmt.Println("firstGlobal: ", objfile.FirstGlobal)

	fmt.Println("symbol table size: ", len(objfile.ElfSyms))

	// 打印出符号表的
	for _, sym := range objfile.ElfSyms {
		fmt.Println(linker.ElfGetName(objfile.SymbolStrtab, sym.Name))
	}

}
