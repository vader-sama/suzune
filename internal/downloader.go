package internal

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"sync"
	"syscall"

	"github.com/schollz/progressbar/v3"
)

type Downloader struct {
	wg               *sync.WaitGroup
	url              string
	availableThreads int
	mu               sync.Mutex
}

func NewDownloader(url string) *Downloader {
	var wg sync.WaitGroup
	return &Downloader{
		&wg,
		url,
		runtime.NumCPU(),
		sync.Mutex{},
	}
}

func (d *Downloader) Download(destination string) {
	res, doesSupportPartial, err := d.getResponse()
	if err != nil {
		fmt.Println("Error occurred at download start\n\t", err)
	}

	filename := path.Base(res.Request.URL.Path)
	saveFileAt := filepath.Join(destination, filename)

	file, err := os.OpenFile(saveFileAt, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error occurred at opening file\n\t", err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Remove(saveFileAt)
		pprof.StopCPUProfile()
		os.Exit(1)
	}()
	defer file.Close()

	mimetype := res.Header.Get("content-type")
	size := res.ContentLength
	bar := progressbar.NewOptions(
		int(size),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetDescription("Downloading..."),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]#[reset]",
			SaucerHead:    "[green]#[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	fmt.Printf("Filename: %s\n", filename)
	fmt.Printf("Mimetype: %s\n", mimetype)
	fmt.Printf("Size: %s\n", formatBytes(float64(size)))

	if !doesSupportPartial {
		d.downloadNormally(res, file, bar)
		return
	}

	fmt.Printf("Available threads: %d\n", d.availableThreads)

	chunkSize := size / int64(d.availableThreads)

	for i := 0; i < d.availableThreads; i++ {
		start := int64(i) * chunkSize
		var end int64
		if i == d.availableThreads-1 {
			end = size - 1
		} else {
			end = (int64(i)+1)*chunkSize - 1
		}
		d.wg.Add(1)
		go d.downloadPartially(start, end, file, bar)
	}

	d.wg.Wait()
}

func (d *Downloader) getResponse() (*http.Response, bool, error) {
	fmt.Print("HEAD request sent, awaiting response...")
	fmt.Println(" Respond!")
	res, err := http.Head(d.url)
	if err == nil {
		doesSupportPartial := len(res.Header.Get("accept-ranges")) > 0
		if doesSupportPartial {
			return res, true, nil
		}
	}

	fmt.Println("HEAD method not supported")
	fmt.Println("Threaded download is unavailable, switched to single-thread download")

	res, err = http.Get(d.url)
	return res, false, err
}

func (d *Downloader) downloadPartially(start int64, end int64, file *os.File, bar *progressbar.ProgressBar) {
	defer d.wg.Done()
	req, err := http.NewRequest("GET", d.url, nil)
	if err != nil {
		fmt.Println("Error occurred at creating request\n\t", err)
		return
	}

	req.Header.Set("range", fmt.Sprintf("bytes=%d-%d", start, end))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error occurred at GET request\n\t", err)
	}

	defer res.Body.Close()
	d.mu.Lock()
	defer d.mu.Unlock()

	_, err = file.Seek(start, io.SeekStart)
	if err != nil {
		if err != nil {
			fmt.Println("Error occurred at file seeking\n\t", err)
		}
	}

	_, err = io.CopyN(io.MultiWriter(file, bar), res.Body, end-start+1)
	if err != nil {
		fmt.Println("Error occurred at writing to file\n\t", err)
	}
}

func (d *Downloader) downloadNormally(res *http.Response, file *os.File, bar *progressbar.ProgressBar) {
	defer res.Body.Close()

	_, err := io.Copy(io.MultiWriter(file, bar), res.Body)
	if err != nil {
		fmt.Println("error occurred at writing file\n\t", err)
	}
}
