package downloader

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"
)

// Downloader ...
type Downloader struct {
	url     *url.URL
	written int
}

// New ...
func New(rawurl string) (*Downloader, error) {
	url, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	return &Downloader{url: url}, nil
}

// Run ...
func (d *Downloader) Run() error {
	file, err := os.Create(path.Base(d.url.Path))
	if err != nil {
		return err
	}
	defer file.Close()
	resp, err := http.Get(d.url.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	buf := make([]byte, 32*1024)
	var recordWritten int
	recordTime := time.Now()
	for {
		n, er := resp.Body.Read(buf)
		if n > 0 {
			if _, ew := file.Write(buf[:n]); ew != nil {
				return ew
			}
			d.written += n
			now := time.Now()
			if seconds := now.Sub(recordTime).Seconds(); seconds > 1 {
				if resp.ContentLength == 0 {
					fmt.Printf(
						"\rPROGRESS: %6.2fMB, SPEED: %6.2fKB/s",
						float64(d.written)/1024/1024,
						float64(d.written-recordWritten)/1024/seconds,
					)
				} else {
					fmt.Printf(
						"\rPROGRESS: %6.2f%%, SPEED: %6.2fKB/s, TOTAL: %dMB",
						float64(d.written)/float64(resp.ContentLength)*100,
						float64(d.written-recordWritten)/1024/seconds,
						resp.ContentLength/1024/1024,
					)
				}
				recordWritten = d.written
				recordTime = now
			}
		}
		if er == nil {
			continue
		}
		if er == io.EOF {
			break
		}
		return er
	}
	return nil
}
