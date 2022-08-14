package main

import (
	"flag"
	"fmt"

	"github.com/pkonkol/random/asgo/pkg/as"
	"github.com/pkonkol/random/asgo/pkg/scan"

	"os"
	"path/filepath"
)

var scanFlag = flag.Bool("scan", false, "start scanning (unimplemented)")
var generateDBOverviewFlag = flag.Bool("generateDBOverview", false, "generate high level overview DB of ASes")
var generateDBDetailsFlag = flag.Bool("generateDBDetails", false, "generate high level overview DB of ASes")
var printDBFlag = flag.Bool("printDB", true, "print DB info")

func main() {
	flag.Parse()

	if *generateDBOverviewFlag {
		as.GenerateDBOverview()
	}

	if *generateDBDetailsFlag {
		as.GenerateDBDetails()
	}

	if *printDBFlag {
		fmt.Println("TODO: DB content overview")
	}

	if *scanFlag {
		fmt.Println("Scan not implemented, exiting")
		os.Exit(1)
		os.MkdirAll(filepath.Join(".", "tmp"), os.ModePerm)
		go scan.Run()
	}
}
