package linker

import (
	"os"

	"github.com/shangc1016/rvld/pkg/utils"
)

type File struct {
	Filename string
	Contents []byte
}

// 这个函数必须正确返回
func MustNewFile(filename string) *File {
	contents, err := os.ReadFile(filename)
	utils.MustNo(err)

	return &File{
		Filename: filename,
		Contents: contents,
	}
}
