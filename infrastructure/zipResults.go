package infrastructure

import (
	"fmt"
	"neTTool/helper"
)

func ZipResults(folder2zip string, resultpath string) {
	if err := helper.ZipSource(folder2zip, resultpath); err != nil {
		fmt.Println(err)
	}
}
