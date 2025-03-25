# Terminal Typing Speed Test

A terminal-based typing speed test—because the world was begging for one more. Built with Go and powered by [Charm](https://charm.sh/) libraries, this TUI tool lets you measure your typing speed in the most hacker-ish way possible.

## Features

- Minimalist terminal UI
- Random typing prompts
- Real-time WPM calculation
- Accuracy tracking
- Restart functionality
- Supports deleting by word (Ctrl+Backspace)
- Cursor positioning and animation

## Installation

```sh
go install github.com/yourusername/typing-speed-test@latest
```

Or clone and build manually:

```sh
git clone https://github.com/yourusername/typing-speed-test.git
cd typing-speed-test
go build -o typing-test
./typing-test
```

## Usage

Run the program in your terminal:

```sh
./typing-test
```

- Type the given text as fast and accurately as possible.
- Press `Backspace` to delete characters.
- Use `Ctrl+Backspace` (or `Cmd+Backspace` on Mac) to delete entire words.
- Press `Space` to skip to the next word when in the middle or end of a word.
- Press `r` to restart after finishing.
- Press `Ctrl+C` to quit.

## Future Plans

- Customizable themes and colors
- Different difficulty levels
- Support for custom text inputs
- Multiplayer mode?

## Contributing

PRs are welcome! Feel free to fork the repo, make improvements, and submit a pull request.

## License

MIT
