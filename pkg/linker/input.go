package linker

import (
	"github.com/shangc1016/rvld/pkg/utils"
)

// 处理remaining中要链接的库文件
func ReadInputFiles(ctx *Context, remaining []string) {
	// 传入的是remaining数组,里面包含了
	// - 对象文件object
	// - 静态链接文件archive
	// - libxxx等静态链接库
	//
	// 本函数处理指定链接库的这种方式
	for _, arg := range remaining {

		if libname, ok := utils.RemovePrefix(arg, "-l"); ok {
			// 如果arg有一个"-l"的前缀，那就是通过指定链接库的名字的方式
			// 这儿需要根据静态链接库的名字，找到对应的xxx.a静态链接库文件
			ReadFile(ctx, FindLibraryByName(ctx, libname))
		} else {
			// 这种情况说明就是正常的objfile，或者是archive静态链接库文件
			// 直接读文件，解析object文件就行
			ReadFile(ctx, MustNewFile(arg))
		}
	}

}

func ReadFile(ctx *Context, file *File) {
	fileType := GetFileType(file.Contents)
	switch fileType {
	case FileTypeObject:
		// object file
		ctx.Objs = append(ctx.Objs, CreateObjectFile(file))
	case FileTypeArchive:
		// archive file
		for _, child := range ReadArchiveMembers(file) {
			// 确保每个文件的类型是objfile
			utils.Assert(GetFileType(child.Contents) == FileTypeObject)
			// 把从archive文件中解析出来的每个obj文件都加入到ctx中
			ctx.Objs = append(ctx.Objs, CreateObjectFile(child))
		}
	default:
		utils.Fatal("unknown file type")
	}
}

func CreateObjectFile(file *File) *ObjFile {
	obj := NewObjFile(file)
	// 解析objfile
	obj.Parse()
	return obj
}
