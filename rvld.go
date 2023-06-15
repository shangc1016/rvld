package main

import (
	"fmt"
	"os"

	"github.com/shangc1016/rvld/pkg/utils"
)

func main() {
	if len(os.Args) < 2 {
		utils.Fatal("wrong args")
	}

	filename := os.Args[1]
	contents, err := os.ReadFile(filename)
	if err != nil {
		utils.Fatal(err.Error())
	}

	fmt.Println(len(contents))

}
