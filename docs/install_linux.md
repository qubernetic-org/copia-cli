# Installing copia-cli on Linux

## Recommended _(Official)_

### Homebrew

```bash
brew install qubernetic/tap/copia-cli
```

To upgrade:

```bash
brew upgrade qubernetic/tap/copia-cli
```

### Debian/Ubuntu (.deb)

Download the `.deb` package from [GitHub Releases](https://github.com/qubernetic/copia-cli/releases/latest):

```bash
# Download the latest .deb (amd64)
curl -LO https://github.com/qubernetic/copia-cli/releases/latest/download/copia-cli_*_linux_amd64.deb

# Install
sudo dpkg -i copia-cli_*_linux_amd64.deb
```

### Fedora/RHEL (.rpm)

Download the `.rpm` package from [GitHub Releases](https://github.com/qubernetic/copia-cli/releases/latest):

```bash
# Download the latest .rpm (amd64)
curl -LO https://github.com/qubernetic/copia-cli/releases/latest/download/copia-cli_*_linux_amd64.rpm

# Install
sudo dnf install -y copia-cli_*_linux_amd64.rpm
```

### Precompiled binaries

[Copia CLI releases](https://github.com/qubernetic/copia-cli/releases/latest) contain precompiled binaries for `amd64`, `arm64`, `386`, and `armv6` architectures.

```bash
# Download and extract
tar xzf copia-cli_*_linux_amd64.tar.gz

# Move to PATH
sudo mv copia-cli /usr/local/bin/
```

## Building from source

See [install_source.md](install_source.md).
