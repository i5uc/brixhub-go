# BrixHub Go

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://golang.org)
[![Release](https://img.shields.io/github/v/tag/i5uc/brixhub-go?label=latest%20release)](https://github.com/i5uc/brixhub-go/releases/latest)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

Un-Official Go client library and CLI for the [BrixHub API](https://brixhub.net) — Search across 11+ billion documents.

## Features

- **Multi-criteria Search** — Name, email, phone, address, gaming IDs, vehicle info
- **Reverse Lookups** — Find profiles by email, phone, or IBAN
- **Account Management** — Quotas, usage history, plan details
- **Rate Limit Handling** — Automatic tracking of remaining quota
- **Typed Errors** — Easy error handling with specific error types
- **CLI Tool** — Command-line interface for quick searches

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [CLI Usage](#cli-usage)
- [Library Usage](#library-usage)
- [Building from Source](#building-from-source)
- [API Reference](#api-reference)
- [Rate Limits](#rate-limits)

## Installation

### Prerequisites

- Go 1.21 or later
- BrixHub API key (starts with `brix_`)

### Get an API Key

Register at [brixhub.site](https://brixhub.site/api-docs) to obtain your API key.

### Configure with .env

Copy the example env file and set your API key:

```bash
cp .env.example .env
# then edit .env and fill in BRIXHUB_API_KEY
```

The CLI and examples automatically load `.env` when present.

### Download Binaries

Binary releases are available on GitHub Releases:

- Linux x64: `brixhub-linux-amd64.tar.gz`
- Linux arm64: `brixhub-linux-arm64.tar.gz`
- macOS x64: `brixhub-darwin-amd64.tar.gz`
- macOS arm64: `brixhub-darwin-arm64.tar.gz`
- Windows x64: `brixhub-windows-amd64.exe.tar.gz`

Download the latest release from:

```text
https://github.com/i5uc/brixhub-go/releases/latest
```

Extract and install:

```bash
tar xzf brixhub-linux-amd64.tar.gz
sudo mv brixhub /usr/local/bin/
```

### Install CLI from source

```bash
go install github.com/i5uc/brixhub-go/cmd/brixhub@latest
```

### Install Library

```bash
go get github.com/i5uc/brixhub-go
```

## Quick Start

Set your API key:

```bash
export BRIXHUB_API_KEY="brix_your_api_key_here"
```

Or pass it directly to the binary:

```bash
brixhub --api-key "brix_your_api_key_here" search --nom "Dupont" --prenom "Jean" --ville "Paris"
```

Run a search:

```bash
brixhub search --nom "Dupont" --prenom "Jean" --ville "Paris"
```

## CLI Usage

You can pass the API key directly to the binary with `--api-key` or rely on `BRIXHUB_API_KEY`.

### Search

```bash
# By name
brixhub search --nom "Dupont" --prenom "Jean"

# With flexible matching (contains vs exact)
brixhub search --nom "Martin" --flexible

# By email
brixhub search --email "jean.dupont@example.com"

# By Discord ID
brixhub search --discord-id "123456789012345678"

# By license plate
brixhub search --immat "AA-123-AA"

# Multiple criteria
brixhub search --nom "Dupont" --prenom "Jean" --ville "Paris" --cp "75001"
```

### Reverse Lookup

```bash
# By email
brixhub lookup email jean.dupont@example.com

# By phone (any French format)
brixhub lookup phone 0612345678
brixhub lookup phone "+33 6 12 34 56 78"

# By IBAN
brixhub lookup iban FR7630006000011234567890189
```

### Account Commands

```bash
# Check account info and quota
brixhub account

# View usage history (Pro+ plans)
brixhub usage --limit 20

# Check API health
brixhub health
```

## Library Usage

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/i5uc/brixhub-go/pkg/brixhub"
)

func main() {
    client, err := brixhub.NewClient(os.Getenv("BRIXHUB_API_KEY"))
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Search
    results, _, err := client.Search(ctx, &brixhub.SearchRequest{
        NomFamille: "Dupont",
        Prenom:     "Jean",
        Ville:      "Paris",
        Flexible:   true,
    })
    if err != nil {
        log.Fatal(err)
    }

    for _, profile := range results.Results {
        fmt.Printf("%s %s - %s\n", 
            profile.Prenom, 
            profile.NomFamille, 
            profile.Email)
    }

    // Check remaining quota
    fmt.Printf("Remaining: %d requests\n", client.GetRemainingQuota())
}
```

## Building from Source

### Clone and Build

```bash
# Clone the repository
git clone https://github.com/i5uc/brixhub-go.git
cd brixhub-go

# Download dependencies
go mod tidy

# Build the CLI binary
go build -o brixhub ./cmd/brixhub

# Run
./brixhub --help
```

### Using Make

```bash
# Build for current platform
make build

# Run tests
make test

# Build for all platforms
make build-all

# Install to $GOPATH/bin
make install

# Clean build artifacts
make clean
```

### Cross-Compilation

The Makefile supports building for multiple platforms:

| Platform | Command Output |
|----------|---------------|
| macOS (Intel) | `brixhub-darwin-amd64` |
| macOS (Apple Silicon) | `brixhub-darwin-arm64` |
| Linux (x64) | `brixhub-linux-amd64` |
| Linux (ARM64) | `brixhub-linux-arm64` |
| Windows | `brixhub-windows-amd64.exe` |

### Running Examples

```bash
# Copy and edit .env for local credentials
cp .env.example .env

# Run basic search example
go run ./examples/basic_search

# Run reverse lookup example
go run ./examples/reverse_lookup

# Run account info example
go run ./examples/account_info
```

## API Reference

### Search Parameters

| Field | Type | Description |
|-------|------|-------------|
| `NomFamille` | string | Last name |
| `Prenom` | string | First name |
| `NomNaissance` | string | Birth name |
| `NomUtilisateur` | string | Username (exact match) |
| `DateNaissance` | string | Birth date (YYYY-MM-DD) |
| `Genre` | string | Gender (M/F) |
| `Email` | string | Email address |
| `Telephone` | string | Phone number |
| `Mobile` | string | Mobile number |
| `AdresseIP` | string | IP address |
| `Adresse` | string | Street address |
| `CodePostal` | string | Postal code |
| `Ville` | string | City |
| `Pays` | string | Country |
| `NIR` | string | Social security number |
| `SIRET` | string | 14-digit business ID |
| `SIREN` | string | 9-digit business ID |
| `SteamID` | string | Steam identifier |
| `DiscordID` | string | Discord identifier |
| `FiveMLicense` | string | FiveM license |
| `Immatriculation` | string | License plate |
| `VINPlaque` | string | VIN number |
| `Page` | int | Page number (default: 1) |
| `PerPage` | int | Results per page (default: 10) |
| `Flexible` | bool | Approximate matching |

### Error Handling

```go
results, _, err := client.Search(ctx, req)
if err != nil {
    if apiErr, ok := err.(*brixhub.APIError); ok {
        switch {
        case apiErr.IsUnauthorized():
            // Invalid or expired API key
        case apiErr.IsRateLimit():
            // Quota exceeded or rate limited
        case apiErr.IsPlanLimited():
            // Feature requires higher plan
        }
    }
}
```

## Rate Limits

| Plan | Daily | Per Minute | Pagination |
|------|-------|------------|------------|
| Starter | 1,000 | 10 | Page 1 only |
| Premium | 5,000 | 20 | Unlimited |
| Pro | 10,000 | 30 | Unlimited |
| Enterprise | 100,000 | 100 | Unlimited |

Rate limit headers are automatically parsed and available via `client.GetRemainingQuota()`.

## License

MIT License — see [LICENSE](LICENSE) file.

## Support

- Website: [brixhub.net](https://brixhub.net)
- Discord: [discord.gg/brixhub](https://discord.gg/brixhub)
- Developer : [i5uc](https://guns.lol/i5uc)
- Email: [i5aac@larp.expert](mailto:i5aac@larp.expert)



