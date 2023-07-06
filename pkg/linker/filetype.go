package linker

import (
	"debug/elf"

	"github.com/shangc1016/rvld/pkg/utils"
)

type FileType uint16

const (
	FileTypeUnknown FileType = iota
	FileTypeEmpty   FileType = iota
	FileTypeObject  FileType = iota // 对象文件   xxx.o
	FileTypeArchive FileType = iota // 静态链接库文件 xxx.a
)

// 首先判断是不是elf格式的文件，
// 然后判断elf header中的type字段
func GetFileType(contents []byte) FileType {
	if len(contents) == 0 {
		return FileTypeEmpty
	}
	if CheckMagic(contents, ElfObjectMagicNumber) {
		// 从elf header的第16个字节开始读，读一个uint16
		// 就是读入elf header中的type
		elfType := elf.Type(utils.Read[uint16](contents[16:]))
		switch elfType {
		case elf.ET_REL:
			return FileTypeObject
		}
		return FileTypeUnknown
	}

	// xxx.a 静态链接文件的魔数是"!<arch>\n"
	if CheckMagic(contents, ElfStaticLinkMagicNumber) {
		return FileTypeArchive
	}

	return FileTypeUnknown
}
