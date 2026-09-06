package main

import (
	"bytes"
	"encoding/binary"
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

	lafsMagic = 0x4C414653
)

var file string

type Superblock struct {
	Magic      uint32
	BlockCount uint32
	BlockSize  uint32

	InodeBitmap uint32
	DataBitmap  uint32

	InodeTableStart uint32
	InodeTableLen   uint32

	DataStart uint32
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Error: Subcommand required (format)")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "format":
		format(os.Args[2:])
	}
}

func format(args []string) {
	fs := flag.NewFlagSet("format", flag.ExitOnError)
	fs.StringVar(&file, "file", "disk.img", "File system")
	fs.Parse(args)

	f, err := os.Create(file)
	check(err, "error creating file system file")
	defer f.Close()

	err = f.Truncate(blockSize * blockCount)
	check(err, "error truncating file system file")

	sb := Superblock{
		Magic:           lafsMagic,
		BlockCount:      blockCount,
		BlockSize:       blockSize,
		InodeBitmap:     inodeBitmap,
		DataBitmap:      dataBitmap,
		InodeTableStart: inodeTableStart,
		InodeTableLen:   inodeTableLength,
		DataStart:       dataStart,
	}
	check(writeSuperblock(sb, f), "error writing superblock")
}

func writeSuperblock(sb Superblock, file *os.File) error {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, sb)
	if err != nil {
		return err
	}
	totalBytes := buf.Len()
	n, err := file.WriteAt(buf.Bytes(), blockOffset(superBlock))
	if err != nil {
		return err
	}
	if totalBytes != n {
		fmt.Printf("[WARN] %d of %d bytes written\n", n, totalBytes)
	}
	return nil
}

func readSuperblock(f *os.File) (Superblock, error) {
	size := binary.Size(Superblock{})
	buf := make([]byte, size)
	_, err := f.ReadAt(buf, blockOffset(superBlock))
	if err != nil {
		return Superblock{}, fmt.Errorf("error reading superblock bytes from file system")
	}
	sb := Superblock{}
	err = binary.Read(bytes.NewReader(buf), binary.LittleEndian, &sb)
	if err != nil {
		return Superblock{}, fmt.Errorf("error parsing superblock bytes")
	}
	if sb.Magic != lafsMagic {
		return Superblock{}, fmt.Errorf("not a lafs file system")
	}
	return sb, nil
}

func check(err error, msg string) {
	if err != nil {
		fmt.Println(msg, err)
		os.Exit(1)
	}
}

func blockOffset(block int) int64 {
	return int64(block * blockSize)
}
