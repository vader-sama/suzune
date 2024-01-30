# Suzune
Multi-threaded file downloader written in Go

## Usage
Download a file from URL
```bash
suzune https://example.com/dummy.mp4
```

Download a file with custom path or name
```bash
suzune -o /samples/custom.mp4 https://example.com/dummy.mp4 
```

## How it works?
First it performs a `HEAD` request to the given URL to find out whatever the given URL supports `206 Partial Content` or not, If it doesn't, **Suzune** switches to single-thread method. However if it does support `206 Partial Content`, **Suzune** downloads the file partially by multiple **Go Routines**.