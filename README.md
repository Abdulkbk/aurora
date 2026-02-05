# Aurora 🌌

A CLI tool for building custom Docker images from GitHub PRs and forks for [Lightning Polar](https://lightningpolar.com/).

Aurora makes it easy for code reviewers to test proposed changes by building Docker images directly from pull requests, without manually cloning repos or dealing with complex build setups.

## Features

- 🔗 **Build from PRs** - Just paste a GitHub PR URL and Aurora fetches the fork and branch automatically
- 🐳 **Docker Integration** - Builds images using embedded Dockerfiles optimized for Polar
- ⚡ **Lightning-focused** - Currently supports LND with more node types coming soon
- 📦 **Single Binary** - No dependencies, just download and run

## Installation

### From Source

```bash
git clone https://github.com/Abdulkbk/aurora.git
cd aurora
make build
```

This creates an `aurora` binary in the current directory.

### Install to PATH

```bash
make install
```

This installs `aurora` to your `$GOPATH/bin`.

## Usage

### Build from a Pull Request

```bash
aurora build --pr https://github.com/lightningnetwork/lnd/pull/1234 --tag my-test-image
```

Aurora will:

1. Parse the PR URL
2. Fetch the fork URL and branch from GitHub API
3. Build a Docker image using the fork's code
4. Tag it with your specified name

### Build from a Fork/Branch

If you have a specific fork and branch (or if the PR branch was deleted):

```bash
aurora build --repo https://github.com/username/lnd --branch feature-branch --tag my-fork-image
```

### Example Output

```
🚀 Aurora Build
===============
📋 PR:     lightningnetwork/lnd#10545
🔍 Fetching PR details from GitHub...
📝 Title:  switchrpc: improve SendOnion error handling
📊 State:  open
🔗 Fork:   https://github.com/calvinrzachman/lnd.git
🌿 Branch: switchrpc-error-handle-combined
📦 Type:   lnd (default)
🏷️  Tag:    sendonion

🔨 Building Docker image...
----------------------------
[... Docker build output ...]
----------------------------
✅ Build complete! Image: sendonion:aurora

To use in Polar, add this as a custom node image.
```

## Using with Polar

After building an image with Aurora, you can use it in Lightning Polar:

1. Open Polar
2. Create a new network or edit an existing one
3. When adding an LND node, select "Managed" and choose your custom image
4. The image will appear with the tag you specified (e.g., `sendonion:aurora`)

## Commands

| Command          | Description                            |
| ---------------- | -------------------------------------- |
| `aurora build`   | Build a Docker image from a PR or fork |
| `aurora version` | Show version information               |
| `aurora help`    | Show help information                  |

### Build Flags

| Flag          | Description                | Required                  |
| ------------- | -------------------------- | ------------------------- |
| `--pr`        | GitHub PR URL              | Either `--pr` or `--repo` |
| `--repo`      | GitHub repository URL      | Either `--pr` or `--repo` |
| `--branch`    | Branch name                | Required with `--repo`    |
| `--tag`       | Docker image tag           | ✅ Yes                    |
| `--node-type` | Node type (default: `lnd`) | No                        |

## Supported Node Types

| Node           | Status         |
| -------------- | -------------- |
| LND            | ✅ Supported   |
| Bitcoin Core   | 🔜 Coming Soon |
| Core Lightning | 🔜 Coming Soon |
| Eclair         | 🔜 Coming Soon |
| LIT            | 🔜 Coming Soon |
| Taproot Assets | 🔜 Coming Soon |

## Development

### Prerequisites

- Go 1.21+
- Docker

### Building

```bash
make build      # Build the binary
make test       # Run tests
make fmt        # Format code
make tidy       # Tidy dependencies
make clean      # Clean build artifacts
```

## How It Works

1. **Parse PR URL** - Extracts owner, repo, and PR number from GitHub URLs
2. **Fetch PR Details** - Calls GitHub API to get the fork's clone URL and branch name
3. **Prepare Dockerfile** - Uses an embedded Dockerfile optimized for the node type
4. **Build Image** - Runs `docker build` with the fork URL and branch as build args
5. **Tag & Output** - Tags the image and provides instructions for Polar

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see [LICENSE](LICENSE) for details.

## Acknowledgments

- [Lightning Polar](https://lightningpolar.com/) - The excellent Lightning Network development tool
- [LND](https://github.com/lightningnetwork/lnd) - Lightning Network Daemon
