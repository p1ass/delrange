package main

import (
	"github.com/p1ass/delrange"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(delrange.Analyzer) }
