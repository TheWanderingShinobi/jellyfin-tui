# Jellyfin TUI

Jellyfin TUI is a terminal user interface for the Jellyfin media server. It allows you to browse, search, and play your media library directly from the command line.

## Features

- Browse your Jellyfin media library
- Search for specific media items
- Play videos using MPV
- Manage playlists
- User-friendly terminal interface

## Prerequisites

Before you begin, ensure you have met the following requirements:

- Go 1.16 or higher
- MPV media player
- A running Jellyfin server

## Installation

1. Clone the repository:

```

git clone https://github.com/yourusername/jellyfin-tui.git
cd jellyfin-tui

```

2. Build the application:

```

make build

```

Or, if you don't have Make installed:

```

go build -o jellyfin-tui ./cmd/jellyfin-tui

```

3. (Optional) Move the binary to a directory in your PATH to run it from anywhere:

```

sudo mv jellyfin-tui /usr/local/bin/

```

## Configuration

On first run, a default configuration file will be created at `~/.config/jellyfin-tui/config.json`. Edit this file to change the default Jellyfin server URL and other settings:

```json
{
  "server_url": "http://your-jellyfin-server:8096",
  "default_user": "",
  "items_per_page": 20
}
```

## Usage

Run the application:

```
jellyfin-tui
```

Use the arrow keys to navigate, Enter to select, and 'q' to quit. Press 'h' for help at any time.

## Controls

- Arrow keys / j,k: Navigate
- Enter: Select/Play
- q: Quit (from main menu) or Go back
- /: Search
- f: Filter (in browse view)
- p: Add to playlist (in detail view)
- h: Help

## Contributing

Contributions to Jellyfin TUI are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgements

- [Jellyfin](https://jellyfin.org/) for the amazing media server
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) for the TUI framework
- [MPV](https://mpv.io/) for video playback
