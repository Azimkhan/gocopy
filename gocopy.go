package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var limit int
var offset int
var src string
var dest string

func init() {
	flag.IntVar(&limit, "limit", 0, "limit")
	flag.IntVar(&offset, "offset", 0, "offset")
	flag.StringVar(&src, "src", "", "source file")
	flag.StringVar(&dest, "dest", "", "destination file")

}
func main() {
	flag.Parse()
	fmt.Printf("%v %v %v %v\n", limit, offset, src, dest)
	err := Copy(src, dest, limit, offset)
	if err != nil {
		log.Fatal(err)
	}
}

func Copy(src string, dest string, limit int, offset int) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("src file open error: %v", err)
	}
	defer srcFile.Close()

	if limit <= 0 {
		return fmt.Errorf("limit must be > 0")
	}

	if offset > 0 {
		_, seekErr := srcFile.Seek(int64(offset), io.SeekStart)
		if seekErr != nil {
			return fmt.Errorf("seek error: %v", seekErr)
		}
	}

	destFile, err2 := os.Create(dest)
	if err2 != nil {
		return fmt.Errorf("dest file open error: %v", err2)
	}
	defer destFile.Close()

	_, copyErr := io.CopyN(destFile, srcFile, int64(limit))
	if copyErr != nil {
		return fmt.Errorf("copy error")
	}
	return nil
}
