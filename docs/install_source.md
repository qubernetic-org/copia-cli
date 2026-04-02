# Installation from source

1. Verify that you have Go 1.26+ installed

   ```bash
   $ go version
   ```

   If `go` is not installed, follow instructions on [the Go website](https://golang.org/doc/install).

2. Clone this repository

   ```bash
   $ git clone https://github.com/qubernetic/copia-cli.git
   $ cd copia-cli
   ```

3. Build and install

   ```bash
   $ make build
   $ sudo cp bin/copia-cli /usr/local/bin/
   ```

4. Verify installation

   ```bash
   $ copia-cli --version
   ```

## Cross-compilation

The Makefile supports cross-compilation via Go environment variables:

```bash
# Linux ARM64
$ GOOS=linux GOARCH=arm64 make build

# macOS ARM64
$ GOOS=darwin GOARCH=arm64 make build

# Windows AMD64
$ GOOS=windows GOARCH=amd64 make build
```
