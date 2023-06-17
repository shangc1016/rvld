package linker

import (
	"github.com/shangc1016/rvld/pkg/utils"
)

// 链接文件的基类
type InputFile struct {
	File        *File
	ElfSections []SHdr // section数组
}

func NewInputFile(file *File) InputFile {
	f := InputFile{File: file}

	// check file len not less than Ehdr's size
	if len(file.Contents) < int(EhdrSize) {
		utils.Fatal("file too small")
	}

	// check magic number
	if !CheckMagic(file.Contents) {
		utils.Fatal("not an ELF format file")
	}

	// read elf header from file
	ehdr := utils.Read[EHdr](f.File.Contents)

	// 切片后移，指向了elf文件的section段的起始地址
	contents := file.Contents[ehdr.ShOff:]

	shdr := utils.Read[SHdr](contents)

	// 在elf header中，ShOff表示section header区域在elf文件中的偏移
	// ShNum表示section区域中有几个section header
	// ShNum的类型是uint16，可能可执行文件的section header比较多，
	// 超过uint16了。这种情况需要ShNum是0，需要先读入第一个section header
	// 然后够再从这个section header中读到真正的section num。
	// 因为section header的里面的size类型是uint64位的，足够大。
	numSections := uint64(ehdr.ShNum)
	if numSections == 0 {
		numSections = shdr.Size
	}

	// 先把第一个shdr放进sction header的数组中
	f.ElfSections = []SHdr{shdr}
	for numSections > 1 {
		// 相当于contents每次往后移动ShdrSize，为了读入下一个Shdr
		contents = contents[SHdrSize:]
		// 然后把这个section header放到文件的section header数组中
		f.ElfSections = append(f.ElfSections, utils.Read[SHdr](contents))
		numSections--
	}

	return f
}
