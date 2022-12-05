```yaml
version: '3'

tasks:
  default:
    cmds:
    - (cd analytics1d && go clean && go build -o ../bin/analytics1d)
    - (cd analytics1w && go clean && go build -o ../bin/analytics1w)
    - (cd analytics1m && go clean && go build -o ../bin/analytics1m)
```

Files are around 16MB

```text
$ ls -l bin
total 100608
-rwxr-xr-x  1 ally  staff  17142624 Dec  5 20:52 analytics1d
-rwxr-xr-x  1 ally  staff  17181360 Dec  5 20:52 analytics1m
-rwxr-xr-x  1 ally  staff  17181360 Dec  5 20:52 analytics1w
```

Asking the oracle that is ChatGPT

```text
How can I make go binaries smaller?
```

There are several ways to make Go binaries smaller. Here are a few options you can try:

- Use the `-ldflags="-s -w"` flag when building your Go binary to strip symbol information and debug information from the binary, which can make it significantly smaller.

```bash
go build -ldflags="-s -w"
```

Building with this flag brings them down to around 12mb

```text
$ ls -l bin 
total 77856
-rwxr-xr-x  1 ally  staff  13276000 Dec  5 20:57 analytics1d
-rwxr-xr-x  1 ally  staff  13290160 Dec  5 20:57 analytics1m
-rwxr-xr-x  1 ally  staff  13290160 Dec  5 20:57 analytics1w
```

- Use the `-trimpath` flag when building your Go binary to remove the file system paths from the binary. This can make the binary smaller by removing information that is not needed at runtime.

```bash
go build -trimpath
```

Building with this flag doesn't really do anything.

- Use the `-buildmode=pie` flag when building your Go binary to build a position-independent executable (PIE). This can make the binary smaller by removing unused code and data that is not needed at runtime.

```bash
go build -buildmode=pie
```

Building with this flag doesn't really do anything.

- Use the `-buildmode=c-shared` flag when building your Go binary to build a C shared library. This can make the binary smaller by removing unused code and data that is not needed at runtime, and by allowing multiple Go programs to share the same compiled code.

```bash
go build -buildmode=c-shared
```

Building with this flag doesn't really do anything.

- Use the `strip` command to remove symbol information and debug information from the binary after it has been built. This can make the binary smaller by removing unnecessary data that is not needed at runtime.

```bash
strip binary
```

For more information about these and other techniques for making Go binaries smaller, you can refer to the Go documentation: https://golang.org/doc/install/source#build.

---

Using all of these flags brings it down to around 13MB which isn't too bad.