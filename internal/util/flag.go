package util

import (
	"flag"
	"fmt"
	"net/url"
	"os"
)

type Flags struct {
	URL            string
	Output         string
	Version        bool
	MaxConcurrency int
	ChunkSize      int
	MaxRetries     int
}

func ParseFlags() *Flags {
	flags := &Flags{}
	flag.StringVar(&flags.Output, "o", "", "Output path")
	flag.StringVar(&flags.Output, "output", "", "Output path (long form)")
	flag.BoolVar(&flags.Version, "v", false, "Show Version")
	flag.BoolVar(&flags.Version, "version", false, "Show Version (long form)")
	flag.IntVar(&flags.MaxConcurrency, "max-concurrency", 4, "Maximum concurrency")
	flag.IntVar(&flags.ChunkSize, "chunk-size", 1024, "Chunk size in bytes")
	flag.IntVar(&flags.MaxRetries, "max-retries", 5, "Maximum number of retries for each chunk")
	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: %s [options] <url>\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	args := flag.Args()
	if flags.Version {
		fmt.Println("Suzune v2.0.1")
		os.Exit(0)
	}
	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	flags.URL = args[0]
	parsed, err := url.ParseRequestURI(flags.URL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		_, _ = fmt.Fprintf(os.Stderr, "Invalid URL: %s\n", flags.URL)
		os.Exit(1)
	}
	return flags
}
