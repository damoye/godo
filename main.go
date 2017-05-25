package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"
)

// const url = "https://github.com/qt/qt/archive/4.8.zip"
// const url = "http://ftp.jaist.ac.jp/pub/qtproject/archive/qt/5.7/5.7.0/qt-opensource-mac-x64-clang-5.7.0.dmg"

func main() {
	url := os.Args[1]
	start := time.Now()

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.ContentLength <= 0 {
		panic("ContentLength <= 0")
	}
	file, err := os.Create(path.Base(url))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var written, recordWritten int64
	recordTime := time.Now()
	buf := make([]byte, 64*1024)
	for {
		n, er := resp.Body.Read(buf)
		if n > 0 {
			if _, ew := file.Write(buf[:n]); ew != nil {
				panic(ew)
			}
			written += int64(n)
			now := time.Now()
			elapsed := now.Sub(recordTime).Seconds()
			if elapsed >= 1 {
				if resp.ContentLength > 0 {
					fmt.Printf(
						"\rPROGRESS: %.2f %%, SPEED: %.2f KB/s, TOTAL: %d MB   ",
						float64(written)/float64(resp.ContentLength)*100,
						float64(written-recordWritten)/elapsed/1024,
						resp.ContentLength/1024/1024,
					)
				} else {
					fmt.Printf(
						"\rWRITTEN: %dMB, SPEED: %.2f",
						written/1024/1024,
						float64(written-recordWritten)/elapsed/1024,
					)
				}
				recordTime, recordWritten = now, written
			}
		}
		if er != nil {
			if er == io.EOF {
				break
			}
			panic(er)
		}
	}
	fmt.Printf("\rAVERAGE_SPEED : %.2f KB/s", float64(resp.ContentLength)/time.Now().Sub(start).Seconds()/1024)
}
