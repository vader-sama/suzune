package main

import (
	"flag"
	"fmt"
	"os"
	"suzune/internal"
)

func main() {
	var destination string

	flag.StringVar(&destination, "o", "./", "Save file at")
	flag.Parse()

	url := flag.Arg(0)
	if url == "" {
		fmt.Println("No url specified")
		fmt.Println("Usage: suzune -url <url>")
		os.Exit(1)
	}

	d := internal.NewDownloader(url)
	d.Download(destination)
	fmt.Println("\nDownload completed!")
}
