package main

import (
	"flag"
	"fmt"

	"github.com/pkonkol/random/asgo/pkg/as"
)

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
}
