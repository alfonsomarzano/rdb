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

```powershell
# Initialize a new RDB repository
rdb init --layout tree --types "text,audio,texture,shader,mesh"

# Add assets
rdb add .\assets\1030002\ --type text --id 1030002 --name "DialogLine_Intro"

# Check status
rdb status

# Commit changes
rdb commit -m "Add intro dialog line"

# Build distributable package
rdb build
```

## Installation

```powershell
# Clone the repository
git clone <repository-url>
cd RDB

# Install dependencies
pip install -r requirements.txt

# Install the CLI tool
pip install -e .
```

## Usage

### Core Commands

- `rdb init` - Initialize a new RDB repository
- `rdb status` - Show working tree status
- `rdb add` - Stage files for commit
- `rdb commit` - Create a new commit
- `rdb log` - Show commit history
- `rdb diff` - Show changes between commits

### Branching & Merging

- `rdb branch` - Manage branches
- `rdb checkout` - Switch branches
- `rdb merge` - Merge branches

### Packaging

- `rdb build` - Create `.rdbdata` package
- `rdb unpack` - Extract `.rdbdata` package

### Remote Operations

- `rdb remote` - Manage remote repositories
- `rdb push` - Push to remote
- `rdb pull` - Pull from remote
- `rdb clone` - Clone a repository

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

## License

MIT License 