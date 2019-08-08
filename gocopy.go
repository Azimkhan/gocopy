package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"
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

	p := make(chan int, 1)
	go func() {
		for i := range p {
			fmt.Printf("Progress: %d\r", i)
		}
	}()
	_, copyErr := CopyN(destFile, srcFile, int64(limit), p)
	close(p)

	if copyErr != nil {
		return fmt.Errorf("copy error")
	}
	return nil
}

func CopyN(dst io.Writer, src io.Reader, n int64, p chan int) (written int64, err error) {
	written, err = io.Copy(dst, &LimitedReaderWithProgress{src, n, p, 0})
	if written == n {
		return n, nil
	}
	if written < n && err == nil {
		// src stopped early; must have been EOF.
		err = io.EOF
	}
	return
}

type LimitedReaderWithProgress struct {
	R    io.Reader // underlying reader
	N    int64     // max bytes remaining
	P    chan int
	read int
}

func (l *LimitedReaderWithProgress) Read(p []byte) (n int, err error) {
	time.Sleep(10 * time.Millisecond)
	if l.N <= 0 {
		return 0, io.EOF
	}
	if int64(len(p)) > l.N {
		p = p[0:l.N]
	}
	n, err = l.R.Read(p)
	l.N -= int64(n)
	l.read += n
	l.P <- l.read
	return
}
