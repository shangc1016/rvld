package linker

import "unsafe"

// elf header size
const EhdrSize = unsafe.Sizeof(EHdr{})

// section header size
const SHdrSize = unsafe.Sizeof(SHdr{})

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
	Name      uint32
	Type      uint32
	Flags     uint64
	Addr      uint64
	Offset    uint64
	Size      uint64
	Link      uint32
	AddrAlign uint64
	EntSize   uint64
}
