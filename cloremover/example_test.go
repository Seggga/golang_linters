package cloremover_test

import (
	"os"

	"github.com/seggga/golang_linters/cloremover"
)

func Example() {
	// read flags
	conf := &cloremover.ConfigType{}
	_ = cloremover.ReadFlags(conf)
	currentDir, _ := os.Getwd()
	fileSystem := os.DirFS(currentDir)
	// collect data
	fileSlice, _ := cloremover.FindClones(conf, nil, fileSystem)
	// display data
	outputMap := cloremover.PrintClones(fileSlice, conf)
	// remove data
	cloremover.Remove(fileSlice, conf, outputMap, nil)
}
