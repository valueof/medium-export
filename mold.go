package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/valueof/mold/parser"
)

func walk(dir string, logger *log.Logger, fn func(string, io.Reader)) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		logger.Fatalf("%s: %v\n", path.Base(dir), err)
	}

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".html") == false {
			logger.Printf("skipped %s: not .html", f.Name())
			continue
		}

		dat, err := os.Open(path.Join(dir, f.Name()))
		if err != nil {
			logger.Fatalf("%s: %v", f.Name(), err)
			continue
		}

		fn(f.Name(), dat)
	}
}

func main() {
	var (
		buf    bytes.Buffer
		logger = log.New(&buf, "mold: ", log.Llongfile)
		root   = "./data"
	)

	dirs, err := ioutil.ReadDir(root)
	if err != nil {
		panic(err)
	}

	for _, d := range dirs {
		if d.IsDir() == false {
			logger.Printf("skipped %s: not a directory", d.Name())
			continue
		}

		switch d.Name() {
		case "blocks":
			users := []parser.MediumUser{}
			walk(path.Join(root, d.Name()), logger, func(fn string, dat io.Reader) {
				part, err := parser.ParseBlocked(dat)
				if err != nil {
					logger.Fatalf("%s: %v", fn, err)
					return
				}
				users = append(users, part...)
			})
			fmt.Printf("blocks: %v\n", users)
		}
	}

	// fmt.Print(&buf)
}