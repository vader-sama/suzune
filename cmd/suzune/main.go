package main

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"suzune/internal/http"
	"suzune/internal/util"
)

func main() {
	flags := util.ParseFlags()
	u, err := url.Parse(flags.URL)
	if err != nil {
		fmt.Println(err)
		return
	}
	src := http.NewFile(flags.URL, flags)
	stat, err := src.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}
	filename := path.Base(u.Path)
	fullPath := filepath.Join(flags.Output, filename)
	dst, err := os.OpenFile(fullPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		panic(fmt.Errorf("failed to open file %s: %w", fullPath, err))
	}
	if err := dst.Truncate(stat.Size); err != nil {
		panic(fmt.Errorf("preallocate failed: %w", err))
	}
	chunks, err := src.Tune()
	if err != nil {
		panic(fmt.Errorf("failed to tune %s: %w", filename, err))
	}
	fmt.Printf("Mime-Type: %s | Size: %s | Chunks: %d\n", stat.ContentType, util.HumanSize(stat.Size), chunks)
	src.AddProgressBar(util.NewProgressBar(stat.Size))
	if err = src.Save(dst); err != nil {
		panic(fmt.Errorf("failed to save %s: %w", fullPath, err))
	}
	fmt.Printf("Done! File saved to %s\n", fullPath)
}
