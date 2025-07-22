# fj - formatjson

I built this simple program to be able to parse JSON without the use of an online tool, even better if it used the terminal. So I built it.

## Features

- Format JSON from files, URLs, pipes or standard input
- Customize indentation spaces
- Sort object keys
- Automatic clipboard integration
- Auto-save formatted JSON to files
- Cross-platform support (macOS, Linux, Windows)
- Simple configuration system

## Installation

### Binary Release

Download the latest binary for your platform from the [GitHub Releases](https://github.com/nicolasalberti00/fj/releases) page.

### To install the project

```bash
# Clone the repository
git clone https://github.com/nicolasalberti00/fj.git
cd fj

# Build the binary
go build -o fj ./cmd/fj

# Install the binary 
go install ./cmd/fj
```

## Building and Running

### Prerequisites

- Go 1.18 or later
- Git

### Building the Project

Note: tests at the moment are a WIP, I will update them whenever they are ready and well integrated for all platforms.

1. Clone the repository:
   ```bash
   git clone https://github.com/nicolasalberti00/fj.git
   cd fj
   ```

2. Build the project:
   ```bash
   go build -o fj ./cmd/fj
   ```

3. Run the tests:
   ```bash
   go test ./...
   ```

## Usage

```bash
# Format JSON from a file
fj file.json

# Format JSON from a URL
fj https://example.com/data.json

# Format JSON from stdin
cat file.json | fj

# Format with 4-space indentation
fj -indent 4 file.json

# Format with sorted keys
fj -sort file.json

# Disable clipboard copy
fj -clipboard=false file.json

# Save current settings as default
fj -indent 4 -sort -save-config
```

## Command-Line Options

- `-indent int`: Number of spaces for indentation (default 2)
- `-sort`: Sort object keys
- `-clipboard`: Copy result to clipboard (default true)
- `-outdir string`: Output directory for saved files
- `-trust-all`: Trust all URLs without prompting
- `-save-config`: Save current flags as default configuration
- `-version`: Show version information
- `-help`: Show help information

## Configuration

fj uses a configuration file stored in:
- Windows: `%USERPROFILE%\fj\config.json`
- macOS and Linux: `~/.config/fj/config.json`

You can save your preferred settings using the `-save-config` flag.

## Upcoming Features

- Interactive mode
- JSON diff functionality
- JSON schema validation
- Internationalization support

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit an issue and/or a Pull Request!

Thanks for the interest and happy formatting!
