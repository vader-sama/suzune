# Suzune ⚡

**Suzune** is a command-line tool to download files efficiently. It supports concurrent downloads, chunked downloads,
retries, and displays progress in a clean, human-readable format. Perfect for grabbing large files quickly from the
terminal.

## Features

- Concurrent downloads with configurable chunk count
- Chunk size control for optimal performance
- Automatic retries for failed chunks
- Lightweight, written in pure Go

## Build from source

```bash
git clone https://github.com/vader-sama/suzune.git
cd suzune
make build
```

The binary will be placed in:

```bash
./bin/suzune
```

## Options

| Flag                | Description                 | Default           |
|---------------------|-----------------------------|-------------------|
| `-o, --output`      | Output file path            | Current directory |
| `-v, --version`     | Show version                | —                 |
| `--max-concurrency` | Number of concurrent chunks | 4                 |
| `--chunk-size`      | Size of each chunk in bytes | 1024              |
| `--max-retries`     | Number of retries per chunk | 5                 |

## Examples

Download a file with default settings:

```bash
suzune https://example.com/file.zip
```

Download a file with custom name:

```bash
suzune -o custom.zip https://example.com/file.zip
```

## Contributing

Contributions are welcome! Feel free to open issues or pull requests.