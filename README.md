# gotext-tools

[![Go Report Card](https://goreportcard.com/badge/github.com/Tom5521/xgotext)](https://goreportcard.com/report/github.com/Tom5521/xgotext)

A Go library and CLI toolkit for working with Gettext `.po` and `.mo` files. This project provides cross-platform tools for extracting, merging, and compiling translation files, with a focus on Go projects.

---

## Features

✅ **Cross-platform** – Written in Go, works on Linux, NetBSD, FreeBSD, OpenBSD, Windows, Darwin(MacOS), Plan9, Solaris and DragonFly. 386, amd64, arm, arm64, ppc64, ppc64le, riscv64.

✅ **Go-native support** – Designed with Go projects in mind.  
✅ **Easy to Install** – Pre-built binaries available, or build from source.

---

## CLI Tools

Located in `cli/`, these tools provide command-line utilities for managing `.po` and `.mo` files.

### `gotext-tools`

It is a wrapper of all the utilities that will be mentioned below, if you don't want to install all of them one by one, or you want to save disk space, I recommend you install this one.

**Usage:**

```sh
gotext-tools [tool] [option]... [arg]...
```

### `msgomerge`

A cross-platform alternative to `msgmerge`, used for updating `.po` files with new translations while preserving existing ones.

**Usage:**

```sh
msgomerge [old.po] [new.po] -o [output.po]
```

[More information here](/cli/msgomerge/README.md)

### `msgofmt`

A cross-platform alternative to `msgfmt`, used for compiling `.po` files into binary `.mo` files.

**Usage:**

```sh
msgofmt [input.po] -o [output.mo]
```

[More information here](/cli/msgofmt/README.md)

### `xgotext`

A Go-compatible version of `xgettext` for extracting translatable strings from Go source code.

**Usage:**

```sh
xgotext -o [output.po] [file1.go] [file2.go] ...
```

[More information here](/cli/xgotext/README.md)

---

📌 **Coming Soon:** More CLI tools for advanced Gettext operations.

---

## Library

The core library is located in `pkg/` and provides structured handling of `.po` and `.mo` files.

### `go/parse`

Extracts Gettext-compatible strings from Go source code. Useful for generating translation templates.

#### Example:

<details>

```go
package main

import (
  goparse"github.com/Tom5521/gotext-tools/v2/pkg/go/parse"
  "fmt"
)

func main(){
  myGolangFile := `package main

  import "fmt" // Import strings are ignored!

  func MyFunc(){
    a := 10
    "My anonymous string"

    switch "a"{
      case "b":
      case "c":
    }
  }`

  file,err := goparse.FromString(myGolangFile,"my-file.go")
  if err != nil{
    panic(err)
  }

  fmt.Println(file.Entries)
}
```

</details>

### `po`

The main package for working with `.po` files. Includes:

- **`Entry` & `Entries`** – Structured representation of translation entries.
- **`File`**
- **Sorting & Comparison** – Easily organize and compare translations.

<details>

```go
package main

import (
  "os"
  "github.com/Tom5521/gotext-tools/v2/pkg/po"
  "github.com/Tom5521/gotext-tools/v2/pkg/po/compiler"
  "github.com/Tom5521/gotext-tools/v2/pkg/po/parse"
)


func main(){
  def,_ := parse.Mo("es.mo")
  ref,_ := parse.Po("en.pot")
  if def.Equal(ref){
    return
  }

  merged := po.Merge(def.Entries,ref.Entries)
  merged = merged.CleanFuzzy().CleanDuplicates()
  compiler.PoToWriter(merged,os.Stdout)
}

```

</details>

### `po/compile`

Compiles parsed `.po` files into `.mo` (binary) or updated `.po` files.

<details>

```go
package main

import (
  "os"
  "github.com/Tom5521/gotext-tools/v2/pkg/po"
  "github.com/Tom5521/gotext-tools/v2/pkg/po/compile"
)

func main(){
  myFile := &po.File{
    Name: "My File!",
    Entries: po.Entries{
      {
        ID: "Hello World!",
        Str: "Hola Mundo!",
      },
      {
        ID: "Bye World!",
        Str: "Adios Mundo!",
      },
    },
  }

  compile.PoToWriter(myFile,os.Stdout)
}
```

</details>

### `po/parse`

Parsers for reading `.po` and `.mo` files into structured Go objects.

<details>

```go
package main

import (
  "github.com/Tom5521/gotext-tools/v2/pkg/po/parse"
)

func main(){
  myPoFile := `msgid "hello"
msgstr "hola"

#, fuzzy
msgid "world"
msgstr "mundo"`

  myFile,_ := parse.PoFromString(myPoFile,"my_po_file.po")
}
```

</details>

---

## Installation

### From Source

```sh
git clone https://github.com/Tom5521/gotext-tools
cd gotext-tools
make local-install APP=gotext-tools
```

#### Or...

```sh
git clone https://github.com/Tom5521/gotext-tools
cd gotext-tools
make go-install-gotext-tools
```

### Pre-built Binaries

Check the [Releases](https://github.com/Tom5521/gotext-tools/releases) page for pre-compiled executables.

---

### Notes

**Recommendation:** the version of this module is go1.18, it will work fine, but if you need better performance I recommend you to use go1.21 as minimum, since several internal functions are manual implementations (since for go1.18 they haven't added some libraries yet) so they are a bit slower than the official implementations.

Here is an example with benchmarks:

#### Go 1.24.3

```
$ just bench ./pkg/po/compile/mo_benchmark_test.go
goos: linux
goarch: amd64
cpu: Intel(R) Core(TM) i5-3330 CPU @ 3.00GHz
BenchmarkMoCompiler/WithHashTable-4               169875              5946 ns/op
BenchmarkMoCompiler/WithoutHashTable-4            218752              5495 ns/op
PASS
ok      command-line-arguments  2.350s
```

#### Go 1.18

```
$ just gocmd=go1.18 bench ./pkg/po/compile/mo_benchmark_test.go
goos: linux
goarch: amd64
cpu: Intel(R) Core(TM) i5-3330 CPU @ 3.00GHz
BenchmarkMoCompiler/WithHashTable-4               135243              7897 ns/op
BenchmarkMoCompiler/WithoutHashTable-4            161006              7837 ns/op
PASS
ok      command-line-arguments  2.507s
```

I have to note that this only applies to the library,
the CLI binaries use a module with version 1.23 and are
compiled with the latest version of golang, so these problems do not affect them.

---

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.

---

## License

This project is open-source under the **MIT License**. See [LICENSE](https://github.com/Tom5521/gotext-tools/blob/main/LICENSE) for details.
