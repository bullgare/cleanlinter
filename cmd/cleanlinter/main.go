package main

import (
	"github.com/bullgare/cleanlinter"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(cleanlinter.NewAnalyzer())
}
