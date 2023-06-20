package linker

import (
	"debug/elf"
)

type ObjFile struct {
	InputFile       // 结构体组合
	SymtabSec *SHdr // `symbol table`这个sh包含了所有sh的类型信息
}

func NewObjFile(file *File) *ObjFile {
	o := &ObjFile{
		InputFile: NewInputFile(file),
	}
	return o
}

func (o *ObjFile) Parse() {
	// 解析对象文件的符号表
	// 可以使用命令
	//
	// 这个parse是找到类型为符号表的一个section header
	//
	o.SymtabSec = o.FindfSection(uint32(elf.SHT_SYMTAB))
	if o.SymtabSec != nil {
		o.FirstGlobal = int64(o.SymtabSec.Info)
		// 根据符号表这一个sh，解析出来符号表信息
		o.FillUpElfSyms(o.SymtabSec)
		// 根据符号表的sh的link字段在符号表数字中的偏移拿到符号表
		// 然后根据拿到的符号表解析褚这个符号表对应的文件内容，写到symbolstrtab中
		o.SymbolStrtab = o.GetBytesFromIdx(int64(o.SymtabSec.Link))
	}
}
