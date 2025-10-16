package http

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"suzune/internal/util"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
)

type File struct {
	*util.Flags
	url         string
	stat        *FileStat
	resp        *http.Response
	progressBar *progressbar.ProgressBar
	chunkSize   int64
	chunks      int64
}

func NewFile(url string, flags *util.Flags) *File {
	return &File{
		url:   url,
		Flags: flags,
	}
}

func (f *File) AddProgressBar(progress *progressbar.ProgressBar) {
	f.progressBar = progress
}

func (f *File) Stat() (*FileStat, error) {
	if f.stat != nil {
		return f.stat, nil
	}
	err := f.request()
	if err != nil {
		return nil, err
	}
	f.stat = &FileStat{
		Size:        f.resp.ContentLength,
		ContentType: f.resp.Header.Get("Content-Type"),
	}
	return f.stat, nil
}

func (f *File) Tune() (int64, error) {
	var err error
	if f.resp == nil {
		err = f.request()
	}
	if f.stat == nil {
		f.stat, err = f.Stat()
	}
	if err != nil {
		return 0, err
	}
	if f.resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("bad status: %s", f.resp.Status)
	}

	f.chunkSize, f.chunks = f.calcChunks()
	if f.resp.Header.Get("Accept-Ranges") != "bytes" {
		f.chunks = 1
	}
	return f.chunks, nil
}

func (f *File) Save(dst *os.File) error {
	defer func(dst *os.File) {
		_ = dst.Close()
	}(dst)
	if f.chunks == 1 {
		_, err := io.Copy(io.MultiWriter(dst, f.progressBar), f.resp.Body)
		if err != nil {
			return err
		}
	}
	_ = f.resp.Body.Close()
	var wg sync.WaitGroup
	errCh := make(chan error, f.chunks)
	for i := int64(0); i < f.chunks; i++ {
		start := i * f.chunkSize
		end := start + f.chunkSize
		if end > f.stat.Size {
			end = f.stat.Size
		}
		size := end - start
		wg.Add(1)
		go func(offset, size int64) {
			defer wg.Done()
			f.savePartialWithRetries(errCh, dst, offset, size)
		}(start, size)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			return err
		}
	}

	return dst.Sync()
}

func (f *File) Close() error {
	if f.resp != nil {
		return f.resp.Body.Close()
	}
	return nil
}

func (f *File) savePartialWithRetries(ch chan<- error, dst *os.File, offset int64, size int64) {
	for attempt := 1; attempt <= f.MaxRetries; attempt++ {
		if err := f.savePartial(dst, offset, size); err != nil {
			if attempt == f.MaxRetries {
				ch <- fmt.Errorf("chunk %d failed after %d retries: %w", offset, attempt, err)
			}
			time.Sleep(time.Duration(attempt) * time.Second)
			continue
		}
		break
	}
}

func (f *File) savePartial(dst *os.File, offset int64, size int64) error {
	req, err := http.NewRequest("GET", f.url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", offset, offset+size-1))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	if resp.StatusCode != http.StatusPartialContent {
		return fmt.Errorf("not partial: %s", resp.Status)
	}
	buf := make([]byte, 32*1024)
	var written int64
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			_, err := dst.WriteAt(buf[:n], offset+written)
			if err != nil {
				return err
			}
			written += int64(n)
		}
		if f.progressBar != nil {
			_ = f.progressBar.Add(n)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *File) calcChunks() (int64, int64) {
	chunks := int(f.stat.Size / int64(f.ChunkSize))
	if chunks == 0 {
		chunks = 1
	} else if chunks > f.MaxConcurrency {
		chunks = f.MaxConcurrency
	}
	return (f.stat.Size + int64(chunks) - 1) / int64(chunks), int64(chunks)
}

func (f *File) request() error {
	resp, err := http.Get(f.url)
	if err != nil {
		return err
	}
	f.resp = resp
	return nil
}
