package linker

import (
	"bytes"
	"unsafe"
)

// elf header size
const EhdrSize = unsafe.Sizeof(EHdr{})

// section header size
const SHdrSize = unsafe.Sizeof(SHdr{})

// 符号表结构的大小
const SymSize = unsafe.Sizeof(Sym{})

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

// 根据在shstrtab这个sh中的偏移量，返回相应sh的名字的字符串
func ElfGetName(strTab []byte, offset uint32) string {
	// 在strTab的offset偏移开始，返回第一个是0的位置
	len := bytes.Index(strTab[offset:], []byte{0})
	return string(strTab[offset : offset+uint32(len)])
}
