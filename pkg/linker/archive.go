package linker

import (
	"github.com/shangc1016/rvld/pkg/utils"
)

func ReadArchiveMembers(file *File) []*File {
	// 确保是静态链接文件
	utils.Assert(GetFileType(file.Contents) == FileTypeArchive)

	//
	// archive文件是多个对象文件的集合。
	// 从archive文件中解析出多个obj文件

	// 静态链接文件xxx.a的格式如下
	// [!<arch>\n][ section ][ section ][ section ][ section ]
	// section 的大小是对2取整的。所以section之间可能会有margin
	// section里面的结构是[ arHdr ] [ 任意长度的数据 ]

	// 这个是archive文件中magic number的偏移
	// len(ElfStaticLinkMagicNumber) = 8
	pos := len(ElfStaticLinkMagicNumber)

	// 需要有一个专门的strTabl保存string table
	var strTab []byte
	var files []*File

	for len(file.Contents)-pos > 1 {
		// 每个section的起始地址会对2取整
		if pos%2 == 1 {
			pos++
		}

		arHdr := utils.Read[ArHdr](file.Contents[pos:])
		dataStart := pos + int(ArHdrSize)
		// pos后移[ section ]，也就是[ arHdr ][ contents... ]
		pos += int(ArHdrSize) + arHdr.GetSize()
		dataEnd := pos
		// 这是其中一段section的数据
		contents := file.Contents[dataStart:dataEnd]
		//

		if arHdr.IsSymtab() {
			continue
		} else if arHdr.IsStrtab() {
			//
			strTab = contents
			continue
		}

		// 既不是strtab、也不是symtab。那就是objfile
		files = append(files, &File{
			Filename: arHdr.ReadName(strTab),
			Contents: contents,
			Parent:   file,
		})

	}

	return files
}
