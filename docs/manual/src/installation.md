# Installation

## Homebrew (macOS/Linux)

```bash
brew install qubernetic-org/tap/copia
```

## Precompiled Binaries

Download the latest release from [GitHub Releases](https://github.com/qubernetic-org/copia-cli/releases/latest).

### Linux

```bash
curl -sL https://github.com/qubernetic-org/copia-cli/releases/latest/download/copia_linux_amd64.tar.gz | tar xz
sudo mv copia /usr/local/bin/
```

### macOS

```bash
curl -sL https://github.com/qubernetic-org/copia-cli/releases/latest/download/copia_darwin_arm64.tar.gz | tar xz
sudo mv copia /usr/local/bin/
```

### Windows

```powershell
Invoke-WebRequest -Uri https://github.com/qubernetic-org/copia-cli/releases/latest/download/copia_windows_amd64.zip -OutFile copia.zip
Expand-Archive copia.zip -DestinationPath "$env:LOCALAPPDATA\Programs\copia"
```

Add `$env:LOCALAPPDATA\Programs\copia` to your PATH.

## Build from Source

Requires [Go 1.26+](https://go.dev/dl/).

```bash
go install github.com/qubernetic-org/copia-cli/cmd/copia@latest
```

## Verify Installation

```bash
copia --version
```
