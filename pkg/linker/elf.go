package linker

import (
	"bytes"
	"strconv"
	"strings"
	"unsafe"

	"github.com/shangc1016/rvld/pkg/utils"
)

// elf header size
const EhdrSize = unsafe.Sizeof(EHdr{})

// section header size
const SHdrSize = unsafe.Sizeof(SHdr{})

// 符号表结构的大小
const SymSize = unsafe.Sizeof(Sym{})

// archive file header size
const ArHdrSize = unsafe.Sizeof(ArHdr{})

const ElfObjectMagicNumber = "\177ELF"
const ElfStaticLinkMagicNumber = "!<arch>\n"

// elf header struct
type EHdr struct {
	Ident     [16]uint8
	Type      uint16
	Machine   uint16
	Version   uint32
	Entry     uint64
	PhOff     uint64
	ShOff     uint64
	Flags     uint32
	EhSize    uint16
	PhEntSize uint16
	PhNum     uint16
	ShEntSize uint16
	ShNum     uint16
	ShStrNdx  uint16
}

// section header struct
type SHdr struct {
	Name      uint32 // 这个name是指再shstrtab这个sh中，当前sh名字的v起始偏移
	Type      uint32
	Flags     uint64
	Addr      uint64
	Offset    uint64
	Size      uint64
	Link      uint32
	Info      uint32
	AddrAlign uint64
	EntSize   uint64
}

// 符号表项
type Sym struct {
	Name  uint32
	Info  uint8
	Other uint8
	Shndx uint16
	Val   uint64
	Size  uint64
}

type ArHdr struct {
	Name [16]byte
	Date [12]byte
	Uid  [6]byte
	Gid  [6]byte
	Mode [8]byte
	Size [10]byte // 保存一个section中数据区域的大小
	Fmag [2]byte
}

func (a *ArHdr) HasPrefix(s string) bool {
	return strings.HasPrefix(string(a.Name[:]), s)
}

// string table
func (a *ArHdr) IsStrtab() bool {
	return a.HasPrefix("// ")
}

// symbol table
func (a *ArHdr) IsSymtab() bool {
	return a.HasPrefix("/ ") || a.HasPrefix("/SYM64/ ")
}

//
// 如果不是string table、也不是symbol table那就是objfile
//

func (a *ArHdr) ReadName(strTab []byte) string {
	// archive的objfile的name有两种。
	// long filename  :
	// short filename : 直接保存在每个section的arHdr的name字段

	if a.HasPrefix("/") {
		// 如果name以/开头，那么就是long filename，否则是short filename
		// long filename的话，是斜杠加上一个数字,后面都是空格。根据这个数字再次在strtab中进行索引。
		start, err := strconv.Atoi(strings.TrimSpace(string(a.Name[1:])))
		utils.MustNo(err)
		// start就是这个section的名字在strtab中的起始偏移位置
		// end结束位置就是在start的位置往后面找到"/\n"
		end := start + bytes.Index(strTab[start:], []byte("/\n"))
		return string(strTab[start:end])
	}

	// short filename
	// short filename 结束的位置是'/'
	end := strings.Index(string(a.Name[:]), "/")
	// 确保一定能找到'/'
	utils.Assert(end != -1)

	return string(a.Name[:end])
}

// 返回archive中一个section的size
func (a *ArHdr) GetSize() int {
	size, err := strconv.Atoi(strings.TrimSpace(string(a.Size[:])))
	utils.MustNo(err)
	return size
}

// 根据在shstrtab这个sh中的偏移量，返回相应sh的名字的字符串
func ElfGetName(strTab []byte, offset uint32) string {
	// 在strTab的offset偏移开始，返回第一个是0的位置
	len := bytes.Index(strTab[offset:], []byte{0})
	return string(strTab[offset : offset+uint32(len)])
}
