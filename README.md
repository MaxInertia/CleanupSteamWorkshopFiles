# BO3 Workshop Content Manager

A command-line tool to manage and clean up Black Ops 3 Workshop content from your Steam installation.

This README was generated by `claude-3.5-sonnet`.

## Prerequisites

### Installing Go

1. Visit the [official Go downloads page](https://golang.org/dl/)
2. Download the installer for your operating system:
   - Windows: Download the `.msi` installer
   - macOS: Download the `.pkg` installer
   - Linux: Follow distribution-specific instructions
3. Run the installer and follow the prompts
4. Verify the installation by opening a terminal/command prompt and running:
   ```bash
   go version
   ```

## Building the Tool

1. Clone this repository or download the source files
2. Open a terminal/command prompt in the project directory
3. Build the executable:
   ```bash
   go build
   ```
   This will create an executable file:
   - Windows: `workshop-manager.exe`
   - macOS/Linux: `workshop-manager`

## Usage

Run the tool by providing your Steam installation path as an argument:

```bash
# Windows
workshop-manager.exe "C:\Program Files (x86)\Steam"

# macOS
workshop-manager "/Users/username/Library/Application Support/Steam"

# Linux
workshop-manager "~/.local/share/Steam"
```

The tool will:
1. List all your BO3 Workshop subscriptions with their sizes
2. Allow you to select items to delete by entering their numbers (space-separated)
3. Show the total space that will be reclaimed
4. Ask for confirmation before deletion

## Example Output
```
Steam Path:
        C:\Program Files (x86)\Steam
BO3 Workshop Subscriptions Path:
        C:\Program Files (x86)\Steam\steamapps\workshop\content\311210

Subscriptions:
        0: Custom Map 1 (map) [1.2 GB]
        1: Custom Map 2 (map) [800 MB]
Space used: 2.0 GB

Delete which? 1
Selected Subscriptions:
        1: Custom Map 2 (C:\Program Files (x86)\Steam\steamapps\workshop\content\311210\123456789)

This will reclaim 800 MB

Delete (y/n)?
```

## Safety Features

- The tool only operates within the BO3 Workshop content directory
- Confirmation is required before deletion
- Path validation ensures deletions only occur in the correct directory
- Symbolic links are detected and skipped to avoid escaping the workshop directory