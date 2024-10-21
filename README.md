# GoTube TUI

GoTube is a terminal-based YouTube video search and play tool built with Bubble
Tea and Lipgloss. You can search for videos directly from the terminal, navigate
through the search results, and play selected videos using `mpv`.

## Installation

1. **Clone the repository**:

```bash
git clone https://github.com/your-username/gotube.git
```

## Keybindings

**Search Mode**:
    Type to enter your search query.
    Press Enter to search for YouTube videos.

**Results Mode**:
    Use `j` or down to move the cursor down through the list of results.
    Use `k` or up to move the cursor up through the list of results.
    Press `Enter` to play the selected video with mpv.
    Press `q` to quit the program.

## Requirements

`Go`: Install the Go programming language from https://golang.org.

`mpv`: Make sure mpv is installed on your system.  
    You can install it via your package manager:  
        - Debian/Ubuntu: `sudo apt install mpv`  
        - Arch: `sudo pacman -S mpv`  
        - macOS: `brew install mpv`  

## Contributing

Feel free to open issues or submit pull requests for improvements and bug fixes.

