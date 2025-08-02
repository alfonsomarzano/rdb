# RDB - Resource Database CLI Tool

A Windows-first CLI tool for managing a local "Resource Database" (RDB) — a directory tree of typed assets — with git-like workflows.

## Features

- **Typed Assets**: Store assets under folders named by ID (e.g., `1030002/`, `1000624/`, `1010042/`)
- **Git-like Workflows**: Familiar commands like `init`, `status`, `add`, `commit`, `branch`, `merge`
- **Content Integrity**: SHA-256 for objects and manifest index for fast lookups
- **Human-readable**: Predictable directory layout
- **Portable**: Package working tree into `.rdbdata` ZIP files
- **Windows-friendly**: PowerShell examples, CRLF handling, UTF-8 paths

## Quick Start

```bash
# Build the CLI tool
make build

# Initialize a new RDB repository
rdb init --layout tree --types "text,audio,texture,shader,mesh"

# Add assets
rdb add ./assets/1030002/ --type text --id 1030002 --name "DialogLine_Intro"

# Check status
rdb status

# Commit changes
rdb commit -m "Add intro dialog line"

# Build distributable package
rdb build
```

## Installation

### From Source

```bash
# Clone the repository
git clone <repository-url>
cd RDB

# Install dependencies
make deps

# Build the CLI tool
make build
```

### Prebuilt Binary

Download the prebuilt binary for your platform from the releases page and place it in your `PATH`.

## Usage

### Core Commands

- `rdb init` - Initialize a new RDB repository
- `rdb status` - Show working tree status
- `rdb add` - Stage files for commit
- `rdb commit` - Create a new commit
- `rdb log` - Show commit history
- `rdb list` - List asset types and folders
- `rdb cd` - Change directory to asset folder
- `rdb build` - Create `.rdbdata` package

### Additional Features

- `rdb init --layout <layout>` - Specify repository layout (`tree` or `flat`)
- `rdb add --type <type> --id <id>` - Specify asset type and ID when adding files
- `rdb commit --amend` - Amend the previous commit
- `rdb log --oneline` - Show abbreviated commit history
- `rdb build --compression <method>` - Specify compression method (`store` or `deflate`)

## Directory Structure

```
<repo_root>/
  .rdb/                         # internal metadata
    config.json                 # repo config
    HEAD                        # current branch ref
    refs/heads/<branch>         # branch pointers
    index                       # staging index
    objects/                    # content-addressed blobs
  assets/                       # top-level assets directory
    1030002/                    # asset id folder (Strings)
      meta.json                 # metadata
      en.txt                    # payload files
      fr.txt
    1000624/                    # asset id folder (Flash Images)
    1010042/                    # asset id folder (Loading Screens)
  .rdbignore                    # ignore rules
```

## Makefile Targets

- `make build` - Build the RDB CLI tool
- `make build-all` - Build for multiple platforms
- `make deps` - Install dependencies
- `make test` - Run tests
- `make test-coverage` - Run tests with coverage
- `make clean` - Clean build artifacts
- `make install` - Install the tool
- `make run` - Run the tool
- `make help` - Show help

## License

MIT License