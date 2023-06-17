package linker

import "bytes"

// 检查文件内容的前四个字符
func CheckMagic(contents []byte) bool {
	return bytes.HasPrefix(contents, []byte("\177ELF"))
}
