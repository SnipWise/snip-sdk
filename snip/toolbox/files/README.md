# Files Package

The `files` package provides utility functions for working with files and directories in Go.

## Overview

This package offers simple, convenient functions for finding files, reading and writing text files, and iterating over files in directories. It simplifies common file operations with a clean, easy-to-use API.

## Functions

### FindFiles

Searches for files with a specific extension in the given root directory and its subdirectories (recursive search).

**Signature:**
```go
func FindFiles(dirPath string, ext string) ([]string, error)
```

**Parameters:**
- `dirPath` - The root directory to start the search from
- `ext` - The file extension to search for (e.g., ".md", ".html", ".txt"), or ".*" to match all files

**Returns:**
- `[]string` - A slice of file paths that match the given extension
- `error` - An error if the search encounters any issues

**Example:**
```go
package main

import (
    "fmt"
    "log"
    "github.com/snipwise/snip-sdk/files"
)

func main() {
    // Find all markdown files in a directory
    mdFiles, err := files.FindFiles("./docs", ".md")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Found markdown files:")
    for _, file := range mdFiles {
        fmt.Println(file)
    }

    // Find all files (any extension)
    allFiles, err := files.FindFiles("./project", ".*")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Total files found: %d\n", len(allFiles))
}
```

---

### ForEachFile

Iterates over all files with a specific extension in a directory and its subdirectories, executing a callback function for each file found.

**Signature:**
```go
func ForEachFile(dirPath string, ext string, callback func(string) error) ([]string, error)
```

**Parameters:**
- `dirPath` - The root directory to start the search from
- `ext` - The file extension to search for, or ".*" to match all files
- `callback` - A function to be called for each file found. If the callback returns an error, the iteration stops

**Returns:**
- `[]string` - A slice of file paths that match the given extension
- `error` - An error if the search encounters any issues or if the callback returns an error

**Example:**
```go
package main

import (
    "fmt"
    "log"
    "github.com/snipwise/snip-sdk/files"
)

func main() {
    // Process each text file
    processedFiles, err := files.ForEachFile("./data", ".txt", func(path string) error {
        fmt.Printf("Processing: %s\n", path)

        // Read and process the file
        content, err := files.ReadTextFile(path)
        if err != nil {
            return err
        }

        fmt.Printf("  Size: %d bytes\n", len(content))
        return nil
    })

    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Processed %d files\n", len(processedFiles))
}
```

**Example with early termination:**
```go
package main

import (
    "fmt"
    "log"
    "errors"
    "github.com/snipwise/snip-sdk/files"
)

func main() {
    // Stop iteration when a specific file is found
    _, err := files.ForEachFile("./logs", ".log", func(path string) error {
        fmt.Printf("Checking: %s\n", path)

        content, err := files.ReadTextFile(path)
        if err != nil {
            return err
        }

        // Stop if we find an error in the log
        if strings.Contains(content, "FATAL") {
            return errors.New("fatal error found in " + path)
        }

        return nil
    })

    if err != nil {
        fmt.Println("Stopped:", err)
    }
}
```

---

### GetContentFiles

Searches for files with a specific extension and reads all their contents into memory.

**Signature:**
```go
func GetContentFiles(dirPath string, ext string) ([]string, error)
```

**Parameters:**
- `dirPath` - The directory path to start the search from
- `ext` - The file extension to search for

**Returns:**
- `[]string` - A slice of file contents (each element is the content of one file)
- `error` - An error if the search or reading encounters any issues

**Example:**
```go
package main

import (
    "fmt"
    "log"
    "strings"
    "github.com/snipwise/snip-sdk/files"
)

func main() {
    // Read all markdown files
    contents, err := files.GetContentFiles("./docs", ".md")
    if err != nil {
        log.Fatal(err)
    }

    // Process all contents
    totalLines := 0
    for i, content := range contents {
        lines := strings.Split(content, "\n")
        totalLines += len(lines)
        fmt.Printf("File %d: %d lines\n", i+1, len(lines))
    }

    fmt.Printf("Total lines across all files: %d\n", totalLines)
}
```

**Note:** Be careful when using this function with large files or many files, as it loads all content into memory at once.

---

### ReadTextFile

Reads the contents of a text file at the given path and returns the contents as a string.

**Signature:**
```go
func ReadTextFile(path string) (string, error)
```

**Parameters:**
- `path` - The path to the text file

**Returns:**
- `string` - The contents of the text file as a string
- `error` - An error if the file cannot be read

**Example:**
```go
package main

import (
    "fmt"
    "log"
    "github.com/snipwise/snip-sdk/files"
)

func main() {
    // Read a configuration file
    content, err := files.ReadTextFile("config.json")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Configuration:")
    fmt.Println(content)

    // Read and process a text file
    data, err := files.ReadTextFile("data.txt")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Read %d bytes\n", len(data))
}
```

---

### WriteTextFile

Writes content to a text file, creating the file if it doesn't exist or overwriting it if it does.

**Signature:**
```go
func WriteTextFile(path, content string) error
```

**Parameters:**
- `path` - The path where the file should be written
- `content` - The string content to write to the file

**Returns:**
- `error` - An error if the file cannot be created or written

**Example:**
```go
package main

import (
    "fmt"
    "log"
    "time"
    "github.com/snipwise/snip-sdk/files"
)

func main() {
    // Write a simple text file
    err := files.WriteTextFile("output.txt", "Hello, World!")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("File written successfully")

    // Generate and write a report
    report := fmt.Sprintf("Report generated at: %s\n", time.Now())
    report += "Status: OK\n"
    report += "Items processed: 42\n"

    err = files.WriteTextFile("report.txt", report)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Report saved")
}
```

---

### GetAllFilesInDirectory

Returns all file paths in the given directory (non-recursive, files only, no subdirectories).

**Signature:**
```go
func GetAllFilesInDirectory(dirPath string) ([]string, error)
```

**Parameters:**
- `dirPath` - The directory path to read

**Returns:**
- `[]string` - A slice of full file paths in the directory
- `error` - An error if the directory cannot be read

**Example:**
```go
package main

import (
    "fmt"
    "log"
    "path/filepath"
    "github.com/snipwise/snip-sdk/files"
)

func main() {
    // Get all files in current directory
    files, err := files.GetAllFilesInDirectory(".")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Files in current directory:")
    for _, file := range files {
        fmt.Printf("  - %s\n", filepath.Base(file))
    }

    // Get files from a specific directory
    dataFiles, err := files.GetAllFilesInDirectory("./data")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Found %d files in ./data\n", len(dataFiles))
}
```

---

## Common Use Cases

### Batch Processing Files

```go
package main

import (
    "fmt"
    "log"
    "strings"
    "github.com/snipwise/snip-sdk/files"
)

func main() {
    // Find all text files and convert to uppercase
    txtFiles, err := files.FindFiles("./documents", ".txt")
    if err != nil {
        log.Fatal(err)
    }

    for _, filePath := range txtFiles {
        // Read file
        content, err := files.ReadTextFile(filePath)
        if err != nil {
            log.Printf("Error reading %s: %v\n", filePath, err)
            continue
        }

        // Process content
        upperContent := strings.ToUpper(content)

        // Write back
        outputPath := strings.Replace(filePath, ".txt", "_upper.txt", 1)
        err = files.WriteTextFile(outputPath, upperContent)
        if err != nil {
            log.Printf("Error writing %s: %v\n", outputPath, err)
            continue
        }

        fmt.Printf("Processed: %s -> %s\n", filePath, outputPath)
    }
}
```

### Building a File Index

```go
package main

import (
    "fmt"
    "log"
    "path/filepath"
    "github.com/snipwise/snip-sdk/files"
)

type FileInfo struct {
    Path string
    Size int
}

func main() {
    var index []FileInfo

    // Build index of all markdown files
    _, err := files.ForEachFile("./docs", ".md", func(path string) error {
        content, err := files.ReadTextFile(path)
        if err != nil {
            return err
        }

        index = append(index, FileInfo{
            Path: filepath.Base(path),
            Size: len(content),
        })

        return nil
    })

    if err != nil {
        log.Fatal(err)
    }

    // Print index
    fmt.Println("File Index:")
    for _, info := range index {
        fmt.Printf("  %s: %d bytes\n", info.Path, info.Size)
    }
}
```

### Directory Statistics

```go
package main

import (
    "fmt"
    "log"
    "path/filepath"
    "github.com/snipwise/snip-sdk/files"
)

func main() {
    // Get statistics for a directory
    allFiles, err := files.FindFiles("./project", ".*")
    if err != nil {
        log.Fatal(err)
    }

    // Count by extension
    extCount := make(map[string]int)
    for _, file := range allFiles {
        ext := filepath.Ext(file)
        if ext == "" {
            ext = "(no extension)"
        }
        extCount[ext]++
    }

    fmt.Println("File statistics:")
    fmt.Printf("Total files: %d\n\n", len(allFiles))
    fmt.Println("By extension:")
    for ext, count := range extCount {
        fmt.Printf("  %s: %d files\n", ext, count)
    }
}
```

### Combining Files

```go
package main

import (
    "fmt"
    "log"
    "strings"
    "github.com/snipwise/snip-sdk/files"
)

func main() {
    // Combine all markdown files into one
    contents, err := files.GetContentFiles("./chapters", ".md")
    if err != nil {
        log.Fatal(err)
    }

    // Join with separator
    combined := strings.Join(contents, "\n\n---\n\n")

    // Write combined file
    err = files.WriteTextFile("book.md", combined)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Combined %d files into book.md\n", len(contents))
}
```

---

## Comparison: FindFiles vs ForEachFile vs GetContentFiles

| Function | Returns | Memory Usage | Use Case |
|----------|---------|--------------|----------|
| `FindFiles` | File paths only | Low | When you need file paths for later processing |
| `ForEachFile` | File paths + callback execution | Low-Medium | When you need to process files one at a time |
| `GetContentFiles` | All file contents | High | When you need all file contents in memory at once |

---

## Usage Notes

- **Recursive vs Non-Recursive:**
  - `FindFiles`, `ForEachFile`, and `GetContentFiles` are **recursive** (search subdirectories)
  - `GetAllFilesInDirectory` is **non-recursive** (only the specified directory)

- **Extension Matching:**
  - Use `".*"` to match all files regardless of extension
  - Extensions are case-sensitive (`.txt` ≠ `.TXT`)
  - Include the dot in the extension (use `".md"`, not `"md"`)

- **Error Handling:**
  - `ForEachFile` stops iteration if the callback returns an error
  - Use this for early termination when a condition is met

- **Memory Considerations:**
  - `GetContentFiles` loads all file contents into memory - be cautious with large files
  - For large file sets, prefer `ForEachFile` to process files one at a time

- **File Paths:**
  - All functions return absolute or relative paths based on the input `dirPath`
  - Use `filepath.Base()` to get just the filename
  - Use `filepath.Dir()` to get the directory part

## Installation

```bash
go get github.com/snipwise/snip-sdk/files
```

## Related Packages

- `path/filepath` - Standard library for path manipulation
- `os` - Standard library for file operations

## License

See the main repository for license information.
