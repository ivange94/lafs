package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	blockCount = 64
	blockSize  = 4 * 1024 // 4KB

	superBlock = 0

	inodeBitmap = 1
	dataBitmap  = 2

	inodeTableStart  = 3
	inodeTableLength = 5

	dataStart = 8
)

var file string

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Error: Subcommand required (format)")
		os.Exit(1)
	}

	formatCmd := flag.NewFlagSet("format", flag.ExitOnError)
	formatCmd.StringVar(&file, "file", "disk.img", "File system")
	switch os.Args[1] {
	case "format":
		formatCmd.Parse(os.Args[2:])
		format(file)
	}
}

func format(filename string) {
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating filesystem", err)
		os.Exit(1)
	}
	defer f.Close()
	err = os.Truncate(file, blockSize*blockCount)
	if err != nil {
		fmt.Println("Error truncating file", err)
		os.Exit(1)
	}
}
