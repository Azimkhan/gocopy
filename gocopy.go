package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
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
	//fmt.Printf("%v %v %v %v\n", limit, offset, src, dest)
	err := Copy(src, dest, limit, offset)
	if err != nil {
		log.Fatal(err)
	}
}

func Copy(src string, dest string, limit int, offset int) error {
	var err error
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("src file open error: %v", err)
	}
	defer srcFile.Close()

	// check offset
	if offset > 0 {
		pos, seekErr := srcFile.Seek(int64(offset), io.SeekStart)
		fmt.Println(pos)
		if seekErr != nil {
			return fmt.Errorf("seek error: %v", seekErr)
		}
	}

	destFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("dest file open error: %v", err)
	}
	defer destFile.Close()

	p := make(chan int, 1)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(pChan chan int) {
		defer wg.Done()
		for i := range pChan {
			fmt.Printf("Bytes copied: %d\r", i)
		}
		fmt.Println("")
	}(p)
	_, err = CopyN(destFile, srcFile, int64(limit), p)
	wg.Wait()

	if err != nil {
		return fmt.Errorf("copy error")
	}
	return nil
}

func CopyN(dst io.Writer, src io.Reader, limit int64, p chan int) (written int64, err error) {
	if limit > 0 {
		src = io.LimitReader(src, limit)
	}
	written, err = io.Copy(dst, NewReaderWithProgress(src, p))
	if written == limit {
		return limit, nil
	}
	if written < limit && err == nil {
		// src stopped early; must have been EOF.
		err = io.EOF
	}
	return
}

func NewReaderWithProgress(reader io.Reader, pChan chan int) io.Reader {
	return &ReaderWithProgress{reader, pChan, 0}
}

type ReaderWithProgress struct {
	R         io.Reader
	pChan     chan int
	bytesRead int
}

func (l *ReaderWithProgress) Read(p []byte) (n int, err error) {
	defer func() {
		if n > 0 {
			l.pChan <- l.bytesRead
		}
		if err == io.EOF {
			close(l.pChan)
		}
	}()
	n, err = l.R.Read(p)
	l.bytesRead += n
	return
}
