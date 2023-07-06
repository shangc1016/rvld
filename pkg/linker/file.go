package linker

import (
	"os"

	"github.com/shangc1016/rvld/pkg/utils"
)

type File struct {
	Filename string
	Contents []byte
	Parent   *File
}

// 这个函数必须正确返回
func MustNewFile(filename string) *File {
	contents, err := os.ReadFile(filename)
	utils.MustNo(err)

	return &File{
		Filename: filename,
		Contents: contents,
	}
}

// 检验文件
func OpenLibrary(filepath string) *File {
	contents, err := os.ReadFile(filepath)
	if err != nil {
		return nil
	}
	return &File{
		Filename: filepath,
		Contents: contents,
	}
}

// 根据静态链接库的名字，找到链接库文件
func FindLibraryByName(ctx *Context, name string) *File {
	// 一般的规则是链接库的名字xxx，那么链接库文件名字就是libxxx.a
	for _, dir := range ctx.Args.LibraryPaths {
		// 查找链接库的路径都在ctx.Args.Libraries里面
		path := dir + "/lib" + name + ".a"
		if f := OpenLibrary(path); f != nil {
			return f
		}
	}
	utils.Fatal("library file not fount")
	return nil
}
