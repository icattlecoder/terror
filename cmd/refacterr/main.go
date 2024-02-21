package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/dave/dst/decorator"
	"github.com/icattlecoder/terrors/cmd/parser"
	"golang.org/x/tools/go/packages"
	"log"
	"os"
)

var (
	dir   = flag.String("dir", "", "project dir")
	pkg   = flag.String("pkg", "", "pkg path")
	write = flag.Bool("w", false, "write back")
)

func main() {
	flag.Parse()

	pkgs, err := decorator.Load(&packages.Config{
		Dir:  *dir,
		Mode: packages.LoadAllSyntax,
	}, *pkg)
	if err != nil {
		log.Fatalln(err)
	}

	if *dir == "" || *pkg == "" {
		flag.PrintDefaults()
		return
	}
	for _, p := range pkgs {
		if err := parser.New(p).Run(*write); err != nil {
			if !*write {
				fmt.Println("found no traced error:")
			}
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
