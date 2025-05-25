# README

A Go package for efficient directory traversal and file listing, with flexible depth control and exclusion options.

Supports both OS-native traversal and fast index-based search using [Everything](https://www.voidtools.com/) (Windows only).

## Features

- Traverse directory trees from a specified root.
- Optionally include files, directories, or both.
- Control the maximum traversal depth.
- Exclude specified directory or file names.
- Two traversal backends:
  - Standard Go filesystem walk (`filepath.WalkDir`)
  - Fast search via [Everything](https://www.voidtools.com/) (if available, Windows only)


### Example

```go
w := &walk.Walker{}
w.Init("C:\\Users\\YourName", false, 2, "node_modules,.git")

// OS-native traversal
found, err := w.Traverse()
if err != nil {
    log.Fatal(err)
}
fmt.Println("Found:", found)

// Everything-based traversal (Windows, Everything required)
found2, err := w.EverythingTraverse()
if err != nil {
    log.Fatal(err)
}
fmt.Println("Found (Everything):", found2)
```

## API

### `func (w *Walker) Init(root string, all bool, depth int, exclude string)`

Initializes the walker.

- `root`: Root directory to start traversal.
- `all`: If true, include files and directories. If false, include only directories.
- `depth`: Maximum traversal depth (`-1` for unlimited).
- `exclude`: Comma-separated list of directory/file names to exclude.

### `func (w Walker) Traverse() ([]string, error)`

Traverse the filesystem using Go's standard library.

Returns a list of paths found.

### `func (w Walker) EverythingTraverse() ([]string, error)`

Traverse using Everything index (much faster, Windows only).

Returns a list of paths found.

## Notes

- Exclusions apply to names (not full paths), and always exclude hidden directories (starting with a dot) and `AppData` by default.
- If using `EverythingTraverse`, ensure the following requirements are met:
  - [Everything](https://www.voidtools.com/) is installed and running.
  - [Everything SDK](https://www.voidtools.com/support/everything/sdk/) is downloaded and extracted.
  - The file `Everything64.dll` from the SDK must be placed **in the same directory as `walk.go`** or **in the same directory as the compiled executable file**.
  - The Go package [`github.com/AWtnb/go-everything`](https://github.com/AWtnb/go-everything) is available.


## License

MIT

## Author

[AWtnb](https://github.com/AWtnb)