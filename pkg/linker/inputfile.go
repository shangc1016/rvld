package linker

import (
	"debug/elf"
	"fmt"

	"github.com/shangc1016/rvld/pkg/utils"
)

// 链接文件的基类
type InputFile struct {
	File         *File
	ElfSections  []SHdr // section数组
	ShStrtab     []byte // 对象文件的strtab sh
	FirstGlobal  int64
	ElfSyms      []Sym  // 符号表数组
	SymbolStrtab []byte // 符号表的str table
}

func NewInputFile(file *File) InputFile {
	f := InputFile{File: file}

	// check file len not less than Ehdr's size
	if len(file.Contents) < int(EhdrSize) {
		utils.Fatal("file too small")
	}

	// check magic number
	if !CheckMagic(file.Contents, ElfObjectMagicNumber) {
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

	// ehdr.ShStrndx这个数字类型是uint16，和上面表示范围一样的问题
	shstrndx := uint32(ehdr.ShStrNdx)
	if ehdr.ShStrNdx == uint16(elf.SHN_XINDEX) {
		// 如果elf头中的shstrndx这个字段为0，那么表示shstrndx这个sh的下标索引在第一个sh的link字段里
		shstrndx = shdr.Link
	}

	// 拿到了对象文件的shstrtab这个sh的字符数组
	f.ShStrtab = f.GetBytesFromIdx(int64(shstrndx))

	// fmt.Printf("%s\n", f.ShStrtab)

	return f
}

// 根据secion header拿到对应的section的内容
func (f *InputFile) GetBytesFromShdr(s *SHdr) []byte {
	end := s.Offset + s.Size
	if uint64(len(f.File.Contents)) < end {
		utils.Fatal(fmt.Sprintf("section header is out of range: %d", s.Offset))
	}
	return f.File.Contents[s.Offset : s.Offset+s.Size]
}

// 根据section header在数组中的偏移，拿到section的内容
func (f *InputFile) GetBytesFromIdx(idx int64) []byte {
	return f.GetBytesFromShdr(&f.ElfSections[idx])
}

// 填充符号表数组 symbol table
func (f *InputFile) FillUpElfSyms(s *SHdr) {
	bs := f.GetBytesFromShdr(s)
	// 计算得到symbol table的数量
	nums := len(bs) / int(SymSize)

	f.ElfSyms = make([]Sym, 0, nums)
	// 把每个symbol  table
	for nums > 0 {
		f.ElfSyms = append(f.ElfSyms, utils.Read[Sym](bs))
		// bs向后移动一个symsize的大小
		bs = bs[SymSize:]
		nums--
	}
}

// 在对象文件的sh数组中，找到相应类型的sh
func (f *InputFile) FindfSection(typ uint32) *SHdr {
	for i := 0; i < len(f.ElfSections); i++ {
		shdr := f.ElfSections[i]
		if shdr.Type == typ {
			return &shdr
		}
	}
	return nil
}
