package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/shangc1016/rvld/pkg/linker"
	"github.com/shangc1016/rvld/pkg/utils"
)

var version string

func main() {

	// ctx接收所有链接器的参数
	ctx := linker.NewContext()
	// remaining里面是-lxxx这种的动态链接，链接库的信息
	// 或者是其他的链接文件xxx.o这种格式
	remaining := parseArgs(ctx)

	// 如果参数中没有解析出machinetype，也就是没有指定链接器生成可执行文件的架构
	if ctx.Args.Emulation == linker.MachineTypeNone {
		// 如果参数中没有指定machineType,
		// 需要从第一个对象文件中解析出来
		for _, filename := range remaining {
			// 这种情况就是去除了-lxxx动态库的情况
			if strings.HasPrefix(filename, "-") {
				continue
			}
			file := linker.MustNewFile(filename)
			ctx.Args.Emulation = linker.GetMachineTypeFromContents(file.Contents)

			// 只要emulation有值，直接break
			if ctx.Args.Emulation != linker.MachineTypeNone {
				break
			}
		}
	}

	// 判断传给rvld的参数中的machine硬件架构是不是riscv64
	if ctx.Args.Emulation != linker.MachinetypeRISCV64 {
		utils.Fatal("unknown emulation type")
	}

	linker.ReadInputFiles(ctx, remaining)

	fmt.Println(len(ctx.Objs))

}

func parseArgs(ctx *linker.Context) []string {

	args := os.Args[1:]

	// 根据参数名字,生成可能的参数匹配项
	dashes := func(name string) []string {
		if len(name) == 1 {
			return []string{"-" + name}
		}
		return []string{"-" + name, "--" + name}
	}

	// 这个地方为啥要这么写,arg放在外面很丑
	// 我想是因为,在遍历args的时候,readArgs,readFlag需要相同的返回值
	arg := ""
	// 解析args,这种是有取值的
	readArgs := func(name string) bool {
		for _, opt := range dashes(name) {
			if args[0] == opt {
				// 因为还要拿到这个参数名的取值
				// 所以要判断一下args的长度，最少是2
				if len(args) == 1 {
					utils.Fatal(fmt.Sprintf("option -%s: argument missing", name))
				}
				arg = args[1]
				args = args[2:]
				return true
			}

			prefix := opt
			// 这儿的判断是为了处理形如-plugin-opt=xxx
			// 参数名字和值之间有一个=等于号的情况
			if len(name) > 1 {
				prefix += "="
			}

			// 这块是处理-melf64lriscv这种参数的解析
			// 参数符号和参数值是连起来的
			if strings.HasPrefix(args[0], prefix) {
				arg = args[0][len(prefix):]
				args = args[1:]
				return true
			}
		}
		return false
	}

	// 解析flags,这种是开关类型的
	readFlag := func(name string) bool {
		for _, opt := range dashes(name) {
			if args[0] == opt {
				args = args[1:]
				return true
			}
		}
		return false
	}

	//
	remaining := make([]string, 0)

	// 解析传递给rvld的参数
	for len(args) > 0 {

		if readFlag("help") {
			fmt.Println("Usage: .....")
			os.Exit(0)
		}

		if readArgs("o") || readArgs("output") {
			// flag about `output`
			ctx.Args.Output = arg
		} else if readFlag("v") || readFlag("version") {
			// flag about `version`
			fmt.Println(version)
			os.Exit(0)
		} else if readArgs("m") {
			// -m<arch-about> 指定架构
			if arg == "elf64lriscv" {
				// 只支持riscv64架构
				ctx.Args.Emulation = linker.MachinetypeRISCV64
			} else {
				utils.Fatal(fmt.Sprintf("unknown -m argument: %s", arg))
			}
		} else if readArgs("L") {
			// -L 指定库文件的路径
			ctx.Args.LibraryPaths = append(ctx.Args.LibraryPaths, arg)
		} else if readArgs("l") {
			// -l参数指定的是动态库，本项目不涉及动态库的链接
			remaining = append(remaining, "-l"+arg)
		} else if readArgs("sysroot") || readArgs("plugin") || readArgs("plugin-opt") ||
			readArgs("hash-style") || readArgs("build-id") || readFlag("static") ||
			readFlag("as-needed") || readFlag("start-group") || readFlag("end-group") ||
			readFlag("s") || readFlag("no-relax") {
			// 本项目不处理这些参数，直接忽略掉
		} else {
			if args[0][0] == '-' {
				utils.Fatal(fmt.Sprintf("unknown command line option: %s", args[0]))
			}
			remaining = append(remaining, args[0])
			args = args[1:]
		}

	}
	return remaining
}
