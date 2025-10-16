package main

import (
	"fmt"
	"io"
	"os"
	"suzune/internal/http"
	"suzune/internal/util"
)

func main() {
	flags := util.ParseFlags()
	fullPath, err := util.ResolveOutputPath(flags.URL, flags.Output)
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
	if stat.Size <= 0 {
		panic(fmt.Errorf("file size must be greater than 0"))
	}
	dst, err := os.OpenFile(fullPath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(fmt.Errorf("failed to open file %s: %w", fullPath, err))
	}
	if _, err := dst.Seek(stat.Size-1, io.SeekStart); err != nil {
		panic(fmt.Errorf("seek failed: %w", err))
	}
	if _, err := dst.Write([]byte{0}); err != nil {
		panic(fmt.Errorf("preallocate failed: %w", err))
	}

	chunks, err := src.Tune()
	if err != nil {
		panic(fmt.Errorf("failed to tune: %w", err))
	}
	fmt.Printf("Mime-Type: %s | Size: %s | Chunks: %d\n", stat.ContentType, util.HumanSize(stat.Size), chunks)
	src.AddProgressBar(util.NewProgressBar(stat.Size))
	if err = src.Save(dst); err != nil {
		panic(fmt.Errorf("failed to save %s: %w", fullPath, err))
	}
	fmt.Printf("\nDone! File saved to %s\n", fullPath)
}
