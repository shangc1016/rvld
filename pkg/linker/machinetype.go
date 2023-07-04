package linker

import (
	"debug/elf"

	"github.com/shangc1016/rvld/pkg/utils"
)

type MachineType = uint8

const (
	MachineTypeNone    MachineType = iota
	MachinetypeRISCV64 MachineType = iota
)

func GetMachineTypeFromContents(contents []byte) MachineType {
	ft := GetFileType(contents)
	switch ft {
	case FileTypeObject:
		// 先判断elf header的machine字段，是不是riscv
		machine := utils.Read[uint16](contents[18:])
		if machine == uint16(elf.EM_RISCV) {
			class := elf.Class(contents[4])
			switch class {
			// 再判断elf header的class字段，看是不是64位的
			case elf.ELFCLASS64:
				// 本项目只处理riscv64
				return MachinetypeRISCV64
			}
			return MachineTypeNone
		}
	}
	return MachineTypeNone
}

type MachineTypeStringer struct {
	MachineType
}

func (m MachineTypeStringer) String() string {
	switch m.MachineType {
	case MachinetypeRISCV64:
		return "riscv64"
	}
	utils.Assert(m.MachineType == MachineTypeNone)
	return "none"
}
