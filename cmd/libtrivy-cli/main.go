package main

import (
	"os"

	"github.com/cjalmeida/libtrivy/pkg/scan"
)

func main() {
	// parse args
	sourceFile := os.Args[1]
	destFile := ""
	if len(os.Args) > 2 {
		destFile = os.Args[2]
	}

	err := scan.Scan(sourceFile, destFile)
	if err != nil {
		panic(err)
	}
}
